import argparse
import os
import sys
import glob
import subprocess
import yaml
from proto_parser import ProtoParser
from registry import ProtoRegistry

def find_protos(root_dir):
    """Recursively finds all .proto files."""
    proto_files = []
    for root, _, files in os.walk(root_dir):
        for file in files:
            if file.endswith(".proto"):
                full_path = os.path.join(root, file)
                # Normalize to forward slashes for protoc consistency
                proto_files.append(full_path.replace("\\", "/"))
    return proto_files

def run_protoc_with_buf(output_descriptor, proto_root):
    """Runs buf to generate the binary descriptor set with dependencies."""
    # Change to proto root to use buf.yaml configuration
    original_dir = os.getcwd()
    
    # Find project root (where buf.yaml is located)
    # proto_root is typically ../../proto from generator directory
    project_root = os.path.abspath(os.path.join(proto_root, ".."))
    
    # Convert output_descriptor to absolute path before changing directory
    abs_output_descriptor = os.path.abspath(output_descriptor)
    
    try:
        os.chdir(project_root)
        
        # Use buf to build descriptor set - this includes all dependencies
        cmd = [
            "buf",
            "build",
            "-o", abs_output_descriptor,
            "--as-file-descriptor-set"
        ]
        
        print(f"Running buf from: {os.getcwd()}")
        print(f"Output descriptor: {abs_output_descriptor}")
        result = subprocess.run(cmd, capture_output=True, text=True)
        
        if result.returncode != 0:
            print(f"Buf build failed with code {result.returncode}")
            print("STDOUT:", result.stdout)
            print("STDERR:", result.stderr)
            raise subprocess.CalledProcessError(result.returncode, cmd, output=result.stdout, stderr=result.stderr)
        else:
            if result.stderr:
                print("Buf warnings:", result.stderr)
    finally:
        os.chdir(original_dir)

