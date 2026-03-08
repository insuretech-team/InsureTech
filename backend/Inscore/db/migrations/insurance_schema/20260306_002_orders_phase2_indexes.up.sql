-- Enhancement: add missing indexes for orders phase-2 columns
-- These indexes correspond to @inject_tag gorm index annotations on order.proto fields 19-32.
-- No CREATE/ALTER/ADD COLUMN — columns were created by proto-driven GORM migration.
-- Date: 2026-03-06

CREATE UNIQUE INDEX IF NOT EXISTS idx_orders_idempotency_key
  ON insurance_schema.orders (idempotency_key)
  WHERE idempotency_key IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_orders_invoice_id
  ON insurance_schema.orders (invoice_id)
  WHERE invoice_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_orders_organisation_id
  ON insurance_schema.orders (organisation_id)
  WHERE organisation_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_orders_purchase_order_id
  ON insurance_schema.orders (purchase_order_id)
  WHERE purchase_order_id IS NOT NULL;
