"""
OpenAPI Spec Validator for InsureTech API
Validates generated OpenAPI specs against apirules.md guidelines
"""

import yaml
import re
import os
from typing import List, Dict, Any, Tuple
from dataclasses import dataclass
from enum import Enum


class Severity(Enum):
    ERROR = "ERROR"
    WARNING = "WARNING"
    INFO = "INFO"


@dataclass
class ValidationIssue:
    severity: Severity
    category: str
    message: str
    location: str
    rule: str
    suggestion: str = ""


class APIRulesValidator:
    """Validates OpenAPI spec against InsureTech API design rules"""
    
    # Verb patterns to detect in names
    VERBS = [
        'Create', 'Update', 'Delete', 'Get', 'List', 'Fetch', 'Retrieve',
        'Cancel', 'Renew', 'Generate', 'Issue', 'Process', 'Submit',
        'Approve', 'Reject', 'Verify', 'Validate', 'Send', 'Receive',
        'Calculate', 'Execute', 'Perform', 'Handle', 'Manage'
    ]
    
    # Verb to noun mappings
    VERB_TO_NOUN = {
        'Create': 'Creation',
        'Update': 'Update',
        'Delete': 'Deletion',
        'Cancel': 'Cancellation',
        'Renew': 'Renewal',
        'Generate': 'Generation',
        'Issue': 'Issuance',
        'Process': 'Processing',
        'Submit': 'Submission',
        'Approve': 'Approval',
        'Reject': 'Rejection',
        'Verify': 'Verification',
        'Validate': 'Validation',
        'Calculate': 'Calculation',
        'Execute': 'Execution',
    }
    
    def __init__(self, spec_path: str):
        """Initialize validator with OpenAPI spec path"""
        self.spec_path = spec_path
        self.spec = None
        self.issues: List[ValidationIssue] = []
    
    def load_spec(self) -> bool:
        """Load OpenAPI spec from file"""
        try:
            with open(self.spec_path, 'r', encoding='utf-8') as f:
                self.spec = yaml.safe_load(f)
            return True
        except Exception as e:
            print(f"Error loading spec: {e}")
            return False
    
    def validate_all(self) -> List[ValidationIssue]:
        """Run all validation checks"""
        if not self.spec:
            if not self.load_spec():
                return self.issues
        
        self.issues = []
        
        # Run all validation checks
        self.validate_naming_conventions()
        self.validate_dto_names()
        self.validate_url_patterns()
        self.validate_action_parameters()
        self.validate_http_methods()
        self.validate_response_codes()
        self.validate_descriptions()
        self.validate_required_fields()
        self.validate_security()
        self.validate_error_responses()
        
        return self.issues
    
    def validate_naming_conventions(self):
        """Check kebab-case in URLs"""
        paths = self.spec.get('paths', {})
        
        for path, path_item in paths.items():
            # Check path uses kebab-case (allow {params})
            path_segments = re.findall(r'[^/{}\?]+', path)
            
            for segment in path_segments:
                if segment.startswith('v') and segment[1:].isdigit():
                    continue  # Skip version segments like v1
                
                # Check for camelCase or PascalCase
                if re.search(r'[a-z][A-Z]', segment):
                    suggestion = self._to_kebab_case(segment)
                    self.issues.append(ValidationIssue(
                        severity=Severity.ERROR,
                        category="Naming Convention",
                        message=f"Path segment '{segment}' uses camelCase, should use kebab-case",
                        location=path,
                        rule="apirules.md: ✅ Kebab-case naming convention",
                        suggestion=f"Use '{suggestion}' instead"
                    ))
    
    def validate_dto_names(self):
        """Check DTO names don't contain verbs"""
        schemas = self.spec.get('components', {}).get('schemas', {})
        
        for schema_name, schema_def in schemas.items():
            # Check if name contains a verb
            for verb in self.VERBS:
                if verb in schema_name:
                    noun = self.VERB_TO_NOUN.get(verb, verb)
                    suggestion = schema_name.replace(verb, noun)
                    
                    self.issues.append(ValidationIssue(
                        severity=Severity.ERROR,
                        category="DTO Naming",
                        message=f"DTO name '{schema_name}' contains verb '{verb}'",
                        location=f"components.schemas.{schema_name}",
                        rule="apirules.md: URIs = Nouns only (applies to DTOs)",
                        suggestion=f"Rename to '{suggestion}'"
                    ))
                    break
    
    def validate_url_patterns(self):
        """Check URLs follow REST conventions - nouns only"""
        paths = self.spec.get('paths', {})
        
        for path in paths.keys():
            # Extract path segments (exclude params and actions)
            segments = re.findall(r'/([^/{\?:]+)', path)
            
            for segment in segments:
                if segment.startswith('v') and segment[1:].isdigit():
                    continue  # Skip version
                
                # Check if segment is a verb
                segment_words = re.findall(r'[A-Z]?[a-z]+', segment)
                for word in segment_words:
                    if word.capitalize() in self.VERBS:
                        self.issues.append(ValidationIssue(
                            severity=Severity.ERROR,
                            category="URL Pattern",
                            message=f"URL segment '{segment}' contains verb '{word}'",
                            location=path,
                            rule="apirules.md: URIs = Nouns only (no verbs)",
                            suggestion=f"Use resource nouns instead of verbs"
                        ))
    
    def validate_action_parameters(self):
        """Check custom actions use query parameters, not URL syntax"""
        paths = self.spec.get('paths', {})
        
        for path, path_item in paths.items():
            # Check for custom action syntax like :cancel, :renew
            if ':' in path and not path.startswith('http'):
                # Extract action
                action_match = re.search(r':([a-zA-Z-]+)', path)
                if action_match:
                    action = action_match.group(1)
                    base_path = path.split(':')[0]
                    
                    self.issues.append(ValidationIssue(
                        severity=Severity.ERROR,
                        category="Action Pattern",
                        message=f"Custom action '{action}' uses URL syntax instead of query parameter",
                        location=path,
                        rule="apirules.md: Actions = Query params (?action=cancel)",
                        suggestion=f"Use '{base_path}?action={self._to_kebab_case(action)}' instead"
                    ))
    
    def validate_http_methods(self):
        """Validate proper use of HTTP methods"""
        paths = self.spec.get('paths', {})
        
        method_rules = {
            'get': 'Read operations only, should be idempotent and safe',
            'post': 'Create or trigger actions',
            'patch': 'Partial updates',
            'put': 'Full replacement',
            'delete': 'Remove resources'
        }
        
        for path, path_item in paths.items():
            for method, operation in path_item.items():
                if method not in method_rules:
                    continue
                
                # Check POST on collection vs item
                if method == 'post' and not re.search(r'\{[^}]+\}', path):
                    # POST on collection - should create
                    summary = operation.get('summary', '').lower()
                    if 'get' in summary or 'list' in summary or 'fetch' in summary:
                        self.issues.append(ValidationIssue(
                            severity=Severity.WARNING,
                            category="HTTP Method",
                            message=f"POST method with read operation summary",
                            location=f"{path}.{method}",
                            rule="POST should create resources or trigger actions",
                            suggestion="Use GET for read operations"
                        ))
    
    def validate_response_codes(self):
        """Validate appropriate HTTP status codes"""
        paths = self.spec.get('paths', {})
        
        for path, path_item in paths.items():
            for method, operation in path_item.items():
                if method not in ['get', 'post', 'patch', 'put', 'delete']:
                    continue
                
                responses = operation.get('responses', {})
                
                # Check for appropriate success codes
                if method == 'post':
                    if '201' not in responses and '200' in responses:
                        self.issues.append(ValidationIssue(
                            severity=Severity.WARNING,
                            category="Response Code",
                            message=f"POST operation returns 200, should return 201 for creation",
                            location=f"{path}.{method}.responses",
                            rule="POST creation should return 201 Created",
                            suggestion="Use 201 status code with Location header"
                        ))
                
                if method == 'delete':
                    if '204' not in responses and '200' in responses:
                        self.issues.append(ValidationIssue(
                            severity=Severity.INFO,
                            category="Response Code",
                            message=f"DELETE operation returns 200, consider 204 No Content",
                            location=f"{path}.{method}.responses",
                            rule="DELETE typically returns 204 No Content",
                            suggestion="Use 204 for successful deletion with no body"
                        ))
                
                # Check for error responses
                required_errors = ['400', '401', '403', '500']
                for error_code in required_errors:
                    if error_code not in responses:
                        self.issues.append(ValidationIssue(
                            severity=Severity.WARNING,
                            category="Error Response",
                            message=f"Missing {error_code} error response",
                            location=f"{path}.{method}.responses",
                            rule="All operations should document error responses",
                            suggestion=f"Add {error_code} response with Error schema"
                        ))
    
    def validate_descriptions(self):
        """Check for missing or generic descriptions"""
        schemas = self.spec.get('components', {}).get('schemas', {})
        paths = self.spec.get('paths', {})
        
        # Check schema descriptions
        for schema_name, schema_def in schemas.items():
            desc = schema_def.get('description', '').strip()
            if not desc:
                self.issues.append(ValidationIssue(
                    severity=Severity.WARNING,
                    category="Documentation",
                    message=f"Schema '{schema_name}' has no description",
                    location=f"components.schemas.{schema_name}",
                    rule="All schemas should have descriptions",
                    suggestion="Add description from api/descriptions/ or proto comments"
                ))
            
            # Check property descriptions
            properties = schema_def.get('properties', {})
            for prop_name, prop_def in properties.items():
                if not prop_def.get('description'):
                    self.issues.append(ValidationIssue(
                        severity=Severity.INFO,
                        category="Documentation",
                        message=f"Property '{prop_name}' in '{schema_name}' has no description",
                        location=f"components.schemas.{schema_name}.properties.{prop_name}",
                        rule="Properties should have descriptions",
                        suggestion="Document field purpose and constraints"
                    ))
        
        # Check operation descriptions
        for path, path_item in paths.items():
            for method, operation in path_item.items():
                if method not in ['get', 'post', 'patch', 'put', 'delete']:
                    continue
                
                summary = operation.get('summary', '').strip()
                description = operation.get('description', '').strip()
                
                if not summary and not description:
                    self.issues.append(ValidationIssue(
                        severity=Severity.ERROR,
                        category="Documentation",
                        message=f"Operation has no summary or description",
                        location=f"{path}.{method}",
                        rule="All operations must be documented",
                        suggestion="Add summary and description from api/descriptions/"
                    ))
                elif summary and len(summary.split()) <= 2:
                    self.issues.append(ValidationIssue(
                        severity=Severity.WARNING,
                        category="Documentation",
                        message=f"Operation summary is too brief: '{summary}'",
                        location=f"{path}.{method}",
                        rule="Summaries should be descriptive",
                        suggestion="Provide more detailed summary"
                    ))
    
    def validate_required_fields(self):
        """Check that required fields are marked"""
        schemas = self.spec.get('components', {}).get('schemas', {})
        
        for schema_name, schema_def in schemas.items():
            if schema_def.get('type') != 'object':
                continue
            
            properties = schema_def.get('properties', {})
            required = schema_def.get('required', [])
            
            # Request DTOs should have required fields
            if 'Request' in schema_name and not required and properties:
                self.issues.append(ValidationIssue(
                    severity=Severity.WARNING,
                    category="Schema Validation",
                    message=f"Request DTO '{schema_name}' has no required fields",
                    location=f"components.schemas.{schema_name}",
                    rule="Request DTOs should mark required fields",
                    suggestion="Add 'required' array based on proto field_behavior"
                ))
    
    def validate_security(self):
        """Check security schemes are properly configured"""
        security_schemes = self.spec.get('components', {}).get('securitySchemes', {})
        global_security = self.spec.get('security', [])
        
        if not security_schemes:
            self.issues.append(ValidationIssue(
                severity=Severity.ERROR,
                category="Security",
                message="No security schemes defined",
                location="components.securitySchemes",
                rule="apirules.md: Security by default",
                suggestion="Define BearerAuth and ApiKeyAuth schemes"
            ))
        
        if not global_security:
            self.issues.append(ValidationIssue(
                severity=Severity.WARNING,
                category="Security",
                message="No global security requirements",
                location="security",
                rule="All endpoints should require authentication by default",
                suggestion="Add global security: [BearerAuth: []]"
            ))
    
    def validate_error_responses(self):
        """Validate error responses use standard Error schema"""
        paths = self.spec.get('paths', {})
        
        for path, path_item in paths.items():
            for method, operation in path_item.items():
                if method not in ['get', 'post', 'patch', 'put', 'delete']:
                    continue
                
                responses = operation.get('responses', {})
                
                # Check 4xx and 5xx responses
                for code, response_def in responses.items():
                    if code.startswith('4') or code.startswith('5'):
                        content = response_def.get('content', {})
                        json_content = content.get('application/json', {})
                        schema = json_content.get('schema', {})
                        ref = schema.get('$ref', '')
                        
                        if 'Error' not in ref:
                            self.issues.append(ValidationIssue(
                                severity=Severity.WARNING,
                                category="Error Response",
                                message=f"Error response {code} doesn't use Error schema",
                                location=f"{path}.{method}.responses.{code}",
                                rule="Error responses should use standard Error schema",
                                suggestion="Use $ref: '#/components/schemas/Error'"
                            ))
    
    def _to_kebab_case(self, text: str) -> str:
        """Convert camelCase or PascalCase to kebab-case"""
        # Insert hyphens before uppercase letters
        s1 = re.sub('(.)([A-Z][a-z]+)', r'\1-\2', text)
        return re.sub('([a-z0-9])([A-Z])', r'\1-\2', s1).lower()
    
    def generate_report(self, output_path: str = None) -> str:
        """Generate validation report"""
        if not self.issues:
            return "✅ No issues found. OpenAPI spec is compliant with API rules."
        
        # Group by severity
        errors = [i for i in self.issues if i.severity == Severity.ERROR]
        warnings = [i for i in self.issues if i.severity == Severity.WARNING]
        info = [i for i in self.issues if i.severity == Severity.INFO]
        
        report = []
        report.append("=" * 80)
        report.append("OpenAPI Spec Validation Report")
        report.append("=" * 80)
        report.append(f"\nSpec: {self.spec_path}")
        report.append(f"\nSummary:")
        report.append(f"  ❌ Errors: {len(errors)}")
        report.append(f"  ⚠️  Warnings: {len(warnings)}")
        report.append(f"  ℹ️  Info: {len(info)}")
        report.append(f"\nTotal Issues: {len(self.issues)}")
        report.append("\n" + "=" * 80)
        
        # Report errors
        if errors:
            report.append("\n\n❌ ERRORS (Must Fix)\n")
            report.append("-" * 80)
            for i, issue in enumerate(errors, 1):
                report.append(f"\n{i}. [{issue.category}] {issue.message}")
                report.append(f"   Location: {issue.location}")
                report.append(f"   Rule: {issue.rule}")
                if issue.suggestion:
                    report.append(f"   💡 Suggestion: {issue.suggestion}")
        
        # Report warnings
        if warnings:
            report.append("\n\n⚠️  WARNINGS (Should Fix)\n")
            report.append("-" * 80)
            for i, issue in enumerate(warnings, 1):
                report.append(f"\n{i}. [{issue.category}] {issue.message}")
                report.append(f"   Location: {issue.location}")
                report.append(f"   Rule: {issue.rule}")
                if issue.suggestion:
                    report.append(f"   💡 Suggestion: {issue.suggestion}")
        
        # Report info
        if info:
            report.append("\n\nℹ️  INFORMATION (Consider Fixing)\n")
            report.append("-" * 80)
            for i, issue in enumerate(info, 1):
                report.append(f"\n{i}. [{issue.category}] {issue.message}")
                report.append(f"   Location: {issue.location}")
                if issue.suggestion:
                    report.append(f"   💡 Suggestion: {issue.suggestion}")
        
        report.append("\n\n" + "=" * 80)
        report.append("End of Report")
        report.append("=" * 80)
        
        report_text = "\n".join(report)
        
        if output_path:
            with open(output_path, 'w', encoding='utf-8') as f:
                f.write(report_text)
            print(f"Report saved to: {output_path}")
        
        return report_text


def main():
    """CLI entry point"""
    import argparse
    
    parser = argparse.ArgumentParser(
        description='Validate OpenAPI spec against InsureTech API rules'
    )
    parser.add_argument(
        'spec_path',
        help='Path to OpenAPI YAML file'
    )
    parser.add_argument(
        '-o', '--output',
        help='Output report file path',
        default=None
    )
    parser.add_argument(
        '-v', '--verbose',
        action='store_true',
        help='Verbose output'
    )
    
    args = parser.parse_args()
    
    validator = APIRulesValidator(args.spec_path)
    
    print(f"Validating: {args.spec_path}")
    print("Loading spec...")
    
    if not validator.load_spec():
        print("❌ Failed to load spec")
        return 1
    
    print("Running validation checks...")
    issues = validator.validate_all()
    
    print(f"\n✅ Validation complete. Found {len(issues)} issues.\n")
    
    report = validator.generate_report(args.output)
    
    if args.verbose or not args.output:
        print(report)
    
    # Exit code based on errors
    errors = [i for i in issues if i.severity == Severity.ERROR]
    return 1 if errors else 0


if __name__ == '__main__':
    exit(main())
