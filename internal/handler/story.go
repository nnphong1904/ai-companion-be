package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

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

// GetActiveStories handles GET /api/stories?cursor=...&limit=...
func (h *StoryHandler) GetActiveStories(w http.ResponseWriter, r *http.Request) {
	var cursor *time.Time
	if cursorStr := r.URL.Query().Get("cursor"); cursorStr != "" {
		t, err := time.Parse(time.RFC3339Nano, cursorStr)
		if err != nil {
			Error(w, http.StatusBadRequest, "invalid cursor format")
			return
		}
		cursor = &t
	}

	limit := 20
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	page, err := h.stories.GetActiveStories(r.Context(), cursor, limit)
	if err != nil {
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
