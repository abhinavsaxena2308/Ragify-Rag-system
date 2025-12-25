import React from 'react';

const EmptyState = ({ 
  title, 
  description, 
  icon = '📄', 
  actionButton = null,
  showIcon = true
}) => {
  return (
    <div className="text-center py-12 px-4">
      {showIcon && (
        <div className="mx-auto flex items-center justify-center h-16 w-16 rounded-full bg-gradient-to-br from-[#DC143C] to-red-700 text-white shadow-lg">
          <span className="text-3xl">{icon}</span>
        </div>
      )}
      <h3 className="mt-6 text-xl font-semibold text-white">{title}</h3>
      <p className="mt-3 text-gray-300 max-w-md mx-auto">
        {description}
      </p>
      {actionButton && (
        <div className="mt-8">
          {actionButton}
        </div>
      )}
    </div>
  );
};

export default EmptyState;