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
	Units         []*UnitDTO   `json:"units"`
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
	Type     MsgType `json:"type"`
	Action   string  `json:"action"`
	Q        int     `json:"q"`
	R        int     `json:"r"`
	CardType int     `json:"cardType"`
}

type SnapshotMsg struct {
	Type          MsgType               `json:"type"`
	Hexes         map[string]*HexDTO    `json:"hexes"`
	Players       map[string]*PlayerDTO `json:"players"`
	Units         []*UnitDTO            `json:"units"`
	Elapsed       float64               `json:"elapsed"`
	Over          bool                  `json:"over"`
	Winner        int                   `json:"winner"`
	WinReason     string                `json:"winReason"`
	Waiting       bool                  `json:"waiting"`
	Paused        bool                  `json:"paused"`
	PauseTimeLeft float64               `json:"pauseTimeLeft"`
}

type HexDTO struct {
	Q           int  `json:"q"`
	R           int  `json:"r"`
	Owner       int  `json:"owner"`
	Capital     bool `json:"capital"`
	HasBuilding bool `json:"hasBuilding"`
}

type CardDTO struct {
	Type int `json:"type"`
}

type BarBoostDTO struct {
	Bonus    float64 `json:"bonus"`
	TimeLeft float64 `json:"timeLeft"`
}

type DeckDTO struct {
	Cards     []CardDTO `json:"cards"`
	AsideCard *CardDTO  `json:"asideCard"`
	AutoTimer float64   `json:"autoTimer"`
}

type PlayerDTO struct {
	ID        int           `json:"id"`
	Bar       float64       `json:"bar"`
	BarBoosts []BarBoostDTO `json:"barBoosts"`
	Deck      DeckDTO       `json:"deck"`
	CapitalHP int           `json:"capitalHp"`
}

type UnitDTO struct {
	ID    int `json:"id"`
	Owner int `json:"owner"`
	Type  int `json:"type"`
	Q     int `json:"q"`
	R     int `json:"r"`
	HP    int `json:"hp"`
}

func BuildSnapshot(state *game.GameState) *SnapshotMsg {
	msg := &SnapshotMsg{
		Type:          MsgSnapshot,
		Hexes:         make(map[string]*HexDTO, len(state.Hexes)),
		Players:       make(map[string]*PlayerDTO, len(state.Players)),
		Units:         make([]*UnitDTO, 0, len(state.Units)),
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
			Q:           hex.Q,
			R:           hex.R,
			Owner:       int(hs.Owner),
			Capital:     hs.Capital,
			HasBuilding: hs.HasBuilding,
		}
	}

	for pid, p := range state.Players {
		key := playerKey(pid)
		msg.Players[key] = playerToDTO(p)
	}

	for _, u := range state.Units {
		msg.Units = append(msg.Units, &UnitDTO{
			ID:    u.ID,
			Owner: int(u.Owner),
			Type:  int(u.Type),
			Q:     u.Pos.Q,
			R:     u.Pos.R,
			HP:    u.HP,
		})
	}

	return msg
}

func playerToDTO(p *game.Player) *PlayerDTO {
	cards := make([]CardDTO, len(p.Deck.Cards))
	for i, c := range p.Deck.Cards {
		cards[i] = CardDTO{Type: int(c.Type)}
	}
	var aside *CardDTO
	if p.Deck.AsideCard != nil {
		d := CardDTO{Type: int(p.Deck.AsideCard.Type)}
		aside = &d
	}
	boosts := make([]BarBoostDTO, len(p.BarBoosts))
	for i, b := range p.BarBoosts {
		boosts[i] = BarBoostDTO{Bonus: b.Bonus, TimeLeft: b.TimeLeft}
	}
	return &PlayerDTO{
		ID:        int(p.ID),
		Bar:       p.Bar,
		BarBoosts: boosts,
		Deck: DeckDTO{
			Cards:     cards,
			AsideCard: aside,
			AutoTimer: p.Deck.AutoTimer,
		},
		CapitalHP: p.CapitalHP,
	}
}

func hexKey(h game.Hex) string {
	return fmt.Sprintf("%d,%d", h.Q, h.R)
}

func playerKey(pid game.PlayerID) string {
	return fmt.Sprintf("%d", pid)
}
