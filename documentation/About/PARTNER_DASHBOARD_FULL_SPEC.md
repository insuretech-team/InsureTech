# LabaidInsuretech Partner & Agent Dashboard Full Specification

**Document Version:** 1.0  
**Date:** 2026-03-18  
**Audience:** Product, Engineering, QA, Operations, Compliance, Partner Stakeholders  
**Status:** Draft for alignment and implementation  
**Derived From:** INSURER_DASHBOARD_FULL_SPEC.md, SRS_v3/SPECS_V3.7, proto definitions, partner_portal codebase, b2b_portal source

---

## 1. Purpose

This document defines the full specification for the **LabaidInsuretech Partner & Agent Dashboard**, a dedicated web portal for partner-side and agent-side operations within the broader LabaidInsuretech platform.

The Partner Portal serves two primary audiences:

1. **Service Partners** — hospitals, clinics, pharmacies, auto service centers, pet clinics, fire inspectors, laptop/mobile repair shops, ambulance providers — who deliver services to insured customers and interact with the platform for cashless billing, claims, and settlement
2. **Agent Networks** — insurance agents and sales representatives who sell policies, onboard customers, track commissions, and manage referrals under partner organizations

The dashboard is intended to give partners and agents a secure operational workspace to:

- receive onboarding, KYB verification, and activation
- manage partner profile, documents, and compliance status
- receive bulk-uploaded partner rosters (cashless/discount) by policy type from admin
- submit bills for cashless or reimbursement claims
- check and track claims submitted by customers or on behalf of customers
- initiate policy purchase on behalf of customers with consent
- manage referral pipelines
- track commissions at partner and agent levels
- manage agents under the partner organization
- view analytics, performance, and operational KPIs

This specification is aligned to the current platform direction documented across `documentation/About`, `documentation/BRD`, `documentation/core_plans`, and `documentation/SRS_v3`.

---

## 2. Strategic Positioning

The Partner Portal is not a standalone product. It is a partner-facing portal inside the **same LabaidInsuretech platform architecture** used by the B2B portal, customer portal, insurer dashboard, and system admin portal.

### 2.1 Partner Types and Policy-Type Mapping

The platform supports cashless and discount partners organized by insurance line of business:

| Partner Type | Policy Lines Served | Service Model | Priority |
|---|---|---|---|
| **Hospital** | Health (Group, Individual, OMP) | Cashless + Reimbursement + Pre-auth | M2 |
| **Clinic** | Health (Individual, Micro) | Cashless + Discount | M2 |
| **Pharmacy** | Health (All) | Discount + Reimbursement | M2 |
| **Auto Service Center** | Motor (Private, Commercial) | Cashless + Estimate Submission | M2 |
| **Pet Clinic / Veterinary** | Pet Insurance, Livestock | Cashless + Reimbursement | M3 |
| **Fire Inspector / Surveyor** | Fire, Property, P&C | Survey + Assessment Report | M3 |
| **Laptop/Mobile Repair** | Device Insurance, Gadget | Cashless + Estimate | M3 |
| **Ambulance Service** | Health (Emergency) | Cashless + Pre-auth | M3 |
| **MFS Provider** | All (Payment Channel) | Premium Collection + Disbursement | M2 |
| **E-commerce** | All (Embedded Insurance) | Widget + API Integration | M3 |
| **Agent Network** | All (Sales Channel) | Policy Sales + Renewals + Referrals | M2 |
| **Corporate** | Group Health, Group Life | B2B Enrollment + Census | M2 |

### 2.2 Revenue Model Alignment

Partners interact with four commission/revenue models defined in the platform:

- **COMMISSION:** Traditional percentage of premium (health, life partners)
- **FLAT_FEE:** Fixed subscription or service fee (P&C, device partners)
- **HYBRID:** Premium percentage + operational fees (auto partners)
- **REINSURANCE:** Risk share (not partner-facing, but visible in reporting)

### 2.3 Admin-Driven Partner Onboarding via Excel Upload

A critical workflow is the **admin bulk upload of cashless/discount partners by policy type**. The system admin or business admin uploads an Excel roster that:

- maps each partner to one or more policy types (e.g., hospital → Health, auto shop → Motor)
- specifies cashless vs. discount designation per policy type
- sets discount percentages or cashless limits per partner-policy combination
- includes KYB reference data (trade license, TIN, bank account)
- triggers onboarding invitations and verification workflows

This enables rapid partner-network expansion across all insurance lines.

---

## 3. Objectives

### 3.1 Business Objectives

- enable rapid partner onboarding across health, motor, fire, pet, and device insurance lines
- replace manual spreadsheet-based partner coordination with auditable digital workflows
- standardize bill submission, claims checking, and settlement tracking across partner types
- provide partner-level visibility into policy referrals, conversions, and commissions
- support agent hierarchy with territory management and commission splits

### 3.2 Operational Objectives

- establish a single partner workspace for onboarding, billing, claims, commission, and agent management
- support admin-driven bulk partner upload with policy-type mapping
- create structured interfaces for bill submission, document upload, and claims tracking
- maintain full auditability for partner and agent actions
- support partial live integration while allowing controlled fallback during rollout

### 3.3 Technical Objectives

- reuse the platform's existing auth/session/gateway/BFF patterns (same as B2B and insurer portals)
- align with proto-first domain models (partner.proto, commission_config.proto, commission_payout.proto)
- remain deployable within the current Docker + nginx production strategy
- deliver responsive desktop-first layout with usable mobile/tablet fallback for field agents

---

## 4. Scope

### 4.1 In Scope

- partner authentication and role-based portal access
- partner onboarding and KYB verification workflow
- partner profile management and document upload
- admin bulk Excel upload of cashless/discount partners by policy type
- partner bill submission for cashless and reimbursement claims
- claims tracking and status monitoring
- customer referral pipeline and policy initiation on behalf of customer
- commission tracking at partner and agent levels
- agent registration, management, and hierarchy
- partner performance analytics and KPIs
- settings and configuration management
- cross-organization collaboration with Labaid Insuretech focal persons

### 4.2 Out of Scope for This Spec

- insurer dashboard operations (covered in INSURER_DASHBOARD_FULL_SPEC.md)
- customer portal experience
- core actuarial engine logic
- payment gateway implementation details
- partner accounting/GL replacement
- full regulatory filing engine implementation

---

## 5. Primary Users and Roles

### 5.1 Partner Portal Roles

| Role | Purpose | Typical Permissions |
|------|---------|---------------------|
| Partner Admin | Owns partner organization configuration and oversight | full partner workspace access, agent management, profile editing, reporting, settings |
| Partner Finance | Manages billing, invoicing, and commission reconciliation | bill submission, commission view, payment history, financial reporting |
| Partner Claims Officer | Handles claim submissions and tracking for partner-initiated claims | claim submission, status tracking, document upload, claims reporting |
| Insurance Agent | Sells policies, onboards customers, tracks referrals and commissions | customer referral, policy initiation, commission view, performance dashboard |
| Regional Manager | Oversees agents within a territory under the partner | agent performance monitoring, territory reporting, escalation |

### 5.2 Internal Counterpart Roles

| Role | Interaction with Partner Portal |
|------|------|
| Focal Person | Partner verification, approval, dispute resolution, performance monitoring |
| Business Admin | Bulk partner upload, commission configuration, partner suspension/activation |
| System Admin | Platform configuration, partner type and policy-type mapping |
| Claims Officer | Claim review, settlement, partner-submitted bill verification |
| Support Team | Partner query resolution, escalation handling |

### 5.3 Proto-Defined Role Hierarchy

As per SRS FR-014 and FR-019, the partner portal operates within this hierarchy:

```
Business Admin
  └── Focal Person (Partner Bridge) ★ KEY ROLE
        └── Partner Admin (Tenant Root)
              ├── Regional Manager
              ├── Insurance Agent
              └── Partner Staff (Finance, Claims)
```

---

## 6. Product Principles

1. **Partner-first operations:** the portal must reflect real partner/agent work — billing, claims, referrals, commissions — not generic admin tables
2. **Policy-type aware:** every partner view must be contextual to the partner's assigned policy types (health, motor, fire, pet, device, etc.)
3. **Bulk-upload friendly:** admin-driven Excel upload must be a first-class workflow, not an afterthought
4. **Document-heavy by design:** partner workflows depend on trade licenses, TIN certificates, MOU uploads, bills, claim evidence, and verification documents
5. **State-driven workflow:** every onboarding step, bill submission, and claim must map to explicit statuses and transitions
6. **Commission transparent:** partners and agents must have clear, auditable visibility into commission calculation, accrual, and payout
7. **Collaborative but auditable:** chat, calls, document requests, and notes tied to entities with timestamps and actor identity
8. **Fallback-ready:** the UI should remain operational while live endpoints are phased in

