-- =====================================================
-- Production Enhancement: insurance_schema.insurer_configs
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
    WHERE i.indrelid = 'insurance_schema.insurer_configs'::regclass AND i.indisprimary
    LIMIT 1;
    
    -- Index on created_at
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'insurance_schema' AND table_name = 'insurer_configs' AND column_name = 'created_at') THEN
        IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'insurance_schema' AND table_name = 'insurer_configs' AND column_name = 'deleted_at') THEN
            EXECUTE 'CREATE INDEX IF NOT EXISTS idx_insurer_configs_created_at ON insurance_schema.insurer_configs(created_at DESC) WHERE deleted_at IS NULL';
        ELSE
            EXECUTE 'CREATE INDEX IF NOT EXISTS idx_insurer_configs_created_at ON insurance_schema.insurer_configs(created_at DESC)';
        END IF;
    END IF;
    
    -- Index on updated_at
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'insurance_schema' AND table_name = 'insurer_configs' AND column_name = 'updated_at') THEN
        IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'insurance_schema' AND table_name = 'insurer_configs' AND column_name = 'deleted_at') THEN
            EXECUTE 'CREATE INDEX IF NOT EXISTS idx_insurer_configs_updated_at ON insurance_schema.insurer_configs(updated_at DESC) WHERE deleted_at IS NULL';
        ELSE
            EXECUTE 'CREATE INDEX IF NOT EXISTS idx_insurer_configs_updated_at ON insurance_schema.insurer_configs(updated_at DESC)';
        END IF;
    END IF;
    
    -- Partial index on active records (if deleted_at exists and we have PK)
    IF pk_col IS NOT NULL AND EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'insurance_schema' AND table_name = 'insurer_configs' AND column_name = 'deleted_at') THEN
        EXECUTE format('CREATE INDEX IF NOT EXISTS idx_insurer_configs_active ON insurance_schema.insurer_configs(%I) WHERE deleted_at IS NULL', pk_col);
    END IF;
END $$;

-- Add updated_at trigger if column exists
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'insurance_schema' AND table_name = 'insurer_configs' AND column_name = 'updated_at') THEN
        EXECUTE 'CREATE OR REPLACE FUNCTION insurance_schema.trg_insurer_configs_updated_at() RETURNS TRIGGER AS $body$ BEGIN NEW.updated_at = CURRENT_TIMESTAMP; RETURN NEW; END; $body$ LANGUAGE plpgsql';
        EXECUTE 'DROP TRIGGER IF EXISTS trg_insurer_configs_update ON insurance_schema.insurer_configs';
        EXECUTE 'CREATE TRIGGER trg_insurer_configs_update BEFORE UPDATE ON insurance_schema.insurer_configs FOR EACH ROW EXECUTE FUNCTION insurance_schema.trg_insurer_configs_updated_at()';
    END IF;
END $$;

COMMENT ON TABLE insurance_schema.insurer_configs IS 'Enhanced with automated triggers and indexes';

COMMIT;
