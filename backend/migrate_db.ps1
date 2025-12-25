Write-Host "Running RAGify database migration..." -ForegroundColor Green

# Prompt for PostgreSQL password
$pgPassword = Read-Host -Prompt "Enter PostgreSQL password for user 'postgres'" -AsSecureString
$pgPasswordPlainText = [Runtime.InteropServices.Marshal]::PtrToStringBSTR([Runtime.InteropServices.Marshal]::SecureStringToBSTR($pgPassword))

# Update the .env file with the provided password temporarily
$envContent = Get-Content .env
$originalContent = $envContent
$envContent = $envContent -replace 'POSTGRES_PASSWORD=.*', "POSTGRES_PASSWORD=$pgPasswordPlainText"
$envContent | Set-Content .env

Write-Host "Updated .env file with new password for migration." -ForegroundColor Green

# Run the Go migration
Write-Host "Running database migration..." -ForegroundColor Yellow
go run migrate.go

if ($LASTEXITCODE -eq 0) {
    Write-Host "Database migration completed successfully!" -ForegroundColor Green
    Write-Host "You can now run the application with: go run cmd/server/main.go" -ForegroundColor Cyan
} else {
    Write-Host "Database migration failed!" -ForegroundColor Red
    exit 1
}

# Restore original .env content (optional - you might want to keep the updated password)
# Uncomment the next two lines if you want to revert the .env file after migration
# Write-Host "Restoring original .env file..." -ForegroundColor Yellow
# $originalContent | Set-Content .env

Write-Host "Press any key to continue..." -ForegroundColor Gray
$null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")