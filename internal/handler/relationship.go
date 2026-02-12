package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"ai-companion-be/internal/middleware"
	"ai-companion-be/internal/models"
	"ai-companion-be/internal/service"
)

// RelationshipHandler handles relationship state endpoints.
type RelationshipHandler struct {
	relationships *service.RelationshipService
}

// NewRelationshipHandler creates a new RelationshipHandler.
func NewRelationshipHandler(relationships *service.RelationshipService) *RelationshipHandler {
	return &RelationshipHandler{relationships: relationships}
}

// SelectCompanion handles POST /api/onboarding/select-companion.
func (h *RelationshipHandler) SelectCompanion(w http.ResponseWriter, r *http.Request) {
	var req models.SelectCompanionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	userID := middleware.GetUserID(r.Context())

	state, err := h.relationships.SelectCompanion(r.Context(), userID, req.CompanionID)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	JSON(w, http.StatusCreated, state)
}

// GetRelationship handles GET /api/companions/{id}/relationship.
func (h *RelationshipHandler) GetRelationship(w http.ResponseWriter, r *http.Request) {
	companionID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid companion id")
		return
	}

	userID := middleware.GetUserID(r.Context())

	state, err := h.relationships.GetRelationship(r.Context(), userID, companionID)
	if err != nil {
		Error(w, http.StatusNotFound, "relationship not found")
		return
	}

	JSON(w, http.StatusOK, state)
}

// GetAllRelationships handles GET /api/relationships.
func (h *RelationshipHandler) GetAllRelationships(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	states, err := h.relationships.GetAllRelationships(r.Context(), userID)
	if err != nil {
		Error(w, http.StatusInternalServerError, "failed to fetch relationships")
		return
	}

	JSON(w, http.StatusOK, states)
}
