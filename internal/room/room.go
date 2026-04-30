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
	mu      sync.RWMutex
	state   *game.GameState
	clients []ClientSender
	stop    chan struct{}
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
		state: state,
		stop:  make(chan struct{}),
	}
}

func (r *Room) State() *game.GameState {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.state
}

func (r *Room) AddClient(c ClientSender) {
	r.mu.Lock()
	r.clients = append(r.clients, c)
	r.mu.Unlock()
	c.SendSnapshot(r.state)
}

func (r *Room) RemoveClient(c ClientSender) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, cl := range r.clients {
		if cl == c {
			r.clients = append(r.clients[:i], r.clients[i+1:]...)
			return
		}
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
