"""
Update individual schema files with descriptions from openapi.yaml
Then regenerate the full spec
"""
import os
import yaml

print("Loading openapi.yaml with descriptions...")
with open('../openapi.yaml', 'r', encoding='utf-8') as f:
    spec = yaml.safe_load(f)

schemas = spec['components']['schemas']
updated_count = 0

# Update individual schema files
print("\nUpdating individual schema files...")

# Update schemas directory
for root, dirs, files in os.walk('../schemas'):
    for file in files:
        if file.endswith('.yaml'):
            file_path = os.path.join(root, file)
            
            with open(file_path, 'r', encoding='utf-8') as f:
                file_data = yaml.safe_load(f)
            
            if not file_data:
                continue
            
            modified = False
            for schema_name, schema_def in file_data.items():
                # Get description from openapi.yaml
                if schema_name in schemas and isinstance(schema_def, dict):
                    openapi_desc = schemas[schema_name].get('description', '').strip()
                    current_desc = schema_def.get('description', '').strip()
                    
                    if openapi_desc and not current_desc:
                        schema_def['description'] = openapi_desc
                        modified = True
                        print(f"  Updated: {schema_name}")
                    
                    # Also update required fields
                    if schema_name in schemas and 'required' in schemas[schema_name]:
                        if 'required' not in schema_def or not schema_def['required']:
                            schema_def['required'] = schemas[schema_name]['required']
                            modified = True
            
            if modified:
                with open(file_path, 'w', encoding='utf-8') as f:
                    yaml.dump(file_data, f, default_flow_style=False, sort_keys=False, allow_unicode=True)
                updated_count += 1

# Update events directory
for root, dirs, files in os.walk('../events'):
    for file in files:
        if file.endswith('.yaml'):
            file_path = os.path.join(root, file)
            
            with open(file_path, 'r', encoding='utf-8') as f:
                file_data = yaml.safe_load(f)
            
            if not file_data:
                continue
            
            modified = False
            for schema_name, schema_def in file_data.items():
                if schema_name in schemas and isinstance(schema_def, dict):
                    openapi_desc = schemas[schema_name].get('description', '').strip()
                    current_desc = schema_def.get('description', '').strip()
                    
                    if openapi_desc and not current_desc:
                        schema_def['description'] = openapi_desc
                        modified = True
                        print(f"  Updated: {schema_name}")
            
            if modified:
                with open(file_path, 'w', encoding='utf-8') as f:
                    yaml.dump(file_data, f, default_flow_style=False, sort_keys=False, allow_unicode=True)
                updated_count += 1

# Update enum files
for file in os.listdir('../enums'):
    if file.endswith('.yaml'):
        file_path = os.path.join('../enums', file)
        
        with open(file_path, 'r', encoding='utf-8') as f:
            file_data = yaml.safe_load(f)
        
        if not file_data:
            continue
        
        modified = False
        for schema_name, schema_def in file_data.items():
            if schema_name in schemas and isinstance(schema_def, dict):
                openapi_desc = schemas[schema_name].get('description', '').strip()
                current_desc = schema_def.get('description', '').strip()
                
                if openapi_desc and not current_desc:
                    schema_def['description'] = openapi_desc
                    modified = True
                    print(f"  Updated: {schema_name}")
        
        if modified:
            with open(file_path, 'w', encoding='utf-8') as f:
                yaml.dump(file_data, f, default_flow_style=False, sort_keys=False, allow_unicode=True)
            updated_count += 1

print(f"\n✓ Updated {updated_count} schema files")
print("\nNow regenerating openapi.yaml...")
