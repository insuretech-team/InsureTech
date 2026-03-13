# InsureTech OpenAPI Specification

This directory contains the OpenAPI 3.1 specification for the LabAid InsureTech platform, generated from Protocol Buffer service definitions.

## Directory Structure

```
api/
├── schemas/          # OpenAPI schema definitions (DTOs/entities)
├── paths/            # API endpoint path definitions
├── components/       # Reusable OpenAPI components (security, responses, etc.)
├── generated/        # Generated OpenAPI specification files
├── generator.py      # Python script to generate OpenAPI from proto
└── openapi.yaml      # Main OpenAPI 3.1 specification
```

## Generation Process

The OpenAPI specification is generated from proto service definitions:

1. **Proto Services** → Read all `*_service.proto` files
2. **Extract HTTP Annotations** → Parse `google.api.http` options
3. **Generate Schemas** → Convert proto messages to OpenAPI schemas
4. **Generate Paths** → Create path definitions from RPCs
5. **Assemble OpenAPI** → Combine into single `openapi.yaml`

## Usage

### Generate OpenAPI Specification

```bash
# From project root
python api/generator.py

# Or using PowerShell
pwsh api/generate-openapi.ps1
```

### View Generated Spec

```bash
# View in Swagger UI (requires docker)
docker run -p 8080:8080 -v ${PWD}/api:/api swaggerapi/swagger-ui

# Or use online viewer
# Upload api/generated/openapi.yaml to https://editor.swagger.io/
```

## Services Included

All proto services are converted to OpenAPI paths:

- **AuthN** - Authentication (login, register, session management)
- **AuthZ** - Authorization (roles, permissions, policies)
- **Policy** - Policy management and lifecycle
- **Claims** - Claims submission and processing
- **Payment** - Payment processing and verification
- **Underwriting** - Quote generation and approval
- **Partner** - Partner management and integration
- **Product** - Insurance product catalog
- **Commission** - Commission calculation and payouts
- **Refund** - Refund processing
- **Notification** - Notification delivery
- **Document** - Document generation
- **Support** - Customer support tickets
- **Fraud** - Fraud detection and rules
- **KYC** - KYC verification
- **Beneficiary** - Beneficiary management
- **Workflow** - Business process workflows
- **Task** - Task management
- **Report** - Report generation
- **Analytics** - Analytics and metrics
- **Audit** - Audit logging
- **Voice** - Voice interaction
- **IoT** - IoT device management
- **MFS** - Mobile Financial Services integration
- **API Key** - API key management
- **Tenant** - Multi-tenancy management
- **Insurer** - Insurer management

## OpenAPI Features

- **OpenAPI Version**: 3.1.0
- **Authentication**: Bearer tokens, API keys
- **Content Types**: `application/json`, `application/grpc+json`
- **Error Responses**: Standardized error schema
- **Security Schemes**: JWT, API Key, OAuth2
- **Examples**: Request/response examples for all operations
- **Tags**: Organized by domain/service
- **Descriptions**: Generated from proto comments

## Customization

### Adding Custom Schemas

Add custom schema files to `api/schemas/`:

```yaml
# api/schemas/custom_error.yaml
CustomError:
  type: object
  properties:
    code:
      type: string
    message:
      type: string
```

### Adding Custom Paths

Add custom path files to `api/paths/`:

```yaml
# api/paths/health.yaml
/health:
  get:
    summary: Health check
    responses:
      '200':
        description: Service is healthy
```

## Validation

Validate the generated OpenAPI spec:

```bash
# Using openapi-generator-cli
docker run --rm -v ${PWD}/api:/api openapitools/openapi-generator-cli validate -i /api/generated/openapi.yaml

# Using spectral
npx @stoplight/spectral-cli lint api/generated/openapi.yaml
```

## Code Generation from OpenAPI

Generate client SDKs from the OpenAPI spec:

```bash
# Generate TypeScript client
openapi-generator-cli generate -i api/generated/openapi.yaml -g typescript-axios -o sdk/typescript

# Generate Python client
openapi-generator-cli generate -i api/generated/openapi.yaml -g python -o sdk/python

# Generate Java client
openapi-generator-cli generate -i api/generated/openapi.yaml -g java -o sdk/java
```

## Notes

- The OpenAPI spec is **read-only** - modify proto files instead
- Regenerate after proto changes: `python api/generator.py`
- HTTP annotations in proto files drive the REST API structure
- gRPC services are exposed via gRPC-Gateway with HTTP/JSON mapping
