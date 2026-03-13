# 6. Detailed Business Functional Requirements (Complete Catalog)

This section enumerates business requirements derived from SRS V3.7 functional requirements.
Each requirement is phrased in business language and retains traceability to the original SRS FR-ID.

Notation
- **BR-ID**: Business Requirement Identifier (for BRD tracking)
- **SRS Trace**: SRS Feature Group and FR ID(s)
- **Priority**: aligned to SRS phase labels (M1/M2/M3/D/S/F)

## 6.1 4.1 User Management & Authentication (FG-001)

Business intent: define the outcomes this capability must deliver for customers/partners/admin teams.

### BR-01-001 — Support phone-based registration (Bangladesh mobile format: +880 1XXX XXXXXX) with OTP validation

- **SRS Trace:** FG-001 / FR-001
- **Priority:** M1
- **Business acceptance (summary):**
  - OTP sent within 60s, 6-digit code valid for 5 minutes
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-01-002 — Send OTP via SMS within 60 seconds with 6-digit code valid for 5 minutes

- **SRS Trace:** FG-001 / FR-002
- **Priority:** M1
- **Business acceptance (summary):**
  - 95% delivery success rate, retry on failure
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-01-003 — Allow maximum 3 OTP resend attempts per 15-minute window

- **SRS Trace:** FG-001 / FR-003
- **Priority:** M1
- **Business acceptance (summary):**
  - Rate limiting enforced, user notified on limit
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-01-004 — Enforce unique mobile number per account and detect duplicate registrations

- **SRS Trace:** FG-001 / FR-004
- **Priority:** M1
- **Business acceptance (summary):**
  - Error message on duplicate, database constraint enforced
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-01-005 — Support email-based registration with email verification link (24-hour validity)

- **SRS Trace:** FG-001 / FR-005
- **Priority:** M2
- **Business acceptance (summary):**
  - Verification email sent within 2 minutes, link expires after 24hrs
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-01-006 — Implement secure password policy: minimum 8 characters, 1 uppercase, 1 number, 1 special character

- **SRS Trace:** FG-001 / FR-006
- **Priority:** M1
- **Business acceptance (summary):**
  - Password strength indicator shown, validation enforced
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-01-007 — Provide biometric authentication (fingerprint/face ID) for mobile users

- **SRS Trace:** FG-001 / FR-007
- **Priority:** D
- **Business acceptance (summary):**
  - Device biometric API integration, fallback to password
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-01-008 — Support password reset via OTP to registered mobile number

- **SRS Trace:** FG-001 / FR-008
- **Priority:** M1
- **Business acceptance (summary):**
  - Reset OTP sent within 60s, new password saved securely
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-01-009 — Implement session management with Secure Token Service (15-minute access, 7-day refresh)

- **SRS Trace:** FG-001 / FR-009
- **Priority:** M1
- **Business acceptance (summary):**
  - Token rotation implemented, refresh token stored securely
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-01-010 — Enforce account lockout after 5 failed login attempts for 30 minutes

- **SRS Trace:** FG-001 / FR-010
- **Priority:** M2
- **Business acceptance (summary):**
  - Lockout triggered automatically, user notified via SMS
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-01-011 — Maintain user profile with: full name, date of birth, gender, occupation, address

- **SRS Trace:** FG-001 / FR-011
- **Priority:** M1
- **Business acceptance (summary):**
  - All mandatory fields validated, profile completeness indicator
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-01-012 — Support profile photo upload with validation (max 5MB, JPEG/PNG, face detection)

- **SRS Trace:** FG-001 / FR-012
- **Priority:** M3
- **Business acceptance (summary):**
  - Image compressed to <2MB, face detection validates single face
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-01-013 — Have stakeholders registration via SAML Identity provider

- **SRS Trace:** FG-001 / FR-013
- **Priority:** D
- **Business acceptance (summary):**
  - SAML 2.0 integration with Azure AD/Okta, SSO enabled
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

[[[PAGEBREAK]]]

## 6.2 4.2 Authorization & Access Control (FG-002)

Business intent: define the outcomes this capability must deliver for customers/partners/admin teams.

### BR-02-014 — Implement Role-Based Access Control (RBAC) with predefined roles: System Admin, Business Admin, Focal Person, Partner Admin, Agent, Customer

- **SRS Trace:** FG-002 / FR-014
- **Priority:** M1
- **Business acceptance (summary):**
  - Roles enforced at API gateway level, permissions validated on each request
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-02-015 — Enforce Attribute-Based Access Control (ABAC) for fine-grained permissions based on user attributes, resource type, and context

- **SRS Trace:** FG-002 / FR-015
- **Priority:** M1
- **Business acceptance (summary):**
  - Dynamic policy evaluation <50ms, audit logs for all authorization decisions
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-02-016 — Implement tenant isolation for partner organizations with data segregation

- **SRS Trace:** FG-002 / FR-016
- **Priority:** M2
- **Business acceptance (summary):**
  - Multi-tenant database architecture, row-level security enforced
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-02-017 — Enforce 2FA (Two-Factor Authentication) for all admin-level access

- **SRS Trace:** FG-002 / FR-017
- **Priority:** M3
- **Business acceptance (summary):**
  - TOTP-based 2FA with 30-second rotation, backup codes provided
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-02-018 — Maintain Access Control Lists (ACL) for resource-level permissions

- **SRS Trace:** FG-002 / FR-018
- **Priority:** M1
- **Business acceptance (summary):**
  - ACL stored in database, cached in Redis for performance
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-02-019 — Implement hierarchical role inheritance (Partner Admin > Agent > Customer)

- **SRS Trace:** FG-002 / FR-019
- **Priority:** D
- **Business acceptance (summary):**
  - Child roles inherit parent permissions, override capability available
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-02-020 — Provide permission audit trail for all sensitive operations

- **SRS Trace:** FG-002 / FR-020
- **Priority:** M3
- **Business acceptance (summary):**
  - Immutable audit logs, queryable by role/user/action/timestamp
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

[[[PAGEBREAK]]]

## 6.3 4.3 Product Management & Catalog (FG-003)

Business intent: define the outcomes this capability must deliver for customers/partners/admin teams.

### BR-03-021 — Provide product catalog with categorization: Health, Life, Motor, Travel, Micro-insurance

- **SRS Trace:** FG-003 / FR-021
- **Priority:** M1
- **Business acceptance (summary):**
  - Products displayed by category, search and filter enabled
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-03-022 — Support product search by name, category, coverage type, and premium range

- **SRS Trace:** FG-003 / FR-022
- **Priority:** M1
- **Business acceptance (summary):**
  - Search results <500ms, fuzzy matching for Bengali text
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-03-023 — Display product details: coverage, premium, tenure, exclusions, terms & conditions

- **SRS Trace:** FG-003 / FR-023
- **Priority:** M2
- **Business acceptance (summary):**
  - All product information visible before purchase, PDF download available
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-03-024 — Provide premium calculator with dynamic inputs (age, sum assured, tenure, riders)

- **SRS Trace:** FG-003 / FR-024
- **Priority:** M3
- **Business acceptance (summary):**
  - Real-time calculation <2s, breakdown of premium components shown
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-03-025 — Support product comparison (side-by-side up to 3 products)

- **SRS Trace:** FG-003 / FR-025
- **Priority:** M3
- **Business acceptance (summary):**
  - Comparison table with key features, coverage, and pricing
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-03-026 — Enable Business Admin to create, update, and deactivate products

- **SRS Trace:** FG-003 / FR-026
- **Priority:** M1
- **Business acceptance (summary):**
  - Product CRUD operations, version history maintained
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-03-027 — Support product variants with configurable riders and add-ons

- **SRS Trace:** FG-003 / FR-027
- **Priority:** M3
- **Business acceptance (summary):**
  - Base product + optional riders, dynamic pricing recalculation
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-03-028 — Cache product catalog in Redis with 5-minute TTL for performance

- **SRS Trace:** FG-003 / FR-028
- **Priority:** M3
- **Business acceptance (summary):**
  - Cache hit rate >80%, automatic invalidation on product updates
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-03-029 — Support multi-language product descriptions (Bengali and English)

- **SRS Trace:** FG-003 / FR-029
- **Priority:** M3
- **Business acceptance (summary):**
  - Language toggle in UI, content stored in i18n format
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

[[[PAGEBREAK]]]

## 6.4 4.4 Policy Lifecycle Management (FG-004)

Business intent: define the outcomes this capability must deliver for customers/partners/admin teams.

