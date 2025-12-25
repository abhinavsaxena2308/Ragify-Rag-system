@echo off
echo Setting up RAGify database...

REM Prompt for PostgreSQL password
set /p PG_PASSWORD=Enter PostgreSQL password for user 'postgres': 

REM Update the .env file with the provided password
powershell -Command "(Get-Content .env) -replace 'POSTGRES_PASSWORD=.*', 'POSTGRES_PASSWORD=%PG_PASSWORD%' | Set-Content .env"

echo Updated .env file with new password.

REM Create the database using psql
echo Creating ragify database...
echo CREATE DATABASE ragify WITH OWNER = postgres ENCODING = 'UTF8' LC_COLLATE = 'en_US.UTF-8' LC_CTYPE = 'en_US.UTF-8' TABLESPACE = pg_default CONNECTION LIMIT = -1; > create_db.sql

psql -h localhost -p 5432 -U postgres -d postgres -f create_db.sql

if %ERRORLEVEL% EQU 0 (
    echo Database created successfully!
    echo Creating tables...
    
    REM Run the database setup script
    psql -h localhost -p 5432 -U postgres -d ragify -f database_setup.sql
    
    if %ERRORLEVEL% EQU 0 (
        echo Database setup completed successfully!
        echo You can now run the application with: go run cmd/server/main.go
    ) else (
        echo Error creating tables!
        exit /b 1
    )
) else (
    echo Error creating database!
    echo Make sure PostgreSQL is running and the credentials are correct.
    exit /b 1
)

REM Clean up temporary file
del create_db.sql

pause