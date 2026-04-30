package net

import (
	"fmt"
	"hexar/internal/game"
)

type MsgType string

const (
	MsgSnapshot MsgType = "snapshot"
)

type SnapshotMsg struct {
	Type    MsgType               `json:"type"`
	Hexes   map[string]*HexDTO    `json:"hexes"`
	Players map[string]*PlayerDTO `json:"players"`
	Elapsed float64               `json:"elapsed"`
}

type HexDTO struct {
	Q        int  `json:"q"`
	R        int  `json:"r"`
	Owner    int  `json:"owner"`
	Building int  `json:"building"`
	Level    int  `json:"level"`
	Capital  bool `json:"capital"`
}

type PlayerDTO struct {
	ID   int     `json:"id"`
	Gold float64 `json:"gold"`
	TP   float64 `json:"tp"`
}

func BuildSnapshot(state *game.GameState) *SnapshotMsg {
	msg := &SnapshotMsg{
		Type:    MsgSnapshot,
		Hexes:   make(map[string]*HexDTO, len(state.Hexes)),
		Players: make(map[string]*PlayerDTO, len(state.Players)),
		Elapsed: state.Elapsed,
	}

	for hex, hs := range state.Hexes {
		key := hexKey(hex)
		msg.Hexes[key] = &HexDTO{
			Q:        hex.Q,
			R:        hex.R,
			Owner:    int(hs.Owner),
			Building: int(hs.Building),
			Level:    hs.Level,
			Capital:  hs.Capital,
		}
	}

	for pid, p := range state.Players {
		key := playerKey(pid)
		msg.Players[key] = &PlayerDTO{
			ID:   int(p.ID),
			Gold: p.Gold,
			TP:   p.TP,
		}
	}

	return msg
}

func hexKey(h game.Hex) string {
	return fmt.Sprintf("%d,%d", h.Q, h.R)
}

func playerKey(pid game.PlayerID) string {
	return fmt.Sprintf("%d", pid)
}
