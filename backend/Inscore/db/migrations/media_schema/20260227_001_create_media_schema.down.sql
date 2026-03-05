-- =====================================================
-- Rollback: Drop media_schema
-- =====================================================

BEGIN;

-- Drop schema (CASCADE will drop all tables)
DROP SCHEMA IF EXISTS media_schema CASCADE;

COMMIT;
