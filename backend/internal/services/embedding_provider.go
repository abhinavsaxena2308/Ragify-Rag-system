package services

import (
	"context"
)

// EmbeddingRequest represents a request to generate embeddings
type EmbeddingRequest struct {
	Texts      []string `json:"texts"`
	ModelName  string   `json:"model_name,omitempty"`
	Dimensions *int     `json:"dimensions,omitempty"`
}

// EmbeddingResponse represents the response from embedding generation
type EmbeddingResponse struct {
	Model   string          `json:"model"`
	Data    []EmbeddingData `json:"data"`
	Usage   Usage           `json:"usage"`
	Created int64           `json:"created"`
}

// EmbeddingData represents a single embedding
type EmbeddingData struct {
	Object    string    `json:"object"`
	Embedding []float32 `json:"embedding"`
	Index     int       `json:"index"`
}

// Usage represents token usage information
type Usage struct {
	PromptTokens int `json:"prompt_tokens"`
	TotalTokens  int `json:"total_tokens"`
}

// EmbeddingProvider interface for different embedding providers
type EmbeddingProvider interface {
	// GenerateEmbedding generates embedding for a single text
	GenerateEmbedding(ctx context.Context, text string) ([]float32, error)

	// GenerateBatchEmbeddings generates embeddings for multiple texts
	GenerateBatchEmbeddings(ctx context.Context, texts []string) ([]EmbeddingData, error)

	// GetModelName returns the configured model name
	GetModelName() string

	// GetDimensions returns the embedding dimensions for the model
	GetDimensions() int

	// IsHealthy checks if the provider is accessible
	IsHealthy(ctx context.Context) error
}
