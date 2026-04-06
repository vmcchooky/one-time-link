package httpapi

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"
)

type healthResponse struct {
	Service   string `json:"service"`
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Version   string `json:"version"`
}

type errorResponse struct {
	Error     string `json:"error"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
	RequestID string `json:"request_id,omitempty"`
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeErrorJSON(w, r, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
		return
	}

	writeJSON(w, http.StatusOK, healthResponse{
		Service:   s.config.ServiceName,
		Status:    "healthy",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Version:   "1.0.0",
	})
}

func (s *Server) handleCreateSecret(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErrorJSON(w, r, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
		return
	}

	// Try to read body to trigger MaxBytesReader error if size exceeded
	_, err := io.ReadAll(r.Body)
	if err != nil {
		// MaxBytesReader returns an error when limit is exceeded
		writeErrorJSON(w, r, http.StatusRequestEntityTooLarge, "payload_too_large", "Request body exceeds 15KB limit")
		return
	}

	writeErrorJSON(w, r, http.StatusNotImplemented, "not_implemented", "Create secret endpoint is scaffolded but not implemented yet")
}

func (s *Server) handleCreateRevealSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErrorJSON(w, r, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
		return
	}

	writeErrorJSON(w, r, http.StatusNotImplemented, "not_implemented", "Reveal session endpoint is scaffolded but not implemented yet")
}

func (s *Server) handleSecretRoutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/secrets/")

	switch {
	case strings.HasSuffix(path, "/status") && r.Method == http.MethodGet:
		writeErrorJSON(w, r, http.StatusNotImplemented, "not_implemented", "Secret status endpoint is scaffolded but not implemented yet")
	case strings.HasSuffix(path, "/consume") && r.Method == http.MethodPost:
		writeErrorJSON(w, r, http.StatusNotImplemented, "not_implemented", "Secret consume endpoint is scaffolded but not implemented yet")
	default:
		writeErrorJSON(w, r, http.StatusNotFound, "not_found", "Route not found")
	}
}

func writeJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeErrorJSON(w http.ResponseWriter, r *http.Request, statusCode int, errorCode, message string) {
	requestID := getRequestID(r.Context())
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(errorResponse{
		Error:     errorCode,
		Message:   message,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		RequestID: requestID,
	})
}
