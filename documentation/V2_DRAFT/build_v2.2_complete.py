#!/usr/bin/env python3
"""
Complete V2.2 Document Builder
Adds remaining sections to V2.2_COMPLETE.md
"""

# Read current content
with open(r"C:\_DEV\GO\InsureTech\V2.2_COMPLETE.md", "r", encoding="utf-8") as f:
    content = f.read()

# Add remaining sections
remaining_sections = """
### 4.15 Business Rules & Workflows (MD FEEDBACK)

| FG | ID | Requirement Description | Priority |
|----|----|-----------------------|----------|
| **FG-015** | | **Business Rules & Workflows** | **M** |
| | FR-129 | Premium Calculation Fallbacks: If insurer API fails, use cached rates (max 24hrs old); if unavailable, queue quote and notify customer within 2 hours | M |
| | FR-130 | Premium Calculation Edge Cases: Handle age-based loading, occupation risk factors, pre-existing conditions with clear messaging | M |
| | FR-131 | Duplicate Policy Detection: Block duplicate policy purchase for same product + same insured person within 30 days; allow cross-product purchases | M |
| | FR-132 | Policy Merge Workflow: Focal Person can merge duplicate accounts after verifying NID, transfer policies, consolidate claims history | M |
| | FR-133 | Claim Status State Machine: Define explicit states (Submitted → Under Review → Documents Requested → Approved/Rejected → Payment Initiated → Settled/Closed) | M |
| | FR-134 | Claim Status Transition Rules: Auto-move to "Documents Requested" if incomplete; require Business Admin+Focal Person approval for >BDT 50K | M |
| | FR-135 | Renewal Timeline: Send first reminder 30 days before expiry, second at 15 days, third at 7 days; allow renewal 45 days before to 30 days after expiry | M |
| | FR-136 | Grace Period Logic: 30-day grace period post-expiry with coverage continued; auto-lapse if unpaid after grace period | M |
| | FR-137 | Lapsed Policy Reinstatement: Allow reinstatement within 90 days of lapse with medical underwriting; require Focal Person approval | D |

---

### 4.16 Integration Details (MD FEEDBACK)

| FG | ID | Requirement Description | Priority |
|----|----|-----------------------|----------|
| **FG-016** | | **Integration & API Contracts** | **M** |
| | FR-138 | API Contract Specification: All Category 3 APIs must provide OpenAPI 3.0 spec with request/response schemas, error codes (400-BadRequest, 401-Unauthorized, 404-NotFound, 409-Conflict, 422-ValidationError, 500-ServerError, 503-ServiceUnavailable), example payloads | M |
| | FR-139 | Insurer API Payloads: Premium Calculation API: Request {productId, age, sumAssured, tenure, occupation} → Response {basePremium, loadings[], discounts[], finalPremium, breakdown[]}; Policy Issuance API: Request {quoteId, customerDetails, paymentProof} → Response {policyNumber, certificateUrl, effectiveDate} | M |
| | FR-140 | Payment Gateway Payloads: Initiate Payment: Request {orderId, amount, currency:BDT, customerPhone, returnUrl, webhookUrl} → Response {paymentUrl, transactionId}; Webhook Callback: {transactionId, status, orderId, amount, signature(HMAC-SHA256)} | M |
| | FR-141 | Retry Logic: Failed API calls retry with exponential backoff: 1s, 2s, 4s, 8s, 16s (max 5 retries); Use circuit breaker pattern (open after 5 consecutive failures, half-open after 30s, close after 3 successes) | M |
| | FR-142 | Idempotency: All payment and policy issuance APIs must accept Idempotency-Key header (UUID); Store idempotency keys for 24 hours; Return cached response for duplicate requests within 24hrs | M |
| | FR-143 | Callback Security: Payment gateway webhooks must include HMAC-SHA256 signature in header; System validates signature using shared secret; Reject unsigned/invalid callbacks; Log all callback attempts | M |
| | FR-144 | EHR Integration Approach - Option A (Preferred): Use LabAid FHIR API with Patient resource matching by NID/phone; Query Encounter resources for IPD admissions; Verify policy coverage in real-time; Pre-authorization workflow with Claim resource | D |
| | FR-145 | EHR Integration Approach - Option B (Fallback): Use LabAid custom REST API with endpoints: GET /patients/{nid}/admissions, POST /preauth/verify, GET /bills/{admission_id}; Secure with mutual TLS + API key | M |
| | FR-146 | EHR Integration Timeout: Set connection timeout 5s, read timeout 15s; If timeout, queue for manual verification; Notify hospital staff via SMS | M |

---

### 4.17 Fraud Detection & Risk Controls (MD FEEDBACK)

| FG | ID | Requirement Description | Priority |
|----|----|-----------------------|----------|
| **FG-017** | | **Fraud Detection & Operations Controls** | **M** |
| | FR-147 | Basic Claim Abuse Checks: Flag claims with (1) Submission within 48hrs of policy purchase, (2) Same claim type >2 times in 12 months, (3) Claim amount = policy limit exactly, (4) Medical provider not in approved network | M |
| | FR-148 | Monitoring Ownership - RACI: Responsible: DevOps Engineer (24x7 monitoring), Accountable: System Admin (incident escalation), Consulted: Business Admin (business impact), Informed: All stakeholders (post-incident report) | M |
| | FR-149 | Escalation RACI for P1 Incidents: Responsible: On-call engineer (immediate response), Accountable: System Admin (decision maker), Consulted: Database Admin + Repository Admin (technical expertise), Informed: Business Admin + Focal Person (business continuity) | M |
| | FR-150 | Fraud Detection Dashboard: Business Admin and Focal Person access real-time dashboard showing flagged transactions, AML alerts, duplicate policy attempts, frequent claim patterns with drill-down capability | M |

---

## 5. External Interface Requirements

### 5.1 User Interfaces

**Mobile Application Requirements:**

| Requirement | Specification | Rationale |
|-------------|---------------|-----------|
| **Screen Design** | Consistent with provided prototype designs | Bangladesh-specific UX considerations |
| **Orientation Support** | Portrait and landscape | Enhanced user experience |
| **Screen Size Optimization** | 4.7" to 6.7" | Coverage of common device range |
| **Theme Support** | Dark mode and light mode | User preference and battery saving |
| **Touch Targets** | Minimum 44px | Accessibility guidelines compliance |
| **Download Size** | Initial: 10MB, Maximum: 100MB | Bangladesh network constraints |
| **Offline Mode** | Critical data persistent view | Network reliability issues |
| **Voice Assistance** | For Type 3 users | Low digital literacy support |

**Web Portal Requirements:**

| Requirement | Specification | Rationale |
|-------------|---------------|-----------|
| **Design Approach** | Desktop-first responsive design | Admin/partner primary use case |
| **Resolution Support** | 1024x768 to 4K displays | Wide device coverage |
| **Multi-tab Support** | Multiple browser tabs | Productivity enhancement |
| **Search & Filter** | Comprehensive capabilities | Data management efficiency |

### 5.2 API Architecture & Communication Protocols

**API Category Structure:**

| API Category | Protocol | Use Case | Security Layer | Performance Target |
|--------------|----------|----------|----------------|-------------------|
| **Category 1** | Protocol Buffer + gRPC | Gateway ↔ Microservices | System Admin Middle Layer | < 100ms |
| **Category 2** | GraphQL + JWT | Gateway ↔ Customer Device | JWT + OAuth v2 | < 2 seconds |
| **Category 3** | RESTful + JSON (OpenAPI) | 3rd Party Integration | Server-side Auth | < 200ms |
| **Public API** | RESTful + JSON (OpenAPI) | Product Search/List | Public Access | < 1 second |

### 5.3 Third-Party System Integrations

**Payment Gateway Integrations:**

| Provider | Integration Method | Scope | Compliance Requirements |
|----------|-------------------|-------|------------------------|
| **bKash** | REST API + Webhook | Payment processing, refunds, transaction status | bKash merchant agreement |
| **Nagad** | REST API + Webhook | Payment processing, refunds, transaction status | Nagad partner certification |
| **Rocket** | REST API + Webhook | Payment processing, refunds, transaction status | DBBL partnership |
| **Card Processors** | PCI-DSS Gateway | Credit/debit card processing | PCI-DSS Level 1 compliance |

**Insurance Company API Integration:**
- Premium calculation and underwriting APIs
- Policy issuance and certificate generation
- Claims processing and settlement APIs
- Regulatory reporting data exchange
- Real-time policy status updates

**LabAid EHR System Integration:**
- Cashless claim verification for IPD services
- Pre-authorization workflows
- Medical record access (with consent)
- Integration via secured HL7/FHIR endpoints or custom APIs
- Timeout handling: 5s connection, 15s read

---

## 6. Non-Functional Requirements

| ID | Requirement Description | Priority |
|----|-------------------------|----------|
| NFR-001 | The system shall ensure 99.5% uptime to meet critical business needs | M |
| NFR-002 | The system shall encrypt (TLS 1.3 and AES-256) all sensitive data in transit and at rest | M |
| NFR-003 | The system shall handle up to 1,000 concurrent users without performance degradation | M |
| NFR-004 | The system shall handle up to 5,000 concurrent users without performance degradation | D |
| NFR-005 | The system shall handle up to 10,000 concurrent users without performance degradation | F |
| NFR-006 | The system shall have a median response time of less than 1 second and not more than 2s for 95% of user actions | M |
| NFR-007 | Disaster Recovery: RTO < 15 minutes for critical services, RPO < 1 hour | M |
| NFR-008 | The system shall have Category 1 API response time < 100ms | D |
| NFR-009 | The system shall have Category 3 API response time < 200ms | D |
| NFR-010 | The system shall have Category 2 API response time < 2 sec | D |
| NFR-011 | The system shall have public API response time < 1 sec | D |
| NFR-012 | The system shall ensure minimum app startup time < 5 sec | D |
| NFR-013 | The system shall support 1 million active users and 200K active policies in 24 months with auto vertical scaling | F |
| NFR-014 | The system shall have automatic switch option for user to different mobile device with data migration | F |
| NFR-015 | The system shall have multi-language support (English + Bengali) | D |

---

## 7. Data Model & Storage Requirements

### 7.1 Database Architecture Distribution

| Storage System | Use Case | Performance Target | Rationale |
|----------------|----------|-------------------|-----------|
| **PostgreSQL V17** | Relational data, user/policy/claims | < 100ms for 95% queries | ACID compliance, JSONB, Bengali full-text search |
| **TigerBeetle** | Financial transactions, insurance records | < 10ms transaction time | Purpose-built for financial accuracy |
| **Redis** | Cache, sessions | < 1ms access time | Product catalog (5min TTL), User sessions (15min) |
| **DynamoDB/MongoDB** | Product catalog, metadata | < 50ms read time | NoSQL flexibility for unstructured data |
| **AWS S3** | Documents, uploaded objects | < 5s upload/download | Scalable object storage with encryption |
| **Pgvector/PineCone** | Tokenized data, AI/ML | < 100ms similarity search | Vector database for ML operations |
| **SQLite** | Mobile app local data | < 10ms local query | Encrypted offline capability |

### 7.2 Enhanced Data Entities (MD FEEDBACK)

**Nominee/Beneficiary Entity:**
```sql
NomineeID (Primary Key, UUID)
PolicyID (Foreign Key)
FullName (VARCHAR 200, NOT NULL)
Relationship (ENUM: Spouse, Child, Parent, Sibling, Other)
DateOfBirth (DATE)
NationalID (VARCHAR 20, UNIQUE)
SharePercentage (DECIMAL 5,2, CHECK >0 AND <=100)
ContactPhone (VARCHAR 15)
Address (TEXT)
CreatedAt, UpdatedAt
CONSTRAINT: SUM(SharePercentage) per PolicyID = 100
```

**KYCDocument Entity:**
```sql
DocumentID (Primary Key, UUID)
UserID (Foreign Key)
DocumentType (ENUM: NID_Front, NID_Back, Passport, Photo, Utility_Bill, Income_Proof, Medical_Report)
FileName (VARCHAR 255)
S3Path (TEXT)
FileHash (SHA256, for integrity)
UploadedAt (TIMESTAMP)
VerificationStatus (ENUM: Pending, Verified, Rejected)
VerifiedBy (Foreign Key → Stakeholder)
VerifiedAt (TIMESTAMP)
RejectionReason (TEXT, NULL)
OCRData (JSONB, extracted fields)
ExpiryDate (DATE, for passports/utility bills)
```

**Quote Entity:**
```sql
QuoteID (Primary Key, UUID)
UserID (Foreign Key, NULL if anonymous)
ProductID (Foreign Key)
RequestedSumAssured (DECIMAL)
RequestedTenure (INT, years)
CustomerAge (INT)
Occupation (VARCHAR 100)
BasePremium (DECIMAL)
Loadings (JSONB, {reason: amount}[])
Discounts (JSONB, {code: amount}[])
FinalPremium (DECIMAL)
QuoteValidUntil (TIMESTAMP, +48hrs)
ConvertedToPolicyID (Foreign Key, NULL)
CreatedAt, UpdatedAt
```

**PartnerCommission Entity:**
```sql
CommissionID (Primary Key, UUID)
PartnerID (Foreign Key)
PolicyID (Foreign Key)
CommissionType (ENUM: Acquisition, Renewal, Claims_Assistance)
CommissionRate (DECIMAL 5,2, percentage)
CommissionAmount (DECIMAL)
PayoutStatus (ENUM: Pending, Approved, Paid, Withheld)
PayoutDate (DATE, NULL)
PayoutReference (VARCHAR 100, payment txn ID)
CreatedAt, UpdatedAt
```

**NotificationLog Entity:**
```sql
NotificationID (Primary Key, UUID)
UserID (Foreign Key, NULL for system-wide)
Channel (ENUM: SMS, Email, Push, In_App)
EventType (ENUM: OTP, Purchase_Confirmation, Claim_Update, Renewal_Reminder, Marketing)
MessageTemplate (VARCHAR 50, template ID)
MessageContent (TEXT)
RecipientPhone (VARCHAR 15, for SMS)
RecipientEmail (VARCHAR 255, for Email)
DeviceToken (TEXT, for Push)
SendStatus (ENUM: Queued, Sent, Failed, Delivered, Bounced)
SentAt (TIMESTAMP)
DeliveredAt (TIMESTAMP, NULL)
FailureReason (TEXT, NULL)
RetryCount (INT, DEFAULT 0)
CreatedAt
```

### 7.3 Data Storage Distribution Strategy

| Data Type | Storage System | Rationale | Performance Target |
|-----------|----------------|-----------|-------------------|
| **User & Policy Data** | PostgreSQL V17 | ACID compliance, complex queries, Bengali support | < 100ms query time |
| **Financial Transactions** | TigerBeetle | Purpose-built for financial accuracy, double-entry bookkeeping | < 10ms transaction time |
| **Product Catalog** | DynamoDB/MongoDB | Flexible schema, fast reads | < 50ms read time |
| **Documents & Files** | AWS S3 | Scalable object storage, encryption at rest | < 5s upload/download |
| **Cache & Sessions** | Redis | Fast in-memory operations | < 1ms access time |
| **Vector/AI Data** | Pgvector/PineCone | ML model data, similarity search | < 100ms similarity search |
| **Mobile App Data** | SQLite | Offline capability, encrypted local storage | < 10ms local query |

---

## 8. Security & Compliance Requirements

### 8.1 Security Infrastructure & Key Management

| ID | Requirement Description | Priority |
|----|-------------------------|----------|
| SEC-001 | The system shall use separate secret vault - AWS KMS/Azure Key Vault/HashiCorp, 90-day key rotation | M |
| SEC-002 | The system shall use Data Masking: NID (last 3 digits), phone (mask middle), email (mask username) | M |
| SEC-003 | The system shall follow PCI-DSS compliance for card flows - Approach: Hosted payment page (redirect model) - DO NOT store card data, Level: SAQ-A (simplest, for redirecting merchants), Requirements: Annual SAQ, quarterly ASV scans, TLS 1.3, Tokenization: Store only gateway tokens for recurring payments | M |
| SEC-004 | The system shall have AML/CFT detection hooks - Transaction Monitoring: 20+ automated rules for AML detection including Rapid purchases (>3 policies in 7 days), High-value premiums (>BDT 5 lakh), Frequent cancellations, Mismatched nominees, Geographic/payment anomalies | D |
| SEC-005 | The system shall have IDRA reporting capabilities following IDRA data format - Monthly Reports: Premium Collection (Form IC-1), Claims Intimation (Form IC-2), Quarterly Reports: Claims Settlement (IC-3), Financial Performance (IC-4), Annual Reports: FCR (Financial Condition Report), CARAMELS Framework Returns, Event-Based: Significant incidents (48hrs), fraud cases (7 days), Platform: Report generator with IDRA Excel templates, audit trail, 20-year archive | D |
| SEC-006 | The system shall have regular penetration testing - Penetration Testing: Pre-launch + annually (SISA InfoSec or international firm) | D |
| SEC-007 | The system shall have regular security audits from various security auditors and regulatory bodies and maintain compliance | D |
| SEC-008 | DAST: OWASP ZAP/Burp Suite (weekly on staging) | D |
| SEC-009 | SAST: SonarQube/Checkmarx (every commit, block critical vulnerabilities) | D |
| SEC-010 | Virus scanning: ClamAV on uploaded files | M |

### 8.2 Enhanced IDRA Compliance (MD FEEDBACK)

| ID | Requirement Description | Priority |
|----|-------------------------|----------|
| SEC-011 | IDRA Monthly Reports: Generate Form IC-1 (Premium Collection) by 10th of each month with breakdown by product line, geographic region, partner channel in Excel format per IDRA template v2024 | M |
| SEC-012 | IDRA Monthly Reports: Generate Form IC-2 (Claims Intimation) by 10th of each month listing all new claims with policy number, claim amount, claim type, date of intimation | M |
| SEC-013 | IDRA Quarterly Reports: Generate Form IC-3 (Claims Settlement) within 15 days of quarter-end showing settlement ratio, average TAT, pending >30 days breakdown | M |
| SEC-014 | IDRA Quarterly Reports: Generate Form IC-4 (Financial Performance) within 20 days of quarter-end with premium earned, claims paid, commission paid, net profit/loss | M |
| SEC-015 | IDRA Annual FCR: Generate Financial Condition Report (FCR) within 90 days of year-end including full CARAMELS framework assessment with external auditor sign-off | M |
| SEC-016 | IDRA Event-Based Reporting: Report significant incidents (fraud >BDT 1L, data breach, system outage >4hrs) within 48 hours via IDRA portal | M |

### 8.3 Enhanced AML/CFT Compliance (MD FEEDBACK)

| ID | Requirement Description | Priority |
|----|-------------------------|----------|
| SEC-017 | AML/CFT Concrete Triggers: Flag transactions matching: (1) >3 policies in 7 days, (2) Premium >BDT 5L without income proof, (3) Nominee mismatch with no relationship doc, (4) Payment from third-party account, (5) Frequent cancellations >2 in 30 days, (6) Geographic anomaly (policy in Dhaka, payment from remote district), (7) Multiple failed KYC attempts >3, (8) PEP match in screening | M |
| SEC-018 | SAR Workflow: (1) System auto-flags suspicious transaction → (2) Compliance Officer reviews within 24hrs → (3) If confirmed suspicious, escalate to Business Admin+Focal Person → (4) Prepare SAR with evidence → (5) Submit to BFIU within 3 business days → (6) Mark account for enhanced monitoring → (7) Do NOT notify customer (tipping off prohibited) | M |
| SEC-019 | Data Deletion Exceptions: Customer data deletion requests processed within 30 days EXCEPT: (a) Active policy holders (deletion after policy expiry+7yrs), (b) Ongoing claims (deletion after settlement+7yrs), (c) Under SAR investigation (deletion prohibited until case closed), (d) Regulatory hold (deletion requires IDRA/BFIU approval) | M |
| SEC-020 | Right to Erasure Workflow: Customer submits deletion request → System validates exceptions → If eligible, anonymize PII while retaining transaction records → Generate deletion certificate → Notify customer within 30 days | D |

### 8.4 Data Protection & Encryption Standards

| Data Classification | Encryption Standard | Key Management | Access Control |
|---------------------|-------------------|----------------|----------------|
| Personally Identifiable Information (PII) | AES-256 | AWS KMS with 90-day rotation | Role-based with audit logging |
| Financial Transaction Data | AES-256 + Additional Hashing | TigerBeetle built-in encryption | Restricted access with MFA |
| KYC Documents | AES-256 with client-side encryption | End-to-end encryption | Compliance officer access only |
| Medical Records | AES-256 with additional anonymization | Healthcare-specific key management | Medical staff + consent-based |
| Audit Logs | AES-256 with immutable storage | Centralized key management | Read-only access for auditors |

---

"""

# Append to file
with open(r"C:\_DEV\GO\InsureTech\V2.2_COMPLETE.md", "a", encoding="utf-8") as f:
    f.write(remaining_sections)

print("✅ Added sections 4.15-4.17, 5, 6, 7, and 8 to V2.2_COMPLETE.md")
print("📄 Current file size:", len(content) + len(remaining_sections), "characters")
