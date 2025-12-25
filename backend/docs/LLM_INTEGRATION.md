# LLM Integration for Ragify RAG System

This document describes the enhanced LLM integration that supports multiple providers with chat completion API, strict system prompts, timeout handling, and error management.

## Supported Providers

### 1. OpenRouter
- **Description**: Cloud-based service providing access to multiple LLM models
- **API Key**: Required (set via `LLM_API_KEY` environment variable)
- **Models**: Supports various models (e.g., `anthropic/claude-3-haiku`, `openai/gpt-4o-mini`)
- **Endpoint**: `https://openrouter.ai/api/v1/chat/completions`

### 2. Ollama
- **Description**: Local LLM server for running models locally
- **API Key**: Not required
- **Models**: Any model available in your Ollama instance (e.g., `llama2`, `mistral`)
- **Endpoint**: `http://localhost:11434` (default)

## Configuration

### Environment Variables

```bash
# Primary Provider Configuration
LLM_PROVIDER=openrouter              # or "ollama"
LLM_MODEL=anthropic/claude-3-haiku   # model name
LLM_API_KEY=your_openrouter_api_key  # required for OpenRouter
LLM_ENDPOINT=http://localhost:11434  # for Ollama

# Advanced Configuration
LLM_FALLBACK=ollama                   # comma-separated fallback providers
LLM_TIMEOUT=120                       # timeout in seconds
LLM_MAX_TOKENS=4000                   # maximum tokens per response
LLM_TEMPERATURE=0.1                   # temperature (0.0-1.0)
LLM_RETRY_COUNT=3                     # number of retries
LLM_HEALTH_CHECK=true                 # enable health checks
```

### Example Configuration Files

#### `.env` for OpenRouter
```env
LLM_PROVIDER=openrouter
LLM_MODEL=anthropic/claude-3-haiku
LLM_API_KEY=sk-or-v1-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
LLM_FALLBACK=ollama
LLM_TIMEOUT=120
LLM_MAX_TOKENS=4000
LLM_TEMPERATURE=0.1
LLM_HEALTH_CHECK=true
```

#### `.env` for Ollama
```env
LLM_PROVIDER=ollama
LLM_MODEL=llama2
LLM_ENDPOINT=http://localhost:11434
LLM_TIMEOUT=120
LLM_MAX_TOKENS=4000
LLM_TEMPERATURE=0.1
LLM_HEALTH_CHECK=true
```

## Usage Examples

### Basic Usage

```go
package main

import (
    "context"
    "log"
    "ragify-backend/internal/config"
    "ragify-backend/internal/services"
)

func main() {
    // Load configuration
    cfg := config.LoadConfig()
    
    // Create LLM provider service
    llmService := services.NewLLMProviderService(cfg.LLM)
    
    // Add fallback providers
    for _, fallback := range cfg.LLM.Fallback {
        llmService.AddFallbackProvider(fallback)
    }
    
    // Generate response
    ctx := context.Background()
    response, err := llmService.GenerateResponse(ctx, "What is AI?", nil)
    if err != nil {
        log.Fatal(err)
    }
    
    log.Println(response)
}
```

### Chat Completion with System Prompts

```go
// Create chat completion request with system prompt
request := &services.ChatCompletionRequest{
    Messages: []services.ChatMessage{
        {
            Role:    "system",
            Content: "You are a helpful AI assistant. Be accurate and concise.",
        },
        {
            Role:    "user", 
            Content: "Explain quantum computing in simple terms.",
        },
    },
    MaxTokens:   300,
    Temperature: 0.1,
    Timeout:     30 * time.Second,
}

response, err := llmService.ChatCompletion(ctx, request)
if err != nil {
    log.Fatal(err)
}

log.Println(response.Choices[0].Message.Content)
log.Printf("Tokens used: %d", response.Usage.TotalTokens)
```

### RAG Integration Example

```go
// Build RAG prompt with context
systemPrompt := `You are a helpful AI assistant that answers questions based ONLY on the provided context.

