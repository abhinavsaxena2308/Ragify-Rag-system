package handlers

import (
	"net/http"
	"ragify-backend/internal/services"
	"ragify-backend/internal/utils"
	"strconv"

	"github.com/labstack/echo/v4"
)

type DocumentHandler struct {
	Service *services.DocumentService
}

func NewDocumentHandler(service *services.DocumentService) *DocumentHandler {
	return &DocumentHandler{Service: service}
}

// UploadDocument handles document upload
func (h *DocumentHandler) UploadDocument(c echo.Context) error {
	// In a real implementation, this would process the file upload
	// and store it in the database and file system
	return utils.SendSuccess(c, "Document upload endpoint", nil, http.StatusOK)
}

// GetDocuments retrieves all documents for the user
func (h *DocumentHandler) GetDocuments(c echo.Context) error {
	// Use the service to get documents
	documents, err := h.Service.GetDocuments()
	if err != nil {
		return utils.InternalError(c, "Failed to retrieve documents", err)
	}

	return utils.SendSuccess(c, "Documents retrieved successfully", documents, http.StatusOK)
}

// GetDocument retrieves a specific document by ID
func (h *DocumentHandler) GetDocument(c echo.Context) error {
	id := c.Param("id")
	docID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return utils.BadRequest(c, "Invalid document ID")
	}

	// Use the service to get the document
	document, err := h.Service.GetDocumentByID(uint(docID))
	if err != nil {
		return utils.NotFound(c, "Document not found")
	}

	return utils.SendSuccess(c, "Document retrieved successfully", document, http.StatusOK)
}

// DeleteDocument deletes a document by ID
func (h *DocumentHandler) DeleteDocument(c echo.Context) error {
	id := c.Param("id")
	docID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return utils.BadRequest(c, "Invalid document ID")
	}

	// Use the service to delete the document
	err = h.Service.DeleteDocument(uint(docID))
	if err != nil {
		return utils.InternalError(c, "Failed to delete document", err)
	}

	return utils.SendSuccess(c, "Document deleted successfully", map[string]interface{}{
		"id": docID,
	}, http.StatusOK)
}
