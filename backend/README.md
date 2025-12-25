# RAGify Backend

The backend service for the RAGify document question answering system, built with Go and Echo framework.

## Project Structure

```
backend/
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point
├── internal/
│   ├── config/                    # Configuration management
│   │   └── config.go
│   ├── handlers/                  # HTTP request handlers
│   │   ├── base.go
│   │   ├── document_handler.go
│   │   ├── chat_handler.go
│   │   └── user_handler.go
│   ├── models/                    # Data models
│   │   ├── document.go
│   │   ├── user.go
│   │   └── chat.go
│   ├── routes/                    # API route definitions
│   │   └── routes.go
│   ├── services/                  # Business logic (to be implemented)
│   ├── repositories/              # Database operations (to be implemented)
│   └── utils/                     # Utility functions
│       ├── response.go
│       └── db.go
├── .env                          # Environment variables
├── go.mod                        # Go module definition
├── go.sum                        # Go module checksums
└── README.md                     # This file
```

## Setup

1. **Prerequisites**
   - Go 1.19+
   - PostgreSQL (for document metadata)
   - Ollama (for local LLM) or access to OpenRouter

2. **Database Setup**
   - Option 1: Use the SQL script
     ```sql
     -- Connect to PostgreSQL as a superuser
     -- Run the commands in database_setup.sql
     ```
   - Option 2: Run the migration
     ```bash
     go run migrate.go
     ```

3. **Installation**
   ```bash
   # Clone the repository
   git clone <repository-url>
   cd ragify-backend

   # Install dependencies
   go mod tidy

   # Set up environment variables
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Environment Variables**
   - `PORT`: Server port (default: 8080)
   - `POSTGRES_HOST`: PostgreSQL host
   - `POSTGRES_PORT`: PostgreSQL port
   - `POSTGRES_USER`: PostgreSQL user
   - `POSTGRES_PASSWORD`: PostgreSQL password
   - `POSTGRES_DB`: PostgreSQL database name
   - `LLM_PROVIDER`: LLM provider (ollama, openrouter)
   - `LLM_ENDPOINT`: LLM API endpoint
   - `LLM_MODEL`: LLM model name
   - `JWT_SECRET`: JWT secret key

## Running the Application

```bash
# Run the server
go run cmd/server/main.go
```

The server will start on `http://localhost:8080` (or the port specified in your environment).

## API Endpoints

- `GET /health` - Health check endpoint
- `POST /api/v1/documents` - Upload a document
- `GET /api/v1/documents` - Get all documents
- `GET /api/v1/documents/:id` - Get a specific document
- `DELETE /api/v1/documents/:id` - Delete a document
- `POST /api/v1/chat/ask` - Ask a question about documents
- `POST /api/v1/chat/session` - Create a chat session
- `GET /api/v1/chat/session/:id` - Get a chat session
- `GET /api/v1/chat/session/:id/messages` - Get messages in a session
- `POST /api/v1/users` - Create a user
- `POST /api/v1/users/login` - User login

## Health Check

The application provides a health check endpoint at `/health` that returns:
```json
{
  "status": "OK",
  "app": "RAGify Backend"
}
```

## Features

- Document upload and management
- RAG-based question answering
- Chat session management
- User authentication
- Environment-based configuration
- Comprehensive error handling
- Structured logging
- CORS support

## Architecture

The backend follows a clean architecture pattern with:
- Handlers for HTTP request processing
- Services for business logic (to be implemented)
- Repositories for data access (to be implemented)
- Models for data structures
- Configuration management
- Utility functions