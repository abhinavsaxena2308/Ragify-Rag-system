package handlers

import (
	"net/http"

	"ragify-backend/internal/models"
	"ragify-backend/internal/services"
	"ragify-backend/internal/utils"

	"github.com/labstack/echo/v4"
)

type ChatHandler struct {
	Service *services.ChatService
}

func NewChatHandler(service *services.ChatService) *ChatHandler {
	return &ChatHandler{Service: service}
}

// AskQuestion handles a question about documents
func (h *ChatHandler) AskQuestion(c echo.Context) error {
	// Request structure
	var req struct {
		SessionID   *string `json:"session_id"`
		Query       string  `json:"query"`
		DocumentIDs []uint  `json:"document_ids,omitempty"`
	}

	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, "Invalid request format")
	}

	if req.Query == "" {
		return utils.BadRequest(c, "Query is required")
	}

	// Handle null session_id by converting to empty string
	sessionID := ""
	if req.SessionID != nil {
		sessionID = *req.SessionID
	}

	// Use the service to process the question
	answer, sources, err := h.Service.AskQuestion(sessionID, req.Query, req.DocumentIDs)
	if err != nil {
		return utils.InternalError(c, "Failed to process question", err)
	}

	response := map[string]interface{}{
		"session_id": sessionID,
		"answer":     answer,
		"sources":    sources,
	}

	return utils.SendSuccess(c, "Question processed successfully", response, http.StatusOK)
}

// CreateSession creates a new chat session
func (h *ChatHandler) CreateSession(c echo.Context) error {
	// Request structure
	var req struct {
		UserID *uint `json:"user_id,omitempty"`
	}

	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, "Invalid request format")
	}

	// Create a new session
	session := &models.ChatSession{
		UserID: req.UserID,
	}

	// Use the service to create the session
	err := h.Service.CreateSession(session)
	if err != nil {
		return utils.InternalError(c, "Failed to create session", err)
	}

	return utils.SendSuccess(c, "Session created successfully", session, http.StatusOK)
}

// GetSession retrieves a chat session
func (h *ChatHandler) GetSession(c echo.Context) error {
	sessionID := c.Param("id")

	// Use the service to get the session
	session, err := h.Service.GetSessionByID(sessionID)
	if err != nil {
		return utils.NotFound(c, "Session not found")
	}

	return utils.SendSuccess(c, "Session retrieved successfully", session, http.StatusOK)
}

// GetSessionMessages retrieves all messages in a chat session
func (h *ChatHandler) GetSessionMessages(c echo.Context) error {
	sessionID := c.Param("id")

	// Use the service to get messages
	messages, err := h.Service.GetMessagesBySessionID(sessionID)
	if err != nil {
		return utils.InternalError(c, "Failed to retrieve messages", err)
	}

	return utils.SendSuccess(c, "Messages retrieved successfully", messages, http.StatusOK)
}
