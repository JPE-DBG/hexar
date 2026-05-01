package room

import (
	"hexar/internal/game"
	"hexar/internal/mapgen"
	"sync"
)

type ClientSender interface {
	SendSnapshot(state *game.GameState)
}

type Room struct {
	mu           sync.RWMutex
	state        *game.GameState
	clients      []ClientSender
	clientPlayer map[ClientSender]game.PlayerID
	nextSlot     int
	actions      chan game.Action
	stop         chan struct{}
}

func New() *Room {
	state := game.NewGameState()
	state.Hexes = mapgen.GenerateTestMap()

	spawns := mapgen.SpawnPositions()
	for i, pos := range spawns {
		pid := game.PlayerID(i + 1)
		state.Players[pid] = &game.Player{ID: pid}
		state.Hexes[pos].Owner = pid
		state.Hexes[pos].Capital = true
	}

	return &Room{
		state:        state,
		clientPlayer: make(map[ClientSender]game.PlayerID),
		actions:      make(chan game.Action, 256),
		stop:         make(chan struct{}),
	}
}

func (r *Room) State() *game.GameState {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.state
}

func (r *Room) AddClient(c ClientSender) game.PlayerID {
	r.mu.Lock()
	r.clients = append(r.clients, c)
	var pid game.PlayerID
	if r.nextSlot < 2 {
		pid = game.PlayerID(r.nextSlot + 1)
		r.nextSlot++
	}
	r.clientPlayer[c] = pid
	r.mu.Unlock()
	c.SendSnapshot(r.state)
	return pid
}

func (r *Room) RemoveClient(c ClientSender) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, cl := range r.clients {
		if cl == c {
			r.clients = append(r.clients[:i], r.clients[i+1:]...)
			break
		}
	}
	delete(r.clientPlayer, c)
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
