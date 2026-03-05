BEGIN;

-- Seed departments (proto-owned table: b2b_schema.departments)
INSERT INTO b2b_schema.departments (
  department_id,
  name,
  business_id,
  employee_no,
  total_premium
)
VALUES
  ('44444444-4444-4444-4444-444444444001', 'HR',               '22222222-2222-2222-2222-222222222001', 2, '{"amount":100000,"currency":"BDT","decimal_amount":1000}'::jsonb),
  ('44444444-4444-4444-4444-444444444002', 'Finance',          '22222222-2222-2222-2222-222222222001', 2, '{"amount":109000,"currency":"BDT","decimal_amount":1090}'::jsonb),
  ('44444444-4444-4444-4444-444444444003', 'IT',               '22222222-2222-2222-2222-222222222001', 2, '{"amount":147000,"currency":"BDT","decimal_amount":1470}'::jsonb),
  ('44444444-4444-4444-4444-444444444004', 'Marketing',        '22222222-2222-2222-2222-222222222001', 1, '{"amount":46000,"currency":"BDT","decimal_amount":460}'::jsonb),
  ('44444444-4444-4444-4444-444444444005', 'Operations',       '22222222-2222-2222-2222-222222222001', 2, '{"amount":84500,"currency":"BDT","decimal_amount":845}'::jsonb),
  ('44444444-4444-4444-4444-444444444006', 'Customer Support', '22222222-2222-2222-2222-222222222001', 1, '{"amount":41000,"currency":"BDT","decimal_amount":410}'::jsonb),
  ('44444444-4444-4444-4444-444444444007', 'Sales',            '22222222-2222-2222-2222-222222222001', 1, '{"amount":39000,"currency":"BDT","decimal_amount":390}'::jsonb),
  ('44444444-4444-4444-4444-444444444008', 'Legal',            '22222222-2222-2222-2222-222222222001', 1, '{"amount":57500,"currency":"BDT","decimal_amount":575}'::jsonb)
ON CONFLICT (department_id) DO UPDATE
SET
  name = EXCLUDED.name,
  business_id = EXCLUDED.business_id,
  employee_no = EXCLUDED.employee_no,
  total_premium = EXCLUDED.total_premium,
  updated_at = NOW();

