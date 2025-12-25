import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import FileUpload from '../components/FileUpload';
import uploadService from '../services/uploadService';

const UploadPage = () => {
  const [selectedFile, setSelectedFile] = useState(null);
  const [uploadProgress, setUploadProgress] = useState(0);
  const [isUploading, setIsUploading] = useState(false);
  const [uploadStatus, setUploadStatus] = useState(null); // 'success', 'error', or null
  const [uploadMessage, setUploadMessage] = useState('');
  const navigate = useNavigate();

  const handleFileSelect = (file) => {
    setSelectedFile(file);
    setUploadProgress(0);
    setUploadStatus(null);
    setUploadMessage('');
  };

  const handleUpload = async () => {
    if (!selectedFile) {
      setUploadStatus('error');
      setUploadMessage('Please select a file first');
      return;
    }

    setIsUploading(true);
    setUploadStatus(null);
    setUploadMessage('');

    try {
      // Use the new upload service with progress tracking
      const response = await uploadService.uploadDocument(selectedFile, (progress) => {
        setUploadProgress(progress);
      });

      setUploadProgress(100);
      setUploadStatus('success');
      setUploadMessage('File uploaded successfully!');
      
      // Wait a moment then navigate to chat
      setTimeout(() => {
        navigate('/chat');
      }, 1500);

    } catch (error) {
      setUploadStatus('error');
      setUploadMessage(error.message || 'Upload failed');
      setIsUploading(false);
    }
  };

  const handleSkip = () => {
    navigate('/chat');
  };

  return (
    <div className="max-w-2xl mx-auto">
      <h1 className="text-2xl font-bold text-gray-900 mb-6">Upload Document</h1>
      
      <div className="bg-white rounded-lg shadow-md p-6">
        <FileUpload 
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
                className="bg-indigo-600 h-2.5 rounded-full transition-all duration-300"
                style={{ width: `${uploadProgress}%` }}
              ></div>
            </div>
          </div>
        )}

        {/* Status messages */}
        {uploadStatus === 'success' && (
          <div className="mt-4 p-3 bg-green-50 border border-green-200 rounded-md">
            <p className="text-green-600">{uploadMessage}</p>
          </div>
        )}
        
        {uploadStatus === 'error' && (
          <div className="mt-4 p-3 bg-red-50 border border-red-200 rounded-md">
            <p className="text-red-600">{uploadMessage}</p>
          </div>
        )}

        {/* Action buttons */}
        <div className="mt-6 flex flex-col sm:flex-row gap-3">
          <button
            onClick={handleUpload}
            disabled={isUploading || !selectedFile}
            className={`px-4 py-2 rounded-md text-white font-medium ${
              isUploading || !selectedFile
                ? 'bg-gray-400 cursor-not-allowed'
                : 'bg-indigo-600 hover:bg-indigo-700'
            }`}
          >
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
      <div className="mt-8 bg-blue-50 border border-blue-200 rounded-lg p-4">
        <h3 className="font-medium text-blue-800 mb-2">About Document Upload</h3>
        <ul className="text-blue-700 text-sm list-disc pl-5 space-y-1">
          <li>Uploaded documents will be processed using RAG technology for intelligent Q&A</li>
          <li>Your documents are processed securely and not stored permanently</li>
          <li>For best results, use documents with clear text content</li>
        </ul>
      </div>
    </div>
  );
};

export default UploadPage;