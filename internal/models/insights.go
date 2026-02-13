package models

import "time"

// CompanionInsights aggregates relationship analytics for a user-companion pair.
type CompanionInsights struct {
	MoodHistory []MoodSnapshot `json:"mood_history"`
	Streak      StreakInfo     `json:"streak"`
	Milestones  []Milestone    `json:"milestones"`
	Stats       InsightStats   `json:"stats"`
}

// MoodSnapshot represents a single day's mood score.
type MoodSnapshot struct {
	Date      string  `json:"date"`
	MoodScore float64 `json:"mood_score"`
	MoodLabel string  `json:"mood_label"`
}

// StreakInfo tracks consecutive days of interaction.
type StreakInfo struct {
	Current int `json:"current"`
	Longest int `json:"longest"`
}

// Milestone represents an achieved relationship milestone.
type Milestone struct {
	Key         string     `json:"key"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	AchievedAt  *time.Time `json:"achieved_at,omitempty"`
	Achieved    bool       `json:"achieved"`
}

// InsightStats holds aggregate counts for the relationship.
type InsightStats struct {
	TotalMessages int        `json:"total_messages"`
	TotalMemories int        `json:"total_memories"`
	FirstMessage  *time.Time `json:"first_message,omitempty"`
	DaysTogether  int        `json:"days_together"`
}

// ReactionSummary aggregates story reaction data for a user-companion pair.
type ReactionSummary struct {
	Total           int              `json:"total"`
	Counts          map[string]int   `json:"counts"`
	Recent          []RecentReaction `json:"recent"`
	DominantEmotion *string          `json:"dominant_emotion"`
}

// RecentReaction represents a single recent reaction.
type RecentReaction struct {
	Reaction  string    `json:"reaction"`
	ReactedAt time.Time `json:"reacted_at"`
}
