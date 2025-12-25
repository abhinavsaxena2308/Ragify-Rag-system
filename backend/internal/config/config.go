package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
	PostgreSQL  PostgreSQLConfig
	FAISS       FAISSConfig
	LLM         LLMConfig
	RAG         RAGConfig
}

type PostgreSQLConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type FAISSConfig struct {
	IndexPath string
}

type LLMConfig struct {
	Provider string
	Endpoint string
	APIKey   string
	Model    string
}

type RAGConfig struct {
	RetrievalTopK       int
	SimilarityThreshold float32
	MaxContextLength    int
	EmbeddingDimensions int
}

func LoadConfig() *Config {
	// Load .env file if it exists
	godotenv.Load()

	return &Config{
		Port:        getEnvOrDefault("PORT", "8080"),
		DatabaseURL: getEnvOrDefault("DATABASE_URL", ""),
		JWTSecret:   getEnvOrDefault("JWT_SECRET", "default_secret_key_change_in_production"),
		PostgreSQL: PostgreSQLConfig{
			Host:     getEnvOrDefault("POSTGRES_HOST", "localhost"),
			Port:     getEnvOrDefault("POSTGRES_PORT", "5432"),
			User:     getEnvOrDefault("POSTGRES_USER", "postgres"),
			Password: getEnvOrDefault("POSTGRES_PASSWORD", "postgres"),
			DBName:   getEnvOrDefault("POSTGRES_DB", "ragify"),
			SSLMode:  getEnvOrDefault("POSTGRES_SSLMODE", "disable"),
		},
		FAISS: FAISSConfig{
			IndexPath: getEnvOrDefault("FAISS_INDEX_PATH", "./data/faiss_index"),
		},
		LLM: LLMConfig{
			Provider: getEnvOrDefault("LLM_PROVIDER", "ollama"),
			Endpoint: getEnvOrDefault("LLM_ENDPOINT", "http://localhost:11434"),
			APIKey:   getEnvOrDefault("LLM_API_KEY", ""),
			Model:    getEnvOrDefault("LLM_MODEL", "llama2"),
		},
		RAG: RAGConfig{
			RetrievalTopK:       getEnvOrDefaultInt("RAG_TOP_K", 5),
			SimilarityThreshold: getEnvOrDefaultFloat32("RAG_SIMILARITY_THRESHOLD", 0.0),
			MaxContextLength:    getEnvOrDefaultInt("RAG_MAX_CONTEXT_LENGTH", 4000),
			EmbeddingDimensions: getEnvOrDefaultInt("RAG_EMBEDDING_DIMENSIONS", 1536),
		},
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvOrDefaultInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		// Parse int from string
		if intValue := parseInt(value); intValue != 0 {
			return intValue
		}
	}
	return defaultValue
}

func getEnvOrDefaultFloat32(key string, defaultValue float32) float32 {
	if value := os.Getenv(key); value != "" {
		// Parse float32 from string
		if floatValue := parseFloat32(value); floatValue != 0 {
			return floatValue
		}
	}
	return defaultValue
}

func parseInt(s string) int {
	var result int
	fmt.Sscanf(s, "%d", &result)
	return result
}

func parseFloat32(s string) float32 {
	var result float32
	fmt.Sscanf(s, "%f", &result)
	return result
}
