package models

import (
	"time"

	"github.com/google/uuid"
)

// Memory represents a curated meaningful moment between a user and a companion.
type Memory struct {
	ID          uuid.UUID  `json:"id"`
	UserID      uuid.UUID  `json:"user_id"`
	CompanionID uuid.UUID  `json:"companion_id"`
	MessageID   *uuid.UUID `json:"message_id,omitempty"`
	Content     string     `json:"content"`
	Tag         *string    `json:"tag,omitempty"`
	Pinned      bool       `json:"pinned"`
	CreatedAt   time.Time  `json:"created_at"`
}

// MemoryPage represents a paginated page of memories.
type MemoryPage struct {
	Memories   []Memory `json:"memories"`
	NextCursor string   `json:"next_cursor,omitempty"`
	HasMore    bool     `json:"has_more"`
}

// CreateMemoryRequest is the payload for creating a new memory.
type CreateMemoryRequest struct {
	MessageID *uuid.UUID `json:"message_id,omitempty"`
	Content   string     `json:"content"`
	Tag       *string    `json:"tag,omitempty"`
}
