import os
import yaml

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
            
            for root, _, files in os.walk(paths_dir):
                for file in files:
                    if file.endswith(".yaml"):
                        file_path = os.path.join(root, file)
                        with open(file_path, 'r', encoding='utf-8') as f:
                            path_data = yaml.safe_load(f)
                            if path_data:
                                for path_url, path_item in path_data.items():
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
            for file in os.listdir(components_schemas_dir):
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
            for root, _, files in os.walk(dtos_dir):
                for file in files:
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
            for root, _, files in os.walk(schemas_dir):
                for file in files:
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
            for root, _, files in os.walk(events_dir):
                for file in files:
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
            for file in os.listdir(enums_dir):
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
        
        # Apply global security if security schemes are defined
        if root_schema['components'].get('securitySchemes'):
            print("\nApplying global security to all operations...")
            # Add global security - all operations require BearerAuth by default
            root_schema['security'] = [
                {'BearerAuth': []}
            ]
            print("  ✓ Global security applied: All operations now require BearerAuth")
        else:
            print("\n⚠️  No security schemes defined - operations will be unprotected")

        return root_schema
