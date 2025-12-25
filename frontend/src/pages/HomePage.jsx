import React from 'react';
import { Link } from 'react-router-dom';

const HomePage = () => {
  return (
    <div className="flex flex-col items-center justify-center h-full px-4 py-12">
      {/* Hero Section */}
      <div className="text-center mb-16 max-w-4xl">
        <div className="mb-8">
          <div className="w-20 h-20 bg-gradient-to-r from-red-600 to-red-700 rounded-full flex items-center justify-center mx-auto mb-6 shadow-2xl">
            <span className="text-white font-bold text-3xl">R</span>
          </div>
          <h1 className="text-5xl md:text-6xl font-bold text-white mb-6 bg-gradient-to-r from-red-400 to-red-300 bg-clip-text text-transparent">
            Welcome to RAGify
          </h1>
          <p className="text-xl text-gray-400 max-w-2xl mx-auto leading-relaxed">
            Your AI-powered Document Question Answering System
          </p>
        </div>
        
        <div className="flex flex-col sm:flex-row gap-4 justify-center">
          <Link 
            to="/upload" 
            className="bg-gradient-to-r from-red-600 to-red-700 text-white font-semibold px-8 py-4 rounded-lg hover:from-red-700 hover:to-red-800 transition-all duration-200 shadow-lg text-lg"
          >
            Upload Documents
          </Link>
          <Link 
            to="/chat" 
            className="bg-gray-800 text-white font-semibold px-8 py-4 rounded-lg hover:bg-gray-700 transition-all duration-200 shadow-lg border border-gray-700 text-lg"
          >
            Start Chatting
          </Link>
        </div>
      </div>
      
      {/* Features Grid */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-8 w-full max-w-6xl mb-16">
        <div className="bg-gray-800/50 backdrop-blur-sm p-8 rounded-2xl border border-gray-700 hover:border-red-500 transition-all duration-300 group hover:shadow-2xl">
          <div className="w-16 h-16 bg-gradient-to-r from-red-600 to-red-700 rounded-xl flex items-center justify-center mb-6 group-hover:scale-110 transition-transform duration-300">
            <span className="text-2xl">📄</span>
          </div>
          <h2 className="text-2xl font-bold text-white mb-4">Upload Documents</h2>
          <p className="text-gray-400 leading-relaxed">
            Upload your documents and let our AI process them for intelligent Q&A with advanced RAG technology
          </p>
        </div>
        
        <div className="bg-gray-800/50 backdrop-blur-sm p-8 rounded-2xl border border-gray-700 hover:border-red-500 transition-all duration-300 group hover:shadow-2xl">
          <div className="w-16 h-16 bg-gradient-to-r from-red-600 to-red-700 rounded-xl flex items-center justify-center mb-6 group-hover:scale-110 transition-transform duration-300">
            <span className="text-2xl">❓</span>
          </div>
          <h2 className="text-2xl font-bold text-white mb-4">Ask Questions</h2>
          <p className="text-gray-400 leading-relaxed">
            Ask questions about your documents and get precise, context-aware answers powered by AI
          </p>
        </div>
        
        <div className="bg-gray-800/50 backdrop-blur-sm p-8 rounded-2xl border border-gray-700 hover:border-red-500 transition-all duration-300 group hover:shadow-2xl">
          <div className="w-16 h-16 bg-gradient-to-r from-red-600 to-red-700 rounded-xl flex items-center justify-center mb-6 group-hover:scale-110 transition-transform duration-300">
            <span className="text-2xl">💡</span>
          </div>
          <h2 className="text-2xl font-bold text-white mb-4">Get Insights</h2>
          <p className="text-gray-400 leading-relaxed">
            Extract valuable insights and information from your document collections with intelligent analysis
          </p>
        </div>
      </div>
      
      {/* Quick Stats or Info */}
      <div className="bg-gradient-to-r from-gray-800/50 to-gray-900/50 backdrop-blur-sm border border-gray-700 rounded-2xl p-8 w-full max-w-4xl">
        <div className="text-center">
          <h2 className="text-3xl font-bold text-white mb-4">Ready to Get Started?</h2>
          <p className="text-gray-400 mb-6 text-lg">
            Upload your first document and start asking questions to unlock the power of AI-powered document analysis.
          </p>
          <div className="flex flex-col sm:flex-row gap-4 justify-center">
            <Link 
              to="/upload" 
              className="bg-gradient-to-r from-red-600 to-red-700 text-white font-semibold px-8 py-3 rounded-lg hover:from-red-700 hover:to-red-800 transition-all duration-200 shadow-lg"
            >
              Upload Document
            </Link>
            <Link 
              to="/documents" 
              className="bg-gray-800 text-white font-semibold px-8 py-3 rounded-lg hover:bg-gray-700 transition-all duration-200 shadow-lg border border-gray-700"
            >
              View Documents
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
};

export default HomePage;