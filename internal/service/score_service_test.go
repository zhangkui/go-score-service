package service

import (
	"context"
	"testing"

	"go-score-service/internal/store"
)

func TestServiceRegister(t *testing.T) {
	svc := NewScoreService(store.NewMemoryStore())
	if err := svc.Register(context.Background(), "u1", "alice"); err != nil {
		t.Fatalf("register failed: %v", err)
	}
}

func TestServiceAddAndGetScore(t *testing.T) {
	svc := NewScoreService(store.NewMemoryStore())
	svc.Register(context.Background(), "u1", "alice")
	score, err := svc.AddScore(context.Background(), "u1", 100)
	if err != nil {
		t.Fatalf("add score failed: %v", err)
	}
	if score != 100 {
		t.Fatalf("expected 100, got %d", score)
	}
	got, err := svc.GetScore(context.Background(), "u1")
	if err != nil {
		t.Fatalf("get score failed: %v", err)
	}
	if got != 100 {
		t.Fatalf("expected 100, got %d", got)
	}
}

func TestServiceDeductScore(t *testing.T) {
	svc := NewScoreService(store.NewMemoryStore())
	svc.Register(context.Background(), "u1", "alice")
	svc.AddScore(context.Background(), "u1", 100)

	// Deducting 50 must decrease the score, not increase it.
	score, err := svc.DeductScore(context.Background(), "u1", 50)
	if err != nil {
		t.Fatalf("deduct score failed: %v", err)
	}
	if score != 50 {
		t.Fatalf("expected score 50 after deduction, got %d", score)
	}

	got, err := svc.GetScore(context.Background(), "u1")
	if err != nil {
		t.Fatalf("get score failed: %v", err)
	}
	if got != 50 {
		t.Fatalf("expected persisted score 50, got %d", got)
	}
}

func TestServiceDeductScoreNotFound(t *testing.T) {
	svc := NewScoreService(store.NewMemoryStore())
	if _, err := svc.DeductScore(context.Background(), "missing", 50); err != store.ErrUserNotFound {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestServiceGetUser(t *testing.T) {
	svc := NewScoreService(store.NewMemoryStore())
	svc.Register(context.Background(), "u1", "alice")
	svc.AddScore(context.Background(), "u1", 50)
	user, err := svc.GetUser(context.Background(), "u1")
	if err != nil {
		t.Fatalf("get user failed: %v", err)
	}
	if user.Name != "alice" || user.Score != 50 {
		t.Fatalf("unexpected user: %+v", user)
	}
}

func TestServiceLeaderboard(t *testing.T) {
	svc := NewScoreService(store.NewMemoryStore())
	svc.Register(context.Background(), "u1", "alice")
	svc.Register(context.Background(), "u2", "bob")
	svc.AddScore(context.Background(), "u1", 100)
	svc.AddScore(context.Background(), "u2", 200)
	board := svc.Leaderboard(context.Background(), 0)
	if len(board) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(board))
	}
	if board[0].Score != 200 {
		t.Fatalf("expected 200 first, got %d", board[0].Score)
	}
}
