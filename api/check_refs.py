import yaml
import re

with open('openapi.yaml') as f:
    spec = yaml.safe_load(f)

# Get all defined schemas
defined = set(spec['components']['schemas'].keys())

# Find all $ref references
refs_found = set()

def find_refs(obj):
    if isinstance(obj, dict):
        if '$ref' in obj:
            ref = obj['$ref']
            if ref.startswith('#/components/schemas/'):
                schema_name = ref.replace('#/components/schemas/', '')
                refs_found.add(schema_name)
        for v in obj.values():
            find_refs(v)
    elif isinstance(obj, list):
        for item in obj:
            find_refs(item)

find_refs(spec)

# Find broken references
broken = refs_found - defined

print(f'Total references: {len(refs_found)}')
print(f'Defined schemas: {len(defined)}')
print(f'Broken references: {len(broken)}')

if broken:
    print('\nBroken references:')
    for ref in sorted(broken):
        print(f'  - {ref}')
