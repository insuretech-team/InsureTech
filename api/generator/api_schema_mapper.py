"""
API to Schema Group Mapper
Maps API endpoints (domain-based) to database schema groups
"""

from typing import Dict, List, Set


class APISchemaMapper:
    """Maps API endpoints to their corresponding database schema groups"""
    
    def __init__(self):
        # Map API domains to schema groups they interact with
        self.domain_to_schemas = {
            # Authentication & Authorization
            'authn': ['authn_schema'],
            'auth': ['authn_schema'],
            'authz': ['authz_schema'],
            
            # Insurance Core
            'policy': ['insurance_schema', 'payment_schema'],
            'claims': ['insurance_schema', 'payment_schema'],
            'underwriting': ['insurance_schema'],
            'beneficiary': ['insurance_schema'],
            'endorsement': ['insurance_schema'],
            'renewal': ['insurance_schema'],
            'refund': ['payment_schema'],
            'products': ['insurance_schema'],
            'insurer': ['insurance_schema'],
            
            # Partners & Agents
            'partner': ['partner_schema', 'payment_schema'],
            'commission': ['payment_schema', 'partner_schema'],
            
            # Payments
            'payment': ['payment_schema'],
            'mfs': ['payment_schema'],
            
            # Documents & Storage
            'document': ['storage_schema'],
            'storage': ['storage_schema'],
            
            # Support
            'support': ['support_schema'],
            'ticket': ['support_schema'],
            'faq': ['support_schema'],
            
            # Notifications
            'notification': ['notification_schema'],
            
            # Analytics & Reporting
            'analytics': ['analytics_schema'],
            'report': ['analytics_schema'],
            
            # AI & Voice
            'ai': ['ai_schema'],
            'voice': ['ai_schema', 'authn_schema'],
            
            # IoT
            'iot': ['iot_schema'],
            
            # Workflow & Tasks
            'workflow': ['workflow_schema'],
            'task': ['workflow_schema'],
            
            # Tenant
            'tenant': ['tenant_schema'],
            
            # Compliance & Audit
            'audit': ['compliance_schema'],
            'compliance': ['compliance_schema'],
            
            # WebRTC
            'webrtc': ['webrtc'],
            
            # B2B
            'department': ['b2b_schema'],
            'employee': ['b2b_schema'],
            
            # KYC
            'kyc': ['authn_schema'],
            
            # Fraud
            'fraud': ['insurance_schema'],
            
            # API Keys
            'apikey': ['authn_schema'],
        }
        
        # Reverse mapping: schema -> domains
        self.schema_to_domains = {}
        for domain, schemas in self.domain_to_schemas.items():
            for schema in schemas:
                if schema not in self.schema_to_domains:
                    self.schema_to_domains[schema] = []
                if domain not in self.schema_to_domains[schema]:
                    self.schema_to_domains[schema].append(domain)
    
    def get_schemas_for_domain(self, domain: str) -> List[str]:
        """Get schema groups that a domain interacts with"""
        return self.domain_to_schemas.get(domain, [])
    
    def get_domains_for_schema(self, schema: str) -> List[str]:
        """Get API domains that interact with a schema group"""
        return self.schema_to_domains.get(schema, [])
    
    def map_apis_to_schemas(self, apis_by_domain: Dict) -> Dict:
        """Map API endpoints to their schema groups"""
        apis_by_schema = {}
        
        for domain, endpoints in apis_by_domain.items():
            schemas = self.get_schemas_for_domain(domain)
            
            for schema in schemas:
                if schema not in apis_by_schema:
                    apis_by_schema[schema] = []
                
                # Add endpoints with domain info
                for endpoint in endpoints:
                    apis_by_schema[schema].append({
                        **endpoint,
                        'domain': domain
                    })
        
        return apis_by_schema
    
    def get_schema_api_count(self, schema: str, apis_by_domain: Dict) -> int:
        """Get count of API endpoints that interact with a schema"""
        count = 0
        domains = self.get_domains_for_schema(schema)
        
        for domain in domains:
            if domain in apis_by_domain:
                count += len(apis_by_domain[domain])
        
        return count
    
    def generate_schema_summary(self, apis_by_domain: Dict, proto_summary: Dict) -> Dict:
        """Generate complete summary with API-to-schema mapping"""
        summary = {
            'schema_groups': {},
            'stats': {
                'total_schemas': 0,
                'total_tables': 0,
                'total_apis': 0
            }
        }
        
        # Get schema groups from proto
        schema_groups = proto_summary.get('schema_groups', {})
        
        for schema_name, tables in schema_groups.items():
            # Get API domains for this schema
            domains = self.get_domains_for_schema(schema_name)
            
            # Get API count
            api_count = self.get_schema_api_count(schema_name, apis_by_domain)
            
            summary['schema_groups'][schema_name] = {
                'table_count': len(tables),
                'tables': tables,
                'api_domains': domains,
                'api_count': api_count
            }
            
            summary['stats']['total_tables'] += len(tables)
            summary['stats']['total_apis'] += api_count
        
        summary['stats']['total_schemas'] = len(schema_groups)
        
        return summary


def main():
    """CLI entry point for testing"""
    import json
    
    mapper = APISchemaMapper()
    
    print("API to Schema Mapping:")
    print("=" * 80)
    
    # Test mapping
    test_domains = ['authn', 'policy', 'payment', 'partner']
    
    for domain in test_domains:
        schemas = mapper.get_schemas_for_domain(domain)
        print(f"\n{domain} → {', '.join(schemas)}")
    
    print("\n" + "=" * 80)
    print("\nSchema to API Domains:")
    print("=" * 80)
    
    # Show reverse mapping
    for schema in sorted(mapper.schema_to_domains.keys()):
        domains = mapper.get_domains_for_schema(schema)
        print(f"\n{schema}:")
        print(f"  Domains: {', '.join(domains)}")


if __name__ == '__main__':
    main()
