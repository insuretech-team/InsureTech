# InsureTech Proto - Core Business Modules

## Overview
This document details the core business domain modules that form the heart of the InsureTech platform.

---

# 1. Authentication Module (insuretech/authn/)

## Purpose
Manages user authentication, sessions, OTP, and identity verification documents.

## Key Files
- `entity/v1/user.proto` - User account
- `entity/v1/user_profile.proto` - User profile information
- `entity/v1/session.proto` - Session management
- `entity/v1/otp.proto` - One-Time Password
- `entity/v1/document_type.proto` - Identity document types
- `entity/v1/user_document.proto` - User's identity documents
- `services/v1/auth_service.proto` - Authentication RPC service

## Key Entities

### User Entity
```
id (UUID)                  - Primary key
email                      - Email address (PII, unique)
phone                      - Phone number (PII)
password_hash              - Hashed password (encrypted, log_redacted)
first_name                 - First name (PII)
last_name                  - Last name (PII)
status                     - ACTIVE, INACTIVE, SUSPENDED, DELETED
last_login_at              - Last login timestamp
password_changed_at        - Password last changed
two_factor_enabled         - 2FA status
roles                      - Array of role IDs
tenant_id                  - FK to tenant
audit_info                 - Standard audit trail
```

### UserProfile Entity
```
user_id (FK)               - Reference to User
date_of_birth              - Date of birth (PII)
gender                     - Gender
nationality                - Nationality
marital_status             - Marital status
occupation                 - Occupation
mother_maiden_name         - For security questions (PII)
profile_photo_url          - Profile photo storage URL
language_preference        - Preferred language
notification_preferences   - Notification settings
address                    - Complete address (Address type)
secondary_phone            - Backup phone (PII)
nid_number                 - National ID number (PII, encrypted)
nid_expiry_date           - NID expiry date
```

### Session Entity
```
id (UUID)                  - Session ID
user_id (FK)               - Reference to User
access_token               - JWT access token (encrypted, log_redacted)
refresh_token              - Refresh token (encrypted, log_redacted)
token_expires_at           - Token expiration
refresh_expires_at         - Refresh token expiration
ip_address                 - Client IP (PII, log_masked)
user_agent                 - Client user agent
device_info                - Device information
last_activity_at           - Last activity timestamp
active                     - Is session active?
```

### OTP Entity
```
id (UUID)                  - Primary key
user_id (FK)               - Reference to User
otp_code                   - One-time password (encrypted, log_redacted)
method                     - SMS, EMAIL, AUTHENTICATOR
channel                    - Phone or email address (PII, log_masked)
expires_at                 - OTP expiration (typically 5-10 minutes)
attempts                   - Failed attempt count
max_attempts               - Maximum attempts (typically 3)
verified                   - Is OTP verified?
verified_at                - Verification timestamp
reason                     - LOGIN, PASSWORD_RESET, TRANSACTION
```

### UserDocument Entity
```
id (UUID)                  - Primary key
user_id (FK)               - Reference to User
document_type              - NID, PASSPORT, DRIVING_LICENSE
document_number            - Document number (PII, encrypted)
issue_date                 - Date issued
expiry_date                - Expiration date
issuing_country            - Country of issuance
document_url               - Storage URL for document image
verified                   - Is document verified?
verification_date          - Date of verification
verification_method        - Manual, automated, biometric
document_front_url         - Front side image URL
document_back_url          - Back side image URL
```

## Auth Service RPC Methods

```protobuf
service AuthService {
  rpc Register(RegisterRequest) returns (RegisterResponse);
  rpc Login(LoginRequest) returns (LoginResponse);
  rpc Logout(LogoutRequest) returns (LogoutResponse);
  rpc RefreshToken(RefreshTokenRequest) returns (RefreshTokenResponse);
  rpc ValidateSession(ValidateSessionRequest) returns (ValidateSessionResponse);
  rpc RequestOTP(RequestOTPRequest) returns (RequestOTPResponse);
  rpc VerifyOTP(VerifyOTPRequest) returns (VerifyOTPResponse);
  rpc ChangePassword(ChangePasswordRequest) returns (ChangePasswordResponse);
  rpc ResetPassword(ResetPasswordRequest) returns (ResetPasswordResponse);
  rpc UpdateProfile(UpdateProfileRequest) returns (UpdateProfileResponse);
  rpc GetUserProfile(GetUserProfileRequest) returns (GetUserProfileResponse);
}
```

---

# 2. Authorization Module (insuretech/authz/)

