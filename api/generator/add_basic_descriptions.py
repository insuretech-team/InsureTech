"""
Add basic auto-generated descriptions and required fields to schemas
"""
import yaml

def generate_description(schema_name, schema_def):
    """Generate a basic description from schema name and type"""
    name = schema_name
    
    # Categorize by suffix/pattern
    if name.endswith('Request'):
        base = name.replace('Request', '')
        return f"Request payload for {base} operation"
    elif name.endswith('Response'):
        base = name.replace('Response', '')
        return f"Response payload for {base} operation"
    elif name.endswith('Event'):
        base = name.replace('Event', '')
        return f"Event emitted when {base} occurs"
    elif name.endswith('Info'):
        return f"{name.replace('Info', '')} information"
    elif name.endswith('Config'):
        return f"Configuration for {name.replace('Config', '')}"
    elif 'Pagination' in name:
        return f"Pagination {name.replace('Pagination', '').lower()} information"
    else:
        # Entity or enum
        if schema_def.get('type') == 'string' and 'enum' in schema_def:
            return f"Enumeration for {name}"
        else:
            return f"{name} entity"

def add_required_fields(schema_def):
    """Mark id fields and key fields as required"""
    if 'properties' not in schema_def:
        return []
    
    required = []
    props = schema_def['properties']
    
    for field_name, field_def in props.items():
        # Mark these as required
        if field_name in ['id', 'tenant_id']:
            required.append(field_name)
        # For request DTOs, mark common required fields
        elif field_name in ['user_id', 'policy_id', 'claim_id'] and 'Request' in str(schema_def):
            required.append(field_name)
    
    return required

print("Loading openapi.yaml...")
with open('../openapi.yaml', 'r', encoding='utf-8') as f:
    spec = yaml.safe_load(f)

schemas = spec['components']['schemas']
desc_added = 0
req_added = 0

for schema_name, schema_def in schemas.items():
    if not isinstance(schema_def, dict):
        continue
    
    # Add description if missing
    if not schema_def.get('description', '').strip():
        desc = generate_description(schema_name, schema_def)
        schema_def['description'] = desc
        desc_added += 1
    
    # Add required fields if it's an object with properties
    if schema_def.get('type') == 'object' and 'properties' in schema_def:
        required = add_required_fields(schema_def)
        if required and 'required' not in schema_def:
            schema_def['required'] = required
            req_added += 1

print(f"Added descriptions to: {desc_added} schemas")
print(f"Added required fields to: {req_added} schemas")

# Write back
print("Writing updated openapi.yaml...")
with open('../openapi.yaml', 'w', encoding='utf-8') as f:
    yaml.dump(spec, f, default_flow_style=False, sort_keys=False, allow_unicode=True, width=120)

# Calculate coverage
with_desc = sum(1 for s in schemas.values() if isinstance(s, dict) and s.get('description', '').strip())
total = len([s for s in schemas.values() if isinstance(s, dict)])
coverage = (with_desc / total * 100) if total > 0 else 0

print(f"\nFinal description coverage: {with_desc}/{total} ({coverage:.1f}%)")
print("Done!")
