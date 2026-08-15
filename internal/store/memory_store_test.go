package store

import (
	"sync"
	"testing"

	"go-score-service/internal/model"
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

func TestLeaderboardTieOrderStable(t *testing.T) {
	s := NewMemoryStore()
	// Register out of lexicographic order so the result cannot depend on
	// registration order or Go's randomized map iteration order.
	s.Register("u3", "carol")
	s.Register("u1", "alice")
	s.Register("u2", "bob")
	// Identical scores → all three users are tied.
	s.AddScore("u1", 100)
	s.AddScore("u2", 100)
	s.AddScore("u3", 100)

	first := s.Leaderboard(0)
	if len(first) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(first))
	}
	// Tied users must be ordered deterministically by UserID ascending.
	want := []string{"u1", "u2", "u3"}
	for i, e := range first {
		if e.UserID != want[i] {
			t.Fatalf("entry %d: expected %s, got %s", i, want[i], e.UserID)
		}
	}
	// Repeated queries must return the exact same ordering every time.
	for n := 0; n < 100; n++ {
		got := s.Leaderboard(0)
		if !equalScoreEntries(first, got) {
			t.Fatalf("leaderboard order changed on repeat %d:\n first=%v\n got  =%v", n, first, got)
		}
	}
}

func equalScoreEntries(a, b []model.ScoreEntry) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
