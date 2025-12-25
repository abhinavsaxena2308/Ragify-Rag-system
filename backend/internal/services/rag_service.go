package services

import (
	"context"
	"fmt"
	"log"
	"strings"
	"text/template"
	"time"

	"ragify-backend/internal/config"
	"ragify-backend/internal/models"
)

// RAGService implements the core RAG pipeline
type RAGService struct {
	embeddingService *EmbeddingService
	vectorStore      VectorStoreService
	config           config.RAGConfig
	promptTemplate   *models.PromptTemplate
	llmService       LLMService
}

// LLMService interface for LLM providers
type LLMService interface {
	GenerateResponse(ctx context.Context, prompt string, options map[string]interface{}) (string, error)
	IsHealthy(ctx context.Context) error
	GetModelName() string
}

// NewRAGService creates a new RAG service
func NewRAGService(
	embeddingService *EmbeddingService,
	vectorStore VectorStoreService,
	llmService LLMService,
	cfg config.RAGConfig,
) *RAGService {
	return &RAGService{
		embeddingService: embeddingService,
		vectorStore:      vectorStore,
		config:           cfg,
		promptTemplate:   models.DefaultPromptTemplate(),
		llmService:       llmService,
	}
}

// ProcessQuestion processes a user question through the RAG pipeline
func (r *RAGService) ProcessQuestion(ctx context.Context, req *models.RAGRequest) (*models.RAGResponse, error) {
	startTime := time.Now()

	// Set default TopK if not provided
	topK := req.TopK
	if topK <= 0 {
		topK = r.config.RetrievalTopK
		if topK <= 0 {
			topK = 5
		}
	}

	log.Printf("Processing question: %s (TopK: %d)", req.Question, topK)

	// Step 1: Generate embedding for the question
	queryEmbedding, err := r.embeddingService.GenerateQueryEmbedding(ctx, req.Question)
	if err != nil {
		return nil, fmt.Errorf("failed to generate query embedding: %w", err)
	}
	log.Printf("Generated query embedding with %d dimensions", len(queryEmbedding))

	// Step 2: Retrieve relevant chunks from vector store
	retrievalStart := time.Now()
	searchResults, err := r.vectorStore.SearchChunks(ctx, queryEmbedding, topK)
	if err != nil {
		return nil, fmt.Errorf("failed to search vector store: %w", err)
	}
	retrievalTime := time.Since(retrievalStart)

	log.Printf("Retrieved %d chunks in %v", len(searchResults), retrievalTime)

	// Filter results by similarity threshold if configured
	if r.config.SimilarityThreshold > 0 {
		filteredResults := make([]SearchResult, 0)
		for _, result := range searchResults {
			if result.Score >= r.config.SimilarityThreshold {
				filteredResults = append(filteredResults, result)
			}
		}
		searchResults = filteredResults
		log.Printf("Filtered to %d chunks above threshold %.2f", len(searchResults), r.config.SimilarityThreshold)
	}

	if len(searchResults) == 0 {
		return &models.RAGResponse{
			Answer:       "I don't have enough information to answer this question based on the available context.",
			Sources:      []models.Source{},
			Question:     req.Question,
			ResponseTime: time.Since(startTime),
			Timestamp:    time.Now(),
			SessionID:    req.SessionID,
		}, nil
	}

	// Step 3: Build context from retrieved chunks
	context := r.buildContext(searchResults)

	// Step 4: Build LLM prompt
	prompt, err := r.buildPrompt(req.Question, context)
	if err != nil {
		return nil, fmt.Errorf("failed to build prompt: %w", err)
	}

	// Step 5: Call LLM
	llmStart := time.Now()
	answer, err := r.llmService.GenerateResponse(ctx, prompt, map[string]interface{}{
		"max_tokens":  r.config.MaxContextLength,
		"temperature": 0.1, // Lower temperature for more factual responses
	})
	llmTime := time.Since(llmStart)

	if err != nil {
		return nil, fmt.Errorf("failed to generate LLM response: %w", err)
	}

	log.Printf("Generated LLM response in %v", llmTime)

	// Step 6: Build response with sources
	sources := r.buildSources(searchResults)

	response := &models.RAGResponse{
		Answer:       answer,
		Sources:      sources,
		Question:     req.Question,
		ResponseTime: time.Since(startTime),
		Timestamp:    time.Now(),
		SessionID:    req.SessionID,
	}

	log.Printf("Completed RAG pipeline in %v (retrieval: %v, LLM: %v)",
		response.ResponseTime, retrievalTime, llmTime)

	return response, nil
}

