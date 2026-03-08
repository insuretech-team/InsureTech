# InsureTech Proto Files - Complete Summary

This document contains a comprehensive overview of all proto files in the E:\Projects\InsureTech\proto directory.

## Directory Structure Overview

The proto directory is organized by domain modules with the following structure:
- Each domain has `entity/v1/`, `events/v1/`, and `services/v1/` subdirectories
- Common utilities are in `insuretech/common/v1/`

## Modules

### 1. AI Module (`insuretech/ai/`)
**Purpose**: AI agents for underwriting, claims assessment, fraud detection

**Files**:
- `entity/v1/agent.proto` - Agent entity definitions
- `events/v1/ai_events.proto` - AI events (agent created, conversation started, decision made, fraud detected)
- `services/v1/ai_service.proto` - AI Service with RPC methods (Chat, StreamChat, AnalyzeDocument, DetectFraud, AssessRisk, EvaluateClaim)

**Key Messages**:
- AIAgent, AgentConfig, AgentPerformance
- AIAgentCreatedEvent, ConversationStartedEvent, AIDecisionMadeEvent, FraudDetectedEvent
- ChatRequest/Response, AnalyzeDocumentRequest/Response, DetectFraudRequest/Response

---

### 2. Analytics Module (`insuretech/analytics/`)
**Purpose**: Metrics, dashboards, reports, and KPI tracking

**Files**:
- `entity/v1/analytics.proto` - Dashboard, Widget, Report, Visualization definitions
- `entity/v1/metric.proto` - MetricType, AggregatedMetric definitions
- `events/v1/analytics_events.proto` - Analytics events (MetricRecorded, ReportGenerated, DashboardAccessed, KPIThresholdBreached)
- `services/v1/analytics_service.proto` - AnalyticsService with RPC methods (GetMetrics, GetDashboard, CreateDashboard, GenerateReport, ScheduleReport, RunQuery)

**Key Messages**:
- MetricRecordedEvent, ReportGeneratedEvent, DashboardAccessedEvent, KPIThresholdBreachedEvent
- GetMetricsRequest/Response, GenerateReportRequest/Response

---

### 3. API Key Module (`insuretech/apikey/`)
**Purpose**: API authentication and usage tracking for insurers and partners

**Files**:
- `entity/v1/api_key.proto` - ApiKey entity (with status, scopes, rate limits)
- `entity/v1/api_key_usage.proto` - ApiKeyUsage entity for tracking requests
- `events/v1/api_key_events.proto` - API key events (Created, Revoked, RateLimitExceeded, Expired)
- `services/v1/api_key_service.proto` - ApiKeyService with RPC methods (GenerateApiKey, RevokeApiKey, RotateApiKey, ValidateApiKey, GetUsageStats)

**Key Messages**:
- ApiKey (with rate limiting, IP whitelist, expiration)
- ApiKeyUsage (tracks endpoint, method, response time, IP, trace_id)
- ApiKeyCreatedEvent, ApiKeyRevokedEvent, ApiKeyRateLimitExceededEvent

---

### 4. Audit Module (`insuretech/audit/`)
**Purpose**: Audit trails, compliance logging, and regulatory compliance

**Files**:
- `entity/v1/audit_event.proto` - AuditEvent for business events
- `entity/v1/audit_log.proto` - AuditLog for system audit trail (FR-153 to FR-158)
- `entity/v1/compliance_log.proto` - ComplianceLog for regulatory compliance
- `events/v1/audit_events.proto` - Audit events
- `services/v1/audit_service.proto` - AuditService with RPC methods

**Key Messages**:
- AuditEvent (with category, severity, entity tracking)
- AuditLog (with old/new values, change tracking, IP, user_agent, trace_id)
- ComplianceLog (with regulation references, evidence)
- EventCategory: SECURITY, COMPLIANCE, BUSINESS, TECHNICAL
- EventSeverity: INFO, WARNING, ERROR, CRITICAL
- AuditAction: CREATE, READ, UPDATE, DELETE, LOGIN, LOGOUT, APPROVE, REJECT, EXPORT

---

### 5. Authentication Module (`insuretech/authn/`)
**Purpose**: User authentication, sessions, OTP, KYC documents

