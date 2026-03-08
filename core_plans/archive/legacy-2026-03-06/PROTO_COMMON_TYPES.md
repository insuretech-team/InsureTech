# InsureTech Proto - Common Types Reference

## Overview
This document contains the complete reference for common types used across all InsureTech services.

---

## db.proto - Database Schema Annotations

### TableOptions (Message-level Options)
Controls database table generation and management:

```protobuf
message TableOptions {
  string table_name = 1;           // Table name in database
  string schema_name = 2;          // Schema name (empty = public)
  int32 migration_order = 3;       // Migration order (dependency management)
  bool is_table = 4;               // Whether to create table
  string comment = 5;              // Table description
  bool soft_delete = 6;            // Enable soft delete (deleted_at)
  bool audit_fields = 7;           // Enable audit fields (created_at, updated_at, created_by, updated_by)
  bool enable_rls = 8;             // Enable row-level security
  PartitionStrategy partition_strategy = 9;
  string partition_column = 10;    // For time-based partitioning
}

enum PartitionStrategy {
  PARTITION_STRATEGY_UNSPECIFIED = 0;
  PARTITION_STRATEGY_NONE = 1;
  PARTITION_STRATEGY_RANGE_YEAR = 2;
  PARTITION_STRATEGY_RANGE_MONTH = 3;
  PARTITION_STRATEGY_LIST = 4;
  PARTITION_STRATEGY_HASH = 5;
}
```

### ColumnOptions (Field-level Options)
Controls database column generation and constraints:

```protobuf
message ColumnOptions {
  string column_name = 1;          // Column name in database
  string sql_type = 2;             // SQL type (VARCHAR(255), TEXT, etc.)
  bool not_null = 3;               // NOT NULL constraint
  bool unique = 4;                 // UNIQUE constraint
  bool primary_key = 5;            // PRIMARY KEY
  bool auto_increment = 6;         // AUTO INCREMENT / SERIAL
  string default_value = 7;        // DEFAULT value (SQL expression)
  string check_constraint = 8;     // CHECK constraint
  ForeignKey foreign_key = 9;      // Foreign key reference
  IndexOptions index = 10;         // Index specification
  string comment = 11;             // Column description
  bool exclude_from_insert = 12;   // Exclude from INSERT operations
  bool exclude_from_update = 13;   // Exclude from UPDATE operations
  bool encrypted = 14;             // Encrypt at rest
  bool is_json = 15;               // JSON/JSONB column
}
```

### ForeignKey Definition
Specifies foreign key constraints:

```protobuf
message ForeignKey {
  string references_table = 1;     // Referenced table name
  string references_column = 2;    // Referenced column name
  string references_schema = 3;    // Referenced schema (empty = same)
  ReferentialAction on_delete = 4; // ON DELETE action
  ReferentialAction on_update = 5; // ON UPDATE action
  string constraint_name = 6;      // Constraint name
}

enum ReferentialAction {
  REFERENTIAL_ACTION_UNSPECIFIED = 0;
  REFERENTIAL_ACTION_NO_ACTION = 1;
  REFERENTIAL_ACTION_RESTRICT = 2;
  REFERENTIAL_ACTION_CASCADE = 3;
  REFERENTIAL_ACTION_SET_NULL = 4;
  REFERENTIAL_ACTION_SET_DEFAULT = 5;
}
```

### IndexOptions
Specifies index configuration:

```protobuf
message IndexOptions {
  string index_name = 1;           // Index name
  IndexType index_type = 2;        // Index type
  bool unique = 3;                 // Is unique index
  repeated string composite_fields = 4; // Fields for composite index
  string index_method = 5;         // btree, hash, gin, gist, etc.
  string where_clause = 6;         // Partial index WHERE clause
}

enum IndexType {
  INDEX_TYPE_UNSPECIFIED = 0;
  INDEX_TYPE_NONE = 1;
  INDEX_TYPE_BTREE = 2;            // Default for most columns
  INDEX_TYPE_HASH = 3;             // For equality searches
  INDEX_TYPE_GIN = 4;              // For JSON, arrays, full-text search
  INDEX_TYPE_GIST = 5;             // For spatial data
  INDEX_TYPE_BRIN = 6;             // For large tables with ordering
}
```

---

## error.proto - Error Handling

### Error Message
Standard error response for all services:

