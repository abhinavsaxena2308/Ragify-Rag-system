package services

import (
	"ragify-backend/internal/models"

	"gorm.io/gorm"
)

type ChatService struct {
	DB *gorm.DB
}

func NewChatService(db *gorm.DB) *ChatService {
	return &ChatService{
		DB: db,
	}
}

func (s *ChatService) CreateSession(session *models.ChatSession) error {
	return s.DB.Create(session).Error
}

func (s *ChatService) GetSessionByID(sessionID string) (*models.ChatSession, error) {
	var session models.ChatSession
	err := s.DB.Where("session_id = ?", sessionID).First(&session).Error
	return &session, err
}

func (s *ChatService) CreateMessage(message *models.Message) error {
	return s.DB.Create(message).Error
}

func (s *ChatService) GetMessagesBySessionID(sessionID string) ([]models.Message, error) {
	var messages []models.Message
	err := s.DB.Where("session_id = ?", sessionID).Order("timestamp ASC").Find(&messages).Error
	return messages, err
}

func (s *ChatService) AskQuestion(sessionID, query string, documentIDs []uint) (string, []string, error) {
	// This is a placeholder implementation
	// In a real implementation, this would:
	// 1. Retrieve relevant document chunks using FAISS
	// 2. Format the context for the LLM
	// 3. Query the LLM for a response
	// 4. Store the conversation in the database

	// For now, return a sample response
	return "This is a sample response to your query: " + query, []string{"document1.pdf", "document2.pdf"}, nil
}
