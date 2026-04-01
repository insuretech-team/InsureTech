# InsureTech Proto Documentation - Complete Guide

## Overview

This directory contains comprehensive documentation of the InsureTech protobuf API definitions. The proto files define the complete data model and service interfaces for the InsureTech insurance platform.

## Documentation Files Generated

### 1. **PROTO_INDEX.md** ⭐ START HERE
**Purpose**: Master index and navigation guide  
**Contains**:
- Complete list of all 39 modules
- Module directory structure
- Key design principles
- Entity relationships diagram
- Database schema conventions
- API conventions and patterns
- Getting started guide
- Troubleshooting section

**Use when**: You need an overview or are looking for a specific module

---

### 2. **PROTO_QUICK_REFERENCE.md** ⭐ FOR DEVELOPERS
**Purpose**: Quick lookup and common patterns  
**Contains**:
- System architecture diagram
- Core data flows (Policy lifecycle, Claim processing)
- Module quick reference table
- Critical entity relationships
- Standard field checklist
- Common status enums
- Error code quick reference
- Service RPC naming conventions
- HTTP mapping patterns
- Database index strategy
- Event publishing patterns
- Security checklist
- Common developer tasks
- Common validation rules
- Performance tuning tips
- Testing considerations
- Deployment checklist
- Useful commands

**Use when**: You need quick answers, patterns, or examples

---

### 3. **PROTO_COMMON_TYPES.md**
**Purpose**: Shared utilities and common definitions  
**Contains**:
- db.proto - Database schema annotations (TableOptions, ColumnOptions, ForeignKey, IndexOptions)
- error.proto - Standard error handling (Error message, FieldViolation, 100+ error codes, ErrorSeverity)
- security.proto - Security and privacy annotations (Field-level annotations, SecurityClassification, DataCategory, SecurityEvent)
- types.proto - Common data types (Money, Address, Phone, Email, Document, AuditInfo, InsuranceType)
- Complete enum definitions
- Field annotation reference

**Use when**: You need to understand common types or add a new entity with proper annotations

---

### 4. **PROTO_CORE_MODULES.md**
**Purpose**: Core business domain documentation  
**Contains**:

#### Authentication Module (authn)
- User entity with 15+ fields
- UserProfile, Session, OTP entities
- UserDocument for identity verification
- Auth Service RPC methods

#### Authorization Module (authz)
- Role entity with hierarchy
- UserRole assignments
- RBAC policy rules
- Casbin integration
- MFA configuration
- Access decision auditing
- Authz Service RPC methods

#### Policy Module (policy)
- Policy entity (core insurance contract)
- Policy status lifecycle and transitions
- Quotation entity for pre-sale
- PolicyServiceRequest for modifications
- Policy Service RPC methods

#### Claims Module (claims)
- Claim entity with full lifecycle
- Claim status transitions
- Assessment and approval workflow
- Fraud check integration
- Claims Service RPC methods

#### Payment Module (payment)
- Payment entity with gateway integration
- Payment status lifecycle
- Reconciliation support
- Refund handling
- Payment Service RPC methods

#### Products Module (products)
- Product entity with underwriting rules
- ProductPlan variants (Silver, Gold, Platinum)
- Rider/add-ons for products
- Pricing configuration
- Product Service RPC methods

**Use when**: You're working with core business logic or need to understand entity relationships

---

### 5. **PROTO_SUPPORT_MODULES.md**
**Purpose**: Support services and infrastructure  
**Contains**:

#### Audit & Compliance (audit)
- AuditEvent for business events
- AuditLog for system audit trail (FR-153 to FR-158)
- ComplianceLog for regulatory compliance
- EventCategory, EventSeverity, AuditAction enums

#### API Key Module (apikey)
- ApiKey entity with security features
- ApiKeyUsage tracking (FR-207)
- Distributed tracing support (FR-208)
- Rate limiting and IP whitelist
- Key rotation support

#### KYC Module (kyc)
- KYCVerification with risk scoring
- DocumentVerification for identity docs
- PEP and AML checks
- Risk categorization

#### Support Module (support)
- Ticket entity with lifecycle
- TicketMessage for communications
- FAQ database
- KnowledgeBase articles

#### Notification Module (notification)
- Notification entity for multi-channel delivery
- Alert entity for system alerts
- Retry and delivery tracking
- Template variable support

