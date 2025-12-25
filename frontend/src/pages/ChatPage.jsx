import React, { useState, useRef, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import chatService from '../services/chatService';

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
        text: 'Thinking...',
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
    <div style={{ display: 'flex', flexDirection: 'column', height: '100%', fontFamily: 'Arial, sans-serif' }}>
      {/* Chat Messages Container */}
      <div style={{ flex: 1, overflowY: 'auto', padding: '20px' }}>
        <div style={{ maxWidth: '800px', margin: '0 auto' }}>
          {messages.length === 1 && messages[0].sender === 'ai' ? (
            // Welcome state
            <div style={{ textAlign: 'center', padding: '40px' }}>
              <div style={{ fontSize: '48px', marginBottom: '20px' }}>R</div>
              <h1 style={{ fontSize: '32px', marginBottom: '16px' }}>Welcome to RAGify</h1>
              <p style={{ color: '#666', fontSize: '18px', marginBottom: '32px', maxWidth: '400px', margin: '0 auto 32px' }}>
                Upload documents and ask questions about them. I'll provide answers based on your uploaded content.
              </p>
              <button
                onClick={() => navigate('/upload')}
                style={{
                  padding: '12px 24px',
                  backgroundColor: '#007bff',
                  color: 'white',
                  border: 'none',
                  borderRadius: '4px',
                  cursor: 'pointer',
                  fontSize: '16px'
                }}
              >
                Upload Documents
              </button>
            </div>
          ) : (
            // Chat messages
            <div>
              {messages.map((message) => (
                <div key={message.id} style={{ 
                  marginBottom: '20px',
                  display: 'flex',
                  justifyContent: message.sender === 'user' ? 'flex-end' : 'flex-start'
                }}>
                  <div style={{
                    maxWidth: '70%',
                    padding: '12px 16px',
                    backgroundColor: message.sender === 'user' ? '#007bff' : '#f1f3f4',
                    color: message.sender === 'user' ? 'white' : 'black',
                    borderRadius: '8px',
                    border: '1px solid #ddd'
                  }}>
                    {message.isLoading ? (
                      <span>{message.text}</span>
                    ) : (
                      <div>
                        <div>{message.text}</div>
                        {message.sources && message.sources.length > 0 && (
                          <div style={{ marginTop: '12px', paddingTop: '12px', borderTop: '1px solid #ddd' }}>
                            <div style={{ fontSize: '12px', fontWeight: 'bold', marginBottom: '8px' }}>Sources:</div>
                            {message.sources.map((source, index) => (
                              <div key={index} style={{ fontSize: '12px', color: '#666', marginBottom: '4px' }}>
                                📄 {source}
                              </div>
                            ))}
                          </div>
                        )}
                      </div>
                    )}
                  </div>
                </div>
              ))}
              <div ref={messagesEndRef} />
            </div>
          )}
        </div>
      </div>
      
      {/* Input Area */}
      <div style={{ borderTop: '1px solid #ddd', backgroundColor: '#f8f9fa', padding: '16px' }}>
        <div style={{ maxWidth: '800px', margin: '0 auto' }}>
          <form onSubmit={handleSubmit} style={{ display: 'flex', gap: '12px' }}>
            <textarea
              ref={inputRef}
              value={inputValue}
              onChange={(e) => setInputValue(e.target.value)}
              onKeyDown={handleKeyDown}
              placeholder="Ask a question about your documents..."
              style={{
                flex: 1,
                padding: '12px',
                border: '1px solid #ddd',
                borderRadius: '4px',
                fontSize: '14px',
                resize: 'none',
                minHeight: '44px',
                maxHeight: '120px'
              }}
              disabled={isLoading}
            />
            <button
              type="submit"
              disabled={isLoading || !inputValue.trim()}
              style={{
                padding: '12px 20px',
                backgroundColor: isLoading || !inputValue.trim() ? '#6c757d' : '#007bff',
                color: 'white',
                border: 'none',
                borderRadius: '4px',
                cursor: isLoading || !inputValue.trim() ? 'not-allowed' : 'pointer',
                fontSize: '14px'
              }}
            >
              {isLoading ? 'Sending...' : 'Send'}
            </button>
          </form>
          <div style={{ display: 'flex', justifyContent: 'space-between', marginTop: '8px', fontSize: '12px', color: '#666' }}>
            <span>RAGify can make mistakes. Consider checking important information.</span>
            <span>{inputValue.length}/4000</span>
          </div>
        </div>
      </div>
    </div>
  );
};

export default ChatPage;