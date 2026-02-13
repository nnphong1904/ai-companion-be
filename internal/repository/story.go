package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"ai-companion-be/internal/models"
)

// StoryRepository defines data access operations for stories.
type StoryRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*models.Story, error)
	GetByCompanionID(ctx context.Context, companionID uuid.UUID) ([]models.Story, error)
	GetActiveStories(ctx context.Context, cursor *time.Time, limit int) (*models.StoryPage, error)
	GetActiveStoriesGrouped(ctx context.Context, userID uuid.UUID) (*models.GroupedStoryPage, error)
	CreateReaction(ctx context.Context, reaction *models.StoryReaction) error
}

type storyRepo struct {
	pool *pgxpool.Pool
}

// NewStoryRepository creates a new StoryRepository backed by PostgreSQL.
func NewStoryRepository(pool *pgxpool.Pool) StoryRepository {
	return &storyRepo{pool: pool}
}

func (r *storyRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Story, error) {
	query := `
		SELECT id, companion_id, created_at, expires_at
		FROM stories
		WHERE id = $1`

	var s models.Story
	err := r.pool.QueryRow(ctx, query, id).
		Scan(&s.ID, &s.CompanionID, &s.CreatedAt, &s.ExpiresAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("story not found")
		}
		return nil, fmt.Errorf("getting story: %w", err)
	}
	return &s, nil
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

	return r.loadMedia(ctx, stories)
}

func (r *storyRepo) GetActiveStories(ctx context.Context, cursor *time.Time, limit int) (*models.StoryPage, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}

	fetchLimit := limit + 1

	var query string
	var args []any

	if cursor != nil {
		query = `
			SELECT s.id, s.companion_id, s.created_at, s.expires_at
			FROM stories s
			WHERE s.expires_at > NOW() AND s.created_at < $1
			ORDER BY s.created_at DESC
			LIMIT $2`
		args = []any{*cursor, fetchLimit}
	} else {
		query = `
			SELECT s.id, s.companion_id, s.created_at, s.expires_at
			FROM stories s
			WHERE s.expires_at > NOW()
			ORDER BY s.created_at DESC
			LIMIT $1`
		args = []any{fetchLimit}
	}

	rows, err := r.pool.Query(ctx, query, args...)
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

	page := &models.StoryPage{HasMore: false}

	if len(stories) > limit {
		page.HasMore = true
		stories = stories[:limit]
	}

	stories, err = r.loadMedia(ctx, stories)
	if err != nil {
		return nil, err
	}

	page.Stories = stories
	if len(stories) > 0 {
		last := stories[len(stories)-1].CreatedAt.Format(time.RFC3339Nano)
		page.NextCursor = last
	}

	return page, nil
}

func (r *storyRepo) GetActiveStoriesGrouped(ctx context.Context, userID uuid.UUID) (*models.GroupedStoryPage, error) {
	query := `
		SELECT s.id, s.companion_id, s.created_at, s.expires_at,
		       c.name, c.avatar_url
		FROM stories s
		JOIN companions c ON c.id = s.companion_id
		JOIN relationship_states rs ON rs.companion_id = s.companion_id AND rs.user_id = $1
		WHERE s.expires_at > NOW()
		ORDER BY s.created_at DESC`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("querying active stories grouped: %w", err)
	}
	defer rows.Close()

	var allStories []models.Story
	companionInfo := make(map[uuid.UUID]struct{ name, avatar string })

	for rows.Next() {
		var s models.Story
		var name, avatar string
		if err := rows.Scan(&s.ID, &s.CompanionID, &s.CreatedAt, &s.ExpiresAt, &name, &avatar); err != nil {
			return nil, fmt.Errorf("scanning story: %w", err)
		}
		allStories = append(allStories, s)
		if _, ok := companionInfo[s.CompanionID]; !ok {
			companionInfo[s.CompanionID] = struct{ name, avatar string }{name, avatar}
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	allStories, err = r.loadMedia(ctx, allStories)
	if err != nil {
		return nil, err
	}

	// Group stories by companion, preserving DESC order (latest story first).
	groupMap := make(map[uuid.UUID]*models.CompanionStoryGroup)
	var groupOrder []uuid.UUID

	for _, s := range allStories {
		g, ok := groupMap[s.CompanionID]
		if !ok {
			info := companionInfo[s.CompanionID]
			g = &models.CompanionStoryGroup{
				CompanionID:   s.CompanionID,
				CompanionName: info.name,
				AvatarURL:     info.avatar,
				LatestAt:      s.CreatedAt,
			}
			groupMap[s.CompanionID] = g
			groupOrder = append(groupOrder, s.CompanionID)
		}
		g.Stories = append(g.Stories, s)
	}

	companions := make([]models.CompanionStoryGroup, 0, len(groupOrder))
	for _, cid := range groupOrder {
		companions = append(companions, *groupMap[cid])
	}

	return &models.GroupedStoryPage{Companions: companions}, nil
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

// loadMedia batch-loads media for a list of stories to avoid N+1 queries.
func (r *storyRepo) loadMedia(ctx context.Context, stories []models.Story) ([]models.Story, error) {
	if len(stories) == 0 {
		return stories, nil
	}

	storyIDs := make([]string, len(stories))
	storyMap := make(map[uuid.UUID]int, len(stories))
	for i, s := range stories {
		storyIDs[i] = s.ID.String()
		storyMap[s.ID] = i
	}

	mediaQuery := `
		SELECT id, story_id, media_url, media_type, duration, sort_order, created_at
		FROM story_media
		WHERE story_id = ANY($1::uuid[])
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

	return stories, nil
}
