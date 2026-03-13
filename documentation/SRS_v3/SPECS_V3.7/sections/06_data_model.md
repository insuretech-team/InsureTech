# 6. Data Model & Persistence

## 6.1 Proto Schema Organization

All data models are defined using Protocol Buffers (proto3) and organized by domain with a consistent structure:

```
proto/insuretech/
├── authn/                          Authentication Domain
│   ├── entity/v1/                  Data entities
│   │   ├── user.proto
│   │   └── session.proto
│   ├── events/v1/                  Domain events
│   │   └── auth_events.proto
│   └── services/v1/                gRPC services
│       └── auth_service.proto
│
├── authz/                          Authorization Domain
│   ├── entity/v1/
│   ├── events/v1/
│   └── services/v1/
│
├── policy/                         Policy Management Domain
│   ├── entity/v1/
│   │   └── policy.proto
│   ├── events/v1/
│   └── services/v1/
│
├── claims/                         Claims Processing Domain
│   ├── entity/v1/
│   │   └── claim.proto
│   ├── events/v1/
│   └── services/v1/
│
├── payment/                        Payment Processing Domain
│   ├── entity/v1/
│   │   └── payment.proto
│   ├── events/v1/
│   └── services/v1/
│
├── partner/                        Partner Management Domain
│   ├── entity/v1/
│   │   └── partner.proto
│   ├── events/v1/
│   └── services/v1/
│
├── products/                       Product Catalog Domain
│   ├── entity/v1/
│   ├── events/v1/
│   └── services/v1/
│
├── notification/                   Notification Domain
│   ├── entity/v1/
│   ├── events/v1/
│   └── services/v1/
│
├── ai/                            AI Engine Domain
│   ├── entity/v1/
│   ├── events/v1/
│   └── services/v1/
│
├── analytics/                      Analytics & BI Domain
│   ├── entity/v1/
│   ├── events/v1/
│   └── services/v1/
│
└── iot/                           IoT Integration Domain
    ├── entity/v1/
    ├── events/v1/
    └── services/v1/
```

**See Appendix A for complete proto definitions with code examples.**

## 6.2 Data Architecture Overview

The LabAid InsureTech Platform uses a **hybrid data architecture** combining Protocol Buffers for service contracts with optimized database schemas for persistence:

**Data Storage Strategy:**
- **PostgreSQL 17+:** Primary transactional data (policies, claims, users, KYC)
- **TimescaleDB:** Time-series data (IoT telemetry, audit logs, analytics)
- **TigerBeetle:** Financial transactions with double-entry bookkeeping
- **DynamoDB:** Product catalog, configuration, session data
- **Redis 7.0+:** Caching, session management, real-time data
- **AWS S3:** Document storage (policy certificates, claims documents, images)
- **Apache Kafka:** Event streaming and audit logs
- **Pgvector:** Vector embeddings for AI/ML operations

## 6.3 Domain Models

### 6.3.1 Core Entities

**Authentication Domain:**
- `User` - Registered users with mobile/email
- `Session` - Active user sessions with JWT tokens
- `OTP` - One-time passwords for verification

**Policy Domain:**
- `Policy` - Insurance policies with coverage details
- `Applicant` - Policyholder information
- `Nominee` - Beneficiaries with share percentages
- `Rider` - Additional coverage options

**Claims Domain:**
- `Claim` - Claim submissions with status tracking
- `ClaimDocument` - Supporting documents (bills, reports)
- `ClaimApproval` - Multi-level approval workflow
- `FraudCheck` - Fraud detection results

**Payment Domain:**
- `Payment` - Financial transactions
- `Transaction` - Double-entry accounting records
- `Refund` - Refund processing

**Partner Domain:**
- `Partner` - Business partners (hospitals, MFS, e-commerce)
- `Agent` - Sales representatives
- `Commission` - Commission structure and calculations

## 6.4 Proto-First Data Model Strategy

### 6.4.1 Why Proto-First?

**Single Source of Truth:**
The LabAid InsureTech Platform adopts a **Proto-First approach** where Protocol Buffer definitions serve as the canonical data model across all layers:

```
Proto Definitions (Source of Truth)
    ├── Code Generation → Go/C#/Python/Node.js structs
    ├── Database Schema → PostgreSQL/TimescaleDB tables
    ├── API Contracts → gRPC/REST endpoints
    ├── Event Schemas → Kafka message formats
    └── Documentation → Auto-generated API docs
```

**Key Benefits:**
- ✅ **Type Safety:** Compile-time validation across all services
- ✅ **Consistency:** Same data structure in app, database, and APIs
- ✅ **Versioning:** Built-in support for backward/forward compatibility
- ✅ **Multi-Language:** Generate code for Go, C#, Python, Node.js from single source
- ✅ **Performance:** Efficient binary serialization with Protocol Buffers
- ✅ **Documentation:** Self-documenting with comments in proto files

