package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// OllamaService implements LLMService for Ollama
type OllamaService struct {
	endpoint string
	model    string
	client   *http.Client
}

// OllamaRequest represents a request to Ollama API
type OllamaRequest struct {
	Model   string                 `json:"model"`
	Prompt  string                 `json:"prompt"`
	Options map[string]interface{} `json:"options,omitempty"`
	Stream  bool                   `json:"stream"`
}

// OllamaResponse represents a response from Ollama API
type OllamaResponse struct {
	Model           string `json:"model"`
	Response        string `json:"response"`
	Done            bool   `json:"done"`
	TotalDur        string `json:"total_duration,omitempty"`
	LoadDur         string `json:"load_duration,omitempty"`
	PromptEvalCount int    `json:"prompt_eval_count,omitempty"`
	PromptEvalDur   string `json:"prompt_eval_duration,omitempty"`
	EvalCount       int    `json:"eval_count,omitempty"`
	EvalDur         string `json:"eval_duration,omitempty"`
}

// NewOllamaService creates a new Ollama service
func NewOllamaService(endpoint, model string) *OllamaService {
	return &OllamaService{
		endpoint: endpoint,
		model:    model,
		client: &http.Client{
			Timeout: 120 * time.Second, // Longer timeout for LLM responses
		},
	}
}

// GenerateResponse generates a response from the LLM
func (o *OllamaService) GenerateResponse(ctx context.Context, prompt string, options map[string]interface{}) (string, error) {
	request := OllamaRequest{
		Model:   o.model,
		Prompt:  prompt,
		Options: options,
		Stream:  false,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/api/generate", o.endpoint)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	log.Printf("Sending request to Ollama: model=%s, prompt_length=%d", o.model, len(prompt))

	resp, err := o.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Ollama returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var ollamaResp OllamaResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	log.Printf("Received response from Ollama: response_length=%d, done=%v", len(ollamaResp.Response), ollamaResp.Done)

	return ollamaResp.Response, nil
}

// IsHealthy checks if the Ollama service is healthy
func (o *OllamaService) IsHealthy(ctx context.Context) error {
	url := fmt.Sprintf("%s/api/tags", o.endpoint)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create health check request: %w", err)
	}

	resp, err := o.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform health check: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Ollama health check failed with status %d", resp.StatusCode)
	}

	return nil
}

// GetModelName returns the current model name
func (o *OllamaService) GetModelName() string {
	return o.model
}

// OpenRouterService implements LLMService for OpenRouter
type OpenRouterService struct {
	apiKey   string
	model    string
	client   *http.Client
	endpoint string
}

// OpenRouterRequest represents a request to OpenRouter API
type OpenRouterRequest struct {
	Model    string                 `json:"model"`
	Messages []OpenRouterMessage    `json:"messages"`
	Options  map[string]interface{} `json:"options,omitempty"`
}

// OpenRouterMessage represents a message in OpenRouter API
type OpenRouterMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenRouterResponse represents a response from OpenRouter API
type OpenRouterResponse struct {
	ID      string `json:"id"`
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		Index        int    `json:"index"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// NewOpenRouterService creates a new OpenRouter service
func NewOpenRouterService(apiKey, model string) *OpenRouterService {
	return &OpenRouterService{
		apiKey:   apiKey,
		model:    model,
		endpoint: "https://openrouter.ai/api/v1/chat/completions",
		client: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

// GenerateResponse generates a response from the LLM
func (o *OpenRouterService) GenerateResponse(ctx context.Context, prompt string, options map[string]interface{}) (string, error) {
	// Split prompt into system and user messages
	messages := []OpenRouterMessage{
		{Role: "user", Content: prompt},
	}

	request := OpenRouterRequest{
		Model:    o.model,
		Messages: messages,
		Options:  options,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", o.endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+o.apiKey)
	req.Header.Set("HTTP-Referer", "https://github.com/abhinavsaxena2308/Ragify-Rag-system")
	req.Header.Set("X-Title", "Ragify RAG System")

	log.Printf("Sending request to OpenRouter: model=%s, prompt_length=%d", o.model, len(prompt))

	resp, err := o.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("OpenRouter returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var openRouterResp OpenRouterResponse
	if err := json.Unmarshal(body, &openRouterResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(openRouterResp.Choices) == 0 {
		return "", fmt.Errorf("no choices returned from OpenRouter")
	}

	response := openRouterResp.Choices[0].Message.Content
	log.Printf("Received response from OpenRouter: response_length=%d", len(response))

	return response, nil
}

// IsHealthy checks if the OpenRouter service is healthy
func (o *OpenRouterService) IsHealthy(ctx context.Context) error {
	if o.apiKey == "" {
		return fmt.Errorf("OpenRouter API key is not configured")
	}

	// Simple health check - try to get available models
	req, err := http.NewRequestWithContext(ctx, "GET", "https://openrouter.ai/api/v1/models", nil)
	if err != nil {
		return fmt.Errorf("failed to create health check request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+o.apiKey)

	resp, err := o.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform health check: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("OpenRouter health check failed with status %d", resp.StatusCode)
	}

	return nil
}

// GetModelName returns the current model name
func (o *OpenRouterService) GetModelName() string {
	return o.model
}
