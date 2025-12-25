import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import LoadingSpinner from '../components/LoadingSpinner';

const Layout = ({ children }) => {
  const [isMenuOpen, setIsMenuOpen] = useState(false);

  const toggleMenu = () => {
    setIsMenuOpen(!isMenuOpen);
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-900 to-black text-gray-100 flex flex-col">
      {/* Header */}
      <header className="bg-gray-900 border-b border-gray-800 backdrop-blur-lg bg-opacity-90 sticky top-0 z-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <div className="flex items-center space-x-3">
              <div className="w-8 h-8 bg-gradient-to-r from-red-600 to-red-700 rounded-lg flex items-center justify-center shadow-lg">
                <span className="text-white font-bold text-lg">R</span>
              </div>
              <h1 className="text-xl font-bold bg-gradient-to-r from-red-500 to-red-400 bg-clip-text text-transparent">RAGify</h1>
            </div>
            
            {/* Desktop Navigation */}
            <nav className="hidden md:flex space-x-1">
              <Link 
                to="/" 
                className="nav-link px-3 py-2 rounded-lg text-sm font-medium transition-all duration-200 hover:bg-gray-800 hover:text-red-400"
              >
                <span className="flex items-center">
                  <span className="mr-2">🏠</span>
                  Home
                </span>
              </Link>
              <Link 
                to="/upload" 
                className="nav-link px-3 py-2 rounded-lg text-sm font-medium transition-all duration-200 hover:bg-gray-800 hover:text-red-400"
              >
                <span className="flex items-center">
                  <span className="mr-2">📤</span>
                  Upload
                </span>
              </Link>
              <Link 
                to="/documents" 
                className="nav-link px-3 py-2 rounded-lg text-sm font-medium transition-all duration-200 hover:bg-gray-800 hover:text-red-400"
              >
                <span className="flex items-center">
                  <span className="mr-2">📁</span>
                  Docs
                </span>
              </Link>
              <Link 
                to="/chat" 
                className="nav-link px-3 py-2 rounded-lg text-sm font-medium transition-all duration-200 hover:bg-gray-800 hover:text-red-400"
              >
                <span className="flex items-center">
                  <span className="mr-2">💬</span>
                  Chat
                </span>
              </Link>
            </nav>
            
            {/* Mobile menu button */}
            <div className="md:hidden flex items-center">
              <button
                onClick={toggleMenu}
                className="inline-flex items-center justify-center p-2 rounded-lg text-gray-400 hover:text-white hover:bg-gray-800 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-red-500 transition-all duration-200"
                aria-expanded="false"
              >
                <span className="sr-only">Open main menu</span>
                <svg 
                  className={`${isMenuOpen ? 'hidden' : 'block'} h-6 w-6`} 
                  xmlns="http://www.w3.org/2000/svg" 
                  fill="none" 
                  viewBox="0 0 24 24" 
                  stroke="currentColor"
                >
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
                </svg>
                <svg 
                  className={`${isMenuOpen ? 'block' : 'hidden'} h-6 w-6`} 
                  xmlns="http://www.w3.org/2000/svg" 
                  fill="none" 
                  viewBox="0 0 24 24" 
                  stroke="currentColor"
                >
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
          </div>
        </div>
        
        {/* Mobile Navigation */}
        <div className={`${isMenuOpen ? 'block' : 'hidden'} md:hidden bg-gray-900 border-b border-gray-800`} id="mobile-menu">
          <div className="px-2 pt-2 pb-3 space-y-1 sm:px-3">
            <Link 
              to="/" 
              className="nav-link block px-3 py-2 rounded-lg text-base font-medium transition-all duration-200 hover:bg-gray-800 hover:text-red-400"
              onClick={() => setIsMenuOpen(false)}
            >
              <span className="flex items-center">
                <span className="mr-2">🏠</span>
                Home
              </span>
            </Link>
            <Link 
              to="/upload" 
              className="nav-link block px-3 py-2 rounded-lg text-base font-medium transition-all duration-200 hover:bg-gray-800 hover:text-red-400"
              onClick={() => setIsMenuOpen(false)}
            >
              <span className="flex items-center">
                <span className="mr-2">📤</span>
                Upload
              </span>
            </Link>
            <Link 
              to="/documents" 
              className="nav-link block px-3 py-2 rounded-lg text-base font-medium transition-all duration-200 hover:bg-gray-800 hover:text-red-400"
              onClick={() => setIsMenuOpen(false)}
            >
              <span className="flex items-center">
                <span className="mr-2">📁</span>
                Docs
              </span>
            </Link>
            <Link 
              to="/chat" 
              className="nav-link block px-3 py-2 rounded-lg text-base font-medium transition-all duration-200 hover:bg-gray-800 hover:text-red-400"
              onClick={() => setIsMenuOpen(false)}
            >
              <span className="flex items-center">
                <span className="mr-2">💬</span>
                Chat
              </span>
            </Link>
          </div>
        </div>
      </header>

      {/* Main Content Area with Sidebar Layout */}
      <div className="flex flex-1 max-w-7xl mx-auto w-full px-4 sm:px-6 lg:px-8 py-6 flex-grow">
        <main className="w-full">
          {children}
        </main>
      </div>

      {/* Footer */}
      <footer className="bg-gray-900 border-t border-gray-800 mt-auto">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
          <div className="flex flex-col md:flex-row justify-between items-center space-y-4 md:space-y-0">
            <p className="text-gray-500 text-sm">
              &copy; {new Date().getFullYear()} RAGify - Document Question Answering System
            </p>
            <div className="flex space-x-6">
              <a href="#" className="text-gray-500 hover:text-red-400 transition-colors duration-200 text-sm">Privacy</a>
              <a href="#" className="text-gray-500 hover:text-red-400 transition-colors duration-200 text-sm">Terms</a>
              <a href="#" className="text-gray-500 hover:text-red-400 transition-colors duration-200 text-sm">Support</a>
            </div>
          </div>
        </div>
      </footer>
    </div>
  );
};

export default Layout;