## Purpose
Manages role-based access control (RBAC), policies, MFA, and access decision auditing.

## Key Files
- `entity/v1/role.proto` - Role definitions
- `entity/v1/user_role.proto` - User role assignments
- `entity/v1/policy_rule.proto` - RBAC policy rules
- `entity/v1/casbin_rule.proto` - Casbin RBAC model
- `entity/v1/access_decision_audit.proto` - Access decision audit
- `entity/v1/role_mfa_config.proto` - MFA configuration
- `entity/v1/token_config.proto` - Token expiration config
- `services/v1/authz_service.proto` - Authorization RPC service

## Key Entities

### Role Entity
```
id (UUID)                  - Primary key
name                       - Role name (ADMIN, USER, AGENT, etc.)
description                - Role description
permissions                - Array of permission codes
parent_role_id             - FK to parent role (for hierarchy)
tenant_id                  - FK to tenant (multi-tenancy)
status                     - ACTIVE, INACTIVE
created_at                 - Creation timestamp
updated_at                 - Last update timestamp
```

### UserRole Entity
```
id (UUID)                  - Primary key
user_id (FK)               - Reference to User
role_id (FK)               - Reference to Role
tenant_id (FK)             - Reference to Tenant
scope                      - Scope of role (SYSTEM, TENANT, DEPARTMENT)
assigned_at                - Assignment timestamp
assigned_by                - Assigned by user ID
expires_at                 - Role expiration (optional)
status                     - ACTIVE, INACTIVE
```

### PolicyRule Entity (RBAC Rules)
```
id (UUID)                  - Primary key
rule_name                  - Rule identifier
subject                    - Subject (role or user)
object                     - Object/Resource (policy, claim, etc.)
action                     - Action (read, write, delete, etc.)
effect                     - ALLOW or DENY
condition                  - Optional condition expression
priority                   - Rule priority for conflicts
tenant_id (FK)             - Reference to Tenant
active                     - Is rule active?
```

### CasbinRule Entity
```
id (UUID)                  - Primary key
p_type                     - Rule type (p, g, g2, etc.)
v0 - v5                    - Rule values
tenant_id (FK)             - Reference to Tenant
```

### RoleMFAConfig Entity
```
id (UUID)                  - Primary key
role_id (FK)               - Reference to Role
mfa_required               - Is MFA required for this role?
mfa_methods                - Array of allowed MFA methods
grace_period_days          - Grace period for MFA setup
enforcement_date           - When MFA becomes mandatory
status                     - ACTIVE, INACTIVE
```

### AccessDecisionAudit Entity
```
id (UUID)                  - Primary key
user_id (FK)               - Reference to User
resource_type              - Type of resource (policy, claim)
resource_id (FK)           - Reference to resource
action                     - Requested action (read, write, etc.)
decision                   - ALLOW or DENY
reason                     - Reason for decision
ip_address                 - Client IP (PII, log_masked)
timestamp                  - Decision timestamp
trace_id                   - Distributed trace ID
```

## Authorization Service RPC Methods

```protobuf
service AuthzService {
  rpc CheckAccess(CheckAccessRequest) returns (CheckAccessResponse);
  rpc ListRoles(ListRolesRequest) returns (ListRolesResponse);
  rpc CreateRole(CreateRoleRequest) returns (CreateRoleResponse);
  rpc UpdateRole(UpdateRoleRequest) returns (UpdateRoleResponse);
  rpc DeleteRole(DeleteRoleRequest) returns (DeleteRoleResponse);
  rpc AssignRole(AssignRoleRequest) returns (AssignRoleResponse);
  rpc RevokeRole(RevokeRoleRequest) returns (RevokeRoleResponse);
  rpc UpdatePolicy(UpdatePolicyRequest) returns (UpdatePolicyResponse);
  rpc GetAccessDecisions(GetAccessDecisionsRequest) returns (GetAccessDecisionsResponse);
}
```

---

# 3. Policy Module (insuretech/policy/)

## Purpose
Core policy management - the central entity in insurance operations.

## Key Files
- `entity/v1/policy.proto` - Main policy entity
- `entity/v1/quotation.proto` - Policy quotation
- `entity/v1/policy_service_request.proto` - Service requests on policies
- `events/v1/policy_events.proto` - Policy lifecycle events
- `services/v1/policy_service.proto` - Policy RPC service

## Key Entities

