# InsureTech Proto - Quick Reference Guide

## System Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────┐
│                         API Gateway Layer                            │
│  (gRPC + REST via grpc-gateway, HTTP/2 for streaming)              │
└────────┬────────────────────────────────────────────────────────────┘
         │
    ┌────┴────┬──────────┬──────────┬──────────┬──────────┬─────────┐
    │          │          │          │          │          │         │
┌───▼──┐  ┌──▼───┐  ┌───▼──┐  ┌───▼──┐  ┌───▼──┐  ┌───▼──┐  ┌───▼──┐
│Auth* │  │Policy│  │Claims│  │Pay*  │  │Prod* │  │Fraud │  │AI    │
│      │  │      │  │      │  │ment  │  │ucts  │  │      │  │      │
└──────┘  └──────┘  └──────┘  └──────┘  └──────┘  └──────┘  └──────┘

┌──────────────────────────────────────────────────────────────────────┐
│              Support & Infrastructure Services                        │
├──────────────────────────────────────────────────────────────────────┤
│  Audit  │ KYC   │ Support │ Notification │ Document │ Media │Workflow
│         │       │         │              │          │       │
└──────────────────────────────────────────────────────────────────────┘

┌──────────────────────────────────────────────────────────────────────┐
│              Data Layer (PostgreSQL Multi-Schema)                     │
├──────────────────────────────────────────────────────────────────────┤
│  authn_schema     │  compliance_schema   │  public schema             │
│  (auth data)      │  (audit/compliance)  │  (business entities)       │
└──────────────────────────────────────────────────────────────────────┘

┌──────────────────────────────────────────────────────────────────────┐
│              Event Streaming Layer (Kafka/RabbitMQ)                  │
│  All services publish events with correlation_id for tracing         │
└──────────────────────────────────────────────────────────────────────┘
```

---

## Core Data Flow

### Policy Lifecycle
```
Customer Inquiry
  ↓
Quotation Created (quotation.proto)
  ↓
Apply for Policy (policy.proto)
  ↓
Underwriting Assessment (underwriting_decision.proto, AI)
  ↓
Policy Approved
  ↓
Payment Processing (payment.proto)
  ↓
Policy Active
  ├→ Claim Submission (claim.proto)
  │   ├→ Assessment (AI claims assessment)
  │   ├→ Fraud Check (fraud_case.proto)
  │   ├→ Approval/Rejection
  │   └→ Settlement (payment.proto)
  │
  ├→ Endorsements (endorsement.proto)
  │
  └→ Renewal (renewal_schedule.proto)
      ↓
      Premium Payment
      ↓
      Policy Renewed or Cancelled
```

### Claim Processing
```
Claim Submitted
  ↓ [ClaimSubmittedEvent published to Kafka]
  ├→ KYC/Verification Check (kyc_verification.proto)
  ├→ Document Processing (document_generation.proto, media.proto)
  ├→ AI Assessment (ai_service.proto - EvaluateClaim)
  ├→ Fraud Detection (fraud_service.proto - DetectFraud)
  │   └→ AI Agent Flagging (ai_events.proto - FraudDetectedEvent)
  │
  ├→ Manual Review (support tickets if needed)
  │
  ├→ Approval/Rejection Decision
  │   └→ [ClaimApprovedEvent or ClaimRejectedEvent]
  │
  └→ Settlement & Payment (payment.proto)
      └→ [PaymentCompletedEvent]
