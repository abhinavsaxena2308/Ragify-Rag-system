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
    <div className="w-full">
      <div className="flex items-center mb-8">
        <div className="w-12 h-12 bg-gradient-to-r from-red-600 to-red-700 rounded-xl flex items-center justify-center mr-4">
          <span className="text-white text-2xl">📁</span>
        </div>
        <div>
          <h1 className="text-2xl font-bold text-white">Documents</h1>
          <p className="text-gray-400 text-sm">Manage your uploaded documents</p>
        </div>
      </div>
      
      <div className="mb-8 bg-gray-800 rounded-xl shadow-xl border border-gray-700 p-6">
        <h2 className="text-lg font-medium text-white mb-4">Upload New Document</h2>
        <div className="flex flex-col sm:flex-row items-start sm:items-center space-y-4 sm:space-y-0 sm:space-x-4">
          <label className={`bg-gradient-to-r from-red-600 to-red-700 text-white px-4 py-3 rounded-lg cursor-pointer flex items-center font-medium transition-all duration-200 ${uploading ? 'opacity-75 cursor-not-allowed' : 'hover:from-red-700 hover:to-red-800 shadow-md'}`}>
            {uploading && <LoadingSpinner size="sm" className="mr-2" />}
            {uploading ? 'Uploading...' : 'Choose File'}
            <input 
              type="file" 
              className="hidden" 
              accept=".pdf,.doc,.docx,.txt,.xls,.xlsx,.ppt,.pptx"
              onChange={handleFileUpload}
              disabled={uploading}
            />
          </label>
          <span className="text-gray-400 text-sm">Supported formats: PDF, DOC, DOCX, TXT, XLS, XLSX, PPT, PPTX</span>
        </div>
      </div>
      
      <div>
        <h2 className="text-lg font-medium text-white mb-6">Your Documents</h2>
        {loading ? (
          <div className="text-center py-12">
            <LoadingSpinner size="lg" className="mx-auto text-red-500" />
            <p className="mt-4 text-gray-400">Loading documents...</p>
          </div>
        ) : documents.length === 0 ? (
          <EmptyState
            title="No documents yet"
            description="Upload your first document to get started with RAGify. Your documents will appear here once uploaded."
            icon="📄"
            actionButton={
              <label className="inline-flex items-center px-4 py-3 border border-transparent text-sm font-medium rounded-lg shadow-sm text-white bg-gradient-to-r from-red-600 to-red-700 hover:from-red-700 hover:to-red-800 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 cursor-pointer transition-all duration-200">
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
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
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