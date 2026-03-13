from google.protobuf.descriptor_pb2 import FileDescriptorSet

from google.protobuf.descriptor_pb2 import FileDescriptorSet, FieldDescriptorProto
import sys
import os

# Ensure generated google protos are importable
gen_path = os.path.join(os.path.dirname(os.path.abspath(__file__)), "gen")

annotations_pb2 = None
try:
    # First check if we have google.protobuf installed
    import google.protobuf
    
    # Add gen path to sys.path for namespace package extension
    if gen_path not in sys.path:
        sys.path.insert(0, gen_path)
    
    # Import google first to set up namespace
    import google
    # Extend the namespace package path to include our gen directory
    google_gen_path = os.path.join(gen_path, "google")
    if hasattr(google, '__path__'):
        # Convert to list if needed and add our path
        if not isinstance(google.__path__, list):
            google.__path__ = list(google.__path__)
        if google_gen_path not in google.__path__:
            google.__path__.insert(0, google_gen_path)
    
    # Now import google.api.annotations_pb2
    from google.api import annotations_pb2
except ImportError as e:
    print(f"Warning: Could not import google.api.annotations_pb2: {e}")
    annotations_pb2 = None
except Exception as e:
    print(f"Warning: Unexpected error loading annotations: {e}")
    annotations_pb2 = None