### Policy Entity (Core)
```
id (UUID)                  - Primary key
policy_number              - Unique policy number
insurer_id (FK)            - Reference to Insurer
product_id (FK)            - Reference to Product
plan_id (FK)               - Reference to ProductPlan
holder_id (FK)             - Reference to policyholder User
status                     - DRAFT, QUOTED, ACTIVE, SUSPENDED, LAPSED, CANCELLED, RENEWED
premium_amount             - Policy premium (Money type)
currency                   - Currency code
start_date                 - Coverage start date
end_date                   - Coverage end date
renewal_date               - Next renewal date
term_months                - Coverage period in months
coverage_details           - Coverage map/JSON
exclusions                 - Array of exclusions
deductible                 - Deductible amount (Money type)
max_coverage_amount        - Maximum coverage (Money type)
riders                     - Array of rider IDs
beneficiaries              - Array of beneficiary information
underwriting_decision      - Underwriting status
underwriting_notes         - Underwriting notes
payment_status             - PENDING, PARTIAL, PAID, OVERDUE
last_premium_paid_date     - Last payment date
next_premium_due_date      - Next premium due
payment_method             - Payment method selected
auto_renew                 - Is auto-renewal enabled?
grace_period_days          - Grace period after lapse
cancellation_reason        - Reason if cancelled
cancellation_date          - Cancellation date
refund_status              - PENDING, APPROVED, REJECTED, COMPLETED
refund_amount              - Refund amount if applicable
tenant_id (FK)             - Reference to Tenant
audit_info                 - Standard audit trail
```

### Quotation Entity
```
id (UUID)                  - Primary key
quote_number               - Unique quote reference
policy_id (FK)             - Reference to Policy (if converted)
insurer_id (FK)            - Reference to Insurer
product_id (FK)            - Reference to Product
holder_id (FK)             - Reference to applicant
premium_amount             - Quoted premium (Money type)
tax_amount                 - Tax/VAT (Money type)
total_amount               - Total quote (Money type)
coverage_details           - Proposed coverage
riders                     - Proposed riders
rate_factors               - Applied rate factors
risk_score                 - Risk assessment score
quote_expiry_date          - Quote validity period
status                     - DRAFT, QUOTED, ACCEPTED, REJECTED, EXPIRED
conversion_count           - How many times quoted
conversion_date            - Date converted to policy
validity_period_days       - Quote validity (typically 30-90 days)
tenant_id (FK)             - Reference to Tenant
```

### PolicyServiceRequest Entity
```
id (UUID)                  - Primary key
policy_id (FK)             - Reference to Policy
service_type               - MODIFICATION, ENDORSEMENT, CANCELLATION, CLAIM, DOCUMENT
request_status             - SUBMITTED, APPROVED, REJECTED, COMPLETED
description                - Service request details
requested_changes          - What needs to change (JSON)
approval_status            - Status of approval
approved_by                - Approved by user ID
approval_date              - Approval timestamp
effective_date             - When change is effective
request_date               - Request submission date
completion_date            - Request completion date
```

## Policy Service RPC Methods

```protobuf
service PolicyService {
  rpc CreateQuotation(CreateQuotationRequest) returns (CreateQuotationResponse);
  rpc GetQuotation(GetQuotationRequest) returns (GetQuotationResponse);
  rpc CreatePolicy(CreatePolicyRequest) returns (CreatePolicyResponse);
  rpc GetPolicy(GetPolicyRequest) returns (GetPolicyResponse);
  rpc ListPolicies(ListPoliciesRequest) returns (ListPoliciesResponse);
  rpc UpdatePolicy(UpdatePolicyRequest) returns (UpdatePolicyResponse);
  rpc CancelPolicy(CancelPolicyRequest) returns (CancelPolicyResponse);
  rpc RenewPolicy(RenewPolicyRequest) returns (RenewPolicyResponse);
  rpc ProcessPayment(ProcessPaymentRequest) returns (ProcessPaymentResponse);
}
```

## Policy Status Transitions
```
DRAFT 
  -> QUOTED (after quotation)
  -> ACTIVE (after payment)

ACTIVE 
  -> SUSPENDED (non-payment)
  -> LAPSED (end of grace period)
  -> CANCELLED (requested cancellation)
  -> RENEWED (automatic renewal)

LAPSED 
  -> ACTIVE (reinstatement)
  -> CANCELLED (decided not to reinstate)

CANCELLED 
  -> (final state, cannot transition)
```

---

# 4. Claims Module (insuretech/claims/)

## Purpose
Manages claim submission, assessment, approval, and settlement.

## Key Files
- `entity/v1/claim.proto` - Claim entity
- `events/v1/claim_events.proto` - Claim lifecycle events
- `services/v1/claim_service.proto` - Claims RPC service

