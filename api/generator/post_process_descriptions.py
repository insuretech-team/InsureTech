"""
Post-process schemas to add descriptions from markdown files and proto comments
"""
import os
import yaml
import re

def load_md_description(md_file):
    """Extract description from markdown file"""
    if not os.path.exists(md_file):
        return None
    
    with open(md_file, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # Try to find proto comment (most reliable)
    match = re.search(r'\*\*Proto Comment\*\*:\s*(.+?)(?:\n\s*\n|\n\*\*)', content, re.DOTALL)
    if match:
        desc = match.group(1).strip()
        # Remove "Maps to..." part
        if 'Maps to' in desc:
            desc = desc.split('Maps to')[0].strip()
        # Clean up newlines
        desc = ' '.join(desc.split())
        if desc and len(desc) > 10:
            return desc
    
    # Fallback: look for any descriptive text
    match = re.search(r'\*\*Purpose\*\*:\s*(.+?)(?:\n\s*\n|\n\*\*)', content, re.DOTALL)
    if match:
        desc = match.group(1).strip()
        desc = ' '.join(desc.split())
        if desc and len(desc) > 10:
            return desc
    
    return None

def add_descriptions_to_openapi():
    """Add descriptions to openapi.yaml from individual schema files and markdown"""
    
    print("Loading openapi.yaml...")
    with open('../openapi.yaml', 'r', encoding='utf-8') as f:
        spec = yaml.safe_load(f)
    
    schemas = spec['components']['schemas']
    updated_count = 0
    
    # Map schema names to their source files
    schema_files = {}
    
    # Scan schemas directory
    for root, dirs, files in os.walk('../schemas'):
        for file in files:
            if file.endswith('.yaml'):
                file_path = os.path.join(root, file)
                with open(file_path, 'r', encoding='utf-8') as f:
                    data = yaml.safe_load(f)
                    if data:
                        for schema_name in data.keys():
                            schema_files[schema_name] = file_path
    
    # Scan events directory
    for root, dirs, files in os.walk('../events'):
        for file in files:
            if file.endswith('.yaml'):
                file_path = os.path.join(root, file)
                with open(file_path, 'r', encoding='utf-8') as f:
                    data = yaml.safe_load(f)
                    if data:
                        for schema_name in data.keys():
                            schema_files[schema_name] = file_path
    
    print(f"Found {len(schema_files)} schema source files")
    
    # Update schemas without descriptions
    for schema_name, schema_def in schemas.items():
        if not isinstance(schema_def, dict):
            continue
        
        # Skip if already has description
        if schema_def.get('description', '').strip():
            continue
        
        # Find corresponding markdown file
        if schema_name in schema_files:
            yaml_path = schema_files[schema_name]
            # Convert yaml path to markdown path
            # ../schemas/insuretech/ai/entity/v1/AIAgent.yaml
            # -> ../descriptions/entity/insuretech/ai/entity/v1/AIAgent.md
            
            rel_path = os.path.relpath(yaml_path, '..')
            
            if rel_path.startswith('schemas'):
                md_path = rel_path.replace('schemas', 'descriptions/entity', 1)
            elif rel_path.startswith('events'):
                md_path = rel_path.replace('events', 'descriptions/event', 1)
            else:
                continue
            
            md_path = md_path.replace('.yaml', '.md')
            md_path = os.path.join('..', md_path)
            
            # Load description from markdown
            desc = load_md_description(md_path)
            if desc:
                schema_def['description'] = desc
                updated_count += 1
                print(f"  Added description to: {schema_name}")
    
    print(f"\nUpdated {updated_count} schemas with descriptions")
    
    # Write back
    print("Writing updated openapi.yaml...")
    with open('../openapi.yaml', 'w', encoding='utf-8') as f:
        yaml.dump(spec, f, default_flow_style=False, sort_keys=False, allow_unicode=True, width=120)
    
    # Calculate new coverage
    with_desc = sum(1 for s in schemas.values() if isinstance(s, dict) and s.get('description', '').strip())
    total = len([s for s in schemas.values() if isinstance(s, dict)])
    coverage = (with_desc / total * 100) if total > 0 else 0
    
    print(f"\nNew description coverage: {with_desc}/{total} ({coverage:.1f}%)")

if __name__ == '__main__':
    add_descriptions_to_openapi()