```

---

## Module Quick Reference Table

| Module | Purpose | Key Entities | Status Enums | Events |
|--------|---------|--------------|-------------|--------|
| **authn** | Authentication | User, Session, OTP, UserDocument | ACTIVE, INACTIVE, SUSPENDED | UserLoginEvent, UserRegisteredEvent |
| **authz** | Authorization | Role, UserRole, PolicyRule, CasbinRule | ACTIVE, INACTIVE | AccessGrantedEvent, AccessDeniedEvent |
| **policy** | Policy Management | Policy, Quotation, PolicyServiceRequest | DRAFT, QUOTED, ACTIVE, LAPSED, CANCELLED | PolicyCreatedEvent, PolicyActivatedEvent, PolicyRenewedEvent |
| **claims** | Claim Management | Claim | SUBMITTED, UNDER_REVIEW, APPROVED, REJECTED, SETTLED | ClaimSubmittedEvent, ClaimApprovedEvent, ClaimSettledEvent |
| **payment** | Payment Processing | Payment | INITIATED, PENDING, COMPLETED, FAILED, CANCELLED | PaymentInitiatedEvent, PaymentCompletedEvent |
| **products** | Product Configuration | Product, ProductPlan, Rider, PricingConfig | ACTIVE, INACTIVE, RETIRED | ProductCreatedEvent, ProductActivatedEvent |
| **audit** | Audit Trail | AuditLog, AuditEvent, ComplianceLog | N/A (events only) | AuditEvent (all operations logged) |
| **apikey** | API Authentication | ApiKey, ApiKeyUsage | ACTIVE, EXPIRED, REVOKED, SUSPENDED, ROTATING | ApiKeyCreatedEvent, ApiKeyRevokedEvent |
| **kyc** | KYC Verification | KYCVerification, DocumentVerification | VERIFIED, FAILED, EXPIRED | KYCVerifiedEvent, KYCFailedEvent |
| **support** | Support Ticketing | Ticket, TicketMessage, FAQ, KnowledgeBase | OPEN, IN_PROGRESS, RESOLVED, CLOSED | TicketCreatedEvent, TicketResolvedEvent |
| **ai** | AI Processing | Agent, AgentConfig, AgentPerformance | ACTIVE, INACTIVE | AIDecisionMadeEvent, FraudDetectedEvent |
| **fraud** | Fraud Detection | FraudAlert, FraudCase, FraudRule | ACTIVE, INACTIVE | FraudAlertCreatedEvent, FraudCaseOpenedEvent |

---

## Critical Entity Relationships

### Master-Detail Relationships
```
Insurer (1) ──────→ (N) Policy
            ──────→ (N) Product
            ──────→ (N) Partner

Product (1) ──────→ (N) ProductPlan
          ──────→ (N) Rider

Policy (1) ──────→ (N) Claim
         ──────→ (N) Payment
         ──────→ (N) Beneficiary
         ──────→ (N) Endorsement
         ──────→ (N) RenewalSchedule
         ──────→ (N) TicketMessage (support)

User (1) ──────→ (N) UserRole
       ──────→ (N) Session
       ──────→ (N) UserDocument
       ──────→ (N) OTP
       ──────→ (N) Policy (as holder)
```

### FK Constraints (DELETE Actions)
```
User (CASCADE)
  ├── UserProfile
  ├── Session
  ├── OTP
  ├── UserDocument
  └── UserRole

Policy (CASCADE)
  ├── Claim
  ├── Payment
  ├── Beneficiary
  ├── Endorsement
  ├── RenewalSchedule
  └── ServiceRequest

ApiKey (CASCADE)
  └── ApiKeyUsage

KYCVerification (CASCADE)
  └── DocumentVerification
```

---

## Standard Fields in Every Entity

### Audit Fields (when audit_fields: true)
```protobuf
string id                              // UUID primary key
string tenant_id                       // Multi-tenancy FK
string status                          // ACTIVE, INACTIVE, etc.
google.protobuf.Timestamp created_at   // Auto-set by DB
google.protobuf.Timestamp updated_at   // Auto-updated by DB
insuretech.common.v1.AuditInfo audit_info = {
  created_by                           // User ID
  created_at                           // Timestamp
  updated_by                           // User ID
  updated_at                           // Timestamp
  deleted_by                           // User ID (if soft-deleted)
  deleted_at                           // Timestamp (if soft-deleted)
}
```

### Security Annotations (Sample)
```protobuf
string email [
  (insuretech.common.v1.pii) = true,
  (insuretech.common.v1.log_masked) = true,
  (insuretech.common.v1.data_purpose) = "User authentication"
]

string password_hash [
  (insuretech.common.v1.sensitive) = true,
  (insuretech.common.v1.log_redacted) = true,
  (insuretech.common.v1.encrypted_security) = true
]
```

---

## Common Status Enums

### Policy Status
```
DRAFT → QUOTED → ACTIVE → {SUSPENDED, LAPSED, CANCELLED, RENEWED}
```

### Claim Status
```
DRAFT → SUBMITTED → UNDER_REVIEW → {APPROVED, REJECTED} → SETTLED
                                                    ↓
                                                APPEALED → {APPROVED, REJECTED}
```

### Payment Status
```
INITIATED → PENDING → COMPLETED
         ↓
        FAILED → {CANCELLED, RETRY}
