package httpapi

import (
	"log"
	"net/http"
	"time"

	"one-time-link/backend/internal/config"
	"one-time-link/backend/internal/secret"
)

type Server struct {
	config        config.Config
	secretService secret.Service
}

func NewServer(cfg config.Config, secretService secret.Service) *Server {
	return &Server{
		config:        cfg,
		secretService: secretService,
	}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", s.handleHealth)
	mux.HandleFunc("POST /api/secrets", s.handleCreateSecret)
	mux.HandleFunc("GET /api/secrets/", s.handleSecretRoutes)
	mux.HandleFunc("POST /api/reveal-sessions", s.handleCreateRevealSession)

	return withCORS(s.config.AllowedOrigin, withJSONHeaders(withRequestLogging(mux)))
}

func withCORS(allowedOrigin string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if allowedOrigin != "" {
			w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
			w.Header().Set("Vary", "Origin")
		}

		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func withJSONHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		next.ServeHTTP(w, r)
	})
}

func withRequestLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf(
			`{"event":"http_request","method":"%s","path":"%s","remote_addr":"%s","duration_ms":%d}`,
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
			time.Since(start).Milliseconds(),
		)
	})
}
