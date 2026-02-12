package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"ai-companion-be/internal/models"
)

// UserRepository defines data access operations for users.
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
}

type userRepo struct {
	pool *pgxpool.Pool
}

// NewUserRepository creates a new UserRepository backed by PostgreSQL.
func NewUserRepository(pool *pgxpool.Pool) UserRepository {
	return &userRepo{pool: pool}
}

func (r *userRepo) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (id, email, password, name, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING created_at, updated_at`

	return r.pool.QueryRow(ctx, query, user.ID, user.Email, user.Password, user.Name).
		Scan(&user.CreatedAt, &user.UpdatedAt)
}

func (r *userRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `SELECT id, email, password, name, created_at, updated_at FROM users WHERE id = $1`

	var user models.User
	err := r.pool.QueryRow(ctx, query, id).
		Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("getting user by id: %w", err)
	}
	return &user, nil
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT id, email, password, name, created_at, updated_at FROM users WHERE email = $1`

	var user models.User
	err := r.pool.QueryRow(ctx, query, email).
		Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("getting user by email: %w", err)
	}
	return &user, nil
}
