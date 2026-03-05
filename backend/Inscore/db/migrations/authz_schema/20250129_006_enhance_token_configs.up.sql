-- =====================================================
-- Production Enhancement: authz_schema.token_configs
-- Proto: authz/entity/v1/token_config.proto → TokenConfig
-- migration_order: 6
-- =====================================================

BEGIN;

DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'authz_schema' AND table_name = 'token_configs') THEN
        -- Partial index: active key only (JWKS endpoint serves only active keys)
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_token_configs_active ON authz_schema.token_configs(kid) WHERE is_active = true';
        -- Index on created_at for key rotation history ordering
        IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'authz_schema' AND table_name = 'token_configs' AND column_name = 'created_at') THEN
            EXECUTE 'CREATE INDEX IF NOT EXISTS idx_token_configs_created_at ON authz_schema.token_configs(created_at DESC)';
        END IF;
        -- Index on rotated_at for rotation audit
        IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'authz_schema' AND table_name = 'token_configs' AND column_name = 'rotated_at') THEN
            EXECUTE 'CREATE INDEX IF NOT EXISTS idx_token_configs_rotated_at ON authz_schema.token_configs(rotated_at DESC) WHERE rotated_at IS NOT NULL';
        END IF;
    END IF;
END $$;

COMMENT ON TABLE authz_schema.token_configs IS 'Enhanced with indexes for active key lookup and rotation history';

COMMIT;
