-- =====================================================
-- Production Enhancement: authz_schema.casbin_rules
-- Proto: authz/entity/v1/casbin_rule.proto → CasbinRule
-- migration_order: 1
-- =====================================================

BEGIN;

-- Composite index for gorm-adapter policy lookups: ptype + v0 + v1
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'authz_schema' AND table_name = 'casbin_rules') THEN
        -- Composite lookup index: ptype + domain (v1) — most common query pattern
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_casbin_rules_ptype_v1 ON authz_schema.casbin_rules(ptype, v1)';
        -- Composite index for subject + domain lookups (role resolution)
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_casbin_rules_v0_v1 ON authz_schema.casbin_rules(v0, v1)';
        -- Index for full tuple lookup (enforce call): ptype, v0, v1, v2, v3
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_casbin_rules_full_tuple ON authz_schema.casbin_rules(ptype, v0, v1, v2, v3)';
    END IF;
END $$;

-- Unique constraint: prevent duplicate casbin rules
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'authz_schema' AND table_name = 'casbin_rules') THEN
        IF NOT EXISTS (
            SELECT 1 FROM information_schema.table_constraints
            WHERE table_schema = 'authz_schema' AND table_name = 'casbin_rules'
              AND constraint_name = 'uq_casbin_rules_tuple'
        ) THEN
            EXECUTE 'ALTER TABLE authz_schema.casbin_rules ADD CONSTRAINT uq_casbin_rules_tuple UNIQUE (ptype, v0, v1, v2, v3, v4, v5)';
        END IF;
    END IF;
END $$;

COMMENT ON TABLE authz_schema.casbin_rules IS 'Enhanced with composite indexes for gorm-adapter Casbin PERM model enforcement';

COMMIT;
