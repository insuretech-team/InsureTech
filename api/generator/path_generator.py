import os
import re
import yaml

# Import our custom modules
try:
    from name_transformer import NameTransformer
    from description_loader import DescriptionLoader
except ImportError:
    NameTransformer = None
    DescriptionLoader = None


# ---------------------------------------------------------------------------
# Example value generator
# Produces realistic example values for OpenAPI schema properties so that
# success responses show actual DTO shape instead of a bare `data: {}`.
# ---------------------------------------------------------------------------

# Field-name → example value (checked before type-based fallbacks)
_FIELD_EXAMPLES = {
    # IDs
    "id":                       "01HGW2PAYMENT0000000000001",
    "user_id":                  "01HGW2USER00000000000000001",
    "session_id":               "01HGW2SESSION000000000000001",
    "policy_id":                "01HGW2POLICY000000000000001",
    "claim_id":                 "01HGW2CLAIM0000000000000001",
    "payment_id":               "01HGW2PAY00000000000000001",
    "order_id":                 "01HGW2ORDER0000000000000001",
    "product_id":               "01HGW2PROD00000000000000001",
    "tenant_id":                "01HGW2TENANT000000000000001",
    "partner_id":               "01HGW2PARTNER00000000000001",
    "quote_id":                 "01HGW2QUOTE0000000000000001",
    "insurer_id":               "01HGW2INSURER00000000000001",
    "document_id":              "01HGW2DOC00000000000000001",
    "invoice_id":               "01HGW2INV00000000000000001",
    "report_id":                "01HGW2RPT00000000000000001",
    "analysis_id":              "01HGW2ANALYSIS0000000000001",
    "api_key_id":               "01HGW2APIKEY000000000000001",
    "task_id":                  "01HGW2TASK00000000000000001",
    "workflow_id":              "01HGW2WF000000000000000001",
    "device_id":                "01HGW2DEVICE000000000000001",
    # Auth / tokens
    "access_token":             "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.example",
    "refresh_token":            "rt_01HGW2REFRESH00000000000001",
    "session_token":            "st_01HGW2SESSION00000000000001",
    "csrf_token":               "csrf_01HGW2CSRF000000000000001",
    "api_key":                  "sk_live_01HGW2APIKEY00000000000001",
    "new_api_key":              "sk_live_01HGW2NEWKEY00000000000001",
    "mfa_session_token":        "mfa_01HGW2MFA000000000000001",
    "token_type":               "Bearer",
    "access_token_expires_in":  900,
    "refresh_token_expires_in": 2592000,
    "expires_in":               900,
    # User / profile
    "email":                    "user@example.com",
    "phone_number":             "+8801712345678",
    "name":                     "John Doe",
    "first_name":               "John",
    "last_name":                "Doe",
    "username":                 "johndoe",
    "avatar_url":               "https://cdn.example.com/avatars/user_01.jpg",
    "session_type":             "JWT",
    "mfa_required":             False,
    "mfa_method":               "TOTP",
    # Scores / confidence
    "confidence":               0.92,
    "confidence_score":         0.87,
    "fraud_score":              12.5,
    "risk_score":               35.0,
    "risk_category":            "LOW",
    "recommendation":           "APPROVE",
    "is_suspicious":            False,
    "verification_passed":      True,
    "conversation_ended":       False,
    # Counts / pagination
    "total_count":              42,
    "row_count":                10,
    "execution_time_ms":        23.4,
    # Timestamps
    "created_at":               "2024-01-15T10:30:00Z",
    "updated_at":               "2024-01-15T10:30:00Z",
    "deleted_at":               None,
    "next_run_at":              "2024-02-01T00:00:00Z",
    "issued_at":                "2024-01-15T10:30:00Z",
    "expires_at":               "2024-01-16T10:30:00Z",
    # Status / type
    "status":                   "ACTIVE",
    "type":                     "INDIVIDUAL",
    "currency":                 "BDT",
    "amount":                   "5000.00",
    "valid":                    True,
    "owner_type":               "USER",
    "owner_id":                 "01HGW2USER00000000000000001",
    # Other
    "conversation_id":          "01HGW2CONV00000000000000001",
    "message":                  "Operation completed successfully",
    "description":              "Example description",
    "url":                      "https://api.labaidinsuretech.com/v1/resource/01HGW2ID",
    "download_url":             "https://cdn.labaidinsuretech.com/docs/example.pdf",
}

