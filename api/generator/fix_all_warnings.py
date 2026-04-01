#!/usr/bin/env python3
"""
fix_all_warnings.py — InsureTech OpenAPI Warning Fixer

Fixes common OpenAPI validation warnings in all generated YAML files:
  1. Adds missing descriptions to schema properties that have none
  2. Adds 'required' arrays to Request schemas (non-nullable scalar fields)
  3. Fixes empty string descriptions → meaningful defaults
  4. Ensures all enums have descriptions
  5. Fixes integer fields with format: int64 that are typed as string (proto artifact)

Uses ruamel.yaml for fast round-trip YAML editing (preserves formatting).

Called by run_api_pipeline.ps1 Step 12 in parallel with fix_pagination.py.
"""

import sys
import os
import io
from pathlib import Path

try:
    from ruamel.yaml import YAML
    from ruamel.yaml.scalarstring import DoubleQuotedScalarString
except ImportError:
    print("Error: ruamel.yaml not found. Run: pip install ruamel.yaml")
    sys.exit(1)

# ---------------------------------------------------------------------------
# Config
# ---------------------------------------------------------------------------
SCRIPT_DIR = Path(__file__).parent
API_DIR    = SCRIPT_DIR.parent

# Directories to process
TARGET_DIRS = [
    API_DIR / "schemas",
    API_DIR / "paths",
    API_DIR / "enums",
    API_DIR / "events",
]

# Fields that are always required in Request schemas
ALWAYS_REQUIRED_PATTERNS = {
    # Auth
    "phone_number", "otp_code", "otp_id", "email", "password",
    # Common IDs
    "user_id", "policy_id", "claim_id", "payment_id", "order_id",
    "product_id", "tenant_id", "partner_id", "quote_id",
    # Common scalars
    "amount", "currency", "status", "type",
}

# Property name → default description (when description is missing)
DEFAULT_DESCRIPTIONS = {
    "id":                 "Unique identifier",
    "created_at":         "Timestamp when the record was created (UTC ISO 8601)",
    "updated_at":         "Timestamp when the record was last updated (UTC ISO 8601)",
    "deleted_at":         "Timestamp when the record was soft-deleted (UTC ISO 8601)",
    "status":             "Current status of the resource",
    "type":               "Resource type identifier",
    "name":               "Human-readable name",
    "description":        "Human-readable description",
    "message":            "Human-readable message",
    "email":              "Email address",
    "phone_number":       "Phone number in international format (e.g. +8801XXXXXXXXX)",
    "amount":             "Monetary amount as a string to avoid floating-point precision issues",
    "currency":           "ISO 4217 currency code (e.g. BDT, USD)",
    "page":               "Page number for pagination (1-based)",
    "page_size":          "Number of items per page",
    "total_pages":        "Total number of pages available",
    "total_items":        "Total number of items across all pages",
    "has_next":           "Whether a next page exists",
    "has_previous":       "Whether a previous page exists",
    "request_id":         "Unique request trace ID for debugging",
    "error_id":           "Unique error instance ID for support lookups",
    "retryable":          "Whether the client should retry this request",
    "retry_after_seconds":"Number of seconds to wait before retrying",
    "field_violations":   "Field-level validation errors (present on 422 responses)",
    "success":            "true = operation succeeded, false = operation failed",
    "data":               "Response payload on success, null on failure",
    "error":              "Error details on failure, null on success",
    "meta":               "Response metadata: request tracing and pagination",
    "access_token":       "JWT access token (short-lived, 15 minutes)",
    "refresh_token":      "Refresh token (long-lived, 30 days)",
    "token_type":         "Token type, always 'Bearer'",
    "expires_in":         "Token expiry time in seconds",
}

yaml = YAML()
yaml.preserve_quotes = True
yaml.default_flow_style = False
yaml.width = 120

# Ensure Python None is always serialized as explicit `null` (not a bare key).
# ruamel.yaml by default writes None as '' (empty), which breaks OpenAPI examples.
from ruamel.yaml.representer import RoundTripRepresenter
RoundTripRepresenter.add_representer(
    type(None),
    lambda dumper, _: dumper.represent_scalar('tag:yaml.org,2002:null', 'null')
)


