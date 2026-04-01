BEGIN;

-- =============================================================================
-- Seed: B2C Customer Media & Storage RBAC Policies
-- Date: 2026-03-27
-- BUG-011 FIX: B2C users need access to their own media (KYC document uploads/views)
-- and storage files. Previously missing from the customer role policy set.
--
-- Casbin p-rule format: ptype='p', v0=role/sub, v1=domain, v2=object, v3=action, v4=effect
-- Domain for B2C: 'b2c:root'
-- =============================================================================

-- Add media + storage p-rules for B2C customer role
INSERT INTO authz_schema.casbin_rules (ptype, v0, v1, v2, v3, v4, v5)
VALUES
  -- ── Media — gateway prefix: svc:media ────────────────────────────────────
  -- B2C users can upload, view, and list their own media (KYC docs, profile photos)
  ('p', 'role:customer', 'b2c:root', 'svc:media/*',            'GET',    'allow', ''),
  ('p', 'role:customer', 'b2c:root', 'svc:media/*',            'POST',   'allow', ''),
  ('p', 'role:customer', 'b2c:root', 'svc:media/*',            'DELETE', 'allow', ''),

  -- ── Storage — gateway prefix: svc:storage ────────────────────────────────
  -- B2C users can upload and retrieve their own files (policy documents, receipts)
  ('p', 'role:customer', 'b2c:root', 'svc:storage/*',          'GET',    'allow', ''),
  ('p', 'role:customer', 'b2c:root', 'svc:storage/*',          'POST',   'allow', ''),

  -- ── Profile photo upload URL ──────────────────────────────────────────────
  -- Already covered by svc:authn/* but explicit for clarity
  ('p', 'role:customer', 'b2c:root', 'svc:authn/auth/users/*', '*',      'allow', '')

ON CONFLICT ON CONSTRAINT uq_casbin_rules_tuple DO NOTHING;

COMMIT;
