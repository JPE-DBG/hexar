package room

import (
	"hexar/internal/game"
	"hexar/internal/mapgen"
	"sync"
	"time"
)

const disconnectGrace = 2 * time.Minute
const actionQueueSize = 256

type ClientSender interface {
	SendSnapshot(state *game.GameState)
}

type Room struct {
	mu               sync.Mutex
	state            *game.GameState
	clients          []ClientSender
	clientPlayer     map[ClientSender]game.PlayerID
	activeClients    map[game.PlayerID]ClientSender
	disconnectTimers map[game.PlayerID]*time.Timer
	actions          chan game.Action
	stop             chan struct{}
	started          bool
}

func New() *Room {
	state := game.NewGameState()
	state.Hexes = mapgen.GenerateTestMap()
	state.Waiting = true

	spawns := mapgen.SpawnPositions()
	for i, pos := range spawns {
		pid := game.PlayerID(i + 1)
		state.Players[pid] = &game.Player{ID: pid}
		state.Hexes[pos].Owner = pid
		state.Hexes[pos].Capital = true
	}

	return &Room{
		state:            state,
		clientPlayer:     make(map[ClientSender]game.PlayerID),
		activeClients:    make(map[game.PlayerID]ClientSender),
		disconnectTimers: make(map[game.PlayerID]*time.Timer),
		actions:          make(chan game.Action, actionQueueSize),
		stop:             make(chan struct{}),
	}
}

func (r *Room) State() *game.GameState {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.state
}

// OnConnect registers a client for the given player and sends a full snapshot.
// Cancels any pending disconnect/forfeit timer for this player.
// Starts the game loop when all players are connected for the first time.
func (r *Room) OnConnect(c ClientSender, pid game.PlayerID) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if timer, ok := r.disconnectTimers[pid]; ok {
		timer.Stop()
		delete(r.disconnectTimers, pid)
	}

	if old, ok := r.activeClients[pid]; ok {
		for i, cl := range r.clients {
			if cl == old {
				r.clients = append(r.clients[:i], r.clients[i+1:]...)
				break
			}
		}
		delete(r.clientPlayer, old)
	}

	r.clients = append(r.clients, c)
	r.clientPlayer[c] = pid
	r.activeClients[pid] = c

	if !r.started && len(r.activeClients) == len(r.state.Players) {
		r.state.Waiting = false
		r.started = true
		for _, cl := range r.clients {
			if cl != c {
				cl.SendSnapshot(r.state)
			}
		}
		go r.Run()
	}

	c.SendSnapshot(r.state)
}

// OnDisconnect removes a client and starts a 30-second forfeit timer.
func (r *Room) OnDisconnect(pid game.PlayerID, c ClientSender) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.activeClients[pid] != c {
		return
	}
	delete(r.activeClients, pid)
	delete(r.clientPlayer, c)
	for i, cl := range r.clients {
		if cl == c {
			r.clients = append(r.clients[:i], r.clients[i+1:]...)
			break
		}
	}

	if !r.state.Over && r.started {
		timer := time.AfterFunc(disconnectGrace, func() {
			r.EnqueueAction(game.Action{Type: game.ActionForfeit, Player: pid})
		})
		r.disconnectTimers[pid] = timer
	}
}

func (r *Room) EnqueueAction(action game.Action) {
	select {
	case r.actions <- action:
	default:
	}
}

func (r *Room) broadcast() {
	for _, c := range r.clients {
		c.SendSnapshot(r.state)
	}
}

func (r *Room) Stop() {
	close(r.stop)
}

// IsConnected returns true if the given player already has an active WebSocket connection.
func (r *Room) IsConnected(pid game.PlayerID) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, ok := r.activeClients[pid]
	return ok
}
