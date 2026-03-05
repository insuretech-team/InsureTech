-- =====================================================
-- Production Enhancement: authz_schema.roles
-- =====================================================

BEGIN;

-- Add performance indexes if columns exist
DO $$
DECLARE
    pk_col TEXT;
BEGIN
    -- Find the primary key column
    SELECT a.attname INTO pk_col
    FROM pg_index i
    JOIN pg_attribute a ON a.attrelid = i.indrelid AND a.attnum = ANY(i.indkey)
    WHERE i.indrelid = 'authz_schema.roles'::regclass AND i.indisprimary
    LIMIT 1;
    
    -- Index on created_at
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'authz_schema' AND table_name = 'roles' AND column_name = 'created_at') THEN
        IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'authz_schema' AND table_name = 'roles' AND column_name = 'deleted_at') THEN
            EXECUTE 'CREATE INDEX IF NOT EXISTS idx_roles_created_at ON authz_schema.roles(created_at DESC) WHERE deleted_at IS NULL';
        ELSE
            EXECUTE 'CREATE INDEX IF NOT EXISTS idx_roles_created_at ON authz_schema.roles(created_at DESC)';
        END IF;
    END IF;
    
    -- Index on updated_at
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'authz_schema' AND table_name = 'roles' AND column_name = 'updated_at') THEN
        IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'authz_schema' AND table_name = 'roles' AND column_name = 'deleted_at') THEN
            EXECUTE 'CREATE INDEX IF NOT EXISTS idx_roles_updated_at ON authz_schema.roles(updated_at DESC) WHERE deleted_at IS NULL';
        ELSE
            EXECUTE 'CREATE INDEX IF NOT EXISTS idx_roles_updated_at ON authz_schema.roles(updated_at DESC)';
        END IF;
    END IF;
    
    -- Partial index on active records (if deleted_at exists and we have PK)
    IF pk_col IS NOT NULL AND EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'authz_schema' AND table_name = 'roles' AND column_name = 'deleted_at') THEN
        EXECUTE format('CREATE INDEX IF NOT EXISTS idx_roles_active ON authz_schema.roles(%I) WHERE deleted_at IS NULL', pk_col);
    END IF;
END $$;

-- Add updated_at trigger if column exists
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'authz_schema' AND table_name = 'roles' AND column_name = 'updated_at') THEN
        EXECUTE 'CREATE OR REPLACE FUNCTION authz_schema.trg_roles_updated_at() RETURNS TRIGGER AS $body$ BEGIN NEW.updated_at = CURRENT_TIMESTAMP; RETURN NEW; END; $body$ LANGUAGE plpgsql';
        EXECUTE 'DROP TRIGGER IF EXISTS trg_roles_update ON authz_schema.roles';
        EXECUTE 'CREATE TRIGGER trg_roles_update BEFORE UPDATE ON authz_schema.roles FOR EACH ROW EXECUTE FUNCTION authz_schema.trg_roles_updated_at()';
    END IF;
END $$;

COMMENT ON TABLE authz_schema.roles IS 'Enhanced with automated triggers and indexes';

COMMIT;