# Type/format → example value fallbacks
_TYPE_EXAMPLES = {
    ("string",  None):          "example_value",
    ("string",  "date-time"):   "2024-01-15T10:30:00Z",
    ("string",  "date"):        "2024-01-15",
    ("string",  "uuid"):        "550e8400-e29b-41d4-a716-446655440000",
    ("string",  "uri"):         "https://example.com",
    ("string",  "email"):       "user@example.com",
    ("string",  "password"):    "••••••••",
    ("integer", None):          1,
    ("integer", "int32"):       1,
    ("integer", "int64"):       1,
    ("number",  None):          1.0,
    ("number",  "double"):      1.0,
    ("number",  "float"):       1.0,
    ("boolean", None):          True,
    ("object",  None):          {},
    ("array",   None):          [],
}

# Suffix patterns → example value (for ID-like fields)
_SUFFIX_EXAMPLES = [
    ("_id",       "01HGW2EXAMPLE0000000000001"),
    ("_at",       "2024-01-15T10:30:00Z"),
    ("_url",      "https://example.com/resource"),
    ("_token",    "tok_example_000000000001"),
    ("_key",      "key_example_000000000001"),
    ("_count",    10),
    ("_score",    0.85),
    ("_amount",   "1000.00"),
    ("_code",     "EXAMPLE_CODE"),
    ("_type",     "EXAMPLE_TYPE"),
    ("_status",   "ACTIVE"),
    ("_name",     "Example Name"),
    ("_email",    "user@example.com"),
    ("_phone",    "+8801712345678"),
]


def _example_value_for_prop(prop_name: str, prop_def: dict, depth: int = 0) -> object:
    """Return a realistic example value for a single schema property."""
    # Direct name match first
    if prop_name in _FIELD_EXAMPLES:
        return _FIELD_EXAMPLES[prop_name]

    # Suffix match
    for suffix, val in _SUFFIX_EXAMPLES:
        if prop_name.endswith(suffix):
            return val

    prop_type   = prop_def.get("type")
    prop_format = prop_def.get("format")
    prop_ref    = prop_def.get("$ref")

    # $ref → return typed placeholder (we can't resolve at path-gen time)
    if prop_ref:
        schema_name = prop_ref.split("/")[-1] if prop_ref else ""
        # Money object
        if "Money" in schema_name:
            return {"amount": "1000.00", "currency": "BDT"}
        # User object
        if schema_name == "User":
            return {"id": "01HGW2USER00000000000000001", "email": "user@example.com",
                    "phone_number": "+8801712345678", "status": "ACTIVE"}
        return {}

    # allOf → use first schema
    if "allOf" in prop_def:
        return {}

    # Array type
    if prop_type == "array":
        items = prop_def.get("items", {})
        items_type = items.get("type")
        items_ref  = items.get("$ref")
        if items_ref:
            schema_name = items_ref.split("/")[-1]
            if "Money" in schema_name:
                return [{"amount": "1000.00", "currency": "BDT"}]
            return [{}]
        if items_type == "string":
            return ["example_item"]
        if items_type == "integer":
            return [1]
        if items_type == "number":
            return [1.0]
        if items_type == "boolean":
            return [True]
        return []

    # Object type with additionalProperties
    if prop_type == "object":
        if prop_def.get("additionalProperties"):
            return {"key": "value"}
        return {}

    # Type + format lookup
    key = (prop_type, prop_format)
    if key in _TYPE_EXAMPLES:
        return _TYPE_EXAMPLES[key]

    # Type only
    key2 = (prop_type, None)
    if key2 in _TYPE_EXAMPLES:
        return _TYPE_EXAMPLES[key2]

    return "example"


