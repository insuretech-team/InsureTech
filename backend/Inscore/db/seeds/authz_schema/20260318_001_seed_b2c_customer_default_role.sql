BEGIN;

-- =============================================================================
-- Seed: B2C Customer Default Role Policies + Auto-Assignment
-- Date: 2026-03-18
-- Purpose: Every USER_TYPE_B2C_CUSTOMER gets the 'customer' role automatically.
--
-- Casbin model used (builtinModel in casbin_enforcer.go):
--   [request_definition]  r = sub, dom, obj, act
--   [policy_definition]   p = sub, dom, obj, act, eft   → v0,v1,v2,v3,v4
--   [role_definition]     g = _, _, _                   → v0=sub, v1=role, v2=dom
--   [matchers] g(r.sub, p.sub, r.dom) && r.dom == p.dom
--              && keyMatch2(r.obj, p.obj) && actionMatch(r.act, p.act)
--
-- Key:
--   p-rule format: ptype='p', v0=role/sub, v1=domain, v2=object, v3=action, v4=effect
--   g-rule format: ptype='g', v0=user:<uid>, v1=role:<name>, v2=domain
--   Domain for B2C: 'b2c:root'   (X-Portal=b2c, X-Tenant-ID=root)
-- =============================================================================

-- ── 1. Ensure B2C customer role exists (idempotent) ──────────────────────────
INSERT INTO authz_schema.roles (
  role_id,
  name,
  portal,
  description,
  is_system,
  is_active
)
VALUES (
  '9c3df4f3-98c4-408a-8b21-112ed8e75d9a',
  'customer',
  'PORTAL_B2C',
  'Default B2C Customer — auto-assigned on registration',
  true,
  true
)
ON CONFLICT (role_id) DO UPDATE
SET
  description = EXCLUDED.description,
  is_system   = EXCLUDED.is_system,
  is_active   = EXCLUDED.is_active,
  updated_at  = NOW();

-- ── 2. Clear any stale customer p-rules and re-seed correctly ─────────────────
DELETE FROM authz_schema.casbin_rules
WHERE ptype = 'p' AND v0 = 'role:customer';

