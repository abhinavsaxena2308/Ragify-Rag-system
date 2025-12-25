import React from 'react';
import LoadingSpinner from './LoadingSpinner';

const ChatMessage = ({ message, isUser, isLoading = false }) => {
  const { text, sender, timestamp, sources = [] } = message;

  // Format timestamp to a readable format
  const formatTime = (date) => {
    if (!date) return '';
    const time = new Date(date);
    return time.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
  };

  return (
    <div className={`group ${isUser ? 'flex justify-end' : 'flex justify-start'} mb-8`}>
      <div className={`flex items-start space-x-3 max-w-3xl ${isUser ? 'flex-row-reverse space-x-reverse' : ''}`}>
        {/* Avatar */}
        <div className={`flex-shrink-0 w-8 h-8 rounded-full flex items-center justify-center text-sm font-medium ${
          isUser 
            ? 'bg-gradient-to-r from-red-600 to-red-700 text-white' 
            : 'bg-gray-700 text-gray-300'
        }`}>
          {isUser ? 'U' : 'R'}
        </div>

        {/* Message Content */}
        <div className={`flex-1 ${isUser ? 'text-right' : 'text-left'}`}>
          <div className={`inline-block px-4 py-3 rounded-2xl ${
            isUser
              ? 'bg-gradient-to-r from-red-600 to-red-700 text-white rounded-br-sm shadow-lg'
              : 'bg-gray-800 text-gray-100 rounded-bl-sm shadow-lg border border-gray-700'
          }`}>
            {/* Message text */}
            <div className="whitespace-pre-wrap break-words leading-relaxed">
              {isLoading ? (
                <div className="flex items-center space-x-2 py-1">
                  <LoadingSpinner size="sm" className="text-gray-400" />
                  <span className="text-gray-400 text-sm">Thinking...</span>
                </div>
              ) : (
                <div className="prose prose-invert max-w-none">
                  {text}
                </div>
              )}
            </div>

            {/* Sources section for AI responses */}
            {!isUser && !isLoading && sources && sources.length > 0 && (
              <div className="mt-3 pt-3 border-t border-gray-700">
                <p className="text-xs font-medium text-gray-300 mb-2 flex items-center">
                  <svg className="w-4 h-4 mr-1 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                  </svg>
                  Sources:
                </p>
                <div className="space-y-1">
                  {sources.map((source, index) => (
                    <div key={index} className="text-xs text-gray-400 bg-gray-900 rounded px-2 py-1">
                      📄 {source}
                    </div>
                  ))}
                </div>
              </div>
            )}
          </div>

          {/* Timestamp */}
          {timestamp && !isLoading && (
            <div className={`text-xs text-gray-500 mt-1 ${isUser ? 'text-right' : 'text-left'}`}>
              {formatTime(timestamp)}
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default ChatMessage;