package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-score-service/internal/store"
)

func TestDeductScoreDecreasesScore(t *testing.T) {
	scoreStore := store.NewMemoryStore()
	if err := scoreStore.Register("u1", "alice"); err != nil {
		t.Fatalf("register user: %v", err)
	}
	if _, err := scoreStore.AddScore("u1", 100); err != nil {
		t.Fatalf("add score: %v", err)
	}
	router := NewRouter(scoreStore)

	deductRequest := httptest.NewRequest(http.MethodPost, "/score/deduct", bytes.NewBufferString(`{"user_id":"u1","amount":50}`))
	deductResponse := httptest.NewRecorder()
	router.ServeHTTP(deductResponse, deductRequest)
	if deductResponse.Code != http.StatusOK {
		t.Fatalf("expected deduct status %d, got %d", http.StatusOK, deductResponse.Code)
	}

	var deductBody map[string]int64
	if err := json.NewDecoder(deductResponse.Body).Decode(&deductBody); err != nil {
		t.Fatalf("decode deduct response: %v", err)
	}
	if deductBody["score"] != 50 {
		t.Fatalf("expected deducted score 50, got %d", deductBody["score"])
	}

	getRequest := httptest.NewRequest(http.MethodGet, "/score/get?user_id=u1", nil)
	getResponse := httptest.NewRecorder()
	router.ServeHTTP(getResponse, getRequest)
	if getResponse.Code != http.StatusOK {
		t.Fatalf("expected get status %d, got %d", http.StatusOK, getResponse.Code)
	}

	var getBody map[string]int64
	if err := json.NewDecoder(getResponse.Body).Decode(&getBody); err != nil {
		t.Fatalf("decode get response: %v", err)
	}
	if getBody["score"] != 50 {
		t.Fatalf("expected stored score 50, got %d", getBody["score"])
	}
}
