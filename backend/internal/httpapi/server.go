package httpapi

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"one-time-link/backend/internal/config"
	"one-time-link/backend/internal/secret"

	"github.com/google/uuid"
)

type Server struct {
	config        config.Config
	secretService secret.Service
}

type contextKey string

const requestIDKey contextKey = "request_id"

func NewServer(cfg config.Config, secretService secret.Service) *Server {
	return &Server{
		config:        cfg,
		secretService: secretService,
	}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", s.handleHealth)
	mux.HandleFunc("/api/secrets", s.handleCreateSecret)
	mux.HandleFunc("/api/secrets/", s.handleSecretRoutes)
	mux.HandleFunc("/api/reveal-sessions", s.handleCreateRevealSession)

	return withRequestID(
		withCORS(s.config.AllowedOrigin,
			withJSONHeaders(
				withRequestSizeLimit(15*1024, // 15KB limit
					withRequestLogging(mux)))))
}

func withRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		w.Header().Set("X-Request-ID", requestID)
		ctx := context.WithValue(r.Context(), requestIDKey, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func withCORS(allowedOrigin string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if allowedOrigin != "" {
			w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
			w.Header().Set("Vary", "Origin")
		}

		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept, X-Request-ID")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Expose-Headers", "X-Request-ID")

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

func withRequestSizeLimit(maxBytes int64, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
			r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
		}
		next.ServeHTTP(w, r)
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func withRequestLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap response writer to capture status code
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(rw, r)

		duration := time.Since(start)
		requestID := r.Context().Value(requestIDKey)
		if requestID == nil {
			requestID = ""
		}

		// Hash IP and User-Agent for privacy
		ipHash := hashString(r.RemoteAddr)
		uaHash := hashString(r.UserAgent())

		logEntry := map[string]interface{}{
			"timestamp":       time.Now().UTC().Format(time.RFC3339),
			"level":           "info",
			"event":           "http_request",
			"request_id":      requestID,
			"method":          r.Method,
			"path":            r.URL.Path,
			"status":          rw.statusCode,
			"duration_ms":     duration.Milliseconds(),
			"ip_hash":         ipHash,
			"user_agent_hash": uaHash,
		}

		logJSON, _ := json.Marshal(logEntry)
		log.Println(string(logJSON))
	})
}

func hashString(s string) string {
	if s == "" {
		return ""
	}
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:8]) // Use first 8 bytes for brevity
}

func getRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDKey).(string); ok {
		return id
	}
	return ""
}
