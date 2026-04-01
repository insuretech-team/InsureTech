-- Per project convention, indexes/constraints stay in db/migrations.
-- Only roll back the helper index here.
-- Structural rollback is handled by protobuf changes plus `dbops migrate --prune`
-- when that workflow is explicitly requested.

DROP INDEX IF EXISTS insurance_schema.idx_quotations_active;
