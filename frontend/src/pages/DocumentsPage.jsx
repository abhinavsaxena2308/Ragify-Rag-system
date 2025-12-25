import React, { useState, useEffect } from 'react';
import uploadService from '../services/uploadService';

const DocumentsPage = () => {
  const [documents, setDocuments] = useState([]);
  const [loading, setLoading] = useState(true);
  const [uploading, setUploading] = useState(false);
  
  // Mock documents for now
  useEffect(() => {
    // Fetch documents from API
    const fetchDocuments = async () => {
      try {
        const data = await uploadService.getDocuments();
        setDocuments(data);
      } catch (error) {
        console.error('Error fetching documents:', error);
        alert('Failed to load documents. Please try again later.');
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
    if (confirm('Are you sure you want to delete this document?')) {
      try {
        // Call the API to delete the document
        await uploadService.deleteDocument(id);
        
        // Update the local state to remove the document
        setDocuments(documents.filter(doc => doc.id !== id));
        alert('Document deleted successfully');
      } catch (error) {
        console.error('Error deleting document:', error);
        alert('Failed to delete document. Please try again.');
      }
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
        alert('Document uploaded successfully');
      } catch (error) {
        console.error('Error uploading document:', error);
        alert('Failed to upload document. Please try again.');
      } finally {
        setUploading(false);
      }
    }
  };

  const formatFileSize = (bytes) => {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

  const formatDate = (dateString) => {
    const date = new Date(dateString);
    return date.toLocaleDateString() + ' ' + date.toLocaleTimeString();
  };
  
  return (
    <div style={{ display: 'flex', flexDirection: 'column', height: '100%', padding: '20px', fontFamily: 'Arial, sans-serif' }}>
      {/* Header */}
      <div style={{ textAlign: 'center', marginBottom: '40px' }}>
        <div style={{ fontSize: '48px', marginBottom: '16px' }}>📁</div>
        <h1 style={{ fontSize: '32px', marginBottom: '8px' }}>My Documents</h1>
        <p style={{ color: '#666' }}>Manage your uploaded documents</p>
      </div>
      
      {/* Upload Section */}
      <div style={{ marginBottom: '40px', maxWidth: '600px', margin: '0 auto 40px' }}>
        <div style={{ border: '1px solid #ddd', borderRadius: '8px', padding: '24px', backgroundColor: '#f9f9f9' }}>
          <h2 style={{ marginBottom: '16px', fontSize: '18px' }}>Upload New Document</h2>
          <div style={{ display: 'flex', alignItems: 'center', gap: '16px', flexWrap: 'wrap' }}>
            <label style={{
              backgroundColor: uploading ? '#6c757d' : '#007bff',
              color: 'white',
              padding: '12px 20px',
              borderRadius: '4px',
              cursor: uploading ? 'not-allowed' : 'pointer',
              fontSize: '16px'
            }}>
              {uploading ? 'Uploading...' : 'Choose File'}
              <input 
                type="file" 
                style={{ display: 'none' }} 
                accept=".pdf,.doc,.docx,.txt,.xls,.xlsx,.ppt,.pptx"
                onChange={handleFileUpload}
                disabled={uploading}
              />
            </label>
            <span style={{ color: '#666', fontSize: '14px' }}>
              PDF, DOC, DOCX, TXT, XLS, XLSX, PPT, PPTX (Max 10MB)
            </span>
          </div>
        </div>
      </div>
      
      {/* Documents List */}
      <div style={{ flex: 1, overflowY: 'auto' }}>
        <div style={{ maxWidth: '1200px', margin: '0 auto' }}>
          {loading ? (
            <div style={{ textAlign: 'center', padding: '40px' }}>
              <div style={{ fontSize: '24px', marginBottom: '16px' }}>⏳</div>
              <p style={{ color: '#666' }}>Loading documents...</p>
            </div>
          ) : documents.length === 0 ? (
            <div style={{ textAlign: 'center', padding: '40px' }}>
              <div style={{ fontSize: '48px', marginBottom: '16px' }}>📄</div>
              <h2 style={{ fontSize: '24px', marginBottom: '16px' }}>No documents yet</h2>
              <p style={{ color: '#666', marginBottom: '24px', maxWidth: '400px', margin: '0 auto 24px' }}>
                Upload your first document to get started with RAGify. Your documents will appear here once uploaded.
              </p>
              <label style={{
                backgroundColor: '#007bff',
                color: 'white',
                padding: '12px 24px',
                borderRadius: '4px',
                cursor: 'pointer',
                fontSize: '16px'
              }}>
                Upload Document
                <input 
                  type="file" 
                  style={{ display: 'none' }} 
                  accept=".pdf,.doc,.docx,.txt,.xls,.xlsx,.ppt,.pptx"
                  onChange={handleFileUpload}
                  disabled={uploading}
                />
              </label>
            </div>
          ) : (
            <div>
              <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: '24px' }}>
                <h2 style={{ fontSize: '20px', margin: 0 }}>Your Documents ({documents.length})</h2>
                <div style={{ display: 'flex', gap: '16px' }}>
                  <button style={{
                    background: 'none',
                    border: '1px solid #ddd',
                    padding: '8px',
                    cursor: 'pointer',
                    borderRadius: '4px'
                  }}>
                    🔄
                  </button>
                  <button style={{
                    background: 'none',
                    border: '1px solid #ddd',
                    padding: '8px',
                    cursor: 'pointer',
                    borderRadius: '4px'
                  }}>
                    ☰
                  </button>
                </div>
              </div>
              <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(280px, 1fr))', gap: '20px' }}>
                {documents.map((document) => (
                  <div key={document.id} style={{
                    border: '1px solid #ddd',
                    borderRadius: '8px',
                    padding: '20px',
                    backgroundColor: '#f9f9f9',
                    transition: 'all 0.3s ease'
                  }}>
                    <div style={{ display: 'flex', alignItems: 'center', marginBottom: '12px' }}>
                      <div style={{
                        fontSize: '32px',
                        marginRight: '12px',
                        backgroundColor: '#e9ecef',
                        padding: '8px',
                        borderRadius: '4px'
                      }}>
                        📄
                      </div>
                      <div style={{ flex: 1, minWidth: 0 }}>
                        <h3 style={{
                          margin: 0,
                          fontSize: '16px',
                          fontWeight: 'bold',
                          whiteSpace: 'nowrap',
                          overflow: 'hidden',
                          textOverflow: 'ellipsis'
                        }}>
                          {document.name}
                        </h3>
                        <p style={{ 
                          margin: '4px 0 0 0', 
                          color: '#666', 
                          fontSize: '14px' 
                        }}>
                          {formatFileSize(document.size)}
                        </p>
                      </div>
                    </div>
                    <div style={{ 
                      fontSize: '12px', 
                      color: '#666', 
                      marginBottom: '16px' 
                    }}>
                      Uploaded: {formatDate(document.uploadDate)}
                    </div>
                    <div style={{ display: 'flex', gap: '8px' }}>
                      <button style={{
                        flex: 1,
                        padding: '8px 12px',
                        backgroundColor: '#007bff',
                        color: 'white',
                        border: 'none',
                        borderRadius: '4px',
                        cursor: 'pointer',
                        fontSize: '14px'
                      }}>
                        View
                      </button>
                      <button
                        onClick={() => handleDeleteDocument(document.id)}
                        style={{
                          flex: 1,
                          padding: '8px 12px',
                          backgroundColor: '#dc3545',
                          color: 'white',
                          border: 'none',
                          borderRadius: '4px',
                          cursor: 'pointer',
                          fontSize: '14px'
                        }}
                      >
                        Delete
                      </button>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default DocumentsPage;