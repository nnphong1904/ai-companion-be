package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"ai-companion-be/internal/service"
)

// CompanionHandler handles companion endpoints.
type CompanionHandler struct {
	companions *service.CompanionService
}

// NewCompanionHandler creates a new CompanionHandler.
func NewCompanionHandler(companions *service.CompanionService) *CompanionHandler {
	return &CompanionHandler{companions: companions}
}

// GetAll handles GET /api/companions.
func (h *CompanionHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	companions, err := h.companions.GetAll(r.Context())
	if err != nil {
		Error(w, http.StatusInternalServerError, "failed to fetch companions")
		return
	}

	JSON(w, http.StatusOK, companions)
}

// GetByID handles GET /api/companions/{id}.
func (h *CompanionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid companion id")
		return
	}

	companion, err := h.companions.GetByID(r.Context(), id)
	if err != nil {
		Error(w, http.StatusNotFound, "companion not found")
		return
	}

	JSON(w, http.StatusOK, companion)
}
