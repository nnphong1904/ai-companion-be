package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"ai-companion-be/internal/models"
)

// InsightsRepository defines data access operations for relationship insights.
type InsightsRepository interface {
	RecordMoodSnapshot(ctx context.Context, userID, companionID uuid.UUID, moodScore float64) error
	GetMoodHistory(ctx context.Context, userID, companionID uuid.UUID, days int) ([]models.MoodSnapshot, error)
	GetMessageDates(ctx context.Context, userID, companionID uuid.UUID) ([]time.Time, error)
	GetStats(ctx context.Context, userID, companionID uuid.UUID) (*models.InsightStats, error)
	GetReactionSummary(ctx context.Context, userID, companionID uuid.UUID) (*models.ReactionSummary, error)
}

type insightsRepo struct {
	pool *pgxpool.Pool
}

// NewInsightsRepository creates a new InsightsRepository backed by PostgreSQL.
func NewInsightsRepository(pool *pgxpool.Pool) InsightsRepository {
	return &insightsRepo{pool: pool}
}

func (r *insightsRepo) RecordMoodSnapshot(ctx context.Context, userID, companionID uuid.UUID, moodScore float64) error {
	query := `
		INSERT INTO mood_history (id, user_id, companion_id, recorded_date, mood_score)
		VALUES (gen_random_uuid(), $1, $2, CURRENT_DATE, $3)
		ON CONFLICT (user_id, companion_id, recorded_date)
		DO UPDATE SET mood_score = $3`

	_, err := r.pool.Exec(ctx, query, userID, companionID, moodScore)
	if err != nil {
		return fmt.Errorf("recording mood snapshot: %w", err)
	}
	return nil
}

func (r *insightsRepo) GetMoodHistory(ctx context.Context, userID, companionID uuid.UUID, days int) ([]models.MoodSnapshot, error) {
	query := `
		SELECT recorded_date, mood_score
		FROM mood_history
		WHERE user_id = $1 AND companion_id = $2
		  AND recorded_date >= CURRENT_DATE - $3::int
		ORDER BY recorded_date ASC`

	rows, err := r.pool.Query(ctx, query, userID, companionID, days)
	if err != nil {
		return nil, fmt.Errorf("querying mood history: %w", err)
	}
	defer rows.Close()

	var history []models.MoodSnapshot
	for rows.Next() {
		var date time.Time
		var score float64
		if err := rows.Scan(&date, &score); err != nil {
			return nil, fmt.Errorf("scanning mood history: %w", err)
		}
		history = append(history, models.MoodSnapshot{
			Date:      date.Format("2006-01-02"),
			MoodScore: score,
			MoodLabel: models.GetMoodLabel(score),
		})
	}

	return history, rows.Err()
}

func (r *insightsRepo) GetMessageDates(ctx context.Context, userID, companionID uuid.UUID) ([]time.Time, error) {
	query := `
		SELECT DISTINCT created_at::date AS msg_date
		FROM messages
		WHERE user_id = $1 AND companion_id = $2 AND role = 'user'
		ORDER BY msg_date DESC`

	rows, err := r.pool.Query(ctx, query, userID, companionID)
	if err != nil {
		return nil, fmt.Errorf("querying message dates: %w", err)
	}
	defer rows.Close()

	var dates []time.Time
	for rows.Next() {
		var d time.Time
		if err := rows.Scan(&d); err != nil {
			return nil, fmt.Errorf("scanning message date: %w", err)
		}
		dates = append(dates, d)
	}

	return dates, rows.Err()
}

func (r *insightsRepo) GetStats(ctx context.Context, userID, companionID uuid.UUID) (*models.InsightStats, error) {
	query := `
		SELECT
			(SELECT count(*) FROM messages WHERE user_id = $1 AND companion_id = $2) AS total_messages,
			(SELECT count(*) FROM memories WHERE user_id = $1 AND companion_id = $2) AS total_memories,
			(SELECT min(created_at) FROM messages WHERE user_id = $1 AND companion_id = $2) AS first_message`

	var stats models.InsightStats
	var firstMsg *time.Time

	err := r.pool.QueryRow(ctx, query, userID, companionID).
		Scan(&stats.TotalMessages, &stats.TotalMemories, &firstMsg)
	if err != nil {
		return nil, fmt.Errorf("querying insight stats: %w", err)
	}

	stats.FirstMessage = firstMsg
	if firstMsg != nil {
		stats.DaysTogether = int(time.Since(*firstMsg).Hours()/24) + 1
	}

	return &stats, nil
}

func (r *insightsRepo) GetReactionSummary(ctx context.Context, userID, companionID uuid.UUID) (*models.ReactionSummary, error) {
	// Counts by reaction type.
	countsQuery := `
		SELECT sr.reaction, COUNT(*) AS count
		FROM story_reactions sr
		JOIN story_media sm ON sr.media_id = sm.id
		JOIN stories s ON sm.story_id = s.id
		WHERE sr.user_id = $1 AND s.companion_id = $2
		GROUP BY sr.reaction`

	rows, err := r.pool.Query(ctx, countsQuery, userID, companionID)
	if err != nil {
		return nil, fmt.Errorf("querying reaction counts: %w", err)
	}
	defer rows.Close()

	counts := map[string]int{
		"love":       0,
		"sad":        0,
		"heart_eyes": 0,
		"angry":      0,
	}
	total := 0

	for rows.Next() {
		var reaction string
		var count int
		if err := rows.Scan(&reaction, &count); err != nil {
			return nil, fmt.Errorf("scanning reaction count: %w", err)
		}
		counts[reaction] = count
		total += count
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating reaction counts: %w", err)
	}

	// Dominant emotion.
	var dominant *string
	maxCount := 0
	for reaction, count := range counts {
		if count > maxCount {
			maxCount = count
			r := reaction
			dominant = &r
		}
	}

	// Recent reactions.
	recentQuery := `
		SELECT sr.reaction, sr.created_at
		FROM story_reactions sr
		JOIN story_media sm ON sr.media_id = sm.id
		JOIN stories s ON sm.story_id = s.id
		WHERE sr.user_id = $1 AND s.companion_id = $2
		ORDER BY sr.created_at DESC
		LIMIT 5`

	recentRows, err := r.pool.Query(ctx, recentQuery, userID, companionID)
	if err != nil {
		return nil, fmt.Errorf("querying recent reactions: %w", err)
	}
	defer recentRows.Close()

	var recent []models.RecentReaction
	for recentRows.Next() {
		var rr models.RecentReaction
		if err := recentRows.Scan(&rr.Reaction, &rr.ReactedAt); err != nil {
			return nil, fmt.Errorf("scanning recent reaction: %w", err)
		}
		recent = append(recent, rr)
	}
	if err := recentRows.Err(); err != nil {
		return nil, fmt.Errorf("iterating recent reactions: %w", err)
	}

	return &models.ReactionSummary{
		Total:           total,
		Counts:          counts,
		Recent:          recent,
		DominantEmotion: dominant,
	}, nil
}
