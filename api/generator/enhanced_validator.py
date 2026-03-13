#!/usr/bin/env python3
"""
Enhanced OpenAPI Validator with Business Rules and Quality Metrics

Validates OpenAPI spec against:
- Business logic rules
- Field constraints
- Security patterns
- Performance guidelines
- Documentation quality

Usage:
    python enhanced_validator.py ../openapi.yaml --report validation_report.html
"""

import argparse
import yaml
import re
import json
from pathlib import Path
from typing import Dict, List, Tuple, Any
from datetime import datetime


class ValidationIssue:
    """Represents a validation issue"""
    
    SEVERITY_ERROR = 'error'
    SEVERITY_WARNING = 'warning'
    SEVERITY_INFO = 'info'
    
    def __init__(self, severity: str, category: str, location: str, message: str, suggestion: str = None):
        self.severity = severity
        self.category = category
        self.location = location
        self.message = message
        self.suggestion = suggestion
    
    def to_dict(self):
        return {
            'severity': self.severity,
            'category': self.category,
            'location': self.location,
            'message': self.message,
            'suggestion': self.suggestion
        }


class EnhancedValidator:
    """Enhanced validator with business rules and quality checks"""
    
    def __init__(self, spec_path: str):
        self.spec_path = spec_path
        self.issues: List[ValidationIssue] = []
        self.metrics = {
            'total_schemas': 0,
            'total_paths': 0,
            'total_operations': 0,
            'schemas_with_descriptions': 0,
            'schemas_with_examples': 0,
            'operations_with_security': 0,
            'operations_with_rate_limits': 0,
            'description_coverage': 0.0,
            'example_coverage': 0.0,
            'security_coverage': 0.0
        }
        
        # Load spec
        with open(spec_path, 'r', encoding='utf-8') as f:
            self.spec = yaml.safe_load(f)
    
    def validate_all(self) -> Dict[str, Any]:
        """Run all validation checks"""
        print("Running enhanced validation...")
        print()
        
        # Structural validation
        self._validate_structure()
        
        # Business rules
        self._validate_business_rules()
        
        # Security patterns
        self._validate_security()
        
        # Documentation quality
        self._validate_documentation()
        
        # Field constraints
        self._validate_constraints()
        
        # Performance patterns
        self._validate_performance()
        
        # Calculate metrics
        self._calculate_metrics()
        
        return self._generate_report()
    
    def _validate_structure(self):
        """Validate basic OpenAPI structure"""
        print("1. Validating structure...")
        
        # Check required fields
        if 'openapi' not in self.spec:
            self.issues.append(ValidationIssue(
                ValidationIssue.SEVERITY_ERROR,
                'Structure',
                'root',
                'Missing required field: openapi',
                'Add openapi version (e.g., 3.1.0)'
            ))
        
        if 'info' not in self.spec:
            self.issues.append(ValidationIssue(
                ValidationIssue.SEVERITY_ERROR,
                'Structure',
                'root',
                'Missing required field: info'
            ))
        
        if 'paths' not in self.spec:
            self.issues.append(ValidationIssue(
                ValidationIssue.SEVERITY_ERROR,
                'Structure',
                'root',
                'Missing required field: paths'
            ))
        
        print(f"   Found {len(self.spec.get('paths', {}))} paths")
        print(f"   Found {len(self.spec.get('components', {}).get('schemas', {}))} schemas")
    
    def _validate_business_rules(self):
        """Validate business logic rules"""
        print("\n2. Validating business rules...")
        
        schemas = self.spec.get('components', {}).get('schemas', {})
        
        for schema_name, schema in schemas.items():
            # Request DTOs should have required fields
            if schema_name.endswith('Request'):
                if 'required' not in schema or not schema['required']:
                    self.issues.append(ValidationIssue(
                        ValidationIssue.SEVERITY_WARNING,
                        'Business Rules',
                        f'schemas.{schema_name}',
                        'Request DTO has no required fields',
                        'Add required field annotations'
                    ))
            
            # Money/Amount fields should have specific types
            if isinstance(schema, dict) and 'properties' in schema:
                for field_name, field_schema in schema['properties'].items():
                    # Skip if field already uses $ref (correct implementation)
                    if isinstance(field_schema, dict) and '$ref' in field_schema:
                        continue
                    
                    # Skip the Money schema itself (causes false positives)
                    if schema_name == 'Money':
                        continue
                    
                    # Only check fields that are clearly money amounts (not modes, calculations, flags)
                    money_indicators = [
                        '_amount', '_premium', '_price', '_fee', '_cost',
                        'sum_insured', 'sum_assured', 'refundable_'
                    ]
                    
                    # Exclude fields that contain these terms (they're not money values)
                    non_money_terms = [
                        'mode', 'method', 'type', 'calculation', 'breakdown',
                        'consumer', 'adjusted', 'paid', 'used', 'required',
                        'modes'  # premium_payment_modes is an array
                    ]
                    
                    is_likely_money = any(indicator in field_name.lower() for indicator in money_indicators)
                    is_non_money = any(term in field_name.lower() for term in non_money_terms)
                    
                    if is_likely_money and not is_non_money:
                        if isinstance(field_schema, dict) and 'type' in field_schema:
                            # Check if it's a numeric type (should be Money)
                            field_type = field_schema.get('type')
                            if field_type in ['number', 'integer', 'string']:
                                self.issues.append(ValidationIssue(
                                    ValidationIssue.SEVERITY_INFO,
                                    'Business Rules',
                                    f'schemas.{schema_name}.{field_name}',
                                    'Money field should use Money type reference',
                                    'Use $ref to #/components/schemas/Money'
                                ))
        
        print(f"   Checked {len(schemas)} schemas for business rules")
    
    def _validate_security(self):
        """Validate security patterns"""
        print("\n3. Validating security patterns...")
        
        paths = self.spec.get('paths', {})
        
        for path, path_item in paths.items():
            for method, operation in path_item.items():
                if method in ['get', 'post', 'put', 'delete', 'patch']:
                    # Check if operation has security
                    if 'security' not in operation and 'security' not in self.spec:
                        self.issues.append(ValidationIssue(
                            ValidationIssue.SEVERITY_WARNING,
                            'Security',
                            f'{method.upper()} {path}',
                            'Operation has no security requirements',
                            'Add security schemes or mark as public'
                        ))
                    
                    # Check for sensitive data in GET params
                    if method == 'get' and 'parameters' in operation:
                        for param in operation['parameters']:
                            param_name = param.get('name', '').lower()
                            
                            # More precise detection: distinguish between identifiers and actual secrets
                            # SAFE: api_key_id, token_id, secret_id (resource identifiers)
                            # UNSAFE: api_key, token, secret, password (actual secret values)
                            is_identifier = param_name.endswith('_id') or param_name.endswith('-id') or param_name.endswith('id')
                            
                            # Only flag if it's a secret value, not an identifier
                            sensitive_patterns = ['password', 'pwd', 'secret', 'credential']
                            # For 'key' and 'token', only flag if NOT an identifier
                            if not is_identifier:
                                sensitive_patterns.extend(['api_key', 'api-key', 'apikey', 'token', 'access_key', 'private_key'])
                            
                            if any(sensitive in param_name for sensitive in sensitive_patterns):
                                self.issues.append(ValidationIssue(
                                    ValidationIssue.SEVERITY_ERROR,
                                    'Security',
                                    f'{method.upper()} {path}',
                                    f'Sensitive parameter "{param_name}" in GET request',
                                    'Use POST request or header for sensitive data'
                                ))
        
        print(f"   Checked {len(paths)} paths for security patterns")
    
    def _validate_documentation(self):
        """Validate documentation quality"""
        print("\n4. Validating documentation...")
        
        schemas = self.spec.get('components', {}).get('schemas', {})
        
        for schema_name, schema in schemas.items():
            if isinstance(schema, dict):
                # Check for description
                if 'description' not in schema or not schema.get('description', '').strip():
                    self.issues.append(ValidationIssue(
                        ValidationIssue.SEVERITY_INFO,
                        'Documentation',
                        f'schemas.{schema_name}',
                        'Schema has no description',
                        'Add meaningful description'
                    ))
                
                # Check field descriptions
                if 'properties' in schema:
                    for field_name, field_schema in schema['properties'].items():
                        if isinstance(field_schema, dict) and '$ref' not in field_schema:
                            if 'description' not in field_schema:
                                self.issues.append(ValidationIssue(
                                    ValidationIssue.SEVERITY_INFO,
                                    'Documentation',
                                    f'schemas.{schema_name}.{field_name}',
                                    'Field has no description'
                                ))
        
        print(f"   Checked documentation for {len(schemas)} schemas")
    
    def _validate_constraints(self):
        """Validate field constraints"""
        print("\n5. Validating field constraints...")
        
        schemas = self.spec.get('components', {}).get('schemas', {})
        constraint_count = 0
        
        for schema_name, schema in schemas.items():
            if isinstance(schema, dict) and 'properties' in schema:
                for field_name, field_schema in schema['properties'].items():
                    if isinstance(field_schema, dict):
                        # Check string fields for patterns
                        if field_schema.get('type') == 'string':
                            if 'email' in field_name.lower() and 'pattern' not in field_schema and 'format' not in field_schema:
                                self.issues.append(ValidationIssue(
                                    ValidationIssue.SEVERITY_INFO,
                                    'Constraints',
                                    f'schemas.{schema_name}.{field_name}',
                                    'Email field should have format: email',
                                    'Add format: email'
                                ))
                            
                            if 'phone' in field_name.lower() and 'pattern' not in field_schema:
                                self.issues.append(ValidationIssue(
                                    ValidationIssue.SEVERITY_INFO,
                                    'Constraints',
                                    f'schemas.{schema_name}.{field_name}',
                                    'Phone field should have pattern validation'
                                ))
                        
                        # Check if constraints exist
                        if any(k in field_schema for k in ['minLength', 'maxLength', 'minimum', 'maximum', 'pattern']):
                            constraint_count += 1
        
        print(f"   Found {constraint_count} fields with constraints")
    
    def _validate_performance(self):
        """Validate performance patterns"""
        print("\n6. Validating performance patterns...")
        
        paths = self.spec.get('paths', {})
        
        for path, path_item in paths.items():
            for method, operation in path_item.items():
                if method in ['get', 'post']:
                    # Check for pagination on list endpoints
                    if 'list' in operation.get('operationId', '').lower() or path.endswith('s'):
                        has_pagination = False
                        if 'parameters' in operation:
                            param_names = [p.get('name', '') for p in operation['parameters']]
                            if any(p in param_names for p in ['page', 'limit', 'offset', 'pageSize']):
                                has_pagination = True
                        
                        if not has_pagination:
                            self.issues.append(ValidationIssue(
                                ValidationIssue.SEVERITY_WARNING,
                                'Performance',
                                f'{method.upper()} {path}',
                                'List endpoint should support pagination',
                                'Add pagination parameters (page, limit)'
                            ))
        
        print(f"   Checked {len(paths)} paths for performance patterns")
    
    def _calculate_metrics(self):
        """Calculate quality metrics"""
        schemas = self.spec.get('components', {}).get('schemas', {})
        paths = self.spec.get('paths', {})
        
        self.metrics['total_schemas'] = len(schemas)
        self.metrics['total_paths'] = len(paths)
        
        # Count schemas with descriptions
        schemas_with_desc = sum(1 for s in schemas.values() 
                               if isinstance(s, dict) and s.get('description', '').strip())
        self.metrics['schemas_with_descriptions'] = schemas_with_desc
        self.metrics['description_coverage'] = (schemas_with_desc / len(schemas) * 100) if schemas else 0
        
        # Count operations
        operation_count = 0
        operations_with_security = 0
        
        for path_item in paths.values():
            for method in ['get', 'post', 'put', 'delete', 'patch']:
                if method in path_item:
                    operation_count += 1
                    if 'security' in path_item[method]:
                        operations_with_security += 1
        
        self.metrics['total_operations'] = operation_count
        self.metrics['operations_with_security'] = operations_with_security
        self.metrics['security_coverage'] = (operations_with_security / operation_count * 100) if operation_count else 0
    
    def _generate_report(self) -> Dict[str, Any]:
        """Generate validation report"""
        # Group issues by severity
        errors = [i for i in self.issues if i.severity == ValidationIssue.SEVERITY_ERROR]
        warnings = [i for i in self.issues if i.severity == ValidationIssue.SEVERITY_WARNING]
        info = [i for i in self.issues if i.severity == ValidationIssue.SEVERITY_INFO]
        
        return {
            'spec_path': self.spec_path,
            'timestamp': datetime.now().isoformat(),
            'summary': {
                'total_issues': len(self.issues),
                'errors': len(errors),
                'warnings': len(warnings),
                'info': len(info)
            },
            'metrics': self.metrics,
            'issues': {
                'errors': [i.to_dict() for i in errors],
                'warnings': [i.to_dict() for i in warnings],
                'info': [i.to_dict() for i in info]
            }
        }
    
    def print_summary(self, report: Dict[str, Any]):
        """Print validation summary"""
        print()
        print("=" * 80)
        print("VALIDATION SUMMARY")
        print("=" * 80)
        print()
        print(f"Spec: {report['spec_path']}")
        print(f"Timestamp: {report['timestamp']}")
        print()
        print(f"📊 Metrics:")
        print(f"  - Total Schemas: {report['metrics']['total_schemas']}")
        print(f"  - Total Paths: {report['metrics']['total_paths']}")
        print(f"  - Total Operations: {report['metrics']['total_operations']}")
        print(f"  - Description Coverage: {report['metrics']['description_coverage']:.1f}%")
        print(f"  - Security Coverage: {report['metrics']['security_coverage']:.1f}%")
        print()
        print(f"🔍 Issues:")
        print(f"  - ❌ Errors: {report['summary']['errors']}")
        print(f"  - ⚠️  Warnings: {report['summary']['warnings']}")
        print(f"  - ℹ️  Info: {report['summary']['info']}")
        print(f"  - Total: {report['summary']['total_issues']}")
        print()


