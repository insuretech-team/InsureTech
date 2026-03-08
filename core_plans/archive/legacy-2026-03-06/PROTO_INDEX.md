# InsureTech Proto - Complete Index & Guide

## Document Overview

This index provides a complete guide to the InsureTech proto files organized in `E:\Projects\InsureTech\proto`.

### Related Documentation Files
1. **PROTO_FILES_SUMMARY.md** - High-level overview of all 39 modules
2. **PROTO_COMMON_TYPES.md** - Common data types, error codes, security annotations
3. **PROTO_CORE_MODULES.md** - Authentication, Authorization, Policy, Claims, Payment, Products
4. **PROTO_SUPPORT_MODULES.md** - Audit, API Keys, KYC, Support, Notifications, Documents, Media, Storage

---

## Quick Navigation

### By Business Domain

#### 1. User Management & Security
- **Authentication (authn)**: User accounts, sessions, OTP, identity documents
- **Authorization (authz)**: Role-based access control, policies, MFA
- **API Key (apikey)**: API authentication and usage tracking
- **Audit (audit)**: Audit trails, compliance logging

#### 2. Core Insurance Operations
- **Policy (policy)**: Policy management, quotations, service requests
- **Claims (claims)**: Claim submission, assessment, settlement
- **Payment (payment)**: Payment processing and reconciliation
- **Products (products)**: Product definitions, plans, riders, pricing

#### 3. Business Support
- **Support (support)**: Ticketing system, FAQ, knowledge base
- **Notification (notification)**: Multi-channel notifications and alerts
- **Document (document)**: Document generation and templating
- **Media (media)**: Media file storage and processing

#### 4. Compliance & Integration
- **KYC (kyc)**: Know Your Customer verification
- **B2B (b2b)**: Business-to-business operations
- **Partner (partner)**: Partner/channel management
- **Tenant (tenant)**: Multi-tenancy support

#### 5. Financial & Operational
- **Billing (billing)**: Invoice management
- **Commission (commission)**: Commission tracking and payouts
- **Refund (refund)**: Refund processing
- **Renewal (renewal)**: Policy renewal management

#### 6. Underwriting & Risk
- **Underwriting (underwriting)**: Underwriting decisions, health declarations
- **Fraud (fraud)**: Fraud detection and case management
- **Beneficiary (beneficiary)**: Beneficiary management
- **Insurer (insurer)**: Insurance company management

#### 7. Advanced Features
- **AI (ai)**: AI agents for underwriting, claims assessment, fraud detection
- **Analytics (analytics)**: Metrics, dashboards, reporting
- **Workflow (workflow)**: Workflow orchestration
- **IoT (iot)**: IoT device management
- **Voice (voice)**: Voice interaction and transcription
- **WebRTC (webrtc)**: Real-time communication
- **MFS (mfs)**: Mobile Financial Services integration

#### 8. Infrastructure
- **Common (common)**: Shared types and utilities
- **Storage (storage)**: File storage management
- **Report (report)**: Business reporting
- **Task (task)**: Task management
- **Services (services)**: Service provider management

#### 9. Other Services
- **Endorsement (endorsement)**: Policy endorsements
- **Orders (orders)**: Order management
- **Insurance Service (insurance)**: High-level insurance operations

---

## Module Directory Structure

Each main module follows this pattern:

```
insuretech/{module}/
├── entity/v1/
│   └── *.proto          # Data models with DB annotations
├── events/v1/
│   └── *.proto          # Event messages for event streaming
└── services/v1/
    └── *.proto          # gRPC service definitions
```

### Example: Policy Module
```
insuretech/policy/
├── entity/v1/
│   ├── policy.proto
│   ├── quotation.proto
│   └── policy_service_request.proto
├── events/v1/
│   └── policy_events.proto
└── services/v1/
    └── policy_service.proto
```

---

## Complete Module List (39 Modules)

