# InsureTech Proto - Support & Infrastructure Modules

## Overview
This document details support services, compliance, and infrastructure modules.

---

# 1. Audit & Compliance Module (insuretech/audit/)

## Purpose
Comprehensive audit trails and regulatory compliance logging.

## Key Files
- `entity/v1/audit_event.proto` - Business event audit
- `entity/v1/audit_log.proto` - System audit trail (FR-153 to FR-158)
- `entity/v1/compliance_log.proto` - Regulatory compliance logging
- `events/v1/audit_events.proto` - Audit events
- `services/v1/audit_service.proto` - Audit RPC service

## Key Entities

### AuditEvent Entity (Business Events)
```
id (UUID)                  - Primary key
category                   - SECURITY, COMPLIANCE, BUSINESS, TECHNICAL
event_type                 - Specific event type
severity                   - INFO, WARNING, ERROR, CRITICAL
description                - Event description
user_id (FK)               - Reference to User who triggered event
entity_type                - Type of affected entity (policy, claim, etc.)
entity_id (FK)             - ID of affected entity
metadata                   - Additional event metadata (JSON)
timestamp                  - Event occurrence time
```

### AuditLog Entity (System Audit Trail)
```
audit_log_id (UUID)        - Primary key
entity_type                - Type of entity (user, policy, claim, etc.)
entity_id (FK)             - ID of entity
action                     - CREATE, READ, UPDATE, DELETE, LOGIN, LOGOUT, APPROVE, REJECT, EXPORT
user_id (FK)               - Reference to User performing action
user_email                 - User email (PII, log_masked)
user_role                  - User's role at time of action
old_values                 - State before change (JSON)
new_values                 - State after change (JSON)
changes                    - Field-level changes (JSON)
ip_address                 - Client IP (PII, log_masked)
user_agent                 - Client user agent
trace_id                   - Distributed tracing ID (FR-208)
timestamp                  - Action timestamp
```

### ComplianceLog Entity (Regulatory Compliance)
```
id (UUID)                  - Primary key
compliance_framework       - GDPR, HIPAA, CCPA, LOCAL, etc.
regulation_reference       - Reference to regulation/article
entity_type                - Type of entity
entity_id (FK)             - ID of entity
action                     - Compliance-related action
description                - Compliance action description
evidence                   - Evidence of compliance (JSON/file refs)
verified_by                - Verification user ID
verification_date          - Verification timestamp
status                     - PENDING, VERIFIED, FAILED, WAIVED
timestamp                  - Compliance log timestamp
```

## Key Enums

### EventCategory
- EVENT_CATEGORY_SECURITY: Security-related events
- EVENT_CATEGORY_COMPLIANCE: Regulatory compliance events
- EVENT_CATEGORY_BUSINESS: Business logic events
- EVENT_CATEGORY_TECHNICAL: System/technical events

### EventSeverity
- EVENT_SEVERITY_INFO: Informational
- EVENT_SEVERITY_WARNING: Warning level
- EVENT_SEVERITY_ERROR: Error level
- EVENT_SEVERITY_CRITICAL: Critical/emergency

### AuditAction
- AUDIT_ACTION_CREATE: Entity created
- AUDIT_ACTION_READ: Entity accessed
- AUDIT_ACTION_UPDATE: Entity modified
- AUDIT_ACTION_DELETE: Entity deleted
- AUDIT_ACTION_LOGIN: User login
- AUDIT_ACTION_LOGOUT: User logout
- AUDIT_ACTION_APPROVE: Action approved
- AUDIT_ACTION_REJECT: Action rejected
- AUDIT_ACTION_EXPORT: Data exported

## Audit Service RPC Methods

```protobuf
service AuditService {
  rpc GetAuditLog(GetAuditLogRequest) returns (GetAuditLogResponse);
  rpc ListAuditLogs(ListAuditLogsRequest) returns (ListAuditLogsResponse);
  rpc SearchAuditLogs(SearchAuditLogsRequest) returns (SearchAuditLogsResponse);
  rpc GetComplianceLog(GetComplianceLogRequest) returns (GetComplianceLogResponse);
  rpc ListComplianceLogs(ListComplianceLogsRequest) returns (ListComplianceLogsResponse);
  rpc ExportAuditReport(ExportAuditReportRequest) returns (ExportAuditReportResponse);
}
```

