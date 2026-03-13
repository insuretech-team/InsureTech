"""
FIXED DTO Name Transformer
Converts verb-based proto message names to noun-based OpenAPI schema names
Following apirules.md: DTOs should use nouns, not verbs

FIXES APPLIED:
1. Preserved entity names (RenewalReminder, RenewalSchedule) - no transformations
2. Added safeguards to prevent string corruption
3. Explicit preservation list for known entities
"""

import re
from typing import Dict, Optional


class NameTransformer:
    """Transforms proto message names to follow API naming conventions"""
    
    # Verb to noun mappings
    VERB_TO_NOUN = {
        'Create': 'Creation',
        'Update': 'Update',  # Already a noun
        'Delete': 'Deletion',
        'Get': 'Retrieval',
        'List': 'Listing',
        'Fetch': 'Fetch',
        'Retrieve': 'Retrieval',
        'Cancel': 'Cancellation',
        'Renew': 'Renewal',  # For RPC methods like RenewPolicy
        'Generate': 'Generation',
        'Issue': 'Issuance',
        'Process': 'Processing',
        'Submit': 'Submission',
        'Approve': 'Approval',
        'Reject': 'Rejection',
        'Verify': 'Verification',
        'Validate': 'Validation',
        'Send': 'Sending',
        'Receive': 'Receipt',
        'Calculate': 'Calculation',
        'Execute': 'Execution',
        'Perform': 'Performance',
        'Handle': 'Handling',
        'Manage': 'Management',
        'Search': 'Search',
        'Query': 'Query',
        'Filter': 'Filter',
        'Sort': 'Sorting',
        'Export': 'Export',
        'Import': 'Import',
        'Upload': 'Upload',
        'Download': 'Download',
        'Analyze': 'Analysis',
        'Evaluate': 'Evaluation',
        'Assess': 'Assessment',
        'Review': 'Review',
        'Confirm': 'Confirmation',
        'Acknowledge': 'Acknowledgment',
        'Notify': 'Notification',
        'Alert': 'Alert',
        'Trigger': 'Trigger',
        'Activate': 'Activation',
        'Deactivate': 'Deactivation',
        'Enable': 'Enablement',
        'Disable': 'Disablement',
        'Suspend': 'Suspension',
        'Resume': 'Resumption',
        'Pause': 'Pause',
        'Start': 'Start',
        'Stop': 'Stop',
        'Restart': 'Restart',
        'Terminate': 'Termination',
        'Complete': 'Completion',
        'Finalize': 'Finalization',
        'Initialize': 'Initialization',
        'Register': 'Registration',
        'Unregister': 'Deregistration',
        'Enroll': 'Enrollment',
        'Unenroll': 'Unenrollment',
        'Assign': 'Assignment',
        'Unassign': 'Unassignment',
        'Attach': 'Attachment',
        'Detach': 'Detachment',
        'Rotate': 'Rotation',
    }
    
    # CRITICAL FIX: Patterns that should remain unchanged (entities, not DTOs)
    PRESERVE_PATTERNS = [
        r'.*Info$',          # UserInfo, PolicyInfo
        r'.*Data$',          # ClaimData, PaymentData
        r'.*Details$',       # OrderDetails
        r'.*Status$',        # PolicyStatus
        r'.*Type$',          # ProductType
        r'.*State$',         # ClaimState
        r'.*Config$',        # SystemConfig
        r'.*Settings$',      # UserSettings
        r'.*Options$',       # SearchOptions
        r'.*Params$',        # QueryParams
        r'.*Reminder$',      # RenewalReminder - FIX: preserve entity
        r'.*Schedule$',      # RenewalSchedule - FIX: preserve entity
        r'.*Definition$',    # ReportDefinition - intentional naming
        r'.*Execution$',     # ReportExecution - intentional naming
        r'.*Metric$',        # AggregatedMetric - preserve
        r'.*Log$',           # AuditLog - preserve
        r'.*Event$',         # All events should preserve names
        r'.*Period$',        # GracePeriod - preserve
    ]
    
    # EXPLICIT PRESERVATION: Known entity names that should NEVER be transformed
    PRESERVE_EXACT = {
        'RenewalReminder',
        'RenewalSchedule', 
        'GracePeriod',
        'ReportDefinition',
        'ReportExecution',
        'ReportSchedule',
        'AggregatedMetric',
        'MetricDefinition',
        'AuditLog',
        'AuditEvent',
        'ComplianceLog',
        'BusinessMetrics',
        'Dashboard',
        'Report',
    }
    
    def __init__(self, preserve_request_response: bool = True):
        """
        Initialize transformer
        
        Args:
            preserve_request_response: If True, keep Request/Response suffix
        """
        self.preserve_request_response = preserve_request_response
    
    def transform(self, proto_name: str) -> str:
        """
        Transform a proto message name to API-compliant name
        
        Args:
            proto_name: Original proto message name (e.g., CreatePolicyRequest)
            
        Returns:
            Transformed name (e.g., PolicyCreationRequest)
        """
        # CRITICAL FIX: Check exact match preservation first
        base_without_suffix = self._remove_suffix(proto_name)
        if base_without_suffix in self.PRESERVE_EXACT:
            return proto_name  # Return as-is, no transformation
        
        # Check if should preserve as-is based on patterns
        for pattern in self.PRESERVE_PATTERNS:
            if re.match(pattern, proto_name):
                return proto_name
        
        # Extract suffix (Request, Response, Event, etc.)
        suffix = None
        base_name = proto_name
        
        for possible_suffix in ['Request', 'Response', 'Event', 'Message', 'DTO']:
            if proto_name.endswith(possible_suffix):
                suffix = possible_suffix
                base_name = proto_name[:-len(possible_suffix)]
                break
        
        # Transform the base name
        transformed_base = self._transform_base_name(base_name)
        
        # Reconstruct with suffix if needed
        if suffix and self.preserve_request_response:
            return transformed_base + suffix
        elif suffix:
            return transformed_base
        else:
            return transformed_base
    
    def _remove_suffix(self, name: str) -> str:
        """Remove Request/Response/Event suffix to get base name"""
        for suffix in ['Request', 'Response', 'Event', 'Message', 'DTO']:
            if name.endswith(suffix):
                return name[:-len(suffix)]
        return name
    
    def _transform_base_name(self, name: str) -> str:
        """
        Transform the base name (without Request/Response suffix)
        
        Examples:
            CreatePolicy → PolicyCreation
            CancelPolicy → PolicyCancellation
            GetPolicyDetails → PolicyDetailsRetrieval
            ListUserPolicies → UserPolicyListing
            
        NOT TRANSFORMED (entities):
            RenewalReminder → RenewalReminder (preserved)
            RenewalSchedule → RenewalSchedule (preserved)
        """
        # SAFETY CHECK: If already in preserve list, return as-is
        if name in self.PRESERVE_EXACT:
            return name
        
        # Pattern 1: VerbNoun (e.g., CreatePolicy)
        for verb, noun in self.VERB_TO_NOUN.items():
            if name.startswith(verb):
                remainder = name[len(verb):]
                if remainder:
                    # CreatePolicy → PolicyCreation
                    # But NOT: RenewPolicy → PolicyRenewal (this is a DTO, OK to transform)
                    return remainder + noun
                else:
                    # Just the verb (rare)
                    return noun
        
        # Pattern 2: VerbAdjectiveNoun (e.g., GetPolicyDetails)
        # Try to find verb at start
        words = self._split_pascal_case(name)
        if words and words[0] in self.VERB_TO_NOUN:
            verb = words[0]
            noun = self.VERB_TO_NOUN[verb]
            remaining_words = words[1:]
            if remaining_words:
                # GetPolicyDetails → PolicyDetailsRetrieval
                return ''.join(remaining_words) + noun
            else:
                return noun
        
        # No transformation needed
        return name
    
    def _split_pascal_case(self, name: str) -> list:
        """Split PascalCase into words"""
        # Insert space before uppercase letters
        spaced = re.sub(r'([A-Z])', r' \1', name).strip()
        return spaced.split()
    
    def _to_kebab_case(self, text: str) -> str:
        """Convert PascalCase/camelCase to kebab-case"""
        s1 = re.sub('(.)([A-Z][a-z]+)', r'\1-\2', text)
        return re.sub('([a-z0-9])([A-Z])', r'\1-\2', s1).lower()
    
    def transform_bulk(self, names: list) -> Dict[str, str]:
        """
        Transform multiple names at once
        
        Args:
            names: List of proto names
            
        Returns:
            Dictionary mapping original to transformed names
        """
        return {name: self.transform(name) for name in names}
    
    def get_transformation_report(self, proto_names: list) -> str:
        """
        Generate a report of all transformations
        
        Args:
            proto_names: List of proto message names
            
        Returns:
            Formatted report string
        """
        transformations = self.transform_bulk(proto_names)
        
        report = []
        report.append("=" * 80)
        report.append("DTO Name Transformation Report")
        report.append("=" * 80)
        report.append(f"\nTotal Messages: {len(proto_names)}")
        
        changed = {k: v for k, v in transformations.items() if k != v}
        unchanged = {k: v for k, v in transformations.items() if k == v}
        
        report.append(f"Transformed: {len(changed)}")
        report.append(f"Unchanged: {len(unchanged)}")
        
        if changed:
            report.append("\n" + "=" * 80)
            report.append("Transformations")
            report.append("=" * 80)
            
            for original, transformed in sorted(changed.items()):
                report.append(f"\n  {original}")
                report.append(f"  └─> {transformed}")
        
        if unchanged:
            report.append("\n" + "=" * 80)
            report.append("Unchanged (Already Compliant or Preserved)")
            report.append("=" * 80)
            
            for name in sorted(unchanged.keys()):
                if name in self.PRESERVE_EXACT or self._remove_suffix(name) in self.PRESERVE_EXACT:
                    report.append(f"  🔒 {name} (entity - preserved)")
                else:
                    report.append(f"  ✓ {name}")
        
        report.append("\n" + "=" * 80)
        
        return "\n".join(report)


def main():
    """CLI entry point for testing transformations"""
    import argparse
    
    parser = argparse.ArgumentParser(
        description='Transform proto DTO names to API-compliant names'
    )
    parser.add_argument(
        'names',
        nargs='+',
        help='Proto message names to transform'
    )
    parser.add_argument(
        '--no-suffix',
        action='store_true',
        help='Remove Request/Response suffix'
    )
    
    args = parser.parse_args()
    
    transformer = NameTransformer(
        preserve_request_response=not args.no_suffix
    )
    
    print("\nDTO Name Transformations:")
    print("=" * 60)
    
    for name in args.names:
        transformed = transformer.transform(name)
        
        if name != transformed:
            print(f"\n{name}")
            print(f"  → {transformed}")
        else:
            base = transformer._remove_suffix(name)
            if base in transformer.PRESERVE_EXACT:
                print(f"\n🔒 {name} (entity - preserved)")
            else:
                print(f"\n✓ {name} (no change needed)")
    
    print("\n" + "=" * 60)


if __name__ == '__main__':
    main()
