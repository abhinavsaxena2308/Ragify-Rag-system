import React, { useState, useEffect } from 'react';
import DocumentCard from '../components/DocumentCard';
import uploadService from '../services/uploadService';

const DocumentsPage = () => {
  const [documents, setDocuments] = useState([]);
  const [loading, setLoading] = useState(true);
  
  // Mock documents for now
  useEffect(() => {
    // Fetch documents from API
    const fetchDocuments = async () => {
      try {
        const data = await uploadService.getDocuments();
        setDocuments(data);
      } catch (error) {
        console.error('Error fetching documents:', error);
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
    } catch (error) {
      console.error('Error deleting document:', error);
      // In a real app, you might want to show an error message to the user
    }
  };
  
  const handleFileUpload = (e) => {
    const file = e.target.files[0];
    if (file) {
      // In a real app, we would call the API:
      // await documentService.uploadDocument(file);
      
      // For now, just add to the list
      const newDoc = {
        id: documents.length + 1,
        name: file.name,
        size: file.size,
        uploadDate: new Date().toISOString(),
        type: file.type
      };
      setDocuments([newDoc, ...documents]);
    }
  };
  
  return (
    <div>
      <h1 className="text-2xl font-bold text-gray-900 mb-6">Documents</h1>
      
      <div className="mb-6 bg-white rounded-lg shadow-md p-6">
        <h2 className="text-lg font-medium text-gray-900 mb-3">Upload New Document</h2>
        <div className="flex items-center space-x-4">
          <label className="bg-indigo-600 text-white px-4 py-2 rounded-md hover:bg-indigo-700 cursor-pointer">
            Choose File
            <input 
              type="file" 
              className="hidden" 
              accept=".pdf,.doc,.docx,.txt,.xls,.xlsx,.ppt,.pptx"
              onChange={handleFileUpload}
            />
          </label>
          <span className="text-gray-500 text-sm">Supported formats: PDF, DOC, DOCX, TXT, XLS, XLSX, PPT, PPTX</span>
        </div>
      </div>
      
      <div>
        <h2 className="text-lg font-medium text-gray-900 mb-4">Your Documents</h2>
        {loading ? (
          <div className="text-center py-8">
            <p className="text-gray-500">Loading documents...</p>
          </div>
        ) : documents.length === 0 ? (
          <div className="text-center py-8">
            <p className="text-gray-500">No documents uploaded yet.</p>
          </div>
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