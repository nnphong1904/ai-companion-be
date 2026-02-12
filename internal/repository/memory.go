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

// MemoryRepository defines data access operations for memories.
type MemoryRepository interface {
	Create(ctx context.Context, memory *models.Memory) error
	GetByUserAndCompanion(ctx context.Context, userID, companionID uuid.UUID, limit int) (*models.MemoryPage, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Memory, error)
	Delete(ctx context.Context, id uuid.UUID) error
	TogglePin(ctx context.Context, id uuid.UUID) (*models.Memory, error)
}

type memoryRepo struct {
	pool *pgxpool.Pool
}

// NewMemoryRepository creates a new MemoryRepository backed by PostgreSQL.
func NewMemoryRepository(pool *pgxpool.Pool) MemoryRepository {
	return &memoryRepo{pool: pool}
}

func (r *memoryRepo) Create(ctx context.Context, memory *models.Memory) error {
	query := `
		INSERT INTO memories (id, user_id, companion_id, content, tag, pinned, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
		RETURNING created_at`

	return r.pool.QueryRow(ctx, query,
		memory.ID, memory.UserID, memory.CompanionID, memory.Content, memory.Tag, memory.Pinned,
	).Scan(&memory.CreatedAt)
}

func (r *memoryRepo) GetByUserAndCompanion(ctx context.Context, userID, companionID uuid.UUID, limit int) (*models.MemoryPage, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	fetchLimit := limit + 1

	query := `
		SELECT id, user_id, companion_id, content, tag, pinned, created_at
		FROM memories
		WHERE user_id = $1 AND companion_id = $2
		ORDER BY pinned DESC, created_at DESC
		LIMIT $3`

	rows, err := r.pool.Query(ctx, query, userID, companionID, fetchLimit)
	if err != nil {
		return nil, fmt.Errorf("querying memories: %w", err)
	}
	defer rows.Close()

	var memories []models.Memory
	for rows.Next() {
		var m models.Memory
		if err := rows.Scan(&m.ID, &m.UserID, &m.CompanionID, &m.Content, &m.Tag, &m.Pinned, &m.CreatedAt); err != nil {
			return nil, fmt.Errorf("scanning memory: %w", err)
		}
		memories = append(memories, m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	page := &models.MemoryPage{HasMore: false}

	if len(memories) > limit {
		page.HasMore = true
		memories = memories[:limit]
	}

	page.Memories = memories
	if len(memories) > 0 {
		last := memories[len(memories)-1].CreatedAt.Format(time.RFC3339Nano)
		page.NextCursor = last
	}

	return page, nil
}

func (r *memoryRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Memory, error) {
	query := `SELECT id, user_id, companion_id, content, tag, pinned, created_at FROM memories WHERE id = $1`

	var m models.Memory
	err := r.pool.QueryRow(ctx, query, id).
		Scan(&m.ID, &m.UserID, &m.CompanionID, &m.Content, &m.Tag, &m.Pinned, &m.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("memory not found")
		}
		return nil, fmt.Errorf("getting memory: %w", err)
	}
	return &m, nil
}

func (r *memoryRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM memories WHERE id = $1`

	tag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("deleting memory: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("memory not found")
	}
	return nil
}

func (r *memoryRepo) TogglePin(ctx context.Context, id uuid.UUID) (*models.Memory, error) {
	query := `
		UPDATE memories SET pinned = NOT pinned
		WHERE id = $1
		RETURNING id, user_id, companion_id, content, tag, pinned, created_at`

	var m models.Memory
	err := r.pool.QueryRow(ctx, query, id).
		Scan(&m.ID, &m.UserID, &m.CompanionID, &m.Content, &m.Tag, &m.Pinned, &m.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("memory not found")
		}
		return nil, fmt.Errorf("toggling pin: %w", err)
	}
	return &m, nil
}
