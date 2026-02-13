package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"ai-companion-be/internal/models"
	"ai-companion-be/internal/repository"
)

// StoryService handles story-related business logic.
type StoryService struct {
	stories       repository.StoryRepository
	relationships repository.RelationshipRepository
	insights      repository.InsightsRepository
}

// NewStoryService creates a new StoryService.
func NewStoryService(stories repository.StoryRepository, relationships repository.RelationshipRepository, insights repository.InsightsRepository) *StoryService {
	return &StoryService{stories: stories, relationships: relationships, insights: insights}
}

// GetByCompanionID returns all active stories for a companion.
func (s *StoryService) GetByCompanionID(ctx context.Context, companionID uuid.UUID) ([]models.Story, error) {
	return s.stories.GetByCompanionID(ctx, companionID)
}

// GetActiveStories returns a paginated list of currently active stories.
func (s *StoryService) GetActiveStories(ctx context.Context, cursor *time.Time, limit int) (*models.StoryPage, error) {
	return s.stories.GetActiveStories(ctx, cursor, limit)
}

// ReactToStory records a user's reaction and updates the relationship state.
func (s *StoryService) ReactToStory(ctx context.Context, userID uuid.UUID, storyID uuid.UUID, req models.ReactToStoryRequest) error {
	validReactions := map[string]bool{
		"love": true, "sad": true, "heart_eyes": true, "angry": true,
	}
	if !validReactions[req.Reaction] {
		return fmt.Errorf("invalid reaction: must be love, sad, heart_eyes, or angry")
	}

	reaction := &models.StoryReaction{
		ID:       uuid.New(),
		UserID:   userID,
		StoryID:  storyID,
		MediaID:  req.MediaID,
		Reaction: req.Reaction,
	}

	if err := s.stories.CreateReaction(ctx, reaction); err != nil {
		return fmt.Errorf("creating reaction: %w", err)
	}

	// Update relationship state: story reactions boost mood and relationship.
	s.updateRelationshipOnReaction(ctx, userID, storyID)

	return nil
}

func (s *StoryService) updateRelationshipOnReaction(ctx context.Context, userID, storyID uuid.UUID) {
	// Look up the single story by ID instead of fetching all active stories.
	story, err := s.stories.GetByID(ctx, storyID)
	if err != nil {
		return
	}

	state, err := s.relationships.GetByUserAndCompanion(ctx, userID, story.CompanionID)
	if err != nil {
		return
	}

	state.MoodScore = clampScore(state.MoodScore + 3)
	state.RelationshipScore = clampScore(state.RelationshipScore + 2)
	_ = s.relationships.Update(ctx, state)

	// Record daily mood snapshot for insights.
	_ = s.insights.RecordMoodSnapshot(ctx, userID, story.CompanionID, state.MoodScore)
}

func clampScore(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 100 {
		return 100
	}
	return v
}
