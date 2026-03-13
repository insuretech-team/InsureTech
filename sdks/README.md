# InsureTech SDKs

This directory contains SDK generators and generated SDKs for the InsureTech API Platform.

## Structure

```
sdks/
├── sdk-generator/          # SDK generators for different languages
│   ├── go/                # Go SDK generator
│   │   ├── templates/     # Go SDK templates
│   │   ├── generator.go   # Generator implementation
│   │   ├── generate.ps1   # Windows generation script
│   │   ├── generate.sh    # Unix generation script
│   │   ├── README.md      # Generator documentation
│   │   └── QUICKSTART.md  # Quick start guide
│   ├── python/            # Python SDK generator (future)
│   ├── java/              # Java SDK generator (future)
│   └── typescript/        # TypeScript SDK generator (future)
├── go-sdk/                # Generated Go SDK
│   ├── pkg/              # SDK packages
│   ├── examples/         # Usage examples
│   ├── docs/            # API documentation
│   └── README.md        # SDK documentation
└── README.md            # This file
```

## Available SDKs

### Go SDK

**Status**: ✅ Available  
**Location**: `./go-sdk`  
**Generator**: `./sdk-generator/go`

The Go SDK is auto-generated from the OpenAPI specification located at `C:\_DEV\GO\InsureTech\api\openapi.yaml`, using the Protocol Buffer definitions in `C:\_DEV\GO\InsureTech\proto` as the source of truth.

#### Quick Start

```bash
# Install the SDK
go get github.com/insuretech/go-sdk

# Use in your code
import insuretech "github.com/insuretech/go-sdk"
```

See [Go SDK Documentation](./go-sdk/README.md) for detailed usage.

### Python SDK

**Status**: 🚧 Planned  
**Location**: `./python-sdk` (future)

### Java SDK

**Status**: 🚧 Planned  
**Location**: `./java-sdk` (future)

### TypeScript SDK

**Status**: 🚧 Planned  
**Location**: `./typescript-sdk` (future)

## Generating SDKs

### Prerequisites

- Access to the OpenAPI specification: `C:\_DEV\GO\InsureTech\api\openapi.yaml`
- Access to Protocol Buffer definitions: `C:\_DEV\GO\InsureTech\proto`

### Generate Go SDK

#### Windows (PowerShell)

```powershell
cd C:\_DEV\GO\InsureTech\sdks\sdk-generator\go
.\generate.ps1
```

#### Unix/Linux/Mac

```bash
cd C:/_DEV/GO/InsureTech/sdks/sdk-generator/go
chmod +x generate.sh
./generate.sh
```

#### Manual Generation

```bash
cd C:\_DEV\GO\InsureTech\sdks\sdk-generator\go
go build -o generator.exe generator.go
./generator.exe
```

## SDK Architecture

All SDKs follow a consistent architecture:

### 1. Client Wrapper
- Main client class/struct
- Configuration options
- Authentication handling
- Request/response interceptors

### 2. Service Clients
- One service client per API domain
- Methods for each API endpoint
- Type-safe request/response handling

### 3. Models
- Auto-generated from OpenAPI schemas
- Validation logic
- Serialization/deserialization

### 4. Error Handling
- Standardized error types
- HTTP status code handling
- Error recovery strategies

### 5. Utilities
- Pagination helpers
- Retry logic
- Rate limiting
- Logging

## Source of Truth

The SDKs are generated from two primary sources:

1. **Protocol Buffers** (`C:\_DEV\GO\InsureTech\proto`)
   - Canonical data structure definitions
   - Service contracts
   - Message types

2. **OpenAPI Specification** (`C:\_DEV\GO\InsureTech\api\openapi.yaml`)
   - REST API endpoints
   - Request/response schemas
   - Authentication schemes
   - API documentation

## Development Workflow

1. **Define Proto Messages** - Update proto files in `proto/` directory
2. **Generate OpenAPI Spec** - Run API generator to create OpenAPI spec
3. **Generate SDKs** - Run SDK generators to create client libraries
4. **Test SDKs** - Run tests against generated SDKs
5. **Publish SDKs** - Publish to package managers (npm, PyPI, Maven, etc.)

## Generator Templates

Each SDK generator uses templates to customize the generated code:

- `client.go.tmpl` - Main client implementation
- `config.go.tmpl` - Configuration structures
- `errors.go.tmpl` - Error handling
- `models.go.tmpl` - Data model modifiers
- `services.go.tmpl` - Service client implementations

Templates are located in `sdk-generator/<language>/templates/`

## Testing

Each SDK includes:

- Unit tests for generated code
- Integration tests against API
- Example programs demonstrating usage

## Documentation

Each SDK includes:

- **README.md** - Overview and installation
- **QUICKSTART.md** - Getting started guide
- **API.md** - Complete API reference
- **Examples** - Code samples for common scenarios

## Maintenance

SDKs are auto-generated and should not be manually edited. To make changes:

1. Update the source (proto files or OpenAPI spec)
2. Regenerate the SDK
3. Test the changes
4. Publish new version

## Version Management

SDK versions follow semantic versioning (SemVer):

- **Major version** - Breaking changes
- **Minor version** - New features (backward compatible)
- **Patch version** - Bug fixes

SDK versions are synchronized with the API version.

## Support

For SDK issues or questions:

- Documentation: https://docs.insuretech.com
- Email: support@insuretech.com
- GitHub Issues: https://github.com/insuretech/sdks/issues

## License

All SDKs are licensed under the MIT License.
