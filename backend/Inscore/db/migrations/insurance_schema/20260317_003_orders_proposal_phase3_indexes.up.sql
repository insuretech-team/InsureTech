-- Enhancement: add missing indexes and metadata for orders proposal-flow columns
-- These indexes correspond to proto-added order fields 33-39.
-- Table and column shape remains owned by proto-first migration.
-- Date: 2026-03-17

BEGIN;

DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = 'insurance_schema'
          AND table_name = 'orders'
          AND column_name = 'proposal_id'
    ) THEN
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_orders_proposal_id ON insurance_schema.orders (proposal_id) WHERE proposal_id IS NOT NULL';
        EXECUTE 'COMMENT ON COLUMN insurance_schema.orders.proposal_id IS ''Proto-managed FK to insurance_schema.insurance_proposals; SQL enhancement adds lookup index only''';
    END IF;

    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = 'insurance_schema'
          AND table_name = 'orders'
          AND column_name = 'refund_id'
    ) THEN
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_orders_refund_id ON insurance_schema.orders (refund_id) WHERE refund_id IS NOT NULL';
        EXECUTE 'COMMENT ON COLUMN insurance_schema.orders.refund_id IS ''Proto-managed FK to payment_schema.refunds for insurer proposal rejection reversals''';
    END IF;

    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = 'insurance_schema'
          AND table_name = 'orders'
          AND column_name = 'insurer_id'
    ) THEN
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_orders_insurer_id ON insurance_schema.orders (insurer_id) WHERE insurer_id IS NOT NULL';
        EXECUTE 'COMMENT ON COLUMN insurance_schema.orders.insurer_id IS ''Proto-managed FK to insurance_schema.insurers; identifies target insurer for proposal review''';
    END IF;

    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = 'insurance_schema'
          AND table_name = 'orders'
          AND column_name = 'proposal_status'
    ) THEN
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_orders_proposal_status ON insurance_schema.orders (proposal_status)';
        EXECUTE 'COMMENT ON COLUMN insurance_schema.orders.proposal_status IS ''Proposal-review dimension of the order lifecycle, managed by proto enum values''';
    END IF;

    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = 'insurance_schema'
          AND table_name = 'orders'
          AND column_name = 'proposal_submitted_at'
    ) THEN
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_orders_proposal_submitted_at ON insurance_schema.orders (proposal_submitted_at DESC) WHERE proposal_submitted_at IS NOT NULL';
        EXECUTE 'COMMENT ON COLUMN insurance_schema.orders.proposal_submitted_at IS ''Timestamp when a paid order is submitted to the insurer as a proposal''';
    END IF;

    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = 'insurance_schema'
          AND table_name = 'orders'
          AND column_name = 'proposal_decided_at'
    ) THEN
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_orders_proposal_decided_at ON insurance_schema.orders (proposal_decided_at DESC) WHERE proposal_decided_at IS NOT NULL';
        EXECUTE 'COMMENT ON COLUMN insurance_schema.orders.proposal_decided_at IS ''Timestamp when the insurer decision is captured''';
    END IF;

    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = 'insurance_schema'
          AND table_name = 'orders'
          AND column_name = 'proposal_decision_reason'
    ) THEN
        EXECUTE 'COMMENT ON COLUMN insurance_schema.orders.proposal_decision_reason IS ''Insurer approval/rejection explanation captured in the order projection''';
    END IF;

    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = 'insurance_schema'
          AND table_name = 'orders'
          AND column_name = 'proposal_status'
    )
    AND EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = 'insurance_schema'
          AND table_name = 'orders'
          AND column_name = 'updated_at'
    ) THEN
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_orders_proposal_status_updated_at ON insurance_schema.orders (proposal_status, updated_at DESC)';
    END IF;

    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = 'insurance_schema'
          AND table_name = 'orders'
          AND column_name = 'insurer_id'
    )
    AND EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = 'insurance_schema'
          AND table_name = 'orders'
          AND column_name = 'proposal_status'
    )
    AND EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = 'insurance_schema'
          AND table_name = 'orders'
          AND column_name = 'updated_at'
    ) THEN
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_orders_insurer_proposal_status_updated_at ON insurance_schema.orders (insurer_id, proposal_status, updated_at DESC) WHERE insurer_id IS NOT NULL';
    END IF;
END $$;

COMMENT ON TABLE insurance_schema.orders IS 'Proto-first order lifecycle table enhanced with proposal-review lookup indexes and metadata';

COMMIT;