### BR-04-030 — Support end-to-end policy purchase flow: product selection → applicant details → nominee details → payment → policy issuance

- **SRS Trace:** FG-004 / FR-030
- **Priority:** M1
- **Business acceptance (summary):**
  - Complete flow in <10 minutes, progress saved at each step
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-04-031 — Collect applicant information: full name, DOB, NID, address, occupation, income, health declaration

- **SRS Trace:** FG-004 / FR-031
- **Priority:** M1
- **Business acceptance (summary):**
  - All mandatory fields validated, conditional fields based on product type
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-04-032 — Support multiple nominee/beneficiary addition with relationship and share percentage (must sum to 100%)

- **SRS Trace:** FG-004 / FR-032
- **Priority:** M1
- **Business acceptance (summary):**
  - Minimum 1 nominee required, share percentage validation enforced
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-04-033 — Validate NID uniqueness across policies to prevent duplicate insurance

- **SRS Trace:** FG-004 / FR-033
- **Priority:** M1
- **Business acceptance (summary):**
  - Database constraint enforced, user notified of existing policies
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-04-034 — Generate unique policy number with format: LBT-YYYY-XXXX-NNNNNN

- **SRS Trace:** FG-004 / FR-034
- **Priority:** M1
- **Business acceptance (summary):**
  - Sequential numbering, year-based prefix, collision prevention
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-04-035 — Issue digital policy document (PDF) with QR code for verification

- **SRS Trace:** FG-004 / FR-035
- **Priority:** M2
- **Business acceptance (summary):**
  - PDF generated within 30s of payment confirmation, QR code scannable
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-04-036 — Send policy document via SMS link and email attachment

- **SRS Trace:** FG-004 / FR-036
- **Priority:** M2
- **Business acceptance (summary):**
  - Delivery within 5 minutes, retry mechanism on failure
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-04-037 — Activate policy immediately upon payment confirmation for instant coverage

- **SRS Trace:** FG-004 / FR-037
- **Priority:** M2
- **Business acceptance (summary):**
  - Policy status updated in real-time, customer notified
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-04-038 — Support policy cooling-off period (15 days from issuance) for full refund

- **SRS Trace:** FG-004 / FR-038
- **Priority:** M3
- **Business acceptance (summary):**
  - Cancellation request processed within 24hrs, refund initiated
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-04-039 — Maintain policy status: Pending Payment, Active, Suspended, Cancelled, Lapsed, Expired

- **SRS Trace:** FG-004 / FR-039
- **Priority:** M1
- **Business acceptance (summary):**
  - Status transitions logged with timestamp, notifications triggered
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-04-040 — Provide customer policy dashboard showing all active and past policies, renewal prompts, and premium payment history

- **SRS Trace:** FG-004 / FR-040
- **Priority:** M1
- **Business acceptance (summary):**
  - Dashboard loads <3s, real-time status updates
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

[[[PAGEBREAK]]]

## 6.5 4.5 Policy Management & Renewals (FG-005)

Business intent: define the outcomes this capability must deliver for customers/partners/admin teams.

### BR-05-084 — Implement 'Family Insurance Wallet' allowing users to group and manage policies for multiple family members under one account

- **SRS Trace:** FG-005 / FR-084
- **Priority:** D
- **Business acceptance (summary):**
  - Unified dashboard, single-click bulk payment, relationship management
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-05-085 — Send renewal reminders: 30 days, 15 days, 7 days, 1 day before expiry via SMS, email, push notification

- **SRS Trace:** FG-005 / FR-085
- **Priority:** M2
- **Business acceptance (summary):**
  - Notifications sent on schedule, delivery confirmation tracked
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-05-086 — Support manual policy renewal with one-click process reusing existing policy data

- **SRS Trace:** FG-005 / FR-086
- **Priority:** M2
- **Business acceptance (summary):**
  - Renewal completed in <3 minutes, updated policy document issued
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-05-087 — Support automatic policy renewal with stored payment method (opt-in by customer)

- **SRS Trace:** FG-005 / FR-087
- **Priority:** M3
- **Business acceptance (summary):**
  - Customer consent recorded, auto-charge 7 days before expiry
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-05-088 — Allow customer to update policy details during renewal: current address, nominee information

- **SRS Trace:** FG-005 / FR-088
- **Priority:** M3
- **Business acceptance (summary):**
  - Limited fields editable, verification required for major changes
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-05-089 — Implement grace period (30 days) for premium payment post-expiry with continued coverage

- **SRS Trace:** FG-005 / FR-089
- **Priority:** M2
- **Business acceptance (summary):**
  - Policy status "Grace Period", coverage continues, daily reminders
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-05-090 — Auto-lapse policy after grace period if payment not received, with reinstatement option

- **SRS Trace:** FG-005 / FR-090
- **Priority:** M2
- **Business acceptance (summary):**
  - Policy status "Lapsed", reinstatement within 90 days with penalty
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-05-091 — Provide policy document download (PDF) with version history for all renewals

- **SRS Trace:** FG-005 / FR-091
- **Priority:** M1
- **Business acceptance (summary):**
  - All versions accessible, clearly marked with issue date
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-05-092 — Track policy lifecycle events: issuance, renewal, lapse, reinstatement, cancellation with audit trail

- **SRS Trace:** FG-005 / FR-092
- **Priority:** M1
- **Business acceptance (summary):**
  - Immutable event log, queryable by date range and policy number
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-05-093 — Support policy cancellation workflow with cancellation request submission by customer/agent/admin

- **SRS Trace:** FG-005 / FR-093
- **Priority:** M1
- **Business acceptance (summary):**
  - Request form with reason dropdown, attachment support
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-05-094 — Implement approval workflow for policy cancellation: Business Admin + Focal Person approval required for policies >30 days old

- **SRS Trace:** FG-005 / FR-094
- **Priority:** M1
- **Business acceptance (summary):**
  - Approval routing, 48hr SLA
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-05-095 — Calculate pro-rata refund: (Premium Paid - Days Covered - Admin Fee - Cancellation Charge) with transparent breakdown

- **SRS Trace:** FG-005 / FR-095
- **Priority:** M1
- **Business acceptance (summary):**
  - Refund calculator, configurable fees
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-05-096 — Process refund within 7 working days via MFS or bank transfer

- **SRS Trace:** FG-005 / FR-096
- **Priority:** M1
- **Business acceptance (summary):**
  - Payment gateway integration, notifications
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-05-097 — Update policy status to CANCELLED and notify all stakeholders

- **SRS Trace:** FG-005 / FR-097
- **Priority:** M1
- **Business acceptance (summary):**
  - Multi-channel notification, IDRA reporting
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-05-098 — Support policy endorsement for: Address, Sum insured, Nominee, Contact changes

- **SRS Trace:** FG-005 / FR-098
- **Priority:** M1
- **Business acceptance (summary):**
  - Amendment forms, validation
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-05-099 — Calculate additional premium for mid-term sum insured increases

- **SRS Trace:** FG-005 / FR-099
- **Priority:** M1
- **Business acceptance (summary):**
  - Premium calculator, payment integration
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-05-100 — Calculate pro-rata refund for sum insured decreases

- **SRS Trace:** FG-005 / FR-100
- **Priority:** M2
- **Business acceptance (summary):**
  - Credit to premium account
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-05-101 — Generate endorsement document with suffix (POL-001/END-01)

- **SRS Trace:** FG-005 / FR-101
- **Priority:** M1
- **Business acceptance (summary):**
  - PDF generation, version tracking
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-05-102 — Require approval for sum insured changes >10%

- **SRS Trace:** FG-005 / FR-102
- **Priority:** M1
- **Business acceptance (summary):**
  - Approval workflow, threshold config
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

[[[PAGEBREAK]]]

## 6.6 4.6 Business Rules & Workflows (FG-06)

Business intent: define the outcomes this capability must deliver for customers/partners/admin teams.

### BR-06-214 — Implement premium calculation fallbacks: If insurer API fails, use cached rates (max 24hrs old); if unavailable, queue quote and notify customer within 2 hours

- **SRS Trace:** FG-06 / FR-214
- **Priority:** M1
- **Business acceptance (summary):**
  - • Fallback logic tested / • Cache validation / • Queue notification works
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-06-215 — Handle premium calculation edge cases: age-based loading, occupation risk factors, pre-existing conditions with clear messaging

- **SRS Trace:** FG-06 / FR-215
- **Priority:** M2
- **Business acceptance (summary):**
  - • All edge cases covered / • Messaging user-friendly / • Actuarial validation
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-06-216 — Implement duplicate policy detection: Block duplicate policy purchase for same product + same insured person within 30 days; allow cross-product purchases

