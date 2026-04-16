package httpapi

import (
	"encoding/json"
	"io"
	"net/http"
	"one-time-link/backend/internal/secret"
	"strings"
	"time"
)

type healthResponse struct {
	Service   string `json:"service"`
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Version   string `json:"version"`
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		RespondError(w, r, ErrMethodNotAllowed("Method not allowed"))
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
		RespondError(w, r, ErrMethodNotAllowed("Method not allowed"))
		return
	}

	// Read body to trigger MaxBytesReader if size exceeded
	body, err := io.ReadAll(r.Body)
	if err != nil {
		// MaxBytesReader will return error if limit exceeded
		RespondError(w, r, ErrPayloadTooLarge("Request body exceeds 15KB limit").
			WithDetail("max_size_bytes", 15*1024))
		return
	}

	// Parse JSON
	var req secret.CreateSecretRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		RespondError(w, r, ErrInvalidRequest("Invalid JSON request body").
			WithError(err))
		return
	}

	// Validate request
	if err := secret.ValidateCreateSecretRequest(req); err != nil {
		// Check if it's a multi-validation error
		if multiErr, ok := err.(*secret.MultiValidationError); ok {
			appErr := ErrValidationFailed("Request validation failed")

			// Add each validation error to details
			for _, valErr := range multiErr.Errors {
				AddValidationError(appErr, valErr.Field, valErr.Message, valErr.Code)
			}

			RespondError(w, r, appErr)
			return
		}

		// Fallback for other validation errors
		RespondError(w, r, ErrInvalidRequest(err.Error()).WithError(err))
		return
	}

	// Create secret using service
	resp, err := s.secretService.CreateSecret(r.Context(), req)
	if err != nil {
		RespondError(w, r, ErrInternal("Failed to create secret").WithError(err))
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

func (s *Server) handleCreateRevealSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		RespondError(w, r, ErrMethodNotAllowed("Method not allowed"))
		return
	}

	RespondError(w, r, ErrNotImplemented("Reveal session endpoint is scaffolded but not implemented yet"))
}

func (s *Server) handleSecretRoutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/secrets/")

	switch {
	case strings.HasSuffix(path, "/status") && r.Method == http.MethodGet:
		s.handleGetSecretStatus(w, r, path)
	case strings.HasSuffix(path, "/consume") && r.Method == http.MethodPost:
		s.handleConsumeSecret(w, r, path)
	default:
		RespondError(w, r, ErrNotFound("Route not found"))
	}
}

func (s *Server) handleGetSecretStatus(w http.ResponseWriter, r *http.Request, path string) {
	// Extract secret ID from path (remove "/status" suffix)
	secretID := strings.TrimSuffix(path, "/status")
	if secretID == "" {
		RespondError(w, r, ErrInvalidSecretID("Secret ID is required"))
		return
	}

	// Get secret status from service
	status, err := s.secretService.GetSecretStatus(r.Context(), secretID)
	if err != nil {
		RespondError(w, r, ErrInternal("Failed to get secret status").WithError(err))
		return
	}

	// If secret not found, return 404
	if status.Status == "not_found" {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(status)
		return
	}

	// Return status
	writeJSON(w, http.StatusOK, status)
}

func (s *Server) handleConsumeSecret(w http.ResponseWriter, r *http.Request, path string) {
	// Extract secret ID from path (remove "/consume" suffix)
	secretID := strings.TrimSuffix(path, "/consume")
	if secretID == "" {
		RespondError(w, r, ErrInvalidSecretID("Secret ID is required"))
		return
	}

	// Consume secret from service
	response, err := s.secretService.ConsumeSecret(r.Context(), secretID)
	if err != nil {
		// Check if secret not found or already consumed
		if strings.Contains(err.Error(), "not found or already consumed") {
			RespondError(w, r, ErrAlreadyConsumed("This secret has already been revealed or has expired.").
				WithDetail("secret_id", secretID))
			return
		}
		RespondError(w, r, ErrInternal("Failed to consume secret").WithError(err))
		return
	}

	// Return consumed secret data
	writeJSON(w, http.StatusOK, response)
}

func writeJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(payload)
}
