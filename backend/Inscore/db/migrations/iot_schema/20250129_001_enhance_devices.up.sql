-- =====================================================
-- Production Enhancement: iot_schema.devices
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
    WHERE i.indrelid = 'iot_schema.devices'::regclass AND i.indisprimary
    LIMIT 1;
    
    -- Index on created_at
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'iot_schema' AND table_name = 'devices' AND column_name = 'created_at') THEN
        IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'iot_schema' AND table_name = 'devices' AND column_name = 'deleted_at') THEN
            EXECUTE 'CREATE INDEX IF NOT EXISTS idx_devices_created_at ON iot_schema.devices(created_at DESC) WHERE deleted_at IS NULL';
        ELSE
            EXECUTE 'CREATE INDEX IF NOT EXISTS idx_devices_created_at ON iot_schema.devices(created_at DESC)';
        END IF;
    END IF;
    
    -- Index on updated_at
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'iot_schema' AND table_name = 'devices' AND column_name = 'updated_at') THEN
        IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'iot_schema' AND table_name = 'devices' AND column_name = 'deleted_at') THEN
            EXECUTE 'CREATE INDEX IF NOT EXISTS idx_devices_updated_at ON iot_schema.devices(updated_at DESC) WHERE deleted_at IS NULL';
        ELSE
            EXECUTE 'CREATE INDEX IF NOT EXISTS idx_devices_updated_at ON iot_schema.devices(updated_at DESC)';
        END IF;
    END IF;
    
    -- Partial index on active records (if deleted_at exists and we have PK)
    IF pk_col IS NOT NULL AND EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'iot_schema' AND table_name = 'devices' AND column_name = 'deleted_at') THEN
        EXECUTE format('CREATE INDEX IF NOT EXISTS idx_devices_active ON iot_schema.devices(%I) WHERE deleted_at IS NULL', pk_col);
    END IF;
END $$;

-- Add updated_at trigger if column exists
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'iot_schema' AND table_name = 'devices' AND column_name = 'updated_at') THEN
        EXECUTE 'CREATE OR REPLACE FUNCTION iot_schema.trg_devices_updated_at() RETURNS TRIGGER AS $body$ BEGIN NEW.updated_at = CURRENT_TIMESTAMP; RETURN NEW; END; $body$ LANGUAGE plpgsql';
        EXECUTE 'DROP TRIGGER IF EXISTS trg_devices_update ON iot_schema.devices';
        EXECUTE 'CREATE TRIGGER trg_devices_update BEFORE UPDATE ON iot_schema.devices FOR EACH ROW EXECUTE FUNCTION iot_schema.trg_devices_updated_at()';
    END IF;
END $$;

COMMENT ON TABLE iot_schema.devices IS 'Enhanced with automated triggers and indexes';

COMMIT;
