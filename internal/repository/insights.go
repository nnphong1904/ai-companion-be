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
