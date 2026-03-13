#!/usr/bin/env python3
"""
Domain Page Generator
Generates individual static HTML pages for each domain with Swagger UI embedded
"""

import yaml
import json
import os
from pathlib import Path
from typing import Dict, List, Any
from doc_generator import DocGenerator

class DomainPageGenerator:
    """Generate individual domain pages with embedded Swagger UI"""
    
    def __init__(self, openapi_path: str, doc_generator: DocGenerator):
        self.openapi_path = openapi_path
        self.doc_gen = doc_generator
        
        with open(openapi_path, 'r', encoding='utf-8') as f:
            self.full_spec = yaml.safe_load(f)
    
    def filter_spec_by_domain(self, domain: str, api_paths: List[str]) -> Dict[str, Any]:
        """Create a filtered OpenAPI spec for a specific domain"""
        
        # Start with base spec structure
        filtered_spec = {
            'openapi': self.full_spec['openapi'],
            'info': {
                'title': f"{self.doc_gen.get_domain_info(domain)['name']} API",
                'version': self.full_spec['info'].get('version', '1.0.0'),
                'description': self.doc_gen.get_domain_info(domain)['description']
            },
            'servers': self.full_spec.get('servers', []),
            'paths': {},
            'components': {
                'schemas': {},
                'securitySchemes': self.full_spec.get('components', {}).get('securitySchemes', {})
            }
        }
        
        # Add only paths for this domain
        all_paths = self.full_spec.get('paths', {})
        for api_path in api_paths:
            if api_path in all_paths:
                filtered_spec['paths'][api_path] = all_paths[api_path]
        
        # Collect all referenced schemas recursively
        all_schemas = self.full_spec.get('components', {}).get('schemas', {})
        referenced_schemas = set()
        self._collect_all_schema_refs(filtered_spec['paths'], referenced_schemas, all_schemas)
        
        # Add all collected schemas to filtered spec
        for schema_name in referenced_schemas:
            if schema_name in all_schemas:
                filtered_spec['components']['schemas'][schema_name] = all_schemas[schema_name]
        
        return filtered_spec
    
    def _collect_all_schema_refs(self, obj: Any, refs: set, all_schemas: Dict, visited: set = None):
        """Recursively collect all $ref schema names including nested ones"""
        if visited is None:
            visited = set()
        
        # Skip if not a dict or list
        if not isinstance(obj, (dict, list, str, int, float, bool, type(None))):
            return
        
        if isinstance(obj, dict):
            if '$ref' in obj:
                ref = obj.get('$ref', '')
                if isinstance(ref, str) and ref.startswith('#/components/schemas/'):
                    schema_name = ref.split('/')[-1]
                    if schema_name and schema_name not in visited:
                        visited.add(schema_name)
                        refs.add(schema_name)
                        # Recursively collect refs from this schema
                        if schema_name in all_schemas:
                            self._collect_all_schema_refs(all_schemas[schema_name], refs, all_schemas, visited)
            # Only iterate over dict values
            for key, value in obj.items():
                if isinstance(value, (dict, list)):
                    self._collect_all_schema_refs(value, refs, all_schemas, visited)
        elif isinstance(obj, list):
            for item in obj:
                if isinstance(item, (dict, list)):
                    self._collect_all_schema_refs(item, refs, all_schemas, visited)
    
    def generate_domain_page(self, domain: str, apis: List[Dict], output_dir: str):
        """Generate an individual domain page"""
        
        domain_info = self.doc_gen.get_domain_info(domain)
        
        # Extract paths from API list
        api_paths = [api['path'] for api in apis]
        
        # Create filtered spec
        filtered_spec = self.filter_spec_by_domain(domain, api_paths)
        
        # Calculate stats
        num_endpoints = len(apis)
        num_schemas = len(filtered_spec.get('components', {}).get('schemas', {}))
        
        # Convert to JSON for embedding
        spec_json = json.dumps(filtered_spec, indent=2)
        
        # Generate HTML (using regular string to avoid f-string issues with CSS)
        html = '''<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>''' + domain_info['name'] + ''' API - InsureTech</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5.10.3/swagger-ui.css">
    <style>
        body {{
            margin: 0;
            padding: 0;
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
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
            display: flex;
            align-items: center;
            gap: 20px;
        }}
        
        .header-icon {{
            font-size: 3em;
        }}
        
        .header-text h1 {{
            margin: 0;
            font-size: 2em;
        }}
        
        .header-text p {{
            margin: 5px 0 0 0;
            opacity: 0.9;
            font-size: 1.1em;
        }}
        
        .header-stats {{
            margin-left: auto;
            display: flex;
            gap: 30px;
        }}
        
        .stat {{
            text-align: center;
        }}
        
        .stat-value {{
            font-size: 2em;
            font-weight: bold;
        }}
        
        .stat-label {{
            font-size: 0.9em;
            opacity: 0.8;
        }}
        
        .nav-bar {{
            background: white;
            border-bottom: 2px solid #e0e0e0;
            padding: 15px 40px;
            box-shadow: 0 2px 5px rgba(0,0,0,0.05);
        }}
        
        .nav-content {{
            max-width: 1400px;
            margin: 0 auto;
            display: flex;
            gap: 20px;
            align-items: center;
        }}
        
        .nav-link {{
            text-decoration: none;
            color: #667eea;
            font-weight: 500;
            padding: 8px 16px;
            border-radius: 6px;
            transition: all 0.3s;
        }}
        
        .nav-link:hover {{
            background: rgba(102, 126, 234, 0.1);
        }}
        
        .swagger-container {{
            max-width: 1400px;
            margin: 0 auto;
            padding: 20px 40px;
        }}
        
        .swagger-ui .topbar {{
            display: none;
        }}
        
        .swagger-ui .info {{
            margin: 20px 0;
        }}
        
        @media (max-width: 768px) {{
            .header-stats {{
                display: none;
            }}
        }}
    </style>
</head>
<body>
    <div class="header">
        <div class="header-content">
            <div class="header-icon">''' + domain_info['icon'] + '''</div>
            <div class="header-text">
                <h1>''' + domain_info['name'] + '''</h1>
                <p>''' + domain_info['description'] + '''</p>
            </div>
            <div class="header-stats">
                <div class="stat">
                    <div class="stat-value">''' + str(num_endpoints) + '''</div>
                    <div class="stat-label">Endpoints</div>
                </div>
                <div class="stat">
                    <div class="stat-value">''' + str(num_schemas) + '''</div>
                    <div class="stat-label">Schemas</div>
                </div>
            </div>
        </div>
    </div>
    
    <div class="nav-bar">
        <div class="nav-content">
            <a href="index.html" class="nav-link">← Back to All Domains</a>
            <a href="../openapi.yaml" class="nav-link" download>📄 Download Full Spec</a>
            <a href="swagger.html" class="nav-link">📘 Full Swagger UI</a>
            <a href="redoc.html" class="nav-link">📗 Full ReDoc</a>
        </div>
    </div>
    
    <div class="swagger-container">
        <div id="swagger-ui"></div>
    </div>
    
    <script src="https://unpkg.com/swagger-ui-dist@5.10.3/swagger-ui-bundle.js"></script>
    <script src="https://unpkg.com/swagger-ui-dist@5.10.3/swagger-ui-standalone-preset.js"></script>
    <script>
        window.onload = function() {{
            // Embedded spec for this domain
            const spec = {spec_json};
            
            const ui = SwaggerUIBundle({
                spec: spec,
                dom_id: '#swagger-ui',
                deepLinking: true,
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIStandalonePreset
                ],
                plugins: [
                    SwaggerUIBundle.plugins.DownloadUrl
                ],
                layout: "BaseLayout",
                defaultModelsExpandDepth: 2,
                defaultModelExpandDepth: 2,
                docExpansion: "list",
                filter: true,
                showExtensions: true,
                showCommonExtensions: true,
                tryItOutEnabled: true
            });
            
            window.ui = ui;
        };
    </script>
</body>
</html>'''
        
        # Write to file
        output_path = os.path.join(output_dir, f'domain_{domain}.html')
        with open(output_path, 'w', encoding='utf-8') as f:
            f.write(html)
        
        return output_path
    
    def generate_all_domain_pages(self, output_dir: str) -> Dict[str, str]:
        """Generate pages for all domains"""
        
        print("Generating individual domain pages...")
        
        # Organize data
        organized = self.doc_gen.organize_data()
        
        # Ensure output directory exists
        os.makedirs(output_dir, exist_ok=True)
        
        # Generate page for each domain
        generated_pages = {}
        for domain in sorted(organized['apis'].keys()):
            apis = organized['apis'][domain]
            if not apis:
                continue
            
            page_path = self.generate_domain_page(domain, apis, output_dir)
            generated_pages[domain] = os.path.basename(page_path)
            print(f"  ✓ Generated {domain}: {len(apis)} endpoints")
        
        print(f"\nGenerated {len(generated_pages)} domain pages")
        return generated_pages


def main():
    """CLI entry point"""
    import argparse
    
    parser = argparse.ArgumentParser(description='Generate domain-specific API pages')
    parser.add_argument('--spec', default='../openapi.yaml', help='Path to OpenAPI spec')
    parser.add_argument('--output-dir', default='../docs', help='Output directory')
    
    args = parser.parse_args()
    
    # Initialize generators
    doc_gen = DocGenerator(args.spec)
    page_gen = DomainPageGenerator(args.spec, doc_gen)
    
    # Generate all pages
    generated_pages = page_gen.generate_all_domain_pages(args.output_dir)
    
    print(f"\n✓ Domain pages generated in {args.output_dir}")
    print(f"  Total: {len(generated_pages)} pages")


if __name__ == '__main__':
    main()
