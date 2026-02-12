package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"ai-companion-be/internal/middleware"
	"ai-companion-be/internal/models"
	"ai-companion-be/internal/service"
)

// MemoryHandler handles memory endpoints.
type MemoryHandler struct {
	memories *service.MemoryService
}

// NewMemoryHandler creates a new MemoryHandler.
func NewMemoryHandler(memories *service.MemoryService) *MemoryHandler {
	return &MemoryHandler{memories: memories}
}

// Create handles POST /api/companions/{id}/memories.
func (h *MemoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	companionID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid companion id")
		return
	}

	var req models.CreateMemoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	userID := middleware.GetUserID(r.Context())

	memory, err := h.memories.Create(r.Context(), userID, companionID, req)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	JSON(w, http.StatusCreated, memory)
}

// GetByCompanion handles GET /api/companions/{id}/memories?limit=...
func (h *MemoryHandler) GetByCompanion(w http.ResponseWriter, r *http.Request) {
	companionID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid companion id")
		return
	}

	userID := middleware.GetUserID(r.Context())

	limit := 50
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	page, err := h.memories.GetByCompanion(r.Context(), userID, companionID, limit)
	if err != nil {
		Error(w, http.StatusInternalServerError, "failed to fetch memories")
		return
	}

	JSON(w, http.StatusOK, page)
}

// Delete handles DELETE /api/memories/{id}.
func (h *MemoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	memoryID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid memory id")
		return
	}

	userID := middleware.GetUserID(r.Context())

	if err := h.memories.Delete(r.Context(), userID, memoryID); err != nil {
		Error(w, http.StatusNotFound, err.Error())
		return
	}

	JSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

// TogglePin handles PATCH /api/memories/{id}/pin.
func (h *MemoryHandler) TogglePin(w http.ResponseWriter, r *http.Request) {
	memoryID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid memory id")
		return
	}

	userID := middleware.GetUserID(r.Context())

	memory, err := h.memories.TogglePin(r.Context(), userID, memoryID)
	if err != nil {
		Error(w, http.StatusNotFound, err.Error())
		return
	}

	JSON(w, http.StatusOK, memory)
}