- **SRS Trace:** FG-06 / FR-216
- **Priority:** M1
- **Business acceptance (summary):**
  - • Detection accurate / • Cross-product allowed / • Clear error message
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-06-217 — Enable policy merge workflow: Focal Person can merge duplicate accounts after verifying NID, transfer policies, consolidate claims history

- **SRS Trace:** FG-06 / FR-217
- **Priority:** M3
- **Business acceptance (summary):**
  - • Merge workflow tested / • Data integrity maintained / • Audit logged
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-06-218 — Define claim status state machine: Submitted → Under Review → Documents Requested → Approved/Rejected → Payment Initiated → Settled/Closed

- **SRS Trace:** FG-06 / FR-218
- **Priority:** M1
- **Business acceptance (summary):**
  - • State machine implemented / • Invalid transitions blocked / • Status tracking accurate
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-06-219 — Enforce claim status transition rules: Auto-move to "Documents Requested" if incomplete; require Business Admin+Focal Person approval for >BDT 50K

- **SRS Trace:** FG-06 / FR-219
- **Priority:** M1
- **Business acceptance (summary):**
  - • Transition rules enforced / • Approval routing correct / • Notifications sent
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-06-220 — Implement gamified renewal rewards program offering discounts or gift vouchers for early renewals

- **SRS Trace:** FG-06 / FR-220
- **Priority:** D
- **Business acceptance (summary):**
  - Points calculation engine, partner voucher integration, leaderboard
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-06-221 — Implement grace period logic: 30-day grace period post-expiry with coverage continued; auto-lapse if unpaid after grace period

- **SRS Trace:** FG-06 / FR-221
- **Priority:** M3
- **Business acceptance (summary):**
  - • Grace period enforced / • Coverage continued / • Auto-lapse works / • Customer notified
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-06-222 — Enable lapsed policy reinstatement: Allow reinstatement within 90 days of lapse with medical underwriting; require Focal Person approval

- **SRS Trace:** FG-06 / FR-222
- **Priority:** D
- **Business acceptance (summary):**
  - • Reinstatement workflow / • Medical underwriting integrated / • Approval required
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

[[[PAGEBREAK]]]

## 6.7 4.7 Payment Processing (FG-007)

Business intent: define the outcomes this capability must deliver for customers/partners/admin teams.

### BR-07-073 — Support multiple payment methods: bKash, Nagad, Rocket, Bank Transfer, Credit/Debit Card, Manual Cash/Cheque

- **SRS Trace:** FG-007 / FR-073
- **Priority:** M1
- **Business acceptance (summary):**
  - All MFS integrated, card via hosted payment page, manual verification
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-07-074 — Integrate bKash payment gateway with production credentials and sandbox for testing

- **SRS Trace:** FG-007 / FR-074
- **Priority:** M1
- **Business acceptance (summary):**
  - Transaction success rate >99%, fallback to manual on failure
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-07-075 — Integrate Nagad and Rocket MFS with tokenization for recurring payments

- **SRS Trace:** FG-007 / FR-075
- **Priority:** M3
- **Business acceptance (summary):**
  - Secure token storage, PCI-DSS Level SAQ-A compliance
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-07-076 — Support manual payment with proof upload (bank receipt, bKash screenshot) for verification

- **SRS Trace:** FG-007 / FR-076
- **Priority:** M1
- **Business acceptance (summary):**
  - Image upload <5MB, admin verification within 24hrs
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-07-077 — Implement payment verification workflow: pending → verified → policy activated OR rejected → refund

- **SRS Trace:** FG-007 / FR-077
- **Priority:** M2
- **Business acceptance (summary):**
  - Admin approval for manual payments, automated for MFS
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-07-078 — Generate payment receipt with transaction ID, amount, date, policy number

- **SRS Trace:** FG-007 / FR-078
- **Priority:** M2
- **Business acceptance (summary):**
  - PDF receipt sent via SMS/email within 5 minutes
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-07-079 — Support partial payment and installment plans for high-premium policies (quarterly, half-yearly, annual)

- **SRS Trace:** FG-007 / FR-079
- **Priority:** M3
- **Business acceptance (summary):**
  - Auto-reminders before due date, grace period 15 days
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-07-080 — Implement payment retry mechanism with exponential backoff for failed transactions

- **SRS Trace:** FG-007 / FR-080
- **Priority:** M2
- **Business acceptance (summary):**
  - Max 3 retries, customer notified on each attempt
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-07-081 — Support refund processing for policy cancellation with configurable refund rules

- **SRS Trace:** FG-007 / FR-081
- **Priority:** M2
- **Business acceptance (summary):**
  - Refund initiated within 7 days, credited to original payment method
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-07-082 — Integrate TigerBeetle for financial transaction recording with double-entry bookkeeping

- **SRS Trace:** FG-007 / FR-082
- **Priority:** M2
- **Business acceptance (summary):**
  - All transactions recorded, real-time balance reconciliation
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-07-083 — Maintain payment audit trail with immutable logs for regulatory compliance

- **SRS Trace:** FG-007 / FR-083
- **Priority:** M1
- **Business acceptance (summary):**
  - PostgreSQL + S3 storage, 20-year retention
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

[[[PAGEBREAK]]]

## 6.8 4.8 Claims Management (FG-008)

Business intent: define the outcomes this capability must deliver for customers/partners/admin teams.

### BR-08-041 — Provide fixed-step claim submission form: policy selection, incident details, claim reason, document upload (images, bills, reports)

- **SRS Trace:** FG-008 / FR-041
- **Priority:** M1
- **Business acceptance (summary):**
  - Form completion <5 minutes, draft saving at each step
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-08-042 — Validate claim eligibility: policy active, within coverage period, claim type covered, no duplicate submission

- **SRS Trace:** FG-008 / FR-042
- **Priority:** M1
- **Business acceptance (summary):**
  - Validation in <3s, clear error messages on rejection
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-08-043 — Generate unique claim number with format: CLM-YYYY-XXXX-NNNNNN and digital hash for submission integrity

- **SRS Trace:** FG-008 / FR-043
- **Priority:** M1
- **Business acceptance (summary):**
  - Collision-free numbering, SHA-256 hash for document integrity
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-08-044 — Automatically notify partner/insurer upon claim submission with shared status dashboard

- **SRS Trace:** FG-008 / FR-044
- **Priority:** M2
- **Business acceptance (summary):**
  - Notification within 60s, dashboard accessible to all stakeholders
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-08-045 — Provide real-time claim status tracking: Submitted, Under Review, Approved, Rejected, Settled

- **SRS Trace:** FG-008 / FR-045
- **Priority:** M3
- **Business acceptance (summary):**
  - Status updates visible in <5s, push notifications on status change
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-08-046 — Implement tiered approval workflow based on claim amount as per Approval Matrix

- **SRS Trace:** FG-008 / FR-046
- **Priority:** M3
- **Business acceptance (summary):**
  - Auto-routing to correct approver, escalation on timeout
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-08-047 — Support document verification with image quality check, OCR extraction, and fraud detection

- **SRS Trace:** FG-008 / FR-047
- **Priority:** M3
- **Business acceptance (summary):**
  - Image validation <10s, OCR accuracy >85%, duplicate detection
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-08-048 — Provide chat interface between customer, partner agent, and focal person for claim discussion

- **SRS Trace:** FG-008 / FR-048
- **Priority:** M3
- **Business acceptance (summary):**
  - Real-time messaging, file attachment support, message history
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-08-049 — Support WebRTC video call for claim verification and inspection

- **SRS Trace:** FG-008 / FR-049
- **Priority:** D
- **Business acceptance (summary):**
  - HD video quality, screen sharing, call recording for audit
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-08-050 — Allow partner to add verification notes and approve/reject with reason

- **SRS Trace:** FG-008 / FR-050
- **Priority:** M2
- **Business acceptance (summary):**
  - Notes timestamped, approval requires mandatory reason field
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-08-051 — Enforce joint approval by Business Admin and Focal Person for claims BDT 50K-2L

- **SRS Trace:** FG-008 / FR-051
- **Priority:** M3
- **Business acceptance (summary):**
  - Both approvals required, timeout escalation after 5 days
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-08-052 — Automate payment process upon claim approval as per customer's selected payment channel

- **SRS Trace:** FG-008 / FR-052
- **Priority:** M3
- **Business acceptance (summary):**
  - Payment initiated within 24hrs, confirmation sent to customer
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-08-053 — Support Zero Human Touch Claims (ZHTC) for auto-verification and payment of small claims (<BDT 10K) with partner pre-agreement

