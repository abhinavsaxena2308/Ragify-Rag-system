package models

import (
	"time"

	_ "gorm.io/gorm"
)

type Document struct {
	ID       uint      `json:"id" gorm:"primaryKey"`
	Name     string    `json:"name" gorm:"not null"`
	FileName string    `json:"file_name" gorm:"not null"`
	Size     int64     `json:"size"`
	Type     string    `json:"type"`
	Path     string    `json:"path" gorm:"not null"`
	Uploaded time.Time `json:"uploaded" gorm:"autoCreateTime"`
	UserID   *uint     `json:"user_id,omitempty"`
	Chunks   []Chunk   `json:"chunks" gorm:"foreignKey:DocumentID"`
}

type Chunk struct {
	ID         uint   `json:"id" gorm:"primaryKey"`
	DocumentID uint   `json:"document_id" gorm:"not null;index"`
	Content    string `json:"content" gorm:"not null;size:10000"`
	Embedding  []byte `json:"-" gorm:"type:bytea"` // Store as bytes, not exposed in JSON
	PageNumber *int   `json:"page_number,omitempty"`
	ChunkIndex int    `json:"chunk_index"`
}