## Key Entities

### Claim Entity (Core)
```
id (UUID)                  - Primary key
claim_number               - Unique claim reference
policy_id (FK)             - Reference to Policy
insurer_id (FK)            - Reference to Insurer
claimant_id (FK)           - Reference to claimant User
status                     - DRAFT, SUBMITTED, UNDER_REVIEW, APPROVED, REJECTED, SETTLED, APPEALED
claim_type                 - Health, Life, Auto, Home, etc.
incident_date              - Date of incident (PII consideration)
submission_date            - Date claim submitted
claim_amount               - Claimed amount (Money type)
approved_amount            - Approved amount (Money type)
currency                   - Currency code
description                - Claim description
documents                  - Array of document URLs
medical_reports            - Array of medical report URLs
police_report_url          - Police report (if applicable)
assessment_notes           - Assessment by adjuster
assessment_date            - Assessment completion date
assessor_id (FK)           - Reference to assessor User
fraud_check_status         - PENDING, PASSED, FLAGGED, UNDER_INVESTIGATION
fraud_score                - Fraud risk score (0-100)
ai_assessment              - AI agent assessment
approval_status            - Approved/Rejected by
approved_by (FK)           - Reference to approver User
approval_date              - Approval timestamp
rejection_reason           - Reason if rejected
settlement_date            - Settlement date
settlement_method          - Bank transfer, check, etc.
appeal_status              - If appealed
appeal_reason              - Reason for appeal
tenant_id (FK)             - Reference to Tenant
audit_info                 - Standard audit trail
```

## Claim Service RPC Methods

```protobuf
service ClaimService {
  rpc SubmitClaim(SubmitClaimRequest) returns (SubmitClaimResponse);
  rpc GetClaim(GetClaimRequest) returns (GetClaimResponse);
  rpc ListClaims(ListClaimsRequest) returns (ListClaimsResponse);
  rpc UpdateClaim(UpdateClaimRequest) returns (UpdateClaimResponse);
  rpc AssessClaim(AssessClaimRequest) returns (AssessClaimResponse);
  rpc ApproveClaim(ApproveClaimRequest) returns (ApproveClaimResponse);
  rpc RejectClaim(RejectClaimRequest) returns (RejectClaimResponse);
  rpc SettleClaim(SettleClaimRequest) returns (SettleClaimResponse);
  rpc AppealClaim(AppealClaimRequest) returns (AppealClaimResponse);
}
```

## Claim Status Transitions
```
DRAFT 
  -> SUBMITTED (submitted for processing)

SUBMITTED 
  -> UNDER_REVIEW (assessment started)

UNDER_REVIEW 
  -> APPROVED (assessment passed)
  -> REJECTED (assessment failed)
  -> FLAGGED (fraud detected)

APPROVED 
  -> SETTLED (payment made)

REJECTED 
  -> APPEALED (customer appeals)

APPEALED 
  -> APPROVED (appeal accepted)
  -> REJECTED (appeal denied)
```

---

# 5. Payment Module (insuretech/payment/)

## Purpose
Manages payment processing, reconciliation, and refunds.

## Key Files
- `entity/v1/payment.proto` - Payment entity
- `events/v1/payment_events.proto` - Payment events
- `services/v1/payment_service.proto` - Payment RPC service

## Key Entities

### Payment Entity
```
id (UUID)                  - Primary key
payment_reference          - Unique payment reference/transaction ID
policy_id (FK)             - Reference to Policy
payer_id (FK)              - Reference to payer User
payment_type               - PREMIUM, CLAIM_SETTLEMENT, REFUND
amount                     - Payment amount (Money type)
currency                   - Currency code
status                     - INITIATED, PENDING, COMPLETED, FAILED, CANCELLED
payment_method             - CARD, BANK_TRANSFER, MOBILE_WALLET, CASH
gateway_reference          - Payment gateway transaction ID
gateway_response           - Payment gateway response (JSON)
initiated_at               - Payment initiation timestamp
completed_at               - Payment completion timestamp
failed_at                  - Failure timestamp (if failed)
failure_reason             - Reason for failure
reconciled                 - Is payment reconciled?
reconciliation_date        - Reconciliation timestamp
receipt_url                - Payment receipt URL
ip_address                 - Payment IP (PII, log_masked)
device_info                - Device information
user_agent                 - Client user agent
three_d_secure             - 3DS verification status
idempotency_key            - Idempotency key (prevent duplicates)
metadata                   - Additional metadata (JSON)
tenant_id (FK)             - Reference to Tenant
audit_info                 - Standard audit trail
```

