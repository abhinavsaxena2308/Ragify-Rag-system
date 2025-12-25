package services

import (
	"strings"
	"testing"

	"ragify-backend/internal/models"
	"ragify-backend/internal/utils"
)

func TestTokenCounter(t *testing.T) {
	tc := utils.NewTokenCounter()

	// Test basic token counting
	text := "This is a simple test sentence with several words."
	tokens := tc.CountTokens(text)

	if tokens <= 0 {
		t.Errorf("Expected positive token count, got %d", tokens)
	}

	// Test empty text
	emptyTokens := tc.CountTokens("")
	if emptyTokens != 0 {
		t.Errorf("Expected 0 tokens for empty text, got %d", emptyTokens)
	}

	// Test token stats
	stats := tc.GetTokenStats(text)
	if stats["token_count"] == nil {
		t.Error("Expected token_count in stats")
	}
}

func TestChunkingService(t *testing.T) {
	config := DefaultChunkingConfig()
	cs := NewChunkingService(config)

	// Test with simple text
	text := `This is the first paragraph. It contains multiple sentences to test the chunking functionality.
	
	This is the second paragraph. It should be separated from the first one.
	
	This is the third paragraph. It contains more content to ensure we have enough text for proper chunking.`

	textWithPages := &utils.TextWithPageNumbers{
		Text:      text,
		PageCount: 1,
		PageInfo: []utils.PageInfo{
			{
				PageNumber: 1,
				Text:       text,
			},
		},
	}

	chunks := cs.ChunkTextWithPageInfo(textWithPages, 1)

	if len(chunks) == 0 {
		t.Error("Expected at least one chunk")
	}

	// Verify chunk metadata
	for i, chunk := range chunks {
		if chunk.DocumentID != 1 {
			t.Errorf("Chunk %d: Expected DocumentID 1, got %d", i, chunk.DocumentID)
		}

		if chunk.PageNumber == nil || *chunk.PageNumber != 1 {
			t.Errorf("Chunk %d: Expected PageNumber 1, got %v", i, chunk.PageNumber)
		}

		if chunk.ChunkIndex != i {
			t.Errorf("Chunk %d: Expected ChunkIndex %d, got %d", i, i, chunk.ChunkIndex)
		}
	}
}

func TestChunkDocument(t *testing.T) {
	config := DefaultChunkingConfig()
	cs := NewChunkingService(config)

	// Create a test document
	document := &models.Document{
		ID:       1,
		Filename: "test.txt",
		TextContent: `This is a test document. It contains multiple paragraphs to test the chunking functionality.
		
		The first paragraph contains some introductory content.
		
		The second paragraph provides more detailed information about the topic being discussed.
		
		The third paragraph concludes the document with a summary of the key points.`,
		PageCount: 1,
	}

	chunks, err := cs.ChunkDocument(document)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(chunks) == 0 {
		t.Error("Expected at least one chunk")
	}

	// Test chunk validation
	issues := cs.ValidateChunks(chunks)
	if len(issues) > 0 {
		t.Errorf("Chunk validation issues: %v", issues)
	}

	// Test chunking stats
	stats := cs.GetChunkingStats(chunks)
	if stats["total_chunks"] == nil {
		t.Error("Expected total_chunks in stats")
	}
}

func TestTextCleaning(t *testing.T) {
	config := DefaultChunkingConfig()
	cs := NewChunkingService(config)

	// Test text with various whitespace issues
	dirtyText := "  This   text   has   \t\n  multiple   \n\n   whitespace   issues.  "
	cleanText := cs.cleanText(dirtyText)

	expected := "This text has multiple whitespace issues."
	if cleanText != expected {
		t.Errorf("Expected '%s', got '%s'", expected, cleanText)
	}
}

func TestParagraphSplitting(t *testing.T) {
	config := DefaultChunkingConfig()
	cs := NewChunkingService(config)

	text := `First paragraph.
	
	Second paragraph.
	
	Third paragraph with more content.`

	paragraphs := cs.splitIntoParagraphs(text)

	if len(paragraphs) != 3 {
		t.Errorf("Expected 3 paragraphs, got %d", len(paragraphs))
	}

	for i, paragraph := range paragraphs {
		if paragraph == "" {
			t.Errorf("Paragraph %d is empty", i)
		}
	}
}

func TestOverlapText(t *testing.T) {
	config := ChunkingConfig{
		ChunkSizeTokens: 10,
		OverlapTokens:   3,
	}
	cs := NewChunkingService(config)

	text := "This is a test sentence for overlap functionality."
	overlap := cs.getOverlapText(text)

	if overlap == "" {
		t.Error("Expected non-empty overlap text")
	}

	// Check that overlap doesn't exceed token limit
	tc := utils.NewTokenCounter()
	overlapTokens := tc.CountTokens(overlap)
	if overlapTokens > config.OverlapTokens {
		t.Errorf("Overlap %d tokens exceeds limit %d", overlapTokens, config.OverlapTokens)
	}
}

func BenchmarkTokenCounter(b *testing.B) {
	tc := utils.NewTokenCounter()
	baseText := `This is a benchmark test for the token counter. It contains multiple sentences and words to test performance. 
	The text should be long enough to provide meaningful benchmark results. We repeat this text several times to simulate real-world usage.
	`
	text := baseText + strings.Repeat(baseText, 10)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tc.CountTokens(text)
	}
}

func BenchmarkChunking(b *testing.B) {
	config := DefaultChunkingConfig()
	cs := NewChunkingService(config)

	// Create large text for benchmarking
	text := strings.Repeat(`This is a test paragraph for benchmarking the chunking service. It contains enough text to create multiple chunks and test the performance of the chunking algorithm. The text should be representative of real-world documents that would be processed by the system. `, 100)

	textWithPages := &utils.TextWithPageNumbers{
		Text:      text,
		PageCount: 1,
		PageInfo: []utils.PageInfo{
			{
				PageNumber: 1,
				Text:       text,
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cs.ChunkTextWithPageInfo(textWithPages, 1)
	}
}