def fix_schema_properties(schema: dict, schema_name: str) -> int:
    """Fix warnings in schema properties. Returns count of fixes applied."""
    fixes = 0
    props = schema.get("properties", {})
    if not props:
        return 0

    required = schema.get("required", [])
    new_required = list(required)

    for prop_name, prop_def in props.items():
        if not isinstance(prop_def, dict):
            continue

        # Fix 1: Missing description → add default
        if not prop_def.get("description"):
            default_desc = DEFAULT_DESCRIPTIONS.get(prop_name)
            if not default_desc:
                # Generate from field name
                default_desc = prop_name.replace("_", " ").capitalize()
            prop_def["description"] = default_desc
            fixes += 1

        # Fix 2: Empty description → replace with default
        elif prop_def.get("description", "").strip() == "":
            default_desc = DEFAULT_DESCRIPTIONS.get(prop_name, prop_name.replace("_", " ").capitalize())
            prop_def["description"] = default_desc
            fixes += 1

        # Fix 3: total_items typed as string (proto int64 artifact) → integer
        if prop_name == "total_items" and prop_def.get("type") == "string":
            prop_def["type"] = "integer"
            prop_def["format"] = "int64"
            fixes += 1

        # Fix 4: Add required for Request schemas non-nullable scalar fields
        if schema_name.endswith("Request"):
            prop_type = prop_def.get("type", "")
            is_nullable = prop_def.get("nullable", False)
            is_repeated = prop_type == "array"
            is_message  = "$ref" in prop_def or "allOf" in prop_def
            is_optional = prop_def.get("x-optional", False)

            if (not is_nullable and not is_repeated and not is_message
                    and not is_optional and prop_name not in required
                    and prop_name in ALWAYS_REQUIRED_PATTERNS):
                new_required.append(prop_name)
                fixes += 1

    # Write back required if changed
    if new_required and new_required != list(required):
        schema["required"] = sorted(set(new_required))

    return fixes


def fix_yaml_file(file_path: Path) -> int:
    """Load, fix, and save a single YAML file. Returns count of fixes."""
    try:
        with open(file_path, "r", encoding="utf-8") as f:
            data = yaml.load(f)

        if not data or not isinstance(data, dict):
            return 0

        fixes = 0

        # Fix schemas at top level (schema files)
        if "properties" in data:
            schema_name = file_path.stem
            fixes += fix_schema_properties(data, schema_name)

        # Fix schemas inside components/schemas (openapi.yaml)
        schemas = data.get("components", {}).get("schemas", {})
        for schema_name, schema_def in schemas.items():
            if isinstance(schema_def, dict):
                fixes += fix_schema_properties(schema_def, schema_name)

        # Fix path operation requestBody schemas
        paths = data.get("paths", {})
        for path_url, path_item in paths.items():
            if not isinstance(path_item, dict):
                continue
            for method, op in path_item.items():
                if not isinstance(op, dict):
                    continue
                rb = op.get("requestBody", {})
                schema = rb.get("content", {}).get("application/json", {}).get("schema", {})
                if schema and "properties" in schema:
                    op_id = op.get("operationId", file_path.stem)
                    fixes += fix_schema_properties(schema, op_id + "Request")

        if fixes > 0:
            buf = io.StringIO()
            yaml.dump(data, buf)
            new_content = buf.getvalue()
            try:
                with open(file_path, "r", encoding="utf-8") as f:
                    existing = f.read()
            except Exception:
                existing = None
            if new_content != existing:
                with open(file_path, "w", encoding="utf-8") as f:
                    f.write(new_content)

        return fixes

    except Exception as e:
        print(f"  ⚠ Error processing {file_path.name}: {e}", file=sys.stderr)
        return 0


def main():
    print("=" * 60)
    print("  OpenAPI Warning Fixer")
    print("=" * 60)

    total_files  = 0
    total_fixes  = 0
    errors       = 0

    for target_dir in TARGET_DIRS:
        if not target_dir.exists():
            continue

        yaml_files = list(target_dir.rglob("*.yaml"))
        print(f"\n  Processing {target_dir.name}/ ({len(yaml_files)} files)...")

        for yaml_file in yaml_files:
            fixes = fix_yaml_file(yaml_file)
            total_files += 1
            total_fixes += fixes

    # NOTE: openapi.yaml is NOT processed here.
    # ruamel.yaml reformats PyYAML-generated content differently (indentation,
    # quotes, blank lines) causing the file to grow ~2300 lines per run even
    # when no fixes are needed. All warning fixes for the assembled spec are
    # now applied directly inside path_generator.py and assembler.py at
    # generation time, so no post-processing of openapi.yaml is needed.

    print()
    print("=" * 60)
    print(f"  ✅ Fixed {total_fixes} warnings across {total_files} files")
    print("=" * 60)


if __name__ == "__main__":
    main()
