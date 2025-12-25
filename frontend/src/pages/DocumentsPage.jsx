import React, { useState, useEffect } from 'react';
import DocumentCard from '../components/DocumentCard';
import uploadService from '../services/uploadService';
import { useToast } from '../components/Toast';
import LoadingSpinner from '../components/LoadingSpinner';
import EmptyState from '../components/EmptyState';

const DocumentsPage = () => {
  const [documents, setDocuments] = useState([]);
  const [loading, setLoading] = useState(true);
  const [uploading, setUploading] = useState(false);
  const { addToast } = useToast();
  
  // Mock documents for now
  useEffect(() => {
    // Fetch documents from API
    const fetchDocuments = async () => {
      try {
        const data = await uploadService.getDocuments();
        setDocuments(data);
      } catch (error) {
        console.error('Error fetching documents:', error);
        addToast('Failed to load documents. Please try again later.', 'error');
      } finally {
        setLoading(false);
      }
    };
    fetchDocuments();
    
    // For now, using mock data
    setDocuments([
      {
        id: 1,
        name: 'Company Handbook.pdf',
        size: 2457600, // 2.4MB
        uploadDate: new Date().toISOString(),
        type: 'application/pdf'
      },
      {
        id: 2,
        name: 'Financial Report Q3.docx',
        size: 1048576, // 1MB
        uploadDate: new Date(Date.now() - 86400000).toISOString(), // yesterday
        type: 'application/vnd.openxmlformats-officedocument.wordprocessingml.document'
      },
      {
        id: 3,
        name: 'Product Specifications.txt',
        size: 51200, // 50KB
        uploadDate: new Date(Date.now() - 172800000).toISOString(), // 2 days ago
        type: 'text/plain'
      }
    ]);
    setLoading(false);
  }, []);
  
  const handleDeleteDocument = async (id) => {
    try {
      // Call the API to delete the document
      await uploadService.deleteDocument(id);
      
      // Update the local state to remove the document
      setDocuments(documents.filter(doc => doc.id !== id));
      addToast('Document deleted successfully', 'success');
    } catch (error) {
      console.error('Error deleting document:', error);
      addToast('Failed to delete document. Please try again.', 'error');
    }
  };
  
  const handleFileUpload = async (e) => {
    const file = e.target.files[0];
    if (file) {
      setUploading(true);
      try {
        // Call the API to upload the document
        await uploadService.uploadDocument(file, (progress) => {
          // Optionally show upload progress
        });
        
        // Refresh the document list
        const data = await uploadService.getDocuments();
        setDocuments(data);
        addToast('Document uploaded successfully', 'success');
      } catch (error) {
        console.error('Error uploading document:', error);
        addToast('Failed to upload document. Please try again.', 'error');
      } finally {
        setUploading(false);
      }
    }
  };
  
  return (
    <div>
      <h1 className="text-2xl font-bold text-gray-900 mb-6">Documents</h1>
      
      <div className="mb-6 bg-white rounded-lg shadow-md p-6">
        <h2 className="text-lg font-medium text-gray-900 mb-3">Upload New Document</h2>
        <div className="flex items-center space-x-4">
          <label className={`bg-indigo-600 text-white px-4 py-2 rounded-md cursor-pointer ${uploading ? 'opacity-75 cursor-not-allowed' : 'hover:bg-indigo-700'}`}>
            {uploading && <LoadingSpinner size="sm" className="mr-2 inline-block" />}
            {uploading ? 'Uploading...' : 'Choose File'}
            <input 
              type="file" 
              className="hidden" 
              accept=".pdf,.doc,.docx,.txt,.xls,.xlsx,.ppt,.pptx"
              onChange={handleFileUpload}
              disabled={uploading}
            />
          </label>
          <span className="text-gray-500 text-sm">Supported formats: PDF, DOC, DOCX, TXT, XLS, XLSX, PPT, PPTX</span>
        </div>
      </div>
      
      <div>
        <h2 className="text-lg font-medium text-gray-900 mb-4">Your Documents</h2>
        {loading ? (
          <div className="text-center py-12">
            <LoadingSpinner size="lg" className="mx-auto" />
            <p className="mt-4 text-gray-500">Loading documents...</p>
          </div>
        ) : documents.length === 0 ? (
          <EmptyState
            title="No documents yet"
            description="Upload your first document to get started with RAGify. Your documents will appear here once uploaded."
            icon="📄"
            actionButton={
              <label className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 cursor-pointer">
                Upload Document
                <input 
                  type="file" 
                  className="hidden" 
                  accept=".pdf,.doc,.docx,.txt,.xls,.xlsx,.ppt,.pptx"
                  onChange={handleFileUpload}
                  disabled={uploading}
                />
              </label>
            }
          />
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {documents.map((document) => (
              <DocumentCard 
                key={document.id} 
                document={document} 
                onDelete={handleDeleteDocument} 
              />
            ))}
          </div>
        )}
      </div>
    </div>
  );
};

export default DocumentsPage;