- **SRS Trace:** FG-008 / FR-053
- **Priority:** D
- **Business acceptance (summary):**
  - 95% automation rate, ML-based fraud check, instant settlement
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-08-054 — Implement fraud detection: frequent claims (>3 in 6 months), duplicate documents, rapid policy-to-claim (<48hrs)

- **SRS Trace:** FG-008 / FR-054
- **Priority:** M3
- **Business acceptance (summary):**
  - Auto-flagging with risk score, manual review queue, customer warning system
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-08-055 — Auto-revoke customer access for confirmed fraud as per InsureTech policy

- **SRS Trace:** FG-008 / FR-055
- **Priority:** M3
- **Business acceptance (summary):**
  - Account suspension after approval, appeal process available
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-08-056 — Maintain balance sheet on Customer, Partner, Agent, and InsureTech level for selected time periods

- **SRS Trace:** FG-008 / FR-056
- **Priority:** M3
- **Business acceptance (summary):**
  - Daily, monthly, quarterly reconciliation, export to Excel/PDF
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-08-057 — Track Turn Around Time (TAT) per approval level and alert on SLA breach

- **SRS Trace:** FG-008 / FR-057
- **Priority:** M3
- **Business acceptance (summary):**
  - Real-time TAT monitoring, email alerts on approaching deadline
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-08-058 — Provide claim history and analytics for risk assessment and premium adjustment

- **SRS Trace:** FG-008 / FR-058
- **Priority:** M3
- **Business acceptance (summary):**
  - Claim frequency report, average claim amount, settlement ratio
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-08-103 — Enforce claims document requirements: PDF/JPG/PNG, max 10MB per file, 50MB total per claim, 300 DPI minimum

- **SRS Trace:** FG-008 / FR-103
- **Priority:** M1
- **Business acceptance (summary):**
  - Client-side validation, OCR quality check
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-08-104 — Calculate co-payment and deductibles: (Claim Amount - Deductible) × Co-payment % with annual deductible tracking

- **SRS Trace:** FG-008 / FR-104
- **Priority:** M1
- **Business acceptance (summary):**
  - Product-level config, breakdown display
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-08-105 — Support claims reimbursement workflow with document review and bank/MFS transfer within 7-15 working days

- **SRS Trace:** FG-008 / FR-105
- **Priority:** M1
- **Business acceptance (summary):**
  - Document verification, payment processing, status notifications
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

[[[PAGEBREAK]]]

## 6.9 4.9 Partner & Agent Management (FG-009)

Business intent: define the outcomes this capability must deliver for customers/partners/admin teams.

### BR-09-059 — Support partner onboarding workflow: application submission, KYB verification, MOU upload, approval by Focal Person

- **SRS Trace:** FG-009 / FR-059
- **Priority:** M2
- **Business acceptance (summary):**
  - Complete onboarding in <7 days, status tracking at each step
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-09-060 — Collect partner information: organization name, type (hospital/MFS/e-commerce/agent), trade license, TIN, bank account, contact details

- **SRS Trace:** FG-009 / FR-060
- **Priority:** M2
- **Business acceptance (summary):**
  - All mandatory fields validated, document verification required
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-09-061 — Implement KYB (Know Your Business) verification with trade license validation and credit check

- **SRS Trace:** FG-009 / FR-061
- **Priority:** M2
- **Business acceptance (summary):**
  - Automated validation where possible, manual review for exceptions
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-09-062 — Provide dedicated partner portal with dashboard showing: leads, conversions, commissions, analytics

- **SRS Trace:** FG-009 / FR-062
- **Priority:** M2
- **Business acceptance (summary):**
  - Dashboard loads <3s, real-time data updates, export functionality
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-09-063 — Calculate and track partner commissions based on configurable rates (acquisition, renewal, claims assistance)

- **SRS Trace:** FG-009 / FR-063
- **Priority:** M2
- **Business acceptance (summary):**
  - Commission calculated on policy activation, monthly payout reports
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-09-064 — Support partner API integration for embedded insurance (e-commerce checkout, hospital admission)

- **SRS Trace:** FG-009 / FR-064
- **Priority:** M3
- **Business acceptance (summary):**
  - RESTful API with sandbox, developer documentation, webhook support
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-09-065 — Enable partner to initiate policy purchase on behalf of customer with consent and authentication

- **SRS Trace:** FG-009 / FR-065
- **Priority:** M2
- **Business acceptance (summary):**
  - Customer OTP verification required, policy linked to customer account
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-09-066 — Provide Focal Person portal for partner management: verification, approval, dispute resolution, performance monitoring

- **SRS Trace:** FG-009 / FR-066
- **Priority:** M1
- **Business acceptance (summary):**
  - Full CRUD operations on partners, approval workflow, audit trail
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-09-067 — Support multi-level agent hierarchy under partners (Partner Admin > Regional Manager > Agent)

- **SRS Trace:** FG-009 / FR-067
- **Priority:** M3
- **Business acceptance (summary):**
  - Hierarchical commission split, territory management, performance tracking
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-09-068 — Track partner performance metrics: policies sold, claim settlement ratio, customer satisfaction, fraud incidents

- **SRS Trace:** FG-009 / FR-068
- **Priority:** M2
- **Business acceptance (summary):**
  - Weekly/monthly reports, performance scoring, alerts on anomalies
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-09-069 — Support partner suspension/termination with graceful policy transfer mechanism

- **SRS Trace:** FG-009 / FR-069
- **Priority:** M2
- **Business acceptance (summary):**
  - Existing policies remain active, new sales blocked, customer notification
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-09-070 — Focal Person shall have authority to verify and approve/reject partner applications within 3 business days

- **SRS Trace:** FG-009 / FR-070
- **Priority:** M1
- **Business acceptance (summary):**
  - Decision recorded with reason, partner notified automatically
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-09-071 — Focal Person shall monitor partner compliance and flag suspicious activities for investigation

- **SRS Trace:** FG-009 / FR-071
- **Priority:** M2
- **Business acceptance (summary):**
  - Real-time dashboard with alerts, escalation to Business Admin
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-09-072 — Focal Person shall resolve partner-customer disputes with documented decision trail

- **SRS Trace:** FG-009 / FR-072
- **Priority:** M2
- **Business acceptance (summary):**
  - Dispute resolution within 7 days, audit log maintained
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

[[[PAGEBREAK]]]

## 6.10 4.10 Partner Portal & Business Intelligence (FG-010)

Business intent: define the outcomes this capability must deliver for customers/partners/admin teams.

### BR-10-141 — Provide hospital partners special dashboard to initiate insurance purchase on behalf of customers

- **SRS Trace:** FG-010 / FR-141
- **Priority:** M2
- **Business acceptance (summary):**
  - Patient data prefill from hospital system, consent capture
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-10-142 — Support API for transferring customer records with authentication token and purchase ID

- **SRS Trace:** FG-010 / FR-142
- **Priority:** D
- **Business acceptance (summary):**
  - RESTful API with OAuth2, data mapping documentation
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-10-143 — Provide e-commerce partners embedded widget for insurance product display at checkout

- **SRS Trace:** FG-010 / FR-143
- **Priority:** M2
- **Business acceptance (summary):**
  - JavaScript SDK, responsive design, cart integration
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-10-144 — Provide sandbox environment for 3rd party developers with test credentials and mock data

- **SRS Trace:** FG-010 / FR-144
- **Priority:** D
- **Business acceptance (summary):**
  - Isolated test environment, sample code, API documentation
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-10-145 — Provide partner analytics: leads generated, conversion rate, commission earned, customer feedback

- **SRS Trace:** FG-010 / FR-145
- **Priority:** M2
- **Business acceptance (summary):**
  - Dashboard with filters, trend charts, export to Excel/PDF
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-10-146 — Provide partner API for retrieving analytics and commission statements programmatically

- **SRS Trace:** FG-010 / FR-146
- **Priority:** D
- **Business acceptance (summary):**
  - RESTful API, pagination support, webhook for new data
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-10-147 — Implement Business Intelligence tool (Metabase/Tableau/Power BI) for advanced analytics

- **SRS Trace:** FG-010 / FR-147
- **Priority:** F
- **Business acceptance (summary):**
  - Read replica connection, pre-built dashboards, scheduled reports
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-10-148 — Provide executive dashboard: daily sales, policy count, claims ratio, revenue, system health

