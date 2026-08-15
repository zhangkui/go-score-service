package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"go-score-service/internal/service"
	"go-score-service/internal/store"
)

func NewRouter(s *store.MemoryStore) http.Handler {
	svc := service.NewScoreService(s)
	mux := http.NewServeMux()
	mux.HandleFunc("/register", handleRegister(svc))
	mux.HandleFunc("/score/add", handleAddScore(svc))
	mux.HandleFunc("/score/deduct", handleDeductScore(svc))
	mux.HandleFunc("/score/get", handleGetScore(svc))
	mux.HandleFunc("/user/get", handleGetUser(svc))
	mux.HandleFunc("/leaderboard", handleLeaderboard(svc))
	mux.HandleFunc("/health", handleHealthCheck(svc))
	return mux
}

func handleRegister(svc *service.ScoreService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid request body")
			return
		}
		if req.ID == "" {
			writeError(w, http.StatusBadRequest, "id is required")
			return
		}
		if err := svc.Register(r.Context(), req.ID, req.Name); err != nil {
			writeError(w, http.StatusConflict, err.Error())
			return
		}
		writeJSON(w, http.StatusCreated, map[string]string{"status": "ok", "id": req.ID})
	}
}

func handleAddScore(svc *service.ScoreService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			UserID string `json:"user_id"`
			Delta  int64  `json:"delta"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid request body")
			return
		}
		score, err := svc.AddScore(r.Context(), req.UserID, req.Delta)
		if err != nil {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]int64{"user_id": 0, "score": score})
	}
}

func handleDeductScore(svc *service.ScoreService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			UserID string `json:"user_id"`
			Amount int64  `json:"amount"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid request body")
			return
		}
		score, err := svc.DeductScore(r.Context(), req.UserID, req.Amount)
		if err != nil {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]int64{"user_id": 0, "score": score})
	}
}

func handleGetScore(svc *service.ScoreService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user_id")
		if userID == "" {
			writeError(w, http.StatusBadRequest, "user_id is required")
			return
		}
		score, err := svc.GetScore(r.Context(), userID)
		if err != nil {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]int64{"user_id": 0, "score": score})
	}
}

func handleGetUser(svc *service.ScoreService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user_id")
		if userID == "" {
			writeError(w, http.StatusBadRequest, "user_id is required")
			return
		}
		user, err := svc.GetUser(r.Context(), userID)
		if err != nil {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, user)
	}
}

func handleHealthCheck(svc *service.ScoreService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "healthy"})
	}
}

func handleLeaderboard(svc *service.ScoreService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limitStr := r.URL.Query().Get("limit")
		limit := 0
		if limitStr != "" {
			if n, err := strconv.Atoi(limitStr); err == nil && n > 0 {
				limit = n
			}
		}
		entries := svc.Leaderboard(r.Context(), limit)
		writeJSON(w, http.StatusOK, entries)
	}
}

func writeJSON(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
