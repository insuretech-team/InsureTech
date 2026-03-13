"""
Enhances generated schemas with descriptions from markdown files
"""
import os
import yaml
import re

class DescriptionEnhancer:
    def __init__(self, descriptions_dir, api_dir):
        self.descriptions_dir = descriptions_dir
        self.api_dir = api_dir
        
    def load_description(self, desc_file):
        """Load description from markdown file"""
        if not os.path.exists(desc_file):
            return None
            
        with open(desc_file, 'r', encoding='utf-8') as f:
            content = f.read()
            
        # Extract main description (after ## Description and before ## Fields)
        match = re.search(r'## Description\s+(.*?)\s+##', content, re.DOTALL)
        if match:
            desc = match.group(1).strip()
            # Remove comment lines and empty sections
            desc = re.sub(r'<!--.*?-->', '', desc, flags=re.DOTALL)
            desc = re.sub(r'\n\s*\n', '\n', desc).strip()
            # Take first paragraph if it's a real description
            if desc and not desc.startswith('**'):
                lines = desc.split('\n')
                # Get first non-empty, non-template line
                for line in lines:
                    line = line.strip()
                    if line and not line.startswith('**') and not line.startswith('-'):
                        return line
        
        # Fallback: look for proto comment or overview
        match = re.search(r'\*\*Proto Comment\*\*:\s*(.+?)(?:\n\n|\n\*\*)', content, re.DOTALL)
        if match:
            desc = match.group(1).strip()
            if desc and 'Maps to' in desc:
                # Take first sentence before "Maps to"
                desc = desc.split('Maps to')[0].strip()
            return desc
            
        return None
    
    def enhance_schema_file(self, schema_file, desc_file):
        """Add description to a schema YAML file"""
        if not os.path.exists(desc_file):
            return False
            
        description = self.load_description(desc_file)
        if not description:
            return False
            
        # Load schema
        with open(schema_file, 'r', encoding='utf-8') as f:
            schema_data = yaml.safe_load(f)
            
        if not schema_data:
            return False
            
        # Add description to schema
        modified = False
        for schema_name, schema_def in schema_data.items():
            if isinstance(schema_def, dict):
                # Only add if no description exists or it's empty
                current_desc = schema_def.get('description', '').strip()
                if not current_desc or current_desc == '':
                    schema_def['description'] = description
                    modified = True
                    
        if modified:
            # Write back
            with open(schema_file, 'w', encoding='utf-8') as f:
                yaml.dump(schema_data, f, default_flow_style=False, sort_keys=False, allow_unicode=True)
            return True
            
        return False
    
    def enhance_all_schemas(self):
        """Enhance all schemas with descriptions"""
        enhanced_count = 0
        
        # Enhance entity schemas
        entities_dir = os.path.join(self.api_dir, 'schemas', 'insuretech')
        desc_entities_dir = os.path.join(self.descriptions_dir, 'entity', 'insuretech')
        
        if os.path.exists(entities_dir):
            for root, dirs, files in os.walk(entities_dir):
                for file in files:
                    if file.endswith('.yaml'):
                        schema_file = os.path.join(root, file)
                        # Determine corresponding description file
                        rel_path = os.path.relpath(root, entities_dir)
                        desc_dir = os.path.join(desc_entities_dir, rel_path)
                        desc_file = os.path.join(desc_dir, file.replace('.yaml', '.md'))
                        
                        if self.enhance_schema_file(schema_file, desc_file):
                            enhanced_count += 1
                            print(f"  Enhanced: {file}")
        
        # Enhance event schemas
        events_dir = os.path.join(self.api_dir, 'events', 'insuretech')
        desc_events_dir = os.path.join(self.descriptions_dir, 'event', 'insuretech')
        
        if os.path.exists(events_dir):
            for root, dirs, files in os.walk(events_dir):
                for file in files:
                    if file.endswith('.yaml'):
                        schema_file = os.path.join(root, file)
                        rel_path = os.path.relpath(root, events_dir)
                        desc_dir = os.path.join(desc_events_dir, rel_path)
                        desc_file = os.path.join(desc_dir, file.replace('.yaml', '.md'))
                        
                        if self.enhance_schema_file(schema_file, desc_file):
                            enhanced_count += 1
                            print(f"  Enhanced: {file}")
        
        # Enhance DTO schemas (request/response)
        # DTOs are in schemas/insuretech/.../services/
        desc_dtos_dir = os.path.join(self.descriptions_dir, 'dto', 'insuretech')
        
        if os.path.exists(desc_dtos_dir):
            for root, dirs, files in os.walk(desc_dtos_dir):
                for file in files:
                    if file.endswith('.md'):
                        desc_file = os.path.join(root, file)
                        # Find corresponding schema file
                        rel_path = os.path.relpath(root, desc_dtos_dir)
                        schema_dir = os.path.join(entities_dir, rel_path)
                        schema_file = os.path.join(schema_dir, file.replace('.md', '.yaml'))
                        
                        if os.path.exists(schema_file):
                            if self.enhance_schema_file(schema_file, desc_file):
                                enhanced_count += 1
                                print(f"  Enhanced: {file}")
        
        return enhanced_count

if __name__ == '__main__':
    descriptions_dir = '../descriptions'
    api_dir = '..'
    
    enhancer = DescriptionEnhancer(descriptions_dir, api_dir)
    print("Enhancing schemas with descriptions...")
    count = enhancer.enhance_all_schemas()
    print(f"\nEnhanced {count} schemas")
