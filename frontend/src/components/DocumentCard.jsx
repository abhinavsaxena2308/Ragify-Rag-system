import React from 'react';

const DocumentCard = ({ document, onDelete }) => {
  const { id, name, size, uploadDate, type } = document;

  // Format file size
  const formatFileSize = (bytes) => {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

  // Get file icon based on type
  const getFileIcon = (type) => {
    if (type.includes('pdf')) return '📄';
    if (type.includes('word') || type.includes('doc')) return '📝';
    if (type.includes('excel') || type.includes('sheet')) return '📊';
    if (type.includes('text') || type.includes('plain')) return '📑';
    if (type.includes('powerpoint') || type.includes('presentation')) return '📽️';
    return '📁';
  };

  return (
    <div className="border border-gray-700 rounded-xl p-4 hover:shadow-xl transition-all duration-300 bg-gray-800 shadow-lg">
      <div className="flex items-start">
        <div className="text-2xl mr-3">{getFileIcon(type)}</div>
        <div className="flex-1">
          <h3 className="font-medium text-white truncate">{name}</h3>
          <div className="mt-1 text-sm text-gray-400">
            <p>{formatFileSize(size)}</p>
            <p>Uploaded: {new Date(uploadDate).toLocaleDateString()}</p>
          </div>
        </div>
        <button
          onClick={() => onDelete(id)}
          className="text-[#DC143C] hover:text-red-400 ml-2 transition-colors duration-200"
          title="Delete document"
        >
          🗑️
        </button>
      </div>
    </div>
  );
};

export default DocumentCard;