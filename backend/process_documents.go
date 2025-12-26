package main

import (
	"log"
	"strings"

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
	chunkingService := services.NewChunkingService(services.DefaultChunkingConfig())

	// Get all documents that have text content but no chunks
	var documents []models.Document
	err = db.Where("text_content IS NOT NULL AND text_content != ''").Find(&documents).Error
	if err != nil {
		log.Fatalf("Failed to fetch documents: %v", err)
	}

	log.Printf("Found %d documents with text content", len(documents))

	processedCount := 0

	for _, doc := range documents {
		if doc.TextContent == "" || strings.TrimSpace(doc.TextContent) == "" {
			log.Printf("Skipping document %d - no text content", doc.ID)
			continue
		}

		// Check if chunks already exist
		var existingChunks int64
		err = db.Model(&models.Chunk{}).Where("document_id = ?", doc.ID).Count(&existingChunks).Error
		if err != nil {
			log.Printf("Failed to check chunks for document %d: %v", doc.ID, err)
			continue
		}

		if existingChunks > 0 {
			log.Printf("Document %d already has %d chunks, skipping", doc.ID, existingChunks)
			continue
		}

		// Create chunks from text content
		chunks, err := chunkingService.ChunkDocument(&doc)
		if err != nil {
			log.Printf("Failed to chunk document %d: %v", doc.ID, err)
			continue
		}

		if len(chunks) == 0 {
			log.Printf("No chunks created for document %d", doc.ID)
			continue
		}

		// Save chunks to database
		for i, chunk := range chunks {
			chunk.ChunkIndex = i
			err = db.Create(&chunk).Error
			if err != nil {
				log.Printf("Failed to save chunk %d for document %d: %v", i, doc.ID, err)
				continue
			}
		}

		log.Printf("Created %d chunks for document %d (%s)", len(chunks), doc.ID, doc.Filename)
		processedCount++
	}

	log.Printf("Processing completed! Created chunks for %d documents", processedCount)
}
