import React, { useState, useEffect, useCallback } from 'react';
import { createPortal } from 'react-dom';

const ToastContext = React.createContext();

const ToastProvider = ({ children }) => {
  const [toasts, setToasts] = useState([]);

  const addToast = useCallback((message, type = 'info', duration = 5000) => {
    const id = Date.now() + Math.random();
    const newToast = {
      id,
      message,
      type,
      duration,
    };
    
    setToasts((prevToasts) => [...prevToasts, newToast]);

    // Auto-remove toast after duration
    if (duration > 0) {
      setTimeout(() => {
        removeToast(id);
      }, duration);
    }

    return id;
  }, []);

  const removeToast = useCallback((id) => {
    setToasts((prevToasts) => prevToasts.filter((toast) => toast.id !== id));
  }, []);

  const value = { addToast, removeToast };

  return (
    <ToastContext.Provider value={value}>
      {children}
      {toasts.length > 0 && createPortal(
        <div className="fixed top-4 right-4 z-50 space-y-2">
          {toasts.map((toast) => (
            <ToastItem key={toast.id} toast={toast} removeToast={removeToast} />
          ))}
        </div>,
        document.body
      )}
    </ToastContext.Provider>
  );
};

const ToastItem = ({ toast, removeToast }) => {
  const { id, message, type } = toast;

  const getToastStyle = () => {
    switch (type) {
      case 'success':
        return 'bg-green-500 text-white';
      case 'error':
        return 'bg-red-500 text-white';
      case 'warning':
        return 'bg-yellow-500 text-gray-900';
      case 'info':
      default:
        return 'bg-blue-500 text-white';
    }
  };

  const getToastIcon = () => {
    switch (type) {
      case 'success':
        return '✅';
      case 'error':
        return '❌';
      case 'warning':
        return '⚠️';
      case 'info':
      default:
        return 'ℹ️';
    }
  };

  return (
    <div
      className={`flex items-center p-4 rounded-lg shadow-lg max-w-xs ${getToastStyle()} animate-fadeIn`}
      role="alert"
    >
      <span className="mr-2 text-lg">{getToastIcon()}</span>
      <span className="flex-grow">{message}</span>
      <button
        onClick={() => removeToast(id)}
        className="ml-4 text-lg hover:opacity-80 focus:outline-none"
        aria-label="Close toast"
      >
        ×
      </button>
    </div>
  );
};

const useToast = () => {
  const context = React.useContext(ToastContext);
  if (!context) {
    throw new Error('useToast must be used within a ToastProvider');
  }
  return context;
};

export { ToastProvider, useToast };