#### Document Module (document)
- DocumentTemplate for reusable templates
- DocumentGeneration with job tracking
- Template variables and formatting
- Document expiry and versioning

#### Media Module (media)
- Media entity for file storage
- ProcessingJob for media transformation
- Virus scanning status
- Storage lifecycle management

#### Storage Module (storage)
- Storage entity for S3/cloud storage
- Encryption and lifecycle policies
- Quota management
- Multi-provider support

**Use when**: You need to implement support features or infrastructure concerns

---

### 6. **PROTO_FILES_SUMMARY.md**
**Purpose**: High-level overview of all modules  
**Contains**:
- Complete 39-module list with brief descriptions
- Module-by-module breakdown
- Key messages and enums for each module
- Common patterns explanation
- Cross-module relationships
- Service layer standards
- Distributed tracing and event-driven architecture details

**Use when**: You need an executive summary or high-level understanding

---

## Directory Structure

```
E:\Projects\InsureTech\proto\
├── insuretech/
│   ├── ai/                    - AI agents (6 files)
│   ├── analytics/             - Metrics & dashboards (4 files)
│   ├── apikey/                - API authentication (4 files)
│   ├── audit/                 - Audit trail (5 files)
│   ├── authn/                 - Authentication (7 files)
│   ├── authz/                 - Authorization (9 files)
│   ├── b2b/                   - B2B operations (4 files)
│   ├── beneficiary/           - Beneficiary management (3 files)
│   ├── billing/               - Invoice management (1 file)
│   ├── claims/                - Claim management (3 files)
│   ├── commission/            - Commission tracking (3 files)
│   ├── common/v1/             - Shared utilities (4 files)
│   ├── document/              - Document generation (3 files)
│   ├── endorsement/           - Endorsements (3 files)
│   ├── fraud/                 - Fraud detection (3 files)
│   ├── insurance/             - Insurance service (1 file)
│   ├── insurer/               - Insurer management (3 files)
│   ├── iot/                   - IoT devices (3 files)
│   ├── kyc/                   - KYC verification (3 files)
│   ├── media/                 - Media storage (3 files)
│   ├── mfs/                   - MFS integration (3 files)
│   ├── notification/          - Notifications (3 files)
│   ├── orders/                - Order management (3 files)
│   ├── partner/               - Partner management (3 files)
│   ├── payment/               - Payment processing (3 files)
│   ├── policy/                - Policy management (3 files)
│   ├── products/              - Product definitions (3 files)
│   ├── refund/                - Refund processing (3 files)
│   ├── renewal/               - Renewal management (3 files)
│   ├── report/                - Reporting (3 files)
│   ├── services/              - Service providers (1 file)
│   ├── storage/               - File storage (3 files)
│   ├── support/               - Customer support (3 files)
│   ├── task/                  - Task management (3 files)
│   ├── tenant/                - Multi-tenancy (3 files)
│   ├── underwriting/          - Underwriting (3 files)
│   ├── voice/                 - Voice interaction (3 files)
│   └── webrtc/v1/             - WebRTC communication (9 files)
│
└── check_migrations.py        - Migration validation script
```

## Quick Start Guide

### Step 1: Understand the Architecture
Read: **PROTO_INDEX.md**
- Get overview of 39 modules
- Understand design principles
- Review entity relationships

### Step 2: Find Your Domain
Choose based on your work:
- **Core Business**: Read PROTO_CORE_MODULES.md
- **Support Services**: Read PROTO_SUPPORT_MODULES.md
- **Common Types**: Read PROTO_COMMON_TYPES.md

### Step 3: Implement
Use: **PROTO_QUICK_REFERENCE.md** for patterns and examples
- RPC method naming
- HTTP mappings
- Error handling
- Validation rules

### Step 4: Review Checklist
Before deployment:
- Security annotations ✓
- Audit logging ✓
- Error handling ✓
- Distributed tracing ✓

## Key Concepts

### Multi-Tenancy
Every business entity includes `tenant_id` for complete data isolation:
```protobuf
string tenant_id = X [(insuretech.common.v1.column) = {
  column_name: "tenant_id"
  foreign_key: { references_table: "tenants" }
}];
```

### Audit Trail
Standard audit information in all entities:
```protobuf
insuretech.common.v1.AuditInfo audit_info = X [
  (insuretech.common.v1.column) = {
    column_name: "audit_info"
    is_json: true
  }
];
```

