-- ============================================================
-- Seed: Align enum values and activate Pragati mapped products
-- Purpose:
--   1. Activate Pragati's five mapped products in both base catalog and
--      insurer-specific catalog
--   2. Replace legacy Pragati generic insurer-product mappings with the real
--      sellable set: PPA, Nibedita, Pet, Travel, SSB
-- Depends on:
--   - 20260327_001_seed_pragati_products.sql
--   - 20260327_002_seed_bangladesh_insurers.sql
--   - 20260327_004_seed_bangladesh_insurer_products.sql
-- ============================================================

-- Activate the five Pragati base products in the sellable product catalog.
UPDATE insurance_schema.products
SET
    status = 'PRODUCT_STATUS_ACTIVE',
    updated_at = now()
WHERE product_code IN ('PPA-001', 'NIB-001', 'PET-001', 'TRV-001', 'SSB-001')
  AND deleted_at IS NULL;

-- Remove obsolete generic Pragati insurer-product rows.
DELETE FROM insurance_schema.insurer_products
WHERE code IN ('PRAGATI-MOTOR', 'PRAGATI-HEALTH');

WITH pragati_insurer AS (
    SELECT insurer_id
    FROM insurance_schema.insurers
    WHERE code = 'PRAGATI'
),
pragati_products AS (
    SELECT product_id, product_code
    FROM insurance_schema.products
    WHERE product_code IN ('PPA-001', 'NIB-001', 'PET-001', 'TRV-001', 'SSB-001')
      AND deleted_at IS NULL
),
seed_data AS (
    SELECT
        v.seed_product_id,
        i.insurer_id,
        p.product_id AS base_product_id,
        v.code,
        v.name,
        v.status,
        v.min_sum_assured,
        v.max_sum_assured,
        v.min_premium,
        v.max_premium,
        v.min_entry_age,
        v.max_entry_age,
        v.max_maturity_age,
        v.min_term_years,
        v.max_term_years,
        v.premium_payment_modes,
        v.medical_required,
        v.medical_threshold,
        v.free_look_period_days,
        v.features,
        v.exclusions,
        v.effective_from,
        v.effective_to,
        v.audit_info
    FROM pragati_insurer i
    JOIN pragati_products p ON 1 = 1
    JOIN (
        VALUES
        (
            'a4f2f420-42d2-4c3a-92a0-100000000001'::uuid,
            'PPA-001',
            'PRAGATI-PPA',
            'People''s Personal Accident',
            'ACTIVE',
            '{"amount":5000000,"currency":"BDT","decimal_amount":50000}'::jsonb,
            '{"amount":10000000,"currency":"BDT","decimal_amount":100000}'::jsonb,
            '{"amount":7400,"currency":"BDT","decimal_amount":74}'::jsonb,
            '{"amount":7400,"currency":"BDT","decimal_amount":74}'::jsonb,
            NULL::integer,
            NULL::integer,
            NULL::integer,
            1,
            1,
            ARRAY['YEARLY']::text[],
            false,
            NULL::jsonb,
            15,
            '{
              "source_url": "https://pragatiinsurance.com/miscellaneous/",
              "public_benefits": [
                "personal accident cover",
                "accidental death benefit",
                "permanent disability benefit"
              ]
            }'::jsonb,
            NULL::jsonb,
            DATE '2026-03-27',
            NULL::date,
            '{}'::jsonb
        ),
        (
            'a4f2f420-42d2-4c3a-92a0-100000000003'::uuid,
            'NIB-001',
            'PRAGATI-NIBEDITA',
            'Nibedita',
            'ACTIVE',
            '{"amount":5000000,"currency":"BDT","decimal_amount":50000}'::jsonb,
            '{"amount":50000000,"currency":"BDT","decimal_amount":500000}'::jsonb,
            '{"amount":10000,"currency":"BDT","decimal_amount":100}'::jsonb,
            '{"amount":30000,"currency":"BDT","decimal_amount":300}'::jsonb,
            NULL::integer,
            NULL::integer,
            NULL::integer,
            1,
            1,
            ARRAY['YEARLY']::text[],
            false,
            NULL::jsonb,
            15,
            '{
              "source_url": "https://pragatiinsurance.com/miscellaneous/",
              "public_benefits": [
                "women-focused personal accident cover",
                "childbirth-related death benefit",
                "household goods damage benefit"
              ]
            }'::jsonb,
            NULL::jsonb,
            DATE '2026-03-27',
            NULL::date,
            '{}'::jsonb
        ),
        (
            'a4f2f420-42d2-4c3a-92a0-100000000004'::uuid,
            'PET-001',
            'PRAGATI-PET',
            'Cat and Dog Insurance',
            'ACTIVE',
            '{"amount":1000000,"currency":"BDT","decimal_amount":10000}'::jsonb,
            '{"amount":3000000,"currency":"BDT","decimal_amount":30000}'::jsonb,
            '{"amount":103500,"currency":"BDT","decimal_amount":1035}'::jsonb,
            '{"amount":310500,"currency":"BDT","decimal_amount":3105}'::jsonb,
            NULL::integer,
            NULL::integer,
            NULL::integer,
            1,
            1,
            ARRAY['YEARLY']::text[],
            false,
            NULL::jsonb,
            15,
            '{
              "source_url": "https://pragatiinsurance.com/miscellaneous/",
              "public_benefits": [
                "pet accident cover",
                "critical illness cover",
                "veterinary reimbursement"
              ]
            }'::jsonb,
            NULL::jsonb,
            DATE '2026-03-27',
            NULL::date,
            '{}'::jsonb
        ),
        (
            'a4f2f420-42d2-4c3a-92a0-100000000002'::uuid,
            'TRV-001',
            'PRAGATI-TRAVEL',
            'Overseas Mediclaim Insurance',
            'ACTIVE',
            '{"amount":5000000,"currency":"BDT","decimal_amount":50000}'::jsonb,
            '{"amount":10000000,"currency":"BDT","decimal_amount":100000}'::jsonb,
            '{"amount":123900,"currency":"BDT","decimal_amount":1239}'::jsonb,
            '{"amount":1170400,"currency":"BDT","decimal_amount":11704}'::jsonb,
            NULL::integer,
            NULL::integer,
            NULL::integer,
            1,
            5,
            ARRAY['SINGLE', 'YEARLY']::text[],
            false,
            NULL::jsonb,
            15,
            '{
              "source_url": "https://pragatiinsurance.com/overseas-mediclaim-insurance/",
              "plan_highlights": [
                "Schengen and worldwide destinations",
                "medical expenses and hospitalization",
                "emergency evacuation and repatriation"
              ]
            }'::jsonb,
            NULL::jsonb,
            DATE '2026-03-27',
            NULL::date,
            '{}'::jsonb
        ),
        (
            'a4f2f420-42d2-4c3a-92a0-100000000024'::uuid,
            'SSB-001',
            'PRAGATI-SSB',
            'Sorbojonin Surokkha Bima',
            'ACTIVE',
            '{"amount":10000000,"currency":"BDT","decimal_amount":100000}'::jsonb,
            '{"amount":20000000,"currency":"BDT","decimal_amount":200000}'::jsonb,
            '{"amount":11500,"currency":"BDT","decimal_amount":115}'::jsonb,
            '{"amount":23000,"currency":"BDT","decimal_amount":230}'::jsonb,
            NULL::integer,
            NULL::integer,
            NULL::integer,
            1,
            1,
            ARRAY['YEARLY']::text[],
            false,
            NULL::jsonb,
            15,
            '{
              "source_url": "https://pragatiinsurance.com/miscellaneous/",
              "public_benefits": [
                "bangladesh-only accident cover",
                "permanent disability benefit",
                "daily hospitalization allowance"
              ]
            }'::jsonb,
            NULL::jsonb,
            DATE '2026-03-27',
            NULL::date,
            '{}'::jsonb
        )
    ) AS v(
        seed_product_id,
        base_product_code,
        code,
        name,
        status,
        min_sum_assured,
        max_sum_assured,
        min_premium,
        max_premium,
        min_entry_age,
        max_entry_age,
        max_maturity_age,
        min_term_years,
        max_term_years,
        premium_payment_modes,
        medical_required,
        medical_threshold,
        free_look_period_days,
        features,
        exclusions,
        effective_from,
        effective_to,
        audit_info
    )
      ON p.product_code = v.base_product_code
)
INSERT INTO insurance_schema.insurer_products (
    product_id,
    insurer_id,
    base_product_id,
    code,
    name,
    status,
    min_sum_assured,
    max_sum_assured,
    min_premium,
    max_premium,
    min_entry_age,
    max_entry_age,
    max_maturity_age,
    min_term_years,
    max_term_years,
    premium_payment_modes,
    medical_required,
    medical_threshold,
    free_look_period_days,
    features,
    exclusions,
    effective_from,
    effective_to,
    audit_info
)
SELECT
    seed_product_id,
    insurer_id,
    base_product_id,
    code,
    name,
    status,
    min_sum_assured,
    max_sum_assured,
    min_premium,
    max_premium,
    min_entry_age,
    max_entry_age,
    max_maturity_age,
    min_term_years,
    max_term_years,
    premium_payment_modes,
    medical_required,
    medical_threshold,
    free_look_period_days,
    features,
    exclusions,
    effective_from,
    effective_to,
    audit_info
FROM seed_data
ON CONFLICT (code) DO UPDATE
SET
    insurer_id = EXCLUDED.insurer_id,
    base_product_id = EXCLUDED.base_product_id,
    name = EXCLUDED.name,
    status = EXCLUDED.status,
    min_sum_assured = EXCLUDED.min_sum_assured,
    max_sum_assured = EXCLUDED.max_sum_assured,
    min_premium = EXCLUDED.min_premium,
    max_premium = EXCLUDED.max_premium,
    min_entry_age = EXCLUDED.min_entry_age,
    max_entry_age = EXCLUDED.max_entry_age,
    max_maturity_age = EXCLUDED.max_maturity_age,
    min_term_years = EXCLUDED.min_term_years,
    max_term_years = EXCLUDED.max_term_years,
    premium_payment_modes = EXCLUDED.premium_payment_modes,
    medical_required = EXCLUDED.medical_required,
    medical_threshold = EXCLUDED.medical_threshold,
    free_look_period_days = EXCLUDED.free_look_period_days,
    features = EXCLUDED.features,
    exclusions = EXCLUDED.exclusions,
    effective_from = EXCLUDED.effective_from,
    effective_to = EXCLUDED.effective_to,
    audit_info = EXCLUDED.audit_info;