### 6.4.2 Schema Generation Workflow

**Step 1: Define Proto Schemas**
```protobuf
// proto/insuretech/policy/entity/v1/policy.proto
message Policy {
  string policy_id = 1;                    // UUID
  string policy_number = 2;                // LBT-YYYY-XXXX-NNNNNN
  string customer_id = 3;
  PolicyStatus status = 4;
  double premium_amount = 5;
  google.protobuf.Timestamp created_at = 6;
}
```

**Step 2: Generate Database Schema**
```bash
# Using buf or protoc-gen-sql
buf generate

# Outputs: migrations/001_initial_schema.sql
CREATE TABLE policies (
    policy_id UUID PRIMARY KEY,
    policy_number VARCHAR(50) UNIQUE NOT NULL,
    customer_id UUID NOT NULL,
    status VARCHAR(50) NOT NULL,
    premium_amount DECIMAL(12,2) NOT NULL,
    created_at TIMESTAMP NOT NULL
);
```

**Step 3: Manual Optimization (Enhancement Layer)**
```sql
-- migrations/002_add_indexes.sql
CREATE INDEX idx_policies_customer_id ON policies(customer_id);
CREATE INDEX idx_policies_status ON policies(status);
CREATE INDEX idx_policies_created_at ON policies(created_at DESC);

-- Add foreign keys
ALTER TABLE policies ADD CONSTRAINT fk_customer 
    FOREIGN KEY (customer_id) REFERENCES users(user_id);

-- Add constraints
ALTER TABLE policies ADD CONSTRAINT chk_premium_positive 
    CHECK (premium_amount > 0);
```

**Step 4: Generate Application Code**
```bash
# Go
protoc --go_out=. --go-grpc_out=. proto/**/*.proto

# C#
protoc --csharp_out=. --grpc_out=. proto/**/*.proto

# Python
python -m grpc_tools.protoc -I. --python_out=. --grpc_python_out=. proto/**/*.proto
```

### 6.4.3 Database Strategy by Type

#### PostgreSQL (Primary Transactional Database)

**Generated Tables:**
```
From proto/insuretech/authn/entity/v1/*.proto:
  ├── users               (User proto → users table)
  ├── user_profiles       (UserProfile proto)
  ├── sessions            (Session proto)
  └── otps                (OTP proto)

From proto/insuretech/policy/entity/v1/*.proto:
  ├── policies            (Policy proto)
  ├── policy_nominees     (Nominee proto)
  └── policy_riders       (Rider proto)

From proto/insuretech/claims/entity/v1/*.proto:
  ├── claims              (Claim proto)
  ├── claim_documents     (ClaimDocument proto)
  ├── claim_approvals     (ClaimApproval proto)
  └── fraud_checks        (FraudCheckResult proto)

From proto/insuretech/payment/entity/v1/*.proto:
  └── payments            (Payment proto)

From proto/insuretech/partner/entity/v1/*.proto:
  ├── partners            (Partner proto)
  └── agents              (Agent proto)
```

**Enhancement Strategy:**
- **Indexes:** Add for frequently queried columns (customer_id, status, dates)
- **Foreign Keys:** Enforce referential integrity
- **Constraints:** Business rules (amount > 0, dates logical)
- **Triggers:** Audit logging, automatic timestamp updates
- **Partitioning:** For large tables (policies by year, claims by month)

#### TimescaleDB (Time-Series Database)

**Generated Hypertables:**
```
From proto/insuretech/iot/entity/v1/device.proto:
  └── telemetry          (Telemetry proto → hypertable on timestamp)

From proto/insuretech/authn/events/v1/*.proto:
  └── audit_logs         (Event protos → hypertable)

From proto/insuretech/analytics/entity/v1/*.proto:
  └── metrics            (Metric proto → hypertable)
```

**Enhancement Strategy:**
```sql
-- Create hypertable
SELECT create_hypertable('telemetry', 'timestamp');

-- Add continuous aggregates for dashboards
CREATE MATERIALIZED VIEW telemetry_hourly
WITH (timescaledb.continuous) AS
SELECT time_bucket('1 hour', timestamp) AS hour,
       device_id,
       AVG(metrics->>'speed') as avg_speed,
       COUNT(*) as data_points
FROM telemetry
GROUP BY hour, device_id;

-- Set retention policy (90 days hot, archive rest)
SELECT add_retention_policy('telemetry', INTERVAL '90 days');
```

#### DynamoDB (NoSQL Document Store)

