-- =====================================================
-- Production Enhancement: authn_schema.otps
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
    WHERE i.indrelid = 'authn_schema.otps'::regclass AND i.indisprimary
    LIMIT 1;
    
    -- Index on created_at
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'authn_schema' AND table_name = 'otps' AND column_name = 'created_at') THEN
        IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'authn_schema' AND table_name = 'otps' AND column_name = 'deleted_at') THEN
            EXECUTE 'CREATE INDEX IF NOT EXISTS idx_otps_created_at ON authn_schema.otps(created_at DESC) WHERE deleted_at IS NULL';
        ELSE
            EXECUTE 'CREATE INDEX IF NOT EXISTS idx_otps_created_at ON authn_schema.otps(created_at DESC)';
        END IF;
    END IF;
    
    -- Index on updated_at
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'authn_schema' AND table_name = 'otps' AND column_name = 'updated_at') THEN
        IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'authn_schema' AND table_name = 'otps' AND column_name = 'deleted_at') THEN
            EXECUTE 'CREATE INDEX IF NOT EXISTS idx_otps_updated_at ON authn_schema.otps(updated_at DESC) WHERE deleted_at IS NULL';
        ELSE
            EXECUTE 'CREATE INDEX IF NOT EXISTS idx_otps_updated_at ON authn_schema.otps(updated_at DESC)';
        END IF;
    END IF;
    
    -- Partial index on active records (if deleted_at exists and we have PK)
    IF pk_col IS NOT NULL AND EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'authn_schema' AND table_name = 'otps' AND column_name = 'deleted_at') THEN
        EXECUTE format('CREATE INDEX IF NOT EXISTS idx_otps_active ON authn_schema.otps(%I) WHERE deleted_at IS NULL', pk_col);
    END IF;
END $$;

-- Add updated_at trigger if column exists
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'authn_schema' AND table_name = 'otps' AND column_name = 'updated_at') THEN
        EXECUTE 'CREATE OR REPLACE FUNCTION authn_schema.trg_otps_updated_at() RETURNS TRIGGER AS $body$ BEGIN NEW.updated_at = CURRENT_TIMESTAMP; RETURN NEW; END; $body$ LANGUAGE plpgsql';
        EXECUTE 'DROP TRIGGER IF EXISTS trg_otps_update ON authn_schema.otps';
        EXECUTE 'CREATE TRIGGER trg_otps_update BEFORE UPDATE ON authn_schema.otps FOR EACH ROW EXECUTE FUNCTION authn_schema.trg_otps_updated_at()';
    END IF;
END $$;

COMMENT ON TABLE authn_schema.otps IS 'Enhanced with automated triggers and indexes';

COMMIT;
