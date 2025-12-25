import React, { useState, useRef, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import ChatMessage from '../components/ChatMessage';
import { chatService } from '../services/api';

const ChatPage = () => {
  const [messages, setMessages] = useState([
    {
      id: 1,
      text: 'Hello! I\'m your RAGify assistant. Upload some documents and ask me questions about them.',
      sender: 'ai',
      timestamp: new Date(),
      sources: []
    }
  ]);
  const [inputValue, setInputValue] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState(null);
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
      setError(null);

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

      // In a real implementation, we would call the API:
      // const response = await chatService.askQuestion(inputValue);
      // For now, simulate an API response with mock data
      setTimeout(() => {
        const aiResponse = {
          id: aiMessageId,
          text: `I understand you're asking about "${inputValue}". This is a simulated response from the RAGify system. In a real implementation, this would be generated based on your uploaded documents.`,
          sender: 'ai',
          timestamp: new Date(),
          sources: ['Company Handbook.pdf', 'Financial Report Q3.docx'] // Mock sources
        };

        // Update the loading message with the actual response
        setMessages(prev => 
          prev.map(msg => 
            msg.id === aiMessageId ? aiResponse : msg
          )
        );
        setIsLoading(false);
      }, 1500);

    } catch (err) {
      // Remove the loading message
      setMessages(prev => prev.filter(msg => !msg.isLoading));
      
      const errorMessage = {
        id: Date.now(),
        text: 'Sorry, I encountered an error processing your request. Please try again.',
        sender: 'ai',
        timestamp: new Date(),
        sources: []
      };
      
      setMessages(prev => [...prev, errorMessage]);
      setError(err.message);
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
    <div className="flex flex-col h-[calc(100vh-200px)] max-w-4xl mx-auto">
      <div className="flex justify-between items-center mb-4">
        <h1 className="text-2xl font-bold text-gray-900">Document Q&A Chat</h1>
        <button
          onClick={() => navigate('/upload')}
          className="bg-indigo-600 text-white px-4 py-2 rounded-md hover:bg-indigo-700 text-sm"
        >
          Upload Documents
        </button>
      </div>
      
      {/* Chat messages container */}
      <div className="flex-grow overflow-y-auto mb-4 pr-2">
        <div className="space-y-2">
          {messages.map((message) => (
            <ChatMessage
              key={message.id}
              message={message}
              isUser={message.sender === 'user'}
              isLoading={message.isLoading}
            />
          ))}
          {error && (
            <div className="flex justify-start mb-4">
              <div className="max-w-[85%] md:max-w-[75%] bg-red-100 text-red-800 rounded-2xl px-4 py-3 rounded-tl-none">
                Error: {error}
              </div>
            </div>
          )}
        </div>
        <div ref={messagesEndRef} />
      </div>
      
      {/* Input area */}
      <form onSubmit={handleSubmit} className="mt-auto">
        <div className="flex flex-col">
          <div className="flex rounded-lg border border-gray-300 overflow-hidden focus-within:ring-2 focus-within:ring-indigo-500 focus-within:border-indigo-500">
            <textarea
              ref={inputRef}
              value={inputValue}
              onChange={(e) => setInputValue(e.target.value)}
              onKeyDown={handleKeyDown}
              placeholder="Ask a question about your documents..."
              className="flex-grow min-h-[60px] max-h-32 px-4 py-2 resize-none focus:outline-none"
              disabled={isLoading}
              rows={1}
            />
            <button
              type="submit"
              disabled={isLoading || !inputValue.trim()}
              className={`self-end m-2 px-4 py-2 rounded-md ${
                isLoading || !inputValue.trim()
                  ? 'bg-gray-300 text-gray-500 cursor-not-allowed'
                  : 'bg-indigo-600 text-white hover:bg-indigo-700'
              }`}
            >
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