**Generated Collections:**
```
From proto/insuretech/products/entity/v1/product.proto:
  └── products_catalog   (Product proto → DynamoDB items)

From proto/insuretech/authn/entity/v1/session.proto:
  └── active_sessions    (Session proto → TTL-enabled items)
```

**Key Design:**
- **Primary Key:** entity_id (from proto)
- **Sort Key:** timestamp or type (from proto fields)
- **TTL:** For session data (expires_at from proto)
- **Global Secondary Indexes:** Based on query patterns

#### TigerBeetle (Financial Ledger)

**Generated Accounts:**
```
From proto/insuretech/payment/entity/v1/payment.proto:
  Account types mapped from PaymentType enum:
  ├── Premium Collection Accounts
  ├── Claims Settlement Accounts
  ├── Commission Payment Accounts
  └── Refund Accounts
```

**Double-Entry Example:**
```
Premium Payment (BDT 1,500):
  Debit:  Customer Account        -1,500 BDT
  Credit: Premium Collection      +1,500 BDT
```

### 6.4.4 Migration Strategy

#### Phase 1: Initial Schema Generation
```bash
# Generate from all entity protos
protoc-gen-sql \
  --out=migrations/001_initial_schema.sql \
  proto/insuretech/*/entity/v1/*.proto

# Apply to database
psql -d insuretech_dev -f migrations/001_initial_schema.sql
```

#### Phase 2: Add Enhancements
```sql
-- migrations/002_add_indexes.sql
CREATE INDEX CONCURRENTLY idx_policies_customer_status 
    ON policies(customer_id, status);

-- migrations/003_add_foreign_keys.sql
ALTER TABLE policies ADD CONSTRAINT fk_customer ...;
ALTER TABLE claims ADD CONSTRAINT fk_policy ...;

-- migrations/004_add_partitioning.sql
CREATE TABLE policies_2025 PARTITION OF policies
    FOR VALUES FROM ('2025-01-01') TO ('2026-01-01');
```

#### Phase 3: Data Migration
```sql
-- migrations/005_migrate_legacy_data.sql
INSERT INTO users (user_id, mobile_number, ...)
SELECT uuid_generate_v4(), phone, ...
FROM legacy_customers;
```

#### Phase 4: Performance Optimization
```sql
-- migrations/006_optimize_queries.sql
CREATE MATERIALIZED VIEW policy_summary AS
SELECT customer_id, COUNT(*) as policy_count, SUM(premium_amount) as total_premium
FROM policies WHERE status = 'ACTIVE'
GROUP BY customer_id;

-- Refresh strategy
REFRESH MATERIALIZED VIEW CONCURRENTLY policy_summary;
```

### 6.4.5 Schema Versioning Strategy

**Proto Evolution:**
```protobuf
// v1 - Initial version
message Policy {
  string policy_id = 1;
  double premium_amount = 2;
}

// v2 - Add new field (backward compatible)
message Policy {
  string policy_id = 1;
  double premium_amount = 2;
  string partner_id = 3;        // New field - optional by default
}

// v3 - Deprecate field (forward compatible)
message Policy {
  string policy_id = 1;
  double premium_amount = 2;
  string partner_id = 3;
  double old_field = 4 [deprecated = true];  // Mark as deprecated
}
```

**Database Migration for Proto Changes:**
```sql
-- When adding field to proto
ALTER TABLE policies ADD COLUMN partner_id UUID;

-- When deprecating field
-- Keep column for backward compatibility, mark in comments
COMMENT ON COLUMN policies.old_field IS 'DEPRECATED: Use new_field instead';

-- After grace period, drop column
ALTER TABLE policies DROP COLUMN old_field;
```

### 6.4.6 Data Consistency Patterns

#### Strong Consistency (PostgreSQL)
```
Policy Creation:
  1. Begin Transaction
  2. Insert into policies
  3. Insert into policy_nominees
  4. Insert into policy_riders
  5. Commit (all or nothing)
```

#### Eventual Consistency (Event-Driven)
```
Policy Issued Event → Kafka:
  ├→ Notification Service (send SMS)
  ├→ Analytics Service (update metrics)
  ├→ Partner Service (calculate commission)
  └→ Document Service (generate PDF)
  
Each service processes independently with retries
```

#### CQRS Pattern
```
Command (Write):
  Proto → Service Logic → PostgreSQL (write) → Kafka Event

Query (Read):
  Proto → Read Model (materialized view) → Fast response
```

### 7.4.7 Backup and Recovery

**PostgreSQL Backup:**
```bash
# Daily full backup
pg_dump insuretech_prod > backup_$(date +%Y%m%d).sql

# Point-in-time recovery (WAL archiving)
archive_command = 'cp %p /archive/%f'
```

