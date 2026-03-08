# PoliSync — C# .NET 8 Insurance Commerce & Policy Engine

> **Proto-first · CQRS · MediatR · gRPC · EF Core 8 · Kafka · PostgreSQL**

PoliSync is the insurance commerce and policy lifecycle engine for the LabAid InsureTech platform. It handles product catalog, quotations, policy issuance, underwriting, claims, renewals, endorsements, and commission calculations.

## Architecture

### Technology Stack

- **.NET 8 LTS** - Runtime
- **gRPC** (Grpc.AspNetCore) - Service communication
- **EF Core 8 + Npgsql** - Database access (schema-mapped, no migrations)
- **MediatR 12** - CQRS pattern implementation
- **FluentValidation 11** - Request validation
- **Confluent.Kafka** - Event streaming
- **Redis** - Distributed caching
- **Serilog** - Structured logging
- **OpenTelemetry** - Observability

### Design Patterns

- **CQRS** - Command Query Responsibility Segregation
- **Domain Events** - Event-driven architecture
- **Repository Pattern** - Data access abstraction
- **Unit of Work** - Transaction management
- **Result Pattern** - No exceptions for domain errors

## Project Structure

```
backend/polisync/
├── PoliSync.sln
├── src/
│   ├── PoliSync.SharedKernel/      # Core abstractions (zero deps)
│   │   ├── Domain/                 # Entity, ValueObject, DomainEvent, Money
│   │   ├── CQRS/                   # ICommand, IQuery, Result<T>
│   │   ├── Messaging/              # IEventBus, IDomainEventHandler
│   │   ├── Persistence/            # IRepository, IUnitOfWork
│   │   ├── Auth/                   # ICurrentUser
│   │   └── Pii/                    # IPiiEncryptor
│   │
│   ├── PoliSync.Infrastructure/    # Infrastructure implementations
│   │   ├── Persistence/            # DbContext, Repository, UnitOfWork
│   │   ├── Messaging/              # KafkaEventBus
│   │   ├── GrpcClients/            # Upstream service clients
│   │   ├── Cache/                  # Redis cache
│   │   ├── Pii/                    # AES-256-GCM encryption
│   │   └── Auth/                   # CurrentUser from JWT
│   │
│   ├── PoliSync.ApiHost/           # Single Kestrel host (all gRPC services)
│   │   ├── Program.cs              # DI wiring, middleware
│   │   ├── appsettings.json        # Configuration
│   │   └── Interceptors/           # JWT, Logging, Validation
│   │
│   ├── PoliSync.Products/          # Product catalog bounded context
│   ├── PoliSync.Quotes/            # Quotation management
│   ├── PoliSync.Orders/            # Order & checkout
│   ├── PoliSync.Policy/            # Core policy lifecycle
│   ├── PoliSync.Endorsement/       # Policy amendments
│   ├── PoliSync.Renewal/           # Renewal schedules & reminders
│   ├── PoliSync.Underwriting/      # Risk assessment
│   ├── PoliSync.Claims/            # Claims & settlement
│   ├── PoliSync.Commission/        # Commission & revenue share
│   └── PoliSync.Refund/            # Refund calculations
│
└── tests/
    ├── PoliSync.Products.Tests/
    ├── PoliSync.Policy.Tests/
    └── PoliSync.Integration.Tests/
```

## Services & Ports

| Service | gRPC Port | HTTP Port | Responsibility |
|---------|-----------|-----------|----------------|
| product-service | 50120 | 50121 | Product catalog, pricing rules |
| quote-service | 50130 | 50131 | Premium calculation, quotations |
| order-service | 50140 | 50141 | Checkout, payment flow |
| commission-service | 50150 | 50151 | Agent/partner commission |
| policy-service | 50160 | 50161 | Policy issuance, lifecycle |
| underwriting-service | 50170 | 50171 | Risk scoring, decisions |
| claim-service | 50210 | 50211 | FNOL, approval, settlement |

## Database Schema

PoliSync uses **read-write access** to PostgreSQL schemas:
- `insurance_schema` - Products, policies, claims, quotes, orders
- `commission_schema` - Commission configs, payouts, revenue share

**Schema is managed by Go migrations** - EF Core is configured for read-write only, no migrations.

## Configuration

### Environment Variables

```bash
DB_PASSWORD=<postgres_password>
REDIS_PASSWORD=<redis_password>
PII_ENCRYPTION_KEY=<base64_32byte_key>
```

### appsettings.json

Key configuration sections:
- `ConnectionStrings` - PostgreSQL, Redis
- `Kafka` - Bootstrap servers, topics
- `GrpcClients` - Upstream Go service URLs
- `Jwt` - Token validation settings
- `Pii` - Encryption key path
- `Commission` - Tax rates, revenue share
- `Renewal` - Reminder schedule, grace period
- `Claims` - Approval thresholds, fraud limits

## Running Locally