Context:
` + ragContext + `

Question: ` + question

request := &services.ChatCompletionRequest{
    Messages: []services.ChatMessage{
        {
            Role:    "system",
            Content: systemPrompt,
        },
        {
            Role:    "user",
            Content: question,
        },
    },
    MaxTokens:   500,
    Temperature: 0.1, // Low temperature for factual responses
}

response, err := llmService.ChatCompletion(ctx, request)
```

## Error Handling

The LLM service provides structured error handling with specific error codes:

```go
response, err := llmService.ChatCompletion(ctx, request)
if err != nil {
    if llmErr, ok := err.(*services.LLMError); ok {
        switch llmErr.Code {
        case services.ErrCodeTimeout:
            log.Println("Request timed out")
        case services.ErrCodeRateLimit:
            log.Println("Rate limit exceeded")
        case services.ErrCodeInvalidAPIKey:
            log.Println("Invalid API key")
        case services.ErrCodeModelNotFound:
            log.Println("Model not found")
        default:
            log.Printf("LLM error: %s", llmErr.Message)
        }
    }
}
```

## Provider Management

### Checking Provider Status

```go
// Get status of all providers
status := llmService.GetProviderStatus()
for provider, providerStatus := range status {
    log.Printf("Provider: %s", provider)
    log.Printf("  Model: %s", providerStatus.Model)
    log.Printf("  Healthy: %t", providerStatus.IsHealthy)
    log.Printf("  Primary: %t", providerStatus.IsPrimary)
}

// Run health checks on all providers
llmService.CheckAllProviders(ctx)
```

### Switching Providers

```go
// Switch primary provider
err := llmService.SetPrimaryProvider("openrouter")
if err != nil {
    log.Printf("Failed to switch provider: %v", err)
}

// Add fallback provider
err = llmService.AddFallbackProvider("ollama")
if err != nil {
    log.Printf("Failed to add fallback: %v", err)
}
```

## Testing

Run the integration tests:

```bash
# Test LLM integration
go run examples/llm_integration_test.go

# Test RAG with LLM integration
go run examples/rag_llm_demo.go
```

## Features

### 1. Multi-Provider Support
- Primary provider with automatic fallback
- Health checking and provider status monitoring
- Dynamic provider switching

### 2. Chat Completion API
- Full support for system and user messages
- Token usage tracking
- Configurable parameters (temperature, max_tokens, etc.)

### 3. Error Handling
- Structured error types with specific codes
- Timeout and retry mechanisms
- Graceful degradation with fallback providers

### 4. Performance Optimization
- Connection pooling and keep-alive
- Configurable timeouts
- Health check caching

### 5. RAG Integration
- Strict system prompt enforcement
- Context-aware responses
- Source citation support

## Best Practices

1. **Use System Prompts**: Always provide clear system prompts for consistent behavior
2. **Set Appropriate Temperature**: Use low temperature (0.1-0.3) for factual responses, higher (0.7-0.9) for creative responses
3. **Configure Timeouts**: Set reasonable timeouts based on model and network conditions
4. **Monitor Health**: Regularly check provider health and handle failures gracefully
5. **Use Fallbacks**: Configure fallback providers for high availability
6. **Track Usage**: Monitor token usage for cost management

## Troubleshooting

### Common Issues

1. **API Key Errors**: Verify `LLM_API_KEY` is set correctly for OpenRouter
2. **Connection Errors**: Check `LLM_ENDPOINT` and network connectivity
3. **Model Not Found**: Verify model name is correct and available
4. **Timeout Errors**: Increase `LLM_TIMEOUT` or check network latency
5. **Rate Limiting**: Implement retry logic or upgrade API plan

### Debug Mode

Enable debug logging by setting log level:

```go
log.SetFlags(log.LstdFlags | log.Lshortfile)
```

### Health Check Issues

If health checks fail:
1. Verify provider endpoints are accessible
2. Check API key validity
3. Ensure model is available
4. Review network/firewall settings