- **SRS Trace:** FG-010 / FR-148
- **Priority:** M2
- **Business acceptance (summary):**
  - Real-time data, drill-down capability, mobile-responsive
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-10-205 — Provide partner-specific branding capability for white-label insurance offerings

- **SRS Trace:** FG-010 / FR-205
- **Priority:** F
- **Business acceptance (summary):**
  - Custom logo, colors, domain mapping, isolated tenant data
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-10-206 — Enable partners to configure commission structures and incentive programs

- **SRS Trace:** FG-010 / FR-206
- **Priority:** D
- **Business acceptance (summary):**
  - Tiered commission, bonus rules, performance-based adjustments
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-10-207 — Log all API requests with payload, headers, timestamps

- **SRS Trace:** FG-010 / FR-207
- **Priority:** M2
- **Business acceptance (summary):**
  - Structured logging, rotation, searchable
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-10-208 — Implement distributed tracing across microservices

- **SRS Trace:** FG-010 / FR-208
- **Priority:** D
- **Business acceptance (summary):**
  - Jaeger integration, trace ID propagation
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

[[[PAGEBREAK]]]

## 6.11 4.11 Customer Support & Helpdesk (FG-011)

Business intent: define the outcomes this capability must deliver for customers/partners/admin teams.

### BR-11-106 — Provide in-app FAQ section with searchable knowledge base covering common queries

- **SRS Trace:** FG-011 / FR-106
- **Priority:** M1
- **Business acceptance (summary):**
  - Search results <1s, categorized by topic, Bengali and English
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-11-107 — Support customer support call initiation from mobile app with call recording

- **SRS Trace:** FG-011 / FR-107
- **Priority:** M3
- **Business acceptance (summary):**
  - Click-to-call integration, call routing to available agent
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-11-108 — Implement ticketing system for customer issues with unique ticket ID and status tracking

- **SRS Trace:** FG-011 / FR-108
- **Priority:** M2
- **Business acceptance (summary):**
  - Ticket creation <30s, status updates via notification
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-11-109 — Provide support agent portal with ticket queue, customer history, and resolution templates

- **SRS Trace:** FG-011 / FR-109
- **Priority:** M2
- **Business acceptance (summary):**
  - Agent dashboard loads <3s, SLA countdown visible
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-11-110 — Auto-record customer support calls and create ticket with call summary

- **SRS Trace:** FG-011 / FR-110
- **Priority:** M3
- **Business acceptance (summary):**
  - Speech-to-text transcription, auto-tag issue category
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-11-111 — Track support metrics: average response time, resolution time, customer satisfaction score

- **SRS Trace:** FG-011 / FR-111
- **Priority:** M2
- **Business acceptance (summary):**
  - Real-time dashboard, weekly reports to management
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-11-112 — Support escalation workflow: Tier 1 (Support) → Tier 2 (Technical) → Tier 3 (Engineering)

- **SRS Trace:** FG-011 / FR-112
- **Priority:** M2
- **Business acceptance (summary):**
  - Auto-escalation after 24hrs unresolved, notification sent
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-11-113 — Provide customer feedback form after ticket resolution with 5-star rating

- **SRS Trace:** FG-011 / FR-113
- **Priority:** M2
- **Business acceptance (summary):**
  - Feedback collected, low ratings flagged for review
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

[[[PAGEBREAK]]]

## 6.12 4.12 Notifications & Communication (FG-012)

Business intent: define the outcomes this capability must deliver for customers/partners/admin teams.

### BR-12-114 — Implement Kafka event-driven notification system with multiple channels: in-app push, SMS, email

- **SRS Trace:** FG-012 / FR-114
- **Priority:** M1
- **Business acceptance (summary):**
  - Event published within 100ms, delivery to all channels coordinated
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-12-115 — Send notifications for: OTP, verification, purchase confirmation, claims updates, renewal reminders, payment confirmations

- **SRS Trace:** FG-012 / FR-115
- **Priority:** M1
- **Business acceptance (summary):**
  - Template-based messages, personalized with customer data
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-12-116 — Support notification preferences with opt-in/opt-out for marketing and promotional messages

- **SRS Trace:** FG-012 / FR-116
- **Priority:** M2
- **Business acceptance (summary):**
  - User preferences stored, GDPR-compliant consent management
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-12-117 — Implement customer mute mode with minimum text notification (avoiding push for low-end devices)

- **SRS Trace:** FG-012 / FR-117
- **Priority:** M2
- **Business acceptance (summary):**
  - Device capability detection, graceful degradation
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-12-118 — Allow partners to create secondary marketing notifications filtered by: age, gender, location, policy type

- **SRS Trace:** FG-012 / FR-118
- **Priority:** M3
- **Business acceptance (summary):**
  - D
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-12-119 — Track notification delivery status: queued, sent, delivered, failed, bounced with retry mechanism

- **SRS Trace:** FG-012 / FR-119
- **Priority:** M2
- **Business acceptance (summary):**
  - Real-time status tracking, max 3 retries with exponential backoff
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-12-120 — Support message templates with dynamic placeholders for personalization

- **SRS Trace:** FG-012 / FR-120
- **Priority:** M2
- **Business acceptance (summary):**
  - Template engine with Bengali/English support, variable substitution
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-12-121 — Implement rate limiting for notifications to prevent spam (max 5 per hour per user)

- **SRS Trace:** FG-012 / FR-121
- **Priority:** M3
- **Business acceptance (summary):**
  - Redis-based rate limiting, exception for critical alerts
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-12-122 — Provide notification history in customer dashboard with read/unread status

- **SRS Trace:** FG-012 / FR-122
- **Priority:** M3
- **Business acceptance (summary):**
  - Last 90 days visible, older notifications archived
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-12-123 — Support rich push notifications with images, action buttons, and deep links

- **SRS Trace:** FG-012 / FR-123
- **Priority:** D
- **Business acceptance (summary):**
  - Platform-specific implementation (iOS/Android), click tracking
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

[[[PAGEBREAK]]]

## 6.13 4.13 IoT Integration & Usage-Based Insurance (FG-013)

Business intent: define the outcomes this capability must deliver for customers/partners/admin teams.

### BR-13-124 — Support IoT device integration for Usage-Based Insurance (UBI) via proprietary protocol

- **SRS Trace:** FG-013 / FR-124
- **Priority:** F
- **Business acceptance (summary):**
  - MQTT/CoAP protocol support, device authentication, encrypted communication
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-13-125 — Collect and process IoT data: location, speed, temperature, health vitals based on insurance type

- **SRS Trace:** FG-013 / FR-125
- **Priority:** D
- **Business acceptance (summary):**
  - Real-time data ingestion, time-series database storage
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-13-126 — Implement risk scoring based on IoT data patterns for dynamic premium adjustment

- **SRS Trace:** FG-013 / FR-126
- **Priority:** F
- **Business acceptance (summary):**
  - ML-based risk model, monthly recalculation, customer notification
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-13-127 — Provide customer dashboard showing IoT insights and risk score with improvement tips

- **SRS Trace:** FG-013 / FR-127
- **Priority:** F
- **Business acceptance (summary):**
  - Visualization with charts, gamification elements, personalized recommendations
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-13-128 — Support telematics integration for motor insurance with driving behavior analysis

- **SRS Trace:** FG-013 / FR-128
- **Priority:** D
- **Business acceptance (summary):**
  - Acceleration, braking, speed monitoring, trip history, safety score
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-13-129 — Integrate with wearable devices for health insurance with fitness tracking

- **SRS Trace:** FG-013 / FR-129
- **Priority:** D
- **Business acceptance (summary):**
  - Steps, heart rate, sleep quality monitoring, wellness rewards program
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-13-130 — Implement data privacy controls allowing customers to pause/resume IoT data collection

- **SRS Trace:** FG-013 / FR-130
- **Priority:** F
- **Business acceptance (summary):**
  - One-click toggle, data deletion option, privacy dashboard
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-13-178 — Integrate with IoT devices: GPS trackers (vehicles), health wearables (fitness bands), smart home sensors (fire/water leak)

- **SRS Trace:** FG-013 / FR-178
- **Priority:** M3
- **Business acceptance (summary):**
  - MQTT/CoAP protocol support, device SDK documentation, API endpoints
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-13-179 — Support IoT device registration, provisioning, and lifecycle management with certificate-based authentication

- **SRS Trace:** FG-013 / FR-179
- **Priority:** M3
- **Business acceptance (summary):**
  - X.509 certificates, device onboarding workflow, status tracking (active/inactive/suspended)
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-13-180 — Process and store IoT telemetry data using MQTT broker with TimescaleDB for time-series storage