-- ── 3. Seed Casbin p-rules for B2C customer role ─────────────────────────────
-- p-rule column mapping: v0=sub, v1=domain, v2=object, v3=action, v4=effect
-- Domain: 'b2c:root' (X-Portal=b2c, X-Tenant-ID=root)
-- Object: plain resource name OR svc:<service>/<resource> (matched by keyMatch2)
-- Action: HTTP verb (GET, POST, PUT, PATCH, DELETE, *) OR semantic (read, create, etc.)
--         actionMatch() supports both forms
INSERT INTO authz_schema.casbin_rules (ptype, v0, v1, v2, v3, v4, v5)
VALUES
  -- ── Gateway route self-access (authz middleware checks this before forwarding) ──
  ('p', 'role:customer', 'b2c:root', 'svc:authz/authz/check',    'POST',   'allow', ''),
  ('p', 'role:customer', 'b2c:root', 'svc:authz/check',          'POST',   'allow', ''),
  ('p', 'role:customer', 'b2c:root', 'svc:authz/roles',          'GET',    'allow', ''),
  ('p', 'role:customer', 'b2c:root', 'svc:authz/portals/*',      'GET',    'allow', ''),

  -- ── Auth service endpoints ────────────────────────────────────────────────────
  ('p', 'role:customer', 'b2c:root', 'svc:authn/auth/otp',       '*',      'allow', ''),
  ('p', 'role:customer', 'b2c:root', 'svc:authn/auth/login',     'POST',   'allow', ''),
  ('p', 'role:customer', 'b2c:root', 'svc:authn/auth/logout',    'POST',   'allow', ''),
  ('p', 'role:customer', 'b2c:root', 'svc:authn/auth/session',   '*',      'allow', ''),
  ('p', 'role:customer', 'b2c:root', 'svc:authn/auth/token',     'POST',   'allow', ''),
  ('p', 'role:customer', 'b2c:root', 'svc:authn/auth/password',  '*',      'allow', ''),
  ('p', 'role:customer', 'b2c:root', 'svc:gateway/me',           'GET',    'allow', ''),

  -- ── Business resources — object format matches gateway AuthZMiddleware ──────────
  -- Gateway builds: svc:<service>/* via buildObject(servicePrefix, resource)
  -- Policy (insurance policy) — gateway prefix: svc:policy
  ('p', 'role:customer', 'b2c:root', 'svc:policy/*',             'GET',    'allow', ''),
  ('p', 'role:customer', 'b2c:root', 'svc:policy/*',             'POST',   'allow', ''),
  ('p', 'role:customer', 'b2c:root', 'svc:policy/*',             'PATCH',  'allow', ''),
  -- Claim — gateway prefix: svc:claim
  ('p', 'role:customer', 'b2c:root', 'svc:claim/*',              'GET',    'allow', ''),
  ('p', 'role:customer', 'b2c:root', 'svc:claim/*',              'POST',   'allow', ''),
  ('p', 'role:customer', 'b2c:root', 'svc:claim/*',              'PATCH',  'allow', ''),
  -- Payment — gateway prefix: svc:payment
  ('p', 'role:customer', 'b2c:root', 'svc:payment/*',            'GET',    'allow', ''),
  ('p', 'role:customer', 'b2c:root', 'svc:payment/*',            'POST',   'allow', ''),
  -- Document/Profile — gateway prefix: svc:document
  ('p', 'role:customer', 'b2c:root', 'svc:document/*',           'GET',    'allow', ''),
  ('p', 'role:customer', 'b2c:root', 'svc:document/*',           'POST',   'allow', ''),
  ('p', 'role:customer', 'b2c:root', 'svc:document/*',           'PATCH',  'allow', ''),
  ('p', 'role:customer', 'b2c:root', 'svc:document/*',           'DELETE', 'allow', ''),
  -- Product (browse only) — gateway prefix: svc:product
  ('p', 'role:customer', 'b2c:root', 'svc:product/*',            'GET',    'allow', ''),
  -- Quote / Order — gateway prefix: svc:quote, svc:order
  ('p', 'role:customer', 'b2c:root', 'svc:quote/*',              'GET',    'allow', ''),
  ('p', 'role:customer', 'b2c:root', 'svc:quote/*',              'POST',   'allow', ''),
  ('p', 'role:customer', 'b2c:root', 'svc:order/*',              'GET',    'allow', ''),
  ('p', 'role:customer', 'b2c:root', 'svc:order/*',              'POST',   'allow', ''),
  -- Notification — gateway prefix: svc:notification
  ('p', 'role:customer', 'b2c:root', 'svc:notification/*',       'GET',    'allow', ''),
  ('p', 'role:customer', 'b2c:root', 'svc:notification/*',       'DELETE', 'allow', ''),
  -- Beneficiary — gateway prefix: svc:beneficiary
  ('p', 'role:customer', 'b2c:root', 'svc:beneficiary/*',        '*',      'allow', ''),
  -- KYC — gateway prefix: svc:kyc
  ('p', 'role:customer', 'b2c:root', 'svc:kyc/*',                '*',      'allow', ''),
  -- Support — gateway prefix: svc:support
  ('p', 'role:customer', 'b2c:root', 'svc:support/*',            'GET',    'allow', ''),
  ('p', 'role:customer', 'b2c:root', 'svc:support/*',            'POST',   'allow', '')

ON CONFLICT ON CONSTRAINT uq_casbin_rules_tuple DO NOTHING;

-- ── 4. Fix any stale g-rules with wrong column order and re-insert correctly ──
-- g-rule format: v0=user:<uid>, v1=role:<name>, v2=domain (NOT v1=domain, v2=role)
DELETE FROM authz_schema.casbin_rules
WHERE ptype = 'g' AND v2 = 'role:customer'; -- stale wrong-order rows

-- ── 5. Auto-assign customer role to ALL existing B2C users ───────────────────
-- g-rule: (user:<uid>, role:customer, b2c:root)
DO $$
DECLARE
  b2c_role_id  UUID := '9c3df4f3-98c4-408a-8b21-112ed8e75d9a';
  system_uid   UUID := '00000000-0000-0000-0000-000000000001';
  rec          RECORD;
BEGIN
  FOR rec IN
    SELECT user_id
    FROM authn_schema.users
    WHERE user_type = 'USER_TYPE_B2C_CUSTOMER'
      AND deleted_at IS NULL
  LOOP
    -- user_roles table
    INSERT INTO authz_schema.user_roles (
      user_role_id, user_id, role_id, domain, assigned_by, assigned_at
    )
    VALUES (
      gen_random_uuid(), rec.user_id, b2c_role_id, 'b2c:root', system_uid, NOW()
    )
    ON CONFLICT ON CONSTRAINT uq_user_roles_user_role_domain DO NOTHING;

    -- Casbin g-rule: CORRECT ORDER (sub, role, domain)
    INSERT INTO authz_schema.casbin_rules (ptype, v0, v1, v2, v3, v4, v5)
    VALUES (
      'g',
      'user:' || rec.user_id::TEXT,
      'role:customer',
      'b2c:root',
      '', '', ''
    )
    ON CONFLICT ON CONSTRAINT uq_casbin_rules_tuple DO NOTHING;

    RAISE NOTICE 'Assigned customer role to B2C user: %', rec.user_id;
  END LOOP;
END $$;

-- ── 6. DB trigger: auto-assign customer role on new B2C user INSERT ───────────
CREATE OR REPLACE FUNCTION authz_schema.fn_auto_assign_b2c_customer_role()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
DECLARE
  b2c_role_id UUID := '9c3df4f3-98c4-408a-8b21-112ed8e75d9a';
  system_uid  UUID := '00000000-0000-0000-0000-000000000001';
BEGIN
  IF NEW.user_type = 'USER_TYPE_B2C_CUSTOMER' AND NEW.deleted_at IS NULL THEN
    INSERT INTO authz_schema.user_roles (
      user_role_id, user_id, role_id, domain, assigned_by, assigned_at
    )
    VALUES (
      gen_random_uuid(), NEW.user_id, b2c_role_id, 'b2c:root', system_uid, NOW()
    )
    ON CONFLICT ON CONSTRAINT uq_user_roles_user_role_domain DO NOTHING;

    -- g-rule: CORRECT ORDER (sub=user:<uid>, role=role:customer, dom=b2c:root)
    INSERT INTO authz_schema.casbin_rules (ptype, v0, v1, v2, v3, v4, v5)
    VALUES (
      'g',
      'user:' || NEW.user_id::TEXT,
      'role:customer',
      'b2c:root',
      '', '', ''
    )
    ON CONFLICT ON CONSTRAINT uq_casbin_rules_tuple DO NOTHING;
  END IF;
  RETURN NEW;
END;
$$;

DROP TRIGGER IF EXISTS trg_auto_assign_b2c_customer_role ON authn_schema.users;

CREATE TRIGGER trg_auto_assign_b2c_customer_role
  AFTER INSERT ON authn_schema.users
  FOR EACH ROW
  EXECUTE FUNCTION authz_schema.fn_auto_assign_b2c_customer_role();

COMMIT;


