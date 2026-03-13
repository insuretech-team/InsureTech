"""
Path Validator - Smart collision detection and prevention
Ensures no data loss during API generation
"""

class PathValidator:
    def __init__(self):
        self.all_paths = {}  # path_url -> {method -> operation_details}
        self.collisions = []
        self.warnings = []
        self.service_contributions = {}  # path_url -> [service_names]
        
    def register_operation(self, path_url, http_method, operation_id, service_name, rpc_name):
        """
        Register an operation and detect collisions
        
        Args:
            path_url: The URL path (e.g., /v1/policies/{policy_id})
            http_method: HTTP method (get, post, patch, delete, put)
            operation_id: Unique operation identifier
            service_name: Name of the service
            rpc_name: Original RPC method name from proto
        """
        # Initialize path tracking
        if path_url not in self.all_paths:
            self.all_paths[path_url] = {}
            self.service_contributions[path_url] = set()
        
        # Track service contribution
        self.service_contributions[path_url].add(service_name)
        
        # Check for collision
        if http_method in self.all_paths[path_url]:
            # COLLISION DETECTED!
            existing = self.all_paths[path_url][http_method]
            collision = {
                'path': path_url,
                'method': http_method,
                'existing_service': existing['service'],
                'existing_operation': existing['operation_id'],
                'existing_rpc': existing['rpc_name'],
                'new_service': service_name,
                'new_operation': operation_id,
                'new_rpc': rpc_name,
            }
            self.collisions.append(collision)
            
            print(f"  ⚠️  COLLISION DETECTED: {http_method.upper()} {path_url}")
            print(f"      Service 1: {existing['service']}.{existing['rpc_name']} → {existing['operation_id']}")
            print(f"      Service 2: {service_name}.{rpc_name} → {operation_id}")
            print(f"      Action: Keeping first operation (Service 1)")
            
            return False  # Collision - do not register
        else:
            # No collision - register operation
            self.all_paths[path_url][http_method] = {
                'operation_id': operation_id,
                'service': service_name,
                'rpc_name': rpc_name
            }
            return True  # Success
    
    def get_path_methods(self, path_url):
        """Get all HTTP methods registered for a path"""
        return list(self.all_paths.get(path_url, {}).keys())
    
    def has_collisions(self):
        """Check if any collisions were detected"""
        return len(self.collisions) > 0
    
    def get_collision_count(self):
        """Get total number of collisions"""
        return len(self.collisions)
    
    def get_multi_service_paths(self):
        """Get paths that have contributions from multiple services"""
        multi_service = {}
        for path_url, services in self.service_contributions.items():
            if len(services) > 1:
                multi_service[path_url] = {
                    'services': list(services),
                    'methods': self.get_path_methods(path_url)
                }
        return multi_service
    
    def print_report(self):
        """Print comprehensive validation report"""
        print("\n" + "="*80)
        print("PATH VALIDATION REPORT")
        print("="*80)
        
        # Summary
        total_paths = len(self.all_paths)
        total_operations = sum(len(methods) for methods in self.all_paths.values())
        
        print(f"\n📊 SUMMARY:")
        print(f"   Total unique paths: {total_paths}")
        print(f"   Total operations: {total_operations}")
        print(f"   Collisions detected: {self.get_collision_count()}")
        
        # Multi-service paths
        multi_service = self.get_multi_service_paths()
        if multi_service:
            print(f"\n🔄 MULTI-SERVICE PATHS: {len(multi_service)}")
            for path_url, info in sorted(multi_service.items()):
                print(f"\n   {path_url}")
                print(f"      Services: {', '.join(info['services'])}")
                print(f"      Methods: {', '.join(m.upper() for m in info['methods'])}")
        
        # Collision details
        if self.collisions:
            print(f"\n⚠️  COLLISION DETAILS:")
            for collision in self.collisions:
                print(f"\n   Path: {collision['path']}")
                print(f"   Method: {collision['method'].upper()}")
                print(f"   Conflict:")
                print(f"      ❌ {collision['new_service']}.{collision['new_rpc']} (REJECTED)")
                print(f"      ✅ {collision['existing_service']}.{collision['existing_rpc']} (KEPT)")
        
        # Path coverage analysis
        print(f"\n📈 PATH COVERAGE ANALYSIS:")
        method_distribution = {}
        for path_url, methods in self.all_paths.items():
            count = len(methods)
            method_distribution[count] = method_distribution.get(count, 0) + 1
        
        for count, num_paths in sorted(method_distribution.items()):
            print(f"   Paths with {count} method(s): {num_paths}")
        
        print("\n" + "="*80)
    
    def validate_standard_crud_pattern(self, path_url):
        """
        Validate if a path follows standard CRUD patterns
        
        Standard patterns:
        - Collection: GET (list), POST (create)
        - Resource: GET (retrieve), PATCH/PUT (update), DELETE (delete)
        """
        methods = self.get_path_methods(path_url)
        
        # Check if it's a collection path (no {id})
        is_collection = '{' not in path_url
        
        if is_collection:
            # Collection should have GET and/or POST
            if 'get' not in methods and 'post' not in methods:
                self.warnings.append({
                    'path': path_url,
                    'type': 'unusual_collection',
                    'message': 'Collection path without GET or POST'
                })
        else:
            # Resource should have GET, PATCH/PUT, or DELETE
            has_standard = any(m in methods for m in ['get', 'patch', 'put', 'delete'])
            if not has_standard and 'post' in methods:
                # Custom action paths are OK with POST
                if ':' not in path_url:
                    self.warnings.append({
                        'path': path_url,
                        'type': 'unusual_resource',
                        'message': 'Resource path with only POST (not a custom action)'
                    })
    
    def get_warnings(self):
        """Get all validation warnings"""
        return self.warnings
