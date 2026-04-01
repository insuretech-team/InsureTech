-- =====================================================
-- Production Enhancement: insurance_schema.insurance_proposals
-- =====================================================
-- Proto-first owns the table/columns/FK definitions.
-- This enhancement layer adds supporting indexes, trigger automation, and metadata only.

BEGIN;

DO $$
DECLARE
    pk_col TEXT;
BEGIN
    IF to_regclass('insurance_schema.insurance_proposals') IS NULL THEN
        RETURN;
    END IF;

    SELECT a.attname INTO pk_col
    FROM pg_index i
    JOIN pg_attribute a ON a.attrelid = i.indrelid AND a.attnum = ANY(i.indkey)
    WHERE i.indrelid = 'insurance_schema.insurance_proposals'::regclass
      AND i.indisprimary
    LIMIT 1;

    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'insurance_schema' AND table_name = 'insurance_proposals' AND column_name = 'created_at') THEN
        IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'insurance_schema' AND table_name = 'insurance_proposals' AND column_name = 'deleted_at') THEN
            EXECUTE 'CREATE INDEX IF NOT EXISTS idx_insurance_proposals_created_at ON insurance_schema.insurance_proposals(created_at DESC) WHERE deleted_at IS NULL';
        ELSE
            EXECUTE 'CREATE INDEX IF NOT EXISTS idx_insurance_proposals_created_at ON insurance_schema.insurance_proposals(created_at DESC)';
        END IF;
    END IF;

    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'insurance_schema' AND table_name = 'insurance_proposals' AND column_name = 'updated_at') THEN
        IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'insurance_schema' AND table_name = 'insurance_proposals' AND column_name = 'deleted_at') THEN
            EXECUTE 'CREATE INDEX IF NOT EXISTS idx_insurance_proposals_updated_at ON insurance_schema.insurance_proposals(updated_at DESC) WHERE deleted_at IS NULL';
        ELSE
            EXECUTE 'CREATE INDEX IF NOT EXISTS idx_insurance_proposals_updated_at ON insurance_schema.insurance_proposals(updated_at DESC)';
        END IF;
    END IF;

    IF pk_col IS NOT NULL
       AND EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'insurance_schema' AND table_name = 'insurance_proposals' AND column_name = 'deleted_at') THEN
        EXECUTE format('CREATE INDEX IF NOT EXISTS idx_insurance_proposals_active ON insurance_schema.insurance_proposals(%I) WHERE deleted_at IS NULL', pk_col);
    END IF;

    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'insurance_schema' AND table_name = 'insurance_proposals' AND column_name = 'insurer_id')
       AND EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'insurance_schema' AND table_name = 'insurance_proposals' AND column_name = 'status')
       AND EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'insurance_schema' AND table_name = 'insurance_proposals' AND column_name = 'submitted_at') THEN
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_insurance_proposals_insurer_status_submitted_at ON insurance_schema.insurance_proposals(insurer_id, status, submitted_at DESC)';
    END IF;

    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'insurance_schema' AND table_name = 'insurance_proposals' AND column_name = 'customer_id')
       AND EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'insurance_schema' AND table_name = 'insurance_proposals' AND column_name = 'status')
       AND EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'insurance_schema' AND table_name = 'insurance_proposals' AND column_name = 'submitted_at') THEN
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_insurance_proposals_customer_status_submitted_at ON insurance_schema.insurance_proposals(customer_id, status, submitted_at DESC)';
    END IF;

    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'insurance_schema' AND table_name = 'insurance_proposals' AND column_name = 'order_id') THEN
        EXECUTE 'CREATE UNIQUE INDEX IF NOT EXISTS idx_insurance_proposals_order_id ON insurance_schema.insurance_proposals(order_id)';
    END IF;

    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'insurance_schema' AND table_name = 'insurance_proposals' AND column_name = 'quotation_id') THEN
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_insurance_proposals_quotation_id ON insurance_schema.insurance_proposals(quotation_id)';
    END IF;

    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'insurance_schema' AND table_name = 'insurance_proposals' AND column_name = 'refund_id') THEN
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_insurance_proposals_refund_id ON insurance_schema.insurance_proposals(refund_id) WHERE refund_id IS NOT NULL';
    END IF;

    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'insurance_schema' AND table_name = 'insurance_proposals' AND column_name = 'correlation_id') THEN
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_insurance_proposals_correlation_id ON insurance_schema.insurance_proposals(correlation_id) WHERE correlation_id IS NOT NULL';
    END IF;