```

### KYC Status
```
PENDING → IN_PROGRESS → {VERIFIED, FAILED, EXPIRED, SUSPENDED}
```

---

## Error Code Quick Reference

### By Category

**Generic (1000-1099)**
```
1000: INTERNAL_ERROR
1001: INVALID_REQUEST
1002: UNAUTHORIZED
1003: FORBIDDEN
1004: NOT_FOUND
1005: ALREADY_EXISTS
1007: RATE_LIMIT_EXCEEDED
1009: TIMEOUT
```

**Validation (1100-1199)**
```
1100: VALIDATION_ERROR
1101: MISSING_REQUIRED_FIELD
1102: INVALID_FIELD_VALUE
1103: INVALID_FIELD_FORMAT
1104: FIELD_OUT_OF_RANGE
```

**Authentication (1200-1299)**
```
1200: AUTHENTICATION_FAILED
1201: INVALID_CREDENTIALS
1202: EXPIRED_TOKEN
1205: OTP_EXPIRED
1206: OTP_INVALID
```

**Business Logic (1400-1499)**
```
1400: BUSINESS_RULE_VIOLATION
1401: INVALID_STATE_TRANSITION
1402: OPERATION_NOT_ALLOWED
1403: QUOTA_EXCEEDED
```

**Policy (2000-2099)**
```
2000: POLICY_NOT_FOUND
2001: POLICY_ALREADY_CANCELLED
2002: POLICY_ALREADY_LAPSED
2003: POLICY_NOT_ACTIVE
```

**Claim (2100-2199)**
```
2100: CLAIM_NOT_FOUND
2101: CLAIM_ALREADY_SETTLED
2102: CLAIM_AMOUNT_EXCEEDS_COVERAGE
2103: CLAIM_OUTSIDE_COVERAGE_PERIOD
```

**Payment (2200-2299)**
```
2200: PAYMENT_NOT_FOUND
2201: PAYMENT_ALREADY_COMPLETED
2202: PAYMENT_FAILED
2203: INSUFFICIENT_FUNDS
```

---

## Service RPC Method Naming Convention

### CRUD Operations
```
Create{Entity}      - Create new instance
Get{Entity}         - Retrieve single instance
List{Entity}s       - List multiple instances
Update{Entity}      - Update instance
Delete{Entity}      - Delete instance
```

### Custom Operations
```
{Verb}{Entity}      - Custom action
Examples:
  - ApproveClaim
  - RejectClaim
  - SettleClaim
  - CancelPolicy
  - RenewPolicy
  - RotateApiKey
  - VerifyUser
  - DetectFraud
```

### Query/Search Operations
```
Search{Entity}s     - Search with filters
Get{Entity}By{Field}  - Find by specific field
List{Entity}sBy{Criteria} - List by criteria
```

---

## HTTP Mapping Patterns

### Standard CRUD
```protobuf
rpc GetEntity(GetEntityRequest) returns (GetEntityResponse) {
  option (google.api.http) = { get: "/v1/entities/{id}" };
}

rpc CreateEntity(CreateEntityRequest) returns (CreateEntityResponse) {
  option (google.api.http) = { post: "/v1/entities" body: "*" };
}

rpc UpdateEntity(UpdateEntityRequest) returns (UpdateEntityResponse) {
  option (google.api.http) = { put: "/v1/entities/{id}" body: "*" };
}

rpc DeleteEntity(DeleteEntityRequest) returns (DeleteEntityResponse) {
  option (google.api.http) = { delete: "/v1/entities/{id}" };
}
```

### Custom Actions
```protobuf
rpc ApproveClaim(ApproveClaimRequest) returns (ApproveClaimResponse) {
  option (google.api.http) = {
    post: "/v1/claims/{claim_id}:approve"
    body: "*"
  };
}

rpc CancelPolicy(CancelPolicyRequest) returns (CancelPolicyResponse) {
  option (google.api.http) = {
    post: "/v1/policies/{policy_id}:cancel"
    body: "*"
  };
}
```

### Batch Operations
```protobuf
rpc BatchCreateClaims(BatchCreateClaimsRequest) returns (BatchCreateClaimsResponse) {
  option (google.api.http) = {
    post: "/v1/claims:batch"
    body: "*"
  };
}
```

---

## Database Index Strategy

### By Module

**authn_schema (Auth Data)**
```
User
  - idx_users_email (BTREE, UNIQUE)
  - idx_users_phone (BTREE, UNIQUE)
  - idx_users_tenant_id (BTREE)