### Security First
Sensitive fields marked with security annotations:
```protobuf
string email [(insuretech.common.v1.pii) = true];
string password [(insuretech.common.v1.log_redacted) = true];
```

### Event-Driven
All major operations publish events to Kafka with correlation_id:
```protobuf
string correlation_id = 8;  // For distributed tracing (FR-208)
google.protobuf.Timestamp timestamp = 9;
```

## Module Statistics

| Category | Count | Examples |
|----------|-------|----------|
| Core Modules | 6 | authn, authz, policy, claims, payment, products |
| Insurance Operations | 17 | beneficiary, billing, commission, fraud, underwriting |
| Support Services | 5 | support, notification, document, media, storage |
| Advanced Features | 4 | ai, analytics, workflow, webrtc |
| Infrastructure | 6 | common, audit, apikey, kyc, iot, mfs |
| **Total Modules** | **39** | |
| **Total Proto Files** | **200+** | Across all modules |

## File Sizes

```
PROTO_INDEX.md                    ~50 KB (Master index)
PROTO_QUICK_REFERENCE.md          ~45 KB (Developer reference)
PROTO_CORE_MODULES.md             ~40 KB (Business logic)
PROTO_SUPPORT_MODULES.md          ~45 KB (Infrastructure)
PROTO_COMMON_TYPES.md             ~35 KB (Shared types)
PROTO_FILES_SUMMARY.md            ~40 KB (Overview)
README_PROTO_DOCS.md              ~15 KB (This file)
```

## Navigation Map

```
START HERE (PROTO_INDEX.md)
    ↓
Choose your focus:
    ├→ Core Business (PROTO_CORE_MODULES.md)
    │   ├→ Authentication/Authorization
    │   ├→ Policy Management
    │   ├→ Claims Processing
    │   ├→ Payments
    │   └→ Products
    │
    ├→ Support Services (PROTO_SUPPORT_MODULES.md)
    │   ├→ Audit & Compliance
    │   ├→ API Keys & Security
    │   ├→ KYC Verification
    │   ├→ Support Ticketing
    │   ├→ Notifications
    │   ├→ Documents & Media
    │   └→ Storage
    │
    ├→ Common Utilities (PROTO_COMMON_TYPES.md)
    │   ├→ Error Handling
    │   ├→ Security Annotations
    │   ├→ Database Options
    │   └→ Shared Data Types
    │
    └→ Quick Reference (PROTO_QUICK_REFERENCE.md)
        ├→ Patterns & Examples
        ├→ Status Enums
        ├→ Error Codes
        ├→ RPC Methods
        └→ Deployment Checklist
```

## For Different Roles

### Backend Developer
1. Read: PROTO_QUICK_REFERENCE.md (patterns)
2. Read: PROTO_CORE_MODULES.md or PROTO_SUPPORT_MODULES.md (relevant domain)
3. Reference: PROTO_COMMON_TYPES.md (for annotations)
4. Check: PROTO_INDEX.md (for relationships)

### Database Administrator
1. Read: PROTO_COMMON_TYPES.md (db.proto section)
2. Read: PROTO_INDEX.md (database section)
3. Reference: PROTO_QUICK_REFERENCE.md (index strategy)
4. Check: Each module's entity definitions

### API Consumer/Client Developer
1. Read: PROTO_QUICK_REFERENCE.md (HTTP mappings)
2. Read: PROTO_INDEX.md (API conventions)
3. Reference: PROTO_CORE_MODULES.md or PROTO_SUPPORT_MODULES.md (specific service)
4. Check: Error codes in PROTO_COMMON_TYPES.md

### QA/Tester
1. Read: PROTO_QUICK_REFERENCE.md (status enums, validation rules)
2. Read: PROTO_CORE_MODULES.md (entity lifecycles)
3. Reference: PROTO_INDEX.md (entity relationships)
4. Check: PROTO_QUICK_REFERENCE.md (testing section)

### Security Officer
1. Read: PROTO_COMMON_TYPES.md (security.proto section)
2. Read: PROTO_SUPPORT_MODULES.md (audit section)
3. Reference: PROTO_QUICK_REFERENCE.md (security checklist)
4. Check: PROTO_INDEX.md (compliance & regulatory)

### DevOps/SRE
1. Read: PROTO_INDEX.md (architecture overview)
2. Read: PROTO_QUICK_REFERENCE.md (deployment checklist)
3. Reference: PROTO_SUPPORT_MODULES.md (infrastructure modules)
4. Check: PROTO_COMMON_TYPES.md (database options)

