@echo off
echo Setting up RAGify PostgreSQL database using terminal commands...

REM Check if psql is available
where psql >nul 2>nul
if %ERRORLEVEL% neq 0 (
    echo Error: psql is not installed or not in PATH
    echo Please install PostgreSQL client tools first
    pause
    exit /b 1
)

REM Prompt for password
set /p PG_PASSWORD=Enter PostgreSQL password for user 'postgres': 

REM Update .env file with the password
powershell -Command "(Get-Content .env) -replace 'POSTGRES_PASSWORD=.*', 'POSTGRES_PASSWORD=%PG_PASSWORD%' | Set-Content .env"
if %ERRORLEVEL% neq 0 (
    echo Error updating .env file
    pause
    exit /b 1
)
echo Updated .env file with new password.

echo Creating ragify database...
echo SELECT \'CREATE DATABASE ragify\' WHERE NOT EXISTS \(SELECT FROM pg_database WHERE datname = \'ragify\'\)\\gexec > create_db_temp.sql

psql -h localhost -p 5432 -U postgres -d postgres -f create_db_temp.sql
if %ERRORLEVEL% equ 0 (
    echo Database 'ragify' created or already exists.
    echo Creating tables in ragify database...
    
    REM Run the database setup script
    psql -h localhost -p 5432 -U postgres -d ragify -f database_setup.sql
    if %ERRORLEVEL% equ 0 (
        echo Database setup completed successfully!
        echo You can now run the application with: go run cmd/server/main.go
    ) else (
        echo Error creating tables!
        del create_db_temp.sql
        pause
        exit /b 1
    )
) else (
    echo Error creating database!
    echo Make sure PostgreSQL is running and the credentials are correct.
    del create_db_temp.sql
    pause
    exit /b 1
)

REM Clean up
del create_db_temp.sql
echo Setup completed successfully!
pause