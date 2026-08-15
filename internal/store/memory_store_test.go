package store

import (
	"sync"
	"testing"
)

func TestRegister(t *testing.T) {
	s := NewMemoryStore()
	if err := s.Register("u1", "alice"); err != nil {
		t.Fatalf("register failed: %v", err)
	}
	if err := s.Register("u1", "alice"); err != ErrUserAlreadyExists {
		t.Fatalf("expected ErrUserAlreadyExists, got %v", err)
	}
}

func TestAddScore(t *testing.T) {
	s := NewMemoryStore()
	s.Register("u1", "alice")
	score, err := s.AddScore("u1", 100)
	if err != nil {
		t.Fatalf("add score failed: %v", err)
	}
	if score != 100 {
		t.Fatalf("expected 100, got %d", score)
	}
	score, _ = s.AddScore("u1", -30)
	if score != 70 {
		t.Fatalf("expected 70, got %d", score)
	}
}

func TestGetScore(t *testing.T) {
	s := NewMemoryStore()
	s.Register("u1", "alice")
	s.AddScore("u1", 50)
	score, err := s.GetScore("u1")
	if err != nil {
		t.Fatalf("get score failed: %v", err)
	}
	if score != 50 {
		t.Fatalf("expected 50, got %d", score)
	}
}

func TestLeaderboard(t *testing.T) {
	s := NewMemoryStore()
	s.Register("u1", "alice")
	s.Register("u2", "bob")
	s.Register("u3", "carol")
	s.AddScore("u1", 100)
	s.AddScore("u2", 200)
	s.AddScore("u3", 50)
	board := s.Leaderboard(2)
	if len(board) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(board))
	}
	if board[0].UserID != "u2" {
		t.Fatalf("expected u2 first, got %s", board[0].UserID)
	}
}

func TestGetUserReturnsIndependentCopy(t *testing.T) {
	s := NewMemoryStore()
	s.Register("u1", "alice")
	s.AddScore("u1", 50)

	user, err := s.GetUser("u1")
	if err != nil {
		t.Fatalf("get user failed: %v", err)
	}
	if user.Name != "alice" || user.Score != 50 {
		t.Fatalf("unexpected user: %+v", user)
	}

	// Mutate the returned user; the store's data must not change.
	user.Name = "mallory"
	user.Score = 9999
	user.ID = "evil"

	stored, err := s.GetUser("u1")
	if err != nil {
		t.Fatalf("second get user failed: %v", err)
	}
	if stored.ID != "u1" {
		t.Fatalf("stored ID changed: got %q, want %q", stored.ID, "u1")
	}
	if stored.Name != "alice" {
		t.Fatalf("stored Name changed: got %q, want %q", stored.Name, "alice")
	}
	if stored.Score != 50 {
		t.Fatalf("stored Score changed: got %d, want %d", stored.Score, 50)
	}

	// A second mutation of the first returned copy must also be isolated.
	user.Score = -1
	stored2, _ := s.GetUser("u1")
	if stored2.Score != 50 {
		t.Fatalf("stored Score changed on second mutation: got %d, want %d", stored2.Score, 50)
	}
}

func TestConcurrentAddScore(t *testing.T) {
	s := NewMemoryStore()
	s.Register("u1", "alice")
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.AddScore("u1", 1)
		}()
	}
	wg.Wait()
	score, _ := s.GetScore("u1")
	if score != 100 {
		t.Fatalf("expected 100, got %d", score)
	}
}
