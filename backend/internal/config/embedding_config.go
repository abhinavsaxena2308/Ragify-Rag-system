package config

import (
	"os"
	"strconv"
	"time"
)

// EmbeddingConfig holds configuration for embedding providers
type EmbeddingConfig struct {
	APIKey     string        `json:"api_key"`
	ModelName  string        `json:"model_name"`
	BaseURL    string        `json:"base_url"`
	Timeout    time.Duration `json:"timeout"`
	MaxRetries int           `json:"max_retries"`
	BatchSize  int           `json:"batch_size"`
	RetryDelay time.Duration `json:"retry_delay"`
}

// DefaultEmbeddingConfig returns default configuration for embedding providers
func DefaultEmbeddingConfig() EmbeddingConfig {
	return EmbeddingConfig{
		ModelName:  "openai/text-embedding-3-small",
		Timeout:    30 * time.Second,
		MaxRetries: 3,
		BatchSize:  100,
		RetryDelay: 1 * time.Second,
	}
}

// EmbeddingConfigLoader handles loading embedding configuration from environment variables
type EmbeddingConfigLoader struct{}

// NewEmbeddingConfigLoader creates a new embedding configuration loader
func NewEmbeddingConfigLoader() *EmbeddingConfigLoader {
	return &EmbeddingConfigLoader{}
}

// LoadFromEnv loads embedding configuration from environment variables
func (l *EmbeddingConfigLoader) LoadFromEnv() EmbeddingConfig {
	config := DefaultEmbeddingConfig()

	// Load API key
	if apiKey := os.Getenv("OPENROUTER_API_KEY"); apiKey != "" {
		config.APIKey = apiKey
	}

	// Load model name
	if modelName := os.Getenv("EMBEDDING_MODEL"); modelName != "" {
		config.ModelName = modelName
	}

	// Load base URL
	if baseURL := os.Getenv("OPENROUTER_BASE_URL"); baseURL != "" {
		config.BaseURL = baseURL
	}

	// Load timeout
	if timeoutStr := os.Getenv("EMBEDDING_TIMEOUT"); timeoutStr != "" {
		if timeout, err := time.ParseDuration(timeoutStr); err == nil {
			config.Timeout = timeout
		}
	}

	// Load max retries
	if maxRetriesStr := os.Getenv("EMBEDDING_MAX_RETRIES"); maxRetriesStr != "" {
		if maxRetries, err := strconv.Atoi(maxRetriesStr); err == nil {
			config.MaxRetries = maxRetries
		}
	}

	// Load batch size
	if batchSizeStr := os.Getenv("EMBEDDING_BATCH_SIZE"); batchSizeStr != "" {
		if batchSize, err := strconv.Atoi(batchSizeStr); err == nil {
			config.BatchSize = batchSize
		}
	}

	// Load retry delay
	if retryDelayStr := os.Getenv("EMBEDDING_RETRY_DELAY"); retryDelayStr != "" {
		if retryDelay, err := time.ParseDuration(retryDelayStr); err == nil {
			config.RetryDelay = retryDelay
		}
	}

	return config
}

// LoadFromFile loads embedding configuration from a file (future enhancement)
func (l *EmbeddingConfigLoader) LoadFromFile(filename string) (EmbeddingConfig, error) {
	// This could be implemented later to support JSON/YAML config files
	return l.LoadFromEnv(), nil
}

// GetSupportedModels returns a list of supported embedding models
func (l *EmbeddingConfigLoader) GetSupportedModels() []EmbeddingModel {
	return []EmbeddingModel{
		{
			Name:        "openai/text-embedding-3-small",
			DisplayName: "OpenAI Text Embedding 3 Small",
			Description: "Fast and efficient embedding model with 1536 dimensions",
			Dimensions:  1536,
			Provider:    "OpenRouter",
		},
		{
			Name:        "openai/text-embedding-3-large",
			DisplayName: "OpenAI Text Embedding 3 Large",
			Description: "High-quality embedding model with 3072 dimensions",
			Dimensions:  3072,
			Provider:    "OpenRouter",
		},
		{
			Name:        "openai/text-embedding-ada-002",
			DisplayName: "OpenAI Text Embedding Ada 002",
			Description: "Legacy embedding model with 1536 dimensions",
			Dimensions:  1536,
			Provider:    "OpenRouter",
		},
		{
			Name:        "cohere/embed-english-v3.0",
			DisplayName: "Cohere Embed English v3.0",
			Description: "English embedding model with 1024 dimensions",
			Dimensions:  1024,
			Provider:    "OpenRouter",
		},
		{
			Name:        "cohere/embed-multilingual-v3.0",
			DisplayName: "Cohere Embed Multilingual v3.0",
			Description: "Multilingual embedding model with 1024 dimensions",
			Dimensions:  1024,
			Provider:    "OpenRouter",
		},
	}
}

// EmbeddingModel represents information about an embedding model
type EmbeddingModel struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
	Dimensions  int    `json:"dimensions"`
	Provider    string `json:"provider"`
}

// ValidateConfig validates the embedding configuration
func (l *EmbeddingConfigLoader) ValidateConfig(config EmbeddingConfig) []string {
	var issues []string

	if config.APIKey == "" {
		issues = append(issues, "API key is required")
	}

	if config.ModelName == "" {
		issues = append(issues, "model name is required")
	} else {
		// Check if model is in supported list
		supported := false
		supportedModels := l.GetSupportedModels()
		for _, model := range supportedModels {
			if model.Name == config.ModelName {
				supported = true
				break
			}
		}
		if !supported {
			issues = append(issues, "model "+config.ModelName+" is not in the supported list")
		}
	}

	if config.Timeout <= 0 {
		issues = append(issues, "timeout must be positive")
	}

	if config.MaxRetries < 0 {
		issues = append(issues, "max retries cannot be negative")
	}

	if config.BatchSize <= 0 {
		issues = append(issues, "batch size must be positive")
	}

	if config.RetryDelay <= 0 {
		issues = append(issues, "retry delay must be positive")
	}

	return issues
}

// GetDefaultConfigForModel returns default configuration for a specific model
func (l *EmbeddingConfigLoader) GetDefaultConfigForModel(modelName string) EmbeddingConfig {
	config := DefaultEmbeddingConfig()
	config.ModelName = modelName

	// Adjust defaults based on model characteristics
	switch modelName {
	case "openai/text-embedding-3-large":
		// Larger model might need more timeout
		config.Timeout = 60 * time.Second
		config.BatchSize = 50 // Smaller batch size for larger model
	case "cohere/embed-english-v3.0", "cohere/embed-multilingual-v3.0":
		// Cohere models might have different optimal settings
		config.BatchSize = 96
	}

	return config
}
