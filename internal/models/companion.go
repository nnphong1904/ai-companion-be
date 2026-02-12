package models

import (
	"time"

	"github.com/google/uuid"
)

// Companion represents an AI companion.
type Companion struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	AvatarURL   string    `json:"avatar_url"`
	Personality string    `json:"personality"`
	CreatedAt   time.Time `json:"created_at"`
}
