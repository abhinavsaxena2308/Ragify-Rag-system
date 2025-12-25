package routes

import (
	"ragify-backend/internal/handlers"
	"ragify-backend/internal/services"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func SetupRoutes(e *echo.Echo, db *gorm.DB) {
	// Health check is already set up in main
	// Public routes
	public := e.Group("/api/v1")

	// Initialize services
	docService := services.NewDocumentService(db)
	chatService := services.NewChatService(db)

	// Document routes
	docHandler := handlers.NewDocumentHandler(docService)
	docGroup := public.Group("/documents")
	docGroup.POST("", docHandler.UploadDocument)
	docGroup.GET("", docHandler.GetDocuments)
	docGroup.GET("/:id", docHandler.GetDocument)
	docGroup.DELETE("/:id", docHandler.DeleteDocument)

	// Chat routes
	chatHandler := handlers.NewChatHandler(chatService)
	chatGroup := public.Group("/chat")
	chatGroup.POST("/ask", chatHandler.AskQuestion)
	chatGroup.POST("/session", chatHandler.CreateSession)
	chatGroup.GET("/session/:id", chatHandler.GetSession)
	chatGroup.GET("/session/:id/messages", chatHandler.GetSessionMessages)

	// User routes
	userHandler := handlers.NewUserHandler()
	userGroup := public.Group("/users")
	userGroup.POST("", userHandler.CreateUser)
	userGroup.POST("/login", userHandler.Login)
}
