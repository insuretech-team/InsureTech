-- ============================================================
-- Seed: Bangladesh Insurer Product Catalogs
-- Purpose: Seed insurer-specific product variants for publicly documented
--          Bangladeshi insurers without activating them for sale
-- Depends on:
--   - 20260327_001_seed_pragati_products.sql
--   - 20260327_002_seed_bangladesh_insurers.sql
-- Notes:
--   - `insurance_schema.insurer_products.status` uses compact DB values like
--     ACTIVE / INACTIVE rather than the longer proto enum names.
--   - Pragati's five web-seeded products are mapped into insurer_products and
--     activated to match the requested sellable product set.
--   - Only publicly visible product names and ranges from official sites/PDFs
--     are seeded. When public pages do not expose a value, the JSONB money
--     field is left NULL instead of inventing data.
-- Sources (official/public, checked 2026-03-27):
--   Pragati Insurance PLC
--     https://pragatiinsurance.com/miscellaneous/
--     https://pragatiinsurance.com/overseas-mediclaim-insurance/
--   Green Delta Insurance PLC
--     https://green-delta.com/personal-accident-insurance/
--     https://green-delta.com/nibedita/
--     https://green-delta.com/shurokkha-365/
--     https://green-delta.com/motor-insurance/
--     https://green-delta.com/travel-insurance/
--     https://green-delta.com/health-insurance/
--   MetLife Bangladesh
--     https://www.metlife.com.bd/solutions/health-and-protection/mciirop/
--     https://www.metlife.com.bd/content/dam/metlifecom/bd/PDFs/Brochure/hospitalcare.pdf
--     https://www.metlife.com.bd/content/dam/metlifecom/bd/PDFs/Brochure/EPP_plus.pdf
--     https://www.metlife.com.bd/content/dam/metlifecom/bd/PDFs/Brochure/mdpssuper.pdf
--     https://www.metlife.com.bd/content/dam/metlifecom/bd/PDFs/Brochure/endowment.pdf
--     https://www.metlife.com.bd/solutions/group-insurance/employee-benefits/group-death-disability/
--   Chartered Life Insurance PLC
--     https://charteredlifebd.org/prospectus.aspx
--     https://charteredlifebd.org/Corporateinformatiom.aspx
-- ============================================================

