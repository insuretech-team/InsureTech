import os
import re
import yaml

# Import our custom modules
try:
    from name_transformer import NameTransformer
    from description_loader import DescriptionLoader
except ImportError:
    NameTransformer = None
    DescriptionLoader = None

class PathGenerator:
    def __init__(self, registry, descriptions_dir=None):
        self.registry = registry
        self.name_transformer = NameTransformer() if NameTransformer else None
        self.description_loader = DescriptionLoader(descriptions_dir) if DescriptionLoader and descriptions_dir else None

    def generate_path_item(self, method_data, service_name):
        """
        Generates an OpenAPI Path Item for a single method.
        Returns: (path_url, verb, path_item_dict)
        """
        http_rule = method_data.get('http_rule')
        if not http_rule:
            return None, None, None

        # Determine Verb and Path
        verb, full_path_url = self._extract_verb_and_path(http_rule)
        if not verb:
            return None, None, None

        # Convert custom actions (:cancel) to query parameters (?action=cancel)
        path_url, action_param = self._process_custom_actions(full_path_url)
        
        # Fix kebab-case in path
        path_url = self._to_kebab_case_path(path_url)
        
        # Split query params from URL pattern
        if '?' in path_url:
            path_url, query_string = path_url.split('?', 1)
        else:
            query_string = ""

        # Extract Parameters from Path
        parameters = self._extract_parameters(path_url)
        
        # Note: Custom methods use colon syntax in path (e.g., /v1/resource/{id}:cancel)
        # This is part of the path itself, not a parameter
        
        # Add query params from the split URL
        if query_string:
            for pair in query_string.split('&'):
                if '=' in pair:
                    k, v = pair.split('=', 1)
                    parameters.append({
                        "name": k,
                        "in": "query",
                        "required": True,
                        "schema": {"type": "string", "enum": [v]} 
                    })
        
        # Load operation description
        method_name = method_data.get('name')
        method_comment = method_data.get('comment', '')
        if self.description_loader:
            desc_data = self.description_loader.load_operation_description(
                service_name, method_name, proto_comment=method_comment
            )
            summary = desc_data['summary']
            description = desc_data.get('description', '')
        else:
            # Fallback to proto comment if available
            if method_comment:
                summary = method_comment.split('\n')[0].strip()
                description = method_comment
            else:
                summary = method_name
                description = ""
        
        # Determine appropriate status code
        success_code = self._get_success_code(verb, method_name)
        
        # Build operation with proper responses
        operation = {
            "summary": summary,
            "operationId": f"{service_name}_{method_name}",
            "responses": self._build_responses(success_code, method_data)
        }
        
        if description:
            operation['description'] = description
        
        if parameters:
            operation['parameters'] = parameters
            
        # Request Body
        if verb in ['post', 'put', 'patch']:
            body_field = http_rule.body
            if body_field:
                schema_ref = self._get_ref(method_data['input_type'])
                operation['requestBody'] = {
                    "content": {
                        "application/json": {
                            "schema": schema_ref
                        }
                    },
                    "required": True
                }
                
        # Construct Path Item (Single Operation)
        path_item = {
            verb: operation
        }
        
        return path_url, verb, path_item
    
    def _process_custom_actions(self, path_url):
        """Keep custom method with colon syntax as per Google API Design Guide
        
        Google API Design Guide specifies custom methods should use colon syntax:
        POST /v1/policies/{id}:cancel
        POST /v1/policies/{id}:renew
        
        This is NOT a query parameter, it's part of the URL path itself.
        OpenAPI 3.x supports colons in paths.
        """
        # Custom methods with colon are kept as-is in the path
        # /v1/policies/{id}:cancel → /v1/policies/{id}:cancel (no transformation)
        if ':' in path_url and not path_url.startswith('http'):
            # Convert action to kebab-case for consistency
            parts = path_url.split(':')
            base_path = parts[0]
            action = parts[1] if len(parts) > 1 else None
            if action and self.name_transformer:
                action = self.name_transformer._to_kebab_case(action)
                return f"{base_path}:{action}", None
        return path_url, None
    
    def _to_kebab_case_path(self, path):
        """Convert camelCase segments in path to kebab-case"""
        # Split by / and process each segment
        segments = path.split('/')
        result = []
        for seg in segments:
            # Skip parameters {id} and version segments v1
            if seg.startswith('{') or (seg.startswith('v') and seg[1:].isdigit()):
                result.append(seg)
            elif self.name_transformer:
                result.append(self.name_transformer._to_kebab_case(seg))
            else:
                result.append(seg)
        return '/'.join(result)
    
    def _get_success_code(self, verb, method_name):
        """Determine appropriate success status code"""
        if verb == 'post':
            # Creation vs action
            if method_name.startswith('Create') or method_name.startswith('Register'):
                return '201'
            return '200'
        elif verb == 'delete':
            return '204'
        return '200'
    
    def _build_responses(self, success_code, method_data):
        """Build complete response dictionary with errors"""
        responses = {}
        
        # Success response
        if success_code == '204':
            responses['204'] = {
                "description": "No content"
            }
        elif success_code == '201':
            responses['201'] = {
                "description": "Resource created successfully",
                "content": {
                    "application/json": {
                        "schema": self._get_ref(method_data['output_type'])
                    }
                },
                "headers": {
                    "Location": {
                        "description": "URL of the created resource",
                        "schema": {"type": "string"}
                    }
                }
            }
        else:
            responses['200'] = {
                "description": "Successful response",
                "content": {
                    "application/json": {
                        "schema": self._get_ref(method_data['output_type'])
                    }
                }
            }
        
        # Error responses
        responses['400'] = {
            "description": "Bad request - Invalid input parameters",
            "content": {
                "application/json": {
                    "schema": {"$ref": "#/components/schemas/Error"}
                }
            }
        }
        responses['401'] = {
            "description": "Unauthorized - Authentication required",
            "content": {
                "application/json": {
                    "schema": {"$ref": "#/components/schemas/Error"}
                }
            }
        }
        responses['403'] = {
            "description": "Forbidden - Insufficient permissions",
            "content": {
                "application/json": {
                    "schema": {"$ref": "#/components/schemas/Error"}
                }
            }
        }
        responses['404'] = {
            "description": "Not found - Resource does not exist",
            "content": {
                "application/json": {
                    "schema": {"$ref": "#/components/schemas/Error"}
                }
            }
        }
        responses['500'] = {
            "description": "Internal server error",
            "content": {
                "application/json": {
                    "schema": {"$ref": "#/components/schemas/Error"}
                }
            }
        }
        
        return responses

    def _extract_verb_and_path(self, http_rule):
        """Extracts the HTTP verb and URL path from HttpRule."""
        if http_rule.HasField('get'):
            return 'get', http_rule.get
        if http_rule.HasField('post'):
            return 'post', http_rule.post
        if http_rule.HasField('put'):
            return 'put', http_rule.put
        if http_rule.HasField('delete'):
            return 'delete', http_rule.delete
        if http_rule.HasField('patch'):
            return 'patch', http_rule.patch
        return None, None

    def _extract_parameters(self, path_url):
        """Extracts path parameters from URL template (e.g. /users/{id})."""
        # Regex to find {var}
        matches = re.findall(r'\{([^\}]+)\}', path_url)
        params = []
        for match in matches:
            # Handle {var=*} syntax if present (though rare in simple matching)
            # just take the var name
            var_name = match.split('=')[0]
            params.append({
                "name": var_name,
                "in": "path",
                "required": True,
                "schema": {"type": "string"} # Default to string for path params
            })
        return params

    def _get_ref(self, type_name):
        """Resolves a proto type name to an OpenAPI schema ref."""
        # Normalize type name: .package.Message -> package.Message
        if type_name.startswith('.'):
            type_name = type_name[1:]
            
        ref = self.registry.get_ref(type_name)
        if ref:
            return {"$ref": ref}
        else:
            return {"type": "object", "description": f"Unresolved type: {type_name}"}
