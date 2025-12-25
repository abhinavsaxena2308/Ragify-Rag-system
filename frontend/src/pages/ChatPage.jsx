import React, { useState, useRef, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import ChatMessage from '../components/ChatMessage';
import chatService from '../services/chatService';
import { useToast } from '../components/Toast';
import LoadingSpinner from '../components/LoadingSpinner';

const ChatPage = () => {
  const [messages, setMessages] = useState([
    {
      id: 1,
      text: "Hello! I'm your RAGify assistant. Upload some documents and ask me questions about them.",
      sender: 'ai',
      timestamp: new Date(),
      sources: []
    }
  ]);
  const [inputValue, setInputValue] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const { addToast } = useToast();
  const messagesEndRef = useRef(null);
  const inputRef = useRef(null);
  const navigate = useNavigate();

  // Auto-scroll to bottom when messages change
  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!inputValue.trim() || isLoading) return;

    try {
      // Add user message to the chat
      const userMessage = {
        id: Date.now(),
        text: inputValue,
        sender: 'user',
        timestamp: new Date()
      };

      setMessages(prev => [...prev, userMessage]);
      setInputValue('');
      setIsLoading(true);

      // Add a temporary loading message for AI response
      const aiMessageId = Date.now() + 1;
      const loadingMessage = {
        id: aiMessageId,
        text: '',
        sender: 'ai',
        timestamp: new Date(),
        sources: [],
        isLoading: true
      };

      setMessages(prev => [...prev, loadingMessage]);

      // Call the actual API service
      const response = await chatService.askQuestion(inputValue);
      
      const aiResponse = {
        id: aiMessageId,
        text: response.answer || 'No response generated',
        sender: 'ai',
        timestamp: new Date(),
        sources: response.sources || [] // Use sources from the API response
      };

      // Update the loading message with the actual response
      setMessages(prev => 
        prev.map(msg => 
          msg.id === aiMessageId ? aiResponse : msg
        )
      );

    } catch (err) {
      // Remove the loading message
      setMessages(prev => prev.filter(msg => msg.id !== aiMessageId));
      
      const errorMessage = {
        id: Date.now(),
        text: 'Sorry, I encountered an error processing your request. Please try again.',
        sender: 'ai',
        timestamp: new Date(),
        sources: []
      };
      
      setMessages(prev => [...prev, errorMessage]);
      addToast(err.message || 'An error occurred while processing your question', 'error');
    } finally {
      setIsLoading(false);
    }
  };

  const handleKeyDown = (e) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSubmit(e);
    }
  };

  return (
    <div className="flex flex-col h-full bg-gray-900">
      {/* Chat Messages Container */}
      <div className="flex-1 overflow-y-auto">
        <div className="max-w-3xl mx-auto px-4 py-6">
          {messages.length === 1 && messages[0].sender === 'ai' ? (
            // Welcome state
            <div className="flex flex-col items-center justify-center h-full text-center py-20">
              <div className="w-16 h-16 bg-gradient-to-r from-red-600 to-red-700 rounded-full flex items-center justify-center mb-6">
                <span className="text-white font-bold text-2xl">R</span>
              </div>
              <h1 className="text-3xl font-bold text-white mb-4">Welcome to RAGify</h1>
              <p className="text-gray-400 text-lg mb-8 max-w-md">
                Upload documents and ask questions about them. I'll provide answers based on your uploaded content.
              </p>
              <button
                onClick={() => navigate('/upload')}
                className="bg-gradient-to-r from-red-600 to-red-700 text-white px-6 py-3 rounded-lg hover:from-red-700 hover:to-red-800 font-medium transition-all duration-200 shadow-lg"
              >
                Upload Documents
              </button>
            </div>
          ) : (
            // Chat messages
            <div className="space-y-6">
              {messages.map((message) => (
                <ChatMessage
                  key={message.id}
                  message={message}
                  isUser={message.sender === 'user'}
                  isLoading={message.isLoading}
                />
              ))}
              <div ref={messagesEndRef} />
            </div>
          )}
        </div>
      </div>
      
      {/* Input Area */}
      <div className="border-t border-gray-800 bg-gray-950">
        <div className="max-w-3xl mx-auto px-4 py-4">
          <form onSubmit={handleSubmit} className="relative">
            <div className="flex items-end space-x-2">
              <div className="flex-1 relative">
                <textarea
                  ref={inputRef}
                  value={inputValue}
                  onChange={(e) => setInputValue(e.target.value)}
                  onKeyDown={handleKeyDown}
                  placeholder="Ask a question about your documents..."
                  className="w-full px-4 py-3 pr-12 bg-gray-800 text-white rounded-lg resize-none focus:outline-none focus:ring-2 focus:ring-red-500 border border-gray-700 placeholder-gray-500 transition-all duration-200"
                  disabled={isLoading}
                  rows={1}
                  style={{
                    minHeight: '52px',
                    maxHeight: '200px',
                    height: 'auto'
                  }}
                  onInput={(e) => {
                    e.target.style.height = 'auto';
                    e.target.style.height = Math.min(e.target.scrollHeight, 200) + 'px';
                  }}
                />
                <button
                  type="button"
                  className="absolute right-2 bottom-2 p-2 text-gray-400 hover:text-gray-200 transition-colors duration-200"
                  title="Attach file"
                >
                  <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15.172 7l-6.586 6.586a2 2 0 102.828 2.828l6.414-6.586a4 4 0 00-5.656-5.656l-6.415 6.585a6 6 0 108.486 8.486L20.5 13" />
                  </svg>
                </button>
              </div>
              
              <button
                type="submit"
                disabled={isLoading || !inputValue.trim()}
                className={`px-4 py-3 rounded-lg flex items-center justify-center font-medium transition-all duration-200 ${
                  isLoading || !inputValue.trim()
                    ? 'bg-gray-700 text-gray-400 cursor-not-allowed'
                    : 'bg-gradient-to-r from-red-600 to-red-700 text-white hover:from-red-700 hover:to-red-800 shadow-lg'
                }`}
              >
                {isLoading ? (
                  <LoadingSpinner size="sm" />
                ) : (
                  <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8" />
                  </svg>
                )}
              </button>
            </div>
            
            <div className="flex items-center justify-between mt-2 text-xs text-gray-500">
              <span>RAGify can make mistakes. Consider checking important information.</span>
              <span>{inputValue.length}/4000</span>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
};

export default ChatPage;