## Common Questions & Answers

### Q: What's the difference between entity, events, and services?
**A**: 
- **entity/v1/**: Data models (what's stored in database)
- **events/v1/**: Event messages (what's published to Kafka)
- **services/v1/**: RPC definitions (what APIs expose)

### Q: How do I add a new field to an entity?
**A**: 
1. Open the entity proto file
2. Add field with unique number
3. Add column annotation with database constraints
4. Add security annotations if PII/sensitive
5. Update AuditLog as needed

### Q: How is multi-tenancy implemented?
**A**: Every business entity has `tenant_id` field:
- Enforced in database with FK constraint
- Used in WHERE clauses to filter data
- Ensures data isolation between customers

### Q: What's the correlation_id for?
**A**: Distributed tracing across services (FR-208):
- Generated for each request
- Passed through event events
- Logged in audit trails
- Enables tracing request through multiple services

### Q: How many error codes are there?
**A**: 100+ codes organized by:
- Generic (1000-1099)
- Validation (1100-1199)
- Authentication (1200-1299)
- Business logic (1400-1499)
- Domain-specific (2000+)

### Q: Can I modify existing proto files?
**A**: Yes, but follow backward compatibility:
- Add new fields with new field numbers
- Don't change existing field numbers
- Don't rename existing fields
- Use deprecation for removal

### Q: What's the difference between soft_delete and hard delete?
**A**: 
- **Soft delete**: Sets deleted_at timestamp, data recoverable
- **Hard delete**: Removes data permanently
- Soft delete is default for compliance/audit

## Troubleshooting

### Proto Compilation Error
- Check all imports are correct
- Verify field numbers don't exceed 1000
- Ensure all referenced messages are defined
- Check for circular imports

### Migration Order Issue
- Lower migration_order runs first
- FK dependencies must be resolved
- Tables referenced by FKs must exist first

### Index Not Being Used
- Verify index is created correctly
- Check if index name matches definition
- Ensure column names are correct
- Review query execution plan

### Data Isolation Issue
- Verify tenant_id is included in WHERE clause
- Check FK constraint on tenant_id
- Review RLS (row-level security) policies
- Check application-level filtering

## Resources

### External Documentation
- [Protocol Buffers Guide](https://developers.google.com/protocol-buffers)
- [gRPC Documentation](https://grpc.io/docs)
- [Google API Design Guide](https://cloud.google.com/apis/design)

### Internal References
- See individual proto files for detailed comments
- Check `check_migrations.py` for migration validation
- Review generated code in `gen/go/insuretech` directory

## Maintenance

### Keeping Documentation Updated
1. When adding new module: Add to PROTO_INDEX.md
2. When modifying entity: Update relevant documentation
3. When changing status enum: Update PROTO_QUICK_REFERENCE.md
4. When adding error code: Update PROTO_COMMON_TYPES.md

### Validation
- All 39 modules documented ✓
- All entity relationships mapped ✓
- All error codes catalogued ✓
- All service methods described ✓
- All security annotations documented ✓

## Support

For questions about proto definitions:
1. Check relevant documentation file
2. Search PROTO_INDEX.md for keywords
3. Review specific proto file comments
4. Check PROTO_QUICK_REFERENCE.md for patterns

---

## Summary of Documentation Files

| File | Purpose | Size | Best For |
|------|---------|------|----------|
| PROTO_INDEX.md | Master index & navigation | ~50KB | Overview, finding things |
| PROTO_QUICK_REFERENCE.md | Patterns & quick lookup | ~45KB | Developers, examples |
| PROTO_CORE_MODULES.md | Business domain docs | ~40KB | Core business logic |
| PROTO_SUPPORT_MODULES.md | Infrastructure docs | ~45KB | Support & infrastructure |
| PROTO_COMMON_TYPES.md | Shared types reference | ~35KB | Annotations, common types |
| PROTO_FILES_SUMMARY.md | High-level overview | ~40KB | Executives, overview |
| README_PROTO_DOCS.md | This file | ~15KB | Getting started |

**Total Documentation**: ~270KB of comprehensive proto API documentation

---

**Last Updated**: 2024  
**Version**: 1.0  
**Status**: Complete

Start with **PROTO_INDEX.md** for full overview or **PROTO_QUICK_REFERENCE.md** for quick answers.
