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
              ? 'border-red-500 bg-red-50'
              : error
              ? 'border-red-500 bg-red-50'
              : 'border-gray-300 hover:border-red-400 bg-gray-50'
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
            <div className="text-4xl mb-3">📁</div>
            <p className="text-lg font-medium text-gray-700 mb-1">
              {isProcessing
                ? 'Processing...'
                : isDragActive
                ? 'Drop your file here'
                : 'Drag & drop your file here'}
            </p>
            <p className="text-gray-600 mb-3">
              {isProcessing ? 'Please wait...' : 'or click to browse files'}
            </p>
            <p className="text-sm text-gray-500">
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
          <div className="mt-3 p-3 bg-red-50 border border-red-200 rounded-lg">
            <p className="text-red-600 text-sm">{error}</p>
          </div>
        )}

        {/* Selected file preview */}
        {selectedFile && (
          <div className="mt-4 p-4 bg-white border border-gray-200 rounded-lg flex items-center justify-between shadow-sm">
            <div className="flex items-center">
              <div className="text-2xl mr-3">{getFileIcon(selectedFile)}</div>
              <div>
                <p className="font-medium text-gray-900 truncate max-w-xs">
                  {selectedFile.name}
                </p>
                <p className="text-sm text-gray-500">
                  {formatFileSize(selectedFile.size)}
                </p>
              </div>
            </div>
            <button
              type="button"
              onClick={handleRemoveFile}
              disabled={isProcessing}
              className={`ml-4 px-3 py-1 rounded-md ${
                isProcessing
                  ? 'text-gray-400 cursor-not-allowed bg-gray-100'
                  : 'text-white bg-red-600 hover:bg-red-700'
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
