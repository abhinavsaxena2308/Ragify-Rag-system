package utils

import (
	"regexp"
	"strings"
)

// TokenCounter provides methods for counting tokens in text
type TokenCounter struct{}

// NewTokenCounter creates a new TokenCounter instance
func NewTokenCounter() *TokenCounter {
	return &TokenCounter{}
}

// CountTokens estimates the number of tokens in text using a simple heuristic
// This approximation works well for English text and common tokenizers
func (tc *TokenCounter) CountTokens(text string) int {
	if text == "" {
		return 0
	}

	// Clean and normalize text first
	cleanText := tc.cleanText(text)

	// Split into tokens using whitespace and punctuation as delimiters
	tokens := tc.tokenizeText(cleanText)

	return len(tokens)
}

// cleanText normalizes text for accurate token counting
func (tc *TokenCounter) cleanText(text string) string {
	// Convert to lowercase for consistency
	text = strings.ToLower(text)

	// Remove extra whitespace
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")

	// Trim leading/trailing whitespace
	text = strings.TrimSpace(text)

	return text
}

// tokenizeText splits text into tokens using word boundaries
func (tc *TokenCounter) tokenizeText(text string) []string {
	// Use regex to find word-like sequences
	// This matches words, numbers, and common punctuation patterns
	pattern := regexp.MustCompile(`\b\w+\b|[^\w\s]`)

	matches := pattern.FindAllString(text, -1)

	// Filter out empty strings and very short tokens
	var tokens []string
	for _, match := range matches {
		match = strings.TrimSpace(match)
		if len(match) > 0 {
			tokens = append(tokens, match)
		}
	}

	return tokens
}

// CountTokensAdvanced provides a more accurate token count for specific models
func (tc *TokenCounter) CountTokensAdvanced(text string, model string) int {
	switch model {
	case "gpt-3.5-turbo", "gpt-4":
		return tc.countOpenAITokens(text)
	case "llama2", "llama3":
		return tc.countLlamaTokens(text)
	default:
		return tc.CountTokens(text) // fallback to simple counting
	}
}

// countOpenAITokens approximates OpenAI tokenization
func (tc *TokenCounter) countOpenAITokens(text string) int {
	// OpenAI's tokenizer is more complex, but this provides a reasonable approximation
	// Generally, OpenAI tokens are shorter than simple word tokens
	tokens := tc.CountTokens(text)

	// Apply a correction factor (OpenAI typically has ~1.3x more tokens than words)
	return int(float64(tokens) * 1.3)
}

// countLlamaTokens approximates Llama tokenization
func (tc *TokenCounter) countLlamaTokens(text string) int {
	// Llama tokenizer uses sentencepiece, this is an approximation
	tokens := tc.CountTokens(text)

	// Apply a correction factor for sentencepiece tokenization
	return int(float64(tokens) * 1.1)
}

// GetTokenStats provides statistics about tokenization
func (tc *TokenCounter) GetTokenStats(text string) map[string]interface{} {
	if text == "" {
		return map[string]interface{}{
			"char_count":       0,
			"word_count":       0,
			"token_count":      0,
			"avg_token_length": 0,
		}
	}

	cleanText := tc.cleanText(text)
	tokens := tc.tokenizeText(cleanText)

	var totalTokenLength int
	for _, token := range tokens {
		totalTokenLength += len(token)
	}

	avgLength := 0
	if len(tokens) > 0 {
		avgLength = totalTokenLength / len(tokens)
	}

	return map[string]interface{}{
		"char_count":       len(text),
		"clean_char_count": len(cleanText),
		"word_count":       len(strings.Fields(cleanText)),
		"token_count":      len(tokens),
		"avg_token_length": avgLength,
	}
}

// IsTokenLimit checks if text exceeds a token limit
func (tc *TokenCounter) IsTokenLimit(text string, limit int) bool {
	return tc.CountTokens(text) > limit
}

// TruncateToTokens truncates text to fit within a token limit
func (tc *TokenCounter) TruncateToTokens(text string, limit int) string {
	if tc.CountTokens(text) <= limit {
		return text
	}

	// Simple approach: split by words and add until limit is reached
	words := strings.Fields(text)
	var result strings.Builder

	for _, word := range words {
		testText := result.String() + " " + word
		if tc.CountTokens(testText) > limit {
			break
		}
		if result.Len() > 0 {
			result.WriteString(" ")
		}
		result.WriteString(word)
	}

	return result.String()
}
