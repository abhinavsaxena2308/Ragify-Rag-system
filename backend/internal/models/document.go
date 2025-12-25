package models

import (
	"time"
)

// Document represents an uploaded document
type Document struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Filename     string    `json:"filename" gorm:"not null"`
	OriginalName string    `json:"original_name" gorm:"not null"`
	ContentType  string    `json:"content_type" gorm:"not null"`
	Size         int64     `json:"size" gorm:"not null"`
	FilePath     string    `json:"file_path" gorm:"not null"`
	TextContent  string    `json:"text_content" gorm:"type:text"`
	PageCount    int       `json:"page_count" gorm:"not null"`
	ExtractedAt  time.Time `json:"extracted_at" gorm:"not null"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	UserID       *uint     `json:"user_id,omitempty"`
	Chunks       []Chunk   `json:"chunks,omitempty" gorm:"foreignKey:DocumentID"`
}

// Chunk represents a text chunk from a document
type Chunk struct {
	ID         uint   `json:"id" gorm:"primaryKey"`
	DocumentID uint   `json:"document_id" gorm:"not null;index"`
	Content    string `json:"content" gorm:"not null;size:10000"`
	Embedding  []byte `json:"-" gorm:"type:bytea"` // Store as bytes, not exposed in JSON
	PageNumber *int   `json:"page_number,omitempty"`
	ChunkIndex int    `json:"chunk_index"`
}

// DocumentResponse represents the response for document operations
type DocumentResponse struct {
	ID           uint   `json:"id"`
	Filename     string `json:"filename"`
	OriginalName string `json:"original_name"`
	ContentType  string `json:"content_type"`
	Size         int64  `json:"size"`
	PageCount    int    `json:"page_count"`
	CreatedAt    string `json:"created_at"`
}

// UploadResponse represents the response after successful upload
type UploadResponse struct {
	Message    string           `json:"message"`
	DocumentID uint             `json:"document_id"`
	Document   DocumentResponse `json:"document"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code,omitempty"`
}
