package main

import (
	"context"
	"fmt"
	"log"

	"ragify-backend/internal/config"
	"ragify-backend/internal/models"
	"ragify-backend/internal/services"
)

// Example demonstrating the RAG pipeline usage
func main() {
	log.Println("Starting RAG Pipeline Example...")

	// Load configuration
	cfg := config.LoadConfig()
	log.Printf("Loaded configuration: RAG TopK=%d, Embedding Dimensions=%d",
		cfg.RAG.RetrievalTopK, cfg.RAG.EmbeddingDimensions)

	// Initialize embedding service
	embeddingConfig := config.DefaultEmbeddingConfig()
	embeddingConfig.ModelName = "openai/text-embedding-3-small"

	// For this example, we'll use a mock embedding service
	// In production, you would use: embeddingService := services.NewOpenRouterEmbeddingService(embeddingConfig)
	embeddingService := createMockEmbeddingService(embeddingConfig)

	// Initialize vector store
	vectorStore, err := services.NewVectorStoreService(cfg.FAISS, cfg.RAG.EmbeddingDimensions)
	if err != nil {
		log.Fatalf("Failed to create vector store: %v", err)
	}
	defer vectorStore.Close()

	// Initialize LLM service
	var llmService services.LLMService
	switch cfg.LLM.Provider {
	case "ollama":
		llmService = services.NewOllamaService(cfg.LLM.Endpoint, cfg.LLM.Model)
	case "openrouter":
		llmService = services.NewOpenRouterService(cfg.LLM.APIKey, cfg.LLM.Model)
	default:
		llmService = services.NewOllamaService(cfg.LLM.Endpoint, cfg.LLM.Model)
	}

	// Initialize RAG service
	ragService := services.NewRAGService(embeddingService, vectorStore, llmService, cfg.RAG)

	// Add some sample data to the vector store
	if err := addSampleData(ragService, embeddingService, vectorStore); err != nil {
		log.Fatalf("Failed to add sample data: %v", err)
	}

	// Test the RAG pipeline
	if err := testRAGPipeline(ragService); err != nil {
		log.Fatalf("RAG pipeline test failed: %v", err)
	}

	log.Println("RAG Pipeline Example completed successfully!")
}

// addSampleData adds sample documents and chunks to the vector store
func addSampleData(ragService *services.RAGService, embeddingService *services.EmbeddingService, vectorStore services.VectorStoreService) error {
	log.Println("Adding sample data to vector store...")

	ctx := context.Background()

	// Sample chunks
	sampleChunks := []*models.Chunk{
		{
			ID:         1,
			DocumentID: 1,
			Content:    "Machine learning is a subset of artificial intelligence that focuses on neural networks and deep learning algorithms. It enables computers to learn from data without being explicitly programmed.",
			ChunkIndex: 0,
			PageNumber: &[]int{1}[0],
		},
		{
			ID:         2,
			DocumentID: 1,
			Content:    "Natural language processing (NLP) is a branch of AI that helps computers understand, interpret, and generate human language. Common applications include chatbots, translation, and sentiment analysis.",
			ChunkIndex: 1,
			PageNumber: &[]int{1}[0],
		},
		{
			ID:         3,
			DocumentID: 2,
			Content:    "Computer vision is an AI field that trains computers to interpret and understand the visual world. It's used in facial recognition, autonomous vehicles, and medical image analysis.",
			ChunkIndex: 0,
			PageNumber: &[]int{1}[0],
		},
		{
			ID:         4,
			DocumentID: 2,
			Content:    "Deep learning is a subset of machine learning that uses neural networks with multiple layers. It has revolutionized fields like image recognition, natural language processing, and game playing.",
			ChunkIndex: 1,
			PageNumber: &[]int{2}[0],
		},
	}

	// Generate embeddings and add to vector store
	for _, chunk := range sampleChunks {
		// Generate embedding (using mock service for this example)
		embedding := make([]float32, 1536)
		for i := range embedding {
			embedding[i] = float32(i+int(chunk.ID)) / float32(1536)
		}

		err := vectorStore.AddChunk(ctx, chunk, embedding)
		if err != nil {
			return fmt.Errorf("failed to add chunk %d: %w", chunk.ID, err)
		}

		log.Printf("Added chunk %d to vector store", chunk.ID)
	}

	// Save the vector store
	if err := vectorStore.Save(); err != nil {
		return fmt.Errorf("failed to save vector store: %w", err)
	}

	// Get statistics
	stats := vectorStore.GetStats()
	log.Printf("Vector store stats: %+v", stats)

	return nil
}

// testRAGPipeline tests the RAG pipeline with sample questions
func testRAGPipeline(ragService *services.RAGService) error {
	log.Println("Testing RAG pipeline...")

	ctx := context.Background()

	// Sample questions
	questions := []string{
		"What is machine learning?",
		"How does natural language processing work?",
		"What are the applications of computer vision?",
		"Explain deep learning and its uses.",
	}

	for i, question := range questions {
		log.Printf("\n--- Question %d: %s ---", i+1, question)

		// Create RAG request
		req := &models.RAGRequest{
			Question: question,
			TopK:     3,
		}

		// Process question
		response, err := ragService.ProcessQuestion(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to process question '%s': %w", question, err)
		}

		// Display results
		fmt.Printf("Question: %s\n", response.Question)
		fmt.Printf("Answer: %s\n", response.Answer)
		fmt.Printf("Sources: %d\n", len(response.Sources))
		fmt.Printf("Response Time: %v\n", response.ResponseTime)

		for j, source := range response.Sources {
			fmt.Printf("  Source %d: Chunk %d (Score: %.3f)\n", j+1, source.ChunkID, source.Score)
			fmt.Printf("    Content: %s\n", source.Content)
		}
		fmt.Println("---")
	}

	return nil
}

// createMockEmbeddingService creates a mock embedding service for demonstration
func createMockEmbeddingService(cfg config.EmbeddingConfig) *services.EmbeddingService {
	// This is a mock implementation for demonstration
	// In production, you would use a real embedding service
	mockProvider := &MockEmbeddingProvider{
		dimensions: 1536,
	}
	return services.NewEmbeddingService(mockProvider, cfg)
}

// MockEmbeddingProvider is a mock implementation for demonstration
type MockEmbeddingProvider struct {
	dimensions int
}

func (m *MockEmbeddingProvider) GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	// Generate deterministic mock embeddings based on text
	embedding := make([]float32, m.dimensions)

	// Simple hash-based embedding generation
	hash := 0
	for _, char := range text {
		hash = hash*31 + int(char)
	}

	for i := range embedding {
		embedding[i] = float32((hash+i)%1000) / float32(1000)
	}

	return embedding, nil
}

func (m *MockEmbeddingProvider) GenerateBatchEmbeddings(ctx context.Context, texts []string) ([]services.EmbeddingData, error) {
	results := make([]services.EmbeddingData, len(texts))
	for i, text := range texts {
		embedding, err := m.GenerateEmbedding(ctx, text)
		if err != nil {
			return nil, err
		}
		results[i] = services.EmbeddingData{
			Embedding: embedding,
			Index:     i,
		}
	}
	return results, nil
}

func (m *MockEmbeddingProvider) GetModelName() string {
	return "mock-embedding-model"
}

func (m *MockEmbeddingProvider) GetDimensions() int {
	return m.dimensions
}

func (m *MockEmbeddingProvider) IsHealthy(ctx context.Context) error {
	return nil
}
