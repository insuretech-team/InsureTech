# InsureTech Microservice Architecture Guide

## Overview

The InsureTech backend is a **hybrid polyglot microservices platform** consisting of:
- **Go (InScore)**: Core infrastructure, data layers, and foundational services
- **C# (PoliSync)**: Business logic, domain models, and policy-related services
- **.NET (Ledger Claw)**: Financial ledger and accounting services

All services communicate via **gRPC** for inter-service calls and **Kafka** for event-driven architecture.

---

## Architecture Layers

### 1. **API Gateway Layer** (HTTP Entry Point)
- **Service**: `gateway` (Port 8080)
- **Backend**: Go (InScore)
- **Purpose**: HTTP-to-gRPC translation, routing, authentication
- **Pattern**: Acts as single entry point for all HTTP clients

### 2. **Core Infrastructure Services** (Go - InScore)
These provide foundational capabilities:

| Service | Port | Purpose | Key Features |
|---------|------|---------|--------------|
| `tenant` | 50050 | Organization & config management | Multi-tenancy |
| `authn` | 50060 | Identity, sessions, OTP | Hybrid auth (email, SMS, OTP) |
| `authz` | 50070 | RBAC/ABAC permissions | Casbin-based policy enforcement |
| `audit` | 50080 | Centralized audit logging | Compliance tracking |
| `kyc` | 50090 | Identity verification orchestrator | External KYC provider integration |
| `workflow` | 50180 | Business process orchestration | State machine engine |

### 3. **User & Partner Management** (Go - InScore)

| Service | Port | Purpose |
|---------|------|---------|
| `partner` | 50100 | Agency/broker portal |
| `beneficiary` | 50110 | Beneficiary & nominee mgmt |
| `b2b` | 50112 | Employee & department mgmt |

### 4. **Data Layer** (Go - InScore)

| Service | Port | Purpose | Notes |
|---------|------|---------|-------|
| `insurance` | 50115 | Insurance schema CRUD | Central data layer for all insurance entities |
| `orders` | 50142 | Order lifecycle data layer | Go layer for PoliSync C# orders service |

### 5. **Business Logic & Commerce** (C# - PoliSync)

| Service | Port | Purpose |
|---------|------|---------|
| `product` | 50120 | Product factory, pricing |
| `quote` | 50130 | Quotation management |
| `order` | 50140 | Cart, checkout, order lifecycle |
| `commission` | 50150 | Revenue share & payouts |
| `policy` | 50160 | Core policy management |
| `underwriting` | 50170 | Risk assessment & decisions |

### 6. **Specialized Services** (Go - InScore)

**Communications & Media**:
- `notification` (50230): SMS, email, push
- `support` (50240): Ticketing & customer support
- `webrtc` (50250): Video signaling
- `media` (50260): Video/image transcoding
- `ocr` (50270): Optical character recognition

**Documents & Storage**:
- `docgen` (50280): PDF generation
- `storage` (50290): S3/blob storage wrapper

**Intelligence**:
- `iot` (50300): Telemetry ingestion
- `analytics` (50310): Reporting & insights
- `ai` (50320): LLM agents & copilot

**Fraud & Financials**:
- `fraud` (50220): Real-time fraud detection
- `payment` (50190): Payment gateway integration
- `ledger` (50200): Double-entry accounting (Ledger Claw)

---

## Directory Structure

