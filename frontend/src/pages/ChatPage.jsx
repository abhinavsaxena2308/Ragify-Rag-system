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
    <div className="flex flex-col h-[calc(100vh-200px)] max-w-4xl mx-auto w-full">
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center mb-6 gap-4">
        <div className="flex items-center space-x-3">
          <div className="w-10 h-10 bg-gradient-to-r from-red-600 to-red-700 rounded-full flex items-center justify-center">
            <span className="text-white font-bold">AI</span>
          </div>
          <div>
            <h1 className="text-xl font-bold text-white">Document Q&A Chat</h1>
            <p className="text-xs text-gray-400">Powered by RAG technology</p>
          </div>
        </div>
        <button
          onClick={() => navigate('/upload')}
          className="bg-gradient-to-r from-red-600 to-red-700 text-white px-4 py-2 rounded-lg hover:from-red-700 hover:to-red-800 text-sm font-medium transition-all duration-200 shadow-md"
        >
          Upload Documents
        </button>
      </div>
      
      {/* Chat messages container */}
      <div className="flex-grow overflow-y-auto mb-4 bg-gray-900 rounded-xl p-4 border border-gray-800">
        <div className="space-y-6 max-w-3xl mx-auto">
          {messages.map((message) => (
            <ChatMessage
              key={message.id}
              message={message}
              isUser={message.sender === 'user'}
              isLoading={message.isLoading}
            />
          ))}
        </div>
        <div ref={messagesEndRef} />
      </div>
      
      {/* Input area */}
      <form onSubmit={handleSubmit} className="mt-auto">
        <div className="flex flex-col">
          <div className="flex rounded-xl border border-gray-700 overflow-hidden focus-within:ring-2 focus-within:ring-red-500 focus-within:border-red-500 bg-gray-800 shadow-lg">
            <textarea
              ref={inputRef}
              value={inputValue}
              onChange={(e) => setInputValue(e.target.value)}
              onKeyDown={handleKeyDown}
              placeholder="Ask a question about your documents..."
              className="flex-grow min-h-[60px] max-h-32 px-4 py-3 resize-none focus:outline-none border-0 focus:ring-0 bg-transparent text-white placeholder-gray-500"
              disabled={isLoading}
              rows={1}
            />
            <button
              type="submit"
              disabled={isLoading || !inputValue.trim()}
              className={`self-end m-2 px-4 py-2 rounded-lg flex items-center justify-center font-medium transition-all duration-200 ${
                isLoading || !inputValue.trim()
                  ? 'bg-gray-700 text-gray-400 cursor-not-allowed'
                  : 'bg-gradient-to-r from-red-600 to-red-700 text-white hover:from-red-700 hover:to-red-800 shadow-md'
              }`}
            >
              {isLoading && <LoadingSpinner size="sm" className="mr-2" />}
              {isLoading ? 'Sending...' : 'Send'}
            </button>
          </div>
          <p className="text-xs text-gray-500 mt-2 text-center">
            Ask questions about your uploaded documents. The AI will reference the relevant sources.
          </p>
        </div>
      </form>
    </div>
  );
};

export default ChatPage;