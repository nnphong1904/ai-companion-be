package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"ai-companion-be/internal/models"
	"ai-companion-be/internal/repository"
)

// InsightsService handles relationship insights business logic.
type InsightsService struct {
	insights      repository.InsightsRepository
	relationships repository.RelationshipRepository
}

// NewInsightsService creates a new InsightsService.
func NewInsightsService(insights repository.InsightsRepository, relationships repository.RelationshipRepository) *InsightsService {
	return &InsightsService{insights: insights, relationships: relationships}
}

// GetInsights computes the full insights payload for a user-companion pair.
func (s *InsightsService) GetInsights(ctx context.Context, userID, companionID uuid.UUID) (*models.CompanionInsights, error) {
	moodHistory, err := s.insights.GetMoodHistory(ctx, userID, companionID, 14)
	if err != nil {
		return nil, err
	}

	dates, err := s.insights.GetMessageDates(ctx, userID, companionID)
	if err != nil {
		return nil, err
	}
	streak := computeStreak(dates)

	stats, err := s.insights.GetStats(ctx, userID, companionID)
	if err != nil {
		return nil, err
	}

	state, _ := s.relationships.GetByUserAndCompanion(ctx, userID, companionID)

	milestones := computeMilestones(stats, state)

	return &models.CompanionInsights{
		MoodHistory: moodHistory,
		Streak:      streak,
		Milestones:  milestones,
		Stats:       *stats,
	}, nil
}

// RecordMood records a daily mood snapshot (called from other services on interaction).
func (s *InsightsService) RecordMood(ctx context.Context, userID, companionID uuid.UUID, moodScore float64) {
	_ = s.insights.RecordMoodSnapshot(ctx, userID, companionID, moodScore)
}

func computeStreak(dates []time.Time) models.StreakInfo {
	if len(dates) == 0 {
		return models.StreakInfo{}
	}

	today := time.Now().UTC().Truncate(24 * time.Hour)

	current := 0
	longest := 0
	streak := 1

	// dates are DESC order. Check if today or yesterday is included for current streak.
	firstDate := dates[0].UTC().Truncate(24 * time.Hour)
	daysSinceFirst := int(today.Sub(firstDate).Hours() / 24)

	if daysSinceFirst > 1 {
		// Last message was more than 1 day ago — no current streak.
		current = 0
	} else {
		current = 1
	}

	for i := 1; i < len(dates); i++ {
		prev := dates[i-1].UTC().Truncate(24 * time.Hour)
		curr := dates[i].UTC().Truncate(24 * time.Hour)
		diff := int(prev.Sub(curr).Hours() / 24)

		if diff == 1 {
			streak++
		} else {
			if streak > longest {
				longest = streak
			}
			streak = 1
		}
	}
	if streak > longest {
		longest = streak
	}

	// Current streak: count consecutive days backward from today/yesterday.
	if current > 0 {
		for i := 1; i < len(dates); i++ {
			prev := dates[i-1].UTC().Truncate(24 * time.Hour)
			curr := dates[i].UTC().Truncate(24 * time.Hour)
			if int(prev.Sub(curr).Hours()/24) == 1 {
				current++
			} else {
				break
			}
		}
	}

	return models.StreakInfo{
		Current: current,
		Longest: longest,
	}
}

func computeMilestones(stats *models.InsightStats, state *models.RelationshipState) []models.Milestone {
	milestones := []models.Milestone{
		{
			Key:         "first_message",
			Title:       "First Words",
			Description: "Sent your first message",
			Achieved:    stats.TotalMessages > 0,
		},
		{
			Key:         "messages_50",
			Title:       "Getting Chatty",
			Description: "Exchanged 50 messages",
			Achieved:    stats.TotalMessages >= 50,
		},
		{
			Key:         "messages_200",
			Title:       "Deep Conversations",
			Description: "Exchanged 200 messages",
			Achieved:    stats.TotalMessages >= 200,
		},
		{
			Key:         "messages_1000",
			Title:       "Inseparable",
			Description: "Exchanged 1,000 messages",
			Achieved:    stats.TotalMessages >= 1000,
		},
		{
			Key:         "first_memory",
			Title:       "First Memory",
			Description: "Saved your first memory together",
			Achieved:    stats.TotalMemories > 0,
		},
		{
			Key:         "memories_10",
			Title:       "Memory Lane",
			Description: "Saved 10 memories together",
			Achieved:    stats.TotalMemories >= 10,
		},
		{
			Key:         "week_together",
			Title:       "One Week Together",
			Description: "Been connected for 7 days",
			Achieved:    stats.DaysTogether >= 7,
		},
		{
			Key:         "month_together",
			Title:       "One Month Together",
			Description: "Been connected for 30 days",
			Achieved:    stats.DaysTogether >= 30,
		},
	}

	if state != nil {
		milestones = append(milestones,
			models.Milestone{
				Key:         "mood_happy",
				Title:       "Warming Up",
				Description: "Reached Happy mood level",
				Achieved:    state.MoodScore >= 50,
			},
			models.Milestone{
				Key:         "mood_attached",
				Title:       "Deeply Attached",
				Description: "Reached Attached mood level",
				Achieved:    state.MoodScore >= 80,
			},
			models.Milestone{
				Key:         "bond_50",
				Title:       "Strong Bond",
				Description: "Relationship score reached 50",
				Achieved:    state.RelationshipScore >= 50,
			},
			models.Milestone{
				Key:         "bond_max",
				Title:       "Soulmates",
				Description: "Reached maximum relationship score",
				Achieved:    state.RelationshipScore >= 100,
			},
		)
	}

	return milestones
}
