package services

import (
	"fmt"
	"regexp"
	"strings"

	"ragify-backend/internal/models"
	"ragify-backend/internal/utils"
)

// ChunkingConfig holds configuration for text chunking
type ChunkingConfig struct {
	ChunkSizeTokens int    `json:"chunk_size_tokens"`
	OverlapTokens   int    `json:"overlap_tokens"`
	MinChunkSize    int    `json:"min_chunk_size"`
	MaxChunkSize    int    `json:"max_chunk_size"`
	Separator       string `json:"separator"`
	Model           string `json:"model"`
}

// DefaultChunkingConfig returns default configuration for chunking
func DefaultChunkingConfig() ChunkingConfig {
	return ChunkingConfig{
		ChunkSizeTokens: 400,
		OverlapTokens:   50,
		MinChunkSize:    50,
		MaxChunkSize:    800,
		Separator:       "\n\n",
		Model:           "llama2",
	}
}

// ChunkingService handles text chunking for RAG
type ChunkingService struct {
	config       ChunkingConfig
	tokenCounter *utils.TokenCounter
}

// NewChunkingService creates a new ChunkingService instance
func NewChunkingService(config ChunkingConfig) *ChunkingService {
	if config.ChunkSizeTokens == 0 {
		config = DefaultChunkingConfig()
	}

	return &ChunkingService{
		config:       config,
		tokenCounter: utils.NewTokenCounter(),
	}
}

// ChunkTextWithPageInfo chunks text with page information
func (cs *ChunkingService) ChunkTextWithPageInfo(textWithPages *utils.TextWithPageNumbers, documentID uint) []models.Chunk {
	var chunks []models.Chunk

	for _, pageInfo := range textWithPages.PageInfo {
		pageChunks := cs.chunkSinglePage(pageInfo.Text, documentID, pageInfo.PageNumber)
		chunks = append(chunks, pageChunks...)
	}

	// Re-index chunks globally
	for i := range chunks {
		chunks[i].ChunkIndex = i
	}

	return chunks
}

// chunkSinglePage chunks text from a single page
func (cs *ChunkingService) chunkSinglePage(text string, documentID uint, pageNumber int) []models.Chunk {
	// Clean and normalize the text
	cleanText := cs.cleanText(text)

	if cleanText == "" {
		return []models.Chunk{}
	}

	// Split text into paragraphs for better semantic chunks
	paragraphs := cs.splitIntoParagraphs(cleanText)

	var chunks []models.Chunk
	var currentChunk strings.Builder
	var currentTokens int

	for _, paragraph := range paragraphs {
		paragraphTokens := cs.tokenCounter.CountTokens(paragraph)

		// If adding this paragraph exceeds chunk size, create a new chunk
		if currentTokens+paragraphTokens > cs.config.ChunkSizeTokens && currentChunk.Len() > 0 {
			// Create chunk from accumulated content
			chunkContent := currentChunk.String()
			if cs.tokenCounter.CountTokens(chunkContent) >= cs.config.MinChunkSize {
				chunk := models.Chunk{
					DocumentID: documentID,
					Content:    chunkContent,
					PageNumber: &pageNumber,
					ChunkIndex: len(chunks), // Will be updated later
				}
				chunks = append(chunks, chunk)
			}

			// Start new chunk with overlap
			currentChunk.Reset()
			currentTokens = 0

			// Add overlap from previous content
			overlapText := cs.getOverlapText(chunkContent)
			if overlapText != "" {
				currentChunk.WriteString(overlapText)
				currentTokens = cs.tokenCounter.CountTokens(overlapText)
			}
		}

		// Add current paragraph
		if currentChunk.Len() > 0 {
			currentChunk.WriteString(cs.config.Separator)
		}
		currentChunk.WriteString(paragraph)
		currentTokens += paragraphTokens
	}

	// Don't forget the last chunk
	if currentChunk.Len() > 0 {
		chunkContent := currentChunk.String()
		if cs.tokenCounter.CountTokens(chunkContent) >= cs.config.MinChunkSize {
			chunk := models.Chunk{
				DocumentID: documentID,
				Content:    chunkContent,
				PageNumber: &pageNumber,
				ChunkIndex: len(chunks), // Will be updated later
			}
			chunks = append(chunks, chunk)
		}
	}

	return chunks
}

// cleanText normalizes and cleans text for chunking
func (cs *ChunkingService) cleanText(text string) string {
	// Remove excessive whitespace
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")

	// Remove control characters but keep newlines and tabs
	text = regexp.MustCompile(`[\x00-\x08\x0B\x0C\x0E-\x1F\x7F]`).ReplaceAllString(text, "")

	// Normalize newlines
	text = regexp.MustCompile(`\r\n|\r`).ReplaceAllString(text, "\n")

	// Remove excessive newlines (more than 2 consecutive)
	text = regexp.MustCompile(`\n{3,}`).ReplaceAllString(text, "\n\n")

	// Trim leading/trailing whitespace
	text = strings.TrimSpace(text)

	return text
}

// splitIntoParagraphs splits text into semantic paragraphs
func (cs *ChunkingService) splitIntoParagraphs(text string) []string {
	// Split by double newlines first (paragraph boundaries)
	paragraphs := regexp.MustCompile(`\n\s*\n`).Split(text, -1)

	var result []string
	for _, paragraph := range paragraphs {
		paragraph = strings.TrimSpace(paragraph)
		if paragraph == "" {
			continue
		}

		// If paragraph is too long, split it further by sentences
		if cs.tokenCounter.CountTokens(paragraph) > cs.config.ChunkSizeTokens {
			sentences := cs.splitIntoSentences(paragraph)
			result = append(result, sentences...)
		} else {
			result = append(result, paragraph)
		}
	}

	return result
}

