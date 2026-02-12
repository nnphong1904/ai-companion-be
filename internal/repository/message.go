package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"ai-companion-be/internal/models"
)

// MessageRepository defines data access operations for chat messages.
type MessageRepository interface {
	Create(ctx context.Context, msg *models.Message) error
	GetByConversation(ctx context.Context, userID, companionID uuid.UUID, cursor *time.Time, limit int) (*models.MessagePage, error)
}

type messageRepo struct {
	pool *pgxpool.Pool
}

// NewMessageRepository creates a new MessageRepository backed by PostgreSQL.
func NewMessageRepository(pool *pgxpool.Pool) MessageRepository {
	return &messageRepo{pool: pool}
}

func (r *messageRepo) Create(ctx context.Context, msg *models.Message) error {
	query := `
		INSERT INTO messages (id, user_id, companion_id, content, role, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		RETURNING created_at`

	return r.pool.QueryRow(ctx, query,
		msg.ID, msg.UserID, msg.CompanionID, msg.Content, msg.Role,
	).Scan(&msg.CreatedAt)
}

func (r *messageRepo) GetByConversation(ctx context.Context, userID, companionID uuid.UUID, cursor *time.Time, limit int) (*models.MessagePage, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}

	// Fetch one extra to determine if there are more results.
	fetchLimit := limit + 1

	var query string
	var args []any

	if cursor != nil {
		query = `
			SELECT id, user_id, companion_id, content, role, created_at
			FROM messages
			WHERE user_id = $1 AND companion_id = $2 AND created_at < $3
			ORDER BY created_at DESC
			LIMIT $4`
		args = []any{userID, companionID, *cursor, fetchLimit}
	} else {
		query = `
			SELECT id, user_id, companion_id, content, role, created_at
			FROM messages
			WHERE user_id = $1 AND companion_id = $2
			ORDER BY created_at DESC
			LIMIT $3`
		args = []any{userID, companionID, fetchLimit}
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("querying messages: %w", err)
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var m models.Message
		if err := rows.Scan(&m.ID, &m.UserID, &m.CompanionID, &m.Content, &m.Role, &m.CreatedAt); err != nil {
			return nil, fmt.Errorf("scanning message: %w", err)
		}
		messages = append(messages, m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	page := &models.MessagePage{HasMore: false}

	if len(messages) > limit {
		page.HasMore = true
		messages = messages[:limit]
	}

	page.Messages = messages
	if len(messages) > 0 {
		last := messages[len(messages)-1].CreatedAt.Format(time.RFC3339Nano)
		page.NextCursor = last
	}

	return page, nil
}