---

# 2. API Key Module (insuretech/apikey/)

## Purpose
API key management, authentication, and usage tracking (FR-207, FR-208).

## Key Files
- `entity/v1/api_key.proto` - API key entity
- `entity/v1/api_key_usage.proto` - Usage tracking
- `events/v1/api_key_events.proto` - API key events
- `services/v1/api_key_service.proto` - API Key RPC service

## Key Entities

### ApiKey Entity
```
id (UUID)                  - Primary key (api_key_id)
key_hash                   - SHA-256 hash of API key (encrypted, log_redacted)
name                       - Human-readable key name
owner_type                 - INSURER, PARTNER, INTERNAL
owner_id (FK)              - Reference to insurer or partner
scopes                     - Array of permission scopes
                             Examples: policy:read, policy:write, claim:read
status                     - ACTIVE, EXPIRED, REVOKED, SUSPENDED, ROTATING
rate_limit_per_minute      - Max requests per minute (default: 60)
expires_at                 - Key expiration date
last_used_at               - Last usage timestamp
ip_whitelist               - Array of allowed IPs (empty = all)
audit_info                 - Standard audit trail
```

### ApiKeyUsage Entity (FR-207)
```
id (UUID)                  - Primary key (usage_id)
api_key_id (FK)            - Reference to ApiKey
endpoint                   - Called endpoint (e.g., /v1/policies)
http_method                - GET, POST, PUT, DELETE, PATCH
status_code                - HTTP status code returned
response_time_ms           - Response time in milliseconds
request_ip                 - Requester IP (PII, log_masked)
user_agent                 - Client user agent
request_payload            - Request body (JSONB, for audit)
response_payload           - Response body (JSONB, for audit)
trace_id                   - Distributed tracing ID (FR-208)
timestamp                  - Request timestamp
```

## Key Enums

### ApiKeyOwnerType
- API_KEY_OWNER_TYPE_INSURER: Owned by insurance company
- API_KEY_OWNER_TYPE_PARTNER: Owned by partner/channel
- API_KEY_OWNER_TYPE_INTERNAL: Internal system key

### ApiKeyStatus
- API_KEY_STATUS_ACTIVE: Key is usable
- API_KEY_STATUS_EXPIRED: Key has expired
- API_KEY_STATUS_REVOKED: Key was revoked
- API_KEY_STATUS_SUSPENDED: Key is temporarily suspended
- API_KEY_STATUS_ROTATING: Key rotation in progress (grace period)

## API Key Service RPC Methods

```protobuf
service ApiKeyService {
  rpc GenerateApiKey(GenerateApiKeyRequest) returns (GenerateApiKeyResponse);
  rpc RevokeApiKey(RevokeApiKeyRequest) returns (RevokeApiKeyResponse);
  rpc RotateApiKey(RotateApiKeyRequest) returns (RotateApiKeyResponse);
  rpc GetApiKey(GetApiKeyRequest) returns (GetApiKeyResponse);
  rpc ListApiKeys(ListApiKeysRequest) returns (ListApiKeysResponse);
  rpc ValidateApiKey(ValidateApiKeyRequest) returns (ValidateApiKeyResponse);
  rpc GetUsageStats(GetUsageStatsRequest) returns (GetUsageStatsResponse);
}
```

### Key Security Features
- Keys are never stored in plaintext (stored as SHA-256 hash)
- Keys shown only once at generation
- Rate limiting prevents abuse
- IP whitelist for additional security
- Audit trail of all API usage (FR-207)
- Distributed tracing support (FR-208)

---

# 3. KYC Module (insuretech/kyc/)

## Purpose
Know Your Customer (KYC) verification and compliance.

## Key Files
- `entity/v1/kyc_verification.proto` - KYC verification status
- `entity/v1/document_verification.proto` - Document-level verification
- `events/v1/kyc_events.proto` - KYC events
- `services/v1/kyc_service.proto` - KYC RPC service

## Key Entities

