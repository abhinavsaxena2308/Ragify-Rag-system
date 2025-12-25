package services

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"ragify-backend/internal/config"
	"ragify-backend/internal/models"
)

// EmbeddingService provides high-level embedding generation functionality
type EmbeddingService struct {
	provider EmbeddingProvider
	config   config.EmbeddingConfig
	mu       sync.RWMutex
}

// NewEmbeddingService creates a new embedding service with the specified provider
func NewEmbeddingService(provider EmbeddingProvider, config config.EmbeddingConfig) *EmbeddingService {
	return &EmbeddingService{
		provider: provider,
		config:   config,
	}
}

// NewOpenRouterEmbeddingService creates a new embedding service using OpenRouter
func NewOpenRouterEmbeddingService(config config.EmbeddingConfig) *EmbeddingService {
	provider := NewOpenRouterProvider(config)
	return NewEmbeddingService(provider, config)
}

// GenerateChunkEmbedding generates embedding for a single chunk
func (s *EmbeddingService) GenerateChunkEmbedding(ctx context.Context, chunk *models.Chunk) ([]float32, error) {
	if chunk == nil {
		return nil, fmt.Errorf("chunk cannot be nil")
	}

	if chunk.Content == "" {
		return nil, fmt.Errorf("chunk content cannot be empty")
	}

	embedding, err := s.generateWithRetry(ctx, chunk.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to generate embedding for chunk %d: %w", chunk.ID, err)
	}

	return embedding, nil
}

// GenerateDocumentEmbeddings generates embeddings for all chunks in a document
func (s *EmbeddingService) GenerateDocumentEmbeddings(ctx context.Context, document *models.Document, chunks []models.Chunk) (map[uint][]float32, error) {
	if document == nil {
		return nil, fmt.Errorf("document cannot be nil")
	}

	if len(chunks) == 0 {
		return make(map[uint][]float32), nil
	}

	// Extract chunk contents
	texts := make([]string, len(chunks))
	chunkIDs := make([]uint, len(chunks))

	for i, chunk := range chunks {
		if chunk.Content == "" {
			return nil, fmt.Errorf("chunk %d has empty content", chunk.ID)
		}
		texts[i] = chunk.Content
		chunkIDs[i] = chunk.ID
	}

	// Generate embeddings in batches
	embeddings, err := s.generateBatchWithRetry(ctx, texts)
	if err != nil {
		return nil, fmt.Errorf("failed to generate embeddings for document %d: %w", document.ID, err)
	}

	// Map embeddings back to chunk IDs
	result := make(map[uint][]float32)
	for i, embedding := range embeddings {
		result[chunkIDs[i]] = embedding.Embedding
	}

	return result, nil
}

// GenerateTextEmbedding generates embedding for arbitrary text
func (s *EmbeddingService) GenerateTextEmbedding(ctx context.Context, text string) ([]float32, error) {
	if text == "" {
		return nil, fmt.Errorf("text cannot be empty")
	}

	return s.generateWithRetry(ctx, text)
}

// GenerateQueryEmbedding generates embedding for a search query
func (s *EmbeddingService) GenerateQueryEmbedding(ctx context.Context, query string) ([]float32, error) {
	if query == "" {
		return nil, fmt.Errorf("query cannot be empty")
	}

	// For queries, we might want to use the same process but could add query-specific preprocessing
	return s.generateWithRetry(ctx, query)
}

// generateWithRetry implements retry logic for single embedding generation
func (s *EmbeddingService) generateWithRetry(ctx context.Context, text string) ([]float32, error) {
	var lastErr error

	for attempt := 0; attempt <= s.config.MaxRetries; attempt++ {
		if attempt > 0 {
			// Add exponential backoff delay
			delay := time.Duration(attempt) * s.config.RetryDelay
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(delay):
			}

			log.Printf("Embedding generation retry attempt %d for text: %.50s...", attempt+1, text)
		}

		embedding, err := s.provider.GenerateEmbedding(ctx, text)
		if err == nil {
			return embedding, nil
		}

		lastErr = err

		// Don't retry on certain errors
		if isNonRetryableError(err) {
			break
		}
	}

	return nil, fmt.Errorf("failed after %d attempts: %w", s.config.MaxRetries+1, lastErr)
}

