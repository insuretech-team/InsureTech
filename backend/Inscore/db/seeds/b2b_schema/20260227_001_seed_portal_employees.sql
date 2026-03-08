BEGIN;

-- Seed the default organisation referenced by all default B2B employees.
INSERT INTO b2b_schema.organisations (
  organisation_id,
  tenant_id,
  name,
  code,
  industry,
  contact_email,
  contact_phone,
  address,
  status,
  total_employees
)
VALUES (
  '22222222-2222-2222-2222-222222222001',
  '00000000-0000-0000-0000-000000000001',
  'Default B2B Organisation',
  'DEFAULT-B2B',
  'Insurance',
  'b2b-admin@lifeplus.local',
  '+8801000000000',
  'Dhaka, Bangladesh',
  'ORGANISATION_STATUS_ACTIVE',
  0
)
ON CONFLICT (organisation_id) DO UPDATE
SET
  tenant_id = EXCLUDED.tenant_id,
  name = EXCLUDED.name,
  code = EXCLUDED.code,
  industry = EXCLUDED.industry,
  contact_email = EXCLUDED.contact_email,
  contact_phone = EXCLUDED.contact_phone,
  address = EXCLUDED.address,
  status = EXCLUDED.status,
  updated_at = NOW();

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

-- Link the first seeded B2B admin user to the default organisation when available.
DO $$
DECLARE
  default_admin_user_id UUID;
  default_admin_email TEXT;
  default_admin_mobile TEXT;
  default_admin_role_id UUID;
BEGIN
  SELECT u.user_id, COALESCE(u.email, 'b2b-admin@lifeplus.local'), COALESCE(u.mobile_number, '+8801000000000')
  INTO default_admin_user_id, default_admin_email, default_admin_mobile
  FROM authn_schema.users u
  WHERE u.deleted_at IS NULL
    AND u.user_type = 'USER_TYPE_B2B_ORG_ADMIN'
  ORDER BY u.created_at ASC
  LIMIT 1;

  IF default_admin_user_id IS NULL THEN
    RAISE NOTICE 'No USER_TYPE_B2B_ORG_ADMIN found; skipping default org admin membership seed';
    RETURN;
  END IF;

  IF EXISTS (
    SELECT 1
    FROM b2b_schema.org_members
    WHERE organisation_id = '22222222-2222-2222-2222-222222222001'
      AND user_id = default_admin_user_id
      AND deleted_at IS NULL
  ) THEN
    UPDATE b2b_schema.org_members
    SET
      role = 'ORG_MEMBER_ROLE_BUSINESS_ADMIN',
      status = 'ORG_MEMBER_STATUS_ACTIVE',
      deleted_at = NULL,
      updated_at = NOW()
    WHERE organisation_id = '22222222-2222-2222-2222-222222222001'
      AND user_id = default_admin_user_id;
  ELSE
    INSERT INTO b2b_schema.org_members (
      member_id,
      organisation_id,
      user_id,
      role,
      status,
      joined_at
    )
    VALUES (
      '33333333-3333-3333-3333-333333333001',
      '22222222-2222-2222-2222-222222222001',
      default_admin_user_id,
      'ORG_MEMBER_ROLE_BUSINESS_ADMIN',
      'ORG_MEMBER_STATUS_ACTIVE',
      NOW()
    );
  END IF;

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
    number_of_dependent,
    email,
    mobile_number,
    date_of_joining,
    user_id
  )
  VALUES (
    '66666666-6666-6666-6666-666666666999',
    'Default B2B Admin',
    'B2BADMIN001',
    '44444444-4444-4444-4444-444444444001',
    '22222222-2222-2222-2222-222222222001',
    'INSURANCE_TYPE_HEALTH',
    '55555555-5555-5555-5555-555555555001',
    '{"amount":10000000,"currency":"BDT","decimal_amount":100000}'::jsonb,
    '{"amount":50000,"currency":"BDT","decimal_amount":500}'::jsonb,
    'EMPLOYEE_STATUS_ACTIVE',
    0,
    default_admin_email,
    default_admin_mobile,
    CURRENT_DATE,
    default_admin_user_id
  )
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
    email = EXCLUDED.email,
    mobile_number = EXCLUDED.mobile_number,
    date_of_joining = EXCLUDED.date_of_joining,
    user_id = EXCLUDED.user_id,
    updated_at = NOW();

  SELECT r.role_id
  INTO default_admin_role_id
  FROM authz_schema.roles r
  WHERE r.deleted_at IS NULL
    AND r.portal = 'PORTAL_B2B'
    AND r.name = 'b2b_org_admin'
  ORDER BY r.created_at ASC
  LIMIT 1;

  IF default_admin_role_id IS NOT NULL THEN
    INSERT INTO authz_schema.user_roles (
      user_role_id,
      user_id,
      role_id,
      domain,
      assigned_by,
      assigned_at
    )
    VALUES (
      '77777777-7777-7777-7777-777777777001',
      default_admin_user_id,
      default_admin_role_id,
      'b2b:22222222-2222-2222-2222-222222222001',
      default_admin_user_id,
      NOW()
    )
    ON CONFLICT ON CONSTRAINT uq_user_roles_user_role_domain DO UPDATE
    SET assigned_at = EXCLUDED.assigned_at;

    INSERT INTO authz_schema.casbin_rules (ptype, v0, v1, v2, v3, v4, v5)
    VALUES (
      'g',
      'user:' || default_admin_user_id::TEXT,
      'b2b:22222222-2222-2222-2222-222222222001',
      'role:b2b_org_admin',
      '',
      '',
      ''
    )
    ON CONFLICT ON CONSTRAINT uq_casbin_rules_tuple DO NOTHING;
  END IF;
END $$;

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
