#!/usr/bin/env python3
"""
Endpoint Page Generator
Generates individual static HTML pages for EACH API endpoint
"""

import yaml
import json
import os
import re
from pathlib import Path
from typing import Dict, Any
from doc_generator import DocGenerator

class EndpointPageGenerator:
    """Generate individual static pages for each endpoint"""
    
    def __init__(self, openapi_path: str, doc_generator: DocGenerator):
        self.openapi_path = openapi_path
        self.doc_gen = doc_generator
        
        with open(openapi_path, 'r', encoding='utf-8') as f:
            self.full_spec = yaml.safe_load(f)
    
    def get_endpoint_id(self, method: str, path: str) -> str:
        """Generate a unique ID for an endpoint"""
        # Convert /v1/authn/login to authn_login_post
        # Also replace : (colon) for custom actions like /v1/auth/otp:send
        clean_path = (path.replace('/v1/', '')
                          .replace('/', '_')
                          .replace('{', '')
                          .replace('}', '')
                          .replace('-', '_')
                          .replace(':', '_'))
        return f"{clean_path}_{method.lower()}"
    
    def get_schema_details(self, schema_ref: str) -> Dict[str, Any]:
        """Get full schema details including properties"""
        if not schema_ref or not schema_ref.startswith('#/components/schemas/'):
            return {}
        
        schema_name = schema_ref.split('/')[-1]
        schemas = self.full_spec.get('components', {}).get('schemas', {})
        return schemas.get(schema_name, {})
    
    def render_schema(self, schema: Dict, level: int = 0) -> str:
        """Render a schema as formatted HTML with full details"""
        if not schema:
            return '<div class="schema-empty">No schema defined</div>'
        
        if '$ref' in schema:
            schema_details = self.get_schema_details(schema['$ref'])
            schema_name = schema['$ref'].split('/')[-1]
            if schema_details:
                return f'<div class="schema-ref-name">{schema_name}</div>' + self.render_schema(schema_details, level)
            return f'<div class="schema-ref-name">{schema_name}</div>'
        
        indent = '  ' * level
        html = '<div class="schema-properties">'
        
        schema_type = schema.get('type', 'object')
        
        if schema_type == 'object' and 'properties' in schema:
            html += '<div class="json-block"><pre>{'
            required = schema.get('required', [])
            
            for prop_name, prop_schema in schema['properties'].items():
                is_required = prop_name in required
                prop_type = prop_schema.get('type', 'unknown')
                prop_desc = prop_schema.get('description', '')
                
                html += f'\n{indent}  "<span class="json-key">{prop_name}</span>": '
                
                if '$ref' in prop_schema:
                    ref_name = prop_schema['$ref'].split('/')[-1]
                    html += f'<span class="json-ref">{ref_name}</span>'
                elif prop_type == 'object':
                    html += self.render_schema(prop_schema, level + 1)
                elif prop_type == 'array':
                    items = prop_schema.get('items', {})
                    if '$ref' in items:
                        ref_name = items['$ref'].split('/')[-1]
                        html += f'[<span class="json-ref">{ref_name}</span>]'
                    else:
                        html += f'[<span class="json-type">{items.get("type", "any")}</span>]'
                else:
                    html += f'<span class="json-type">"{prop_type}"</span>'
                
                html += ','
                
                if prop_desc:
                    html += f'  <span class="json-comment">// {prop_desc}</span>'
                if is_required:
                    html += ' <span class="json-required">(required)</span>'
            
            html += f'\n{indent}}}</pre></div>'
        elif schema_type == 'array':
            items = schema.get('items', {})
            html += '[' + self.render_schema(items, level) + ']'
        else:
            html += f'<span class="json-type">{schema_type}</span>'
        
        html += '</div>'
        return html
    
    def generate_endpoint_page(self, method: str, path: str, operation: Dict, domain: str, output_dir: str) -> str:
        """Generate a static HTML page for a single endpoint"""
        
        endpoint_id = self.get_endpoint_id(method, path)
        domain_info = self.doc_gen.get_domain_info(domain)
        
        # Extract operation details
        operation_id = operation.get('operationId', '')
        summary = operation.get('summary', '')
        description = operation.get('description', '')
        parameters = operation.get('parameters', [])
        request_body = operation.get('requestBody', {})
        responses = operation.get('responses', {})
        security = operation.get('security', [])
        
        # Method color
        method_colors = {
            'GET': '#4caf50',
            'POST': '#2196f3',
            'PUT': '#ff9800',
            'DELETE': '#f44336',
            'PATCH': '#9c27b0'
        }
        method_color = method_colors.get(method.upper(), '#666')
        
        # Build parameters HTML
        params_html = ''
        if parameters:
            params_html = '<h2>Parameters</h2><div class="params-list">'
            for param in parameters:
                param_name = param.get('name', '')
                param_in = param.get('in', '')
                param_required = '(required)' if param.get('required') else '(optional)'
                param_desc = param.get('description', 'No description')
                param_schema = param.get('schema', {})
                param_type = param_schema.get('type', 'string')
                
                params_html += f'''
                <div class="param">
                    <div class="param-header">
                        <span class="param-name">{param_name}</span>
                        <span class="param-badge badge-{param_in}">{param_in}</span>
                        <span class="param-type">{param_type}</span>
                        <span class="param-required">{param_required}</span>
                    </div>
                    <div class="param-description">{param_desc}</div>
                </div>
                '''
            params_html += '</div>'
        
        # Build request body HTML with full schema details
        request_html = ''
        if request_body:
            request_desc = request_body.get('description', '')
            content = request_body.get('content', {})
            if 'application/json' in content:
                schema = content['application/json'].get('schema', {})
                schema_html = self.render_schema(schema)
                
                request_html = f'''
                <h2>Request Body</h2>
                {f'<div class="section-description">{request_desc}</div>' if request_desc else ''}
                <div class="request-body">
                    <div class="content-type-badge">application/json</div>
                    {schema_html}
                </div>
                '''
        
        # Build responses HTML with full schema details
        responses_html = '<h2>Responses</h2><div class="responses-list">'
        for status_code, response in responses.items():
            response_desc = response.get('description', '')
            content = response.get('content', {})
            schema_html = ''
            
            if 'application/json' in content:
                schema = content['application/json'].get('schema', {})
                schema_html = f'<div class="response-body">{self.render_schema(schema)}</div>'
            
            status_class = 'success' if status_code.startswith('2') else 'error' if status_code.startswith(('4', '5')) else 'info'
            
            responses_html += f'''
            <div class="response response-{status_class}">
                <div class="response-header">
                    <span class="response-code">{status_code}</span>
                    <span class="response-desc">{response_desc}</span>
                </div>
                {schema_html}
            </div>
            '''
        responses_html += '</div>'
        
        # Generate HTML
        html = f'''<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{method.upper()} {path} - InsureTech API</title>
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
        
        .endpoint-title {{
            display: flex;
            align-items: center;
            gap: 15px;
            margin-bottom: 10px;
        }}
        
        .method-badge {{
            background: {method_color};
            padding: 8px 16px;
            border-radius: 6px;
            font-weight: 600;
            font-size: 1.1em;
        }}
        
        .endpoint-path {{
            font-size: 1.8em;
            font-weight: 600;
            font-family: 'Courier New', monospace;
        }}
        
        .endpoint-summary {{
            font-size: 1.2em;
            opacity: 0.95;
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
        
        .param {{
            background: #f8f9fa;
            border-left: 4px solid #667eea;
            padding: 15px;
            margin-bottom: 15px;
            border-radius: 6px;
        }}
        
        .param-header {{
            display: flex;
            gap: 10px;
            align-items: center;
            margin-bottom: 8px;
        }}
        
        .param-name {{
            font-family: 'Courier New', monospace;
            font-weight: 600;
            font-size: 1.1em;
            color: #333;
        }}
        
        .param-badge {{
            padding: 3px 8px;
            border-radius: 4px;
            font-size: 0.75em;
            font-weight: 600;
            text-transform: uppercase;
        }}
        
        .badge-query {{ background: #e3f2fd; color: #1565c0; }}
        .badge-path {{ background: #fff3e0; color: #ef6c00; }}
        .badge-header {{ background: #f3e5f5; color: #7b1fa2; }}
        .badge-body {{ background: #e8f5e9; color: #2e7d32; }}
        
        .param-type {{
            color: #666;
            font-family: 'Courier New', monospace;
            font-size: 0.9em;
        }}
        
        .param-required {{
            color: #f44336;
            font-size: 0.85em;
        }}
        
        .param-description {{
            color: #666;
            margin-top: 5px;
        }}
        
        .response {{
            background: #f8f9fa;
            border-left: 4px solid #666;
            padding: 15px;
            margin-bottom: 15px;
            border-radius: 6px;
        }}
        
        .response-success {{ border-color: #4caf50; }}
        .response-error {{ border-color: #f44336; }}
        .response-info {{ border-color: #2196f3; }}
        
        .response-header {{
            display: flex;
            gap: 15px;
            align-items: center;
        }}
        
        .response-code {{
            font-family: 'Courier New', monospace;
            font-weight: 600;
            font-size: 1.2em;
            color: #333;
        }}
        
        .response-desc {{
            color: #666;
        }}
        
        .response-schema {{
            margin-top: 10px;
            color: #666;
        }}
        
        .schema-ref {{
            font-family: 'Courier New', monospace;
            background: #e3f2fd;
            padding: 10px;
            border-radius: 4px;
            margin-bottom: 10px;
        }}
        
        .description {{
            color: #555;
            line-height: 1.8;
            margin-bottom: 20px;
        }}
        
        .json-block {{
            background: #1e1e1e;
            color: #d4d4d4;
            padding: 20px;
            border-radius: 8px;
            overflow-x: auto;
            font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
            font-size: 0.95em;
            line-height: 1.6;
            margin: 15px 0;
        }}
        
        .json-block pre {{
            margin: 0;
            color: #d4d4d4;
        }}
        
        .json-key {{
            color: #9cdcfe;
            font-weight: 600;
        }}
        
        .json-type {{
            color: #ce9178;
        }}
        
        .json-ref {{
            color: #4ec9b0;
            font-style: italic;
        }}
        
        .json-comment {{
            color: #6a9955;
            font-style: italic;
        }}
        
        .json-required {{
            color: #f48771;
            font-weight: 600;
            font-size: 0.85em;
        }}
        
        .schema-ref-name {{
            font-family: 'Courier New', monospace;
            font-weight: 600;
            color: #667eea;
            font-size: 1.1em;
            margin-bottom: 10px;
            padding: 8px 12px;
            background: #f0f4ff;
            border-radius: 6px;
            display: inline-block;
        }}
        
        .schema-properties {{
            margin: 10px 0;
        }}
        
        .content-type-badge {{
            display: inline-block;
            background: #4caf50;
            color: white;
            padding: 4px 12px;
            border-radius: 4px;
            font-size: 0.85em;
            font-weight: 600;
            margin-bottom: 15px;
        }}
        
        .section-description {{
            background: #f8f9fa;
            padding: 15px;
            border-left: 4px solid #667eea;
            margin-bottom: 20px;
            border-radius: 4px;
            color: #555;
            line-height: 1.6;
        }}
        
        .response-body {{
            margin-top: 15px;
        }}
        
        .request-body {{
            margin-top: 15px;
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
                <a href="index.html">← Back to Documentation</a> / 
                <a href="index.html#{domain}">{domain_info['name']}</a>
            </div>
            <div class="endpoint-title">
                <span class="method-badge">{method.upper()}</span>
                <span class="endpoint-path">{path}</span>
            </div>
            <div class="endpoint-summary">{summary}</div>
        </div>
    </div>
    
    <div class="container">
        {f'<div class="card"><div class="description">{description}</div></div>' if description else ''}
        
        {f'<div class="card">{params_html}</div>' if params_html else ''}
        
        {f'<div class="card">{request_html}</div>' if request_html else ''}
        
        <div class="card">
            {responses_html}
        </div>
        
        <div style="text-align: center; margin-top: 40px;">
            <a href="index.html" class="back-link">← Back to Documentation Hub</a>
        </div>
    </div>
</body>
</html>'''
        
        # Write to file
        output_path = os.path.join(output_dir, f'endpoint_{endpoint_id}.html')
        with open(output_path, 'w', encoding='utf-8') as f:
            f.write(html)
        
        return endpoint_id
    
    def generate_all_endpoint_pages(self, output_dir: str) -> Dict[str, str]:
        """Generate pages for all endpoints"""
        
        print("Generating individual endpoint pages...")
        
        # Organize data
        organized = self.doc_gen.organize_data()
        
        # Ensure output directory exists
        os.makedirs(output_dir, exist_ok=True)
        
        # Generate page for each endpoint
        endpoint_pages = {}  # Maps "domain_method_path" -> "endpoint_xxx.html"
        
        paths = self.full_spec.get('paths', {})
        for path, methods in paths.items():
            for method in ['get', 'post', 'put', 'delete', 'patch']:
                if method not in methods:
                    continue
                
                operation = methods[method]
                
                # Determine domain
                parts = path.strip('/').split('/')
                if len(parts) >= 2 and parts[0] == 'v1':
                    raw_domain = parts[1]
                    # Use same logic as doc_generator to get canonical domain
                    domain = self._normalize_domain(raw_domain)
                else:
                    domain = 'common'
                
                endpoint_id = self.generate_endpoint_page(method, path, operation, domain, output_dir)
                endpoint_pages[f"{method}_{path}"] = f"endpoint_{endpoint_id}.html"
                
        print(f"  ✓ Generated {len(endpoint_pages)} endpoint pages")
        return endpoint_pages
    
    def _normalize_domain(self, raw_domain: str) -> str:
        """Normalize domain name (same logic as doc_generator)"""
        if ':' in raw_domain:
            base = raw_domain.split(':')[0]
            hyphen_map = {
                'kyc-verifications': 'kyc',
                'notification-templates': 'notification',
                'api-keys': 'apikey'
            }
            if base in hyphen_map:
                return hyphen_map[base]
            plural_map = {
                'notifications': 'notification',
                'payments': 'payment',
                'users': 'authn',
                'beneficiaries': 'beneficiary',
                'policies': 'policy'
            }
            return plural_map.get(base, base)
        elif '-' in raw_domain:
            domain_map = {
                'kyc-verifications': 'kyc',
                'api-keys': 'apikey',
                'notification-templates': 'notification'
            }
            return domain_map.get(raw_domain, raw_domain.split('-')[0])
        else:
            singular_map = {
                'auth': 'authn',
                'users': 'authn',
                'notifications': 'notification',
                'payments': 'payment',
                'policies': 'policy',
                'claims': 'claims',
                'products': 'products'
            }
            return singular_map.get(raw_domain, raw_domain)


def main():
    """CLI entry point"""
    import argparse
    
    parser = argparse.ArgumentParser(description='Generate endpoint-specific pages')
    parser.add_argument('--spec', default='../openapi.yaml', help='Path to OpenAPI spec')
    parser.add_argument('--output-dir', default='../docs', help='Output directory')
    
    args = parser.parse_args()
    
    # Initialize generators
    doc_gen = DocGenerator(args.spec)
    page_gen = EndpointPageGenerator(args.spec, doc_gen)
    
    # Generate all pages
    endpoint_pages = page_gen.generate_all_endpoint_pages(args.output_dir)
    
    print(f"\n✓ Endpoint pages generated in {args.output_dir}")
    print(f"  Total: {len(endpoint_pages)} pages")


if __name__ == '__main__':
    main()
