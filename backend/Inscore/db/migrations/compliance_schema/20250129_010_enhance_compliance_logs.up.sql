-- =====================================================
-- Production Enhancement: compliance_schema.compliance_logs
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
    WHERE i.indrelid = 'compliance_schema.compliance_logs'::regclass AND i.indisprimary
    LIMIT 1;
    
    -- Index on created_at
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'compliance_schema' AND table_name = 'compliance_logs' AND column_name = 'created_at') THEN
        IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'compliance_schema' AND table_name = 'compliance_logs' AND column_name = 'deleted_at') THEN
            EXECUTE 'CREATE INDEX IF NOT EXISTS idx_compliance_logs_created_at ON compliance_schema.compliance_logs(created_at DESC) WHERE deleted_at IS NULL';
        ELSE
            EXECUTE 'CREATE INDEX IF NOT EXISTS idx_compliance_logs_created_at ON compliance_schema.compliance_logs(created_at DESC)';
        END IF;
    END IF;
    
    -- Index on updated_at
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'compliance_schema' AND table_name = 'compliance_logs' AND column_name = 'updated_at') THEN
        IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'compliance_schema' AND table_name = 'compliance_logs' AND column_name = 'deleted_at') THEN
            EXECUTE 'CREATE INDEX IF NOT EXISTS idx_compliance_logs_updated_at ON compliance_schema.compliance_logs(updated_at DESC) WHERE deleted_at IS NULL';
        ELSE
            EXECUTE 'CREATE INDEX IF NOT EXISTS idx_compliance_logs_updated_at ON compliance_schema.compliance_logs(updated_at DESC)';
        END IF;
    END IF;
    
    -- Partial index on active records (if deleted_at exists and we have PK)
    IF pk_col IS NOT NULL AND EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'compliance_schema' AND table_name = 'compliance_logs' AND column_name = 'deleted_at') THEN
        EXECUTE format('CREATE INDEX IF NOT EXISTS idx_compliance_logs_active ON compliance_schema.compliance_logs(%I) WHERE deleted_at IS NULL', pk_col);
    END IF;
END $$;

-- Add updated_at trigger if column exists
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'compliance_schema' AND table_name = 'compliance_logs' AND column_name = 'updated_at') THEN
        EXECUTE 'CREATE OR REPLACE FUNCTION compliance_schema.trg_compliance_logs_updated_at() RETURNS TRIGGER AS $body$ BEGIN NEW.updated_at = CURRENT_TIMESTAMP; RETURN NEW; END; $body$ LANGUAGE plpgsql';
        EXECUTE 'DROP TRIGGER IF EXISTS trg_compliance_logs_update ON compliance_schema.compliance_logs';
        EXECUTE 'CREATE TRIGGER trg_compliance_logs_update BEFORE UPDATE ON compliance_schema.compliance_logs FOR EACH ROW EXECUTE FUNCTION compliance_schema.trg_compliance_logs_updated_at()';
    END IF;
END $$;

COMMENT ON TABLE compliance_schema.compliance_logs IS 'Enhanced with automated triggers and indexes';

COMMIT;
