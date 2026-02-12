package service

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"

	"ai-companion-be/internal/models"
	"ai-companion-be/internal/repository"
)

const (
	defaultMoodScore         = 50.0
	defaultRelationshipScore = 0.0
	moodDecayPerHour         = 0.5 // mood points lost per hour of inactivity
)

// RelationshipService handles relationship state business logic.
type RelationshipService struct {
	relationships repository.RelationshipRepository
}

// NewRelationshipService creates a new RelationshipService.
func NewRelationshipService(relationships repository.RelationshipRepository) *RelationshipService {
	return &RelationshipService{relationships: relationships}
}

// SelectCompanion creates the initial relationship state during onboarding.
func (s *RelationshipService) SelectCompanion(ctx context.Context, userID, companionID uuid.UUID) (*models.RelationshipState, error) {
	state := &models.RelationshipState{
		ID:                uuid.New(),
		UserID:            userID,
		CompanionID:       companionID,
		MoodScore:         defaultMoodScore,
		RelationshipScore: defaultRelationshipScore,
	}

	if err := s.relationships.Create(ctx, state); err != nil {
		return nil, fmt.Errorf("creating relationship: %w", err)
	}

	state.MoodLabel = models.GetMoodLabel(state.MoodScore)
	return state, nil
}

// GetRelationship returns the relationship state with time-decayed mood.
func (s *RelationshipService) GetRelationship(ctx context.Context, userID, companionID uuid.UUID) (*models.RelationshipState, error) {
	state, err := s.relationships.GetByUserAndCompanion(ctx, userID, companionID)
	if err != nil {
		return nil, err
	}

	applyTimeDecay(state)
	state.MoodLabel = models.GetMoodLabel(state.MoodScore)
	return state, nil
}

// GetAllRelationships returns all relationship states for a user with time decay applied.
func (s *RelationshipService) GetAllRelationships(ctx context.Context, userID uuid.UUID) ([]models.RelationshipState, error) {
	states, err := s.relationships.GetAllByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	for i := range states {
		applyTimeDecay(&states[i])
		states[i].MoodLabel = models.GetMoodLabel(states[i].MoodScore)
	}

	return states, nil
}

// applyTimeDecay reduces mood_score based on hours since last interaction.
// Mood is recalculated on read, not persisted.
func applyTimeDecay(state *models.RelationshipState) {
	hoursSinceInteraction := time.Since(state.LastInteraction).Hours()
	decay := hoursSinceInteraction * moodDecayPerHour
	state.MoodScore = math.Max(0, state.MoodScore-decay)
}
