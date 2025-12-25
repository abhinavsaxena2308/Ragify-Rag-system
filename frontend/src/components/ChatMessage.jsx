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
    <div className={`flex ${isUser ? 'justify-end' : 'justify-start'} mb-4`}>
      <div
        className={`max-w-[85%] md:max-w-[75%] rounded-2xl px-4 py-3 ${
          isUser
            ? 'bg-red-600 text-white rounded-tr-none shadow-md'
            : 'bg-white text-gray-800 rounded-tl-none shadow-md border border-gray-200'
        }`}
      >
        {/* Message content */}
        <div className="whitespace-pre-wrap break-words">
          {isLoading ? (
            <div className="flex items-center">
              <LoadingSpinner size="sm" className={isUser ? 'text-white' : 'text-gray-800'} />
            </div>
          ) : (
            text
          )}
        </div>

        {/* Timestamp */}
        <div className={`text-xs mt-1 ${isUser ? 'text-red-200' : 'text-gray-500'}`}>
          {timestamp && !isLoading ? formatTime(timestamp) : ''}
        </div>

        {/* Sources section for AI responses */}
        {!isUser && !isLoading && sources && sources.length > 0 && (
          <div className="mt-2 pt-2 border-t border-gray-200">
            <p className="text-xs font-medium text-gray-600 mb-1 flex items-center">
              <span className="mr-1">📄</span>
              Sources:
            </p>
            <ul className="text-xs text-gray-500 space-y-1">
              {sources.map((source, index) => (
                <li key={index} className="truncate">
                  • {source}
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