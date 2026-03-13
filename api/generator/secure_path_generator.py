#!/usr/bin/env python3
"""
Secure Path Generator - Security-First API Route Generation

Prevents common security issues:
1. Sensitive data in GET parameters
2. Direct ID-based access to sensitive resources
3. Missing authentication
4. Predictable identifiers

Usage:
    This module wraps path_generator.py with security checks and fixes
"""

from security_analyzer import SecurityAnalyzer, SecurityIssue
from typing import Dict, Any, List
import re


class SecurePathGenerator:
    """Generates secure API paths with built-in security checks"""
    
    # Resources that should NEVER allow direct GET by ID
    RESTRICTED_RESOURCES = [
        'api-keys', 'api_keys', 'apikeys',
        'secrets', 'credentials', 'tokens',
        'passwords', 'private-keys', 'encryption-keys',
        'auth-tokens', 'refresh-tokens'
    ]
    
    # Sensitive parameter names that should not be in URLs
    SENSITIVE_PARAMS = [
        'api_key', 'api-key', 'apikey',
        'secret', 'token', 'password',
        'credential', 'private_key', 'session_id'
    ]
    
    def __init__(self, registry, output_dir: str):
        self.registry = registry
        self.output_dir = output_dir
        self.security_analyzer = SecurityAnalyzer()
        self.security_issues = []
    
    def should_restrict_direct_access(self, resource_name: str) -> bool:
        """Check if resource should have restricted direct access"""
        resource_lower = resource_name.lower().replace('_', '-')
        return any(restricted in resource_lower for restricted in self.RESTRICTED_RESOURCES)
    
    def transform_insecure_path(self, path: str, method: str, service_name: str) -> tuple:
        """
        Transform insecure paths to secure alternatives
        
        Returns:
            (new_path, new_method, should_skip, reason)
        """
        # Extract resource name from path
        # /v1/api-keys/{api_key_id} -> api-keys
        match = re.search(r'/v1/([^/]+)/\{[^}]+\}', path)
        if not match:
            return path, method, False, None
        
        resource = match.group(1)
        
        # Check if this is a restricted resource
        if self.should_restrict_direct_access(resource):
            
            # Rule 1: No GET by ID for sensitive resources
            if method.lower() == 'get' and '{' in path:
                # Transform to session-based endpoint
                new_path = f'/v1/{resource}:list-current'
                reason = f"Security: GET {path} is insecure. Use authenticated session-based listing instead."
                
                self.security_issues.append({
                    'severity': 'critical',
                    'original_path': f'{method.upper()} {path}',
                    'new_path': f'GET {new_path}',
                    'reason': reason
                })
                
                # Skip generating the insecure endpoint
                return new_path, 'get', True, reason
            
            # Rule 2: Convert to POST with verification for retrieval
            if method.lower() == 'get' and '{' in path:
                # Transform GET /v1/api-keys/{id} -> POST /v1/api-keys:retrieve
                base_path = path.split('/{')[0]
                new_path = f'{base_path}:retrieve'
                
                reason = f"Security: Converted to POST {new_path} with body verification"
                return new_path, 'post', False, reason
        
        return path, method, False, None
    
    def add_security_annotations(self, operation: Dict[str, Any], path: str, method: str) -> Dict[str, Any]:
        """Add security requirements and annotations to operation"""
        
        # Always add authentication for sensitive resources
        for resource in self.RESTRICTED_RESOURCES:
            if resource in path.lower():
                if 'security' not in operation:
                    operation['security'] = [
                        {'bearerAuth': []},  # OAuth2/JWT
                        {'apiKeyAuth': []}    # API Key fallback
                    ]
                
                # Add security note in description
                if 'description' not in operation or not operation['description']:
                    operation['description'] = ''
                
                operation['description'] += '\n\n**Security**: This endpoint requires authentication and authorized access only.'
                break
        
        # Add rate limiting for sensitive operations
        if method.lower() in ['post', 'put', 'delete']:
            if 'x-rate-limit' not in operation:
                operation['x-rate-limit'] = {
                    'limit': 100,
                    'period': '1h',
                    'scope': 'user'
                }
        
        return operation
    
    def validate_path_security(self, path: str, method: str, operation: Dict[str, Any]) -> List[SecurityIssue]:
        """Validate path for security issues"""
        return self.security_analyzer.analyze_path(path, method, operation)
    
    def generate_secure_alternative(self, insecure_path: str, method: str) -> Dict[str, Any]:
        """
        Generate secure alternative endpoint specification
        
        For example:
        - GET /v1/api-keys/{api_key_id} 
        - Becomes: GET /v1/api-keys:list-current (returns user's own keys)
        """
        alternatives = {}
        
        # Pattern: GET /v1/api-keys/{api_key_id}
        if '/api-keys/' in insecure_path or '/api_keys/' in insecure_path:
            alternatives['path'] = '/v1/api-keys:list-current'
            alternatives['method'] = 'GET'
            alternatives['description'] = 'List API keys for the authenticated user. Returns only keys belonging to the current session.'
            alternatives['security'] = [{'bearerAuth': []}]
            alternatives['responses'] = {
                '200': {
                    'description': 'List of user\'s API keys',
                    'content': {
                        'application/json': {
                            'schema': {
                                'type': 'object',
                                'properties': {
                                    'api_keys': {
                                        'type': 'array',
                                        'items': {'$ref': '#/components/schemas/ApiKey'}
                                    }
                                }
                            }
                        }
                    }
                }
            }
        
        # Pattern: GET /v1/secrets/{secret_id}
        elif 'secret' in insecure_path.lower() and '{' in insecure_path:
            alternatives['path'] = insecure_path.split('/{')[0] + ':retrieve'
            alternatives['method'] = 'POST'
            alternatives['description'] = 'Retrieve secret with additional verification (2FA, OTP)'
            alternatives['requestBody'] = {
                'required': True,
                'content': {
                    'application/json': {
                        'schema': {
                            'type': 'object',
                            'required': ['verification_token'],
                            'properties': {
                                'verification_token': {
                                    'type': 'string',
                                    'description': '2FA or OTP token for verification'
                                }
                            }
                        }
                    }
                }
            }
        
        return alternatives
    
    def get_security_report(self) -> str:
        """Generate security transformation report"""
        if not self.security_issues:
            return "✅ No security issues found"
        
        report = ["=" * 80]
        report.append("SECURITY TRANSFORMATIONS APPLIED")
        report.append("=" * 80)
        report.append("")
        report.append(f"Total insecure patterns prevented: {len(self.security_issues)}")
        report.append("")
        
        for issue in self.security_issues:
            report.append(f"🔒 {issue['severity'].upper()}")
            report.append(f"   Original: {issue['original_path']}")
            report.append(f"   Secure:   {issue['new_path']}")
            report.append(f"   Reason:   {issue['reason']}")
            report.append("")
        
        return "\n".join(report)


