"""
Updated OpenAPI Generator with all integrations
Integrates: file_organizer, name_transformer, description_loader, improved path_generator
"""
import argparse
import os
import sys
import glob
import subprocess
import yaml
from proto_parser import ProtoParser
from registry import ProtoRegistry
from schema_generator import SchemaGenerator
from path_generator import PathGenerator
from assembler import OpenAPIAssembler

def find_protos(root_dir):
    """Recursively finds all .proto files."""
    proto_files = []
    for root, _, files in os.walk(root_dir):
        for file in files:
            if file.endswith(".proto"):
                full_path = os.path.join(root, file)
                proto_files.append(full_path.replace("\\", "/"))
    return proto_files

def run_protoc(proto_files, output_descriptor, import_root):
    """Runs protoc to generate the binary descriptor set."""
    args_file_path = os.path.join(os.path.dirname(output_descriptor), "protoc_args.txt")
    with open(args_file_path, "w") as f:
        for p in proto_files:
            f.write(f"{p}\n")
            
    cmd = [
        "protoc",
        f"--descriptor_set_out={output_descriptor}",
        "--include_imports",
        "--include_source_info",
        f"--proto_path={import_root}",
        f"--proto_path=C:\\Users\\faruk\\.gemini\\tools\\include",
        f"@{args_file_path}"
    ]
    
    result = subprocess.run(cmd, capture_output=True, text=True)
    
    if result.returncode != 0:
        print(f"❌ Protoc failed with code {result.returncode}")
        print("STDERR:", result.stderr)
        raise subprocess.CalledProcessError(result.returncode, cmd, output=result.stdout, stderr=result.stderr)
    else:
        if result.stderr:
            print("⚠️  Protoc warnings:", result.stderr)

