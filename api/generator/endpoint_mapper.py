"""
Endpoint Mapper - Maps all HTTP methods and custom actions per endpoint
Provides complete visibility into API structure
"""

class EndpointMapper:
    def __init__(self):
        self.endpoints = {}  # path -> endpoint info
        
    def add_endpoint(self, path_url, http_method, operation_details):
        """
        Add an endpoint mapping
        
        Args:
            path_url: The URL path
            http_method: HTTP method (get, post, patch, delete, put)
            operation_details: Dict with operation info
        """
        if path_url not in self.endpoints:
            self.endpoints[path_url] = {
                'path': path_url,
                'is_collection': '{' not in path_url,
                'is_custom_action': ':' in path_url,
                'methods': {},
                'services': set(),
                'resource_name': self._extract_resource_name(path_url)
            }
        
        # Add method
        self.endpoints[path_url]['methods'][http_method] = operation_details
        self.endpoints[path_url]['services'].add(operation_details.get('service', 'unknown'))
        
        # Extract custom action if present
        if ':' in path_url:
            action = path_url.split(':')[-1]
            self.endpoints[path_url]['custom_action'] = action
    
    def _extract_resource_name(self, path_url):
        """Extract the resource name from path"""
        # Remove version prefix and query params
        path = path_url.split('?')[0]
        parts = [p for p in path.split('/') if p and not p.startswith('{') and not p.startswith('v')]
        
        # Get last non-parameter part before custom action
        if ':' in path:
            parts = path.split(':')[0].split('/')
            parts = [p for p in parts if p and not p.startswith('{') and not p.startswith('v')]
        
        return parts[-1] if parts else 'unknown'
    
    def get_endpoint_info(self, path_url):
        """Get complete information about an endpoint"""
        return self.endpoints.get(path_url)
    
    def get_all_endpoints(self):
        """Get all endpoints sorted by path"""
        return dict(sorted(self.endpoints.items()))
    
    def get_endpoints_by_resource(self, resource_name):
        """Get all endpoints for a specific resource"""
        return {
            path: info for path, info in self.endpoints.items()
            if info['resource_name'] == resource_name
        }
    
    def get_custom_action_endpoints(self):
        """Get all custom action endpoints"""
        return {
            path: info for path, info in self.endpoints.items()
            if info['is_custom_action']
        }
    
    def export_endpoint_map(self):
        """
        Export complete endpoint map in structured format
        Useful for documentation and API review
        """
        map_data = []
        
        for path, info in sorted(self.endpoints.items()):
            endpoint_data = {
                'path': path,
                'resource': info['resource_name'],
                'type': 'custom_action' if info['is_custom_action'] else ('collection' if info['is_collection'] else 'resource'),
                'methods': {},
                'services': list(info['services'])
            }
            
            if 'custom_action' in info:
                endpoint_data['custom_action'] = info['custom_action']
            
            # Add method details
            for method, op_details in info['methods'].items():
                endpoint_data['methods'][method.upper()] = {
                    'operation_id': op_details.get('operation_id'),
                    'summary': op_details.get('summary'),
                    'service': op_details.get('service')
                }
            
            map_data.append(endpoint_data)
        
        return map_data
    
    def print_endpoint_map(self):
        """Print human-readable endpoint map"""
        print("\n" + "="*80)
        print("COMPLETE ENDPOINT MAP")
        print("="*80)
        
        # Group by resource
        by_resource = {}
        for path, info in self.endpoints.items():
            resource = info['resource_name']
            if resource not in by_resource:
                by_resource[resource] = []
            by_resource[resource].append((path, info))
        
        # Print by resource
        for resource in sorted(by_resource.keys()):
            endpoints = by_resource[resource]
            print(f"\n📦 RESOURCE: {resource.upper()}")
            print(f"   Endpoints: {len(endpoints)}")
            
            for path, info in sorted(endpoints):
                # Determine endpoint type
                if info['is_custom_action']:
                    icon = "⚡"
                    type_str = f"CUSTOM ACTION: {info.get('custom_action', 'unknown')}"
                elif info['is_collection']:
                    icon = "📋"
                    type_str = "COLLECTION"
                else:
                    icon = "📄"
                    type_str = "RESOURCE"
                
                print(f"\n   {icon} {path}")
                print(f"      Type: {type_str}")
                
                # Show methods
                methods = info['methods']
                print(f"      Methods: {', '.join(m.upper() for m in sorted(methods.keys()))}")
                
                # Show operations
                for method in sorted(methods.keys()):
                    op = methods[method]
                    service = op.get('service', 'unknown')
                    summary = op.get('summary', 'N/A')
                    print(f"         {method.upper():6} → {service}.{summary}")
                
                # Show contributing services
                if len(info['services']) > 1:
                    print(f"      ⚠️  Multi-service: {', '.join(sorted(info['services']))}")
        
        # Statistics
        total_endpoints = len(self.endpoints)
        total_operations = sum(len(info['methods']) for info in self.endpoints.values())
        custom_actions = len([e for e in self.endpoints.values() if e['is_custom_action']])
        collections = len([e for e in self.endpoints.values() if e['is_collection'] and not e['is_custom_action']])
        resources = len([e for e in self.endpoints.values() if not e['is_collection'] and not e['is_custom_action']])
        
        print(f"\n📊 STATISTICS:")
        print(f"   Total endpoints: {total_endpoints}")
        print(f"   Total operations: {total_operations}")
        print(f"   Collections: {collections}")
        print(f"   Resources: {resources}")
        print(f"   Custom actions: {custom_actions}")
        
        print("\n" + "="*80)
    
    def export_to_markdown(self, output_file):
        """Export endpoint map to markdown file for documentation"""
        with open(output_file, 'w', encoding='utf-8') as f:
            f.write("# API Endpoint Map\n\n")
            f.write("Auto-generated endpoint documentation\n\n")
            
            # Group by resource
            by_resource = {}
            for path, info in self.endpoints.items():
                resource = info['resource_name']
                if resource not in by_resource:
                    by_resource[resource] = []
                by_resource[resource].append((path, info))
            
            # Write table of contents
            f.write("## Table of Contents\n\n")
            for resource in sorted(by_resource.keys()):
                f.write(f"- [{resource.title()}](#{resource.lower()})\n")
            
            # Write details
            for resource in sorted(by_resource.keys()):
                endpoints = by_resource[resource]
                f.write(f"\n## {resource.title()}\n\n")
                
                # Create table
                f.write("| Path | Type | Methods | Operations |\n")
                f.write("|------|------|---------|------------|\n")
                
                for path, info in sorted(endpoints):
                    # Type
                    if info['is_custom_action']:
                        type_str = f"Custom Action: `{info.get('custom_action')}`"
                    elif info['is_collection']:
                        type_str = "Collection"
                    else:
                        type_str = "Resource"
                    
                    # Methods
                    methods_str = ", ".join(f"`{m.upper()}`" for m in sorted(info['methods'].keys()))
                    
                    # Operations
                    ops = []
                    for method in sorted(info['methods'].keys()):
                        op = info['methods'][method]
                        ops.append(f"{method.upper()}: {op.get('summary', 'N/A')}")
                    ops_str = "<br>".join(ops)
                    
                    f.write(f"| `{path}` | {type_str} | {methods_str} | {ops_str} |\n")
            
            # Statistics
            f.write(f"\n## Statistics\n\n")
            total_endpoints = len(self.endpoints)
            total_operations = sum(len(info['methods']) for info in self.endpoints.values())
            custom_actions = len([e for e in self.endpoints.values() if e['is_custom_action']])
            collections = len([e for e in self.endpoints.values() if e['is_collection'] and not e['is_custom_action']])
            resources = len([e for e in self.endpoints.values() if not e['is_collection'] and not e['is_custom_action']])
            
            f.write(f"- **Total Endpoints:** {total_endpoints}\n")
            f.write(f"- **Total Operations:** {total_operations}\n")
            f.write(f"- **Collections:** {collections}\n")
            f.write(f"- **Resources:** {resources}\n")
            f.write(f"- **Custom Actions:** {custom_actions}\n")
