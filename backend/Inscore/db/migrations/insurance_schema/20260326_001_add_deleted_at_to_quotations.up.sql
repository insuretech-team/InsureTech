-- Quotation table structure is owned by protobuf descriptors.
-- Per project convention, indexes/constraints stay in db/migrations.
-- Keep only the partial active-record index here; `deleted_at` itself is reconciled from proto.

DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = 'insurance_schema'
          AND table_name = 'quotations'
          AND column_name = 'deleted_at'
    ) THEN
        CREATE INDEX IF NOT EXISTS idx_quotations_active
            ON insurance_schema.quotations (quotation_id)
            WHERE deleted_at IS NULL;
    ELSE
        RAISE NOTICE 'insurance_schema.quotations.deleted_at does not exist yet; skipping idx_quotations_active';
    END IF;
END $$;
