import React from 'react';

const EmptyState = ({ 
  title, 
  description, 
  icon = '📄', 
  actionButton = null,
  showIcon = true
}) => {
  return (
    <div className="text-center py-12">
      {showIcon && (
        <div className="mx-auto flex items-center justify-center h-16 w-16 rounded-full bg-gray-100">
          <span className="text-3xl">{icon}</span>
        </div>
      )}
      <h3 className="mt-4 text-lg font-medium text-gray-900">{title}</h3>
      <p className="mt-2 text-sm text-gray-500 max-w-md mx-auto">
        {description}
      </p>
      {actionButton && (
        <div className="mt-6">
          {actionButton}
        </div>
      )}
    </div>
  );
};

export default EmptyState;