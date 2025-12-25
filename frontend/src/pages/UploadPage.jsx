import React, { useState, useRef } from 'react';
import { useNavigate } from 'react-router-dom';
import uploadService from '../services/uploadService';

const UploadPage = () => {
  const [selectedFile, setSelectedFile] = useState(null);
  const [uploadProgress, setUploadProgress] = useState(0);
  const [isUploading, setIsUploading] = useState(false);
  const fileUploadRef = useRef(null);
  const navigate = useNavigate();

  const handleFileSelect = (file) => {
    setSelectedFile(file);
    setUploadProgress(0);
  };

  const handleUpload = async () => {
    if (!selectedFile) {
      alert('Please select a file first');
      return;
    }

    setIsUploading(true);
    
    // Set the file upload component to processing state
    if (fileUploadRef.current && fileUploadRef.current.setProcessingState) {
      fileUploadRef.current.setProcessingState(true);
    }

    try {
      // Use the new upload service with progress tracking
      const response = await uploadService.uploadDocument(selectedFile, (progress) => {
        setUploadProgress(progress);
      });

      setUploadProgress(100);
      alert('File uploaded successfully!');
      
      // Navigate to chat after a short delay
      setTimeout(() => {
        navigate('/chat');
      }, 1500);

    } catch (error) {
      alert(error.message || 'Upload failed');
    } finally {
      setIsUploading(false);
      // Reset the file upload component processing state
      if (fileUploadRef.current && fileUploadRef.current.setProcessingState) {
        fileUploadRef.current.setProcessingState(false);
      }
    }
  };

  const handleSkip = () => {
    navigate('/chat');
  };

  const handleFileChange = (e) => {
    const file = e.target.files[0];
    if (file) {
      handleFileSelect(file);
    }
  };

  return (
    <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', height: '100%', padding: '20px', fontFamily: 'Arial, sans-serif' }}>
      <div style={{ width: '100%', maxWidth: '600px' }}>
        {/* Header */}
        <div style={{ textAlign: 'center', marginBottom: '40px' }}>
          <div style={{ fontSize: '48px', marginBottom: '16px' }}>📤</div>
          <h1 style={{ fontSize: '32px', marginBottom: '8px' }}>Upload Document</h1>
          <p style={{ color: '#666' }}>Upload your documents for AI-powered Q&A</p>
        </div>
        
        {/* Upload Area */}
        <div style={{ border: '1px solid #ddd', borderRadius: '8px', padding: '32px', backgroundColor: '#f9f9f9', marginBottom: '24px' }}>
          <div style={{ border: '2px dashed #ccc', borderRadius: '8px', padding: '40px', textAlign: 'center', backgroundColor: '#fff' }}>
            <div style={{ fontSize: '48px', marginBottom: '16px' }}>📁</div>
            <p style={{ marginBottom: '16px', color: '#666' }}>
              {selectedFile ? selectedFile.name : 'Click to select or drag and drop your file here'}
            </p>
            <input
              type="file"
              ref={fileUploadRef}
              onChange={handleFileChange}
              accept=".pdf,.docx,.txt"
              style={{ display: 'none' }}
            />
            <button
              onClick={() => fileUploadRef.current?.click()}
              style={{
                padding: '12px 24px',
                backgroundColor: '#007bff',
                color: 'white',
                border: 'none',
                borderRadius: '4px',
                cursor: 'pointer',
                fontSize: '16px'
              }}
            >
              Choose File
            </button>
          </div>
          
          {/* Upload progress */}
          {isUploading && (
            <div style={{ marginTop: '24px' }}>
              <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '8px' }}>
                <span>Uploading...</span>
                <span>{uploadProgress}%</span>
              </div>
              <div style={{ width: '100%', backgroundColor: '#e9ecef', borderRadius: '4px', height: '8px' }}>
                <div
                  style={{
                    backgroundColor: '#007bff',
                    height: '100%',
                    borderRadius: '4px',
                    transition: 'width 0.3s ease',
                    width: `${uploadProgress}%`
                  }}
                ></div>
              </div>
            </div>
          )}

          {/* Action buttons */}
          <div style={{ marginTop: '24px', display: 'flex', gap: '12px', justifyContent: 'center' }}>
            <button
              onClick={handleUpload}
              disabled={isUploading || !selectedFile}
              style={{
                padding: '12px 24px',
                backgroundColor: isUploading || !selectedFile ? '#6c757d' : '#007bff',
                color: 'white',
                border: 'none',
                borderRadius: '4px',
                cursor: isUploading || !selectedFile ? 'not-allowed' : 'pointer',
                fontSize: '16px'
              }}
            >
              {isUploading ? 'Uploading...' : 'Upload Document'}
            </button>
            
            <button
              onClick={handleSkip}
              disabled={isUploading}
              style={{
                padding: '12px 24px',
                backgroundColor: '#6c757d',
                color: 'white',
                border: '1px solid #ccc',
                borderRadius: '4px',
                cursor: isUploading ? 'not-allowed' : 'pointer',
                fontSize: '16px'
              }}
            >
              Skip for now
            </button>
          </div>
        </div>

        {/* Info section */}
        <div style={{ border: '1px solid #ddd', borderRadius: '8px', padding: '24px', backgroundColor: '#f9f9f9' }}>
          <h3 style={{ marginBottom: '16px', fontSize: '18px' }}>About Document Upload</h3>
          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(250px, 1fr))', gap: '16px' }}>
            <div style={{ display: 'flex', alignItems: 'flex-start', gap: '12px' }}>
              <div style={{ fontSize: '24px' }}>🔍</div>
              <div>
                <h4 style={{ margin: '0 0 4px 0', fontSize: '16px' }}>Smart Processing</h4>
                <p style={{ margin: 0, color: '#666', fontSize: '14px' }}>
                  Documents are processed using RAG technology for intelligent Q&A
                </p>
              </div>
            </div>
            <div style={{ display: 'flex', alignItems: 'flex-start', gap: '12px' }}>
              <div style={{ fontSize: '24px' }}>🔒</div>
              <div>
                <h4 style={{ margin: '0 0 4px 0', fontSize: '16px' }}>Secure Upload</h4>
                <p style={{ margin: 0, color: '#666', fontSize: '14px' }}>
                  Your documents are processed securely and not stored permanently
                </p>
              </div>
            </div>
            <div style={{ display: 'flex', alignItems: 'flex-start', gap: '12px' }}>
              <div style={{ fontSize: '24px' }}>📄</div>
              <div>
                <h4 style={{ margin: '0 0 4px 0', fontSize: '16px' }}>Multiple Formats</h4>
                <p style={{ margin: 0, color: '#666', fontSize: '14px' }}>
                  Support for PDF, DOCX, and TXT files up to 10MB
                </p>
              </div>
            </div>
            <div style={{ display: 'flex', alignItems: 'flex-start', gap: '12px' }}>
              <div style={{ fontSize: '24px' }}>⚡</div>
              <div>
                <h4 style={{ margin: '0 0 4px 0', fontSize: '16px' }}>Fast Processing</h4>
                <p style={{ margin: 0, color: '#666', fontSize: '14px' }}>
                  Quick analysis and immediate availability for Q&A
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default UploadPage;