import React from 'react';

const HomePage = () => {
  return (
    <div className="text-center">
      <h1 className="text-3xl font-bold text-black mb-4">Welcome to RAGify</h1>
      <p className="text-lg text-gray-700 mb-8">
        Your AI-powered Document Question Answering System
      </p>
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="bg-white p-6 rounded-lg shadow-lg border border-gray-200 hover:shadow-xl transition-shadow duration-300">
          <div className="text-4xl mb-4 text-red-600">📄</div>
          <h2 className="text-xl font-semibold text-black mb-2">Upload Documents</h2>
          <p className="text-gray-600">
            Upload your documents and let our AI process them for intelligent Q&A
          </p>
        </div>
        <div className="bg-white p-6 rounded-lg shadow-lg border border-gray-200 hover:shadow-xl transition-shadow duration-300">
          <div className="text-4xl mb-4 text-red-600">❓</div>
          <h2 className="text-xl font-semibold text-black mb-2">Ask Questions</h2>
          <p className="text-gray-600">
            Ask questions about your documents and get precise answers
          </p>
        </div>
        <div className="bg-white p-6 rounded-lg shadow-lg border border-gray-200 hover:shadow-xl transition-shadow duration-300">
          <div className="text-4xl mb-4 text-red-600">💡</div>
          <h2 className="text-xl font-semibold text-black mb-2">Get Insights</h2>
          <p className="text-gray-600">
            Extract valuable insights from your document collections
          </p>
        </div>
      </div>
      <div className="mt-12 max-w-2xl mx-auto">
        <div className="bg-gradient-to-r from-red-600 to-red-700 text-white p-8 rounded-xl shadow-lg">
          <h2 className="text-2xl font-bold mb-4">Get Started Today</h2>
          <p className="mb-6">Upload your first document and start asking questions to unlock the power of AI-powered document analysis.</p>
          <a 
            href="/upload" 
            className="inline-block bg-white text-red-600 font-bold px-6 py-3 rounded-lg hover:bg-gray-100 transition-colors duration-200"
          >
            Upload Document
          </a>
        </div>
      </div>
    </div>
  );
};

export default HomePage;