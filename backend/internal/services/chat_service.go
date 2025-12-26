package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
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

// generateSessionID generates a random session ID
func generateSessionID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
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

// MessageResponse represents a message with deserialized source docs
type MessageResponse struct {
	ID         uint      `json:"id"`
	SessionID  string    `json:"session_id"`
	Role       string    `json:"role"`
	Content    string    `json:"content"`
	Timestamp  time.Time `json:"timestamp"`
	SourceDocs []string  `json:"source_docs"`
}

func (s *ChatService) GetMessagesBySessionID(sessionID string) ([]MessageResponse, error) {
	var messages []models.Message
	err := s.DB.Where("session_id = ?", sessionID).Order("timestamp ASC").Find(&messages).Error
	if err != nil {
		return nil, err
	}

	var messageResponses []MessageResponse
	for _, msg := range messages {
		var sourceDocs []string
		if msg.SourceDocs != "" {
			json.Unmarshal([]byte(msg.SourceDocs), &sourceDocs)
		}

		messageResponses = append(messageResponses, MessageResponse{
			ID:         msg.ID,
			SessionID:  msg.SessionID,
			Role:       msg.Role,
			Content:    msg.Content,
			Timestamp:  msg.Timestamp,
			SourceDocs: sourceDocs,
		})
	}

	return messageResponses, nil
}

func (s *ChatService) AskQuestion(sessionID, query string, documentIDs []uint) (string, []string, error) {
	// Create session if none provided
	if sessionID == "" {
		session := &models.ChatSession{
			SessionID: generateSessionID(),
		}
		if err := s.CreateSession(session); err != nil {
			return "", nil, err
		}
		sessionID = session.SessionID
	}

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
	sourceDocsJSON, _ := json.Marshal(sourceDocuments)
	aiMessage := &models.Message{
		SessionID:  sessionID,
		Content:    response,
		Role:       "assistant",
		Timestamp:  time.Now(),
		SourceDocs: string(sourceDocsJSON),
	}
	if err := s.CreateMessage(aiMessage); err != nil {
		return "", nil, err
	}

	return response, sourceDocuments, nil
}
