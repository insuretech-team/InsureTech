-- Migration: 20260327_002_make_orders_plan_id_nullable
-- Reason: B2C users creating orders without a specific plan should be allowed.
--         plan_id is optional in many insurance product flows (single-plan products).
--         The FK constraint remains to enforce referential integrity when a plan IS provided.
-- Applied manually on 2026-03-27 during B2C API testing.

-- Make plan_id nullable in orders table
ALTER TABLE insurance_schema.orders
    ALTER COLUMN plan_id DROP NOT NULL;

-- Add a partial index for non-null plan_id lookups (preserve query performance)
CREATE INDEX IF NOT EXISTS idx_orders_plan_id_non_null
    ON insurance_schema.orders (plan_id)
    WHERE plan_id IS NOT NULL;
