-- =====================================================
-- Production Enhancement: authz_schema.role_mfa_configs
-- Proto: authz/entity/v1/role_mfa_config.proto → RoleMFAConfig
-- migration_order: 8
-- =====================================================

BEGIN;

DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'authz_schema' AND table_name = 'role_mfa_configs') THEN
        -- Index on mfa_required for filtering MFA-required roles (login path check)
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_role_mfa_configs_required ON authz_schema.role_mfa_configs(role_id) WHERE mfa_required = true';
        -- Index on updated_at for cache invalidation
        IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'authz_schema' AND table_name = 'role_mfa_configs' AND column_name = 'updated_at') THEN
            EXECUTE 'CREATE INDEX IF NOT EXISTS idx_role_mfa_configs_updated_at ON authz_schema.role_mfa_configs(updated_at DESC)';
        END IF;
    END IF;
END $$;

-- updated_at trigger
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'authz_schema' AND table_name = 'role_mfa_configs' AND column_name = 'updated_at') THEN
        EXECUTE 'CREATE OR REPLACE FUNCTION authz_schema.trg_role_mfa_configs_updated_at() RETURNS TRIGGER AS $body$ BEGIN NEW.updated_at = CURRENT_TIMESTAMP; RETURN NEW; END; $body$ LANGUAGE plpgsql';
        EXECUTE 'DROP TRIGGER IF EXISTS trg_role_mfa_configs_update ON authz_schema.role_mfa_configs';
        EXECUTE 'CREATE TRIGGER trg_role_mfa_configs_update BEFORE UPDATE ON authz_schema.role_mfa_configs FOR EACH ROW EXECUTE FUNCTION authz_schema.trg_role_mfa_configs_updated_at()';
    END IF;
END $$;

COMMENT ON TABLE authz_schema.role_mfa_configs IS 'Enhanced with indexes and updated_at trigger for MFA config lookups on login';

COMMIT;