WITH insurer_lookup AS (
    SELECT insurer_id, code
    FROM insurance_schema.insurers
    WHERE code IN ('PRAGATI', 'GREENDELTA', 'METLIFE', 'CHARTERED')
),
product_lookup AS (
    SELECT product_id, product_code
    FROM insurance_schema.products
    WHERE product_code IN ('PPA-001', 'NIB-001', 'PET-001', 'TRV-001', 'SSB-001', 'HLT-001', 'LIF-001', 'MTR-001')
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
    FROM insurer_lookup i
    JOIN product_lookup p
      ON 1 = 1
    JOIN (
        VALUES
        (
            'a4f2f420-42d2-4c3a-92a0-100000000001'::uuid,
            'PRAGATI',
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
              "product_family": "miscellaneous",
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
            'a4f2f420-42d2-4c3a-92a0-100000000002'::uuid,
            'PRAGATI',
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
            'a4f2f420-42d2-4c3a-92a0-100000000025'::uuid,
            'PRAGATI',
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
            'a4f2f420-42d2-4c3a-92a0-100000000026'::uuid,
            'PRAGATI',
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
            'a4f2f420-42d2-4c3a-92a0-100000000024'::uuid,
            'PRAGATI',
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
        ),
        (
            'a4f2f420-42d2-4c3a-92a0-100000000005'::uuid,
            'GREENDELTA',
            'PPA-001',
            'GREENDELTA-PPA',
            'People''s Personal Accident Insurance',
            'INACTIVE',
            '{"amount":10000000,"currency":"BDT","decimal_amount":100000}'::jsonb,
            '{"amount":10000000,"currency":"BDT","decimal_amount":100000}'::jsonb,
            '{"amount":6000,"currency":"BDT","decimal_amount":60}'::jsonb,
            '{"amount":6000,"currency":"BDT","decimal_amount":60}'::jsonb,
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
              "source_url": "https://green-delta.com/personal-accident-insurance/",
              "public_benefits": [
                "accidental death",
                "permanent total disability",
                "permanent partial disability"
              ],
              "public_quote": "Taka 60 annual premium for Taka 100000 coverage"
            }'::jsonb,
            NULL::jsonb,
            DATE '2026-03-27',
            NULL::date,
            '{}'::jsonb
        ),
        (
            'a4f2f420-42d2-4c3a-92a0-100000000006'::uuid,
            'GREENDELTA',
            'NIB-001',
            'GREENDELTA-NIBEDITA',
            'Nibedita',
            'INACTIVE',
            '{"amount":15000000,"currency":"BDT","decimal_amount":150000}'::jsonb,
            '{"amount":45000000,"currency":"BDT","decimal_amount":450000}'::jsonb,
            '{"amount":23000,"currency":"BDT","decimal_amount":230}'::jsonb,
            '{"amount":69000,"currency":"BDT","decimal_amount":690}'::jsonb,
            18,
            50,
            NULL::integer,
            1,
            1,
            ARRAY['YEARLY']::text[],
            false,
            NULL::jsonb,
            15,
            '{
              "source_url": "https://green-delta.com/nibedita/",
              "plans": [
                {"premium_bdt": 230, "coverage_bdt": 150000},
                {"premium_bdt": 460, "coverage_bdt": 300000},
                {"premium_bdt": 690, "coverage_bdt": 450000}
              ],
              "public_benefits": [
                "death due to childbirth",
                "road accident benefit",
                "damage to household goods"
              ]
            }'::jsonb,
            NULL::jsonb,
            DATE '2026-03-27',
            NULL::date,
            '{}'::jsonb
        ),
        (
            'a4f2f420-42d2-4c3a-92a0-100000000007'::uuid,
            'GREENDELTA',
            'HLT-001',
            'GREENDELTA-S365',
            'Shurokkha 365',
            'INACTIVE',
            '{"amount":10000000,"currency":"BDT","decimal_amount":100000}'::jsonb,
            '{"amount":30000000,"currency":"BDT","decimal_amount":300000}'::jsonb,
            '{"amount":36500,"currency":"BDT","decimal_amount":365}'::jsonb,
            '{"amount":109500,"currency":"BDT","decimal_amount":1095}'::jsonb,
            18,
            59,
            NULL::integer,
            1,
            1,
            ARRAY['YEARLY']::text[],
            false,
            NULL::jsonb,
            15,
            '{
              "source_url": "https://green-delta.com/shurokkha-365/",
              "plans": [
                {"premium_bdt": 365, "daily_hospital_cash_bdt": 500, "annual_limit_bdt": 100000},
                {"premium_bdt": 730, "daily_hospital_cash_bdt": 1000, "annual_limit_bdt": 200000},
                {"premium_bdt": 1095, "daily_hospital_cash_bdt": 1500, "annual_limit_bdt": 300000}
              ],
              "public_benefits": [
                "hospital daily cash",
                "accidental death",
                "partial permanent disability"
              ]
            }'::jsonb,
            NULL::jsonb,
            DATE '2026-03-27',
            NULL::date,
            '{}'::jsonb
        ),
        (
            'a4f2f420-42d2-4c3a-92a0-100000000008'::uuid,
            'GREENDELTA',
            'TRV-001',
            'GREENDELTA-TRAVEL',
            'Travel Insurance',
            'INACTIVE',
            NULL::jsonb,
            NULL::jsonb,
            NULL::jsonb,
            NULL::jsonb,
            NULL::integer,
            NULL::integer,
            NULL::integer,
            1,
            1,
            ARRAY[]::text[],
            false,
            NULL::jsonb,
            15,
            '{
              "source_url": "https://green-delta.com/travel-insurance/",
              "product_family": "travel"
            }'::jsonb,
            NULL::jsonb,
            DATE '2026-03-27',
            NULL::date,
            '{}'::jsonb
        ),
        (
            'a4f2f420-42d2-4c3a-92a0-100000000009'::uuid,
            'GREENDELTA',
            'MTR-001',
            'GREENDELTA-MOTOR',
            'Motor Insurance',
            'INACTIVE',
            NULL::jsonb,
            NULL::jsonb,
            NULL::jsonb,
            NULL::jsonb,
            NULL::integer,
            NULL::integer,
            NULL::integer,
            1,
            1,
            ARRAY[]::text[],
            false,
            NULL::jsonb,
            15,
            '{
              "source_url": "https://green-delta.com/motor-insurance/",
              "product_family": "motor"
            }'::jsonb,
            NULL::jsonb,
            DATE '2026-03-27',
            NULL::date,
            '{}'::jsonb
        ),
        (
            'a4f2f420-42d2-4c3a-92a0-100000000010'::uuid,
            'GREENDELTA',
            'HLT-001',
            'GREENDELTA-HEALTH',
            'Health Insurance',
            'INACTIVE',
            NULL::jsonb,
            NULL::jsonb,
            NULL::jsonb,
            NULL::jsonb,
            NULL::integer,
            NULL::integer,
            NULL::integer,
            1,
            1,
            ARRAY[]::text[],
            false,
            NULL::jsonb,
            15,
            '{
              "source_url": "https://green-delta.com/health-insurance/",
              "product_family": "health"
            }'::jsonb,
            NULL::jsonb,
            DATE '2026-03-27',
            NULL::date,
            '{}'::jsonb
        ),
        (
            'a4f2f420-42d2-4c3a-92a0-100000000011'::uuid,
            'METLIFE',
            'HLT-001',
            'METLIFE-MCIIROP',
            'Critical Illness Insurance with Return of Premium',
            'INACTIVE',
            '{"amount":30000000,"currency":"BDT","decimal_amount":300000}'::jsonb,
            '{"amount":200000000,"currency":"BDT","decimal_amount":2000000}'::jsonb,
            '{"amount":140200,"currency":"BDT","decimal_amount":1402}'::jsonb,
            NULL::jsonb,
            18,
            55,
            65,
            10,
            10,
            ARRAY['ONE_TIME', 'QUARTERLY', 'SEMI_ANNUAL', 'ANNUAL']::text[],
            false,
            NULL::jsonb,
            15,
            '{
              "source_url": "https://www.metlife.com.bd/solutions/health-and-protection/mciirop/",
              "launch_source_url": "https://www.metlife.com.bd/about-us/newsroom/2025/august/metlife-launches-mciirop/",
              "public_benefits": [
                "10 critical illnesses",
                "100 percent return of premium on claim-free maturity",
                "accidental death benefit"
              ],
              "premium_example": "BDT 1402 monthly example from official launch note"
            }'::jsonb,
            NULL::jsonb,
            DATE '2026-03-27',
            NULL::date,
            '{}'::jsonb
        ),
        (
            'a4f2f420-42d2-4c3a-92a0-100000000012'::uuid,
            'METLIFE',
            'HLT-001',
            'METLIFE-HCARE',
            'Hospital Care',
            'INACTIVE',
            NULL::jsonb,
            NULL::jsonb,
            NULL::jsonb,
            NULL::jsonb,
            18,
            55,
            60,
            1,
            5,
            ARRAY['ANNUAL', 'SEMI_ANNUAL', 'QUARTERLY', 'MONTHLY']::text[],
            false,
            NULL::jsonb,
            15,
            '{
              "source_url": "https://www.metlife.com.bd/content/dam/metlifecom/bd/PDFs/Brochure/hospitalcare.pdf",
              "public_benefits": [
                "daily hospital cash from BDT 500 to BDT 5000",
                "ICU benefit from BDT 1000 to BDT 10000",
                "surgery benefit from BDT 10000 to BDT 150000"
              ]
            }'::jsonb,
            NULL::jsonb,
            DATE '2026-03-27',
            NULL::date,
            '{}'::jsonb
        ),
        (
            'a4f2f420-42d2-4c3a-92a0-100000000013'::uuid,
            'METLIFE',
            'LIF-001',
            'METLIFE-EPPPLUS',
            'Education Protection Plan Plus',
            'INACTIVE',
            NULL::jsonb,
            NULL::jsonb,
            NULL::jsonb,
            NULL::jsonb,
            NULL::integer,
            NULL::integer,
            NULL::integer,
            10,
            20,
            ARRAY[]::text[],
            false,
            NULL::jsonb,
            15,
            '{
              "source_url": "https://www.metlife.com.bd/content/dam/metlifecom/bd/PDFs/Brochure/EPP_plus.pdf",
              "public_benefits": [
                "child age from 30 days to 15 years",
                "education benefit from age 18 to 22",
                "maturity benefit equals face amount"
              ],
              "premium_example": "Official brochure example shows BDT 38580 yearly for 20-year term and face amount BDT 300000"
            }'::jsonb,
            NULL::jsonb,
            DATE '2026-03-27',
            NULL::date,
            '{}'::jsonb
        ),
        (
            'a4f2f420-42d2-4c3a-92a0-100000000014'::uuid,
            'METLIFE',
            'LIF-001',
            'METLIFE-MDPS',
            'Monthly Deposit Protection Scheme',
            'INACTIVE',
            '{"amount":50000000,"currency":"BDT","decimal_amount":500000}'::jsonb,
            NULL::jsonb,
            '{"amount":100000,"currency":"BDT","decimal_amount":1000}'::jsonb,
            NULL::jsonb,
            18,
            55,
            60,
            5,
            20,
            ARRAY['ANNUAL', 'HALF_YEARLY', 'QUARTERLY', 'MONTHLY']::text[],
            false,
            NULL::jsonb,
            15,
            '{
              "source_url": "https://www.metlife.com.bd/content/dam/metlifecom/bd/PDFs/Brochure/mdpssuper.pdf",
              "public_benefits": [
                "in-hospital daily benefit",
                "return of fund value at maturity",
                "critical illness coverage available on gold plan"
              ],
              "premium_floor_bdt": 1000,
              "public_face_amount_floor_bdt": 500000
            }'::jsonb,
            NULL::jsonb,
            DATE '2026-03-27',
            NULL::date,
            '{}'::jsonb
        ),
        (
            'a4f2f420-42d2-4c3a-92a0-100000000015'::uuid,
            'METLIFE',
            'LIF-001',
            'METLIFE-ENDOWTH',
            'MetLife Endowment - Growth',
            'INACTIVE',
            NULL::jsonb,
            NULL::jsonb,
            NULL::jsonb,
            NULL::jsonb,
            NULL::integer,
            NULL::integer,
            NULL::integer,
            16,
            30,
            ARRAY['ANNUAL', 'QUARTERLY', 'SEMI_ANNUAL', 'MONTHLY']::text[],
            false,
            NULL::jsonb,
            15,
            '{
              "source_url": "https://www.metlife.com.bd/content/dam/metlifecom/bd/PDFs/Brochure/endowment.pdf",
              "public_benefits": [
                "terms of 16 to 20 years, 25 years or 30 years",
                "cash value accumulation",
                "maturity payout options"
              ]
            }'::jsonb,
            NULL::jsonb,
            DATE '2026-03-27',
            NULL::date,
            '{}'::jsonb
        ),
        (
            'a4f2f420-42d2-4c3a-92a0-100000000016'::uuid,
            'METLIFE',
            'LIF-001',
            'METLIFE-GROUPDD',
            'Group Death and Disability',
            'INACTIVE',
            '{"amount":10000000,"currency":"BDT","decimal_amount":100000}'::jsonb,
            '{"amount":50000000,"currency":"BDT","decimal_amount":500000}'::jsonb,
            NULL::jsonb,
            NULL::jsonb,
            NULL::integer,
            NULL::integer,
            NULL::integer,
            1,
            1,
            ARRAY[]::text[],
            false,
            NULL::jsonb,
            15,
            '{
              "source_url": "https://www.metlife.com.bd/solutions/group-insurance/employee-benefits/group-death-disability/",
              "public_benefits": [
                "group life coverage",
                "total permanent disability protection",
                "employer-funded staff benefit"
              ],
              "public_coverage_range_bdt": "100000 to 500000"
            }'::jsonb,
            NULL::jsonb,
            DATE '2026-03-27',
            NULL::date,
            '{}'::jsonb
        ),
        (
            'a4f2f420-42d2-4c3a-92a0-100000000017'::uuid,
            'CHARTERED',
            'LIF-001',
            'CHARTERED-CEP1',
            'Endowment Plan-1',
            'INACTIVE',
            '{"amount":5000000,"currency":"BDT","decimal_amount":50000}'::jsonb,
            NULL::jsonb,
            NULL::jsonb,
            NULL::jsonb,
            18,
            55,
            75,
            10,
            35,
            ARRAY[]::text[],
            false,
            NULL::jsonb,
            15,
            '{
              "source_url": "https://charteredlifebd.org/prospectus.aspx",
              "public_benefits": [
                "minimum sum assured BDT 50000",
                "term from 10 to 35 years",
                "entry age from 18 to 55 years"
              ]
            }'::jsonb,
            NULL::jsonb,
            DATE '2026-03-27',
            NULL::date,
            '{}'::jsonb
        ),
        (
            'a4f2f420-42d2-4c3a-92a0-100000000018'::uuid,
            'CHARTERED',
            'LIF-001',
            'CHARTERED-CMSP',
            'Monthly Savings Plan',
            'INACTIVE',
            '{"amount":10000000,"currency":"BDT","decimal_amount":100000}'::jsonb,
            '{"amount":1000000000,"currency":"BDT","decimal_amount":10000000}'::jsonb,
            '{"amount":100000,"currency":"BDT","decimal_amount":1000}'::jsonb,
            NULL::jsonb,
            20,
            55,
            75,
            10,
            25,
            ARRAY['MONTHLY']::text[],
            false,
            NULL::jsonb,
            15,
            '{
              "source_url": "https://charteredlifebd.org/prospectus.aspx",
              "public_benefits": [
                "minimum monthly premium BDT 1000",
                "sum assured from BDT 100000 to BDT 10000000",
                "term from 10 to 25 years"
              ]
            }'::jsonb,
            NULL::jsonb,
            DATE '2026-03-27',
            NULL::date,
            '{}'::jsonb
        ),
        (
            'a4f2f420-42d2-4c3a-92a0-100000000019'::uuid,
            'CHARTERED',
            'LIF-001',
            'CHARTERED-CSP',
            'Single Premium Plan',
            'INACTIVE',
            '{"amount":2000000,"currency":"BDT","decimal_amount":20000}'::jsonb,
            NULL::jsonb,
            '{"amount":2900000,"currency":"BDT","decimal_amount":29000}'::jsonb,
            NULL::jsonb,
            18,
            60,
            NULL::integer,
            6,
            16,
            ARRAY['SINGLE']::text[],
            false,
            NULL::jsonb,
            15,
            '{
              "source_url": "https://charteredlifebd.org/prospectus.aspx",
              "public_benefits": [
                "minimum premium BDT 29000",
                "minimum sum assured BDT 20000",
                "policy term from 6 to 16 years"
              ]
            }'::jsonb,
            NULL::jsonb,
            DATE '2026-03-27',
            NULL::date,
            '{}'::jsonb
        ),
        (
            'a4f2f420-42d2-4c3a-92a0-100000000020'::uuid,
            'CHARTERED',
            'LIF-001',
            'CHARTERED-NIR3',
            'Nirapotta - 3 Years',
            'INACTIVE',
            '{"amount":5000000,"currency":"BDT","decimal_amount":50000}'::jsonb,
            '{"amount":15000000,"currency":"BDT","decimal_amount":150000}'::jsonb,
            '{"amount":300000,"currency":"BDT","decimal_amount":3000}'::jsonb,
            '{"amount":300000,"currency":"BDT","decimal_amount":3000}'::jsonb,
            NULL::integer,
            NULL::integer,
            NULL::integer,
            3,
            3,
            ARRAY[]::text[],
            false,
            NULL::jsonb,
            15,
            '{
              "source_url": "https://charteredlifebd.org/prospectus.aspx",
              "public_benefits": [
                "natural death benefit BDT 50000",
                "accidental death benefit BDT 150000",
                "hospitalization benefit BDT 30000"
              ],
              "public_premium_bdt": 3000
            }'::jsonb,
            NULL::jsonb,
            DATE '2026-03-27',
            NULL::date,
            '{}'::jsonb
        ),
        (
            'a4f2f420-42d2-4c3a-92a0-100000000021'::uuid,
            'CHARTERED',
            'LIF-001',
            'CHARTERED-NIR5',
            'Nirapotta - 5 Years',
            'INACTIVE',
            '{"amount":5000000,"currency":"BDT","decimal_amount":50000}'::jsonb,
            '{"amount":15000000,"currency":"BDT","decimal_amount":150000}'::jsonb,
            '{"amount":500000,"currency":"BDT","decimal_amount":5000}'::jsonb,
            '{"amount":500000,"currency":"BDT","decimal_amount":5000}'::jsonb,
            NULL::integer,
            NULL::integer,
            NULL::integer,
            5,
            5,
            ARRAY[]::text[],
            false,
            NULL::jsonb,
            15,
            '{
              "source_url": "https://charteredlifebd.org/prospectus.aspx",
              "public_benefits": [
                "natural death benefit BDT 50000",
                "accidental death benefit BDT 150000",
                "hospitalization benefit BDT 30000"
              ],
              "public_premium_bdt": 5000
            }'::jsonb,
            NULL::jsonb,
            DATE '2026-03-27',
            NULL::date,
            '{}'::jsonb
        ),
        (
            'a4f2f420-42d2-4c3a-92a0-100000000022'::uuid,
            'CHARTERED',
            'LIF-001',
            'CHARTERED-GROUPLIFE',
            'Group Life Insurance',
            'INACTIVE',
            NULL::jsonb,
            NULL::jsonb,
            NULL::jsonb,
            NULL::jsonb,
            NULL::integer,
            NULL::integer,
            NULL::integer,
            1,
            1,
            ARRAY[]::text[],
            false,
            NULL::jsonb,
            15,
            '{
              "source_url": "https://charteredlifebd.org/Corporateinformatiom.aspx",
              "product_family": "group life"
            }'::jsonb,
            NULL::jsonb,
            DATE '2026-03-27',
            NULL::date,
            '{}'::jsonb
        ),
        (
            'a4f2f420-42d2-4c3a-92a0-100000000023'::uuid,
            'CHARTERED',
            'HLT-001',
            'CHARTERED-GROUPHEALTH',
            'Group Health Insurance',
            'INACTIVE',
            NULL::jsonb,
            NULL::jsonb,
            NULL::jsonb,
            NULL::jsonb,
            NULL::integer,
            NULL::integer,
            NULL::integer,
            1,
            1,
            ARRAY[]::text[],
            false,
            NULL::jsonb,
            15,
            '{
              "source_url": "https://charteredlifebd.org/Corporateinformatiom.aspx",
              "product_family": "group health"
            }'::jsonb,
            NULL::jsonb,
            DATE '2026-03-27',
            NULL::date,
            '{}'::jsonb
        )
    ) AS v(
        seed_product_id,
        insurer_code,
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
      ON i.code = v.insurer_code
     AND p.product_code = v.base_product_code
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
