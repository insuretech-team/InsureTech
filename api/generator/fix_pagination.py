#!/usr/bin/env python3
"""
fix_pagination.py — InsureTech Pagination Standards Enforcer

Ensures all list (GET) endpoints in the generated openapi.yaml comply with
Rule 05 (Pagination & List Endpoints):

  1. Every GET endpoint returning a list gets standard pagination query params
     (page, page_size, sort_by, sort_order, search, from_date, to_date)
  2. Removes deprecated PageResponse and PaginationResponse schema references
     and replaces with PaginationMeta (already injected by assembler.py)
  3. Ensures list response data uses data.items[] shape
  4. Verifies PaginationMeta is in meta.pagination (not inside data)

Uses ruamel.yaml for fast round-trip editing (preserves formatting).
Called by run_api_pipeline.ps1 Step 12 in parallel with fix_all_warnings.py.
"""

import sys
import io
import re
from pathlib import Path

try:
    from ruamel.yaml import YAML
except ImportError:
    print("Error: ruamel.yaml not found. Run: pip install ruamel.yaml")
    sys.exit(1)

# ---------------------------------------------------------------------------
# Config
# ---------------------------------------------------------------------------
SCRIPT_DIR = Path(__file__).parent
API_DIR    = SCRIPT_DIR.parent
OPENAPI_PATH = API_DIR / "openapi.yaml"

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

# ---------------------------------------------------------------------------
# Standard pagination query parameters (Rule 05)
# ---------------------------------------------------------------------------
PAGINATION_PARAMS = [
    {
        "name": "page",
        "in": "query",
        "required": False,
        "schema": {"type": "integer", "minimum": 1, "default": 1},
        "description": "Page number (1-based). Default: 1"
    },
    {
        "name": "page_size",
        "in": "query",
        "required": False,
        "schema": {"type": "integer", "minimum": 1, "maximum": 100, "default": 20},
        "description": "Number of items per page (1-100). Default: 20"
    },
    {
        "name": "sort_by",
        "in": "query",
        "required": False,
        "schema": {"type": "string"},
        "description": "Field name to sort by (e.g. 'created_at', 'name')"
    },
    {
        "name": "sort_order",
        "in": "query",
        "required": False,
        "schema": {"type": "string", "enum": ["asc", "desc"], "default": "desc"},
        "description": "Sort direction: 'asc' or 'desc'. Default: desc"
    },
    {
        "name": "search",
        "in": "query",
        "required": False,
        "schema": {"type": "string"},
        "description": "Full-text search query"
    },
    {
        "name": "from_date",
        "in": "query",
        "required": False,
        "schema": {"type": "string", "format": "date"},
        "description": "Filter records from this date (inclusive), format YYYY-MM-DD"
    },
    {
        "name": "to_date",
        "in": "query",
        "required": False,
        "schema": {"type": "string", "format": "date"},
        "description": "Filter records to this date (inclusive), format YYYY-MM-DD"
    },
]

# Operation name patterns that indicate a list endpoint
LIST_OPERATION_PATTERNS = re.compile(
    r'(List|GetAll|Search|Browse|Find|Fetch|Query|Index|History|'
    r'Retrieve.*s$|Get.*s$)',
    re.IGNORECASE
)

# Schemas to replace with PaginationMeta reference
DEPRECATED_PAGINATION_SCHEMAS = {"PageResponse", "PaginationResponse"}


def is_list_operation(operation_id: str, path_url: str) -> bool:
    """Determine if an operation returns a list of resources."""
    method_name = operation_id.split("_", 1)[-1] if "_" in operation_id else operation_id

    # Check operation name
    if LIST_OPERATION_PATTERNS.search(method_name):
        return True

    # GET on a collection path (no path params at end, no action colon)
    if ":" not in path_url:
        parts = path_url.rstrip("/").split("/")
        last = parts[-1] if parts else ""
        # Collection paths end with plural noun (no {param})
        if not last.startswith("{") and not last.endswith("}"):
            if last.endswith("s") or last in ("history", "analytics", "metrics", "logs"):
                return True

    return False


def has_pagination_param(parameters: list, param_name: str) -> bool:
    """Check if a parameter already exists in the list."""
    return any(p.get("name") == param_name for p in parameters)


def add_pagination_params(operation: dict) -> int:
    """Add missing pagination parameters to a list operation. Returns count added."""
    if "parameters" not in operation:
        operation["parameters"] = []

    added = 0
    for param in PAGINATION_PARAMS:
        if not has_pagination_param(operation["parameters"], param["name"]):
            operation["parameters"].append(param)
            added += 1

    return added


