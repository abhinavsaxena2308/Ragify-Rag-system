package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	DatabaseURL  string
	JWTSecret    string
	PostgreSQL   PostgreSQLConfig
	FAISS        FAISSConfig
	LLM          LLMConfig
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
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}