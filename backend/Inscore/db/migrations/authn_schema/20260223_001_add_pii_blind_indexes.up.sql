-- =====================================================
-- PII Blind-Index Indexes: authn_schema.users
-- Columns are defined in user.proto (fields 29, 30, 31).
-- This migration only creates the unique partial indexes
-- for equality lookups on HMAC-SHA256 blind index columns.
-- migration_order: 53 (after user table baseline)
-- =====================================================

BEGIN;

-- Unique index on mobile_number_idx (one account per mobile)
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_mobile_number_idx
    ON authn_schema.users(mobile_number_idx)
    WHERE mobile_number_idx IS NOT NULL AND deleted_at IS NULL;

-- Unique index on email_idx (one account per email)
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email_idx
    ON authn_schema.users(email_idx)
    WHERE email_idx IS NOT NULL AND deleted_at IS NULL;

-- Unique index on biometric_token_idx (one biometric per device)
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_biometric_token_idx
    ON authn_schema.users(biometric_token_idx)
    WHERE biometric_token_idx IS NOT NULL AND deleted_at IS NULL;

COMMIT;
