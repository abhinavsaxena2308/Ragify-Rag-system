import axiosClient from './axiosClient';

const uploadService = {
  // Upload a document
  uploadDocument: async (file, onUploadProgress) => {
    if (!file) {
      throw new Error('No file provided for upload');
    }

    // Validate file type - expanded to match backend
    const allowedTypes = [
      'application/pdf',
      'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
      'application/msword',
      'text/plain',
      'application/vnd.ms-excel',
      'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
      'application/vnd.ms-powerpoint',
      'application/vnd.openxmlformats-officedocument.presentationml.presentation'
    ];
    
    if (!allowedTypes.includes(file.type)) {
      throw new Error(`Invalid file type. Allowed types: PDF, DOC, DOCX, TXT, XLS, XLSX, PPT, PPTX`);
    }

    // Validate file size (10MB limit)
    const maxSize = 10 * 1024 * 1024; // 10MB in bytes
    if (file.size > maxSize) {
      throw new Error('File size exceeds 10MB limit');
    }

    const formData = new FormData();
    formData.append('file', file);

    try {
      const response = await axiosClient.post('/documents', formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
        onUploadProgress: (progressEvent) => {
          if (onUploadProgress && typeof onUploadProgress === 'function') {
            const progress = Math.round((progressEvent.loaded * 100) / progressEvent.total);
            onUploadProgress(progress);
          }
        },
      });

      return response.data;
    } catch (error) {
      // Handle specific error cases
      if (error.response) {
        // Server responded with error status
        const { status, data } = error.response;
        if (status === 413) {
          throw new Error('File too large');
        } else if (status === 415) {
          throw new Error('Unsupported file type');
        } else if (data && data.message) {
          throw new Error(data.message);
        } else if (data && data.error) {
          throw new Error(data.error);
        } else {
          throw new Error(`Upload failed with status ${status}: ${data || 'Unknown error'}`);
        }
      } else if (error.request) {
        // Request was made but no response received
        throw new Error('Network error: Unable to reach server. Make sure the backend is running on port 8080.');
      } else {
        // Something else happened
        throw new Error(error.message || 'An error occurred during upload');
      }
    }
  },

  // Get all documents
  getDocuments: async () => {
    try {
      const response = await axiosClient.get('/documents');
      return response.data;
    } catch (error) {
      if (error.response) {
        throw new Error(error.response.data.message || `Failed to fetch documents: ${error.response.status}`);
      } else if (error.request) {
        throw new Error('Network error: Unable to reach server. Make sure the backend is running on port 8080.');
      } else {
        throw new Error(error.message || 'An error occurred while fetching documents');
      }
    }
  },

  // Delete a document
  deleteDocument: async (documentId) => {
    if (!documentId) {
      throw new Error('Document ID is required for deletion');
    }

    try {
      const response = await axiosClient.delete(`/documents/${documentId}`);
      return response.data;
    } catch (error) {
      if (error.response) {
        throw new Error(error.response.data.message || `Failed to delete document: ${error.response.status}`);
      } else if (error.request) {
        throw new Error('Network error: Unable to reach server. Make sure the backend is running on port 8080.');
      } else {
        throw new Error(error.message || 'An error occurred while deleting document');
      }
    }
  },

  // Get document details
  getDocument: async (documentId) => {
    if (!documentId) {
      throw new Error('Document ID is required');
    }

    try {
      const response = await axiosClient.get(`/documents/${documentId}`);
      return response.data;
    } catch (error) {
      if (error.response) {
        throw new Error(error.response.data.message || `Failed to fetch document: ${error.response.status}`);
      } else if (error.request) {
        throw new Error('Network error: Unable to reach server. Make sure the backend is running on port 8080.');
      } else {
        throw new Error(error.message || 'An error occurred while fetching document');
      }
    }
  },
};

export default uploadService;