"""
Schema-Based Documentation Generator
Generates index.html with cards organized by database schema groups (not domains)
"""

import json
import os
from proto_schema_analyzer import ProtoSchemaAnalyzer


class SchemaBasedDocGenerator:
    """Generates documentation organized by database schema groups"""
    
    def __init__(self, proto_root: str = "../../proto"):
        self.proto_root = proto_root
        self.analyzer = ProtoSchemaAnalyzer(proto_root)
        self.summary = None
        
    def analyze_proto(self):
        """Analyze proto files to get schema structure"""
        print("Analyzing proto files...")
        self.summary = self.analyzer.scan_all_protos()
        self.analyzer.print_summary(self.summary)
        return self.summary
    
    def get_schema_icon(self, schema_name: str) -> str:
        """Get icon for schema group"""
        icon_map = {
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
            'b2b_schema': '🏢'
        }
        return icon_map.get(schema_name, '📦')
    
    def get_schema_description(self, schema_name: str) -> str:
        """Get description for schema group"""
        desc_map = {
            'authn_schema': 'Authentication and user management tables',
            'authz_schema': 'Authorization, roles, and permissions',
            'insurance_schema': 'Core insurance business logic tables',
            'payment_schema': 'Payment processing and financial transactions',
            'partner_schema': 'Partner and agent management',
            'notification_schema': 'Notification and alert system',
            'storage_schema': 'Document and file storage',
            'support_schema': 'Customer support and ticketing',
            'analytics_schema': 'Analytics, metrics, and reporting',
            'ai_schema': 'AI agents and voice processing',
            'iot_schema': 'IoT devices and telemetry data',
            'workflow_schema': 'Workflow and task management',
            'tenant_schema': 'Multi-tenant configuration',
            'compliance_schema': 'Audit logs and compliance tracking',
            'webrtc': 'Real-time communication (WebRTC)',
            'b2b_schema': 'B2B group insurance (departments, employees)'
        }
        return desc_map.get(schema_name, 'Database schema group')
    
    def generate_schema_cards_html(self) -> str:
        """Generate HTML cards for schema groups"""
        if not self.summary:
            return '<div class="empty-state">No schema groups found</div>'
        
        cards_html = []
        schema_groups = self.summary['schema_groups']
        
        for schema_name in sorted(schema_groups.keys()):
            tables = schema_groups[schema_name]
            icon = self.get_schema_icon(schema_name)
            description = self.get_schema_description(schema_name)
            
            # Generate table list for modal
            table_list = '<br>'.join([f"• {t['table_name']}" for t in tables[:10]])
            if len(tables) > 10:
                table_list += f"<br>• ... and {len(tables) - 10} more"
            
            card_html = f'''
                <div class="domain-card" onclick="showSchemaDetail('{schema_name}')">
                    <div class="domain-icon">{icon}</div>
                    <div class="domain-name">{schema_name.replace('_', ' ').title()}</div>
                    <div class="domain-description">{description}</div>
                    <span class="domain-count">{len(tables)} tables</span>
                </div>
            '''
            cards_html.append(card_html)
        
        return '\n'.join(cards_html)
    
    def generate_stats_html(self) -> dict:
        """Generate statistics for the stats bar"""
        if not self.summary:
            return {
                'schema_groups': 0,
                'total_tables': 0,
                'enums': 0,
                'dtos': 0,
                'events': 0
            }
        
        stats = self.summary['stats']
        return {
            'schema_groups': stats['schema_groups_count'],
            'total_tables': stats['total_tables'],
            'enums': stats['enums_count'],
            'dtos': stats['dtos_count'],
            'events': stats['events_count']
        }
    
    def generate_schema_detail_data(self) -> str:
        """Generate JavaScript data for schema details modal"""
        if not self.summary:
            return '{}'
        
        detail_data = {}
        schema_groups = self.summary['schema_groups']
        
        for schema_name, tables in schema_groups.items():
            detail_data[schema_name] = {
                'name': schema_name.replace('_', ' ').title(),
                'icon': self.get_schema_icon(schema_name),
                'description': self.get_schema_description(schema_name),
                'tables': [
                    {
                        'table_name': t['table_name'],
                        'message_name': t['message_name'],
                        'migration_order': t['migration_order'],
                        'package': t['package']
                    }
                    for t in tables
                ]
            }
        
        return json.dumps(detail_data, indent=2)


def main():
    """CLI entry point"""
    import argparse
    
    parser = argparse.ArgumentParser(description='Generate schema-based documentation')
    parser.add_argument('--proto-root', default='../../proto', help='Root directory for proto files')
    parser.add_argument('--output', default='../docs/schema_groups.json', help='Output JSON file')
    
    args = parser.parse_args()
    
    generator = SchemaBasedDocGenerator(args.proto_root)
    generator.analyze_proto()
    
    # Export schema detail data
    detail_data = generator.generate_schema_detail_data()
    with open(args.output, 'w', encoding='utf-8') as f:
        f.write(detail_data)
    
    print(f"\nSchema detail data exported to: {args.output}")
    
    # Print stats
    stats = generator.generate_stats_html()
    print(f"\nStats for index.html:")
    print(f"  Schema Groups: {stats['schema_groups']}")
    print(f"  Total Tables: {stats['total_tables']}")
    print(f"  Enums: {stats['enums']}")
    print(f"  DTOs: {stats['dtos']}")
    print(f"  Events: {stats['events']}")


if __name__ == '__main__':
    main()
