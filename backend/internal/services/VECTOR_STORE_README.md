# VectorStoreService

A pure Go implementation of vector storage and similarity search for the Ragify RAG system. This service provides FAISS-like functionality without requiring CGO dependencies.

## Features

- **Vector Storage**: Store embeddings with metadata
- **Similarity Search**: Fast cosine similarity search with top-k results
- **Persistence**: Save/load vector store from disk
- **Metadata Support**: Store chunk and document metadata alongside vectors
- **Type-specific Search**: Search only chunks or only documents
- **Thread-safe**: Concurrent access with proper locking
- **Health Monitoring**: Built-in health checks and statistics

## Quick Start

```go
package main

import (
    "context"
    "log"
    
    "ragify-backend/internal/config"
    "ragify-backend/internal/models"
    "ragify-backend/internal/services"
)

func main() {
    // Initialize vector store
    cfg := config.FAISSConfig{
        IndexPath: "./data/my_vector_store",
    }
    
    vectorStore, err := services.NewVectorStoreService(cfg, 1536) // 1536 dimensions
    if err != nil {
        log.Fatal(err)
    }
    defer vectorStore.Close()
    
    // Add a chunk with embedding
    chunk := &models.Chunk{
        ID:      1,
        Content: "This is a sample document chunk.",
    }
    
    embedding := make([]float32, 1536) // Your embedding here
    err = vectorStore.AddChunk(context.Background(), chunk, embedding)
    if err != nil {
        log.Fatal(err)
    }
    
    // Search for similar content
    queryEmbedding := make([]float32, 1536) // Your query embedding here
    results, err := vectorStore.Search(context.Background(), queryEmbedding, 5)
    if err != nil {
        log.Fatal(err)
    }
    
    for _, result := range results {
        if result.Vector.Type == "chunk" {
            log.Printf("Found: %s (Score: %.4f)", 
                result.Vector.Chunk.Content, result.Score)
        }
    }
}
```

## API Reference

### VectorStoreService Interface

```go
type VectorStoreService interface {
    AddVector(ctx context.Context, embedding []float32, metadata VectorMetadata) error
    AddChunk(ctx context.Context, chunk *models.Chunk, embedding []float32) error
    AddDocument(ctx context.Context, document *models.Document, embedding []float32) error
    Search(ctx context.Context, queryEmbedding []float32, k int) ([]SearchResult, error)
    SearchChunks(ctx context.Context, queryEmbedding []float32, k int) ([]SearchResult, error)
    SearchDocuments(ctx context.Context, queryEmbedding []float32, k int) ([]SearchResult, error)
    Save() error
    Load() error
    GetStats() map[string]interface{}
    Delete(ctx context.Context, id uint) error
    Clear(ctx context.Context) error
    Close() error
    IsHealthy(ctx context.Context) error
}
```

### Core Methods

#### AddVector
Adds a vector with custom metadata to the store.

```go
metadata := VectorMetadata{
    ID:   123,
    Type: "custom",
    Extra: map[string]interface{}{
        "category": "technology",
        "tags":     []string{"AI", "ML"},
    },
}

err := vectorStore.AddVector(ctx, embedding, metadata)
```

#### AddChunk
Convenience method to add a chunk with its embedding.

```go
err := vectorStore.AddChunk(ctx, chunk, embedding)
```

#### AddDocument
Convenience method to add a document with its embedding.

```go
err := vectorStore.AddDocument(ctx, document, embedding)
```

#### Search
Performs similarity search and returns top-k most similar vectors.

```go
results, err := vectorStore.Search(ctx, queryEmbedding, 10)
```

#### SearchChunks
Searches only among chunk vectors.

```go
results, err := vectorStore.SearchChunks(ctx, queryEmbedding, 5)
```

#### SearchDocuments
Searches only among document vectors.

```go
results, err := vectorStore.SearchDocuments(ctx, queryEmbedding, 3)
```

### Data Structures

#### VectorMetadata
```go
type VectorMetadata struct {
    ID       uint                   `json:"id"`
    Type     string                 `json:"type"` // "chunk", "document", or custom
    Document *models.Document       `json:"document,omitempty"`
    Chunk    *models.Chunk          `json:"chunk,omitempty"`
    Extra    map[string]interface{} `json:"extra,omitempty"`
}
```

#### SearchResult
```go
type SearchResult struct {
    Vector     VectorMetadata `json:"vector"`
    Score      float32        `json:"score"`      // Similarity score (higher is better)
    Distance   float32        `json:"distance"`   // Cosine distance (lower is better)
}
```

## Configuration

The vector store is configured via the `FAISSConfig` struct:

```go
type FAISSConfig struct {
    IndexPath string // Base path for storing index files
}
```

Environment variables:
- `FAISS_INDEX_PATH`: Path to store the vector index (default: "./data/faiss_index")

## Performance

The pure Go implementation uses:
- **Cosine Similarity**: For vector similarity calculation
- **Selection Sort**: For top-k selection (efficient for small k)
- **In-memory Storage**: Fast access with JSON persistence

Benchmark results (5000 vectors, 384 dimensions):
```
BenchmarkVectorStoreSearch-4    8    158628325 ns/op
```

## Integration with EmbeddingService

```go
// Initialize embedding service
embeddingConfig := config.DefaultEmbeddingConfig()
embeddingConfig.ModelName = "openai/text-embedding-3-small"
embeddingService := services.NewOpenRouterEmbeddingService(embeddingConfig)

// Get dimensions from embedding provider
dimension := embeddingService.GetDimensions()

// Initialize vector store with correct dimensions
vectorStore, err := services.NewVectorStoreService(cfg, dimension)

// Process documents
for _, document := range documents {
    // Generate embeddings
    embeddings, err := embeddingService.GenerateDocumentEmbeddings(ctx, document, chunks)
    if err != nil {
        continue
    }
    
    // Add to vector store
    for chunkID, embedding := range embeddings {
        chunk := chunks[chunkID]
        err := vectorStore.AddChunk(ctx, chunk, embedding)
        if err != nil {
            continue
        }
    }
}
```

## File Storage

The vector store creates two files:
- `{IndexPath}.json`: Main vector data
- `{IndexPath}_meta.json`: Metadata backup

Example:
```
./data/faiss_index.json      # Vector data
./data/faiss_index_meta.json # Metadata backup
```

## Testing

Run the tests:
```bash
go test ./internal/services -v -run TestVectorStoreService
```

Run benchmarks:
```bash
go test ./internal/services -v -bench=BenchmarkVectorStoreSearch
```

## Limitations

1. **Memory Usage**: All vectors are kept in memory
2. **Search Algorithm**: Uses simple selection sort (O(nk) for top-k)
3. **Deletion**: Vector deletion is O(n) operation
4. **Persistence**: Full JSON serialization on save

## Future Enhancements

- [ ] Add support for different distance metrics (Euclidean, Manhattan)
- [ ] Implement more efficient search algorithms (KD-tree, Annoy)
- [ ] Add incremental persistence
- [ ] Support for vector compression
- [ ] Batch operations
- [ ] Vector indexing strategies

## Migration from FAISS

If you need to switch from the original FAISS implementation:

1. The API remains the same - just change the initialization
2. Vector format is compatible (same []float32)
3. Metadata structure is identical
4. Search results have the same format

```go
// Before (FAISS with CGO)
// vectorStore, err := NewFAISSVectorStore(cfg, dimension)

// After (Pure Go)
vectorStore, err := NewVectorStoreService(cfg, dimension)
```

The pure Go implementation provides better portability and easier deployment while maintaining the same API.
