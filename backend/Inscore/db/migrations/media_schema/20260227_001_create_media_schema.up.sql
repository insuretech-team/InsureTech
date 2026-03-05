-- =====================================================
-- Create media_schema for media file management
-- Tables are auto-generated from proto definitions
-- This migration creates the schema and any additional objects
-- =====================================================

BEGIN;

-- Create media_schema if it doesn't exist
CREATE SCHEMA IF NOT EXISTS media_schema;

-- Grant permissions (only if role exists)
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'inscore_app') THEN
        GRANT USAGE ON SCHEMA media_schema TO inscore_app;
        GRANT ALL ON ALL TABLES IN SCHEMA media_schema TO inscore_app;
        GRANT ALL ON ALL SEQUENCES IN SCHEMA media_schema TO inscore_app;
        ALTER DEFAULT PRIVILEGES IN SCHEMA media_schema GRANT ALL ON TABLES TO inscore_app;
        ALTER DEFAULT PRIVILEGES IN SCHEMA media_schema GRANT ALL ON SEQUENCES TO inscore_app;
    END IF;
END $$;

COMMENT ON SCHEMA media_schema IS 'Media file management with processing, validation, and virus scanning';

COMMIT;
