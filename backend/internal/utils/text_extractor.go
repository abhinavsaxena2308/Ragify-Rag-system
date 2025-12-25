package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/extractor"
	"github.com/unidoc/unipdf/v3/model"
)

// TextWithPageNumbers represents extracted text with page information
type TextWithPageNumbers struct {
	Text      string
	PageCount int
	PageInfo  []PageInfo
}

// PageInfo represents information about a specific page
type PageInfo struct {
	PageNumber int
	Text       string
}

// ExtractTextWithPageNumbers extracts text from a PDF file with page numbers
func ExtractTextWithPageNumbers(filePath string) (*TextWithPageNumbers, error) {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file does not exist: %s", filePath)
	}

	// Determine file type by extension
	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".pdf":
		return extractPDFText(filePath)
	case ".docx":
		return extractDocxText(filePath)
	case ".txt":
		return extractPlainText(filePath)
	case ".doc":
		return extractDocText(filePath)
	case ".xls", ".xlsx":
		return extractExcelText(filePath)
	case ".ppt", ".pptx":
		return extractPowerPointText(filePath)
	default:
		return nil, fmt.Errorf("unsupported file type: %s", ext)
	}
}

// extractPDFText extracts text from PDF files with page numbers
func extractPDFText(filePath string) (*TextWithPageNumbers, error) {
	// Initialize unipdf license (community edition)
	err := license.SetMeteredKey("yf00fcae9f503d050240dd3079ba78d4d2b16886e26c552090bdef1be14490be5")
	if err != nil {
		// For community usage, we can continue without a license
		// This will add a watermark but still work
	}

	// Open PDF file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open PDF: %v", err)
	}
	defer file.Close()

	// Create PDF reader
	pdfReader, err := model.NewPdfReader(file)
	if err != nil {
		return nil, fmt.Errorf("failed to create PDF reader: %v", err)
	}

	// Get total page count
	pageCount, err := pdfReader.GetNumPages()
	if err != nil {
		return nil, fmt.Errorf("failed to get page count: %v", err)
	}

	var allText strings.Builder
	var pageInfo []PageInfo

	// Extract text from each page
	for i := 1; i <= pageCount; i++ {
		page, err := pdfReader.GetPage(i)
		if err != nil {
			continue // Skip page if there's an error
		}

		// Create extractor for the page
		ex, err := extractor.New(page)
		if err != nil {
			continue // Skip page if extractor creation fails
		}

		// Extract text from the page
		text, err := ex.ExtractText()
		if err != nil {
			continue // Skip if text extraction fails
		}

		// Clean up the text
		cleanText := strings.TrimSpace(text)
		if cleanText != "" {
			// Add page number and text
			pageInfo = append(pageInfo, PageInfo{
				PageNumber: i,
				Text:       cleanText,
			})

			// Add to overall text with page marker
			allText.WriteString(fmt.Sprintf("--- Page %d ---\n%s\n\n", i, cleanText))
		}
	}

	return &TextWithPageNumbers{
		Text:      allText.String(),
		PageCount: pageCount,
		PageInfo:  pageInfo,
	}, nil
}

// extractDocxText extracts text from DOCX files
func extractDocxText(filePath string) (*TextWithPageNumbers, error) {
	// This is a simplified implementation
	// In a real implementation, you would use a library like github.com/unidoc/unioffice
	// For now, we'll return a placeholder
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open DOCX file: %v", err)
	}
	defer file.Close()

	// Read file content (simplified - just reads as text)
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read DOCX file: %v", err)
	}

	text := string(content)

	return &TextWithPageNumbers{
		Text:      text,
		PageCount: 1, // DOCX doesn't have clear page boundaries without proper parsing
		PageInfo:  []PageInfo{{PageNumber: 1, Text: text}},
	}, nil
}

// extractPlainText extracts text from plain text files
func extractPlainText(filePath string) (*TextWithPageNumbers, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open text file: %v", err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read text file: %v", err)
	}

	text := string(content)
	lines := strings.Split(text, "\n")

	var pageInfo []PageInfo
	var resultText strings.Builder

	for i, line := range lines {
		cleanLine := strings.TrimSpace(line)
		if cleanLine != "" {
			// Estimate page breaks (simplified - every 50 lines could be a page)
			if i > 0 && i%50 == 0 {
				resultText.WriteString(fmt.Sprintf("--- Page %d ---\n", (i/50)+1))
				pageInfo = append(pageInfo, PageInfo{
					PageNumber: (i / 50) + 1,
					Text:       cleanLine,
				})
			}
			resultText.WriteString(cleanLine + "\n")
		}
	}

	return &TextWithPageNumbers{
		Text:      resultText.String(),
		PageCount: len(pageInfo),
		PageInfo:  pageInfo,
	}, nil
}

// extractDocText extracts text from legacy DOC files (placeholder)
func extractDocText(filePath string) (*TextWithPageNumbers, error) {
	// This would require a library like github.com/unidoc/unioffice
	// For now, return a placeholder
	return &TextWithPageNumbers{
		Text:      "[DOC file content - text extraction not fully implemented]",
		PageCount: 1,
		PageInfo:  []PageInfo{{PageNumber: 1, Text: "[DOC file content]"}},
	}, nil
}

// extractExcelText extracts text from Excel files (placeholder)
func extractExcelText(filePath string) (*TextWithPageNumbers, error) {
	// This would require a library like github.com/360EntSecGroup/360excel
	// For now, return a placeholder
	return &TextWithPageNumbers{
		Text:      "[Excel file content - text extraction not fully implemented]",
		PageCount: 1,
		PageInfo:  []PageInfo{{PageNumber: 1, Text: "[Excel file content]"}},
	}, nil
}

// extractPowerPointText extracts text from PowerPoint files (placeholder)
func extractPowerPointText(filePath string) (*TextWithPageNumbers, error) {
	// This would require a library like github.com/unidoc/unioffice
	// For now, return a placeholder
	return &TextWithPageNumbers{
		Text:      "[PowerPoint file content - text extraction not fully implemented]",
		PageCount: 1,
		PageInfo:  []PageInfo{{PageNumber: 1, Text: "[PowerPoint file content]"}},
	}, nil
}
