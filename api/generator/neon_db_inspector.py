"""
Neon DB Inspector
Connects to Neon PostgreSQL database to get actual table structure
"""

import os
import psycopg2
from typing import Dict, List
from dotenv import load_dotenv


class NeonDBInspector:
    """Inspects Neon database to get actual table structure"""
    
    def __init__(self, env_file: str = "../../.env"):
        """Initialize with database connection from .env"""
        load_dotenv(env_file)
        
        self.host = os.getenv('PGHOST2')
        self.database = os.getenv('PGDATABASE2')
        self.user = os.getenv('PGUSER2')
        self.password = os.getenv('PGPASSWORD2')
        self.port = os.getenv('PGPORT2', '5432')  # Default Neon port
        self.sslmode = os.getenv('PGSSLMODE', 'require')
        
        self.conn = None
        self.cursor = None
        
    def connect(self):
        """Connect to Neon database"""
        try:
            self.conn = psycopg2.connect(
                host=self.host,
                database=self.database,
                user=self.user,
                password=self.password,
                port=self.port,
                sslmode=self.sslmode
            )
            self.cursor = self.conn.cursor()
            print(f"Connected to Neon DB: {self.database}")
            return True
        except Exception as e:
            print(f"Failed to connect to Neon DB: {e}")
            return False
    
    def disconnect(self):
        """Disconnect from database"""
        if self.cursor:
            self.cursor.close()
        if self.conn:
            self.conn.close()
        print("Disconnected from Neon DB")
    
    def get_all_schemas(self) -> List[str]:
        """Get all database schemas (excluding system schemas)"""
        query = """
        SELECT schema_name 
        FROM information_schema.schemata 
        WHERE schema_name NOT IN ('pg_catalog', 'information_schema', 'pg_toast')
        ORDER BY schema_name;
        """
        
        self.cursor.execute(query)
        schemas = [row[0] for row in self.cursor.fetchall()]
        return schemas
    
    def get_tables_in_schema(self, schema_name: str) -> List[Dict]:
        """Get all tables in a specific schema"""
        query = """
        SELECT 
            table_name,
            (SELECT COUNT(*) FROM information_schema.columns 
             WHERE table_schema = t.table_schema AND table_name = t.table_name) as column_count
        FROM information_schema.tables t
        WHERE table_schema = %s
        AND table_type = 'BASE TABLE'
        ORDER BY table_name;
        """
        
        self.cursor.execute(query, (schema_name,))
        tables = []
        for row in self.cursor.fetchall():
            tables.append({
                'table_name': row[0],
                'column_count': row[1],
                'schema_name': schema_name
            })
        return tables
    
    def get_table_columns(self, schema_name: str, table_name: str) -> List[Dict]:
        """Get columns for a specific table"""
        query = """
        SELECT 
            column_name,
            data_type,
            is_nullable,
            column_default
        FROM information_schema.columns
        WHERE table_schema = %s AND table_name = %s
        ORDER BY ordinal_position;
        """
        
        self.cursor.execute(query, (schema_name, table_name))
        columns = []
        for row in self.cursor.fetchall():
            columns.append({
                'column_name': row[0],
                'data_type': row[1],
                'is_nullable': row[2] == 'YES',
                'column_default': row[3]
            })
        return columns
    
    def get_complete_db_structure(self) -> Dict:
        """Get complete database structure"""
        if not self.connect():
            return {'error': 'Failed to connect to database'}
        
        try:
            structure = {
                'schemas': {},
                'stats': {
                    'total_schemas': 0,
                    'total_tables': 0,
                    'total_columns': 0
                }
            }
            
            schemas = self.get_all_schemas()
            structure['stats']['total_schemas'] = len(schemas)
            
            for schema_name in schemas:
                tables = self.get_tables_in_schema(schema_name)
                structure['schemas'][schema_name] = {
                    'table_count': len(tables),
                    'tables': tables
                }
                structure['stats']['total_tables'] += len(tables)
                structure['stats']['total_columns'] += sum(t['column_count'] for t in tables)
            
            return structure
            
        except Exception as e:
            print(f"Error getting DB structure: {e}")
            return {'error': str(e)}
        finally:
            self.disconnect()
    
    def compare_with_proto(self, proto_summary: Dict) -> Dict:
        """Compare database structure with proto definitions"""
        db_structure = self.get_complete_db_structure()
        
        if 'error' in db_structure:
            return {'error': db_structure['error']}
        
        comparison = {
            'proto_only': [],  # Tables in proto but not in DB
            'db_only': [],     # Tables in DB but not in proto
            'matched': [],     # Tables in both
            'schema_mismatches': []  # Tables in different schemas
        }
        
        # Get proto tables
        proto_tables = {}
        for schema_name, tables in proto_summary.get('schema_groups', {}).items():
            for table in tables:
                proto_tables[table['table_name']] = {
                    'schema': schema_name,
                    'message': table['message_name']
                }
        
        # Get DB tables
        db_tables = {}
        for schema_name, schema_data in db_structure.get('schemas', {}).items():
            for table in schema_data['tables']:
                db_tables[table['table_name']] = {
                    'schema': schema_name,
                    'columns': table['column_count']
                }
        
        # Compare
        proto_table_names = set(proto_tables.keys())
        db_table_names = set(db_tables.keys())
        
        comparison['proto_only'] = list(proto_table_names - db_table_names)
        comparison['db_only'] = list(db_table_names - proto_table_names)
        
        # Check matched tables for schema mismatches
        for table_name in proto_table_names & db_table_names:
            proto_schema = proto_tables[table_name]['schema']
            db_schema = db_tables[table_name]['schema']
            
            if proto_schema == db_schema:
                comparison['matched'].append({
                    'table_name': table_name,
                    'schema': proto_schema
                })
            else:
                comparison['schema_mismatches'].append({
                    'table_name': table_name,
                    'proto_schema': proto_schema,
                    'db_schema': db_schema
                })
        
        return comparison


