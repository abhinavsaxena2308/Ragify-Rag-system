// API service for RAGify backend
const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api';

// Document-related API calls
export const documentService = {
  // Upload a document
  uploadDocument: async (file) => {
    const formData = new FormData();
    formData.append('file', file);

    const response = await fetch(`${API_BASE_URL}/documents/upload`, {
      method: 'POST',
      body: formData,
    });

    if (!response.ok) {
      throw new Error(`Upload failed: ${response.statusText}`);
    }

    return response.json();
  },

  // Get all documents
  getDocuments: async () => {
    const response = await fetch(`${API_BASE_URL}/documents`);

    if (!response.ok) {
      throw new Error(`Failed to fetch documents: ${response.statusText}`);
    }

    return response.json();
  },

  // Delete a document
  deleteDocument: async (documentId) => {
    const response = await fetch(`${API_BASE_URL}/documents/${documentId}`, {
      method: 'DELETE',
    });

    if (!response.ok) {
      throw new Error(`Failed to delete document: ${response.statusText}`);
    }

    return response.json();
  },
};

// Chat/Q&A API calls
export const chatService = {
  // Send a question and get an answer
  askQuestion: async (question, documentIds = []) => {
    const response = await fetch(`${API_BASE_URL}/chat`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        question,
        documentIds,
      }),
    });

    if (!response.ok) {
      throw new Error(`Failed to get answer: ${response.statusText}`);
    }

    return response.json();
  },
};