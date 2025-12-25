package services

import (
	"fmt"
	"ragify-backend/internal/models"
	"ragify-backend/internal/utils"
	"time"

	"gorm.io/gorm"
)

// DocumentService handles document operations
type DocumentService struct {
	db *gorm.DB
}

// NewDocumentService creates a new document service
func NewDocumentService(db *gorm.DB) *DocumentService {
	return &DocumentService{db: db}
}

// UploadDocument handles file upload, validation, and text extraction
func (s *DocumentService) UploadDocument(filename string, filePath string, contentType string, size int64) (*models.Document, error) {
	// Create document record
	document := &models.Document{
		Filename:     filename,
		OriginalName: filename,
		ContentType:  contentType,
		Size:         size,
		FilePath:     filePath,
		CreatedAt:    time.Now(),
	}

	// Save document to database
	if err := s.db.Create(document).Error; err != nil {
		return nil, fmt.Errorf("failed to save document: %v", err)
	}

	// Extract text content asynchronously
	go func() {
		textWithPages, err := utils.ExtractTextWithPageNumbers(filePath)
		if err != nil {
			fmt.Printf("Failed to extract text from %s: %v\n", filename, err)
			return
		}

		// Update document with extracted text
		document.TextContent = textWithPages.Text
		document.PageCount = textWithPages.PageCount
		document.ExtractedAt = time.Now()

		if err := s.db.Save(document).Error; err != nil {
			fmt.Printf("Failed to update document with extracted text: %v\n", err)
		}

		fmt.Printf("Successfully extracted %d pages from %s\n", textWithPages.PageCount, filename)
	}()

	return document, nil
}

// GetDocuments retrieves all documents for a user
func (s *DocumentService) GetDocuments() ([]models.Document, error) {
	var documents []models.Document

	// Get documents ordered by creation date (newest first)
	if err := s.db.Order("created_at DESC").Find(&documents).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve documents: %v", err)
	}

	return documents, nil
}

// GetDocument retrieves a specific document by ID
func (s *DocumentService) GetDocument(id uint) (*models.Document, error) {
	var document models.Document

	if err := s.db.First(&document, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("document not found")
		}
		return nil, fmt.Errorf("failed to retrieve document: %v", err)
	}

	return &document, nil
}

// DeleteDocument removes a document and its associated file
func (s *DocumentService) DeleteDocument(id uint) error {
	// First get the document to get the file path
	var document models.Document
	if err := s.db.First(&document, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("document not found")
		}
		return fmt.Errorf("failed to retrieve document: %v", err)
	}

	// Delete the physical file
	if err := utils.DeleteFile(document.FilePath); err != nil {
		fmt.Printf("Failed to delete file %s: %v\n", document.FilePath, err)
		// Continue with database deletion even if file deletion fails
	}

	// Delete the document record and associated chunks
	if err := s.db.Select("Chunks").Delete(&document).Error; err != nil {
		return fmt.Errorf("failed to delete document: %v", err)
	}

	return nil
}

// GetDocumentsByUser retrieves documents for a specific user
func (s *DocumentService) GetDocumentsByUser(userID uint) ([]models.Document, error) {
	var documents []models.Document

	if err := s.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&documents).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve user documents: %v", err)
	}

	return documents, nil
}

// UpdateDocument updates an existing document
func (s *DocumentService) UpdateDocument(document *models.Document) error {
	if err := s.db.Save(document).Error; err != nil {
		return fmt.Errorf("failed to update document: %v", err)
	}

	return nil
}

// SearchDocuments searches documents by filename or content
func (s *DocumentService) SearchDocuments(query string, userID uint) ([]models.Document, error) {
	var documents []models.Document

	// Search in filename and text content
	searchPattern := "%" + query + "%"
	if err := s.db.Where(
		"(user_id = ? OR user_id IS NULL) AND (filename LIKE ? OR text_content LIKE ?)",
		userID, searchPattern, searchPattern,
	).Order("created_at DESC").Find(&documents).Error; err != nil {
		return nil, fmt.Errorf("failed to search documents: %v", err)
	}

	return documents, nil
}