- **SRS Trace:** FG-013 / FR-180
- **Priority:** M3
- **Business acceptance (summary):**
  - Handle 10,000 devices, 1 msg/min/device average, data retention policy (90 days hot, 2 years warm)
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-13-181 — Generate real-time alerts based on IoT data thresholds: aggressive driving (>80km/h in city), health anomalies (heart rate), home incidents

- **SRS Trace:** FG-013 / FR-181
- **Priority:** M3
- **Business acceptance (summary):**
  - Rule engine for threshold monitoring, push notifications, SMS alerts, configurable rules
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-13-182 — Support Usage-Based Insurance (UBI) pricing calculation based on IoT data: driving score (speed, braking, time-of-day), step count, heart rate variability

- **SRS Trace:** FG-013 / FR-182
- **Priority:** M3
- **Business acceptance (summary):**
  - Dynamic premium adjustment algorithm, monthly recalculation, transparent scoring dashboard
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-13-183 — Provide IoT device management portal for partners to monitor connected devices, data streams, and device health

- **SRS Trace:** FG-013 / FR-183
- **Priority:** M3
- **Business acceptance (summary):**
  - Real-time device status, data visualization charts, anomaly detection, bulk operations
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-13-184 — Support batch and real-time IoT data processing with configurable collection frequencies (1min to 1hour intervals)

- **SRS Trace:** FG-013 / FR-184
- **Priority:** M3
- **Business acceptance (summary):**
  - Stream processing (Kafka Streams), batch jobs, data quality checks, deduplication
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-13-185 — Maintain IoT device inventory with status tracking (online/offline/maintenance/decommissioned) and metadata

- **SRS Trace:** FG-013 / FR-185
- **Priority:** M3
- **Business acceptance (summary):**
  - Device registry, heartbeat monitoring (5min timeout), auto-offline detection, firmware version tracking
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

[[[PAGEBREAK]]]

## 6.14 4.14 AI & Automation Features (FG-014)

Business intent: define the outcomes this capability must deliver for customers/partners/admin teams.

### BR-14-164 — Implement AI chatbot for customer assistance during product search, selection, purchase, and claims

- **SRS Trace:** FG-014 / FR-164
- **Priority:** F
- **Business acceptance (summary):**
  - Bengali NLP support, 80% query resolution, human handoff capability
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-14-165 — Implement LLM multi-agent network for intelligent document processing and validation

- **SRS Trace:** FG-014 / FR-165
- **Priority:** F
- **Business acceptance (summary):**
  - OCR integration, field extraction accuracy >90%, fraud detection
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-14-166 — Implement AI-powered fraud detection using pattern recognition and anomaly detection

- **SRS Trace:** FG-014 / FR-166
- **Priority:** D
- **Business acceptance (summary):**
  - ML model with continuous learning, risk scoring, false positive <10%
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-14-167 — Support predictive analytics for risk assessment and premium optimization

- **SRS Trace:** FG-014 / FR-167
- **Priority:** F
- **Business acceptance (summary):**
  - Historical data analysis, model retraining, A/B testing capability
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-14-168 — Implement voice-assisted workflow for Type 3 users (rural/low digital literacy)

- **SRS Trace:** FG-014 / FR-168
- **Priority:** F
- **Business acceptance (summary):**
  - Bengali speech recognition, step-by-step guidance, voice commands
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-14-169 — Provide AI-based document verification with face matching and NID validation

- **SRS Trace:** FG-014 / FR-169
- **Priority:** M3
- **Business acceptance (summary):**
  - Liveness detection, face match confidence >95%, automated approval flow
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

[[[PAGEBREAK]]]

## 6.15 4.15 Voice-Assisted Features (FG-015)

Business intent: define the outcomes this capability must deliver for customers/partners/admin teams.

### BR-15-170 — Support Bengali speech-to-text (STT) with 90%+ accuracy for standard dialects (Dhaka, Chittagong, Sylhet)

- **SRS Trace:** FG-015 / FR-170
- **Priority:** M2
- **Business acceptance (summary):**
  - ASR model integration (Google/AWS/local), <2s latency, multi-dialect support
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-15-171 — Provide voice-guided policy purchase workflow with step-by-step audio instructions in Bengali

- **SRS Trace:** FG-015 / FR-171
- **Priority:** M2
- **Business acceptance (summary):**
  - Complete policy purchase via voice, TTS integration, progress tracking
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-15-172 — Support voice-based claims submission with automated transcription and field validation

- **SRS Trace:** FG-015 / FR-172
- **Priority:** M3
- **Business acceptance (summary):**
  - Voice recording up to 5min, transcription accuracy >85%, auto-populate claim form
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-15-173 — Provide text-to-speech (TTS) for Bengali language with natural-sounding voice

- **SRS Trace:** FG-015 / FR-173
- **Priority:** M2
- **Business acceptance (summary):**
  - Natural prosody, <1s response time, caching for common phrases, offline fallback
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-15-174 — Support voice navigation throughout mobile app for accessibility (elderly/visually impaired users)

- **SRS Trace:** FG-015 / FR-174
- **Priority:** D
- **Business acceptance (summary):**
  - Voice commands for all major functions, screen reader compatibility
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-15-175 — Provide voice command taxonomy: "buy policy", "file claim", "check status", "pay premium", "call agent"

- **SRS Trace:** FG-015 / FR-175
- **Priority:** M2
- **Business acceptance (summary):**
  - Intent recognition with 85%+ accuracy, contextual understanding, error handling
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-15-176 — Support seamless fallback to human agent when voice recognition confidence is below 80%

- **SRS Trace:** FG-015 / FR-176
- **Priority:** M3
- **Business acceptance (summary):**
  - Confidence scoring, automatic handoff with context transfer, queue management
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-15-177 — Log and analyze voice interactions for continuous improvement with user consent

- **SRS Trace:** FG-015 / FR-177
- **Priority:** D
- **Business acceptance (summary):**
  - Voice data collection opt-in, anonymization, model retraining pipeline, performance metrics
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

[[[PAGEBREAK]]]

## 6.16 4.16 Fraud Detection & Risk Controls (FG-016)

Business intent: define the outcomes this capability must deliver for customers/partners/admin teams.

### BR-16-186 — Flag claims submitted within 48hrs of policy purchase for manual review

- **SRS Trace:** FG-016 / FR-186
- **Priority:** M2
- **Business acceptance (summary):**
  - Auto-flagging with notification to Claims Officer, review queue
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-16-187 — Detect same claim type >2 times in 12 months and flag for pattern analysis

- **SRS Trace:** FG-016 / FR-187
- **Priority:** M2
- **Business acceptance (summary):**
  - Historical claim analysis, risk scoring, enhanced verification
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-16-188 — Flag claims where amount exactly matches policy limit (100% of coverage)

- **SRS Trace:** FG-016 / FR-188
- **Priority:** M2
- **Business acceptance (summary):**
  - Suspicious pattern detection, additional document requirements
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-16-189 — Validate medical provider against approved network list and flag non-network claims

- **SRS Trace:** FG-016 / FR-189
- **Priority:** M2
- **Business acceptance (summary):**
  - Provider database, real-time validation, approval workflow
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-16-190 — Implement device fingerprinting to detect multiple accounts from same device (>3 accounts)

- **SRS Trace:** FG-016 / FR-190
- **Priority:** M3
- **Business acceptance (summary):**
  - Browser/mobile device ID tracking, IP analysis, account linking
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-16-191 — Provide fraud detection dashboard for Business Admin and Focal Person with drill-down capability

- **SRS Trace:** FG-016 / FR-191
- **Priority:** M2
- **Business acceptance (summary):**
  - Real-time alerts, risk score visualization, action buttons
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-16-192 — Implement RACI for monitoring and incident escalation per defined roles

- **SRS Trace:** FG-016 / FR-192
- **Priority:** M1
- **Business acceptance (summary):**
  - Responsibility matrix enforced, escalation triggers, notification system
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

[[[PAGEBREAK]]]

## 6.17 4.17 Admin & Reporting (FG-017)

Business intent: define the outcomes this capability must deliver for customers/partners/admin teams.

### BR-17-131 — Provide role-based admin dashboards for: System Admin, Business Admin, Focal Person, Database Admin, Repository Admin

- **SRS Trace:** FG-017 / FR-131
- **Priority:** M1
- **Business acceptance (summary):**
  - Dynamic content based on role, real-time data updates
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-17-132 — Enforce strict 2FA for all admin-level access with TOTP authentication

