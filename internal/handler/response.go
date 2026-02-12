package handler

import "ai-companion-be/internal/response"

// Re-export for convenience within handler package.
var (
	JSON  = response.JSON
	Error = response.Error
)
