import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Layout from './layouts/Layout';
import HomePage from './pages/HomePage';
import DocumentsPage from './pages/DocumentsPage';
import ChatPage from './pages/ChatPage';
import UploadPage from './pages/UploadPage';
import { ToastProvider } from './components/Toast';
import ErrorBoundary from './components/ErrorBoundary';

function App() {
  return (
    <ToastProvider>
      <Router>
        <ErrorBoundary>
          <Layout>
            <Routes>
              <Route path="/" element={<HomePage />} />
              <Route path="/documents" element={<DocumentsPage />} />
              <Route path="/chat" element={<ChatPage />} />
              <Route path="/upload" element={<UploadPage />} />
            </Routes>
          </Layout>
        </ErrorBoundary>
      </Router>
    </ToastProvider>
  );
}

export default App
