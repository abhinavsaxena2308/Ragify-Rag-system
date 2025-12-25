package services

import (
	"context"
	"testing"

	"ragify-backend/internal/config"
	"ragify-backend/internal/models"
)

// ExampleVectorStoreService demonstrates how to use the VectorStoreService
func ExampleVectorStoreService() {
	// Initialize configuration
	cfg := config.FAISSConfig{
		IndexPath: "./data/example_faiss_index",
	}

	// Create vector store service (assuming 384 dimensions for embeddings)
	dimension := 384
	vectorStore, err := NewVectorStoreService(cfg, dimension)
	if err != nil {
		panic(err)
	}
	defer vectorStore.Close()

	// Example: Add a chunk with embedding
	ctx := context.Background()

	chunk := &models.Chunk{
		ID:      1,
		Content: "This is a sample chunk of text about artificial intelligence and machine learning.",
		// ... other chunk fields
	}

	// Example embedding (normally you'd generate this with an embedding service)
	embedding := make([]float32, dimension)
	for i := range embedding {
		embedding[i] = float32(i) / float32(dimension) // Dummy embedding
	}

	// Add the chunk to the vector store
	err = vectorStore.AddChunk(ctx, chunk, embedding)
	if err != nil {
		panic(err)
	}

	// Example: Search for similar chunks
	queryEmbedding := make([]float32, dimension)
	for i := range queryEmbedding {
		queryEmbedding[i] = float32(i*2) / float32(dimension) // Dummy query embedding
	}

	results, err := vectorStore.Search(ctx, queryEmbedding, 5)
	if err != nil {
		panic(err)
	}

	// Process search results
	for _, result := range results {
		if result.Vector.Type == "chunk" {
			chunk := result.Vector.Chunk
			println("Found chunk:", chunk.Content)
			println("Score:", result.Score)
			println("Distance:", result.Distance)
		}
	}

	// Save the index to disk
	err = vectorStore.Save()
	if err != nil {
		panic(err)
	}

	// Get statistics
	stats := vectorStore.GetStats()
	println("Total vectors:", stats["total_vectors"])
	println("Dimension:", stats["dimension"])
}

// TestVectorStoreService is a basic test for the VectorStoreService
func TestVectorStoreService(t *testing.T) {
	// Setup
	cfg := config.FAISSConfig{
		IndexPath: "./data/test_faiss_index",
	}

	dimension := 128
	vectorStore, err := NewVectorStoreService(cfg, dimension)
	if err != nil {
		t.Fatalf("Failed to create vector store: %v", err)
	}
	defer vectorStore.Close()

	ctx := context.Background()

	// Test adding vectors
	chunk := &models.Chunk{
		ID:      1,
		Content: "Test chunk content",
	}

	embedding := make([]float32, dimension)
	for i := range embedding {
		embedding[i] = 0.1
	}

	err = vectorStore.AddChunk(ctx, chunk, embedding)
	if err != nil {
		t.Fatalf("Failed to add chunk: %v", err)
	}

	// Test search
	queryEmbedding := make([]float32, dimension)
	for i := range queryEmbedding {
		queryEmbedding[i] = 0.1
	}

	results, err := vectorStore.Search(ctx, queryEmbedding, 1)
	if err != nil {
		t.Fatalf("Failed to search: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	if results[0].Vector.Chunk.ID != chunk.ID {
		t.Errorf("Expected chunk ID %d, got %d", chunk.ID, results[0].Vector.Chunk.ID)
	}

	// Test save/load
	err = vectorStore.Save()
	if err != nil {
		t.Fatalf("Failed to save: %v", err)
	}

	// Create new instance to test loading
	vectorStore2, err := NewVectorStoreService(cfg, dimension)
	if err != nil {
		t.Fatalf("Failed to create vector store for loading: %v", err)
	}
	defer vectorStore2.Close()

	results2, err := vectorStore2.Search(ctx, queryEmbedding, 1)
	if err != nil {
		t.Fatalf("Failed to search after loading: %v", err)
	}

	if len(results2) != 1 {
		t.Fatalf("Expected 1 result after loading, got %d", len(results2))
	}
}

// BenchmarkVectorStoreSearch benchmarks the search functionality
func BenchmarkVectorStoreSearch(b *testing.B) {
	cfg := config.FAISSConfig{
		IndexPath: "./data/benchmark_faiss_index",
	}

	dimension := 384
	vectorStore, err := NewVectorStoreService(cfg, dimension)
	if err != nil {
		b.Fatalf("Failed to create vector store: %v", err)
	}
	defer vectorStore.Close()

	ctx := context.Background()

	// Add some test data
	for i := 0; i < 1000; i++ {
		chunk := &models.Chunk{
			ID:      uint(i),
			Content: "Test chunk content",
		}

		embedding := make([]float32, dimension)
		for j := range embedding {
			embedding[j] = float32(i+j) / float32(dimension)
		}

		err = vectorStore.AddChunk(ctx, chunk, embedding)
		if err != nil {
			b.Fatalf("Failed to add chunk %d: %v", i, err)
		}
	}

	// Benchmark search
	queryEmbedding := make([]float32, dimension)
	for i := range queryEmbedding {
		queryEmbedding[i] = 0.5
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := vectorStore.Search(ctx, queryEmbedding, 10)
		if err != nil {
			b.Fatalf("Failed to search: %v", err)
		}
	}
}

// ExampleNewVectorStoreService shows integration with the embedding service
func ExampleNewVectorStoreService() {
	// This would typically be done in your main application setup

	// 1. Initialize embedding service
	/*
		embeddingConfig := config.DefaultEmbeddingConfig()
		embeddingConfig.ModelName = "openai/text-embedding-3-small"
		embeddingService := NewOpenRouterEmbeddingService(embeddingConfig)

		// 2. Initialize vector store
		faissConfig := config.FAISSConfig{
			IndexPath: "./data/production_faiss_index",
		}

		// Get dimensions from the embedding provider
		dimension := embeddingService.GetDimensions()
		vectorStore, err := NewVectorStoreService(faissConfig, dimension)
		if err != nil {
			panic(err)
		}
		defer vectorStore.Close()

		// 3. Process documents and add to vector store
		ctx := context.Background()

		for _, document := range documents {
			// Generate embeddings for chunks
			embeddings, err := embeddingService.GenerateDocumentEmbeddings(ctx, document, chunks)
			if err != nil {
				log.Printf("Failed to generate embeddings for document %d: %v", document.ID, err)
				continue
			}

			// Add each chunk embedding to vector store
			for chunkID, embedding := range embeddings {
				chunk := chunks[chunkID]
				err := vectorStore.AddChunk(ctx, chunk, embedding)
				if err != nil {
					log.Printf("Failed to add chunk %d to vector store: %v", chunkID, err)
					continue
				}
			}
		}

		// 4. Save the index
		err = vectorStore.Save()
		if err != nil {
			log.Printf("Failed to save vector store: %v", err)
		}
	*/
}
