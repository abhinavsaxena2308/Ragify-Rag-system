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

// OllamaService implements LLMClient for Ollama
type OllamaService struct {
	endpoint string
	model    string
	client   *http.Client
}

// OllamaRequest represents a request to Ollama API
type OllamaRequest struct {
	Model    string                 `json:"model"`
	Messages []OllamaMessage        `json:"messages,omitempty"`
	Prompt   string                 `json:"prompt,omitempty"`
	Options  map[string]interface{} `json:"options,omitempty"`
	Stream   bool                   `json:"stream"`
}

// OllamaMessage represents a message in Ollama API
type OllamaMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
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
			Transport: &http.Transport{
				MaxIdleConns:       10,
				IdleConnTimeout:    30 * time.Second,
				DisableCompression: false,
			},
		},
	}
}

// ChatCompletion sends a chat completion request with proper system prompt support
func (o *OllamaService) ChatCompletion(ctx context.Context, request *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	// Set default timeout if not provided
	timeout := request.Timeout
	if timeout == 0 {
		timeout = 120 * time.Second
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Build request
	reqModel := request.Model
	if reqModel == "" {
		reqModel = o.model
	}

	// Convert ChatMessages to Ollama format
	messages := make([]OllamaMessage, len(request.Messages))
	for i, msg := range request.Messages {
		messages[i] = OllamaMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	// Ollama chat API request
	ollamaReq := OllamaRequest{
		Model:    reqModel,
		Messages: messages,
		Options:  request.Options,
		Stream:   request.Stream,
	}

	// Add max_tokens and temperature to options if not already present
	if ollamaReq.Options == nil {
		ollamaReq.Options = make(map[string]interface{})
	}
	if request.MaxTokens > 0 {
		ollamaReq.Options["num_predict"] = request.MaxTokens
	}
	if request.Temperature > 0 {
		ollamaReq.Options["temperature"] = request.Temperature
	}

	jsonData, err := json.Marshal(ollamaReq)
	if err != nil {
		return nil, &LLMError{
			Code:    ErrCodeInvalidRequest,
			Message: fmt.Sprintf("failed to marshal request: %v", err),
			Type:    "client_error",
		}
	}

	url := fmt.Sprintf("%s/api/chat", o.endpoint)

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, &LLMError{
			Code:    ErrCodeNetworkError,
			Message: fmt.Sprintf("failed to create request: %v", err),
			Type:    "network_error",
		}
	}

	httpReq.Header.Set("Content-Type", "application/json")

	log.Printf("Sending chat completion to Ollama: model=%s, messages=%d, max_tokens=%d",
		reqModel, len(messages), request.MaxTokens)

	resp, err := o.client.Do(httpReq)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, &LLMError{
				Code:    ErrCodeTimeout,
				Message: "request timeout exceeded",
				Type:    "timeout_error",
			}
		}
		return nil, &LLMError{
			Code:    ErrCodeNetworkError,
			Message: fmt.Sprintf("failed to send request: %v", err),
			Type:    "network_error",
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &LLMError{
			Code:    ErrCodeNetworkError,
			Message: fmt.Sprintf("failed to read response: %v", err),
			Type:    "network_error",
		}
	}

	if resp.StatusCode != http.StatusOK {
		return o.handleHTTPError(resp.StatusCode, body)
	}

	var ollamaResp OllamaChatResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return nil, &LLMError{
			Code:    ErrCodeInternalError,
			Message: fmt.Sprintf("failed to unmarshal response: %v", err),
			Type:    "parse_error",
		}
	}

	// Convert to standard ChatCompletionResponse
	response := &ChatCompletionResponse{
		ID:      fmt.Sprintf("ollama-%d", time.Now().Unix()),
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   reqModel,
		Choices: []ChatChoice{
			{
				Index: 0,
				Message: ChatMessage{
					Role:    "assistant",
					Content: ollamaResp.Message.Content,
				},
				FinishReason: "stop",
			},
		},
		Usage: UsageInfo{
			PromptTokens:     ollamaResp.PromptEvalCount,
			CompletionTokens: ollamaResp.EvalCount,
			TotalTokens:      ollamaResp.PromptEvalCount + ollamaResp.EvalCount,
		},
	}

	log.Printf("Received chat completion from Ollama: response_length=%d, tokens=%d",
		len(ollamaResp.Message.Content), response.Usage.TotalTokens)

	return response, nil
}

