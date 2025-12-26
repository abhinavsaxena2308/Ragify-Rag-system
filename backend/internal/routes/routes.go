package routes

import (
	"context"
	"math/rand"
	"ragify-backend/internal/config"
	"ragify-backend/internal/handlers"
	"ragify-backend/internal/services"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// MockEmbeddingProvider for testing
type MockEmbeddingProvider struct{}

func (m *MockEmbeddingProvider) GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	// Generate a simple mock embedding of 1536 dimensions (matching OpenAI default)
	embedding := make([]float32, 1536)
	for i := range embedding {
		embedding[i] = rand.Float32()*2 - 1 // Random values between -1 and 1
	}
	return embedding, nil
}

func (m *MockEmbeddingProvider) GenerateBatchEmbeddings(ctx context.Context, texts []string) ([]services.EmbeddingData, error) {
	var embeddings []services.EmbeddingData
	for i, text := range texts {
		embedding, err := m.GenerateEmbedding(ctx, text)
		if err != nil {
			return nil, err
		}
		embeddings = append(embeddings, services.EmbeddingData{
			Object:    "embedding",
			Embedding: embedding,
			Index:     i,
		})
	}
	return embeddings, nil
}

func (m *MockEmbeddingProvider) GetModelName() string {
	return "mock-embedding-model"
}

func (m *MockEmbeddingProvider) GetDimensions() int {
	return 1536
}

func (m *MockEmbeddingProvider) IsHealthy(ctx context.Context) error {
	return nil
}

// SetupRoutes configures all application routes
func SetupRoutes(e *echo.Echo, db *gorm.DB) {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize embedding service (using mock for testing)
	mockProvider := &MockEmbeddingProvider{}
	embeddingConfig := config.DefaultEmbeddingConfig()
	embeddingService := services.NewEmbeddingService(mockProvider, embeddingConfig)

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