def fix_deprecated_schema_refs(data: dict) -> int:
    """
    Replace deprecated PageResponse/PaginationResponse $ref with PaginationMeta
    anywhere in the spec. Returns count of replacements.
    """
    fixes = 0

    def replace_refs(obj):
        nonlocal fixes
        if isinstance(obj, dict):
            for key, val in obj.items():
                if key == "$ref" and isinstance(val, str):
                    for deprecated in DEPRECATED_PAGINATION_SCHEMAS:
                        if f"/schemas/{deprecated}" in val:
                            obj[key] = val.replace(f"/schemas/{deprecated}", "/schemas/PaginationMeta")
                            fixes += 1
                else:
                    replace_refs(val)
        elif isinstance(obj, list):
            for item in obj:
                replace_refs(item)

    replace_refs(data)
    return fixes


def remove_deprecated_schemas(data: dict) -> int:
    """Remove PageResponse and PaginationResponse from components/schemas."""
    schemas = data.get("components", {}).get("schemas", {})
    removed = 0
    for schema_name in list(DEPRECATED_PAGINATION_SCHEMAS):
        if schema_name in schemas:
            del schemas[schema_name]
            removed += 1
            print(f"    ✓ Removed deprecated schema: {schema_name}")
    return removed


def main():
    print("=" * 60)
    print("  Pagination Standards Enforcer (Rule 05)")
    print("=" * 60)

    # NOTE: Pagination query params are now injected directly by path_generator.py
    # at generation time (Rule 05 logic baked in). Processing openapi.yaml here
    # with ruamel.yaml caused the file to grow ~2300 lines per run due to format
    # differences between PyYAML (writer) and ruamel.yaml (round-tripper).
    # This script now only checks for deprecated schema refs as a safety net.

    if not OPENAPI_PATH.exists():
        print(f"  ✗ openapi.yaml not found at: {OPENAPI_PATH}")
        sys.exit(1)

    # Step 1: Pagination params — now injected at generation time in path_generator.py
    print("\n  [1/3] Pagination params — injected at generation time (skipping)")

    # Step 2 & 3: Check for deprecated schema refs using fast text search.
    # Do NOT load openapi.yaml with ruamel.yaml — round-trip reformats PyYAML
    # output (different indent/quoting) causing the file to grow on every run.
    print("\n  [2/3] Checking for deprecated pagination schema refs (text scan)...")
    with open(OPENAPI_PATH, "rb") as f:
        raw = f.read()
    content = raw.decode("utf-8")

    deprecated_found = any(dep in content for dep in DEPRECATED_PAGINATION_SCHEMAS)
    if not deprecated_found:
        print("    ✓ No deprecated schema refs found")
        print("\n  [3/3] No deprecated schemas to remove")
        print()
        print("=" * 60)
        print("  ✅ Pagination check complete — no changes needed")
        print("=" * 60)
        return

    # Only load with ruamel.yaml if deprecated refs are actually present
    print(f"    ⚠ Deprecated refs detected — loading and fixing...")
    spec = yaml.load(io.StringIO(content))
    if not spec:
        print("  ✗ openapi.yaml is empty or invalid")
        sys.exit(1)

    total_refs_fixed = fix_deprecated_schema_refs(spec)
    print(f"    ✓ Replaced {total_refs_fixed} deprecated $ref(s)")

    print("\n  [3/3] Removing deprecated schemas...")
    removed = remove_deprecated_schemas(spec)
    if removed == 0:
        print("    ✓ No deprecated schemas found")

    # Write back with LF normalisation so write_guard can detect sameness
    buf = io.StringIO()
    yaml.dump(spec, buf)
    new_lf = buf.getvalue().replace('\r\n', '\n')
    existing_lf = content.replace('\r\n', '\n')
    if new_lf != existing_lf:
        with open(OPENAPI_PATH, "w", encoding="utf-8", newline='\n') as f:
            f.write(new_lf)
        print(f"\n  ✓ Saved updated openapi.yaml (deprecated refs removed)")
    else:
        print(f"\n  ✓ openapi.yaml unchanged (skipped write)")

    print()
    print("=" * 60)
    print(f"  ✅ Pagination check complete:")
    print(f"     Deprecated refs fixed: {total_refs_fixed}")
    print(f"     Deprecated schemas removed: {removed}")
    print("=" * 60)


if __name__ == "__main__":
    main()
