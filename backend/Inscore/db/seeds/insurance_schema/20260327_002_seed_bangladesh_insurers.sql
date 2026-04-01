-- ============================================================
-- Seed: Bangladesh Insurers
-- Purpose: Seed core insurer master data used by insurer proto/repository
-- Scope: Pragati Insurance PLC, Green Delta Insurance PLC,
--        MetLife Bangladesh, Chartered Life Insurance PLC
-- Sources (official/public, checked 2026-03-27):
--   Pragati Insurance PLC
--     https://pragatiinsurance.com/company-profile/
--     https://pragatiinsurance.com/contact/
--     https://pragatiinsurance.com/wp-content/uploads/2025/07/PIL-Annual-Report-2024.pdf
--   Green Delta Insurance PLC
--     https://green-delta.com/contact-us/
--     https://green-delta.com/wp-content/uploads/2025/04/Financial-Highlight-2024-1.pdf
--     https://green-delta.com/wp-content/uploads/2025/03/Annual-Report-2024.pdf
--   MetLife Bangladesh
--     https://www.metlife.com.bd/contact-us/
--     https://www.metlife.com.bd/content/dam/metlifecom/bd/PDFs/others/fact-sheet-en-2025.pdf
--   Chartered Life Insurance PLC
--     https://charteredlifebd.org/Corporateinformatiom.aspx
-- Notes:
--   - Publicly unavailable regulatory IDs (trade license, TIN, IDRA license no.)
--     are left NULL.
--   - Paid-up capital values use BDT in smallest unit (paisa) inside JSONB,
--     with decimal_amount stored in whole BDT for display consistency.
--   - MetLife Bangladesh public pages expose rating/contact information but not a
--     local paid-up capital figure on the pages reviewed, so paid_up_capital is {}.
-- ============================================================

