package models

import (
	"time"

	"github.com/google/uuid"
)

// Story represents a companion's story containing multiple media slides.
type Story struct {
	ID          uuid.UUID    `json:"id"`
	CompanionID uuid.UUID    `json:"companion_id"`
	CreatedAt   time.Time    `json:"created_at"`
	ExpiresAt   time.Time    `json:"expires_at"`
	Media       []StoryMedia `json:"media,omitempty"`
}

// StoryMedia represents a single slide within a story.
type StoryMedia struct {
	ID        uuid.UUID `json:"id"`
	StoryID   uuid.UUID `json:"story_id"`
	MediaURL  string    `json:"media_url"`
	MediaType string    `json:"media_type"` // "image" or "video"
	Duration  int       `json:"duration"`   // display duration in seconds
	SortOrder int       `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
}

// StoryReaction represents a user's emoji reaction to a story slide.
type StoryReaction struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	StoryID   uuid.UUID `json:"story_id"`
	MediaID   uuid.UUID `json:"media_id"`
	Reaction  string    `json:"reaction"` // "love", "sad", "heart_eyes", "angry"
	CreatedAt time.Time `json:"created_at"`
}

// StoryPage represents a cursor-paginated page of stories.
type StoryPage struct {
	Stories    []Story `json:"stories"`
	NextCursor string  `json:"next_cursor,omitempty"`
	HasMore    bool    `json:"has_more"`
}

// ReactToStoryRequest is the payload for reacting to a story slide.
type ReactToStoryRequest struct {
	MediaID  uuid.UUID `json:"media_id"`
	Reaction string    `json:"reaction"`
}