Session
  - idx_sessions_user_id (BTREE)
  - idx_sessions_expires_at (BTREE)

ApiKey
  - idx_api_keys_key_hash (BTREE, UNIQUE)
  - idx_api_keys_owner_id (BTREE)
  - idx_api_keys_status (BTREE)

ApiKeyUsage
  - idx_api_key_usage_key_id (BTREE)
  - idx_api_key_usage_endpoint (BTREE)
  - idx_api_key_usage_timestamp (BTREE) [PARTITION HASH]
```

**compliance_schema (Audit/Compliance)**
```
AuditLog
  - idx_audit_logs_entity_id (BTREE)
  - idx_audit_logs_user_id (BTREE)
  - idx_audit_logs_timestamp (BTREE) [PARTITION RANGE_YEAR]
  - idx_audit_logs_action (BTREE)

ComplianceLog
  - idx_compliance_logs_entity_id (BTREE)
  - idx_compliance_logs_timestamp (BTREE) [PARTITION RANGE_MONTH]

AuditEvent
  - idx_audit_events_entity_id (BTREE)
  - idx_audit_events_category (BTREE)
  - idx_audit_events_timestamp (BTREE) [PARTITION RANGE_MONTH]
```

**public schema (Business Entities)**
```
Policy
  - idx_policies_insurer_id (BTREE)
  - idx_policies_holder_id (BTREE)
  - idx_policies_status (BTREE)
  - idx_policies_renewal_date (BTREE)

Claim
  - idx_claims_policy_id (BTREE)
  - idx_claims_status (BTREE)
  - idx_claims_submission_date (BTREE) [PARTITION RANGE_MONTH]

Payment
  - idx_payments_policy_id (BTREE)
  - idx_payments_status (BTREE)
  - idx_payments_timestamp (BTREE) [PARTITION RANGE_MONTH]
```

---

## Event Publishing Pattern

### Standard Event Structure
```protobuf
message EntityEvent {
  string event_id = 1;                    // UUID, unique per event
  string entity_id = 2;                   // Entity being changed
  string action = 3;                      // CREATED, UPDATED, DELETED
  google.protobuf.Timestamp timestamp = 4; // Event time
  string correlation_id = 5;              // For distributed tracing (FR-208)
  map<string, string> metadata = 6;       // Extra context
}
```

### Kafka Topic Naming
```
insuretech.{module}.{action}     Flat structure
Example:
  - insuretech.policy.created
  - insuretech.claim.submitted
  - insuretech.payment.completed
  - insuretech.fraud.detected