def _build_data_example(schema_ref: dict, schemas_dir: str = None) -> dict:
    """
    Build a realistic `data` example object from a schema $ref.
    Reads the corresponding YAML file from the schemas directory and walks
    its properties to produce field-level example values.

    Falls back to {} if the schema cannot be resolved.
    """
    if not schema_ref:
        return {}

    ref = schema_ref.get("$ref", "")
    if not ref:
        return {}

    # Extract schema name: '#/components/schemas/LoginResponse' → 'LoginResponse'
    schema_name = ref.split("/")[-1]

    if not schemas_dir or not os.path.isdir(schemas_dir):
        return {}

    # Search for schema YAML file by name
    schema_file = None
    for root, dirs, files in os.walk(schemas_dir):
        dirs.sort()
        for fname in files:
            if fname == f"{schema_name}.yaml":
                schema_file = os.path.join(root, fname)
                break
        if schema_file:
            break

    if not schema_file:
        return {}

    try:
        with open(schema_file, "r", encoding="utf-8") as f:
            data = yaml.safe_load(f)
    except Exception:
        return {}

    if not data or not isinstance(data, dict):
        return {}

    # Schema file format: { SchemaName: { type: object, properties: {...} } }
    schema_def = data.get(schema_name, data)
    props = schema_def.get("properties", {})
    if not props:
        return {}

    example = {}
    for prop_name, prop_def in props.items():
        if not isinstance(prop_def, dict):
            continue
        example[prop_name] = _example_value_for_prop(prop_name, prop_def)

    return example

# ---------------------------------------------------------------------------
# Rule 05: Pagination — standard query params injected at generation time
# so fix_pagination.py does NOT need to post-process openapi.yaml.
# ---------------------------------------------------------------------------
_LIST_OPERATION_RE = re.compile(
    r'(List|GetAll|Search|Browse|Find|Fetch|Query|Index|History|'
    r'Retrieve.*s$|Get.*s$)',
    re.IGNORECASE
)

_PAGINATION_PARAMS = [
    {"name": "page",       "in": "query", "required": False,
     "schema": {"type": "integer", "minimum": 1, "default": 1},
     "description": "Page number (1-based). Default: 1"},
    {"name": "page_size",  "in": "query", "required": False,
     "schema": {"type": "integer", "minimum": 1, "maximum": 100, "default": 20},
     "description": "Number of items per page (1-100). Default: 20"},
    {"name": "sort_by",    "in": "query", "required": False,
     "schema": {"type": "string"},
     "description": "Field name to sort by (e.g. 'created_at', 'name')"},
    {"name": "sort_order", "in": "query", "required": False,
     "schema": {"type": "string", "enum": ["asc", "desc"], "default": "desc"},
     "description": "Sort direction: 'asc' or 'desc'. Default: desc"},
    {"name": "search",     "in": "query", "required": False,
     "schema": {"type": "string"},
     "description": "Full-text search query"},
    {"name": "from_date",  "in": "query", "required": False,
     "schema": {"type": "string", "format": "date"},
     "description": "Filter records from this date (inclusive), format YYYY-MM-DD"},
    {"name": "to_date",    "in": "query", "required": False,
     "schema": {"type": "string", "format": "date"},
     "description": "Filter records to this date (inclusive), format YYYY-MM-DD"},
]


def _is_list_operation(method_name: str, path_url: str) -> bool:
    """Return True if this GET operation returns a collection (list)."""
    if _LIST_OPERATION_RE.search(method_name):
        return True
    # GET on a collection path (no path params at end, no action colon)
    if ':' not in path_url:
        parts = path_url.rstrip('/').split('/')
        last = parts[-1] if parts else ''
        if not last.startswith('{') and not last.endswith('}'):
            if last.endswith('s') or last in ('history', 'analytics', 'metrics', 'logs'):
                return True
    return False


# ---------------------------------------------------------------------------
# Rule 02: Status Code Classification
# Verbs that CREATE a new persistent resource → 201 Created
# ---------------------------------------------------------------------------
RESOURCE_CREATION_VERBS = (
    'Create', 'Register', 'Submit', 'Initiate', 'Request',
    'Add', 'Upload', 'Start', 'Open', 'File',
)

# ---------------------------------------------------------------------------
# Rule 04: Security Classification
# Endpoints that require NO authentication (security: [])
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

# Webhook endpoints — authenticated via HMAC signature, not bearer token
WEBHOOK_OPERATIONS = {
    'PaymentService_HandleGatewayWebhook',
    'MFSService_ProcessWebhook',
}

