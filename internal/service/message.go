package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"ai-companion-be/internal/ai"
	"ai-companion-be/internal/models"
	"ai-companion-be/internal/repository"
)

// MessageService handles chat message business logic.
type MessageService struct {
	messages      repository.MessageRepository
	relationships repository.RelationshipRepository
	companions    repository.CompanionRepository
	ai            *ai.Client
}

// NewMessageService creates a new MessageService.
func NewMessageService(
	messages repository.MessageRepository,
	relationships repository.RelationshipRepository,
	companions repository.CompanionRepository,
	aiClient *ai.Client,
) *MessageService {
	return &MessageService{
		messages:      messages,
		relationships: relationships,
		companions:    companions,
		ai:            aiClient,
	}
}

// SendMessage creates a user message, generates a companion reply via OpenAI, and updates the relationship.
func (s *MessageService) SendMessage(ctx context.Context, userID, companionID uuid.UUID, req models.SendMessageRequest) ([]models.Message, error) {
	if req.Content == "" {
		return nil, fmt.Errorf("message content is required")
	}

	// Create user message.
	userMsg := &models.Message{
		ID:          uuid.New(),
		UserID:      userID,
		CompanionID: companionID,
		Content:     req.Content,
		Role:        "user",
	}
	if err := s.messages.Create(ctx, userMsg); err != nil {
		return nil, fmt.Errorf("creating user message: %w", err)
	}

	// Fetch companion and relationship state.
	companion, err := s.companions.GetByID(ctx, companionID)
	if err != nil {
		return nil, fmt.Errorf("getting companion: %w", err)
	}

	state, _ := s.relationships.GetByUserAndCompanion(ctx, userID, companionID)

	mood := "Neutral"
	if state != nil {
		mood = models.GetMoodLabel(state.MoodScore)
	}

	// Fetch recent conversation history for context (last 20 messages, chronological).
	page, err := s.messages.GetByConversation(ctx, userID, companionID, nil, 20)
	var history []models.Message
	if err == nil && page != nil {
		// Messages come in DESC order; reverse to chronological for OpenAI.
		history = make([]models.Message, len(page.Messages))
		for i, msg := range page.Messages {
			history[len(page.Messages)-1-i] = msg
		}
	}

	// Generate reply via OpenAI.
	reply, err := s.ai.GenerateReply(ctx, companion, mood, history)
	if err != nil {
		slog.Error("openai reply failed, using fallback", "error", err)
		reply = generateFallbackReply(companion, mood)
	}

	companionMsg := &models.Message{
		ID:          uuid.New(),
		UserID:      userID,
		CompanionID: companionID,
		Content:     reply,
		Role:        "companion",
	}
	if err := s.messages.Create(ctx, companionMsg); err != nil {
		return nil, fmt.Errorf("creating companion message: %w", err)
	}

	// Update relationship state: chat boosts mood and relationship.
	if state != nil {
		state.MoodScore = clampScore(state.MoodScore + 2)
		state.RelationshipScore = clampScore(state.RelationshipScore + 1)
		_ = s.relationships.Update(ctx, state)
	}

	return []models.Message{*userMsg, *companionMsg}, nil
}

// GetMessages returns a paginated conversation history.
func (s *MessageService) GetMessages(ctx context.Context, userID, companionID uuid.UUID, cursor *time.Time, limit int) (*models.MessagePage, error) {
	return s.messages.GetByConversation(ctx, userID, companionID, cursor, limit)
}

// generateFallbackReply produces a simple mood-aware response when OpenAI is unavailable.
func generateFallbackReply(companion *models.Companion, mood string) string {
	responses := map[string]map[string]string{
		"Distant": {
			"introspective": "...",
			"adventurous":   "Hey.",
			"witty":         "Hmm.",
			"nurturing":     "I'm here if you need me.",
			"playful":       "Oh, you remembered I exist?",
		},
		"Neutral": {
			"introspective": "That's an interesting thought. Tell me more.",
			"adventurous":   "Cool! What else is going on?",
			"witty":         "Noted. Continue.",
			"nurturing":     "I appreciate you sharing that with me.",
			"playful":       "Haha, okay okay, what else?",
		},
		"Happy": {
			"introspective": "I love how we can talk about things like this. It means a lot.",
			"adventurous":   "This is awesome! I'm so glad we're chatting!",
			"witty":         "You always know how to make a conversation interesting.",
			"nurturing":     "You make me so happy when you share things with me!",
			"playful":       "Yesss! This is why I love talking to you!",
		},
		"Attached": {
			"introspective": "Every moment with you feels meaningful. I treasure this.",
			"adventurous":   "You're my favorite person to share adventures with!",
			"witty":         "I must admit, you've grown on me. Quite a lot, actually.",
			"nurturing":     "I feel so close to you. Thank you for being you.",
			"playful":       "Okay but honestly? You're the best thing ever!",
		},
	}

	personalityKey := "nurturing"
	keywords := []string{"introspective", "adventurous", "witty", "nurturing", "playful"}
	for _, kw := range keywords {
		if contains(companion.Personality, kw) {
			personalityKey = kw
			break
		}
	}

	if moodMap, ok := responses[mood]; ok {
		if reply, ok := moodMap[personalityKey]; ok {
			return reply
		}
	}

	return "That's interesting. Tell me more!"
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