```

### Event Consumer Pattern
```
Services subscribe to events relevant to their domain:
- Policy service → publishes policy.* events, subscribes to payment.completed
- Claims service → subscribes to policy.created, publishes claim.* events
- Fraud service → subscribes to claim.submitted, publishes fraud.detected
- Analytics service → subscribes to all *.* events
```

---

## Security Quick Checklist

### Field-Level Security
- [ ] PII fields marked with `@(pii = true)`
- [ ] Sensitive fields with `@(sensitive = true)`
- [ ] Passwords/keys with `@(log_redacted = true)`
- [ ] Phone/email with `@(log_masked = true)`
- [ ] Payment data with `@(encrypted_security = true)`

### Entity-Level Security
- [ ] `audit_fields: true` for most entities
- [ ] `soft_delete: true` for data retention compliance
- [ ] Foreign key constraints with ON DELETE CASCADE/SET_NULL
- [ ] Row-level security (RLS) enabled for sensitive tables

### API Security
- [ ] API keys hashed (SHA-256) in storage
- [ ] Rate limiting per API key
- [ ] IP whitelist support
- [ ] Scopes/permissions system
- [ ] All responses include error field

### Audit & Compliance
- [ ] AuditLog tracks all CRUD operations
- [ ] ComplianceLog for regulatory requirements
- [ ] Distributed tracing with correlation_id
- [ ] User ID and IP address logged
- [ ] Immutable audit trail

---

## Migration Order Ranges

```
0-10:     System initialization
11-20:    authn_schema tables (users, sessions, roles)
21-30:    authz_schema tables (permissions, policies)
31-50:    compliance_schema (audit, compliance logs)
51-70:    Core business entities (policies, claims)
71-90:    Supporting entities (payments, beneficiaries)
91-100:   Event/streaming tables
101+:     Future/custom tables
```

---

## Common Developer Tasks

### Task 1: Create a New API Endpoint
1. Add RPC method to service.proto
2. Define request and response messages
3. Add HTTP mapping
4. Include error field in response
5. Implement handler in service

### Task 2: Add Audit Trail to Entity
1. Set `audit_fields: true` in table options
2. Mark sensitive fields with security annotations
3. AuditLog automatically tracks changes
4. No additional code needed

### Task 3: Add Distributed Tracing
1. Include correlation_id in event
2. Pass trace_id in audit logs
3. Propagate correlation_id in service calls
4. Analytics can correlate across services

### Task 4: Implement Rate Limiting
1. Create ApiKey entity with rate_limit_per_minute
2. Check ApiKeyUsage count in time window
3. Return 429 (Too Many Requests) if exceeded
4. Publish ApiKeyRateLimitExceededEvent

---

## Common Validation Rules

### Policy Validation
- `end_date` must be > `start_date`
- `premium_amount` must be > 0
- `max_coverage_amount` must be ≥ `deductible`
- `term_months` must be > 0

### Claim Validation
- `claim_amount` must be > 0
- `claim_amount` ≤ `policy.max_coverage_amount`
- Claim can only be submitted if policy is ACTIVE
- `incident_date` must be ≤ claim submission date
- Incident must be within coverage period

### Payment Validation
- `amount` must match policy premium or claim settlement
- Payment status transitions: INITIATED → PENDING → COMPLETED
- Cannot process payment twice
- Refund only possible within 30 days

### User Validation
- `email` must be unique per tenant
- `phone` must be valid format
- `password` must be ≥ 8 chars, include upper, lower, number
- `date_of_birth` must make user ≥ 18 years old

---

## Performance Tuning Tips

### Query Optimization
```
Use composite indexes for:
  - (tenant_id, status)
  - (policy_id, status, created_at)
  - (entity_type, entity_id, action, timestamp)
```

### Partitioning Strategy
```
Partition large tables by date:
  - audit_logs: RANGE MONTH
  - compliance_logs: RANGE YEAR
  - api_key_usage: HASH (by api_key_id)
  - claims: RANGE MONTH (by submission_date)
```

### Connection Pooling
```
Recommended:
  - Min connections: 10
  - Max connections: 50
  - Max idle: 5 minutes
  - Connection timeout: 10 seconds
```

---

## Testing Considerations

### Unit Test Coverage
- RPC request/response marshalling
- Field validation logic
- Error handling paths
- Status transition rules

### Integration Tests
- Service-to-service communication
- Database FK constraints
- Audit trail creation
- Event publishing

### End-to-End Tests
- Complete policy lifecycle
- Claim from submission to settlement
- Multi-step workflows
- Error recovery scenarios

---

## Deployment Checklist

- [ ] All proto files compiled successfully
- [ ] All external dependencies imported
- [ ] Database migrations in correct order
- [ ] Indexes created for performance
- [ ] Audit logging configured
- [ ] Event streaming configured
- [ ] API gateway configured
- [ ] Rate limiting configured
- [ ] CORS policies configured
- [ ] Security annotations documented
- [ ] API documentation generated
- [ ] Load testing completed

---

## Useful Commands

### Proto Compilation
```bash
# Compile all protos
protoc --go_out=. --go-grpc_out=. \
       --grpc-gateway_out=. insuretech/**/*.proto

# Generate with specific output
protoc --go_out=./gen/go --go-grpc_out=./gen/go insuretech/**/*.proto

# Generate for multiple languages
protoc --go_out=. --java_out=. --python_out=. insuretech/**/*.proto
```

### Database Setup
```sql
-- Create schemas
CREATE SCHEMA authn_schema;
CREATE SCHEMA compliance_schema;

-- Run migrations in migration_order
-- Lower numbers first

-- Create indexes after tables
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_policies_policy_number ON policies(policy_number);
```

### Testing gRPC Services
```bash
# Using grpcurl
grpcurl -plaintext localhost:50051 list
grpcurl -plaintext -d '{"id":"123"}' localhost:50051 insuretech.policy.services.v1.PolicyService/GetPolicy

# Interactive client
grpcui -plaintext localhost:50051
```

---

End of Quick Reference Guide