def should_skip_insecure_endpoint(path: str, method: str) -> tuple:
    """
    Quick check if endpoint should be skipped for security reasons
    
    Returns:
        (should_skip: bool, reason: str)
    """
    # Pattern 1: GET /v1/api-keys/{api_key_id}
    if method.lower() == 'get' and '/api-keys/' in path and '{api_key_id}' in path:
        return True, "Direct GET access to API keys by ID is a critical security vulnerability"
    
    # Pattern 2: GET /v1/secrets/{secret_id}
    if method.lower() == 'get' and '/secrets/' in path and '{' in path:
        return True, "Direct GET access to secrets by ID is a critical security vulnerability"
    
    # Pattern 3: GET /v1/tokens/{token_id}
    if method.lower() == 'get' and '/tokens/' in path and '{' in path:
        return True, "Direct GET access to tokens by ID is a critical security vulnerability"
    
    # Pattern 4: GET /v1/credentials/*
    if method.lower() == 'get' and '/credentials/' in path and '{' in path:
        return True, "Direct GET access to credentials by ID is a critical security vulnerability"
    
    return False, None


if __name__ == '__main__':
    print("Secure Path Generator - Testing")
    print()
    
    # Test insecure patterns
    test_paths = [
        ('GET', '/v1/api-keys/{api_key_id}'),
        ('GET', '/v1/api-keys/{api_key_id}/usage'),
        ('GET', '/v1/secrets/{secret_id}'),
        ('POST', '/v1/policies/{policy_id}/cancel'),  # OK
    ]
    
    for method, path in test_paths:
        should_skip, reason = should_skip_insecure_endpoint(path, method)
        status = "🔴 BLOCK" if should_skip else "✅ ALLOW"
        print(f"{status} {method} {path}")
        if reason:
            print(f"     Reason: {reason}")
