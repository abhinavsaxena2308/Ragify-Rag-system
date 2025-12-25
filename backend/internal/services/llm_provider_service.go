package services

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"ragify-backend/internal/config"
)

// LLMProviderService manages multiple LLM providers and provides failover functionality
type LLMProviderService struct {
	providers    map[string]LLMClient
	primary      string
	fallback     []string
	config       config.LLMConfig
	mu           sync.RWMutex
	healthStatus map[string]bool
	lastCheck    map[string]time.Time
}

// NewLLMProviderService creates a new LLM provider service
func NewLLMProviderService(cfg config.LLMConfig) *LLMProviderService {
	service := &LLMProviderService{
		providers:    make(map[string]LLMClient),
		primary:      cfg.Provider,
		fallback:     []string{}, // Can be configured later
		config:       cfg,
		healthStatus: make(map[string]bool),
		lastCheck:    make(map[string]time.Time),
	}

	// Initialize providers based on configuration
	service.initializeProviders()

	return service
}

// initializeProviders sets up LLM providers based on configuration
func (s *LLMProviderService) initializeProviders() {
	// Initialize Ollama if configured
	if s.config.Endpoint != "" {
		ollamaService := NewOllamaService(s.config.Endpoint, s.config.Model)
		s.providers["ollama"] = ollamaService
		log.Printf("Initialized Ollama provider: %s", s.config.Model)
	}

	// Initialize OpenRouter if API key is provided
	if s.config.APIKey != "" {
		openRouterService := NewOpenRouterService(s.config.APIKey, s.config.Model)
		s.providers["openrouter"] = openRouterService
		log.Printf("Initialized OpenRouter provider: %s", s.config.Model)
	}

	// Set initial health status
	for provider := range s.providers {
		s.healthStatus[provider] = true // Assume healthy initially
		s.lastCheck[provider] = time.Time{}
	}
}

// ChatCompletion sends a chat completion request with failover support
func (s *LLMProviderService) ChatCompletion(ctx context.Context, request *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Try primary provider first
	if primaryClient, exists := s.providers[s.primary]; exists {
		response, err := s.tryProvider(ctx, primaryClient, s.primary, request)
		if err == nil {
			return response, nil
		}
		log.Printf("Primary provider %s failed: %v", s.primary, err)
	}

	// Try fallback providers
	for _, providerName := range s.fallback {
		if fallbackClient, exists := s.providers[providerName]; exists {
			response, err := s.tryProvider(ctx, fallbackClient, providerName, request)
			if err == nil {
				log.Printf("Fallback provider %s succeeded", providerName)
				return response, nil
			}
			log.Printf("Fallback provider %s failed: %v", providerName, err)
		}
	}

	// Try any other available providers
	for providerName, client := range s.providers {
		if providerName != s.primary && !containsSlice(s.fallback, providerName) {
			response, err := s.tryProvider(ctx, client, providerName, request)
			if err == nil {
				log.Printf("Alternative provider %s succeeded", providerName)
				return response, nil
			}
		}
	}

	return nil, &LLMError{
		Code:    ErrCodeInternalError,
		Message: "all LLM providers failed",
		Type:    "provider_error",
	}
}

// tryProvider attempts to use a specific provider with error handling
func (s *LLMProviderService) tryProvider(ctx context.Context, client LLMClient, providerName string, request *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	// Check if provider is healthy (if we haven't checked recently, assume healthy)
	if lastCheck, exists := s.lastCheck[providerName]; exists {
		if time.Since(lastCheck) < 5*time.Minute && !s.healthStatus[providerName] {
			return nil, &LLMError{
				Code:    ErrCodeInternalError,
				Message: fmt.Sprintf("provider %s is marked as unhealthy", providerName),
				Type:    "provider_error",
			}
		}
	}

	// Try the request
	response, err := client.ChatCompletion(ctx, request)
	if err != nil {
		// Mark provider as unhealthy
		s.mu.Lock()
		s.healthStatus[providerName] = false
		s.lastCheck[providerName] = time.Now()
		s.mu.Unlock()

		// Log the error
		log.Printf("Provider %s failed: %v", providerName, err)
		return nil, err
	}

	// Mark provider as healthy
	s.mu.Lock()
	s.healthStatus[providerName] = true
	s.lastCheck[providerName] = time.Now()
	s.mu.Unlock()

	return response, nil
}