### KYCVerification Entity
```
id (UUID)                  - Primary key
user_id (FK)               - Reference to User
verification_type          - BASIC, STANDARD, ENHANCED
status                     - PENDING, IN_PROGRESS, VERIFIED, FAILED, EXPIRED, SUSPENDED
verification_level        - 0 (not verified) to 5 (fully verified)
risk_score                - KYC risk score (0-100)
risk_category             - LOW, MEDIUM, HIGH
pep_check                 - Is user a PEP (Politically Exposed Person)?
aml_check                 - AML (Anti-Money Laundering) check status
sanctions_check           - Sanctions list check status
documents_verified        - Number of verified documents
documents_required        - Number of documents required
verification_date         - Date of verification
expiry_date              - Verification expiry date
verified_by              - Reference to verification agent
notes                    - Verification notes
metadata                 - Additional metadata (JSON)
tenant_id (FK)           - Reference to Tenant
audit_info               - Standard audit trail
```

### DocumentVerification Entity
```
id (UUID)                  - Primary key
kyc_verification_id (FK)   - Reference to KYCVerification
document_id (FK)           - Reference to UserDocument
document_type              - NID, PASSPORT, DRIVING_LICENSE, etc.
verification_status        - PENDING, VERIFIED, REJECTED, EXPIRED
verification_method        - MANUAL, AUTOMATED, BIOMETRIC_MATCH
confidence_score           - Verification confidence (0-100)
verified_by               - Verification agent ID
verification_date         - Verification timestamp
expiry_date              - Document validity check
issues_found             - Any issues found (array)
front_image_url          - Front side image URL
back_image_url           - Back side image URL
extracted_data           - Data extracted from document (JSON)
```

## KYC Service RPC Methods

```protobuf
service KYCService {
  rpc InitiateKYC(InitiateKYCRequest) returns (InitiateKYCResponse);
  rpc VerifyUser(VerifyUserRequest) returns (VerifyUserResponse);
  rpc VerifyDocument(VerifyDocumentRequest) returns (VerifyDocumentResponse);
  rpc GetVerificationStatus(GetVerificationStatusRequest) returns (GetVerificationStatusResponse);
  rpc UpdateVerification(UpdateVerificationRequest) returns (UpdateVerificationResponse);
  rpc ReverifyUser(ReverifyUserRequest) returns (ReverifyUserResponse);
  rpc CheckSanctions(CheckSanctionsRequest) returns (CheckSanctionsResponse);
}
```

---

# 4. Support Module (insuretech/support/)

## Purpose
Customer support ticketing, FAQ, and knowledge base.

## Key Files
- `entity/v1/ticket.proto` - Support ticket
- `entity/v1/ticket_message.proto` - Ticket messages/comments
- `entity/v1/faq.proto` - FAQ entries
- `entity/v1/knowledge_base.proto` - Knowledge base articles
- `events/v1/support_events.proto` - Support events
- `services/v1/support_service.proto` - Support RPC service

## Key Entities

### Ticket Entity
```
id (UUID)                  - Primary key (ticket_id)
ticket_number              - Unique ticket reference
customer_id (FK)           - Reference to customer User
agent_id (FK)              - Assigned support agent ID
category                   - BILLING, CLAIMS, POLICY, TECHNICAL, GENERAL
priority                   - LOW, MEDIUM, HIGH, CRITICAL
status                     - OPEN, IN_PROGRESS, WAITING_CUSTOMER, RESOLVED, CLOSED, REOPENED
subject                    - Ticket subject
description                - Initial problem description
resolution                 - Resolution provided
resolution_date            - Resolution timestamp
satisfaction_rating        - Customer satisfaction (1-5)
satisfaction_comment       - Customer comment
created_at                 - Creation timestamp
updated_at                 - Last update timestamp
resolved_at                - Resolution timestamp
closed_at                  - Closure timestamp
related_policy_id (FK)     - Related policy (if applicable)
related_claim_id (FK)      - Related claim (if applicable)
tags                       - Array of tags for categorization
attachments                - Array of attachment URLs
metadata                   - Additional metadata (JSON)
tenant_id (FK)             - Reference to Tenant
audit_info                 - Standard audit trail
```

### TicketMessage Entity
```
id (UUID)                  - Primary key
ticket_id (FK)             - Reference to Ticket
sender_id (FK)             - Reference to User (customer or agent)
sender_type                - CUSTOMER, AGENT, SYSTEM
message_text               - Message content
attachments                - Array of attachment URLs
is_internal               - Is this an internal note?
created_at                 - Creation timestamp
edited_at                  - Last edit timestamp
edited_by                  - User who edited
read_by                    - Array of users who read message
```

