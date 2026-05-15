package net

import (
	"context"
	"encoding/json"
	"hexar/internal/lobby"
	"log"
	"net/http"

	"github.com/coder/websocket"
)

type Server struct {
	lobby   *lobby.Lobby
	mux     *http.ServeMux
	version string
}

func NewServer(lob *lobby.Lobby, clientDir string, version string) *Server {
	s := &Server{lobby: lob, mux: http.NewServeMux(), version: version}
	s.mux.Handle("/", http.FileServer(http.Dir(clientDir)))
	s.mux.HandleFunc("/health", s.handleHealth)
	s.mux.HandleFunc("/lobby/create", s.handleCreate)
	s.mux.HandleFunc("/lobby/join", s.handleJoin)
	s.mux.HandleFunc("/ws", s.handleWS)
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok", "version": s.version})
}

func (s *Server) handleCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST required", http.StatusMethodNotAllowed)
		return
	}
	code, token, pid := s.lobby.Create()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"code":     code,
		"token":    token,
		"playerId": int(pid),
	})
}

func (s *Server) handleJoin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST required", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Code string `json:"code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	token, pid, err := s.lobby.Join(req.Code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"token":    token,
		"playerId": int(pid),
	})
}

func (s *Server) handleWS(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	token := r.URL.Query().Get("token")

	entry, ok := s.lobby.Get(code)
	if !ok {
		http.Error(w, "room not found", http.StatusNotFound)
		return
	}
	pid, ok := entry.Validate(token)
	if !ok {
		http.Error(w, "invalid token", http.StatusForbidden)
		return
	}

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		log.Printf("ws accept error: %v", err)
		return
	}

	ctx := r.Context()
	rm := entry.Room

	if rm.IsConnected(pid) {
		conn.Close(websocket.StatusCode(4001), "already connected in another tab")
		return
	}

	client := NewClient(conn, rm)
	client.playerID = pid

	rm.OnConnect(client, pid)

	welcome, _ := json.Marshal(WelcomeMsg{Type: MsgWelcome, PlayerID: int(pid)})
	conn.Write(ctx, websocket.MessageText, welcome)

	go client.WritePump(ctx)
	client.ReadPump(ctx)

	rm.OnDisconnect(pid, client)
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
