package handlers

import (
	"fmt"
	"net/http"
	"ragify-backend/internal/models"
	"ragify-backend/internal/services"
	"ragify-backend/internal/utils"
	"strconv"

	"github.com/labstack/echo/v4"
)

// DocumentHandler handles document-related HTTP requests
type DocumentHandler struct {
	documentService *services.DocumentService
}

// NewDocumentHandler creates a new document handler
func NewDocumentHandler(documentService *services.DocumentService) *DocumentHandler {
	return &DocumentHandler{documentService: documentService}
}

// UploadDocument handles file upload requests
func (h *DocumentHandler) UploadDocument(c echo.Context) error {
	// Get file from form
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Failed to get file from form",
			Message: "No file provided or invalid form data",
		})
	}

	// Validate file
	validation := utils.ValidateFile(file)
	if !validation.IsValid {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Validation failed",
			Message: validation.Error,
		})
	}

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to open uploaded file",
			Message: err.Error(),
		})
	}
	defer src.Close()

	// Save file to local storage
	filePath, err := utils.SaveFile(file, src)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to save file",
			Message: err.Error(),
		})
	}

	// Create document record in database
	document, err := h.documentService.UploadDocument(file.Filename, filePath, file.Header.Get("Content-Type"), file.Size)
	if err != nil {
		// Clean up file if database operation fails
		utils.DeleteFile(filePath)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to save document",
			Message: err.Error(),
		})
	}

	// Prepare response
	response := models.UploadResponse{
		Message:    "Document uploaded successfully",
		DocumentID: document.ID,
		Document: models.DocumentResponse{
			ID:           document.ID,
			Filename:     document.Filename,
			OriginalName: document.OriginalName,
			ContentType:  document.ContentType,
			Size:         document.Size,
			PageCount:    document.PageCount,
			CreatedAt:    document.CreatedAt.Format("2006-01-02 15:04:05"),
		},
	}

	return c.JSON(http.StatusCreated, response)
}

// GetDocuments retrieves all documents
func (h *DocumentHandler) GetDocuments(c echo.Context) error {
	// Get user ID from context (if available)
	userID := c.Get("user_id")

	var documents []models.Document
	var err error

	// If user ID is available, get user's documents, otherwise get all documents
	if userID != nil {
		documents, err = h.documentService.GetDocumentsByUser(userID.(uint))
	} else {
		documents, err = h.documentService.GetDocuments()
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to retrieve documents",
			Message: err.Error(),
		})
	}

	// Convert to response format
	var documentResponses []models.DocumentResponse
	for _, doc := range documents {
		documentResponses = append(documentResponses, models.DocumentResponse{
			ID:           doc.ID,
			Filename:     doc.Filename,
			OriginalName: doc.OriginalName,
			ContentType:  doc.ContentType,
			Size:         doc.Size,
			PageCount:    doc.PageCount,
			CreatedAt:    doc.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return c.JSON(http.StatusOK, documentResponses)
}

// GetDocument retrieves a specific document by ID
func (h *DocumentHandler) GetDocument(c echo.Context) error {
	// Get document ID from URL parameter
	idParam := c.Param("id")
	documentID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid document ID",
			Message: "Document ID must be a valid number",
		})
	}

	// Get document from database
	document, err := h.documentService.GetDocument(uint(documentID))
	if err != nil {
		if err.Error() == "document not found" {
			return c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error:   "Document not found",
				Message: fmt.Sprintf("Document with ID %s not found", idParam),
			})
		}
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to retrieve document",
			Message: err.Error(),
		})
	}

	// Prepare response
	response := models.DocumentResponse{
		ID:           document.ID,
		Filename:     document.Filename,
		OriginalName: document.OriginalName,
		ContentType:  document.ContentType,
		Size:         document.Size,
		PageCount:    document.PageCount,
		CreatedAt:    document.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	return c.JSON(http.StatusOK, response)
}

// DeleteDocument removes a document
func (h *DocumentHandler) DeleteDocument(c echo.Context) error {
	// Get document ID from URL parameter
	idParam := c.Param("id")
	documentID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid document ID",
			Message: "Document ID must be a valid number",
		})
	}

	// Delete document
	err = h.documentService.DeleteDocument(uint(documentID))
	if err != nil {
		if err.Error() == "document not found" {
			return c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error:   "Document not found",
				Message: fmt.Sprintf("Document with ID %s not found", idParam),
			})
		}
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to delete document",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Document deleted successfully",
	})
}