class PathGenerator:
    def __init__(self, registry, descriptions_dir=None, schemas_dir=None):
        self.registry = registry
        self.name_transformer = NameTransformer() if NameTransformer else None
        self.description_loader = DescriptionLoader(descriptions_dir) if DescriptionLoader and descriptions_dir else None
        # schemas_dir is used by the example generator to walk DTO YAML files
        self.schemas_dir = schemas_dir

    def generate_path_item(self, method_data, service_name):
        """
        Generates an OpenAPI Path Item for a single method.
        Returns: (path_url, verb, path_item_dict)
        """
        http_rule = method_data.get('http_rule')
        if not http_rule:
            return None, None, None

        # Determine Verb and Path
        verb, full_path_url = self._extract_verb_and_path(http_rule)
        if not verb:
            return None, None, None

        # Convert custom actions (:cancel) to query parameters (?action=cancel)
        path_url, action_param = self._process_custom_actions(full_path_url)
        
        # Fix kebab-case in path
        path_url = self._to_kebab_case_path(path_url)
        
        # Split query params from URL pattern
        # Special case: ?action=<verb> → convert to colon-suffix path per Google API Design Guide
        # e.g. /v1/policies/{id}?action=renew → /v1/policies/{id}:renew
        # This ensures each action gets a unique path (no collision between approve/reject/etc.)
        if '?' in path_url:
            path_url, query_string = path_url.split('?', 1)
            action_value = None
            remaining_params = []
            for pair in query_string.split('&'):
                if '=' in pair:
                    k, v = pair.split('=', 1)
                    if k == 'action':
                        action_value = v
                    else:
                        remaining_params.append((k, v))
            if action_value:
                # Convert ?action=approve → :approve suffix on path
                kebab_action = self.name_transformer._to_kebab_case(action_value) if self.name_transformer else action_value
                path_url = f"{path_url}:{kebab_action}"
            query_string = '&'.join(f"{k}={v}" for k, v in remaining_params)
        else:
            query_string = ""

        # Extract Parameters from Path
        parameters = self._extract_parameters(path_url)

        # Add remaining (non-action) query params
        if query_string:
            for pair in query_string.split('&'):
                if '=' in pair:
                    k, v = pair.split('=', 1)
                    parameters.append({
                        "name": k,
                        "in": "query",
                        "required": True,
                        "schema": {"type": "string", "enum": [v]}
                    })
        
        # Load operation description
        method_name = method_data.get('name')
        method_comment = method_data.get('comment', '')
        if self.description_loader:
            desc_data = self.description_loader.load_operation_description(
                service_name, method_name, proto_comment=method_comment
            )
            summary = desc_data['summary']
            description = desc_data.get('description', '')
        else:
            # Fallback to proto comment if available
            if method_comment:
                summary = method_comment.split('\n')[0].strip()
                description = method_comment
            else:
                summary = method_name
                description = ""
        
        # Determine appropriate status code
        success_code = self._get_success_code(verb, method_name)
        
        # Build operation with proper responses
        operation_id = f"{service_name}_{method_name}"
        operation = {
            "summary": summary,
            "operationId": operation_id,
            "responses": self._build_responses(success_code, method_data, verb, path_url, operation_id),
            # Rule 04: per-endpoint security declaration
            "security": self._get_security(operation_id, path_url),
        }
        
        if description:
            operation['description'] = description
        
        # Rule 05: inject pagination query params for GET list operations at
        # generation time — avoids fix_pagination.py having to post-process
        # the assembled openapi.yaml with ruamel.yaml (which reformats it).
        if verb == 'get' and _is_list_operation(method_name, path_url):
            existing_names = {p.get('name') for p in parameters}
            for param in _PAGINATION_PARAMS:
                if param['name'] not in existing_names:
                    parameters.append(param)

        if parameters:
            operation['parameters'] = parameters
            
        # Request Body
        if verb in ['post', 'put', 'patch']:
            body_field = http_rule.body
            if body_field:
                schema_ref = self._get_ref(method_data['input_type'])
                operation['requestBody'] = {
                    "content": {
                        "application/json": {
                            "schema": schema_ref
                        }
                    },
                    "required": True
                }
                
        # Construct Path Item (Single Operation)
        path_item = {
            verb: operation
        }
        
        return path_url, verb, path_item
    
    def _process_custom_actions(self, path_url):
        """Keep custom method with colon syntax as per Google API Design Guide
        
        Google API Design Guide specifies custom methods should use colon syntax:
        POST /v1/policies/{id}:cancel
        POST /v1/policies/{id}:renew
        
        This is NOT a query parameter, it's part of the URL path itself.
        OpenAPI 3.x supports colons in paths.
        """
        # Custom methods with colon are kept as-is in the path
        # /v1/policies/{id}:cancel → /v1/policies/{id}:cancel (no transformation)
        if ':' in path_url and not path_url.startswith('http'):
            # Convert action to kebab-case for consistency
            parts = path_url.split(':')
            base_path = parts[0]
            action = parts[1] if len(parts) > 1 else None
            if action and self.name_transformer:
                action = self.name_transformer._to_kebab_case(action)
                return f"{base_path}:{action}", None
        return path_url, None
    
    def _to_kebab_case_path(self, path):
        """Convert camelCase segments in path to kebab-case"""
        # Split by / and process each segment
        segments = path.split('/')
        result = []
        for seg in segments:
            # Skip parameters {id} and version segments v1
            if seg.startswith('{') or (seg.startswith('v') and seg[1:].isdigit()):
                result.append(seg)
            elif self.name_transformer:
                result.append(self.name_transformer._to_kebab_case(seg))
            else:
                result.append(seg)
        return '/'.join(result)
    
    def _get_success_code(self, verb, method_name):
        """
        Rule 02: Determine correct HTTP success status code.
        - DELETE → 204 No Content
        - Custom action paths (colon syntax) → 200 (handled in caller via path_url check)
        - POST + resource creation verb → 201 Created
        - Everything else → 200 OK
        """
        if verb == 'delete':
            return '204'
        if verb == 'post':
            for creation_verb in RESOURCE_CREATION_VERBS:
                if method_name.startswith(creation_verb):
                    return '201'
        return '200'

    def _get_security(self, operation_id, path_url):
        """
        Rule 04: Return per-endpoint security declaration.
        - Public endpoints → [] (no auth)
        - Webhook endpoints → [] (HMAC authenticated)
        - Everything else → [BearerAuth]
        """
        if operation_id in PUBLIC_OPERATIONS:
            return []
        if operation_id in WEBHOOK_OPERATIONS:
            return []
        return [{"BearerAuth": []}]

    def _response_entry(self, description, inner_ref=None, headers=None, is_error=False, http_status_code=200):
        """
        Build a response entry with the unified ApiResponse envelope.

        ALL responses — success AND error — use the IDENTICAL schema shape so
        that Go / TS / Kotlin / Swift codegen produces ONE type per endpoint:

            ApiResponse<T>
              success: Bool
              data:    T?       ← populated on success, null on error
              error:   Error?   ← populated on error,   null on success
              meta:    ResponseMeta?

        When inner_ref is provided, `data` is typed to the concrete DTO.
        When inner_ref is None (204 endpoints), plain ApiResponse is used.

        The `example` block always shows ALL four fields explicitly so that
        Swagger UI / Redoc / Apidog render a complete, consistent preview:
          - success responses: success=true,  data={DTO fields...}, error=null, meta={...}
          - error responses:   success=false, data=null,            error={...}, meta={...}
        """
        if inner_ref:
            schema = {
                "allOf": [
                    {"$ref": "#/components/schemas/ApiResponse"},
                    {
                        "type": "object",
                        "properties": {
                            "data": inner_ref
                        }
                    }
                ]
            }
        else:
            schema = {"$ref": "#/components/schemas/ApiResponse"}

        # Per-status-code example so Swagger/Redoc/Apidog shows the correct
        # success flag and populated/null data vs error fields.
        # ALL four envelope fields are always present and explicit (never bare/missing).
        if is_error:
            example = {
                "success": False,
                "data": None,
                "error": {
                    "code": "ERROR_CODE",
                    "message": description,
                    "error_id": "err_example",
                    "http_status_code": http_status_code,
                    "retryable": False,
                    "field_violations": []
                },
                "meta": {
                    "request_id": "req_example",
                    "timestamp": "2024-01-15T10:30:00Z"
                }
            }
        else:
            # Build realistic DTO example by walking the schema's properties.
            # This replaces the generic `data: {}` with actual field-level values
            # so Swagger UI / Redoc / Apidog show a meaningful response preview.
            data_example = _build_data_example(inner_ref, self.schemas_dir) if inner_ref else {}
            example = {
                "success": True,
                "data": data_example,
                "error": None,
                "meta": {
                    "request_id": "req_example",
                    "timestamp": "2024-01-15T10:30:00Z",
                    "pagination": None
                }
            }

        entry = {
            "description": description,
            "content": {
                "application/json": {
                    "schema": schema,
                    "example": example
                }
            }
        }
        if headers:
            entry["headers"] = headers
        return entry

    def _build_responses(self, success_code, method_data, verb, path_url, operation_id):
        """
        Rule 02 + 03: Build complete response dictionary.

        CRITICAL: Every status code for a given endpoint uses the IDENTICAL
        schema shape — ApiResponse<T> — so codegen across Go / TS / Kotlin /
        Swift produces exactly ONE type.  The client checks `success` to know
        whether `data` or `error` is populated.
        """
        responses = {}

        # Resolve the typed inner ref once — shared by ALL status codes.
        if success_code == '204':
            inner_ref = None
        else:
            inner_ref = self._get_ref(method_data['output_type'])

        # --- Success response ---
        if success_code == '204':
            responses['204'] = {"description": "No content - Operation completed successfully"}

        elif success_code == '201':
            location_path = self._infer_resource_location(path_url)
            responses['201'] = self._response_entry(
                "Resource created successfully", inner_ref,
                headers={
                    "Location": {
                        "description": f"URL of the newly created resource (e.g. {location_path})",
                        "schema": {"type": "string"}
                    }
                }
            )
        else:
            responses['200'] = self._response_entry("Successful response", inner_ref)

        # --- Error responses (same schema shape, correct HTTP status code in example) ---
        #
        # Error response rules:
        #
        # 400 Bad Request    — ALL endpoints: malformed body / invalid params
        #
        # 401 Unauthorized   — ALL endpoints (including public ones):
        #                       Public endpoints: wrong credentials / expired OTP / invalid token
        #                       Protected endpoints: missing or invalid Bearer token
        #                       Exception: webhook endpoints use HMAC, not Bearer — skip.
        #
        # 403 Forbidden      — ALL endpoints:
        #                       Public endpoints: account locked / suspended / not yet verified
        #                       Protected endpoints: insufficient RBAC permissions
        #                       Exception: webhook endpoints — skip.
        #
        # 404 Not Found      — Endpoints with path parameters ({id}, {user_id}, etc.)
        #                       Also DELETE endpoints (even without params) to signal missing resource
        #
        # 409 Conflict       — POST endpoints that create a new resource (201 Created)
        #                       Also state-transition endpoints (:approve, :cancel, :activate, etc.)
        #
        # 422 Unprocessable  — POST / PUT / PATCH: business-rule violations (valid shape, wrong data)
        #
        # 429 Too Many Req   — ALL endpoints (rate limiting applies everywhere)
        #
        # 500 Internal Error — ALL endpoints

        is_webhook = operation_id in WEBHOOK_OPERATIONS

        # 400 — always
        responses['400'] = self._response_entry(
            "Bad request - Malformed request body or invalid parameters", inner_ref,
            is_error=True, http_status_code=400
        )

        # 401 — ALL endpoints including webhooks:
        #   Webhooks: invalid or missing HMAC signature
        #   Public:   wrong credentials / expired OTP / invalid token
        #   Protected: missing or expired Bearer token
        if is_webhook:
            desc_401 = "Unauthorized - Invalid or missing HMAC webhook signature"
        elif operation_id in PUBLIC_OPERATIONS:
            desc_401 = "Unauthorized - Invalid credentials, expired OTP, or invalid token"
        else:
            desc_401 = "Unauthorized - Valid authentication token required"
        responses['401'] = self._response_entry(
            desc_401, inner_ref, is_error=True, http_status_code=401
        )

        # 403 — ALL endpoints including webhooks:
        #   Webhooks: webhook source IP not whitelisted / provider not recognised
        #   Public:   account locked, suspended, or not yet verified
        #   Protected: insufficient RBAC permissions
        if is_webhook:
            desc_403 = "Forbidden - Webhook source not whitelisted or provider not recognised"
        elif operation_id in PUBLIC_OPERATIONS:
            desc_403 = "Forbidden - Account locked, suspended, or not yet verified"
        else:
            desc_403 = "Forbidden - Insufficient permissions for this operation"
        responses['403'] = self._response_entry(
            desc_403, inner_ref, is_error=True, http_status_code=403
        )

        # 404 — endpoints with path params OR DELETE (resource must exist)
        has_path_params = '{' in path_url
        is_state_action = ':' in path_url  # :approve, :cancel, :activate, etc.
        if has_path_params or verb == 'delete' or is_state_action:
            responses['404'] = self._response_entry(
                "Not found - The requested resource does not exist", inner_ref,
                is_error=True, http_status_code=404
            )

        # 409 — resource creation (201) AND state-transition endpoints
        if success_code == '201' or is_state_action:
            responses['409'] = self._response_entry(
                "Conflict - Resource already exists or state conflict", inner_ref,
                is_error=True, http_status_code=409
            )

        # 422 — POST / PUT / PATCH (business validation)
        if verb in ('post', 'put', 'patch'):
            responses['422'] = self._response_entry(
                "Unprocessable Entity - Request is valid but failed business validation. "
                "See error.field_violations for field-level details", inner_ref,
                is_error=True, http_status_code=422
            )

        # 429 — all endpoints (rate limiting applies everywhere)
        responses['429'] = self._response_entry(
            "Too Many Requests - Rate limit exceeded. Retry after the indicated delay", inner_ref,
            is_error=True, http_status_code=429
        )

        # 500 — always
        responses['500'] = self._response_entry(
            "Internal server error - Unexpected server-side error", inner_ref,
            is_error=True, http_status_code=500
        )

        return responses

    def _infer_resource_location(self, path_url):
        """
        Rule 02: Infer the Location header URL pattern for a created resource.
        e.g. /v1/policies → /v1/policies/{policy_id}
             /v1/users/{user_id}/documents → /v1/users/{user_id}/documents/{document_id}
        """
        # Strip action suffix (colon notation)
        base_path = path_url.split(':')[0].rstrip('/')
        # If path already ends with a parameter, use as-is
        if base_path.endswith('}'):
            return base_path
        # Otherwise append /{singular_id}
        last_segment = base_path.split('/')[-1]
        # Simple singularization: remove trailing 's', handle 'ies' → 'y'
        if last_segment.endswith('ies'):
            singular = last_segment[:-3] + 'y'
        elif last_segment.endswith('s') and not last_segment.endswith('ss'):
            singular = last_segment[:-1]
        else:
            singular = last_segment
        return f"{base_path}/{{{singular}_id}}"

    def _extract_verb_and_path(self, http_rule):
        """Extracts the HTTP verb and URL path from HttpRule."""
        if http_rule.HasField('get'):
            return 'get', http_rule.get
        if http_rule.HasField('post'):
            return 'post', http_rule.post
        if http_rule.HasField('put'):
            return 'put', http_rule.put
        if http_rule.HasField('delete'):
            return 'delete', http_rule.delete
        if http_rule.HasField('patch'):
            return 'patch', http_rule.patch
        return None, None

    def _extract_parameters(self, path_url):
        """Extracts path parameters from URL template (e.g. /users/{id})."""
        # Regex to find {var}
        matches = re.findall(r'\{([^\}]+)\}', path_url)
        params = []
        for match in matches:
            # Handle {var=*} syntax if present (though rare in simple matching)
            # just take the var name
            var_name = match.split('=')[0]
            params.append({
                "name": var_name,
                "in": "path",
                "required": True,
                "schema": {"type": "string"} # Default to string for path params
            })
        return params

    def _get_ref(self, type_name):
        """Resolves a proto type name to an OpenAPI schema ref."""
        # Normalize type name: .package.Message -> package.Message
        if type_name.startswith('.'):
            type_name = type_name[1:]
            
        ref = self.registry.get_ref(type_name)
        if ref:
            return {"$ref": ref}
        else:
            return {"type": "object", "description": f"Unresolved type: {type_name}"}
