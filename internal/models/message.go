package models

import (
	"time"

	"github.com/google/uuid"
)

// Message represents a chat message between a user and a companion.
type Message struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	CompanionID uuid.UUID `json:"companion_id"`
	Content     string    `json:"content"`
	Role        string    `json:"role"` // "user" or "companion"
	CreatedAt   time.Time `json:"created_at"`
	IsMemorized bool      `json:"is_memorized"`
}

// SendMessageRequest is the payload for sending a chat message.
type SendMessageRequest struct {
	Content string `json:"content"`
}

// MessagePage represents a cursor-paginated page of messages.
type MessagePage struct {
	Messages []Message `json:"messages"`
	NextCursor string  `json:"next_cursor,omitempty"`
	HasMore    bool    `json:"has_more"`
}
