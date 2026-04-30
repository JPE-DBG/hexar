package room

import (
	"hexar/internal/game"
	"time"
)

func (r *Room) Run() {
	ticker := time.NewTicker(game.TickRate * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			r.mu.Lock()
			game.RunTick(r.state, game.TickDt)
			r.broadcast()
			r.mu.Unlock()
		case <-r.stop:
			return
		}
	}
}
