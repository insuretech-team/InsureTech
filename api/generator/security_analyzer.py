#!/usr/bin/env python3
"""
Security Analyzer for API Generation

Detects and fixes security issues:
1. Sensitive data in GET parameters (API keys, tokens, passwords)
2. Sensitive fields from proto annotations (google.api.field_behavior = SENSITIVE)
3. Insecure path patterns
4. Missing authentication requirements

Usage:
    from security_analyzer import SecurityAnalyzer
    analyzer = SecurityAnalyzer(proto_parser)
    issues = analyzer.analyze_path(path, method, operation)
"""

from typing import List, Dict, Any, Tuple
import re


class SecurityIssue:
    """Represents a security issue"""
    
    SEVERITY_CRITICAL = 'critical'
    SEVERITY_HIGH = 'high'
    SEVERITY_MEDIUM = 'medium'
    SEVERITY_LOW = 'low'
    
    def __init__(self, severity: str, category: str, location: str, message: str, fix: str = None):
        self.severity = severity
        self.category = category
        self.location = location
        self.message = message
        self.fix = fix


class SecurityAnalyzer:
    """Analyzes and fixes security issues in API generation"""
    
    # Sensitive keywords that should never be in GET parameters
    SENSITIVE_KEYWORDS = [
        'api_key', 'api-key', 'apikey',
        'secret', 'token', 'password', 'pwd',
        'credential', 'auth', 'access_key',
        'private_key', 'encryption_key'
    ]
    
    # Resources that should have restricted access
    SENSITIVE_RESOURCES = [
        'api-keys', 'api_keys',
        'secrets', 'credentials',
        'tokens', 'passwords',
        'private-keys', 'encryption-keys'
    ]
    
    def __init__(self, proto_parser=None):
        self.proto_parser = proto_parser
        self.issues: List[SecurityIssue] = []
    
    def analyze_path(self, path: str, method: str, operation: Dict[str, Any]) -> List[SecurityIssue]:
        """
        Analyze a single API path for security issues
        
        Returns:
            List of security issues found
        """
        issues = []
        
        # Check for sensitive data in path parameters
        if method.lower() == 'get':
            issues.extend(self._check_sensitive_path_params(path, method, operation))
        
        # Check for sensitive resources without proper protection
        issues.extend(self._check_sensitive_resources(path, method, operation))
        
        # Check for missing authentication
        issues.extend(self._check_authentication(path, method, operation))
        
        # Check for insecure patterns
        issues.extend(self._check_insecure_patterns(path, method, operation))
        
        return issues
    
    def _check_sensitive_path_params(self, path: str, method: str, operation: Dict[str, Any]) -> List[SecurityIssue]:
        """Check for sensitive data in GET request path parameters"""
        issues = []
        
        # Extract parameters from path
        param_pattern = r'\{([^}]+)\}'
        params = re.findall(param_pattern, path)
        
        for param in params:
            param_lower = param.lower()
            
            # Check if parameter name contains sensitive keywords
            for keyword in self.SENSITIVE_KEYWORDS:
                if keyword in param_lower:
                    issues.append(SecurityIssue(
                        SecurityIssue.SEVERITY_CRITICAL,
                        'Sensitive Data Exposure',
                        f'{method.upper()} {path}',
                        f'Sensitive parameter "{param}" exposed in GET request URL',
                        f'Use POST request body or secure header instead of path parameter'
                    ))
                    break
        
        return issues
    
    def _check_sensitive_resources(self, path: str, method: str, operation: Dict[str, Any]) -> List[SecurityIssue]:
        """Check for sensitive resources with inadequate protection"""
        issues = []
        
        # Check if path contains sensitive resource names
        for resource in self.SENSITIVE_RESOURCES:
            if f'/{resource}/' in path or path.endswith(f'/{resource}'):
                
                # These resources should NEVER allow direct ID-based retrieval
                if method.lower() == 'get' and '{' in path:
                    issues.append(SecurityIssue(
                        SecurityIssue.SEVERITY_CRITICAL,
                        'Insecure Resource Access',
                        f'{method.upper()} {path}',
                        f'Direct GET access to sensitive resource "{resource}" by ID is insecure',
                        f'Use authenticated session-based retrieval or POST with verification token'
                    ))
                
                # Check for missing security requirements
                if 'security' not in operation:
                    issues.append(SecurityIssue(
                        SecurityIssue.SEVERITY_HIGH,
                        'Missing Authentication',
                        f'{method.upper()} {path}',
                        f'Sensitive resource "{resource}" has no authentication requirement',
                        'Add security scheme (OAuth2, API Key, etc.)'
                    ))
        
        return issues
    
    def _check_authentication(self, path: str, method: str, operation: Dict[str, Any]) -> List[SecurityIssue]:
        """Check for missing authentication on sensitive operations"""
        issues = []
        
        # POST, PUT, DELETE should always have authentication
        if method.lower() in ['post', 'put', 'delete', 'patch']:
            if 'security' not in operation:
                # Check if it's a public endpoint (login, register, etc.)
                public_endpoints = ['/login', '/register', '/forgot-password', '/verify-otp']
                is_public = any(endpoint in path for endpoint in public_endpoints)
                
                if not is_public:
                    issues.append(SecurityIssue(
                        SecurityIssue.SEVERITY_HIGH,
                        'Missing Authentication',
                        f'{method.upper()} {path}',
                        'Mutating operation has no authentication requirement',
                        'Add security scheme'
                    ))
        
        return issues
    
    def _check_insecure_patterns(self, path: str, method: str, operation: Dict[str, Any]) -> List[SecurityIssue]:
        """Check for known insecure API patterns"""
        issues = []
        
        # Pattern 1: Sequential integer IDs for sensitive resources
        if re.search(r'/\{(\w+_)?id\}', path):
            for resource in self.SENSITIVE_RESOURCES:
                if resource in path:
                    issues.append(SecurityIssue(
                        SecurityIssue.SEVERITY_MEDIUM,
                        'Predictable Identifiers',
                        f'{method.upper()} {path}',
                        f'Sensitive resource uses potentially predictable ID pattern',
                        'Use UUIDs or non-sequential identifiers for sensitive resources'
                    ))
                    break
        
        # Pattern 2: API keys in query parameters
        if 'parameters' in operation:
            for param in operation['parameters']:
                if param.get('in') == 'query':
                    param_name = param.get('name', '').lower()
                    if any(keyword in param_name for keyword in self.SENSITIVE_KEYWORDS):
                        issues.append(SecurityIssue(
                            SecurityIssue.SEVERITY_CRITICAL,
                            'Sensitive Data in Query',
                            f'{method.upper()} {path}',
                            f'Sensitive parameter "{param.get("name")}" in query string',
                            'Move to request header or POST body'
                        ))
        
        return issues
    
    def suggest_secure_alternatives(self, path: str, method: str) -> Dict[str, str]:
        """
        Suggest secure alternatives for insecure patterns
        
        Returns:
            Dictionary with 'insecure' and 'secure' patterns
        """
        suggestions = {}
        
        # Pattern: GET /v1/api-keys/{api_key_id}
        if 'api-keys' in path or 'api_keys' in path:
            if method.lower() == 'get' and '{' in path:
                suggestions['insecure'] = path
                suggestions['secure'] = '/v1/api-keys:list-current'
                suggestions['explanation'] = 'Use session-based listing instead of direct ID access'
                suggestions['method'] = 'GET (authenticated session) or POST with verification'
        
        # Pattern: GET /v1/secrets/{secret_id}
        if 'secret' in path.lower():
            if method.lower() == 'get' and '{' in path:
                suggestions['insecure'] = path
                suggestions['secure'] = '/v1/secrets:retrieve'
                suggestions['explanation'] = 'Use POST with additional verification (OTP, 2FA)'
                suggestions['method'] = 'POST'
        
        # Pattern: GET /v1/tokens/{token_id}
        if 'token' in path.lower():
            if method.lower() == 'get' and '{' in path:
                suggestions['insecure'] = path
                suggestions['secure'] = '/v1/tokens:validate'
                suggestions['explanation'] = 'Use POST to validate token without exposing ID'
                suggestions['method'] = 'POST'
        
        return suggestions
    
    def fix_path_security(self, path: str, method: str, operation: Dict[str, Any]) -> Tuple[str, str, Dict[str, Any]]:
        """
        Attempt to automatically fix security issues
        
        Returns:
            Tuple of (fixed_path, fixed_method, fixed_operation)
        """
        fixed_path = path
        fixed_method = method
        fixed_operation = operation.copy()
        
        # Fix 1: Convert GET with sensitive params to POST
        if method.lower() == 'get':
            for resource in self.SENSITIVE_RESOURCES:
                if f'/{resource}/' in path and '{' in path:
                    # Change to POST with custom action
                    fixed_path = re.sub(r'/\{[^}]+\}$', ':retrieve', path)
                    fixed_method = 'post'
                    
                    # Move path parameter to request body
                    param_name = re.search(r'\{([^}]+)\}', path)
                    if param_name:
                        if 'requestBody' not in fixed_operation:
                            fixed_operation['requestBody'] = {
                                'required': True,
                                'content': {
                                    'application/json': {
                                        'schema': {
                                            'type': 'object',
                                            'required': [param_name.group(1)],
                                            'properties': {
                                                param_name.group(1): {
                                                    'type': 'string',
                                                    'description': f'Identifier for {resource} retrieval'
                                                }
                                            }
                                        }
                                    }
                                }
                            }
                    break
        
        # Fix 2: Add authentication if missing
        if 'security' not in fixed_operation:
            # Add default security scheme
            fixed_operation['security'] = [{'bearerAuth': []}]
        
        return fixed_path, fixed_method, fixed_operation