```protobuf
message Error {
  string code = 1;                              // Machine-readable code (UPPER_SNAKE_CASE)
  string message = 2;                          // Human-readable message
  map<string, string> details = 3;             // Key-value pairs for additional details
  repeated FieldViolation field_violations = 4; // Field-specific validation errors
  bool retryable = 5;                          // Can operation be retried?
  int32 retry_after_seconds = 6;               // Suggested retry delay
  int32 http_status_code = 7;                  // HTTP status code equivalent
  string error_id = 8;                         // Unique error instance ID
  string documentation_url = 9;                // Link to error documentation
}
```

### FieldViolation
Field-level validation error:

```protobuf
message FieldViolation {
  string field = 1;                // Field path (e.g., "applicant.date_of_birth")
  string code = 2;                 // Field-specific error code
  string description = 3;          // Human-readable description
  string rejected_value = 4;       // The invalid value (if safe)
}
```

### ErrorCode Enum
Comprehensive error codes for all scenarios:

```protobuf
enum ErrorCode {
  // Generic errors (1000-1099)
  ERROR_CODE_INTERNAL_ERROR = 1000;
  ERROR_CODE_INVALID_REQUEST = 1001;
  ERROR_CODE_UNAUTHORIZED = 1002;
  ERROR_CODE_FORBIDDEN = 1003;
  ERROR_CODE_NOT_FOUND = 1004;
  ERROR_CODE_ALREADY_EXISTS = 1005;
  ERROR_CODE_PRECONDITION_FAILED = 1006;
  ERROR_CODE_RATE_LIMIT_EXCEEDED = 1007;
  ERROR_CODE_SERVICE_UNAVAILABLE = 1008;
  ERROR_CODE_TIMEOUT = 1009;
  
  // Validation errors (1100-1199)
  ERROR_CODE_VALIDATION_ERROR = 1100;
  ERROR_CODE_MISSING_REQUIRED_FIELD = 1101;
  ERROR_CODE_INVALID_FIELD_VALUE = 1102;
  ERROR_CODE_INVALID_FIELD_FORMAT = 1103;
  ERROR_CODE_FIELD_OUT_OF_RANGE = 1104;
  ERROR_CODE_INVALID_ENUM_VALUE = 1105;
  
  // Authentication errors (1200-1299)
  ERROR_CODE_AUTHENTICATION_FAILED = 1200;
  ERROR_CODE_INVALID_CREDENTIALS = 1201;
  ERROR_CODE_EXPIRED_TOKEN = 1202;
  ERROR_CODE_INVALID_TOKEN = 1203;
  ERROR_CODE_SESSION_EXPIRED = 1204;
  ERROR_CODE_OTP_EXPIRED = 1205;
  ERROR_CODE_OTP_INVALID = 1206;
  
  // Authorization errors (1300-1399)
  ERROR_CODE_PERMISSION_DENIED = 1300;
  ERROR_CODE_INSUFFICIENT_PRIVILEGES = 1301;
  ERROR_CODE_RESOURCE_ACCESS_DENIED = 1302;
  
  // Business logic errors (1400-1499)
  ERROR_CODE_BUSINESS_RULE_VIOLATION = 1400;
  ERROR_CODE_INVALID_STATE_TRANSITION = 1401;
  ERROR_CODE_OPERATION_NOT_ALLOWED = 1402;
  ERROR_CODE_QUOTA_EXCEEDED = 1403;
  ERROR_CODE_DUPLICATE_OPERATION = 1404;
  
  // Policy-specific errors (2000-2099)
  ERROR_CODE_POLICY_NOT_FOUND = 2000;
  ERROR_CODE_POLICY_ALREADY_CANCELLED = 2001;
  ERROR_CODE_POLICY_ALREADY_LAPSED = 2002;
  ERROR_CODE_POLICY_NOT_ACTIVE = 2003;
  ERROR_CODE_INVALID_POLICY_STATUS = 2004;
  
  // Claim-specific errors (2100-2199)
  ERROR_CODE_CLAIM_NOT_FOUND = 2100;
  ERROR_CODE_CLAIM_ALREADY_SETTLED = 2101;
  ERROR_CODE_CLAIM_AMOUNT_EXCEEDS_COVERAGE = 2102;
  ERROR_CODE_CLAIM_OUTSIDE_COVERAGE_PERIOD = 2103;
  
  // Payment-specific errors (2200-2299)
  ERROR_CODE_PAYMENT_NOT_FOUND = 2200;
  ERROR_CODE_PAYMENT_ALREADY_COMPLETED = 2201;
  ERROR_CODE_PAYMENT_FAILED = 2202;
  ERROR_CODE_INSUFFICIENT_FUNDS = 2203;
  ERROR_CODE_PAYMENT_GATEWAY_ERROR = 2204;
  ERROR_CODE_INVALID_PAYMENT_METHOD = 2205;
  
  // Underwriting errors (2300-2399)
  ERROR_CODE_QUOTE_NOT_FOUND = 2300;
  ERROR_CODE_QUOTE_EXPIRED = 2301;
  ERROR_CODE_UNDERWRITING_DECLINED = 2302;
  ERROR_CODE_MEDICAL_EXAM_REQUIRED = 2303;
  
  // KYC/Verification errors (2400-2499)
  ERROR_CODE_KYC_VERIFICATION_FAILED = 2400;
  ERROR_CODE_INVALID_NID = 2401;
  ERROR_CODE_DOCUMENT_VERIFICATION_FAILED = 2402;
  ERROR_CODE_BIOMETRIC_VERIFICATION_FAILED = 2403;
  
  // Third-party integration errors (2500-2599)
  ERROR_CODE_EXTERNAL_SERVICE_ERROR = 2500;
  ERROR_CODE_MFS_PROVIDER_ERROR = 2501;
  ERROR_CODE_SMS_GATEWAY_ERROR = 2502;
  ERROR_CODE_EMAIL_SERVICE_ERROR = 2503;
}

enum ErrorSeverity {
  ERROR_SEVERITY_UNSPECIFIED = 0;
  ERROR_SEVERITY_INFO = 1;       // Informational, no action needed
  ERROR_SEVERITY_WARNING = 2;    // Warning, operation succeeded with issues
  ERROR_SEVERITY_ERROR = 3;      // Error, operation failed
  ERROR_SEVERITY_CRITICAL = 4;   // Critical error, system-level issue
}
```

