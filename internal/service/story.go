package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"ai-companion-be/internal/models"
	"ai-companion-be/internal/repository"
)

// StoryService handles story-related business logic.
type StoryService struct {
	stories       repository.StoryRepository
	relationships repository.RelationshipRepository
}

// NewStoryService creates a new StoryService.
func NewStoryService(stories repository.StoryRepository, relationships repository.RelationshipRepository) *StoryService {
	return &StoryService{stories: stories, relationships: relationships}
}

// GetByCompanionID returns all active stories for a companion.
func (s *StoryService) GetByCompanionID(ctx context.Context, companionID uuid.UUID) ([]models.Story, error) {
	return s.stories.GetByCompanionID(ctx, companionID)
}

// GetActiveStories returns all currently active stories across companions.
func (s *StoryService) GetActiveStories(ctx context.Context) ([]models.Story, error) {
	return s.stories.GetActiveStories(ctx)
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
	// Best-effort: don't fail the reaction if relationship update fails.
	// We need the companion ID from the story. For now, we skip if it fails.
	stories, err := s.stories.GetActiveStories(ctx)
	if err != nil {
		return
	}

	var companionID uuid.UUID
	for _, st := range stories {
		if st.ID == storyID {
			companionID = st.CompanionID
			break
		}
	}
	if companionID == uuid.Nil {
		return
	}

	state, err := s.relationships.GetByUserAndCompanion(ctx, userID, companionID)
	if err != nil {
		return
	}

	state.MoodScore = clampScore(state.MoodScore + 3)
	state.RelationshipScore = clampScore(state.RelationshipScore + 2)
	_ = s.relationships.Update(ctx, state)
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