// GenerateResponse generates a response from the LLM (backward compatibility)
func (o *OllamaService) GenerateResponse(ctx context.Context, prompt string, options map[string]interface{}) (string, error) {
	request := &ChatCompletionRequest{
		Messages: []ChatMessage{
			{Role: "user", Content: prompt},
		},
		Options: options,
	}

	response, err := o.ChatCompletion(ctx, request)
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

// OllamaChatResponse represents a chat response from Ollama API
type OllamaChatResponse struct {
	Model           string        `json:"model"`
	CreatedAt       time.Time     `json:"created_at"`
	Message         OllamaMessage `json:"message"`
	Done            bool          `json:"done"`
	TotalDur        time.Duration `json:"total_duration,omitempty"`
	LoadDur         time.Duration `json:"load_duration,omitempty"`
	PromptEvalCount int           `json:"prompt_eval_count,omitempty"`
	PromptEvalDur   time.Duration `json:"prompt_eval_duration,omitempty"`
	EvalCount       int           `json:"eval_count,omitempty"`
	EvalDur         time.Duration `json:"eval_duration,omitempty"`
}

// handleHTTPError handles HTTP error responses from Ollama
func (o *OllamaService) handleHTTPError(statusCode int, body []byte) (*ChatCompletionResponse, error) {
	bodyStr := string(body)

	switch statusCode {
	case 404:
		return nil, &LLMError{
			Code:    ErrCodeModelNotFound,
			Message: "model not found",
			Type:    "invalid_request_error",
		}
	case 500:
		return nil, &LLMError{
			Code:    ErrCodeInternalError,
			Message: fmt.Sprintf("Ollama internal error: %s", bodyStr),
			Type:    "server_error",
		}
	default:
		return nil, &LLMError{
			Code:    ErrCodeInternalError,
			Message: fmt.Sprintf("Ollama returned status %d: %s", statusCode, bodyStr),
			Type:    "api_error",
		}
	}
}

// GetProvider returns the provider name
func (o *OllamaService) GetProvider() string {
	return "ollama"
}

// GetModelName returns the current model name
func (o *OllamaService) GetModelName() string {
	return o.model
}

// OpenRouterService implements LLMClient for OpenRouter
type OpenRouterService struct {
	apiKey   string
	model    string
	client   *http.Client
	endpoint string
}

// OpenRouterRequest represents a request to OpenRouter API
type OpenRouterRequest struct {
	Model       string                 `json:"model"`
	Messages    []OpenRouterMessage    `json:"messages"`
	MaxTokens   int                    `json:"max_tokens,omitempty"`
	Temperature float64                `json:"temperature,omitempty"`
	Options     map[string]interface{} `json:"options,omitempty"`
	Stream      bool                   `json:"stream,omitempty"`
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
			Transport: &http.Transport{
				MaxIdleConns:       10,
				IdleConnTimeout:    30 * time.Second,
				DisableCompression: false,
			},
		},
	}
}

