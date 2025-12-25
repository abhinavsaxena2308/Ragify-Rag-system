package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"ragify-backend/internal/models"
	"ragify-backend/internal/services"

	"github.com/labstack/echo/v4"
)

// AskHandler handles RAG question-answering requests
type AskHandler struct {
	ragService *services.RAGService
}

// NewAskHandler creates a new AskHandler
func NewAskHandler(ragService *services.RAGService) *AskHandler {
	return &AskHandler{
		ragService: ragService,
	}
}

// Ask handles POST /api/ask requests
func (h *AskHandler) Ask(c echo.Context) error {
	startTime := time.Now()

	// Parse request
	var req models.RAGRequest
	if err := c.Bind(&req); err != nil {
		log.Printf("Failed to bind request: %v", err)
		return c.JSON(http.StatusBadRequest, ErrorResponse("Invalid request format", err))
	}

	// Validate request
	if err := h.validateRequest(&req); err != nil {
		log.Printf("Request validation failed: %v", err)
		return c.JSON(http.StatusBadRequest, ErrorResponse("Validation failed", err))
	}

	log.Printf("Processing ask request: question=%s, top_k=%d", req.Question, req.TopK)

	// Process through RAG pipeline
	ctx, cancel := context.WithTimeout(c.Request().Context(), 120*time.Second)
	defer cancel()

	response, err := h.ragService.ProcessQuestion(ctx, &req)
	if err != nil {
		log.Printf("RAG processing failed: %v", err)
		return c.JSON(http.StatusInternalServerError, ErrorResponse("Failed to process question", err))
	}

	// Log response time
	totalTime := time.Since(startTime)
	log.Printf("Ask request completed in %v: sources=%d", totalTime, len(response.Sources))

	return c.JSON(http.StatusOK, SuccessResponse("Question processed successfully", response))
}

// AskStream handles streaming responses for real-time answers
func (h *AskHandler) AskStream(c echo.Context) error {
	// Parse request
	var req models.RAGRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse("Invalid request format", err))
	}

	// Validate request
	if err := h.validateRequest(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse("Validation failed", err))
	}

	// Set SSE headers
	c.Response().Header().Set("Content-Type", "text/event-stream")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")

	// Process through RAG pipeline
	ctx, cancel := context.WithTimeout(c.Request().Context(), 120*time.Second)
	defer cancel()

	response, err := h.ragService.ProcessQuestion(ctx, &req)
	if err != nil {
		// Send error event
		c.Response().Write([]byte("event: error\n"))
		c.Response().Write([]byte("data: " + err.Error() + "\n\n"))
		return nil
	}

	// Send response as JSON
	c.Response().Write([]byte("event: response\n"))
	jsonData := `{"answer": "` + response.Answer + `", "sources": ` +
		strconv.Itoa(len(response.Sources)) + `}\n\n`
	c.Response().Write([]byte(jsonData))
	c.Response().Write([]byte("event: done\n\n"))

	return nil
}

// Health handles GET /api/health requests for RAG service
func (h *AskHandler) Health(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	// Check RAG service health
	err := h.ragService.IsHealthy(ctx)
	if err != nil {
		return c.JSON(http.StatusServiceUnavailable, ErrorResponse("RAG service unhealthy", err))
	}

	// Get service stats
	stats := h.ragService.GetStats()
	config := h.ragService.GetConfig()

	healthData := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now(),
		"services": map[string]interface{}{
			"rag": map[string]interface{}{
				"status": "healthy",
				"stats":  stats,
				"config": config,
			},
		},
		"version": "1.0.0",
	}

	return c.JSON(http.StatusOK, SuccessResponse("RAG service is healthy", healthData))
}

// Stats handles GET /api/stats requests
func (h *AskHandler) Stats(c echo.Context) error {
	stats := h.ragService.GetStats()
	config := h.ragService.GetConfig()
	promptTemplate := h.ragService.GetPromptTemplate()

	statsData := map[string]interface{}{
		"rag_stats":       stats,
		"rag_config":      config,
		"prompt_template": promptTemplate,
		"timestamp":       time.Now(),
	}

	return c.JSON(http.StatusOK, SuccessResponse("RAG statistics retrieved", statsData))
}

// UpdatePrompt handles PUT /api/prompt requests to update the prompt template
func (h *AskHandler) UpdatePrompt(c echo.Context) error {
	var template models.PromptTemplate
	if err := c.Bind(&template); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse("Invalid request format", err))
	}

	// Validate template
	if template.SystemPrompt == "" || template.UserPrompt == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse("System prompt and user prompt are required", nil))
	}

	// Update template
	h.ragService.UpdatePromptTemplate(&template)

	log.Printf("Updated prompt template")

	return c.JSON(http.StatusOK, SuccessResponse("Prompt template updated successfully", template))
}

// GetPrompt handles GET /api/prompt requests
func (h *AskHandler) GetPrompt(c echo.Context) error {
	template := h.ragService.GetPromptTemplate()
	return c.JSON(http.StatusOK, SuccessResponse("Prompt template retrieved", template))
}

// validateRequest validates the RAG request
func (h *AskHandler) validateRequest(req *models.RAGRequest) error {
	if req.Question == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Question is required")
	}

	if len(req.Question) > 1000 {
		return echo.NewHTTPError(http.StatusBadRequest, "Question is too long (max 1000 characters)")
	}

	if req.TopK < 0 || req.TopK > 20 {
		return echo.NewHTTPError(http.StatusBadRequest, "TopK must be between 0 and 20")
	}

	return nil
}

// AskWithSession handles POST /api/ask/session requests for conversation-aware responses
func (h *AskHandler) AskWithSession(c echo.Context) error {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse("Session ID is required", nil))
	}

	// Parse request
	var req models.RAGRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse("Invalid request format", err))
	}

	// Add session ID to request
	req.SessionID = sessionID

	// Validate
	if err := h.validateRequest(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse("Validation failed", err))
	}

	// Process through RAG pipeline
	ctx, cancel := context.WithTimeout(c.Request().Context(), 120*time.Second)
	defer cancel()

	response, err := h.ragService.ProcessQuestion(ctx, &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse("Failed to process question", err))
	}

	return c.JSON(http.StatusOK, SuccessResponse("Question processed successfully", response))
}
