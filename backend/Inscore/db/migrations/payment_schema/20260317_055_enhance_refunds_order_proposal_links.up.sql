-- Enhancement: add proposal-rejection refund indexes and metadata
-- Proto-first owns the refunds table shape. This file adds supporting indexes/comments only.
-- Table and column shape remains owned by proto-first migration.

BEGIN;

DO $$
BEGIN
    IF to_regclass('payment_schema.refunds') IS NULL THEN
        RETURN;
    END IF;

    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = 'payment_schema'
          AND table_name = 'refunds'
          AND column_name = 'order_id'
    ) THEN
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_refunds_order_id ON payment_schema.refunds (order_id) WHERE order_id IS NOT NULL';
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_refunds_order_status ON payment_schema.refunds (order_id, status) WHERE order_id IS NOT NULL';
        EXECUTE 'COMMENT ON COLUMN payment_schema.refunds.order_id IS ''Order reference for refunds created before policy issuance (e.g. insurer proposal rejection)''';
    END IF;

    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = 'payment_schema'
          AND table_name = 'refunds'
          AND column_name = 'proposal_id'
    ) THEN
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_refunds_proposal_id ON payment_schema.refunds (proposal_id) WHERE proposal_id IS NOT NULL';
        EXECUTE 'COMMENT ON COLUMN payment_schema.refunds.proposal_id IS ''Proposal reference that caused the refund; FK remains proto-managed''';
    END IF;

    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = 'payment_schema'
          AND table_name = 'refunds'
          AND column_name = 'reason'
    )
    AND EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = 'payment_schema'
          AND table_name = 'refunds'
          AND column_name = 'status'
    ) THEN
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_refunds_reason_status ON payment_schema.refunds (reason, status)';
    END IF;

    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = 'payment_schema'
          AND table_name = 'refunds'
          AND column_name = 'proposal_id'
    )
    AND EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = 'payment_schema'
          AND table_name = 'refunds'
          AND column_name = 'status'
    ) THEN
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_refunds_proposal_status ON payment_schema.refunds (proposal_id, status) WHERE proposal_id IS NOT NULL';
    END IF;
END $$;

COMMENT ON TABLE payment_schema.refunds IS 'Proto-first refund table enhanced with order/proposal reversal indexes and metadata';

COMMIT;