```
backend/
├── inscore/                          # Go monorepo
│   ├── cmd/                          # Entry points (wrapper scripts)
│   │   ├── authn/main.go            # Authn wrapper
│   │   ├── authz/main.go            # Authz wrapper
│   │   ├── orders/main.go           # Orders wrapper
│   │   ├── gateway/main.go          # Gateway entry
│   │   ├── dbmanager/               # Database management CLI
│   │   └── dbops/                   # DB operations (migrations, sync)
│   ├── microservices/                # Actual service implementations
│   │   ├── authn/
│   │   │   ├── cmd/server/main.go   # Real authn service
│   │   │   └── internal/
│   │   │       ├── config/          # Configuration
│   │   │       ├── service/         # Business logic
│   │   │       ├── repository/      # Data access
│   │   │       ├── grpc/            # gRPC handlers
│   │   │       ├── domain/          # Domain interfaces
│   │   │       ├── events/          # Event publishing
│   │   │       ├── consumers/       # Event handling
│   │   │       └── middleware/      # Middleware
│   │   ├── orders/
│   │   │   └── internal/
│   │   │       ├── config/
│   │   │       ├── service/
│   │   │       ├── repository/
│   │   │       ├── grpc/
│   │   │       ├── domain/
│   │   │       ├── events/
│   │   │       └── consumers/
│   │   └── [other services]/
│   ├── db/                           # Shared database layer
│   │   ├── config.go                # Database configuration
│   │   ├── db.go                    # Manager & connection logic
│   │   ├── manager.go               # Database manager
│   │   ├── migrations/              # SQL migrations by schema
│   │   │   ├── authn_schema/
│   │   │   ├── insurance_schema/
│   │   │   ├── payment_schema/
│   │   │   └── [other schemas]/
│   │   └── seeds/                   # Seed data
│   ├── pkg/                          # Shared packages
│   │   ├── logger/                  # Logging
│   │   ├── kafka/                   # Kafka producer/consumer
│   │   ├── crypto/                  # Encryption utilities
│   │   ├── interceptors/            # gRPC interceptors
│   │   ├── webrtc/                  # WebRTC support
│   │   └── [other utilities]/
│   └── configs/                      # Configuration files
│       ├── services.yaml            # Service registry & ports
│       ├── database.yaml            # Database config
│       ├── kafka.yaml               # Kafka topics
│       └── [other configs]/
├── polisync/                        # C# monorepo (PoliSync)
│   ├── src/
│   │   ├── PoliSync.ApiHost/       # API host & composition
│   │   ├── PoliSync.Products/      # Product domain
│   │   ├── PoliSync.Policy/        # Policy domain
│   │   ├── PoliSync.Claims/        # Claims domain
│   │   ├── PoliSync.Underwriting/  # Underwriting domain
│   │   ├── PoliSync.Orders/        # Orders business logic
│   │   ├── PoliSync.Infrastructure/ # Shared infrastructure
│   │   │   ├── Clients/            # gRPC client factories
│   │   │   ├── GrpcClients/        # Generated gRPC stubs
│   │   │   ├── Persistence/        # Database repositories
│   │   │   ├── Messaging/          # Event bus (Kafka)
│   │   │   └── Pii/                # PII encryption
│   │   └── PoliSync.SharedKernel/  # Shared domain models
│   └── tests/
└── ledgerclaw/                      # Financial ledger service
    └── [ledger implementation]/
```

---

## Microservice Internal Structure

Every Go microservice follows this consistent pattern:

### Standard Directory Layout

```
microservices/{service}/
├── cmd/
│   └── server/
│       └── main.go                 # Service startup
├── internal/
│   ├── config/
│   │   └── config.go              # Configuration loading
│   ├── domain/
│   │   └── interfaces.go           # Business interfaces
│   ├── service/
│   │   └── {service}_service.go   # Business logic
│   ├── repository/
│   │   └── repository.go           # Data access (GORM)
│   ├── grpc/
│   │   ├── server.go              # gRPC server setup
│   │   ├── {service}_handler.go   # gRPC handlers
│   │   └── errors.go              # Error mapping
│   ├── events/
│   │   ├── publisher.go           # Event publishing
│   │   └── topics.go              # Kafka topics
│   ├── consumers/
│   │   └── consumer.go            # Event handling
│   └── middleware/
│       └── [middleware files]/
```

### Execution Flow

