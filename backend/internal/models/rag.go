package models

import (
	"time"
)

// RAGRequest represents a request to the RAG system
type RAGRequest struct {
	Question  string `json:"question" binding:"required"`
	TopK      int    `json:"top_k,omitempty"`      // Number of chunks to retrieve, defaults to 5
	SessionID string `json:"session_id,omitempty"` // Optional session ID for conversation context
}

// RAGResponse represents the response from the RAG system
type RAGResponse struct {
	Answer       string        `json:"answer"`
	Sources      []Source      `json:"sources"`
	Question     string        `json:"question"`
	ResponseTime time.Duration `json:"response_time"`
	Timestamp    time.Time     `json:"timestamp"`
	SessionID    string        `json:"session_id,omitempty"`
}

// Source represents a source document/chunk used in the answer
type Source struct {
	ChunkID      uint    `json:"chunk_id"`
	DocumentID   uint    `json:"document_id"`
	Content      string  `json:"content"`
	Score        float32 `json:"score"`
	PageNumber   *int    `json:"page_number,omitempty"`
	DocumentName string  `json:"document_name"`
}

// RAGPipeline represents the RAG pipeline configuration
type RAGPipeline struct {
	EmbeddingDimensions int     `json:"embedding_dimensions"`
	RetrievalTopK       int     `json:"retrieval_top_k"`
	SimilarityThreshold float32 `json:"similarity_threshold"`
	MaxContextLength    int     `json:"max_context_length"`
	LLMModel            string  `json:"llm_model"`
}

// RAGStats represents statistics about the RAG system
type RAGStats struct {
	TotalQuestions       int64         `json:"total_questions"`
	AverageResponseTime  time.Duration `json:"average_response_time"`
	TotalDocuments       int           `json:"total_documents"`
	TotalChunks          int           `json:"total_chunks"`
	AverageRetrievalTime time.Duration `json:"average_retrieval_time"`
	CacheHitRate         float64       `json:"cache_hit_rate"`
}

// PromptTemplate represents a template for building LLM prompts
type PromptTemplate struct {
	SystemPrompt string `json:"system_prompt"`
	UserPrompt   string `json:"user_prompt"`
}

// DefaultPromptTemplate returns the default prompt template
func DefaultPromptTemplate() *PromptTemplate {
	return &PromptTemplate{
		SystemPrompt: `You are a helpful AI assistant that answers questions based on the provided context. 
Use only the information from the provided context to answer the question. 
If the context doesn't contain enough information to answer the question, say "I don't have enough information to answer this question based on the provided context."
Do not make up information or use external knowledge.
Be concise and accurate in your responses.`,
		UserPrompt: `Context:
{{.Context}}

Question: {{.Question}}

Answer based on the context above:`,
	}
}

// ConversationContext represents conversation history for context-aware responses
type ConversationContext struct {
	SessionID string    `json:"session_id"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens"`
}

// HealthCheck represents the health check response
type HealthCheck struct {
	Status    string                   `json:"status"`
	Timestamp time.Time                `json:"timestamp"`
	Services  map[string]ServiceHealth `json:"services"`
	Version   string                   `json:"version"`
	Uptime    time.Duration            `json:"uptime"`
}

// ServiceHealth represents the health of a specific service
type ServiceHealth struct {
	Status  string        `json:"status"`
	Latency time.Duration `json:"latency"`
	Error   string        `json:"error,omitempty"`
}
