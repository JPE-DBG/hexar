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
	conn     *websocket.Conn
	room     *room.Room
	playerID game.PlayerID
	send     chan []byte
}

func NewClient(conn *websocket.Conn, r *room.Room) *Client {
	return &Client{
		conn: conn,
		room: r,
		send: make(chan []byte, 64),
	}
}

func (c *Client) SendSnapshot(state *game.GameState) {
	msg := BuildSnapshot(state)
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
			Type   string `json:"type"`
			Action string `json:"action"`
			Q      int    `json:"q"`
			R      int    `json:"r"`
		}
		if json.Unmarshal(data, &raw) != nil {
			continue
		}
		if raw.Type == string(MsgAction) && raw.Action == "claim" {
			c.room.EnqueueAction(game.ClaimAction{
				Player: c.playerID,
				Target: game.Hex{Q: raw.Q, R: raw.R},
			})
		}
	}
}
