package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"ai-companion-be/internal/middleware"
	"ai-companion-be/internal/models"
	"ai-companion-be/internal/service"
)

// StoryHandler handles story endpoints.
type StoryHandler struct {
	stories *service.StoryService
}

// NewStoryHandler creates a new StoryHandler.
func NewStoryHandler(stories *service.StoryService) *StoryHandler {
	return &StoryHandler{stories: stories}
}

// GetActiveStories handles GET /api/stories — returns stories grouped by companion.
func (h *StoryHandler) GetActiveStories(w http.ResponseWriter, r *http.Request) {
	page, err := h.stories.GetActiveStoriesGrouped(r.Context())
	if err != nil {
		slog.Error("GetActiveStories failed", "error", err)
		Error(w, http.StatusInternalServerError, "failed to fetch stories")
		return
	}

	JSON(w, http.StatusOK, page)
}

// GetByCompanion handles GET /api/companions/{id}/stories.
func (h *StoryHandler) GetByCompanion(w http.ResponseWriter, r *http.Request) {
	companionID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid companion id")
		return
	}

	stories, err := h.stories.GetByCompanionID(r.Context(), companionID)
	if err != nil {
		Error(w, http.StatusInternalServerError, "failed to fetch stories")
		return
	}

	JSON(w, http.StatusOK, stories)
}

// React handles POST /api/stories/{id}/react.
func (h *StoryHandler) React(w http.ResponseWriter, r *http.Request) {
	storyID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid story id")
		return
	}

	var req models.ReactToStoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	userID := middleware.GetUserID(r.Context())

	if err := h.stories.ReactToStory(r.Context(), userID, storyID, req); err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