- **SRS Trace:** FG-017 / FR-132
- **Priority:** M1
- **Business acceptance (summary):**
  - Google Authenticator/Authy compatible, backup codes provided
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-17-133 — Provide user management module: create, update, suspend, delete users with audit trail

- **SRS Trace:** FG-017 / FR-133
- **Priority:** M2
- **Business acceptance (summary):**
  - Full CRUD operations, role assignment, activity logs
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-17-134 — Provide product management module: create, update, activate/deactivate insurance products

- **SRS Trace:** FG-017 / FR-134
- **Priority:** M1
- **Business acceptance (summary):**
  - Version control, effective date management, pricing configuration
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-17-135 — Provide claims management dashboard with filtering: status, amount range, date, partner

- **SRS Trace:** FG-017 / FR-135
- **Priority:** M2
- **Business acceptance (summary):**
  - Advanced search, bulk actions, export functionality
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-17-136 — Provide task management system with assignment to internal users and deadline tracking

- **SRS Trace:** FG-017 / FR-136
- **Priority:** D
- **Business acceptance (summary):**
  - Task creation, assignment, status updates, notification on overdue
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-17-137 — Generate standard reports: daily sales, claims ratio, partner performance, policy counts, revenue

- **SRS Trace:** FG-017 / FR-137
- **Priority:** M2
- **Business acceptance (summary):**
  - M
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-17-138 — Provide custom report builder with drag-drop interface for business users

- **SRS Trace:** FG-017 / FR-138
- **Priority:** D
- **Business acceptance (summary):**
  - Visual query builder, chart generation, saved report templates
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-17-139 — Track KPIs aligned to business plan: policy acquisition rate, claim settlement ratio, customer retention

- **SRS Trace:** FG-017 / FR-139
- **Priority:** M3
- **Business acceptance (summary):**
  - M
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-17-140 — Provide system health monitoring dashboard: server status, API response times, error rates

- **SRS Trace:** FG-017 / FR-140
- **Priority:** M2
- **Business acceptance (summary):**
  - Integration with Prometheus/Grafana, alert configuration
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

[[[PAGEBREAK]]]

## 6.18 4.18 Analytics & Reporting (FG-018)

Business intent: define the outcomes this capability must deliver for customers/partners/admin teams.

### BR-18-149 — Track user behavior analytics: page views, feature usage, drop-off points, conversion funnel

- **SRS Trace:** FG-018 / FR-149
- **Priority:** D
- **Business acceptance (summary):**
  - Integration with analytics platform (Google Analytics/Mixpanel)
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-18-150 — Provide predictive analytics for customer churn, claim likelihood, policy renewal probability

- **SRS Trace:** FG-018 / FR-150
- **Priority:** F
- **Business acceptance (summary):**
  - F
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-18-151 — Generate customer segmentation reports: demographics, policy type, risk profile, lifetime value

- **SRS Trace:** FG-018 / FR-151
- **Priority:** D
- **Business acceptance (summary):**
  - Automated segmentation, export for marketing campaigns
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-18-152 — Provide geographic analytics: policy distribution by district, claims heatmap, agent performance by region

- **SRS Trace:** FG-018 / FR-152
- **Priority:** D
- **Business acceptance (summary):**
  - Map visualization, district-level drill-down, comparative analysis
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-18-202 — Provide geospatial risk visualization overlaying claims data on regional maps for heatmap analysis

- **SRS Trace:** FG-018 / FR-202
- **Priority:** D
- **Business acceptance (summary):**
  - Mapbox/Google Maps integration, district-level aggregation, color-coded risk zones
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-18-203 — Provide pre-built dashboards: Executive, Operations, Compliance with drill-down

- **SRS Trace:** FG-018 / FR-203
- **Priority:** D
- **Business acceptance (summary):**
  - Interactive charts, export capability, scheduled email delivery
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-18-204 — Track compliance metrics: AML flags, IDRA report status, audit logs access

- **SRS Trace:** FG-018 / FR-204
- **Priority:** M2
- **Business acceptance (summary):**
  - Real-time compliance dashboard, alerts on violations
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

[[[PAGEBREAK]]]

## 6.19 4.19 Audit & Logging (FG-019)

Business intent: define the outcomes this capability must deliver for customers/partners/admin teams.

### BR-19-153 — Maintain immutable audit logs for critical actions: policy issue, claim approval, payment, dispute resolution

- **SRS Trace:** FG-019 / FR-153
- **Priority:** M1
- **Business acceptance (summary):**
  - PostgreSQL with append-only tables, tamper detection
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-19-154 — Implement data retention policy with 20-year minimum for regulatory compliance

- **SRS Trace:** FG-019 / FR-154
- **Priority:** M2
- **Business acceptance (summary):**
  - Tiered storage (hot/warm/cold), automated archival, retrieval SLA
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-19-155 — Track all logged-in user actions with IP address, device info, timestamp, action type

- **SRS Trace:** FG-019 / FR-155
- **Priority:** M3
- **Business acceptance (summary):**
  - Comprehensive logging, queryable audit trail, GDPR compliance
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-19-156 — Allow partners to maintain additional logs as per MOU agreement with InsureTech

- **SRS Trace:** FG-019 / FR-156
- **Priority:** F
- **Business acceptance (summary):**
  - Partner-specific log tables, data isolation, access controls
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-19-157 — Provide regulatory portal for IDRA/BFIU to access requested data as per law

- **SRS Trace:** FG-019 / FR-157
- **Priority:** M2
- **Business acceptance (summary):**
  - Secure portal, report generation, audit trail of data access
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-19-158 — Implement log aggregation and analysis with alerting on suspicious patterns

- **SRS Trace:** FG-019 / FR-158
- **Priority:** M2
- **Business acceptance (summary):**
  - ELK stack/CloudWatch integration, anomaly detection, real-time alerts
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

[[[PAGEBREAK]]]

## 6.20 4.20 System Interface Architecture (FG-020)

Business intent: define the outcomes this capability must deliver for customers/partners/admin teams.

### BR-20-193 — Implement High-Performance Internal API for gateway-microservices communication with low latency guarantees

- **SRS Trace:** FG-020 / FR-193
- **Priority:** M1
- **Business acceptance (summary):**
  - <100ms response time, circuit breaker pattern, retry logic
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-20-194 — Implement Client-Optimized API for gateway-customer device communication with efficient data fetching

- **SRS Trace:** FG-020 / FR-194
- **Priority:** M1
- **Business acceptance (summary):**
  - <2s response time, query optimization, field-level authorization
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-20-195 — Implement Standard Integration API for 3rd party partners with comprehensive documentation

- **SRS Trace:** FG-020 / FR-195
- **Priority:** D
- **Business acceptance (summary):**
  - <200ms response time, standardized docs, sandbox environment
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-20-196 — Provide public Public Discovery API for product search and listing with rate limiting

- **SRS Trace:** FG-020 / FR-196
- **Priority:** M1
- **Business acceptance (summary):**
  - <1s response time, request limiting, caching enabled
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-20-197 — Expose only Cloudflare proxy and NGINX entry node to public, blocking direct microservice access

- **SRS Trace:** FG-020 / FR-197
- **Priority:** M1
- **Business acceptance (summary):**
  - Firewall rules configured, internal IPs hidden, DDoS protection
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-20-198 — Implement Real-Time Connection capability for instant updates (notifications, claims status)

- **SRS Trace:** FG-020 / FR-198
- **Priority:** D
- **Business acceptance (summary):**
  - Persistent connection management, automatic reconnection, heartbeat
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-20-199 — Use Efficient Binary Protocol for IoT data extraction and data binding

- **SRS Trace:** FG-020 / FR-199
- **Priority:** F
- **Business acceptance (summary):**
  - Custom binary formatting, data compression, low latency
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-20-200 — Consolidate, annotate and process data for AI agent training within regulatory limits

- **SRS Trace:** FG-020 / FR-200
- **Priority:** F
- **Business acceptance (summary):**
  - Data anonymization, consent management, audit trail
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-20-201 — Generate statistics and predictions based on big data for partner insights

- **SRS Trace:** FG-020 / FR-201
- **Priority:** F
- **Business acceptance (summary):**
  - ML pipeline, data lake architecture, API for insights delivery
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-20-159 — Implement Blockchain-based shared ledger for automated reinsurance settlements and smart contract execution

- **SRS Trace:** FG-020 / FR-159
- **Priority:** D
- **Business acceptance (summary):**
  - Immutable ledger, transparency audit trail
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-20-160 — Implement AI-driven dynamic premium discounting based on real-time risk assessment and loyalty scoring

