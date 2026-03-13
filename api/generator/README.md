# Proto to OpenAPI Generator

This directory contains the Python-based generator that converts Protocol Buffer definitions into OpenAPI 3.0 specifications.

## Prerequisites

1. **Python 3.8+** with the following packages:
   ```bash
   pip install protobuf pyyaml ruamel.yaml
   ```

2. **Buf CLI** - For proto compilation
   ```bash
   # Install buf from https://buf.build/docs/installation
   ```

3. **Google API Proto Files (Python)**
   ```powershell
   # Run this once to generate Python files for google.api annotations
   .\generate_google_api_protos.ps1
   ```

## Directory Structure

```
generator/
├── main.py                          # Main entry point
├── proto_parser.py                  # Parses proto descriptors
├── schema_generator.py              # Generates OpenAPI schemas
├── path_generator.py                # Generates API paths from HTTP annotations
├── assembler.py                     # Assembles final OpenAPI spec
├── name_transformer.py              # Transforms proto names to OpenAPI names
├── registry.py                      # Type registry and collision detection
├── proto_source_parser.py           # Parses proto source for annotations
├── path_validator.py                # Validates generated paths
├── endpoint_mapper.py               # Maps endpoints for documentation
├── post_processor.py                # Post-processes OpenAPI spec
├── fix_all_warnings.py              # Fixes validation warnings
├── fix_pagination.py                # Adds pagination to list endpoints
├── generate_google_api_protos.ps1   # Setup script for google.api
└── gen/                             # Generated Python proto files
    └── google/
        └── api/
            ├── annotations_pb2.py   # Google API annotations
            └── http_pb2.py          # HTTP rule definitions
```

## How It Works

### 1. Descriptor Compilation
```python
# Uses buf to compile all protos + dependencies into a single descriptor set
buf build -o descriptors.pb --as-file-descriptor-set
```

This creates a binary file containing:
- All your proto definitions
- All dependencies (google.api, google.protobuf, etc.)
- Source code info (comments, locations)

### 2. Parsing HTTP Annotations
```protobuf
service AuthService {
  rpc Register(RegisterRequest) returns (RegisterResponse) {
    option (google.api.http) = {
      post: "/v1/auth/register"
      body: "*"
    };
  }
}
```

The parser extracts:
- HTTP method (POST, GET, etc.)
- Path pattern (`/v1/auth/register`)
- Body field mapping
- Path parameters

### 3. Schema Generation
- Proto messages → OpenAPI schemas
- Proto enums → OpenAPI enums
- Proto field types → OpenAPI types
- Proto comments → OpenAPI descriptions

### 4. Path Generation
- Service methods with HTTP annotations → OpenAPI paths
- Path parameters extracted from URL patterns
- Request body from body field
- Response schema from output type

### 5. Assembly
Combines all generated files into a single `openapi.yaml`:
```yaml
openapi: 3.0.0
info: {...}
paths:
  /v1/auth/register:
    post: {...}
components:
  schemas: {...}
```

## Google API Annotations Setup

### Why Needed?
The generator needs to parse `google.api.http` options from proto files. These are defined in:
- `google/api/annotations.proto`
- `google/api/http.proto`

### How It Works
Python's protobuf library needs compiled Python files (`*_pb2.py`) to parse extensions.

**The Setup Script:**
1. Uses buf to generate Python files from googleapis
2. Places them in `gen/google/api/`
3. Creates `__init__.py` files for proper Python package structure

**The Import Code (in proto_parser.py):**
```python
# Extend Python's google namespace package to include our generated files
import google
google.__path__ = list(google.__path__)  # Convert from _NamespacePath
google.__path__.insert(0, "gen/google")  # Add our gen directory

# Now we can import
from google.api import annotations_pb2
```

### Troubleshooting

**Error: `No module named 'google.api'`**
```powershell
# Solution: Generate the Python proto files
cd api\generator
.\generate_google_api_protos.ps1
```

**Error: `ModuleNotFoundError: No module named 'google.protobuf'`**
```bash
# Solution: Install protobuf
pip install protobuf
```

