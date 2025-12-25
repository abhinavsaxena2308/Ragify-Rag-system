package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"sync"
	"time"

	"ragify-backend/internal/config"
	"ragify-backend/internal/models"
)

// VectorMetadata stores metadata associated with a vector
type VectorMetadata struct {
	ID       uint                   `json:"id"`
	Type     string                 `json:"type"` // "chunk" or "document"
	Document *models.Document       `json:"document,omitempty"`
	Chunk    *models.Chunk          `json:"chunk,omitempty"`
	Extra    map[string]interface{} `json:"extra,omitempty"`
}

// SearchResult represents a similarity search result
type SearchResult struct {
	Vector   VectorMetadata `json:"vector"`
	Score    float32        `json:"score"`
	Distance float32        `json:"distance"`
}

// PureGoVectorStore provides a pure Go implementation of vector storage and similarity search
// This implementation uses basic in-memory storage and cosine similarity for search
type PureGoVectorStore struct {
	vectors      []VectorWithMetadata
	mu           sync.RWMutex
	config       config.FAISSConfig
	dimension    int
	indexPath    string
	metadataPath string
}

// VectorWithMetadata combines a vector with its metadata
type VectorWithMetadata struct {
	Vector   []float32      `json:"vector"`
	Metadata VectorMetadata `json:"metadata"`
	AddedAt  time.Time      `json:"added_at"`
}

// NewPureGoVectorStore creates a new pure Go vector store
func NewPureGoVectorStore(cfg config.FAISSConfig, dimension int) (*PureGoVectorStore, error) {
	store := &PureGoVectorStore{
		config:       cfg,
		dimension:    dimension,
		vectors:      make([]VectorWithMetadata, 0),
		indexPath:    cfg.IndexPath + ".json",
		metadataPath: cfg.IndexPath + "_meta.json",
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(cfg.IndexPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create index directory: %w", err)
	}

	// Try to load existing data
	if err := store.load(); err != nil {
		log.Printf("Warning: failed to load existing vector store: %v", err)
	}

	log.Printf("Initialized pure Go vector store with %d vectors", len(store.vectors))
	return store, nil
}

// load loads vectors from disk
func (s *PureGoVectorStore) load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.indexPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // No file yet
		}
		return err
	}

	return json.Unmarshal(data, &s.vectors)
}

// save saves vectors to disk
func (s *PureGoVectorStore) save() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, err := json.MarshalIndent(s.vectors, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.indexPath, data, 0644)
}

// AddVector adds a vector with metadata to the store
func (s *PureGoVectorStore) AddVector(ctx context.Context, embedding []float32, metadata VectorMetadata) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(embedding) != s.dimension {
		return fmt.Errorf("embedding dimension mismatch: expected %d, got %d", s.dimension, len(embedding))
	}

	vector := VectorWithMetadata{
		Vector:   make([]float32, len(embedding)),
		Metadata: metadata,
		AddedAt:  time.Now(),
	}
	copy(vector.Vector, embedding)

	s.vectors = append(s.vectors, vector)

	log.Printf("Added vector %d to store", len(s.vectors)-1)
	return nil
}

// AddChunk adds a chunk embedding to the vector store
func (s *PureGoVectorStore) AddChunk(ctx context.Context, chunk *models.Chunk, embedding []float32) error {
	metadata := VectorMetadata{
		ID:    chunk.ID,
		Type:  "chunk",
		Chunk: chunk,
	}

	return s.AddVector(ctx, embedding, metadata)
}

// AddDocument adds a document embedding to the vector store
func (s *PureGoVectorStore) AddDocument(ctx context.Context, document *models.Document, embedding []float32) error {
	metadata := VectorMetadata{
		ID:       document.ID,
		Type:     "document",
		Document: document,
	}

	return s.AddVector(ctx, embedding, metadata)
}

// cosineDistance computes cosine distance between two vectors
func (s *PureGoVectorStore) cosineDistance(a, b []float32) float32 {
	if len(a) != len(b) {
		return math.MaxFloat32
	}

	var dotProduct, normA, normB float32

	for i := 0; i < len(a); i++ {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return math.MaxFloat32
	}

	cosineSimilarity := dotProduct / (float32(math.Sqrt(float64(normA))) * float32(math.Sqrt(float64(normB))))

	// Convert similarity to distance (lower is better)
	return 1.0 - cosineSimilarity
}