### Core Modules
1. **ai/** - AI agents and intelligent processing
2. **analytics/** - Metrics, dashboards, reports
3. **apikey/** - API authentication and usage tracking
4. **audit/** - Audit trails and compliance logging
5. **authn/** - User authentication and sessions
6. **authz/** - Authorization and role-based access control

### Insurance Operations
7. **b2b/** - Business-to-business operations
8. **beneficiary/** - Beneficiary management
9. **billing/** - Invoice management
10. **claims/** - Claim management
11. **commission/** - Commission tracking and payouts
12. **endorsement/** - Policy endorsements and amendments
13. **fraud/** - Fraud detection and case management
14. **insurer/** - Insurance company management
15. **kyc/** - Know Your Customer verification
16. **orders/** - Order management
17. **partner/** - Partner/channel management
18. **payment/** - Payment processing
19. **policy/** - Core policy management
20. **products/** - Product definitions and configuration
21. **refund/** - Refund processing
22. **renewal/** - Policy renewal management
23. **underwriting/** - Underwriting decisions and risk assessment

### Support & Customer Services
24. **document/** - Document generation and templating
25. **media/** - Media file storage and processing
26. **notification/** - Multi-channel notifications
27. **support/** - Customer support ticketing
28. **voice/** - Voice interaction and transcription

### Advanced Features
29. **iot/** - IoT device management
30. **mfs/** - Mobile Financial Services integration
31. **workflow/** - Workflow orchestration
32. **webrtc/** - Real-time WebRTC communication

### Infrastructure & Utilities
33. **common/** - Shared types and utilities
34. **insurance/** - High-level insurance operations
35. **report/** - Business reporting
36. **services/** - Service provider management
37. **storage/** - File storage management
38. **task/** - Task management
39. **tenant/** - Multi-tenancy support

---

## File Statistics

### Total Files
- **Proto Files**: 200+ files
- **Documentation Files**: 5 markdown files
- **Helper Scripts**: 1 migration script

### Module Breakdown
- **Core Modules**: 6 modules (18 core entity files)
- **Insurance Operations**: 17 modules (40+ files)
- **Support Services**: 5 modules (15+ files)
- **Advanced Features**: 4 modules (20+ files)
- **Infrastructure**: 6 modules (10+ files)
- **Miscellaneous**: 1 module (4 files)

---

## Common Patterns & Standards

### 1. File Organization
```
package insuretech.{module}.{category}.v1;

option go_package = "github.com/newage-saint/insuretech/gen/go/insuretech/{module}/{category}/v1";
option csharp_namespace = "Insuretech.{Module}.{Category}.V1";
```

### 2. Entity Annotations
All entities include:
- **Database mapping**: Table name, schema, migration order
- **Audit support**: soft_delete, audit_fields
- **Security**: PII, encryption, masking flags
- **Indexing**: Composite indexes, foreign keys

### 3. Standard Entity Fields
```protobuf
message Entity {
  string id = 1;                    // UUID primary key
  string tenant_id = 2;             // Multi-tenancy
  string status = 3;                // ACTIVE, INACTIVE, etc.
  google.protobuf.Timestamp created_at = 4;
  google.protobuf.Timestamp updated_at = 5;
  insuretech.common.v1.AuditInfo audit_info = 6;
}
```

### 4. Event Pattern
```protobuf
message EntityEvent {
  string event_id = 1;
  string entity_id = 2;
  string entity_type = 3;
  string action = 4;                // CREATED, UPDATED, DELETED
  google.protobuf.Timestamp timestamp = 5;
  string correlation_id = 6;        // Distributed tracing
  map<string, string> metadata = 7;
}
```

### 5. Service RPC Pattern
```protobuf
service EntityService {
  rpc CreateEntity(CreateEntityRequest) returns (CreateEntityResponse) {
    option (google.api.http) = { post: "/v1/entities" body: "*" };
  }
  
  rpc GetEntity(GetEntityRequest) returns (GetEntityResponse) {
    option (google.api.http) = { get: "/v1/entities/{id}" };
  }
  
  // ... more methods
}
```

### 6. Response Pattern
```protobuf
message EntityResponse {
  Entity entity = 1;
  insuretech.common.v1.Error error = 2;  // Standard error handling
}

message ListResponse {
  repeated Entity items = 1;
  int32 total_count = 2;
  int32 page = 3;
  int32 page_size = 4;
  insuretech.common.v1.Error error = 5;
}
```

---

## Key Design Principles

### 1. Multi-Tenancy
- Every entity with business data includes `tenant_id`
- Data isolation at database and application level
- Tenant-scoped indexes for performance

### 2. Audit & Compliance
- Comprehensive audit logging with `AuditLog` entity
- Field-level change tracking (old_values, new_values)
- Regulatory compliance logging with compliance framework references
- User and IP tracking for security

### 3. Security
- PII fields marked with `@(pii = true)`
- Sensitive data encrypted at rest
- Log masking and redaction for sensitive fields
- API key hashing (SHA-256) for secure storage
- Distributed tracing with correlation_id (FR-208)

### 4. Database Optimization
- Composite indexes for common queries
- Partition strategies for large tables
- Soft delete for data retention
- Foreign key constraints for referential integrity

### 5. Event-Driven Architecture
- Events for all significant business operations
- Kafka event streaming capability
- Distributed tracing across services
- Real-time updates and service communication

### 6. Error Handling
- Standardized error codes (1000+ codes)
- Field-level validation errors
- Retryable vs non-retryable errors
- HTTP status code mapping

---

## Key Entity Relationships

### Policy-Centric Model
```
Insurer
  └── Product
      ├── ProductPlan
      └── Rider
  └── Policy
      ├── Quotation
      ├── Claim
      ├── Payment
      ├── Beneficiary
      ├── Endorsement
      └── RenewalSchedule
```

### User-Centric Model
```
User
  ├── UserProfile
  ├── UserDocument
  ├── Session
  ├── OTP
  ├── UserRole
  │   └── Role
  │       └── PolicyRule (RBAC)
  └── [Audit trails in all operations]
```

### Transaction Model
```
Policy
  ├── Quotation (pre-sale)
  ├── Payment (premium collection)
  │   └── Refund (if applicable)
  ├── Claim (incident handling)
  │   ├── PaymentSettlement
  │   └── FraudDetection
  └── Renewal (policy continuation)
```

---

## Distributed Tracing

All operations support distributed tracing via:
- **correlation_id**: Used in events and audit logs
- **trace_id**: Captured in audit logs and API usage
- **span tracking**: Through service call chains

### Implementation
- Reference: FR-208 (Distributed tracing support)
- Applied in: ApiKeyUsage.trace_id, AuditLog.trace_id, all events
- Format: VARCHAR(64) field for trace ID

---

## Multi-Language Support

Proto files support code generation for:
- **Go**: `option go_package`
- **C#**: `option csharp_namespace`
- **Java**: Standard protobuf generation
- **Python**: Standard protobuf generation
- **TypeScript/JavaScript**: Via grpc-web or similar

Generated code locations:
- Go: `github.com/newage-saint/insuretech/gen/go/insuretech/{module}/{category}/v1`
- C#: `Insuretech.{Module}.{Category}.V1`

---

## Feature References

### Functional Requirements
- **FR-153 to FR-158**: Audit trail implementation (audit_log.proto)
- **FR-207**: API usage tracking (api_key_usage.proto)
- **FR-208**: Distributed tracing (trace_id fields)

### Cross-Guideline References
- **CG-6**: Distributed tracing with correlation_id

---

## Database Schema Conventions

### Schema Names
- `authn_schema`: Authentication entities (users, sessions, api_keys)
- `compliance_schema`: Audit and compliance logs
- `public`: Default schema for most business entities

### Migration Order
- Lower numbers run first
- Ensures FK dependencies are met
- Example: authn_schema tables run before policy tables

### Partition Strategies
- **RANGE_MONTH**: Time-series data (logs, events)
- **RANGE_YEAR**: Historical data (audit logs, compliance logs)
- **HASH**: Large transaction tables (api_key_usage)

---

## Security Classifications

### Field Annotations
- **pii**: Personally identifiable information
- **sensitive**: Highly sensitive data requiring audit
- **encrypted_security**: Encrypted at rest
- **log_masked**: Masked in logs (e.g., 017****5678)
- **log_redacted**: Completely redacted in logs
- **requires_consent**: GDPR consent required
- **data_purpose**: Purpose for data collection

### Security Levels
- **PUBLIC**: No restrictions
- **INTERNAL**: Internal use only
- **CONFIDENTIAL**: Restricted access
- **HIGHLY_CONFIDENTIAL**: Strictly controlled

---

## API Conventions

### HTTP Verbs
- **POST**: Create new resource or perform action
- **GET**: Retrieve resource
- **PUT**: Replace entire resource
- **PATCH**: Partial update
- **DELETE**: Delete resource

### URL Patterns
```
POST   /v1/{resources}                    - Create
GET    /v1/{resources}                    - List
GET    /v1/{resources}/{id}               - Get
PUT    /v1/{resources}/{id}               - Update
DELETE /v1/{resources}/{id}               - Delete
POST   /v1/{resources}/{id}:action        - Custom action
POST   /v1/{resources}:batch              - Batch operation
```

### Example from policyservice.proto
```protobuf
rpc CreateQuotation(CreateQuotationRequest) returns (CreateQuotationResponse) {
  option (google.api.http) = { post: "/v1/policies/quotations" body: "*" };
}

rpc GetPolicy(GetPolicyRequest) returns (GetPolicyResponse) {
  option (google.api.http) = { get: "/v1/policies/{policy_id}" };
}

rpc CancelPolicy(CancelPolicyRequest) returns (CancelPolicyResponse) {
  option (google.api.http) = {
    post: "/v1/policies/{policy_id}:cancel"
    body: "*"
  };
}
```

---

## Pagination & Filtering

### List Requests
```protobuf
message ListRequest {
  int32 page = 1;              // 1-based page number
  int32 page_size = 2;         // Items per page (max 100)
  string sort_by = 3;          // Field to sort by
  string sort_order = 4;       // ASC or DESC
  map<string, string> filters = 5; // Filter conditions
}
```

### List Responses
```protobuf
message ListResponse {
  repeated Entity items = 1;
  int32 total_count = 2;
  int32 page = 3;
  int32 page_size = 4;
  string next_page_token = 5;  // Optional cursor
  insuretech.common.v1.Error error = 6;
}
```

---

## Versioning Strategy

### API Versions
- All services use **v1** (first API version)
- Future versions: v2, v3, etc.
- Backward compatibility maintained for minor updates

### Proto Package Versioning
```
package insuretech.{module}.{category}.v1;
```

---

## Documentation References

### External Standards
- **Protobuf 3**: [https://developers.google.com/protocol-buffers](https://developers.google.com/protocol-buffers)
- **gRPC**: [https://grpc.io](https://grpc.io)
- **Google API Design Guide**: [https://cloud.google.com/apis/design](https://cloud.google.com/apis/design)

### Security Standards
- **GDPR**: Field annotations for compliance
- **HIPAA**: Health data handling (critical_illness, health information)
- **PCI DSS**: Payment data security (payment_method, card data)

---

## Getting Started with Proto Files

### 1. Setup
```bash
# Install protoc compiler
# Add grpc-gateway for HTTP/REST support
# Set up code generation tools
```

### 2. Generate Code
```bash
# Generate Go code
protoc --go_out=. --go-grpc_out=. insuretech/**/*.proto

