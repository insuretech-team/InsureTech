import os
import yaml
from pathlib import Path

# Import our custom modules
try:
    from file_organizer import FileOrganizer
    from name_transformer import NameTransformer
    from description_loader import DescriptionLoader
except ImportError:
    # Fallback if not available
    FileOrganizer = None
    NameTransformer = None
    DescriptionLoader = None

# Proto Field Types
TYPE_DOUBLE = 1
TYPE_FLOAT = 2
TYPE_INT64 = 3
TYPE_UINT64 = 4
TYPE_INT32 = 5
TYPE_FIXED64 = 6
TYPE_FIXED32 = 7
TYPE_BOOL = 8
TYPE_STRING = 9
TYPE_GROUP = 10
TYPE_MESSAGE = 11
TYPE_BYTES = 12
TYPE_UINT32 = 13
TYPE_ENUM = 14
TYPE_SFIXED32 = 15
TYPE_SFIXED64 = 16
TYPE_SINT32 = 17
TYPE_SINT64 = 18

class SchemaGenerator:
    def __init__(self, registry, output_dir=None, descriptions_dir=None):
        self.registry = registry
        self.output_dir = output_dir
        
        # Initialize optional components
        self.file_organizer = FileOrganizer(output_dir) if FileOrganizer and output_dir else None
        self.name_transformer = NameTransformer() if NameTransformer else None
        self.description_loader = DescriptionLoader(descriptions_dir) if DescriptionLoader and descriptions_dir else None
        
        # Ensure directories exist
        if self.file_organizer:
            self.file_organizer.ensure_directories()

    def generate_schema(self, message_data, parser):
        """
        Generates an OpenAPI schema dictionary for a message.
        args:
            message_data: dict from ProtoParser ({'descriptor': ..., 'full_name': ..., 'comment': ...})
            parser: ProtoParser instance to extract fields
        Returns:
            Tuple of (transformed_name, schema, output_path)
        """
        full_name = message_data['full_name']
        original_name = message_data['descriptor'].name
        
        # Transform name (verbs to nouns)
        transformed_name = original_name
        if self.name_transformer:
            transformed_name = self.name_transformer.transform(original_name)
            if transformed_name != original_name:
                print(f"  Transformed: {original_name} → {transformed_name}")
        
        # Load rich description
        proto_comment = message_data.get('comment', '')
        if self.description_loader:
            description = self.description_loader.load_schema_description(full_name, proto_comment)
        else:
            description = proto_comment.strip()
        
        # If no description from proto, generate intelligent one
        if not description or not description.strip():
            description = self._generate_smart_description(transformed_name, message_data)
        
        schema = {
            "type": "object",
            "properties": {},
        }
        
        # Only add description if it's not empty
        if description and description.strip():
            schema["description"] = description.strip()
        
        fields = parser.extract_fields(message_data)
        required_fields = []

        for field in fields:
            prop_name = field['name']
            prop_schema = self._map_type(field)
            has_ref = '$ref' in prop_schema
            default_value_to_add = None
            additional_properties = {}  # Properties to add with allOf if has_ref (reset per field)
            
            # Load field description - only add if not empty
            field_comment = field.get('comment', '')
            field_desc = None
            if self.description_loader:
                field_desc = self.description_loader.load_field_description(
                    transformed_name, prop_name, field_comment
                )
            elif field_comment:
                field_desc = field_comment.strip()
            
            # Only add description if it has content
            if field_desc and field_desc.strip():
                prop_schema['description'] = field_desc.strip()
            
            # Add security annotations as OpenAPI extensions
            # If field has $ref, collect these in additional_properties instead
            def add_property(key, value):
                if has_ref:
                    additional_properties[key] = value
                else:
                    prop_schema[key] = value
            
            if field.get('is_sensitive'):
                add_property('x-sensitive', True)
                # Sensitive fields should be write-only (not returned in responses)
                add_property('writeOnly', True)
            
            if field.get('is_pii'):
                add_property('x-pii', True)
            
            if field.get('is_encrypted'):
                add_property('x-encrypted', True)
            
            if field.get('is_log_masked'):
                add_property('x-log-masked', True)
            
            if field.get('is_log_redacted'):
                add_property('x-log-redacted', True)
            
            if field.get('requires_consent'):
                add_property('x-requires-consent', True)
            
            if field.get('data_purpose'):
                add_property('x-data-purpose', field.get('data_purpose'))
            
            if field.get('retention_days'):
                add_property('x-retention-days', field.get('retention_days'))
            
            # Google API field behaviors
            if field.get('is_output_only'):
                add_property('readOnly', True)
            
            if field.get('is_input_only'):
                add_property('writeOnly', True)
            
            if field.get('is_immutable'):
                add_property('x-immutable', True)
                # Immutable fields can be set on create but not updated
                add_property('x-immutable-after-creation', True)
            
            # Database column annotations (relevant for OpenAPI)
            if field.get('unique'):
                add_property('x-unique', True)
            
            if field.get('primary_key'):
                add_property('x-primary-key', True)
            
            if field.get('auto_increment'):
                # Auto-increment fields are typically read-only
                if not prop_schema.get('readOnly') and not has_ref:
                    prop_schema['readOnly'] = True
                add_property('x-auto-increment', True)
            
            if field.get('default_value'):
                # Convert default value to correct type based on field type
                default_val = field.get('default_value')
                field_type = prop_schema.get('type') if not has_ref else None
                
                # Skip SQL functions and database-specific defaults
                # These are for database schema, not API defaults
                sql_functions = ['now()', 'current_timestamp', 'gen_random_uuid()', 
                                'uuid_generate_v4()', 'current_date', 'current_time']
                if any(func in default_val.lower() for func in sql_functions):
                    # Add as extension for documentation but not as default
                    add_property('x-database-default', default_val)
                elif has_ref:
                    # For $ref (enum), NEVER add default to API spec with allOf
                    # Enum defaults from database schema often use short names (e.g., 'IDLE')
                    # but OpenAPI needs full enum names (e.g., 'AGENT_STATUS_IDLE')
                    # Since these are database defaults, not API defaults, we only add as extension
                    # IMPORTANT: Do NOT add to additional_properties to avoid allOf + default issues
                    add_property('x-database-default', default_val.strip("'\""))
                elif field_type == 'integer':
                    try:
                        prop_schema['default'] = int(default_val)
                    except (ValueError, TypeError):
                        # If conversion fails, skip default
                        pass
                elif field_type == 'number':
                    try:
                        prop_schema['default'] = float(default_val)
                    except (ValueError, TypeError):
                        pass
                elif field_type == 'boolean':
                    if isinstance(default_val, str):
                        prop_schema['default'] = default_val.lower() in ('true', '1', 'yes')
                    else:
                        prop_schema['default'] = bool(default_val)
                else:
                    # String or other types - but validate it's not a SQL function
                    if not any(char in default_val for char in ['(', ')', '{', '}']):
                        prop_schema['default'] = default_val
                    else:
                        # Looks like a function call, add as extension
                        add_property('x-database-default', default_val)
            
            if field.get('db_encrypted'):
                # Database-level encryption (different from transport encryption)
                add_property('x-encrypted-at-rest', True)
            
            # Check if required (parse field_behavior annotations or not_null)
            if self._is_required(field) or field.get('not_null'):
                required_fields.append(prop_name)
            
            # If we have a $ref with additional properties, use allOf pattern
            # IMPORTANT: Never add 'default' with allOf as it causes OpenAPI spec errors
            if has_ref and additional_properties:
                # OpenAPI doesn't allow additional properties alongside $ref
                # Use allOf to combine them
                original_ref = prop_schema['$ref']
                prop_schema = {
                    "allOf": [{"$ref": original_ref}]
                }
                # Add all additional properties at the same level as allOf
                # but filter out 'default' to prevent spec errors
                filtered_props = {k: v for k, v in additional_properties.items() if k != 'default'}
                prop_schema.update(filtered_props)
            
            # Handle Map fields (map<K,V> in proto)
            # Maps are represented as repeated message fields with special Entry type
            if field.get('is_map'):
                # Map fields should be converted to object with additionalProperties
                # Extract description before converting
                item_desc = prop_schema.get('description', 'Map field (key-value pairs)')
                
                prop_schema = {
                    "type": "object",
                    "additionalProperties": {"type": "string"},  # Default to string
                    "description": item_desc
                }
            # Handle normal Repeated fields (arrays)
            elif field['label'] == 3:  # LABEL_REPEATED
                prop_schema = {
                    "type": "array",
                    "items": prop_schema
                }
                if 'description' in prop_schema.get('items', {}):
                    prop_schema['description'] = prop_schema['items']['description']

            schema['properties'][prop_name] = prop_schema
        
        # Add required array
        if required_fields:
            schema['required'] = required_fields
        
        # Determine output path
        output_path = None
        if self.file_organizer and self.output_dir:
            output_path, category = self.file_organizer.get_output_path(full_name, transformed_name)
        
        return transformed_name, schema, output_path
    
    def _is_required(self, field):
        """
        Check if field has REQUIRED annotation from google.api.field_behavior
        
        Parses field options to find field_behavior extension with REQUIRED value
        """
        # Check if field has 'options' key (from proto parser)
        if 'options' not in field or not field['options']:
            return False
        
        options = field['options']
        
        # Check for field_behavior annotation
        # The field_behavior extension is in google.api and has number 1052
        # It's a repeated enum with values: OPTIONAL=0, REQUIRED=1, OUTPUT_ONLY=2, etc.
        
        try:
            # Try to access field_behavior from options
            # This depends on how the proto parser exposes options
            if hasattr(options, 'Extensions'):
                # If using google.protobuf descriptor
                from google.api import field_behavior_pb2
                behaviors = options.Extensions[field_behavior_pb2.field_behavior]
                # REQUIRED = 2 in the enum
                return field_behavior_pb2.FieldBehavior.REQUIRED in behaviors
            
            # Alternative: Check if proto parser provides it as dict
            if isinstance(options, dict):
                field_behavior = options.get('field_behavior', [])
                # Check if 'REQUIRED' or value 2 is in the list
                return 'REQUIRED' in field_behavior or 2 in field_behavior
                
        except (ImportError, AttributeError, KeyError):
            # If google.api not available or field_behavior not found, check field comments
            # Sometimes proto comments indicate required fields
            field_comment = field.get('comment', '').lower()
            if 'required' in field_comment:
                return True
        
        return False

    def generate_enum_schema(self, enum_data):
        """
        Generates an OpenAPI schema for an Enum.
        Returns:
            Tuple of (enum_name, schema, output_path)
        """
        descriptor = enum_data['descriptor']
        full_name = enum_data['full_name']
        enum_name = descriptor.name
        
        # Enums don't need name transformation
        values = [v.name for v in descriptor.value]
        
        # Load description
        proto_comment = enum_data.get('comment', '')
        if self.description_loader:
            description = self.description_loader.load_schema_description(full_name, proto_comment)
        else:
            description = proto_comment.strip()
        
        # If no description, generate smart one
        if not description or not description.strip():
            description = self._generate_enum_description(enum_name, values)
        
        schema = {
            "type": "string",
            "enum": values,
            "description": description
        }
        
        # Determine output path
        output_path = None
        if self.file_organizer and self.output_dir:
            output_path, category = self.file_organizer.get_output_path(full_name, enum_name)
        
        return enum_name, schema, output_path
    
    def write_schema_file(self, name, schema, output_path):
        """Write schema to file"""
        if not output_path:
            return
        
        # Convert any tuples to lists before writing (enums can create tuples)
        def sanitize_for_yaml(obj):
            if isinstance(obj, tuple):
                return list(obj)
            elif isinstance(obj, dict):
                return {k: sanitize_for_yaml(v) for k, v in obj.items()}
            elif isinstance(obj, list):
                return [sanitize_for_yaml(item) for item in obj]
            return obj
        
        sanitized_schema = sanitize_for_yaml(schema)
        
        # Ensure directory exists
        Path(output_path).parent.mkdir(parents=True, exist_ok=True)
        
        # Write YAML file
        with open(output_path, 'w', encoding='utf-8') as f:
            yaml.dump({name: sanitized_schema}, f, sort_keys=False, allow_unicode=True)

    def _map_type(self, field):
        """Maps proto field type to OpenAPI type."""
        t = field['type']
        
        if t in (TYPE_DOUBLE, TYPE_FLOAT):
            return {"type": "number", "format": "double" if t == TYPE_DOUBLE else "float"}
        
        if t in (TYPE_INT32, TYPE_UINT32, TYPE_SINT32, TYPE_FIXED32, TYPE_SFIXED32):
            return {"type": "integer", "format": "int32"}
            
        if t in (TYPE_INT64, TYPE_UINT64, TYPE_SINT64, TYPE_FIXED64, TYPE_SFIXED64):
            return {"type": "string", "format": "int64"} # Google APIs often use string for 64-bit int
            # Or use type: integer, format: int64
            
        if t == TYPE_BOOL:
            return {"type": "boolean"}
            
        if t == TYPE_STRING:
            return {"type": "string"}
            
        if t == TYPE_BYTES:
            return {"type": "string", "format": "byte"}
            
        if t == TYPE_MESSAGE or t == TYPE_ENUM:
            # Resolve Reference
            # field['type_name'] is like ".insuretech.common.v1.Error"
            full_name = field['type_name']
            if full_name.startswith('.'):
                full_name = full_name[1:]
                
            # Wrapper types or WKTs check FIRST
            if full_name == "google.protobuf.Timestamp":
                 return {"type": "string", "format": "date-time"}
            if full_name == "google.protobuf.Struct":
                 return {"type": "object"}
            if full_name == "google.protobuf.Duration":
                 return {"type": "string"}
            if full_name == "google.protobuf.Any":
                 return {"type": "object"}
            # FieldMask, Empty, etc.
                 
            ref = self.registry.get_ref(full_name)
            if ref:
                return {"$ref": ref}
            else:
                # Debug: why did lookup fail?
                print(f"    WARNING: No ref found for {full_name}")
                # Fallback if not found (e.g. map entry or external like google.protobuf.Struct)
                if full_name == "google.protobuf.Struct":
                     return {"type": "object"}
                if full_name == "google.protobuf.Timestamp":
                     return {"type": "string", "format": "date-time"}
                
                # Check if it's a Map Entry (auto-generated by protobuf for map<K,V>)
                # Map entries have names like "SomeMessage.SomeFieldEntry"
                # These should be converted to proper OpenAPI objects with additionalProperties
                if full_name.endswith('Entry') and '.' in full_name:
                    # This is a protobuf map entry - convert to object with additionalProperties
                    # Map<string, string> becomes: type: object, additionalProperties: {type: string}
                    # For now, we'll use a generic object with string values
                    # In the future, we could parse the actual key/value types
                    return {
                        "type": "object",
                        "additionalProperties": {"type": "string"},
                        "description": f"Map field (key-value pairs)"
                    }
                
                # Unknown type - return generic object
                return {"type": "object", "description": f"Unknown type: {full_name}"}

        return {"type": "string"} # Default fallback (e.g. unknown)
    
    def _generate_smart_description(self, schema_name, message_data):
        """Generate intelligent description based on naming patterns"""
        import re
        
        # For Request DTOs
        if schema_name.endswith('Request'):
            base = schema_name.replace('Request', '')
            words = re.findall(r'[A-Z][a-z]+|[A-Z]+(?=[A-Z][a-z]|\b)', base)
            operation = ' '.join(words).lower()
            
            if 'Create' in base or 'Creation' in base:
                return f"Request payload for creating a new {operation.replace('create ', '').replace('creation ', '')}. Contains all required fields and optional parameters for initialization."
            elif 'Update' in base:
                return f"Request payload for updating an existing {operation.replace('update ', '')}. Contains fields to be modified."
            elif 'Get' in base or 'Retrieve' in base or 'Retrieval' in base:
                return f"Request payload for retrieving {operation.replace('get ', '').replace('retrieve ', '').replace('retrieval ', '')} information. May include filters and pagination parameters."
            elif 'List' in base:
                return f"Request payload for listing {operation.replace('list ', '')} items. Supports filtering, sorting, and pagination."
            else:
                return f"Request payload for {operation} operation. Contains parameters required to execute the operation."
        
        # For Response DTOs
        elif schema_name.endswith('Response'):
            base = schema_name.replace('Response', '')
            words = re.findall(r'[A-Z][a-z]+|[A-Z]+(?=[A-Z][a-z]|\b)', base)
            operation = ' '.join(words).lower()
            return f"Response payload for {operation} operation. Returns operation results and status."
        
        # For Events
        elif schema_name.endswith('Event'):
            base = schema_name.replace('Event', '')
            words = re.findall(r'[A-Z][a-z]+|[A-Z]+(?=[A-Z][a-z]|\b)', base)
            if len(words) >= 2:
                action = words[-1].lower()
                entity = ' '.join(words[:-1]).lower()
                action_map = {
                    'created': 'is created', 'updated': 'is updated', 'deleted': 'is deleted',
                    'activated': 'is activated', 'deactivated': 'is deactivated',
                    'approved': 'is approved', 'rejected': 'is rejected',
                    'completed': 'is completed', 'failed': 'fails', 'started': 'starts',
                    'ended': 'ends', 'cancelled': 'is cancelled', 'expired': 'expires'
                }
                action_phrase = action_map.get(action, f'undergoes {action}')
                return f"Event emitted when {entity} {action_phrase}. Published to event stream for downstream processing and audit trail."
            return f"Event message for {' '.join(words).lower()} notification."
        
        # For Entities
        else:
            words = re.findall(r'[A-Z][a-z]+|[A-Z]+(?=[A-Z][a-z]|\b)', schema_name)
            entity_name = ' '.join(words).lower()
            return f"Domain entity representing {entity_name}. Core business object in the system."
    
    def _generate_enum_description(self, enum_name, values):
        """Generate description for enum based on name and values"""
        import re
        
        words = re.findall(r'[A-Z][a-z]+|[A-Z]+(?=[A-Z][a-z]|\b)', enum_name)
        field_name = ' '.join(words).lower()
        
        # Create preview of values (first 5)
        value_preview = ', '.join(values[:5])
        if len(values) > 5:
            value_preview += f' (and {len(values) - 5} more)'
        
        # Generate description based on suffix
        if enum_name.endswith('Status'):
            base = field_name.replace(' status', '')
            return f"Status values for {base}. Defines lifecycle states: {value_preview}."
        elif enum_name.endswith('Type'):
            base = field_name.replace(' type', '')
            return f"Type categorization for {base}. Options: {value_preview}."
        elif enum_name.endswith('Category'):
            base = field_name.replace(' category', '')
            return f"Category classification for {base}. Categories: {value_preview}."
        elif enum_name.endswith('Method'):
            base = field_name.replace(' method', '')
            return f"Available methods for {base}. Methods: {value_preview}."
        elif enum_name.endswith('Level'):
            base = field_name.replace(' level', '')
            return f"Severity or priority levels for {base}. Levels: {value_preview}."
        elif enum_name.endswith('Priority'):
            base = field_name.replace(' priority', '')
            return f"Priority levels for {base}. Priorities: {value_preview}."
        elif enum_name.endswith('Action'):
            base = field_name.replace(' action', '')
            return f"Available actions for {base}. Actions: {value_preview}."
        elif enum_name.endswith('Role'):
            base = field_name.replace(' role', '')
            return f"User roles for {base}. Roles: {value_preview}."
        else:
            return f"Enumeration of {field_name} values. Valid options: {value_preview}."