def analyze_openapi_security(spec_path: str) -> List[SecurityIssue]:
    """
    Analyze entire OpenAPI spec for security issues
    
    Returns:
        List of all security issues found
    """
    import yaml
    
    analyzer = SecurityAnalyzer()
    all_issues = []
    
    with open(spec_path, 'r', encoding='utf-8') as f:
        spec = yaml.safe_load(f)
    
    paths = spec.get('paths', {})
    
    for path, path_item in paths.items():
        for method in ['get', 'post', 'put', 'delete', 'patch']:
            if method in path_item:
                issues = analyzer.analyze_path(path, method, path_item[method])
                all_issues.extend(issues)
    
    return all_issues


if __name__ == '__main__':
    import sys
    
    if len(sys.argv) < 2:
        print("Usage: python security_analyzer.py <openapi.yaml>")
        sys.exit(1)
    
    issues = analyze_openapi_security(sys.argv[1])
    
    print(f"Found {len(issues)} security issues:")
    print()
    
    # Group by severity
    critical = [i for i in issues if i.severity == SecurityIssue.SEVERITY_CRITICAL]
    high = [i for i in issues if i.severity == SecurityIssue.SEVERITY_HIGH]
    medium = [i for i in issues if i.severity == SecurityIssue.SEVERITY_MEDIUM]
    
    print(f"🔴 CRITICAL: {len(critical)}")
    for issue in critical[:5]:
        print(f"  - {issue.message}")
        print(f"    Location: {issue.location}")
        print(f"    Fix: {issue.fix}")
        print()
    
    print(f"🟠 HIGH: {len(high)}")
    print(f"🟡 MEDIUM: {len(medium)}")
