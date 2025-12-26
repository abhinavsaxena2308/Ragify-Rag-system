-- Fix source_docs column from TEXT[] to TEXT
-- This will change the column from PostgreSQL array to regular text

-- First, drop the existing column if it exists as array
ALTER TABLE messages DROP COLUMN IF EXISTS source_docs;

-- Add the new column as TEXT
ALTER TABLE messages ADD COLUMN source_docs TEXT;

-- Update the table structure
-- This will allow storing JSON serialized array as text
