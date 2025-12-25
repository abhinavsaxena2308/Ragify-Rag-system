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
    <div className={`flex ${isUser ? 'justify-end' : 'justify-start'} mb-6`}>
      <div
        className={`max-w-[85%] md:max-w-[75%] rounded-2xl px-4 py-3 ${
          isUser
            ? 'bg-gradient-to-r from-red-700 to-red-800 text-white rounded-tr-none shadow-lg ml-auto'
            : 'bg-gray-800 text-gray-100 rounded-tl-none shadow-lg border border-gray-700 mr-auto'
        }`}
        style={{ minWidth: '200px' }}
      >
        {/* Message content */}
        <div className="whitespace-pre-wrap break-words">
          {isLoading ? (
            <div className="flex items-center py-2">
              <LoadingSpinner size="sm" className={isUser ? 'text-white' : 'text-gray-100'} />
            </div>
          ) : (
            text
          )}
        </div>

        {/* Timestamp */}
        <div className={`text-xs mt-2 ${isUser ? 'text-red-300' : 'text-gray-400'}`}>
          {timestamp && !isLoading ? formatTime(timestamp) : ''}
        </div>

        {/* Sources section for AI responses */}
        {!isUser && !isLoading && sources && sources.length > 0 && (
          <div className="mt-3 pt-3 border-t border-gray-700">
            <p className="text-xs font-medium text-gray-300 mb-2 flex items-center">
              <span className="mr-2 text-red-400">📄</span>
              Sources:
            </p>
            <ul className="text-xs text-gray-400 space-y-1">
              {sources.map((source, index) => (
                <li key={index} className="truncate">
                  <span className="text-red-400">•</span> {source}
                </li>
              ))}
            </ul>
          </div>
        )}
      </div>
    </div>
  );
};

export default ChatMessage;