package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"ai-companion-be/internal/models"
)

// RelationshipRepository defines data access operations for relationship states.
type RelationshipRepository interface {
	Create(ctx context.Context, state *models.RelationshipState) error
	GetByUserAndCompanion(ctx context.Context, userID, companionID uuid.UUID) (*models.RelationshipState, error)
	GetAllByUser(ctx context.Context, userID uuid.UUID) ([]models.RelationshipState, error)
	Update(ctx context.Context, state *models.RelationshipState) error
}

type relationshipRepo struct {
	pool *pgxpool.Pool
}

// NewRelationshipRepository creates a new RelationshipRepository backed by PostgreSQL.
func NewRelationshipRepository(pool *pgxpool.Pool) RelationshipRepository {
	return &relationshipRepo{pool: pool}
}

func (r *relationshipRepo) Create(ctx context.Context, state *models.RelationshipState) error {
	query := `
		INSERT INTO relationship_states (id, user_id, companion_id, mood_score, relationship_score, last_interaction, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING last_interaction, updated_at`

	return r.pool.QueryRow(ctx, query,
		state.ID, state.UserID, state.CompanionID, state.MoodScore, state.RelationshipScore,
	).Scan(&state.LastInteraction, &state.UpdatedAt)
}

func (r *relationshipRepo) GetByUserAndCompanion(ctx context.Context, userID, companionID uuid.UUID) (*models.RelationshipState, error) {
	query := `
		SELECT id, user_id, companion_id, mood_score, relationship_score, last_interaction, updated_at
		FROM relationship_states
		WHERE user_id = $1 AND companion_id = $2`

	var s models.RelationshipState
	err := r.pool.QueryRow(ctx, query, userID, companionID).
		Scan(&s.ID, &s.UserID, &s.CompanionID, &s.MoodScore, &s.RelationshipScore, &s.LastInteraction, &s.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("relationship not found")
		}
		return nil, fmt.Errorf("getting relationship state: %w", err)
	}
	return &s, nil
}

func (r *relationshipRepo) GetAllByUser(ctx context.Context, userID uuid.UUID) ([]models.RelationshipState, error) {
	query := `
		SELECT id, user_id, companion_id, mood_score, relationship_score, last_interaction, updated_at
		FROM relationship_states
		WHERE user_id = $1
		ORDER BY updated_at DESC`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("querying relationships: %w", err)
	}
	defer rows.Close()

	var states []models.RelationshipState
	for rows.Next() {
		var s models.RelationshipState
		if err := rows.Scan(&s.ID, &s.UserID, &s.CompanionID, &s.MoodScore, &s.RelationshipScore, &s.LastInteraction, &s.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scanning relationship: %w", err)
		}
		states = append(states, s)
	}

	return states, rows.Err()
}

func (r *relationshipRepo) Update(ctx context.Context, state *models.RelationshipState) error {
	query := `
		UPDATE relationship_states
		SET mood_score = $1, relationship_score = $2, last_interaction = NOW(), updated_at = NOW()
		WHERE id = $3
		RETURNING last_interaction, updated_at`

	return r.pool.QueryRow(ctx, query, state.MoodScore, state.RelationshipScore, state.ID).
		Scan(&state.LastInteraction, &state.UpdatedAt)
}
