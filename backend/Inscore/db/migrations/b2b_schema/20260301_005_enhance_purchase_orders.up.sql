-- =====================================================
-- Production Enhancement: b2b_schema.purchase_orders
-- Tables + columns created by proto-first migration generator.
-- This file: FK constraints, indexes, triggers, RLS ONLY.
-- Do NOT add columns here — add fields to purchase_order.proto instead.
-- All DDL guarded by table/column existence checks. No ::regclass casts.
-- =====================================================

BEGIN;

DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'b2b_schema' AND table_name = 'purchase_orders') THEN
    RAISE NOTICE 'b2b_schema.purchase_orders does not exist yet — skipping enhance. Run proto migration first.';
    RETURN;
  END IF;

  -- ── FK: business_id → organisations ──────────────────────────────────────
  IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'b2b_schema' AND table_name = 'organisations')
  AND NOT EXISTS (
    SELECT 1 FROM pg_constraint c
    JOIN pg_class t ON t.oid = c.conrelid
    JOIN pg_namespace n ON n.oid = t.relnamespace
    WHERE n.nspname = 'b2b_schema' AND t.relname = 'purchase_orders' AND c.conname = 'fk_purchase_orders_business_id'
  ) THEN
    IF NOT EXISTS (
      SELECT 1 FROM b2b_schema.purchase_orders po
      WHERE po.business_id IS NOT NULL
        AND NOT EXISTS (SELECT 1 FROM b2b_schema.organisations o WHERE o.organisation_id = po.business_id)
    ) THEN
      EXECUTE 'ALTER TABLE b2b_schema.purchase_orders ADD CONSTRAINT fk_purchase_orders_business_id FOREIGN KEY (business_id) REFERENCES b2b_schema.organisations(organisation_id) ON DELETE RESTRICT';
    ELSE
      RAISE NOTICE 'Skipping fk_purchase_orders_business_id — orphan business_id rows found in purchase_orders';
    END IF;
  END IF;

  -- ── FK: department_id → departments ──────────────────────────────────────
  IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'b2b_schema' AND table_name = 'departments')
  AND NOT EXISTS (
    SELECT 1 FROM pg_constraint c
    JOIN pg_class t ON t.oid = c.conrelid
    JOIN pg_namespace n ON n.oid = t.relnamespace
    WHERE n.nspname = 'b2b_schema' AND t.relname = 'purchase_orders' AND c.conname = 'fk_purchase_orders_department_id'
  ) THEN
    IF NOT EXISTS (
      SELECT 1 FROM b2b_schema.purchase_orders po
      WHERE po.department_id IS NOT NULL
        AND NOT EXISTS (SELECT 1 FROM b2b_schema.departments d WHERE d.department_id = po.department_id)
    ) THEN
      EXECUTE 'ALTER TABLE b2b_schema.purchase_orders ADD CONSTRAINT fk_purchase_orders_department_id FOREIGN KEY (department_id) REFERENCES b2b_schema.departments(department_id) ON DELETE RESTRICT';
    ELSE
      RAISE NOTICE 'Skipping fk_purchase_orders_department_id — orphan department_id rows found in purchase_orders';
    END IF;
  END IF;

  -- ── Unique: PO number per organisation ───────────────────────────────────
  IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE schemaname = 'b2b_schema' AND tablename = 'purchase_orders' AND indexname = 'uq_purchase_orders_number_per_org') THEN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'b2b_schema' AND table_name = 'purchase_orders' AND column_name = 'deleted_at') THEN
      EXECUTE 'CREATE UNIQUE INDEX uq_purchase_orders_number_per_org ON b2b_schema.purchase_orders(business_id, purchase_order_number) WHERE deleted_at IS NULL';
    ELSE
      EXECUTE 'CREATE UNIQUE INDEX uq_purchase_orders_number_per_org ON b2b_schema.purchase_orders(business_id, purchase_order_number)';
    END IF;
  END IF;

  -- ── Indexes ───────────────────────────────────────────────────────────────
  IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'b2b_schema' AND table_name = 'purchase_orders' AND column_name = 'deleted_at') THEN
    EXECUTE 'CREATE INDEX IF NOT EXISTS idx_purchase_orders_active ON b2b_schema.purchase_orders(purchase_order_id) WHERE deleted_at IS NULL';
    EXECUTE 'CREATE INDEX IF NOT EXISTS idx_purchase_orders_business_id ON b2b_schema.purchase_orders(business_id) WHERE deleted_at IS NULL';
    EXECUTE 'CREATE INDEX IF NOT EXISTS idx_purchase_orders_department_id ON b2b_schema.purchase_orders(department_id) WHERE deleted_at IS NULL';
    EXECUTE 'CREATE INDEX IF NOT EXISTS idx_purchase_orders_status ON b2b_schema.purchase_orders(status) WHERE deleted_at IS NULL';
    EXECUTE 'CREATE INDEX IF NOT EXISTS idx_purchase_orders_business_status ON b2b_schema.purchase_orders(business_id, status) WHERE deleted_at IS NULL';
    EXECUTE 'CREATE INDEX IF NOT EXISTS idx_purchase_orders_plan_id ON b2b_schema.purchase_orders(plan_id) WHERE deleted_at IS NULL';
  ELSE
    EXECUTE 'CREATE INDEX IF NOT EXISTS idx_purchase_orders_business_id ON b2b_schema.purchase_orders(business_id)';
    EXECUTE 'CREATE INDEX IF NOT EXISTS idx_purchase_orders_department_id ON b2b_schema.purchase_orders(department_id)';
    EXECUTE 'CREATE INDEX IF NOT EXISTS idx_purchase_orders_status ON b2b_schema.purchase_orders(status)';
  END IF;

  -- ── timestamps ────────────────────────────────────────────────────────────
  IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'b2b_schema' AND table_name = 'purchase_orders' AND column_name = 'created_at') THEN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'b2b_schema' AND table_name = 'purchase_orders' AND column_name = 'deleted_at') THEN
      EXECUTE 'CREATE INDEX IF NOT EXISTS idx_purchase_orders_created_at ON b2b_schema.purchase_orders(created_at DESC) WHERE deleted_at IS NULL';
    ELSE
      EXECUTE 'CREATE INDEX IF NOT EXISTS idx_purchase_orders_created_at ON b2b_schema.purchase_orders(created_at DESC)';
    END IF;
  END IF;

  IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'b2b_schema' AND table_name = 'purchase_orders' AND column_name = 'updated_at') THEN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'b2b_schema' AND table_name = 'purchase_orders' AND column_name = 'deleted_at') THEN
      EXECUTE 'CREATE INDEX IF NOT EXISTS idx_purchase_orders_updated_at ON b2b_schema.purchase_orders(updated_at DESC) WHERE deleted_at IS NULL';
    ELSE
      EXECUTE 'CREATE INDEX IF NOT EXISTS idx_purchase_orders_updated_at ON b2b_schema.purchase_orders(updated_at DESC)';
    END IF;
  END IF;

  -- ── updated_at trigger ────────────────────────────────────────────────────
  IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'b2b_schema' AND table_name = 'purchase_orders' AND column_name = 'updated_at') THEN
    EXECUTE 'CREATE OR REPLACE FUNCTION b2b_schema.trg_purchase_orders_updated_at() RETURNS TRIGGER AS $body$ BEGIN NEW.updated_at = CURRENT_TIMESTAMP; RETURN NEW; END; $body$ LANGUAGE plpgsql';
    EXECUTE 'DROP TRIGGER IF EXISTS trg_purchase_orders_update ON b2b_schema.purchase_orders';
    EXECUTE 'CREATE TRIGGER trg_purchase_orders_update BEFORE UPDATE ON b2b_schema.purchase_orders FOR EACH ROW EXECUTE FUNCTION b2b_schema.trg_purchase_orders_updated_at()';
  END IF;

  -- ── RLS ───────────────────────────────────────────────────────────────────
  EXECUTE 'ALTER TABLE b2b_schema.purchase_orders ENABLE ROW LEVEL SECURITY';
  EXECUTE 'DROP POLICY IF EXISTS purchase_orders_isolation ON b2b_schema.purchase_orders';
  EXECUTE $pol$
    CREATE POLICY purchase_orders_isolation ON b2b_schema.purchase_orders
    USING (
      current_setting('app.current_organisation_id', TRUE) IS NULL
      OR current_setting('app.current_organisation_id', TRUE) = ''
      OR business_id::TEXT = current_setting('app.current_organisation_id', TRUE)
    )
  $pol$;

  EXECUTE $cmt$COMMENT ON TABLE b2b_schema.purchase_orders IS 'Group insurance purchase orders — enhanced with FK constraints, unique PO number per org, status indexes, RLS.'$cmt$;

END $$;

COMMIT;
