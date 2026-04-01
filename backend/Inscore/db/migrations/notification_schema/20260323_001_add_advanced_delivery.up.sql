BEGIN;

DO $$
BEGIN
    IF to_regclass('notification_schema.push_device_tokens') IS NOT NULL THEN
        EXECUTE 'CREATE UNIQUE INDEX IF NOT EXISTS uq_push_device_tokens_provider_token ON notification_schema.push_device_tokens(provider, device_token)';
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_push_device_tokens_user_active ON notification_schema.push_device_tokens(user_id, is_active)';
        EXECUTE '' ||
            'CREATE OR REPLACE FUNCTION notification_schema.trg_push_device_tokens_updated_at() ' ||
            'RETURNS TRIGGER AS $body$ ' ||
            'BEGIN NEW.updated_at = CURRENT_TIMESTAMP; RETURN NEW; END; ' ||
            '$body$ LANGUAGE plpgsql';
        EXECUTE 'DROP TRIGGER IF EXISTS trg_push_device_tokens_update ON notification_schema.push_device_tokens';
        EXECUTE 'CREATE TRIGGER trg_push_device_tokens_update BEFORE UPDATE ON notification_schema.push_device_tokens FOR EACH ROW EXECUTE FUNCTION notification_schema.trg_push_device_tokens_updated_at()';
        EXECUTE 'COMMENT ON TABLE notification_schema.push_device_tokens IS ''Registered mobile and web push tokens for notification delivery''';
    END IF;
END $$;

DO $$
BEGIN
    IF to_regclass('notification_schema.webhook_subscriptions') IS NOT NULL THEN
        EXECUTE '' ||
            'CREATE OR REPLACE FUNCTION notification_schema.trg_webhook_subscriptions_updated_at() ' ||
            'RETURNS TRIGGER AS $body$ ' ||
            'BEGIN NEW.updated_at = CURRENT_TIMESTAMP; RETURN NEW; END; ' ||
            '$body$ LANGUAGE plpgsql';
        EXECUTE 'DROP TRIGGER IF EXISTS trg_webhook_subscriptions_update ON notification_schema.webhook_subscriptions';
        EXECUTE 'CREATE TRIGGER trg_webhook_subscriptions_update BEFORE UPDATE ON notification_schema.webhook_subscriptions FOR EACH ROW EXECUTE FUNCTION notification_schema.trg_webhook_subscriptions_updated_at()';
        EXECUTE 'COMMENT ON TABLE notification_schema.webhook_subscriptions IS ''External subscribers that receive signed notification lifecycle webhooks''';
    END IF;
END $$;

DO $$
BEGIN
    IF to_regclass('notification_schema.webhook_delivery_attempts') IS NOT NULL THEN
        EXECUTE 'CREATE UNIQUE INDEX IF NOT EXISTS uq_webhook_delivery_attempts_subscription_notification_event ON notification_schema.webhook_delivery_attempts(subscription_id, notification_id, lifecycle_event)';
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_webhook_delivery_attempts_due ON notification_schema.webhook_delivery_attempts(status, scheduled_at, created_at)';
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_webhook_delivery_attempts_subscription_created_at ON notification_schema.webhook_delivery_attempts(subscription_id, created_at DESC)';
        EXECUTE '' ||
            'CREATE OR REPLACE FUNCTION notification_schema.trg_webhook_delivery_attempts_updated_at() ' ||
            'RETURNS TRIGGER AS $body$ ' ||
            'BEGIN NEW.updated_at = CURRENT_TIMESTAMP; RETURN NEW; END; ' ||
            '$body$ LANGUAGE plpgsql';
        EXECUTE 'DROP TRIGGER IF EXISTS trg_webhook_delivery_attempts_update ON notification_schema.webhook_delivery_attempts';
        EXECUTE 'CREATE TRIGGER trg_webhook_delivery_attempts_update BEFORE UPDATE ON notification_schema.webhook_delivery_attempts FOR EACH ROW EXECUTE FUNCTION notification_schema.trg_webhook_delivery_attempts_updated_at()';
        EXECUTE 'COMMENT ON TABLE notification_schema.webhook_delivery_attempts IS ''Retryable outbound webhook deliveries emitted by the notification service''';
    END IF;
END $$;

COMMIT;
