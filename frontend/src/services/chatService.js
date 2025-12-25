import axiosClient from './axiosClient';

const chatService = {
  // Send a question and get an answer from the AI
  askQuestion: async (question, documentIds = []) => {
    if (!question || typeof question !== 'string' || question.trim().length === 0) {
      throw new Error('Question is required and must be a non-empty string');
    }

    try {
      const response = await axiosClient.post('/chat', {
        question: question.trim(),
        documentIds: Array.isArray(documentIds) ? documentIds : []
      });

      return response.data;
    } catch (error) {
      if (error.response) {
        // Server responded with error status
        const { status, data } = error.response;
        if (status === 400) {
          throw new Error(data.message || 'Invalid request: Please check your question');
        } else if (status === 404) {
          throw new Error('No documents found to answer your question. Please upload documents first.');
        } else if (status === 500) {
          throw new Error(data.message || 'Internal server error: Unable to process your request');
        } else {
          throw new Error(data.message || `Request failed with status ${status}`);
        }
      } else if (error.request) {
        // Request was made but no response received
        throw new Error('Network error: Unable to reach the server');
      } else {
        // Something else happened
        throw new Error(error.message || 'An error occurred while processing your question');
      }
    }
  },

  // Get chat history (if the backend supports it)
  getChatHistory: async (limit = 10, offset = 0) => {
    try {
      const response = await axiosClient.get('/chat/history', {
        params: {
          limit,
          offset
        }
      });

      return response.data;
    } catch (error) {
      if (error.response) {
        throw new Error(error.response.data.message || `Failed to fetch chat history: ${error.response.status}`);
      } else if (error.request) {
        throw new Error('Network error: Unable to reach the server');
      } else {
        throw new Error(error.message || 'An error occurred while fetching chat history');
      }
    }
  },

  // Get sources for a specific question (if the backend supports it)
  getSources: async (questionId) => {
    if (!questionId) {
      throw new Error('Question ID is required');
    }

    try {
      const response = await axiosClient.get(`/chat/${questionId}/sources`);
      return response.data;
    } catch (error) {
      if (error.response) {
        throw new Error(error.response.data.message || `Failed to fetch sources: ${error.response.status}`);
      } else if (error.request) {
        throw new Error('Network error: Unable to reach the server');
      } else {
        throw new Error(error.message || 'An error occurred while fetching sources');
      }
    }
  },

  // Stream chat response (for real-time streaming if supported by backend)
  streamQuestion: async (question, documentIds = [], onMessage) => {
    if (!question || typeof question !== 'string' || question.trim().length === 0) {
      throw new Error('Question is required and must be a non-empty string');
    }

    if (typeof onMessage !== 'function') {
      throw new Error('onMessage callback function is required for streaming');
    }

    // Note: This is a placeholder implementation
    // Actual streaming implementation would depend on backend capabilities
    // (e.g., Server-Sent Events, WebSocket, or fetch with ReadableStream)
    
    try {
      // For now, we'll simulate streaming by calling the regular askQuestion endpoint
      // In a real implementation, this would use a different endpoint that supports streaming
      const response = await axiosClient.post('/chat/stream', {
        question: question.trim(),
        documentIds: Array.isArray(documentIds) ? documentIds : []
      }, {
        // This would be configured differently for actual streaming
        responseType: 'stream'
      });

      // Process the stream response
      // This is a simplified example - actual implementation would depend on the backend
      return response.data;
    } catch (error) {
      if (error.response) {
        throw new Error(error.response.data.message || `Streaming request failed: ${error.response.status}`);
      } else if (error.request) {
        throw new Error('Network error: Unable to reach the server');
      } else {
        throw new Error(error.message || 'An error occurred during streaming');
      }
    }
  }
};

export default chatService;