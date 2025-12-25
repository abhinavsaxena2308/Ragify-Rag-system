Write-Host "Setting up RAGify database..." -ForegroundColor Green

# Prompt for PostgreSQL password
$pgPassword = Read-Host -Prompt "Enter PostgreSQL password for user 'postgres'" -AsSecureString
$pgPasswordPlainText = [Runtime.InteropServices.Marshal]::PtrToStringBSTR([Runtime.InteropServices.Marshal]::SecureStringToBSTR($pgPassword))

# Update the .env file with the provided password
$envContent = Get-Content .env
$envContent = $envContent -replace 'POSTGRES_PASSWORD=.*', "POSTGRES_PASSWORD=$pgPasswordPlainText"
$envContent | Set-Content .env

Write-Host "Updated .env file with new password." -ForegroundColor Green

# Create the database using psql
Write-Host "Creating ragify database..." -ForegroundColor Yellow

$dbCreationScript = "CREATE DATABASE ragify WITH OWNER = postgres ENCODING = 'UTF8' LC_COLLATE = 'en_US.UTF-8' LC_CTYPE = 'en_US.UTF-8' TABLESPACE = pg_default CONNECTION LIMIT = -1;"

# Write to temporary file
$dbCreationScript | Out-File -FilePath "create_db.sql" -Encoding UTF8

# Execute the SQL command to create the database
$result = psql -h localhost -p 5432 -U postgres -d postgres -f create_db.sql 2>&1

if ($LASTEXITCODE -eq 0) {
    Write-Host "Database created successfully!" -ForegroundColor Green
    Write-Host "Creating tables..." -ForegroundColor Yellow
    
    # Run the database setup script
    $tableResult = psql -h localhost -p 5432 -U postgres -d ragify -f database_setup.sql 2>&1
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "Database setup completed successfully!" -ForegroundColor Green
        Write-Host "You can now run the application with: go run cmd/server/main.go" -ForegroundColor Cyan
    } else {
        Write-Host "Error creating tables!" -ForegroundColor Red
        Write-Host $tableResult
        exit 1
    }
} else {
    Write-Host "Error creating database!" -ForegroundColor Red
    Write-Host "Make sure PostgreSQL is running and the credentials are correct." -ForegroundColor Yellow
    Write-Host $result
    exit 1
}

# Clean up temporary file
Remove-Item "create_db.sql"

Write-Host "Press any key to continue..." -ForegroundColor Gray
$null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")