package main

import (
	"log"
	"net/http"

	"go-score-service/internal/handler"
	"go-score-service/internal/store"
)

func main() {
	store := store.NewMemoryStore()
	mux := handler.NewRouter(store)
	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	log.Println("go-score-service starting on :8080")
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()
	select {}
}