// ChatCompletion sends a chat completion request with proper system prompt support
func (o *OpenRouterService) ChatCompletion(ctx context.Context, request *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	// Set default timeout if not provided
	timeout := request.Timeout
	if timeout == 0 {
		timeout = 120 * time.Second
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Convert ChatMessages to OpenRouter format
	messages := make([]OpenRouterMessage, len(request.Messages))
	for i, msg := range request.Messages {
		messages[i] = OpenRouterMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	// Build request
	reqModel := request.Model
	if reqModel == "" {
		reqModel = o.model
	}

	openRouterReq := OpenRouterRequest{
		Model:       reqModel,
		Messages:    messages,
		MaxTokens:   request.MaxTokens,
		Temperature: request.Temperature,
		Options:     request.Options,
		Stream:      request.Stream,
	}

	jsonData, err := json.Marshal(openRouterReq)
	if err != nil {
		return nil, &LLMError{
			Code:    ErrCodeInvalidRequest,
			Message: fmt.Sprintf("failed to marshal request: %v", err),
			Type:    "client_error",
		}
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", o.endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, &LLMError{
			Code:    ErrCodeNetworkError,
			Message: fmt.Sprintf("failed to create request: %v", err),
			Type:    "network_error",
		}
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+o.apiKey)
	httpReq.Header.Set("HTTP-Referer", "https://github.com/abhinavsaxena2308/Ragify-Rag-system")
	httpReq.Header.Set("X-Title", "Ragify RAG System")

	log.Printf("Sending chat completion to OpenRouter: model=%s, messages=%d, max_tokens=%d",
		reqModel, len(messages), request.MaxTokens)

	resp, err := o.client.Do(httpReq)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, &LLMError{
				Code:    ErrCodeTimeout,
				Message: "request timeout exceeded",
				Type:    "timeout_error",
			}
		}
		return nil, &LLMError{
			Code:    ErrCodeNetworkError,
			Message: fmt.Sprintf("failed to send request: %v", err),
			Type:    "network_error",
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &LLMError{
			Code:    ErrCodeNetworkError,
			Message: fmt.Sprintf("failed to read response: %v", err),
			Type:    "network_error",
		}
	}

	if resp.StatusCode != http.StatusOK {
		return o.handleHTTPError(resp.StatusCode, body)
	}

	var openRouterResp OpenRouterResponse
	if err := json.Unmarshal(body, &openRouterResp); err != nil {
		return nil, &LLMError{
			Code:    ErrCodeInternalError,
			Message: fmt.Sprintf("failed to unmarshal response: %v", err),
			Type:    "parse_error",
		}
	}

	if len(openRouterResp.Choices) == 0 {
		return nil, &LLMError{
			Code:    ErrCodeInvalidRequest,
			Message: "no choices returned from OpenRouter",
			Type:    "api_error",
		}
	}

	// Convert to standard ChatCompletionResponse
	response := &ChatCompletionResponse{
		ID:      openRouterResp.ID,
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   reqModel,
		Choices: make([]ChatChoice, len(openRouterResp.Choices)),
		Usage: UsageInfo{
			PromptTokens:     openRouterResp.Usage.PromptTokens,
			CompletionTokens: openRouterResp.Usage.CompletionTokens,
			TotalTokens:      openRouterResp.Usage.TotalTokens,
		},
	}

	for i, choice := range openRouterResp.Choices {
		response.Choices[i] = ChatChoice{
			Index: choice.Index,
			Message: ChatMessage{
				Role:    choice.Message.Role,
				Content: choice.Message.Content,
			},
			FinishReason: choice.FinishReason,
		}
	}

	log.Printf("Received chat completion from OpenRouter: response_length=%d, tokens=%d",
		len(response.Choices[0].Message.Content), response.Usage.TotalTokens)

	return response, nil
}

// GenerateResponse generates a response from the LLM (backward compatibility)
func (o *OpenRouterService) GenerateResponse(ctx context.Context, prompt string, options map[string]interface{}) (string, error) {
	request := &ChatCompletionRequest{
		Messages: []ChatMessage{
			{Role: "user", Content: prompt},
		},
		Options: options,
	}

	response, err := o.ChatCompletion(ctx, request)
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

// handleHTTPError handles HTTP error responses from OpenRouter
func (o *OpenRouterService) handleHTTPError(statusCode int, body []byte) (*ChatCompletionResponse, error) {
	bodyStr := string(body)

	switch statusCode {
	case 401:
		return nil, &LLMError{
			Code:    ErrCodeInvalidAPIKey,
			Message: "invalid API key",
			Type:    "authentication_error",
		}
	case 429:
		return nil, &LLMError{
			Code:    ErrCodeRateLimit,
			Message: "rate limit exceeded",
			Type:    "rate_limit_error",
		}
	case 404:
		return nil, &LLMError{
			Code:    ErrCodeModelNotFound,
			Message: "model not found",
			Type:    "invalid_request_error",
		}
	case 402:
		return nil, &LLMError{
			Code:    ErrCodeInsufficientQuota,
			Message: "insufficient quota",
			Type:    "insufficient_quota_error",
		}
	default:
		return nil, &LLMError{
			Code:    ErrCodeInternalError,
			Message: fmt.Sprintf("OpenRouter returned status %d: %s", statusCode, bodyStr),
			Type:    "api_error",
		}
	}
}

// GetProvider returns the provider name
func (o *OpenRouterService) GetProvider() string {
	return "openrouter"
}

// GetModelName returns the current model name
func (o *OpenRouterService) GetModelName() string {
	return o.model
}
