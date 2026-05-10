package net

import (
	"context"
	"encoding/json"
	"hexar/internal/game"
	"hexar/internal/room"
	"log"

	"github.com/coder/websocket"
)

type Client struct {
	conn         *websocket.Conn
	room         *room.Room
	playerID     game.PlayerID
	send         chan []byte
	prevSnapshot *SnapshotMsg
}

func NewClient(conn *websocket.Conn, r *room.Room) *Client {
	return &Client{
		conn: conn,
		room: r,
		send: make(chan []byte, 64),
	}
}

func (c *Client) SendSnapshot(state *game.GameState) {
	curr := BuildSnapshot(state)
	var msg any
	if c.prevSnapshot == nil {
		msg = curr
	} else {
		msg = buildDelta(c.prevSnapshot, curr)
	}
	c.prevSnapshot = curr

	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("marshal error: %v", err)
		return
	}
	select {
	case c.send <- data:
	default:
	}
}

func (c *Client) WritePump(ctx context.Context) {
	for {
		select {
		case msg, ok := <-c.send:
			if !ok {
				return
			}
			err := c.conn.Write(ctx, websocket.MessageText, msg)
			if err != nil {
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

func (c *Client) ReadPump(ctx context.Context) {
	for {
		_, data, err := c.conn.Read(ctx)
		if err != nil {
			return
		}
		var raw struct {
			Type     string `json:"type"`
			Action   string `json:"action"`
			Q        int    `json:"q"`
			R        int    `json:"r"`
			Building string `json:"building"`
			TechID   int    `json:"techId"`
		}
		if json.Unmarshal(data, &raw) != nil {
			continue
		}
		if raw.Type != string(MsgAction) {
			continue
		}

		target := game.Hex{Q: raw.Q, R: raw.R}
		var action game.Action

		switch raw.Action {
		case "claim":
			action = game.Action{Type: game.ActionClaim, Player: c.playerID, Target: target}
		case "upgrade":
			bt := parseBuildingType(raw.Building)
			action = game.Action{Type: game.ActionUpgrade, Player: c.playerID, Target: target, Building: bt}
		case "demolish":
			action = game.Action{Type: game.ActionDemolish, Player: c.playerID, Target: target}
		case "attack":
			action = game.Action{Type: game.ActionAttack, Player: c.playerID, Target: target}
		case "counter-spend":
			action = game.Action{Type: game.ActionCounterSpend, Player: c.playerID, Target: target}
		case "unlock-tech":
			action = game.Action{Type: game.ActionUnlockTech, Player: c.playerID, TechID: game.TechID(raw.TechID)}
		case "drop-hex":
			action = game.Action{Type: game.ActionDropHex, Player: c.playerID, Target: target}
		case "fortify":
			action = game.Action{Type: game.ActionFortify, Player: c.playerID, Target: target}
		default:
			continue
		}

		c.room.EnqueueAction(action)
	}
}

func parseBuildingType(s string) game.BuildingType {
	switch s {
	case "gold":
		return game.BuildingGold
	case "power":
		return game.BuildingPower
	case "research":
		return game.BuildingResearch
	default:
		return game.BuildingNone
	}
}
