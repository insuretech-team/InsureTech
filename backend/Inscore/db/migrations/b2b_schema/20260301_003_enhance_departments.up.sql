-- =====================================================
-- Production Enhancement: b2b_schema.departments
-- Tables + columns created by proto-first migration generator.
-- This file: FK constraints, indexes, triggers, RLS ONLY.
-- Do NOT add columns here — add fields to department.proto instead.
-- All DDL guarded by table/column existence checks. No ::regclass casts.
-- =====================================================

BEGIN;

DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'b2b_schema' AND table_name = 'departments') THEN
    RAISE NOTICE 'b2b_schema.departments does not exist yet — skipping enhance. Run proto migration first.';
    RETURN;
  END IF;

  -- ── FK: business_id → organisations ──────────────────────────────────────
  IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'b2b_schema' AND table_name = 'organisations')
  AND NOT EXISTS (
    SELECT 1 FROM pg_constraint c
    JOIN pg_class t ON t.oid = c.conrelid
    JOIN pg_namespace n ON n.oid = t.relnamespace
    WHERE n.nspname = 'b2b_schema' AND t.relname = 'departments' AND c.conname = 'fk_departments_business_id'
  ) THEN
    -- Only add FK if no orphan rows exist
    IF NOT EXISTS (
      SELECT 1 FROM b2b_schema.departments d
      WHERE d.business_id IS NOT NULL
        AND NOT EXISTS (SELECT 1 FROM b2b_schema.organisations o WHERE o.organisation_id = d.business_id)
    ) THEN
      EXECUTE 'ALTER TABLE b2b_schema.departments ADD CONSTRAINT fk_departments_business_id FOREIGN KEY (business_id) REFERENCES b2b_schema.organisations(organisation_id) ON DELETE RESTRICT';
    ELSE
      RAISE NOTICE 'Skipping fk_departments_business_id — orphan business_id rows found in departments';
    END IF;
  END IF;

  -- ── Active records index ──────────────────────────────────────────────────
  IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'b2b_schema' AND table_name = 'departments' AND column_name = 'deleted_at') THEN
    EXECUTE 'CREATE INDEX IF NOT EXISTS idx_departments_active ON b2b_schema.departments(department_id) WHERE deleted_at IS NULL';
  END IF;

  -- ── business_id index (most common query) ─────────────────────────────────
  IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'b2b_schema' AND table_name = 'departments' AND column_name = 'business_id') THEN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'b2b_schema' AND table_name = 'departments' AND column_name = 'deleted_at') THEN
      EXECUTE 'CREATE INDEX IF NOT EXISTS idx_departments_business_id ON b2b_schema.departments(business_id) WHERE deleted_at IS NULL';
    ELSE
      EXECUTE 'CREATE INDEX IF NOT EXISTS idx_departments_business_id ON b2b_schema.departments(business_id)';
    END IF;
  END IF;

  -- ── business_id + name composite index ───────────────────────────────────
  IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'b2b_schema' AND table_name = 'departments' AND column_name = 'name') THEN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'b2b_schema' AND table_name = 'departments' AND column_name = 'deleted_at') THEN
      EXECUTE 'CREATE INDEX IF NOT EXISTS idx_departments_business_name ON b2b_schema.departments(business_id, name) WHERE deleted_at IS NULL';
    ELSE
      EXECUTE 'CREATE INDEX IF NOT EXISTS idx_departments_business_name ON b2b_schema.departments(business_id, name)';
    END IF;
  END IF;

  -- ── timestamps ────────────────────────────────────────────────────────────
  IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'b2b_schema' AND table_name = 'departments' AND column_name = 'created_at') THEN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'b2b_schema' AND table_name = 'departments' AND column_name = 'deleted_at') THEN
      EXECUTE 'CREATE INDEX IF NOT EXISTS idx_departments_created_at ON b2b_schema.departments(created_at DESC) WHERE deleted_at IS NULL';
    ELSE
      EXECUTE 'CREATE INDEX IF NOT EXISTS idx_departments_created_at ON b2b_schema.departments(created_at DESC)';
    END IF;
  END IF;

  IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'b2b_schema' AND table_name = 'departments' AND column_name = 'updated_at') THEN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'b2b_schema' AND table_name = 'departments' AND column_name = 'deleted_at') THEN
      EXECUTE 'CREATE INDEX IF NOT EXISTS idx_departments_updated_at ON b2b_schema.departments(updated_at DESC) WHERE deleted_at IS NULL';
    ELSE
      EXECUTE 'CREATE INDEX IF NOT EXISTS idx_departments_updated_at ON b2b_schema.departments(updated_at DESC)';
    END IF;
  END IF;

  -- ── updated_at trigger ────────────────────────────────────────────────────
  IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'b2b_schema' AND table_name = 'departments' AND column_name = 'updated_at') THEN
    EXECUTE 'CREATE OR REPLACE FUNCTION b2b_schema.trg_departments_updated_at() RETURNS TRIGGER AS $body$ BEGIN NEW.updated_at = CURRENT_TIMESTAMP; RETURN NEW; END; $body$ LANGUAGE plpgsql';
    EXECUTE 'DROP TRIGGER IF EXISTS trg_departments_update ON b2b_schema.departments';
    EXECUTE 'CREATE TRIGGER trg_departments_update BEFORE UPDATE ON b2b_schema.departments FOR EACH ROW EXECUTE FUNCTION b2b_schema.trg_departments_updated_at()';
  END IF;

  -- NOTE: trg_sync_dept_employee_no trigger on b2b_schema.employees is installed
  -- in 20260301_004_enhance_employees.up.sql AFTER the employee proto migration
  -- adds deleted_at to the employees table. migration_order for departments (91)
  -- runs before employee proto migration (92) adds deleted_at.

  -- ── RLS ───────────────────────────────────────────────────────────────────
  EXECUTE 'ALTER TABLE b2b_schema.departments ENABLE ROW LEVEL SECURITY';
  EXECUTE 'DROP POLICY IF EXISTS departments_isolation ON b2b_schema.departments';
  EXECUTE $pol$
    CREATE POLICY departments_isolation ON b2b_schema.departments
    USING (
      current_setting('app.current_organisation_id', TRUE) IS NULL
      OR current_setting('app.current_organisation_id', TRUE) = ''
      OR business_id::TEXT = current_setting('app.current_organisation_id', TRUE)
    )
  $pol$;

  EXECUTE $cmt$COMMENT ON TABLE b2b_schema.departments IS 'B2B departments — enhanced with FK to organisations, employee_no sync trigger, RLS.'$cmt$;

END $$;

COMMIT;