def main():
    """CLI entry point"""
    import argparse
    import json
    
    parser = argparse.ArgumentParser(description='Inspect Neon database structure')
    parser.add_argument('--env-file', default='../../.env', help='Path to .env file')
    parser.add_argument('--output', default='../neon_db_structure.json', help='Output JSON file')
    parser.add_argument('--compare-proto', help='Path to proto_schema_summary.json for comparison')
    
    args = parser.parse_args()
    
    inspector = NeonDBInspector(args.env_file)
    
    if args.compare_proto:
        # Compare with proto
        with open(args.compare_proto, 'r') as f:
            proto_summary = json.load(f)
        
        print(f"\nComparing database with proto definitions...")
        comparison = inspector.compare_with_proto(proto_summary)
        
        if 'error' in comparison:
            print(f"Error: {comparison['error']}")
            return
        
        print(f"\nMatched tables: {len(comparison['matched'])}")
        print(f"Proto only: {len(comparison['proto_only'])}")
        print(f"DB only: {len(comparison['db_only'])}")
        print(f"Schema mismatches: {len(comparison['schema_mismatches'])}")
        
        if comparison['proto_only']:
            print("\nTables in proto but not in DB:")
            for table in comparison['proto_only'][:10]:
                print(f"  - {table}")
        
        if comparison['db_only']:
            print("\nTables in DB but not in proto:")
            for table in comparison['db_only'][:10]:
                print(f"  - {table}")
        
        if comparison['schema_mismatches']:
            print("\nSchema mismatches:")
            for mismatch in comparison['schema_mismatches']:
                print(f"  - {mismatch['table_name']}: proto={mismatch['proto_schema']}, db={mismatch['db_schema']}")
        
        # Export comparison
        output_file = args.output.replace('.json', '_comparison.json')
        with open(output_file, 'w') as f:
            json.dump(comparison, f, indent=2)
        print(f"\nComparison exported to: {output_file}")
    else:
        # Just get DB structure
        print("\nInspecting database structure...")
        structure = inspector.get_complete_db_structure()
        
        if 'error' in structure:
            print(f"Error: {structure['error']}")
            return
        
        print(f"\nDatabase Structure:")
        print(f"   Schemas: {structure['stats']['total_schemas']}")
        print(f"   Tables: {structure['stats']['total_tables']}")
        print(f"   Columns: {structure['stats']['total_columns']}")
        
        print(f"\nSchemas:")
        for schema_name, schema_data in sorted(structure['schemas'].items()):
            print(f"   {schema_name}: {schema_data['table_count']} tables")
        
        # Export structure
        with open(args.output, 'w') as f:
            json.dump(structure, f, indent=2)
        print(f"\nStructure exported to: {args.output}")


if __name__ == '__main__':
    main()
