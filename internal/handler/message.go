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

// MessageHandler handles chat message endpoints.
type MessageHandler struct {
	messages *service.MessageService
}

// NewMessageHandler creates a new MessageHandler.
func NewMessageHandler(messages *service.MessageService) *MessageHandler {
	return &MessageHandler{messages: messages}
}

// Send handles POST /api/companions/{id}/messages.
func (h *MessageHandler) Send(w http.ResponseWriter, r *http.Request) {
	companionID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid companion id")
		return
	}

	var req models.SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	userID := middleware.GetUserID(r.Context())

	messages, err := h.messages.SendMessage(r.Context(), userID, companionID, req)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	JSON(w, http.StatusCreated, messages)
}

// GetHistory handles GET /api/companions/{id}/messages?cursor=...&limit=...
func (h *MessageHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	companionID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid companion id")
		return
	}

	userID := middleware.GetUserID(r.Context())

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

	page, err := h.messages.GetMessages(r.Context(), userID, companionID, cursor, limit)
	if err != nil {
		Error(w, http.StatusInternalServerError, "failed to fetch messages")
		return
	}

	JSON(w, http.StatusOK, page)
}
