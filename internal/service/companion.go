package service

import (
	"context"

	"github.com/google/uuid"

	"ai-companion-be/internal/models"
	"ai-companion-be/internal/repository"
)

// CompanionService handles companion-related business logic.
type CompanionService struct {
	companions repository.CompanionRepository
}

// NewCompanionService creates a new CompanionService.
func NewCompanionService(companions repository.CompanionRepository) *CompanionService {
	return &CompanionService{companions: companions}
}

// GetAll returns all available companions.
func (s *CompanionService) GetAll(ctx context.Context) ([]models.Companion, error) {
	return s.companions.GetAll(ctx)
}

// GetByID returns a single companion.
func (s *CompanionService) GetByID(ctx context.Context, id uuid.UUID) (*models.Companion, error) {
	return s.companions.GetByID(ctx, id)
}
