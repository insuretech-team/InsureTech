import os
import yaml

# ---------------------------------------------------------------------------
# Rule 01 + 05: Canonical schemas injected into every generated spec.
# These are the standard envelope and pagination schemas all clients depend on.
# ---------------------------------------------------------------------------
CANONICAL_SCHEMAS = {
    "ApiResponse": {
        "type": "object",
        "required": ["success"],
        "description": (
            "Standard API response envelope. Identical shape for EVERY status "
            "code \u2014 200, 201, 400, 401, 500, etc. Client code always decodes to "
            "ApiResponse<T> and checks `success`.\n\n"
            "success=true  \u2192 `data` has the typed payload, `error` is null.\n"
            "success=false \u2192 `data` is null, `error` has details."
        ),
        "properties": {
            "success": {
                "type": "boolean",
                "description": "true when the operation succeeded, false on any error."
            },
            "data": {
                "description": (
                    "Typed response payload on success. "
                    "null on failure or no-content actions (logout, delete, etc.). "
                    "Concrete type is specified per-endpoint via allOf composition."
                )
            },
            "error": {
                "$ref": "#/components/schemas/Error",
                "description": "Error details on failure. Always null on success."
            },
            "meta": {
                "$ref": "#/components/schemas/ResponseMeta",
                "description": "Response metadata: request tracing, pagination, timestamps."
            }
        }
    },
    "ResponseMeta": {
        "type": "object",
        "description": "Metadata attached to every API response for tracing and pagination.",
        "properties": {
            "request_id": {
                "type": "string",
                "description": "Unique request trace ID. Use for support lookups and debugging."
            },
            "timestamp": {
                "type": "string",
                "format": "date-time",
                "description": "Server timestamp when the response was generated (UTC ISO 8601)."
            },
            "api_version": {
                "type": "string",
                "description": "API version that served this response."
            },
            "pagination": {
                "allOf": [{"$ref": "#/components/schemas/PaginationMeta"}],
                "nullable": True,
                "description": "Pagination info for list responses. null for non-list endpoints."
            }
        }
    },
    "PaginationMeta": {
        "type": "object",
        "required": ["page", "page_size", "total_pages", "total_items", "has_next", "has_previous"],
        "description": "Standard pagination metadata for all list endpoints.",
        "properties": {
            "page": {
                "type": "integer",
                "description": "Current page number (1-based)."
            },
            "page_size": {
                "type": "integer",
                "description": "Number of items returned in this page."
            },
            "total_pages": {
                "type": "integer",
                "description": "Total number of pages available."
            },
            "total_items": {
                "type": "integer",
                "format": "int64",
                "description": "Total number of items across all pages."
            },
            "has_next": {
                "type": "boolean",
                "description": "Whether a next page exists."
            },
            "has_previous": {
                "type": "boolean",
                "description": "Whether a previous page exists."
            },
            "next_page_token": {
                "type": "string",
                "nullable": True,
                "description": "Cursor token for cursor-based pagination (optional)."
            }
        }
    }
}

