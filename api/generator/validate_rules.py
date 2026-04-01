#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
validate_rules.py — InsureTech API Rule Validator

Validates that the generated openapi.yaml complies with all rules defined in
E:/Projects/InsureTech/rules/

Usage:
    python validate_rules.py [openapi.yaml path]
    python validate_rules.py ../openapi.yaml

Exit codes:
    0 = all rules pass
    1 = violations found (use in CI to fail the build)
"""

import sys
import yaml

# ---------------------------------------------------------------------------
# Rule 04: Public endpoints — must declare security: []
# ---------------------------------------------------------------------------
PUBLIC_OPERATIONS = {
    'AuthService_Register',
    'AuthService_Login',
    'AuthService_SendOTP',
    'AuthService_VerifyOTP',
    'AuthService_ResendOTP',
    'AuthService_ValidateCSRF',
    'AuthService_RegisterEmailUser',
    'AuthService_EmailLogin',
    'AuthService_SendEmailOTP',
    'AuthService_VerifyEmail',
    'AuthService_RequestPasswordResetByEmail',
    'AuthService_ResetPasswordByEmail',
    'AuthService_ResetPassword',
    'AuthService_GetJWKS',
    'AuthService_BiometricAuthenticate',
    'ProductService_ListProducts',
    'ProductService_GetProduct',
    'ProductService_SearchProducts',
    'ProductService_CalculatePremium',
}

# ---------------------------------------------------------------------------
# Rule 02: Verbs that create a resource → must return 201
# ---------------------------------------------------------------------------
RESOURCE_CREATION_VERBS = (
    'Create', 'Register', 'Submit', 'Initiate', 'Request',
    'Add', 'Upload', 'Start', 'Open', 'File',
)

def is_creation_method(method_name: str) -> bool:
    for verb in RESOURCE_CREATION_VERBS:
        if method_name.startswith(verb):
            return True
    return False


def validate_rules(spec_path: str) -> int:
    """
    Validate all API rules against the given OpenAPI spec.
    Returns number of violations found.
    """
    print(f"\n{'='*60}")
    print(f"  InsureTech API Rule Validator")
    print(f"  Spec: {spec_path}")
    print(f"{'='*60}\n")

    with open(spec_path, encoding='utf-8') as f:
        spec = yaml.safe_load(f)

    paths = spec.get('paths', {})
    schemas = spec.get('components', {}).get('schemas', {})
    security_schemes = spec.get('components', {}).get('securitySchemes', {})

    errors = []
    warnings = []
    passed = []

    # -----------------------------------------------------------------------
    # Rule 01: Canonical schemas must exist
    # -----------------------------------------------------------------------
    for canonical in ('ApiResponse', 'ResponseMeta', 'PaginationMeta'):
        if canonical in schemas:
            passed.append(f"Rule 01: Canonical schema '{canonical}' present ✓")
        else:
            errors.append(f"Rule 01: Missing canonical schema '{canonical}'")

    # -----------------------------------------------------------------------
    # Rule 04: BearerAuth security scheme must be defined
    # -----------------------------------------------------------------------
    if 'BearerAuth' in security_schemes:
        passed.append("Rule 04: BearerAuth security scheme defined ✓")
    else:
        errors.append("Rule 04: Missing BearerAuth in securitySchemes")

    # -----------------------------------------------------------------------
    # Rule 05: PaginationMeta must exist, PaginationResponse must NOT exist
    # -----------------------------------------------------------------------
    if 'PaginationMeta' in schemas:
        passed.append("Rule 05: PaginationMeta schema present ✓")
    else:
        errors.append("Rule 05: Missing PaginationMeta schema")

    if 'PaginationResponse' in schemas:
        errors.append("Rule 05: Duplicate PaginationResponse schema must be removed (use PaginationMeta)")
    else:
        passed.append("Rule 05: PaginationResponse removed ✓")

    # Rule 05: PageResponse.total_items must be integer not string
    if 'PageResponse' in schemas:
        pr_props = schemas['PageResponse'].get('properties', {})
        ti = pr_props.get('total_items', {})
        if ti.get('type') == 'integer':
            passed.append("Rule 05: PageResponse.total_items is integer ✓")
        else:
            errors.append(f"Rule 05: PageResponse.total_items type='{ti.get('type')}' must be integer")

    # -----------------------------------------------------------------------
    # Per-endpoint checks
    # -----------------------------------------------------------------------
    rule02_errors = []
    rule02_passed = 0
    rule03_errors = []
    rule03_passed = 0
    rule04_errors = []
    rule04_passed = 0

    for path_url, path_item in paths.items():
        for method, op in path_item.items():
            if method not in ('get', 'post', 'put', 'patch', 'delete'):
                continue

            op_id = op.get('operationId', '')
            method_name = op_id.split('_', 1)[-1] if '_' in op_id else op_id
            responses = op.get('responses', {})
            security = op.get('security')  # None = not declared, [] = public, [...] = protected

            # ----------------------------------------------------------------
            # Rule 02: Correct status codes
            # ----------------------------------------------------------------
            # POST to collection path (no colon action) with creation verb → must have 201
            if method == 'post' and ':' not in path_url:
                if is_creation_method(method_name):
                    if '201' in responses:
                        rule02_passed += 1
                    else:
                        rule02_errors.append(
                            f"  [{op_id}] POST {path_url} should return 201 (creation verb: {method_name})"
                        )

            # PUT/PATCH must have 200 not 201
            if method in ('put', 'patch') and '201' in responses:
                rule02_errors.append(
                    f"  [{op_id}] {method.upper()} {path_url} should return 200 not 201"
                )

            # DELETE must have 204
            if method == 'delete' and '204' not in responses:
                rule02_errors.append(
                    f"  [{op_id}] DELETE {path_url} missing 204 response"
                )

            # ----------------------------------------------------------------
            # Rule 03: Required error responses
            # ----------------------------------------------------------------
            # 400 must always be present
            if '400' not in responses:
                rule03_errors.append(f"  [{op_id}] {method.upper()} {path_url} missing 400")

            # 422 must be present for all write endpoints
            if method in ('post', 'put', 'patch') and '422' not in responses:
                rule03_errors.append(f"  [{op_id}] {method.upper()} {path_url} missing 422")
            else:
                if method in ('post', 'put', 'patch'):
                    rule03_passed += 1

            # 401/403 must be present for authenticated endpoints
            if op_id not in PUBLIC_OPERATIONS:
                if '401' not in responses:
                    rule03_errors.append(f"  [{op_id}] {method.upper()} {path_url} missing 401")
                if '403' not in responses:
                    rule03_errors.append(f"  [{op_id}] {method.upper()} {path_url} missing 403")

            # 201 must include Location header
            if '201' in responses:
                headers = responses['201'].get('headers', {})
                if 'Location' not in headers:
                    rule02_errors.append(
                        f"  [{op_id}] POST {path_url} returns 201 but missing Location header"
                    )

            # 500 must always be present
            if '500' not in responses:
                rule03_errors.append(f"  [{op_id}] {method.upper()} {path_url} missing 500")

            # ----------------------------------------------------------------
            # Rule 04: Security must be explicitly declared on every operation
            # ----------------------------------------------------------------
            if security is None:
                rule04_errors.append(
                    f"  [{op_id}] {method.upper()} {path_url} has no security declaration"
                )
            else:
                rule04_passed += 1

    # -----------------------------------------------------------------------
    # Rule 08: No 'error' field inside *Response schemas
    # Exclude canonical ApiResponse — it intentionally has an error field
    # as part of the standard envelope definition.
    # -----------------------------------------------------------------------
    CANONICAL_SCHEMA_NAMES = {'ApiResponse', 'ResponseMeta', 'PaginationMeta'}
    rule08_errors = []
    rule08_passed = 0
    for schema_name, schema_def in schemas.items():
        if schema_name.endswith('Response') and schema_name not in CANONICAL_SCHEMA_NAMES:
            props = schema_def.get('properties', {})
            if 'error' in props:
                rule08_errors.append(
                    f"  Schema '{schema_name}' has 'error' field — must be removed (Rule 01/03)"
                )
            else:
                rule08_passed += 1

    # -----------------------------------------------------------------------
    # Print results
    # -----------------------------------------------------------------------
    total_errors = (
        len(errors) + len(rule02_errors) + len(rule03_errors) +
        len(rule04_errors) + len(rule08_errors)
    )

    if rule02_errors:
        print(f"❌ Rule 02 — HTTP Status Codes ({len(rule02_errors)} violations):")
        for e in rule02_errors: print(e)
        print()
    else:
        print(f"✅ Rule 02 — HTTP Status Codes: {rule02_passed} checks passed")

    if rule03_errors:
        print(f"\n❌ Rule 03 — Error Responses ({len(rule03_errors)} violations):")
        for e in rule03_errors: print(e)
        print()
    else:
        print(f"✅ Rule 03 — Error Responses: {rule03_passed} write endpoints have 422")

    if rule04_errors:
        print(f"\n❌ Rule 04 — Security Declarations ({len(rule04_errors)} violations):")
        for e in rule04_errors[:20]: print(e)
        if len(rule04_errors) > 20:
            print(f"  ... and {len(rule04_errors)-20} more")
        print()
    else:
        print(f"✅ Rule 04 — Security: {rule04_passed} endpoints have security declared")

    if rule08_errors:
        print(f"\n❌ Rule 08 — Error in Response Schemas ({len(rule08_errors)} violations):")
        for e in rule08_errors[:20]: print(e)
        if len(rule08_errors) > 20:
            print(f"  ... and {len(rule08_errors)-20} more")
        print()
    else:
        print(f"✅ Rule 08 — Response Schemas: {rule08_passed} schemas have no embedded error field")

    if errors:
        print(f"\n❌ General Rule Violations ({len(errors)}):")
        for e in errors: print(f"  {e}")
        print()
    else:
        for p in passed:
            print(f"✅ {p}")

    print(f"\n{'='*60}")
    if total_errors == 0:
        print(f"  ✅ ALL RULES PASSED — {len(passed)} checks validated")
    else:
        print(f"  ❌ {total_errors} RULE VIOLATION(S) FOUND")
    print(f"{'='*60}\n")

    return total_errors


if __name__ == '__main__':
    # Force UTF-8 output on Windows (avoids CP1252 UnicodeEncodeError with emoji)
    if sys.stdout.encoding and sys.stdout.encoding.lower() != 'utf-8':
        sys.stdout.reconfigure(encoding='utf-8', errors='replace')
    spec_path = sys.argv[1] if len(sys.argv) > 1 else '../openapi.yaml'
    violations = validate_rules(spec_path)
    sys.exit(1 if violations > 0 else 0)
