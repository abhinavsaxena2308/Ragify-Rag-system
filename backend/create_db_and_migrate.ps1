Write-Host "Setting up RAGify database and running migrations..." -ForegroundColor Green

# Prompt for PostgreSQL password
$pgPassword = Read-Host -Prompt "Enter PostgreSQL password for user 'postgres'" -AsSecureString
$pgPasswordPlainText = [Runtime.InteropServices.Marshal]::PtrToStringBSTR([Runtime.InteropServices.Marshal]::SecureStringToBSTR($pgPassword))

# Update the .env file with the provided password
$envContent = Get-Content .env
$envContent = $envContent -replace 'POSTGRES_PASSWORD=.*', "POSTGRES_PASSWORD=$pgPasswordPlainText"
$envContent | Set-Content .env

Write-Host "Updated .env file with new password." -ForegroundColor Green

Write-Host "Attempting to create database 'ragify'..." -ForegroundColor Yellow

# Try to create the database
$dbCreationScript = "SELECT 'CREATE DATABASE ragify' WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'ragify')\\gexec"
$creationResult = psql -h localhost -p 5432 -U postgres -d postgres -c $dbCreationScript 2>$null

if ($LASTEXITCODE -eq 0) {
    Write-Host "Database 'ragify' created or already exists." -ForegroundColor Green
} else {
    Write-Host "Failed to create database. Please ensure PostgreSQL is running and credentials are correct." -ForegroundColor Red
    Write-Host "You may need to create the database manually first." -ForegroundColor Yellow
    pause
    exit 1
}

Write-Host "Running database migrations..." -ForegroundColor Yellow

# Run the Go migration
$goModTidyResult = go mod tidy
if ($LASTEXITCODE -ne 0) {
    Write-Host "Error running 'go mod tidy'" -ForegroundColor Red
    exit 1
}

# Attempt to run the migration
$migrationResult = go run migrate.go 2>&1

if ($LASTEXITCODE -eq 0) {
    Write-Host "Database migration completed successfully!" -ForegroundColor Green
    Write-Host "You can now run the application with: go run cmd/server/main.go" -ForegroundColor Cyan
} else {
    Write-Host "Database migration failed!" -ForegroundColor Red
    Write-Host $migrationResult
    Write-Host ""
    Write-Host "Possible issues:" -ForegroundColor Yellow
    Write-Host "1. PostgreSQL may not be running" -ForegroundColor Yellow
    Write-Host "2. Database 'ragify' may not exist" -ForegroundColor Yellow
    Write-Host "3. Credentials in .env may be incorrect" -ForegroundColor Yellow
    Write-Host "4. Required PostgreSQL extensions may not be installed" -ForegroundColor Yellow
    exit 1
}

Write-Host "Press any key to continue..." -ForegroundColor Gray
$null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")