WITH insurer_seed (
    insurer_id,
    code,
    name,
    name_bn,
    type,
    status,
    trade_license_number,
    tin_number,
    idra_license_number,
    idra_license_expiry,
    contact_info,
    registered_address,
    head_office_address,
    logo_url,
    website_url,
    financial_rating,
    paid_up_capital,
    audit_info
) AS (
    VALUES
    (
        'c4d7d5b8-f43f-4e8c-8b41-c0d8f7146001'::uuid,
        'PRAGATI',
        'Pragati Insurance PLC',
        NULL::date,
        'NON_LIFE',
        'ACTIVE',
        NULL::date,
        NULL::date,
        NULL::date,
        NULL,
        '{
          "mobile_number": "09613115511",
          "email": "info@pragatiinsurance.com",
          "landline": "+880-2-55012680-2"
        }'::jsonb,
        '{
          "address_line1": "Pragati Insurance Bhaban, 20-21 Kawran Bazar",
          "city": "Dhaka",
          "district": "Dhaka",
          "division": "Dhaka",
          "postal_code": "1215",
          "country": "Bangladesh"
        }'::jsonb,
        '{
          "address_line1": "Pragati Insurance Bhaban, 20-21 Kawran Bazar",
          "city": "Dhaka",
          "district": "Dhaka",
          "division": "Dhaka",
          "postal_code": "1215",
          "country": "Bangladesh"
        }'::jsonb,
        NULL,
        'https://pragatiinsurance.com/',
        'AAA',
        '{
          "amount": 73691000000,
          "currency": "BDT",
          "decimal_amount": 736910000
        }'::jsonb,
        '{}'::jsonb
    ),
    (
        'c4d7d5b8-f43f-4e8c-8b41-c0d8f7146002'::uuid,
        'GREENDELTA',
        'Green Delta Insurance PLC',
        NULL,
        'NON_LIFE',
        'ACTIVE',
        NULL,
        NULL,
        NULL,
        NULL,
        '{
          "mobile_number": "+8801730031888",
          "email": "info@green-delta.com",
          "landline": "+8809613444888"
        }'::jsonb,
        '{
          "address_line1": "Green Delta Aims Tower (5th Floor), 51-52 Mohakhali",
          "city": "Dhaka",
          "district": "Dhaka",
          "division": "Dhaka",
          "postal_code": "1212",
          "country": "Bangladesh"
        }'::jsonb,
        '{
          "address_line1": "Green Delta Aims Tower (5th Floor), 51-52 Mohakhali",
          "city": "Dhaka",
          "district": "Dhaka",
          "division": "Dhaka",
          "postal_code": "1212",
          "country": "Bangladesh"
        }'::jsonb,
        NULL,
        'https://green-delta.com/',
        'AAA',
        '{
          "amount": 100200000000,
          "currency": "BDT",
          "decimal_amount": 1002000000
        }'::jsonb,
        '{}'::jsonb
    ),
    (
        'c4d7d5b8-f43f-4e8c-8b41-c0d8f7146003'::uuid,
        'METLIFE',
        'MetLife Bangladesh',
        NULL,
        'LIFE',
        'ACTIVE',
        NULL,
        NULL,
        NULL,
        NULL,
        '{
          "mobile_number": "09666716344",
          "email": "customer.services@metlife.com.bd",
          "landline": "16344"
        }'::jsonb,
        '{
          "address_line1": "MetLife Building, 18-20 Motijheel C.A.",
          "address_line2": "P.O. Box 9",
          "city": "Dhaka",
          "district": "Dhaka",
          "division": "Dhaka",
          "postal_code": "1000",
          "country": "Bangladesh"
        }'::jsonb,
        '{
          "address_line1": "MetLife Building, 18-20 Motijheel C.A.",
          "address_line2": "P.O. Box 9",
          "city": "Dhaka",
          "district": "Dhaka",
          "division": "Dhaka",
          "postal_code": "1000",
          "country": "Bangladesh"
        }'::jsonb,
        NULL,
        'https://www.metlife.com.bd/',
        'AAA',
        '{}'::jsonb,
        '{}'::jsonb
    ),
    (
        'c4d7d5b8-f43f-4e8c-8b41-c0d8f7146004'::uuid,
        'CHARTERED',
        'Chartered Life Insurance PLC',
        NULL,
        'LIFE',
        'ACTIVE',
        NULL,
        NULL,
        NULL,
        NULL,
        '{
          "mobile_number": "+8801777770990",
          "email": "mail@charteredlifebd.com",
          "landline": "+8802-55128956-7"
        }'::jsonb,
        '{
          "address_line1": "Islam Tower (8th Floor), 464/H D.I.T Road, West Rampura",
          "city": "Dhaka",
          "district": "Dhaka",
          "division": "Dhaka",
          "postal_code": "1219",
          "country": "Bangladesh"
        }'::jsonb,
        '{
          "address_line1": "Islam Tower (8th Floor), 464/H D.I.T Road, West Rampura",
          "city": "Dhaka",
          "district": "Dhaka",
          "division": "Dhaka",
          "postal_code": "1219",
          "country": "Bangladesh"
        }'::jsonb,
        NULL,
        'https://charteredlifebd.org/',
        'AA-',
        '{
          "amount": 37500000000,
          "currency": "BDT",
          "decimal_amount": 375000000
        }'::jsonb,
        '{}'::jsonb
    )
)
INSERT INTO insurance_schema.insurers (
    insurer_id,
    code,
    name,
    name_bn,
    type,
    status,
    trade_license_number,
    tin_number,
    idra_license_number,
    idra_license_expiry,
    contact_info,
    registered_address,
    head_office_address,
    logo_url,
    website_url,
    financial_rating,
    paid_up_capital,
    audit_info
)
SELECT
    insurer_id,
    code,
    name,
    name_bn,
    type,
    status,
    trade_license_number,
    tin_number,
    idra_license_number,
    idra_license_expiry::date,
    contact_info,
    registered_address,
    head_office_address,
    logo_url,
    website_url,
    financial_rating,
    paid_up_capital,
    audit_info
FROM insurer_seed
ON CONFLICT (code) DO UPDATE
SET
    name = EXCLUDED.name,
    name_bn = EXCLUDED.name_bn,
    type = EXCLUDED.type,
    status = EXCLUDED.status,
    trade_license_number = EXCLUDED.trade_license_number,
    tin_number = EXCLUDED.tin_number,
    idra_license_number = EXCLUDED.idra_license_number,
    idra_license_expiry = EXCLUDED.idra_license_expiry,
    contact_info = EXCLUDED.contact_info,
    registered_address = EXCLUDED.registered_address,
    head_office_address = EXCLUDED.head_office_address,
    logo_url = EXCLUDED.logo_url,
    website_url = EXCLUDED.website_url,
    financial_rating = EXCLUDED.financial_rating,
    paid_up_capital = EXCLUDED.paid_up_capital;
