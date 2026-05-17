package net

import (
	"fmt"
	"hexar/internal/game"
)

type MsgType string

const (
	MsgSnapshot MsgType = "snapshot"
	MsgDelta    MsgType = "delta"
	MsgWelcome  MsgType = "welcome"
	MsgAction   MsgType = "action"
)

type DeltaMsg struct {
	Type          MsgType      `json:"type"`
	Players       []*PlayerDTO `json:"players"`
	Battles       []*BattleDTO `json:"battles"`
	Elapsed       float64      `json:"elapsed"`
	Over          bool         `json:"over"`
	Winner        int          `json:"winner,omitempty"`
	WinReason     string       `json:"winReason,omitempty"`
	Waiting       bool         `json:"waiting"`
	Paused        bool         `json:"paused"`
	PauseTimeLeft float64      `json:"pauseTimeLeft"`
	HexChanges    []*HexDTO    `json:"hexChanges,omitempty"`
}

type WelcomeMsg struct {
	Type     MsgType `json:"type"`
	PlayerID int     `json:"playerId"`
}

type ActionMsg struct {
	Type   MsgType `json:"type"`
	Action string  `json:"action"`
	Q      int     `json:"q"`
	R      int     `json:"r"`
}

type SnapshotMsg struct {
	Type          MsgType               `json:"type"`
	Hexes         map[string]*HexDTO    `json:"hexes"`
	Players       map[string]*PlayerDTO `json:"players"`
	Battles       []*BattleDTO          `json:"battles"`
	Elapsed       float64               `json:"elapsed"`
	Over          bool                  `json:"over"`
	Winner        int                   `json:"winner"`
	WinReason     string                `json:"winReason"`
	Waiting       bool                  `json:"waiting"`
	Paused        bool                  `json:"paused"`
	PauseTimeLeft float64               `json:"pauseTimeLeft"`
}

type HexDTO struct {
	Q             int     `json:"q"`
	R             int     `json:"r"`
	Owner         int     `json:"owner"`
	Building      int     `json:"building"`
	Level         int     `json:"level"`
	Capital       bool    `json:"capital"`
	FortifyTimer  float64 `json:"fortifyTimer"`
	PreviousOwner int     `json:"previousOwner"`
}

type PlayerDTO struct {
	ID             int     `json:"id"`
	Gold           float64 `json:"gold"`
	TP             float64 `json:"tp"`
	Tech           []bool  `json:"tech,omitempty"`
	AutoDropActive bool    `json:"autoDropActive"`
	AutoDropGrace  float64 `json:"autoDropGrace"`
	VanguardTimer  float64 `json:"vanguardTimer"`
}

type BattleDTO struct {
	AQ           int     `json:"aq"`
	AR           int     `json:"ar"`
	DQ           int     `json:"dq"`
	DR           int     `json:"dr"`
	TimeLeft     float64 `json:"timeLeft"`
	Attacker     int     `json:"attacker"`
	Defender     int     `json:"defender"`
	CounterBoost int     `json:"counterBoost"`
}

func BuildSnapshot(state *game.GameState) *SnapshotMsg {
	msg := &SnapshotMsg{
		Type:          MsgSnapshot,
		Hexes:         make(map[string]*HexDTO, len(state.Hexes)),
		Players:       make(map[string]*PlayerDTO, len(state.Players)),
		Battles:       make([]*BattleDTO, 0, len(state.Battles)),
		Elapsed:       state.Elapsed,
		Over:          state.Over,
		Winner:        int(state.Winner),
		WinReason:     state.WinReason,
		Waiting:       state.Waiting,
		Paused:        state.Paused,
		PauseTimeLeft: state.PauseTimeLeft,
	}

	for hex, hs := range state.Hexes {
		key := hexKey(hex)
		msg.Hexes[key] = &HexDTO{
			Q:             hex.Q,
			R:             hex.R,
			Owner:         int(hs.Owner),
			Building:      int(hs.Building),
			Level:         hs.Level,
			Capital:       hs.Capital,
			FortifyTimer:  hs.FortifyTimer,
			PreviousOwner: int(hs.PreviousOwner),
		}
	}

	for pid, p := range state.Players {
		key := playerKey(pid)
		tech := make([]bool, game.TechCount)
		copy(tech, p.Tech[:])
		msg.Players[key] = &PlayerDTO{
			ID:             int(p.ID),
			Gold:           p.Gold,
			TP:             p.TP,
			Tech:           tech,
			AutoDropActive: p.AutoDropActive,
			AutoDropGrace:  p.AutoDropGrace,
			VanguardTimer:  p.VanguardTimer,
		}
	}

	for _, b := range state.Battles {
		msg.Battles = append(msg.Battles, &BattleDTO{
			AQ:           b.AttackerHex.Q,
			AR:           b.AttackerHex.R,
			DQ:           b.DefenderHex.Q,
			DR:           b.DefenderHex.R,
			TimeLeft:     b.TimeLeft,
			Attacker:     int(b.Attacker),
			Defender:     int(b.Defender),
			CounterBoost: b.CounterBoost,
		})
	}

	return msg
}

func hexKey(h game.Hex) string {
	return fmt.Sprintf("%d,%d", h.Q, h.R)
}

func playerKey(pid game.PlayerID) string {
	return fmt.Sprintf("%d", pid)
}
