# InsuranceEngine — Modular .NET 8 Policy & Claims Core

> **Modular Monolith · CQRS · MediatR · gRPC · EF Core 8 · Kafka · PostgreSQL · OpenTelemetry**

InsuranceEngine is the high-performance core for policy lifecycle and claims management within the LabAid InsureTech platform. It provides a modular, scalable architecture for handling the complexities of modern insurance products.

## Architecture

### Technology Stack

- **.NET 8 LTS** - High-performance runtime
- **gRPC** (Grpc.AspNetCore) - Low-latency service communication
- **EF Core 8 + Npgsql** - Database access (external SQL migrations)
- **MediatR 12** - CQRS & Mediator pattern implementation
- **FluentValidation 11** - Type-safe request validation
- **Confluent.Kafka** - Resilient event streaming
- **Redis** - Distributed caching & state management
- **Serilog** - Structured observability logging
- **OpenTelemetry** - Distributed tracing and metrics

### Design Patterns

- **Modular Monolith** - Logical separation of bounded contexts
- **CQRS** - Command Query Responsibility Segregation
- **Domain Events** - Decoupled side-effects via event-driven architecture
- **Repository Pattern** - Abstracted data persistence
- **Validation Pipeline** - Automated request validation via MediatR behaviors

## Project Structure

```
backend/insurance_engine/
├── InsuranceEngine.sln
├── src/
│   ├── InsuranceEngine.SharedKernel/     # Shared Domain/Building Blocks
│   ├── InsuranceEngine.ApiHost/          # Unified Kestrel/gRPC Host
│   ├── InsuranceEngine.Products/         # Product Catalog & Pricing
│   ├── InsuranceEngine.Beneficiary/      # Beneficiary & KYC Context
│   ├── InsuranceEngine.Policy/           # Policy Issuance & Lifecycle
│   ├── InsuranceEngine.Underwriting/     # Risk Assessment & Decisioning
│   ├── InsuranceEngine.Claims/           # FNOL & Settlement
│   ├── InsuranceEngine.Fraud/            # Fraud Detection & Rules
│   ├── InsuranceEngine.Partners/         # Partner Management
│   └── InsuranceEngine.Commission/       # Revenue & Commission Logic
└── tests/
    ├── InsuranceEngine.*.Tests/          # Context-specific Unit Tests
    └── InsuranceEngine.Integration.Tests/ # Cross-module Integration Tests
```

## Services & Ports

| Service | gRPC Port | HTTP Port | Responsibility |
|---------|-----------|-----------|----------------|
| Insurance API | 50051 | 5000 | Main entry point for Insurance operations |
| Product Service | - | - | Catalog management & pricing engine |
| Policy Service | - | - | Issuance, endorsement, and lifecycle |
| Claims Service | - | - | Claim filing, approval, and payout |

## Database Management

InsuranceEngine follows the **PoliSync Pattern** for database schema management:
- **Schema**: Dedicated `insurance_schema` in PostgreSQL.
- **Migrations**: Externally managed via the `inscore` (Go) project.
- **Access**: EF Core 8 is used for data access with explicit schema mapping.

## Getting Started

Refer to the [QUICKSTART.md](QUICKSTART.md) for local setup, dependency configuration, and development workflows.

## Monitoring

- **Traces**: OpenTelemetry traces exported to Console/OTLP.
- **Metrics**: OpenTelemetry metrics for system performance.
- **Health**: `/health` endpoint for readiness/liveness checks.

## License

Proprietary - LabAid InsureTech Platform