**Files**:
- `entity/v1/user.proto` - User entity with auth details
- `entity/v1/user_profile.proto` - User profile information
- `entity/v1/session.proto` - Session management
- `entity/v1/otp.proto` - One-Time Password handling
- `entity/v1/document_type.proto` - Document types for identity verification
- `entity/v1/user_document.proto` - User documents (NID, passport, etc.)
- `entity/v1/enums.proto` - Authentication enums
- `events/v1/core.proto` - Authentication events
- `services/v1/auth_service.proto` - AuthService (Login, Logout, Register, ValidateSession, RefreshToken, RequestOTP, VerifyOTP)
- `services/v1/core.proto` - Core auth service messages

**Key Entities**:
- User (with status, roles, last_login)
- UserProfile (personal info, preferences)
- Session (with expiration, refresh token)
- OTP (one-time password)

---

### 6. Authorization Module (`insuretech/authz/`)
**Purpose**: Role-based access control, policies, MFA configuration

**Files**:
- `entity/v1/role.proto` - Role definitions
- `entity/v1/user_role.proto` - User role assignments
- `entity/v1/policy_rule.proto` - Policy rules
- `entity/v1/casbin_rule.proto` - Casbin RBAC rules
- `entity/v1/access_decision_audit.proto` - Access decision tracking
- `entity/v1/role_mfa_config.proto` - MFA configuration per role
- `entity/v1/token_config.proto` - Token expiration config
- `entity/v1/portal_config.proto` - Portal configuration
- `entity/v1/enums.proto` - Authorization enums
- `events/v1/core.proto` - Authorization events
- `services/v1/authz_service.proto` - AuthzService (CheckAccess, ListRoles, AssignRole, RevokeRole, UpdatePolicy)
- `services/v1/core.proto` - Core authz messages

**Key Concepts**:
- RBAC with role hierarchy
- Casbin-based policy engine
- MFA support (OTP, SMS, Email, Biometric)
- Access decision audit trail

---

### 7. B2B Module (`insuretech/b2b/`)
**Purpose**: Business-to-business functionality (departments, employees, organizations)

**Files**:
- `entity/v1/organisation.proto` - Organization/Company details
- `entity/v1/department.proto` - Department structure
- `entity/v1/employee.proto` - Employee information
- `entity/v1/purchase_order.proto` - Purchase orders
- `events/v1/core.proto` - B2B events
- `services/v1/b2b_service.proto` - B2BService with RPC methods

**Key Entities**:
- Organisation (with registration, contact details)
- Department (with hierarchy)
- Employee (with roles, status)
- PurchaseOrder (for B2B transactions)

---

### 8. Beneficiary Module (`insuretech/beneficiary/`)
**Purpose**: Beneficiary management for policies

**Files**:
- `entity/v1/beneficiary.proto` - Base Beneficiary entity
- `entity/v1/individual.proto` - Individual beneficiary
- `entity/v1/business.proto` - Business beneficiary
- `events/v1/beneficiary_events.proto` - Beneficiary events
- `services/v1/beneficiary_service.proto` - BeneficiaryService

**Key Entities**:
- Beneficiary (base with relationship to policy holder)
- IndividualBeneficiary (with personal info)
- BusinessBeneficiary (with company details)

---

### 9. Billing Module (`insuretech/billing/`)
**Purpose**: Invoice management

**Files**:
- `entity/v1/invoice.proto` - Invoice entity with line items, taxes

---

### 10. Claims Module (`insuretech/claims/`)
**Purpose**: Claim management from submission to settlement

**Files**:
- `entity/v1/claim.proto` - Claim entity with status tracking
- `events/v1/claim_events.proto` - Claim events (Submitted, Approved, Rejected, Settled)
- `services/v1/claim_service.proto` - ClaimService (SubmitClaim, GetClaim, UpdateClaim, ApproveClaim, RejectClaim)

**Key Entities**:
- Claim (with policy reference, amount, status)
- ClaimStatus: DRAFT, SUBMITTED, UNDER_REVIEW, APPROVED, REJECTED, SETTLED, APPEALED

---

### 11. Commission Module (`insuretech/commission/`)
**Purpose**: Commission tracking and payout management

