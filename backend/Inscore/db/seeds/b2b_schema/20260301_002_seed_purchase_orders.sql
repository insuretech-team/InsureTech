BEGIN;

INSERT INTO b2b_schema.purchase_orders (
  purchase_order_id,
  purchase_order_number,
  business_id,
  department_id,
  product_id,
  plan_id,
  insurance_category,
  employee_count,
  number_of_dependents,
  coverage_amount,
  estimated_premium,
  status,
  requested_by,
  notes
)
VALUES
  (
    '77777777-7777-7777-7777-777777777001',
    'PO-20260227-A1B2C3D4',
    '22222222-2222-2222-2222-222222222001',
    '44444444-4444-4444-4444-444444444001',
    'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa1',
    '55555555-5555-5555-5555-555555555001',
    'INSURANCE_TYPE_HEALTH',
    32,
    18,
    '{"amount":320000000,"currency":"BDT","decimal_amount":3200000}'::jsonb,
    '{"amount":1600000,"currency":"BDT","decimal_amount":16000}'::jsonb,
    'APPROVED',
    'ccca65ad-ae2c-4d42-8ccc-2122db78d617',
    'Annual employee health coverage expansion for the HR team.'
  ),
  (
    '77777777-7777-7777-7777-777777777002',
    'PO-20260227-E5F6G7H8',
    '22222222-2222-2222-2222-222222222001',
    '44444444-4444-4444-4444-444444444003',
    'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbb2',
    '55555555-5555-5555-5555-555555555004',
    'INSURANCE_TYPE_LIFE',
    20,
    26,
    '{"amount":500000000,"currency":"BDT","decimal_amount":5000000}'::jsonb,
    '{"amount":1700000,"currency":"BDT","decimal_amount":17000}'::jsonb,
    'SUBMITTED',
    'ccca65ad-ae2c-4d42-8ccc-2122db78d617',
    'Life cover upgrade for senior engineering employees.'
  ),
  (
    '77777777-7777-7777-7777-777777777003',
    'PO-20260227-I9J0K1L2',
    '22222222-2222-2222-2222-222222222001',
    '44444444-4444-4444-4444-444444444005',
    'cccccccc-cccc-cccc-cccc-ccccccccccc3',
    '55555555-5555-5555-5555-555555555002',
    'INSURANCE_TYPE_HEALTH',
    14,
    9,
    '{"amount":119000000,"currency":"BDT","decimal_amount":1190000}'::jsonb,
    '{"amount":602000,"currency":"BDT","decimal_amount":6020}'::jsonb,
    'FULFILLED',
    'ccca65ad-ae2c-4d42-8ccc-2122db78d617',
    'Operations team onboarding order completed and activated.'
  )
ON CONFLICT (purchase_order_id) DO UPDATE
SET
  purchase_order_number = EXCLUDED.purchase_order_number,
  business_id = EXCLUDED.business_id,
  department_id = EXCLUDED.department_id,
  product_id = EXCLUDED.product_id,
  plan_id = EXCLUDED.plan_id,
  insurance_category = EXCLUDED.insurance_category,
  employee_count = EXCLUDED.employee_count,
  number_of_dependents = EXCLUDED.number_of_dependents,
  coverage_amount = EXCLUDED.coverage_amount,
  estimated_premium = EXCLUDED.estimated_premium,
  status = EXCLUDED.status,
  requested_by = EXCLUDED.requested_by,
  notes = EXCLUDED.notes,
  updated_at = NOW();

COMMIT;
