#!/usr/bin/env python3
"""
Schema and Enum Page Generator
Generates individual static HTML pages for schemas and enums
"""

import yaml
import json
import os
from typing import Dict, Any, List
from doc_generator import DocGenerator

class SchemaEnumPageGenerator:
    """Generate individual pages for schemas and enums"""
    
    def __init__(self, openapi_path: str, doc_generator: DocGenerator):
        self.openapi_path = openapi_path
        self.doc_gen = doc_generator
        
        with open(openapi_path, 'r', encoding='utf-8') as f:
            self.full_spec = yaml.safe_load(f)
    
    def render_schema_properties_as_table(self, schema: Dict) -> str:
        """Render schema properties as database-style table"""
        if not schema or not isinstance(schema, dict):
            return ''
        
        properties = schema.get('properties', {})
        required = schema.get('required', [])
        
        if not properties:
            return '<p class="no-props">No fields defined</p>'
        
        # Build table rows
        rows_html = ''
        for prop_name, prop_schema in properties.items():
            is_required = prop_name in required
            prop_type = prop_schema.get('type', 'unknown')
            prop_desc = prop_schema.get('description', 'No description')
            prop_format = prop_schema.get('format', '')
            
            # Constraints
            constraints = []
            if 'minLength' in prop_schema:
                constraints.append(f"min: {prop_schema['minLength']}")
            if 'maxLength' in prop_schema:
                constraints.append(f"max: {prop_schema['maxLength']}")
            if 'minimum' in prop_schema:
                constraints.append(f"min: {prop_schema['minimum']}")
            if 'maximum' in prop_schema:
                constraints.append(f"max: {prop_schema['maximum']}")
            if 'pattern' in prop_schema:
                constraints.append(f"pattern: {prop_schema['pattern'][:30]}...")
            if 'enum' in prop_schema:
                constraints.append(f"enum ({len(prop_schema['enum'])} values)")
            
            constraint_text = ', '.join(constraints) if constraints else '—'
            
            # Handle refs and determine type display
            ref_link = ''
            if '$ref' in prop_schema:
                ref_name = prop_schema['$ref'].split('/')[-1]
                clean_ref = ref_name.replace('.', '_').replace(':', '_').lower()
                ref_link = f'<a href="schema_{clean_ref}.html" class="ref-link">→ {ref_name}</a>'
                prop_type = 'reference'
                prop_format = ref_name
            elif prop_type == 'array':
                items = prop_schema.get('items', {})
                if '$ref' in items:
                    ref_name = items['$ref'].split('/')[-1]
                    clean_ref = ref_name.replace('.', '_').replace(':', '_').lower()
                    ref_link = f'<a href="schema_{clean_ref}.html" class="ref-link">→ {ref_name}[]</a>'
                    prop_format = f'array of {ref_name}'
                else:
                    prop_format = f'array of {items.get("type", "unknown")}'
            
            # Nullable
            nullable = prop_schema.get('nullable', False)
            nullable_text = 'YES' if nullable or not is_required else 'NO'
            
            # Required indicator
            required_icon = '🔴' if is_required else '⚪'
            
            rows_html += f'''
            <tr>
                <td class="field-name">
                    {required_icon} <strong>{prop_name}</strong>
                </td>
                <td class="field-type">
                    <code>{prop_type}</code>
                    {f'<br><small>{prop_format}</small>' if prop_format else ''}
                    {ref_link}
                </td>
                <td class="field-nullable">{nullable_text}</td>
                <td class="field-constraints"><small>{constraint_text}</small></td>
                <td class="field-description">{prop_desc}</td>
            </tr>
            '''
        
        html = f'''
        <div class="db-table-container">
            <table class="db-table">
                <thead>
                    <tr>
                        <th style="width: 20%;">Field Name</th>
                        <th style="width: 15%;">Type</th>
                        <th style="width: 8%;">Nullable</th>
                        <th style="width: 20%;">Constraints</th>
                        <th style="width: 37%;">Description</th>
                    </tr>
                </thead>
                <tbody>
                    {rows_html}
                </tbody>
            </table>
        </div>
        '''
        
        return html
    
    def generate_schema_page(self, schema_name: str, schema: Dict, output_dir: str) -> str:
        """Generate a static HTML page for a schema"""
        
        schema_type = schema.get('type', 'object')
        description = schema.get('description', 'No description available')
        
        # Render properties as database table
        properties_html = self.render_schema_properties_as_table(schema)
        
        # Count properties
        prop_count = len(schema.get('properties', {}))
        required_count = len(schema.get('required', []))
        
        # Generate HTML
        html = f'''<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{schema_name} Schema - InsureTech API</title>
    <style>
        * {{ margin: 0; padding: 0; box-sizing: border-box; }}
        
        body {{
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: #f5f7fa;
            line-height: 1.6;
        }}
        
        .header {{
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 30px 40px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }}
        
        .header-content {{
            max-width: 1200px;
            margin: 0 auto;
        }}
        
        .breadcrumb {{
            opacity: 0.9;
            margin-bottom: 15px;
            font-size: 0.9em;
        }}
        
        .breadcrumb a {{
            color: white;
            text-decoration: none;
        }}
        
        .breadcrumb a:hover {{
            text-decoration: underline;
        }}
        
        .schema-title {{
            display: flex;
            align-items: center;
            gap: 15px;
            margin-bottom: 10px;
        }}
        
        .type-badge {{
            background: #4caf50;
            padding: 8px 16px;
            border-radius: 6px;
            font-weight: 600;
            font-size: 0.9em;
        }}
        
        .schema-name {{
            font-size: 2em;
            font-weight: 600;
            font-family: 'Courier New', monospace;
        }}
        
        .container {{
            max-width: 1200px;
            margin: 40px auto;
            padding: 0 40px;
        }}
        
        .card {{
            background: white;
            border-radius: 12px;
            padding: 30px;
            margin-bottom: 30px;
            box-shadow: 0 2px 8px rgba(0,0,0,0.05);
        }}
        
        h2 {{
            color: #333;
            margin-bottom: 20px;
            padding-bottom: 10px;
            border-bottom: 2px solid #e0e0e0;
        }}
        
        .description {{
            color: #555;
            line-height: 1.8;
            font-size: 1.1em;
            margin-bottom: 20px;
        }}
        
        .db-table-container {{
            background: white;
            border-radius: 8px;
            overflow: hidden;
            box-shadow: 0 1px 3px rgba(0,0,0,0.1);
        }}
        
        .db-table {{
            width: 100%;
            border-collapse: collapse;
            font-size: 0.95em;
        }}
        
        .db-table thead {{
            background: #2c3e50;
            color: white;
        }}
        
        .db-table th {{
            padding: 12px 15px;
            text-align: left;
            font-weight: 600;
            font-size: 0.9em;
            text-transform: uppercase;
            letter-spacing: 0.5px;
        }}
        
        .db-table tbody tr {{
            border-bottom: 1px solid #e0e0e0;
            transition: background 0.2s;
        }}
        
        .db-table tbody tr:hover {{
            background: #f8f9fa;
        }}
        
        .db-table tbody tr:last-child {{
            border-bottom: none;
        }}
        
        .db-table td {{
            padding: 12px 15px;
            vertical-align: top;
        }}
        
        .field-name {{
            font-family: 'Courier New', monospace;
            font-weight: 600;
            color: #2c3e50;
        }}
        
        .field-type code {{
            background: #e8f4f8;
            color: #0277bd;
            padding: 2px 6px;
            border-radius: 3px;
            font-size: 0.9em;
            font-weight: 600;
        }}
        
        .field-type small {{
            color: #666;
            display: block;
            margin-top: 4px;
        }}
        
        .field-nullable {{
            text-align: center;
            font-weight: 600;
            color: #666;
        }}
        
        .field-constraints {{
            color: #666;
            font-family: 'Courier New', monospace;
        }}
        
        .field-description {{
            color: #555;
            line-height: 1.5;
        }}
        
        .ref-link {{
            display: inline-block;
            margin-top: 4px;
            color: #667eea;
            text-decoration: none;
            font-weight: 600;
            font-size: 0.9em;
        }}
        
        .ref-link:hover {{
            text-decoration: underline;
        }}
        
        .no-props {{
            color: #999;
            font-style: italic;
            padding: 20px;
            text-align: center;
        }}
        
        .stats-grid {{
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
            gap: 15px;
            margin-bottom: 20px;
        }}
        
        .stat-box {{
            background: #f8f9fa;
            padding: 15px;
            border-radius: 8px;
            text-align: center;
            border-left: 4px solid #667eea;
        }}
        
        .stat-value {{
            font-size: 2em;
            font-weight: bold;
            color: #667eea;
        }}
        
        .stat-label {{
            font-size: 0.85em;
            color: #666;
            margin-top: 5px;
        }}
        
        .back-link {{
            display: inline-block;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 12px 24px;
            border-radius: 8px;
            text-decoration: none;
            font-weight: 600;
            transition: transform 0.3s;
        }}
        
        .back-link:hover {{
            transform: translateY(-2px);
            box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
        }}
    </style>
</head>
<body>
    <div class="header">
        <div class="header-content">
            <div class="breadcrumb">
                <a href="index.html">← Back to Documentation</a> / Schemas
            </div>
            <div class="schema-title">
                <span class="type-badge">{schema_type}</span>
                <span class="schema-name">{schema_name}</span>
            </div>
        </div>
    </div>
    
    <div class="container">
        <div class="card">
            <div class="description">{description}</div>
        </div>
        
        <div class="card">
            <div class="stats-grid">
                <div class="stat-box">
                    <div class="stat-value">{prop_count}</div>
                    <div class="stat-label">Total Fields</div>
                </div>
                <div class="stat-box">
                    <div class="stat-value">{required_count}</div>
                    <div class="stat-label">Required</div>
                </div>
                <div class="stat-box">
                    <div class="stat-value">{prop_count - required_count}</div>
                    <div class="stat-label">Optional</div>
                </div>
            </div>
        </div>
        
        <div class="card">
            <h2>📋 Schema Fields</h2>
            {properties_html}
        </div>
        
        <div style="text-align: center; margin-top: 40px;">
            <a href="index.html" class="back-link">← Back to Documentation Hub</a>
        </div>
    </div>
</body>
</html>'''
        
        # Write to file
        clean_name = schema_name.replace('.', '_').replace(':', '_').lower()
        output_path = os.path.join(output_dir, f'schema_{clean_name}.html')
        with open(output_path, 'w', encoding='utf-8') as f:
            f.write(html)
        
        return f'schema_{clean_name}.html'
    
    def generate_enum_page(self, enum_name: str, enum_schema: Dict, output_dir: str) -> str:
        """Generate a static HTML page for an enum"""
        
        description = enum_schema.get('description', 'No description available')
        enum_values = enum_schema.get('enum', [])
        
        # Generate enum values HTML
        values_html = '<div class="enum-values">'
        for value in enum_values:
            values_html += f'<div class="enum-value"><code>{value}</code></div>'
        values_html += '</div>'
        
        # Generate HTML
        html = f'''<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{enum_name} Enum - InsureTech API</title>
    <style>
        * {{ margin: 0; padding: 0; box-sizing: border-box; }}
        
        body {{
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: #f5f7fa;
            line-height: 1.6;
        }}
        
        .header {{
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 30px 40px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }}
        
        .header-content {{
            max-width: 1200px;
            margin: 0 auto;
        }}
        
        .breadcrumb {{
            opacity: 0.9;
            margin-bottom: 15px;
            font-size: 0.9em;
        }}
        
        .breadcrumb a {{
            color: white;
            text-decoration: none;
        }}
        
        .breadcrumb a:hover {{
            text-decoration: underline;
        }}
        
        .enum-title {{
            display: flex;
            align-items: center;
            gap: 15px;
            margin-bottom: 10px;
        }}
        
        .type-badge {{
            background: #ff9800;
            padding: 8px 16px;
            border-radius: 6px;
            font-weight: 600;
            font-size: 0.9em;
        }}
        
        .enum-name {{
            font-size: 2em;
            font-weight: 600;
            font-family: 'Courier New', monospace;
        }}
        
        .container {{
            max-width: 1200px;
            margin: 40px auto;
            padding: 0 40px;
        }}
        
        .card {{
            background: white;
            border-radius: 12px;
            padding: 30px;
            margin-bottom: 30px;
            box-shadow: 0 2px 8px rgba(0,0,0,0.05);
        }}
        
        h2 {{
            color: #333;
            margin-bottom: 20px;
            padding-bottom: 10px;
            border-bottom: 2px solid #e0e0e0;
        }}
        
        .description {{
            color: #555;
            line-height: 1.8;
            font-size: 1.1em;
            margin-bottom: 20px;
        }}
        
        .enum-values {{
            display: grid;
            grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
            gap: 15px;
        }}
        
        .enum-value {{
            background: #f8f9fa;
            border-left: 4px solid #ff9800;
            padding: 15px;
            border-radius: 6px;
        }}
        
        .enum-value code {{
            font-family: 'Courier New', monospace;
            font-size: 1.1em;
            font-weight: 600;
            color: #333;
        }}
        
        .back-link {{
            display: inline-block;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 12px 24px;
            border-radius: 8px;
            text-decoration: none;
            font-weight: 600;
            transition: transform 0.3s;
        }}
        
        .back-link:hover {{
            transform: translateY(-2px);
            box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
        }}
    </style>
</head>
<body>
    <div class="header">
        <div class="header-content">
            <div class="breadcrumb">
                <a href="index.html">← Back to Documentation</a> / Enums
            </div>
            <div class="enum-title">
                <span class="type-badge">ENUM</span>
                <span class="enum-name">{enum_name}</span>
            </div>
        </div>
    </div>
    
    <div class="container">
        <div class="card">
            <div class="description">{description}</div>
        </div>
        
        <div class="card">
            <h2>Possible Values ({len(enum_values)})</h2>
            {values_html}
        </div>
        
        <div style="text-align: center; margin-top: 40px;">
            <a href="index.html" class="back-link">← Back to Documentation Hub</a>
        </div>
    </div>
</body>
</html>'''
        
        # Write to file
        clean_name = enum_name.replace('.', '_').replace(':', '_').lower()
        output_path = os.path.join(output_dir, f'enum_{clean_name}.html')
        with open(output_path, 'w', encoding='utf-8') as f:
            f.write(html)
        
        return f'enum_{clean_name}.html'
    
    def generate_all_pages(self, output_dir: str) -> Dict[str, Dict[str, str]]:
        """Generate pages for all schemas and enums"""
        
        print("Generating schema and enum pages...")
        
        os.makedirs(output_dir, exist_ok=True)
        
        schemas = self.full_spec.get('components', {}).get('schemas', {})
        schema_pages = {}
        enum_pages = {}
        
        for schema_name, schema_def in schemas.items():
            if isinstance(schema_def, dict):
                # Check if it's an enum
                if schema_def.get('type') == 'string' and 'enum' in schema_def:
                    page_name = self.generate_enum_page(schema_name, schema_def, output_dir)
                    enum_pages[schema_name] = page_name
                else:
                    page_name = self.generate_schema_page(schema_name, schema_def, output_dir)
                    schema_pages[schema_name] = page_name
        
        print(f"  ✓ Generated {len(schema_pages)} schema pages")
        print(f"  ✓ Generated {len(enum_pages)} enum pages")
        
        return {'schemas': schema_pages, 'enums': enum_pages}


def main():
    """CLI entry point"""
    import argparse
    
    parser = argparse.ArgumentParser(description='Generate schema and enum pages')
    parser.add_argument('--spec', default='../openapi.yaml', help='Path to OpenAPI spec')
    parser.add_argument('--output-dir', default='../docs', help='Output directory')
    
    args = parser.parse_args()
    
    doc_gen = DocGenerator(args.spec)
    page_gen = SchemaEnumPageGenerator(args.spec, doc_gen)
    
    pages = page_gen.generate_all_pages(args.output_dir)
    
    print(f"\n✓ Pages generated in {args.output_dir}")
    print(f"  Schemas: {len(pages['schemas'])}")
    print(f"  Enums: {len(pages['enums'])}")


if __name__ == '__main__':
    main()
