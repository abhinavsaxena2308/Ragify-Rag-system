import React from 'react';
import { Link } from 'react-router-dom';

const HomePage = () => {
  return (
    <div className="text-center w-full">
      <div className="mb-12">
        <h1 className="text-4xl md:text-5xl font-bold text-white mb-4 bg-gradient-to-r from-red-400 to-red-300 bg-clip-text text-transparent">Welcome to RAGify</h1>
        <p className="text-lg text-gray-400 max-w-2xl mx-auto">
          Your AI-powered Document Question Answering System
        </p>
      </div>
      
      <div className="grid grid-cols-1 md:grid-cols-3 gap-8 mb-16">
        <div className="bg-gray-800 p-8 rounded-2xl border border-gray-700 hover:border-red-500 transition-all duration-300 group">
          <div className="text-5xl mb-6 text-red-500 group-hover:text-red-400 transition-colors duration-300">📄</div>
          <h2 className="text-xl font-bold text-white mb-4">Upload Documents</h2>
          <p className="text-gray-400 group-hover:text-gray-300 transition-colors duration-300">
            Upload your documents and let our AI process them for intelligent Q&A
          </p>
        </div>
        <div className="bg-gray-800 p-8 rounded-2xl border border-gray-700 hover:border-red-500 transition-all duration-300 group">
          <div className="text-5xl mb-6 text-red-500 group-hover:text-red-400 transition-colors duration-300">❓</div>
          <h2 className="text-xl font-bold text-white mb-4">Ask Questions</h2>
          <p className="text-gray-400 group-hover:text-gray-300 transition-colors duration-300">
            Ask questions about your documents and get precise answers
          </p>
        </div>
        <div className="bg-gray-800 p-8 rounded-2xl border border-gray-700 hover:border-red-500 transition-all duration-300 group">
          <div className="text-5xl mb-6 text-red-500 group-hover:text-red-400 transition-colors duration-300">💡</div>
          <h2 className="text-xl font-bold text-white mb-4">Get Insights</h2>
          <p className="text-gray-400 group-hover:text-gray-300 transition-colors duration-300">
            Extract valuable insights from your document collections
          </p>
        </div>
      </div>
      
      <div className="max-w-3xl mx-auto">
        <div className="bg-gradient-to-r from-gray-800 to-gray-900 border border-gray-700 rounded-2xl p-8 shadow-xl">
          <div className="flex flex-col md:flex-row items-center justify-between gap-6">
            <div className="text-left">
              <h2 className="text-2xl font-bold text-white mb-2">Get Started Today</h2>
              <p className="text-gray-400 mb-4">Upload your first document and start asking questions to unlock the power of AI-powered document analysis.</p>
              <Link 
                to="/upload" 
                className="inline-block bg-gradient-to-r from-red-600 to-red-700 text-white font-bold px-6 py-3 rounded-lg hover:from-red-700 hover:to-red-800 transition-all duration-200 shadow-lg"
              >
                Upload Document
              </Link>
            </div>
            <div className="text-6xl">🚀</div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default HomePage;