### FAQ Entity
```
id (UUID)                  - Primary key
question                   - FAQ question
answer                     - FAQ answer (supports markdown)
category                   - FAQ category
search_keywords            - Array of search keywords
helpful_count              - Times marked as helpful
unhelpful_count            - Times marked as unhelpful
view_count                 - Number of views
status                     - PUBLISHED, DRAFT, ARCHIVED
created_by                 - Created by user ID
updated_by                 - Updated by user ID
created_at                 - Creation timestamp
updated_at                 - Last update timestamp
related_articles           - Array of related article IDs
priority                   - Display priority
language                   - Language code
```

### KnowledgeBase Entity
```
id (UUID)                  - Primary key
title                      - Article title
slug                       - URL-friendly slug
content                    - Article content (markdown)
category                   - Knowledge base category
subcategory                - Subcategory
tags                       - Array of tags
author_id                  - Article author
status                     - PUBLISHED, DRAFT, ARCHIVED
view_count                 - Number of views
helpful_count              - Times marked as helpful
created_at                 - Creation timestamp
updated_at                 - Last update timestamp
published_at               - Publication timestamp
related_articles           - Array of related article IDs
language                   - Language code
attachments                - Array of attachment URLs
```

## Support Service RPC Methods

```protobuf
service SupportService {
  rpc CreateTicket(CreateTicketRequest) returns (CreateTicketResponse);
  rpc GetTicket(GetTicketRequest) returns (GetTicketResponse);
  rpc ListTickets(ListTicketsRequest) returns (ListTicketsResponse);
  rpc UpdateTicket(UpdateTicketRequest) returns (UpdateTicketResponse);
  rpc AddMessage(AddMessageRequest) returns (AddMessageResponse);
  rpc CloseTicket(CloseTicketRequest) returns (CloseTicketResponse);
  rpc SearchFAQ(SearchFAQRequest) returns (SearchFAQResponse);
  rpc GetArticle(GetArticleRequest) returns (GetArticleResponse);
  rpc RateTicket(RateTicketRequest) returns (RateTicketResponse);
}
```

---

# 5. Notification Module (insuretech/notification/)

## Purpose
Multi-channel notifications and alerts.

## Key Files
- `entity/v1/notification.proto` - Notification entity
- `entity/v1/alert.proto` - Alert entity
- `events/v1/notification_events.proto` - Notification events
- `services/v1/notification_service.proto` - Notification RPC service

## Key Entities

### Notification Entity
```
id (UUID)                  - Primary key
recipient_id (FK)          - Reference to recipient User
notification_type          - EMAIL, SMS, PUSH, IN_APP, WEBHOOK
channel                    - Communication channel
subject                    - Notification subject (for email)
body                       - Notification body/message
data                       - Additional data (JSON)
status                     - PENDING, SENT, DELIVERED, FAILED, BOUNCED
delivery_timestamp         - When message was delivered
read                       - Has recipient read notification?
read_timestamp             - When message was read
priority                   - LOW, NORMAL, HIGH, URGENT
retry_count                - Number of delivery attempts
max_retries                - Maximum retry attempts
next_retry_at              - When to retry next
error_message              - Error details (if failed)
template_id                - Email/SMS template used
template_variables         - Template variable values (JSON)
related_entity_type        - Type of related entity (policy, claim)
related_entity_id          - ID of related entity
created_at                 - Creation timestamp
expires_at                 - Notification expiry (for in-app)
metadata                   - Additional metadata (JSON)
tenant_id (FK)             - Reference to Tenant
```

### Alert Entity
```
id (UUID)                  - Primary key
alert_type                 - POLICY, CLAIM, PAYMENT, COMPLIANCE, SYSTEM
severity                   - INFO, WARNING, ERROR, CRITICAL
title                      - Alert title
description                - Alert description
affected_user_ids          - Array of affected user IDs
affected_entity_type       - Entity type affected
affected_entity_id         - Entity ID affected
action_required            - Is action required?
required_action_type       - Type of action needed
action_deadline            - Deadline for action (if required)
action_completed           - Has action been taken?
status                     - ACTIVE, RESOLVED, DISMISSED, EXPIRED
created_at                 - Creation timestamp
resolved_at                - Resolution timestamp
metadata                   - Additional metadata (JSON)
tenant_id (FK)             - Reference to Tenant
```

