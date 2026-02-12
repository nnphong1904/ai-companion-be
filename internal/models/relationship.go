package models

import (
	"time"

	"github.com/google/uuid"
)

// RelationshipState tracks the emotional state between a user and a companion.
type RelationshipState struct {
	ID                uuid.UUID `json:"id"`
	UserID            uuid.UUID `json:"user_id"`
	CompanionID       uuid.UUID `json:"companion_id"`
	MoodScore         float64   `json:"mood_score"`
	RelationshipScore float64   `json:"relationship_score"`
	MoodLabel         string    `json:"mood_label"`
	LastInteraction   time.Time `json:"last_interaction"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// GetMoodLabel returns a human-readable mood label for the given score.
//
//	<20  → Distant
//	20–50 → Neutral
//	50–80 → Happy
//	80+  → Attached
func GetMoodLabel(score float64) string {
	switch {
	case score < 20:
		return "Distant"
	case score < 50:
		return "Neutral"
	case score < 80:
		return "Happy"
	default:
		return "Attached"
	}
}

// SelectCompanionRequest is the payload for onboarding companion selection.
type SelectCompanionRequest struct {
	CompanionID uuid.UUID `json:"companion_id"`
}