def main():
    print("Starting generator...")
    parser = argparse.ArgumentParser(description="OpenAPI Generator from Proto")
    parser.add_argument("--discover", action="store_true", help="Auto-discover all protos")
    parser.add_argument("--descriptor", default="../input/descriptors.pb", help="Path to descriptors.pb")
    parser.add_argument("--proto-root", default="../../proto", help="Root directory for proto imports")
    parser.add_argument("--api-root", default="..", help="Root directory for API generation")
    # parser.add_argument("--output-dir", default="api/schemas", help="Output directory for schemas") # Removed in favor of api-root
    
    args = parser.parse_args()

    # 1. Discovery & Compilation
    if args.discover:
        print(f"Scanning for protos in {args.proto_root}...")
        all_protos = find_protos(args.proto_root)
        print(f"Found {len(all_protos)} proto files.")
        
        if not all_protos:
            print("No proto files found.")
            sys.exit(0)
            
        print("Compiling descriptors with buf (includes dependencies)...")
        try:
            os.makedirs(os.path.dirname(args.descriptor), exist_ok=True)
            run_protoc_with_buf(args.descriptor, args.proto_root)
            print("Compilation successful.")
        except subprocess.CalledProcessError as e:
            print(f"Buf compilation failed: {e}")
            sys.exit(1)

    # 2. Parsing
    # Initialize source parser for security annotations
    from proto_source_parser import ProtoSourceParser
    source_parser = ProtoSourceParser(args.proto_root)
    source_parser.scan_all_protos()
    
    parser = ProtoParser(proto_source_parser=source_parser)
    try:
        if not os.path.exists(args.descriptor):
             print(f"Descriptor file not found: {args.descriptor}. Run with --discover first.")
             sys.exit(1)
             
        parser.load_descriptor_set(args.descriptor)
        print("Successfully loaded descriptor set.")
    except Exception as e:
        print(f"Failed to parse: {e}")
        sys.exit(1)

    # 3. Registry Building with Name Transformation
    from name_transformer import NameTransformer
    
    registry = ProtoRegistry()
    name_transformer = NameTransformer(preserve_request_response=True)
    messages = parser.get_messages()
    
    print(f"Registering {len(messages)} messages and {len(parser.get_enums())} enums with transformed names...")
    
    # Register messages with transformed names
    for msg in messages:
        original_name = msg['descriptor'].name
        transformed_name = name_transformer.transform(original_name)
        registry.register_message(
            full_name=msg['full_name'],
            file_package=msg['package'],
            message_name=transformed_name  # Use transformed name!
        )
    
    # Register enums (no transformation needed)
    for enum in parser.get_enums():
         registry.register_message(
            full_name=enum['full_name'],
            file_package=enum['package'],
            message_name=enum['descriptor'].name
        )
        
    print(f"Registry populated with {len(registry._type_map)} types (with transformed names).")
    
    # Print collision report
    if registry.has_collisions():
        print()
        print("="*80)
        print(registry.get_collision_report())
        print("="*80)
        print()
    
    print("Registry built.")
    
    # 4. Schema Generation
    from schema_generator import SchemaGenerator
    import yaml
    
    print()
    print("Generating schemas...")
    schema_gen = SchemaGenerator(registry)
    
    # Generate Messages
    for msg in messages:
        if msg['package'].startswith('google.protobuf'):
            continue
            
        # generate_schema returns (name, schema, output_path) tuple
        _, schema, _ = schema_gen.generate_schema(msg, parser)
        rel_path = registry.get_file_path(msg['full_name'])
        
        if not rel_path:
             print(f"Skipping {msg['full_name']} (no path resolved)")
             continue
        
        # Determine output folder: events for Event messages, schemas for others
        # This prevents duplicates - events should ONLY be in events/ folder
        msg_name = msg['descriptor'].name
        if msg_name.endswith('Event'):
            base_folder = "events"
        else:
            base_folder = "schemas"
             
        out_path = os.path.join(args.api_root, base_folder, rel_path)
        os.makedirs(os.path.dirname(out_path), exist_ok=True)
        
        # Get the actual schema name from registry (includes collision prefix if any)
        key_name = registry.get_schema_name(msg['full_name'])
        if not key_name:
            # Fallback to transformation
            key_name = name_transformer.transform(msg['descriptor'].name)
        
        # Sanitize to remove Python tuples before writing
        def sanitize_for_yaml(obj):
            if isinstance(obj, tuple):
                return list(obj)
            elif isinstance(obj, dict):
                return {k: sanitize_for_yaml(v) for k, v in obj.items()}
            elif isinstance(obj, list):
                return [sanitize_for_yaml(item) for item in obj]
            return obj
        
        wrapped_schema = { key_name: sanitize_for_yaml(schema) }
        
        with open(out_path, 'w') as f:
            yaml.dump(wrapped_schema, f, sort_keys=False)
            
    # Generate Enums
    for enum in parser.get_enums():
        if enum['package'].startswith('google.protobuf'): continue
        
        # generate_enum_schema returns (name, schema, path) tuple
        original_enum_name, enum_schema, enum_path = schema_gen.generate_enum_schema(enum)
        
        # Get actual enum name from registry (includes collision prefix if any)
        enum_name = registry.get_schema_name(enum['full_name'])
        if not enum_name:
            enum_name = original_enum_name
        
        # Build enum path - FLAT structure in enums/ folder
        # No nested directories, just EnumName.yaml
        enums_dir = os.path.join(args.api_root, "enums")
        os.makedirs(enums_dir, exist_ok=True)  # Create enums folder
        out_path = os.path.join(enums_dir, f'{enum_name}.yaml')
        
        wrapped_schema = { enum_name: enum_schema }
        
        # Sanitize to remove Python tuples before writing
        def sanitize_for_yaml(obj):
            if isinstance(obj, tuple):
                return list(obj)
            elif isinstance(obj, dict):
                return {k: sanitize_for_yaml(v) for k, v in obj.items()}
            elif isinstance(obj, list):
                return [sanitize_for_yaml(item) for item in obj]
            return obj
        
        sanitized = sanitize_for_yaml(wrapped_schema)
        
        with open(out_path, 'w') as f:
            yaml.dump(sanitized, f, sort_keys=False)

    print(f"Schemas generated in {args.api_root}/schemas")

    # 5. Path Generation with Validation
    from path_generator import PathGenerator
    from path_validator import PathValidator
    from endpoint_mapper import EndpointMapper
    
    print("Generating paths with validation...")
    descriptions_dir = os.path.join(args.api_root, "descriptions")
    path_gen = PathGenerator(registry, descriptions_dir=descriptions_dir)
    path_validator = PathValidator()
    endpoint_mapper = EndpointMapper()
    
    services = parser.get_services()
    for service_data in services:
        service_paths = {} # url -> path_item
        
        service_name = service_data['descriptor'].name
        for method in service_data['methods']:
            path_url, verb, path_item_op = path_gen.generate_path_item(method, service_name)
            
            if path_url and verb:
                # Get RPC name
                rpc_name = method.get('name', 'unknown')
                operation_id = path_item_op[verb].get('operationId', 'unknown')
                
                # VALIDATE: Register operation and check for collisions
                is_valid = path_validator.register_operation(
                    path_url=path_url,
                    http_method=verb,
                    operation_id=operation_id,
                    service_name=service_name,
                    rpc_name=rpc_name
                )
                
                # MAP: Add to endpoint mapper
                endpoint_mapper.add_endpoint(
                    path_url=path_url,
                    http_method=verb,
                    operation_details={
                        'operation_id': operation_id,
                        'service': service_name,
                        'rpc_name': rpc_name,
                        'summary': path_item_op[verb].get('summary', 'N/A')
                    }
                )
                
                # Merge into service_paths only if valid (no collision within same service)
                if path_url not in service_paths:
                    service_paths[path_url] = {}
                
                if verb in service_paths[path_url]:
                    print(f"  Warning: Duplicate {verb.upper()} for {path_url} within {service_name}")
                    print(f"    Existing: {service_paths[path_url][verb].get('operationId')}")
                    print(f"    New: {operation_id}")
                    print(f"    Keeping first operation")
                else:
                    service_paths[path_url].update(path_item_op)
                
        # Write Service Paths file
        if service_paths:
            # Construct File Path: api/paths/<package>/<ServiceName>.yaml
            package_path = service_data['package'].replace('.', '/')
            out_path = os.path.join(args.api_root, "paths", package_path, f"{service_data['descriptor'].name}.yaml")
            os.makedirs(os.path.dirname(out_path), exist_ok=True)
            
            # Write key-value: PathURL: PathItem
            with open(out_path, 'w') as f:
                yaml.dump(service_paths, f, sort_keys=False)
                
    print("Path generation complete.")
    
    # Print validation report
    path_validator.print_report()
    
    # Print endpoint map
    endpoint_mapper.print_endpoint_map()
    
    # Export endpoint map to markdown
    endpoint_map_file = os.path.join(args.api_root, "ENDPOINT_MAP.md")
    endpoint_mapper.export_to_markdown(endpoint_map_file)
    print(f"\n✅ Endpoint map exported to: {endpoint_map_file}")

    # 6. Assembly
    from assembler import OpenAPIAssembler
    print("Assembling openapi.yaml...")
    assembler = OpenAPIAssembler(registry, args.api_root)
    root_spec = assembler.assemble()
    
    # 6.5. Post-processing (add required fields, pagination, etc.)
    from post_processor import post_process_openapi_spec
    root_spec = post_process_openapi_spec(root_spec)
    
    output_file = os.path.join(args.api_root, "openapi.yaml")
    with open(output_file, 'w') as f:
        yaml.dump(root_spec, f, sort_keys=False)
        
    print(f"OpenAPI Spec generated: {output_file}")
    
    # 7. Generate JSON version for JavaScript consumption (Schema Visualizer)
    print("\nGenerating openapi.json for Schema Visualizer...")
    json_output_file = os.path.join(args.api_root, "docs", "openapi.json")
    os.makedirs(os.path.dirname(json_output_file), exist_ok=True)
    
    import json
    with open(json_output_file, 'w', encoding='utf-8') as f:
        json.dump(root_spec, f, indent=2, ensure_ascii=False)
    
    schema_count = len(root_spec.get('components', {}).get('schemas', {}))
    print(f"✅ OpenAPI JSON generated: {json_output_file}")
    print(f"   📊 {schema_count} schemas available for visualization")
    
    # 8. Copy Schema Visualizer JavaScript files to docs directory
    print("\nCopying Schema Visualizer files...")
    import shutil
    
    generator_dir = os.path.dirname(os.path.abspath(__file__))
    docs_dir = os.path.join(args.api_root, "docs")
    
    visualizer_files = [
        'schema-visualizer.js',
        'openapi-loader.js'
    ]
    
    for js_file in visualizer_files:
        src = os.path.join(generator_dir, js_file)
        dst = os.path.join(docs_dir, js_file)
        
        if os.path.exists(src):
            shutil.copy2(src, dst)
            print(f"   ✅ Copied {js_file}")
        else:
            print(f"   ⚠️  Warning: {js_file} not found in generator directory")
    
    print("✅ Schema Visualizer files ready")
    
    # 9. Analyze Proto Schema Structure
    print("\n" + "="*80)
    print("ANALYZING PROTO SCHEMA STRUCTURE")
    print("="*80)
    from proto_schema_analyzer import ProtoSchemaAnalyzer
    
    proto_analyzer = ProtoSchemaAnalyzer(args.proto_root)
    proto_summary = proto_analyzer.scan_all_protos()
    proto_analyzer.print_summary(proto_summary)
    
    # Export for doc generator
    schema_summary_file = os.path.join(args.api_root, "proto_schema_summary.json")
    proto_analyzer.export_to_json(schema_summary_file, proto_summary)
    
    # 10. (Optional) Validate with Neon DB
    print("\n" + "="*80)
    print("VALIDATING WITH NEON DATABASE")
    print("="*80)
    try:
        from neon_db_inspector import NeonDBInspector
        
        # Find .env file in project root
        env_file = os.path.join(args.api_root, "..", ".env")
        if not os.path.exists(env_file):
            print(f"⚠️  .env file not found at {env_file}")
            print("   Skipping DB validation (using proto data only)")
        else:
            db_inspector = NeonDBInspector(env_file=env_file)
            comparison = db_inspector.compare_with_proto(proto_summary)
            
            if 'error' not in comparison:
                print(f"✅ DB Connection successful")
                print(f"   Matched: {len(comparison['matched'])} tables")
                print(f"   Proto only: {len(comparison['proto_only'])} tables")
                print(f"   DB only: {len(comparison['db_only'])} tables")
                
                # Export comparison
                comparison_file = os.path.join(args.api_root, "db_proto_comparison.json")
                with open(comparison_file, 'w') as f:
                    json.dump(comparison, f, indent=2)
                print(f"   Comparison exported: {comparison_file}")
            else:
                print(f"⚠️  DB validation skipped: {comparison['error']}")
                print("   Using proto data only")
    except Exception as e:
        print(f"⚠️  DB validation skipped: {e}")
        print("   Using proto data only")
    
    # 11. Map APIs to Schema Groups
    print("\n" + "="*80)
    print("MAPPING APIs TO SCHEMA GROUPS")
    print("="*80)
    from api_schema_mapper import APISchemaMapper
    
    # Extract APIs by domain from OpenAPI spec
    apis_by_domain = {}
    for path, methods in root_spec.get('paths', {}).items():
        # Extract domain from path: /v1/{domain}/...
        parts = path.strip('/').split('/')
        if len(parts) >= 2 and parts[0] == 'v1':
            # Handle action suffixes and hyphens
            raw_domain = parts[1]
            
            if ':' in raw_domain:
                base = raw_domain.split(':')[0]
                # Map hyphenated domains
                hyphen_map = {
                    'kyc-verifications': 'kyc',
                    'notification-templates': 'notification',
                    'api-keys': 'apikey'
                }
                domain = hyphen_map.get(base, base)
            elif '-' in raw_domain:
                domain_map = {
                    'api-keys': 'apikey',
                    'audit-events': 'audit',
                    'audit-logs': 'audit',
                    'compliance-logs': 'audit',
                    'notification-templates': 'notification',
                    'commission-payouts': 'commission',
                    'document-templates': 'document',
                    'fraud-alerts': 'fraud',
                    'fraud-cases': 'fraud',
                    'kyc-verifications': 'kyc',
                    'knowledge-base': 'support',
                    'renewal-schedules': 'renewal',
                    'report-definitions': 'report',
                    'report-schedules': 'report',
                    'voice-sessions': 'voice',
                    'workflow-definitions': 'workflow',
                    'workflow-instances': 'workflow',
                    'workflow-tasks': 'workflow'
                }
                domain = domain_map.get(raw_domain, raw_domain.split('-')[0])
            else:
                # Map plurals and aliases
                singular_map = {
                    'auth': 'authn',
                    'users': 'authn',
                    'roles': 'authz',
                    'beneficiaries': 'beneficiary',
                    'claims': 'claims',
                    'commissions': 'commission',
                    'documents': 'document',
                    'endorsements': 'endorsement',
                    'insurers': 'insurer',
                    'notifications': 'notification',
                    'partners': 'partner',
                    'payments': 'payment',
                    'policies': 'policy',
                    'products': 'products',
                    'quotes': 'underwriting',
                    'refunds': 'refund',
                    'renewals': 'renewal',
                    'reports': 'report',
                    'tasks': 'task',
                    'tenants': 'tenant',
                    'faqs': 'support',
                    'tickets': 'support'
                }
                domain = singular_map.get(raw_domain, raw_domain)
            
            if domain not in apis_by_domain:
                apis_by_domain[domain] = []
            
            for method, operation in methods.items():
                if method in ['get', 'post', 'put', 'delete', 'patch']:
                    apis_by_domain[domain].append({
                        'path': path,
                        'method': method.upper(),
                        'operation_id': operation.get('operationId', ''),
                        'summary': operation.get('summary', '')
                    })
    
    mapper = APISchemaMapper()
    schema_api_summary = mapper.generate_schema_summary(apis_by_domain, proto_summary)
    
    # Export for doc generator
    schema_api_file = os.path.join(args.api_root, "schema_api_mapping.json")
    with open(schema_api_file, 'w') as f:
        json.dump(schema_api_summary, f, indent=2)
    
    print(f"✅ API-Schema mapping complete")
    print(f"   Schema groups: {schema_api_summary['stats']['total_schemas']}")
    print(f"   Total tables: {schema_api_summary['stats']['total_tables']}")
    print(f"   Total APIs: {schema_api_summary['stats']['total_apis']}")
    print(f"   Mapping exported: {schema_api_file}")
    
    # 12. Generate Enhanced Documentation
    print("\n" + "="*80)
    print("GENERATING ENHANCED DOCUMENTATION")
    print("="*80)
    from doc_generator import DocGenerator
    
    doc_gen = DocGenerator(output_file, schema_summary_file)
    index_html = os.path.join(args.api_root, "docs", "index.html")
    doc_gen.generate_documentation(index_html)
    
    print("\n" + "="*80)
    print("✅ PIPELINE COMPLETE")
    print("="*80)

if __name__ == "__main__":
    main()