class ProtoParser:
    def __init__(self, proto_source_parser=None):
        self.file_desc_set = FileDescriptorSet()
        self.services = []
        self.messages = []
        self.enums = []
        self._comment_maps = {} # file_name -> { path_tuple: comment_str }
        self.proto_source_parser = proto_source_parser  # For extracting custom annotations

    def load_descriptor_set(self, filepath):
        """Loads a binary FileDescriptorSet from disk."""
        with open(filepath, 'rb') as f:
            self.file_desc_set.ParseFromString(f.read())
        self._extract_data()

    def _extract_data(self):
        """Iterates over files in the set and extracts messages/services."""
        for file_proto in self.file_desc_set.file:
            package = file_proto.package
            self._build_comment_map(file_proto)
            
            for service_index, service in enumerate(file_proto.service):
                methods = []
                for method_index, method in enumerate(service.method):
                    http_rule = None
                    is_sensitive = False
                    
                    # Try to extract google.api.http extension
                    if annotations_pb2 and method.options.HasExtension(annotations_pb2.http):
                        http_rule = method.options.Extensions[annotations_pb2.http]
                    
                    # Extract method comment from source_code_info
                    # Path for method: (6, service_index, 2, method_index)
                    # 6 = service field number in FileDescriptorProto
                    # 2 = method field number in ServiceDescriptorProto
                    method_path = (6, service_index, 2, method_index)
                    method_comment = self._get_comment(file_proto.name, method_path)
                    
                    # Check if method comment contains security markers
                    if method_comment and 'sensitive' in method_comment.lower():
                        is_sensitive = True
                    
                    methods.append({
                        'descriptor': method,
                        'http_rule': http_rule,
                        'name': method.name,
                        'input_type': method.input_type,
                        'output_type': method.output_type,
                        'is_sensitive': is_sensitive,
                        'comment': method_comment
                    })

                self.services.append({
                    'descriptor': service,
                    'package': package,
                    'file': file_proto.name,
                    'full_name': f"{package}.{service.name}" if package else service.name,
                    'methods': methods
                })
            
            # Extract Top-Level Enums
            for i, enum in enumerate(file_proto.enum_type):
                path = (5, i) # 5 = enum_type
                comment = self._get_comment(file_proto.name, path)
                self.enums.append({
                    'descriptor': enum,
                    'package': package,
                    'file': file_proto.name,
                    'full_name': f"{package}.{enum.name}" if package else enum.name,
                    'comment': comment,
                    'path': path
                })

            # Use index tracking for comment lookup
            for i, message in enumerate(file_proto.message_type):
                path = (4, i) # 4 = message_type
                comment = self._get_comment(file_proto.name, path)
                full_name = f"{package}.{message.name}" if package else message.name
                
                self.messages.append({
                    'descriptor': message,
                    'package': package,
                    'file': file_proto.name,
                    'full_name': full_name,
                    'comment': comment,
                    'path': path # location path for children
                })
                
                # Extract Nested Enums
                for j, enum in enumerate(message.enum_type):
                    enum_path = path + (4, j) # 4 = enum_type within message
                    enum_comment = self._get_comment(file_proto.name, enum_path)
                    self.enums.append({
                        'descriptor': enum,
                        'package': package,
                        'file': file_proto.name,
                        'full_name': f"{full_name}.{enum.name}",
                        'comment': enum_comment,
                        'path': enum_path
                    })

    def _build_comment_map(self, file_proto):
        """Builds a map of location paths to comments for a file."""
        m = {}
        if file_proto.source_code_info:
            for location in file_proto.source_code_info.location:
                if location.leading_comments or location.trailing_comments:
                    # comments can be leading or trailing, usually leading is the docblock
                    comment = (location.leading_comments + (location.trailing_comments or "")).strip()
                    m[tuple(location.path)] = comment
        self._comment_maps[file_proto.name] = m

    def _get_comment(self, filename, path):
        """Retrieves comments for a specific path."""
        return self._comment_maps.get(filename, {}).get(path, "")

    def get_services(self):
        return self.services

    def get_messages(self):
        return self.messages

    def get_enums(self):
        return self.enums

    def extract_fields(self, message_wrapper):
        """Extracts field details from a message wrapper."""
        msg = message_wrapper['descriptor']
        filename = message_wrapper['file']
        parent_path = message_wrapper['path']
        
        fields = []
        for i, field in enumerate(msg.field):
            # Field path: parent_path + (2, i) where 2 = field
            field_path = parent_path + (2, i)
            comment = self._get_comment(filename, field_path)
            
            # Extract field_behavior annotations
            field_behaviors = []
            is_sensitive = False
            is_output_only = False
            is_input_only = False
            is_immutable = False
            is_required = False
            is_pii = False
            is_encrypted = False
            is_log_masked = False
            is_log_redacted = False
            
            try:
                # Try to import field_behavior extension (google.api)
                from google.api import field_behavior_pb2
                # field_behavior is a REPEATED extension - check if it has it first
                # Use ListFields() to safely check for repeated extensions
                for fd, value in field.options.ListFields():
                    if fd.name == 'field_behavior':
                        # It's a repeated field, iterate through values
                        for behavior in value:
                            field_behaviors.append(behavior)
                            if behavior == field_behavior_pb2.FieldBehavior.SENSITIVE:
                                is_sensitive = True
                            elif behavior == field_behavior_pb2.FieldBehavior.OUTPUT_ONLY:
                                is_output_only = True
                            elif behavior == field_behavior_pb2.FieldBehavior.INPUT_ONLY:
                                is_input_only = True
                            elif behavior == field_behavior_pb2.FieldBehavior.IMMUTABLE:
                                is_immutable = True
                            elif behavior == field_behavior_pb2.FieldBehavior.REQUIRED:
                                is_required = True
                        break
            except (ImportError, AttributeError, KeyError, TypeError) as e:
                pass
            
            # Extract custom insuretech security annotations
            # These use extension numbers 50010-50017 from insuretech.common.v1.security
            try:
                for fd, value in field.options.ListFields():
                    # Check for custom security annotations by name
                    if fd.name == 'sensitive' and value:
                        is_sensitive = True
                    elif fd.name == 'pii' and value:
                        is_pii = True
                    elif fd.name == 'encrypted_security' and value:
                        is_encrypted = True
                    elif fd.name == 'log_masked' and value:
                        is_log_masked = True
                    elif fd.name == 'log_redacted' and value:
                        is_log_redacted = True
            except (AttributeError, KeyError, TypeError) as e:
                pass
            
            # Fallback: Use proto source parser to extract custom annotations from source files
            requires_consent = False
            data_purpose = None
            retention_days = None
            not_null = False
            unique = False
            primary_key = False
            auto_increment = False
            default_value = None
            sql_type = None
            db_encrypted = False
            
            if self.proto_source_parser:
                msg_full_name = message_wrapper['full_name']
                source_annotations = self.proto_source_parser.get_field_annotations(msg_full_name, field.name)
                
                # Security annotations
                if source_annotations.get('sensitive'):
                    is_sensitive = True
                if source_annotations.get('pii'):
                    is_pii = True
                if source_annotations.get('encrypted'):
                    is_encrypted = True
                if source_annotations.get('log_masked'):
                    is_log_masked = True
                if source_annotations.get('log_redacted'):
                    is_log_redacted = True
                if source_annotations.get('requires_consent'):
                    requires_consent = True
                if source_annotations.get('data_purpose'):
                    data_purpose = source_annotations.get('data_purpose')
                if source_annotations.get('retention_days'):
                    retention_days = source_annotations.get('retention_days')
                
                # Database annotations
                if source_annotations.get('not_null'):
                    not_null = True
                if source_annotations.get('unique'):
                    unique = True
                if source_annotations.get('primary_key'):
                    primary_key = True
                if source_annotations.get('auto_increment'):
                    auto_increment = True
                if source_annotations.get('default_value'):
                    default_value = source_annotations.get('default_value')
                if source_annotations.get('sql_type'):
                    sql_type = source_annotations.get('sql_type')
                if source_annotations.get('db_encrypted'):
                    db_encrypted = True
            
            # Final fallback: check comment for behavior keywords
            if not is_sensitive and 'sensitive' in comment.lower():
                is_sensitive = True
            if not is_output_only and ('output_only' in comment.lower() or 'read-only' in comment.lower()):
                is_output_only = True
            if not is_required and 'required' in comment.lower():
                is_required = True
            
            # Detect map fields
            # In protobuf, map<K,V> is syntactic sugar for:
            # repeated MapFieldEntry { key = 1; value = 2; }
            # where MapFieldEntry has map_entry option set to true
            is_map_field = False
            map_key_type = None
            map_value_type = None
            
            if field.label == 3 and field.type == 11:  # REPEATED MESSAGE
                # Check if the referenced message is a map entry
                type_name = field.type_name
                if type_name.startswith('.'):
                    type_name = type_name[1:]
                
                # Map entries are named like "MessageName.FieldNameEntry"
                if type_name.endswith('Entry') and '.' in type_name:
                    # Look up the message to check if it's a map entry
                    # For now, we'll use a simple heuristic
                    is_map_field = True
                    # Default to string types - will be refined when we look up the actual message
                    map_key_type = 'string'
                    map_value_type = 'string'
            
            fields.append({
                'name': field.name,
                'number': field.number,
                'label': field.label, # 1=OPTIONAL, 2=REQUIRED, 3=REPEATED
                'type': field.type,   # 1=DOUBLE, 9=STRING, 11=MESSAGE, etc.
                'type_name': field.type_name, # For messages/enums
                'comment': comment,
                'oneof_index': field.oneof_index if field.HasField('oneof_index') else None,
                'behaviors': field_behaviors,
                # Map field detection
                'is_map': is_map_field,
                'map_key_type': map_key_type,
                'map_value_type': map_value_type,
                # Security annotations
                'is_sensitive': is_sensitive,
                'is_output_only': is_output_only,
                'is_input_only': is_input_only,
                'is_immutable': is_immutable,
                'is_required': is_required,
                'is_pii': is_pii,
                'is_encrypted': is_encrypted,
                'is_log_masked': is_log_masked,
                'is_log_redacted': is_log_redacted,
                'requires_consent': requires_consent,
                'data_purpose': data_purpose,
                'retention_days': retention_days,
                # Database annotations
                'not_null': not_null,
                'unique': unique,
                'primary_key': primary_key,
                'auto_increment': auto_increment,
                'default_value': default_value,
                'sql_type': sql_type,
                'db_encrypted': db_encrypted
            })
        return fields

