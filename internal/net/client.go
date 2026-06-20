package net

import (
	"context"
	"encoding/json"
	"hexar/internal/game"
	"hexar/internal/room"
	"log"

	"github.com/coder/websocket"
)

const (
	sendBufferSize = 64  // outbound message channel capacity per client
	maxDeltaBytes  = 500 // M6 success criterion: steady-state delta must stay under this
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
		send: make(chan []byte, sendBufferSize),
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
	if _, isDelta := msg.(*DeltaMsg); isDelta && len(data) > maxDeltaBytes {
		log.Printf("delta over budget: %d bytes", len(data))
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
			CardType int    `json:"cardType"`
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
		case "play-top":
			action = game.Action{Type: game.ActionPlayTopCard, Player: c.playerID}
		case "play-aside":
			action = game.Action{Type: game.ActionPlayAsideCard, Player: c.playerID}
		case "push-aside":
			action = game.Action{Type: game.ActionPushAside, Player: c.playerID}
		case "buy-card":
			action = game.Action{Type: game.ActionBuyCard, Player: c.playerID, CardType: game.CardType(raw.CardType)}
		case "remove-top":
			action = game.Action{Type: game.ActionRemoveTopCard, Player: c.playerID}
		case "remove-aside":
			action = game.Action{Type: game.ActionRemoveAsideCard, Player: c.playerID}
		case "claim-hex":
			action = game.Action{Type: game.ActionClaimHex, Player: c.playerID, Target: target}
		default:
			continue
		}

		c.room.EnqueueAction(action)
	}
}
