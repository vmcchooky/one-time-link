package main

import (
	"log"
	"net/http"

	"one-time-link/backend/internal/config"
	"one-time-link/backend/internal/httpapi"
	"one-time-link/backend/internal/secret"
)

func main() {
	cfg := config.Load()
	secretService := secret.NewInMemoryService()
	server := httpapi.NewServer(cfg, secretService)

	log.Printf("starting %s on %s", cfg.ServiceName, cfg.ListenAddress())

	if err := http.ListenAndServe(cfg.ListenAddress(), server.Handler()); err != nil {
		log.Fatal(err)
	}
}
