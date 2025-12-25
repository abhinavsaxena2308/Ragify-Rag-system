package chunking

import (
	"fmt"
	"log"

	"ragify-backend/internal/models"
	"ragify-backend/internal/services"
	"ragify-backend/internal/utils"
)

func RunChunkingExample() {
	// Create chunking service with default configuration
	config := services.DefaultChunkingConfig()
	chunkingService := services.NewChunkingService(config)

	// Example document text
	sampleText := `Introduction to RAG Systems

Retrieval-Augmented Generation (RAG) is a powerful approach that combines retrieval systems with generative models. This allows the model to access and incorporate external knowledge when generating responses.

How RAG Works

The RAG process involves several key steps:
1. Document ingestion and preprocessing
2. Text chunking and embedding generation
3. Vector storage and indexing
4. Retrieval during query time
5. Context-augmented generation

Benefits of RAG

RAG systems offer several advantages over traditional approaches:
- Reduced hallucination in generated responses
- Access to up-to-date information
- Improved factual accuracy
- Better source attribution
- Domain-specific knowledge integration

Implementation Considerations

When implementing a RAG system, consider:
- Chunk size and overlap parameters
- Embedding model selection
- Vector database choice
- Retrieval strategies
- Performance optimization

Conclusion

RAG represents a significant advancement in making AI systems more reliable and knowledgeable. By combining the strengths of retrieval and generation, we can create systems that are both accurate and helpful.`

	// Create a sample document
	document := &models.Document{
		ID:          1,
		Filename:    "rag_introduction.txt",
		TextContent: sampleText,
		PageCount:   1,
	}

	// Chunk the document
	chunks, err := chunkingService.ChunkDocument(document)
	if err != nil {
		log.Fatalf("Error chunking document: %v", err)
	}

	// Display results
	fmt.Printf("Document chunked into %d chunks\n\n", len(chunks))

	// Create token counter for analysis
	tokenCounter := utils.NewTokenCounter()

	for i, chunk := range chunks {
		fmt.Printf("=== Chunk %d ===\n", i+1)
		fmt.Printf("Document ID: %d\n", chunk.DocumentID)
		if chunk.PageNumber != nil {
			fmt.Printf("Page Number: %d\n", *chunk.PageNumber)
		}
		fmt.Printf("Chunk Index: %d\n", chunk.ChunkIndex)

		// Count tokens in this chunk
		tokens := tokenCounter.CountTokens(chunk.Content)
		fmt.Printf("Token Count: %d\n", tokens)
		fmt.Printf("Content Preview: %.100s...\n", chunk.Content)
		fmt.Printf("Content Length: %d characters\n", len(chunk.Content))
		fmt.Println()
	}

	// Get chunking statistics
	stats := chunkingService.GetChunkingStats(chunks)
	fmt.Println("=== Chunking Statistics ===")
	fmt.Printf("Total Chunks: %v\n", stats["total_chunks"])
	fmt.Printf("Total Tokens: %v\n", stats["total_tokens"])
	fmt.Printf("Average Chunk Size: %v tokens\n", stats["avg_chunk_size"])
	fmt.Printf("Min Chunk Size: %v tokens\n", stats["min_chunk_size"])
	fmt.Printf("Max Chunk Size: %v tokens\n", stats["max_chunk_size"])

	// Validate chunks
	issues := chunkingService.ValidateChunks(chunks)
	if len(issues) > 0 {
		fmt.Println("\n=== Validation Issues ===")
		for _, issue := range issues {
			fmt.Printf("- %s\n", issue)
		}
	} else {
		fmt.Println("\n=== Validation ===")
		fmt.Println("All chunks passed validation!")
	}
}
