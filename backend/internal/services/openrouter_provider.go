package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"ragify-backend/internal/config"
)

// OpenRouterProvider implements EmbeddingProvider for OpenRouter API
type OpenRouterProvider struct {
	config     config.EmbeddingConfig
	httpClient *http.Client
}

// NewOpenRouterProvider creates a new OpenRouter embedding provider
func NewOpenRouterProvider(config config.EmbeddingConfig) *OpenRouterProvider {
	if config.BaseURL == "" {
		config.BaseURL = "https://openrouter.ai/api/v1"
	}
	if config.ModelName == "" {
		config.ModelName = "openai/text-embedding-3-small"
	}

	return &OpenRouterProvider{
		config: config,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

// openRouterRequest represents the request structure for OpenRouter API
type openRouterRequest struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
}

// openRouterResponse represents the response structure from OpenRouter API
type openRouterResponse struct {
	Model string `json:"model"`
	Data  []struct {
		Object    string    `json:"object"`
		Embedding []float32 `json:"embedding"`
		Index     int       `json:"index"`
	} `json:"data"`
	Usage struct {
		PromptTokens int `json:"prompt_tokens"`
		TotalTokens  int `json:"total_tokens"`
	} `json:"usage"`
	Created int64 `json:"created"`
}

// openRouterError represents error response from OpenRouter
type openRouterError struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error"`
}

// GenerateEmbedding generates embedding for a single text
func (p *OpenRouterProvider) GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	embeddings, err := p.GenerateBatchEmbeddings(ctx, []string{text})
	if err != nil {
		return nil, err
	}

	if len(embeddings) == 0 {
		return nil, fmt.Errorf("no embeddings returned")
	}

	return embeddings[0].Embedding, nil
}

// GenerateBatchEmbeddings generates embeddings for multiple texts with retry logic
func (p *OpenRouterProvider) GenerateBatchEmbeddings(ctx context.Context, texts []string) ([]EmbeddingData, error) {
	if len(texts) == 0 {
		return []EmbeddingData{}, nil
	}

	// Process in batches if needed
	batchSize := p.config.BatchSize
	if batchSize <= 0 {
		batchSize = 100
	}

	var allEmbeddings []EmbeddingData

	for i := 0; i < len(texts); i += batchSize {
		end := i + batchSize
		if end > len(texts) {
			end = len(texts)
		}

		batch := texts[i:end]
		embeddings, err := p.generateBatchEmbeddingsWithRetry(ctx, batch)
		if err != nil {
			return nil, fmt.Errorf("failed to generate embeddings for batch %d: %w", i/batchSize, err)
		}

		// Adjust indices for global position
		for j := range embeddings {
			embeddings[j].Index = i + j
		}

		allEmbeddings = append(allEmbeddings, embeddings...)
	}

	return allEmbeddings, nil
}

// generateBatchEmbeddingsWithRetry implements retry logic for batch embedding generation
func (p *OpenRouterProvider) generateBatchEmbeddingsWithRetry(ctx context.Context, texts []string) ([]EmbeddingData, error) {
	var lastErr error

	for attempt := 0; attempt <= p.config.MaxRetries; attempt++ {
		if attempt > 0 {
			// Add delay before retry
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(p.config.RetryDelay):
			}
		}

		embeddings, err := p.generateBatchEmbeddings(ctx, texts)
		if err == nil {
			return embeddings, nil
		}

		lastErr = err

		// Don't retry on certain errors
		if isNonRetryableError(err) {
			break
		}
	}

	return nil, lastErr
}

// generateBatchEmbeddings performs the actual API call
func (p *OpenRouterProvider) generateBatchEmbeddings(ctx context.Context, texts []string) ([]EmbeddingData, error) {
	reqBody := openRouterRequest{
		Model: p.config.ModelName,
		Input: texts,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", p.config.BaseURL+"/embeddings", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.config.APIKey)
	req.Header.Set("HTTP-Referer", "https://github.com/abhinavsaxena2308/Ragify-Rag-system")
	req.Header.Set("X-Title", "Ragify RAG System")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errorResp openRouterError
		if err := json.Unmarshal(body, &errorResp); err == nil && errorResp.Error.Message != "" {
			return nil, fmt.Errorf("OpenRouter API error: %s (code: %s)", errorResp.Error.Message, errorResp.Error.Code)
		}
		return nil, fmt.Errorf("OpenRouter API error: status %d, body: %s", resp.StatusCode, string(body))
	}

	var apiResp openRouterResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Convert to our internal format
	embeddings := make([]EmbeddingData, len(apiResp.Data))
	for i, data := range apiResp.Data {
		embeddings[i] = EmbeddingData{
			Object:    data.Object,
			Embedding: data.Embedding,
			Index:     data.Index,
		}
	}

	return embeddings, nil
}

// GetModelName returns the configured model name
func (p *OpenRouterProvider) GetModelName() string {
	return p.config.ModelName
}

// GetDimensions returns the embedding dimensions for the model
func (p *OpenRouterProvider) GetDimensions() int {
	// Common embedding dimensions for popular models
	switch p.config.ModelName {
	case "openai/text-embedding-3-small":
		return 1536
	case "openai/text-embedding-3-large":
		return 3072
	case "openai/text-embedding-ada-002":
		return 1536
	case "cohere/embed-english-v3.0":
		return 1024
	case "cohere/embed-multilingual-v3.0":
		return 1024
	default:
		// Default to 1536 for unknown models
		return 1536
	}
}

// IsHealthy checks if the provider is accessible
func (p *OpenRouterProvider) IsHealthy(ctx context.Context) error {
	// Test with a simple embedding request
	_, err := p.GenerateEmbedding(ctx, "health check")
	return err
}

// isNonRetryableError determines if an error should not be retried
func isNonRetryableError(err error) bool {
	if err == nil {
		return false
	}

	errStr := err.Error()
	// Don't retry on authentication errors, invalid model, etc.
	nonRetryableErrors := []string{
		"invalid api key",
		"unauthorized",
		"forbidden",
		"invalid model",
		"not found",
		"rate limit exceeded",
	}

	for _, nonRetryable := range nonRetryableErrors {
		if contains(errStr, nonRetryable) {
			return true
		}
	}

	return false
}

// contains checks if a string contains a substring (case-insensitive)
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			len(s) > len(substr) &&
				(s[:len(substr)] == substr ||
					s[len(s)-len(substr):] == substr ||
					containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
