package vectorstore

import (
	"context"
	"fmt"
	"log"

	"ragify-backend/internal/config"
	"ragify-backend/internal/models"
	"ragify-backend/internal/services"
)

// Simple integration example showing how to use VectorStoreService with EmbeddingService
func RunVectorStoreExample() {
	// Initialize configuration
	cfg := config.LoadConfig()

	// Initialize embedding service
	embeddingConfig := config.DefaultEmbeddingConfig()
	embeddingConfig.ModelName = "openai/text-embedding-3-small"

	// Note: You'll need to set up your OpenRouter API key in environment variables
	// embeddingService := services.NewOpenRouterEmbeddingService(embeddingConfig)

	// Initialize vector store
	vectorStore, err := services.NewVectorStoreService(cfg.FAISS, 1536) // 1536 dimensions for text-embedding-3-small
	if err != nil {
		log.Fatalf("Failed to create vector store: %v", err)
	}
	defer vectorStore.Close()

	ctx := context.Background()

	// Example: Add some sample chunks
	sampleChunks := []*models.Chunk{
		{
			ID:      1,
			Content: "Machine learning is a subset of artificial intelligence that focuses on neural networks and deep learning algorithms.",
		},
		{
			ID:      2,
			Content: "Natural language processing (NLP) is a branch of AI that helps computers understand and interpret human language.",
		},
		{
			ID:      3,
			Content: "Computer vision is an AI field that trains computers to interpret and understand the visual world.",
		},
	}

	// Generate embeddings and add to vector store
	for _, chunk := range sampleChunks {
		// In a real implementation, you would generate embeddings using the embedding service
		// embedding, err := embeddingService.GenerateChunkEmbedding(ctx, chunk)
		// if err != nil {
		//     log.Printf("Failed to generate embedding for chunk %d: %v", chunk.ID, err)
		//     continue
		// }

		// For this example, we'll use dummy embeddings
		embedding := make([]float32, 1536)
		for i := range embedding {
			embedding[i] = float32(i+int(chunk.ID)) / float32(1536)
		}

		err := vectorStore.AddChunk(ctx, chunk, embedding)
		if err != nil {
			log.Printf("Failed to add chunk %d to vector store: %v", chunk.ID, err)
			continue
		}

		fmt.Printf("Added chunk %d to vector store\n", chunk.ID)
	}

	// Save the vector store
	err = vectorStore.Save()
	if err != nil {
		log.Fatalf("Failed to save vector store: %v", err)
	}

	// Get statistics
	stats := vectorStore.GetStats()
	fmt.Printf("Vector store stats: %+v\n", stats)

	// Example search
	queryEmbedding := make([]float32, 1536)
	for i := range queryEmbedding {
		queryEmbedding[i] = float32(i) / float32(1536) // Dummy query embedding
	}

	results, err := vectorStore.Search(ctx, queryEmbedding, 3)
	if err != nil {
		log.Fatalf("Failed to search: %v", err)
	}

	fmt.Printf("\nSearch results:\n")
	for i, result := range results {
		if result.Vector.Type == "chunk" {
			fmt.Printf("%d. Score: %.4f, Distance: %.4f\n", i+1, result.Score, result.Distance)
			fmt.Printf("   Content: %s\n", result.Vector.Chunk.Content)
		}
	}

	// Health check
	err = vectorStore.IsHealthy(ctx)
	if err != nil {
		log.Printf("Vector store health check failed: %v", err)
	} else {
		fmt.Println("Vector store is healthy!")
	}
}
