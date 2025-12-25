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
	log.Println("Starting LLM Integration Test...")

	// Load configuration
	cfg := config.LoadConfig()
	log.Printf("Loaded configuration: Primary=%s, Model=%s", cfg.LLM.Provider, cfg.LLM.Model)

	// Create LLM provider service
	llmService := services.NewLLMProviderService(cfg.LLM)

	// Add fallback providers if configured
	for _, fallback := range cfg.LLM.Fallback {
		if err := llmService.AddFallbackProvider(fallback); err != nil {
			log.Printf("Warning: Could not add fallback provider %s: %v", fallback, err)
		}
	}

	// Test basic functionality
	if err := testBasicFunctionality(llmService); err != nil {
		log.Fatalf("Basic functionality test failed: %v", err)
	}

	// Test chat completion with system prompt
	if err := testChatCompletion(llmService); err != nil {
		log.Fatalf("Chat completion test failed: %v", err)
	}

	// Test error handling and fallback
	if err := testErrorHandling(llmService); err != nil {
		log.Fatalf("Error handling test failed: %v", err)
	}

	// Test provider status and health checks
	if err := testProviderStatus(llmService); err != nil {
		log.Fatalf("Provider status test failed: %v", err)
	}

	log.Println("All LLM Integration Tests completed successfully!")
}

func testBasicFunctionality(llmService *services.LLMProviderService) error {
	log.Println("\n=== Testing Basic Functionality ===")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test simple prompt generation
	prompt := "What is machine learning? Explain in one sentence."
	response, err := llmService.GenerateResponse(ctx, prompt, map[string]interface{}{
		"temperature": 0.1,
		"max_tokens":  100,
	})

	if err != nil {
		return fmt.Errorf("failed to generate response: %w", err)
	}

	log.Printf("Prompt: %s", prompt)
	log.Printf("Response: %s", response)
	log.Printf("Provider: %s", llmService.GetProvider())
	log.Printf("Model: %s", llmService.GetModelName())

	return nil
}

func testChatCompletion(llmService *services.LLMProviderService) error {
	log.Println("\n=== Testing Chat Completion with System Prompt ===")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test with system and user messages
	request := &services.ChatCompletionRequest{
		Messages: []services.ChatMessage{
			{
				Role:    "system",
				Content: "You are a helpful AI assistant. Be concise and accurate.",
			},
			{
				Role:    "user",
				Content: "What are the main types of neural networks?",
			},
		},
		MaxTokens:   200,
		Temperature: 0.1,
		Timeout:     30 * time.Second,
	}

	response, err := llmService.ChatCompletion(ctx, request)
	if err != nil {
		return fmt.Errorf("failed to complete chat: %w", err)
	}

	log.Printf("Messages sent: %d", len(request.Messages))
	log.Printf("Response: %s", response.Choices[0].Message.Content)
	log.Printf("Tokens used: %d (prompt: %d, completion: %d)",
		response.Usage.TotalTokens, response.Usage.PromptTokens, response.Usage.CompletionTokens)
	log.Printf("Finish reason: %s", response.Choices[0].FinishReason)

	return nil
}

func testErrorHandling(llmService *services.LLMProviderService) error {
	log.Println("\n=== Testing Error Handling and Timeout ===")

	// Test with very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	request := &services.ChatCompletionRequest{
		Messages: []services.ChatMessage{
			{
				Role:    "user",
				Content: "Write a detailed essay about artificial intelligence.",
			},
		},
		Timeout: 1 * time.Millisecond, // Very short timeout
	}

	_, err := llmService.ChatCompletion(ctx, request)
	if err == nil {
		return fmt.Errorf("expected timeout error, but got none")
	}

	// Check if it's a timeout error
	if llmErr, ok := err.(*services.LLMError); ok {
		if llmErr.Code == services.ErrCodeTimeout {
			log.Printf("Successfully caught timeout error: %s", llmErr.Message)
		} else {
			log.Printf("Got LLM error (not timeout): %s - %s", llmErr.Code, llmErr.Message)
		}
	} else {
		log.Printf("Got general error: %v", err)
	}

	return nil
}

func testProviderStatus(llmService *services.LLMProviderService) error {
	log.Println("\n=== Testing Provider Status and Health Checks ===")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test overall health
	if err := llmService.IsHealthy(ctx); err != nil {
		log.Printf("Health check failed: %v", err)
	} else {
		log.Println("Overall health check passed")
	}

	// Get detailed provider status
	status := llmService.GetProviderStatus()
	log.Printf("Number of providers: %d", len(status))

	for providerName, providerStatus := range status {
		log.Printf("Provider: %s", providerName)
		log.Printf("  Model: %s", providerStatus.Model)
		log.Printf("  Healthy: %t", providerStatus.IsHealthy)
		log.Printf("  Primary: %t", providerStatus.IsPrimary)
		log.Printf("  Fallback: %t", providerStatus.IsFallback)
		if !providerStatus.LastCheck.IsZero() {
			log.Printf("  Last Check: %s ago", time.Since(providerStatus.LastCheck).Round(time.Second))
		}
	}

	// Run health checks on all providers
	log.Println("Running health checks on all providers...")
	llmService.CheckAllProviders(ctx)

	// Get updated status
	updatedStatus := llmService.GetProviderStatus()
	for providerName, providerStatus := range updatedStatus {
		log.Printf("Updated health status for %s: %t", providerName, providerStatus.IsHealthy)
	}

	return nil
}

// testProviderSwitching demonstrates switching between providers
func testProviderSwitching(llmService *services.LLMProviderService) error {
	log.Println("\n=== Testing Provider Switching ===")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get current provider
	currentProvider := llmService.GetProvider()
	log.Printf("Current primary provider: %s", currentProvider)

	// Try to switch to another provider if available
	status := llmService.GetProviderStatus()
	for providerName := range status {
		if providerName != currentProvider {
			log.Printf("Attempting to switch to provider: %s", providerName)

			if err := llmService.SetPrimaryProvider(providerName); err != nil {
				log.Printf("Could not switch to %s: %v", providerName, err)
				continue
			}

			// Test the new provider
			response, err := llmService.GenerateResponse(ctx, "Hello, who are you?", nil)
			if err != nil {
				log.Printf("Test failed with %s: %v", providerName, err)
				// Switch back
				llmService.SetPrimaryProvider(currentProvider)
				continue
			}

			log.Printf("Successfully switched to %s and got response: %s", providerName, response)

			// Switch back to original
			llmService.SetPrimaryProvider(currentProvider)
			break
		}
	}

	return nil
}
