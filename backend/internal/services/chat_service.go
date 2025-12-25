package services

import (
	"context"
	"ragify-backend/internal/models"
	"time"

	"gorm.io/gorm"
)

type ChatService struct {
	DB         *gorm.DB
	RAGService *RAGService
}

func NewChatService(db *gorm.DB, ragService *RAGService) *ChatService {
	return &ChatService{
		DB:         db,
		RAGService: ragService,
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
	// Use the RAG service to process the question
	req := &models.RAGRequest{
		Question:  query,
		SessionID: sessionID,
	}

	ragResponse, err := s.RAGService.ProcessQuestion(context.Background(), req)
	if err != nil {
		return "", nil, err
	}

	response := ragResponse.Answer
	var sourceDocuments []string
	for _, source := range ragResponse.Sources {
		sourceDocuments = append(sourceDocuments, source.DocumentName)
	}

	// Store the conversation in the database
	// Create user message
	userMessage := &models.Message{
		SessionID: sessionID,
		Content:   query,
		Role:      "user",
		Timestamp: time.Now(),
	}
	if err := s.CreateMessage(userMessage); err != nil {
		return "", nil, err
	}

	// Create AI message
	aiMessage := &models.Message{
		SessionID: sessionID,
		Content:   response,
		Role:      "assistant",
		Timestamp: time.Now(),
	}
	if err := s.CreateMessage(aiMessage); err != nil {
		return "", nil, err
	}

	return response, sourceDocuments, nil
}
