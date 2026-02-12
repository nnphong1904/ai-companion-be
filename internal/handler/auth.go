package handler

import (
	"encoding/json"
	"net/http"

	"ai-companion-be/internal/middleware"
	"ai-companion-be/internal/models"
	"ai-companion-be/internal/service"
)

// AuthHandler handles authentication endpoints.
type AuthHandler struct {
	auth *service.AuthService
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(auth *service.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

// Signup handles POST /api/auth/signup.
func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var req models.SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.auth.Signup(r.Context(), req)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	JSON(w, http.StatusCreated, resp)
}

// Login handles POST /api/auth/login.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.auth.Login(r.Context(), req)
	if err != nil {
		Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	JSON(w, http.StatusOK, resp)
}

// Me handles GET /api/auth/me.
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	user, err := h.auth.GetUser(r.Context(), userID)
	if err != nil {
		Error(w, http.StatusNotFound, "user not found")
		return
	}

	JSON(w, http.StatusOK, user)
}
