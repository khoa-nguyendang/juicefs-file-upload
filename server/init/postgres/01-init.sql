-- Initialize PostgreSQL for JuiceFS metadata storage
-- This script runs automatically when PostgreSQL container starts

-- Create schema for JuiceFS if needed
CREATE SCHEMA IF NOT EXISTS public;

-- Grant all privileges to juicefs user
GRANT ALL PRIVILEGES ON DATABASE juicefs TO juicefs;
GRANT ALL ON SCHEMA public TO juicefs;

-- Create any initial tables if needed (JuiceFS will handle its own schema)
-- JuiceFS will automatically create its required tables when formatting

-- Optional: Create a table for tracking file uploads (application-specific)
CREATE TABLE IF NOT EXISTS file_uploads (
    id SERIAL PRIMARY KEY,
    filename VARCHAR(255) NOT NULL,
    file_path TEXT NOT NULL,
    file_size BIGINT,
    content_type VARCHAR(100),
    uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    uploaded_by VARCHAR(100),
    metadata JSONB
);

-- Create index for faster queries
CREATE INDEX IF NOT EXISTS idx_file_uploads_filename ON file_uploads(filename);
CREATE INDEX IF NOT EXISTS idx_file_uploads_uploaded_at ON file_uploads(uploaded_at);

-- Grant permissions on the uploads table
GRANT ALL PRIVILEGES ON TABLE file_uploads TO juicefs;
GRANT USAGE, SELECT ON SEQUENCE file_uploads_id_seq TO juicefs;

-- Add some initial metadata
INSERT INTO file_uploads (filename, file_path, file_size, content_type, uploaded_by, metadata)
VALUES
    ('README.md', '/uploads/README.md', 1024, 'text/markdown', 'system', '{"init": true}')
ON CONFLICT DO NOTHING;