import React, { useState } from 'react';
import { Link, useLocation } from 'react-router-dom';

const Layout = ({ children }) => {
  const [isSidebarOpen, setIsSidebarOpen] = useState(false);
  const location = useLocation();

  const toggleSidebar = () => {
    setIsSidebarOpen(!isSidebarOpen);
  };

  const navigation = [
    { name: 'New Chat', path: '/chat' },
    { name: 'Upload Documents', path: '/upload' },
    { name: 'My Documents', path: '/documents' },
    { name: 'Home', path: '/' },
  ];

  return (
    <div style={{ display: 'flex', height: '100vh', fontFamily: 'Arial, sans-serif', backgroundColor: '#fff' }}>
      {/* Sidebar */}
      <aside style={{
        position: 'fixed',
        left: '0',
        width: '256px',
        height: '100%',
        backgroundColor: '#f5f5f5',
        border: '1px solid #ddd',
        zIndex: 40
      }}>
        {/* Sidebar Header */}
        <div style={{ padding: '16px', borderBottom: '1px solid #ddd' }}>
          <h1 style={{ margin: 0, fontSize: '18px', fontWeight: 'bold' }}>RAGify</h1>
        </div>

        {/* Navigation */}
        <nav style={{ padding: '16px' }}>
          {navigation.map((item) => (
            <Link
              key={item.name}
              to={item.path}
              onClick={() => setIsSidebarOpen(false)}
              style={{
                display: 'block',
                padding: '8px 12px',
                marginBottom: '8px',
                textDecoration: 'none',
                color: location.pathname === item.path ? '#000' : '#666',
                backgroundColor: location.pathname === item.path ? '#e0e0e0' : 'transparent',
                borderRadius: '4px'
              }}
            >
              {item.name}
            </Link>
          ))}
        </nav>
      </aside>

      {/* Main Content */}
      <div style={{ flex: 1, marginLeft: '256px' }}>
        {/* Top Bar */}
        <header style={{ 
          backgroundColor: '#fff', 
          borderBottom: '1px solid #ddd', 
          padding: '12px 16px',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'space-between'
        }}>
          <button
            onClick={toggleSidebar}
            style={{ 
              background: 'none', 
              border: '1px solid #ccc', 
              padding: '8px',
              cursor: 'pointer',
              display: isSidebarOpen ? 'none' : 'block'
            }}
          >
            Menu
          </button>
          
          <h2 style={{ margin: 0, fontSize: '16px' }}>
            {navigation.find(item => item.path === location.pathname)?.name || 'RAGify'}
          </h2>

          <div style={{ width: '40px' }}></div>
        </header>

        {/* Main Content Area */}
        <main style={{ padding: '20px', overflow: 'auto', height: 'calc(100vh - 60px)' }}>
          {children}
        </main>
      </div>

      {/* Mobile Sidebar Overlay */}
      {isSidebarOpen && (
        <div
          style={{
            position: 'fixed',
            top: 0,
            left: 0,
            right: 0,
            bottom: 0,
            backgroundColor: 'rgba(0,0,0,0.5)',
            zIndex: 30
          }}
          onClick={toggleSidebar}
        />
      )}
    </div>
  );
};

export default Layout;