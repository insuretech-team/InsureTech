class ProtoRegistry:
    def __init__(self):
        self._type_map = {}  # full_message_name -> openapi_ref
        self._schema_locations = {} # full_message_name -> relative_file_path
        self._collisions = {}  # schema_key -> list of full_names that collide
        self._name_to_full = {}  # simple_name -> full_name (for collision tracking)

    def register_message(self, full_name, file_package, message_name):
        """
        Registers a message and determines its target OpenAPI location.
        Detects and handles name collisions with namespace prefixing.
        
        Args:
            full_name: The full proto name (e.g. insuretech.policy.entity.v1.Policy)
            file_package: The package of the file (e.g. insuretech.policy.entity.v1)
            message_name: The simple message name (e.g. Policy)
        """
        # Determine schema path based on package structure
        package_parts = file_package.split('.')
        relative_path = "/".join(package_parts) + f"/{message_name}.yaml"
        
        # Check for collision
        schema_key = message_name
        
        if schema_key in self._name_to_full:
            # COLLISION DETECTED!
            existing_full_name = self._name_to_full[schema_key]
            
            # Track collision
            if schema_key not in self._collisions:
                self._collisions[schema_key] = [existing_full_name]
            self._collisions[schema_key].append(full_name)
            
            # Apply namespace prefixing
            schema_key = self._create_namespaced_key(full_name, message_name)
            
            # Also update the previously registered one with namespace
            existing_namespaced = self._create_namespaced_key(existing_full_name, message_name)
            self._type_map[existing_full_name] = f"#/components/schemas/{existing_namespaced}"
            
            # Collision auto-resolved via namespace prefix — logged at debug level only.
            # These are expected for insurance_service.proto which aggregates domain types.
            # Only print if namespaced keys ALSO collide (truly unresolvable).
            if existing_namespaced == schema_key:
                print(f"❌ UNRESOLVABLE COLLISION: '{message_name}' → both resolve to '{schema_key}'")
                print(f"    From: {existing_full_name}")
                print(f"    From: {full_name}")
        else:
            # No collision, register normally
            self._name_to_full[schema_key] = full_name
        
        self._type_map[full_name] = f"#/components/schemas/{schema_key}"
        self._schema_locations[full_name] = relative_path

    def get_ref(self, full_name):
        return self._type_map.get(full_name)

    def get_schema_name(self, full_name):
        """
        Get the actual schema name (with collision prefix if any)
        
        Returns schema name string like 'ClaimsDocumentUploadRequest'
        """
        ref = self.get_ref(full_name)
        if ref:
            return ref.split('/')[-1]
        return None

    def get_file_path(self, full_name):
        return self._schema_locations.get(full_name)

    def _create_namespaced_key(self, full_name, message_name):
        """
        Create a namespaced schema key from full proto name
        
        Examples:
            insuretech.claims.services.v1.DocumentUploadRequest → ClaimsDocumentUploadRequest
            insuretech.kyc.services.v1.DocumentUploadRequest → KYCDocumentUploadRequest
            insuretech.policy.services.v1.RenewPolicyRequest → PolicyRenewPolicyRequest
            insuretech.renewal.services.v1.RenewPolicyRequest → RenewalRenewPolicyRequest
        
        Strategy: Use the service/domain name as prefix
        """
        parts = full_name.split('.')
        # parts example: ['insuretech', 'claims', 'services', 'v1', 'DocumentUploadRequest']
        
        if len(parts) >= 3:
            # Get the domain/service name (2nd part after 'insuretech')
            domain = parts[1].capitalize()  # 'claims' -> 'Claims'
            return f"{domain}{message_name}"
        
        # Fallback: use full path
        return '_'.join(parts[1:])  # Skip 'insuretech' prefix
    
    def get_collisions(self):
        """Return dictionary of all detected collisions"""
        return self._collisions
    
    def has_collisions(self):
        """Check if any collisions were detected"""
        return len(self._collisions) > 0
    
    def get_collision_report(self):
        """Generate a human-readable collision report — only shows unresolvable collisions."""
        if not self.has_collisions():
            return "✓ No schema name collisions detected."

        # Separate resolved vs unresolvable
        unresolvable = {}
        resolved_count = 0
        for schema_name, full_names in self._collisions.items():
            namespaced_keys = [self._create_namespaced_key(fn, schema_name) for fn in full_names]
            if len(set(namespaced_keys)) < len(namespaced_keys):
                # Two messages resolve to the same namespaced key — truly unresolvable
                unresolvable[schema_name] = full_names
            else:
                resolved_count += 1

        report = []
        if unresolvable:
            report.append("=" * 80)
            report.append(f"Schema Name Collision Report — {len(unresolvable)} UNRESOLVABLE")
            report.append("=" * 80)
            for schema_name, full_names in sorted(unresolvable.items()):
                report.append(f"\n❌ '{schema_name}' ({len(full_names)} conflicts)")
                for full_name in full_names:
                    namespaced = self._create_namespaced_key(full_name, schema_name)
                    report.append(f"   - {full_name} → {namespaced}")
            report.append("\n" + "=" * 80)
        
        report.append(f"✓ {resolved_count} collision group(s) auto-resolved via namespace prefixing.")
        return "\n".join(report)
    
    def get_relative_ref(self, source_full_name, target_full_name):
        """
        Calculates the relative file path from source message to target message.
        Used for ensuring individual schema files are valid references to each other.
        """
        source_path = self.get_file_path(source_full_name)
        target_path = self.get_file_path(target_full_name)
        
        if not source_path or not target_path:
            return None
            
        # source_path: e.g. insuretech/policy/v1/Policy.yaml
        # target_path: e.g. insuretech/common/v1/Money.yaml
        
        # We need relative path from DIRECTORY of source to target FILE
        source_dir = os.path.dirname(source_path)
        rel_path = os.path.relpath(target_path, source_dir)
        
        # Ensure forward slashes for YAML refs
        return rel_path.replace(os.path.sep, '/')
