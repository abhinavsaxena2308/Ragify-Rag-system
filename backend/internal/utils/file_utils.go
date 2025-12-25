package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

// FileValidationResult represents the result of file validation
type FileValidationResult struct {
	IsValid bool
	Error   string
}

// ValidateFile checks if the uploaded file meets the requirements
func ValidateFile(header *multipart.FileHeader) FileValidationResult {
	// Check file size (10MB limit)
	const maxSize = 10 * 1024 * 1024 // 10MB in bytes
	if header.Size > maxSize {
		return FileValidationResult{
			IsValid: false,
			Error:   "File size exceeds 10MB limit",
		}
	}

	// Check file type
	contentType := header.Header.Get("Content-Type")
	allowedTypes := map[string]bool{
		"application/pdf": true,
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
		"application/msword":       true,
		"text/plain":               true,
		"application/vnd.ms-excel": true,
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":         true,
		"application/vnd.ms-powerpoint":                                             true,
		"application/vnd.openxmlformats-officedocument.presentationml.presentation": true,
	}

	if !allowedTypes[contentType] {
		return FileValidationResult{
			IsValid: false,
			Error:   fmt.Sprintf("Invalid file type: %s. Allowed types: PDF, DOC, DOCX, TXT, XLS, XLSX, PPT, PPTX", contentType),
		}
	}

	return FileValidationResult{
		IsValid: true,
		Error:   "",
	}
}

// SaveFile saves the uploaded file to local storage
func SaveFile(header *multipart.FileHeader, file multipart.File) (string, error) {
	// Generate unique filename
	ext := filepath.Ext(header.Filename)
	uniqueID := uuid.New().String()
	filename := fmt.Sprintf("%s%s", uniqueID, ext)

	// Create uploads directory if it doesn't exist
	uploadsDir := "./uploads"
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create uploads directory: %v", err)
	}

	// Create the file on disk
	filePath := filepath.Join(uploadsDir, filename)
	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %v", err)
	}
	defer dst.Close()

	// Copy the uploaded file to the destination
	if _, err := io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("failed to save file: %v", err)
	}

	return filePath, nil
}

// DeleteFile removes a file from the file system
func DeleteFile(filePath string) error {
	if filePath == "" {
		return nil // No file to delete
	}

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil // File doesn't exist, nothing to delete
	}

	// Delete the file
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete file: %v", err)
	}

	return nil
}

// GetFileExtension returns the file extension from MIME type
func GetFileExtension(contentType string) string {
	extensions := map[string]string{
		"application/pdf": ".pdf",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": ".docx",
		"application/msword":       ".doc",
		"text/plain":               ".txt",
		"application/vnd.ms-excel": ".xls",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":         ".xlsx",
		"application/vnd.ms-powerpoint":                                             ".ppt",
		"application/vnd.openxmlformats-officedocument.presentationml.presentation": ".pptx",
	}

	if ext, exists := extensions[contentType]; exists {
		return ext
	}
	return ""
}

// FormatFileSize formats file size in human readable format
func FormatFileSize(bytes int64) string {
	if bytes == 0 {
		return "0 B"
	}

	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), int64(0)
	for bytes >= unit && exp < 3 {
		div, exp = unit, exp+1
		bytes /= unit
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// GetContentTypeFromExtension returns MIME type from file extension
func GetContentTypeFromExtension(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	mimeTypes := map[string]string{
		".pdf":  "application/pdf",
		".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		".doc":  "application/msword",
		".txt":  "text/plain",
		".xls":  "application/vnd.ms-excel",
		".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		".ppt":  "application/vnd.ms-powerpoint",
		".pptx": "application/vnd.openxmlformats-officedocument.presentationml.presentation",
	}

	if mimeType, exists := mimeTypes[ext]; exists {
		return mimeType
	}
	return "application/octet-stream"
}
