#!/bin/bash

# Script to set up PostgreSQL database for RAGify using terminal commands
# This script requires psql to be available in your PATH

echo "Setting up RAGify PostgreSQL database..."

# Function to prompt for password securely
read_password() {
    read -s -p "Enter PostgreSQL password for user 'postgres': " pg_password
    echo
}

# Check if psql is available
if ! command -v psql &> /dev/null; then
    echo "Error: psql is not installed or not in PATH"
    echo "Please install PostgreSQL client tools first"
    exit 1
fi

# Prompt for password
read_password

# Update .env file with the password
if [ -f .env ]; then
    sed -i "s/POSTGRES_PASSWORD=.*/POSTGRES_PASSWORD=$pg_password/" .env
    echo "Updated .env file with new password"
else
    echo "Error: .env file not found"
    exit 1
fi

# Create the database using psql
echo "Creating ragify database..."
CREATE_DB_QUERY="SELECT 'CREATE DATABASE ragify' WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'ragify')\\gexec"
psql -h localhost -p 5432 -U postgres -d postgres -c "$CREATE_DB_QUERY"

if [ $? -eq 0 ]; then
    echo "Database 'ragify' created or already exists."
    echo "Creating tables in ragify database..."
    
    # Run the database setup script
    psql -h localhost -p 5432 -U postgres -d ragify -f database_setup.sql
    
    if [ $? -eq 0 ]; then
        echo "Database setup completed successfully!"
        echo "You can now run the application with: go run cmd/server/main.go"
    else
        echo "Error creating tables!"
        exit 1
    fi
else
    echo "Error creating database!"
    echo "Make sure PostgreSQL is running and the credentials are correct."
    exit 1
fi

echo "Setup completed successfully!"