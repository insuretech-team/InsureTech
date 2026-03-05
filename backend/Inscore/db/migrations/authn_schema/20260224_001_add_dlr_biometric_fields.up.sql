-- =====================================================
-- DLR + Biometric Indexes: authn_schema.otps + users
-- Columns are defined in otp.proto (fields 12-20) and
-- user.proto (fields 25, 29). This migration only creates
-- the performance indexes for DLR webhook lookups.
-- migration_order: 54 (after pii blind indexes)
-- =====================================================

BEGIN;

-- Index for DLR webhook lookups by provider_message_id (otps)
CREATE INDEX IF NOT EXISTS idx_otps_provider_message_id
    ON authn_schema.otps (provider_message_id)
    WHERE provider_message_id IS NOT NULL;

-- Index for DLR status filtering (analytics / BTRC compliance reports)
CREATE INDEX IF NOT EXISTS idx_otps_dlr_status_updated
    ON authn_schema.otps (dlr_status, dlr_updated_at)
    WHERE dlr_status IS NOT NULL;

COMMIT;