---

## security.proto - Security & Privacy Annotations

### Field-Level Annotations
Used to mark sensitive fields in protobuf messages:

```protobuf
extend google.protobuf.FieldOptions {
  bool pii = 50010;                    // Personally Identifiable Information (PII)
  bool encrypted_security = 50011;     // Encrypt at rest
  bool log_masked = 50012;             // Mask in logs (e.g., 017****5678)
  bool log_redacted = 50013;           // Complete redaction (e.g., [REDACTED])
  bool sensitive = 50014;              // Highly sensitive data
  bool requires_consent = 50015;       // GDPR consent required
  string data_purpose = 50016;         // Purpose of data collection
  int32 retention_days = 50017;        // Retention period (0 = indefinite)
}
```

### Security Classifications

```protobuf
enum SecurityClassification {
  SECURITY_CLASSIFICATION_UNSPECIFIED = 0;
  SECURITY_CLASSIFICATION_PUBLIC = 1;              // Public data, no restrictions
  SECURITY_CLASSIFICATION_INTERNAL = 2;           // Internal use only
  SECURITY_CLASSIFICATION_CONFIDENTIAL = 3;       // Restricted access
  SECURITY_CLASSIFICATION_HIGHLY_CONFIDENTIAL = 4; // Strictly controlled
}
```

### Data Categories

```protobuf
enum DataCategory {
  DATA_CATEGORY_UNSPECIFIED = 0;
  DATA_CATEGORY_PERSONAL_IDENTIFIER = 1;      // Name, NID, passport
  DATA_CATEGORY_CONTACT_INFORMATION = 2;      // Email, phone, address
  DATA_CATEGORY_FINANCIAL_INFORMATION = 3;    // Account numbers, transactions
  DATA_CATEGORY_HEALTH_INFORMATION = 4;       // Medical records, health data
  DATA_CATEGORY_AUTHENTICATION_CREDENTIALS = 5; // Passwords, tokens, keys
  DATA_CATEGORY_BIOMETRIC_DATA = 6;           // Fingerprints, face recognition
  DATA_CATEGORY_LOCATION_DATA = 7;            // GPS, IP addresses
}
```

### SecurityEvent Message
For audit logging of security-related events:

