package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"ragify-backend/internal/config"
	"ragify-backend/internal/services"
)

func main() {
	log.Println("Starting RAG with Enhanced LLM Integration...")

	// Load configuration
	cfg := config.LoadConfig()

	// Create LLM provider service
	llmService := services.NewLLMProviderService(cfg.LLM)

	// Add fallback providers if configured
	for _, fallback := range cfg.LLM.Fallback {
		if err := llmService.AddFallbackProvider(fallback); err != nil {
			log.Printf("Warning: Could not add fallback provider %s: %v", fallback, err)
		}
	}

	// Test the LLM service with system prompts
	if err := testLLMWithSystemPrompts(llmService); err != nil {
		log.Fatalf("LLM system prompt test failed: %v", err)
	}

	// Example of how to integrate with RAG
	if err := demonstrateRAGIntegration(llmService); err != nil {
		log.Fatalf("RAG integration demo failed: %v", err)
	}

	log.Println("Enhanced LLM Integration demo completed successfully!")
}

func testLLMWithSystemPrompts(llmService *services.LLMProviderService) error {
	log.Println("\n=== Testing LLM with System Prompts ===")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test 1: Strict system prompt for factual responses
	request := &services.ChatCompletionRequest{
		Messages: []services.ChatMessage{
			{
				Role:    "system",
				Content: "You are a factual AI assistant. Only provide information that is verifiable and accurate. If you're not sure about something, say 'I don't have enough information to answer this accurately.'",
			},
			{
				Role:    "user",
				Content: "What is the population of Tokyo?",
			},
		},
		MaxTokens:   150,
		Temperature: 0.1, // Low temperature for factual responses
	}

	response, err := llmService.ChatCompletion(ctx, request)
	if err != nil {
		return fmt.Errorf("factual system prompt test failed: %w", err)
	}

	log.Printf("Factual Response: %s", response.Choices[0].Message.Content)

	// Test 2: Creative system prompt
	request = &services.ChatCompletionRequest{
		Messages: []services.ChatMessage{
			{
				Role:    "system",
				Content: "You are a creative storyteller. Be imaginative and engaging in your responses. Use descriptive language and create vivid imagery.",
			},
			{
				Role:    "user",
				Content: "Tell me a short story about a robot discovering music for the first time.",
			},
		},
		MaxTokens:   200,
		Temperature: 0.8, // Higher temperature for creative responses
	}

	response, err = llmService.ChatCompletion(ctx, request)
	if err != nil {
		return fmt.Errorf("creative system prompt test failed: %w", err)
	}

	log.Printf("Creative Response: %s", response.Choices[0].Message.Content)

	return nil
}

func demonstrateRAGIntegration(llmService *services.LLMProviderService) error {
	log.Println("\n=== Demonstrating RAG Integration ===")

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Simulate RAG context
	ragContext := `[Source 1] Machine learning is a subset of artificial intelligence that focuses on the use of data and algorithms to imitate the way that humans learn, gradually improving its accuracy. (Document ID: 1, Page: 1)

[Source 2] Deep learning is a subset of machine learning that uses neural networks with multiple layers to analyze various factors in the data. These neural networks attempt to simulate the behavior of the human brain. (Document ID: 2, Page: 1)

[Source 3] Natural Language Processing (NLP) is a branch of artificial intelligence that helps computers understand, interpret and manipulate human language. (Document ID: 3, Page: 1)`

	question := "What is the relationship between machine learning and deep learning?"

	// Build RAG prompt with strict system prompt
	systemPrompt := `You are a helpful AI assistant that answers questions based ONLY on the provided context. Follow these rules strictly:
1. Use only information from the provided sources
2. If the context doesn't contain the answer, say "I don't have enough information to answer this question based on the provided context"
3. Cite your sources using [Source X] references
4. Be accurate and factual
5. Keep your answer concise but complete

Context:
` + ragContext + `

Question: ` + question

	// Create chat completion request with RAG context
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
		MaxTokens:   300,
		Temperature: 0.1, // Low temperature for factual RAG responses
		Timeout:     30 * time.Second,
	}

	response, err := llmService.ChatCompletion(ctx, request)
	if err != nil {
		return fmt.Errorf("RAG integration test failed: %w", err)
	}

	log.Printf("RAG Question: %s", question)
	log.Printf("RAG Answer: %s", response.Choices[0].Message.Content)
	log.Printf("Tokens Used: %d", response.Usage.TotalTokens)

	// Demonstrate error handling with invalid request
	log.Println("\n--- Testing Error Handling ---")

	// Test with timeout
	shortCtx, shortCancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer shortCancel()

	_, err = llmService.ChatCompletion(shortCtx, request)
	if err != nil {
		if llmErr, ok := err.(*services.LLMError); ok {
			log.Printf("Caught LLM error: %s - %s", llmErr.Code, llmErr.Message)
		} else {
			log.Printf("Caught general error: %v", err)
		}
	}

	return nil
}

// Example of how to create a mock LLM service for testing
func createMockLLMService() services.LLMClient {
	return &MockLLMService{}
}

// MockLLMService is a simple mock implementation for testing
type MockLLMService struct{}

func (m *MockLLMService) ChatCompletion(ctx context.Context, request *services.ChatCompletionRequest) (*services.ChatCompletionResponse, error) {
	// Simulate some processing time
	time.Sleep(100 * time.Millisecond)

	// Generate a mock response
	response := &services.ChatCompletionResponse{
		ID:      "mock-response-123",
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   "mock-model",
		Choices: []services.ChatChoice{
			{
				Index: 0,
				Message: services.ChatMessage{
					Role:    "assistant",
					Content: "This is a mock response from the mock LLM service. In a real implementation, this would be generated by an actual language model.",
				},
				FinishReason: "stop",
			},
		},
		Usage: services.UsageInfo{
			PromptTokens:     50,
			CompletionTokens: 25,
			TotalTokens:      75,
		},
	}

	return response, nil
}

func (m *MockLLMService) GenerateResponse(ctx context.Context, prompt string, options map[string]interface{}) (string, error) {
	response, err := m.ChatCompletion(ctx, &services.ChatCompletionRequest{
		Messages: []services.ChatMessage{
			{Role: "user", Content: prompt},
		},
		Options: options,
	})

	if err != nil {
		return "", err
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no choices returned")
	}

	return response.Choices[0].Message.Content, nil
}

func (m *MockLLMService) IsHealthy(ctx context.Context) error {
	return nil // Always healthy
}

func (m *MockLLMService) GetModelName() string {
	return "mock-model"
}

func (m *MockLLMService) GetProvider() string {
	return "mock"
}