```
cmd/{service}/main.go (wrapper)
    ↓
    └─ Resolves project root
    └─ Loads .env
    └─ Optionally runs DB migrations
    └─ Starts DB sync sidecar (if needed)
    └─ Delegates to actual service

microservices/{service}/cmd/server/main.go (actual service)
    ↓
    ├─ 1. Initialize Logger
    ├─ 2. Load Configuration (services.yaml, database.yaml)
    ├─ 3. Initialize Database (shared db.Manager)
    ├─ 4. Initialize Kafka (producer + consumer)
    ├─ 5. Create Repository (GORM database access)
    ├─ 6. Create Service (business logic)
    ├─ 7. Create gRPC Server & Handlers
    ├─ 8. Register gRPC Services
    ├─ 9. Start Listening
    └─ 10. Shutdown on signal (graceful)
```

---

## Key Patterns & Conventions

### 1. **Dependency Injection**
Each layer depends on the layer below (top-down):
```
Handler (gRPC) → Service (Business) → Repository (Data)
```

### 2. **Domain-Driven Design**
- `domain/interfaces.go`: Defines contracts (Repository, Service)
- `service/`: Implements business rules
- `repository/`: Implements data persistence
- `grpc/`: Adapts to external protocol

### 3. **Configuration Management**
Services load configuration in this order:
1. Load `.env` file (via `ops/env.Load()`)
2. Read `services.yaml` for service ports
3. Read `database.yaml` for database connections
4. Override with environment variables (optional)

### 4. **Error Handling**
- Define custom errors in `service/errors.go`
- Map to gRPC status codes in `grpc/errors.go`
- Example:
  ```go
  // service layer
  return fmt.Errorf("invalid order: %w", ErrInvalidArgument)
  
  // handler layer (maps to gRPC)
  case errors.Is(err, service.ErrInvalidArgument):
      return status.Error(codes.InvalidArgument, err.Error())
  ```

### 5. **Event-Driven Architecture**
Services publish and consume events via Kafka:
```go
// Publishing
publisher.PublishOrderCreated(ctx, order)

// Consuming
consumer := NewEventConsumer()
kafkaConsumer.Start(ctx) // Handles payment.completed, payment.failed
```

---

## Database Architecture

### Multi-Database Setup
The system supports **primary-backup failover**:

```yaml
database:
  primary:
    provider: "digitalocean"  # Production database
  backup:
    provider: "neon"          # Failover database
  failover:
    enabled: true
    health_check_interval: 5s
  sync:
    enabled: true
    interval: 15m  # Keep databases synchronized
```

### Configuration Files
- **database.yaml**: Connection strings, pool sizes, failover settings
- **services.yaml**: Service registry (name, ports, backend)
- **kafka.yaml**: Kafka brokers and topic definitions

### Schema Organization
Migrations are organized by schema (domain):
```
db/migrations/
├── authn_schema/          # Authentication tables
├── authz_schema/          # Authorization tables
├── insurance_schema/      # Insurance domain tables
├── payment_schema/        # Payment domain tables
├── order_schema/          # Order domain tables
└── [other schemas]/
```

Each migration file follows naming:
```
{YYYYMMDD}_{sequence}_{description}.up.sql
```

### Shared Database Manager
All Go services use `db.Manager` (singleton):
```go
// Initialize once
db.InitializeManagerForService(dbConfigPath)

// Use anywhere
database := db.GetDB()  // Returns *gorm.DB

// Automatic failover & health checks happen transparently
```

---

## Service Communication Patterns

### 1. **gRPC Service-to-Service**
Direct synchronous calls between services:
```go
// C# calling Go service
using var channel = GrpcChannel.ForAddress("http://localhost:50115");
var client = new InsuranceService.InsuranceServiceClient(channel);
var response = await client.CreateProductAsync(request);
```

### 2. **Kafka Event Publishing**
Async event-driven communication:
```go
// Publisher (Orders service)
publisher.PublishOrderCreated(ctx, order)

// Subscriber (Payment service listening to orders.created)
consumer.HandleOrderCreated(ctx, orderData)
```

### 3. **Health Checks**
Services expose gRPC health check endpoints:
```
grpc.health.v1.Health/Check
```

---

## Configuration Files Reference