def main():
    print("\n" + "="*80)
    print("OpenAPI Generator (Updated with all improvements)")
    print("="*80 + "\n")
    
    parser = argparse.ArgumentParser(description="OpenAPI Generator from Proto")
    parser.add_argument("--discover", action="store_true", help="Auto-discover all protos")
    parser.add_argument("--descriptor", default="api/input/descriptors.pb", help="Path to descriptors.pb")
    parser.add_argument("--proto-root", default="proto", help="Root directory for proto imports")
    parser.add_argument("--api-root", default="api", help="Root directory for API generation")
    parser.add_argument("--descriptions-dir", default="api/descriptions", help="Directory with markdown descriptions")
    parser.add_argument("--skip-generation", action="store_true", help="Skip schema/path generation, only assemble")
    
    args = parser.parse_args()

    # 1. Discovery & Compilation (if requested)
    if args.discover:
        print(f"📁 Scanning for protos in {args.proto_root}...")
        all_protos = find_protos(args.proto_root)
        print(f"   Found {len(all_protos)} proto files.")
        
        if not all_protos:
            print("❌ No proto files found.")
            sys.exit(0)
            
        print("🔨 Compiling descriptors...")
        try:
            os.makedirs(os.path.dirname(args.descriptor), exist_ok=True)
            run_protoc(all_protos, args.descriptor, args.proto_root)
            print("✅ Compilation successful.")
        except subprocess.CalledProcessError as e:
            print(f"❌ Protoc compilation failed: {e}")
            sys.exit(1)

    # 2. Parsing
    proto_parser = ProtoParser()
    try:
        if not os.path.exists(args.descriptor):
             print(f"❌ Descriptor file not found: {args.descriptor}")
             print("   Run with --discover first.")
             sys.exit(1)
             
        print(f"📖 Loading descriptor set from {args.descriptor}...")
        proto_parser.load_descriptor_set(args.descriptor)
        print("✅ Successfully loaded descriptor set.")
    except Exception as e:
        print(f"❌ Failed to parse: {e}")
        sys.exit(1)

    # 3. Registry Building
    print("\n📋 Building registry...")
    registry = ProtoRegistry()
    messages = proto_parser.get_messages()
    enums = proto_parser.get_enums()
    services = proto_parser.get_services()
    
    print(f"   Messages: {len(messages)}")
    print(f"   Enums: {len(enums)}")
    print(f"   Services: {len(services)}")
    
    for msg in messages:
        registry.register_message(
            full_name=msg['full_name'],
            file_package=msg['package'],
            message_name=msg['descriptor'].name
        )
        
    for enum in enums:
         registry.register_message(
            full_name=enum['full_name'],
            file_package=enum['package'],
            message_name=enum['descriptor'].name
        )
        
    print(f"✅ Registry populated with {len(registry._type_map)} types.\n")
    
    if args.skip_generation:
        print("⏭️  Skipping generation (--skip-generation flag)")
    else:
        # 4. Schema Generation with all improvements
        print("=" * 80)
        print("Generating Schemas")
        print("=" * 80)
        
        schema_gen = SchemaGenerator(
            registry=registry,
            output_dir=args.api_root,
            descriptions_dir=args.descriptions_dir
        )
        
        schema_count = 0
        
        # Generate Messages
        print("\n📝 Generating message schemas...")
        for msg in messages:
            if msg['package'].startswith('google.protobuf'):
                continue
            
            try:
                transformed_name, schema, output_path = schema_gen.generate_schema(msg, proto_parser)
                
                if output_path:
                    schema_gen.write_schema_file(transformed_name, schema, output_path)
                    schema_count += 1
                else:
                    print(f"⚠️  No output path for {msg['full_name']}")
                    
            except Exception as e:
                print(f"❌ Error generating {msg['full_name']}: {e}")
                
        print(f"✅ Generated {schema_count} message schemas")
        
        # Generate Enums
        print("\n📝 Generating enum schemas...")
        enum_count = 0
        for enum in enums:
            if enum['package'].startswith('google.protobuf'):
                continue
            
            try:
                enum_name, schema, output_path = schema_gen.generate_enum_schema(enum)
                
                if output_path:
                    schema_gen.write_schema_file(enum_name, schema, output_path)
                    enum_count += 1
                    
            except Exception as e:
                print(f"❌ Error generating {enum['full_name']}: {e}")
                
        print(f"✅ Generated {enum_count} enum schemas")
        
        # 5. Path Generation with improvements
        print("\n" + "=" * 80)
        print("Generating Paths")
        print("=" * 80 + "\n")
        
        path_gen = PathGenerator(
            registry=registry,
            descriptions_dir=args.descriptions_dir
        )
        
        path_count = 0
        
        for service in services:
            service_name = service['full_name'].split('.')[-1]
            print(f"📝 Generating paths for {service_name}...")
            
            paths_dict = {}
            
            for method in service['methods']:
                try:
                    path_url, verb, path_item = path_gen.generate_path_item(method, service_name)
                    
                    if path_url and verb and path_item:
                        # Merge into paths_dict
                        if path_url not in paths_dict:
                            paths_dict[path_url] = {}
                        paths_dict[path_url].update(path_item)
                        path_count += 1
                        
                except Exception as e:
                    print(f"❌ Error generating path for {method['name']}: {e}")
            
            # Write service paths to file
            if paths_dict:
                service_path_file = os.path.join(
                    args.api_root, 
                    "paths", 
                    service['package'].replace('.', '/'),
                    f"{service_name}.yaml"
                )
                os.makedirs(os.path.dirname(service_path_file), exist_ok=True)
                
                with open(service_path_file, 'w', encoding='utf-8') as f:
                    yaml.dump(paths_dict, f, sort_keys=False, allow_unicode=True)
                    
        print(f"✅ Generated {path_count} paths")
    
    # 6. Assembly
    print("\n" + "=" * 80)
    print("Assembling OpenAPI Spec")
    print("=" * 80 + "\n")
    
    assembler = OpenAPIAssembler(registry, args.api_root)
    spec = assembler.assemble()
    
    # Write final spec
    output_file = os.path.join(args.api_root, 'openapi.yaml')
    with open(output_file, 'w', encoding='utf-8') as f:
        yaml.dump(spec, f, sort_keys=False, allow_unicode=True)
    
    print(f"\n✅ OpenAPI spec written to: {output_file}")
    print(f"   Paths: {len(spec['paths'])}")
    print(f"   Schemas: {len(spec['components']['schemas'])}")
    
    print("\n" + "=" * 80)
    print("Generation Complete!")
    print("=" * 80)
    print("\nNext steps:")
    print("  1. Run validator: python validator.py ../openapi.yaml")
    print("  2. Review validation report")
    print("  3. Populate descriptions in api/descriptions/")
    print("\n")

if __name__ == '__main__':
    try:
        main()
    except KeyboardInterrupt:
        print("\n\n⚠️  Generation interrupted by user")
        sys.exit(1)
    except Exception as e:
        print(f"\n❌ Fatal error: {e}")
        import traceback
        traceback.print_exc()
        sys.exit(1)