// GenerateResponse generates a response from the LLM (backward compatibility)
func (s *LLMProviderService) GenerateResponse(ctx context.Context, prompt string, options map[string]interface{}) (string, error) {
	request := &ChatCompletionRequest{
		Messages: []ChatMessage{
			{Role: "user", Content: prompt},
		},
		Options: options,
	}

	response, err := s.ChatCompletion(ctx, request)
	if err != nil {
		return "", err
	}

	if len(response.Choices) == 0 {
		return "", &LLMError{
			Code:    ErrCodeInvalidRequest,
			Message: "no choices returned",
			Type:    "api_error",
		}
	}

	return response.Choices[0].Message.Content, nil
}

// IsHealthy checks if the LLM service is healthy
func (s *LLMProviderService) IsHealthy(ctx context.Context) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Check if we have any providers
	if len(s.providers) == 0 {
		return fmt.Errorf("no LLM providers configured")
	}

	// Check if primary provider is healthy
	if primaryClient, exists := s.providers[s.primary]; exists {
		if err := primaryClient.IsHealthy(ctx); err == nil {
			return nil
		}
	}

	// Check if any provider is healthy
	for providerName, client := range s.providers {
		if err := client.IsHealthy(ctx); err == nil {
			log.Printf("Provider %s is healthy", providerName)
			return nil
		}
	}

	return fmt.Errorf("no healthy LLM providers available")
}

// GetModelName returns the current model name
func (s *LLMProviderService) GetModelName() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if client, exists := s.providers[s.primary]; exists {
		return fmt.Sprintf("%s:%s", s.primary, client.GetModelName())
	}

	return "unknown"
}

// GetProvider returns the primary provider name
func (s *LLMProviderService) GetProvider() string {
	return s.primary
}

// SetPrimaryProvider sets the primary provider
func (s *LLMProviderService) SetPrimaryProvider(provider string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.providers[provider]; !exists {
		return fmt.Errorf("provider %s is not configured", provider)
	}

	s.primary = provider
	log.Printf("Primary provider changed to: %s", provider)
	return nil
}

// AddFallbackProvider adds a fallback provider
func (s *LLMProviderService) AddFallbackProvider(provider string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.providers[provider]; !exists {
		return fmt.Errorf("provider %s is not configured", provider)
	}

	if !containsSlice(s.fallback, provider) {
		s.fallback = append(s.fallback, provider)
		log.Printf("Added fallback provider: %s", provider)
	}

	return nil
}

// GetProviderStatus returns the status of all providers
func (s *LLMProviderService) GetProviderStatus() map[string]ProviderStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()

	status := make(map[string]ProviderStatus)
	for providerName, client := range s.providers {
		status[providerName] = ProviderStatus{
			Name:       providerName,
			Model:      client.GetModelName(),
			IsHealthy:  s.healthStatus[providerName],
			LastCheck:  s.lastCheck[providerName],
			IsPrimary:  providerName == s.primary,
			IsFallback: containsSlice(s.fallback, providerName),
		}
	}

	return status
}

// CheckAllProviders performs health checks on all providers
func (s *LLMProviderService) CheckAllProviders(ctx context.Context) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for providerName, client := range s.providers {
		err := client.IsHealthy(ctx)
		s.healthStatus[providerName] = err == nil
		s.lastCheck[providerName] = time.Now()

		if err != nil {
			log.Printf("Health check failed for provider %s: %v", providerName, err)
		} else {
			log.Printf("Health check passed for provider %s", providerName)
		}
	}
}

// ProviderStatus represents the status of a provider
type ProviderStatus struct {
	Name       string    `json:"name"`
	Model      string    `json:"model"`
	IsHealthy  bool      `json:"is_healthy"`
	LastCheck  time.Time `json:"last_check"`
	IsPrimary  bool      `json:"is_primary"`
	IsFallback bool      `json:"is_fallback"`
}

// containsSlice checks if a string slice contains a specific string
func containsSlice(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
