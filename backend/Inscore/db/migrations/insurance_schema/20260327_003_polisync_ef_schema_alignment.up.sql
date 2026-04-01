-- Migration: 20260327_003_polisync_ef_schema_alignment
-- Reason: Align the insurance_schema with PoliSync EF Core model.
--         These changes were applied manually during B2C API testing (2026-03-27)
--         and are now made permanent in the migration history.
--
-- Changes applied:
--   1. Widened quotations.insurer_name from VARCHAR(255) to TEXT
--      Root cause: PoliSync stores serialized quotation metadata JSON in this column.
--   2. Made orders.plan_id nullable
--      Root cause: B2C purchase flows may not always select a specific plan.
--   3. Seeded 3 initial insurance products (HLT-001, LIF-001, MTR-001)
--   4. Seeded default Standard Plans for each product
--
-- Note: Items 1 and 2 were already applied via manual ALTER TABLE.
--       This migration uses IF EXISTS / DO NOTHING guards to be idempotent.

-- 1. Widen quotations.insurer_name to TEXT (already applied manually)
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'insurance_schema'
          AND table_name = 'quotations'
          AND column_name = 'insurer_name'
          AND data_type = 'character varying'
    ) THEN
        ALTER TABLE insurance_schema.quotations ALTER COLUMN insurer_name TYPE TEXT;
        RAISE NOTICE 'Widened insurance_schema.quotations.insurer_name to TEXT';
    ELSE
        RAISE NOTICE 'insurance_schema.quotations.insurer_name is already TEXT — skipping';
    END IF;
END $$;

-- 2. Make orders.plan_id nullable (already applied manually)
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'insurance_schema'
          AND table_name = 'orders'
          AND column_name = 'plan_id'
          AND is_nullable = 'NO'
    ) THEN
        ALTER TABLE insurance_schema.orders ALTER COLUMN plan_id DROP NOT NULL;
        CREATE INDEX IF NOT EXISTS idx_orders_plan_id_non_null
            ON insurance_schema.orders (plan_id)
            WHERE plan_id IS NOT NULL;
        RAISE NOTICE 'Made insurance_schema.orders.plan_id nullable';
    ELSE
        RAISE NOTICE 'insurance_schema.orders.plan_id is already nullable — skipping';
    END IF;
END $$;

-- 3. Seed initial products when a valid creator exists.
DO $$
DECLARE
    seed_creator UUID;
BEGIN
    SELECT u.user_id
    INTO seed_creator
    FROM authn_schema.users u
    WHERE u.deleted_at IS NULL
    ORDER BY u.created_at ASC, u.user_id ASC
    LIMIT 1;

    IF seed_creator IS NULL THEN
        RAISE NOTICE 'No authn user available for products.created_by; skipping default product seed';
    ELSE
        INSERT INTO insurance_schema.products (
            product_id, product_code, product_name, category, description,
            base_premium, min_sum_insured, max_sum_insured,
            min_tenure_months, max_tenure_months,
            status, base_premium_currency, min_sum_insured_currency, max_sum_insured_currency,
            plans, available_riders, pricing_config, created_by
        ) VALUES
        (
            gen_random_uuid(), 'LIF-001', 'LifeShield Term Plan', 'LIFE',
            'Comprehensive term life insurance with death benefit coverage. Ideal for income earners seeking affordable protection for their family.',
            250000, 50000000, 5000000000, 12, 240, 'ACTIVE', 'BDT', 'BDT', 'BDT',
            '[{"plan_id":"basic","plan_name":"Basic"},{"plan_id":"premium","plan_name":"Premium"}]'::jsonb,
            '[{"rider_id":"acc","rider_name":"Accidental Death","premium_amount":50000}]'::jsonb,
            '{}'::jsonb,
            seed_creator
        ),
        (
            gen_random_uuid(), 'HLT-001', 'MediCare Individual Health', 'HEALTH',
            'Individual health insurance covering hospitalization, surgery and critical illness.',
            180000, 10000000, 500000000, 12, 60, 'ACTIVE', 'BDT', 'BDT', 'BDT',
            '[{"plan_id":"silver","plan_name":"Silver"},{"plan_id":"gold","plan_name":"Gold"}]'::jsonb,
            '[{"rider_id":"mat","rider_name":"Maternity Benefit","premium_amount":30000}]'::jsonb,
            '{}'::jsonb,
            seed_creator
        ),
        (
            gen_random_uuid(), 'MTR-001', 'AutoGuard Comprehensive Motor', 'MOTOR',
            'Motor vehicle insurance covering own damage, third-party liability, theft and natural disasters.',
            350000, 30000000, 2000000000, 12, 12, 'ACTIVE', 'BDT', 'BDT', 'BDT',
            '[{"plan_id":"tpo","plan_name":"Third Party Only"},{"plan_id":"comp","plan_name":"Comprehensive"}]'::jsonb,
            '[]'::jsonb,
            '{}'::jsonb,
            seed_creator
        )
        ON CONFLICT (product_code) DO NOTHING;
    END IF;
END $$;

-- 4. Seed default Standard Plans for each product
INSERT INTO insurance_schema.product_plans (plan_id, product_id, plan_name, plan_description, premium_amount, min_sum_insured, max_sum_insured)
SELECT
    gen_random_uuid(),
    p.product_id,
    'Standard Plan',
    'Standard coverage plan for ' || p.product_name,
    p.base_premium,
    p.min_sum_insured,
    p.max_sum_insured
FROM insurance_schema.products p
WHERE p.product_code IN ('LIF-001', 'HLT-001', 'MTR-001')
  AND p.deleted_at IS NULL
  AND NOT EXISTS (
    SELECT 1 FROM insurance_schema.product_plans pp WHERE pp.product_id = p.product_id
);
