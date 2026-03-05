-- Align legacy storage.stored_files with proto-first storage_schema.files.
-- Safe to run multiple times.

DO $$
BEGIN
    -- Nothing to migrate if legacy table doesn't exist.
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.tables
        WHERE table_schema = 'storage'
          AND table_name = 'stored_files'
    ) THEN
        RETURN;
    END IF;

    -- Ensure destination schema/table exists (table itself is proto-managed).
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.tables
        WHERE table_schema = 'storage_schema'
          AND table_name = 'files'
    ) THEN
        RETURN;
    END IF;

    INSERT INTO storage_schema.files (
        file_id,
        tenant_id,
        filename,
        content_type,
        size_bytes,
        storage_key,
        bucket,
        url,
        cdn_url,
        file_type,
        reference_id,
        reference_type,
        is_public,
        expires_at,
        uploaded_by,
        created_at,
        updated_at
    )
    SELECT
        sf.file_id,
        sf.tenant_id,
        sf.filename,
        COALESCE(NULLIF(sf.content_type, ''), 'application/octet-stream') AS content_type,
        COALESCE(sf.size_bytes, 0) AS size_bytes,
        sf.s3_key AS storage_key,
        sf.bucket,
        COALESCE(NULLIF(sf.url, ''), 'https://invalid.local/' || sf.file_id::text) AS url,
        sf.cdn_url,
        CASE sf.file_type::text
            WHEN '0' THEN 'FILE_TYPE_UNSPECIFIED'
            WHEN '1' THEN 'FILE_TYPE_IMAGE'
            WHEN '2' THEN 'FILE_TYPE_PDF'
            WHEN '3' THEN 'FILE_TYPE_DOCUMENT'
            WHEN '4' THEN 'FILE_TYPE_SHIPPING_LABEL'
            WHEN '5' THEN 'FILE_TYPE_INVOICE'
            WHEN '6' THEN 'FILE_TYPE_RECEIPT'
            WHEN '7' THEN 'FILE_TYPE_PRODUCT_IMAGE'
            WHEN '8' THEN 'FILE_TYPE_AVATAR'
            WHEN '9' THEN 'FILE_TYPE_VIDEO'
            WHEN 'FILE_TYPE_UNSPECIFIED' THEN 'FILE_TYPE_UNSPECIFIED'
            WHEN 'FILE_TYPE_IMAGE' THEN 'FILE_TYPE_IMAGE'
            WHEN 'FILE_TYPE_PDF' THEN 'FILE_TYPE_PDF'
            WHEN 'FILE_TYPE_DOCUMENT' THEN 'FILE_TYPE_DOCUMENT'
            WHEN 'FILE_TYPE_SHIPPING_LABEL' THEN 'FILE_TYPE_SHIPPING_LABEL'
            WHEN 'FILE_TYPE_INVOICE' THEN 'FILE_TYPE_INVOICE'
            WHEN 'FILE_TYPE_RECEIPT' THEN 'FILE_TYPE_RECEIPT'
            WHEN 'FILE_TYPE_PRODUCT_IMAGE' THEN 'FILE_TYPE_PRODUCT_IMAGE'
            WHEN 'FILE_TYPE_AVATAR' THEN 'FILE_TYPE_AVATAR'
            WHEN 'FILE_TYPE_VIDEO' THEN 'FILE_TYPE_VIDEO'
            ELSE 'FILE_TYPE_OTHER'
        END AS file_type,
        CASE
            WHEN sf.reference_id IS NULL OR sf.reference_id = '' THEN NULL
            WHEN sf.reference_id ~* '^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$' THEN sf.reference_id::uuid
            ELSE NULL
        END AS reference_id,
        sf.reference_type,
        COALESCE(sf.is_public, false) AS is_public,
        sf.expires_at,
        COALESCE(sf.uploaded_by, sf.tenant_id) AS uploaded_by,
        COALESCE(sf.created_at, NOW()) AS created_at,
        COALESCE(sf.updated_at, NOW()) AS updated_at
    FROM storage.stored_files sf
    ON CONFLICT (file_id) DO NOTHING;
END $$;

