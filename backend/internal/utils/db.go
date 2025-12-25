package utils

import (
	"fmt"
	"ragify-backend/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DBConnect creates a connection to the database
func DBConnect(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.PostgreSQL.Host,
		cfg.PostgreSQL.User,
		cfg.PostgreSQL.Password,
		cfg.PostgreSQL.DBName,
		cfg.PostgreSQL.Port,
		cfg.PostgreSQL.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
