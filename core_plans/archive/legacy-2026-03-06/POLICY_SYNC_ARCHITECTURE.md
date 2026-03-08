# PoliSync Architecture

**Version:** 1.0.0  
**Last Updated:** March 3, 2026

---

## System Overview

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         LabAid InsureTech Platform                       │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                           │
│  ┌──────────────┐      ┌──────────────┐      ┌──────────────┐          │
│  │   Gateway    │      │   InScore    │      │  PoliSync    │          │
│  │   (Go)       │◄────►│   (Go)       │◄────►│   (C#)       │          │
│  │   :50000     │      │   :50050+    │      │   :50120+    │          │
│  └──────────────┘      └──────────────┘      └──────────────┘          │
│         │                      │                      │                  │
│         └──────────────────────┴──────────────────────┘                  │
│                                │                                          │
│         ┌──────────────────────┴──────────────────────┐                  │
│         │                                              │                  │
│    ┌────▼────┐    ┌──────┐    ┌───────┐    ┌────────┐                  │
│    │ Postgres│    │ Redis│    │ Kafka │    │  JWT   │                  │
│    │  :5432  │    │ :6379│    │ :9092 │    │ Tokens │                  │
│    └─────────┘    └──────┘    └───────┘    └────────┘                  │
│                                                                           │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## PoliSync Internal Architecture

### Layered Architecture

```
┌─────────────────────────────────────────────────────────────────────────┐
│                          PoliSync.ApiHost                                 │
│                    (Single Kestrel Host - 7 Services)                     │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                           │
│  ┌─────────────────────────────────────────────────────────────────┐    │
│  │                    gRPC Services Layer                           │    │
│  ├─────────────────────────────────────────────────────────────────┤    │
│  │  ProductGrpcService    QuoteGrpcService    OrderGrpcService     │    │
│  │  PolicyGrpcService     ClaimGrpcService    CommissionGrpcService│    │
│  │  UnderwritingGrpcService                                        │    │
│  └─────────────────────────────────────────────────────────────────┘    │
│                                │                                          │
│  ┌─────────────────────────────▼──────────────────────────────────┐    │
│  │                    Interceptors Layer                           │    │
│  ├─────────────────────────────────────────────────────────────────┤    │
│  │  JwtAuthInterceptor  │  LoggingInterceptor  │  ValidationInterceptor│
│  └─────────────────────────────────────────────────────────────────┘    │
│                                │                                          │
│  ┌─────────────────────────────▼──────────────────────────────────┐    │
│  │                    MediatR (CQRS Bus)                           │    │
│  └─────────────────────────────┬──────────────────────────────────┘    │
│                                │                                          │
│         ┌──────────────────────┴──────────────────────┐                  │
│         │                                              │                  │
│  ┌──────▼──────┐                              ┌───────▼──────┐          │
│  │  Commands   │                              │   Queries    │          │
│  │  (Write)    │                              │   (Read)     │          │
│  └──────┬──────┘                              └───────┬──────┘          │
│         │                                              │                  │
│  ┌──────▼──────────────────────────────────────────────▼──────┐        │
│  │                    Application Layer                         │        │
│  ├──────────────────────────────────────────────────────────────┤        │
│  │  Command Handlers  │  Query Handlers  │  Event Handlers     │        │
│  │  Validators        │  DTOs            │  Mappings           │        │
│  └──────────────────────────┬───────────────────────────────────┘        │
│                             │                                             │
│  ┌──────────────────────────▼───────────────────────────────────┐        │
│  │                    Domain Layer                               │        │
│  ├──────────────────────────────────────────────────────────────┤        │
│  │  Aggregates  │  Entities  │  Value Objects  │  Domain Events │        │
│  │  Business Rules  │  State Machines  │  Invariants           │        │
│  └──────────────────────────┬───────────────────────────────────┘        │
│                             │                                             │
│  ┌──────────────────────────▼───────────────────────────────────┐        │
│  │                    Infrastructure Layer                       │        │
│  ├──────────────────────────────────────────────────────────────┤        │
│  │  EF Core DbContext  │  Repositories  │  Unit of Work        │        │
│  │  Kafka EventBus     │  Redis Cache   │  PII Encryptor       │        │
│  │  gRPC Clients       │  CurrentUser   │  Configurations      │        │
│  └──────────────────────────┬───────────────────────────────────┘        │
│                             │                                             │
│  ┌──────────────────────────▼───────────────────────────────────┐        │
│  │                    External Systems                           │        │
│  ├──────────────────────────────────────────────────────────────┤        │
│  │  PostgreSQL  │  Redis  │  Kafka  │  Go Services (11)        │        │
│  └──────────────────────────────────────────────────────────────┘        │
│                                                                           │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## CQRS Flow

### Command Flow (Write Operations)

```
┌──────────┐     ┌──────────┐     ┌──────────┐     ┌──────────┐
│  gRPC    │────►│ Command  │────►│ Command  │────►│ Domain   │
│ Request  │     │          │     │ Handler  │     │ Aggregate│
└──────────┘     └──────────┘     └──────────┘     └────┬─────┘
                                                         │
                                                         │ Raise
                                                         │ Domain
                                                         │ Events
                                                         │
                                  ┌──────────┐     ┌────▼─────┐
                                  │   Kafka  │◄────│ Event    │
                                  │  Topics  │     │ Handlers │
                                  └──────────┘     └──────────┘
                                       │
                                       │ Publish
                                       │
                        ┌──────────────┴──────────────┐
                        │                             │
                   ┌────▼────┐                  ┌─────▼────┐
                   │ Go      │                  │ Other    │
                   │ Services│                  │ C#       │
                   │         │                  │ Handlers │
                   └─────────┘                  └──────────┘
```

### Query Flow (Read Operations)

```
┌──────────┐     ┌──────────┐     ┌──────────┐     ┌──────────┐
│  gRPC    │────►│  Query   │────►│  Query   │────►│  Redis   │
│ Request  │     │          │     │ Handler  │     │  Cache   │
└──────────┘     └──────────┘     └──────────┘     └────┬─────┘
                                                         │
                                                         │ Cache
                                                         │ Miss
                                                         │
                                                    ┌────▼─────┐
                                                    │ DbContext│
                                                    │ (Read)   │
                                                    └────┬─────┘
                                                         │
                                                         │ Project
                                                         │ to DTO
                                                         │
                                                    ┌────▼─────┐
                                                    │ Response │
                                                    └──────────┘
```

---

## Domain Model

### Bounded Contexts

```
┌─────────────────────────────────────────────────────────────────────────┐
│                          PoliSync Domain Model                            │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                           │
│  ┌──────────────────┐  ┌──────────────────┐  ┌──────────────────┐      │
│  │    Products      │  │     Quotes       │  │     Orders       │      │
│  ├──────────────────┤  ├──────────────────┤  ├──────────────────┤      │
│  │ • Product        │  │ • Quotation      │  │ • Order          │      │
│  │ • ProductPlan    │  │ • Premium Calc   │  │ • Payment Flow   │      │
│  │ • PricingConfig  │  │ • Expiry (30d)   │  │ • 5 Steps        │      │
│  │ • Rider          │  │                  │  │                  │      │
│  └──────────────────┘  └──────────────────┘  └──────────────────┘      │
│                                                                           │
│  ┌──────────────────┐  ┌──────────────────┐  ┌──────────────────┐      │
│  │     Policy       │  │   Endorsement    │  │    Renewal       │      │
│  ├──────────────────┤  ├──────────────────┤  ├──────────────────┤      │
│  │ • Policy         │  │ • Endorsement    │  │ • Schedule       │      │
│  │ • Nominee (PII)  │  │ • Approval       │  │ • Reminder       │      │
│  │ • State Machine  │  │ • Policy Diff    │  │ • Grace Period   │      │
│  │ • LP-YYYY-SEQ    │  │ • Pro-rata       │  │ • Auto-lapse     │      │
│  └──────────────────┘  └──────────────────┘  └──────────────────┘      │
│                                                                           │
│  ┌──────────────────┐  ┌──────────────────┐  ┌──────────────────┐      │
│  │  Underwriting    │  │     Claims       │  │   Commission     │      │
│  ├──────────────────┤  ├──────────────────┤  ├──────────────────┤      │
│  │ • Health Decl    │  │ • Claim (FNOL)   │  │ • Config         │      │
│  │ • Risk Score     │  │ • 4-Tier Approval│  │ • Payout         │      │
│  │ • Loading Factor │  │ • Fraud Check    │  │ • Revenue Share  │      │
│  │ • Decision       │  │ • Settlement     │  │ • Tax (10%)      │      │
│  └──────────────────┘  └──────────────────┘  └──────────────────┘      │
│                                                                           │
│  ┌──────────────────┐                                                    │
│  │     Refund       │                                                    │
│  ├──────────────────┤                                                    │
│  │ • Calculation    │                                                    │
│  │ • Pro-rata       │                                                    │
│  │ • Penalty (10%)  │                                                    │
│  └──────────────────┘                                                    │
│                                                                           │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## Event Flow

### Policy Issuance Event Chain

```
Order.PaymentConfirmed
        │
        ▼
Policy.Issue()
        │
        ├──► PolicyIssuedEvent
        │         │
        │         ├──► CommissionHandler ──► Calculate Payout
        │         │
        │         ├──► DocgenGrpcClient ──► Generate PDF
        │         │         │
        │         │         └──► StorageGrpcClient ──► Upload
        │         │
        │         ├──► NotificationService ──► SMS/Email/Push
        │         │
        │         ├──► RenewalScheduler ──► Create Schedule
        │         │
        │         └──► AuditService ──► Log Event
        │
        └──► Kafka: insuretech.policy.issued.v1
                  │
                  ├──► Analytics Service
                  ├──► Reporting Service
                  └──► External Systems
```

### Claim Filing Event Chain

```
Claim.File()
        │
        ├──► ClaimFiledEvent
        │         │
        │         ├──► FraudGrpcClient ──► Check Fraud Score
        │         │         │
        │         │         └──► If score > 0.75 ──► Flag for Investigation
        │         │
        │         ├──► NotificationService ──► Notify Customer
        │         │
        │         └──► WorkflowService ──► Assign Reviewer
        │
        └──► Kafka: insuretech.claim.filed.v1
                  │
                  └──► Analytics Service
```

---

## Data Flow

### Database Schema Relationships

```
insurance_schema
├── products
│   ├── product_plans (FK: product_id)
│   ├── pricing_configs (FK: product_id)
│   └── riders (FK: product_id)
│
├── quotations
│   ├── FK: product_id
│   ├── FK: plan_id
│   └── health_declarations (FK: quotation_id)
│
├── orders
│   ├── FK: quotation_id
│   └── FK: policy_id (after issuance)
│
├── policies
│   ├── FK: product_id, plan_id, quotation_id, order_id
│   ├── nominees (FK: policy_id) [PII ENCRYPTED]
│   ├── policy_riders (FK: policy_id, rider_id)
│   ├── endorsements (FK: policy_id)
│   └── renewal_schedules (FK: policy_id)
│
└── claims
    ├── FK: policy_id
    ├── claim_documents (FK: claim_id)
    ├── claim_approvals (FK: claim_id)
    └── settlements (FK: claim_id, policy_id)

commission_schema
├── commission_configs
│   └── FK: partner_id, product_id
│
├── commission_payouts
│   └── FK: partner_id, policy_id
│
└── revenue_shares
    └── FK: partner_id
```

---

## Security Architecture

### Authentication Flow

```
┌──────────┐     ┌──────────┐     ┌──────────┐     ┌──────────┐
│  Client  │────►│ Gateway  │────►│  Authn   │────►│   JWT    │
│          │     │  :50000  │     │ Service  │     │  Token   │
└──────────┘     └──────────┘     │  :50060  │     └────┬─────┘
                                   └──────────┘          │
                                                         │
                                                         │ Include
                                                         │ in Header
                                                         │
┌──────────┐     ┌──────────┐     ┌──────────┐     ┌───▼──────┐
│ PoliSync │◄────│  JWT     │◄────│ Gateway  │◄────│  Client  │
│ Services │     │ Verify   │     │          │     │          │
└──────────┘     └──────────┘     └──────────┘     └──────────┘
     │
     │ Extract Claims
     │
     ▼
┌──────────┐
│ Current  │
│  User    │
│ Context  │
└──────────┘
```

### Authorization Flow

```
┌──────────┐     ┌──────────┐     ┌──────────┐
│  gRPC    │────►│  Authz   │────►│  RBAC    │
│ Handler  │     │ Client   │     │  Check   │
└──────────┘     └──────────┘     └────┬─────┘
                                        │
                                        │ Roles
                                        │ Permissions
                                        │
                                   ┌────▼─────┐
                                   │  ABAC    │
                                   │  Check   │
                                   └────┬─────┘
                                        │
                                        │ Tenant
                                        │ Partner
                                        │ Context
                                        │
                                   ┌────▼─────┐
                                   │ Allow /  │
                                   │  Deny    │
                                   └──────────┘
```

### PII Encryption

```
┌──────────┐     ┌──────────┐     ┌──────────┐
│ Nominee  │────►│   AES    │────►│ Database │
│   NID    │     │ 256-GCM  │     │ (cipher) │
│  Phone   │     │ Encrypt  │     └──────────┘
└──────────┘     └──────────┘
                      │
                      │ 32-byte key
                      │ from secret
                      │
                 ┌────▼─────┐
                 │  Secret  │
                 │  Manager │
                 └──────────┘
```

---

## Deployment Architecture

### Kubernetes Deployment

```
┌─────────────────────────────────────────────────────────────────────────┐
│                          Kubernetes Cluster                               │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                           │
│  ┌──────────────────────────────────────────────────────────────────┐   │
│  │                    Ingress Controller                             │   │
│  │                    (NGINX / Traefik)                              │   │
│  └────────────────────────────┬─────────────────────────────────────┘   │
│                               │                                           │
│  ┌────────────────────────────▼─────────────────────────────────────┐   │
│  │                    Gateway Service                                │   │
│  │                    (Go - :50000)                                  │   │
│  │                    Replicas: 3                                    │   │
│  └────────────────────────────┬─────────────────────────────────────┘   │
│                               │                                           │
│         ┌─────────────────────┴─────────────────────┐                    │
│         │                                            │                    │
│  ┌──────▼──────┐                              ┌─────▼──────┐            │
│  │  InScore    │                              │ PoliSync   │            │
│  │  Services   │                              │ Services   │            │
│  │  (Go)       │                              │ (C#)       │            │
│  │  Replicas:3 │                              │ Replicas:3 │            │
│  └──────┬──────┘                              └─────┬──────┘            │
│         │                                            │                    │
│         └────────────────────┬───────────────────────┘                    │
│                              │                                            │
│  ┌───────────────────────────▼────────────────────────────────────┐     │
│  │                    StatefulSets                                 │     │
│  ├─────────────────────────────────────────────────────────────────┤     │
│  │  PostgreSQL (Primary + Replicas)                                │     │
│  │  Redis (Cluster Mode)                                           │     │
│  │  Kafka (3 Brokers + Zookeeper)                                  │     │
│  └─────────────────────────────────────────────────────────────────┘     │
│                                                                           │
│  ┌─────────────────────────────────────────────────────────────────┐     │
│  │                    Observability Stack                           │     │
│  ├─────────────────────────────────────────────────────────────────┤     │
│  │  Prometheus  │  Grafana  │  Jaeger  │  Loki  │  AlertManager   │     │
│  └─────────────────────────────────────────────────────────────────┘     │
│                                                                           │
└─────────────────────────────────────────────────────────────────────────┘
```

### Horizontal Pod Autoscaling

```
┌──────────────────────────────────────────────────────────────┐
│                    HPA Configuration                          │
├──────────────────────────────────────────────────────────────┤
│                                                               │
│  PoliSync Services:                                          │
│  • Min Replicas: 2                                           │
│  • Max Replicas: 10                                          │
│  • Target CPU: 70%                                           │
│  • Target Memory: 80%                                        │
│  • Custom Metric: gRPC connections > 1000                   │
│                                                               │
│  Scale Up:                                                   │
│  • Add 1 pod every 30 seconds                               │
│  • Max 3 pods per scale event                               │
│                                                               │
│  Scale Down:                                                 │
│  • Remove 1 pod every 5 minutes                             │
│  • Stabilization window: 5 minutes                          │
│                                                               │
└──────────────────────────────────────────────────────────────┘
```

---

## Performance Characteristics

### Latency Targets

```
┌──────────────────────────────────────────────────────────────┐
│                    Latency Targets (p95)                      │
├──────────────────────────────────────────────────────────────┤
│                                                               │
│  CalculatePremium      ████████░░░░░░░░░░  < 200ms          │
│  GetProduct (cached)   ███░░░░░░░░░░░░░░░  < 100ms          │
│  GetPolicy             ████████░░░░░░░░░░░  < 150ms          │
│  IssuePolicy (sync)    ████████████████░░░  < 500ms          │
│  FileClaim             ██████████░░░░░░░░░  < 300ms          │
│  ListPolicies          ████████████░░░░░░░  < 400ms          │
│                                                               │
└──────────────────────────────────────────────────────────────┘
```

### Throughput Targets

```
┌──────────────────────────────────────────────────────────────┐
│                    Throughput Targets                         │
├──────────────────────────────────────────────────────────────┤
│                                                               │
│  Quote Requests:       1,000 concurrent                      │
│  Policy Issuances:     500 per minute                        │
│  Claim Submissions:    200 per hour                          │
│  Premium Calculations: 2,000 per minute                      │
│                                                               │
└──────────────────────────────────────────────────────────────┘
```

---

## Technology Stack Summary

```
┌─────────────────────────────────────────────────────────────────────────┐
│                          Technology Stack                                 │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                           │
│  Runtime:          .NET 8 LTS                                            │
│  Language:         C# 12                                                 │
│  Communication:    gRPC (Grpc.AspNetCore 2.60)                          │
│  Database:         PostgreSQL 15 + EF Core 8 + Npgsql                   │
│  Caching:          Redis 7 + StackExchangeRedis                         │
│  Messaging:        Kafka 3.x + Confluent.Kafka 2.3                      │
│  CQRS:             MediatR 12                                            │
│  Validation:       FluentValidation 11                                   │
│  Logging:          Serilog 3.1                                           │
│  Tracing:          OpenTelemetry 1.7                                     │
│  Metrics:          Prometheus                                            │
│  Testing:          xUnit + Moq + Testcontainers                         │
│  Container:        Docker + Kubernetes                                   │
│                                                                           │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## Design Principles

1. **Domain-Driven Design**
   - Bounded contexts
   - Aggregates with invariants
   - Domain events
   - Ubiquitous language

2. **CQRS**
   - Command/Query separation
   - Different models for read/write
   - Event sourcing ready

3. **Event-Driven**
   - Domain events
   - Kafka integration
   - Eventual consistency

4. **Clean Architecture**
   - Dependency inversion
   - Infrastructure at edges
   - Domain at center

5. **Microservices**
   - Single responsibility
   - Independent deployment
   - Shared nothing (except DB schemas)

6. **Security First**
   - JWT authentication
   - PII encryption
   - Audit trail
   - Tenant isolation

7. **Observability**
   - Structured logging
   - Distributed tracing
   - Metrics
   - Health checks

8. **Testability**
   - Unit tests
   - Integration tests
   - Load tests
   - Testcontainers

---

**For implementation details, see:**
- [POLISYNC_REFERENCE.md](POLISYNC_REFERENCE.md) - Complete reference, status, state machines & implementation guide
- [QUICKSTART.md](QUICKSTART.md) - Getting started in 5 minutes
