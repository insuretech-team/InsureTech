-- ============================================================
-- Seed: Bangladesh Insurer Configs
-- Purpose: Seed baseline insurer integration config for onboarding flows
-- Depends on: 20260327_002_seed_bangladesh_insurers.sql
-- Notes:
--   - Public carrier API endpoints, auth credentials, and webhook contracts are
--     not publicly disclosed, so these rows intentionally seed operational
--     defaults for manual onboarding and later insurer-specific integration.
--   - `auth_type` is left as AUTHENTICATION_TYPE_UNSPECIFIED until a real
--     carrier integration contract is configured.
--   - `business_model` and `payment_terms` are JSONB defaults tuned by insurer
--     category and currently available products in this repo.
-- ============================================================

WITH insurer_lookup AS (
    SELECT code, insurer_id
    FROM insurance_schema.insurers
    WHERE code IN ('PRAGATI', 'GREENDELTA', 'METLIFE', 'CHARTERED')
),
config_seed AS (
    SELECT
        v.config_id,
        i.insurer_id,
        v.api_base_url,
        v.api_version,
        v.auth_type,
        v.auth_credentials,
        v.webhook_url,
        v.webhook_secret,
        v.business_model,
        v.auto_underwriting_enabled,
        v.underwriting_threshold,
        v.real_time_claim_notification,
        v.claim_settlement_days,
        v.payment_terms,
        v.audit_info
    FROM insurer_lookup i
    JOIN (
        VALUES
        (
            'PRAGATI',
            '18c2f763-370a-49b4-a6d8-966b17c63001'::uuid,
            NULL::text,
            NULL::text,
            'UNSPECIFIED',
            NULL::text,
            NULL::text,
            NULL::text,
            '{
              "integration_mode": "manual",
              "revenue_model": "commission",
              "portfolio_type": "non_life",
              "product_focus": ["personal_accident", "travel", "pet"],
              "settlement_basis": "monthly_statement"
            }'::jsonb,
            false,
            0,
            false,
            15,
            '{
              "currency": "BDT",
              "collection_mode": "platform_collected",
              "settlement_cycle": "monthly",
              "settlement_day_of_month": 7,
              "refund_window_days": 7,
              "supported_payment_methods": ["bank_transfer", "card", "mfs"]
            }'::jsonb,
            '{}'::jsonb
        ),
        (
            'GREENDELTA',
            '18c2f763-370a-49b4-a6d8-966b17c63002'::uuid,
            NULL::text,
            NULL::text,
            'UNSPECIFIED',
            NULL::text,
            NULL::text,
            NULL::text,
            '{
              "integration_mode": "manual",
              "revenue_model": "commission",
              "portfolio_type": "non_life",
              "product_focus": ["motor", "travel", "health", "property"],
              "settlement_basis": "monthly_statement"
            }'::jsonb,
            false,
            0,
            false,
            15,
            '{
              "currency": "BDT",
              "collection_mode": "platform_collected",
              "settlement_cycle": "monthly",
              "settlement_day_of_month": 7,
              "refund_window_days": 7,
              "supported_payment_methods": ["bank_transfer", "card", "mfs"]
            }'::jsonb,
            '{}'::jsonb
        ),
        (
            'METLIFE',
            '18c2f763-370a-49b4-a6d8-966b17c63003'::uuid,
            NULL::text,
            NULL::text,
            'UNSPECIFIED',
            NULL::text,
            NULL::text,
            NULL::text,
            '{
              "integration_mode": "manual",
              "revenue_model": "commission",
              "portfolio_type": "life",
              "product_focus": ["term_life", "group_life", "savings"],
              "settlement_basis": "monthly_statement"
            }'::jsonb,
            false,
            0,
            false,
            10,
            '{
              "currency": "BDT",
              "collection_mode": "platform_collected",
              "settlement_cycle": "monthly",
              "settlement_day_of_month": 7,
              "refund_window_days": 10,
              "supported_payment_methods": ["bank_transfer", "card", "mfs", "direct_debit"]
            }'::jsonb,
            '{}'::jsonb
        ),
        (
            'CHARTERED',
            '18c2f763-370a-49b4-a6d8-966b17c63004'::uuid,
            NULL::text,
            NULL::text,
            'UNSPECIFIED',
            NULL::text,
            NULL::text,
            NULL::text,
            '{
              "integration_mode": "manual",
              "revenue_model": "commission",
              "portfolio_type": "life",
              "product_focus": ["term_life", "group_life", "microinsurance"],
              "settlement_basis": "monthly_statement"
            }'::jsonb,
            false,
            0,
            false,
            15,
            '{
              "currency": "BDT",
              "collection_mode": "platform_collected",
              "settlement_cycle": "monthly",
              "settlement_day_of_month": 7,
              "refund_window_days": 10,
              "supported_payment_methods": ["bank_transfer", "card", "mfs"]
            }'::jsonb,
            '{}'::jsonb
        )
    ) AS v(
        code,
        config_id,
        api_base_url,
        api_version,
        auth_type,
        auth_credentials,
        webhook_url,
        webhook_secret,
        business_model,
        auto_underwriting_enabled,
        underwriting_threshold,
        real_time_claim_notification,
        claim_settlement_days,
        payment_terms,
        audit_info
    )
        ON i.code = v.code
)
INSERT INTO insurance_schema.insurer_configs (
    config_id,
    insurer_id,
    api_base_url,
    api_version,
    auth_type,
    auth_credentials,
    webhook_url,
    webhook_secret,
    business_model,
    auto_underwriting_enabled,
    underwriting_threshold,
    real_time_claim_notification,
    claim_settlement_days,
    payment_terms,
    audit_info
)
SELECT
    config_id,
    insurer_id,
    api_base_url,
    api_version,
    auth_type,
    auth_credentials,
    webhook_url,
    webhook_secret,
    business_model,
    auto_underwriting_enabled,
    underwriting_threshold,
    real_time_claim_notification,
    claim_settlement_days,
    payment_terms,
    audit_info
FROM config_seed
ON CONFLICT (insurer_id) DO UPDATE
SET
    api_base_url = EXCLUDED.api_base_url,
    api_version = EXCLUDED.api_version,
    auth_type = EXCLUDED.auth_type,
    auth_credentials = EXCLUDED.auth_credentials,
    webhook_url = EXCLUDED.webhook_url,
    webhook_secret = EXCLUDED.webhook_secret,
    business_model = EXCLUDED.business_model,
    auto_underwriting_enabled = EXCLUDED.auto_underwriting_enabled,
    underwriting_threshold = EXCLUDED.underwriting_threshold,
    real_time_claim_notification = EXCLUDED.real_time_claim_notification,
    claim_settlement_days = EXCLUDED.claim_settlement_days,
    payment_terms = EXCLUDED.payment_terms;
