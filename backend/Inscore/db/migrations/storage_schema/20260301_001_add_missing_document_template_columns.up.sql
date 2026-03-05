-- Proto-first rule:
-- Column/table shape must come from
-- proto/insuretech/document/entity/v1/document_template.proto.
-- This file keeps only SQL enhancements (constraints/triggers/indexes).

BEGIN;

-- Ensure name has unique constraint for ON CONFLICT(name) upsert.
DO $$
BEGIN
    IF to_regclass('storage_schema.document_templates') IS NOT NULL
       AND EXISTS (
         SELECT 1
         FROM information_schema.columns
         WHERE table_schema = 'storage_schema'
           AND table_name = 'document_templates'
           AND column_name = 'name'
       )
       AND NOT EXISTS (
         SELECT 1
         FROM pg_constraint
         WHERE conrelid = 'storage_schema.document_templates'::regclass
           AND contype = 'u'
           AND conname = 'uq_document_templates_name'
       ) THEN
      ALTER TABLE storage_schema.document_templates
          ADD CONSTRAINT uq_document_templates_name UNIQUE (name);
    END IF;
END $$;

-- Keep updated_at maintenance trigger for update paths.
DO $$
BEGIN
    IF to_regclass('storage_schema.document_templates') IS NOT NULL
       AND EXISTS (
         SELECT 1
         FROM information_schema.columns
         WHERE table_schema = 'storage_schema'
           AND table_name = 'document_templates'
           AND column_name = 'updated_at'
       ) THEN
      EXECUTE 'CREATE OR REPLACE FUNCTION storage_schema.trg_document_templates_updated_at() RETURNS TRIGGER AS $body$ BEGIN NEW.updated_at = CURRENT_TIMESTAMP; RETURN NEW; END; $body$ LANGUAGE plpgsql';
      EXECUTE 'DROP TRIGGER IF EXISTS trg_document_templates_update ON storage_schema.document_templates';
      EXECUTE 'CREATE TRIGGER trg_document_templates_update BEFORE UPDATE ON storage_schema.document_templates FOR EACH ROW EXECUTE FUNCTION storage_schema.trg_document_templates_updated_at()';
    END IF;
END $$;

-- Performance indexes.
DO $$
BEGIN
    IF to_regclass('storage_schema.document_templates') IS NOT NULL
       AND EXISTS (
         SELECT 1
         FROM information_schema.columns
         WHERE table_schema = 'storage_schema'
           AND table_name = 'document_templates'
           AND column_name = 'created_at'
       ) THEN
      EXECUTE 'CREATE INDEX IF NOT EXISTS idx_document_templates_created_at ON storage_schema.document_templates (created_at DESC)';
    END IF;

    IF to_regclass('storage_schema.document_templates') IS NOT NULL
       AND EXISTS (
         SELECT 1
         FROM information_schema.columns
         WHERE table_schema = 'storage_schema'
           AND table_name = 'document_templates'
           AND column_name = 'updated_at'
       ) THEN
      EXECUTE 'CREATE INDEX IF NOT EXISTS idx_document_templates_updated_at ON storage_schema.document_templates (updated_at DESC)';
    END IF;
END $$;

COMMIT;
