// Legacy API service for RAGify backend
// This file is kept for backward compatibility
// New implementations should use the individual service files:
// - uploadService.js
// - chatService.js

export { default as documentService } from './uploadService';
export { default as chatService } from './chatService';
