package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"ai-companion-be/internal/models"
)

// CompanionRepository defines data access operations for companions.
type CompanionRepository interface {
	GetAll(ctx context.Context) ([]models.Companion, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Companion, error)
}

type companionRepo struct {
	pool *pgxpool.Pool
}

// NewCompanionRepository creates a new CompanionRepository backed by PostgreSQL.
func NewCompanionRepository(pool *pgxpool.Pool) CompanionRepository {
	return &companionRepo{pool: pool}
}

func (r *companionRepo) GetAll(ctx context.Context) ([]models.Companion, error) {
	query := `SELECT id, name, description, avatar_url, personality, created_at FROM companions ORDER BY name`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("querying companions: %w", err)
	}
	defer rows.Close()

	var companions []models.Companion
	for rows.Next() {
		var c models.Companion
		if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.AvatarURL, &c.Personality, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("scanning companion: %w", err)
		}
		companions = append(companions, c)
	}

	return companions, rows.Err()
}

func (r *companionRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Companion, error) {
	query := `SELECT id, name, description, avatar_url, personality, created_at FROM companions WHERE id = $1`

	var c models.Companion
	err := r.pool.QueryRow(ctx, query, id).
		Scan(&c.ID, &c.Name, &c.Description, &c.AvatarURL, &c.Personality, &c.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("companion not found")
		}
		return nil, fmt.Errorf("getting companion: %w", err)
	}
	return &c, nil
}
