"""
Generate intelligent descriptions for all schemas based on naming patterns and structure
"""
import yaml
import re

class SmartDescriptionGenerator:
    def __init__(self, openapi_file):
        self.openapi_file = openapi_file
        self.spec = None
        self.schemas = None
    
    def load_spec(self):
        """Load OpenAPI spec"""
        with open(self.openapi_file, 'r', encoding='utf-8') as f:
            self.spec = yaml.safe_load(f)
        self.schemas = self.spec['components']['schemas']
    
    def generate_request_description(self, name):
        """Generate description for Request DTOs"""
        # Remove 'Request' suffix
        base = name.replace('Request', '')
        
        # Split camelCase to words
        words = re.findall(r'[A-Z][a-z]+|[A-Z]+(?=[A-Z][a-z]|\b)', base)
        operation = ' '.join(words).lower()
        
        # Generate based on common patterns
        if 'Create' in base or 'Creation' in base:
            return f"Request payload for creating a new {operation.replace('create ', '').replace('creation ', '')}. Contains all required fields and optional parameters for initialization."
        elif 'Update' in base or 'Modification' in base:
            return f"Request payload for updating an existing {operation.replace('update ', '').replace('modification ', '')}. Contains fields to be modified."
        elif 'Delete' in base or 'Removal' in base:
            return f"Request payload for deleting a {operation.replace('delete ', '').replace('removal ', '')}. Requires identifier for the resource to be removed."
        elif 'Get' in base or 'Retrieve' in base or 'Retrieval' in base:
            return f"Request payload for retrieving {operation.replace('get ', '').replace('retrieve ', '').replace('retrieval ', '')} information. May include filters and pagination parameters."
        elif 'List' in base:
            return f"Request payload for listing {operation.replace('list ', '')} items. Supports filtering, sorting, and pagination."
        elif 'Search' in base:
            return f"Request payload for searching {operation.replace('search ', '')} records. Includes search criteria and filters."
        elif 'Generate' in base or 'Generation' in base:
            return f"Request payload for generating {operation.replace('generate ', '').replace('generation ', '')}. Contains generation parameters and options."
        elif 'Calculate' in base or 'Calculation' in base:
            return f"Request payload for calculating {operation.replace('calculate ', '').replace('calculation ', '')}. Provides input values for computation."
        elif 'Evaluate' in base or 'Evaluation' in base:
            return f"Request payload for evaluating {operation.replace('evaluate ', '').replace('evaluation ', '')}. Contains criteria for assessment."
        elif 'Analyze' in base or 'Analysis' in base:
            return f"Request payload for analyzing {operation.replace('analyze ', '').replace('analysis ', '')}. Includes data to be analyzed."
        elif 'Process' in base or 'Processing' in base:
            return f"Request payload for processing {operation.replace('process ', '').replace('processing ', '')}. Contains items to be processed."
        else:
            return f"Request payload for {operation} operation. Contains parameters required to execute the operation."
    
    def generate_response_description(self, name):
        """Generate description for Response DTOs"""
        base = name.replace('Response', '')
        words = re.findall(r'[A-Z][a-z]+|[A-Z]+(?=[A-Z][a-z]|\b)', base)
        operation = ' '.join(words).lower()
        
        if 'Create' in base or 'Creation' in base:
            return f"Response payload for {operation} operation. Returns the created resource with assigned identifiers and timestamps."
        elif 'Update' in base or 'Modification' in base:
            return f"Response payload for {operation} operation. Returns the updated resource with new values."
        elif 'Delete' in base or 'Removal' in base:
            return f"Response payload for {operation} operation. Confirms successful deletion and returns status."
        elif 'Get' in base or 'Retrieve' in base or 'Retrieval' in base:
            return f"Response payload for {operation} operation. Returns the requested resource with all details."
        elif 'List' in base:
            return f"Response payload for {operation} operation. Returns a paginated list of items with metadata."
        elif 'Search' in base:
            return f"Response payload for {operation} operation. Returns search results with relevance scores."
        elif 'Generate' in base or 'Generation' in base:
            return f"Response payload for {operation} operation. Returns the generated content or reference."
        elif 'Calculate' in base or 'Calculation' in base:
            return f"Response payload for {operation} operation. Returns calculated values and breakdown."
        elif 'Evaluate' in base or 'Evaluation' in base:
            return f"Response payload for {operation} operation. Returns evaluation results and metrics."
        elif 'Analyze' in base or 'Analysis' in base:
            return f"Response payload for {operation} operation. Returns analysis results with insights."
        else:
            return f"Response payload for {operation} operation. Returns operation results and status."
    
    def generate_event_description(self, name):
        """Generate description for Event messages"""
        base = name.replace('Event', '')
        words = re.findall(r'[A-Z][a-z]+|[A-Z]+(?=[A-Z][a-z]|\b)', base)
        
        # Extract entity and action
        if len(words) >= 2:
            action = words[-1].lower()
            entity = ' '.join(words[:-1]).lower()
            
            action_map = {
                'created': 'is created',
                'updated': 'is updated',
                'deleted': 'is deleted',
                'activated': 'is activated',
                'deactivated': 'is deactivated',
                'approved': 'is approved',
                'rejected': 'is rejected',
                'completed': 'is completed',
                'failed': 'fails',
                'started': 'starts',
                'ended': 'ends',
                'cancelled': 'is cancelled',
                'expired': 'expires'
            }
            
            action_phrase = action_map.get(action, f'undergoes {action}')
            return f"Event emitted when {entity} {action_phrase}. Published to event stream for downstream processing and audit trail."
        
        return f"Event message for {' '.join(words).lower()} notification. Published when the event occurs."
    
    def generate_entity_description(self, name, schema_def):
        """Generate description for Entity schemas"""
        words = re.findall(r'[A-Z][a-z]+|[A-Z]+(?=[A-Z][a-z]|\b)', name)
        entity_name = ' '.join(words).lower()
        
        # Check if it has typical entity fields
        props = schema_def.get('properties', {})
        has_id = 'id' in props or any('_id' in p for p in props)
        has_timestamps = 'created_at' in props or 'updated_at' in props
        
        if has_id and has_timestamps:
            return f"Domain entity representing {entity_name}. Persisted in database with full audit trail and lifecycle management."
        elif has_id:
            return f"Domain entity representing {entity_name}. Core business object in the system."
        else:
            return f"Data structure for {entity_name}. Used for data transfer and processing."
    
    def generate_enum_description(self, name, schema_def):
        """Generate description for Enum schemas"""
        words = re.findall(r'[A-Z][a-z]+|[A-Z]+(?=[A-Z][a-z]|\b)', name)
        field_name = ' '.join(words).lower()
        
        enum_values = schema_def.get('enum', [])
        value_preview = ', '.join(enum_values[:5])
        if len(enum_values) > 5:
            value_preview += f', ... ({len(enum_values)} total)'
        
        if 'Status' in name:
            return f"Enumeration of possible status values for {field_name.replace(' status', '')}. Defines lifecycle states: {value_preview}."
        elif 'Type' in name:
            return f"Enumeration of {field_name.replace(' type', '')} types. Categorizes by: {value_preview}."
        elif 'Category' in name:
            return f"Enumeration of {field_name.replace(' category', '')} categories. Available categories: {value_preview}."
        elif 'Method' in name:
            return f"Enumeration of {field_name.replace(' method', '')} methods. Supported methods: {value_preview}."
        elif 'Level' in name:
            return f"Enumeration of {field_name.replace(' level', '')} levels. Defines severity or priority: {value_preview}."
        else:
            return f"Enumeration of possible {field_name} values. Valid options: {value_preview}."
    
    def generate_description(self, schema_name, schema_def):
        """Generate appropriate description based on schema type"""
        if schema_name.endswith('Request'):
            return self.generate_request_description(schema_name)
        elif schema_name.endswith('Response'):
            return self.generate_response_description(schema_name)
        elif schema_name.endswith('Event'):
            return self.generate_event_description(schema_name)
        elif schema_def.get('type') == 'string' and 'enum' in schema_def:
            return self.generate_enum_description(schema_name, schema_def)
        elif schema_def.get('type') == 'object':
            return self.generate_entity_description(schema_name, schema_def)
        else:
            return f"{schema_name} data structure"
    
    def add_required_fields(self, schema_name, schema_def):
        """Add required fields based on common patterns"""
        if schema_def.get('type') != 'object' or 'properties' not in schema_def:
            return []
        
        required = []
        props = schema_def['properties']
        
        # Always require id if present
        if 'id' in props:
            required.append('id')
        
        # Tenant isolation
        if 'tenant_id' in props:
            required.append('tenant_id')
        
        # For Request DTOs
        if schema_name.endswith('Request'):
            # Require foreign keys
            for field in props:
                if field.endswith('_id') and field not in ['id', 'tenant_id']:
                    # Make primary entity IDs required
                    if any(key in field for key in ['user', 'policy', 'claim', 'agent', 'product']):
                        required.append(field)
            
            # Require message/content fields
            if 'message' in props:
                required.append('message')
            if 'content' in props:
                required.append('content')
        
        # Remove duplicates
        return list(set(required))
    
    def process_all_schemas(self):
        """Process all schemas and add descriptions"""
        desc_added = 0
        req_added = 0
        
        for schema_name, schema_def in self.schemas.items():
            if not isinstance(schema_def, dict):
                continue
            
            # Add description if missing
            current_desc = schema_def.get('description', '').strip()
            if not current_desc:
                desc = self.generate_description(schema_name, schema_def)
                schema_def['description'] = desc
                desc_added += 1
            
            # Add required fields
            required = self.add_required_fields(schema_name, schema_def)
            if required and 'required' not in schema_def:
                schema_def['required'] = required
                req_added += 1
        
        return desc_added, req_added
    
    def save_spec(self):
        """Save updated spec"""
        with open(self.openapi_file, 'w', encoding='utf-8') as f:
            yaml.dump(self.spec, f, default_flow_style=False, sort_keys=False, 
                     allow_unicode=True, width=120)
    
    def calculate_coverage(self):
        """Calculate description coverage"""
        with_desc = 0
        total = 0
        
        for schema_def in self.schemas.values():
            if isinstance(schema_def, dict):
                total += 1
                if schema_def.get('description', '').strip():
                    with_desc += 1
        
        coverage = (with_desc / total * 100) if total > 0 else 0
        return with_desc, total, coverage

if __name__ == '__main__':
    generator = SmartDescriptionGenerator('../openapi.yaml')
    
    print("Loading OpenAPI spec...")
    generator.load_spec()
    
    print(f"Processing {len(generator.schemas)} schemas...")
    desc_added, req_added = generator.process_all_schemas()
    
    print(f"\n✓ Added descriptions to: {desc_added} schemas")
    print(f"✓ Added required fields to: {req_added} schemas")
    
    print("\nSaving updated spec...")
    generator.save_spec()
    
    with_desc, total, coverage = generator.calculate_coverage()
    print(f"\n✓ Final coverage: {with_desc}/{total} ({coverage:.1f}%)")
    print("\n✓ Smart description generation complete!")
