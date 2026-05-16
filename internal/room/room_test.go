package room

import (
	"sync"
	"testing"
	"time"

	game "hexar/internal/game"
)

type MockClientSender struct {
	mu        sync.Mutex
	snapshots []*game.GameState
}

func (m *MockClientSender) SendSnapshot(state *game.GameState) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.snapshots = append(m.snapshots, state)
}

func newTestRoom(t *testing.T) *Room {
	room := New()
	if room == nil {
		t.Fatalf("failed to create room")
	}
	return room
}

func TestWaitingState(t *testing.T) {
	room := newTestRoom(t)
	pid1 := game.PlayerID(1)
	pid2 := game.PlayerID(2)

	client1 := &MockClientSender{}
	room.OnConnect(client1, pid1)

	if !room.state.Waiting {
		t.Error("expected Waiting=true after 1st player connects")
	}
	if room.started {
		t.Error("expected game loop not started with only 1 player")
	}
	if len(client1.snapshots) == 0 {
		t.Error("expected snapshot sent to player 1")
	}

	client2 := &MockClientSender{}
	room.OnConnect(client2, pid2)

	if room.state.Waiting {
		t.Error("expected Waiting=false after both players connect")
	}
	if !room.started {
		t.Error("expected game loop started after both players connect")
	}
	if len(client1.snapshots) < 2 {
		t.Error("expected snapshot sent to player 1 when player 2 joins")
	}
}

func TestWaitingStateNoGoldAccrual(t *testing.T) {
	room := newTestRoom(t)
	pid1 := game.PlayerID(1)

	client1 := &MockClientSender{}
	room.OnConnect(client1, pid1)

	goldBefore := room.state.Players[pid1].Gold
	time.Sleep(150 * time.Millisecond)
	goldAfter := room.state.Players[pid1].Gold

	if goldAfter != goldBefore {
		t.Errorf("gold changed while waiting: %.2f → %.2f (expected no change)", goldBefore, goldAfter)
	}
}

func TestPauseOnDisconnect(t *testing.T) {
	room := newTestRoom(t)
	pid1 := game.PlayerID(1)
	pid2 := game.PlayerID(2)

	client1 := &MockClientSender{}
	client2 := &MockClientSender{}
	room.OnConnect(client1, pid1)
	room.OnConnect(client2, pid2)

	time.Sleep(150 * time.Millisecond)

	room.OnDisconnect(pid1, client1)

	if !room.state.Paused {
		t.Error("expected Paused=true after disconnect")
	}
	if room.state.PauseTimeLeft <= 0 || room.state.PauseTimeLeft > 120.5 {
		t.Errorf("expected PauseTimeLeft ~120s, got %.2f", room.state.PauseTimeLeft)
	}
	if len(client2.snapshots) < 2 {
		t.Error("expected snapshot sent to remaining player on disconnect")
	}
}

func TestUnpauseOnReconnect(t *testing.T) {
	room := newTestRoom(t)
	pid1 := game.PlayerID(1)
	pid2 := game.PlayerID(2)

	client1 := &MockClientSender{}
	client2 := &MockClientSender{}
	room.OnConnect(client1, pid1)
	room.OnConnect(client2, pid2)

	time.Sleep(150 * time.Millisecond)

	room.OnDisconnect(pid1, client1)
	if !room.state.Paused {
		t.Fatal("expected Paused=true after disconnect")
	}

	client1New := &MockClientSender{}
	room.OnConnect(client1New, pid1)

	if room.state.Paused {
		t.Error("expected Paused=false after reconnect")
	}
	if room.state.PauseTimeLeft != 0 {
		t.Errorf("expected PauseTimeLeft=0 after unpause, got %.2f", room.state.PauseTimeLeft)
	}
	if _, ok := room.remainingGrace[pid1]; !ok {
		t.Error("expected remainingGrace to be saved for pid1 after reconnect")
	}
}