// buildContext creates a context string from search results
func (r *RAGService) buildContext(searchResults []SearchResult) string {
	var contextParts []string

	for i, result := range searchResults {
		if result.Vector.Type == "chunk" && result.Vector.Chunk != nil {
			chunk := result.Vector.Chunk
			contextPart := fmt.Sprintf("[Source %d] %s", i+1, chunk.Content)

			// Add document reference if available
			if chunk.DocumentID != 0 {
				contextPart += fmt.Sprintf(" (Document ID: %d", chunk.DocumentID)
				if chunk.PageNumber != nil {
					contextPart += fmt.Sprintf(", Page: %d", *chunk.PageNumber)
				}
				contextPart += ")"
			}

			contextParts = append(contextParts, contextPart)
		}
	}

	return strings.Join(contextParts, "\n\n")
}

// buildPrompt creates the LLM prompt using the template
func (r *RAGService) buildPrompt(question, context string) (string, error) {
	// Combine system and user prompts
	fullPrompt := r.promptTemplate.SystemPrompt + "\n\n" + r.promptTemplate.UserPrompt

	// Create template and execute
	tmpl, err := template.New("prompt").Parse(fullPrompt)
	if err != nil {
		return "", fmt.Errorf("failed to parse prompt template: %w", err)
	}

	var promptBuilder strings.Builder
	err = tmpl.Execute(&promptBuilder, map[string]interface{}{
		"Context":  context,
		"Question": question,
	})
	if err != nil {
		return "", fmt.Errorf("failed to execute prompt template: %w", err)
	}

	return promptBuilder.String(), nil
}

// buildSources creates source objects from search results
func (r *RAGService) buildSources(searchResults []SearchResult) []models.Source {
	sources := make([]models.Source, 0, len(searchResults))

	for _, result := range searchResults {
		if result.Vector.Type == "chunk" && result.Vector.Chunk != nil {
			chunk := result.Vector.Chunk
			source := models.Source{
				ChunkID:      chunk.ID,
				DocumentID:   chunk.DocumentID,
				Content:      chunk.Content,
				Score:        result.Score,
				PageNumber:   chunk.PageNumber,
				DocumentName: fmt.Sprintf("Document_%d", chunk.DocumentID), // Will be updated with actual name
			}
			sources = append(sources, source)
		}
	}

	return sources
}
func (r *RAGService) GetStats() *models.RAGStats {
	// Get vector store stats
	vectorStats := r.vectorStore.GetStats()

	totalChunks := 0
	if typeCounts, ok := vectorStats["type_counts"].(map[string]int); ok {
		totalChunks = typeCounts["chunk"]
	}

	return &models.RAGStats{
		TotalDocuments:       0, // Will be updated from database
		TotalChunks:          totalChunks,
		AverageResponseTime:  0, // Will be calculated from usage metrics
		AverageRetrievalTime: 0, // Will be calculated from usage metrics
		CacheHitRate:         0, // Will be calculated from cache metrics
		// Add new fields here
		// ...
	}
}

// IsHealthy checks if the RAG service is healthy
func (r *RAGService) IsHealthy(ctx context.Context) error {
	// Check embedding service
	if err := r.embeddingService.IsHealthy(ctx); err != nil {
		return fmt.Errorf("embedding service unhealthy: %w", err)
	}

	// Check vector store
	if err := r.vectorStore.IsHealthy(ctx); err != nil {
		return fmt.Errorf("vector store unhealthy: %w", err)
	}

	// Check LLM service
	if err := r.llmService.IsHealthy(ctx); err != nil {
		return fmt.Errorf("LLM service unhealthy: %w", err)
	}

	return nil
}

// UpdatePromptTemplate updates the prompt template
func (r *RAGService) UpdatePromptTemplate(template *models.PromptTemplate) {
	r.promptTemplate = template
	log.Printf("Updated prompt template")
}

// GetPromptTemplate returns the current prompt template
func (r *RAGService) GetPromptTemplate() *models.PromptTemplate {
	return r.promptTemplate
}

// SetLLMService updates the LLM service
func (r *RAGService) SetLLMService(llmService LLMService) {
	r.llmService = llmService
	log.Printf("Updated LLM service to: %s", llmService.GetModelName())
}

// GetConfig returns the current RAG configuration
func (r *RAGService) GetConfig() config.RAGConfig {
	return r.config
}