def main():
    parser = argparse.ArgumentParser(description='Enhanced OpenAPI validator')
    parser.add_argument('spec', help='Path to OpenAPI spec file')
    parser.add_argument('--report', help='Output report file (JSON)', default='validation_report.json')
    parser.add_argument('--html', help='Output HTML report', default=None)
    
    args = parser.parse_args()
    
    # Run validation
    validator = EnhancedValidator(args.spec)
    report = validator.validate_all()
    
    # Print summary
    validator.print_summary(report)
    
    # Save JSON report
    with open(args.report, 'w', encoding='utf-8') as f:
        json.dump(report, f, indent=2)
    print(f"✅ JSON report saved: {args.report}")
    
    # Generate HTML report if requested
    if args.html:
        generate_html_report(report, args.html)
        print(f"✅ HTML report saved: {args.html}")
    
    print()
    
    # Exit with error code if there are errors
    if report['summary']['errors'] > 0:
        print("❌ Validation failed with errors")
        return 1
    else:
        print("✅ Validation passed")
        return 0


def generate_html_report(report: Dict[str, Any], output_path: str):
    """Generate HTML report"""
    html = f"""<!DOCTYPE html>
<html>
<head>
    <title>OpenAPI Validation Report</title>
    <style>
        body {{ font-family: Arial, sans-serif; margin: 20px; background: #f5f5f5; }}
        .container {{ max-width: 1200px; margin: 0 auto; background: white; padding: 30px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }}
        h1 {{ color: #333; }}
        .metrics {{ display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 20px; margin: 20px 0; }}
        .metric {{ background: #f8f9fa; padding: 20px; border-radius: 8px; text-align: center; }}
        .metric-value {{ font-size: 32px; font-weight: bold; color: #007bff; }}
        .metric-label {{ color: #666; margin-top: 10px; }}
        .issue {{ padding: 15px; margin: 10px 0; border-left: 4px solid; border-radius: 4px; }}
        .issue-error {{ background: #fff5f5; border-color: #dc3545; }}
        .issue-warning {{ background: #fff9e6; border-color: #ffc107; }}
        .issue-info {{ background: #e7f3ff; border-color: #17a2b8; }}
        .issue-header {{ font-weight: bold; margin-bottom: 5px; }}
        .issue-location {{ color: #666; font-size: 14px; }}
        .summary {{ display: flex; gap: 20px; margin: 20px 0; }}
        .summary-box {{ flex: 1; padding: 20px; border-radius: 8px; text-align: center; }}
        .summary-errors {{ background: #fff5f5; color: #dc3545; }}
        .summary-warnings {{ background: #fff9e6; color: #856404; }}
        .summary-info {{ background: #e7f3ff; color: #004085; }}
    </style>
</head>
<body>
    <div class="container">
        <h1>🔍 OpenAPI Validation Report</h1>
        <p><strong>Spec:</strong> {report['spec_path']}</p>
        <p><strong>Generated:</strong> {report['timestamp']}</p>
        
        <h2>📊 Metrics</h2>
        <div class="metrics">
            <div class="metric">
                <div class="metric-value">{report['metrics']['total_schemas']}</div>
                <div class="metric-label">Schemas</div>
            </div>
            <div class="metric">
                <div class="metric-value">{report['metrics']['total_paths']}</div>
                <div class="metric-label">Paths</div>
            </div>
            <div class="metric">
                <div class="metric-value">{report['metrics']['total_operations']}</div>
                <div class="metric-label">Operations</div>
            </div>
            <div class="metric">
                <div class="metric-value">{report['metrics']['description_coverage']:.1f}%</div>
                <div class="metric-label">Description Coverage</div>
            </div>
            <div class="metric">
                <div class="metric-value">{report['metrics']['security_coverage']:.1f}%</div>
                <div class="metric-label">Security Coverage</div>
            </div>
        </div>
        
        <h2>🎯 Summary</h2>
        <div class="summary">
            <div class="summary-box summary-errors">
                <h3>{report['summary']['errors']}</h3>
                <p>Errors</p>
            </div>
            <div class="summary-box summary-warnings">
                <h3>{report['summary']['warnings']}</h3>
                <p>Warnings</p>
            </div>
            <div class="summary-box summary-info">
                <h3>{report['summary']['info']}</h3>
                <p>Info</p>
            </div>
        </div>
        
        <h2>❌ Errors ({len(report['issues']['errors'])})</h2>
        {''.join([f'<div class="issue issue-error"><div class="issue-header">{i["message"]}</div><div class="issue-location">{i["location"]}</div><div>{i.get("suggestion", "")}</div></div>' for i in report['issues']['errors']])}
        
        <h2>⚠️ Warnings ({len(report['issues']['warnings'])})</h2>
        {''.join([f'<div class="issue issue-warning"><div class="issue-header">{i["message"]}</div><div class="issue-location">{i["location"]}</div><div>{i.get("suggestion", "")}</div></div>' for i in report['issues']['warnings']])}
        
        <h2>ℹ️ Info ({len(report['issues']['info'])})</h2>
        {''.join([f'<div class="issue issue-info"><div class="issue-header">{i["message"]}</div><div class="issue-location">{i["location"]}</div><div>{i.get("suggestion", "")}</div></div>' for i in report['issues']['info'][:20]])}
        {f'<p>... and {len(report["issues"]["info"]) - 20} more</p>' if len(report['issues']['info']) > 20 else ''}
    </div>
</body>
</html>"""
    
    with open(output_path, 'w', encoding='utf-8') as f:
        f.write(html)


if __name__ == '__main__':
    exit(main())
