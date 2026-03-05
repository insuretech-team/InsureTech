BEGIN;
DROP INDEX IF EXISTS authn_schema.idx_otps_provider_message_id;
DROP INDEX IF EXISTS authn_schema.idx_otps_dlr_status_updated;
COMMIT;