---

## 7. Architectural Fit

### 7.1 Platform Alignment

The partner dashboard shall be implemented as a **React/Next.js App Router portal** using the same portal/BFF pattern:

```
Browser Client → Next.js Portal Frontend → Next.js Server Routes/BFF Layer → Gateway/API Layer → Downstream Microservices (REST/gRPC)
```

### 7.2 Authentication and Session Model

The dashboard shall align with existing portal session conventions:

**Cookies:**
- `session_token` (HttpOnly) — session identifier
- `csrf_token` (HttpOnly) — CSRF protection
- `portal_role` — user role (PARTNER_ADMIN, AGENT, etc.)
- `portal_user_id` — user UUID
- `portal_biz_id` — partner organization UUID
- `portal_email` — contact email
- `portal_mobile` — contact mobile

**BFF Headers for Downstream Calls:**
- `x-portal: PORTAL_PARTNER`
- `x-user-id`
- `x-business-id` (partner_id)
- `x-tenant-id`

### 7.3 Service Dependencies

The partner dashboard integrates with:

| Service | Port (gRPC/HTTP) | Purpose |
|---------|---------|---------|
| Auth Service (Go) | 50060/50061 | Authentication, session, OTP |
| Authorization (Go) | 50070/50071 | RBAC/ABAC, permissions |
| Partner Service (Go) | 50100/50101 | Partner CRUD, agent management |
| Insurance Engine (C#) | 50120–50171 | Product catalog, quotes, orders, policies |
| Commission (C#) | 50150/50151 | Commission calculation, payout tracking |
| Claim (C#) | 50210/50211 | Claim submission, approval, settlement |
| Payment (Node.js) | 50190/50191 | Payment processing, settlement |
| Ledger (Go) | 50200/50201 | Double-entry financial records |
| Notification (Go) | 50230/50231 | SMS, email, push notifications |
| Storage (Go) | 50290/50291 | S3/blob document storage |
| Document Gen (Go) | 50280/50281 | PDF generation |
| Workflow (Go) | 50180/50181 | Process orchestration |
| Fraud (Go) | 50220/50221 | Fraud detection scoring |

### 7.4 Event and Saga Alignment

The portal should be compatible with the platform's event-driven architecture:

**Published Events:**
- `insuretech.partner.onboarded.v1` — partner activation
- `insuretech.partner.agent_registered.v1` — new agent under partner
- `insuretech.partner.bill_submitted.v1` — cashless/reimbursement bill
- `insuretech.partner.referral.created.v1` — customer referral

**Consumed Events:**
- `insuretech.claim.submitted.v1` — claim filed by customer
- `insuretech.claim.status_changed.v1` — claim status update
- `insuretech.commission.calculated.v1` — commission accrual notification
- `insuretech.commission.payout.processed.v1` — payout completion
- `insuretech.policy.issued.v1` — policy activation from referral

### 7.5 Proto Data Model Alignment

The portal's data operations align directly with the proto-defined entities:

**Partner Domain (`proto/insuretech/partner/entity/v1/partner.proto`):**
- `Partner` — organization entity with type, status, trade_license, TIN, bank_account, commission rates, benefits (cashless/discount config)
- `Agent` — agent entity with partner_id FK, NID, commission_rate, status
- `PartnerType` enum — HOSPITAL, PHARMACY, DOCTOR, AMBULANCE, AUTO_REPAIR, LAPTOP_REPAIR, MOBILE_REPAIR, MFS, ECOMMERCE, AGENT_NETWORK, CORPORATE
- `PartnerStatus` enum — PENDING_VERIFICATION, ACTIVE, SUSPENDED, TERMINATED
- `InsuranceAgentStatus` enum — ACTIVE, INACTIVE, SUSPENDED
- `PartnerBenefits` — discount_enabled, discount_percentage, cashless_enabled, cashless_limit, auto_approval_threshold, required_documents, service_locations
- `CommissionStructure` — acquisition_rate, renewal_rate, claims_assistance_rate

**Commission Domain (`proto/insuretech/commission/entity/v1/`):**
- `CommissionConfig` — per-product commission rules with revenue_model, acquisition_rate, renewal_rate, agent_split_config, performance_tiers
- `CommissionPayout` — batch payout with recipient_type (PARTNER/AGENT), period, amount, status, payment_method
- `RevenueShare` — platform/insurer revenue split per policy

---

## 8. Information Architecture

The partner dashboard primary navigation:

1. **Dashboard** — KPI overview and activity feed
2. **Partners** — partner roster, policy-type mapping, cashless/discount status (admin-uploaded)
3. **Agents** — agent management, hierarchy, territory assignment
4. **Claims** — bill submission, claims tracking, document management
5. **Policies** — referred policies, active coverage, renewal tracking
6. **Commission** — earnings, payouts, commission statements
7. **Referrals** — customer referral pipeline and conversion tracking
8. **Documents** — KYB documents, verification status, MOU management
9. **Reports** — performance analytics, financial summaries
10. **Settings** — profile, notifications, preferences, security

---

## 9. Detailed Module Specification

### 9.1 Dashboard Home

#### Purpose

Provide an executive and operational summary for the partner workspace.

#### Core Widgets

- **Active Policies Summary** — total referred/serviced policies by product line
- **Claims Pipeline** — claims by status (submitted, under review, approved, settled)
- **Commission Summary** — current month earned, pending payout, last payout amount
- **Agent Performance** — top agents by policy count and commission earned
- **Referral Pipeline** — leads, conversions, conversion rate
- **Pending Actions** — onboarding steps, document requests, bills awaiting response
- **Recent Activity Feed** — last 20 actions across claims, policies, commissions
- **Announcements** — platform notices, product updates, partner communications

#### KPI Cards

| KPI | Description |
|-----|-------------|
| Total Active Policies | Policies sold/serviced by this partner |
| Claims Submitted (MTD) | Claims submitted month to date |
| Claims Settled (MTD) | Claims settled with amount |
| Average Claim TAT | Average turnaround time for partner claims |
| Commission Earned (MTD) | Commissions earned current month |
| Commission Pending Payout | Approved but not yet disbursed |
| Active Agents | Count of active agents under partner |
| Referral Conversion Rate | Leads converted to policies (%) |

#### Required Behaviors

- default partner context loaded from session (partner_id, partner_type)
- all widgets filterable by date range and product line
- KPI cards support drill-down into relevant work queue
- tiles support empty states and fallback data mode
- dashboard loads in under 3 seconds

---

### 9.2 Partners Module (Admin-Uploaded Roster)

#### Purpose

Display the partner network roster, manage policy-type assignments, and support admin-driven bulk Excel upload of cashless/discount partners.

#### Admin Excel Upload Workflow

This is a critical first-class workflow. Business Admins upload Excel files to onboard partner networks:

**Upload process:**
1. Admin navigates to Partner Management → Bulk Upload
2. Downloads Excel template with required columns
3. Fills roster: organization name, type, trade license, TIN, contact, bank details, policy types served, cashless/discount designation, limits
4. Uploads completed Excel file
5. System validates:
   - required fields present
   - trade license uniqueness
   - TIN format (12-digit)
   - Bangladesh mobile format (+880 1XXXXXXXXX)
   - email format
   - valid partner type enum
   - valid policy-type mapping
6. Success/failure report generated per row
7. Valid partners created with status `PENDING_VERIFICATION`
8. Onboarding invitations dispatched via SMS/email
9. Focal Person notified for KYB verification queue

**Excel Template Columns:**

| Column | Type | Required | Validation |
|--------|------|----------|------------|
| Organization Name | text | yes | max 255 chars |
| Partner Type | enum | yes | HOSPITAL, CLINIC, PHARMACY, AUTO_REPAIR, PET_CLINIC, FIRE_INSPECTOR, LAPTOP_REPAIR, MOBILE_REPAIR, AMBULANCE, MFS, ECOMMERCE, AGENT_NETWORK, CORPORATE |
| Trade License Number | text | yes | unique across system |
| TIN Number | text | yes | 12-digit numeric |
| Contact Email | email | yes | valid email format |
| Contact Phone | text | yes | +880 format |
| Bank Account Number | text | yes | AES-256-GCM encrypted at rest |
| Bank Name | text | yes | |
| Bank Branch | text | no | |
| Policy Types Served | comma-separated | yes | e.g., "Health,Motor" |
| Cashless Enabled | boolean | yes | TRUE/FALSE per policy type |
| Discount Enabled | boolean | conditional | TRUE/FALSE |
| Discount Percentage | decimal | conditional | 0–100, required if discount enabled |
| Cashless Limit (BDT) | integer | conditional | required if cashless enabled |
| Auto-Approval Threshold (BDT) | integer | no | claims below this auto-approved |
| Pre-Authorization Required | boolean | no | default TRUE for hospitals |
| Service Locations | comma-separated | no | cities/districts covered |
| Nationwide Coverage | boolean | no | default FALSE |
| Acquisition Commission Rate (%) | decimal | yes | 0–100 |
| Renewal Commission Rate (%) | decimal | yes | 0–100 |
| Claims Assistance Rate (%) | decimal | no | 0–100 |

#### Partner List View

- filterable table by partner type, policy type, status, location
- search by organization name, trade license, TIN
- status badge: Pending Verification, Active, Suspended, Terminated
- inline policy-type tags showing assigned insurance lines
- cashless/discount indicator per policy type
- quick actions: view detail, suspend, contact

#### Partner Detail View

- organization information and contact details
- KYB document status (trade license, TIN, bank verification, MOU)
- assigned policy types with cashless/discount configuration per type
- commission structure (acquisition, renewal, claims assistance rates)
- benefits configuration (PartnerBenefits from proto)
- agent count and hierarchy summary
- performance metrics (policies, claims, conversion rate)
- audit trail of status changes
- focal person assignment

#### Partner Type-Specific Configuration

Each partner type has specific data requirements:

**Hospital/Clinic:**
- facility registration number
- bed count (hospitals)
- specialty departments
- operating hours
- emergency capability flag
- DGHS registration (where applicable)
- pre-authorization workflow configuration

**Pharmacy:**
- drug license number
- DGDA registration
- approved formulary compliance flag

**Auto Service Center:**
- workshop registration
- BRTA certification
- service categories (bodywork, mechanical, electrical, paint)
- estimate submission workflow

**Pet Clinic / Veterinary:**
- veterinary license
- BVA registration
- species coverage (dogs, cats, livestock, poultry)

**Fire Inspector / Surveyor:**
- surveyor license number
- IDRA empanelment status
- survey types supported (fire, property, marine, engineering)

---

### 9.3 Agents Module

#### Purpose

Manage insurance agents under the partner organization with hierarchy, territory, and performance tracking.

#### Agent Registration

Agent registration flow per SRS FR-060 and proto `Agent` entity:

1. Partner Admin initiates agent registration
2. Required fields collected:
   - Full name
   - Phone number (+880 format, unique)
   - Email (unique, optional)
   - NID number (10, 13, or 17 digits, unique)
   - Commission rate (0–100%)
3. System creates user account via Auth Service
4. Agent status set to ACTIVE
5. Welcome SMS/email with portal access credentials
6. Agent appears in partner's agent roster

#### Agent Bulk Upload

Similar to partner bulk upload, Partner Admins can upload agents via Excel:

**Template columns:** Full Name, Phone, Email, NID, Commission Rate (%), Territory/Region

**Validation:** phone uniqueness, NID format and uniqueness, commission rate range

#### Agent List View

- filterable table by status, territory, commission tier
- search by name, phone, NID
- status badge: Active, Inactive, Suspended
- performance columns: policies sold, commission earned, referrals
- quick actions: view detail, suspend, contact

#### Agent Detail View

- personal information (name, phone, email, NID)
- partner assignment and hierarchy position
- commission structure and current rates
- territory/region assignment
- performance dashboard:
  - policies sold (daily, weekly, monthly, yearly)
  - commission earned and pending
  - referral count and conversion rate
  - claim assistance count
- activity timeline
- status change history

#### Agent Hierarchy (FR-067)

The portal supports multi-level agent hierarchy:

```
Partner Admin
  └── Regional Manager
        └── Insurance Agent
```

- hierarchical commission split visible at each level
- territory management with district/upazila assignment
- performance tracking rolled up by hierarchy level
- regional manager sees aggregated agent metrics

---

### 9.4 Claims Module

#### Purpose

Enable partners to submit bills for cashless/reimbursement claims, track claim status, and manage claim-related documentation.

#### Bill Submission Workflow

Partners interact with claims in two primary modes:

**Mode 1: Cashless Bill Submission (Hospital, Clinic, Auto Service, Pet Clinic)**

1. Customer presents policy at partner facility
2. Partner verifies policy via portal (policy search by number or customer NID/phone)
3. Partner initiates cashless claim:
   - selects policy
   - enters incident/treatment details
   - specifies claim type (consultation, surgery, accident repair, etc.)
   - enters itemized bill breakdown
   - uploads supporting documents (bills, prescriptions, reports, photos)
4. System validates:
   - policy active and within coverage period
   - claim type covered under policy
   - partner is listed as cashless partner for this policy type
   - no duplicate submission
   - cashless limit not exceeded
5. If bill is below auto-approval threshold → auto-approved for settlement
6. If above threshold → routed to claims approval workflow
7. Partner tracks claim through status lifecycle

**Mode 2: Reimbursement Bill Assistance**

1. Customer visits partner and pays out of pocket
2. Partner generates itemized bill with required documentation
3. Partner uploads bill and documents to assist customer's reimbursement claim
4. Customer files reimbursement claim referencing partner-uploaded bill
5. Partner tracks reimbursement status

#### Claims List View

- filterable by status, claim type, product line, date range, amount band
- status tabs: All, Submitted, Under Review, Approved, Rejected, Settled
- search by claim number, policy number, customer name/NID
- columns: Claim ID, Customer, Policy Type, Amount, Status, Submitted Date, TAT
- export to Excel/PDF

#### Claim Detail View

- claim summary: ID, status, submitted date, product line
- customer and policy information (name, policy number, coverage)
- incident/treatment details
- itemized bill breakdown:
  - for health: accommodation, consultant fees, medicines, surgical, diagnostic, ancillary
  - for motor: parts, labor, paint, bodywork, electrical, glass
  - for device: parts, labor, diagnostic
  - for pet: consultation, treatment, medication, surgery
- document gallery with checklist status
- co-payment and deductible calculation display
- approval timeline with decision notes
- settlement amount and payment status
- communication thread (chat with focal person / claims officer)

#### Product-Specific Claims Fields

**Health Claims (Hospital/Clinic/Pharmacy):**

| Field Group | Fields |
|---|---|
| Patient Details | Name, age, gender, policy member relation |
| Admission Details | Date of admission, date of discharge, ward type, room number |
| Diagnosis | Primary diagnosis (ICD-10), secondary diagnoses |
| Treatment | Procedures performed, treating doctor |
| Bill Breakdown | Accommodation, consultant, medicines, surgical, diagnostic, ambulance, other |
| Documents Required | Discharge summary, prescription, investigation reports, pharmacy bills, doctor certificate |

**Motor Claims (Auto Service Center):**

| Field Group | Fields |
|---|---|
| Vehicle Details | Registration, make, model, year |
| Incident Details | Date, location, description, police report number |
| Damage Assessment | Damage photos (min 4 angles), surveyor notes |
| Repair Estimate | Parts list with part number and cost, labor hours, paint/body, total estimate |
| Documents Required | 3 repair estimates from different workshops, MVI report, driver statement, FIR (if theft/accident), driving license copy |

**Fire/Property Claims (Fire Inspector):**

| Field Group | Fields |
|---|---|
| Property Details | Location, construction type, occupancy |
| Incident Details | Date, cause, fire brigade report number |
| Damage Assessment | Area affected, estimated loss, salvage value |
| Survey Report | Surveyor findings, photographs, estimated payable |
| Documents Required | Fire brigade report, FIR, stock verification (90 days), electrical compliance, municipal tax receipt, building plan |

**Pet/Livestock Claims (Vet):**

| Field Group | Fields |
|---|---|
| Animal Details | Species, breed, ear tag/microchip, age, weight |
| Incident Details | Date, condition/illness, treating vet, clinic |
| Treatment | Procedures, medication, hospitalization days |
| Documents Required | Vet certificate, treatment records, vaccination history, death certificate (if applicable) |

#### Claims Approval Matrix Reference

Partners should see their claim progress through the platform approval tiers:

| Claimed Amount | Approval Level | Maximum TAT |
|---|---|---|
| BDT 0–10K | L1 Auto/Officer | 24 Hours |
| BDT 10K–50K | L2 Manager | 3 Days |
| BDT 50K–2L | L3 Head (Joint: Business Admin + Focal Person) | 7 Days |
| BDT 2L+ | Board + Insurer Approval | 15 Days |

---

### 9.5 Policies Module

#### Purpose

Allow partners to view referred/serviced policies, initiate policy purchase on behalf of customers, and track renewals.

#### Policy Initiation on Behalf of Customer (FR-065)

1. Agent or partner staff initiates policy purchase from portal
2. Searches or creates customer profile
3. Selects insurance product from catalog (filtered by partner's assigned policy types)
4. Fills applicant details, nominee/beneficiary information
5. System sends OTP to customer's mobile for consent verification
6. Customer confirms with OTP
7. Payment processed (MFS, bank transfer, or cash with proof upload)
8. Policy issued and linked to both customer account and partner referral

**Consent flow is mandatory** — no policy can be issued without customer OTP verification.

#### Policy List View

- filterable by product type, status, date range
- search by policy number, customer name, enrollment ID
- columns: Policy Number, Customer, Product, Coverage, Premium, Status, Issued Date, Next Renewal
- status indicators: Active, Expiring (30-day window), Lapsed, Cancelled
- renewal action for expiring policies

#### Policy Detail View

- policy summary: number, product, coverage amount, premium, tenure
- customer/insured member details
- nominee/beneficiary list with share percentages
- payment history (premiums paid)
- claims linked to this policy
- endorsement history
- renewal history and upcoming renewal date
- associated agent and commission earned

#### Product Catalog Access

Partners see a filtered product catalog based on their assigned policy types:

- Health: Group Health, Individual Health, OMP/Travel Medical, Micro Health
- Motor: Private Vehicle, Commercial Vehicle
- Life: Term Life, Group Life, Micro Life
- Fire/Property: Fire, Property, Marine
- Pet: Pet Insurance, Livestock/Cattle
- Device: Laptop, Mobile, Gadget
- Travel: Travel Insurance, Schengen/Non-Schengen

Each product displays: coverage details, premium range, tenure options, exclusions, required documents.

---

### 9.6 Commission Module

#### Purpose

Provide transparent, auditable commission tracking for partners and agents.

#### Commission Dashboard

- **Total Earned** — lifetime commission earned
- **Current Month** — commission earned this month
- **Pending Payout** — approved but not yet disbursed
- **Last Payout** — most recent payout amount and date
- **Next Payout Date** — estimated next payout cycle

#### Commission List View

- filterable by type (acquisition, renewal, claims_assistance), date range, agent, policy, status
- columns: Commission ID, Policy, Type, Premium Amount, Rate (%), Commission Amount, Status, Date
- status badges: Calculated, Approved, Paid, Reversed
- total row at bottom

#### Commission Detail View

- linked policy details (number, customer, product, premium)
- commission type and rate applied
- calculation breakdown:
  - for acquisition: `Premium × Acquisition Rate`
  - for renewal: `Premium × Renewal Rate`
  - for claims assistance: `Settlement Amount × Claims Assistance Rate`
- agent split (if hierarchical): partner share vs. agent share
- payout batch reference
- payment method and reference
- audit trail

#### Commission Payout History

Aligned with `CommissionPayout` proto entity:

- payout list by period (monthly)
- columns: Payout Number, Period (Start–End), Commission Count, Total Amount, Status, Payment Method, Paid Date
- status: PENDING, APPROVED, PROCESSING, PAID, FAILED, CANCELLED
- payout detail with line-item breakdown
- downloadable commission statement (PDF/Excel)

#### Agent Commission View

For Partner Admins viewing agent commissions:

- agent-level commission summary
- per-agent breakdown by policy and type
- hierarchical roll-up (Regional Manager sees all agents, Partner Admin sees all)
- performance-based bonus tier visibility (from `CommissionConfig.performance_tiers`)

#### Commission Configuration Visibility

Partners can view (not edit) their commission structure:

- acquisition commission rate
- renewal commission rate
- claims assistance rate
- agent split configuration
- performance tier thresholds and bonus rates
- effective dates

---

### 9.7 Referrals Module

#### Purpose

Track customer referral pipeline from lead generation through policy conversion.

#### Referral Pipeline

```
Lead Created → Customer Contacted → Quote Generated → OTP Verified → Payment Completed → Policy Issued
```

#### Referral List View

- filterable by status, product type, agent, date range
- columns: Referral ID, Customer Name, Phone, Product Interest, Agent, Status, Created Date
- conversion funnel summary at top

#### Referral Detail View

- customer contact details
- product interest and selected plan
- agent who created referral
- status history with timestamps
- linked quote (if generated)
- linked policy (if converted)
- commission linked (if applicable)

#### Referral Creation

1. Agent or partner staff creates referral with customer details
2. System checks for existing customer account
3. If new customer — creates lead record
4. If existing — links to customer profile
5. Agent follows up and initiates policy purchase when ready

---

### 9.8 Documents Module

#### Purpose

Manage partner KYB verification documents, MOU agreements, and operational documents.

#### Document Categories

| Category | Documents | Upload By |
|---|---|---|
| **KYB Verification** | Trade license, TIN certificate, NID (owner), bank statement | Partner during onboarding |
| **MOU/Agreement** | Partnership MOU, service agreement, NDA | Both parties |
| **Agent KYC** | Agent NID, photo, education certificate | Agent during registration |
| **Operational** | Rate cards, formulary, service menu | Partner |
| **Claim Support** | Bills, prescriptions, reports, photos | Partner during claim |

#### Onboarding Document Checklist

For each partner type, a KYB document checklist is enforced:

**All Partners (Mandatory):**
- Trade license (image/PDF, max 5MB)
- TIN certificate (image/PDF, max 5MB)
- Owner/authorized signatory NID (front and back)
- Bank account proof (cancelled cheque or bank statement)
- Partnership MOU (signed)

**Hospital/Clinic (Additional):**
- DGHS registration certificate
- Facility photos (exterior, reception, key areas)
- Director/superintendent NID

**Pharmacy (Additional):**
- DGDA drug license
- Pharmacist license

**Auto Service Center (Additional):**
- BRTA workshop registration
- Workshop photos
- Equipment inventory

**Fire Inspector/Surveyor (Additional):**
- IDRA surveyor empanelment letter
- Professional indemnity insurance (if applicable)

**Pet Clinic/Veterinary (Additional):**
- BVA registration
- Veterinarian license

#### Document Management Features

- document library with status: Pending, Verified, Rejected, Expired
- upload with drag-and-drop, max 5MB per file, JPEG/PNG/PDF
- document preview
- version tracking
- expiry date tracking with renewal reminders
- Focal Person verification workflow
- bulk document download

---

### 9.9 Reports and Analytics

#### Purpose

Provide partner-side operational and financial reporting.

#### Performance Reports

| Report | Content |
|---|---|
| **Sales Performance** | Policies sold by product, agent, period; conversion rates |
| **Claims Performance** | Claims submitted, approved, rejected, settled; average TAT |
| **Commission Statement** | Monthly/quarterly commission earned, paid, pending by agent |
| **Agent Leaderboard** | Top agents by policies sold, commissions, customer satisfaction |
| **Revenue Dashboard** | Total premium generated, commission earned, payout history |
| **Partner Scorecard** | Overall performance score combining sales, claims, compliance |

#### Financial Reports

- commission earned vs. paid reconciliation
- outstanding payables
- claim settlement amounts by period
- premium collection summary (for agents)

#### Compliance Reports

- document expiry status
- KYB verification status
- agent license/NID verification status
- audit action log

#### Export Options

- PDF and Excel export for all reports
- scheduled report delivery via email (weekly/monthly)
- date range, product type, and agent filters on all reports

---

### 9.10 Settings Module

#### Purpose

Partner profile, notification preferences, security, and configuration management.

#### Profile Settings

- organization details (name, type, contact information)
- bank account details (masked display, verified status)
- service locations and coverage area
- operating hours
- partner logo upload

#### Notification Preferences

- channel preferences per notification type (SMS, email, in-app push)
- opt-in/opt-out for marketing and promotional notifications
- claim status update notifications
- commission payout notifications
- policy renewal reminders for referred policies

#### Security Settings

- password change
- two-factor authentication (TOTP) management
- active sessions list with revocation
- audit log of login activity

#### Team Management

- invite team members (finance, claims staff)
- assign roles
- manage permissions
- deactivate team members

---

## 10. Admin-Driven Bulk Partner Upload — Detailed Specification

This section provides the complete specification for the admin Excel upload workflow, which is the primary mechanism for building out the partner network.

### 10.1 Upload Flow

```
Admin Dashboard → Partner Management → Bulk Upload
    ↓
Download Template (XLSX)
    ↓
Fill Template (Partner-Type Specific Sheets)
    ↓
Upload File (Validation - Client-Side + Server-Side)
    ↓
Validation Report (Row-by-Row Success/Failure)
    ↓
Confirm Import (Preview Valid Rows)
    ↓
Create Partners (Status: PENDING_VERIFICATION)
    ↓
Dispatch Onboarding Invitations (SMS + Email)
    ↓
Focal Person Queue (KYB Verification)
    ↓
Partner Activated (Status: ACTIVE)
```

### 10.2 Template Design

The Excel template should have separate sheets per partner type:

1. **Hospitals** — healthcare partners for Health insurance
2. **Clinics** — smaller healthcare partners for Health/Micro insurance
3. **Pharmacies** — medicine retailers for Health insurance
4. **Auto Service Centers** — repair shops for Motor insurance
5. **Pet Clinics / Veterinary** — animal care for Pet/Livestock insurance
6. **Fire Inspectors / Surveyors** — assessment for Fire/Property insurance
7. **Device Repair** — laptop/mobile repair for Device insurance
8. **Agent Networks** — sales agents for all policy types
9. **Corporate Partners** — B2B partners for Group insurance

Each sheet has:
- header row with column names and validation comments
- data validation dropdowns for enum fields
- conditional formatting for required fields
- example row showing expected data format

### 10.3 Validation Rules

**Field-Level Validation:**

| Field | Rule |
|---|---|
| Organization Name | Non-empty, max 255 chars, trimmed |
| Partner Type | Must match PartnerType enum |
| Trade License | Non-empty, unique across all partners in system |
| TIN | 12-digit numeric, unique |
| Contact Phone | Regex: `^\+880[1][0-9]{9}$` |
| Contact Email | Valid email format |
| Bank Account | Non-empty (encrypted on ingestion) |
| Cashless Limit | Required if cashless enabled, positive integer |
| Discount Percentage | Required if discount enabled, 0–100 |
| Commission Rates | 0–100, sum check |

**Cross-Row Validation:**
- no duplicate trade license within the same upload
- no duplicate TIN within the same upload
- no duplicate contact phone within the same upload

**System-Level Validation:**
- trade license, TIN, phone, email uniqueness against existing database records
- partner type compatibility with assigned policy types

### 10.4 Error Handling

- validation report generated as downloadable Excel with error column
- rows with errors are skipped, valid rows proceed
- admin can fix errors and re-upload only failed rows
- upload is idempotent — re-uploading an already-created partner (by trade license) updates instead of duplicating

### 10.5 Policy-Type Assignment

Each partner row specifies which policy types the partner serves. This mapping drives:
- which products the partner sees in their catalog
- which claims the partner can submit
- which commission rates apply
- which cashless/discount rules are active

**Example mappings by partner type:**

| Partner Type | Default Policy Types |
|---|---|
| Hospital | Health (Group, Individual, OMP), Life (Group Life with health rider) |
| Clinic | Health (Individual, Micro) |
| Pharmacy | Health (All) — discount only |
| Auto Service Center | Motor (Private, Commercial) |
| Pet Clinic | Pet Insurance, Livestock Insurance |
| Fire Inspector | Fire, Property, Marine, Engineering |
| Laptop/Mobile Repair | Device Insurance |
| Agent Network | All policy types |
| Corporate | Group Health, Group Life |

---

## 11. Partner Onboarding Workflow — State Machine

### 11.1 Onboarding States

```
INVITED → REGISTRATION_STARTED → DOCUMENTS_UPLOADED → PENDING_VERIFICATION → VERIFICATION_IN_PROGRESS → ACTIVE
                                                                   ↓
                                                       DOCUMENTS_REQUESTED (loop)
                                                                   ↓
                                                        REJECTED (terminal)
```

### 11.2 Onboarding Steps

| Step | Actor | Action | Next State |
|---|---|---|---|
| 1 | Admin | Uploads Excel or manually creates partner | INVITED |
| 2 | Partner | Receives SMS/email, clicks registration link | REGISTRATION_STARTED |
| 3 | Partner | Creates portal account (phone + OTP), sets password | REGISTRATION_STARTED |
| 4 | Partner | Uploads KYB documents (trade license, TIN, bank proof, MOU) | DOCUMENTS_UPLOADED |
| 5 | System | Notifies Focal Person of new partner in verification queue | PENDING_VERIFICATION |
| 6 | Focal Person | Reviews documents, verifies trade license, TIN, bank | VERIFICATION_IN_PROGRESS |
| 7a | Focal Person | Approves partner → status = ACTIVE, access granted | ACTIVE |
| 7b | Focal Person | Requests additional documents → partner notified | DOCUMENTS_REQUESTED |
| 7c | Focal Person | Rejects partner → reason recorded, partner notified | REJECTED |

### 11.3 SLA Requirements

- onboarding invitation dispatched within 1 hour of Excel upload
- partner registration link valid for 7 days
- Focal Person must act within 3 business days (FR-070)
- total onboarding completion target: < 7 days (FR-059)

### 11.4 Partner Suspension and Termination (FR-069)

- Partner Admin or Focal Person can request suspension
- Suspension blocks new sales/claims but existing policies remain active
- Termination triggers graceful policy transfer mechanism
- Customer notified of partner change
- All commission accrued before suspension is still payable

---

## 12. Cross-Cutting Collaboration Features

### 12.1 Chat

- contextual conversation threads tied to claim ID, policy ID, or onboarding request
- participant labels by organization and role (Partner, Focal Person, Claims Officer)
- file attachment support (images, PDFs, documents)
- read timestamps
- searchable history
- persistent storage — not ephemeral

### 12.2 Notifications

Partners receive notifications for:

| Event | Channels |
|---|---|
| Onboarding invitation | SMS + Email |
| Document verification outcome | SMS + In-app |
| Claim status change | In-app + SMS |
| Claim document request | In-app + Email |
| Commission calculated | In-app |
| Commission payout processed | SMS + In-app |
| Policy renewal approaching (referred) | In-app |
| Partner suspension/reactivation | SMS + Email |
| Platform announcements | In-app |

### 12.3 Activity and Audit Timeline

Every entity (partner, agent, claim, commission) has an audit timeline showing:
- who performed the action
- when the action occurred
- before/after status (if applicable)
- reason/comment
- related document or entity reference

---

## 13. Data and Entity Model Expectations

The partner dashboard operates on or projects data from these primary entities:

### 13.1 Core Entities

- `Partner` — partner organization (from partner.proto)
- `Agent` — insurance agent under partner (from partner.proto)
- `Commission` — commission records (from partner.proto, partitioned by month)
- `CommissionConfig` — commission rules per product (from commission_config.proto)
- `CommissionPayout` — batch payout records (from commission_payout.proto)
- `Policy` — insurance policies referred by partner
- `Claim` — claims submitted by or for partner
- `ClaimDocument` — supporting claim documents
- `ClaimApproval` — approval workflow records
- `User` — portal user accounts
- `Session` — active sessions
- `Product` — insurance product catalog
- `ProductPlan` — plan variants
- `Payment` — payment transactions
- `Notification` — notification records

### 13.2 Partner-Specific View Entities

Views or projections needed for the partner portal:

- **PartnerDashboardSummary** — aggregated KPIs for dashboard home
- **AgentPerformance** — per-agent metrics over time
- **CommissionStatement** — periodic commission summary for download
- **ClaimsByPartner** — claims filtered to this partner's context
- **PoliciesByPartner** — policies sold/serviced by this partner
- **ReferralPipeline** — lead-to-conversion tracking
- **OnboardingStatus** — document and verification progress

---

## 14. Security, Compliance, and Audit

### 14.1 Access Control

- RBAC enforced via Authorization Service (Go, port 50070)
- partner isolation: partners can only see their own data
- agent isolation: agents see only their own portfolio unless they are Regional Manager or Partner Admin
- Focal Person has cross-partner visibility within their assignment
- all API calls validated at gateway level with partner context headers

### 14.2 Data Protection

- PII fields encrypted at rest (AES-256-GCM): bank_account, contact_phone, NID
- PII masked in logs: email, phone, NID
- consent required for: phone, email, NID, bank data (per proto annotations)
- document storage in S3 with encryption at rest, presigned URLs for access
- file upload validation: max 5MB, JPEG/PNG/PDF (no executables)
- virus/malware scan on uploaded documents

### 14.3 Session Security

- session token with 15-minute access token, 7-day refresh token
- CSRF protection via csrf_token cookie
- session revocation on logout
- 2FA enforcement for Partner Admin role (FR-017)
- account lockout after 5 failed login attempts for 30 minutes (FR-010)

### 14.4 Audit Requirements

- immutable audit logs for: partner status changes, commission calculations, claim submissions, agent registration, document uploads
- audit logs stored in PostgreSQL with append-only tables
- 20-year retention for regulatory compliance (IDRA requirement)
- all API requests logged with payload hash, timestamp, actor identity

### 14.5 Compliance

- IDRA reporting traceability for partner-related transactions
- BFIU/AML hooks for suspicious commission patterns
- no silent data deletion for records under regulatory hold
- partner MOU and agreement versioning

---

## 15. Non-Functional Requirements

### 15.1 Performance Targets

| Metric | Target |
|---|---|
| Dashboard load time | < 3 seconds |
| List/table page load | < 2 seconds |
| Search/filter response | < 500ms |
| Excel upload processing (500 rows) | < 30 seconds |
| Claim submission form | < 5 minutes to complete |
| Commission calculation | Real-time on policy activation |

### 15.2 Reliability

- graceful degradation when downstream services are unavailable
- fallback data or retry guidance instead of blank failure states
- autosave for claim submission and bill entry forms
- safe handling of partially integrated modules
- 99.5% uptime target

### 15.3 Scalability

- support 100+ active partners and 10,000+ agents (SRS target)
- support 50,000+ policies serviced across partner network
- bulk upload handles up to 1,000 partner rows per file
- pagination and virtual scrolling for large tables

### 15.4 Accessibility and Localization

- WCAG 2.1 AA compliance
- Bengali (primary) and English (secondary) language support
- responsive desktop-first layout with usable tablet/mobile fallback
- field agents should be able to use key workflows on mobile browsers

---

## 16. UI and UX Expectations

### 16.1 Branding

- platform branding: **LabaidInsuretech**
- login and shell: partner-oriented, not generic marketing
- partner context visible in header (organization name, partner type badge)
- LabaidInsuretech brand color: `#8C34C7` (purple accent)

### 16.2 Layout

- desktop-first operational dashboard
- card/grid system for KPI widgets
- high-density tables with readable spacing for claims, policies, commission lists
- sidebar navigation matching the 10-item information architecture
- consistent with B2B portal patterns (shared `web_shared` components where applicable)

### 16.3 Interaction Standards

- tab labels contained within their parent card or responsive grid (no overflow)
- no placeholder alerts for core actions
- modal, drawer, or panel patterns for document upload, bill submission, chat
- empty, loading, and error states must be purposeful and informative
- confirmation dialogs for destructive actions (suspend partner/agent)
- toast notifications for successful actions

---

## 17. Proposed Screen Set

1. Login / OTP Verification
2. Dashboard Home
3. Partner Roster (Admin View)
4. Partner Bulk Upload
5. Partner Detail
6. Partner Onboarding Wizard
7. Agent List
8. Agent Detail
9. Agent Bulk Upload
10. Claims List / Queue
11. Claim Submission Form (Bill Entry)
12. Claim Detail
13. Policy List
14. Policy Detail
15. Policy Initiation (on behalf of customer)
16. Product Catalog (filtered by partner type)
17. Commission Dashboard
18. Commission List
19. Commission Statement (downloadable)
20. Payout History
21. Referral Pipeline
22. Referral Detail
23. Documents Library
24. Document Upload
25. Reports Dashboard
26. Report Viewer (with export)
27. Settings — Profile
28. Settings — Notifications
29. Settings — Security (Password, 2FA, Sessions)
30. Settings — Team Management

---

## 18. Rollout Recommendation

### Phase A: Core Partner Operations (Aligned to M2)

- Login and partner workspace context
- Dashboard with KPIs
- Partner profile and document management
- Admin bulk Excel upload of partners
- Partner onboarding workflow
- Agent registration and management
- Basic claims list and status tracking
- Commission dashboard (view-only)
- Settings

### Phase B: Operational Maturity (Aligned to M3)

- Full claim bill submission workflow (cashless + reimbursement)
- Product-specific claim forms (health, motor, fire, pet, device)
- Policy initiation on behalf of customer with OTP consent
- Referral pipeline management
- Commission payout history and statement download
- Agent bulk upload
- Agent hierarchy and territory management
- Chat integration for claim communication
- Notification preferences

### Phase C: Intelligence and Scale (Aligned to D/S)

- Advanced reports and analytics
- Partner performance scorecard
- Agent leaderboard with gamification
- Predictive analytics (churn, claim likelihood)
- Geographic analytics (policy distribution, claims heatmap)
- Partner API for external integration
- White-label/branding capability (FR-205)
- Commission structure self-configuration (FR-206)
- BI tool integration (Metabase/Tableau/Power BI)

---

## 19. Open Integration Dependencies

The following dependencies should be tracked during delivery:

| Dependency | Owner | Status |
|---|---|---|
| Partner Service CRUD APIs (Go) | InScore team | Active — port 50100/50101 |
| Agent registration and management APIs | InScore team | Active — port 50100/50101 |
| Commission calculation engine (C#) | PoliSync team | gRPC contract ✅, domain logic ~10% |
| Claims submission and approval APIs (C#) | PoliSync team | gRPC contract ✅, domain logic ~10% |
| Policy lifecycle APIs (C#) | PoliSync team | gRPC contract ✅, domain logic ~35% |
| Excel upload parsing and validation | Partner Portal BFF | To be built |
| Document storage and retrieval | Storage Service (Go) | Active — port 50290/50291 |
| Notification dispatch | Notification Service (Go) | Active — port 50230/50231 |
| Payment processing for policy initiation | Payment Service (Node.js) | Active — port 50190/50191 |
| Commission payout processing | Payment + Ledger | To be wired |
| Fraud scoring integration | Fraud Service (Go) | Active — port 50220/50221 |
| PDF generation for statements | DocGen Service (Go) | Active — port 50280/50281 |

---

## 20. Current Partner Portal Codebase Status

### 20.1 Existing Implementation

The current `partner_portal/` directory is a fork of the `b2b_portal/` with the following state:

**Functional pages:**
- Claims — uses `ClaimsPage` component with mock data (Health, Auto, Life claims)
- Policies — uses `PoliciesPage` component with mock data (Health, Auto, Life policies)
- Settings — uses `Settings` component (reused from B2B)

**Stub pages (placeholders only):**
- Dashboard Home — only `StatsCards` and `OverviewActivity` active; PolicyOverview, QuickAccess, UpcomingPayments commented out
- Partners — empty stub, comment says "Partner list will be implemented here"
- Agents — empty stub, comment says "Agent list will be implemented here"
- Commission — empty stub, comment says "Commission tracking will be implemented here"

**Existing navigation (7 items):**
Dashboard, Partners, Agents, Claims, Policies, Commission, Settings

**Components inherited from B2B that need repurposing or removal:**
- `billing-invoices/` — B2B invoicing, not applicable
- `departments/` — B2B departments, not applicable
- `employees/` — B2B employees, not applicable for partners
- `purchase-orders/` — B2B procurement, not applicable
- `quotations/` — generic quotations, may be partially reusable

**Components that may be reusable:**
- `stats-cards/` — adaptable for partner KPIs
- `overview-activity/` — adaptable for partner activity feed
- `settings/` — fully reusable
- `dashboard-layout.tsx` — shell/layout reusable
- `dashboard-sidebar.tsx` — navigation sidebar reusable
- `dashboard-header.tsx` — header reusable

**Auth flow:**
- Login with OTP via `(auth)/login` and `(auth)/otp` routes — functional

### 20.2 Navigation Update Required

Current navigation needs expansion from 7 to 10 items per Section 8 of this spec:

| Current | Proposed |
|---|---|
| Dashboard | Dashboard (enhance) |
| Partners | Partners (build) |
| Agents | Agents (build) |
| Claims | Claims (enhance with bill submission) |
| Policies | Policies (enhance with initiation) |
| Commission | Commission (build) |
| — | Referrals (new) |
| — | Documents (new) |
| — | Reports (new) |
| Settings | Settings (enhance) |

---

## 21. Success Criteria

The partner dashboard should be considered successful when:

- partners can complete onboarding through the portal (registration → KYB → activation)
- admin can bulk-upload partner networks via Excel with policy-type mapping
- partners can submit cashless/reimbursement bills through the portal
- partners can track claim status from submission to settlement
- agents can be registered, managed, and tracked under partner organizations
- commission is transparently visible with calculation breakdown and payout history
- partners can initiate policy purchase on behalf of customers with OTP consent
- referral pipeline tracks lead-to-conversion journey
- KYB documents are managed with verification workflow
- reports provide actionable operational and financial intelligence
- the portal remains aligned with the existing platform architecture and deployment model
- field agents can perform key workflows on mobile browsers
- SLA targets are met: onboarding < 7 days, claim TAT per approval matrix, commission payout monthly

---

## 22. Source Basis

This specification was derived from:

- `documentation/About/INSURER_DASHBOARD_FULL_SPEC.md` — insurer dashboard spec as structural template
- `documentation/About/ARCHITECTURE_OVERVIEW.md` — B2B portal architecture and auth flow
- `documentation/About/API_ROUTES_SUMMARY.md` — API route conventions
- `documentation/About/POLISYNC_REFERENCE.md` — C# insurance engine reference
- `documentation/SRS_v3/SPECS_V3.7/sections/02_system_overview.md` — system context and objectives
- `documentation/SRS_v3/SPECS_V3.7/sections/03_architecture.md` — microservices architecture and service registry
- `documentation/SRS_v3/SPECS_V3.7/sections/04_functional_requirements.md` — FR-059 through FR-072 (Partner & Agent Management), FR-141 through FR-148 (Partner Portal & BI)
- `documentation/SRS_v3/SPECS_V3.7/sections/06_data_model.md` — proto-first data model strategy
- `documentation/SRS_v3/SPECS_V3.7/sections/08_integration.md` — gRPC service contracts, Kafka events
- `proto/insuretech/partner/entity/v1/partner.proto` — Partner, Agent, Commission, PartnerBenefits, PartnerType, PartnerStatus entities
- `proto/insuretech/commission/entity/v1/commission_config.proto` — CommissionConfig, RevenueModel
- `proto/insuretech/commission/entity/v1/commission_payout.proto` — CommissionPayout, PayoutStatus
- `proto/insuretech/commission/entity/v1/revenue_share.proto` — RevenueShare
- `proto/insuretech/partner/events/v1/partner_events.proto` — domain events
- `partner_portal/` codebase — current implementation state (mostly stubs)
- `b2b_portal/` codebase — source fork with architecture patterns

---

## APPENDIX A: Partner Excel Upload Template — Column Reference

### A.1 Hospital Partners Sheet

| Column | Type | Required | Validation | Example |
|--------|------|----------|------------|---------|
| Organization Name | text | yes | max 255 | LabAid Hospital Dhanmondi |
| Trade License | text | yes | unique | TL-DHK-2024-001234 |
| TIN Number | text | yes | 12-digit | 123456789012 |
| DGHS Registration | text | yes | | DGHS/DHK/H-2024-456 |
| Contact Email | email | yes | valid email | admin@labaidhospital.com |
| Contact Phone | text | yes | +880 format | +8801712345678 |
| Bank Account | text | yes | encrypted | 1234567890123 |
| Bank Name | text | yes | | Dutch Bangla Bank |
| Bank Branch | text | no | | Dhanmondi Branch |
| Bed Count | integer | yes | > 0 | 250 |
| Emergency Capable | boolean | yes | | TRUE |
| Specialty Departments | text | no | comma-separated | Cardiology,Orthopedics,Neurology |
| Cashless Enabled | boolean | yes | | TRUE |
| Cashless Limit (BDT) | integer | conditional | positive | 500000 |
| Auto-Approval Threshold | integer | no | positive | 10000 |
| Pre-Auth Required | boolean | no | default TRUE | TRUE |
| Discount Enabled | boolean | no | | TRUE |
| Discount Percentage | decimal | conditional | 0–100 | 15.00 |
| Service Locations | text | no | | Dhaka,Chittagong |
| Nationwide | boolean | no | | FALSE |
| Acquisition Rate (%) | decimal | yes | 0–100 | 20.00 |
| Renewal Rate (%) | decimal | yes | 0–100 | 10.00 |
| Claims Assist Rate (%) | decimal | no | 0–100 | 5.00 |

### A.2 Auto Service Center Partners Sheet

| Column | Type | Required | Validation | Example |
|--------|------|----------|------------|---------|
| Organization Name | text | yes | max 255 | Dhaka Motors Workshop |
| Trade License | text | yes | unique | TL-DHK-2024-005678 |
| TIN Number | text | yes | 12-digit | 987654321098 |
| BRTA Registration | text | yes | | BRTA/WS/2024-789 |
| Contact Email | email | yes | | info@dhakamotors.com |
| Contact Phone | text | yes | +880 format | +8801812345678 |
| Bank Account | text | yes | encrypted | 9876543210123 |
| Bank Name | text | yes | | Sonali Bank |
| Service Types | text | yes | comma-separated | Bodywork,Mechanical,Electrical,Paint |
| Cashless Enabled | boolean | yes | | TRUE |
| Cashless Limit (BDT) | integer | conditional | positive | 200000 |
| Auto-Approval Threshold | integer | no | positive | 15000 |
| Discount Enabled | boolean | no | | FALSE |
| Service Locations | text | no | | Dhaka |
| Acquisition Rate (%) | decimal | yes | 0–100 | 15.00 |
| Renewal Rate (%) | decimal | yes | 0–100 | 8.00 |

### A.3 Pet Clinic / Veterinary Partners Sheet

| Column | Type | Required | Validation | Example |
|--------|------|----------|------------|---------|
| Organization Name | text | yes | max 255 | Dhaka Vet Care |
| Trade License | text | yes | unique | TL-DHK-2024-009012 |
| TIN Number | text | yes | 12-digit | 456789012345 |
| BVA Registration | text | yes | | BVA/DHK/2024-123 |
| Vet License Number | text | yes | | VET-2024-456 |
| Contact Email | email | yes | | info@dhakavet.com |
| Contact Phone | text | yes | +880 format | +8801912345678 |
| Bank Account | text | yes | encrypted | 4567890123456 |
| Bank Name | text | yes | | Janata Bank |
| Species Coverage | text | yes | comma-separated | Dogs,Cats,Birds |
| Cashless Enabled | boolean | yes | | TRUE |
| Cashless Limit (BDT) | integer | conditional | positive | 100000 |
| Discount Enabled | boolean | no | | TRUE |
| Discount Percentage | decimal | conditional | 0–100 | 10.00 |
| Service Locations | text | no | | Dhaka,Gazipur |
| Acquisition Rate (%) | decimal | yes | 0–100 | 12.00 |
| Renewal Rate (%) | decimal | yes | 0–100 | 6.00 |

### A.4 Fire Inspector / Surveyor Partners Sheet

| Column | Type | Required | Validation | Example |
|--------|------|----------|------------|---------|
| Organization Name | text | yes | max 255 | Bangladesh Survey Associates |
| Trade License | text | yes | unique | TL-DHK-2024-003456 |
| TIN Number | text | yes | 12-digit | 345678901234 |
| IDRA Empanelment | text | yes | | IDRA/SV/2024-789 |
| Surveyor License | text | yes | | SL-2024-321 |
| Contact Email | email | yes | | info@bdsurvey.com |
| Contact Phone | text | yes | +880 format | +8801612345678 |
| Bank Account | text | yes | encrypted | 3456789012345 |
| Bank Name | text | yes | | Agrani Bank |
| Survey Types | text | yes | comma-separated | Fire,Property,Marine,Engineering |
| Service Locations | text | no | | Dhaka,Chittagong,Sylhet |
| Nationwide | boolean | no | | TRUE |
| Acquisition Rate (%) | decimal | yes | 0–100 | 10.00 |

---

## APPENDIX B: Cashless vs. Discount Partner Configurations

### B.1 Cashless Partners

Cashless partners provide services to insured customers where the insurer/platform settles the bill directly with the partner. The customer pays no out-of-pocket cost (or only the co-pay/deductible).

**Applicable partner types:** Hospital, Clinic, Auto Service Center, Pet Clinic, Laptop/Mobile Repair, Ambulance

**Configuration per partner-policy-type combination (from `PartnerBenefits` proto):**

| Config Field | Purpose | Example |
|---|---|---|
| `cashless_enabled` | Toggle cashless for this partner | TRUE |
| `cashless_limit` | Maximum cashless amount (paisa) | 50000000 (BDT 500,000) |
| `auto_approval_threshold` | Claims below this auto-approved | 1000000 (BDT 10,000) |
| `pre_authorization_required` | Require pre-auth before service | TRUE (hospitals) |
| `authorization_validity_days` | Pre-auth validity period | 30 |
| `required_documents` | Documents needed for cashless claim | ["discharge_summary", "prescription", "bills"] |
| `service_locations` | Covered locations | ["Dhaka", "Chittagong"] |
| `nationwide_coverage` | Available nationwide | FALSE |

**Cashless claim flow:**
1. Customer presents at partner facility
2. Partner checks eligibility via portal
3. Pre-authorization obtained (if required)
4. Service delivered
5. Partner submits itemized bill via portal
6. Platform verifies and routes to approval
7. Settlement paid directly to partner's bank account

### B.2 Discount Partners

Discount partners provide services at a reduced rate to insured customers. The customer pays a discounted price; the discount amount is either absorbed by the partner or claimed from the insurer.

**Applicable partner types:** Pharmacy, Clinic, Doctor, some Auto Service Centers

**Configuration per partner-policy-type combination:**

| Config Field | Purpose | Example |
|---|---|---|
| `discount_enabled` | Toggle discount for this partner | TRUE |
| `discount_percentage` | Standard discount rate | 15.00% |
| `min_discount` | Minimum discount offered | 5.00% |
| `max_discount` | Maximum discount offered | 25.00% |
| `discount_type` | Discount applies to | SERVICE, PRODUCT, CONSULTATION |

**Discount claim flow:**
1. Customer presents policy at partner facility
2. Partner verifies eligibility via portal
3. Partner applies configured discount percentage
4. Customer pays discounted amount
5. Partner may submit discounted amount for platform tracking/reporting
6. No direct settlement from platform to partner for discount model

---

## APPENDIX C: Commission Calculation Examples

### C.1 Acquisition Commission

**Scenario:** Agent sells a Health Insurance policy (Premium: BDT 14,500)

- Acquisition commission rate (agent): 20%
- Partner split: 5% (of agent commission)
- Agent earns: BDT 14,500 × 20% × 95% = **BDT 2,755**
- Partner earns: BDT 14,500 × 20% × 5% = **BDT 145**
- Total acquisition commission: **BDT 2,900**

### C.2 Renewal Commission

**Scenario:** Same policy renews (Renewal premium: BDT 14,500)

- Renewal commission rate: 10%
- Agent earns: BDT 14,500 × 10% × 95% = **BDT 1,377.50**
- Partner earns: BDT 14,500 × 10% × 5% = **BDT 72.50**
- Total renewal commission: **BDT 1,450**

### C.3 Claims Assistance Commission

**Scenario:** Partner assists customer with Health claim settlement (Settled: BDT 45,000)

- Claims assistance rate: 5%
- Partner earns: BDT 45,000 × 5% = **BDT 2,250**
- Agent: no claims assistance commission (partner-level only)

### C.4 Commission Payout Cycle

- commission calculated on policy activation or claim settlement event
- accrued daily
- payout batch generated monthly (1st of following month)
- approval workflow: System → Business Admin approval
- payment via bank transfer (primary) or MFS (fallback)
- commission statement available for download by 5th of each month

---

## APPENDIX D: Partner Portal API Routes (Proposed)

### D.1 Authentication Routes (Reuse from B2B Pattern)

| Method | Route | Purpose |
|--------|-------|---------|
| POST | `/api/auth/login` | Partner login |
| POST | `/api/auth/logout` | End session |
| GET | `/api/auth/session` | Validate session |
| POST | `/api/auth/refresh` | Refresh token |
| GET | `/api/auth/profile` | Get profile |
| PATCH | `/api/auth/profile` | Update profile |
| POST | `/api/auth/send-otp` | Request OTP |
| POST | `/api/auth/verify-otp` | Verify OTP |

### D.2 Partner Routes

| Method | Route | Purpose |
|--------|-------|---------|
| GET | `/api/partners` | List partners (admin) |
| GET | `/api/partners/me` | Get current partner profile |
| PATCH | `/api/partners/me` | Update partner profile |
| POST | `/api/partners/bulk-upload` | Excel bulk upload |
| GET | `/api/partners/bulk-upload/template` | Download Excel template |
| GET | `/api/partners/[id]` | Get partner detail |
| PATCH | `/api/partners/[id]` | Update partner |
| POST | `/api/partners/[id]/suspend` | Suspend partner |
| POST | `/api/partners/[id]/activate` | Activate partner |
| GET | `/api/partners/[id]/documents` | List partner documents |
| POST | `/api/partners/[id]/documents` | Upload document |

### D.3 Agent Routes

| Method | Route | Purpose |
|--------|-------|---------|
| GET | `/api/agents` | List agents under partner |
| POST | `/api/agents` | Register new agent |
| POST | `/api/agents/bulk-upload` | Excel bulk upload |
| GET | `/api/agents/bulk-upload/template` | Download template |
| GET | `/api/agents/[id]` | Get agent detail |
| PATCH | `/api/agents/[id]` | Update agent |
| POST | `/api/agents/[id]/suspend` | Suspend agent |
| POST | `/api/agents/[id]/activate` | Activate agent |

### D.4 Claims Routes

| Method | Route | Purpose |
|--------|-------|---------|
| GET | `/api/claims` | List claims for partner |
| POST | `/api/claims` | Submit new claim/bill |
| GET | `/api/claims/[id]` | Get claim detail |
| POST | `/api/claims/[id]/documents` | Upload claim documents |
| GET | `/api/claims/[id]/documents` | List claim documents |
| POST | `/api/claims/[id]/notes` | Add claim note |

### D.5 Policy Routes

| Method | Route | Purpose |
|--------|-------|---------|
| GET | `/api/policies` | List policies for partner |
| GET | `/api/policies/[id]` | Get policy detail |
| POST | `/api/policies/initiate` | Initiate policy on behalf of customer |
| GET | `/api/products` | List products (filtered by partner type) |
| GET | `/api/products/[id]` | Get product detail |

### D.6 Commission Routes

| Method | Route | Purpose |
|--------|-------|---------|
| GET | `/api/commissions` | List commissions |
| GET | `/api/commissions/summary` | Commission dashboard summary |
| GET | `/api/commissions/[id]` | Get commission detail |
| GET | `/api/payouts` | List payout history |
| GET | `/api/payouts/[id]` | Get payout detail |
| GET | `/api/payouts/[id]/statement` | Download commission statement |

### D.7 Referral Routes

| Method | Route | Purpose |
|--------|-------|---------|
| GET | `/api/referrals` | List referrals |
| POST | `/api/referrals` | Create referral |
| GET | `/api/referrals/[id]` | Get referral detail |
| PATCH | `/api/referrals/[id]` | Update referral status |

### D.8 Reports Routes

| Method | Route | Purpose |
|--------|-------|---------|
| GET | `/api/reports/sales` | Sales performance report |
| GET | `/api/reports/claims` | Claims performance report |
| GET | `/api/reports/agents` | Agent performance report |
| GET | `/api/reports/commission` | Commission report |
| GET | `/api/reports/export` | Export report (PDF/Excel) |

### D.9 Dashboard Routes

| Method | Route | Purpose |
|--------|-------|---------|
| GET | `/api/dashboard/summary` | Dashboard KPIs |
| GET | `/api/dashboard/activity` | Recent activity feed |
| GET | `/api/dashboard/announcements` | Platform announcements |
