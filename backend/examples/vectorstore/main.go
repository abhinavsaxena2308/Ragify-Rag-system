package main

import (
	"context"
	"fmt"
	"log"

	"ragify-backend/internal/config"
	"ragify-backend/internal/models"
	"ragify-backend/internal/services"
)

func main() {
	ctx := context.Background()

	// Initialize configuration
	cfg := config.LoadConfig()

	// Initialize embedding service
	embeddingConfig := config.DefaultEmbeddingConfig()
	embeddingConfig.ModelName = "openai/text-embedding-3-small"
	embeddingService := services.NewEmbeddingService(embeddingConfig)

	// Initialize vector store
	vectorStore, err := services.NewVectorStoreService(cfg.FAISS, 1536) // OpenAI embeddings are 1536 dimensions
	if err != nil {
		log.Fatalf("Failed to create vector store: %v", err)
	}

	// Check current stats
	stats := vectorStore.GetStats()
	fmt.Printf("Initial vector store stats: %+v\n", stats)

	// Test adding a chunk
	testChunk := &models.Chunk{
		ID:         1,
		DocumentID: 1,
		Content:    "This is a test chunk about machine learning and artificial intelligence.",
		PageNumber: &[]int{1}[0],
	}

	// Generate embedding for the chunk
	embedding, err := embeddingService.GenerateChunkEmbedding(ctx, testChunk)
	if err != nil {
		log.Fatalf("Failed to generate embedding: %v", err)
	}

	// Add to vector store
	err = vectorStore.AddChunk(ctx, testChunk, embedding)
	if err != nil {
		log.Fatalf("Failed to add chunk: %v", err)
	}

	fmt.Printf("Added chunk to vector store\n")

	// Save vector store
	err = vectorStore.Save()
	if err != nil {
		log.Fatalf("Failed to save vector store: %v", err)
	}

	// Test search
	queryEmbedding, err := embeddingService.GenerateQueryEmbedding(ctx, "What is machine learning?")
	if err != nil {
		log.Fatalf("Failed to generate query embedding: %v", err)
	}

	results, err := vectorStore.SearchChunks(ctx, queryEmbedding, 5)
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
