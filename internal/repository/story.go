package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"ai-companion-be/internal/models"
)

// StoryRepository defines data access operations for stories.
type StoryRepository interface {
	GetByCompanionID(ctx context.Context, companionID uuid.UUID) ([]models.Story, error)
	GetActiveStories(ctx context.Context) ([]models.Story, error)
	CreateReaction(ctx context.Context, reaction *models.StoryReaction) error
}

type storyRepo struct {
	pool *pgxpool.Pool
}

// NewStoryRepository creates a new StoryRepository backed by PostgreSQL.
func NewStoryRepository(pool *pgxpool.Pool) StoryRepository {
	return &storyRepo{pool: pool}
}

func (r *storyRepo) GetByCompanionID(ctx context.Context, companionID uuid.UUID) ([]models.Story, error) {
	query := `
		SELECT s.id, s.companion_id, s.created_at, s.expires_at
		FROM stories s
		WHERE s.companion_id = $1 AND s.expires_at > NOW()
		ORDER BY s.created_at DESC`

	rows, err := r.pool.Query(ctx, query, companionID)
	if err != nil {
		return nil, fmt.Errorf("querying stories: %w", err)
	}
	defer rows.Close()

	var stories []models.Story
	for rows.Next() {
		var s models.Story
		if err := rows.Scan(&s.ID, &s.CompanionID, &s.CreatedAt, &s.ExpiresAt); err != nil {
			return nil, fmt.Errorf("scanning story: %w", err)
		}
		stories = append(stories, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Load media for each story to avoid N+1: batch query.
	if len(stories) > 0 {
		storyIDs := make([]uuid.UUID, len(stories))
		storyMap := make(map[uuid.UUID]int, len(stories))
		for i, s := range stories {
			storyIDs[i] = s.ID
			storyMap[s.ID] = i
		}

		mediaQuery := `
			SELECT id, story_id, media_url, media_type, duration, sort_order, created_at
			FROM story_media
			WHERE story_id = ANY($1)
			ORDER BY sort_order`

		mediaRows, err := r.pool.Query(ctx, mediaQuery, storyIDs)
		if err != nil {
			return nil, fmt.Errorf("querying story media: %w", err)
		}
		defer mediaRows.Close()

		for mediaRows.Next() {
			var m models.StoryMedia
			if err := mediaRows.Scan(&m.ID, &m.StoryID, &m.MediaURL, &m.MediaType, &m.Duration, &m.SortOrder, &m.CreatedAt); err != nil {
				return nil, fmt.Errorf("scanning story media: %w", err)
			}
			if idx, ok := storyMap[m.StoryID]; ok {
				stories[idx].Media = append(stories[idx].Media, m)
			}
		}
		if err := mediaRows.Err(); err != nil {
			return nil, err
		}
	}

	return stories, nil
}

func (r *storyRepo) GetActiveStories(ctx context.Context) ([]models.Story, error) {
	query := `
		SELECT s.id, s.companion_id, s.created_at, s.expires_at
		FROM stories s
		WHERE s.expires_at > NOW()
		ORDER BY s.created_at DESC`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("querying active stories: %w", err)
	}
	defer rows.Close()

	var stories []models.Story
	for rows.Next() {
		var s models.Story
		if err := rows.Scan(&s.ID, &s.CompanionID, &s.CreatedAt, &s.ExpiresAt); err != nil {
			return nil, fmt.Errorf("scanning story: %w", err)
		}
		stories = append(stories, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Batch-load media.
	if len(stories) > 0 {
		storyIDs := make([]uuid.UUID, len(stories))
		storyMap := make(map[uuid.UUID]int, len(stories))
		for i, s := range stories {
			storyIDs[i] = s.ID
			storyMap[s.ID] = i
		}

		mediaQuery := `
			SELECT id, story_id, media_url, media_type, duration, sort_order, created_at
			FROM story_media
			WHERE story_id = ANY($1)
			ORDER BY sort_order`

		mediaRows, err := r.pool.Query(ctx, mediaQuery, storyIDs)
		if err != nil {
			return nil, fmt.Errorf("querying story media: %w", err)
		}
		defer mediaRows.Close()

		for mediaRows.Next() {
			var m models.StoryMedia
			if err := mediaRows.Scan(&m.ID, &m.StoryID, &m.MediaURL, &m.MediaType, &m.Duration, &m.SortOrder, &m.CreatedAt); err != nil {
				return nil, fmt.Errorf("scanning story media: %w", err)
			}
			if idx, ok := storyMap[m.StoryID]; ok {
				stories[idx].Media = append(stories[idx].Media, m)
			}
		}
		if err := mediaRows.Err(); err != nil {
			return nil, err
		}
	}

	return stories, nil
}

func (r *storyRepo) CreateReaction(ctx context.Context, reaction *models.StoryReaction) error {
	query := `
		INSERT INTO story_reactions (id, user_id, story_id, media_id, reaction, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		ON CONFLICT (user_id, media_id) DO UPDATE SET reaction = $5
		RETURNING created_at`

	return r.pool.QueryRow(ctx, query,
		reaction.ID, reaction.UserID, reaction.StoryID, reaction.MediaID, reaction.Reaction,
	).Scan(&reaction.CreatedAt)
}
