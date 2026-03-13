"""
Proto Source Parser - Extract custom security annotations from .proto source files

Since custom extensions (insuretech.common.v1.sensitive, pii, etc.) are not accessible
via the descriptor unless they're compiled into Python, we parse the source .proto files
directly to extract these annotations.
"""

import os
import re
from typing import Dict, Set, Tuple


class ProtoSourceParser:
    """Parse .proto source files to extract security annotations"""
    
    def __init__(self, proto_root: str):
        self.proto_root = proto_root
        self.field_annotations = {}  # message_full_name.field_name -> {sensitive, pii, etc.}
    
    def scan_all_protos(self):
        """Scan all .proto files in proto_root"""
        print("Scanning proto source files for security annotations...")
        count = 0
        for root, dirs, files in os.walk(self.proto_root):
            for file in files:
                if file.endswith('.proto'):
                    filepath = os.path.join(root, file)
                    self._parse_proto_file(filepath)
                    count += 1
        print(f"  Scanned {count} proto files")
        print(f"  Found {len(self.field_annotations)} fields with security annotations")
    
    def _parse_proto_file(self, filepath: str):
        """Parse a single .proto file"""
        try:
            with open(filepath, 'r', encoding='utf-8') as f:
                content = f.read()
            
            # Extract package name
            package_match = re.search(r'package\s+([\w.]+)\s*;', content)
            if not package_match:
                return
            package = package_match.group(1)
            
            # Find all messages
            # Use simpler approach: find message blocks, then parse fields line by line
            message_pattern = r'message\s+(\w+)\s*\{'
            
            for msg_match in re.finditer(message_pattern, content):
                message_name = msg_match.group(1)
                full_message_name = f"{package}.{message_name}"
                
                # Find the message body by counting braces
                start_pos = msg_match.end()
                brace_count = 1
                end_pos = start_pos
                
                while brace_count > 0 and end_pos < len(content):
                    if content[end_pos] == '{':
                        brace_count += 1
                    elif content[end_pos] == '}':
                        brace_count -= 1
                    end_pos += 1
                
                message_body = content[start_pos:end_pos-1]
                
                # Parse fields - look for pattern: type name = number [...];
                # Match from field type to semicolon, handling nested braces
                lines = message_body.split('\n')
                i = 0
                while i < len(lines):
                    line = lines[i].strip()
                    
                    # Check if line starts a field definition
                    field_start_match = re.match(r'([\w.]+)\s+(\w+)\s*=\s*\d+\s*\[', line)
                    if field_start_match:
                        field_name = field_start_match.group(2)
                        
                        # Collect the full field definition until we find ];
                        field_def = line
                        while not field_def.rstrip().endswith('];') and i < len(lines) - 1:
                            i += 1
                            field_def += '\n' + lines[i]
                        
                        # Extract security annotations
                        annotations = self._extract_annotations(field_def)
                        
                        if annotations:
                            key = f"{full_message_name}.{field_name}"
                            self.field_annotations[key] = annotations
                    
                    i += 1
        
        except Exception as e:
            print(f"  Warning: Failed to parse {filepath}: {e}")
    
    def _extract_annotations(self, annotations_str: str) -> Dict[str, any]:
        """Extract security and database annotations from field options string"""
        annotations = {}
        
        # ===== SECURITY ANNOTATIONS (insuretech.common.v1.security) =====
        if 'insuretech.common.v1.sensitive' in annotations_str and '= true' in annotations_str:
            annotations['sensitive'] = True
        if 'insuretech.common.v1.pii' in annotations_str and '= true' in annotations_str:
            annotations['pii'] = True
        if 'insuretech.common.v1.encrypted_security' in annotations_str and '= true' in annotations_str:
            annotations['encrypted'] = True
        if 'insuretech.common.v1.log_masked' in annotations_str and '= true' in annotations_str:
            annotations['log_masked'] = True
        if 'insuretech.common.v1.log_redacted' in annotations_str and '= true' in annotations_str:
            annotations['log_redacted'] = True
        if 'insuretech.common.v1.requires_consent' in annotations_str and '= true' in annotations_str:
            annotations['requires_consent'] = True
        
        # String annotations - extract the value
        data_purpose_match = re.search(r'insuretech\.common\.v1\.data_purpose\)\s*=\s*"([^"]+)"', annotations_str)
        if data_purpose_match:
            annotations['data_purpose'] = data_purpose_match.group(1)
        
        # Integer annotations - extract the value
        retention_days_match = re.search(r'insuretech\.common\.v1\.retention_days\)\s*=\s*(\d+)', annotations_str)
        if retention_days_match:
            annotations['retention_days'] = int(retention_days_match.group(1))
        
        # ===== DATABASE ANNOTATIONS (insuretech.common.v1.db) =====
        # These are relevant for OpenAPI generation
        if 'insuretech.common.v1.column' in annotations_str:
            # Extract column options that map to OpenAPI
            if 'not_null: true' in annotations_str or 'not_null = true' in annotations_str:
                annotations['not_null'] = True
            
            if 'unique: true' in annotations_str or 'unique = true' in annotations_str:
                annotations['unique'] = True
            
            if 'primary_key: true' in annotations_str or 'primary_key = true' in annotations_str:
                annotations['primary_key'] = True
            
            if 'auto_increment: true' in annotations_str or 'auto_increment = true' in annotations_str:
                annotations['auto_increment'] = True
            
            # Extract default value
            default_match = re.search(r'default_value:\s*"([^"]+)"', annotations_str)
            if default_match:
                # Strip any extra quotes from the default value
                # Proto might have: default_value: "'ACTIVE'" 
                # We want: ACTIVE (not 'ACTIVE')
                default_val = default_match.group(1).strip("'\"")
                annotations['default_value'] = default_val
            
            # Extract SQL type for format hints
            sql_type_match = re.search(r'sql_type:\s*"([^"]+)"', annotations_str)
            if sql_type_match:
                annotations['sql_type'] = sql_type_match.group(1)
            
            # Database-level encryption (different from security encryption)
            if 'encrypted: true' in annotations_str:
                annotations['db_encrypted'] = True
        
        return annotations
    
    def get_field_annotations(self, message_full_name: str, field_name: str) -> Dict[str, any]:
        """Get security annotations for a specific field"""
        key = f"{message_full_name}.{field_name}"
        return self.field_annotations.get(key, {})
    
    def has_annotation(self, message_full_name: str, field_name: str, annotation: str) -> bool:
        """Check if a field has a specific annotation"""
        annotations = self.get_field_annotations(message_full_name, field_name)
        return annotations.get(annotation, False)
