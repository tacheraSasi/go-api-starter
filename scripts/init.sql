-- This script runs automatically when using Docker Compose

-- Create extensions if needed
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

DO $$
BEGIN
    RAISE NOTICE 'Database initialized successfully for Go API Starter';
END $$;