// Search performs similarity search and returns top-k results
func (s *PureGoVectorStore) Search(ctx context.Context, queryEmbedding []float32, k int) ([]SearchResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(queryEmbedding) != s.dimension {
		return nil, fmt.Errorf("query embedding dimension mismatch: expected %d, got %d", s.dimension, len(queryEmbedding))
	}

	if len(s.vectors) == 0 {
		return []SearchResult{}, nil
	}

	if k <= 0 || k > len(s.vectors) {
		k = len(s.vectors)
	}

	// Calculate distances for all vectors
	type resultWithIndex struct {
		index    int
		distance float32
	}

	results := make([]resultWithIndex, len(s.vectors))
	for i, vector := range s.vectors {
		distance := s.cosineDistance(queryEmbedding, vector.Vector)
		results[i] = resultWithIndex{index: i, distance: distance}
	}

	// Simple selection sort for top-k (for small k, this is efficient)
	for i := 0; i < k && i < len(results); i++ {
		minIdx := i
		for j := i + 1; j < len(results); j++ {
			if results[j].distance < results[minIdx].distance {
				minIdx = j
			}
		}
		results[i], results[minIdx] = results[minIdx], results[i]
	}

	// Build final results
	searchResults := make([]SearchResult, 0, k)
	for i := 0; i < k && i < len(results); i++ {
		vector := s.vectors[results[i].index]

		// Convert distance to similarity score (higher is better)
		score := 1.0 - results[i].distance

		searchResults = append(searchResults, SearchResult{
			Vector:   vector.Metadata,
			Score:    score,
			Distance: results[i].distance,
		})
	}

	return searchResults, nil
}

// SearchChunks performs similarity search and returns only chunk results
func (s *PureGoVectorStore) SearchChunks(ctx context.Context, queryEmbedding []float32, k int) ([]SearchResult, error) {
	allResults, err := s.Search(ctx, queryEmbedding, k*2) // Get more results to filter
	if err != nil {
		return nil, err
	}

	// Filter for chunks only
	chunkResults := make([]SearchResult, 0)
	for _, result := range allResults {
		if result.Vector.Type == "chunk" {
			chunkResults = append(chunkResults, result)
			if len(chunkResults) >= k {
				break
			}
		}
	}

	return chunkResults, nil
}

// SearchDocuments performs similarity search and returns only document results
func (s *PureGoVectorStore) SearchDocuments(ctx context.Context, queryEmbedding []float32, k int) ([]SearchResult, error) {
	allResults, err := s.Search(ctx, queryEmbedding, k*2) // Get more results to filter
	if err != nil {
		return nil, err
	}

	// Filter for documents only
	docResults := make([]SearchResult, 0)
	for _, result := range allResults {
		if result.Vector.Type == "document" {
			docResults = append(docResults, result)
			if len(docResults) >= k {
				break
			}
		}
	}

	return docResults, nil
}

// Save saves the vector store to disk
func (s *PureGoVectorStore) Save() error {
	err := s.save()
	if err != nil {
		return fmt.Errorf("failed to save vector store: %w", err)
	}

	log.Printf("Saved vector store with %d vectors to %s", len(s.vectors), s.indexPath)
	return nil
}

// Load loads the vector store from disk
func (s *PureGoVectorStore) Load() error {
	return s.load()
}

// GetStats returns statistics about the vector store
func (s *PureGoVectorStore) GetStats() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stats := map[string]interface{}{
		"total_vectors": len(s.vectors),
		"dimension":     s.dimension,
		"index_type":    "pure_go_cosine",
		"index_path":    s.indexPath,
		"metadata_path": s.metadataPath,
	}

	// Count by type
	typeCounts := make(map[string]int)
	for _, vector := range s.vectors {
		typeCounts[vector.Metadata.Type]++
	}
	stats["type_counts"] = typeCounts

	return stats
}

// Delete removes a vector from the store
func (s *PureGoVectorStore) Delete(ctx context.Context, id uint) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, vector := range s.vectors {
		if vector.Metadata.ID == id {
			// Remove from slice
			s.vectors = append(s.vectors[:i], s.vectors[i+1:]...)
			log.Printf("Removed vector with ID %d from store", id)
			return nil
		}
	}

	return fmt.Errorf("vector with ID %d not found", id)
}

// Clear removes all vectors from the store
func (s *PureGoVectorStore) Clear(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.vectors = make([]VectorWithMetadata, 0)
	log.Printf("Cleared all vectors from store")
	return nil
}

// Close cleans up resources
func (s *PureGoVectorStore) Close() error {
	// Save before closing
	if err := s.Save(); err != nil {
		log.Printf("Warning: failed to save vector store on close: %v", err)
	}

	return nil
}

// IsHealthy checks if the vector store is healthy
func (s *PureGoVectorStore) IsHealthy(ctx context.Context) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Basic health check - verify dimension consistency
	for _, vector := range s.vectors {
		if len(vector.Vector) != s.dimension {
			return fmt.Errorf("vector dimension mismatch: expected %d, got %d", s.dimension, len(vector.Vector))
		}
	}

	return nil
}

// VectorStoreService interface to allow switching between implementations
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

// NewVectorStoreService creates a vector store service (using pure Go implementation for compatibility)
func NewVectorStoreService(cfg config.FAISSConfig, dimension int) (VectorStoreService, error) {
	// Use pure Go implementation for better compatibility
	return NewPureGoVectorStore(cfg, dimension)
}
