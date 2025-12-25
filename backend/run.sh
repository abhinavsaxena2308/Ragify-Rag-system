#!/bin/bash

# Script to run the RAGify backend server
# Make sure PostgreSQL is running and the database is created before running this

echo "Starting RAGify Backend Server..."

# Build the application
go build -o bin/server cmd/server/main.go

if [ $? -eq 0 ]; then
    echo "Build successful!"
    echo "Make sure PostgreSQL is running and the 'ragify' database exists."
    echo "Run: go run cmd/server/main.go"
    go run cmd/server/main.go
else
    echo "Build failed!"
    exit 1
fi