func TestCumulativeGrace(t *testing.T) {
	room := newTestRoom(t)
	pid1 := game.PlayerID(1)
	pid2 := game.PlayerID(2)

	client1 := &MockClientSender{}
	client2 := &MockClientSender{}
	room.OnConnect(client1, pid1)
	room.OnConnect(client2, pid2)

	time.Sleep(150 * time.Millisecond)

	// First disconnect: let grace count down briefly
	room.OnDisconnect(pid1, client1)
	time.Sleep(300 * time.Millisecond) // pause ticks decrement PauseTimeLeft

	// Reconnect: saves remaining grace (should be < 120s)
	client1New := &MockClientSender{}
	room.OnConnect(client1New, pid1)
	savedGrace := room.remainingGrace[pid1]

	if savedGrace >= 120.0 {
		t.Errorf("expected savedGrace < 120s after countdown, got %.2f", savedGrace)
	}

	// Second disconnect: budget must use saved grace, not reset to 120s
	room.OnDisconnect(pid1, client1New)
	pauseTime2 := room.state.PauseTimeLeft

	if pauseTime2 > savedGrace+0.2 {
		t.Errorf("2nd disconnect reset grace to %.2f instead of using saved %.2f", pauseTime2, savedGrace)
	}
}

func TestForfeitEnqueued(t *testing.T) {
	room := newTestRoom(t)
	pid1 := game.PlayerID(1)
	pid2 := game.PlayerID(2)

	client1 := &MockClientSender{}
	client2 := &MockClientSender{}
	room.OnConnect(client1, pid1)
	room.OnConnect(client2, pid2)

	time.Sleep(150 * time.Millisecond)

	room.OnDisconnect(pid1, client1)

	room.mu.Lock()
	room.state.PauseTimeLeft = 0.15
	room.mu.Unlock()

	time.Sleep(300 * time.Millisecond)

	if !room.state.Over {
		t.Errorf("expected game Over=true after forfeit timeout, got false")
	}
	if room.state.WinReason != "forfeit" {
		t.Errorf("expected WinReason='forfeit', got '%s'", room.state.WinReason)
	}
}

// TestDuplicateConnectReplaces verifies that a second OnConnect for the same player
// replaces the existing client (evicts the old connection, registers the new one).
func TestDuplicateConnectReplaces(t *testing.T) {
	room := newTestRoom(t)
	pid1 := game.PlayerID(1)
	pid2 := game.PlayerID(2)

	client1 := &MockClientSender{}
	client2 := &MockClientSender{}
	room.OnConnect(client1, pid1)
	room.OnConnect(client2, pid2)

	time.Sleep(150 * time.Millisecond)

	dupClient := &MockClientSender{}
	room.OnConnect(dupClient, pid1)

	// New client is now the active one
	if room.activeClients[pid1] != dupClient {
		t.Error("expected duplicate connection to become the active client")
	}
	// Old client is evicted
	if _, ok := room.clientPlayer[client1]; ok {
		t.Error("expected original client to be evicted from clientPlayer map")
	}
	// New client received a snapshot
	if len(dupClient.snapshots) == 0 {
		t.Error("expected new client to receive a snapshot on connect")
	}
}

func TestDuplicateConnectDuringPause(t *testing.T) {
	room := newTestRoom(t)
	pid1 := game.PlayerID(1)
	pid2 := game.PlayerID(2)

	client1 := &MockClientSender{}
	client2 := &MockClientSender{}
	room.OnConnect(client1, pid1)
	room.OnConnect(client2, pid2)

	time.Sleep(150 * time.Millisecond)

	room.OnDisconnect(pid1, client1)

	reconnectClient := &MockClientSender{}
	room.OnConnect(reconnectClient, pid1)

	if !room.IsConnected(pid1) {
		t.Error("expected player 1 to be connected after reconnect")
	}
	if room.state.Paused {
		t.Error("expected Paused=false after reconnect")
	}
}
