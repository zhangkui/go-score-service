package service

import (
	"context"

	"go-score-service/internal/model"
	"go-score-service/internal/store"
)

type ScoreService struct {
	store *store.MemoryStore
}

func NewScoreService(s *store.MemoryStore) *ScoreService {
	return &ScoreService{store: s}
}

func (svc *ScoreService) Register(ctx context.Context, id, name string) error {
	return svc.store.Register(id, name)
}

func (svc *ScoreService) AddScore(ctx context.Context, userID string, delta int64) (int64, error) {
	return svc.store.AddScore(userID, delta)
}

func (svc *ScoreService) DeductScore(ctx context.Context, userID string, amount int64) (int64, error) {
	return svc.store.AddScore(userID, amount)
}

func (svc *ScoreService) GetScore(ctx context.Context, userID string) (int64, error) {
	return svc.store.GetScore(userID)
}

func (svc *ScoreService) GetUser(ctx context.Context, userID string) (*model.User, error) {
	return svc.store.GetUser(userID)
}

func (svc *ScoreService) Leaderboard(ctx context.Context, limit int) []model.ScoreEntry {
	return svc.store.Leaderboard(limit)
}
