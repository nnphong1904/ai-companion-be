package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"ai-companion-be/internal/middleware"
	"ai-companion-be/internal/service"
)

// InsightsHandler handles relationship insights endpoints.
type InsightsHandler struct {
	insights *service.InsightsService
}

// NewInsightsHandler creates a new InsightsHandler.
func NewInsightsHandler(insights *service.InsightsService) *InsightsHandler {
	return &InsightsHandler{insights: insights}
}

// GetInsights handles GET /api/companions/{id}/insights.
func (h *InsightsHandler) GetInsights(w http.ResponseWriter, r *http.Request) {
	companionID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid companion id")
		return
	}

	userID := middleware.GetUserID(r.Context())

	insights, err := h.insights.GetInsights(r.Context(), userID, companionID)
	if err != nil {
		Error(w, http.StatusInternalServerError, "failed to fetch insights")
		return
	}

	JSON(w, http.StatusOK, insights)
}

// GetReactionSummary handles GET /api/companions/{id}/reactions/summary.
func (h *InsightsHandler) GetReactionSummary(w http.ResponseWriter, r *http.Request) {
	companionID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid companion id")
		return
	}

	userID := middleware.GetUserID(r.Context())

	summary, err := h.insights.GetReactionSummary(r.Context(), userID, companionID)
	if err != nil {
		Error(w, http.StatusInternalServerError, "failed to fetch reaction summary")
		return
	}

	JSON(w, http.StatusOK, summary)
}