# Generate with gRPC Gateway
protoc --go_out=. --go-grpc_out=. \
       --grpc-gateway_out=. insuretech/**/*.proto
```

### 3. Use in Services
```go
import "insuretech/policy/entity/v1"
import "insuretech/policy/services/v1"

// Use generated types and services
policy := &policyv1.Policy{...}
client := policyv1.NewPolicyServiceClient(conn)
```

---

## Common Development Tasks

### Adding a New Entity
1. Create `entity/v1/{entity}.proto`
2. Define message with DB annotations
3. Include standard fields (id, tenant_id, audit_info)
4. Mark sensitive fields with security annotations
5. Add to existing service or create new one

### Adding a New Service
1. Create `services/v1/{service}.proto`
2. Define service interface (rpc methods)
3. Define request/response message types
4. Add HTTP/gRPC mappings
5. Include error handling

### Adding a New Event
1. Create event message in `events/v1/{module}_events.proto`
2. Include event_id, timestamp, correlation_id
3. Include entity-specific payload
4. Use consistent naming (EntityActionEvent)

### Adding Audit Trail
1. Annotate fields with audit_fields in table options
2. Mark sensitive fields with pii/sensitive/encrypted
3. AuditLog will automatically track changes
4. Add ComplianceLog entries for regulatory events

---

## Troubleshooting

### Common Issues

**Issue**: Foreign key constraint failures
- **Solution**: Check migration_order values; lower numbers should run first

**Issue**: Duplicate index names
- **Solution**: Ensure index_name is unique across all tables

**Issue**: Proto compilation errors
- **Solution**: Verify all imports are correctly specified; check package names

**Issue**: API fields not appearing
- **Solution**: Verify field numbers don't exceed 1000 (reserved for framework); check proto3 syntax

---

## Performance Considerations

### Index Strategies
- **BTREE**: Default for most columns, good for ranges
- **HASH**: For equality searches only
- **GIN**: For JSON, arrays, full-text search
- **BRIN**: For large tables with natural ordering

### Partition Strategies
- **RANGE_MONTH**: Logs, events (time-series)
- **RANGE_YEAR**: Historical data (compliance, audit)
- **HASH**: Transaction tables (api_key_usage)

### Query Optimization
- Use composite indexes for common filter combinations
- Denormalize for read-heavy operations
- Archive old audit logs to separate partitions
- Use appropriate data types (INT vs VARCHAR)

---

## Compliance & Regulatory

### Implemented Controls
- ✅ Audit logging (AuditLog, AuditEvent)
- ✅ Compliance logging (ComplianceLog)
- ✅ Distributed tracing (FR-208)
- ✅ Data retention policies
- ✅ GDPR consent tracking
- ✅ PII protection and encryption
- ✅ Role-based access control
- ✅ Audit trail immutability

### Regulatory Frameworks Supported
- GDPR (General Data Protection Regulation)
- HIPAA (Health Insurance Portability and Accountability Act)
- PCI DSS (Payment Card Industry Data Security Standard)
- SOX (Sarbanes-Oxley Act)
- Local regulatory requirements (via compliance_framework)

---

## Contact & Support

For questions or issues with the proto definitions:
1. Refer to PROTO_FILES_SUMMARY.md for high-level overview
2. Refer to PROTO_COMMON_TYPES.md for shared types
3. Refer to PROTO_CORE_MODULES.md for business logic
4. Refer to PROTO_SUPPORT_MODULES.md for support services
5. Check individual proto files for detailed comments

---

**Generated**: Comprehensive guide to InsureTech proto directory structure and conventions  
**Version**: 1.0  
**Last Updated**: 2024

---

End of Index
