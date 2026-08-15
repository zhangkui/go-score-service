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

func TestLeaderboardOrdersTiedScoresByUserID(t *testing.T) {
	s := NewMemoryStore()
	for _, userID := range []string{"u3", "u1", "u2", "u4"} {
		if err := s.Register(userID, userID); err != nil {
			t.Fatalf("register %s failed: %v", userID, err)
		}
	}
	for _, userID := range []string{"u1", "u2", "u3"} {
		if _, err := s.AddScore(userID, 100); err != nil {
			t.Fatalf("add score for %s failed: %v", userID, err)
		}
	}

	want := []string{"u1", "u2", "u3"}
	for attempt := 0; attempt < 100; attempt++ {
		board := s.Leaderboard(3)
		for index, userID := range want {
			if board[index].UserID != userID {
				t.Fatalf("attempt %d: entry %d: expected %s, got %s", attempt, index, userID, board[index].UserID)
			}
		}
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
