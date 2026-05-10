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
	mu             sync.Mutex
	state          *game.GameState
	clients        []ClientSender
	clientPlayer   map[ClientSender]game.PlayerID
	activeClients  map[game.PlayerID]ClientSender
	remainingGrace map[game.PlayerID]float64
	pausedPlayer   game.PlayerID
	actions        chan game.Action
	stop           chan struct{}
	started        bool
}

func New() *Room {
	state := game.NewGameState()
	state.Hexes = mapgen.GenerateTestMap()
	state.Waiting = true

	spawns := mapgen.SpawnPositions()
	remainingGrace := make(map[game.PlayerID]float64, len(spawns))
	for i, pos := range spawns {
		pid := game.PlayerID(i + 1)
		state.Players[pid] = &game.Player{ID: pid}
		state.Hexes[pos].Owner = pid
		state.Hexes[pos].Capital = true
		remainingGrace[pid] = disconnectGrace.Seconds()
	}

	return &Room{
		state:          state,
		clientPlayer:   make(map[ClientSender]game.PlayerID),
		activeClients:  make(map[game.PlayerID]ClientSender),
		remainingGrace: remainingGrace,
		actions:        make(chan game.Action, actionQueueSize),
		stop:           make(chan struct{}),
	}
}

func (r *Room) State() *game.GameState {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.state
}

// OnConnect registers a client for the given player and sends a full snapshot.
// Starts the game loop when all players are connected for the first time.
// Unpauses and saves remaining grace time when a disconnected player returns.
func (r *Room) OnConnect(c ClientSender, pid game.PlayerID) {
	r.mu.Lock()
	defer r.mu.Unlock()

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
	} else if r.started && r.state.Paused && r.pausedPlayer == pid {
		r.remainingGrace[pid] = r.state.PauseTimeLeft
		r.state.Paused = false
		r.state.PauseTimeLeft = 0
		r.pausedPlayer = 0
		for _, cl := range r.clients {
			if cl != c {
				cl.SendSnapshot(r.state)
			}
		}
	}

	c.SendSnapshot(r.state)
}

// OnDisconnect removes a client, pauses the game, and starts the loop-driven forfeit countdown.
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
		r.state.Paused = true
		r.state.PauseTimeLeft = r.remainingGrace[pid]
		r.pausedPlayer = pid
		r.broadcast()
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
