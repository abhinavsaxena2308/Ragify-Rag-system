import React from 'react';
import { Link } from 'react-router-dom';

const HomePage = () => {
  return (
    <div style={{ padding: '40px', textAlign: 'center', fontFamily: 'Arial, sans-serif' }}>
      <div style={{ marginBottom: '40px' }}>
        <h1 style={{ fontSize: '36px', marginBottom: '16px' }}>Welcome to RAGify</h1>
        <p style={{ fontSize: '18px', color: '#666', maxWidth: '600px', margin: '0 auto' }}>
          Your AI-powered Document Question Answering System
        </p>
      </div>
      
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(300px, 1fr))', gap: '24px', marginBottom: '40px', maxWidth: '900px', margin: '0 auto 40px' }}>
        <div style={{ border: '1px solid #ddd', padding: '24px', borderRadius: '8px', backgroundColor: '#f9f9f9' }}>
          <div style={{ fontSize: '32px', marginBottom: '16px' }}>📄</div>
          <h2 style={{ fontSize: '20px', marginBottom: '12px' }}>Upload Documents</h2>
          <p style={{ color: '#666', lineHeight: '1.5' }}>
            Upload your documents and let our AI process them for intelligent Q&A with advanced RAG technology
          </p>
        </div>
        
        <div style={{ border: '1px solid #ddd', padding: '24px', borderRadius: '8px', backgroundColor: '#f9f9f9' }}>
          <div style={{ fontSize: '32px', marginBottom: '16px' }}>❓</div>
          <h2 style={{ fontSize: '20px', marginBottom: '12px' }}>Ask Questions</h2>
          <p style={{ color: '#666', lineHeight: '1.5' }}>
            Ask questions about your documents and get precise, context-aware answers powered by AI
          </p>
        </div>
        
        <div style={{ border: '1px solid #ddd', padding: '24px', borderRadius: '8px', backgroundColor: '#f9f9f9' }}>
          <div style={{ fontSize: '32px', marginBottom: '16px' }}>💡</div>
          <h2 style={{ fontSize: '20px', marginBottom: '12px' }}>Get Insights</h2>
          <p style={{ color: '#666', lineHeight: '1.5' }}>
            Extract valuable insights and information from your document collections with intelligent analysis
          </p>
        </div>
      </div>
      
      <div style={{ border: '1px solid #ddd', padding: '32px', borderRadius: '8px', backgroundColor: '#f9f9f9', maxWidth: '600px', margin: '0 auto' }}>
        <h2 style={{ fontSize: '24px', marginBottom: '16px' }}>Ready to Get Started?</h2>
        <p style={{ color: '#666', marginBottom: '24px', lineHeight: '1.5' }}>
          Upload your first document and start asking questions to unlock the power of AI-powered document analysis.
        </p>
        <div style={{ display: 'flex', gap: '16px', justifyContent: 'center' }}>
          <Link 
            to="/upload" 
            style={{ 
              padding: '12px 24px', 
              backgroundColor: '#007bff', 
              color: 'white', 
              textDecoration: 'none', 
              borderRadius: '4px',
              border: 'none',
              cursor: 'pointer'
            }}
          >
            Upload Document
          </Link>
          <Link 
            to="/documents" 
            style={{ 
              padding: '12px 24px', 
              backgroundColor: '#6c757d', 
              color: 'white', 
              textDecoration: 'none', 
              borderRadius: '4px',
              border: 'none',
              cursor: 'pointer'
            }}
          >
            View Documents
          </Link>
        </div>
      </div>
    </div>
  );
};

export default HomePage;