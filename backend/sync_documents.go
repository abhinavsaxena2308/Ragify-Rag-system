package main

import (
	"context"
	"log"

	"ragify-backend/internal/config"
	"ragify-backend/internal/models"
	"ragify-backend/internal/services"
	"ragify-backend/internal/utils"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	db, err := utils.DBConnect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize services
	embeddingConfig := config.DefaultEmbeddingConfig()
	embeddingService := services.NewOpenRouterEmbeddingService(embeddingConfig)

	vectorStore, err := services.NewVectorStoreService(cfg.FAISS, 1536)
	if err != nil {
		log.Fatalf("Failed to create vector store: %v", err)
	}

	// Get all documents from database
	var documents []models.Document
	err = db.Find(&documents).Error
	if err != nil {
		log.Fatalf("Failed to fetch documents: %v", err)
	}

	log.Printf("Found %d documents in database", len(documents))

	ctx := context.Background()
	syncedCount := 0

	// Sync each document
	for _, doc := range documents {
		// Get chunks for this document
		var chunks []models.Chunk
		err = db.Where("document_id = ?", doc.ID).Find(&chunks).Error
		if err != nil {
			log.Printf("Failed to fetch chunks for document %d: %v", doc.ID, err)
			continue
		}

		// Add each chunk to vector store
		for _, chunk := range chunks {
			embedding, err := embeddingService.GenerateChunkEmbedding(ctx, &chunk)
			if err != nil {
				log.Printf("Failed to generate embedding for chunk %d: %v", chunk.ID, err)
				continue
			}

			err = vectorStore.AddChunk(ctx, &chunk, embedding)
			if err != nil {
				log.Printf("Failed to add chunk %d to vector store: %v", chunk.ID, err)
				continue
			}

			syncedCount++
		}

		log.Printf("Synced document %d (%s) with %d chunks", doc.ID, doc.Filename, len(chunks))
	}

	// Save vector store
	err = vectorStore.Save()
	if err != nil {
		log.Fatalf("Failed to save vector store: %v", err)
	}

	log.Printf("Sync completed! Added %d chunks to vector store", syncedCount)
}
