package lobby

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"hexar/internal/game"
	"hexar/internal/room"
	"strings"
	"sync"
)

type Session struct {
	PlayerID game.PlayerID
}

type Entry struct {
	Room     *room.Room
	mu       sync.Mutex
	sessions map[string]*Session
}

func (e *Entry) AddSession() (token string, pid game.PlayerID, err error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if len(e.sessions) >= 2 {
		return "", 0, errors.New("room full")
	}
	pid = game.PlayerID(len(e.sessions) + 1)
	token = randToken()
	e.sessions[token] = &Session{PlayerID: pid}
	return
}

func (e *Entry) Validate(token string) (game.PlayerID, bool) {
	e.mu.Lock()
	defer e.mu.Unlock()
	s, ok := e.sessions[token]
	if !ok {
		return 0, false
	}
	return s.PlayerID, true
}

type Lobby struct {
	mu    sync.Mutex
	rooms map[string]*Entry
}

func New() *Lobby {
	return &Lobby{rooms: make(map[string]*Entry)}
}

func (l *Lobby) Create() (code, token string, pid game.PlayerID) {
	l.mu.Lock()
	defer l.mu.Unlock()
	code = randCode()
	r := room.New()
	go r.Run()
	entry := &Entry{
		Room:     r,
		sessions: make(map[string]*Session),
	}
	l.rooms[code] = entry
	token, pid, _ = entry.AddSession()
	return
}

func (l *Lobby) Join(code string) (token string, pid game.PlayerID, err error) {
	l.mu.Lock()
	entry, ok := l.rooms[code]
	l.mu.Unlock()
	if !ok {
		return "", 0, errors.New("room not found")
	}
	return entry.AddSession()
}

func (l *Lobby) Get(code string) (*Entry, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	e, ok := l.rooms[code]
	return e, ok
}

func randCode() string {
	b := make([]byte, 2)
	rand.Read(b)
	return strings.ToUpper(hex.EncodeToString(b))
}

func randToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
