-- Migration: 20260327_001_widen_quotations_insurer_name
-- Reason: PoliSync stores serialized quotation metadata JSON in the insurer_name column
--         which can exceed VARCHAR(255). Widen to TEXT to prevent truncation errors.
-- Applied manually on 2026-03-27 during B2C API testing (see bug_fix.md BUG-CREATE-QUOTE).
-- Reference: insurance_service.go:1698 "value too long for type character varying(255)"

ALTER TABLE insurance_schema.quotations
    ALTER COLUMN insurer_name TYPE TEXT;
