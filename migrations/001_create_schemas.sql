-- 001_create_schemas.sql
-- Create schemas for each service

CREATE SCHEMA IF NOT EXISTS identity;
CREATE SCHEMA IF NOT EXISTS catalog;
CREATE SCHEMA IF NOT EXISTS booking;
CREATE SCHEMA IF NOT EXISTS notification;

-- Grant permissions (adjust as needed for production)
GRANT ALL ON SCHEMA identity TO postgres;
GRANT ALL ON SCHEMA catalog TO postgres;
GRANT ALL ON SCHEMA booking TO postgres;
GRANT ALL ON SCHEMA notification TO postgres;
