package store

import (
	"sync"

	"go-score-service/internal/model"
)

type MemoryStore struct {
	mu    sync.RWMutex
	users map[string]*model.User
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		users: make(map[string]*model.User),
	}
}

func (s *MemoryStore) Register(id, name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.users[id]; exists {
		return ErrUserAlreadyExists
	}
	s.users[id] = &model.User{ID: id, Name: name, Score: 0}
	return nil
}

func (s *MemoryStore) AddScore(userID string, delta int64) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	user, exists := s.users[userID]
	if !exists {
		return 0, ErrUserNotFound
	}
	user.Score += delta
	return user.Score, nil
}

func (s *MemoryStore) GetScore(userID string) (int64, error) {
	user, exists := s.users[userID]
	if !exists {
		return 0, ErrUserNotFound
	}
	return user.Score, nil
}

func (s *MemoryStore) GetUser(userID string) (*model.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	user, exists := s.users[userID]
	if !exists {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (s *MemoryStore) Leaderboard(limit int) []model.ScoreEntry {
	entries := make([]model.ScoreEntry, 0, len(s.users))
	for _, u := range s.users {
		entries = append(entries, model.ScoreEntry{UserID: u.ID, Score: u.Score})
	}
	sortDesc(entries)
	if limit > 0 && limit < len(entries) {
		entries = entries[:limit]
	}
	return entries
}

func sortDesc(entries []model.ScoreEntry) {
	for i := 1; i < len(entries); i++ {
		for j := i; j > 0 && entries[j].Score >= entries[j-1].Score; j-- {
			entries[j], entries[j-1] = entries[j-1], entries[j]
		}
	}
}