```protobuf
message SecurityEvent {
  string event_id = 1;
  string event_type = 2;        // ACCESS, MODIFICATION, DELETION, EXPORT
  string user_id = 3;
  string resource_type = 4;     // USER, POLICY, CLAIM
  string resource_id = 5;
  string action = 6;            // READ, UPDATE, DELETE
  bool authorized = 7;
  string ip_address = 8;
  string user_agent = 9;
  google.protobuf.Timestamp timestamp = 10;
  map<string, string> metadata = 11;
}
```

---

## types.proto - Common Data Types

### Money Type
Represents monetary amounts with currency:

```protobuf
message Money {
  string currency_code = 1;     // ISO 4217 code (USD, BDT, etc.)
  int64 amount_cents = 2;       // Amount in cents to avoid floating point
}
```

### Address Type
Represents physical addresses:

```protobuf
message Address {
  string street_line1 = 1;
  string street_line2 = 2;      // Optional
  string city = 3;
  string state_province = 4;
  string postal_code = 5;
  string country_code = 6;      // ISO 3166-1 alpha-2
  double latitude = 7;          // Optional
  double longitude = 8;         // Optional
}
```

### Phone Type
Represents phone numbers:

```protobuf
message Phone {
  string country_code = 1;      // e.g., "+880", "+1"
  string number = 2;
  PhoneType type = 3;           // MOBILE, LANDLINE, WORK
  bool primary = 4;             // Is this the primary phone?
  bool verified = 5;            // Is phone number verified?
}

enum PhoneType {
  PHONE_TYPE_UNSPECIFIED = 0;
  PHONE_TYPE_MOBILE = 1;
  PHONE_TYPE_LANDLINE = 2;
  PHONE_TYPE_WORK = 3;
}
```

### Email Type
Represents email addresses:

```protobuf
message Email {
  string address = 1;
  bool primary = 2;             // Is this the primary email?
  bool verified = 3;            // Is email verified?
}
```

### Document Type
Represents identity documents:

```protobuf
message Document {
  DocumentType type = 1;        // NID, PASSPORT, DRIVING_LICENSE
  string number = 2;
  string issuing_country = 3;
  google.protobuf.Timestamp issue_date = 4;
  google.protobuf.Timestamp expiry_date = 5;
}

enum DocumentType {
  DOCUMENT_TYPE_UNSPECIFIED = 0;
  DOCUMENT_TYPE_NID = 1;
  DOCUMENT_TYPE_PASSPORT = 2;
  DOCUMENT_TYPE_DRIVING_LICENSE = 3;
  DOCUMENT_TYPE_BIRTH_CERTIFICATE = 4;
}
```

### AuditInfo Type
Standard audit information embedded in entities:

```protobuf
message AuditInfo {
  string created_by = 1;
  google.protobuf.Timestamp created_at = 2;
  string updated_by = 3;
  google.protobuf.Timestamp updated_at = 4;
  string deleted_by = 5;
  google.protobuf.Timestamp deleted_at = 6;
}
```

### InsuranceType Enum
Types of insurance products:

```protobuf
enum InsuranceType {
  INSURANCE_TYPE_UNSPECIFIED = 0;
  INSURANCE_TYPE_HEALTH = 1;
  INSURANCE_TYPE_LIFE = 2;
  INSURANCE_TYPE_AUTO = 3;
  INSURANCE_TYPE_HOME = 4;
  INSURANCE_TYPE_TRAVEL = 5;
  INSURANCE_TYPE_DENTAL = 6;
  INSURANCE_TYPE_VISION = 7;
  INSURANCE_TYPE_DISABILITY = 8;
  INSURANCE_TYPE_CRITICAL_ILLNESS = 9;
}
```

---

## Summary of Common Patterns

### Standard Entity Fields
Most entities include:
- `id` (UUID): Primary key
- `AuditInfo`: created_by, created_at, updated_by, updated_at
- `google.protobuf.Timestamp`: For time-based fields
- Status enums: ACTIVE, INACTIVE, PENDING, etc.

### Standard Event Fields
All events include:
- `event_id`: Unique identifier
- `correlation_id`: For distributed tracing
- `timestamp`: Event occurrence time
- Entity-specific payload

### Database Mapping
All entity fields are annotated with:
- Column metadata (name, type, constraints)
- Index specifications
- Foreign key relationships
- Soft delete support

### Error Handling
All service RPCs return responses with optional `Error` field

### Security
Sensitive fields are annotated with PII, encryption, and logging controls

---

End of Common Types Reference
