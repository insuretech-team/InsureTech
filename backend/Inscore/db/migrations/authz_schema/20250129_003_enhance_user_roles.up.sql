-- =====================================================
-- Production Enhancement: authz_schema.user_roles
-- Proto: authz/entity/v1/user_role.proto → UserRole
-- migration_order: 3
-- =====================================================

BEGIN;

DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'authz_schema' AND table_name = 'user_roles') THEN
        -- Composite index: user_id + domain (primary lookup pattern for role resolution)
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_user_roles_user_domain ON authz_schema.user_roles(user_id, domain)';
        -- Index on domain alone (list all users in domain)
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_user_roles_domain ON authz_schema.user_roles(domain)';
        -- Index on role_id (list all users with a given role)
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_user_roles_role_id ON authz_schema.user_roles(role_id)';
        -- Index on expires_at for background cleanup job (partial index on active rows omitted: CURRENT_TIMESTAMP is not IMMUTABLE)
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_user_roles_expires_at ON authz_schema.user_roles(expires_at) WHERE expires_at IS NOT NULL';
    END IF;
END $$;

-- Unique constraint: one role per user per domain (prevents duplicate assignments)
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'authz_schema' AND table_name = 'user_roles') THEN
        IF NOT EXISTS (
            SELECT 1 FROM information_schema.table_constraints
            WHERE table_schema = 'authz_schema' AND table_name = 'user_roles'
              AND constraint_name = 'uq_user_roles_user_role_domain'
        ) THEN
            EXECUTE 'ALTER TABLE authz_schema.user_roles ADD CONSTRAINT uq_user_roles_user_role_domain UNIQUE (user_id, role_id, domain)';
        END IF;
    END IF;
END $$;

COMMENT ON TABLE authz_schema.user_roles IS 'Enhanced with composite indexes for fast role resolution and expiry management';

COMMIT;
