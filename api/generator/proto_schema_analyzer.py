"""
Proto Schema Analyzer
Recursively scans proto files to extract database schema groups, tables, enums, and message types.
"""

import os
import re
from collections import defaultdict
from typing import Dict, List, Set, Tuple
from pathlib import Path


class ProtoSchemaAnalyzer:
    """Analyzes proto files to extract schema structure"""
    
    def __init__(self, proto_root: str = "../../proto"):
        self.proto_root = proto_root
        self.schema_groups = defaultdict(list)  # schema_name -> [table_names]
        self.enums = []  # List of enum definitions
        self.entities = []  # Messages with is_table: true
        self.dtos = []  # Request/Response messages
        self.events = []  # Event messages
        self.other_messages = []  # Other message types
        
    def scan_all_protos(self) -> Dict:
        """Recursively scan all proto files and extract schema information"""
        print(f"Scanning proto files in: {self.proto_root}")
        
        proto_files = list(Path(self.proto_root).rglob("*.proto"))
        print(f"Found {len(proto_files)} proto files")
        
        for proto_file in proto_files:
            self._analyze_proto_file(proto_file)
        
        return self._build_summary()
    
    def _analyze_proto_file(self, proto_file: Path):
        """Analyze a single proto file"""
        try:
            with open(proto_file, 'r', encoding='utf-8') as f:
                content = f.read()
            
            # Extract package name
            package_match = re.search(r'package\s+([\w.]+);', content)
            package = package_match.group(1) if package_match else 'unknown'
            
            # Find all message definitions
            self._extract_messages(content, package, str(proto_file))
            
            # Find all enum definitions
            self._extract_enums(content, package, str(proto_file))
            
        except Exception as e:
            print(f"Error analyzing {proto_file}: {e}")
    
    def _extract_messages(self, content: str, package: str, file_path: str):
        """Extract message definitions from proto content"""
        
        # Split content into lines for better parsing
        lines = content.split('\n')
        i = 0
        
        while i < len(lines):
            line = lines[i].strip()
            
            # Look for message definition
            if line.startswith('message ') and '{' in line:
                message_name = line.split()[1].replace('{', '').strip()
                
                # Find the closing brace for this message
                brace_count = 1
                message_start = i
                i += 1
                message_lines = []
                
                while i < len(lines) and brace_count > 0:
                    current_line = lines[i]
                    message_lines.append(current_line)
                    brace_count += current_line.count('{') - current_line.count('}')
                    i += 1
                
                message_body = '\n'.join(message_lines)
                
                # Check if it's a table (entity)
                if 'option (insuretech.common.v1.table)' in message_body or 'option(insuretech.common.v1.table)' in message_body:
                    # Extract table options (multi-line)
                    table_name = None
                    schema_name = 'public'
                    migration_order = 999
                    is_table = False
                    
                    for msg_line in message_lines:
                        if 'table_name' in msg_line and ':' in msg_line:
                            match = re.search(r'table_name\s*:\s*"([^"]+)"', msg_line)
                            if match:
                                table_name = match.group(1)
                        elif 'schema_name' in msg_line and ':' in msg_line:
                            match = re.search(r'schema_name\s*:\s*"([^"]+)"', msg_line)
                            if match:
                                schema_name = match.group(1)
                        elif 'migration_order' in msg_line and ':' in msg_line:
                            match = re.search(r'migration_order\s*:\s*(\d+)', msg_line)
                            if match:
                                migration_order = int(match.group(1))
                        elif 'is_table' in msg_line and ':' in msg_line:
                            match = re.search(r'is_table\s*:\s*(true|false)', msg_line)
                            if match:
                                is_table = match.group(1) == 'true'
                    
                    if is_table and table_name:
                        # Add to schema groups
                        self.schema_groups[schema_name].append({
                            'table_name': table_name,
                            'message_name': message_name,
                            'package': package,
                            'migration_order': migration_order,
                            'file_path': file_path
                        })
                        
                        # Add to entities list
                        self.entities.append({
                            'name': message_name,
                            'table_name': table_name,
                            'schema_name': schema_name,
                            'package': package,
                            'migration_order': migration_order,
                            'file_path': file_path
                        })
                else:
                    # Categorize non-table messages
                    if message_name.endswith('Request'):
                        self.dtos.append({
                            'name': message_name,
                            'type': 'request',
                            'package': package,
                            'file_path': file_path
                        })
                    elif message_name.endswith('Response'):
                        self.dtos.append({
                            'name': message_name,
                            'type': 'response',
                            'package': package,
                            'file_path': file_path
                        })
                    elif message_name.endswith('Event'):
                        self.events.append({
                            'name': message_name,
                            'package': package,
                            'file_path': file_path
                        })
                    else:
                        # Other message types (embedded types, etc.)
                        self.other_messages.append({
                            'name': message_name,
                            'package': package,
                            'file_path': file_path
                        })
            else:
                i += 1
    
    def _extract_enums(self, content: str, package: str, file_path: str):
        """Extract enum definitions from proto content"""
        
        # Pattern to match enum blocks
        enum_pattern = r'enum\s+(\w+)\s*\{([^}]+)\}'
        
        for match in re.finditer(enum_pattern, content, re.MULTILINE):
            enum_name = match.group(1)
            enum_body = match.group(2)
            
            # Extract enum values
            value_pattern = r'(\w+)\s*=\s*(\d+);'
            values = []
            for value_match in re.finditer(value_pattern, enum_body):
                values.append({
                    'name': value_match.group(1),
                    'number': int(value_match.group(2))
                })
            
            self.enums.append({
                'name': enum_name,
                'package': package,
                'values': values,
                'value_count': len(values),
                'file_path': file_path
            })
    
    def _build_summary(self) -> Dict:
        """Build a summary of all extracted information"""
        
        # Sort tables within each schema by migration_order
        for schema_name in self.schema_groups:
            self.schema_groups[schema_name].sort(key=lambda x: x['migration_order'])
        
        summary = {
            'schema_groups': dict(self.schema_groups),
            'entities': sorted(self.entities, key=lambda x: x['migration_order']),
            'enums': sorted(self.enums, key=lambda x: x['name']),
            'dtos': sorted(self.dtos, key=lambda x: x['name']),
            'events': sorted(self.events, key=lambda x: x['name']),
            'other_messages': sorted(self.other_messages, key=lambda x: x['name']),
            'stats': {
                'schema_groups_count': len(self.schema_groups),
                'total_tables': sum(len(tables) for tables in self.schema_groups.values()),
                'enums_count': len(self.enums),
                'dtos_count': len(self.dtos),
                'events_count': len(self.events),
                'entities_count': len(self.entities),
                'other_messages_count': len(self.other_messages)
            }
        }
        
        return summary
    
    def print_summary(self, summary: Dict):
        """Print a formatted summary"""
        print("\n" + "="*80)
        print("PROTO SCHEMA ANALYSIS SUMMARY")
        print("="*80)
        
        stats = summary['stats']
        print(f"\nStatistics:")
        print(f"   Database Schema Groups: {stats['schema_groups_count']}")
        print(f"   Total Tables (Entities): {stats['total_tables']}")
        print(f"   Enums: {stats['enums_count']}")
        print(f"   DTOs (Request/Response): {stats['dtos_count']}")
        print(f"   Events: {stats['events_count']}")
        print(f"   Other Messages: {stats['other_messages_count']}")
        
        print(f"\nDatabase Schema Groups:")
        for schema_name in sorted(summary['schema_groups'].keys()):
            tables = summary['schema_groups'][schema_name]
            print(f"\n   {schema_name} ({len(tables)} tables)")
            for table in tables:
                print(f"      - {table['table_name']} (migration_order: {table['migration_order']})")
        
        print("\n" + "="*80)
    
    def export_to_json(self, output_file: str, summary: Dict):
        """Export summary to JSON file"""
        import json
        with open(output_file, 'w', encoding='utf-8') as f:
            json.dump(summary, f, indent=2, ensure_ascii=False)
        print(f"\n✅ Exported to: {output_file}")


def main():
    """CLI entry point"""
    import argparse
    
    parser = argparse.ArgumentParser(description='Analyze proto files for schema structure')
    parser.add_argument('--proto-root', default='../../proto', help='Root directory for proto files')
    parser.add_argument('--output', default='../proto_schema_summary.json', help='Output JSON file')
    
    args = parser.parse_args()
    
    analyzer = ProtoSchemaAnalyzer(args.proto_root)
    summary = analyzer.scan_all_protos()
    analyzer.print_summary(summary)
    analyzer.export_to_json(args.output, summary)


if __name__ == '__main__':
    main()
