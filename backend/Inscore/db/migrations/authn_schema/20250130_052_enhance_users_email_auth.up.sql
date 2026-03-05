-- =====================================================
-- Email Auth Enhancement: authn_schema.users
-- Indexes, constraints, and comments for:
--   user_type, email_verified, email_verified_at,
--   email_login_attempts, email_locked_until
-- Columns are created by proto-generated migration.
-- migration_order: 52 (after existing user enhancements)
-- =====================================================

BEGIN;

-- Index on user_type
CREATE INDEX IF NOT EXISTS idx_users_user_type
    ON authn_schema.users(user_type);

COMMENT ON COLUMN authn_schema.users.user_type IS
    'User type: B2C_CUSTOMER, AGENT, BUSINESS_BENEFICIARY, SYSTEM_USER. Controls auth method routing.';

-- Partial index on email_verified (only non-deleted users)
CREATE INDEX IF NOT EXISTS idx_users_email_verified
    ON authn_schema.users(email_verified)
    WHERE deleted_at IS NULL;

COMMENT ON COLUMN authn_schema.users.email_verified IS
    'Email address verification status. Must be true before email OTP login is allowed.';

COMMENT ON COLUMN authn_schema.users.email_verified_at IS
    'Timestamp when email was first verified.';

COMMENT ON COLUMN authn_schema.users.email_login_attempts IS
    'Failed email OTP login counter. Locked out after 5 failed attempts.';

-- Partial index on email_locked_until (only future lockouts)
CREATE INDEX IF NOT EXISTS idx_users_email_locked_until
    ON authn_schema.users(email_locked_until)
    WHERE email_locked_until IS NOT NULL;

COMMENT ON COLUMN authn_schema.users.email_locked_until IS
    'Email auth lockout expiry. Set to NOW()+30min after 5 failed email OTP attempts.';

-- Enforce: BUSINESS_BENEFICIARY and SYSTEM_USER must have an email
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'chk_users_email_required_for_business'
          AND conrelid = 'authn_schema.users'::regclass
    ) THEN
        ALTER TABLE authn_schema.users
            ADD CONSTRAINT chk_users_email_required_for_business
            CHECK (
                user_type NOT IN ('BUSINESS_BENEFICIARY', 'SYSTEM_USER')
                OR (email IS NOT NULL AND email <> '')
            );

        RAISE NOTICE 'Constraint chk_users_email_required_for_business added';
    ELSE
        RAISE NOTICE 'Constraint already exists, skipping';
    END IF;
END $$;

COMMENT ON TABLE authn_schema.users IS 'Enhanced with email auth fields: user_type, email_verified, email_verified_at, email_login_attempts, email_locked_until';

COMMIT;