**Proto Schema Backup:**
```bash
# Proto files are version controlled in Git
git tag v3.7-schemas
git push origin v3.7-schemas

# Can regenerate schemas from any tagged version
git checkout v3.7-schemas
buf generate
```

**Data Retention:**
| Data Type | Hot (PostgreSQL) | Warm (S3) | Cold (Glacier) | Total Retention |
|-----------|------------------|-----------|----------------|-----------------|
| Active Policies | Lifetime | - | - | Lifetime |
| Expired Policies | 1 year | 5 years | 20 years | 20 years |
| Claims | 2 years | 5 years | 20 years | 20 years |
| Audit Logs | 90 days | 1 year | 7 years | 7 years |
| Telemetry | 90 days | 1 year | Deleted | 1 year |

## 6.5 Data Migration Strategy

**Proto-to-Database Mapping:**
1. Proto definitions serve as canonical data models
2. Database schemas generated from proto files
3. Migrations managed via version-controlled SQL scripts
4. Backward compatibility maintained through proto versioning

**Migration Phases:**
- **Phase 1:** Core entities (User, Policy, Claim, Payment)
- **Phase 2:** Extended entities (Partner, Agent, Product)
- **Phase 3:** Advanced features (IoT, AI, Analytics)

## 6.6 CQRS Implementation

**Command Side (Write):**
- Commands update primary PostgreSQL database
- Events published to Kafka
- Strong consistency guarantees

**Query Side (Read):**
- Materialized views for complex queries
- Read replicas for reporting
- Eventual consistency acceptable
- Cached frequently accessed data in Redis

**Example Flow:**
```
CreatePolicy Command → PostgreSQL INSERT → Kafka PolicyCreated Event
                    ↓
             Read Model Update (async)
                    ↓
        Policy Query → Redis Cache → Read Replica
```

## 6.7 Data Retention & Archival

| Data Type | Hot Storage | Warm Storage | Cold Storage | Retention |
|-----------|-------------|--------------|--------------|-----------|
| Active Policies | PostgreSQL | - | - | Policy lifetime |
| Expired Policies | PostgreSQL (1 year) | S3 (5 years) | Glacier (20 years) | 20 years |
| Claims Data | PostgreSQL | S3 after settlement | Glacier (20 years) | 20 years |
| Audit Logs | TimescaleDB (90 days) | S3 (1 year) | Glacier (7 years) | 7 years |
| IoT Telemetry | TimescaleDB (90 days) | S3 (1 year) | Deleted | 1 year |
| User Sessions | Redis (7 days) | - | - | 7 days |

---

**For complete proto definitions with code examples, see Appendix A.**

## 6.8 System Interface Architecture

### 6.8.1 API Category Specifications

| API Category | Protocol | Use Case | Security Layer | Performance Target |
|-------------|----------|----------|---------------|-------------------|
| **Category 1** | **Protocol Buffer + gRPC** | Gateway ↔ Microservices | System Admin Middle Layer | < 100ms |
| **Category 2** | **GraphQL + JWT** | Gateway ↔ Customer Device | JWT + OAuth v2 | < 2 seconds |
| **Category 3** | **RESTful + JSON (OpenAPI)** | 3rd Party Integration | Server-side Auth | < 200ms |
| **Public API** | **RESTful + JSON (OpenAPI)** | Product Search/List | Public Access | < 1 second |

### 6.8.2 System Interface Diagram

\\\
┌─────────────────────────────────────────────────────────────┐
│                    CLOUDFLARE PROXY                         │
└─────────────────────┬───────────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────────┐
│                    NGINX GATEWAY                            │
└─────────────────────┬───────────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────────┐
│  PUBLIC API (REST/JSON) │ CATEGORY 2 (GraphQL + JWT)        │
│  - Product Search       │ - Customer Device                  │  
│  - Product List         │ - Mobile Apps                      │
└─────────────────────────┼───────────────────────────────────┘
                         │
              ┌─────────▼─────────┐
              │  API GATEWAY      │
              │  (OAuth2 + JWT)   │
              └─────────┬─────────┘
                       │
    ┌─────────────────┼─────────────────┐
    │                 │                 │
┌───▼────┐ ┌─────────▼──────────┐ ┌─────▼─────────┐
│CAT 3   │ │     CATEGORY 1     │ │   INTERNAL    │
│REST API│ │  (gRPC + ProtoBuf) │ │  MICROSERVICES│
│3rd     │ │  Microservices     │ │   (Kafka      │
│Party   │ │  Communication     │ │   Orchestrated)│
└────────┘ └────────────────────┘ └───────────────┘
\\\

