-- =====================================================
-- Production Enhancement: authz_schema.access_decision_audits
-- Proto: authz/entity/v1/access_decision_audit.proto → AccessDecisionAudit
-- migration_order: 7
-- =====================================================

BEGIN;

DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'authz_schema' AND table_name = 'access_decision_audits') THEN
        -- Composite index: user_id + decided_at (primary audit query: all decisions for a user)
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_ada_user_decided ON authz_schema.access_decision_audits(user_id, decided_at DESC)';
        -- Composite index: domain + decided_at (tenant-scoped audit queries)
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_ada_domain_decided ON authz_schema.access_decision_audits(domain, decided_at DESC)';
        -- Partial index: DENY decisions only (SIEM / anomaly detection queries)
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_ada_deny_decisions ON authz_schema.access_decision_audits(user_id, decided_at DESC) WHERE decision = ''POLICY_EFFECT_DENY''';
        -- Index on decided_at for time-range queries and retention cleanup
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_ada_decided_at ON authz_schema.access_decision_audits(decided_at DESC)';
    END IF;
END $$;

-- Table-level constraint: audit rows are immutable (prevent UPDATE)
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'authz_schema' AND table_name = 'access_decision_audits') THEN
        -- Immutability enforced via trigger
        EXECUTE 'CREATE OR REPLACE FUNCTION authz_schema.trg_access_decision_audits_immutable() RETURNS TRIGGER AS $body$ BEGIN RAISE EXCEPTION ''access_decision_audits rows are immutable''; END; $body$ LANGUAGE plpgsql';
        EXECUTE 'DROP TRIGGER IF EXISTS trg_ada_immutable ON authz_schema.access_decision_audits';
        EXECUTE 'CREATE TRIGGER trg_ada_immutable BEFORE UPDATE ON authz_schema.access_decision_audits FOR EACH ROW EXECUTE FUNCTION authz_schema.trg_access_decision_audits_immutable()';
    END IF;
END $$;

COMMENT ON TABLE authz_schema.access_decision_audits IS 'Enhanced with audit indexes, DENY partial index for SIEM, and immutability trigger';

COMMIT;
