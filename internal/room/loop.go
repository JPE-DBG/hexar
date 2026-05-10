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
			if r.state.Paused {
				r.state.PauseTimeLeft -= game.TickDt
				if r.state.PauseTimeLeft <= 0 {
					r.state.PauseTimeLeft = 0
					r.state.Paused = false
					pid := r.pausedPlayer
					r.pausedPlayer = 0
					r.actions <- game.Action{Type: game.ActionForfeit, Player: pid}
				}
			} else {
				actions := r.drainActions()
				game.RunTick(r.state, game.TickDt, actions)
			}
			r.broadcast()
			r.mu.Unlock()
		case <-r.stop:
			return
		}
	}
}

func (r *Room) drainActions() []game.Action {
	var pending []game.Action
	for {
		select {
		case a := <-r.actions:
			pending = append(pending, a)
		default:
			return pending
		}
	}
}