### services.yaml
```yaml
services:
  orders:
    name: orders-service
    backend: inscore
    description: "Order Lifecycle (Go data layer for PoliSync orders)"
    ports:
      grpc: 50142
      http: 50143
```

### database.yaml
```yaml
database:
  primary:
    provider: "digitalocean"
    host: "${PGHOST}"
    port: "${PGPORT}"
    database: "${PGDATABASE}"
    username: "${PGUSER}"
    password: "${PGPASSWORD}"
    ssl_mode: "${PGSSLMODE}"
    max_open_conns: 15
    max_idle_conns: 5
    conn_max_lifetime: "30m"
```

---

## Common Service Startup Checklist

When implementing a new microservice, follow this pattern:

### 1. **Create Entry Point**
```
cmd/{service}/main.go  # Wrapper script
microservices/{service}/cmd/server/main.go  # Actual service
```

### 2. **Define Domain Interfaces**
```go
// internal/domain/interfaces.go
type {Entity}Repository interface { ... }
type {Service}Service interface { ... }
```

### 3. **Implement Repository**
```go
// internal/repository/repository.go
type {Entity}RepositoryImpl struct { db *gorm.DB }
```

### 4. **Implement Service**
```go
// internal/service/{service}_service.go
type {Service}ServiceImpl struct { repo domain.{Entity}Repository }
```

### 5. **Create gRPC Handler**
```go
// internal/grpc/{service}_handler.go
type {Service}Handler struct { svc domain.{Service}Service }
```

### 6. **Wire Everything in main.go**
```go
repo := repository.New{Entity}Repository(db)
svc := service.New{Service}Service(repo)
handler := grpc.New{Service}Handler(svc)
orderservicev1.Register{Service}Server(grpcServer, handler)
```

---

## Environment Variables

### Standard Variables
```bash
# Database
PGHOST=xxx.cloud.digitalocean.com
PGPORT=25060
PGDATABASE=insuretech
PGUSER=user
PGPASSWORD=pass
PGSSLMODE=require

# Failover (optional)
PGHOST2=xxx.neon.tech
PGPORT2=5432
PGDATABASE2=insuretech
PGUSER2=user
PGPASSWORD2=pass
PGSSLMODE2=require

# Kafka
KAFKA_BROKERS=localhost:9092

# Service-specific
AUTHN_PORT=50060
AUTHZ_PORT=50070
# ... etc
```

---

## Debugging & Operations

### View Service Status
```powershell
# Check if service is running
grpcurl -plaintext localhost:50142 list

# Health check
grpcurl -plaintext localhost:50142 grpc.health.v1.Health/Check
```

### Database Operations
```powershell
cd backend/inscore

# Run migrations
go run ./cmd/dbmanager migrate --target=primary

# Check migration status
go run ./cmd/dbmanager sql --sql="SELECT * FROM schema_migrations" --target=primary

# Export data
go run ./cmd/dbmanager csv-backup --table=users --source=primary
```

### Start Service
```powershell
# From project root
cd backend/inscore
go run ./cmd/{service}/main.go

# With migrations
AUTHN_RUN_MIGRATIONS=true go run ./cmd/authn/main.go
```

---

## Design Principles

1. **Separation of Concerns**: Each layer has a single responsibility
2. **Proto-First**: All APIs are defined in `.proto` files (single source of truth)
3. **Interface-Based**: Components depend on interfaces, not implementations
4. **Shared Infrastructure**: Common packages in `pkg/` avoid duplication
5. **Graceful Shutdown**: All services handle `SIGTERM` and drain connections
6. **Health Checks**: Services expose liveness and readiness probes
7. **Event-Driven**: Loose coupling via Kafka for cross-service communication
8. **Database Resilience**: Failover & sync for data consistency

---

## Next Steps for New Services

1. Add service configuration to `configs/services.yaml`
2. Create database migrations in `db/migrations/{schema}/`
3. Generate proto files: `buf generate`
4. Create service structure following the standard pattern
5. Wire dependencies in `cmd/server/main.go`
6. Add Dockerfile to `infra/docker/{service}/`
7. Register with docker-compose if needed