## Notification Service RPC Methods

```protobuf
service NotificationService {
  rpc SendNotification(SendNotificationRequest) returns (SendNotificationResponse);
  rpc SendAlert(SendAlertRequest) returns (SendAlertResponse);
  rpc GetNotification(GetNotificationRequest) returns (GetNotificationResponse);
  rpc ListNotifications(ListNotificationsRequest) returns (ListNotificationsResponse);
  rpc MarkAsRead(MarkAsReadRequest) returns (MarkAsReadResponse);
  rpc GetAlerts(GetAlertsRequest) returns (GetAlertsResponse);
  rpc DismissAlert(DismissAlertRequest) returns (DismissAlertResponse);
}
```

---

# 6. Document Module (insuretech/document/)

## Purpose
Document generation, templating, and management.

## Key Files
- `entity/v1/document_template.proto` - Document template
- `entity/v1/document_generation.proto` - Document generation job
- `events/v1/document_events.proto` - Document events
- `services/v1/document_service.proto` - Document RPC service

## Key Entities

### DocumentTemplate Entity
```
id (UUID)                  - Primary key
name                       - Template name
type                       - POLICY, QUOTE, CLAIM_FORM, RECEIPT, etc.
content                    - Template content (HTML/PDF template)
format                     - HTML, PDF, DOCX, MARKDOWN
variables                  - Array of variable placeholders
description                - Template description
status                     - ACTIVE, INACTIVE, ARCHIVED
version                    - Version number
language                   - Template language
insurer_id (FK)            - Reference to Insurer (if specific)
created_by                 - Created by user ID
updated_by                 - Updated by user ID
created_at                 - Creation timestamp
updated_at                 - Last update timestamp
preview_url                - Preview/sample document
metadata                   - Additional metadata (JSON)
```

### DocumentGeneration Entity
```
id (UUID)                  - Primary key
generation_id              - Unique generation reference
template_id (FK)           - Reference to DocumentTemplate
document_type              - Document type
related_entity_type        - Entity type (policy, claim, etc.)
related_entity_id          - Entity ID
generated_for_user         - User for whom document is generated
variables                  - Template variables used (JSON)
status                     - PENDING, IN_PROGRESS, COMPLETED, FAILED
document_url               - Generated document URL
file_size_bytes            - Generated file size
generation_time_ms         - Generation duration
error_message              - Error details (if failed)
retry_count                - Number of retry attempts
created_at                 - Request timestamp
completed_at               - Completion timestamp
expires_at                 - Document expiry (if applicable)
download_count             - Number of downloads
last_downloaded_at         - Last download timestamp
metadata                   - Additional metadata (JSON)
tenant_id (FK)             - Reference to Tenant
```

## Document Service RPC Methods

```protobuf
service DocumentService {
  rpc GenerateDocument(GenerateDocumentRequest) returns (GenerateDocumentResponse);
  rpc GetTemplate(GetTemplateRequest) returns (GetTemplateResponse);
  rpc ListTemplates(ListTemplatesRequest) returns (ListTemplatesResponse);
  rpc CreateTemplate(CreateTemplateRequest) returns (CreateTemplateResponse);
  rpc UpdateTemplate(UpdateTemplateRequest) returns (UpdateTemplateResponse);
  rpc PreviewDocument(PreviewDocumentRequest) returns (PreviewDocumentResponse);
  rpc DownloadDocument(DownloadDocumentRequest) returns (DownloadDocumentResponse);
}
```

---

# 7. Media Module (insuretech/media/)

## Purpose
Media file storage, retrieval, and processing.

## Key Files
- `entity/v1/media.proto` - Media entity
- `entity/v1/processing_job.proto` - Media processing job
- `events/v1/media_events.proto` - Media events
- `services/v1/media_service.proto` - Media RPC service

## Key Entities

