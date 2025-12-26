package models

import (
	"time"

	_ "gorm.io/gorm"
)

type ChatSession struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    *uint     `json:"user_id,omitempty"`
	SessionID string    `json:"session_id" gorm:"not null;uniqueIndex"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Messages  []Message `json:"messages" gorm:"foreignKey:SessionID"`
}

type Message struct {
	ID          uint        `json:"id" gorm:"primaryKey"`
	SessionID   string      `json:"session_id" gorm:"not null;index"`
	Role        string      `json:"role" gorm:"not null"` // 'user' or 'assistant'
	Content     string      `json:"content" gorm:"not null"`
	Timestamp   time.Time   `json:"timestamp" gorm:"autoCreateTime"`
	SourceDocs  string      `json:"source_docs" gorm:"type:text"` // JSON serialized source documents
	ChatSession ChatSession `json:"-" gorm:"foreignKey:SessionID;references:SessionID"`
}
