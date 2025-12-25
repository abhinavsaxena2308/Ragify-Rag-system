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
    <div className="flex flex-col items-center justify-center h-full px-4 py-8">
      <div className="w-full max-w-2xl">
        {/* Header */}
        <div className="text-center mb-10">
          <div className="w-16 h-16 bg-gradient-to-r from-red-600 to-red-700 rounded-full flex items-center justify-center mx-auto mb-4">
            <span className="text-white text-2xl">📤</span>
          </div>
          <h1 className="text-3xl font-bold text-white mb-2">Upload Document</h1>
          <p className="text-gray-400">Upload your documents for AI-powered Q&A</p>
        </div>
        
        {/* Upload Area */}
        <div className="bg-gray-800/50 backdrop-blur-sm rounded-2xl shadow-xl border border-gray-700 p-8 mb-6">
          <FileUpload 
            ref={fileUploadRef}
            onFileSelect={handleFileSelect}
            allowedTypes={['.pdf', '.docx', '.txt']}
            maxSize={10 * 1024 * 1024} // 10MB
          />
          
          {/* Upload progress */}
          {isUploading && (
            <div className="mt-6">
              <div className="flex justify-between mb-2">
                <span className="text-sm font-medium text-gray-300">Uploading...</span>
                <span className="text-sm font-medium text-gray-300">{uploadProgress}%</span>
              </div>
              <div className="w-full bg-gray-700 rounded-full h-3">
                <div
                  className="bg-gradient-to-r from-red-600 to-red-500 h-3 rounded-full transition-all duration-300"
                  style={{ width: `${uploadProgress}%` }}
                ></div>
              </div>
            </div>
          )}

          {/* Action buttons */}
          <div className="mt-8 flex flex-col sm:flex-row gap-3">
            <button
              onClick={handleUpload}
              disabled={isUploading || !selectedFile}
              className={`flex-1 px-6 py-3 rounded-lg text-white font-medium flex items-center justify-center transition-all duration-200 ${
                isUploading || !selectedFile
                  ? 'bg-gray-700 text-gray-400 cursor-not-allowed'
                  : 'bg-gradient-to-r from-red-600 to-red-700 hover:from-red-700 hover:to-red-800 shadow-lg'
              }`}
            >
              {isUploading && <LoadingSpinner size="sm" className="mr-2" />}
              {isUploading ? 'Uploading...' : 'Upload Document'}
            </button>
            
            <button
              onClick={handleSkip}
              disabled={isUploading}
              className="px-6 py-3 rounded-lg border border-gray-600 text-gray-300 font-medium hover:bg-gray-700 disabled:opacity-50 transition-colors duration-200"
            >
              Skip for now
            </button>
          </div>
        </div>

        {/* Info section */}
        <div className="bg-gray-800/30 backdrop-blur-sm border border-gray-700 rounded-xl p-6">
          <h3 className="font-semibold text-white mb-4 flex items-center">
            <svg className="w-5 h-5 mr-2 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            About Document Upload
          </h3>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="flex items-start space-x-3">
              <div className="w-8 h-8 bg-red-600/20 rounded-lg flex items-center justify-center flex-shrink-0">
                <span className="text-red-400 text-sm">🔍</span>
              </div>
              <div>
                <h4 className="text-white font-medium mb-1">Smart Processing</h4>
                <p className="text-gray-400 text-sm">Documents are processed using RAG technology for intelligent Q&A</p>
              </div>
            </div>
            <div className="flex items-start space-x-3">
              <div className="w-8 h-8 bg-red-600/20 rounded-lg flex items-center justify-center flex-shrink-0">
                <span className="text-red-400 text-sm">🔒</span>
              </div>
              <div>
                <h4 className="text-white font-medium mb-1">Secure Upload</h4>
                <p className="text-gray-400 text-sm">Your documents are processed securely and not stored permanently</p>
              </div>
            </div>
            <div className="flex items-start space-x-3">
              <div className="w-8 h-8 bg-red-600/20 rounded-lg flex items-center justify-center flex-shrink-0">
                <span className="text-red-400 text-sm">📄</span>
              </div>
              <div>
                <h4 className="text-white font-medium mb-1">Multiple Formats</h4>
                <p className="text-gray-400 text-sm">Support for PDF, DOCX, and TXT files up to 10MB</p>
              </div>
            </div>
            <div className="flex items-start space-x-3">
              <div className="w-8 h-8 bg-red-600/20 rounded-lg flex items-center justify-center flex-shrink-0">
                <span className="text-red-400 text-sm">⚡</span>
              </div>
              <div>
                <h4 className="text-white font-medium mb-1">Fast Processing</h4>
                <p className="text-gray-400 text-sm">Quick analysis and immediate availability for Q&A</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default UploadPage;