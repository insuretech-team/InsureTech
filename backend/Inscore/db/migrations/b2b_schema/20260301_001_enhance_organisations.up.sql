-- =====================================================
-- Production Enhancement: b2b_schema.organisations
-- Tables + columns created by proto-first migration generator.
-- This file: indexes, unique constraints, triggers, RLS ONLY.
-- Do NOT add columns here — add fields to organisation.proto instead.
-- All DDL guarded by table/column existence checks using information_schema.
-- No ::regclass casts (fails if table does not yet exist).
-- =====================================================

BEGIN;

DO $$
BEGIN
  -- Only run if the table exists (proto generator runs before this)
  IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'b2b_schema' AND table_name = 'organisations') THEN
    RAISE NOTICE 'b2b_schema.organisations does not exist yet — skipping enhance. Run proto migration first.';
    RETURN;
  END IF;

  -- ── Unique constraint: code globally unique ───────────────────────────────
  IF NOT EXISTS (
    SELECT 1 FROM pg_constraint c
    JOIN pg_class t ON t.oid = c.conrelid
    JOIN pg_namespace n ON n.oid = t.relnamespace
    WHERE n.nspname = 'b2b_schema' AND t.relname = 'organisations' AND c.conname = 'uq_organisations_code'
  ) THEN
    EXECUTE 'ALTER TABLE b2b_schema.organisations ADD CONSTRAINT uq_organisations_code UNIQUE (code)';
  END IF;

  -- ── Active records index ──────────────────────────────────────────────────
  IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'b2b_schema' AND table_name = 'organisations' AND column_name = 'deleted_at') THEN
    EXECUTE 'CREATE INDEX IF NOT EXISTS idx_organisations_active ON b2b_schema.organisations(organisation_id) WHERE deleted_at IS NULL';
  END IF;

  -- ── tenant_id index ───────────────────────────────────────────────────────
  IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'b2b_schema' AND table_name = 'organisations' AND column_name = 'tenant_id') THEN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'b2b_schema' AND table_name = 'organisations' AND column_name = 'deleted_at') THEN
      EXECUTE 'CREATE INDEX IF NOT EXISTS idx_organisations_tenant_id ON b2b_schema.organisations(tenant_id) WHERE deleted_at IS NULL';
    ELSE
      EXECUTE 'CREATE INDEX IF NOT EXISTS idx_organisations_tenant_id ON b2b_schema.organisations(tenant_id)';
    END IF;
  END IF;

  -- ── status index ──────────────────────────────────────────────────────────
  IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'b2b_schema' AND table_name = 'organisations' AND column_name = 'status') THEN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'b2b_schema' AND table_name = 'organisations' AND column_name = 'deleted_at') THEN
      EXECUTE 'CREATE INDEX IF NOT EXISTS idx_organisations_status ON b2b_schema.organisations(status) WHERE deleted_at IS NULL';
    ELSE
      EXECUTE 'CREATE INDEX IF NOT EXISTS idx_organisations_status ON b2b_schema.organisations(status)';
    END IF;
  END IF;

  -- ── created_at index ──────────────────────────────────────────────────────
  IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'b2b_schema' AND table_name = 'organisations' AND column_name = 'created_at') THEN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'b2b_schema' AND table_name = 'organisations' AND column_name = 'deleted_at') THEN
      EXECUTE 'CREATE INDEX IF NOT EXISTS idx_organisations_created_at ON b2b_schema.organisations(created_at DESC) WHERE deleted_at IS NULL';
    ELSE
      EXECUTE 'CREATE INDEX IF NOT EXISTS idx_organisations_created_at ON b2b_schema.organisations(created_at DESC)';
    END IF;
  END IF;

  -- ── updated_at index ──────────────────────────────────────────────────────
  IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'b2b_schema' AND table_name = 'organisations' AND column_name = 'updated_at') THEN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'b2b_schema' AND table_name = 'organisations' AND column_name = 'deleted_at') THEN
      EXECUTE 'CREATE INDEX IF NOT EXISTS idx_organisations_updated_at ON b2b_schema.organisations(updated_at DESC) WHERE deleted_at IS NULL';
    ELSE
      EXECUTE 'CREATE INDEX IF NOT EXISTS idx_organisations_updated_at ON b2b_schema.organisations(updated_at DESC)';
    END IF;
  END IF;

  -- ── updated_at trigger ────────────────────────────────────────────────────
  IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'b2b_schema' AND table_name = 'organisations' AND column_name = 'updated_at') THEN
    EXECUTE 'CREATE OR REPLACE FUNCTION b2b_schema.trg_organisations_updated_at() RETURNS TRIGGER AS $body$ BEGIN NEW.updated_at = CURRENT_TIMESTAMP; RETURN NEW; END; $body$ LANGUAGE plpgsql';
    EXECUTE 'DROP TRIGGER IF EXISTS trg_organisations_update ON b2b_schema.organisations';
    EXECUTE 'CREATE TRIGGER trg_organisations_update BEFORE UPDATE ON b2b_schema.organisations FOR EACH ROW EXECUTE FUNCTION b2b_schema.trg_organisations_updated_at()';
  END IF;

  -- NOTE: trg_sync_org_employee_count trigger on b2b_schema.employees is installed
  -- in 20260301_004_enhance_employees.up.sql AFTER the employee proto migration
  -- adds deleted_at to the employees table. Installing it here would fail if
  -- employees.deleted_at does not yet exist when this migration runs (migration_order 88
  -- runs before employee proto migration at order 92).

  -- ── RLS ───────────────────────────────────────────────────────────────────
  EXECUTE 'ALTER TABLE b2b_schema.organisations ENABLE ROW LEVEL SECURITY';
  EXECUTE 'DROP POLICY IF EXISTS organisations_isolation ON b2b_schema.organisations';
  EXECUTE $pol$
    CREATE POLICY organisations_isolation ON b2b_schema.organisations
    USING (
      current_setting('app.current_organisation_id', TRUE) IS NULL
      OR current_setting('app.current_organisation_id', TRUE) = ''
      OR organisation_id::TEXT = current_setting('app.current_organisation_id', TRUE)
    )
  $pol$;

  EXECUTE $cmt$COMMENT ON TABLE b2b_schema.organisations IS 'B2B corporate client organisations — enhanced with indexes, triggers, RLS.'$cmt$;

END $$;

COMMIT;