END $$;

DO $$
BEGIN
    IF to_regclass('insurance_schema.insurance_proposals') IS NULL THEN
        RETURN;
    END IF;

    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'insurance_schema' AND table_name = 'insurance_proposals' AND column_name = 'updated_at') THEN
        EXECUTE 'CREATE OR REPLACE FUNCTION insurance_schema.trg_insurance_proposals_updated_at() RETURNS TRIGGER AS $body$ BEGIN NEW.updated_at = CURRENT_TIMESTAMP; RETURN NEW; END; $body$ LANGUAGE plpgsql';
        EXECUTE 'DROP TRIGGER IF EXISTS trg_insurance_proposals_update ON insurance_schema.insurance_proposals';
        EXECUTE 'CREATE TRIGGER trg_insurance_proposals_update BEFORE UPDATE ON insurance_schema.insurance_proposals FOR EACH ROW EXECUTE FUNCTION insurance_schema.trg_insurance_proposals_updated_at()';
    END IF;
END $$;

DO $$
BEGIN
    IF to_regclass('insurance_schema.insurance_proposals') IS NULL THEN
        RETURN;
    END IF;

    EXECUTE 'COMMENT ON TABLE insurance_schema.insurance_proposals IS ''Proto-first insurer proposal table enhanced with composite lookup indexes, updated_at trigger, and metadata''';

    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'insurance_schema' AND table_name = 'insurance_proposals' AND column_name = 'order_id') THEN
        EXECUTE 'COMMENT ON COLUMN insurance_schema.insurance_proposals.order_id IS ''Proto-managed one-to-one reference back to insurance_schema.orders''';
    END IF;
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'insurance_schema' AND table_name = 'insurance_proposals' AND column_name = 'quotation_id') THEN
        EXECUTE 'COMMENT ON COLUMN insurance_schema.insurance_proposals.quotation_id IS ''Source quotation converted into a paid order before insurer review''';
    END IF;
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'insurance_schema' AND table_name = 'insurance_proposals' AND column_name = 'insurer_id') THEN
        EXECUTE 'COMMENT ON COLUMN insurance_schema.insurance_proposals.insurer_id IS ''Target insurer for accept/reject decisioning''';
    END IF;
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'insurance_schema' AND table_name = 'insurance_proposals' AND column_name = 'status') THEN
        EXECUTE 'COMMENT ON COLUMN insurance_schema.insurance_proposals.status IS ''Proposal state machine managed by proto enum ProposalStatus''';
    END IF;
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'insurance_schema' AND table_name = 'insurance_proposals' AND column_name = 'submission_payload') THEN
        EXECUTE 'COMMENT ON COLUMN insurance_schema.insurance_proposals.submission_payload IS ''Serialized outbound insurer request payload for audit and replay''';
    END IF;
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'insurance_schema' AND table_name = 'insurance_proposals' AND column_name = 'insurer_response_payload') THEN
        EXECUTE 'COMMENT ON COLUMN insurance_schema.insurance_proposals.insurer_response_payload IS ''Serialized inbound insurer decision payload''';
    END IF;
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'insurance_schema' AND table_name = 'insurance_proposals' AND column_name = 'decision_reason') THEN
        EXECUTE 'COMMENT ON COLUMN insurance_schema.insurance_proposals.decision_reason IS ''Insurer explanation for approval/rejection or referral outcome''';
    END IF;
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'insurance_schema' AND table_name = 'insurance_proposals' AND column_name = 'approved_policy_id') THEN
        EXECUTE 'COMMENT ON COLUMN insurance_schema.insurance_proposals.approved_policy_id IS ''Policy issued after proposal approval; FK remains proto-managed''';
    END IF;
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'insurance_schema' AND table_name = 'insurance_proposals' AND column_name = 'refund_id') THEN
        EXECUTE 'COMMENT ON COLUMN insurance_schema.insurance_proposals.refund_id IS ''Refund linked after paid proposal rejection; FK remains proto-managed''';
    END IF;
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'insurance_schema' AND table_name = 'insurance_proposals' AND column_name = 'correlation_id') THEN
        EXECUTE 'COMMENT ON COLUMN insurance_schema.insurance_proposals.correlation_id IS ''Cross-service tracing and saga correlation id''';
    END IF;
END $$;

COMMIT;