// splitIntoSentences splits text into sentences
func (cs *ChunkingService) splitIntoSentences(text string) []string {
	// Simple sentence splitting - can be improved with NLP libraries
	sentences := regexp.MustCompile(`[.!?]+\s+`).Split(text, -1)

	var result []string
	for _, sentence := range sentences {
		sentence = strings.TrimSpace(sentence)
		if sentence != "" {
			result = append(result, sentence)
		}
	}

	return result
}

// getOverlapText extracts overlapping text for chunk continuity
func (cs *ChunkingService) getOverlapText(chunkContent string) string {
	if cs.config.OverlapTokens <= 0 {
		return ""
	}

	words := strings.Fields(chunkContent)
	if len(words) == 0 {
		return ""
	}

	// Work backwards to find overlap that doesn't exceed token limit
	var overlapWords []string

	// Start from the end and work backwards
	for i := len(words) - 1; i >= 0; i-- {
		testWords := words[i:]
		testText := strings.Join(testWords, " ")
		testTokens := cs.tokenCounter.CountTokens(testText)

		if testTokens <= cs.config.OverlapTokens {
			overlapWords = testWords
		} else {
			break
		}
	}

	if len(overlapWords) > 0 {
		return strings.Join(overlapWords, " ")
	}

	return ""
}

// ChunkDocument processes a complete document and returns chunks
func (cs *ChunkingService) ChunkDocument(document *models.Document) ([]models.Chunk, error) {
	if document == nil {
		return nil, fmt.Errorf("document cannot be nil")
	}

	if document.TextContent == "" {
		return nil, fmt.Errorf("document has no text content")
	}

	// Create TextWithPageNumbers structure
	textWithPages := &utils.TextWithPageNumbers{
		Text:      document.TextContent,
		PageCount: document.PageCount,
		PageInfo:  []utils.PageInfo{},
	}

	// If we have page information in the text, extract it
	if strings.Contains(document.TextContent, "--- Page ") {
		textWithPages.PageInfo = cs.extractPageInfo(document.TextContent)
	} else {
		// Treat entire document as single page
		textWithPages.PageInfo = []utils.PageInfo{
			{
				PageNumber: 1,
				Text:       document.TextContent,
			},
		}
	}

	// Chunk the text
	chunks := cs.ChunkTextWithPageInfo(textWithPages, document.ID)

	return chunks, nil
}

// extractPageInfo extracts page information from text with page markers
func (cs *ChunkingService) extractPageInfo(text string) []utils.PageInfo {
	var pageInfo []utils.PageInfo

	// Split by page markers
	pagePattern := regexp.MustCompile(`--- Page (\d+) ---\s*`)

	matches := pagePattern.FindAllStringSubmatchIndex(text, -1)

	for i, match := range matches {
		if len(match) < 4 {
			continue
		}

		pageStart := match[1] // End of page marker
		var pageEnd int

		if i+1 < len(matches) {
			pageEnd = matches[i+1][0] // Start of next page marker
		} else {
			pageEnd = len(text) // End of document
		}

		pageText := text[pageStart:pageEnd]
		pageText = strings.TrimSpace(pageText)

		if pageText != "" {
			pageInfo = append(pageInfo, utils.PageInfo{
				PageNumber: i + 1, // Use sequential numbering
				Text:       pageText,
			})
		}
	}

	return pageInfo
}

// GetChunkingStats returns statistics about the chunking process
func (cs *ChunkingService) GetChunkingStats(chunks []models.Chunk) map[string]interface{} {
	if len(chunks) == 0 {
		return map[string]interface{}{
			"total_chunks":   0,
			"total_tokens":   0,
			"avg_chunk_size": 0,
			"min_chunk_size": 0,
			"max_chunk_size": 0,
		}
	}

	var totalTokens int
	var minTokens = int(^uint(0) >> 1) // Max int
	var maxTokens int

	for _, chunk := range chunks {
		tokens := cs.tokenCounter.CountTokens(chunk.Content)
		totalTokens += tokens

		if tokens < minTokens {
			minTokens = tokens
		}
		if tokens > maxTokens {
			maxTokens = tokens
		}
	}

	avgTokens := totalTokens / len(chunks)

	return map[string]interface{}{
		"total_chunks":   len(chunks),
		"total_tokens":   totalTokens,
		"avg_chunk_size": avgTokens,
		"min_chunk_size": minTokens,
		"max_chunk_size": maxTokens,
		"config":         cs.config,
	}
}

// ValidateChunks checks if chunks meet quality criteria
func (cs *ChunkingService) ValidateChunks(chunks []models.Chunk) []string {
	var issues []string

	for i, chunk := range chunks {
		tokens := cs.tokenCounter.CountTokens(chunk.Content)

		if tokens < cs.config.MinChunkSize {
			issues = append(issues, fmt.Sprintf("Chunk %d: Too small (%d tokens < %d min)", i, tokens, cs.config.MinChunkSize))
		}

		if tokens > cs.config.MaxChunkSize {
			issues = append(issues, fmt.Sprintf("Chunk %d: Too large (%d tokens > %d max)", i, tokens, cs.config.MaxChunkSize))
		}

		if strings.TrimSpace(chunk.Content) == "" {
			issues = append(issues, fmt.Sprintf("Chunk %d: Empty content", i))
		}
	}

	return issues
}
