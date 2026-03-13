#!/usr/bin/env python3
"""
Enhanced Documentation Generator
Generates organized API documentation with tabs, groups, and modern UI
"""

import yaml
import json
import os
from pathlib import Path
from typing import Dict, List, Any
from collections import defaultdict

class DocGenerator:
    """Generate enhanced documentation structure"""
    
    def __init__(self, openapi_path: str, schema_summary_path: str = None):
        self.openapi_path = openapi_path
        with open(openapi_path, 'r', encoding='utf-8') as f:
            self.spec = yaml.safe_load(f)
        
        # Load schema summary if available
        self.schema_summary = None
        if schema_summary_path and os.path.exists(schema_summary_path):
            with open(schema_summary_path, 'r', encoding='utf-8') as f:
                self.schema_summary = json.load(f)
                print(f"✅ Loaded schema summary from {schema_summary_path}")
        
        # Schema group icons and descriptions (fixed mapping)
        self.schema_icons = {
            'authn_schema': '🔐',
            'authz_schema': '🛡️',
            'insurance_schema': '🏥',
            'payment_schema': '💳',
            'partner_schema': '🤝',
            'notification_schema': '🔔',
            'storage_schema': '📁',
            'support_schema': '💬',
            'analytics_schema': '📊',
            'ai_schema': '🤖',
            'iot_schema': '📡',
            'workflow_schema': '⚙️',
            'tenant_schema': '🏢',
            'compliance_schema': '📋',
            'webrtc': '📹',
            'b2b_schema': '🏢',
            'public': '🌐'  # Default public schema
        }
        
        self.schema_descriptions = {
            'authn_schema': 'User authentication, sessions, and identity management',
            'authz_schema': 'Role-based access control and permissions',
            'insurance_schema': 'Core insurance entities: policies, claims, beneficiaries',
            'payment_schema': 'Payment processing, transactions, and TigerBeetle integration',
            'partner_schema': 'Business partners, agents, and affiliates',
            'notification_schema': 'Push notifications, SMS, and email delivery',
            'storage_schema': 'File storage, documents, and S3/Spaces integration',
            'support_schema': 'Customer support tickets and knowledge base',
            'analytics_schema': 'Business intelligence, metrics, and reporting',
            'ai_schema': 'AI agents, assistants, and intelligent automation',
            'iot_schema': 'IoT devices, telematics, and sensor data',
            'workflow_schema': 'Business process automation and task management',
            'tenant_schema': 'Multi-tenant and white-label management',
            'compliance_schema': 'Audit logs, compliance tracking, and regulatory reports',
            'webrtc': 'Real-time video/audio communication',
            'b2b_schema': 'B2B operations, departments, and employees',
            'public': 'Public schema tables'
        }
        
        # Domain metadata with friendly names and descriptions
        self.domain_info = {
            'ai': {'name': 'AI & Agents', 'icon': '🤖', 'description': 'AI assistants and intelligent automation'},
            'analytics': {'name': 'Analytics', 'icon': '📊', 'description': 'Metrics, reports, and business intelligence'},
            'apikey': {'name': 'API Keys', 'icon': '🔑', 'description': 'API key management and authentication'},
            'audit': {'name': 'Audit & Compliance', 'icon': '📋', 'description': 'Audit logs and compliance tracking'},
            'authn': {'name': 'Authentication', 'icon': '🔐', 'description': 'Hybrid authentication: Server-side (web) + JWT (mobile)'},
            'authz': {'name': 'Authorization', 'icon': '🛡️', 'description': 'Roles, permissions, and access control (RBAC)'},
            'beneficiary': {'name': 'Beneficiaries', 'icon': '👥', 'description': 'Policy beneficiaries and nominees'},
            'claims': {'name': 'Claims', 'icon': '📝', 'description': 'Insurance claims processing'},
            'commission': {'name': 'Commissions', 'icon': '💰', 'description': 'Agent commissions and payouts'},
            'common': {'name': 'Common Types', 'icon': '🔧', 'description': 'Shared types, enums, and utilities'},
            'document': {'name': 'Documents', 'icon': '📄', 'description': 'Document generation and templates'},
            'endorsement': {'name': 'Endorsements', 'icon': '✍️', 'description': 'Policy modifications and amendments'},
            'fraud': {'name': 'Fraud Detection', 'icon': '🚨', 'description': 'Fraud detection and prevention'},
            'insurer': {'name': 'Insurers', 'icon': '🏢', 'description': 'Insurance companies and products'},
            'iot': {'name': 'IoT & Telematics', 'icon': '📡', 'description': 'IoT devices and sensor data'},
            'kyc': {'name': 'KYC', 'icon': '✅', 'description': 'Know Your Customer verification'},
            'mfs': {'name': 'Mobile Finance', 'icon': '📱', 'description': 'Mobile financial services (bKash, Nagad, Rocket)'},
            'notification': {'name': 'Notifications', 'icon': '🔔', 'description': 'Push, SMS, and email notifications'},
            'partner': {'name': 'Partners', 'icon': '🤝', 'description': 'Business partners and affiliates'},
            'payment': {'name': 'Payments', 'icon': '💳', 'description': 'Payment processing and TigerBeetle integration'},
            'policy': {'name': 'Policies', 'icon': '📑', 'description': 'Insurance policies lifecycle'},
            'products': {'name': 'Products', 'icon': '🎁', 'description': 'Insurance products catalog'},
            'refund': {'name': 'Refunds', 'icon': '↩️', 'description': 'Payment refunds and cancellations'},
            'renewal': {'name': 'Renewals', 'icon': '🔄', 'description': 'Policy renewal management'},
            'report': {'name': 'Reports', 'icon': '📈', 'description': 'Custom reports and schedules'},
            'storage': {'name': 'Storage & Files', 'icon': '💾', 'description': 'File storage and S3/Spaces integration'},
            'support': {'name': 'Support', 'icon': '💬', 'description': 'Customer support and tickets'},
            'task': {'name': 'Tasks', 'icon': '✓', 'description': 'Task management and workflow'},
            'tenant': {'name': 'Tenants', 'icon': '🏠', 'description': 'Multi-tenant and white-label management'},
            'underwriting': {'name': 'Underwriting', 'icon': '⚖️', 'description': 'Risk assessment and quotes'},
            'voice': {'name': 'Voice', 'icon': '🎤', 'description': 'Voice commands and assistance'},
            'webrtc': {'name': 'WebRTC', 'icon': '📹', 'description': 'Real-time video/audio communication'},
            'workflow': {'name': 'Workflows', 'icon': '⚙️', 'description': 'Business process automation'}
        }
    
    def organize_data(self) -> Dict[str, Any]:
        """Organize OpenAPI spec into structured groups"""
        
        organized = {
            'apis': defaultdict(list),
            'schemas': defaultdict(list),
            'enums': [],
            'dtos': defaultdict(list)
        }
        
        # Organize paths/APIs by domain
        paths = self.spec.get('paths', {})
        for path, methods in paths.items():
            # Extract domain from path: /v1/{domain}/...
            parts = path.strip('/').split('/')
            if len(parts) >= 2 and parts[0] == 'v1':
                # Extract base domain: notifications:mark-as-read -> notifications
                # api-keys -> apikey, users -> authn, etc.
                raw_domain = parts[1]
                
                # Handle action suffixes (notifications:mark-as-read, payments:reconcile, kyc-verifications:pending)
                if ':' in raw_domain:
                    base = raw_domain.split(':')[0]
                    
                    # First check if base needs hyphen mapping
                    hyphen_map = {
                        'kyc-verifications': 'kyc',
                        'notification-templates': 'notification',
                        'api-keys': 'apikey'
                    }
                    if base in hyphen_map:
                        domain = hyphen_map[base]
                    else:
                        # Then normalize plurals
                        plural_to_singular = {
                            'notifications': 'notification',
                            'payments': 'payment',
                            'users': 'authn',
                            'beneficiaries': 'beneficiary',
                            'policies': 'policy',
                            'products': 'products',  # Keep as is
                            'claims': 'claims'  # Keep as is
                        }
                        domain = plural_to_singular.get(base, base)
                # Handle hyphenated domains
                elif '-' in raw_domain:
                    # Map known multi-word domains
                    domain_map = {
                        'api-keys': 'apikey',
                        'audit-events': 'audit',
                        'audit-logs': 'audit',
                        'compliance-logs': 'audit',
                        'compliance-reports': 'audit',
                        'notification-templates': 'notification',
                        'commission-payouts': 'commission',
                        'document-templates': 'document',
                        'fraud-alerts': 'fraud',
                        'fraud-cases': 'fraud',
                        'fraud-checks': 'fraud',
                        'fraud-rules': 'fraud',
                        'insurer-products': 'insurer',
                        'kyc-verifications': 'kyc',
                        'kyc-verifications:pending': 'kyc',
                        'knowledge-base': 'support',
                        'renewal-schedules': 'renewal',
                        'report-definitions': 'report',
                        'report-executions': 'report',
                        'report-schedules': 'report',
                        'voice-sessions': 'voice',
                        'workflow-definitions': 'workflow',
                        'workflow-instances': 'workflow',
                        'workflow-tasks': 'workflow'
                    }
                    domain = domain_map.get(raw_domain, raw_domain.split('-')[0])
                else:
                    # Map plural paths and aliases to canonical domains
                    singular_map = {
                        # Auth related
                        'auth': 'authn',  # auth -> authn (authentication service)
                        'users': 'authn',
                        'roles': 'authz',
                        
                        # Beneficiaries
                        'beneficiaries': 'beneficiary',
                        'entities': 'beneficiary',
                        'recipients': 'beneficiary',
                        
                        # Plural to singular
                        'claims': 'claims',  # Keep as is (service name)
                        'commissions': 'commission',
                        'documents': 'document',
                        'endorsements': 'endorsement',
                        'insurers': 'insurer',
                        'notifications': 'notification',
                        'partners': 'partner',
                        'payments': 'payment',
                        'policies': 'policy',
                        'products': 'products',  # Keep as is (service name)
                        'quotes': 'underwriting',
                        'refunds': 'refund',
                        'renewals': 'renewal',
                        'reports': 'report',
                        'tasks': 'task',
                        'tenants': 'tenant',
                        
                        # Support related
                        'faqs': 'support',
                        'tickets': 'support'
                    }
                    domain = singular_map.get(raw_domain, raw_domain)
                
                for method, operation in methods.items():
                    if method in ['get', 'post', 'put', 'delete', 'patch']:
                        api_info = {
                            'path': path,
                            'method': method.upper(),
                            'operation_id': operation.get('operationId', ''),
                            'summary': operation.get('summary', ''),
                            'description': operation.get('description', '')
                        }
                        organized['apis'][domain].append(api_info)
        
        # Organize schemas by domain
        schemas = self.spec.get('components', {}).get('schemas', {})
        for schema_name, schema_def in schemas.items():
            # Extract domain from schema name (package prefix)
            if '.' in schema_name:
                parts = schema_name.split('.')
                # Handle insuretech.authn.User -> authn
                domain = parts[1] if len(parts) > 1 else 'common'
            else:
                # Infer from name patterns
                domain = self._infer_domain(schema_name)
            
            schema_info = {
                'name': schema_name,
                'type': schema_def.get('type', 'object'),
                'description': schema_def.get('description', ''),
                'properties_count': len(schema_def.get('properties', {}))
            }
            
            # Categorize as DTO or Schema
            if schema_name.endswith(('Request', 'Response')):
                organized['dtos'][domain].append(schema_info)
            elif schema_def.get('type') == 'string' and 'enum' in schema_def:
                organized['enums'].append(schema_info)
            else:
                organized['schemas'][domain].append(schema_info)
        
        return organized
    
    def _infer_domain(self, schema_name: str) -> str:
        """Infer domain from schema name"""
        lower_name = schema_name.lower()
        
        keywords = {
            'auth': 'authn',
            'user': 'authn',
            'session': 'authn',
            'role': 'authz',
            'permission': 'authz',
            'claim': 'claims',
            'policy': 'policy',
            'payment': 'payment',
            'product': 'products',
            'partner': 'partner',
            'notification': 'notification',
            'document': 'document',
            'task': 'task',
            'workflow': 'workflow'
        }
        
        for keyword, domain in keywords.items():
            if keyword in lower_name:
                return domain
        
        return 'common'
    
    def get_domain_info(self, domain: str) -> Dict[str, str]:
        """Get metadata for a domain"""
        return self.domain_info.get(domain, {
            'name': domain.capitalize(),
            'icon': '📦',
            'description': f'{domain.capitalize()} services'
        })
    
    def generate_documentation(self, output_path: str, endpoint_pages: Dict[str, str] = None):
        """Generate the complete documentation HTML"""
        from doc_templates import get_main_template
        
        print("Organizing API documentation...")
        organized = self.organize_data()
        
        # Generate domain cards for APIs
        apis_html_parts = []
        for domain in sorted(organized['apis'].keys()):
            items = organized['apis'][domain]
            if not items:
                continue
            
            info = self.get_domain_info(domain)
            
            # Always use modal for cards - domain pages are linked from inside the modal
            onclick = f"showDetail('{domain}', 'apis')"
            
            apis_html_parts.append(f'''
                <div class="domain-card" onclick="{onclick}">
                    <div class="domain-icon">{info['icon']}</div>
                    <div class="domain-name">{info['name']}</div>
                    <div class="domain-description">{info['description']}</div>
                    <span class="domain-count">{len(items)} endpoints</span>
                </div>
            ''')
        
        # Generate schema group cards (NEW - from proto analysis)
        schema_groups_html_parts = []
        if self.schema_summary:
            # Load API mapping if available
            schema_api_mapping = {}
            mapping_file = os.path.join(os.path.dirname(self.openapi_path), "schema_api_mapping.json")
            if os.path.exists(mapping_file):
                with open(mapping_file, 'r') as f:
                    schema_api_mapping = json.load(f)
            
            for schema_name in sorted(self.schema_summary['schema_groups'].keys()):
                tables = self.schema_summary['schema_groups'][schema_name]
                icon = self.schema_icons.get(schema_name, '📦')
                description = self.schema_descriptions.get(schema_name, 'Database schema group')
                
                # Get API count from mapping
                api_count = 0
                if schema_api_mapping and 'schema_groups' in schema_api_mapping:
                    api_count = schema_api_mapping['schema_groups'].get(schema_name, {}).get('api_count', 0)
                
                onclick = f"showSchemaDetail('{schema_name}')"
                
                schema_groups_html_parts.append(f'''
                    <div class="domain-card" onclick="{onclick}">
                        <div class="domain-icon">{icon}</div>
                        <div class="domain-name">{schema_name.replace('_', ' ').title()}</div>
                        <div class="domain-description">{description}</div>
                        <span class="domain-count">{len(tables)} tables • {api_count} APIs</span>
                    </div>
                ''')
        
        # Generate domain cards for Schemas - link to table pages (too many for modal)
        schemas_html_parts = []
        for domain in sorted(organized['schemas'].keys()):
            items = organized['schemas'][domain]
            if not items:
                continue
            
            info = self.get_domain_info(domain)
            table_page = f'schemas_{domain}.html'
            schemas_html_parts.append(f'''
                <div class="domain-card" onclick="window.location.href='{table_page}'">
                    <div class="domain-icon">{info['icon']}</div>
                    <div class="domain-name">{info['name']}</div>
                    <div class="domain-description">{info['description']}</div>
                    <span class="domain-count">{len(items)} schemas</span>
                </div>
            ''')
        
        # Generate enum items with links
        enums_html_parts = []
        for enum in sorted(organized['enums'], key=lambda x: x['name']):
            clean_name = enum['name'].replace('.', '_').replace(':', '_').lower()
            enum_page = f'enum_{clean_name}.html'
            enums_html_parts.append(f'''
                <div class="domain-card" onclick="window.location.href='{enum_page}'" style="cursor: pointer;">
                    <div class="domain-icon">🔢</div>
                    <div class="domain-name">{enum['name']}</div>
                    <div class="domain-description">{enum['description'] or 'No description'}</div>
                    <span class="domain-count">Enum</span>
                </div>
            ''')
        
        # Generate domain cards for DTOs - link to table pages (too many for modal)
        dtos_html_parts = []
        for domain in sorted(organized['dtos'].keys()):
            items = organized['dtos'][domain]
            if not items:
                continue
            
            info = self.get_domain_info(domain)
            table_page = f'dtos_{domain}.html'
            dtos_html_parts.append(f'''
                <div class="domain-card" onclick="window.location.href='{table_page}'">
                    <div class="domain-icon">{info['icon']}</div>
                    <div class="domain-name">{info['name']}</div>
                    <div class="domain-description">{info['description']}</div>
                    <span class="domain-count">{len(items)} DTOs</span>
                </div>
            ''')
        
        # Calculate totals - use proto data if available
        if self.schema_summary:
            stats = self.schema_summary['stats']
            total_schema_groups = stats['schema_groups_count']
            total_tables = stats['total_tables']
            total_enums = stats['enums_count']
            total_dtos = stats['dtos_count']
            total_events = stats['events_count']
        else:
            # Fallback to OpenAPI data
            total_schema_groups = len(set(list(organized['apis'].keys()) + list(organized['schemas'].keys()) + list(organized['dtos'].keys())))
            total_tables = sum(len(items) for items in organized['schemas'].values())
            total_enums = len(organized['enums'])
            total_dtos = sum(len(items) for items in organized['dtos'].values())
            total_events = 0
        
        total_apis = sum(len(items) for items in organized['apis'].values())
        
        # Prepare detail data for JavaScript (filter out empty domains)
        detail_data = {
            'apis': {domain: {'items': items} for domain, items in organized['apis'].items() if items},
            'schemas': {domain: {'items': items} for domain, items in organized['schemas'].items() if items},
            'dtos': {domain: {'items': items} for domain, items in organized['dtos'].items() if items}
        }
        
        # Add schema groups detail data
        if self.schema_summary:
            detail_data['schema_groups'] = {}
            for schema_name, tables in self.schema_summary['schema_groups'].items():
                detail_data['schema_groups'][schema_name] = {
                    'tables': tables,
                    'icon': self.schema_icons.get(schema_name, '📦'),
                    'description': self.schema_descriptions.get(schema_name, 'Database schema group')
                }
        
        # Generate HTML
        template = get_main_template()
        html = template.format(
            total_schema_groups=total_schema_groups,
            total_tables=total_tables,
            total_apis=total_apis,
            total_enums=total_enums,
            total_dtos=total_dtos,
            total_events=total_events,
            apis_content=''.join(apis_html_parts) if apis_html_parts else '<div class="empty-state"><div class="empty-state-icon">📭</div><p>No API endpoints found</p></div>',
            schema_groups_content=''.join(schema_groups_html_parts) if schema_groups_html_parts else '<div class="empty-state"><div class="empty-state-icon">📭</div><p>No schema groups found</p></div>',
            schemas_content=''.join(schemas_html_parts) if schemas_html_parts else '<div class="empty-state"><div class="empty-state-icon">📭</div><p>No schemas found</p></div>',
            enums_content=''.join(enums_html_parts) if enums_html_parts else '<div class="empty-state"><div class="empty-state-icon">📭</div><p>No enums found</p></div>',
            dtos_content=''.join(dtos_html_parts) if dtos_html_parts else '<div class="empty-state"><div class="empty-state-icon">📭</div><p>No DTOs found</p></div>',
            detail_data_json=json.dumps(detail_data),
            domain_info_json=json.dumps(self.domain_info)
        )
        
        # Write to file
        with open(output_path, 'w', encoding='utf-8') as f:
            f.write(html)
        
        print(f"✓ Documentation generated: {output_path}")
        print(f"  - {total_schema_groups} schema groups")
        print(f"  - {total_tables} tables")
        print(f"  - {total_apis} API endpoints")
        print(f"  - {total_enums} enums")
        print(f"  - {total_dtos} DTOs")
        if self.schema_summary:
            print(f"  - {total_events} events")


def main():
    """CLI entry point"""
    import argparse
    
    parser = argparse.ArgumentParser(description='Generate enhanced API documentation')
    parser.add_argument('--spec', default='../openapi.yaml', help='Path to OpenAPI spec')
    parser.add_argument('--schema-summary', default='../proto_schema_summary.json', help='Path to proto schema summary')
    parser.add_argument('--output', default='../docs/index.html', help='Output HTML file')
    parser.add_argument('--generate-endpoint-pages', action='store_true', help='Generate individual endpoint pages')
    
    args = parser.parse_args()
    
    generator = DocGenerator(args.spec, args.schema_summary)
    
    # Generate endpoint pages if requested
    endpoint_pages = None
    if args.generate_endpoint_pages:
        from endpoint_page_generator import EndpointPageGenerator
        output_dir = os.path.dirname(args.output)
        page_gen = EndpointPageGenerator(args.spec, generator)
        endpoint_pages = page_gen.generate_all_endpoint_pages(output_dir)
    
    generator.generate_documentation(args.output, endpoint_pages)


if __name__ == '__main__':
    main()
