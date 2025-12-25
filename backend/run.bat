@echo off
echo Starting RAGify Backend Server...

REM Build the application
go build -o bin/server.exe cmd/server/main.go

if %ERRORLEVEL% EQU 0 (
    echo Build successful!
    echo Make sure PostgreSQL is running and the 'ragify' database exists.
    echo Starting server...
    go run cmd/server/main.go
) else (
    echo Build failed!
    exit /b 1
)