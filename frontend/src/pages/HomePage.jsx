import React from 'react';

const HomePage = () => {
  return (
    <div className="text-center">
      <h1 className="text-3xl font-bold text-gray-900 mb-4">Welcome to RAGify</h1>
      <p className="text-lg text-gray-600 mb-8">
        Your AI-powered Document Question Answering System
      </p>
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="bg-white p-6 rounded-lg shadow-md">
          <h2 className="text-xl font-semibold mb-2">Upload Documents</h2>
          <p className="text-gray-600">
            Upload your documents and let our AI process them for intelligent Q&A
          </p>
        </div>
        <div className="bg-white p-6 rounded-lg shadow-md">
          <h2 className="text-xl font-semibold mb-2">Ask Questions</h2>
          <p className="text-gray-600">
            Ask questions about your documents and get precise answers
          </p>
        </div>
        <div className="bg-white p-6 rounded-lg shadow-md">
          <h2 className="text-xl font-semibold mb-2">Get Insights</h2>
          <p className="text-gray-600">
            Extract valuable insights from your document collections
          </p>
        </div>
      </div>
    </div>
  );
};

export default HomePage;