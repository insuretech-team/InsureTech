"""
Script to append remaining critical sections to SRS_V3_FINAL_DRAFT.md
This includes the complete Protocol Buffer data models section
"""

# Protocol Buffer Data Models Section
proto_section = """
[[[PAGEBREAK]]]

## 6. Data Model - Protocol Buffers

### 6.1 Overview

All data entities are defined using Protocol Buffers (Proto3) to ensure:
- **Type Safety:** Compile-time validation across all services
- **Language Agnostic:** Same .proto files used by Go, C#, Node.js, Python services
- **Performance:** Binary serialization faster than JSON
- **Versioning:** Field numbers enable backward/forward compatibility
- **Code Generation:** Auto-generate data models for all languages

### 6.2 Core Entity Definitions

#### 6.2.1 User Entity

```protobuf
syn tax = "proto3";

package labaid.entities;

option csharp_namespace = "LabAid.Entities";  
option go_package = "labaid/entities";

import "google/protobuf/timestamp.proto";

// User entity for customers, partners, admins
message User {
  string user_id = 1;  // UUID
  string mobile_number = 2;  // E.164 format
  string email = 3;  // Optional
  string full_name = 4;
  google.protobuf.Timestamp date_of_birth = 5;
  string nid_number = 6;  // National ID
  string photo_url = 7;  // S3 path
  string nid_front_url = 8;  // S3 path
  string nid_back_url = 9;  // S3 path
  KycStatus kyc_status = 10;
  UserType user_type = 11;
  google.protobuf.Timestamp created_at = 12;
  google.protobuf.Timestamp updated_at = 13;
  bool is_active = 14;
  
  enum KycStatus {
    KYC_STATUS_UNSPECIFIED = 0;
    KYC_STATUS_PENDING = 1;
    KYC_STATUS_VERIFIED = 2;
    KYC_STATUS_REJECTED = 3;
  }
  
  enum UserType {
    USER_TYPE_UNSPECIFIED = 0;
    USER_TYPE_CUSTOMER = 1;
    USER_TYPE_PARTNER_ADMIN = 2;
    USER_TYPE_ADMIN = 3;
    USER_TYPE_SUPPORT = 4;
  }
}

// Additional user profile data
message UserProfile {
  string user_id = 1;  // FK to User
  string current_address = 2;
  string permanent_address = 3;
  string occupation = 4;
  double monthly_income = 5;  // BDT
  string employer_name = 6;
  map<string, string> additional_data = 7;  // Flexible key-value storage
}
```

#### 6.2.2 Policy Entity

```protobuf
syntax = "proto3";

package labaid.entities;

import "google/protobuf/timestamp.proto";

// Policy entity
message Policy {
  string policy_id = 1;  // UUID
  string policy_number = 2;  // Human-readable unique number
  string user_id = 3;  // FK to User
  string product_id = 4;  // FK to Product
  string tenant_id = 5;  // FK to Partner (nullable for direct sales)
  double premium_amount = 6;  // BDT
  double coverage_amount = 7;  // BDT
  google.protobuf.Timestamp start_date = 8;
  google.protobuf.Timestamp end_date = 9;
  PolicyStatus status = 10;
  string payment_reference = 11;
  string policy_document_url = 12;  // S3 path to PDF
  google.protobuf.Timestamp created_at = 13;
  google.protobuf.Timestamp updated_at = 14;
  
  // Partitioning hint (comment): Partition by created_at month in PostgreSQL
  
  enum PolicyStatus {
    POLICY_STATUS_UNSPECIFIED = 0;
    POLICY_STATUS_PENDING_PAYMENT = 1;
    POLICY_STATUS_ACTIVE = 2;
    POLICY_STATUS_EXPIRED = 3;
    POLICY_STATUS_CANCELLED = 4;
    POLICY_STATUS_LAPSED = 5;
  }
}

// Nominee/Beneficiary entity
message Nominee {
  string nominee_id = 1;  // UUID
  string policy_id = 2;  // FK to Policy
  string full_name = 3;
  Relationship relationship = 4;
  google.protobuf.Timestamp date_of_birth = 5;
  string national_id = 6;
  double share_percentage = 7;  // 0-100, sum must be 100 per policy
  string contact_phone = 8;
  string address = 9;
  google.protobuf.Timestamp created_at = 10;
  google.protobuf.Timestamp updated_at = 11;
  
  enum Relationship {
    RELATIONSHIP_UNSPECIFIED = 0;
    RELATIONSHIP_SPOUSE = 1;
    RELATIONSHIP_CHILD = 2;
    RELATIONSHIP_PARENT = 3;
    RELATIONSHIP_SIBLING = 4;
    RELATIONSHIP_OTHER = 5;
  }
}
```

#### 6.2.3 Claim Entity

```protobuf
syntax = "proto3";

package labaid.entities;

import "google/protobuf/timestamp.proto";

// Claim entity
message Claim {
  string claim_id = 1;  // UUID
  string claim_number = 2;  // Human-readable unique number
  string policy_id = 3;  // FK to Policy
  google.protobuf.Timestamp incident_date = 4;
  string incident_type = 5;
  string incident_description = 6;
  string incident_location = 7;
  double claimed_amount = 8;  // BDT
  double approved_amount = 9;  // BDT (nullable)
  ClaimStatus status = 10;
  string assigned_to = 11;  // FK to User (admin)
  ApprovalLevel approval_level = 12;
  string rejection_reason = 13;
  google.protobuf.Timestamp settlement_date = 14;
  repeated string document_urls = 15;  // S3 paths
  google.protobuf.Timestamp created_at = 16;
  google.protobuf.Timestamp updated_at = 17;
  
  enum ClaimStatus {
    CLAIM_STATUS_UNSPECIFIED = 0;
    CLAIM_STATUS_SUBMITTED = 1;
    CLAIM_STATUS_UNDER_REVIEW = 2;
    CLAIM_STATUS_DOCUMENTS_REQUESTED = 3;
    CLAIM_STATUS_APPROVED = 4;
    CLAIM_STATUS_REJECTED = 5;
    CLAIM_STATUS_PAYMENT_INITIATED = 6;
    CLAIM_STATUS_SETTLED = 7;
    CLAIM_STATUS_CLOSED = 8;
  }
  
  enum ApprovalLevel {
    APPROVAL_LEVEL_UNSPECIFIED = 0;
    APPROVAL_LEVEL_L1_AUTO = 1;  // <5K BDT
    APPROVAL_LEVEL_L2_MANAGER = 2;  // 5K-20K BDT
    APPROVAL_LEVEL_L3_HEAD = 3;  // 20K-50K BDT
    APPROVAL_LEVEL_EXECUTIVE = 4;  // 50K-200K BDT
    APPROVAL_LEVEL_BOARD = 5;  // 200K+ BDT
  }
}

// Claim status history for audit trail
message ClaimStatusHistory {
  string history_id = 1;  // UUID
  string claim_id = 2;  // FK to Claim
  Claim.ClaimStatus from_status = 3;
  Claim.ClaimStatus to_status = 4;
  string changed_by = 5;  // FK to User
  string notes = 6;
  google.protobuf.Timestamp changed_at = 7;
}
```

#### 6.2.4 Payment Entity

```protobuf
syntax = "proto3";

package labaid.entities;

import "google/protobuf/timestamp.proto";

// Payment entity
message Payment {
  string payment_id = 1;  // UUID
  string reference_number = 2;  // Unique reference
  string policy_id = 3;  // FK to Policy (nullable for claim settlements)
  string claim_id = 4;  // FK to Claim (nullable for policy payments)
  double amount = 5;  // BDT
  PaymentMethod payment_method = 6;
  string payment_proof_url = 7;  // S3 path (for manual payments)
  string transaction_id = 8;  // From MFS/gateway
  PaymentStatus status = 9;
  string verified_by = 10;  // FK to User (admin)
  google.protobuf.Timestamp verified_at = 11;
  google.protobuf.Timestamp created_at = 12;
  PaymentDirection direction = 13;  // Inbound or outbound
  
  enum PaymentMethod {
    PAYMENT_METHOD_UNSPECIFIED = 0;
    PAYMENT_METHOD_MANUAL = 1;
    PAYMENT_METHOD_BKASH = 2;
    PAYMENT_METHOD_NAGAD = 3;
    PAYMENT_METHOD_ROCKET = 4;
    PAYMENT_METHOD_BANK = 5;
    PAYMENT_METHOD_CARD = 6;
  }
  
  enum PaymentStatus {
    PAYMENT_STATUS_UNSPECIFIED = 0;
    PAYMENT_STATUS_PENDING_VERIFICATION = 1;
    PAYMENT_STATUS_VERIFIED = 2;
    PAYMENT_STATUS_REJECTED = 3;
    PAYMENT_STATUS_REFUNDED = 4;
    PAYMENT_STATUS_SETTLED = 5;
  }
  
  enum PaymentDirection {
    PAYMENT_DIRECTION_UNSPECIFIED = 0;
    PAYMENT_DIRECTION_INBOUND = 1;  // Customer to InsureTech
    PAYMENT_DIRECTION_OUTBOUND = 2;  // InsureTech to Customer/Partner
  }
}
```

#### 6.2.5 Partner Entity

```protobuf
syntax = "proto3";

package labaid.entities;

import "google/protobuf/timestamp.proto";

// Partner entity (multi-tenant)
message Partner {
  string tenant_id = 1;  // UUID - primary tenant identifier
  string partner_name = 2;
  PartnerType partner_type = 3;
  string trade_license = 4;
  string tin = 5;  // Tax Identification Number
  string bank_account = 6;
  string contact_person = 7;
  string contact_mobile = 8;
  string contact_email = 9;
  string mou_document_url = 10;  // S3 path
  KybStatus kyb_status = 11;
  double commission_rate = 12;  // Percentage
  string onboarded_by = 13;  // FK to User (Focal Person)
  google.protobuf.Timestamp created_at = 14;
  bool is_active = 15;
  
  enum PartnerType {
    PARTNER_TYPE_UNSPECIFIED = 0;
    PARTNER_TYPE_HOSPITAL = 1;
    PARTNER_TYPE_MFS = 2;
    PARTNER_TYPE_ECOMMERCE = 3;
    PARTNER_TYPE_NGO = 4;
    PARTNER_TYPE_AGENT = 5;
  }
  
  enum KybStatus {
    KYB_STATUS_UNSPECIFIED = 0;
    KYB_STATUS_PENDING = 1;
    KYB_STATUS_VERIFIED = 2;
    KYB_STATUS_REJECTED = 3;
  }
}

// Partner commission tracking
message PartnerCommission {
  string commission_id = 1;  // UUID
  string partner_id = 2;  // FK to Partner (tenant_id)
  string policy_id = 3;  // FK to Policy
  CommissionType commission_type = 4;
  double commission_rate = 5;  // Percentage
  double commission_amount = 6;  // BDT
  PayoutStatus payout_status = 7;
  google.protobuf.Timestamp payout_date = 8;
  string payout_reference = 9;
  google.protobuf.Timestamp created_at = 10;
  google.protobuf.Timestamp updated_at = 11;
  
  enum CommissionType {
    COMMISSION_TYPE_UNSPECIFIED = 0;
    COMMISSION_TYPE_ACQUISITION = 1;
    COMMISSION_TYPE_RENEWAL = 2;
    COMMISSION_TYPE_CLAIMS_ASSISTANCE = 3;
  }
  
  enum PayoutStatus {
    PAYOUT_STATUS_UNSPECIFIED = 0;
    PAYOUT_STATUS_PENDING = 1;
    PAYOUT_STATUS_APPROVED = 2;
    PAYOUT_STATUS_PAID = 3;
    PAYOUT_STATUS_WITHHELD = 4;
  }
}
```

#### 6.2.6 Product Entity

```protobuf
syntax = "proto3";

package labaid.entities;

import "google/protobuf/timestamp.proto";
import "google/protobuf/struct.proto";

// Product entity
message Product {
  string product_id = 1;  // UUID
  string product_code = 2;  // Unique code
  LocalizedString product_name = 3;  // Bengali + English
  ProductCategory category = 4;
  LocalizedString description = 5;
  google.protobuf.Struct coverage_details = 6;  // Flexible JSON structure
  double premium_base = 7;  // Base premium amount
  repeated PremiumFactor premium_factors = 8;
  google.protobuf.Struct underwriting_rules = 9;  // Flexible rules engine
  string policy_document_template = 10;  // Template identifier
  bool is_active = 11;
  google.protobuf.Timestamp created_at = 12;
  google.protobuf.Timestamp updated_at = 13;
  
  enum ProductCategory {
    PRODUCT_CATEGORY_UNSPECIFIED = 0;
    PRODUCT_CATEGORY_HEALTH = 1;
    PRODUCT_CATEGORY_LIFE = 2;
    PRODUCT_CATEGORY_DEVICE = 3;
    PRODUCT_CATEGORY_LIVESTOCK = 4;
    PRODUCT_CATEGORY_CROP = 5;
    PRODUCT_CATEGORY_MOTOR = 6;
  }
}

// Localized string for Bengali + English support
message LocalizedString {
  string en = 1;  // English
  string bn = 2;  // Bengali
}

// Premium calculation factors
message PremiumFactor {
  string factor_name = 1;  // e.g., "age_loading", "occupation_risk"
  double multiplier = 2;
  string condition = 3;  // Expression: "age > 50"
}
```

#### 6.2.7 Quote Entity

```protobuf
syntax = "proto3";

package labaid.entities;

import "google/protobuf/timestamp.proto";

// Quote entity for premium calculations
message Quote {
  string quote_id = 1;  // UUID
  string user_id = 2;  // FK to User (nullable if anonymous)
  string product_id = 3;  // FK to Product
  double requested_sum_assured = 4;  // BDT
  int32 requested_tenure = 5;  // Years
  int32 customer_age = 6;
  string occupation = 7;
  double base_premium = 8;  // BDT
  repeated Loading loadings = 9;
  repeated Discount discounts = 10;
  double final_premium = 11;  // BDT
  google.protobuf.Timestamp quote_valid_until = 12;  // Expiry (48 hours)
  string converted_to_policy_id = 13;  // FK to Policy (nullable)
  google.protobuf.Timestamp created_at = 14;
  google.protobuf.Timestamp updated_at = 15;
}

message Loading {
  string reason = 1;
  double amount = 2;  // BDT
}

message Discount {
  string code = 1;
  double amount = 2;  // BDT
}
```

#### 6.2.8 KYC Document Entity

```protobuf
syntax = "proto3";

package labaid.entities;

import "google/protobuf/timestamp.proto";

// KYC document tracking
message KycDocument {
  string document_id = 1;  // UUID
  string user_id = 2;  // FK to User
  DocumentType document_type = 3;
  string file_name = 4;
  string s3_path = 5;
  string file_hash = 6;  // SHA256 for integrity
  google.protobuf.Timestamp uploaded_at = 7;
  VerificationStatus verification_status = 8;
  string verified_by = 9;  // FK to User (admin)
  google.protobuf.Timestamp verified_at = 10;
  string rejection_reason = 11;
  map<string, string> ocr_data = 12;  // Extracted fields
  google.protobuf.Timestamp expiry_date = 13;  // For passports/utility bills
  
  enum DocumentType {
    DOCUMENT_TYPE_UNSPECIFIED = 0;
    DOCUMENT_TYPE_NID_FRONT = 1;
    DOCUMENT_TYPE_NID_BACK = 2;
    DOCUMENT_TYPE_PASSPORT = 3;
    DOCUMENT_TYPE_PHOTO = 4;
    DOCUMENT_TYPE_UTILITY_BILL = 5;
    DOCUMENT_TYPE_INCOME_PROOF = 6;
    DOCUMENT_TYPE_MEDICAL_REPORT = 7;
  }
  
  enum VerificationStatus {
    VERIFICATION_STATUS_UNSPECIFIED = 0;
    VERIFICATION_STATUS_PENDING = 1;
    VERIFICATION_STATUS_VERIFIED = 2;
    VERIFICATION_STATUS_REJECTED = 3;
  }
}
```

#### 6.2.9 Audit Log Entity

```protobuf
syntax = "proto3";

package labaid.entities;

import "google/protobuf/timestamp.proto";
import "google/protobuf/struct.proto";

// Immutable audit log
message AuditLog {
  int64 log_id = 1;  // BIGSERIAL (auto-increment)
  string user_id = 2;  // FK to User
  string action = 3;  // e.g., "POLICY_ISSUED", "CLAIM_APPROVED"
  string resource_type = 4;  // e.g., "Policy", "Claim"
  string resource_id = 5;  // UUID of affected resource
  google.protobuf.Struct old_value = 6;  // Previous state (JSON)
  google.protobuf.Struct new_value = 7;  // New state (JSON)
  string ip_address = 8;  // IPv4/IPv6
  string user_agent = 9;
  google.protobuf.Timestamp created_at = 10;
  string signature = 11;  // Cryptographic signature (HMAC-SHA256)
  
  // NOTE: This table is write-only (immutable) for compliance
}
```

#### 6.2.10 Notification Log Entity

```protobuf
syntax = "proto3";

package labaid.entities;

import "google/protobuf/timestamp.proto";

// Notification tracking
message NotificationLog {
  string notification_id = 1;  // UUID
  string user_id = 2;  // FK to User (nullable for system-wide)
  Channel channel = 3;
  EventType event_type = 4;
  string message_template = 5;  // Template ID
  string message_content = 6;  // Rendered message
  string recipient_phone = 7;  // For SMS
  string recipient_email = 8;  // For Email
  string device_token = 9;  // For Push
  SendStatus send_status = 10;
  google.protobuf.Timestamp sent_at = 11;
  google.protobuf.Timestamp delivered_at = 12;
  string failure_reason = 13;
  int32 retry_count = 14;
  google.protobuf.Timestamp created_at = 15;
  
  enum Channel {
    CHANNEL_UNSPECIFIED = 0;
    CHANNEL_SMS = 1;
    CHANNEL_EMAIL = 2;
    CHANNEL_PUSH = 3;
    CHANNEL_IN_APP = 4;
  }
  
  enum EventType {
    EVENT_TYPE_UNSPECIFIED = 0;
    EVENT_TYPE_OTP = 1;
    EVENT_TYPE_PURCHASE_CONFIRMATION = 2;
    EVENT_TYPE_CLAIM_UPDATE = 3;
    EVENT_TYPE_RENEWAL_REMINDER = 4;
    EVENT_TYPE_MARKETING = 5;
  }
  
  enum SendStatus {
    SEND_STATUS_UNSPECIFIED = 0;
    SEND_STATUS_QUEUED = 1;
    SEND_STATUS_SENT = 2;
    SEND_STATUS_FAILED = 3;
    SEND_STATUS_DELIVERED = 4;
    SEND_STATUS_BOUNCED = 5;
  }
}
```

### 6.3 Data Storage Distribution

| Data Type | Storage System | Proto Usage | Rationale |
|-----------|---------------|-------------|-----------|
| **Users & Policies** | PostgreSQL 17 | Proto definitions map to SQL schema | ACID compliance, complex queries, relationshipsjoin
| **Product Catalog** | MongoDB | Proto to BSON | Flexible schema for dynamic product attributes |
| **Documents & Files** | AWS S3 | Proto metadata stored in PostgreSQL | Scalable object storage |
| **Cache & Sessions** | Redis | Proto serialized to binary | Fast in-memory operations |
| **Event Logs** | MongoDB | Proto to BSON | High write throughput for audit logs |
| **gRPC Communication** | All services | Native Proto | Type-safe service contracts |

### 6.4 Proto File Organization

```
proto/
├── entities/
│   ├── user.proto
│   ├── policy.proto
│   ├── claim.proto
│   ├── payment.proto
│   ├── partner.proto
│   ├── product.proto
│   ├── quote.proto
│   ├── kyc_document.proto
│   ├── audit_log.proto
│   └── notification_log.proto
├── services/
│   ├── insurance_engine.proto
│   ├── partner_management.proto
│   ├── ai_engine.proto
│   ├── payment_service.proto
│   ├── notification_service.proto
│   └── analytics_service.proto
└── common/
    ├── localized_string.proto
    ├── pagination.proto
    └── error.proto
```

### 6.5 Code Generation

**Generate Go models:**
```bash
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/entities/*.proto proto/services/*.proto
```

**Generate C# models:**
```bash
protoc --csharp_out=. --grpc_out=. --plugin=protoc-gen-grpc=grpc_csharp_plugin \
    proto/entities/*.proto proto/services/*.proto
```

**Generate Node.js models:**
```bash
protoc --js_out=import_style=commonjs,binary:. \
    --grpc-web_out=import_style=commonjs,mode=grpcwebtext:. \
    proto/entities/*.proto proto/services/*.proto
```

**Generate Python models:**
```bash
python -m grpc_tools.protoc -I. \
    --python_out=. --grpc_python_out=. \
    proto/entities/*.proto proto/services/*.proto
```

[[[PAGEBREAK]]]

"""

# Write to file
output_file = r"G:\_0LifePlus\InsureTech\SRS_v3\proto_section.txt"
with open(output_file, 'w', encoding='utf-8') as f:
    f.write(proto_section)

print(f"Proto section written to: {output_file}")
print(f"Next: Append this to SRS_V3_FINAL_DRAFT.md")