- **SRS Trace:** FG-020 / FR-160
- **Priority:** D
- **Business acceptance (summary):**
  - Risk model integration, real-time calculation, customer notification
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-20-161 — Integrate with SMS Gateway for OTP and notifications

- **SRS Trace:** FG-020 / FR-161
- **Priority:** M1
- **Business acceptance (summary):**
  - Delivery rate >95%, delivery status tracking, cost optimization
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-20-162 — Integrate with Email Service for transactional and marketing emails

- **SRS Trace:** FG-020 / FR-162
- **Priority:** M1
- **Business acceptance (summary):**
  - Template management, bounce handling, unsubscribe management
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-20-163 — Provide Webhook System for real-time event notifications to external systems

- **SRS Trace:** FG-020 / FR-163
- **Priority:** M2
- **Business acceptance (summary):**
  - Event filtering, retry mechanism, authentication, payload signing
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-20-223 — Provide API contract specification: All Category 3 APIs must provide OpenAPI 3.0 spec with request/response schemas, error codes, example payloads

- **SRS Trace:** FG-020 / FR-223
- **Priority:** M3
- **Business acceptance (summary):**
  - • OpenAPI spec complete / • Error codes documented / • Examples provided
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-20-224 — Define insurer API payloads: Premium Calculation API, Policy Issuance API with standardized request/response formats

- **SRS Trace:** FG-020 / FR-224
- **Priority:** M1
- **Business acceptance (summary):**
  - • Payload formats defined / • Validation rules clear / • Sample payloads provided
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-20-225 — Define payment gateway payloads: Initiate Payment, Webhook Callback with HMAC-SHA256 signature validation

- **SRS Trace:** FG-020 / FR-225
- **Priority:** M1
- **Business acceptance (summary):**
  - • Payment payloads defined / • Signature validation implemented / • Security tested
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-20-226 — Implement retry logic: Failed API calls retry with exponential backoff: 1s, 2s, 4s, 8s, 16s (max 5 retries); Use circuit breaker pattern

- **SRS Trace:** FG-020 / FR-226
- **Priority:** M1
- **Business acceptance (summary):**
  - • Retry logic tested / • Exponential backoff works / • Circuit breaker functional
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-20-227 — Implement idempotency: All payment and policy issuance APIs must accept Idempotency-Key header (UUID); Store keys for 24 hours; Return cached response for duplicates

- **SRS Trace:** FG-020 / FR-227
- **Priority:** M1
- **Business acceptance (summary):**
  - • Idempotency enforced / • Key storage works / • Duplicate handling correct
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-20-228 — Implement callback security: Payment gateway webhooks must include HMAC-SHA256 signature in header; Validate signature; Reject unsigned/invalid callbacks; Log all attempts

- **SRS Trace:** FG-020 / FR-228
- **Priority:** M2
- **Business acceptance (summary):**
  - • Signature validation works / • Invalid callbacks rejected / • Logging comprehensive
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-20-229 — Support EHR integration approach - Option A (Preferred): Use LabAid FHIR API with Patient resource matching by NID/phone; Query Encounter resources; Pre-authorization workflow

- **SRS Trace:** FG-020 / FR-229
- **Priority:** S
- **Business acceptance (summary):**
  - • FHIR API integrated / • Patient matching accurate / • Pre-auth workflow functional
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-20-230 — Support EHR integration approach - Option B (Fallback): Use LabAid custom REST API with endpoints for patient admissions, pre-auth verification, bills; Secure with mutual TLS + API key

- **SRS Trace:** FG-020 / FR-230
- **Priority:** D
- **Business acceptance (summary):**
  - • Custom API integrated / • mTLS configured / • API key management
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-20-231 — Handle EHR integration timeout: Set connection timeout 5s, read timeout 15s; If timeout, queue for manual verification; Notify hospital staff via SMS

- **SRS Trace:** FG-020 / FR-231
- **Priority:** D
- **Business acceptance (summary):**
  - • Timeout handling works / • Manual queue functional / • Notifications sent
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

[[[PAGEBREAK]]]

## 6.21 4.22 Data Storage (FG-022)

Business intent: define the outcomes this capability must deliver for customers/partners/admin teams.

### BR-21-232 — Use PostgreSQL V17 for structured data with JSON support and full-text search capability

- **SRS Trace:** FG-022 / FR-232
- **Priority:** M1
- **Business acceptance (summary):**
  - Primary database setup, performance optimization, localization
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-21-233 — Implement read replicas for reporting and analytics workloads

- **SRS Trace:** FG-022 / FR-233
- **Priority:** M3
- **Business acceptance (summary):**
  - Read scaling, data consistency, performance monitoring
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-21-234 — Implement Graph Database (Neo4j/Amazon Neptune) for visualizing complex fraud relationships and entity resolution

- **SRS Trace:** FG-022 / FR-234
- **Priority:** D
- **Business acceptance (summary):**
  - Graph schema defined, node relationship mapping, query performance <1s
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-21-235 — Use Redis for session management and high-frequency real-time data

- **SRS Trace:** FG-022 / FR-235
- **Priority:** M3
- **Business acceptance (summary):**
  - Performance optimization, session management, cache strategies
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-21-236 — Implement data partitioning for policies and claims tables by month

- **SRS Trace:** FG-022 / FR-236
- **Priority:** M3
- **Business acceptance (summary):**
  - Scalability, query performance, maintenance efficiency
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-21-237 — Use S3-compatible Object Storage for document files with encryption at rest

- **SRS Trace:** FG-022 / FR-237
- **Priority:** M1
- **Business acceptance (summary):**
  - Secure document storage, lifecycle management, CDN integration
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-21-238 — Store product catalog and metadata in Document-Oriented NoSQL Database

- **SRS Trace:** FG-022 / FR-238
- **Priority:** M3
- **Business acceptance (summary):**
  - Flexible schema, high availability, global distribution
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-21-239 — Upload data policy - Client-side compression: 5MB → 1-2MB (JPEG 80% quality, 1920x1080 max resolution), Chunked upload: 1MB chunks with resume capability (tus.io protocol), Presigned S3 URLs: Direct upload, 30-minute expiry

- **SRS Trace:** FG-022 / FR-239
- **Priority:** M1
- **Business acceptance (summary):**
  - check upload >5MB fails,<5MB passes
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-21-240 — Backup: Daily full, 6-hour incremental, continuous transaction logs

- **SRS Trace:** FG-022 / FR-240
- **Priority:** M1
- **Business acceptance (summary):**
  - Check new backup after 6hour
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-21-241 — Store app native encrypted data in user device in SQLite

- **SRS Trace:** FG-022 / FR-241
- **Priority:** M2
- **Business acceptance (summary):**
  - Check sqlitefiles
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-21-242 — Process tokenized data on Vector Database for AI embeddings

- **SRS Trace:** FG-022 / FR-242
- **Priority:** D
- **Business acceptance (summary):**
  - Similarity search latency check
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-21-243 — Implement Columnar Database (ClickHouse/Druid) for high-performance real-time analytics and reporting

- **SRS Trace:** FG-022 / FR-243
- **Priority:** D
- **Business acceptance (summary):**
  - OLAP query performance <500ms, data compression, scalability
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

[[[PAGEBREAK]]]

## 6.22 4.23 User Interface Requirements (FG-023)

Business intent: define the outcomes this capability must deliver for customers/partners/admin teams.

### BR-22-244 — Maintain consistent UI across Android and iOS using React Native

- **SRS Trace:** FG-023 / FR-244
- **Priority:** M1
- **Business acceptance (summary):**
  - Shared codebase >90%
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-22-245 — Provide smart data widgets for mobile users

- **SRS Trace:** FG-023 / FR-245
- **Priority:** D
- **Business acceptance (summary):**
  - Customizable dashboard
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-22-246 — Support desktop-first responsive design for portals

- **SRS Trace:** FG-023 / FR-246
- **Priority:** M1
- **Business acceptance (summary):**
  - 1024px minimum width
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-22-247 — Request minimum device permissions

- **SRS Trace:** FG-023 / FR-247
- **Priority:** M1
- **Business acceptance (summary):**
  - Camera, SMS read only
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

### BR-22-248 — Support Bengali and English with toggle

- **SRS Trace:** FG-023 / FR-248
- **Priority:** M1
- **Business acceptance (summary):**
  - i18n framework implemented
- **Primary portals impacted:** (to be confirmed during UX)
  - Customer App / Web, Partner Portal, Admin Portals as applicable

[[[PAGEBREAK]]]
