package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type healthResponse struct {
	Service string `json:"service"`
	Status  string `json:"status"`
	Time    string `json:"time"`
	Store   string `json:"store"`
	Mode    string `json:"mode"`
}

type errorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	status := s.secretService.Health(r.Context())

	writeJSON(w, http.StatusOK, healthResponse{
		Service: s.config.ServiceName,
		Status:  "ok",
		Time:    time.Now().UTC().Format(time.RFC3339),
		Store:   status.Store,
		Mode:    status.Mode,
	})
}

func (s *Server) handleCreateSecret(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusNotImplemented, errorResponse{
		Error:   "not_implemented",
		Message: "Create secret endpoint is scaffolded but not implemented yet.",
	})
}

func (s *Server) handleCreateRevealSession(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusNotImplemented, errorResponse{
		Error:   "not_implemented",
		Message: "Reveal session endpoint is scaffolded but not implemented yet.",
	})
}

func (s *Server) handleSecretRoutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/secrets/")

	switch {
	case strings.HasSuffix(path, "/status") && r.Method == http.MethodGet:
		writeJSON(w, http.StatusNotImplemented, errorResponse{
			Error:   "not_implemented",
			Message: "Secret status endpoint is scaffolded but not implemented yet.",
		})
	case strings.HasSuffix(path, "/consume") && r.Method == http.MethodPost:
		writeJSON(w, http.StatusNotImplemented, errorResponse{
			Error:   "not_implemented",
			Message: "Secret consume endpoint is scaffolded but not implemented yet.",
		})
	default:
		writeJSON(w, http.StatusNotFound, errorResponse{
			Error:   "not_found",
			Message: "Route not found.",
		})
	}
}

func writeJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(payload)
}