class OpenAPIAssembler:
    def __init__(self, registry, output_dir):
        self.registry = registry
        self.output_dir = output_dir

    def assemble(self):
        """
        Assembles the root openapi.yaml file with external references.
        
        Strategy: 
        - Paths are inlined (they reference schemas)
        - Schemas are EXTERNAL references to individual files
        - Only enums are inlined for simplicity
        - Common components can be inlined
        """
        root_schema = {
            "openapi": "3.1.0",
            "info": {
                "title": "InsureTech API",
                "version": "1.0.0",
                "description": "Auto-generated OpenAPI v3.1 specification from Protocol Buffers."
            },
            "servers": [
                {"url": "https://api.labaidinsuretech.com", "description": "Production Server"},
                {"url": "https://staging-api.labaidinsuretech.com", "description": "Staging Server"}
            ],
            "paths": {},
            "components": {
                "schemas": {},
                "securitySchemes": {
                    "BearerAuth": {
                        "type": "http",
                        "scheme": "bearer",
                        "bearerFormat": "JWT"
                    }
                }
            },
            "security": [
                {"BearerAuth": []}
            ]
        }

        # 1. Populate Paths
        print("Loading paths...")
        paths_dir = os.path.join(self.output_dir, "paths")
        if os.path.exists(paths_dir):
            path_count = 0
            method_count = 0
            conflicts_detected = 0
            
            for root, dirs, files in os.walk(paths_dir):
                dirs.sort()
                for file in sorted(files):
                    if file.endswith(".yaml"):
                        file_path = os.path.join(root, file)
                        with open(file_path, 'r', encoding='utf-8') as f:
                            path_data = yaml.safe_load(f)
                            if path_data:
                                for path_url, path_item in sorted(path_data.items()):
                                    # Check if path already exists - MERGE methods instead of overwriting
                                    if path_url in root_schema['paths']:
                                        # Path exists - merge HTTP methods
                                        for http_method, operation in path_item.items():
                                            if http_method in root_schema['paths'][path_url]:
                                                # Conflict detected - same path and method
                                                conflicts_detected += 1
                                                existing_op_id = root_schema['paths'][path_url][http_method].get('operationId', 'unknown')
                                                new_op_id = operation.get('operationId', 'unknown')
                                                print(f"  ⚠️  CONFLICT: {http_method.upper()} {path_url}")
                                                print(f"      Existing: {existing_op_id} (from {file})")
                                                print(f"      New: {new_op_id}")
                                                print(f"      → Keeping existing operation")
                                            else:
                                                # New method for existing path - safe to add
                                                root_schema['paths'][path_url][http_method] = operation
                                                method_count += 1
                                    else:
                                        # First time seeing this path - add all methods
                                        root_schema['paths'][path_url] = path_item
                                        method_count += len(path_item)
                                        path_count += 1
            
            print(f"  Loaded {path_count} unique paths with {method_count} operations")
            if conflicts_detected > 0:
                print(f"  ⚠️  Detected and resolved {conflicts_detected} method conflicts")

        # 2. Populate Components/Schemas from components/ folder
        print("Loading common components...")
        components_schemas_dir = os.path.join(self.output_dir, "components", "schemas")
        component_count = 0
        if os.path.exists(components_schemas_dir):
            for file in sorted(os.listdir(components_schemas_dir)):
                if file.endswith(".yaml"):
                    file_path = os.path.join(components_schemas_dir, file)
                    with open(file_path, 'r', encoding='utf-8') as f:
                        schema_data = yaml.safe_load(f)
                        if schema_data:
                            # File format: { ComponentName: { ... } }
                            for schema_name, schema_def in schema_data.items():
                                root_schema['components']['schemas'][schema_name] = schema_def
                                component_count += 1
        print(f"  Loaded {component_count} common components")

        # 3. Create external references for DTOs, schemas, and events
        # Instead of inlining, we create $ref pointers to external files
        print("Creating external references for schemas...")
        
        # DTOs
        dtos_dir = os.path.join(self.output_dir, "dtos")
        dto_count = 0
        if os.path.exists(dtos_dir):
            for root, dirs, files in os.walk(dtos_dir):
                dirs.sort()
                for file in sorted(files):
                    if file.endswith(".yaml"):
                        file_path = os.path.join(root, file)
                        with open(file_path, 'r', encoding='utf-8') as f:
                            schema_data = yaml.safe_load(f)
                            if schema_data:
                                for schema_name, schema_def in schema_data.items():
                                    # Inline the schema definition
                                    if schema_name not in root_schema['components']['schemas']:
                                        root_schema['components']['schemas'][schema_name] = schema_def
                                        dto_count += 1
        print(f"  Loaded {dto_count} DTO schemas")

        # Entities
        schemas_dir = os.path.join(self.output_dir, "schemas")
        entity_count = 0
        if os.path.exists(schemas_dir):
            for root, dirs, files in os.walk(schemas_dir):
                dirs.sort()
                for file in sorted(files):
                    if file.endswith(".yaml"):
                        if 'google' in root:
                            continue
                        
                        file_path = os.path.join(root, file)
                        with open(file_path, 'r', encoding='utf-8') as f:
                            schema_data = yaml.safe_load(f)
                            if schema_data:
                                for schema_name, schema_def in schema_data.items():
                                    if schema_name in root_schema['components']['schemas']:
                                        print(f"  Warning: Duplicate schema '{schema_name}' - keeping first")
                                    else:
                                        # Inline the schema definition
                                        root_schema['components']['schemas'][schema_name] = schema_def
                                        entity_count += 1
        print(f"  Loaded {entity_count} entity schemas")

        # Events
        events_dir = os.path.join(self.output_dir, "events")
        event_count = 0
        if os.path.exists(events_dir):
            for root, dirs, files in os.walk(events_dir):
                dirs.sort()
                for file in sorted(files):
                    if file.endswith(".yaml"):
                        file_path = os.path.join(root, file)
                        with open(file_path, 'r', encoding='utf-8') as f:
                            schema_data = yaml.safe_load(f)
                            if schema_data:
                                for schema_name, schema_def in schema_data.items():
                                    if schema_name in root_schema['components']['schemas']:
                                        print(f"  Warning: Duplicate event schema '{schema_name}' - keeping first")
                                    else:
                                        # Inline the schema definition
                                        root_schema['components']['schemas'][schema_name] = schema_def
                                        event_count += 1
        print(f"  Loaded {event_count} event schemas")

        # 6. Populate Schemas from enums/ folder (flat structure)
        print("Loading enum schemas...")
        enums_dir = os.path.join(self.output_dir, "enums")
        enum_count = 0
        if os.path.exists(enums_dir):
            for file in sorted(os.listdir(enums_dir)):
                if file.endswith(".yaml"):
                    file_path = os.path.join(enums_dir, file)
                    with open(file_path, 'r', encoding='utf-8') as f:
                        schema_data = yaml.safe_load(f)
                        if schema_data:
                            for schema_name, schema_def in schema_data.items():
                                # Check for duplicates
                                if schema_name in root_schema['components']['schemas']:
                                    print(f"  Warning: Duplicate enum schema '{schema_name}' - keeping first")
                                else:
                                    root_schema['components']['schemas'][schema_name] = schema_def
                                    enum_count += 1
        print(f"  Loaded {enum_count} enum schemas")

        # Summary
        total_schemas = component_count + dto_count + entity_count + event_count + enum_count
        print(f"\nTotal schemas in components: {total_schemas}")
        print(f"  - Common components: {component_count}")
        print(f"  - DTOs: {dto_count}")
        print(f"  - Entities: {entity_count}")
        print(f"  - Events: {event_count}")
        print(f"  - Enums: {enum_count}")
        
        # Rule 04: Extend securitySchemes with ApiKeyAuth for B2B/partner integrations
        root_schema['components']['securitySchemes']['ApiKeyAuth'] = {
            "type": "apiKey",
            "in": "header",
            "name": "X-API-Key",
            "description": (
                "API Key for B2B/partner integrations. "
                "Obtained from POST /v1/auth/api-keys or POST /v1/partners/{id}/credentials:rotate"
            )
        }

        # Rule 04: Global security default — BearerAuth for all endpoints.
        # Individual endpoints override this with security: [] for public routes.
        # Per-endpoint security is set by PathGenerator._get_security().
        root_schema['security'] = [{'BearerAuth': []}]
        print("\n✓ Global security default set: BearerAuth (public endpoints declare security: [])")

        # Rule 01 + 05: Inject canonical schemas FIRST so they always exist
        # and are not overwritten by generated schemas.
        print("\nInjecting canonical schemas (ApiResponse, ResponseMeta, PaginationMeta)...")
        for schema_name, schema_def in CANONICAL_SCHEMAS.items():
            root_schema['components']['schemas'][schema_name] = schema_def
            print(f"  ✓ Injected: {schema_name}")

        # Rule 05: Fix PageResponse.total_items type (was string, must be integer)
        # and retire the duplicate PaginationResponse schema in favour of PaginationMeta.
        if 'PageResponse' in root_schema['components']['schemas']:
            pr = root_schema['components']['schemas']['PageResponse']
            props = pr.get('properties', {})
            if 'total_items' in props:
                props['total_items'] = {
                    "type": "integer",
                    "format": "int64",
                    "description": "Total number of items across all pages."
                }
            print("  ✓ Fixed PageResponse.total_items type: string → integer/int64")

        # Rule 05: Replace all $ref to PaginationResponse with PaginationMeta FIRST,
        # then remove the schema so no dangling refs remain.
        def _replace_pagination_refs(obj):
            """Recursively replace PaginationResponse $refs with PaginationMeta."""
            if isinstance(obj, dict):
                for k, v in obj.items():
                    if k == '$ref' and isinstance(v, str) and '/schemas/PaginationResponse' in v:
                        obj[k] = v.replace('/schemas/PaginationResponse', '/schemas/PaginationMeta')
                    else:
                        _replace_pagination_refs(v)
            elif isinstance(obj, list):
                for item in obj:
                    _replace_pagination_refs(item)

        _replace_pagination_refs(root_schema)

        if 'PaginationResponse' in root_schema['components']['schemas']:
            del root_schema['components']['schemas']['PaginationResponse']
            print("  ✓ Replaced all PaginationResponse $refs → PaginationMeta and removed schema")

        # Sort paths and schemas for deterministic output across runs
        root_schema['paths'] = dict(sorted(root_schema['paths'].items()))
        root_schema['components']['schemas'] = dict(sorted(root_schema['components']['schemas'].items()))

        return root_schema
