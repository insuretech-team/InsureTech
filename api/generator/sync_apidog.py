#!/usr/bin/env python3
"""
Apidog OpenAPI Import Script
Imports OpenAPI 3.x specification to Apidog project

⚠️ NOTE: This script needs the correct API payload structure from Apidog documentation.
Currently getting 422 "Parameter is missing" error.

TODO: Check your Apidog account's API documentation for the exact import-data payload format:
- Log in to Apidog
- Go to Settings > Open API
- Check the import-data endpoint documentation
- Update the payload structure in import_to_apidog() function

Official Documentation: https://help.apidog.com/api-references/import-data
API Endpoint: POST /api/v1/projects/{projectId}/import-data

Environment Variables Required:
- API_DOG_TOKEN: Your Apidog API token (from Settings > Personal Access Token)
- APIDOG_PROJECT_ID: Your Apidog project ID (from project URL)
"""

import os
import sys
import json
from pathlib import Path

try:
    import requests
except ImportError:
    print("Error: 'requests' module not found. Install with: pip install requests")
    sys.exit(1)

try:
    import yaml
except ImportError:
    print("Error: 'pyyaml' module not found. Install with: pip install pyyaml")
    sys.exit(1)


def load_openapi_spec(spec_path: Path) -> dict:
    """Load and validate OpenAPI specification"""
    if not spec_path.exists():
        raise FileNotFoundError(f"OpenAPI spec not found: {spec_path}")
    
    with open(spec_path, 'r', encoding='utf-8') as f:
        spec = yaml.safe_load(f)
    
    if not spec.get('openapi'):
        raise ValueError("Invalid OpenAPI spec: 'openapi' field missing")
    
    return spec


def import_to_apidog(api_token: str, project_id: str, openapi_spec: dict, spec_path: Path) -> bool:
    """
    Import OpenAPI specification to Apidog project
    
    API Documentation: https://help.apidog.com/api-references/import-data
    Endpoint: POST https://api.apidog.com/api/v1/projects/{projectId}/import-data
    """
    
    url = f"https://api.apidog.com/api/v1/projects/{project_id}/import-data"
    
    headers = {
        "Authorization": f"Bearer {api_token}"
    }
    
    print(f"Importing to Apidog project: {project_id}")
    print(f"OpenAPI version: {openapi_spec.get('openapi')}")
    print(f"API title: {openapi_spec.get('info', {}).get('title', 'N/A')}")
    
    # Try file upload approach - many APIs prefer this for large specs
    try:
        # Read OpenAPI spec as YAML string
        with open(spec_path, 'r', encoding='utf-8') as f:
            openapi_yaml_content = f.read()
        
        # Try different payload structures
        # Structure 1: Direct spec in 'input' field
        payload = {
            "input": {
                "type": "openapi",
                "data": openapi_yaml_content
            }
        }
        
        print("\nAttempting import with payload structure 1...")
        response = requests.post(
            url, 
            headers={**headers, "Content-Type": "application/json"}, 
            json=payload, 
            timeout=60
        )
        
        if response.status_code == 422:
            # Structure 2: Flat structure
            print("Structure 1 failed, trying structure 2...")
            payload = {
                "type": "openapi",
                "data": openapi_yaml_content
            }
            response = requests.post(
                url, 
                headers={**headers, "Content-Type": "application/json"}, 
                json=payload, 
                timeout=60
            )
        
        if response.status_code == 200:
            result = response.json()
            print("\n✓ Import successful!")
            
            # Parse import results
            if 'data' in result:
                data = result['data']
                print(f"  - APIs imported: {data.get('apiCount', 0)}")
                print(f"  - Schemas imported: {data.get('schemaCount', 0)}")
            
            return True
            
        elif response.status_code == 401:
            print("\n✗ Authentication failed")
            print("  Check your API_DOG_TOKEN is valid")
            return False
            
        elif response.status_code == 404:
            print("\n✗ Project not found")
            print(f"  Check your APIDOG_PROJECT_ID: {project_id}")
            return False
            
        else:
            print(f"\n✗ Import failed: HTTP {response.status_code}")
            try:
                error_data = response.json()
                print(f"  Error: {error_data.get('message', response.text)}")
            except:
                print(f"  Response: {response.text[:500]}")
            return False
            
    except requests.exceptions.Timeout:
        print("\n✗ Request timeout (>60s)")
        print("  The OpenAPI spec might be too large")
        return False
        
    except requests.exceptions.RequestException as e:
        print(f"\n✗ Request failed: {str(e)}")
        return False


def main():
    """Main execution"""
    print("=" * 60)
    print("  Apidog OpenAPI Import")
    print("=" * 60)
    print()
    
    # Get configuration from environment variables
    api_token = os.getenv('API_DOG_TOKEN')
    project_id = os.getenv('APIDOG_PROJECT_ID')
    
    # Validate configuration
    if not api_token:
        print("✗ Error: API_DOG_TOKEN environment variable not set")
        print()
        print("To fix this:")
        print("  1. Get your token from: Apidog > Settings > Personal Access Token")
        print("  2. Add to .env file: API_DOG_TOKEN=your_token_here")
        sys.exit(1)
    
    if not project_id:
        print("✗ Error: APIDOG_PROJECT_ID environment variable not set")
        print()
        print("To fix this:")
        print("  1. Open your project in Apidog")
        print("  2. Get project ID from URL: apidog.com/project/{PROJECT_ID}")
        print("  3. Add to .env file: APIDOG_PROJECT_ID=your_project_id")
        sys.exit(1)
    
    # Locate OpenAPI spec
    api_root = Path(__file__).parent.parent
    openapi_path = api_root / "openapi.yaml"
    
    print(f"OpenAPI spec: {openapi_path}")
    print(f"Project ID:   {project_id}")
    print(f"Token:        {api_token[:10]}..." if api_token else "None")
    print()
    
    try:
        # Load OpenAPI specification
        print("[1/2] Loading OpenAPI specification...")
        openapi_spec = load_openapi_spec(openapi_path)
        
        # Count endpoints and schemas
        paths_count = len(openapi_spec.get('paths', {}))
        schemas_count = len(openapi_spec.get('components', {}).get('schemas', {}))
        
        print(f"  ✓ Loaded: {paths_count} paths, {schemas_count} schemas")
        print()
        
        # Import to Apidog
        print("[2/2] Importing to Apidog...")
        success = import_to_apidog(api_token, project_id, openapi_spec, openapi_path)
        
        print()
        print("=" * 60)
        
        if success:
            print("✓ Import completed successfully")
            print()
            print("View in Apidog:")
            print(f"  https://apidog.com/project/{project_id}")
            sys.exit(0)
        else:
            print("✗ Import failed")
            sys.exit(1)
            
    except FileNotFoundError as e:
        print(f"✗ Error: {str(e)}")
        sys.exit(1)
        
    except ValueError as e:
        print(f"✗ Error: {str(e)}")
        sys.exit(1)
        
    except Exception as e:
        print(f"✗ Unexpected error: {str(e)}")
        import traceback
        traceback.print_exc()
        sys.exit(1)


if __name__ == "__main__":
    main()
