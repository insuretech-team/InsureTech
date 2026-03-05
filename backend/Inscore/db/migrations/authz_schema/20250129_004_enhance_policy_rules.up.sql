-- =====================================================
-- Production Enhancement: authz_schema.policy_rules
-- Proto: authz/entity/v1/policy_rule.proto → PolicyRule
-- migration_order: 4
-- =====================================================

BEGIN;

DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'authz_schema' AND table_name = 'policy_rules') THEN
        -- Composite index: subject + domain (primary lookup for syncing to casbin_rules)
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_policy_rules_subject_domain ON authz_schema.policy_rules(subject, domain)';
        -- Index on domain alone (list all rules for a domain)
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_policy_rules_domain ON authz_schema.policy_rules(domain)';
        -- Partial index: active rules only (most frequent query)
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_policy_rules_active ON authz_schema.policy_rules(domain, subject) WHERE is_active = true AND deleted_at IS NULL';
        -- Index on created_at for audit/ordering
        IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'authz_schema' AND table_name = 'policy_rules' AND column_name = 'created_at') THEN
            EXECUTE 'CREATE INDEX IF NOT EXISTS idx_policy_rules_created_at ON authz_schema.policy_rules(created_at DESC) WHERE deleted_at IS NULL';
        END IF;
        -- Index on updated_at
        IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'authz_schema' AND table_name = 'policy_rules' AND column_name = 'updated_at') THEN
            EXECUTE 'CREATE INDEX IF NOT EXISTS idx_policy_rules_updated_at ON authz_schema.policy_rules(updated_at DESC) WHERE deleted_at IS NULL';
        END IF;
    END IF;
END $$;

-- updated_at trigger
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'authz_schema' AND table_name = 'policy_rules' AND column_name = 'updated_at') THEN
        EXECUTE 'CREATE OR REPLACE FUNCTION authz_schema.trg_policy_rules_updated_at() RETURNS TRIGGER AS $body$ BEGIN NEW.updated_at = CURRENT_TIMESTAMP; RETURN NEW; END; $body$ LANGUAGE plpgsql';
        EXECUTE 'DROP TRIGGER IF EXISTS trg_policy_rules_update ON authz_schema.policy_rules';
        EXECUTE 'CREATE TRIGGER trg_policy_rules_update BEFORE UPDATE ON authz_schema.policy_rules FOR EACH ROW EXECUTE FUNCTION authz_schema.trg_policy_rules_updated_at()';
    END IF;
END $$;

COMMENT ON TABLE authz_schema.policy_rules IS 'Enhanced with composite indexes for Casbin policy sync and active-rule filtering';

COMMIT;
