-- RAGify Database Setup Script
-- Run this script in pgAdmin4 to create the database and tables

-- Create the database (only if it doesn't exist)
-- Note: You'll need to connect to the 'postgres' database first to create a new database
CREATE DATABASE ragify
    WITH 
    OWNER = postgres
    ENCODING = 'UTF8'
    LC_COLLATE = 'en_US.UTF-8'
    LC_CTYPE = 'en_US.UTF-8'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1;

-- After creating the database, connect to the 'ragify' database and run the following commands:

-- Enable UUID extension if needed (uncomment if you plan to use UUIDs)
-- CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    name VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create documents table
CREATE TABLE documents (
    id SERIAL PRIMARY KEY,
    filename VARCHAR(255) NOT NULL,
    original_name VARCHAR(500) NOT NULL,
    size BIGINT NOT NULL,
    content_type VARCHAR(100) NOT NULL,
    file_path TEXT NOT NULL,
    text_content TEXT,
    page_count INTEGER NOT NULL DEFAULT 1,
    extracted_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    user_id INTEGER REFERENCES users(id) ON DELETE SET NULL
);

-- Create chunks table (for document content chunks)
CREATE TABLE chunks (
    id SERIAL PRIMARY KEY,
    document_id INTEGER NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    embedding BYTEA, -- Store embedding vectors as bytes
    page_number INTEGER,
    chunk_index INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create chat_sessions table
CREATE TABLE chat_sessions (
    id SERIAL PRIMARY KEY,
    session_id VARCHAR(255) UNIQUE NOT NULL,
    user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create messages table
CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    session_id VARCHAR(255) NOT NULL REFERENCES chat_sessions(session_id) ON DELETE CASCADE,
    role VARCHAR(20) NOT NULL, -- 'user' or 'assistant'
    content TEXT NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    source_docs TEXT[] -- Array of source document references
);

-- Create indexes for better performance
CREATE INDEX idx_documents_user_id ON documents(user_id);
CREATE INDEX idx_chunks_document_id ON chunks(document_id);
CREATE INDEX idx_messages_session_id ON messages(session_id);
CREATE INDEX idx_messages_timestamp ON messages(timestamp);
CREATE INDEX idx_chat_sessions_user_id ON chat_sessions(user_id);

-- Insert sample data (optional)
-- Insert a sample user
INSERT INTO users (email, password, name) VALUES 
('admin@example.com', '$2a$10$8K1p/aW.20mx5ZHw8o47suZpM.ei5.Ye.gFQx0.NhC1Yz7j4YEWnG', 'Admin User') -- password: 'password' (hashed)
ON CONFLICT (email) DO NOTHING;

-- Insert a sample document
INSERT INTO documents (filename, original_name, size, content_type, file_path, user_id) VALUES
('sample.pdf', 'Sample Document', 102400, 'application/pdf', '/uploads/sample.pdf', 1)
ON CONFLICT (id) DO NOTHING;

-- Insert a sample chat session
INSERT INTO chat_sessions (session_id, user_id) VALUES
('session_12345', 1)
ON CONFLICT (session_id) DO NOTHING;

-- Insert a sample message
INSERT INTO messages (session_id, role, content, source_docs) VALUES
('session_12345', 'assistant', 'Hello! I am your RAGify assistant. How can I help you today?', ARRAY['sample.pdf'])
ON CONFLICT (id) DO NOTHING;

-- Grant permissions (adjust as needed for your security requirements)
GRANT ALL PRIVILEGES ON TABLE users TO postgres;
GRANT ALL PRIVILEGES ON TABLE documents TO postgres;
GRANT ALL PRIVILEGES ON TABLE chunks TO postgres;
GRANT ALL PRIVILEGES ON TABLE chat_sessions TO postgres;
GRANT ALL PRIVILEGES ON TABLE messages TO postgres;

GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO postgres;

-- Function to update the updated_at timestamp automatically
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers to automatically update the updated_at column
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_documents_updated_at BEFORE UPDATE ON documents FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_chat_sessions_updated_at BEFORE UPDATE ON chat_sessions FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- End of script