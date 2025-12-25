# Quick setup script for RAGify database
# This script demonstrates the database setup process

Write-Host "Setting up RAGify database..." -ForegroundColor Green

# For demonstration purposes, we'll show what would happen
Write-Host "1. Updating .env file with PostgreSQL credentials..." -ForegroundColor Yellow

# Read the current .env file
if (Test-Path ".env") {
    $envContent = Get-Content .env
    
    # Update the password line (you would replace 'your_actual_password' with the real password)
    $envContent = $envContent -replace 'POSTGRES_PASSWORD=.*', 'POSTGRES_PASSWORD=your_actual_postgres_password'
    
    # Write back to the file
    $envContent | Set-Content .env
    
    Write-Host "   .env file updated successfully!" -ForegroundColor Green
} else {
    Write-Host "   Error: .env file not found!" -ForegroundColor Red
    exit 1
}

Write-Host "2. Checking if PostgreSQL is accessible..." -ForegroundColor Yellow

# Test database connection by checking if psql is available
$psqlAvailable = Get-Command psql -ErrorAction SilentlyContinue
if ($psqlAvailable) {
    Write-Host "   psql command is available" -ForegroundColor Green
    
    # In a real scenario, you would run:
    # psql -h localhost -p 5432 -U postgres -d postgres -c "SELECT version();"
    Write-Host "   (In real scenario, would test PostgreSQL connection)" -ForegroundColor Cyan
} else {
    Write-Host "   psql command is NOT available - please install PostgreSQL client tools" -ForegroundColor Red
    exit 1
}

Write-Host "3. Creating ragify database (if it doesn't exist)..." -ForegroundColor Yellow

# In a real scenario, you would run:
# psql -h localhost -p 5432 -U postgres -d postgres -c "CREATE DATABASE ragify WITH OWNER = postgres ENCODING = 'UTF8'..."

Write-Host "   (In real scenario, would create 'ragify' database)" -ForegroundColor Cyan

Write-Host "4. Creating database tables..." -ForegroundColor Yellow

# In a real scenario, you would run:
# psql -h localhost -p 5432 -U postgres -d ragify -f database_setup.sql

Write-Host "   (In real scenario, would create tables using database_setup.sql)" -ForegroundColor Cyan

Write-Host ""
Write-Host "Database setup process completed!" -ForegroundColor Green
Write-Host ""
Write-Host "To perform the actual setup, you need to:" -ForegroundColor Cyan
Write-Host "1. Replace 'your_actual_postgres_password' in .env with your real PostgreSQL password" -ForegroundColor Cyan
Write-Host "2. Run the setup_postgres.bat script (Windows) or setup_postgres.sh (Linux/Mac)" -ForegroundColor Cyan
Write-Host "3. Enter your PostgreSQL password when prompted" -ForegroundColor Cyan
Write-Host ""
Write-Host "All required files have been created and are ready for database setup!" -ForegroundColor Green