**Files**:
- `entity/v1/commission_config.proto` - Commission configuration
- `entity/v1/commission_payout.proto` - Commission payout details
- `entity/v1/revenue_share.proto` - Revenue sharing configuration
- `events/v1/commission_events.proto` - Commission events
- `services/v1/commission_service.proto` - CommissionService

**Key Entities**:
- CommissionConfig (rates, tiers)
- CommissionPayout (amount, status)
- RevenueShare (percentage, allocation)

---

### 12. Common Module (`insuretech/common/v1/`)
**Purpose**: Shared utilities and base definitions

**Files**:
- `db.proto` - Database options (TableOptions, ColumnOptions, ForeignKey, IndexOptions)
  - Extensions for message-level and field-level options
  - PartitionStrategy enum (NONE, RANGE_YEAR, RANGE_MONTH, LIST, HASH)
  - ReferentialAction enum (CASCADE, SET_NULL, RESTRICT, NO_ACTION)
  - IndexType enum (BTREE, HASH, GIN, GIST, BRIN)

- `error.proto` - Standard error handling
  - Error message with code, details, field_violations, retryable flag
  - FieldViolation for validation errors
  - Comprehensive ErrorCode enum (1000+ error codes)
  - ErrorSeverity enum (INFO, WARNING, ERROR, CRITICAL)

- `security.proto` - Security and privacy annotations
  - Field-level annotations: pii, encrypted_security, log_masked, log_redacted, sensitive, requires_consent
  - SecurityClassification enum (PUBLIC, INTERNAL, CONFIDENTIAL, HIGHLY_CONFIDENTIAL)
  - DataCategory enum (PERSONAL_IDENTIFIER, CONTACT_INFORMATION, FINANCIAL_INFORMATION, HEALTH_INFORMATION, etc.)
  - SecurityEvent message

- `types.proto` - Common data types
  - Money (amount, currency)
  - Address (street, city, country, postal_code, lat/lon)
  - Phone (country_code, number, type)
  - Email (address, verified)
  - Document (type, number, issue_date, expiry_date)
  - AuditInfo (created_by, created_at, updated_by, updated_at, deleted_at)
  - InsuranceType enum (HEALTH, LIFE, AUTO, HOME, TRAVEL, etc.)

---

### 13. Document Module (`insuretech/document/`)
**Purpose**: Document generation and templating

**Files**:
- `entity/v1/document_template.proto` - Document template definitions
- `entity/v1/document_generation.proto` - Document generation jobs
- `events/v1/document_events.proto` - Document events (Generated, Sent, Signed)
- `services/v1/document_service.proto` - DocumentService (GenerateDocument, GetTemplate, ListTemplates)

---

### 14. Endorsement Module (`insuretech/endorsement/`)
**Purpose**: Policy endorsements and amendments

**Files**:
- `entity/v1/endorsement.proto` - Endorsement entity
- `events/v1/endorsement_events.proto` - Endorsement events
- `services/v1/endorsement_service.proto` - EndorsementService

---

### 15. Fraud Module (`insuretech/fraud/`)
**Purpose**: Fraud detection and case management

**Files**:
- `entity/v1/fraud_alert.proto` - Fraud alert entity
- `entity/v1/fraud_case.proto` - Fraud case management
- `entity/v1/fraud_rule.proto` - Fraud detection rules
- `events/v1/fraud_events.proto` - Fraud events
- `services/v1/fraud_service.proto` - FraudService

**Key Entities**:
- FraudAlert (with severity, score)
- FraudCase (with status, investigation details)
- FraudRule (detection rules)

---

### 16. Insurance Service Module (`insuretech/insurance/`)
**Purpose**: High-level insurance operations

**Files**:
- `services/v1/insurance_service.proto` - InsuranceService

---

### 17. Insurer Module (`insuretech/insurer/`)
**Purpose**: Insurance company management

**Files**:
- `entity/v1/insurer.proto` - Insurer company details
- `entity/v1/insurer_config.proto` - Insurer configuration (underwriting rules, claim rules)
- `entity/v1/insurer_product.proto` - Insurer product offerings
- `events/v1/insurer_events.proto` - Insurer events
- `services/v1/insurer_service.proto` - InsurerService

---

### 18. IoT Module (`insuretech/iot/`)
**Purpose**: IoT device management

