import React, { useState, useRef, useCallback, forwardRef, useImperativeHandle } from 'react';
import LoadingSpinner from './LoadingSpinner';

const FileUpload = forwardRef(
  (
    {
      onFileSelect,
      allowedTypes = ['.pdf', '.docx', '.txt'],
      maxSize = 10 * 1024 * 1024,
      initialIsProcessing = false
    },
    ref
  ) => {
    const [isDragActive, setIsDragActive] = useState(false);
    const [selectedFile, setSelectedFile] = useState(null);
    const [error, setError] = useState('');
    const [isProcessing, setIsProcessing] = useState(initialIsProcessing);
    const fileInputRef = useRef(null);

    // Expose the setProcessingState function to parent components
    useImperativeHandle(ref, () => ({
      setProcessingState: (processing) => {
        setIsProcessing(processing);
      }
    }));

    const validTypes = allowedTypes.map((type) => type.toLowerCase());

    const validateFile = (file) => {
      // Check file type
      const fileExtension = '.' + file.name.split('.').pop().toLowerCase();
      if (!validTypes.includes(fileExtension)) {
        setError(`Invalid file type. Allowed types: ${allowedTypes.join(', ')}`);
        return false;
      }

      // Check file size
      if (file.size > maxSize) {
        setError(`File size too large. Maximum size: ${maxSize / (1024 * 1024)}MB`);
        return false;
      }

      setError('');
      return true;
    };

    const handleFileChange = (file) => {
      if (!file) return;

      if (validateFile(file)) {
        setSelectedFile(file);
        onFileSelect(file);
      } else {
        setSelectedFile(null);
        onFileSelect(null);
      }
    };

    const handleFileInput = (e) => {
      const file = e.target.files[0];
      handleFileChange(file);
    };

    const handleDrag = useCallback((e) => {
      e.preventDefault();
      e.stopPropagation();
    }, []);

    const handleDragIn = useCallback((e) => {
      e.preventDefault();
      e.stopPropagation();
      if (e.dataTransfer.items && e.dataTransfer.items.length > 0) {
        setIsDragActive(true);
      }
    }, []);

    const handleDragOut = useCallback((e) => {
      e.preventDefault();
      e.stopPropagation();
      setIsDragActive(false);
    }, []);

    const handleDrop = useCallback(
      (e) => {
        e.preventDefault();
        e.stopPropagation();
        setIsDragActive(false);

        if (e.dataTransfer.files && e.dataTransfer.files.length > 0) {
          const file = e.dataTransfer.files[0];
          handleFileChange(file);
          e.dataTransfer.clearData();
        }
      },
      [handleFileChange]
    );

    const handleRemoveFile = () => {
      setSelectedFile(null);
      setError('');
      onFileSelect(null);
      if (fileInputRef.current) {
        fileInputRef.current.value = '';
      }
    };

    const getFileIcon = (file) => {
      const extension = '.' + file.name.split('.').pop().toLowerCase();
      if (extension === '.pdf') return '📄';
      if (extension === '.docx') return '📝';
      if (extension === '.txt') return '📑';
      return '📁';
    };

    const formatFileSize = (bytes) => {
      if (bytes === 0) return '0 Bytes';
      const k = 1024;
      const sizes = ['Bytes', 'KB', 'MB', 'GB'];
      const i = Math.floor(Math.log(bytes) / Math.log(k));
      return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
    };

    return (
      <div className="w-full">
        {/* File selection area */}
        <div
          className={`border-2 border-dashed rounded-xl p-8 text-center cursor-pointer transition-all duration-300 ${
            isDragActive
              ? 'border-[#DC143C] bg-gradient-to-br from-[#DC143C]/10 to-red-900/20'
              : error
              ? 'border-[#DC143C] bg-gradient-to-br from-[#DC143C]/10 to-red-900/20'
              : 'border-gray-600 hover:border-[#DC143C] bg-gradient-to-br from-gray-800 to-gray-900'
          }`}
          onDragEnter={handleDragIn}
          onDragLeave={handleDragOut}
          onDragOver={handleDrag}
          onDrop={handleDrop}
          onClick={() => !isProcessing && fileInputRef.current?.click()}
        >
          <input
            type="file"
            ref={fileInputRef}
            onChange={handleFileInput}
            className="hidden"
            accept={allowedTypes.join(',')}
            disabled={isProcessing}
          />

          <div className="flex flex-col items-center justify-center">
            <div className="text-4xl mb-3 text-white">📁</div>
            <p className="text-lg font-medium text-white mb-1">
              {isProcessing
                ? 'Processing...'
                : isDragActive
                ? 'Drop your file here'
                : 'Drag & drop your file here'}
            </p>
            <p className="text-gray-300 mb-3">
              {isProcessing ? 'Please wait...' : 'or click to browse files'}
            </p>
            <p className="text-sm text-gray-400">
              Supported formats: {allowedTypes.join(', ')} | Max size:{' '}
              {maxSize / (1024 * 1024)}MB
            </p>
            {isProcessing && (
              <div className="mt-2">
                <LoadingSpinner size="md" />
              </div>
            )}
          </div>
        </div>

        {/* Error message */}
        {error && (
          <div className="mt-3 p-3 bg-gradient-to-r from-[#DC143C]/20 to-red-900/30 border border-[#DC143C]/50 rounded-lg backdrop-blur-sm">
            <p className="text-[#DC143C] text-sm font-medium">{error}</p>
          </div>
        )}

        {/* Selected file preview */}
        {selectedFile && (
          <div className="mt-4 p-4 bg-gradient-to-br from-gray-800 to-gray-900 border border-gray-700 rounded-xl flex items-center justify-between shadow-lg">
            <div className="flex items-center">
              <div className="text-2xl mr-3 text-white">{getFileIcon(selectedFile)}</div>
              <div>
                <p className="font-medium text-white truncate max-w-xs">
                  {selectedFile.name}
                </p>
                <p className="text-sm text-gray-400">
                  {formatFileSize(selectedFile.size)}
                </p>
              </div>
            </div>
            <button
              type="button"
              onClick={handleRemoveFile}
              disabled={isProcessing}
              className={`ml-4 px-4 py-2 rounded-lg font-medium transition-all duration-200 ${
                isProcessing
                  ? 'text-gray-500 cursor-not-allowed bg-gray-700'
                  : 'text-white bg-gradient-to-r from-[#DC143C] to-red-700 hover:opacity-90'
              }`}
            >
              Remove
            </button>
          </div>
        )}
      </div>
    );
  }
);

export default FileUpload;
