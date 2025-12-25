package routes

import (
	"ragify-backend/internal/config"
	"ragify-backend/internal/handlers"
	"ragify-backend/internal/services"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// SetupRoutes configures all application routes
func SetupRoutes(e *echo.Echo, db *gorm.DB) {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize embedding service
	embeddingConfig := config.DefaultEmbeddingConfig()
	embeddingService := services.NewOpenRouterEmbeddingService(embeddingConfig)

	// Initialize vector store
	vectorStore, err := services.NewVectorStoreService(cfg.FAISS, cfg.RAG.EmbeddingDimensions)
	if err != nil {
		panic("failed to create vector store: " + err.Error())
	}

	// Initialize LLM service
	var llmService services.LLMService
	switch cfg.LLM.Provider {
	case "ollama":
		llmService = services.NewOllamaService(cfg.LLM.Endpoint, cfg.LLM.Model)
	case "openrouter":
		llmService = services.NewOpenRouterService(cfg.LLM.APIKey, cfg.LLM.Model)
	default:
		llmService = services.NewOllamaService(cfg.LLM.Endpoint, cfg.LLM.Model)
	}

	// Initialize services
	documentService := services.NewDocumentService(db)
	ragService := services.NewRAGService(embeddingService, vectorStore, llmService, cfg.RAG)
	chatService := services.NewChatService(db, ragService)

	// Initialize handlers
	documentHandler := handlers.NewDocumentHandler(documentService)
	chatHandler := handlers.NewChatHandler(chatService)
	askHandler := handlers.NewAskHandler(ragService)

	// API v1 routes
	api := e.Group("/api/v1")

	// Document routes
	documents := api.Group("/documents")
	documents.POST("", documentHandler.UploadDocument)
	documents.GET("", documentHandler.GetDocuments)
	documents.GET("/:id", documentHandler.GetDocument)
	documents.PUT("/:id", documentHandler.UpdateDocument)
	documents.DELETE("/:id", documentHandler.DeleteDocument)

	// Chat routes
	chat := api.Group("/chat")
	chat.POST("/sessions", chatHandler.CreateSession)
	chat.GET("/sessions/:id", chatHandler.GetSession)
	chat.GET("/sessions/:id/messages", chatHandler.GetSessionMessages)
	chat.POST("/ask", chatHandler.AskQuestion)

	// RAG routes
	rag := api.Group("/rag")
	rag.POST("/ask", askHandler.Ask)
}
