package services

import (
	"ragify-backend/internal/models"

	"gorm.io/gorm"
)

type DocumentService struct {
	DB *gorm.DB
}

func NewDocumentService(db *gorm.DB) *DocumentService {
	return &DocumentService{
		DB: db,
	}
}

func (s *DocumentService) CreateDocument(doc *models.Document) error {
	return s.DB.Create(doc).Error
}

func (s *DocumentService) GetDocuments() ([]models.Document, error) {
	var documents []models.Document
	err := s.DB.Preload("Chunks").Find(&documents).Error
	return documents, err
}

func (s *DocumentService) GetDocumentByID(id uint) (*models.Document, error) {
	var document models.Document
	err := s.DB.Preload("Chunks").First(&document, id).Error
	return &document, err
}

func (s *DocumentService) DeleteDocument(id uint) error {
	return s.DB.Delete(&models.Document{}, id).Error
}
