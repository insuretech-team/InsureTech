-- =====================================================
-- Indexes and comments for session_token_lookup
-- Enables production-grade validation: lookup by SHA-256(token) + verify bcrypt(token_hash)
-- Column is created by proto-generated migration.
-- =====================================================

BEGIN;

-- Hash index for fast lookup
CREATE INDEX IF NOT EXISTS idx_sessions_token_lookup
    ON authn_schema.sessions USING hash (session_token_lookup)
    WHERE session_type = 'SERVER_SIDE';

-- Partial unique index: one token maps to one session
CREATE UNIQUE INDEX IF NOT EXISTS uq_sessions_token_lookup
    ON authn_schema.sessions (session_token_lookup)
    WHERE session_type = 'SERVER_SIDE' AND session_token_lookup IS NOT NULL;

COMMENT ON COLUMN authn_schema.sessions.session_token_lookup IS
    'Deterministic SHA-256 hex of the session token (server-side sessions only). Used for lookup; actual token is stored as bcrypt in session_token_hash.';

COMMIT;