**Files**:
- `entity/v1/device.proto` - IoT device entity
- `events/v1/iot_events.proto` - IoT events (Connected, Disconnected, DataReceived)
- `services/v1/iot_service.proto` - IoTService

---

### 19. KYC Module (`insuretech/kyc/`)
**Purpose**: Know Your Customer verification

**Files**:
- `entity/v1/kyc_verification.proto` - KYC verification status
- `entity/v1/document_verification.proto` - Document verification details
- `events/v1/kyc_events.proto` - KYC events (Verified, Failed, Completed)
- `services/v1/kyc_service.proto` - KYCService (VerifyUser, VerifyDocument, GetVerificationStatus)

---

### 20. Media Module (`insuretech/media/`)
**Purpose**: Media file storage and processing

**Files**:
- `entity/v1/media.proto` - Media entity (with URL, type, size)
- `entity/v1/processing_job.proto` - Media processing jobs
- `events/v1/media_events.proto` - Media events (Uploaded, Processed)
- `services/v1/media_service.proto` - MediaService (UploadMedia, GetMedia, ProcessMedia)

---

### 21. MFS Module (`insuretech/mfs/`)
**Purpose**: Mobile Financial Services integration

**Files**:
- `entity/v1/mfs_integration.proto` - MFS provider integration config
- `entity/v1/mfs_transaction.proto` - MFS transaction details
- `entity/v1/mfs_webhook.proto` - MFS webhook handling
- `events/v1/mfs_events.proto` - MFS events
- `services/v1/mfs_service.proto` - MFSService

---

### 22. Notification Module (`insuretech/notification/`)
**Purpose**: Notifications and alerts

**Files**:
- `entity/v1/notification.proto` - Notification entity
- `entity/v1/alert.proto` - Alert entity
- `events/v1/notification_events.proto` - Notification events
- `services/v1/notification_service.proto` - NotificationService (SendNotification, SendAlert, GetNotifications)

---

### 23. Orders Module (`insuretech/orders/`)
**Purpose**: Order management

**Files**:
- `entity/v1/order.proto` - Order entity
- `events/v1/order_events.proto` - Order events
- `services/v1/order_service.proto` - OrderService

---

### 24. Partner Module (`insuretech/partner/`)
**Purpose**: Partner/channel management

**Files**:
- `entity/v1/partner.proto` - Partner entity
- `events/v1/partner_events.proto` - Partner events
- `services/v1/partner_service.proto` - PartnerService

---

### 25. Payment Module (`insuretech/payment/`)
**Purpose**: Payment processing

**Files**:
- `entity/v1/payment.proto` - Payment entity
- `events/v1/payment_events.proto` - Payment events (Initiated, Completed, Failed)
- `services/v1/payment_service.proto` - PaymentService (InitiatePayment, GetPayment, RefundPayment)

---

### 26. Policy Module (`insuretech/policy/`)
**Purpose**: Core policy management

**Files**:
- `entity/v1/policy.proto` - Policy entity (main insurance contract)
- `entity/v1/quotation.proto` - Quotation for policy
- `entity/v1/policy_service_request.proto` - Service requests on policies
- `events/v1/policy_events.proto` - Policy events (Created, Activated, Renewed, Cancelled)
- `services/v1/policy_service.proto` - PolicyService (CreatePolicy, GetPolicy, UpdatePolicy, CancelPolicy, RenewPolicy)

**Key Entities**:
- Policy (with coverage details, premium, status)
- PolicyStatus: DRAFT, QUOTED, ACTIVE, SUSPENDED, LAPSED, CANCELLED, RENEWED

---

### 27. Products Module (`insuretech/products/`)
**Purpose**: Insurance product definitions

**Files**:
- `entity/v1/product.proto` - Insurance product
- `entity/v1/product_plan.proto` - Product plan variants
- `entity/v1/rider.proto` - Riders/add-ons for products
- `entity/v1/pricing_config.proto` - Pricing configuration
- `events/v1/product_events.proto` - Product events
- `services/v1/product_service.proto` - ProductService (GetProduct, ListProducts, CreateProduct)

---

### 28. Refund Module (`insuretech/refund/`)
**Purpose**: Refund processing

