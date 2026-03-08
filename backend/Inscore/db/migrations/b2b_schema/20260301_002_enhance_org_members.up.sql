-- =====================================================
-- Production Enhancement: b2b_schema.org_members
-- Tables + columns created by proto-first migration generator.
-- This file: indexes, unique constraints, triggers, RLS ONLY.
-- Do NOT add columns here — add fields to organisation.proto instead.
-- All DDL guarded by table/column existence checks. No ::regclass casts.
-- =====================================================

BEGIN;

DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'b2b_schema' AND table_name = 'org_members') THEN
    RAISE NOTICE 'b2b_schema.org_members does not exist yet — skipping enhance. Run proto migration first.';
    RETURN;
  END IF;

  -- ── Unique: one active membership per user per org ────────────────────────
  -- This is the constraint that prevents duplicate memberships.
  IF NOT EXISTS (
    SELECT 1 FROM pg_indexes
    WHERE schemaname = 'b2b_schema' AND tablename = 'org_members' AND indexname = 'uq_org_members_user_org_active'
  ) THEN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'b2b_schema' AND table_name = 'org_members' AND column_name = 'deleted_at') THEN
      EXECUTE 'CREATE UNIQUE INDEX uq_org_members_user_org_active ON b2b_schema.org_members(organisation_id, user_id) WHERE deleted_at IS NULL';
    ELSE
      EXECUTE 'CREATE UNIQUE INDEX uq_org_members_user_org_active ON b2b_schema.org_members(organisation_id, user_id)';
    END IF;
  END IF;

  -- ── Active records index ──────────────────────────────────────────────────
  IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'b2b_schema' AND table_name = 'org_members' AND column_name = 'deleted_at') THEN
    EXECUTE 'CREATE INDEX IF NOT EXISTS idx_org_members_active ON b2b_schema.org_members(member_id) WHERE deleted_at IS NULL';
  END IF;

  -- ── CRITICAL: user_id lookup (used on every authenticated request) ────────
  IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'b2b_schema' AND table_name = 'org_members' AND column_name = 'user_id') THEN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'b2b_schema' AND table_name = 'org_members' AND column_name = 'deleted_at') THEN
      EXECUTE 'CREATE INDEX IF NOT EXISTS idx_org_members_user_id ON b2b_schema.org_members(user_id) WHERE deleted_at IS NULL AND status = ''ORG_MEMBER_STATUS_ACTIVE''';
    ELSE
      EXECUTE 'CREATE INDEX IF NOT EXISTS idx_org_members_user_id ON b2b_schema.org_members(user_id)';
    END IF;
  END IF;

  -- ── organisation_id index ─────────────────────────────────────────────────
  IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'b2b_schema' AND table_name = 'org_members' AND column_name = 'organisation_id') THEN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'b2b_schema' AND table_name = 'org_members' AND column_name = 'deleted_at') THEN
      EXECUTE 'CREATE INDEX IF NOT EXISTS idx_org_members_organisation_id ON b2b_schema.org_members(organisation_id) WHERE deleted_at IS NULL';
    ELSE
      EXECUTE 'CREATE INDEX IF NOT EXISTS idx_org_members_organisation_id ON b2b_schema.org_members(organisation_id)';
    END IF;
  END IF;

  -- ── role index (find all admins of an org) ────────────────────────────────
  IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'b2b_schema' AND table_name = 'org_members' AND column_name = 'role') THEN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'b2b_schema' AND table_name = 'org_members' AND column_name = 'deleted_at') THEN
      EXECUTE 'CREATE INDEX IF NOT EXISTS idx_org_members_org_role ON b2b_schema.org_members(organisation_id, role) WHERE deleted_at IS NULL AND status = ''ORG_MEMBER_STATUS_ACTIVE''';
    ELSE
      EXECUTE 'CREATE INDEX IF NOT EXISTS idx_org_members_org_role ON b2b_schema.org_members(organisation_id, role)';
    END IF;
  END IF;

  -- ── joined_at index (ResolveMyOrganisation ORDER BY joined_at ASC) ────────
  IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'b2b_schema' AND table_name = 'org_members' AND column_name = 'joined_at') THEN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'b2b_schema' AND table_name = 'org_members' AND column_name = 'deleted_at') THEN
      EXECUTE 'CREATE INDEX IF NOT EXISTS idx_org_members_joined_at ON b2b_schema.org_members(joined_at ASC) WHERE deleted_at IS NULL';
    ELSE
      EXECUTE 'CREATE INDEX IF NOT EXISTS idx_org_members_joined_at ON b2b_schema.org_members(joined_at ASC)';
    END IF;
  END IF;

  -- ── created_at / updated_at indexes ──────────────────────────────────────
  IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'b2b_schema' AND table_name = 'org_members' AND column_name = 'created_at') THEN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'b2b_schema' AND table_name = 'org_members' AND column_name = 'deleted_at') THEN
      EXECUTE 'CREATE INDEX IF NOT EXISTS idx_org_members_created_at ON b2b_schema.org_members(created_at DESC) WHERE deleted_at IS NULL';
    ELSE
      EXECUTE 'CREATE INDEX IF NOT EXISTS idx_org_members_created_at ON b2b_schema.org_members(created_at DESC)';
    END IF;
  END IF;

  IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'b2b_schema' AND table_name = 'org_members' AND column_name = 'updated_at') THEN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'b2b_schema' AND table_name = 'org_members' AND column_name = 'deleted_at') THEN
      EXECUTE 'CREATE INDEX IF NOT EXISTS idx_org_members_updated_at ON b2b_schema.org_members(updated_at DESC) WHERE deleted_at IS NULL';
    ELSE
      EXECUTE 'CREATE INDEX IF NOT EXISTS idx_org_members_updated_at ON b2b_schema.org_members(updated_at DESC)';
    END IF;
  END IF;

  -- ── updated_at trigger ────────────────────────────────────────────────────
  IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'b2b_schema' AND table_name = 'org_members' AND column_name = 'updated_at') THEN
    EXECUTE 'CREATE OR REPLACE FUNCTION b2b_schema.trg_org_members_updated_at() RETURNS TRIGGER AS $body$ BEGIN NEW.updated_at = CURRENT_TIMESTAMP; RETURN NEW; END; $body$ LANGUAGE plpgsql';
    EXECUTE 'DROP TRIGGER IF EXISTS trg_org_members_update ON b2b_schema.org_members';
    EXECUTE 'CREATE TRIGGER trg_org_members_update BEFORE UPDATE ON b2b_schema.org_members FOR EACH ROW EXECUTE FUNCTION b2b_schema.trg_org_members_updated_at()';
  END IF;

  -- ── RLS ───────────────────────────────────────────────────────────────────
  EXECUTE 'ALTER TABLE b2b_schema.org_members ENABLE ROW LEVEL SECURITY';
  EXECUTE 'DROP POLICY IF EXISTS org_members_isolation ON b2b_schema.org_members';
  EXECUTE $pol$
    CREATE POLICY org_members_isolation ON b2b_schema.org_members
    USING (
      current_setting('app.current_organisation_id', TRUE) IS NULL
      OR current_setting('app.current_organisation_id', TRUE) = ''
      OR organisation_id::TEXT = current_setting('app.current_organisation_id', TRUE)
    )
  $pol$;

  EXECUTE $cmt$COMMENT ON TABLE b2b_schema.org_members IS 'Organisation membership — enhanced with critical user_id index for session-based business_id resolution, unique active membership constraint, RLS.'$cmt$;

END $$;

COMMIT;
