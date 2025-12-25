package main

import (
	"context"
	"fmt"
	"log"

	"ragify-backend/internal/config"
	"ragify-backend/internal/models"
	"ragify-backend/internal/services"
)

// Simple test without external dependencies
func main() {
	log.Println("Starting Simple RAG Test...")

	// Load configuration
	cfg := config.LoadConfig()
	log.Printf("Loaded configuration: RAG TopK=%d, Embedding Dimensions=%d",
		cfg.RAG.RetrievalTopK, cfg.RAG.EmbeddingDimensions)

	// Initialize mock services
	embeddingService := createMockEmbeddingService()
	vectorStore := createMockVectorStore(cfg.FAISS, cfg.RAG.EmbeddingDimensions)
	llmService := createMockLLMService()

	// Initialize RAG service
	ragService := services.NewRAGService(embeddingService, vectorStore, llmService, cfg.RAG)

	// Add sample data
	ctx := context.Background()
	if err := addSampleDataToVectorStore(ctx, vectorStore); err != nil {
		log.Fatalf("Failed to add sample data: %v", err)
	}

	// Test the RAG pipeline
	if err := testSimpleRAGPipeline(ctx, ragService); err != nil {
		log.Fatalf("RAG pipeline test failed: %v", err)
	}

	log.Println("Simple RAG Test completed successfully!")
}

func createMockEmbeddingService() *services.EmbeddingService {
	mockProvider := &MockEmbeddingProvider{
		dimensions: 1536,
	}
	cfg := config.DefaultEmbeddingConfig()
	return services.NewEmbeddingService(mockProvider, cfg)
}

func createMockVectorStore(cfg config.FAISSConfig, dimensions int) services.VectorStoreService {
	store, err := services.NewPureGoVectorStore(cfg, dimensions)
	if err != nil {
		log.Fatalf("Failed to create vector store: %v", err)
	}
	return store
}

func createMockLLMService() services.LLMService {
	return &MockLLMService{}
}

func addSampleDataToVectorStore(ctx context.Context, vectorStore services.VectorStoreService) error {
	log.Println("Adding sample data to vector store...")

	// Sample chunks
	sampleChunks := []*models.Chunk{
		{
			ID:         1,
			DocumentID: 1,
			Content:    "Machine learning is a subset of artificial intelligence that focuses on neural networks and deep learning algorithms.",
			ChunkIndex: 0,
			PageNumber: &[]int{1}[0],
		},
		{
			ID:         2,
			DocumentID: 1,
			Content:    "Natural language processing (NLP) is a branch of AI that helps computers understand and interpret human language.",
			ChunkIndex: 1,
			PageNumber: &[]int{1}[0],
		},
		{
			ID:         3,
			DocumentID: 2,
			Content:    "Computer vision is an AI field that trains computers to interpret and understand the visual world.",
			ChunkIndex: 0,
			PageNumber: &[]int{1}[0],
		},
	}

	// Generate embeddings and add to vector store
	for _, chunk := range sampleChunks {
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

	stats := vectorStore.GetStats()
	log.Printf("Vector store stats: %+v", stats)

	return nil
}

func testSimpleRAGPipeline(ctx context.Context, ragService *services.RAGService) error {
	log.Println("Testing RAG pipeline...")

	// Test question
	question := "What is machine learning?"
	log.Printf("Question: %s", question)

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
	fmt.Printf("\n=== RAG Pipeline Results ===\n")
	fmt.Printf("Question: %s\n", response.Question)
	fmt.Printf("Answer: %s\n", response.Answer)
	fmt.Printf("Sources: %d\n", len(response.Sources))
	fmt.Printf("Response Time: %v\n", response.ResponseTime)
	fmt.Printf("Timestamp: %v\n", response.Timestamp)

	for i, source := range response.Sources {
		fmt.Printf("  Source %d: Chunk %d (Score: %.3f)\n", i+1, source.ChunkID, source.Score)
		fmt.Printf("    Document ID: %d\n", source.DocumentID)
		fmt.Printf("    Content: %s\n", source.Content)
		if source.PageNumber != nil {
			fmt.Printf("    Page: %d\n", *source.PageNumber)
		}
	}

	// Test health check
	fmt.Printf("\n=== Health Check ===\n")
	if err := ragService.IsHealthy(ctx); err != nil {
		fmt.Printf("Health check failed: %v\n", err)
	} else {
		fmt.Printf("All services are healthy!\n")
	}

	// Test stats
	fmt.Printf("\n=== Statistics ===\n")
	stats := ragService.GetStats()
	fmt.Printf("Total Chunks: %d\n", stats.TotalChunks)
	fmt.Printf("Total Documents: %d\n", stats.TotalDocuments)

	return nil
}

// MockEmbeddingProvider is a mock implementation for demonstration
type MockEmbeddingProvider struct {
	dimensions int
}

func (m *MockEmbeddingProvider) GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	embedding := make([]float32, m.dimensions)

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

// MockLLMService is a mock implementation for demonstration
type MockLLMService struct{}

func (m *MockLLMService) GenerateResponse(ctx context.Context, prompt string, options map[string]interface{}) (string, error) {
	// Simple mock response based on prompt content
	if len(prompt) > 100 {
		return "Based on the provided context, machine learning is a subset of artificial intelligence that focuses on neural networks and deep learning algorithms. It enables computers to learn from data without being explicitly programmed.", nil
	}
	return "This is a mock response from the LLM service.", nil
}

func (m *MockLLMService) IsHealthy(ctx context.Context) error {
	return nil
}

func (m *MockLLMService) GetModelName() string {
	return "mock-llm-model"
}