**Files**:
- `entity/v1/refund.proto` - Refund entity
- `events/v1/refund_events.proto` - Refund events
- `services/v1/refund_service.proto` - RefundService

---

### 29. Renewal Module (`insuretech/renewal/`)
**Purpose**: Policy renewal management

**Files**:
- `entity/v1/renewal_schedule.proto` - Renewal schedule
- `entity/v1/renewal_reminder.proto` - Renewal reminders
- `entity/v1/grace_period.proto` - Grace period after expiry
- `events/v1/renewal_events.proto` - Renewal events
- `services/v1/renewal_service.proto` - RenewalService

---

### 30. Report Module (`insuretech/report/`)
**Purpose**: Business reporting

**Files**:
- `entity/v1/report_definition.proto` - Report definition
- `entity/v1/report_schedule.proto` - Report scheduling
- `entity/v1/report_execution.proto` - Report execution tracking
- `events/v1/report_events.proto` - Report events
- `services/v1/report_service.proto` - ReportService

---

### 31. Services Module (`insuretech/services/`)
**Purpose**: Service provider management

**Files**:
- `entity/v1/service_provider.proto` - Service provider entity

---

### 32. Storage Module (`insuretech/storage/`)
**Purpose**: File storage management

**Files**:
- `entity/v1/storage.proto` - Storage entity (S3, bucket references)
- `events/v1/storage_events.proto` - Storage events
- `service/v1/storage_service.proto` - StorageService (UploadFile, DownloadFile, DeleteFile)

---

### 33. Support Module (`insuretech/support/`)
**Purpose**: Customer support ticketing

**Files**:
- `entity/v1/ticket.proto` - Support ticket
- `entity/v1/ticket_message.proto` - Ticket messages/comments
- `entity/v1/faq.proto` - FAQ database
- `entity/v1/knowledge_base.proto` - Knowledge base articles
- `events/v1/support_events.proto` - Support events
- `services/v1/support_service.proto` - SupportService (CreateTicket, GetTicket, AddMessage, CloseTicket)

---

### 34. Task Module (`insuretech/task/`)
**Purpose**: Task management and workflow

**Files**:
- `entity/v1/task.proto` - Task entity
- `events/v1/task_events.proto` - Task events
- `services/v1/task_service.proto` - TaskService

---

### 35. Tenant Module (`insuretech/tenant/`)
**Purpose**: Multi-tenancy support

**Files**:
- `entity/v1/tenant.proto` - Tenant (organization) entity
- `entity/v1/tenant_config.proto` - Tenant configuration
- `events/v1/tenant_events.proto` - Tenant events
- `services/v1/tenant_service.proto` - TenantService

---

### 36. Underwriting Module (`insuretech/underwriting/`)
**Purpose**: Underwriting decisions and risk assessment

**Files**:
- `entity/v1/quote.proto` - Quote for underwriting
- `entity/v1/health_declaration.proto` - Health information for underwriting
- `entity/v1/underwriting_decision.proto` - Underwriting decision
- `events/v1/underwriting_events.proto` - Underwriting events
- `services/v1/underwriting_service.proto` - UnderwritingService (GetQuote, SubmitForUnderwriting, MakeDecision)

---

### 37. Voice Module (`insuretech/voice/`)
**Purpose**: Voice interaction and transcription

**Files**:
- `entity/v1/voice_session.proto` - Voice session
- `entity/v1/voice_command.proto` - Voice commands
- `entity/v1/voice_transcript.proto` - Voice transcript
- `events/v1/voice_events.proto` - Voice events
- `services/v1/voice_service.proto` - VoiceService

---

### 38. WebRTC Module (`insuretech/webrtc/v1/`)
**Purpose**: Real-time communication via WebRTC

**Files**:
- **Entity**:
  - `entity/peer.proto` - Peer connection entity
  - `entity/room.proto` - Conference room entity
  - `entity/session.proto` - WebRTC session
  - `entity/signal.proto` - Signaling messages
  - `entity/track.proto` - Media track

- **Events**:
  - `events/peer_events.proto` - Peer connection events
  - `events/room_events.proto` - Room events
  - `events/signal_events.proto` - Signaling events
  - `events/track_events.proto` - Track events

