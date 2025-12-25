package services

import (
	"context"
	"time"
)

// LLMClient defines the interface for LLM providers with chat completion support
type LLMClient interface {
	// ChatCompletion sends a chat completion request with system and user messages
	ChatCompletion(ctx context.Context, request *ChatCompletionRequest) (*ChatCompletionResponse, error)

	// GenerateResponse generates a response from a simple prompt (backward compatibility)
	GenerateResponse(ctx context.Context, prompt string, options map[string]interface{}) (string, error)

	// IsHealthy checks if the LLM service is healthy
	IsHealthy(ctx context.Context) error

	// GetModelName returns the current model name
	GetModelName() string

	// GetProvider returns the provider name
	GetProvider() string
}

// ChatCompletionRequest represents a chat completion request
type ChatCompletionRequest struct {
	Messages    []ChatMessage          `json:"messages"`
	Model       string                 `json:"model,omitempty"`
	MaxTokens   int                    `json:"max_tokens,omitempty"`
	Temperature float64                `json:"temperature,omitempty"`
	Options     map[string]interface{} `json:"options,omitempty"`
	Stream      bool                   `json:"stream,omitempty"`
	Timeout     time.Duration          `json:"-"`
}

// ChatMessage represents a message in chat completion
type ChatMessage struct {
	Role    string `json:"role"` // system, user, assistant
	Content string `json:"content"`
}

// ChatCompletionResponse represents a chat completion response
type ChatCompletionResponse struct {
	ID      string       `json:"id"`
	Object  string       `json:"object"`
	Created int64        `json:"created"`
	Model   string       `json:"model"`
	Choices []ChatChoice `json:"choices"`
	Usage   UsageInfo    `json:"usage"`
}

// ChatChoice represents a choice in chat completion response
type ChatChoice struct {
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message"`
	FinishReason string      `json:"finish_reason"`
}

// UsageInfo represents token usage information
type UsageInfo struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// LLMError represents an LLM service error
type LLMError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Type    string `json:"type"`
}

func (e *LLMError) Error() string {
	return e.Message
}

// Common error codes
const (
	ErrCodeTimeout           = "timeout"
	ErrCodeRateLimit         = "rate_limit"
	ErrCodeInvalidAPIKey     = "invalid_api_key"
	ErrCodeModelNotFound     = "model_not_found"
	ErrCodeInsufficientQuota = "insufficient_quota"
	ErrCodeNetworkError      = "network_error"
	ErrCodeInvalidRequest    = "invalid_request"
	ErrCodeInternalError     = "internal_error"
)