### Media Entity
```
id (UUID)                  - Primary key
file_name                  - Original file name
file_type                  - File MIME type (image/pdf, etc.)
file_size_bytes            - File size in bytes
storage_path               - S3 path or similar
url                        - Public/private URL
upload_by                  - User who uploaded
upload_date                - Upload timestamp
status                     - UPLOADED, PROCESSING, AVAILABLE, DELETED, QUARANTINED
virus_scan_status          - PENDING, CLEAN, INFECTED
document_type              - Type of document (policy, claim, etc.)
related_entity_type        - Related entity type
related_entity_id          - Related entity ID
tags                       - Array of tags
metadata                   - Additional metadata (JSON)
retention_period_days      - Data retention period
expires_at                 - Expiry date
storage_class              - Storage class (standard, archive, etc.)
version_control            - File version
encryption_key_id          - Encryption key reference
tenant_id (FK)             - Reference to Tenant
```

### ProcessingJob Entity
```
id (UUID)                  - Primary key
job_id                     - Unique job reference
media_id (FK)              - Reference to Media
processing_type            - IMAGE_RESIZE, PDF_EXTRACT, OCR, COMPRESSION, etc.
input_parameters           - Processing parameters (JSON)
status                     - QUEUED, IN_PROGRESS, COMPLETED, FAILED
progress_percentage        - Job progress (0-100)
output_media_id (FK)       - Reference to output Media (if created)
error_message              - Error details (if failed)
started_at                 - Processing start timestamp
completed_at               - Processing completion timestamp
duration_ms                - Processing duration
retry_count                - Number of retries
priority                   - Job priority
result_metadata            - Result metadata (JSON)
```

## Media Service RPC Methods

```protobuf
service MediaService {
  rpc UploadMedia(UploadMediaRequest) returns (UploadMediaResponse);
  rpc GetMedia(GetMediaRequest) returns (GetMediaResponse);
  rpc ListMedia(ListMediaRequest) returns (ListMediaResponse);
  rpc ProcessMedia(ProcessMediaRequest) returns (ProcessMediaResponse);
  rpc DeleteMedia(DeleteMediaRequest) returns (DeleteMediaResponse);
  rpc GetProcessingStatus(GetProcessingStatusRequest) returns (GetProcessingStatusResponse);
  rpc GeneratePresignedUrl(GeneratePresignedUrlRequest) returns (GeneratePresignedUrlResponse);
}
```

---

# 8. Storage Module (insuretech/storage/)

## Purpose
File storage management and operations.

## Key Files
- `entity/v1/storage.proto` - Storage entity
- `events/v1/storage_events.proto` - Storage events
- `service/v1/storage_service.proto` - Storage RPC service

## Key Entities

### Storage Entity
```
id (UUID)                  - Primary key
storage_name               - Storage name/bucket
storage_type               - S3, GCS, AZURE_BLOB, LOCAL_FILESYSTEM
provider                   - Cloud provider
region                     - Storage region
bucket_name                - Bucket/container name
access_level               - PUBLIC, PRIVATE, AUTHENTICATED
encryption_enabled         - Is encryption enabled?
encryption_algorithm       - Encryption method
lifecycle_policy           - Data lifecycle rules (JSON)
total_size_bytes           - Total storage used
max_size_bytes             - Maximum storage limit
quotas                     - Storage quotas per entity type (JSON)
status                     - ACTIVE, INACTIVE, ARCHIVED
created_at                 - Creation timestamp
updated_at                 - Last update timestamp
metadata                   - Additional metadata (JSON)
tenant_id (FK)             - Reference to Tenant
```

## Storage Service RPC Methods

```protobuf
service StorageService {
  rpc UploadFile(UploadFileRequest) returns (UploadFileResponse);
  rpc DownloadFile(DownloadFileRequest) returns (stream DownloadFileResponse);
  rpc DeleteFile(DeleteFileRequest) returns (DeleteFileResponse);
  rpc ListFiles(ListFilesRequest) returns (ListFilesResponse);
  rpc GetFileMetadata(GetFileMetadataRequest) returns (GetFileMetadataResponse);
  rpc GeneratePresignedUrl(GeneratePresignedUrlRequest) returns (GeneratePresignedUrlResponse);
  rpc GetStorageStats(GetStorageStatsRequest) returns (GetStorageStatsResponse);
}
```

---

End of Support & Infrastructure Modules Reference
