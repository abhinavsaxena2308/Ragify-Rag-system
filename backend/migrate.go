package main

import (
	"ragify-backend/internal/config"
	"ragify-backend/internal/models"
	"ragify-backend/internal/utils"
)

// This script can be used to migrate the database
// Run with: go run migrate.go
func main() {
	cfg := config.LoadConfig()

	db, err := utils.DBConnect(cfg)
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	// Migrate the schema
	err = db.AutoMigrate(
		&models.User{},
		&models.Document{},
		&models.Chunk{},
		&models.ChatSession{},
		&models.Message{},
	)
	if err != nil {
		panic("failed to migrate database: " + err.Error())
	}

	println("Database migrated successfully!")
}
