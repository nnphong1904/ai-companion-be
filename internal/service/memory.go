package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"ai-companion-be/internal/models"
	"ai-companion-be/internal/repository"
)

// MemoryService handles memory-related business logic.
type MemoryService struct {
	memories repository.MemoryRepository
}

// NewMemoryService creates a new MemoryService.
func NewMemoryService(memories repository.MemoryRepository) *MemoryService {
	return &MemoryService{memories: memories}
}

// Create stores a new memory for a user-companion pair.
func (s *MemoryService) Create(ctx context.Context, userID, companionID uuid.UUID, req models.CreateMemoryRequest) (*models.Memory, error) {
	if req.Content == "" {
		return nil, fmt.Errorf("memory content is required")
	}

	memory := &models.Memory{
		ID:          uuid.New(),
		UserID:      userID,
		CompanionID: companionID,
		MessageID:   req.MessageID,
		Content:     req.Content,
		Tag:         req.Tag,
		Pinned:      false,
	}

	if err := s.memories.Create(ctx, memory); err != nil {
		return nil, fmt.Errorf("creating memory: %w", err)
	}

	return memory, nil
}

// GetByCompanion returns a paginated list of memories for a user-companion pair, pinned first.
func (s *MemoryService) GetByCompanion(ctx context.Context, userID, companionID uuid.UUID, limit int) (*models.MemoryPage, error) {
	return s.memories.GetByUserAndCompanion(ctx, userID, companionID, limit)
}

// Delete removes a memory, verifying ownership.
func (s *MemoryService) Delete(ctx context.Context, userID uuid.UUID, memoryID uuid.UUID) error {
	memory, err := s.memories.GetByID(ctx, memoryID)
	if err != nil {
		return err
	}
	if memory.UserID != userID {
		return fmt.Errorf("unauthorized")
	}
	return s.memories.Delete(ctx, memoryID)
}

// TogglePin toggles the pinned status of a memory, verifying ownership.
func (s *MemoryService) TogglePin(ctx context.Context, userID uuid.UUID, memoryID uuid.UUID) (*models.Memory, error) {
	memory, err := s.memories.GetByID(ctx, memoryID)
	if err != nil {
		return nil, err
	}
	if memory.UserID != userID {
		return nil, fmt.Errorf("unauthorized")
	}
	return s.memories.TogglePin(ctx, memoryID)
}