// generateBatchWithRetry implements retry logic for batch embedding generation
func (s *EmbeddingService) generateBatchWithRetry(ctx context.Context, texts []string) ([]EmbeddingData, error) {
	var lastErr error

	for attempt := 0; attempt <= s.config.MaxRetries; attempt++ {
		if attempt > 0 {
			// Add exponential backoff delay
			delay := time.Duration(attempt) * s.config.RetryDelay
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(delay):
			}

			log.Printf("Batch embedding generation retry attempt %d for %d texts", attempt+1, len(texts))
		}

		embeddings, err := s.provider.GenerateBatchEmbeddings(ctx, texts)
		if err == nil {
			return embeddings, nil
		}

		lastErr = err

		// Don't retry on certain errors
		if isNonRetryableError(err) {
			break
		}
	}

	return nil, fmt.Errorf("failed after %d attempts: %w", s.config.MaxRetries+1, lastErr)
}

// GetProvider returns the current embedding provider
func (s *EmbeddingService) GetProvider() EmbeddingProvider {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.provider
}

// SetProvider updates the embedding provider
func (s *EmbeddingService) SetProvider(provider EmbeddingProvider) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.provider = provider
}

// GetModelName returns the current model name
func (s *EmbeddingService) GetModelName() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.provider.GetModelName()
}

// GetDimensions returns the embedding dimensions
func (s *EmbeddingService) GetDimensions() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.provider.GetDimensions()
}

// IsHealthy checks if the embedding service is healthy
func (s *EmbeddingService) IsHealthy(ctx context.Context) error {
	s.mu.RLock()
	provider := s.provider
	s.mu.RUnlock()

	return provider.IsHealthy(ctx)
}

// GetStats returns service statistics
func (s *EmbeddingService) GetStats() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return map[string]interface{}{
		"provider_type": fmt.Sprintf("%T", s.provider),
		"model_name":    s.provider.GetModelName(),
		"dimensions":    s.provider.GetDimensions(),
		"timeout":       s.config.Timeout.String(),
		"max_retries":   s.config.MaxRetries,
		"batch_size":    s.config.BatchSize,
		"retry_delay":   s.config.RetryDelay.String(),
	}
}

// ValidateEmbedding validates an embedding vector
func (s *EmbeddingService) ValidateEmbedding(embedding []float32) error {
	expectedDimensions := s.GetDimensions()

	if embedding == nil {
		return fmt.Errorf("embedding cannot be nil")
	}

	if len(embedding) != expectedDimensions {
		return fmt.Errorf("embedding dimension mismatch: expected %d, got %d", expectedDimensions, len(embedding))
	}

	// Check for invalid values
	for i, val := range embedding {
		if val != val { // NaN check
			return fmt.Errorf("embedding contains NaN at index %d", i)
		}
	}

	return nil
}

// EmbeddingGenerationOptions provides options for embedding generation
type EmbeddingGenerationOptions struct {
	Timeout    time.Duration
	MaxRetries int
	BatchSize  int
	ModelName  string
	Dimensions *int
}

// WithOptions applies options to the embedding service
func (s *EmbeddingService) WithOptions(opts EmbeddingGenerationOptions) *EmbeddingService {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Create a copy of the config
	newConfig := s.config

	if opts.Timeout > 0 {
		newConfig.Timeout = opts.Timeout
	}
	if opts.MaxRetries > 0 {
		newConfig.MaxRetries = opts.MaxRetries
	}
	if opts.BatchSize > 0 {
		newConfig.BatchSize = opts.BatchSize
	}
	if opts.ModelName != "" {
		newConfig.ModelName = opts.ModelName
	}

	// Create new service with updated config
	return &EmbeddingService{
		provider: s.provider,
		config:   newConfig,
	}
}
