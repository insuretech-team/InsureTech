-- =====================================================
-- Production Enhancement: authz_schema.portal_configs
-- Proto: authz/entity/v1/portal_config.proto → PortalConfig
-- migration_order: 5
-- =====================================================

BEGIN;

DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'authz_schema' AND table_name = 'portal_configs') THEN
        -- Index on portal (primary lookup — one row per portal)
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_portal_configs_portal ON authz_schema.portal_configs(portal)';
        -- Index on updated_at for cache invalidation queries
        IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'authz_schema' AND table_name = 'portal_configs' AND column_name = 'updated_at') THEN
            EXECUTE 'CREATE INDEX IF NOT EXISTS idx_portal_configs_updated_at ON authz_schema.portal_configs(updated_at DESC)';
        END IF;
    END IF;
END $$;

-- updated_at trigger
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'authz_schema' AND table_name = 'portal_configs' AND column_name = 'updated_at') THEN
        EXECUTE 'CREATE OR REPLACE FUNCTION authz_schema.trg_portal_configs_updated_at() RETURNS TRIGGER AS $body$ BEGIN NEW.updated_at = CURRENT_TIMESTAMP; RETURN NEW; END; $body$ LANGUAGE plpgsql';
        EXECUTE 'DROP TRIGGER IF EXISTS trg_portal_configs_update ON authz_schema.portal_configs';
        EXECUTE 'CREATE TRIGGER trg_portal_configs_update BEFORE UPDATE ON authz_schema.portal_configs FOR EACH ROW EXECUTE FUNCTION authz_schema.trg_portal_configs_updated_at()';
    END IF;
END $$;

-- Seed default portal configs for all 6 portals (INSERT ... ON CONFLICT DO NOTHING)
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'authz_schema' AND table_name = 'portal_configs') THEN
        -- system: MFA required, short TTLs
        INSERT INTO authz_schema.portal_configs
            (portal, mfa_required, mfa_methods, access_token_ttl_seconds, refresh_token_ttl_seconds, session_ttl_seconds, idle_timeout_seconds, allow_concurrent_sessions, max_concurrent_sessions)
        VALUES
            ('PORTAL_SYSTEM',     true,  ARRAY['totp'],              900,    604800, 28800, 1800, false, 1),
            ('PORTAL_BUSINESS',   true,  ARRAY['email_otp'],         900,    604800, 28800, 1800, true,  3),
            ('PORTAL_B2B',        true,  ARRAY['totp','email_otp'],  900,    604800, 28800, 1800, true,  3),
            ('PORTAL_AGENT',      false, ARRAY['sms_otp'],           1800,   604800, 28800, 3600, true,  5),
            ('PORTAL_REGULATOR',  true,  ARRAY['totp'],              900,    604800, 28800, 1800, false, 1),
            ('PORTAL_B2C',        false, ARRAY['sms_otp'],           3600,   2592000,86400, 7200, true,  10)
        ON CONFLICT (portal) DO NOTHING;
    END IF;
END $$;

COMMENT ON TABLE authz_schema.portal_configs IS 'Enhanced with indexes, updated_at trigger, and seeded default per-portal configs';

COMMIT;