**Warning: `http_rule: None` for all methods**
```
# This means annotations aren't being parsed. Check:
1. Does gen/google/api/annotations_pb2.py exist?
2. Is google.protobuf installed?
3. Try regenerating: .\generate_google_api_protos.ps1
```

## Running the Generator

### From API Pipeline Script (Recommended):
```powershell
cd <project-root>
.\run_api_pipeline.ps1 -Fast
```

### Directly (For Development):
```powershell
cd api\generator
python main.py --discover --proto-root=../../proto --api-root=..
```

### Arguments:
- `--discover`: Auto-discover and compile all proto files
- `--descriptor`: Path to descriptor file (default: `../input/descriptors.pb`)
- `--proto-root`: Root directory for proto files (default: `../../proto`)
- `--api-root`: Output directory for API files (default: `..`)

## Output Files

### Schemas (672 files)
```
api/schemas/insuretech/authn/entity/v1/User.yaml
```

### Events (170 files)
```
api/events/insuretech/authn/events/v1/UserRegisteredEvent.yaml
```

### Enums (134 files)
```
api/enums/PolicyStatus.yaml
```

### Paths (30 files)
```
api/paths/insuretech/authn/services/v1/AuthService.yaml
```

## Features

### Name Collision Detection
When multiple protos have the same message name, the generator:
1. Detects the collision
2. Adds a prefix based on package name
3. Reports all collisions

Example:
```
⚠️  COLLISION DETECTED: 'Status'
    Existing: insuretech.policy.entity.v1.Status → PolicyStatus
    New:      insuretech.claims.entity.v1.Status → ClaimsStatus
```

### Request/Response Transformation
Transforms proto naming conventions to REST conventions:
```
GetUserRequest  → UserRetrievalRequest
GetUserResponse → UserRetrievalResponse
CreateUserRequest → UserCreationRequest
```

### Pagination Support
Automatically adds pagination parameters to list endpoints:
```yaml
- name: page_size
  in: query
  schema:
    type: integer
    default: 50
- name: page_token
  in: query
  schema:
    type: string
```

### Security Annotations
Extracts field-level security from proto comments:
```protobuf
message User {
  string email = 1; // @security: pii, email
  string ssn = 2;   // @security: sensitive, pii
}
```

Generates:
```yaml
email:
  type: string
  x-security: ["pii", "email"]
```

## Validation

The generator includes validation for:
- Duplicate operation IDs
- Path parameter consistency
- Schema reference integrity
- Required fields on request DTOs

Run validation:
```powershell
.\run_api_pipeline.ps1  # Includes validation step
```

## Debugging

Enable debug output:
```python
# In main.py, add:
import logging
logging.basicConfig(level=logging.DEBUG)
```

Check descriptor contents:
```python
python -c "
from proto_parser import ProtoParser
p = ProtoParser()
p.load_descriptor_set('../input/descriptors.pb')
services = p.get_services()
print(f'Found {len(services)} services')
for s in services[:5]:
    print(f'  - {s[\"descriptor\"].name}')
"
```

## Maintenance

### Adding New Proto Files
1. Add proto file to `proto/` directory
2. Run `buf generate` to update Go/C#/TS
3. Run `.\run_api_pipeline.ps1` to update OpenAPI

### Updating Google API Protos
```powershell
cd api\generator
Remove-Item -Recurse -Force gen
.\generate_google_api_protos.ps1
```

### Cleaning Generated Files
```powershell
cd api
Remove-Item -Recurse -Force schemas, events, enums, paths
Remove-Item openapi.yaml, openapi.json
```

## Performance

Typical generation time for InsureTech project:
- Proto compilation: ~2 seconds
- Schema generation: ~5 seconds
- Path generation: ~1 second
- Assembly: ~1 second
- **Total: ~10 seconds**

Fast mode (skips some validation):
```powershell
.\run_api_pipeline.ps1 -Fast
```

## Contributing

When modifying the generator:
1. Test with existing protos first
2. Check for name collisions in output
3. Validate generated OpenAPI spec
4. Update this README if adding new features

## Support

For issues:
1. Check that `gen/google/api/` exists and has Python files
2. Verify `buf generate` works
3. Look for error messages in pipeline output
4. Check validation report: `api/validation_report.html`
