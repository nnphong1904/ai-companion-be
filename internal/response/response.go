package response

import (
	"encoding/json"
	"net/http"
)

// APIResponse is the standard JSON envelope for all API responses.
type APIResponse struct {
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

// JSON writes a success JSON response.
func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(APIResponse{Data: data})
}

// Error writes an error JSON response.
func Error(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(APIResponse{Error: message})
}
