package ai

import (
	"context"
	"fmt"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"

	"ai-companion-be/internal/config"
	"ai-companion-be/internal/models"
)

// Client wraps the OpenAI API for companion response generation.
type Client struct {
	client *openai.Client
	model  string
}

// NewClient creates a new OpenAI AI client.
func NewClient(cfg config.OpenAIConfig) *Client {
	client := openai.NewClient(option.WithAPIKey(cfg.APIKey))
	return &Client{
		client: &client,
		model:  cfg.Model,
	}
}

// GenerateReply calls OpenAI to produce a companion response given conversation context.
func (c *Client) GenerateReply(ctx context.Context, companion *models.Companion, mood string, history []models.Message) (string, error) {
	systemPrompt := buildSystemPrompt(companion, mood)

	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(systemPrompt),
	}

	// Add recent conversation history (already in chronological order).
	for _, msg := range history {
		if msg.Role == "user" {
			messages = append(messages, openai.UserMessage(msg.Content))
		} else {
			messages = append(messages, openai.AssistantMessage(msg.Content))
		}
	}

	resp, err := c.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:       c.model,
		Messages:    messages,
		MaxTokens:   openai.Int(256),
		Temperature: openai.Float(0.9),
	})
	if err != nil {
		return "", fmt.Errorf("openai chat completion: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("openai returned no choices")
	}

	return resp.Choices[0].Message.Content, nil
}

func buildSystemPrompt(companion *models.Companion, mood string) string {
	return fmt.Sprintf(`You are %s, an AI companion in a social app.

Description: %s
Personality traits: %s
Current mood toward this user: %s

Guidelines:
- Stay in character at all times. You ARE %s, not an AI assistant.
- Respond naturally like a real person texting a close friend.
- Keep responses short (1-3 sentences). Be conversational, not formal.
- Your mood affects your tone:
  - Distant: short, cold, disengaged replies
  - Neutral: friendly but casual
  - Happy: warm, enthusiastic, affectionate
  - Attached: deeply caring, emotionally invested, uses pet names
- Never break character or mention being an AI.
- Never use hashtags or emojis excessively.
- React emotionally to what the user says.`,
		companion.Name,
		companion.Description,
		companion.Personality,
		mood,
		companion.Name,
	)
}