### Prerequisites

- .NET 8 SDK
- PostgreSQL 15+ (with insurance_schema created by Go migrations)
- Redis 7+
- Kafka 3.x
- Go backend services running (authn, authz, payment, etc.)

### Build

```bash
cd backend/polisync
dotnet restore
dotnet build
```

### Run

```bash
cd src/PoliSync.ApiHost
dotnet run
```

Services will start on configured ports (50120-50211).

### Health Check

```bash
curl http://localhost:50121/health
```

## Development Workflow

### Adding a New Command

1. Create command in `Application/Commands/`
2. Create command handler implementing `IRequestHandler<TCommand, Result<T>>`
3. Create validator implementing `AbstractValidator<TCommand>`
4. Register in DI (auto-discovered by MediatR)

### Adding a New Domain Event

1. Create event inheriting `DomainEvent` in `Domain/Events/`
2. Raise event in aggregate: `RaiseDomainEvent(new MyEvent(...))`
3. Create event handler implementing `IDomainEventHandler<TEvent>`
4. Publish to Kafka in handler via `IEventBus`

### Adding a New gRPC Service

1. Implement `*ServiceBase` from generated proto
2. Inject `IMediator` and dispatch to commands/queries
3. Map errors to gRPC status codes
4. Register in `Program.cs`: `app.MapGrpcService<MyGrpcService>()`

## Testing

```bash
# Unit tests
dotnet test tests/PoliSync.Products.Tests

# Integration tests (requires Testcontainers)
dotnet test tests/PoliSync.Integration.Tests
```

## Kafka Topics

| Topic | Producer | Purpose |
|-------|----------|---------|
| `insuretech.policy.issued.v1` | policy-service | Policy issuance notification |
| `insuretech.policy.cancelled.v1` | policy-service | Policy cancellation |
| `insuretech.quotation.submitted.v1` | quote-service | Quote submission |
| `insuretech.claim.filed.v1` | claim-service | Claim FNOL |
| `insuretech.commission.payout_created.v1` | commission-service | Commission payout |
| `insuretech.renewal.reminder_sent.v1` | renewal-service | Renewal reminder |

## Security

- **JWT Authentication** - All requests validated via `JwtAuthInterceptor`
- **Tenant Isolation** - All queries filtered by `tenant_id`
- **PII Encryption** - NID, phone numbers encrypted with AES-256-GCM
- **Audit Trail** - All writes logged to audit-service
- **RBAC/ABAC** - Permission checks via authz-service

## Performance Targets

| Endpoint | p95 Target |
|----------|------------|
| CalculatePremium | < 200ms |
| GetProduct | < 100ms (cached) |
| IssuePolicy | < 500ms (sync) |
| FileClaim | < 300ms |

## Deployment

### Docker

```bash
docker build -t polisync:latest .
docker run -p 50120-50211:50120-50211 polisync:latest
```

### Kubernetes

See `k8s/` directory for manifests.

## Monitoring

- **Logs** - Serilog JSON to stdout + file
- **Metrics** - Prometheus endpoint on `:9090/metrics`
- **Traces** - OpenTelemetry to OTLP collector
- **Health** - `/health` endpoint

## Implementation Status

### Phase 1 - Infrastructure Foundation ✅
- [x] Solution structure
- [x] SharedKernel (Entity, ValueObject, Money, Result, CQRS interfaces)
- [x] Infrastructure (DbContext, Repository, UnitOfWork, Kafka, Redis, PII)
- [x] ApiHost (Program.cs, interceptors, configuration)

### Phase 2 - Products & Pricing 🚧
- [ ] Product aggregate
- [ ] Pricing rule engine
- [ ] Product gRPC service
- [ ] Redis caching

### Phase 3 - Quotes & Underwriting 🚧
- [ ] Quotation aggregate
- [ ] Premium calculation
- [ ] Underwriting risk scoring
- [ ] Quote gRPC service

### Phase 4 - Orders & Policy Issuance 🚧
- [ ] Order aggregate
- [ ] Policy aggregate
- [ ] Payment integration
- [ ] Policy document generation

### Phase 5 - Endorsement & Renewal 🚧
- [ ] Endorsement aggregate
- [ ] Renewal scheduler
- [ ] Grace period management

### Phase 6 - Claims & Fraud 🚧
- [ ] Claim aggregate
- [ ] 4-tier approval matrix
- [ ] Fraud integration
- [ ] Settlement

### Phase 7 - Commission & Refund 🚧
- [ ] Commission calculation
- [ ] Revenue share
- [ ] Refund calculation

### Phase 8 - Hardening 🚧
- [ ] Integration tests
- [ ] Load testing
- [ ] Security review
- [ ] Production deployment

## Contributing

See [POLISYNC_PLAN.md](POLISYNC_PLAN.md) for detailed implementation plan.

## License

Proprietary - LabAid InsureTech Platform