-- Seed employees (proto-owned table: b2b_schema.employees)
INSERT INTO b2b_schema.employees (
  employee_uuid,
  name,
  employee_id,
  department_id,
  business_id,
  insurance_category,
  assigned_plan_id,
  coverage_amount,
  premium_amount,
  status,
  number_of_dependent
)
VALUES
  ('66666666-6666-6666-6666-666666666001','John Doe','LPBL10032','44444444-4444-4444-4444-444444444001','22222222-2222-2222-2222-222222222001','INSURANCE_TYPE_HEALTH','55555555-5555-5555-5555-555555555001','{"amount":10000000,"currency":"BDT","decimal_amount":100000}'::jsonb,'{"amount":50000,"currency":"BDT","decimal_amount":500}'::jsonb,'ACTIVE',2),
  ('66666666-6666-6666-6666-666666666002','Jane Smith','LPBL10033','44444444-4444-4444-4444-444444444002','22222222-2222-2222-2222-222222222001','INSURANCE_TYPE_HEALTH','55555555-5555-5555-5555-555555555001','{"amount":10000000,"currency":"BDT","decimal_amount":100000}'::jsonb,'{"amount":50000,"currency":"BDT","decimal_amount":500}'::jsonb,'INACTIVE',1),
  ('66666666-6666-6666-6666-666666666003','Bob Johnson','LPBL10034','44444444-4444-4444-4444-444444444003','22222222-2222-2222-2222-222222222001','INSURANCE_TYPE_LIFE','55555555-5555-5555-5555-555555555004','{"amount":25000000,"currency":"BDT","decimal_amount":250000}'::jsonb,'{"amount":85000,"currency":"BDT","decimal_amount":850}'::jsonb,'ACTIVE',3),
  ('66666666-6666-6666-6666-666666666004','Alice Williams','LPBL10035','44444444-4444-4444-4444-444444444004','22222222-2222-2222-2222-222222222001','INSURANCE_TYPE_HEALTH','55555555-5555-5555-5555-555555555001','{"amount":9000000,"currency":"BDT","decimal_amount":90000}'::jsonb,'{"amount":46000,"currency":"BDT","decimal_amount":460}'::jsonb,'ACTIVE',2),
  ('66666666-6666-6666-6666-666666666005','Rafiul Karim','LPBL10036','44444444-4444-4444-4444-444444444005','22222222-2222-2222-2222-222222222001','INSURANCE_TYPE_HEALTH','55555555-5555-5555-5555-555555555002','{"amount":8500000,"currency":"BDT","decimal_amount":85000}'::jsonb,'{"amount":43000,"currency":"BDT","decimal_amount":430}'::jsonb,'ACTIVE',1),
  ('66666666-6666-6666-6666-666666666006','Nusrat Jahan','LPBL10037','44444444-4444-4444-4444-444444444006','22222222-2222-2222-2222-222222222001','INSURANCE_TYPE_HEALTH','55555555-5555-5555-5555-555555555002','{"amount":8000000,"currency":"BDT","decimal_amount":80000}'::jsonb,'{"amount":41000,"currency":"BDT","decimal_amount":410}'::jsonb,'ACTIVE',4),
  ('66666666-6666-6666-6666-666666666007','Sabbir Hossain','LPBL10038','44444444-4444-4444-4444-444444444002','22222222-2222-2222-2222-222222222001','INSURANCE_TYPE_HEALTH','55555555-5555-5555-5555-555555555003','{"amount":14000000,"currency":"BDT","decimal_amount":140000}'::jsonb,'{"amount":59000,"currency":"BDT","decimal_amount":590}'::jsonb,'ACTIVE',2),
  ('66666666-6666-6666-6666-666666666008','Mahi Rahman','LPBL10039','44444444-4444-4444-4444-444444444001','22222222-2222-2222-2222-222222222001','INSURANCE_TYPE_HEALTH','55555555-5555-5555-5555-555555555001','{"amount":10000000,"currency":"BDT","decimal_amount":100000}'::jsonb,'{"amount":50000,"currency":"BDT","decimal_amount":500}'::jsonb,'INACTIVE',0),
  ('66666666-6666-6666-6666-666666666009','Tanvir Ahmed','LPBL10040','44444444-4444-4444-4444-444444444003','22222222-2222-2222-2222-222222222001','INSURANCE_TYPE_LIFE','55555555-5555-5555-5555-555555555004','{"amount":15000000,"currency":"BDT","decimal_amount":150000}'::jsonb,'{"amount":62000,"currency":"BDT","decimal_amount":620}'::jsonb,'ACTIVE',5),
  ('66666666-6666-6666-6666-666666666010','Farzana Akter','LPBL10041','44444444-4444-4444-4444-444444444007','22222222-2222-2222-2222-222222222001','INSURANCE_TYPE_HEALTH','55555555-5555-5555-5555-555555555002','{"amount":7500000,"currency":"BDT","decimal_amount":75000}'::jsonb,'{"amount":39000,"currency":"BDT","decimal_amount":390}'::jsonb,'ACTIVE',1),
  ('66666666-6666-6666-6666-666666666011','Imran Kabir','LPBL10042','44444444-4444-4444-4444-444444444005','22222222-2222-2222-2222-222222222001','INSURANCE_TYPE_HEALTH','55555555-5555-5555-5555-555555555002','{"amount":8200000,"currency":"BDT","decimal_amount":82000}'::jsonb,'{"amount":41500,"currency":"BDT","decimal_amount":415}'::jsonb,'ACTIVE',2),
  ('66666666-6666-6666-6666-666666666012','Sharmin Nahar','LPBL10043','44444444-4444-4444-4444-444444444008','22222222-2222-2222-2222-222222222001','INSURANCE_TYPE_HEALTH','55555555-5555-5555-5555-555555555003','{"amount":13000000,"currency":"BDT","decimal_amount":130000}'::jsonb,'{"amount":57500,"currency":"BDT","decimal_amount":575}'::jsonb,'ACTIVE',3)
ON CONFLICT (employee_uuid) DO UPDATE
SET
  name = EXCLUDED.name,
  employee_id = EXCLUDED.employee_id,
  department_id = EXCLUDED.department_id,
  business_id = EXCLUDED.business_id,
  insurance_category = EXCLUDED.insurance_category,
  assigned_plan_id = EXCLUDED.assigned_plan_id,
  coverage_amount = EXCLUDED.coverage_amount,
  premium_amount = EXCLUDED.premium_amount,
  status = EXCLUDED.status,
  number_of_dependent = EXCLUDED.number_of_dependent,
  updated_at = NOW();

COMMIT;
