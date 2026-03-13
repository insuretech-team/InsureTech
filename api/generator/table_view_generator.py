#!/usr/bin/env python3
"""
Table View Generator
Generates table view pages for schemas and DTOs grouped by domain
"""

import yaml
import json
import os
from typing import Dict, List, Any
from doc_generator import DocGenerator

class TableViewGenerator:
    """Generate table view pages for data models"""
    
    def __init__(self, openapi_path: str, doc_generator: DocGenerator):
        self.openapi_path = openapi_path
        self.doc_gen = doc_generator
        
        with open(openapi_path, 'r', encoding='utf-8') as f:
            self.full_spec = yaml.safe_load(f)
    
    def generate_schema_table_page(self, domain: str, schemas: List[Dict], output_dir: str, is_dto: bool = False) -> str:
        """Generate a table view page for schemas or DTOs in a domain"""
        
        domain_info = self.doc_gen.get_domain_info(domain)
        item_type = 'DTOs' if is_dto else 'Schemas'
        file_prefix = 'dtos' if is_dto else 'schemas'
        
        # Generate table rows
        rows_html = ''
        for schema in sorted(schemas, key=lambda x: x['name']):
            clean_name = schema['name'].replace('.', '_').replace(':', '_').lower()
            schema_page = f'schema_{clean_name}.html'
            
            badge = ''
            if schema['name'].endswith('Request'):
                badge = '<span class="badge badge-request">Request</span>'
            elif schema['name'].endswith('Response'):
                badge = '<span class="badge badge-response">Response</span>'
            else:
                badge = f'<span class="badge badge-type">{schema.get("type", "object")}</span>'
            
            rows_html += f'''
            <tr onclick="window.location.href='{schema_page}'" style="cursor: pointer;">
                <td class="name-cell">
                    <span class="schema-name">{schema['name']}</span>
                    {badge}
                </td>
                <td class="desc-cell">{schema.get('description', 'No description')}</td>
                <td class="props-cell">{schema.get('properties_count', 0)} properties</td>
            </tr>
            '''
        
        # Generate HTML
        html = f'''<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{domain_info['name']} {item_type} - InsureTech API</title>
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
            max-width: 1400px;
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
        
        .domain-title {{
            display: flex;
            align-items: center;
            gap: 15px;
            margin-bottom: 10px;
        }}
        
        .domain-icon {{
            font-size: 2.5em;
        }}
        
        .domain-name {{
            font-size: 2em;
            font-weight: 600;
        }}
        
        .domain-description {{
            font-size: 1.1em;
            opacity: 0.95;
        }}
        
        .stats {{
            display: flex;
            gap: 30px;
            margin-top: 20px;
        }}
        
        .stat {{
            background: rgba(255,255,255,0.2);
            padding: 10px 20px;
            border-radius: 8px;
        }}
        
        .stat-value {{
            font-size: 1.8em;
            font-weight: bold;
        }}
        
        .stat-label {{
            font-size: 0.9em;
            opacity: 0.9;
        }}
        
        .container {{
            max-width: 1400px;
            margin: 40px auto;
            padding: 0 40px;
        }}
        
        .search-box {{
            margin-bottom: 30px;
            position: relative;
        }}
        
        .search-box input {{
            width: 100%;
            padding: 15px 50px 15px 20px;
            border: 2px solid #e0e0e0;
            border-radius: 10px;
            font-size: 1em;
            background: white;
        }}
        
        .search-box input:focus {{
            outline: none;
            border-color: #667eea;
        }}
        
        .search-icon {{
            position: absolute;
            right: 20px;
            top: 50%;
            transform: translateY(-50%);
            font-size: 1.2em;
            color: #999;
        }}
        
        .table-container {{
            background: white;
            border-radius: 12px;
            box-shadow: 0 2px 8px rgba(0,0,0,0.05);
            overflow: hidden;
        }}
        
        table {{
            width: 100%;
            border-collapse: collapse;
        }}
        
        thead {{
            background: #f8f9fa;
        }}
        
        th {{
            padding: 15px 20px;
            text-align: left;
            font-weight: 600;
            color: #333;
            border-bottom: 2px solid #e0e0e0;
        }}
        
        tbody tr {{
            border-bottom: 1px solid #f0f0f0;
            transition: background 0.2s;
        }}
        
        tbody tr:hover {{
            background: #f8f9fa;
        }}
        
        td {{
            padding: 15px 20px;
        }}
        
        .name-cell {{
            font-family: 'Courier New', monospace;
            font-weight: 600;
            color: #333;
        }}
        
        .schema-name {{
            font-size: 1.05em;
        }}
        
        .desc-cell {{
            color: #666;
        }}
        
        .props-cell {{
            color: #999;
            font-size: 0.9em;
            text-align: center;
        }}
        
        .badge {{
            display: inline-block;
            padding: 3px 8px;
            border-radius: 4px;
            font-size: 0.75em;
            font-weight: 600;
            margin-left: 10px;
            text-transform: uppercase;
        }}
        
        .badge-request {{
            background: #e3f2fd;
            color: #1565c0;
        }}
        
        .badge-response {{
            background: #e8f5e9;
            color: #2e7d32;
        }}
        
        .badge-type {{
            background: #f3e5f5;
            color: #7b1fa2;
        }}
        
        .back-link {{
            display: inline-block;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 12px 24px;
            border-radius: 8px;
            text-decoration: none;
            font-weight: 600;
            margin-bottom: 30px;
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
                <a href="index.html">← Back to Documentation</a> / {item_type}
            </div>
            <div class="domain-title">
                <span class="domain-icon">{domain_info['icon']}</span>
                <div>
                    <div class="domain-name">{domain_info['name']} {item_type}</div>
                    <div class="domain-description">{domain_info['description']}</div>
                </div>
            </div>
            <div class="stats">
                <div class="stat">
                    <div class="stat-value">{len(schemas)}</div>
                    <div class="stat-label">{item_type}</div>
                </div>
            </div>
        </div>
    </div>
    
    <div class="container">
        <a href="index.html" class="back-link">← Back to Documentation Hub</a>
        
        <div class="search-box">
            <input type="text" id="search" placeholder="Search schemas..." onkeyup="filterTable()">
            <span class="search-icon">🔍</span>
        </div>
        
        <div class="table-container">
            <table id="schemaTable">
                <thead>
                    <tr>
                        <th style="width: 35%;">Schema Name</th>
                        <th style="width: 50%;">Description</th>
                        <th style="width: 15%;">Properties</th>
                    </tr>
                </thead>
                <tbody>
                    {rows_html}
                </tbody>
            </table>
        </div>
    </div>
    
    <script>
        function filterTable() {{
            const input = document.getElementById('search');
            const filter = input.value.toLowerCase();
            const table = document.getElementById('schemaTable');
            const rows = table.getElementsByTagName('tr');
            
            for (let i = 1; i < rows.length; i++) {{
                const row = rows[i];
                const text = row.textContent || row.innerText;
                if (text.toLowerCase().indexOf(filter) > -1) {{
                    row.style.display = '';
                }} else {{
                    row.style.display = 'none';
                }}
            }}
        }}
    </script>
</body>
</html>'''
        
        # Write to file
        output_path = os.path.join(output_dir, f'{file_prefix}_{domain}.html')
        with open(output_path, 'w', encoding='utf-8') as f:
            f.write(html)
        
        return f'{file_prefix}_{domain}.html'
    
    def generate_dto_table_page(self, domain: str, dtos: List[Dict], output_dir: str) -> str:
        """Generate a table view page for DTOs in a domain"""
        # Just call schema generator but change the title - DTOs and schemas are same structure
        return self.generate_schema_table_page(domain, dtos, output_dir, is_dto=True)
    
    def generate_all_table_pages(self, output_dir: str) -> Dict[str, Dict[str, str]]:
        """Generate all table view pages"""
        
        print("Generating table view pages for schemas and DTOs...")
        
        organized = self.doc_gen.organize_data()
        
        schema_pages = {}
        dto_pages = {}
        
        # Generate schema table pages
        for domain in sorted(organized['schemas'].keys()):
            schemas = organized['schemas'][domain]
            if schemas:
                page_name = self.generate_schema_table_page(domain, schemas, output_dir)
                schema_pages[domain] = page_name
        
        # Generate DTO table pages
        for domain in sorted(organized['dtos'].keys()):
            dtos = organized['dtos'][domain]
            if dtos:
                page_name = self.generate_dto_table_page(domain, dtos, output_dir)
                dto_pages[domain] = page_name
        
        print(f"  ✓ Generated {len(schema_pages)} schema table pages")
        print(f"  ✓ Generated {len(dto_pages)} DTO table pages")
        
        return {'schemas': schema_pages, 'dtos': dto_pages}


def main():
    """CLI entry point"""
    import argparse
    
    parser = argparse.ArgumentParser(description='Generate table view pages')
    parser.add_argument('--spec', default='../openapi.yaml', help='Path to OpenAPI spec')
    parser.add_argument('--output-dir', default='../docs', help='Output directory')
    
    args = parser.parse_args()
    
    doc_gen = DocGenerator(args.spec)
    table_gen = TableViewGenerator(args.spec, doc_gen)
    
    pages = table_gen.generate_all_table_pages(args.output_dir)
    
    print(f"\n✓ Table pages generated in {args.output_dir}")
    print(f"  Schema tables: {len(pages['schemas'])}")
    print(f"  DTO tables: {len(pages['dtos'])}")


if __name__ == '__main__':
    main()
