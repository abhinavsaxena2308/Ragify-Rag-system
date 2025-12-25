import React, { useState, useRef } from 'react';
import { useNavigate } from 'react-router-dom';
import FileUpload from '../components/FileUpload';
import uploadService from '../services/uploadService';
import { useToast } from '../components/Toast';
import LoadingSpinner from '../components/LoadingSpinner';

const UploadPage = () => {
  const [selectedFile, setSelectedFile] = useState(null);
  const [uploadProgress, setUploadProgress] = useState(0);
  const [isUploading, setIsUploading] = useState(false);
  const { addToast } = useToast();
  const fileUploadRef = useRef(null);
  const navigate = useNavigate();

  const handleFileSelect = (file) => {
    setSelectedFile(file);
    setUploadProgress(0);
  };

  const handleUpload = async () => {
    if (!selectedFile) {
      addToast('Please select a file first', 'error');
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
      addToast('File uploaded successfully!', 'success');
      
      // Navigate to chat after a short delay
      setTimeout(() => {
        navigate('/chat');
      }, 1500);

    } catch (error) {
      addToast(error.message || 'Upload failed', 'error');
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

  return (
    <div className="max-w-2xl mx-auto">
      <h1 className="text-2xl font-bold text-black mb-6">Upload Document</h1>
      
      <div className="bg-white rounded-lg shadow-xl border border-gray-200 p-6">
        <FileUpload 
          ref={fileUploadRef}
          onFileSelect={handleFileSelect}
          allowedTypes={['.pdf', '.docx', '.txt']}
          maxSize={10 * 1024 * 1024} // 10MB
        />
        
        {/* Upload progress */}
        {isUploading && (
          <div className="mt-6">
            <div className="flex justify-between mb-1">
              <span className="text-sm font-medium text-gray-700">Uploading...</span>
              <span className="text-sm font-medium text-gray-700">{uploadProgress}%</span>
            </div>
            <div className="w-full bg-gray-200 rounded-full h-2.5">
              <div
                className="bg-red-600 h-2.5 rounded-full transition-all duration-300"
                style={{ width: `${uploadProgress}%` }}
              ></div>
            </div>
          </div>
        )}

        {/* Action buttons */}
        <div className="mt-6 flex flex-col sm:flex-row gap-3">
          <button
            onClick={handleUpload}
            disabled={isUploading || !selectedFile}
            className={`px-4 py-2 rounded-md text-white font-medium flex items-center justify-center ${
              isUploading || !selectedFile
                ? 'bg-gray-400 cursor-not-allowed'
                : 'bg-red-600 hover:bg-red-700'
            }`}
          >
            {isUploading && <LoadingSpinner size="sm" className="mr-2" />}
            {isUploading ? 'Uploading...' : 'Upload Document'}
          </button>
          
          <button
            onClick={handleSkip}
            disabled={isUploading}
            className="px-4 py-2 rounded-md border border-gray-300 text-gray-700 font-medium hover:bg-gray-50 disabled:opacity-50"
          >
            Skip for now
          </button>
        </div>
      </div>

      {/* Info section */}
      <div className="mt-8 bg-gradient-to-r from-red-50 to-red-100 border border-red-200 rounded-lg p-4">
        <h3 className="font-medium text-red-800 mb-2 flex items-center">
          <span className="mr-2">ℹ️</span>
          About Document Upload
        </h3>
        <ul className="text-red-700 text-sm list-disc pl-5 space-y-1">
          <li>Uploaded documents will be processed using RAG technology for intelligent Q&A</li>
          <li>Your documents are processed securely and not stored permanently</li>
          <li>For best results, use documents with clear text content</li>
        </ul>
      </div>
    </div>
  );
};

export default UploadPage;