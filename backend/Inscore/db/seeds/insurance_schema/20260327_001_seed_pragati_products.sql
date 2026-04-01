-- ============================================================
-- Seed: Pragati Insurance Products & Plans
-- Source: documentation/KBank/Non-life-products + insurer-portal/lib/pragati-workbook-data.json
-- Products: PPA-001, NIB-001, PET-001, TRV-001, SSB-001
-- Constraint: product_code MUST match ^[A-Z]{3}-[0-9]{3}$
-- Status: PRODUCT_STATUS_ACTIVE
-- Category: explicit proto enum strings to match live `products` values
-- Note: created_by references a real authn_schema.users.user_id (system admin)
-- ============================================================

DO $$
DECLARE
    v_created_by uuid;
BEGIN
    -- Use first available user or skip product seed on empty authn user table
    SELECT user_id INTO v_created_by FROM authn_schema.users LIMIT 1;

    IF v_created_by IS NULL THEN
        RAISE NOTICE 'No authn user available for Pragati product seed; skipping products and plans';
    ELSE
        -- ── Products ───────────────────────────────────────────────
        INSERT INTO insurance_schema.products (
            product_id, product_code, product_name, category, description,
            base_premium, min_sum_insured, max_sum_insured,
            min_tenure_months, max_tenure_months, status, created_by, created_at
        ) VALUES
        (gen_random_uuid(), 'PPA-001', 'People''s Personal Accident', 'PRODUCT_CATEGORY_LIFE',
         'Annual accidental death and permanent disability cover for mass market. Fixed BDT 100,000 sum insured per person. Covers accidental death, permanent total/partial disability.',
         7400, 5000000, 10000000, 1, 12, 'PRODUCT_STATUS_ACTIVE', v_created_by, now()),

        (gen_random_uuid(), 'NIB-001', 'Nibedita', 'PRODUCT_CATEGORY_LIFE',
         'Women-focused personal accident: accidental death, disability, childbirth-related death, trauma allowance and household goods damage. Pragati Insurance flagship women product.',
         10000, 5000000, 50000000, 1, 12, 'PRODUCT_STATUS_ACTIVE', v_created_by, now()),

        (gen_random_uuid(), 'PET-001', 'Cat and Dog Insurance', 'PRODUCT_CATEGORY_PET',
         'Pet accident and critical illness cover for cats and dogs aged 8 weeks to 10 years. Covers hospitalization, surgery, diagnostics and vet consultations. Reimbursement-based claims.',
         103500, 1000000, 3000000, 1, 12, 'PRODUCT_STATUS_ACTIVE', v_created_by, now()),

        (gen_random_uuid(), 'TRV-001', 'Overseas Mediclaim Travel', 'PRODUCT_CATEGORY_TRAVEL',
         'Medical expenses and hospitalization abroad for frequent travelers. Coverage for Schengen and worldwide destinations including USA and Canada. Emergency evacuation and repatriation included.',
         123900, 5000000, 10000000, 1, 60, 'PRODUCT_STATUS_ACTIVE', v_created_by, now()),

        (gen_random_uuid(), 'SSB-001', 'Sorbojonin Surokkha Bima', 'PRODUCT_CATEGORY_LIFE',
         'Bangladesh-only comprehensive accident protection up to BDT 200,000. Permanent total and partial disability benefits, daily hospitalization allowance. Pragati mass-market PA product.',
         11500, 10000000, 20000000, 1, 12, 'PRODUCT_STATUS_ACTIVE', v_created_by, now())

        ON CONFLICT (product_code) DO UPDATE
        SET
            product_name = EXCLUDED.product_name,
            category = EXCLUDED.category,
            description = EXCLUDED.description,
            base_premium = EXCLUDED.base_premium,
            min_sum_insured = EXCLUDED.min_sum_insured,
            max_sum_insured = EXCLUDED.max_sum_insured,
            min_tenure_months = EXCLUDED.min_tenure_months,
            max_tenure_months = EXCLUDED.max_tenure_months,
            status = EXCLUDED.status,
            created_by = EXCLUDED.created_by,
            updated_at = now();

        -- ── Product Plans ──────────────────────────────────────────
        INSERT INTO insurance_schema.product_plans (
            plan_id, product_id, plan_name, plan_description,
            premium_amount, min_sum_insured, max_sum_insured
        )
        SELECT
            gen_random_uuid(),
            p.product_id,
            plans.plan_name,
            plans.plan_description,
            plans.premium_amount,
            plans.min_si,
            plans.max_si
        FROM insurance_schema.products p
        JOIN (VALUES
            ('PPA-001', 'Standard Plan',           'Fixed BDT 100,000 per person annual accident cover',        7400,    5000000,  10000000),
            ('NIB-001', 'Bronze - Women Basic',    'Capital sum BDT 50,000 for accidental events',              10000,   5000000,  10000000),
            ('NIB-001', 'Silver - Women Enhanced', 'Capital sum BDT 100,000 with trauma allowance',             20000,   10000000, 20000000),
            ('NIB-001', 'Gold - Women Premium',    'Capital sum BDT 200,000 full benefit schedule',             30000,   20000000, 50000000),
            ('PET-001', 'Basic Pet Cover',         'BDT 10,000 - accidental injury only',                       103500,  1000000,  1500000),
            ('PET-001', 'Standard Pet Cover',      'BDT 20,000 - accident and critical illness',                207000,  1500000,  2500000),
            ('PET-001', 'Premium Pet Cover',       'BDT 30,000 - accident, illness and surgery',                310500,  2500000,  3000000),
            ('TRV-001', 'Plan A Non-Schengen',     'Worldwide excl USA/Canada - USD 50,000 medical cover',      123900,  5000000,  7000000),
            ('TRV-001', 'Plan B Worldwide',        'Worldwide incl USA/Canada - USD 100,000 medical cover',     226900,  7000000,  10000000),
            ('TRV-001', 'Schengen Plan A',         'Schengen countries - EUR 30,000 nil deductible',            154900,  5000000,  7000000),
            ('TRV-001', 'CFT Annual',              'Corporate Frequent Traveler annual - 30 days max per trip', 1170400, 7000000,  10000000),
            ('SSB-001', 'Standard Plan',           'BDT 100,000 annual accident cover',                         11500,   10000000, 15000000),
            ('SSB-001', 'Enhanced Plan',           'BDT 200,000 with permanent disability benefit',             23000,   15000000, 20000000)
        ) AS plans(product_code, plan_name, plan_description, premium_amount, min_si, max_si)
        ON (p.product_code = plans.product_code AND p.deleted_at IS NULL)
        WHERE NOT EXISTS (
            SELECT 1 FROM insurance_schema.product_plans pp
            WHERE pp.product_id = p.product_id AND pp.plan_name = plans.plan_name
        );
    END IF;
END $$;
