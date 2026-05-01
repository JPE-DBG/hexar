package net

import (
	"context"
	"encoding/json"
	"hexar/internal/room"
	"log"
	"net/http"

	"github.com/coder/websocket"
)

type Server struct {
	room *room.Room
	mux  *http.ServeMux
}

func NewServer(r *room.Room, clientDir string) *Server {
	s := &Server{room: r, mux: http.NewServeMux()}

	s.mux.Handle("/", http.FileServer(http.Dir(clientDir)))
	s.mux.HandleFunc("/ws", s.handleWS)

	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		log.Printf("ws accept error: %v", err)
		return
	}

	ctx := r.Context()
	client := NewClient(conn, s.room)
	pid := s.room.AddClient(client)
	client.playerID = pid

	welcome, _ := json.Marshal(WelcomeMsg{Type: MsgWelcome, PlayerID: int(pid)})
	conn.Write(ctx, websocket.MessageText, welcome)

	go client.WritePump(ctx)

	client.ReadPump(ctx)

	s.room.RemoveClient(client)
	close(client.send)
	conn.Close(websocket.StatusNormalClosure, "")
}

func ListenAndServe(addr string, srv *Server) error {
	log.Printf("listening on %s", addr)
	return http.ListenAndServe(addr, srv)
}

func ListenAndServeContext(ctx context.Context, addr string, srv *Server) error {
	httpSrv := &http.Server{Addr: addr, Handler: srv}
	go func() {
		<-ctx.Done()
		httpSrv.Close()
	}()
	log.Printf("listening on %s", addr)
	return httpSrv.ListenAndServe()
}
