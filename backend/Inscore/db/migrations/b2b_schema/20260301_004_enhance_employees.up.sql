-- =====================================================
-- Production Enhancement: b2b_schema.employees
-- Tables + columns created by proto-first migration generator.
-- This file: FK constraints, indexes, triggers, RLS ONLY.
-- Do NOT add columns here — add fields to employee.proto instead.
-- All DDL guarded by table/column existence checks. No ::regclass casts.
-- New columns (email, mobile_number, date_of_birth, date_of_joining,
-- gender, user_id) are defined in employee.proto fields 14-19.
-- =====================================================

BEGIN;

DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'b2b_schema' AND table_name = 'employees') THEN
    RAISE NOTICE 'b2b_schema.employees does not exist yet — skipping enhance. Run proto migration first.';
    RETURN;
  END IF;

  -- ── FK: business_id → organisations ──────────────────────────────────────
  IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'b2b_schema' AND table_name = 'organisations')
  AND NOT EXISTS (
    SELECT 1 FROM pg_constraint c
    JOIN pg_class t ON t.oid = c.conrelid
    JOIN pg_namespace n ON n.oid = t.relnamespace
    WHERE n.nspname = 'b2b_schema' AND t.relname = 'employees' AND c.conname = 'fk_employees_business_id'
  ) THEN
    IF NOT EXISTS (
      SELECT 1 FROM b2b_schema.employees e
      WHERE e.business_id IS NOT NULL
        AND NOT EXISTS (SELECT 1 FROM b2b_schema.organisations o WHERE o.organisation_id = e.business_id)
    ) THEN
      EXECUTE 'ALTER TABLE b2b_schema.employees ADD CONSTRAINT fk_employees_business_id FOREIGN KEY (business_id) REFERENCES b2b_schema.organisations(organisation_id) ON DELETE RESTRICT';
    ELSE
      RAISE NOTICE 'Skipping fk_employees_business_id — orphan business_id rows found in employees';
    END IF;
  END IF;

  -- ── FK: department_id → departments ──────────────────────────────────────
  IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'b2b_schema' AND table_name = 'departments')
  AND NOT EXISTS (
    SELECT 1 FROM pg_constraint c
    JOIN pg_class t ON t.oid = c.conrelid
    JOIN pg_namespace n ON n.oid = t.relnamespace
    WHERE n.nspname = 'b2b_schema' AND t.relname = 'employees' AND c.conname = 'fk_employees_department_id'
  ) THEN
    IF NOT EXISTS (
      SELECT 1 FROM b2b_schema.employees e
      WHERE e.department_id IS NOT NULL
        AND NOT EXISTS (SELECT 1 FROM b2b_schema.departments d WHERE d.department_id = e.department_id)
    ) THEN
      EXECUTE 'ALTER TABLE b2b_schema.employees ADD CONSTRAINT fk_employees_department_id FOREIGN KEY (department_id) REFERENCES b2b_schema.departments(department_id) ON DELETE RESTRICT';
    ELSE
      RAISE NOTICE 'Skipping fk_employees_department_id — orphan department_id rows found in employees';
    END IF;
  END IF;

  -- ── Unique: employee_id per organisation ──────────────────────────────────
  IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE schemaname = 'b2b_schema' AND tablename = 'employees' AND indexname = 'uq_employees_id_per_org') THEN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'b2b_schema' AND table_name = 'employees' AND column_name = 'deleted_at') THEN
      EXECUTE 'CREATE UNIQUE INDEX uq_employees_id_per_org ON b2b_schema.employees(business_id, employee_id) WHERE deleted_at IS NULL';
    ELSE
      EXECUTE 'CREATE UNIQUE INDEX uq_employees_id_per_org ON b2b_schema.employees(business_id, employee_id)';
    END IF;
  END IF;

  -- ── Active records index ──────────────────────────────────────────────────
  IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'b2b_schema' AND table_name = 'employees' AND column_name = 'deleted_at') THEN
    EXECUTE 'CREATE INDEX IF NOT EXISTS idx_employees_active ON b2b_schema.employees(employee_uuid) WHERE deleted_at IS NULL';
  END IF;

  -- ── business_id (tenant isolation — most common WHERE clause) ─────────────
  IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'b2b_schema' AND table_name = 'employees' AND column_name = 'deleted_at') THEN
    EXECUTE 'CREATE INDEX IF NOT EXISTS idx_employees_business_id ON b2b_schema.employees(business_id) WHERE deleted_at IS NULL';
    EXECUTE 'CREATE INDEX IF NOT EXISTS idx_employees_department_id ON b2b_schema.employees(department_id) WHERE deleted_at IS NULL';
    EXECUTE 'CREATE INDEX IF NOT EXISTS idx_employees_business_status ON b2b_schema.employees(business_id, status) WHERE deleted_at IS NULL';
    EXECUTE 'CREATE INDEX IF NOT EXISTS idx_employees_dept_status ON b2b_schema.employees(department_id, status) WHERE deleted_at IS NULL';
  ELSE
    EXECUTE 'CREATE INDEX IF NOT EXISTS idx_employees_business_id ON b2b_schema.employees(business_id)';
    EXECUTE 'CREATE INDEX IF NOT EXISTS idx_employees_department_id ON b2b_schema.employees(department_id)';
  END IF;

  -- ── user_id index (B2C bridge: find employee for a logged-in B2C user) ────
  IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'b2b_schema' AND table_name = 'employees' AND column_name = 'user_id') THEN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'b2b_schema' AND table_name = 'employees' AND column_name = 'deleted_at') THEN
      EXECUTE 'CREATE INDEX IF NOT EXISTS idx_employees_user_id ON b2b_schema.employees(user_id) WHERE user_id IS NOT NULL AND deleted_at IS NULL';
    ELSE
      EXECUTE 'CREATE INDEX IF NOT EXISTS idx_employees_user_id ON b2b_schema.employees(user_id) WHERE user_id IS NOT NULL';
    END IF;
  END IF;

  -- ── assigned_plan_id index ────────────────────────────────────────────────
  IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'b2b_schema' AND table_name = 'employees' AND column_name = 'assigned_plan_id') THEN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'b2b_schema' AND table_name = 'employees' AND column_name = 'deleted_at') THEN
      EXECUTE 'CREATE INDEX IF NOT EXISTS idx_employees_assigned_plan ON b2b_schema.employees(assigned_plan_id) WHERE assigned_plan_id IS NOT NULL AND deleted_at IS NULL';
    ELSE
      EXECUTE 'CREATE INDEX IF NOT EXISTS idx_employees_assigned_plan ON b2b_schema.employees(assigned_plan_id) WHERE assigned_plan_id IS NOT NULL';
    END IF;
  END IF;

  -- ── timestamps ────────────────────────────────────────────────────────────
  IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'b2b_schema' AND table_name = 'employees' AND column_name = 'created_at') THEN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'b2b_schema' AND table_name = 'employees' AND column_name = 'deleted_at') THEN
      EXECUTE 'CREATE INDEX IF NOT EXISTS idx_employees_created_at ON b2b_schema.employees(created_at DESC) WHERE deleted_at IS NULL';
    ELSE
      EXECUTE 'CREATE INDEX IF NOT EXISTS idx_employees_created_at ON b2b_schema.employees(created_at DESC)';
    END IF;
  END IF;

  IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'b2b_schema' AND table_name = 'employees' AND column_name = 'updated_at') THEN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'b2b_schema' AND table_name = 'employees' AND column_name = 'deleted_at') THEN
      EXECUTE 'CREATE INDEX IF NOT EXISTS idx_employees_updated_at ON b2b_schema.employees(updated_at DESC) WHERE deleted_at IS NULL';
    ELSE
      EXECUTE 'CREATE INDEX IF NOT EXISTS idx_employees_updated_at ON b2b_schema.employees(updated_at DESC)';
    END IF;
  END IF;

  -- ── updated_at trigger ────────────────────────────────────────────────────
  IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'b2b_schema' AND table_name = 'employees' AND column_name = 'updated_at') THEN
    EXECUTE 'CREATE OR REPLACE FUNCTION b2b_schema.trg_employees_updated_at() RETURNS TRIGGER AS $body$ BEGIN NEW.updated_at = CURRENT_TIMESTAMP; RETURN NEW; END; $body$ LANGUAGE plpgsql';
    EXECUTE 'DROP TRIGGER IF EXISTS trg_employees_update ON b2b_schema.employees';
    EXECUTE 'CREATE TRIGGER trg_employees_update BEFORE UPDATE ON b2b_schema.employees FOR EACH ROW EXECUTE FUNCTION b2b_schema.trg_employees_updated_at()';
  END IF;

  -- ── departments.employee_no auto-sync trigger ─────────────────────────────
  -- Installed HERE (not in 003) because employees.deleted_at now exists.
  -- Handles: insert, soft-delete, un-delete, department transfer.
  IF EXISTS (SELECT 1 FROM information_schema.tables  WHERE table_schema = 'b2b_schema' AND table_name = 'departments')
  AND EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'b2b_schema' AND table_name = 'departments' AND column_name = 'employee_no')
  AND EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'b2b_schema' AND table_name = 'employees'  AND column_name = 'deleted_at') THEN
    EXECUTE $func$
      CREATE OR REPLACE FUNCTION b2b_schema.trg_sync_dept_employee_no()
      RETURNS TRIGGER AS $trig$
      BEGIN
        IF TG_OP = 'INSERT' AND NEW.deleted_at IS NULL THEN
          UPDATE b2b_schema.departments SET employee_no = employee_no + 1 WHERE department_id = NEW.department_id;
        ELSIF TG_OP = 'UPDATE' THEN
          IF NEW.deleted_at IS NOT NULL AND OLD.deleted_at IS NULL THEN
            UPDATE b2b_schema.departments SET employee_no = GREATEST(0, employee_no - 1) WHERE department_id = NEW.department_id;
          ELSIF NEW.deleted_at IS NULL AND OLD.deleted_at IS NOT NULL THEN
            UPDATE b2b_schema.departments SET employee_no = employee_no + 1 WHERE department_id = NEW.department_id;
          ELSIF NEW.department_id <> OLD.department_id AND NEW.deleted_at IS NULL THEN
            UPDATE b2b_schema.departments SET employee_no = GREATEST(0, employee_no - 1) WHERE department_id = OLD.department_id;
            UPDATE b2b_schema.departments SET employee_no = employee_no + 1             WHERE department_id = NEW.department_id;
          END IF;
        END IF;
        RETURN NEW;
      END;
      $trig$ LANGUAGE plpgsql
    $func$;
    EXECUTE 'DROP TRIGGER IF EXISTS trg_dept_employee_no_sync ON b2b_schema.employees';
    EXECUTE 'CREATE TRIGGER trg_dept_employee_no_sync AFTER INSERT OR UPDATE OF deleted_at, department_id ON b2b_schema.employees FOR EACH ROW EXECUTE FUNCTION b2b_schema.trg_sync_dept_employee_no()';
  END IF;

  -- ── organisations.total_employees auto-sync trigger ───────────────────────
  -- Installed HERE (not in 001) because employees.deleted_at now exists.
  IF EXISTS (SELECT 1 FROM information_schema.tables  WHERE table_schema = 'b2b_schema' AND table_name = 'organisations')
  AND EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'b2b_schema' AND table_name = 'organisations' AND column_name = 'total_employees')
  AND EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'b2b_schema' AND table_name = 'employees'    AND column_name = 'deleted_at') THEN
    EXECUTE $func$
      CREATE OR REPLACE FUNCTION b2b_schema.trg_sync_org_employee_count()
      RETURNS TRIGGER AS $trig$
      BEGIN
        IF TG_OP = 'INSERT' AND NEW.deleted_at IS NULL THEN
          UPDATE b2b_schema.organisations SET total_employees = total_employees + 1 WHERE organisation_id = NEW.business_id;
        ELSIF TG_OP = 'UPDATE' THEN
          IF NEW.deleted_at IS NOT NULL AND OLD.deleted_at IS NULL THEN
            UPDATE b2b_schema.organisations SET total_employees = GREATEST(0, total_employees - 1) WHERE organisation_id = NEW.business_id;
          ELSIF NEW.deleted_at IS NULL AND OLD.deleted_at IS NOT NULL THEN
            UPDATE b2b_schema.organisations SET total_employees = total_employees + 1 WHERE organisation_id = NEW.business_id;
          END IF;
        END IF;
        RETURN NEW;
      END;
      $trig$ LANGUAGE plpgsql
    $func$;
    EXECUTE 'DROP TRIGGER IF EXISTS trg_employee_count_sync ON b2b_schema.employees';
    EXECUTE 'CREATE TRIGGER trg_employee_count_sync AFTER INSERT OR UPDATE OF deleted_at ON b2b_schema.employees FOR EACH ROW EXECUTE FUNCTION b2b_schema.trg_sync_org_employee_count()';
  END IF;

  -- ── RLS ───────────────────────────────────────────────────────────────────
  EXECUTE 'ALTER TABLE b2b_schema.employees ENABLE ROW LEVEL SECURITY';
  EXECUTE 'DROP POLICY IF EXISTS employees_isolation ON b2b_schema.employees';
  EXECUTE $pol$
    CREATE POLICY employees_isolation ON b2b_schema.employees
    USING (
      current_setting('app.current_organisation_id', TRUE) IS NULL
      OR current_setting('app.current_organisation_id', TRUE) = ''
      OR business_id::TEXT = current_setting('app.current_organisation_id', TRUE)
    )
  $pol$;

  EXECUTE $cmt$COMMENT ON TABLE b2b_schema.employees IS 'B2B employees — enhanced with PII columns (email, mobile, dob, doj, gender, user_id), FK constraints, composite indexes, RLS.'$cmt$;
  IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'b2b_schema' AND table_name = 'employees' AND column_name = 'user_id') THEN
    EXECUTE $cmt2$COMMENT ON COLUMN b2b_schema.employees.user_id IS 'FK to authn_schema.users — set when B2C portal access is granted. Cross-schema FK enforced at application level.'$cmt2$;
  END IF;

END $$;

COMMIT;