- **Services**:
  - `service/peer_service.proto` - PeerService (CreatePeer, ClosePeer, UpdatePeer)
  - `service/room_service.proto` - RoomService (CreateRoom, JoinRoom, LeaveRoom)
  - `service/signaling_service.proto` - SignalingService (SendSignal, ReceiveSignal)
  - `service/stats_service.proto` - StatsService (GetStats)
  - `service/track_service.proto` - TrackService (AddTrack, RemoveTrack)

---

### 39. Workflow Module (`insuretech/workflow/`)
**Purpose**: Workflow orchestration and automation

**Files**:
- `entity/v1/workflow_definition.proto` - Workflow definition (DAG-based)
- `entity/v1/workflow_instance.proto` - Workflow instance execution
- `entity/v1/workflow_task.proto` - Workflow task
- `entity/v1/workflow_config.proto` - Workflow configuration
- `events/v1/workflow_events.proto` - Workflow events (Started, Completed, Failed)
- `services/v1/workflow_service.proto` - WorkflowService (ExecuteWorkflow, GetStatus, CancelWorkflow)

---

## Common Patterns

### 1. Module Organization
Every main module has:
- **entity/v1/**: Data models with database annotations
- **events/v1/**: Event messages for event streaming (Kafka)
- **services/v1/**: gRPC service definitions with HTTP mappings

### 2. Standard Fields
Most entities include:
- `id`: UUID primary key
- Timestamps: `created_at`, `updated_at`
- Soft delete: `deleted_at`
- Audit info: `created_by`, `updated_by`

### 3. Event Patterns
Every event includes:
- `event_id`: Unique event identifier
- `correlation_id`: For tracing across services (CG-6)
- `timestamp`: Event occurrence time
- Entity-specific fields

### 4. Error Handling
All service RPC responses include optional `Error` field from `common/v1/error.proto`

### 5. Database Annotations
Fields are annotated with:
- Column options (name, type, constraints)
- Index specifications
- Foreign key references
- Soft delete and audit field support

### 6. Security Annotations
Sensitive fields are marked with:
- `pii`: Personally identifiable information
- `sensitive`: Highly sensitive data
- `log_masked`: Mask in logs
- `log_redacted`: Complete redaction
- `encrypted_security`: Encrypted at rest
- `data_purpose`: GDPR compliance

---

## Key Cross-Module Relationships

1. **Policy (Core)**: Referenced by Claims, Payments, Renewals, Endorsements
2. **User (Authentication)**: Referenced by all audit trails, audit logs
3. **Tenant (Multi-tenancy)**: All entities belong to a tenant
4. **Error (Common)**: Used in all service responses
5. **AuditInfo (Common)**: Standard audit trail in all entities
6. **Money/Address/Phone (Common)**: Standard types used across modules

---

## Service Layer Standards

### HTTP Conventions
- POST: Create/Action
- GET: Retrieve
- PUT/PATCH: Update
- DELETE: Delete

### RPC Naming
- `Get{Entity}`: Retrieve single entity
- `List{Entity}s`: List multiple entities
- `Create{Entity}`: Create entity
- `Update{Entity}`: Update entity
- `Delete{Entity}`: Delete entity

### Response Pattern
All responses include optional `error` field for consistent error handling

---

## Migration and Database Management

The `db.proto` file provides comprehensive database schema options:
- **Table-level**: Name, schema, migration order, soft delete, audit fields, RLS
- **Field-level**: Column name, type, constraints, indexes, foreign keys
- **Partitioning**: Support for range, list, and hash partitioning

---

## Security and Compliance

### Privacy Annotations (security.proto)
- Field-level privacy controls
- GDPR compliance support
- Data retention policies
- Purpose tracking for data collection

### Audit Trail
- Comprehensive audit logging (audit_log.proto)
- Compliance logging (compliance_log.proto)
- Event category and severity tracking
- User and IP address tracking

### Access Control
- RBAC with role hierarchy
- Casbin policy engine integration
- MFA configuration per role
- Access decision audit trail

---

## Event-Driven Architecture

All modules publish events to Kafka for:
- Real-time updates
- Service-to-service communication
- Event sourcing
- Audit trail

Events include correlation_id for distributed tracing across services.

---

Generated: Auto-generated summary of InsureTech proto directory structure