## Payment Service RPC Methods

```protobuf
service PaymentService {
  rpc InitiatePayment(InitiatePaymentRequest) returns (InitiatePaymentResponse);
  rpc GetPayment(GetPaymentRequest) returns (GetPaymentResponse);
  rpc ListPayments(ListPaymentsRequest) returns (ListPaymentsResponse);
  rpc ConfirmPayment(ConfirmPaymentRequest) returns (ConfirmPaymentResponse);
  rpc RefundPayment(RefundPaymentRequest) returns (RefundPaymentResponse);
  rpc ReconcilePayment(ReconcilePaymentRequest) returns (ReconcilePaymentResponse);
  rpc GetPaymentStatus(GetPaymentStatusRequest) returns (GetPaymentStatusResponse);
}
```

---

# 6. Product Module (insuretech/products/)

## Purpose
Defines insurance products, plans, riders, and pricing.

## Key Files
- `entity/v1/product.proto` - Product entity
- `entity/v1/product_plan.proto` - Product plans/variants
- `entity/v1/rider.proto` - Riders/add-ons
- `entity/v1/pricing_config.proto` - Pricing rules
- `events/v1/product_events.proto` - Product events
- `services/v1/product_service.proto` - Product RPC service

## Key Entities

### Product Entity
```
id (UUID)                  - Primary key
code                       - Product code (e.g., "HEALTH-001")
name                       - Product name
description                - Product description
insurance_type             - Health, Life, Auto, Home, Travel, etc.
insurer_id (FK)            - Reference to Insurer
status                     - ACTIVE, INACTIVE, RETIRED
category                   - Product category
target_market              - Target customer segment
min_age                    - Minimum applicant age
max_age                    - Maximum applicant age
min_sum_assured             - Minimum sum assured (Money type)
max_sum_assured             - Maximum sum assured (Money type)
base_premium               - Base premium (Money type)
underwriting_rules         - Underwriting rules (JSON)
exclusions                 - Array of exclusions
waiting_period_days        - Waiting period
coverage_period_months     - Coverage duration
renewable                  - Is product renewable?
riders_available           - Array of available rider IDs
image_url                  - Product image URL
document_urls              - Policy document URLs
launch_date                - Product launch date
sunset_date                - Product retirement date (if applicable)
tenant_id (FK)             - Reference to Tenant
audit_info                 - Standard audit trail
```

### ProductPlan Entity
```
id (UUID)                  - Primary key
product_id (FK)            - Reference to Product
plan_name                  - Plan variant name (e.g., "Silver", "Gold", "Platinum")
description                - Plan description
coverage_amount            - Sum assured (Money type)
premium_amount             - Plan premium (Money type)
term_months                - Term in months
deductible                 - Deductible (Money type)
co_insurance_percentage    - Co-insurance percentage
max_claim_frequency        - Max claims per year
status                     - ACTIVE, INACTIVE
order                      - Display order
```

### Rider Entity
```
id (UUID)                  - Primary key
code                       - Rider code
name                       - Rider name
description                - Rider description
product_ids (FK)           - Applicable to products
premium_addition           - Additional premium (Money type)
coverage_details           - Coverage details (JSON)
optional                   - Is rider optional?
status                     - ACTIVE, INACTIVE
```

### PricingConfig Entity
```
id (UUID)                  - Primary key
product_id (FK)            - Reference to Product
age_brackets               - Age-based rate factors
gender_factors             - Gender-based factors
health_factors             - Health-based factors
location_factors           - Location-based factors
occupation_factors         - Occupation-based factors
base_rate                  - Base rate percentage
markup_percentage          - Markup percentage
discount_slabs             - Volume discount slabs
tenant_id (FK)             - Reference to Tenant
```

---

## Key Cross-Module Relationships

```
User (authn)
  ├── UserProfile (authn)
  ├── Session (authn)
  ├── UserRole (authz)
  └── [ALL audit trails]

Tenant (tenant)
  └── [All other entities]

Insurer (insurer)
  ├── Product (products)
  │   ├── ProductPlan (products)
  │   └── Rider (products)
  └── Policy (policy)
      ├── Quotation (policy)
      ├── Claim (claims)
      ├── Payment (payment)
      ├── Beneficiary (beneficiary)
      ├── Endorsement (endorsement)
      └── RenewalSchedule (renewal)

Policy (policy)
  ├── Claim (claims)
  ├── Payment (payment)
  ├── Beneficiary (beneficiary)
  └── ServiceRequest (policy)
```

---

End of Core Modules Reference
