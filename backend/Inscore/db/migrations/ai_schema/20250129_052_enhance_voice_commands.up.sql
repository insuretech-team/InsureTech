-- =====================================================
-- Production Enhancement: ai_schema.voice_commands
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
    WHERE i.indrelid = 'ai_schema.voice_commands'::regclass AND i.indisprimary
    LIMIT 1;
    
    -- Index on created_at
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'ai_schema' AND table_name = 'voice_commands' AND column_name = 'created_at') THEN
        IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'ai_schema' AND table_name = 'voice_commands' AND column_name = 'deleted_at') THEN
            EXECUTE 'CREATE INDEX IF NOT EXISTS idx_voice_commands_created_at ON ai_schema.voice_commands(created_at DESC) WHERE deleted_at IS NULL';
        ELSE
            EXECUTE 'CREATE INDEX IF NOT EXISTS idx_voice_commands_created_at ON ai_schema.voice_commands(created_at DESC)';
        END IF;
    END IF;
    
    -- Index on updated_at
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'ai_schema' AND table_name = 'voice_commands' AND column_name = 'updated_at') THEN
        IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'ai_schema' AND table_name = 'voice_commands' AND column_name = 'deleted_at') THEN
            EXECUTE 'CREATE INDEX IF NOT EXISTS idx_voice_commands_updated_at ON ai_schema.voice_commands(updated_at DESC) WHERE deleted_at IS NULL';
        ELSE
            EXECUTE 'CREATE INDEX IF NOT EXISTS idx_voice_commands_updated_at ON ai_schema.voice_commands(updated_at DESC)';
        END IF;
    END IF;
    
    -- Partial index on active records (if deleted_at exists and we have PK)
    IF pk_col IS NOT NULL AND EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'ai_schema' AND table_name = 'voice_commands' AND column_name = 'deleted_at') THEN
        EXECUTE format('CREATE INDEX IF NOT EXISTS idx_voice_commands_active ON ai_schema.voice_commands(%I) WHERE deleted_at IS NULL', pk_col);
    END IF;
END $$;

-- Add updated_at trigger if column exists
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'ai_schema' AND table_name = 'voice_commands' AND column_name = 'updated_at') THEN
        EXECUTE 'CREATE OR REPLACE FUNCTION ai_schema.trg_voice_commands_updated_at() RETURNS TRIGGER AS $body$ BEGIN NEW.updated_at = CURRENT_TIMESTAMP; RETURN NEW; END; $body$ LANGUAGE plpgsql';
        EXECUTE 'DROP TRIGGER IF EXISTS trg_voice_commands_update ON ai_schema.voice_commands';
        EXECUTE 'CREATE TRIGGER trg_voice_commands_update BEFORE UPDATE ON ai_schema.voice_commands FOR EACH ROW EXECUTE FUNCTION ai_schema.trg_voice_commands_updated_at()';
    END IF;
END $$;

COMMENT ON TABLE ai_schema.voice_commands IS 'Enhanced with automated triggers and indexes';

COMMIT;
