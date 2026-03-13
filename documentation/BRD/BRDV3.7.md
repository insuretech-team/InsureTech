# LabAid InsureTech Platform

Business Requirements Document (BRD)

**Version:** 2.1  
**Date:** Dec 2025  
**Status:** INITIAL DRAFT (Derived from SRS V3.7)  
**Company:** LabAid InsureTech  
**Technology Partner:** LifePlus

[[[PAGEBREAK]]]

## Revision History

| Version | Date | Revised By | Description |
|---------|------|------------|-------------|
|2.0 | Dec 2025 | Faruk Hannan| BRD created from SRS V3.7; business framing, user stories, and grouped functional requirements |
|2.1| Dec 2025 | AI Engine| Enhancements,format, correction,fact check |

[[[PAGEBREAK]]]

## Executive Summary

The LabAid InsureTech Platform is a digital-first insurance distribution and operations platform for Bangladesh. It enables end-to-end journeys including onboarding & KYC, product discovery, quotation and purchase, policy issuance and lifecycle management, premium collection and reconciliation, claims submission/approval/settlement, partner & agent management, customer support, notifications, and regulatory-grade audit/compliance.

This BRD translates the **System Requirements Specification (SRS) V3.7** into business-facing requirements, stakeholder goals, and user stories. It is written for business owners, product managers, delivery teams, and partner stakeholders.

Key business outcomes
- Expand insurance access via micro-insurance products and digital channels.
- Enable high-throughput partner distribution (MFS, hospitals, e-commerce, agent networks).
- Reduce operational cost through workflow automation, self-service, and event-driven notifications.
- Achieve compliance readiness (IDRA, BFIU/AML-CFT) with auditability and reporting.

---

[[[PAGEBREAK]]]

## Table of Contents

1. [Business Overview](#1-business-overview)
2. [Scope](#2-scope)
3. [Stakeholders & Personas](#3-stakeholders--personas)
4. [User Stories](#4-user-stories)
5. [Business Processes](#5-business-processes)
6. [Business Requirements Summary (Functional)](#6-business-requirements-summary-functional)
7. [Project Assumptions & Dependencies](#7-project-assumptions--dependencies)
8. [Acceptance Criteria (Business-Level)](#8-acceptance-criteria-business-level)
9. [Appendices](#9-appendices)

---
[[[PAGEBREAK]]]


---
[[[PAGEBREAK]]]


# Executive Summary

## Document Purpose

This Business Requirements Document (BRD) translates the **LabAid InsureTech Platform System Requirements Specification (SRS) V3.7** into business-facing requirements, user stories, and acceptance criteria. It is intended for:

- **Business Executives:** Understand strategic value and investment rationale
- **Product Managers:** Define features, prioritize roadmap, manage stakeholders
- **Delivery Teams:** Translate business needs into technical implementation
- **Partners & Investors:** Understand platform capabilities and business model

## Platform Vision

The LabAid InsureTech Platform is a **digital-first insurance distribution and operations platform** designed for the Bangladesh market. It enables:

- **Customers:** Discover, purchase, manage, and claim insurance digitally (mobile-first, Bengali/English)
- **Partners:** Distribute insurance products via embedded channels (MFS, hospitals, e-commerce, agent networks)
- **Insurers:** Operate efficiently with automated workflows, compliance-ready audit trails, and risk controls
- **Regulators:** Access required reports and data with full auditability (IDRA, BFIU/AML-CFT)

## Business Value Proposition

| Stakeholder | Value Delivered |
|------------|-----------------|
| **Customers** | Fast onboarding, instant policy issuance, transparent claims, multi-language support, digital-first convenience |
| **Distribution Partners** | Embedded insurance offerings, commission tracking, assisted sales tools, white-label capability (future) |
| **InsureTech Operations** | Reduced manual work, automated workflows, fraud detection, real-time dashboards, scalable infrastructure |
| **Regulators** | Compliance-ready reporting, immutable audit trails, AML/CFT monitoring, long-term record retention |

## Key Capabilities (High-Level)

1. **Customer Onboarding & Authentication** — Mobile-first registration, OTP, biometric login, profile management
2. **Authorization & Multi-Tenancy** — Role-based access control, partner data isolation, admin MFA
3. **Product Catalog Management** — Multi-language products, dynamic pricing, lifecycle management
4. **Policy Purchase & Issuance** — End-to-end digital purchase, nominee management, instant digital policy documents
5. **Policy Servicing** — Renewals, endorsements, cancellations, refunds with transparent workflows
6. **Premium Collection & Reconciliation** — MFS/bank/card payments, manual payment verification, receipts, balance sheets
7. **Claims Management** — Digital submission, document uploads, tiered approvals, fraud detection, settlement
8. **Partner & Agent Management** — KYB onboarding, tenant isolation, commission tracking, performance dashboards
9. **Customer Support** — Ticketing, FAQ, escalation workflows, CSAT tracking
10. **Notifications** — SMS/email/push notifications with consent management and anti-spam controls
11. **Fraud Detection & Risk Controls** — Configurable rules, review queues, audit trails
12. **Regulatory Compliance** — Audit logging, data retention, lawful access workflows, AML/CFT monitoring, IDRA reporting readiness
13. **Analytics & Reporting** — Executive dashboards, operational reports, compliance extracts
14. **Integrations** — NID verification, payment gateways, SMS/email, hospital/EHR systems

## Success Metrics (Business KPIs)

| KPI | What It Measures | Target |
|-----|------------------|--------|
| **Policy Issuance Volume** | Policies issued per month | Growth >20% MoM (early stage) |
| **Customer Acquisition Cost (CAC)** | Cost to acquire one policyholder | <500 BDT (via partners) |
| **Claims Settlement TAT** | Time from submission to settlement | <7 days (simple claims <48 hours) |
| **Payment Success Rate** | % of payment attempts that succeed | >95% |
| **Partner Activation** | Active partners selling policies | 50+ partners within 12 months |
| **Customer Satisfaction (CSAT)** | Post-support/claims satisfaction score | >4.2/5 |
| **Fraud Prevention** | % of fraudulent claims caught before payout | >90% detection rate |

## Investment Rationale

- **Market Opportunity:** Low insurance penetration in Bangladesh; digital channels unlock mass-market micro-insurance
- **Regulatory Readiness:** Built-in compliance reduces risk and accelerates approvals
- **Scalable Platform:** Multi-tenant architecture supports hundreds of partners without proportional cost increase
- **Revenue Streams:** Commission on premiums, value-added services (analytics, IoT risk programs), partner white-label (future)

## Risks & Mitigation

| Risk | Mitigation |
|------|------------|
| Regulatory changes (IDRA, BFIU) | Configurable rules, versioned compliance logic, proactive engagement with regulators |
| Payment gateway downtime | Multi-provider strategy, fallback to manual verification, health monitoring |
| Customer trust/adoption barriers | Transparent processes, Bengali language support, agent-assisted onboarding, digital policy verification |
| Fraud/abuse | Multi-layered fraud detection (rules + manual review), audit trails, account lockouts |

## Document Structure

This BRD is organized as follows:

1. **Business Context** — Market, stakeholders, regulatory environment
2. **Portal & Channel Definitions** — What we are building and for whom
3. **Feature Group Requirements** — Detailed user stories, business rules, workflows, and acceptance criteria for each capability (23 feature groups)
4. **Non-Functional Requirements** — Performance, availability, security, scalability
5. **Security & Compliance** — Controls, AML/CFT operations, IDRA readiness
6. **Integrations** — External system dependencies and data flows
7. **Traceability** — Mapping back to SRS V3.7 for requirements coverage

[[[PAGEBREAK]]]


---
[[[PAGEBREAK]]]


# Business Context

## Market Overview

### Insurance Landscape in Bangladesh

- **Penetration:** <1% of GDP (among lowest in South Asia)
- **Challenges:** Limited distribution, low awareness, trust barriers, paper-heavy processes
- **Opportunities:** Mobile penetration >100%, MFS adoption (bKash, Nagad, Rocket), growing middle class, government push for financial inclusion

### Target Segments

| Segment | Profile | Insurance Needs | Distribution Channel |
|---------|---------|-----------------|---------------------|
| **Urban Middle Class** | Salaried, smartphone users | Health, motor, term life | Direct (web/app) + MFS partners |
| **Rural/Semi-Urban** | Low digital literacy, MFS users | Micro health, accident | Agent-assisted, hospital partnerships |
| **SME Owners** | Business insurance | Business continuity, liability | Partner banks, e-commerce platforms |

## Regulatory Environment

### Key Regulators

1. **Insurance Development and Regulatory Authority (IDRA)**
   - Regulates insurance companies and products
   - Requires product disclosure, policy document standards, financial solvency reporting
   - Audit and inspection rights

2. **Bangladesh Financial Intelligence Unit (BFIU)**
   - Anti-Money Laundering (AML) and Countering the Financing of Terrorism (CFT)
   - Transaction monitoring, suspicious transaction reporting (STR/SAR)
   - Record retention and reporting obligations

### Compliance Requirements (Summary)

- **Product Approval:** Insurance products require IDRA pre-approval
- **Policy Documentation:** Standardized policy documents, mandatory disclosures
- **Customer Identity:** KYC/NID verification for policies above thresholds
- **AML/CFT Monitoring:** Transaction monitoring, threshold-based alerts, STR filing
- **Data Retention:** Long-term retention of policies, payments, claims, customer communications
- **Audit Readiness:** Ability to provide records within defined timelines for regulatory requests

## Stakeholder Ecosystem

### Internal Stakeholders

| Stakeholder | Role | Primary Concerns |
|------------|------|------------------|
| **Business Executives** | Strategy, P&L ownership | Growth, profitability, regulatory risk |
| **Product & Operations** | Product definition, day-to-day ops | Feature velocity, operational efficiency, customer satisfaction |
| **Compliance & Risk** | Regulatory adherence, fraud prevention | Audit readiness, fraud loss, regulatory penalties |
| **Customer Support** | Issue resolution, customer satisfaction | Ticket volumes, resolution time, escalation handling |
| **Technology Leadership** | Platform delivery and operations | Uptime, scalability, security, technical debt |

### External Stakeholders

| Stakeholder | Relationship | Expectations |
|------------|--------------|--------------|
| **Customers** | End users (policyholders) | Fast onboarding, transparent processes, timely claims, Bengali support |
| **Distribution Partners** | MFS providers, hospitals, e-commerce, agent orgs | Easy integration, commission transparency, branded experience (future) |
| **Insurance Companies** | Underwriters (if platform acts as intermediary) | Accurate data, risk controls, regulatory compliance |
| **Regulators** | IDRA, BFIU | Timely reporting, lawful access, compliance with rules |

## Competitive Landscape

### Current Players

- **Traditional Insurers:** Paper-heavy, branch-based distribution, slow claims
- **Emerging Digital Platforms:** Limited to specific products (motor, health), lack comprehensive platform features

### Competitive Advantages (LabAid InsureTech Platform)

| Advantage | How We Deliver |
|-----------|----------------|
| **Digital-First** | Mobile app, web, instant policy issuance |
| **Multi-Channel Distribution** | Partner API, agent tools, embedded insurance |
| **Regulatory-Ready** | Built-in compliance, audit trails, AML/CFT |
| **Bengali Language Support** | Native language for mass-market adoption |
| **Transparent Claims** | Real-time status, digital submission, fast TAT |
| **Scalable Multi-Tenancy** | Support hundreds of partners without linear cost growth |

## Business Model

### Revenue Streams

1. **Commission on Premiums:** Primary revenue from policies sold
2. **Partner Subscription Fees:** (Future) Monthly fees for partner portal access and white-label options
3. **Value-Added Services:** (Future) Analytics, IoT risk monitoring, AI underwriting (where regulatory-permissible)

### Go-To-Market Strategy

- **Phase 1 (M1-M2):** Launch with 5-10 anchor partners (MFS, 2-3 hospitals); focus on health and accident micro-insurance
- **Phase 2 (M3):** Expand to motor insurance, e-commerce partnerships, agent network expansion
- **Phase 3 (M4+):** Advanced features (AI, IoT, voice), white-label partner offerings, cross-border exploration (future)

[[[PAGEBREAK]]]


---
[[[PAGEBREAK]]]


# 4. Portals & Channels (What We Are Building)

This BRD explicitly defines the portals/channels required by SRS V3.7 (FG-023 and cross-cutting requirements).
Each portal definition below includes its business purpose, primary users, and must-have capabilities.

## 4.1 Customer Mobile App

**Primary users:** Customers/Policyholders  
**Business purpose:** Self-service onboarding, discovery, purchase, policy servicing, claims, support.

**Must-have capabilities (business view):**
- Role-appropriate login and account safety controls
- Clear dashboards for primary tasks and statuses
- Auditability for sensitive actions
- Multi-language support where customer-facing (Bengali/English)

## 4.2 Customer Web Portal (PWA)

**Primary users:** Customers/Policyholders  
**Business purpose:** Web equivalent of the customer journey (campaign traffic, desktop users).

**Must-have capabilities (business view):**
- Role-appropriate login and account safety controls
- Clear dashboards for primary tasks and statuses
- Auditability for sensitive actions
- Multi-language support where customer-facing (Bengali/English)

## 4.3 Partner Admin Portal

**Primary users:** Partner admins (MFS/hospital/e-commerce/agent org)  
**Business purpose:** Partner onboarding, agent management, assisted sales, commission/analytics, operational tooling.

**Must-have capabilities (business view):**
- Role-appropriate login and account safety controls
- Clear dashboards for primary tasks and statuses
- Auditability for sensitive actions
- Multi-language support where customer-facing (Bengali/English)

## 4.4 Agent App (Mobile)

**Primary users:** Agents operating under a partner  
**Business purpose:** Assisted onboarding/purchase, lead handling, commissions, basic support tooling.

**Must-have capabilities (business view):**
- Role-appropriate login and account safety controls
- Clear dashboards for primary tasks and statuses
- Auditability for sensitive actions
- Multi-language support where customer-facing (Bengali/English)

## 4.5 Focal Person Portal

**Primary users:** InsureTech focal persons  
**Business purpose:** Partner KYB verification/approval, dispute resolution, monitoring, escalations.

**Must-have capabilities (business view):**
- Role-appropriate login and account safety controls
- Clear dashboards for primary tasks and statuses
- Auditability for sensitive actions
- Multi-language support where customer-facing (Bengali/English)

## 4.6 Business Admin Portal

**Primary users:** InsureTech business ops  
**Business purpose:** Product governance, workflow approvals, claims controls, reporting, operational configuration.

**Must-have capabilities (business view):**
- Role-appropriate login and account safety controls
- Clear dashboards for primary tasks and statuses
- Auditability for sensitive actions
- Multi-language support where customer-facing (Bengali/English)

## 4.7 System Admin Portal

**Primary users:** Platform/system administrators  
**Business purpose:** Security configuration, roles, system health, incident tooling, audit.

**Must-have capabilities (business view):**
- Role-appropriate login and account safety controls
- Clear dashboards for primary tasks and statuses
- Auditability for sensitive actions
- Multi-language support where customer-facing (Bengali/English)

## 4.8 Support/Call Centre Portal

**Primary users:** Customer support agents  
**Business purpose:** Ticketing, customer history, escalation workflows, communication tools.

**Must-have capabilities (business view):**
- Role-appropriate login and account safety controls
- Clear dashboards for primary tasks and statuses
- Auditability for sensitive actions
- Multi-language support where customer-facing (Bengali/English)

## 4.9 Regulatory Access Portal (controlled)

**Primary users:** IDRA/BFIU or auditors (as per lawful request)  
**Business purpose:** Controlled delivery of requested reports/data with full audit trail.

**Must-have capabilities (business view):**
- Role-appropriate login and account safety controls
- Clear dashboards for primary tasks and statuses
- Auditability for sensitive actions
- Multi-language support where customer-facing (Bengali/English)

## 4.10 IoT Device Management Portal (partner/internal)

**Primary users:** IoT partners / internal ops  
**Business purpose:** Device onboarding, health monitoring, telemetry visibility (where IoT program is active).

**Must-have capabilities (business view):**
- Role-appropriate login and account safety controls
- Clear dashboards for primary tasks and statuses
- Auditability for sensitive actions
- Multi-language support where customer-facing (Bengali/English)

**Traceability:** SRS FG-023 (FR-244..FR-248) and portal-related capabilities across FG-001..FG-019.

[[[PAGEBREAK]]]


---
[[[PAGEBREAK]]]


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


---
[[[PAGEBREAK]]]


# 7. Non-Functional Requirements (NFR) — Business-Grade Detail

NFRs are non-negotiable business constraints because they define customer experience, reliability of money movement, regulatory readiness, and operational cost.

## 7.1 NFR Catalog (Derived from SRS Section 5)

### NFR-046 — Database Technology

- **Business requirement:** Maintain relational data integrity using **PostgreSQL V17** with JSONB support
- **Target/Measurement:** ACID compliance tests
- **Priority:** M1
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-047 — Caching & Session

- **Business requirement:** Use **Redis** for distributed caching and session management
- **Target/Measurement:** Cache hit ratio monitoring
- **Priority:** M1
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-048 — API Protocol

- **Business requirement:** Microservices communication shall use **gRPC with Protocol Buffers** (Category 1)
- **Target/Measurement:** Inter-service latency metrics
- **Priority:** M1
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-049 — Client API

- **Business requirement:** Client-facing APIs shall use **REST (OpenAPI 3.0)** with **JWT** authentication
- **Target/Measurement:** Schema validation, Token checks
- **Priority:** M1
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-050 — Public Integration

- **Business requirement:** External integrations shall use **RESTful APIs** with **OpenAPI 3.0** specifications (Category 3)
- **Target/Measurement:** Swagger validator pass
- **Priority:** D
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-051 — Search Engine

- **Business requirement:** Full-text search capabilities shall be implemented using **PostgreSQL Full-Text Search** or dedicated engine
- **Target/Measurement:** Query performance <200ms
- **Priority:** M1
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-052 — Object Storage

- **Business requirement:** Document and static asset storage shall use **S3-compatible storage** (AWS/DigitalOcean)
- **Target/Measurement:** Upload/Download latency
- **Priority:** M1
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-053 — Message Broker

- **Business requirement:** Asynchronous event processing shall be handled by **Apache Kafka**
- **Target/Measurement:** Throughput monitoring
- **Priority:** M1
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-054 — Time-Series Data

- **Business requirement:** IoT telemetry data shall be stored in **TimescaleDB**
- **Target/Measurement:** Ingestion rate monitoring
- **Priority:** M2
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-055 — Vector Database

- **Business requirement:** Vector embeddings for AI features shall be stored in **Pgvector** or **Pinecone**
- **Target/Measurement:** Similarity search latency
- **Priority:** D
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-056 — Graph Database

- **Business requirement:** Fraud visualization and relationship mapping shall use **Neo4j** or **Amazon Neptune**
- **Target/Measurement:** Graph traversal depth/speed
- **Priority:** D
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-057 — Columnar Database

- **Business requirement:** High-performance analytics shall use **ClickHouse** or **Druid**
- **Target/Measurement:** Analytical query speed
- **Priority:** D
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-058 — Financial Ledger

- **Business requirement:** Double-entry bookkeeping shall be enforced using **TigerBeetle**
- **Target/Measurement:** Ledger reconciliation check
- **Priority:** M3
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-059 — Mobile Framework

- **Business requirement:** Cross-platform mobile application shall be built using **React Native**
- **Target/Measurement:** Code reuse >80%
- **Priority:** M1
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-060 — CDN & Security

- **Business requirement:** Public entry points shall be secured via **Cloudflare** proxy
- **Target/Measurement:** WAF block rate
- **Priority:** M1
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-001 — API response time for policy operations

- **Business requirement:** < 500ms (95th percentile)
- **Target/Measurement:** Application performance monitoring
- **Priority:** M1
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-002 — Database query response time

- **Business requirement:** < 100ms (average)
- **Target/Measurement:** Database monitoring tools
- **Priority:** M1
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-003 — Mobile app startup time

- **Business requirement:** < 3 seconds
- **Target/Measurement:** App performance analytics
- **Priority:** M1
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-004 — Web portal page load time

- **Business requirement:** < 2 seconds
- **Target/Measurement:** Browser performance tools
- **Priority:** M1
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-005 — Payment processing time

- **Business requirement:** < 10 seconds end-to-end
- **Target/Measurement:** Payment gateway analytics
- **Priority:** M1
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-006 — Claim processing automation

- **Business requirement:** 80% straight-through processing
- **Target/Measurement:** Business process monitoring
- **Priority:** M2
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-007 — Report generation time

- **Business requirement:** < 30 seconds for standard reports
- **Target/Measurement:** Reporting system metrics
- **Priority:** M2
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-008 — Search functionality response

- **Business requirement:** < 200ms for basic searches
- **Target/Measurement:** Search performance monitoring
- **Priority:** M2
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-009 — Concurrent user support

- **Business requirement:** 10,000 active users
- **Target/Measurement:** Load testing and monitoring
- **Priority:** M1
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-010 — Transaction throughput

- **Business requirement:** 1,000 TPS (policies + claims)
- **Target/Measurement:** Performance testing
- **Priority:** M2
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-011 — Database scalability

- **Business requirement:** 100 million policy records
- **Target/Measurement:** Database performance testing
- **Priority:** M2
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-012 — Auto-scaling capability

- **Business requirement:** Scale out/in based on load
- **Target/Measurement:** Infrastructure monitoring
- **Priority:** M2
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-013 — Peak load handling

- **Business requirement:** 5x normal load during campaigns
- **Target/Measurement:** Stress testing
- **Priority:** M3
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-014 — Storage scalability

- **Business requirement:** 10TB+ document storage
- **Target/Measurement:** Cloud storage metrics
- **Priority:** M3
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-015 — System availability

- **Business requirement:** 99.5% uptime (M1), 99.9% (M2)
- **Target/Measurement:** Infrastructure monitoring
- **Priority:** M1
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-016 — Recovery Time Objective (RTO)

- **Business requirement:** 4 hours maximum
- **Target/Measurement:** Disaster recovery testing
- **Priority:** M1
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-017 — Recovery Point Objective (RPO)

- **Business requirement:** 1 hour maximum data loss
- **Target/Measurement:** Backup and recovery testing
- **Priority:** M1
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-018 — Mean Time To Recovery (MTTR)

- **Business requirement:** < 2 hours
- **Target/Measurement:** Incident response metrics
- **Priority:** M2
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-019 — Service degradation handling

- **Business requirement:** Graceful degradation during outages
- **Target/Measurement:** Chaos engineering testing
- **Priority:** M2
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-020 — Data backup frequency

- **Business requirement:** Real-time replication + daily backups
- **Target/Measurement:** Backup monitoring
- **Priority:** M1
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-029 — User satisfaction score

- **Business requirement:** 4.5+ stars on app stores
- **Target/Measurement:** User feedback and ratings
- **Priority:** M2
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-030 — Task completion rate

- **Business requirement:** 95% for critical user journeys
- **Target/Measurement:** User experience analytics
- **Priority:** M1
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-031 — Learning curve

- **Business requirement:** New users complete first task < 5 minutes
- **Target/Measurement:** User onboarding metrics
- **Priority:** M2
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-032 — Error recovery

- **Business requirement:** Clear error messages with action guidance
- **Target/Measurement:** Error tracking and analysis
- **Priority:** M1
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-033 — Accessibility compliance

- **Business requirement:** WCAG 2.1 AA compliance
- **Target/Measurement:** Accessibility testing tools
- **Priority:** M2
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-034 — Multi-language support

- **Business requirement:** Bengali and English localization
- **Target/Measurement:** Localization testing
- **Priority:** M1
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-035 — Code coverage

- **Business requirement:** 80% unit test coverage
- **Target/Measurement:** Automated testing reports
- **Priority:** M2
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-036 — Deployment frequency

- **Business requirement:** Daily deployments capability
- **Target/Measurement:** CI/CD pipeline metrics
- **Priority:** M2
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-037 — Mean Time To Deploy

- **Business requirement:** < 30 minutes for hotfixes
- **Target/Measurement:** Deployment automation metrics
- **Priority:** M2
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-038 — Monitoring coverage

- **Business requirement:** 100% critical path monitoring
- **Target/Measurement:** Observability platform
- **Priority:** M1
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-039 — Log aggregation

- **Business requirement:** Centralized logging for all services
- **Target/Measurement:** Logging platform metrics
- **Priority:** M1
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-040 — Documentation currency

- **Business requirement:** API documentation auto-generated
- **Target/Measurement:** Documentation automation
- **Priority:** M2
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

[[[PAGEBREAK]]]


---
[[[PAGEBREAK]]]


# 8. Security, Privacy, and Compliance Requirements (Detailed)

Security and compliance are business requirements: they protect customers, protect funds, enable partner trust, and satisfy IDRA/BFIU expectations.
This section translates SRS Section 7 controls into business-operational requirements.

## 8.1 Security Control Catalog (SEC)

### SEC-001

- **Business control requirement:** Use separate secret vault - AWS KMS/Azure Key Vault/HashiCorp, 90-day key rotation
- **Priority:** M1
- **Evidence (examples):** policies, logs, key rotation records, access reviews, penetration test reports, audit extracts.

### SEC-002

- **Business control requirement:** Use Data Masking: NID (last 3 digits), phone (mask middle), email (mask username)
- **Priority:** M2
- **Evidence (examples):** policies, logs, key rotation records, access reviews, penetration test reports, audit extracts.

### SEC-003

- **Business control requirement:** Follow PCI-DSS compliance for card flows - Approach: Hosted payment page (redirect model) - DO NOT store card data, Level: SAQ-A (simplest, for redirecting merchants), Requirements: Annual SAQ, quarterly ASV scans, TLS 1.3, Tokenization: Store only gateway tokens for recurring payments
- **Priority:** M2
- **Evidence (examples):** policies, logs, key rotation records, access reviews, penetration test reports, audit extracts.

### SEC-004

- **Business control requirement:** Have AML/CFT detection hooks - Transaction Monitoring: 20+ automated rules for AML detection including Rapid purchases (>3 policies in 7 days), High-value premiums (>BDT 5 lakh), Frequent cancellations, Mismatched nominees, Geographic/payment anomalies
- **Priority:** D
- **Evidence (examples):** policies, logs, key rotation records, access reviews, penetration test reports, audit extracts.

### SEC-005

- **Business control requirement:** Have IDRA reporting capabilities following IDRA data format - Monthly Reports: Premium Collection (Form IC-1), Claims Intimation (Form IC-2), Quarterly Reports: Claims Settlement (IC-3), Financial Performance (IC-4), Annual Reports: FCR (Financial Condition Report), CARAMELS Framework Returns, Event-Based: Significant incidents (48hrs), fraud cases (7 days), Platform: Report generator with IDRA Excel templates, audit trail, 20-year archive
- **Priority:** D
- **Evidence (examples):** policies, logs, key rotation records, access reviews, penetration test reports, audit extracts.

### SEC-006

- **Business control requirement:** Have regular penetration testing - Penetration Testing: Pre-launch + annually (SISA InfoSec or international firm)
- **Priority:** D
- **Evidence (examples):** policies, logs, key rotation records, access reviews, penetration test reports, audit extracts.

### SEC-007

- **Business control requirement:** Have regular security audits from various security auditors and regulatory bodies and maintain compliance
- **Priority:** D
- **Evidence (examples):** policies, logs, key rotation records, access reviews, penetration test reports, audit extracts.

### SEC-008

- **Business control requirement:** DAST: OWASP ZAP/Burp Suite (weekly on staging)
- **Priority:** D
- **Evidence (examples):** policies, logs, key rotation records, access reviews, penetration test reports, audit extracts.

### SEC-009

- **Business control requirement:** SAST: SonarQube/Checkmarx (every commit, block critical vulnerabilities)
- **Priority:** D
- **Evidence (examples):** policies, logs, key rotation records, access reviews, penetration test reports, audit extracts.

### SEC-010

- **Business control requirement:** Virus scanning: ClamAV on uploaded files
- **Priority:** M
- **Evidence (examples):** policies, logs, key rotation records, access reviews, penetration test reports, audit extracts.

### SEC-021

- **Business control requirement:** Implement API rate limiting per user/IP: 1000 requests/hour for authenticated users, 100 requests/hour for anonymous
- **Priority:** M2
- **Evidence (examples):** policies, logs, key rotation records, access reviews, penetration test reports, audit extracts.

### SEC-022

- **Business control requirement:** Maintain separate encryption keys for different data types with hierarchical key management
- **Priority:** M2
- **Evidence (examples):** policies, logs, key rotation records, access reviews, penetration test reports, audit extracts.

### SEC-023

- **Business control requirement:** Implement real-time security incident response with automated threat isolation
- **Priority:** M2
- **Evidence (examples):** policies, logs, key rotation records, access reviews, penetration test reports, audit extracts.

### SEC-024

- **Business control requirement:** Perform continuous vulnerability assessment with automated patching for critical vulnerabilities
- **Priority:** D
- **Evidence (examples):** policies, logs, key rotation records, access reviews, penetration test reports, audit extracts.

### SEC-025

- **Business control requirement:** Implement zero-trust network architecture with microsegmentation
- **Priority:** D
- **Evidence (examples):** policies, logs, key rotation records, access reviews, penetration test reports, audit extracts.

### SEC-011

- **Business control requirement:** IDRA Monthly Reports: Generate Form IC-1 (Premium Collection) by 10th of each month with breakdown by product line, geographic region, partner channel in Excel format per IDRA template v2024
- **Priority:** M2
- **Evidence (examples):** policies, logs, key rotation records, access reviews, penetration test reports, audit extracts.

### SEC-012

- **Business control requirement:** IDRA Monthly Reports: Generate Form IC-2 (Claims Intimation) by 10th of each month listing all new claims with policy number, claim amount, claim type, date of intimation
- **Priority:** M2
- **Evidence (examples):** policies, logs, key rotation records, access reviews, penetration test reports, audit extracts.

### SEC-013

- **Business control requirement:** IDRA Quarterly Reports: Generate Form IC-3 (Claims Settlement) within 15 days of quarter-end showing settlement ratio, average TAT, pending >30 days breakdown
- **Priority:** M2
- **Evidence (examples):** policies, logs, key rotation records, access reviews, penetration test reports, audit extracts.

### SEC-014

- **Business control requirement:** IDRA Quarterly Reports: Generate Form IC-4 (Financial Performance) within 20 days of quarter-end with premium earned, claims paid, commission paid, net profit/loss
- **Priority:** M3
- **Evidence (examples):** policies, logs, key rotation records, access reviews, penetration test reports, audit extracts.

### SEC-015

- **Business control requirement:** IDRA Annual FCR: Generate Financial Condition Report (FCR) within 90 days of year-end including full CARAMELS framework assessment with external auditor sign-off
- **Priority:** M3
- **Evidence (examples):** policies, logs, key rotation records, access reviews, penetration test reports, audit extracts.

### SEC-016

- **Business control requirement:** IDRA Event-Based Reporting: Report significant incidents (fraud >BDT 1L, data breach, system outage >4hrs) within 48 hours via IDRA portal
- **Priority:** M3
- **Evidence (examples):** policies, logs, key rotation records, access reviews, penetration test reports, audit extracts.

### SEC-017

- **Business control requirement:** AML/CFT Concrete Triggers: Flag transactions matching: (1) >3 policies in 7 days, (2) Premium >BDT 5L without income proof, (3) Nominee mismatch with no relationship doc, (4) Payment from third-party account, (5) Frequent cancellations >2 in 30 days, (6) Geographic anomaly (policy in Dhaka, payment from remote district), (7) Multiple failed KYC attempts >3, (8) PEP match in screening
- **Priority:** M3
- **Evidence (examples):** policies, logs, key rotation records, access reviews, penetration test reports, audit extracts.

### SEC-018

- **Business control requirement:** SAR Workflow: (1) System auto-flags suspicious transaction → (2) Compliance Officer reviews within 24hrs → (3) If confirmed suspicious, escalate to Business Admin+Focal Person → (4) Prepare SAR with evidence → (5) Submit to BFIU within 3 business days → (6) Mark account for enhanced monitoring → (7) Do NOT notify customer (tipping off prohibited)
- **Priority:** M3
- **Evidence (examples):** policies, logs, key rotation records, access reviews, penetration test reports, audit extracts.

### SEC-019

- **Business control requirement:** Data Deletion Exceptions: Customer data deletion requests processed within 30 days EXCEPT: (a) Active policy holders (deletion after policy expiry+7yrs), (b) Ongoing claims (deletion after settlement+7yrs), (c) Under SAR investigation (deletion prohibited until case closed), (d) Regulatory hold (deletion requires IDRA/BFIU approval)
- **Priority:** M3
- **Evidence (examples):** policies, logs, key rotation records, access reviews, penetration test reports, audit extracts.

### SEC-020

- **Business control requirement:** Right to Erasure Workflow: Customer submits deletion request → System validates exceptions → If eligible, anonymize PII while retaining transaction records → Generate deletion certificate → Notify customer within 30 days
- **Priority:** D
- **Evidence (examples):** policies, logs, key rotation records, access reviews, penetration test reports, audit extracts.

## 8.2 AML/CFT Operating Model (Business View)

The platform must support configurable AML monitoring rules, alerting, investigation workflow, and STR/SAR filing support with strict auditability.
(See SRS Section 7.7.x for rule tables and workflows.)

## 8.3 IDRA Reporting and Record-Keeping (Business View)

The platform must retain and produce long-term records (policies, payments, claims, cancellations, approvals, customer communications) with retrieval capability within required SLAs.

[[[PAGEBREAK]]]


---
[[[PAGEBREAK]]]


# 9. Traceability (BRD ↔ SRS)

This matrix ensures every SRS requirement is accounted for in business terms.

## 9.1 Functional Traceability (FG/FR → BRD Coverage)

| SRS Feature Group | Covered in BRD Section(s) |
|---|---|
| FG-001 | Section 6 (Detailed Functional Catalog), plus portal/process sections where relevant |
| FG-002 | Section 6 (Detailed Functional Catalog), plus portal/process sections where relevant |
| FG-003 | Section 6 (Detailed Functional Catalog), plus portal/process sections where relevant |
| FG-004 | Section 6 (Detailed Functional Catalog), plus portal/process sections where relevant |
| FG-005 | Section 6 (Detailed Functional Catalog), plus portal/process sections where relevant |
| FG-06 | Section 6 (Detailed Functional Catalog), plus portal/process sections where relevant |
| FG-007 | Section 6 (Detailed Functional Catalog), plus portal/process sections where relevant |
| FG-008 | Section 6 (Detailed Functional Catalog), plus portal/process sections where relevant |
| FG-009 | Section 6 (Detailed Functional Catalog), plus portal/process sections where relevant |
| FG-010 | Section 6 (Detailed Functional Catalog), plus portal/process sections where relevant |
| FG-011 | Section 6 (Detailed Functional Catalog), plus portal/process sections where relevant |
| FG-012 | Section 6 (Detailed Functional Catalog), plus portal/process sections where relevant |
| FG-013 | Section 6 (Detailed Functional Catalog), plus portal/process sections where relevant |
| FG-014 | Section 6 (Detailed Functional Catalog), plus portal/process sections where relevant |
| FG-015 | Section 6 (Detailed Functional Catalog), plus portal/process sections where relevant |
| FG-016 | Section 6 (Detailed Functional Catalog), plus portal/process sections where relevant |
| FG-017 | Section 6 (Detailed Functional Catalog), plus portal/process sections where relevant |
| FG-018 | Section 6 (Detailed Functional Catalog), plus portal/process sections where relevant |
| FG-019 | Section 6 (Detailed Functional Catalog), plus portal/process sections where relevant |
| FG-020 | Section 6 (Detailed Functional Catalog), plus portal/process sections where relevant |
| FG-022 | Section 6 (Detailed Functional Catalog), plus portal/process sections where relevant |
| FG-023 | Section 6 (Detailed Functional Catalog), plus portal/process sections where relevant |

## 9.2 Non-Functional Traceability (NFR → BRD)

| SRS NFR ID | Covered in BRD Section |
|---|---|
| NFR-001 | Section 7 |
| NFR-002 | Section 7 |
| NFR-003 | Section 7 |
| NFR-004 | Section 7 |
| NFR-005 | Section 7 |
| NFR-006 | Section 7 |
| NFR-007 | Section 7 |
| NFR-008 | Section 7 |
| NFR-009 | Section 7 |
| NFR-010 | Section 7 |
| NFR-011 | Section 7 |
| NFR-012 | Section 7 |
| NFR-013 | Section 7 |
| NFR-014 | Section 7 |
| NFR-015 | Section 7 |
| NFR-016 | Section 7 |
| NFR-017 | Section 7 |
| NFR-018 | Section 7 |
| NFR-019 | Section 7 |
| NFR-020 | Section 7 |
| NFR-029 | Section 7 |
| NFR-030 | Section 7 |
| NFR-031 | Section 7 |
| NFR-032 | Section 7 |
| NFR-033 | Section 7 |
| NFR-034 | Section 7 |
| NFR-035 | Section 7 |
| NFR-036 | Section 7 |
| NFR-037 | Section 7 |
| NFR-038 | Section 7 |
| NFR-039 | Section 7 |
| NFR-040 | Section 7 |
| NFR-046 | Section 7 |
| NFR-047 | Section 7 |
| NFR-048 | Section 7 |
| NFR-049 | Section 7 |
| NFR-050 | Section 7 |
| NFR-051 | Section 7 |
| NFR-052 | Section 7 |
| NFR-053 | Section 7 |
| NFR-054 | Section 7 |
| NFR-055 | Section 7 |
| NFR-056 | Section 7 |
| NFR-057 | Section 7 |
| NFR-058 | Section 7 |
| NFR-059 | Section 7 |
| NFR-060 | Section 7 |

## 9.3 Security Traceability (SEC → BRD)

| SRS SEC ID | Covered in BRD Section |
|---|---|
| SEC-001 | Section 8 |
| SEC-002 | Section 8 |
| SEC-003 | Section 8 |
| SEC-004 | Section 8 |
| SEC-005 | Section 8 |
| SEC-006 | Section 8 |
| SEC-007 | Section 8 |
| SEC-008 | Section 8 |
| SEC-009 | Section 8 |
| SEC-010 | Section 8 |
| SEC-011 | Section 8 |
| SEC-012 | Section 8 |
| SEC-013 | Section 8 |
| SEC-014 | Section 8 |
| SEC-015 | Section 8 |
| SEC-016 | Section 8 |
| SEC-017 | Section 8 |
| SEC-018 | Section 8 |
| SEC-019 | Section 8 |
| SEC-020 | Section 8 |
| SEC-021 | Section 8 |
| SEC-022 | Section 8 |
| SEC-023 | Section 8 |
| SEC-024 | Section 8 |
| SEC-025 | Section 8 |

[[[PAGEBREAK]]]


---
[[[PAGEBREAK]]]


# Feature Group: User Management & Authentication (FG-001)

## Business Objective

Enable secure, low-friction onboarding and authentication for customers, partners, and agents across all digital channels. The platform must support Bangladesh mobile-first users while meeting enterprise security standards.

**Business Value:**
- Minimize drop-off during registration (target: <5% abandonment)
- Prevent fraud through strong identity controls
- Support multi-channel access (mobile app, web, partner portals)
- Meet regulatory requirements for customer identity verification

## Actors & Portals

| Actor | Portal(s) | Primary Use Cases |
|-------|-----------|-------------------|
| Customer | Mobile App, Web PWA | Self-registration, login, profile management |
| Agent | Agent Mobile App, Partner Portal | Assisted registration, customer lookup |
| Partner Admin | Partner Portal | User management, access control |
| Business Admin | Admin Portal | User monitoring, manual verification |

## User Stories

### US-FG-001-01: Customer Self-Registration (Mobile First)

**As a** potential customer  
**I want** to register using only my mobile number  
**So that** I can quickly start exploring insurance products

**Acceptance Criteria:**
- User enters +880 Bangladesh mobile number
- OTP is sent within 60 seconds (FR-002)
- OTP is 6-digit, valid for 5 minutes
- User can resend OTP (max 3 attempts per 15 min window - FR-003)
- Duplicate phone numbers are rejected (FR-004)
- Profile completion required before first purchase

**Flow:**
1. User opens app/web → "Register" screen
2. Enters mobile number → validates format
3. Taps "Send OTP" → backend sends SMS
4. User enters OTP → validates → creates account
5. System prompts for profile completion (name, DOB, gender, etc.)

![User Registration Flow](images/flow_registration_otp.png)

**Exception Paths:**
- OTP not received → "Resend" option (rate limited)
- Invalid OTP → error message, retry (3 attempts then lockout)
- Phone already registered → redirect to login

**Related FRs:** FR-001, FR-002, FR-003, FR-004, FR-005

### US-FG-001-02: Secure Password & Biometric Login

**As a** returning customer  
**I want** to use biometric login on my phone  
**So that** I can access my account quickly and securely

**Acceptance Criteria:**
- Password must meet policy: 8+ chars, 1 uppercase, 1 number, 1 special (FR-006)
- Biometric (fingerprint/face ID) available on supported devices (FR-007)
- Fallback to password if biometric fails
- Session tokens managed securely (15-min access, 7-day refresh - FR-009)

![Login Options](images/flow_login_biometric.png)

**Related FRs:** FR-006, FR-007, FR-009

### US-FG-001-03: Password Recovery

**As a** customer who forgot their password  
**I want** to reset it using my registered mobile  
**So that** I can regain access without calling support

**Acceptance Criteria:**
- User taps "Forgot Password"
- OTP sent to registered mobile (FR-008)
- User enters OTP + new password
- Password policy enforced
- Success confirmation

**Related FRs:** FR-008

### US-FG-001-04: Account Protection (Lockout & MFA)

**As a** business owner  
**I want** accounts locked after repeated failed logins  
**So that** brute-force attacks are prevented

**Acceptance Criteria:**
- 5 failed login attempts → 30-minute lockout (FR-010)
- User notified via SMS
- Admin users require MFA (future: FR-017 references admin MFA)

**Related FRs:** FR-010, FR-017

## Business Rules

| Rule ID | Description |
|---------|-------------|
| BR-AUTH-01 | One mobile number = one customer account (uniqueness enforced) |
| BR-AUTH-02 | OTP valid for 5 minutes; max 3 resends per 15-min window |
| BR-AUTH-03 | Password policy: 8 chars, 1 upper, 1 number, 1 special |
| BR-AUTH-04 | Session access token: 15 min; refresh token: 7 days |
| BR-AUTH-05 | Account lockout: 5 failed attempts → 30 min block |

## Key Workflows

### Registration → First Login
1. User registers (mobile + OTP)
2. System creates user record (status: pending_profile)
3. User completes profile (FR-011: name, DOB, gender, occupation, address)
4. Profile validated → account status: active
5. User can now login and purchase

### Login (Existing User)
1. User enters mobile/email + password (or biometric)
2. System validates credentials
3. System generates access + refresh tokens
4. User redirected to dashboard

### Password Reset
1. User taps "Forgot Password"
2. System sends OTP to registered mobile
3. User validates OTP + sets new password
4. System logs event for audit

## Data Model Notes

**User Entity (SRS Proto: insuretech.authn.entity.v1.User)**
- user_id (UUID)
- mobile_number (unique)
- email (unique, optional)
- password_hash
- profile_complete (boolean)
- account_status (active, suspended, locked)
- created_at, updated_at

**Session Entity (SRS Proto: insuretech.authn.entity.v1.Session)**
- session_id
- user_id
- access_token (JWT, 15 min)
- refresh_token (7 days)
- device_info

## Integration Touchpoints

| External System | Purpose | Failure Handling |
|----------------|---------|------------------|
| SMS Gateway | OTP delivery | Retry 3x; manual verification queue |
| NID Verification (optional) | Identity verification | Async; manual review if API down |

## Security & Privacy

- Passwords: bcrypt hashed, never logged
- OTPs: single-use, time-bound, rate-limited
- Session tokens: JWT with secure signing
- PII (mobile, email): encrypted at rest
- Audit: all login attempts logged (timestamp, device, IP, outcome)

## NFR Constraints

| NFR | Target | Why It Matters |
|-----|--------|----------------|
| Availability | 99.9% | Registration/login downtime = lost customers |
| OTP Delivery | <60s, 95% success | Slow OTP = abandonment |
| Password Hash | bcrypt work factor ≥12 | Protect against breach |

## Acceptance Criteria (Business-Level)

- [ ] Customer can register via mobile + OTP end-to-end
- [ ] Duplicate registrations are blocked
- [ ] OTP rate limiting prevents abuse
- [ ] Password policy is enforced with clear error messages
- [ ] Biometric login works on iOS/Android (for supported devices)
- [ ] Account lockout triggers after 5 failed attempts
- [ ] All auth events are auditable

## Traceability

**SRS Reference:** FG-001 — 4.1 User Management & Authentication  
**Functional Requirements:** FR-001, FR-002, FR-003, FR-004, FR-005, FR-006, FR-007, FR-008, FR-009, FR-010, FR-011, FR-012, FR-013

[[[PAGEBREAK]]]


---
[[[PAGEBREAK]]]


# Feature Group: Authorization & Multi-Tenancy (FG-002)

## Business Objective

Implement role-based access control (RBAC) and attribute-based access control (ABAC) to ensure users see only what they're authorized to access. Support multi-tenant architecture so partner data is strictly isolated.

**Business Value:**
- Protect sensitive data (policies, claims, commissions) from unauthorized access
- Enable partner/agent hierarchies with granular permissions
- Support compliance audits (who accessed what, when)
- Scale to hundreds of partners without cross-contamination

## Actors & Portals

| Actor | Portal(s) | Role Hierarchy |
|-------|-----------|----------------|
| Customer | Mobile App, Web | Customer role |
| Agent | Agent Mobile App | Agent role (inherits Customer) |
| Partner Admin | Partner Portal | Partner Admin (inherits Agent) |
| Focal Person | Admin Portal | Focal Person (cross-partner view) |
| Business Admin | Admin Portal | Business Admin (full access) |
| System Admin | Admin Portal | System Admin (platform config) |

## User Stories

### US-FG-002-01: Role-Based Dashboard

**As a** Partner Admin  
**I want** to see only my partner's data (agents, customers, policies)  
**So that** I cannot access competitors' information

**Acceptance Criteria:**
- Partner Admin logs in → sees partner-scoped dashboard
- Cannot query other partners' data (enforced at API level)
- All queries include tenant_id filter
- Audit log records user + tenant context

![Partner Admin Dashboard](images/dashboard_partner_admin.png)

**Related FRs:** FR-014, FR-015, FR-016

### US-FG-002-02: Hierarchical Role Inheritance

**As a** system designer  
**I want** roles to inherit permissions (Partner Admin > Agent > Customer)  
**So that** permission management is simple and consistent

**Acceptance Criteria:**
- Partner Admin can do everything an Agent can + partner management
- Agent can do everything a Customer can + assisted sales
- Permissions checked at API gateway + service layer

**Related FRs:** FR-019

### US-FG-002-03: Admin Multi-Factor Authentication

**As a** security officer  
**I want** admin accounts to require MFA  
**So that** privileged access is protected

**Acceptance Criteria:**
- Admin login requires OTP (SMS or TOTP app) after password
- MFA setup enforced on first admin login
- Backup codes provided for account recovery

**Related FRs:** FR-017

## Business Rules

| Rule ID | Description |
|---------|-------------|
| BR-AUTHZ-01 | All API requests include user_id + tenant_id + role |
| BR-AUTHZ-02 | Partner data isolated by tenant_id (cannot cross-query) |
| BR-AUTHZ-03 | Focal Person can view multiple tenants (special privilege) |
| BR-AUTHZ-04 | System Admin actions require approval workflow (future) |
| BR-AUTHZ-05 | Role changes logged and require Business Admin approval |

## Key Workflows

### Permission Check (Every API Call)
1. User sends request with access token
2. API Gateway validates token → extracts user_id, role, tenant_id
3. Service checks permissions: "Can [role] perform [action] on [resource]?"
4. If authorized → process; else → 403 Forbidden
5. Audit log records: user, action, resource, outcome, timestamp

### Partner Admin Onboarding
1. Focal Person creates partner account (tenant_id assigned)
2. Partner Admin invited (email + temp password)
3. Partner Admin logs in → forced password change + MFA setup
4. Partner Admin can now create agents under their tenant

## Data Model Notes

**Role Entity (SRS Proto: insuretech.authz.entity.v1.Role)**
- role_id
- role_name (Customer, Agent, Partner Admin, Focal Person, etc.)
- permissions (list of actions)
- inherits_from (role hierarchy)

**Tenant Entity**
- tenant_id (partner_id)
- tenant_name
- status (active, suspended)

**User-Role-Tenant Mapping**
- user_id + role_id + tenant_id (many-to-many)

## Integration Touchpoints

| System | Purpose |
|--------|---------|
| API Gateway | Enforces token validation, role checks |
| Audit Service | Logs all access decisions |

## Security & Privacy

- Permissions enforced at both API Gateway and service layer (defense in depth)
- Tenant isolation validated in automated tests
- Admin MFA mandatory for production access
- Audit logs immutable and retained per compliance period

## NFR Constraints

| NFR | Target |
|-----|--------|
| Authorization Latency | <50ms per permission check |
| Audit Log Availability | 99.99% (critical for compliance) |

## Acceptance Criteria

- [ ] Partner Admin can only see their own partner's data
- [ ] Agents can perform assisted sales within their partner
- [ ] Focal Person can view cross-partner data for oversight
- [ ] Admin MFA is enforced and cannot be bypassed
- [ ] All permission checks are audited

## Traceability

**SRS Reference:** FG-002 — 4.2 Authorization & Access Control  
**Functional Requirements:** FR-014, FR-015, FR-016, FR-017, FR-018, FR-019, FR-020

[[[PAGEBREAK]]]


---
[[[PAGEBREAK]]]


# Feature Group: Product Catalog Management (FG-003)

## Business Objective

Enable business users to define, configure, and manage insurance products without developer involvement. Support multi-language product descriptions, dynamic pricing rules, and product lifecycle management.

**Business Value:**
- Accelerate time-to-market for new products (target: <1 week from approval to live)
- Enable A/B testing and seasonal campaigns
- Support regulatory-compliant product disclosures
- Multi-language support for Bangladesh market (Bengali + English)

## Actors & Portals

| Actor | Portal(s) | Responsibilities |
|-------|-----------|------------------|
| Business Admin | Admin Portal | Create/update products, pricing, coverage rules |
| Customer | Mobile App, Web | Browse, search, compare products |
| Agent | Agent Mobile App | Product catalog for assisted sales |
| Partner Admin | Partner Portal | View partner-authorized products |

## User Stories

### US-FG-003-01: Product Lifecycle Management

**As a** Business Admin  
**I want** to create a new insurance product with coverage details, pricing, and terms  
**So that** customers can purchase it immediately after approval

**Acceptance Criteria:**
- Admin navigates to "Products" → "Create New"
- Fills: product name, category, type, coverage amount, premium base, terms
- Adds Bengali + English descriptions (FR-029)
- Sets status: Draft → Pending Approval → Active
- Version history maintained (FR-028)

![Product Creation Flow](images/admin_product_create.png)

**Related FRs:** FR-021, FR-023, FR-028, FR-029

### US-FG-003-02: Customer Product Discovery

**As a** customer  
**I want** to browse products by category and search by keyword  
**So that** I can find the right insurance for my needs

**Acceptance Criteria:**
- Homepage shows product categories (Health, Motor, Travel, etc.)
- Customer taps category → sees filtered list
- Search bar supports Bengali/English keywords
- Product card shows: name, coverage summary, starting price
- Tap product → detailed product page with full terms

![Product Catalog - Customer View](images/customer_product_catalog.png)

**Related FRs:** FR-024, FR-025, FR-026

### US-FG-003-03: Product Comparison

**As a** customer  
**I want** to compare up to 3 products side-by-side  
**So that** I can make an informed purchase decision

**Acceptance Criteria:**
- Customer selects 2-3 products → "Compare" button
- Comparison table shows: coverage, exclusions, premium, deductible, co-pay
- Highlights differences
- CTA: "Choose Plan" → redirects to purchase

![Product Comparison](images/customer_product_compare.png)

**Related FRs:** FR-027

## Business Rules

| Rule ID | Description |
|---------|-------------|
| BR-PROD-01 | Only Active products visible to customers |
| BR-PROD-02 | Product changes create new version (audit trail) |
| BR-PROD-03 | Multi-language descriptions mandatory (Bengali + English) |
| BR-PROD-04 | Product deactivation requires Business Admin approval |
| BR-PROD-05 | Pricing rules support: flat, tiered, age-based, risk-based |

## Key Workflows

### Product Creation & Approval
1. Business Admin creates product (status: Draft)
2. Reviews internally → status: Pending Approval
3. Focal Person/Business Admin approves → status: Active
4. Product appears in customer catalog immediately
5. All changes versioned

### Customer Product Discovery → Purchase
1. Customer opens app → browses category or searches
2. Views product detail page (coverage, exclusions, terms)
3. Taps "Get Quote" → proceeds to purchase flow (FG-004)

## Data Model Notes

**Product Entity (SRS Proto: insuretech.products.entity.v1.Product)**
- product_id
- product_name (multilang: en, bn)
- category (HEALTH, MOTOR, TRAVEL, etc.)
- product_type (TERM, WHOLE_LIFE, VEHICLE, etc.)
- coverage_amount (min, max)
- base_premium
- coverage_details (JSON)
- exclusions (JSON)
- terms_and_conditions (multilang)
- status (DRAFT, ACTIVE, INACTIVE)
- version_history

## Integration Touchpoints

| System | Purpose |
|--------|---------|
| CMS (future) | Manage rich product content (images, videos) |
| Pricing Engine | Calculate dynamic premiums based on product rules |

## Security & Privacy

- Product versioning prevents accidental data loss
- Only authorized Business Admin can activate/deactivate products
- Audit log for all product changes

## NFR Constraints

| NFR | Target |
|-----|--------|
| Product Catalog Load Time | <2s for customer-facing pages |
| Search Response Time | <500ms |
| Multilingual Support | Bengali, English (future: more) |

## Acceptance Criteria

- [ ] Business Admin can create/update/deactivate products
- [ ] Customers can browse, search, and compare products
- [ ] Multi-language descriptions work correctly
- [ ] Product version history is maintained
- [ ] Only Active products are visible to customers

## Traceability

**SRS Reference:** FG-003 — 4.3 Product Management & Catalog  
**Functional Requirements:** FR-021, FR-022, FR-023, FR-024, FR-025, FR-026, FR-027, FR-028, FR-029

[[[PAGEBREAK]]]


---
[[[PAGEBREAK]]]


# Feature Group: Policy Purchase & Issuance (FG-004)

## Business Objective

Enable end-to-end policy purchase from quote to digital issuance. Support customer self-service and agent-assisted flows. Ensure regulatory compliance with mandatory disclosures and nominee management.

**Business Value:**
- Minimize time-to-policy (target: <5 minutes from quote to issuance)
- Digital-first reduces operational cost (no paper, no manual keying)
- Regulatory compliance built-in (disclosures, nominee validation)
- Real-time policy document generation and delivery

## Actors & Portals

| Actor | Portal(s) | Primary Use Cases |
|-------|-----------|-------------------|
| Customer | Mobile App, Web PWA | Self-service purchase, nominee setup |
| Agent | Agent Mobile App | Assisted purchase with customer consent |
| Partner Admin | Partner Portal | Monitor conversion rates |
| Business Admin | Admin Portal | Override pricing, manual policy issuance |

## User Stories

### US-FG-004-01: End-to-End Policy Purchase

**As a** customer  
**I want** to complete the entire purchase in one session  
**So that** I get instant coverage

**Acceptance Criteria:**
- Customer selects product → enters details (applicant, insured, nominee)
- System validates nominee (shares sum to 100%, relationship valid)
- Premium calculated and shown
- Customer agrees to T&C (mandatory disclosure shown)
- Payment initiated → policy issued on payment confirmation
- Digital policy document delivered via SMS/email

![Policy Purchase Flow](images/flow_policy_purchase_e2e.png)

**Related FRs:** FR-030, FR-031, FR-032, FR-033, FR-034, FR-035, FR-036, FR-037

### US-FG-004-02: Nominee Management

**As a** customer  
**I want** to designate one or more nominees  
**So that** my beneficiaries receive claims if needed

**Acceptance Criteria:**
- Customer can add 1-5 nominees (FR-032)
- Each nominee: name, relationship, share percentage, contact
- System validates shares sum to 100%
- Nominee changes require approval (endorsement flow)

![Nominee Setup](images/flow_nominee_management.png)

**Related FRs:** FR-032, FR-033

### US-FG-004-03: Digital Policy Document

**As a** customer  
**I want** a verifiable digital policy with QR code  
**So that** I can prove coverage anytime

**Acceptance Criteria:**
- Policy document generated as PDF
- Contains QR code linking to verification portal
- Sent via SMS (link) + email (attachment)
- Accessible in app under "My Policies"

![Digital Policy Sample](images/document_digital_policy_sample.png)

**Related FRs:** FR-035, FR-036, FR-037

## Business Rules

| Rule ID | Description |
|---------|-------------|
| BR-POL-01 | Nominee shares must sum to 100% |
| BR-POL-02 | Maximum 5 nominees per policy |
| BR-POL-03 | Policy issued only after payment confirmation |
| BR-POL-04 | Digital policy sent within 5 minutes of payment |
| BR-POL-05 | Policy status: Pending Payment → Active (on payment) |

## Key Workflows

### Self-Service Purchase
1. Customer selects product → "Get Quote"
2. Enters applicant/insured details
3. Adds nominee(s)
4. Reviews premium and T&C
5. Initiates payment
6. Payment confirmed → policy issued → document delivered

### Agent-Assisted Purchase
1. Agent opens customer purchase flow (with customer consent)
2. Agent enters customer details on their behalf
3. Customer reviews and approves via OTP
4. Payment and issuance as above

## Data Model Notes

**Policy Entity (SRS Proto: insuretech.policy.entity.v1.Policy)**
- policy_id
- policy_number (unique, human-readable)
- product_id
- customer_id
- applicant (details)
- insured (details)
- nominees (list)
- coverage_amount
- premium
- start_date, end_date
- status (PENDING_PAYMENT, ACTIVE, SUSPENDED, CANCELLED, LAPSED, EXPIRED)
- digital_document_url

## Integration Touchpoints

| External System | Purpose | Failure Handling |
|----------------|---------|------------------|
| Payment Gateway | Premium collection | Retry, manual verification queue |
| Document Generator | PDF creation | Queue for retry, manual generation |
| SMS/Email Gateway | Policy delivery | Retry, customer can download from app |

## Security & Privacy

- Policy documents contain PII → access controlled
- QR codes use signed tokens for verification
- All purchase actions audited (who, when, what)

## NFR Constraints

| NFR | Target |
|-----|--------|
| Policy Issuance Time | <5 minutes from payment to document delivery |
| Document Generation | <30s per policy |
| Availability | 99.9% (purchase downtime = revenue loss) |

## Acceptance Criteria

- [ ] Customer can complete self-service purchase end-to-end
- [ ] Nominee validation works correctly
- [ ] Digital policy is generated and delivered within SLA
- [ ] QR code verification works
- [ ] Agent-assisted flow requires customer consent

## Traceability

**SRS Reference:** FG-004 — 4.4 Policy Lifecycle Management  
**Functional Requirements:** FR-030, FR-031, FR-032, FR-033, FR-034, FR-035, FR-036, FR-037, FR-038, FR-039, FR-040

[[[PAGEBREAK]]]


---
[[[PAGEBREAK]]]


# Feature Group: Policy Renewals (FG-005)

## Business Objective

Enable seamless policy renewals with minimal customer effort, automated reminders, and flexible payment options. Support both manual and auto-renewal flows while maintaining regulatory compliance and customer consent.

**Business Value:**
- Increase retention rate (target: >85% renewal rate)
- Reduce lapse due to forgotten renewals (automated reminders)
- Lower operational cost (auto-renewal vs manual processing)
- Improve cash flow predictability

## Actors & Portals

| Actor | Portal(s) | Responsibilities |
|-------|-----------|------------------|
| Customer | Mobile App, Web | Renew policies manually, opt-in to auto-renewal, manage payment methods |
| System (Automated) | Background Jobs | Send renewal reminders, execute auto-renewals |
| Agent | Agent Mobile App | Assist customers with renewals |
| Business Admin | Admin Portal | Monitor renewal rates, configure reminder schedules |

## User Stories

### US-FG-005-01: One-Click Manual Renewal

**As a** customer  
**I want** to renew my policy with one click without re-entering all my details  
**So that** renewal is fast and I don't lose coverage

**Acceptance Criteria:**
- 30 days before expiry, customer sees "Renew Now" in app
- Tapping "Renew" → shows current policy details (pre-filled)
- Customer can update payment method if needed
- Customer confirms → payment initiated → policy renewed for next term
- New policy document generated and delivered

![One-Click Renewal Flow](images/flow_renewal_one_click.png)

**Related FRs:** FR-084, FR-085, FR-086, FR-087

### US-FG-005-02: Renewal Reminder Cadence

**As a** customer  
**I want** reminders before my policy expires  
**So that** I don't accidentally lose coverage

**Acceptance Criteria:**
- Reminders sent at: 30 days, 15 days, 7 days, 1 day before expiry
- Channels: SMS, email, push notification
- Each reminder includes: policy number, expiry date, renewal link
- Customer can snooze reminders (future enhancement)

**Related FRs:** FR-086, FR-087

### US-FG-005-03: Auto-Renewal Opt-In

**As a** customer  
**I want** my policy to renew automatically  
**So that** I never have a coverage gap

**Acceptance Criteria:**
- Customer opts in during purchase or policy management
- Must provide consent explicitly (checkbox + T&C acceptance)
- Customer saves payment method (tokenized card/MFS)
- System attempts auto-renewal 7 days before expiry
- If payment succeeds → policy renewed + notification
- If payment fails → retry 3x, then notify customer to pay manually

![Auto-Renewal Flow](images/flow_renewal_auto.png)

**Related FRs:** FR-088, FR-089, FR-090

### US-FG-005-04: Auto-Renewal Cancellation

**As a** customer  
**I want** to cancel auto-renewal anytime  
**So that** I control whether my policy continues

**Acceptance Criteria:**
- Customer navigates to policy → "Manage Auto-Renewal"
- Toggles auto-renewal OFF
- Receives confirmation notification
- System will not charge next term
- Customer receives expiry reminder as usual (manual renewal option)

**Related FRs:** FR-091, FR-092

## Business Rules

| Rule ID | Description |
|---------|-------------|
| BR-REN-01 | Renewal reminders start 30 days before expiry |
| BR-REN-02 | Auto-renewal requires explicit customer consent |
| BR-REN-03 | Auto-renewal payment attempted 7 days before expiry |
| BR-REN-04 | Failed auto-renewal retries: 3 attempts over 5 days |
| BR-REN-05 | Customer can cancel auto-renewal anytime before renewal date |
| BR-REN-06 | Manual renewal allowed up to 30 days after expiry (grace period, if product allows) |

## Key Workflows

### Manual Renewal Flow
1. System identifies policies expiring within 30 days
2. Sends renewal reminder to customer
3. Customer taps "Renew Now" in app/email link
4. Pre-filled renewal form shown (can update payment method)
5. Customer confirms → payment processed
6. Policy renewed → new term starts → document delivered

### Auto-Renewal Flow
1. System identifies auto-renewal policies 7 days before expiry
2. Attempts payment via saved payment method
3. **If successful:** Policy renewed, customer notified
4. **If failed:** Retry after 24 hours (max 3 retries)
5. **If all retries fail:** Customer notified to pay manually, policy lapses if no action

### Grace Period (if applicable)
1. Policy expires but product allows grace period (e.g., 30 days)
2. Customer can still renew during grace period
3. Coverage may be suspended during grace (product-specific)

## Data Model Notes

**Renewal Configuration (per Policy)**
- auto_renew_enabled (boolean)
- renewal_consent_timestamp
- payment_method_token
- renewal_reminder_sent (dates)

**Renewal Transaction**
- renewal_id
- policy_id
- renewal_type (MANUAL, AUTO)
- payment_transaction_id
- status (PENDING, SUCCESS, FAILED)

## Integration Touchpoints

| External System | Purpose | Failure Handling |
|----------------|---------|------------------|
| Payment Gateway | Process renewal payments | Retry logic, notify customer |
| SMS/Email Gateway | Renewal reminders | Queue for retry |
| Document Generator | New policy document | Retry, customer can download from app |

## Security & Privacy

- Auto-renewal consent recorded with timestamp
- Payment tokens stored securely (PCI-DSS compliant)
- Customer can view auto-renewal status and cancel anytime
- All renewal actions audited

## NFR Constraints

| NFR | Target |
|-----|--------|
| Reminder Delivery | 95% success rate, <1 hour from trigger |
| Auto-Renewal Processing Time | <5 minutes per policy |
| Payment Retry Interval | 24 hours between retries |

## Acceptance Criteria

- [ ] Customer receives renewal reminders at defined intervals
- [ ] One-click renewal works without re-entering details
- [ ] Auto-renewal opt-in requires explicit consent
- [ ] Auto-renewal payment failure triggers retries and customer notification
- [ ] Customer can cancel auto-renewal anytime
- [ ] All renewal transactions are auditable

## Traceability

**SRS Reference:** FG-005 — 4.5 Policy Management & Renewals  
**Functional Requirements:** FR-084, FR-085, FR-086, FR-087, FR-088, FR-089, FR-090, FR-091, FR-092, FR-093, FR-094, FR-095, FR-096, FR-097, FR-098, FR-099, FR-100, FR-101, FR-102

[[[PAGEBREAK]]]


---
[[[PAGEBREAK]]]


# Feature Group: Policy Endorsements & Cancellations (FG-06)

## Business Objective

Enable customers to update policy details (endorsements/amendments) and cancel policies when needed, with transparent refund calculations and regulatory-compliant approval workflows.

**Business Value:**
- Customer flexibility improves satisfaction and retention
- Controlled approval workflows prevent fraud
- Transparent refund calculations reduce disputes
- Regulatory compliance (refund rules, audit trails)

## Actors & Portals

| Actor | Portal(s) | Responsibilities |
|-------|-----------|------------------|
| Customer | Mobile App, Web | Request endorsements, cancel policies, view refund estimates |
| Business Admin | Admin Portal | Approve/reject endorsement requests, monitor cancellation patterns |
| Focal Person | Admin Portal | Override approvals (exceptional cases) |
| System (Automated) | Background Jobs | Calculate pro-rata refunds, process approvals |

## User Stories

### US-FG-06-01: Policy Endorsement (Amendments)

**As a** customer  
**I want** to update my policy details (address, nominee, coverage)  
**So that** my policy reflects my current situation

**Acceptance Criteria:**
- Customer navigates to policy → "Request Changes"
- Selects change type: address, nominee, coverage amount, beneficiary details
- Some changes require approval (e.g., coverage increase, nominee change)
- Other changes auto-approved (e.g., address update)
- Approval workflow triggered if required
- Customer notified of approval/rejection
- Approved changes reflected in policy + new endorsement document generated

![Endorsement Request Flow](images/flow_endorsement_request.png)

**Related FRs:** 

### US-FG-06-02: Policy Cancellation with Refund

**As a** customer  
**I want** to cancel my policy and receive a refund  
**So that** I can discontinue coverage I no longer need

**Acceptance Criteria:**
- Customer navigates to policy → "Cancel Policy"
- System shows refund estimate (pro-rata calculation based on unused term)
- Customer selects cancellation reason (dropdown)
- Customer confirms cancellation
- Approval workflow triggered (Business Admin approves)
- On approval: policy status → CANCELLED, refund initiated
- Refund processed to original payment method or customer-selected account
- Cancellation confirmation + refund receipt delivered

![Cancellation and Refund Flow](images/flow_cancellation_refund.png)

**Related FRs:** 

### US-FG-06-03: Transparent Refund Calculation

**As a** customer  
**I want** to see exactly how my refund is calculated  
**So that** I trust the amount I'm receiving

**Acceptance Criteria:**
- Refund calculation shown before cancellation confirmation
- Formula displayed: `Refund = Premium × (Unused Days / Total Days) - Admin Fee`
- Breakdown shows: premium paid, days used, days remaining, admin fee, final refund
- Calculation complies with product rules and regulatory guidelines

**Related FRs:** FR-095, FR-096

## Business Rules

| Rule ID | Description |
|---------|-------------|
| BR-END-01 | Address/contact updates: auto-approved |
| BR-END-02 | Nominee changes: require Business Admin approval |
| BR-END-03 | Coverage increase: requires underwriting review (future: AI-assisted) |
| BR-CAN-01 | Cancellation refund: pro-rata based on unused term minus admin fee |
| BR-CAN-02 | Cancellation within 15 days of purchase: full refund (cooling-off period) |
| BR-CAN-03 | Cancellations require Business Admin approval for policies >100k BDT |
| BR-CAN-04 | Refunds processed within 7 business days of approval |

## Key Workflows

### Endorsement Flow (Approval Required)
1. Customer submits endorsement request
2. System validates request (e.g., nominee shares sum to 100%)
3. Request routed to Business Admin approval queue
4. Business Admin reviews → approves/rejects
5. **If approved:** Policy updated, endorsement document generated, customer notified
6. **If rejected:** Customer notified with reason

### Endorsement Flow (Auto-Approved)
1. Customer submits low-risk change (address, phone)
2. System validates and auto-approves
3. Policy updated immediately, customer notified

### Cancellation Flow
1. Customer requests cancellation → selects reason
2. System calculates pro-rata refund and shows estimate
3. Customer confirms
4. Request routed to approval (if high-value policy)
5. Business Admin approves
6. Policy status → CANCELLED
7. Refund initiated via payment gateway
8. Customer receives confirmation + refund receipt

## Data Model Notes

**Endorsement Entity**
- endorsement_id
- policy_id
- change_type (ADDRESS, NOMINEE, COVERAGE, etc.)
- old_value, new_value
- status (PENDING, APPROVED, REJECTED)
- requested_by, approved_by, timestamps

**Cancellation Entity**
- cancellation_id
- policy_id
- reason
- refund_amount
- refund_calculation_breakdown (JSON)
- status (PENDING, APPROVED, REFUNDED)

## Integration Touchpoints

| External System | Purpose | Failure Handling |
|----------------|---------|------------------|
| Payment Gateway | Process refunds | Retry logic, manual refund queue |
| Document Generator | Generate endorsement documents | Retry, customer can download from app |

## Security & Privacy

- All endorsement/cancellation requests logged with actor and timestamp
- Refund calculations auditable (stored with breakdown)
- High-value cancellations require dual approval (future enhancement)

## NFR Constraints

| NFR | Target |
|-----|--------|
| Endorsement Approval TAT | <24 hours for standard changes |
| Cancellation Approval TAT | <48 hours |
| Refund Processing Time | <7 business days |

## Acceptance Criteria

- [ ] Customer can request policy endorsements via app
- [ ] Approval workflows route correctly based on change type
- [ ] Cancellation refund calculation is transparent and accurate
- [ ] Refunds are processed within SLA
- [ ] All changes are auditable and traceable

## Traceability

**SRS Reference:** FG-06 — 4.6 Business Rules & Workflows  
**Functional Requirements:** FR-214, FR-215, FR-216, FR-217, FR-218, FR-219, FR-220, FR-221, FR-222

[[[PAGEBREAK]]]


---
[[[PAGEBREAK]]]


# Feature Group: Premium Collection & Payment Management (FG-007)

## Business Objective

Enable seamless premium collection through multiple payment channels (MFS, bank, card) with robust fallback mechanisms, receipt generation, and reconciliation. Minimize payment failures and revenue leakage.

**Business Value:**
- Maximize payment success rate (target: >95%)
- Support Bangladesh-preferred payment methods (bKash, Nagad, Rocket priority)
- Reduce manual intervention through automated reconciliation
- Provide audit-ready payment records for compliance

## Actors & Portals

| Actor | Portal(s) | Responsibilities |
|-------|-----------|------------------|
| Customer | Mobile App, Web | Pay premiums, view payment history, download receipts |
| Agent | Agent Mobile App | Initiate payments on behalf of customer (with consent) |
| Business Admin | Admin Portal | Verify manual payments, reconcile transactions, dispute resolution |
| Finance Team | Admin Portal | Generate balance sheets, accounting reports |

## User Stories

### US-FG-007-01: Multi-Channel Premium Payment

**As a** customer  
**I want** to pay using my preferred payment method  
**So that** payment is convenient and successful

**Acceptance Criteria:**
- Payment options shown: bKash, Nagad, Rocket, bank transfer, credit/debit card
- Customer selects method → redirected to payment gateway
- Payment confirmed → policy issued/renewed
- Payment receipt generated automatically
- Payment status visible in app ("My Policies" → payment history)

![Multi-Channel Payment](images/flow_payment_channels.png)

**Related FRs:** FR-073, FR-074, FR-075, FR-076, FR-077

### US-FG-007-02: Payment Confirmation via Webhook

**As a** system  
**I want** to receive real-time payment confirmation from the gateway  
**So that** policies are issued immediately

**Acceptance Criteria:**
- Payment gateway sends webhook on payment success/failure
- System validates webhook signature (security)
- On success: policy status updated, document generated
- On failure: customer notified to retry
- Webhook failures trigger fallback: poll gateway API every 60s (max 10 attempts)

**Related FRs:** FR-075, FR-079

### US-FG-007-03: Manual Payment Verification

**As a** customer in a remote area  
**I want** to pay via bank deposit and upload proof  
**So that** I can still purchase even if digital payment fails

**Acceptance Criteria:**
- Customer selects "Bank Deposit" or "Manual Payment"
- Uploads payment proof (bank slip photo, transaction screenshot)
- Business Admin receives verification request
- Admin verifies payment → marks as confirmed → policy issued
- If rejected → customer notified with reason

![Manual Payment Verification](images/flow_payment_manual.png)

**Related FRs:** FR-076, FR-077

### US-FG-007-04: Payment Receipt and History

**As a** customer  
**I want** a receipt for every payment  
**So that** I have records for disputes and accounting

**Acceptance Criteria:**
- Receipt generated immediately after payment confirmation
- Receipt includes: transaction ID, policy number, amount, date, payment method, merchant ref
- Receipt accessible in app ("Payments" tab) and sent via email
- Receipt downloadable as PDF

![Payment Receipt](images/document_payment_receipt.png)

**Related FRs:** FR-078

## Business Rules

| Rule ID | Description |
|---------|-------------|
| BR-PAY-01 | Payment confirmation required before policy issuance/renewal |
| BR-PAY-02 | Webhook signature validation mandatory (prevent fraud) |
| BR-PAY-03 | Manual payments verified within 24 hours (business days) |
| BR-PAY-04 | Payment retries: max 3 attempts for failed payments |
| BR-PAY-05 | Receipts must include all regulatory-required fields |
| BR-PAY-06 | Payment reconciliation runs daily (match gateway records vs internal) |

## Key Workflows

### Digital Payment Flow (MFS/Card)
1. Customer initiates payment → selects payment method
2. System generates payment order (amount, reference)
3. Customer redirected to gateway (bKash/Nagad/etc.)
4. Customer completes payment in gateway app/web
5. Gateway sends webhook → system validates and records
6. Policy issued → receipt generated → customer notified

### Manual Payment Flow
1. Customer selects "Manual Payment" → sees bank account details
2. Customer deposits money and uploads proof
3. Business Admin sees verification queue
4. Admin verifies → marks as confirmed (or rejects)
5. On confirm: policy issued → customer notified

### Payment Failure Handling
1. Gateway returns failure (insufficient funds, timeout, etc.)
2. Customer notified with error message
3. "Retry Payment" option shown
4. System tracks retry attempts (max 3)
5. If all retries fail → policy remains pending, customer can pay later

### Reconciliation (Daily Job)
1. System fetches gateway transaction report
2. Matches internal payment records vs gateway records
3. Flags mismatches (missing confirmations, duplicate payments)
4. Finance team reviews and resolves discrepancies

## Data Model Notes

**Payment Entity (SRS Proto: insuretech.payment.entity.v1.Payment)**
- payment_id
- policy_id
- customer_id
- amount
- currency (BDT)
- payment_method (BKASH, NAGAD, ROCKET, BANK_TRANSFER, CARD)
- status (PENDING, SUCCESS, FAILED, REFUNDED)
- gateway_transaction_id
- gateway_response (JSON)
- receipt_url
- timestamps (created, confirmed)

**Manual Payment Verification**
- verification_id
- payment_id
- proof_document_url
- verified_by (admin user_id)
- verification_status (PENDING, APPROVED, REJECTED)

## Integration Touchpoints

| External System | Purpose | Failure Handling |
|----------------|---------|------------------|
| bKash/Nagad/Rocket APIs | MFS payments | Retry webhook poll, manual verification queue |
| Bank Payment Gateway | Card/bank transfers | Retry, fallback to manual |
| SMS/Email Gateway | Payment receipts | Queue for retry |

## Security & Privacy

- Webhook signature validation (HMAC/RSA) prevents payment spoofing
- Payment card details never stored (PCI-DSS: use gateway tokenization)
- Payment receipts contain no sensitive card/account info
- All payment transactions logged for audit

## NFR Constraints

| NFR | Target |
|-----|--------|
| Payment Success Rate | >95% (excluding customer-side failures) |
| Webhook Processing Time | <5s from gateway callback to policy issuance |
| Payment Receipt Generation | <10s |
| Manual Payment Verification TAT | <24 hours (business days) |

## Acceptance Criteria

- [ ] Customer can pay using bKash, Nagad, Rocket, card, bank transfer
- [ ] Payment confirmation received via webhook or polling
- [ ] Manual payment verification workflow functions correctly
- [ ] Payment receipts generated and delivered
- [ ] Payment history accessible in customer app
- [ ] Reconciliation identifies and flags mismatches

## Traceability

**SRS Reference:** FG-007 — 4.7 Payment Processing  
**Functional Requirements:** FR-073, FR-074, FR-075, FR-076, FR-077, FR-078, FR-079, FR-080, FR-081, FR-082, FR-083

[[[PAGEBREAK]]]


---
[[[PAGEBREAK]]]


# Feature Group: Claims Management (FG-008)

## Business Objective

Enable customers to submit claims digitally, track status transparently, and receive timely settlements. Implement governance controls (tiered approvals, fraud checks) while maintaining customer trust.

**Business Value:**
- Reduce claims turnaround time (TAT) — target: <7 days for simple claims
- Minimize fraud through automated eligibility checks and fraud rules
- Transparency builds customer trust and retention
- Auditability for regulatory compliance

## Actors & Portals

| Actor | Portal(s) | Responsibilities |
|-------|-----------|------------------|
| Customer | Mobile App, Web | Submit claim, upload documents, track status |
| Agent | Agent Mobile App | Assist customer with claim submission |
| Claims Adjuster | Admin Portal | Review, approve/reject claims |
| Focal Person | Admin Portal | Second-level approval (high-value claims) |
| Business Admin | Admin Portal | Override, dispute resolution |

## User Stories

### US-FG-008-01: Digital Claims Submission

**As a** customer  
**I want** to submit a claim from my phone with photos  
**So that** I don't have to visit a branch

**Acceptance Criteria:**
- Customer selects policy → "File Claim"
- Enters incident details (date, type, amount requested)
- Uploads required documents (ID, bills, photos) — max 10MB per file
- System validates eligibility (policy active, within coverage, no duplicate)
- Claim submitted → status: SUBMITTED → UNDER_REVIEW

![Claims Submission Flow](images/flow_claims_submission.png)

**Related FRs:** FR-041, FR-042, FR-043, FR-044, FR-045

### US-FG-008-02: Transparent Claims Tracking

**As a** customer  
**I want** to see my claim status in real-time  
**So that** I know what's happening and when I'll get paid

**Acceptance Criteria:**
- Dashboard shows claim status (Submitted, Under Review, Approved, Settled, Rejected)
- Timeline shows each status change with timestamp
- Notifications sent at each milestone
- If additional documents needed, customer notified with clear instructions

![Claims Dashboard - Customer](images/dashboard_claims_customer.png)

**Related FRs:** FR-047, FR-048

### US-FG-008-03: Tiered Approval Workflow

**As a** Business Admin  
**I want** claims routed to the right approvers based on amount  
**So that** high-value claims have appropriate oversight

**Acceptance Criteria:**
- Claims <50k BDT → single approver
- Claims 50k-200k BDT → two approvers (tiered)
- Claims >200k BDT → joint approval (two approvers must approve simultaneously)
- Approval matrix configurable by Business Admin

**Related FRs:** FR-046, FR-051, FR-219

### US-FG-008-04: Claims Settlement

**As a** customer  
**I want** my approved claim paid to my chosen account  
**So that** I receive funds quickly

**Acceptance Criteria:**
- Customer selects payout method during claim submission (bank transfer, MFS)
- On approval, settlement initiated automatically
- Customer receives payment confirmation notification
- Settlement recorded in claim history and balance sheet

![Settlement Flow](images/flow_claims_settlement.png)

**Related FRs:** FR-052

## Business Rules

| Rule ID | Description |
|---------|-------------|
| BR-CLM-01 | Claims can only be filed for ACTIVE policies |
| BR-CLM-02 | Duplicate claims (same policy + incident date) rejected |
| BR-CLM-03 | Document uploads: max 10MB per file, supported types: JPG, PNG, PDF |
| BR-CLM-04 | Tiered approvals: <50k (1 approver), 50k-200k (2 tier), >200k (joint) |
| BR-CLM-05 | Fraud flags trigger manual review before approval |
| BR-CLM-06 | SLA: simple claims reviewed within 48 hours |

## Key Workflows

### Claims Submission → Settlement
1. Customer submits claim with documents
2. System validates eligibility (policy status, coverage, duplicates)
3. Fraud detection rules run → flags if suspicious
4. Claim enters review queue (routed by approval matrix)
5. Approver(s) review → approve/reject/request more info
6. If approved → settlement initiated
7. Customer receives payment + notification
8. Claim status: SETTLED

### Exception Handling
- **Documents missing:** status → PENDING_DOCS, customer notified
- **Fraud flag:** status → UNDER_INVESTIGATION, manual review
- **Rejected:** customer notified with reason, appeal option available

## Data Model Notes

**Claim Entity (SRS Proto: insuretech.claims.entity.v1.Claim)**
- claim_id
- policy_id
- customer_id
- incident_date
- claim_type (DEATH, HOSPITALIZATION, ACCIDENT, etc.)
- claim_amount
- documents (list of uploaded files)
- status (SUBMITTED, UNDER_REVIEW, PENDING_DOCS, APPROVED, REJECTED, SETTLED, CANCELLED)
- approvers (list)
- fraud_flags (list)
- settlement_details
- audit_trail

## Integration Touchpoints

| External System | Purpose | Failure Handling |
|----------------|---------|------------------|
| Payment Gateway | Settlement disbursement | Retry, manual processing |
| Document Storage (S3/Minio) | Store claim documents | Redundant storage, backup |
| Fraud Detection Service | Run fraud rules | Queue for manual review if service down |
| Hospital/Provider Network | Validate provider (if applicable) | Manual validation queue |

## Security & Privacy

- Claim documents contain sensitive medical/financial info → encrypted at rest
- Access controlled by role (customer sees only their claims)
- All approvals and rejections logged with reason
- PII redacted in analytics reports

## NFR Constraints

| NFR | Target |
|-----|--------|
| Claims Submission Uptime | 99.9% |
| Document Upload | <10s for 5MB file |
| Fraud Check Latency | <2s per claim |
| Simple Claims TAT | <7 days (target: 48 hours for review) |

## Acceptance Criteria

- [ ] Customer can submit claim with documents from mobile app
- [ ] Eligibility checks prevent invalid claims
- [ ] Fraud detection flags suspicious claims for review
- [ ] Tiered approval workflow routes claims correctly
- [ ] Approved claims settle within SLA
- [ ] Customer can track status in real-time

## Traceability

**SRS Reference:** FG-008 — 4.8 Claims Management  
**Functional Requirements:** FR-041, FR-042, FR-043, FR-044, FR-045, FR-046, FR-047, FR-048, FR-049, FR-050, FR-051, FR-052, FR-053, FR-054, FR-055, FR-056, FR-057, FR-058, FR-103, FR-104, FR-105

[[[PAGEBREAK]]]


---
[[[PAGEBREAK]]]


# Feature Group: Partner Management (FG-009)

## Business Objective

Enable partner onboarding, verification (KYB), tenant isolation, and performance monitoring. Support scalable distribution via MFS providers, hospitals, e-commerce platforms, and agent organizations.

**Business Value:**
- Scale distribution without proportional operational cost increase
- Ensure partner data isolation (regulatory and competitive requirement)
- Enable commission-based revenue sharing with transparent tracking
- Monitor partner performance and fraud patterns

## Actors & Portals

| Actor | Portal(s) | Responsibilities |
|-------|-----------|------------------|
| Partner Admin | Partner Portal | Submit KYB documents, manage agents, view commissions, assisted sales |
| Focal Person | Admin Portal | Verify KYB, approve/reject partner applications, monitor partner compliance |
| Business Admin | Admin Portal | Configure partner commission rules, dispute resolution |
| Agent | Agent Mobile App | Sell under partner umbrella |

## User Stories

### US-FG-009-01: Partner Onboarding & KYB

**As a** potential partner organization  
**I want** to apply for platform access by submitting business verification documents  
**So that** I can start selling insurance products

**Acceptance Criteria:**
- Partner submits application form: business name, type, registration number, address, contact
- Uploads KYB documents: business license, tax certificate, bank statement, director IDs
- Application routed to Focal Person for verification
- Focal Person reviews documents → approves/rejects
- **If approved:** tenant_id assigned, Partner Admin account created, credentials sent
- **If rejected:** Partner notified with reason

![Partner Onboarding Flow](images/flow_partner_onboarding.png)

**Related FRs:** FR-059, FR-061, FR-066, FR-070

### US-FG-009-02: Tenant Isolation & Data Security

**As a** Focal Person  
**I want** each partner to see only their own data  
**So that** competitive information is protected

**Acceptance Criteria:**
- Every partner has unique tenant_id
- All API queries filtered by tenant_id
- Partner A cannot access Partner B's customers, policies, agents, or commissions
- Database queries enforce tenant_id in WHERE clause
- Automated tests validate tenant isolation

**Related FRs:** FR-014, FR-015, FR-016 (from FG-002, but critical for partner management)

### US-FG-009-03: Partner Performance Dashboard

**As a** Partner Admin  
**I want** a dashboard showing my sales, conversion rates, and commissions  
**So that** I can optimize my distribution strategy

**Acceptance Criteria:**
- Dashboard shows: leads, quotes, policies issued, conversion rate, commission earned (pending/paid)
- Filterable by date range, product, agent
- Exportable as CSV/PDF for accounting
- Real-time or near-real-time updates

![Partner Dashboard](images/dashboard_partner_performance.png)

**Related FRs:** FR-062

### US-FG-009-04: Partner Suspension & Reactivation

**As a** Focal Person  
**I want** to suspend a partner for policy violations  
**So that** compliance is enforced

**Acceptance Criteria:**
- Focal Person navigates to partner profile → "Suspend"
- Selects reason (fraud, non-compliance, contract breach)
- On suspension: partner loses access, agents cannot sell, customers can still claim/renew
- Partner notified of suspension reason
- Focal Person can reactivate after resolution

**Related FRs:** FR-064, FR-065

## Business Rules

| Rule ID | Description |
|---------|-------------|
| BR-PTR-01 | Each partner assigned unique tenant_id (never reused) |
| BR-PTR-02 | KYB documents required: business license, tax cert, bank statement, director IDs |
| BR-PTR-03 | Focal Person must approve partner within 72 hours of submission |
| BR-PTR-04 | Suspended partners: agents cannot sell, but customer servicing (claims/renewals) continues |
| BR-PTR-05 | Partner data isolated at API, service, and database levels |

## Key Workflows

### Partner Onboarding
1. Partner submits application + KYB documents via web form
2. System validates completeness
3. Application routed to Focal Person queue
4. Focal Person verifies documents (may request additional docs)
5. Focal Person approves → system creates tenant_id + Partner Admin account
6. Partner receives credentials + onboarding guide
7. Partner can now create agents and start selling

### Partner Performance Monitoring
1. Business Admin reviews partner dashboards monthly
2. Identifies low performers or fraud patterns
3. Engages with partner for improvement or investigation
4. May suspend partner if violations found

### Partner Suspension
1. Focal Person initiates suspension (reason required)
2. System immediately revokes partner/agent access
3. Customer-facing operations (renewals, claims) remain available
4. Partner notified with appeal process
5. After resolution, Focal Person can reactivate

## Data Model Notes

**Partner Entity (SRS Proto: insuretech.partner.entity.v1.Partner)**
- partner_id (tenant_id)
- partner_name
- partner_type (MFS, HOSPITAL, ECOMMERCE, AGENT_ORG)
- registration_number
- kyb_documents (list)
- status (PENDING_APPROVAL, ACTIVE, SUSPENDED, TERMINATED)
- focal_person_id (assigned verifier)
- commission_config
- onboarded_at, verified_at

## Integration Touchpoints

| External System | Purpose | Failure Handling |
|----------------|---------|------------------|
| Partner APIs (MFS, EHR) | Data exchange for embedded insurance | Circuit breaker, manual queue |
| KYB Verification Service (future) | Automated business verification | Manual review if API unavailable |

## Security & Privacy

- Partner KYB documents encrypted at rest
- Tenant isolation validated in security testing
- Partner access logs audited
- Focal Person actions (approve/suspend) logged

## NFR Constraints

| NFR | Target |
|-----|--------|
| KYB Verification TAT | <72 hours |
| Partner Dashboard Load Time | <3s |
| Tenant Isolation Test Coverage | 100% of multi-tenant APIs |

## Acceptance Criteria

- [ ] Partner can submit onboarding application with KYB documents
- [ ] Focal Person can verify and approve/reject partners
- [ ] Tenant isolation prevents cross-partner data access
- [ ] Partner dashboard shows accurate performance metrics
- [ ] Suspended partners lose access immediately

## Traceability

**SRS Reference:** FG-009 — 4.9 Partner & Agent Management  
**Functional Requirements:** FR-059, FR-060, FR-061, FR-062, FR-063, FR-064, FR-065, FR-066, FR-067, FR-068, FR-069, FR-070, FR-071, FR-072

[[[PAGEBREAK]]]


---
[[[PAGEBREAK]]]


# Feature Group: Agent & Sub-Partner Management (FG-010)

## Business Objective

Enable partners to manage agent hierarchies, track agent performance, and ensure agents operate under partner governance. Support commission allocation and agent-assisted sales with customer consent.

**Business Value:**
- Scale distribution via agent networks (last-mile reach)
- Enable performance-based incentives for agents
- Maintain accountability (agents tied to partners)
- Support rural/low-literacy customer segments (agent-assisted)

## Actors & Portals

| Actor | Portal(s) | Responsibilities |
|-------|-----------|------------------|
| Partner Admin | Partner Portal | Create/manage agents, assign territories, set commission splits |
| Agent | Agent Mobile App | Assisted sales, customer onboarding, lead management |
| Customer | Mobile App (indirect) | Consent to agent-assisted purchase |
| Business Admin | Admin Portal | Monitor agent compliance, resolve disputes |

## User Stories

### US-FG-010-01: Agent Onboarding by Partner

**As a** Partner Admin  
**I want** to create agent accounts under my partner organization  
**So that** agents can sell on my behalf

**Acceptance Criteria:**
- Partner Admin navigates to "Agents" → "Add New Agent"
- Enters agent details: name, mobile, NID, territory (optional)
- Assigns commission split (partner vs agent)
- Agent receives credentials via SMS
- Agent can now login to Agent Mobile App

![Agent Onboarding](images/flow_agent_onboarding.png)

**Related FRs:** 

### US-FG-010-02: Agent-Assisted Customer Onboarding

**As an** agent  
**I want** to onboard a customer on their behalf  
**So that** I can help customers with low digital literacy

**Acceptance Criteria:**
- Agent opens "New Customer" flow in Agent App
- Enters customer details (with customer present)
- Sends OTP to customer's mobile for consent
- Customer confirms OTP → account created under agent linkage
- Agent can now guide customer through product selection and purchase

**Related FRs:** FR-067, FR-068

### US-FG-010-03: Agent Commission Tracking

**As an** agent  
**I want** to see my earned commissions  
**So that** I can track my income

**Acceptance Criteria:**
- Agent dashboard shows: policies sold, commission earned (pending/paid), payment history
- Filterable by date, product
- Commission calculation transparent (shown as % of premium)
- Agent can export commission statement

![Agent Commission Dashboard](images/dashboard_agent_commission.png)

**Related FRs:** FR-062, FR-063, FR-141..FR-148

### US-FG-010-04: Partner Hierarchy & Commission Splits

**As a** Partner Admin  
**I want** to configure commission splits between my organization and agents  
**So that** revenue sharing is automated

**Acceptance Criteria:**
- Partner Admin sets default commission split (e.g., Partner 60%, Agent 40%)
- Can override per agent if needed
- Commission calculated automatically on policy issuance
- Both partner and agent see their respective shares in dashboards

**Related FRs:** FR-143, FR-144

## Business Rules

| Rule ID | Description |
|---------|-------------|
| BR-AGT-01 | Agents must be linked to a partner (no independent agents) |
| BR-AGT-02 | Agent-assisted sales require customer OTP consent |
| BR-AGT-03 | Commission split configurable per agent (default: partner 60%, agent 40%) |
| BR-AGT-04 | Commission paid after policy premium collected and reconciled |
| BR-AGT-05 | Suspended agents cannot create new sales but existing policies remain valid |

## Key Workflows

### Agent Creation & Activation
1. Partner Admin creates agent account (details + commission split)
2. System generates agent credentials
3. Agent receives SMS with app download link + login credentials
4. Agent logs in, completes onboarding training (optional checklist)
5. Agent status: ACTIVE → can start selling

### Agent-Assisted Sale
1. Agent meets customer → opens Agent App
2. Agent enters customer details, sends OTP to customer mobile
3. Customer confirms OTP (consent)
4. Agent guides customer through product selection
5. Agent initiates purchase (customer confirms payment)
6. Policy issued → commission allocated to agent + partner

### Commission Payout (Monthly/Bi-weekly)
1. System calculates total commissions per agent and partner
2. Finance team reviews and approves payout batch
3. Payments initiated to agent bank/MFS accounts
4. Agents receive payout notification + statement

## Data Model Notes

**Agent Entity**
- agent_id
- partner_id (linked to partner/tenant)
- agent_name
- mobile_number
- nid_number
- territory (optional)
- commission_split (partner_pct, agent_pct)
- status (ACTIVE, SUSPENDED, TERMINATED)
- created_by (partner admin)

**Agent Commission Transaction**
- transaction_id
- agent_id
- policy_id
- premium_amount
- commission_amount
- commission_pct
- status (PENDING, PAID)
- paid_at

## Integration Touchpoints

| External System | Purpose | Failure Handling |
|----------------|---------|------------------|
| Payment Gateway | Commission payouts | Retry, manual payment queue |
| SMS Gateway | Agent credentials delivery | Queue for retry |

## Security & Privacy

- Agent actions (customer onboarding, sales) logged with customer consent (OTP)
- Agent credentials secured (password + biometric)
- Partner Admin can suspend agents immediately
- Commission calculations auditable

## NFR Constraints

| NFR | Target |
|-----|--------|
| Agent App Responsiveness | Works on low-end Android devices (2GB RAM) |
| Commission Calculation Accuracy | 100% (auditable via test cases) |
| Agent Performance Dashboard | <2s load time |

## Acceptance Criteria

- [ ] Partner Admin can create and manage agent accounts
- [ ] Agent-assisted sales require customer OTP consent
- [ ] Agents can track their commissions in real-time
- [ ] Commission splits are configurable and enforced correctly
- [ ] Suspended agents lose sales access immediately

## Traceability

**SRS Reference:** FG-010 — 4.10 Partner Portal & Business Intelligence  
**Functional Requirements:** FR-141, FR-142, FR-143, FR-144, FR-145, FR-146, FR-147, FR-148, FR-205, FR-206, FR-207, FR-208

[[[PAGEBREAK]]]


---
[[[PAGEBREAK]]]


# Feature Group: Customer Support & Ticketing (FG-011)

## Business Objective

Provide multi-channel customer support with self-service FAQs, ticketing, escalation workflows, and CSAT tracking. Reduce support costs while improving customer satisfaction.

**Business Value:**
- Self-service FAQs deflect >40% of support queries (target)
- Ticketing ensures accountability and SLA tracking
- CSAT feedback drives continuous improvement
- Regulatory compliance (customer complaint handling)

## Actors & Portals

| Actor | Portal(s) | Responsibilities |
|-------|-----------|------------------|
| Customer | Mobile App, Web | Search FAQs, create tickets, track status, provide feedback |
| Support Agent (L1) | Support Portal | Respond to tickets, escalate complex issues |
| Support Manager (L2) | Support Portal | Handle escalations, approve resolutions, review CSAT |
| Business Admin | Admin Portal | Configure FAQs, review support metrics, dispute resolution |

## User Stories

### US-FG-011-01: Self-Service FAQ

**As a** customer  
**I want** to search for answers to common questions  
**So that** I can resolve issues quickly without waiting for support

**Acceptance Criteria:**
- FAQ section accessible from app home + help menu
- Search bar with keyword matching (Bengali + English)
- FAQs categorized: Registration, Purchase, Payments, Claims, Policy Management
- Each FAQ has: question, answer, helpful/not helpful voting
- Most helpful FAQs surfaced at top

![FAQ Interface](images/ui_faq_search.png)

**Related FRs:** FR-106

### US-FG-011-02: Create Support Ticket

**As a** customer  
**I want** to create a support ticket when FAQ doesn't help  
**So that** a human agent can assist me

**Acceptance Criteria:**
- Customer taps "Contact Support" → "Create Ticket"
- Selects category: Account, Payment, Policy, Claims, Other
- Enters description (text + optional attachment)
- Ticket created with unique ID
- Customer receives confirmation SMS/email with ticket number
- Customer can track ticket status in app

![Ticket Creation](images/flow_ticket_creation.png)

**Related FRs:** FR-108, FR-109

### US-FG-011-03: Support Agent Response & Escalation

**As a** Support Agent (L1)  
**I want** to respond to customer tickets  
**So that** issues are resolved quickly

**Acceptance Criteria:**
- Agent sees ticket queue in Support Portal (ordered by priority/age)
- Agent can view ticket details, customer history, related policies
- Agent responds via canned responses or custom message
- If issue complex → agent escalates to L2 Support Manager
- Ticket status updates trigger customer notifications

**Related FRs:** FR-110, FR-111

### US-FG-011-04: CSAT Feedback

**As a** customer  
**I want** to rate my support experience  
**So that** the platform improves

**Acceptance Criteria:**
- After ticket resolution, customer receives CSAT survey (1-5 stars + optional comment)
- Survey sent via SMS link or in-app notification
- Customer feedback recorded and visible to Support Manager
- Low CSAT (<3 stars) triggers manager review

![CSAT Survey](images/ui_csat_survey.png)

**Related FRs:** FR-112, FR-113

## Business Rules

| Rule ID | Description |
|---------|-------------|
| BR-SUP-01 | L1 agents must respond to tickets within 4 hours (business hours) |
| BR-SUP-02 | Escalated tickets (L2) must be resolved within 24 hours |
| BR-SUP-03 | Customer notified at each ticket status change (Assigned, In Progress, Resolved, Closed) |
| BR-SUP-04 | CSAT survey sent within 1 hour of ticket resolution |
| BR-SUP-05 | Tickets auto-close after 7 days of inactivity (with customer notification) |

## Key Workflows

### Self-Service Flow
1. Customer searches FAQ → finds answer → issue resolved (no ticket)

### Ticket Creation → Resolution
1. Customer creates ticket (category + description)
2. Ticket auto-assigned to L1 agent (round-robin or skill-based)
3. Agent reviews ticket + customer history
4. Agent responds (resolves directly or escalates to L2)
5. If escalated: L2 manager takes ownership → resolves
6. Ticket marked "Resolved" → customer notified
7. CSAT survey sent → customer provides feedback
8. Ticket auto-closes after 7 days if no further customer response

### Escalation Flow
1. L1 agent determines issue requires L2 (policy dispute, payment issue, claim exception)
2. Agent escalates with notes
3. L2 manager reviews → may involve Business Admin or Finance team
4. L2 resolves or approves exception
5. Customer notified of resolution

## Data Model Notes

**Ticket Entity**
- ticket_id
- customer_id
- category (ACCOUNT, PAYMENT, POLICY, CLAIMS, OTHER)
- description
- attachments (URLs)
- status (OPEN, ASSIGNED, IN_PROGRESS, RESOLVED, CLOSED, ESCALATED)
- priority (LOW, MEDIUM, HIGH, URGENT)
- assigned_to (agent_id)
- created_at, resolved_at

**CSAT Feedback**
- feedback_id
- ticket_id
- rating (1-5 stars)
- comment (optional)
- submitted_at

## Integration Touchpoints

| External System | Purpose | Failure Handling |
|----------------|---------|------------------|
| SMS/Email Gateway | Ticket notifications, CSAT surveys | Queue for retry |
| Knowledge Base (future) | AI-powered FAQ suggestions | Fallback to manual search |

## Security & Privacy

- Tickets contain PII → access controlled by role
- Support agents see only assigned tickets (or queue)
- Support managers can view all tickets
- Ticket history immutable (audit trail)

## NFR Constraints

| NFR | Target |
|-----|--------|
| FAQ Search Response Time | <500ms |
| Ticket Creation | <2s |
| L1 Response SLA | <4 hours (business hours) |
| L2 Resolution SLA | <24 hours |
| CSAT Survey Delivery | <1 hour after resolution |

## Acceptance Criteria

- [ ] Customer can search and view FAQs in Bengali/English
- [ ] Customer can create tickets with attachments
- [ ] Support agents can respond and escalate tickets
- [ ] Ticket status updates trigger customer notifications
- [ ] CSAT feedback collected and tracked

## Traceability

**SRS Reference:** FG-011 — 4.11 Customer Support & Helpdesk  
**Functional Requirements:** FR-106, FR-107, FR-108, FR-109, FR-110, FR-111, FR-112, FR-113

[[[PAGEBREAK]]]


---
[[[PAGEBREAK]]]


# Feature Group: Notifications & Communication (FG-012)

## Business Objective

Deliver timely, relevant notifications via SMS, email, and push channels with customer consent management, anti-spam controls, and regulatory compliance (marketing opt-in/opt-out).

**Business Value:**
- Keep customers informed (OTPs, policy events, claim updates, renewals)
- Drive engagement (renewal reminders, product offers)
- Regulatory compliance (consent for marketing, opt-out mechanism)
- Cost optimization (rate limiting, channel prioritization)

## Actors & Portals

| Actor | Portal(s) | Responsibilities |
|-------|-----------|------------------|
| Customer | Mobile App, Web | Manage notification preferences (consent, channels, frequency) |
| System (Automated) | Background Jobs | Trigger notifications based on events |
| Business Admin | Admin Portal | Configure notification templates, monitor delivery rates |
| Marketing Team | Admin Portal | Send promotional campaigns (with consent enforcement) |

## User Stories

### US-FG-012-01: Transactional Notifications

**As a** customer  
**I want** to receive critical notifications (OTP, purchase, claim updates)  
**So that** I stay informed about important events

**Acceptance Criteria:**
- OTP sent within 60s of request (SMS mandatory)
- Policy purchase confirmation sent via SMS + email
- Claim status updates sent at each milestone (Submitted, Under Review, Approved, Settled)
- Payment receipts sent immediately after successful payment
- Notifications include: timestamp, reference ID, next action (if applicable)

**Related FRs:** FR-115

### US-FG-012-02: Renewal Reminders

**As a** customer  
**I want** reminders before my policy expires  
**So that** I don't lose coverage

**Acceptance Criteria:**
- Reminders sent at: 30 days, 15 days, 7 days, 1 day before expiry
- Multi-channel: SMS (primary), email (secondary), push (if app installed)
- Each reminder includes: policy number, expiry date, renewal link
- Customer can snooze/disable renewal reminders in preferences

![Renewal Reminder Notification](images/notification_renewal_reminder.png)

**Related FRs:** FR-086, FR-087, FR-116

### US-FG-012-03: Marketing Communication with Consent

**As a** customer  
**I want** to control whether I receive promotional messages  
**So that** I'm not spammed

**Acceptance Criteria:**
- During registration, customer opts in/out of marketing (checkbox, default: opt-out)
- Customer can change preference anytime in app settings
- Marketing messages sent only to opted-in customers
- Every marketing message includes unsubscribe link
- Regulatory compliance: record consent timestamp

![Notification Preferences](images/ui_notification_preferences.png)

**Related FRs:** FR-122, FR-123

### US-FG-012-04: Anti-Spam & Rate Limiting

**As a** business owner  
**I want** rate limits to prevent notification abuse  
**So that** customers don't get overwhelmed and SMS costs are controlled

**Acceptance Criteria:**
- OTP rate limit: max 3 per 15 min per user
- Marketing SMS: max 1 per day per user
- Transactional notifications: no limit (critical)
- System monitors notification volume per user, flags anomalies
- Rate limit exceeded → user sees "Try again in X minutes"

**Related FRs:** FR-003 (OTP rate limit), FR-120 (general rate limiting)

## Business Rules

| Rule ID | Description |
|---------|-------------|
| BR-NOT-01 | Transactional notifications (OTP, purchase, claims) always sent (no consent required) |
| BR-NOT-02 | Marketing notifications require explicit opt-in consent |
| BR-NOT-03 | OTP rate limit: 3 per 15 min per user |
| BR-NOT-04 | Marketing SMS: max 1 per day per user |
| BR-NOT-05 | Notification failures retried 3x, then logged for manual review |
| BR-NOT-06 | Notification preferences sync across all user devices |

## Key Workflows

### Transactional Notification Flow
1. Event triggers notification (e.g., policy purchased)
2. System selects template based on event type
3. Populates template with dynamic data (customer name, policy number, etc.)
4. Sends via SMS + email (parallel)
5. Tracks delivery status (sent, delivered, failed)
6. Logs notification for audit

### Marketing Campaign Flow
1. Marketing team creates campaign in Admin Portal
2. Selects target audience (filters by product, region, etc.)
3. System applies consent filter (only opted-in users)
4. System applies rate limits (skip users who received SMS today)
5. Campaign queued and sent in batches
6. Delivery report generated (sent, delivered, clicked, unsubscribed)

### Notification Preference Update
1. Customer opens app → Settings → Notifications
2. Toggles preferences: marketing (on/off), channels (SMS/email/push), frequency
3. Preference saved → synced to backend
4. Future notifications respect updated preferences

## Data Model Notes

**Notification Entity**
- notification_id
- user_id
- notification_type (OTP, PURCHASE, CLAIM_UPDATE, RENEWAL_REMINDER, MARKETING)
- channel (SMS, EMAIL, PUSH)
- template_id
- content (rendered message)
- status (PENDING, SENT, DELIVERED, FAILED, BOUNCED)
- sent_at, delivered_at

**Notification Preferences**
- user_id
- marketing_consent (boolean)
- channels_enabled (SMS, EMAIL, PUSH)
- consent_timestamp
- last_updated

## Integration Touchpoints

| External System | Purpose | Failure Handling |
|----------------|---------|------------------|
| SMS Gateway (Twilio, local provider) | SMS delivery | Retry 3x, log failure, fallback to email |
| Email Service (SendGrid, SES) | Email delivery | Retry 3x, log failure |
| Push Notification Service (FCM, APNS) | Mobile app push | Best effort (no retry) |

## Security & Privacy

- Marketing consent recorded with timestamp (regulatory requirement)
- PII in notifications minimized (use first name only, masked phone/account numbers)
- Notification logs retained per compliance period
- Unsubscribe links include secure tokens (prevent abuse)

## NFR Constraints

| NFR | Target |
|-----|--------|
| OTP Delivery Time | <60s, 95% success rate |
| Transactional SMS Delivery | <2 min, 95% success rate |
| Email Delivery | <5 min, 90% success rate |
| Push Notification Delivery | <10s (best effort) |
| Notification Throughput | 10,000 SMS/min (campaign bursts) |

## Acceptance Criteria

- [ ] Customers receive transactional notifications (OTP, purchase, claims) reliably
- [ ] Renewal reminders sent at defined intervals
- [ ] Marketing notifications respect customer consent
- [ ] Customers can manage notification preferences
- [ ] Rate limits prevent spam and abuse
- [ ] Notification delivery tracked and logged

## Traceability

**SRS Reference:** FG-012 — 4.12 Notifications & Communication  
**Functional Requirements:** FR-114, FR-115, FR-116, FR-117, FR-118, FR-119, FR-120, FR-121, FR-122, FR-123

[[[PAGEBREAK]]]


---
[[[PAGEBREAK]]]


# Feature Group: Voice-Assisted Flows (FG-013)

## Business Objective

Enable voice-based interactions for customers with low digital literacy or accessibility needs, expanding market reach to rural, elderly, and visually-impaired segments. Support Bengali voice commands and voice-driven policy purchase, claims, and support flows.

**Business Value:**
- Expand addressable market by 20-30% (rural, low-literacy, elderly)
- Reduce agent dependency for simple transactions
- Accessibility compliance (inclusive insurance)
- Differentiation in Bangladesh market

## Actors & Portals

| Actor | Portal(s) | Responsibilities |
|-------|-----------|------------------|
| Customer (voice user) | Mobile App (voice interface), IVR | Use voice commands for policy queries, purchases, claims |
| System (Voice AI) | Backend Voice Service | Speech-to-text, natural language understanding, text-to-speech |
| Support Agent | Admin Portal | Handle voice flow escalations |
| Business Admin | Admin Portal | Configure voice scripts, train voice models |

## User Stories

### US-FG-013-01: Voice-Driven Policy Search

**As a** customer with low digital literacy  
**I want** to search for insurance products using voice in Bengali  
**So that** I can find products without reading text

**Acceptance Criteria:**
- Customer taps microphone icon → speaks query in Bengali
- System converts speech to text → matches products
- System reads product names and summaries aloud
- Customer can say "More details" or "Buy this" for next steps
- Fallback: if speech unclear, system asks for clarification

![Voice Search Flow](images/flow_voice_search.png)

**Related FRs:** FR-124, FR-125, FR-126

### US-FG-013-02: Voice-Assisted Policy Purchase

**As a** elderly customer  
**I want** to complete policy purchase using voice commands  
**So that** I don't need to type complex forms

**Acceptance Criteria:**
- System guides customer step-by-step via voice prompts
- Customer provides details verbally (name, age, coverage amount)
- System confirms each input ("You said [value], is that correct?")
- Payment initiated after voice confirmation
- Fallback to agent if customer gets stuck

**Related FRs:** FR-178, FR-179

### US-FG-013-03: Voice-Based Claims Status Query

**As a** customer  
**I want** to check my claim status by asking "What's my claim status?"  
**So that** I get instant updates without navigating menus

**Acceptance Criteria:**
- Customer says "Check my claim" or "Claim status"
- System identifies customer (voice biometric or fallback to OTP)
- Reads claim status aloud with next steps
- Customer can ask follow-up questions ("When will I get paid?")

![Voice Claim Query](images/flow_voice_claim_query.png)

**Related FRs:** FR-180, FR-181

### US-FG-013-04: Voice Authentication (Future)

**As a** customer  
**I want** to login using my voice  
**So that** access is fast and secure

**Acceptance Criteria:**
- Customer says passphrase or customer ID
- System matches voice biometric (if enrolled)
- Fallback to OTP if voice match fails or not enrolled
- Voice print stored securely (encrypted)

**Related FRs:** FR-182

## Business Rules

| Rule ID | Description |
|---------|-------------|
| BR-VOI-01 | Voice flows available in Bengali and English |
| BR-VOI-02 | Voice unclear → system asks for clarification (max 3 attempts) |
| BR-VOI-03 | Complex transactions → offer agent escalation |
| BR-VOI-04 | Voice authentication requires explicit enrollment consent |
| BR-VOI-05 | All voice interactions logged for quality and compliance |

## Key Workflows

### Voice-Driven Purchase
1. Customer opens app → taps "Voice Assistant"
2. System: "How can I help you today?"
3. Customer: "I want health insurance"
4. System: matches products → reads options
5. Customer: "Tell me more about [product name]"
6. System: reads coverage details
7. Customer: "I want to buy"
8. System: guides through data collection (name, age, etc.)
9. Customer confirms each step verbally
10. System initiates payment → policy issued

### Voice Escalation to Agent
1. Customer stuck or confused
2. System detects frustration or repeated clarification requests
3. System: "Would you like me to connect you to an agent?"
4. Customer: "Yes"
5. System transfers to live agent with context (transcript, intent)

## Data Model Notes

**Voice Interaction Log**
- interaction_id
- customer_id
- transcript (speech-to-text output)
- intent (NLU classification)
- language (bn, en)
- timestamp
- outcome (SUCCESS, ESCALATED, ABANDONED)

## Integration Touchpoints

| External System | Purpose | Failure Handling |
|----------------|---------|------------------|
| Speech-to-Text Service (Google/AWS) | Convert voice to text | Queue for manual processing, fallback to text input |
| NLU Engine | Understand intent | Rule-based fallback |
| Text-to-Speech Service | Read responses | Display text on screen as fallback |

## Security & Privacy

- Voice recordings encrypted at rest
- Customer consent required for voice data storage
- Voice biometric data stored separately, high encryption
- Transcripts retained per data retention policy

## NFR Constraints

| NFR | Target |
|-----|--------|
| Speech Recognition Accuracy | >90% for Bengali (trained model) |
| Response Latency | <3s from voice input to system response |
| Voice Session Timeout | 2 min inactivity → session ends |

## Acceptance Criteria

- [ ] Customer can search products using Bengali voice
- [ ] Voice-driven purchase completes end-to-end for simple products
- [ ] Claims status queries answered via voice
- [ ] Voice flows escalate to agent when customer stuck
- [ ] All voice interactions logged and auditable

## Traceability

**SRS Reference:** FG-013 — 4.13 IoT Integration & Usage-Based Insurance  
**Functional Requirements:** FR-124, FR-125, FR-126, FR-127, FR-128, FR-129, FR-130, FR-178, FR-179, FR-180, FR-181, FR-182, FR-183, FR-184, FR-185

[[[PAGEBREAK]]]


---
[[[PAGEBREAK]]]


# Feature Group: AI/ML Integration (FG-014)

## Business Objective

Leverage AI/ML for fraud detection pattern recognition, underwriting assistance, customer risk scoring, personalized product recommendations, and operational optimization while maintaining transparency and regulatory compliance.

**Business Value:**
- Reduce fraud losses through ML pattern detection (10-20% improvement over rules)
- Accelerate underwriting decisions (future: instant approval for low-risk)
- Increase conversion through personalized recommendations
- Optimize operations (predict claim volumes, staff allocation)

## Actors & Portals

| Actor | Portal(s) | Responsibilities |
|-------|-----------|------------------|
| System (AI/ML Models) | Background ML Services | Run predictions, scoring, recommendations |
| Business Admin | Admin Portal | Monitor model performance, approve model deployments |
| Data Scientist (future) | ML Ops Portal | Train models, tune hyperparameters |
| Customer (indirect) | Mobile App | Receives personalized recommendations |

## User Stories

### US-FG-014-01: ML-Based Fraud Scoring

**As a** fraud detection system  
**I want** ML models to score transactions for fraud risk  
**So that** subtle patterns are caught beyond rule-based detection

**Acceptance Criteria:**
- ML model runs in parallel with rule-based fraud detection
- Outputs fraud risk score (0-100)
- High-risk scores (>80) trigger manual review
- Model retrained monthly on new fraud patterns
- Model performance tracked (precision, recall, false positives)

![ML Fraud Scoring](images/flow_ml_fraud_scoring.png)

**Related FRs:** FR-166, FR-167

### US-FG-014-02: Personalized Product Recommendations

**As a** customer  
**I want** to see insurance products recommended for me  
**So that** I find relevant coverage quickly

**Acceptance Criteria:**
- ML model considers: customer age, occupation, past purchases, browsing history
- Recommendations shown on homepage and product pages
- Explanations provided ("Recommended because...")
- Customer can hide recommendations or provide feedback
- Recommendations comply with fairness constraints (no discriminatory patterns)

![Product Recommendations](images/ui_ml_recommendations.png)

**Related FRs:** FR-164, FR-165

### US-FG-014-03: Underwriting Assistance (Future)

**As an** underwriter  
**I want** AI to suggest approval/rejection with confidence score  
**So that** low-risk applications are fast-tracked

**Acceptance Criteria:**
- ML model trained on historical underwriting decisions
- Outputs: APPROVE/REJECT/MANUAL_REVIEW + confidence (0-100%)
- High-confidence approvals (>95%) can be auto-approved (subject to limits)
- Low-confidence cases routed to human underwriter
- Model explainability: shows key factors in decision

**Related FRs:** FR-166, FR-167

### US-FG-014-04: Claims Amount Prediction

**As a** Business Admin  
**I want** to predict total claims volume for next month  
**So that** I can allocate reserves and staff appropriately

**Acceptance Criteria:**
- ML model trained on historical claims data
- Outputs: predicted claim count, predicted total payout
- Confidence intervals provided
- Actual vs predicted tracked monthly for model accuracy

**Related FRs:** FR-168

## Business Rules

| Rule ID | Description |
|---------|-------------|
| BR-AI-01 | ML models are supplementary; critical decisions require human-in-the-loop |
| BR-AI-02 | Model predictions must be explainable (no pure black-box for regulated decisions) |
| BR-AI-03 | Models retrained periodically (monthly/quarterly) to adapt to new patterns |
| BR-AI-04 | Fairness constraints enforced (no discrimination by protected attributes) |
| BR-AI-05 | Model deployment requires Business Admin approval and A/B testing |

## Key Workflows

### ML Model Training & Deployment
1. Data Scientist trains model on historical data
2. Model validated on holdout test set (accuracy, fairness metrics)
3. Model deployed to staging for A/B test
4. Business Admin reviews A/B results
5. If successful → promote to production
6. Model performance monitored continuously

### ML-Augmented Fraud Detection
1. Transaction submitted (purchase, claim, payment)
2. Rule-based fraud detection runs first
3. ML model runs in parallel → outputs risk score
4. **If rules flag OR ML score >threshold:** Route to manual review
5. **If both clean:** Proceed
6. Fraud Analyst reviews flagged cases, provides feedback (true positive/false positive)
7. Feedback used to retrain model

## Data Model Notes

**ML Model Metadata**
- model_id
- model_name (fraud_scorer_v3, product_recommender_v2)
- model_type (CLASSIFICATION, REGRESSION, RANKING)
- version
- training_date
- performance_metrics (accuracy, precision, recall, AUC)
- status (TRAINING, STAGING, PRODUCTION, RETIRED)

**ML Prediction Log**
- prediction_id
- model_id
- input_features (JSON)
- prediction_output
- confidence_score
- timestamp

## Integration Touchpoints

| External System | Purpose | Failure Handling |
|----------------|---------|------------------|
| ML Training Infrastructure (Cloud) | Train models | Fallback to previous model version |
| Feature Store | Real-time feature retrieval | Cache frequently-used features |

## Security & Privacy

- ML models trained on anonymized/pseudonymized data where possible
- Customer consent for using data in ML training (aggregated, non-identifiable)
- Model predictions logged for audit and bias detection
- Access to ML training data restricted

## NFR Constraints

| NFR | Target |
|-----|--------|
| ML Prediction Latency | <500ms for real-time scoring |
| Model Retraining Frequency | Monthly (or triggered by performance degradation) |
| Model Explainability | SHAP/LIME values available for regulated decisions |

## Acceptance Criteria

- [ ] ML fraud scoring augments rule-based detection
- [ ] Product recommendations increase conversion (A/B tested)
- [ ] Underwriting assistance reduces TAT for low-risk cases
- [ ] Model performance tracked and retrained regularly
- [ ] Fairness and explainability requirements met

## Traceability

**SRS Reference:** FG-014 — 4.14 AI & Automation Features  
**Functional Requirements:** FR-164, FR-165, FR-166, FR-167, FR-168, FR-169

[[[PAGEBREAK]]]


---
[[[PAGEBREAK]]]


# Feature Group: IoT & Telematics Integration (FG-015)

## Business Objective

Enable usage-based insurance products through IoT device integration (telematics for vehicles, wearables for health, smart home sensors) to enable dynamic pricing, proactive risk management, and differentiated product offerings.

**Business Value:**
- Enable pay-as-you-drive motor insurance (competitive differentiation)
- Activity-based health insurance premiums (reward healthy behavior)
- Reduce claims through proactive risk alerts (driver coaching, health alerts)
- Data-driven pricing (fair premiums based on actual behavior)

## Actors & Portals

| Actor | Portal(s) | Responsibilities |
|-------|-----------|------------------|
| Customer | Mobile App | Pair IoT devices, view usage data, receive risk alerts |
| Partner (IoT Provider) | IoT Management Portal | Onboard devices, monitor health, share telemetry |
| Business Admin | Admin Portal | Configure IoT product rules, monitor device adoption |
| System (IoT Platform) | Background IoT Service | Collect telemetry, calculate usage scores, trigger alerts |

## User Stories

### US-FG-015-01: Telematics Device Pairing (Motor Insurance)

**As a** customer with motor insurance  
**I want** to pair a telematics device to my policy  
**So that** my premium reflects my safe driving

**Acceptance Criteria:**
- Customer purchases motor policy with telematics discount option
- Receives OBD-II device or app-based telematics
- Customer activates device → pairs with policy via QR/code
- System starts collecting driving data (speed, braking, acceleration, mileage)
- Customer can view driving score in app

![Telematics Pairing](images/flow_telematics_pairing.png)

**Related FRs:** FR-174

### US-FG-015-02: Dynamic Premium Adjustment

**As a** customer  
**I want** my premium adjusted based on my safe driving  
**So that** I pay less if I drive safely

**Acceptance Criteria:**
- System calculates driving score monthly (0-100)
- Score based on: mileage, speed adherence, harsh braking events, time of day
- Premium discount applied at renewal: 0-30% based on score
- Customer receives monthly score report with tips for improvement
- Discount calculation transparent and auditable

**Related FRs:** FR-171, FR-172

### US-FG-015-03: Proactive Risk Alerts

**As a** customer  
**I want** alerts if my driving pattern is risky  
**So that** I can improve and avoid accidents

**Acceptance Criteria:**
- System detects risk patterns (repeated harsh braking, frequent speeding)
- Sends push notification: "You had 5 harsh braking events this week. Drive carefully!"
- Customer can view driving history and patterns in app
- Alerts optional (customer can opt out)

![Risk Alerts](images/notification_iot_risk_alert.png)

**Related FRs:** FR-173, FR-174

### US-FG-015-04: Wearable Integration (Health Insurance)

**As a** customer with health insurance  
**I want** to connect my fitness tracker  
**So that** my active lifestyle is rewarded

**Acceptance Criteria:**
- Customer pairs wearable (Fitbit, Apple Watch, Mi Band)
- System collects: steps, heart rate, active minutes (with consent)
- Activity score calculated monthly
- Discount applied at renewal or wellness rewards earned
- Customer privacy: raw data not shared with insurer, only aggregated score

**Related FRs:** FR-175, FR-176

## Business Rules

| Rule ID | Description |
|---------|-------------|
| BR-IOT-01 | IoT device pairing requires explicit customer consent |
| BR-IOT-02 | Customer can unpair device anytime (but loses usage-based discount) |
| BR-IOT-03 | Raw telemetry data encrypted in transit and at rest |
| BR-IOT-04 | Premium adjustments capped (max 30% discount, no penalties for unpair) |
| BR-IOT-05 | Device health monitored; inactive devices flagged for replacement |

## Key Workflows

### Telematics Onboarding & Scoring
1. Customer purchases motor policy with telematics option
2. Receives device (mailed or partner pickup)
3. Installs device → activates via app (QR scan or code entry)
4. System starts collecting trips and driving events
5. Monthly score calculated → customer notified
6. At renewal: premium adjusted based on score

### Proactive Risk Alert
1. IoT system detects risk pattern (e.g., 3+ harsh braking events in 1 day)
2. System triggers risk alert notification
3. Customer receives push + in-app message with driving tips
4. If pattern persists → escalate to support outreach (optional coaching)

### Device Health Monitoring
1. System monitors device connectivity (last data received)
2. If device offline >7 days → notify customer "Device may be disconnected"
3. Customer troubleshoots or requests replacement
4. If device fails → customer loses discount until replacement active

## Data Model Notes

**IoT Device Entity (SRS Proto: insuretech.iot.entity.v1.Device)**
- device_id
- device_type (OBD_TELEMATICS, WEARABLE, SMART_HOME)
- policy_id (linked policy)
- customer_id
- activation_date
- last_data_received
- status (ACTIVE, INACTIVE, FAILED)

**Telemetry Data**
- telemetry_id
- device_id
- timestamp
- data_payload (JSON: speed, location, heart_rate, etc.)
- privacy_level (AGGREGATED, DETAILED)

**Usage Score**
- score_id
- device_id
- policy_id
- period (month)
- score (0-100)
- discount_applied (%)

## Integration Touchpoints

| External System | Purpose | Failure Handling |
|----------------|---------|------------------|
| IoT Device Providers | Device onboarding, telemetry ingestion | Queue data, batch upload when connection restored |
| Wearable APIs (Fitbit, Apple Health) | Sync activity data | Best effort, customer can manually trigger sync |
| Mapping/Geo Services | Validate trip routes | Fallback to device GPS |

## Security & Privacy

- Telemetry data encrypted end-to-end
- Customer location data anonymized/aggregated (no real-time tracking exposed)
- Customer can delete IoT data history (subject to policy period)
- Consent management: customer controls what data is shared

## NFR Constraints

| NFR | Target |
|-----|--------|
| Telemetry Ingestion Throughput | 10,000 events/sec |
| Scoring Calculation Latency | <1 hour after data batch received |
| Device Pairing Time | <2 min from activation to first data |

## Acceptance Criteria

- [ ] Customer can pair telematics device to motor policy
- [ ] Driving score calculated and displayed monthly
- [ ] Premium discount applied at renewal based on score
- [ ] Proactive risk alerts sent for unsafe patterns
- [ ] Wearable integration works for health insurance (future)
- [ ] Customer privacy and consent managed properly

## Traceability

**SRS Reference:** FG-015 — 4.15 Voice-Assisted Features  
**Functional Requirements:** FR-170, FR-171, FR-172, FR-173, FR-174, FR-175, FR-176, FR-177

[[[PAGEBREAK]]]


---
[[[PAGEBREAK]]]


# Feature Group: Fraud Detection & Risk Controls (FG-016)

## Business Objective

Detect and prevent fraudulent activities across policy issuance, claims, and payments through configurable rule engines, anomaly detection, and manual review workflows. Protect revenue, maintain customer trust, and support regulatory reporting obligations.

**Business Value:**
- Reduce fraud losses by >90% (target detection rate)
- Protect honest customers from premium increases caused by fraud
- Enable rapid response to emerging fraud patterns
- Support AML/CFT compliance and suspicious activity reporting

## Actors & Portals

| Actor | Portal(s) | Responsibilities |
|-------|-----------|------------------|
| System (Automated) | Background Jobs | Run fraud detection rules on transactions |
| Fraud Analyst | Admin Portal | Review flagged cases, investigate patterns, approve/reject |
| Business Admin | Admin Portal | Configure fraud rules, set thresholds, approve actions |
| Compliance Officer | Admin Portal | Review suspicious activity reports, file STR/SAR |
| Customer (indirect) | Mobile App | May be asked for additional verification if flagged |

## User Stories

### US-FG-016-01: Real-Time Fraud Rule Execution

**As a** fraud detection system  
**I want** to run fraud rules on every transaction in real-time  
**So that** fraudulent activity is caught before payout

**Acceptance Criteria:**
- Fraud rules execute on: policy purchase, claims submission, payments, endorsements
- Rules include: duplicate detection, velocity checks, anomaly patterns, blacklist matching
- Flagged transactions routed to review queue (do not auto-approve)
- Non-flagged transactions proceed normally
- Rule execution latency: <2s

![Fraud Detection Flow](images/flow_fraud_detection.png)

**Related FRs:** FR-187, FR-190, FR-191

### US-FG-016-02: Configurable Fraud Rules

**As a** Business Admin  
**I want** to configure fraud detection rules without developer involvement  
**So that** we can adapt to new fraud patterns quickly

**Acceptance Criteria:**
- Admin can create/edit/disable fraud rules via UI
- Rule types: threshold-based (e.g., claim amount >X), pattern-based (e.g., multiple policies same NID), velocity (e.g., 3+ purchases in 1 hour)
- Rules versioned with audit trail
- Test mode available (flag but don't block)
- Rules can be product-specific or global

![Fraud Rule Configuration](images/admin_fraud_rules.png)

**Related FRs:** FR-187, FR-188

### US-FG-016-03: Fraud Review Queue & Investigation

**As a** Fraud Analyst  
**I want** a prioritized queue of flagged transactions  
**So that** I can investigate and decide quickly

**Acceptance Criteria:**
- Queue shows: transaction type, customer, amount, rule(s) triggered, risk score
- Analyst can view full transaction history, customer profile, related policies/claims
- Analyst actions: Approve (false positive), Block (fraud confirmed), Escalate (needs compliance review)
- All actions logged with reason
- SLA: high-risk cases reviewed within 4 hours

![Fraud Review Queue](images/admin_fraud_queue.png)

**Related FRs:** FR-189, FR-190

### US-FG-016-04: Blacklist Management

**As a** Fraud Analyst  
**I want** to maintain blacklists of fraudulent customers/NIDs/devices  
**So that** repeat offenders are auto-blocked

**Acceptance Criteria:**
- Blacklists: NID, mobile number, email, device ID, payment card
- Blacklist entries include: reason, date added, added by
- System auto-rejects transactions matching blacklist
- Blacklist entries can be removed (e.g., false positive resolved)
- Audit trail for all blacklist changes

**Related FRs:** FR-054, FR-191

## Business Rules

| Rule ID | Description |
|---------|-------------|
| BR-FRD-01 | Flagged transactions require manual review before approval |
| BR-FRD-02 | Fraud rules execute before payment confirmation or claim settlement |
| BR-FRD-03 | Blacklisted customers/NIDs auto-rejected with clear error message |
| BR-FRD-04 | Fraud rule changes require Business Admin approval |
| BR-FRD-05 | High-risk fraud cases (>100k BDT) escalated to Compliance Officer |
| BR-FRD-06 | Fraud analyst actions logged with timestamp and reason |

## Key Workflows

### Real-Time Fraud Detection (Purchase/Claim)
1. Customer submits transaction (purchase, claim, payment)
2. System runs fraud detection rules (parallel execution)
3. **If no flags:** Transaction proceeds normally
4. **If flagged:** Transaction held in review queue, customer notified "under review"
5. Fraud Analyst reviews case within SLA
6. **If approved:** Transaction released, customer notified
7. **If blocked:** Transaction rejected, customer notified, entry logged for investigation

### Fraud Rule Configuration
1. Business Admin navigates to "Fraud Rules" in Admin Portal
2. Creates new rule: name, type (threshold/pattern/velocity), parameters, severity (LOW/MEDIUM/HIGH)
3. Tests rule in sandbox mode (flags but doesn't block)
4. Activates rule → applies to all new transactions
5. Rule version history maintained

### Fraud Investigation & Escalation
1. Fraud Analyst reviews flagged case
2. Views customer history: past policies, claims, payment patterns
3. If pattern suspicious → marks as fraud, adds to blacklist
4. If high-value case → escalates to Compliance Officer
5. Compliance Officer may file STR/SAR with BFIU
6. Customer account suspended if fraud confirmed

## Data Model Notes

**Fraud Rule Entity**
- rule_id
- rule_name
- rule_type (THRESHOLD, PATTERN, VELOCITY, BLACKLIST)
- rule_logic (JSON/DSL)
- severity (LOW, MEDIUM, HIGH, CRITICAL)
- status (ACTIVE, DISABLED, TEST_MODE)
- version, created_by, updated_at

**Fraud Alert Entity**
- alert_id
- transaction_id (policy/claim/payment)
- customer_id
- rules_triggered (list)
- risk_score (1-100)
- status (PENDING_REVIEW, APPROVED, BLOCKED, ESCALATED)
- reviewed_by, reviewed_at
- resolution_notes

**Blacklist Entity**
- blacklist_id
- blacklist_type (NID, MOBILE, EMAIL, DEVICE_ID, CARD)
- value
- reason
- added_by, added_at

## Integration Touchpoints

| External System | Purpose | Failure Handling |
|----------------|---------|------------------|
| AML/CFT Monitoring | Share fraud patterns for AML compliance | Queue for manual filing if integration down |
| Customer Risk Scoring (future) | ML-based risk scores | Fallback to rule-based detection |

## Security & Privacy

- Fraud detection runs in isolated environment (prevent gaming)
- Fraud rules not exposed to customers (security through obscurity)
- Blacklist data encrypted and access-controlled
- All fraud decisions auditable (who, when, why)

## NFR Constraints

| NFR | Target |
|-----|--------|
| Rule Execution Latency | <2s per transaction |
| False Positive Rate | <10% (balance sensitivity vs customer friction) |
| Fraud Review SLA | <4 hours for high-risk, <24 hours for medium |
| Rule Configuration Downtime | <1 min (hot reload) |

## Acceptance Criteria

- [ ] Fraud rules execute in real-time on all critical transactions
- [ ] Business Admin can configure rules via UI without developer
- [ ] Fraud Analysts have prioritized review queue
- [ ] Blacklist auto-blocks repeat offenders
- [ ] All fraud actions auditable with full history
- [ ] False positive rate tracked and optimized

## Traceability

**SRS Reference:** FG-016 — 4.16 Fraud Detection & Risk Controls  
**Functional Requirements:** FR-186, FR-187, FR-188, FR-189, FR-190, FR-191, FR-192

[[[PAGEBREAK]]]


---
[[[PAGEBREAK]]]


# Feature Group: Admin Operations & Configuration (FG-017)

## Business Objective

Provide business and system administrators with tools to manage platform configuration, workflows, user roles, product rules, and operational overrides in a controlled and auditable manner. Enable business agility without requiring developer involvement.

**Business Value:**
- Accelerate time-to-market for product changes (config vs code)
- Enable business owners to manage rules and workflows directly
- Maintain operational safety through approval workflows
- Auditability of all administrative actions

## Actors & Portals

| Actor | Portal(s) | Responsibilities |
|-------|-----------|------------------|
| Business Admin | Admin Portal | Configure products, pricing, workflows, approval rules, manage users |
| System Admin | Admin Portal | System configuration, role management, security settings, incident tooling |
| Focal Person | Admin Portal | Partner approvals, dispute resolution, cross-partner oversight |
| Operations Team | Admin Portal | Monitor queues, resolve stuck workflows, operational overrides |

## User Stories

### US-FG-017-01: Product Configuration Management

**As a** Business Admin  
**I want** to configure product rules without developer help  
**So that** new products or changes go live quickly

**Acceptance Criteria:**
- Admin can create/edit/deactivate products
- Configure coverage rules, exclusions, deductibles, co-pay
- Set pricing rules (flat, tiered, age-based)
- Version history maintained
- Changes require approval before going live

![Product Configuration](images/admin_product_config.png)

**Related FRs:** FR-134

### US-FG-017-02: Workflow Approval Matrix Configuration

**As a** Business Admin  
**I want** to configure approval workflows (who approves what)  
**So that** governance is enforced automatically

**Acceptance Criteria:**
- Admin defines approval rules: endorsements, cancellations, claims (by amount)
- Tiered approvals: L1, L2, joint approvals for high-value
- Rules stored with version history
- Test mode available (simulate approval without executing)

**Related FRs:** FR-134, FR-135

### US-FG-017-03: User Role Management

**As a** System Admin  
**I want** to manage user roles and permissions  
**So that** access is controlled and auditable

**Acceptance Criteria:**
- Admin can create/edit/disable user accounts
- Assign roles: Customer, Agent, Partner Admin, Focal Person, Business Admin, System Admin
- Role-permission mappings configurable
- MFA enforcement for admin roles
- All role changes logged

![User Role Management](images/admin_user_roles.png)

**Related FRs:** FR-014, FR-015, FR-019

### US-FG-017-04: Operational Override (Emergency)

**As a** Business Admin  
**I want** to override stuck workflows or approve exceptions  
**So that** urgent cases don't block operations

**Acceptance Criteria:**
- Admin can manually approve/reject stuck items (claims, endorsements)
- Override requires: reason, evidence (attachments), approval from higher authority
- All overrides logged and flagged for audit review
- Override capability time-limited (expires after X hours)

**Related FRs:** FR-136, FR-137

## Business Rules

| Rule ID | Description |
|---------|-------------|
| BR-ADM-01 | Admin configuration changes require approval before going live |
| BR-ADM-02 | System Admin actions logged with actor, timestamp, reason |
| BR-ADM-03 | Operational overrides require Business Admin or higher authority |
| BR-ADM-04 | Role changes trigger notification to affected users |
| BR-ADM-05 | Admin MFA mandatory for production environment |

## Key Workflows

### Product Configuration Change
1. Business Admin navigates to "Products" → selects product
2. Edits rules (coverage, pricing, exclusions)
3. Saves as draft → submits for approval
4. Focal Person or senior Business Admin reviews
5. If approved → goes live (version incremented)
6. If rejected → Admin notified with reason

### Operational Override (Stuck Claim)
1. Operations team identifies stuck claim (e.g., approval timed out)
2. Business Admin opens override panel
3. Reviews case details, reason for override
4. Enters override justification + attaches evidence
5. Submits for higher approval (if required by policy)
6. Override logged and audit-flagged
7. Claim released/approved

### User Role Assignment
1. System Admin creates new user or edits existing
2. Assigns role(s) + tenant (if partner/agent)
3. Sets MFA requirement (mandatory for admin roles)
4. User notified of account creation/change
5. Role assignment logged

## Data Model Notes

**Admin Configuration Entity**
- config_id
- config_type (PRODUCT, WORKFLOW, APPROVAL_MATRIX, SYSTEM)
- config_value (JSON)
- version
- status (DRAFT, PENDING_APPROVAL, ACTIVE, RETIRED)
- created_by, approved_by
- timestamps

**Admin Action Log**
- action_id
- admin_user_id
- action_type (CONFIG_CHANGE, OVERRIDE, ROLE_ASSIGNMENT)
- entity_id (affected entity)
- reason
- approval_required (boolean)
- approved_by
- timestamp

## Integration Touchpoints

| External System | Purpose | Failure Handling |
|----------------|---------|------------------|
| Audit Log Service | Log all admin actions | Queue locally if service down, batch upload |
| Notification Service | Notify users of role changes | Retry, fallback to email |

## Security & Privacy

- Admin actions require authentication + MFA
- Admin access logs monitored for suspicious patterns
- Configuration changes versioned and rollback-capable
- Override actions flagged for compliance review

## NFR Constraints

| NFR | Target |
|-----|--------|
| Admin Portal Availability | 99.5% (less critical than customer-facing) |
| Configuration Change Propagation | <5 min from approval to live |
| Admin Action Audit Log Retention | 7 years (compliance requirement) |

## Acceptance Criteria

- [ ] Business Admin can configure products without developer
- [ ] Approval workflows configurable and enforced
- [ ] System Admin can manage user roles and permissions
- [ ] Operational overrides logged and auditable
- [ ] All admin actions traceable with actor and reason

## Traceability

**SRS Reference:** FG-017 — 4.17 Admin & Reporting  
**Functional Requirements:** FR-131, FR-132, FR-133, FR-134, FR-135, FR-136, FR-137, FR-138, FR-139, FR-140

[[[PAGEBREAK]]]


---
[[[PAGEBREAK]]]


# Feature Group: Analytics & Reporting (FG-018)

## Business Objective

Deliver real-time dashboards, operational reports, and regulatory extracts to enable data-driven decision-making, performance monitoring, and compliance reporting. Support business users, partners, and regulatory stakeholders.

**Business Value:**
- Visibility into business performance (policies, claims, revenue)
- Operational insights (conversion funnel, bottlenecks, TAT tracking)
- Partner/agent performance tracking
- Regulatory reporting readiness (IDRA, BFIU)

## Actors & Portals

| Actor | Portal(s) | Responsibilities |
|-------|-----------|------------------|
| Business Executive | Admin Portal | View executive dashboards (KPIs, trends) |
| Business Admin | Admin Portal | Operational reports, product performance, claim analytics |
| Partner Admin | Partner Portal | Partner-specific dashboards (sales, commissions) |
| Compliance Officer | Admin Portal | Regulatory reports, suspicious activity extracts |
| Agent | Agent Mobile App | Agent performance, commission statements |

## User Stories

### US-FG-018-01: Executive KPI Dashboard

**As a** Business Executive  
**I want** a real-time dashboard showing key metrics  
**So that** I can monitor business health at a glance

**Acceptance Criteria:**
- Dashboard shows: policies issued (today/week/month), premium collected, claims paid, conversion rate, partner performance
- Filterable by date range, product, region
- Trend charts (line graphs showing growth)
- Drill-down capability (tap metric → detailed view)

![Executive Dashboard](images/dashboard_executive_kpi.png)

**Related FRs:** FR-203, FR-204

### US-FG-018-02: Claims Analytics Report

**As a** Business Admin  
**I want** detailed claims analytics  
**So that** I can identify fraud patterns and optimize TAT

**Acceptance Criteria:**
- Report shows: total claims, approval rate, rejection reasons, average settlement time, top claim types
- Filterable by product, date, status
- Export as CSV/Excel
- Visualizations: bar charts (claim types), pie charts (approval/rejection), trend lines (TAT over time)

![Claims Analytics](images/report_claims_analytics.png)

**Related FRs:** FR-159, FR-160

### US-FG-018-03: Partner Performance Report

**As a** Partner Admin  
**I want** to see my organization's performance  
**So that** I can optimize my distribution strategy

**Acceptance Criteria:**
- Report shows: leads, quotes generated, policies issued, conversion rate, commission earned (pending/paid)
- Agent-level breakdown
- Filterable by date range, product
- Exportable for accounting

**Related FRs:** FR-145, FR-146, FR-147

### US-FG-018-04: Regulatory Reporting (IDRA/BFIU)

**As a** Compliance Officer  
**I want** to generate regulatory-compliant reports  
**So that** audits and submissions are streamlined

**Acceptance Criteria:**
- Pre-configured report templates (IDRA quarterly, BFIU suspicious activity)
- Reports include: policy counts, premium volumes, claim payouts, cancellations, customer demographics
- Reports exportable in required formats (CSV, Excel, PDF)
- Audit trail of report generation (who, when, what data)

![Regulatory Report](images/report_regulatory_extract.png)

**Related FRs:** FR-161, FR-162

## Business Rules

| Rule ID | Description |
|---------|-------------|
| BR-RPT-01 | Dashboards refresh every 15 min (real-time for critical metrics) |
| BR-RPT-02 | Reports accessible only to authorized roles (RBAC enforced) |
| BR-RPT-03 | Sensitive data (PII) redacted in analytics unless explicitly authorized |
| BR-RPT-04 | Report generation logged for audit |
| BR-RPT-05 | Large reports (>10k rows) exported asynchronously with email notification |

## Key Workflows

### Dashboard Access
1. User logs in to Admin/Partner Portal
2. Navigates to "Dashboards" → selects dashboard type
3. Dashboard loads with latest data (cached for performance)
4. User applies filters (date, product, region)
5. Dashboard updates dynamically

### Report Generation
1. User navigates to "Reports" → selects report template
2. Applies filters (date range, product, etc.)
3. Clicks "Generate Report"
4. **If small:** Report generated immediately, displayed in browser
5. **If large:** Report queued, user notified via email when ready
6. User downloads report (CSV/Excel/PDF)

### Scheduled Reports (Future)
1. User sets up scheduled report (weekly/monthly)
2. System generates report automatically on schedule
3. Report emailed to user or saved to shared folder

## Data Model Notes

**Dashboard Configuration**
- dashboard_id
- dashboard_name (Executive KPI, Claims Analytics, Partner Performance)
- widgets (list of chart/metric widgets)
- default_filters
- authorized_roles

**Report Template**
- template_id
- template_name
- report_type (OPERATIONAL, REGULATORY, PARTNER)
- query_logic (SQL/aggregation rules)
- export_formats (CSV, EXCEL, PDF)

**Report Generation Log**
- log_id
- user_id
- template_id
- filters_applied (JSON)
- generated_at
- download_count

## Integration Touchpoints

| External System | Purpose | Failure Handling |
|----------------|---------|------------------|
| Data Warehouse / OLAP | Aggregated data for analytics | Cache last successful query results |
| Export Service | Generate Excel/PDF | Queue and retry |

## Security & Privacy

- Dashboards/reports enforce role-based access (Partner Admins see only their data)
- PII redacted in aggregate reports
- Report downloads logged for audit
- Sensitive reports (regulatory) require additional approval

## NFR Constraints

| NFR | Target |
|-----|--------|
| Dashboard Load Time | <3s for executive dashboards |
| Report Generation (small) | <10s for <1k rows |
| Report Generation (large) | Async with email notification |
| Data Freshness | 15 min lag acceptable for dashboards |

## Acceptance Criteria

- [ ] Executive dashboards show real-time KPIs
- [ ] Claims analytics report provides drill-down capability
- [ ] Partner performance reports available to Partner Admins
- [ ] Regulatory reports generated in compliant formats
- [ ] All report access logged and auditable

## Traceability

**SRS Reference:** FG-018 — 4.18 Analytics & Reporting  
**Functional Requirements:** FR-149, FR-150, FR-151, FR-152, FR-202, FR-203, FR-204

[[[PAGEBREAK]]]


---
[[[PAGEBREAK]]]


# Feature Group: Audit & Logging (FG-019)

## Business Objective

Maintain immutable audit logs for all critical operations to support compliance audits, fraud investigations, regulatory inquiries, and dispute resolution. Ensure long-term retention and retrieval capabilities.

**Business Value:**
- Regulatory compliance (audit trail requirements)
- Fraud investigation support (who did what, when)
- Dispute resolution evidence (customer complaints, payment disputes)
- Operational transparency and accountability

## Actors & Portals

| Actor | Portal(s) | Responsibilities |
|-------|-----------|------------------|
| Compliance Officer | Admin Portal | Review audit logs for compliance audits |
| Fraud Analyst | Admin Portal | Investigate suspicious activities using logs |
| System Admin | Admin Portal | Monitor system health, troubleshoot issues using logs |
| Regulator (via controlled access) | Regulatory Access Portal | Request and review audit extracts |

## User Stories

### US-FG-019-01: Comprehensive Audit Logging

**As a** compliance officer  
**I want** all critical actions logged immutably  
**So that** audits and investigations have complete evidence

**Acceptance Criteria:**
- Logged events include: user actions (login, purchase, claim submission), admin actions (config changes, overrides), system events (payment confirmations, policy issuance)
- Each log entry contains: timestamp, user_id, action, entity_id, before/after state, IP address, result (success/failure)
- Logs stored in append-only format (immutable)
- Logs retained per regulatory requirement (typically 7 years)

![Audit Log Viewer](images/admin_audit_log_viewer.png)

**Related FRs:** FR-153, FR-155, FR-156

### US-FG-019-02: Audit Log Search & Filter

**As a** fraud analyst  
**I want** to search audit logs by user, action, date  
**So that** I can investigate suspicious patterns

**Acceptance Criteria:**
- Search by: user_id, action_type, entity_id (policy/claim/payment), date range
- Results paginated (1000 entries per page)
- Export search results as CSV
- Advanced filters: IP address, device ID, result (success/failure)

**Related FRs:** FR-154, FR-155

### US-FG-019-03: Regulatory Audit Extract

**As a** compliance officer  
**I want** to generate audit extracts for regulator requests  
**So that** compliance is streamlined

**Acceptance Criteria:**
- Officer selects date range + entity types (policies, claims, payments)
- System generates extract with all relevant logs
- Extract includes: summary report + detailed CSV
- Extract generation logged (who requested, what data, when)
- Regulator access requires approval workflow

![Regulatory Audit Extract](images/report_audit_extract.png)

**Related FRs:** FR-156, FR-157

### US-FG-019-04: Data Access Logging

**As a** system administrator  
**I want** all data access logged  
**So that** unauthorized access is detected

**Acceptance Criteria:**
- Every database query logged (user, query type, entity accessed)
- Admin portal access logged (which admin accessed which section)
- Customer PII access flagged for review
- Access logs retained separately with high security

**Related FRs:** FR-158

## Business Rules

| Rule ID | Description |
|---------|-------------|
| BR-AUD-01 | All critical actions logged (policy, claim, payment, admin config) |
| BR-AUD-02 | Logs immutable (append-only, no deletion) |
| BR-AUD-03 | Log retention: 7 years minimum (configurable per jurisdiction) |
| BR-AUD-04 | PII access logged separately with additional security |
| BR-AUD-05 | Audit log access itself audited (who accessed audit logs) |

## Key Workflows

### Automated Audit Logging
1. User/system performs action (e.g., claim submission)
2. Action intercepted by audit logging service
3. Log entry created: timestamp, actor, action, entity, context, result
4. Log entry written to append-only storage
5. Log indexed for search

### Fraud Investigation Using Logs
1. Fraud Analyst identifies suspicious claim
2. Opens audit log viewer → searches by claim_id
3. Reviews all actions on that claim (who submitted, reviewed, approved)
4. Checks related actions (same user, same device, IP patterns)
5. Exports relevant logs for investigation report

### Regulatory Audit Request
1. Regulator requests audit data (via formal letter)
2. Compliance Officer creates extraction request in system
3. Focal Person or Business Admin approves
4. System generates extract (date range, entity types)
5. Extract reviewed internally → delivered to regulator
6. Extraction logged with regulator identity and data delivered

## Data Model Notes

**Audit Log Entry (SRS Proto: insuretech.common.v1.AuditLog)**
- log_id (UUID)
- timestamp
- actor_id (user/system)
- actor_type (CUSTOMER, AGENT, ADMIN, SYSTEM)
- action_type (LOGIN, PURCHASE, CLAIM_SUBMIT, CONFIG_CHANGE, etc.)
- entity_type (POLICY, CLAIM, PAYMENT, USER, CONFIG)
- entity_id
- before_state (JSON snapshot)
- after_state (JSON snapshot)
- result (SUCCESS, FAILURE, PENDING)
- ip_address, device_id
- tenant_id

**Audit Log Access Log (meta-audit)**
- access_log_id
- accessor_user_id
- search_filters (JSON)
- accessed_at
- reason

## Integration Touchpoints

| External System | Purpose | Failure Handling |
|----------------|---------|------------------|
| Log Storage (S3/Minio, append-only) | Long-term log retention | Redundant storage, backup |
| SIEM (future) | Security monitoring and alerting | Best effort integration |

## Security & Privacy

- Audit logs encrypted at rest
- Access to audit logs restricted (Compliance Officer, System Admin only)
- PII in logs encrypted with separate key
- Log tampering detected via checksums/hashing

## NFR Constraints

| NFR | Target |
|-----|--------|
| Log Write Latency | <100ms (asynchronous to avoid blocking operations) |
| Log Retention | 7 years minimum |
| Log Search Performance | <5s for typical queries (<1M entries) |
| Log Storage Growth | ~10GB/month (estimate, varies by transaction volume) |

## Acceptance Criteria

- [ ] All critical actions logged with complete context
- [ ] Logs immutable and retained for required period
- [ ] Audit logs searchable and exportable
- [ ] Regulatory audit extracts generated with approval workflow
- [ ] Data access logging prevents unauthorized access

## Traceability

**SRS Reference:** FG-019 — 4.19 Audit & Logging  
**Functional Requirements:** FR-153, FR-154, FR-155, FR-156, FR-157, FR-158

[[[PAGEBREAK]]]


---
[[[PAGEBREAK]]]


# Feature Group: Fallback & Resilience (FG-020)

## Business Objective

Ensure business continuity when external dependencies (payment gateways, NID verification, pricing APIs) are unavailable through manual fallback workflows, queue management, and graceful degradation. Minimize revenue loss and maintain customer experience during outages.

**Business Value:**
- Minimize revenue loss during partner/provider outages
- Maintain customer experience (fallback to manual verification vs hard failure)
- Operational resilience (business continues with degraded performance)
- SLA accountability (track external provider reliability)

## Actors & Portals

| Actor | Portal(s) | Responsibilities |
|-------|-----------|------------------|
| Business Admin | Admin Portal | Manage fallback queues, manual verification workflows |
| Operations Team | Admin Portal | Monitor external service health, trigger fallback modes |
| Customer (indirect) | Mobile App | Experience graceful degradation (manual verification instead of failure) |

## User Stories

### US-FG-020-01: Payment Gateway Fallback

**As a** customer  
**I want** to complete my purchase even if the payment gateway is down  
**So that** I don't lose the opportunity to buy

**Acceptance Criteria:**
- If payment gateway unavailable → system offers "Manual Payment" option
- Customer can upload bank slip or transaction screenshot
- Business Admin verifies payment manually within 24 hours
- On verification → policy issued
- Customer notified of manual verification timeline

![Payment Fallback Flow](images/flow_payment_fallback.png)

**Related FRs:** FR-225, FR-227, FR-228

### US-FG-020-02: NID Verification Fallback

**As a** customer  
**I want** my policy approved even if NID verification API is down  
**So that** my purchase is not blocked

**Acceptance Criteria:**
- If NID API unavailable → policy marked "Pending Verification"
- Customer allowed to proceed with purchase
- Business Admin manually verifies NID documents later
- On verification → policy status updated to "Active"
- Customer notified of verification completion

**Related FRs:** FR-209, FR-210

### US-FG-020-03: Pricing API Fallback (Quote Generation)

**As a** system  
**I want** to use cached pricing when pricing API is down  
**So that** customers can still get quotes

**Acceptance Criteria:**
- If pricing API unavailable → use last known pricing (cached)
- Quote marked "Indicative (pending final price)"
- When API recovers → pricing re-validated before policy issuance
- If price difference >5% → customer notified and offered to cancel/proceed

**Related FRs:** FR-214, FR-215, FR-216

### US-FG-020-04: Queue Management for Fallback Items

**As a** Business Admin  
**I want** a queue of items requiring manual verification  
**So that** fallback workflows are tracked and resolved

**Acceptance Criteria:**
- Queue shows: item type (payment, NID, endorsement), customer, date submitted, priority
- Admin can process items in queue (verify, approve, reject)
- Queue filtered by status (PENDING, IN_PROGRESS, RESOLVED)
- SLA tracking (items in queue >24 hours flagged)

![Fallback Queue](images/admin_fallback_queue.png)

**Related FRs:** FR-211, FR-212

## Business Rules

| Rule ID | Description |
|---------|-------------|
| BR-FBK-01 | External service failure triggers fallback mode automatically |
| BR-FBK-02 | Fallback items queued for manual processing within 24 hours |
| BR-FBK-03 | Cached pricing valid for 24 hours max |
| BR-FBK-04 | Customer notified when fallback mode activated |
| BR-FBK-05 | System monitors external service health (ping/heartbeat every 5 min) |

## Key Workflows

### Payment Gateway Fallback
1. Customer initiates payment → gateway timeout/error
2. System detects failure → offers manual payment option
3. Customer uploads payment proof
4. Business Admin receives verification queue notification
5. Admin verifies proof → marks payment as confirmed
6. Policy issued → customer notified

### NID Verification Fallback
1. Customer submits policy purchase with NID
2. System calls NID API → API unavailable
3. System marks policy "Pending NID Verification" → allows purchase to proceed
4. Business Admin manually reviews NID documents
5. Admin approves → policy status → "Active"

### Service Health Monitoring
1. System pings external services (NID, payment, pricing) every 5 min
2. If service down → system switches to fallback mode
3. Operations Team notified
4. When service recovers → system switches back to normal mode
5. Fallback queue items re-processed automatically

## Data Model Notes

**Fallback Queue Item**
- queue_id
- item_type (PAYMENT_VERIFICATION, NID_VERIFICATION, ENDORSEMENT_APPROVAL)
- entity_id (policy/payment/claim)
- customer_id
- status (PENDING, IN_PROGRESS, RESOLVED, REJECTED)
- submitted_at
- resolved_at
- resolved_by (admin_user_id)

**Service Health Status**
- service_name (NID_API, PAYMENT_GATEWAY, PRICING_API)
- status (UP, DOWN, DEGRADED)
- last_check_at
- failure_count
- fallback_mode_active (boolean)

## Integration Touchpoints

| External System | Purpose | Failure Handling |
|----------------|---------|------------------|
| Payment Gateways | Premium collection | Manual payment verification queue |
| NID Verification API | Identity verification | Manual document review queue |
| Pricing APIs | Calculate premiums | Cached pricing (24 hour validity) |

## Security & Privacy

- Fallback queues access-controlled (only authorized admins)
- Manual verification actions logged for audit
- Customer notified when entering fallback mode (transparency)

## NFR Constraints

| NFR | Target |
|-----|--------|
| Service Health Check Frequency | Every 5 min |
| Fallback Mode Activation | <1 min from service failure detection |
| Manual Verification SLA | <24 hours (business days) |
| Cached Pricing Validity | 24 hours max |

## Acceptance Criteria

- [ ] Payment gateway failure triggers manual verification flow
- [ ] NID verification failure allows purchase with pending status
- [ ] Pricing API failure uses cached pricing with customer notification
- [ ] Fallback queue visible and processable by Business Admin
- [ ] Service health monitored and fallback mode automated

## Traceability

**SRS Reference:** FG-020 — 4.20 System Interface Architecture  
**Functional Requirements:** FR-193, FR-194, FR-195, FR-196, FR-197, FR-198, FR-199, FR-200, FR-201, FR-159, FR-160, FR-161, FR-162, FR-163, FR-223, FR-224, FR-225, FR-226, FR-227, FR-228, FR-229, FR-230, FR-231

[[[PAGEBREAK]]]


---
[[[PAGEBREAK]]]


# Feature Group: Document Management (FG-022)

## Business Objective

Manage policy documents, receipts, endorsements, claims documents, and other records with versioning, secure storage, retrieval, and verification capabilities. Ensure long-term retention for compliance and customer access.

**Business Value:**
- Digital-first reduces paper costs and operational overhead
- Verifiable documents (QR codes) build customer trust
- Long-term retention supports compliance (7+ years)
- Fast retrieval improves customer experience and dispute resolution

## Actors & Portals

| Actor | Portal(s) | Responsibilities |
|-------|-----------|------------------|
| Customer | Mobile App, Web | View/download policy documents, receipts, claim documents |
| System (Automated) | Background Services | Generate documents (policy, receipt, endorsement), store securely |
| Business Admin | Admin Portal | Manage document templates, manual document generation/upload |
| Compliance Officer | Admin Portal | Audit document retention, retrieve archived documents |

## User Stories

### US-FG-022-01: Digital Policy Document Generation

**As a** customer  
**I want** my policy document generated and delivered instantly  
**So that** I have proof of coverage immediately

**Acceptance Criteria:**
- Policy document generated within 5 minutes of payment confirmation
- Document includes: policy number, customer details, coverage, premium, terms, QR code
- QR code links to verification portal
- Document delivered via SMS (link) + email (PDF attachment)
- Document accessible in app under "My Policies"

![Policy Document Sample](images/document_policy_sample.png)

**Related FRs:** FR-237, FR-238

### US-FG-022-02: Document Verification via QR Code

**As a** hospital staff or service provider  
**I want** to verify a customer's policy by scanning QR code  
**So that** I can confirm coverage before providing service

**Acceptance Criteria:**
- QR code scanned → redirects to verification portal
- Verification portal shows: policy number, customer name, coverage status (ACTIVE/LAPSED), coverage details
- No sensitive PII exposed (address, NID not shown)
- Verification logged for audit

**Related FRs:** FR-220, FR-221

### US-FG-022-03: Document Version History

**As a** customer  
**I want** to access previous versions of my policy document  
**So that** I can see what changed after endorsements

**Acceptance Criteria:**
- Policy detail page shows version history
- Customer can download any previous version
- Version metadata: date, change type (original, endorsement, renewal)
- Endorsement documents linked to base policy

**Related FRs:** FR-222

### US-FG-022-04: Long-Term Document Archival

**As a** compliance officer  
**I want** documents retained and retrievable for 7+ years  
**So that** regulatory requirements are met

**Acceptance Criteria:**
- All documents stored in redundant, durable storage (S3/Minio with replication)
- Documents indexed for fast retrieval by policy_id, customer_id, date
- Archived documents accessible via Admin Portal
- Retrieval time: <30s for documents within retention period

![Document Archive Viewer](images/admin_document_archive.png)

**Related FRs:** FR-223, FR-224

## Business Rules

| Rule ID | Description |
|---------|-------------|
| BR-DOC-01 | Policy documents generated within 5 min of issuance |
| BR-DOC-02 | Documents versioned on every change (endorsement, renewal) |
| BR-DOC-03 | QR codes use signed tokens (tamper-proof, expiry: 5 years) |
| BR-DOC-04 | Document retention: 7 years minimum (configurable per jurisdiction) |
| BR-DOC-05 | Customer can download documents anytime during and after policy period |

## Key Workflows

### Policy Document Generation & Delivery
1. Policy issued (payment confirmed)
2. Document generation service triggered
3. System populates template with policy data
4. PDF generated with embedded QR code
5. Document stored in secure storage
6. SMS + email sent to customer with download link
7. Document visible in app

### Document Verification (QR Scan)
1. Service provider scans QR code on policy document
2. QR redirects to verification portal
3. System validates QR signature → retrieves policy data
4. Verification page displays: policy status, coverage summary
5. Verification event logged

### Document Archival & Retrieval
1. Documents older than 1 year moved to cold storage
2. Indexed metadata retained in database
3. Compliance Officer searches for archived document
4. System retrieves from cold storage → presents for download
5. Retrieval logged for audit

## Data Model Notes

**Document Entity**
- document_id
- document_type (POLICY, RECEIPT, ENDORSEMENT, CLAIM_DOCUMENT, CANCELLATION)
- policy_id (or claim_id, payment_id)
- customer_id
- version
- status (DRAFT, FINAL, ARCHIVED)
- storage_url (S3/Minio path)
- qr_code_token (signed)
- generated_at
- archived_at

**Document Template**
- template_id
- template_name (policy_document_v2, receipt_template)
- template_type (PDF, HTML)
- template_content (with placeholders)
- version

## Integration Touchpoints

| External System | Purpose | Failure Handling |
|----------------|---------|------------------|
| Document Storage (S3/Minio) | Secure, durable storage | Redundant storage, automatic retry |
| PDF Generation Service | Create PDFs from templates | Queue for retry, fallback to pre-generated templates |
| SMS/Email Gateway | Deliver document links | Queue for retry |

## Security & Privacy

- Documents encrypted at rest (AES-256)
- Access controlled by role and ownership (customer sees only their docs)
- QR codes use signed tokens (prevent forgery)
- Document downloads logged for audit

## NFR Constraints

| NFR | Target |
|-----|--------|
| Document Generation Time | <30s per document |
| Document Delivery | <5 min from generation to customer receipt |
| Document Retrieval (active) | <5s |
| Document Retrieval (archived) | <30s |
| Storage Durability | 99.999999999% (S3 standard) |

## Acceptance Criteria

- [ ] Policy documents generated and delivered within 5 minutes
- [ ] QR code verification works for service providers
- [ ] Document version history accessible to customers
- [ ] Documents retained and retrievable for 7+ years
- [ ] All document access logged

## Traceability

**SRS Reference:** FG-022 — 4.22 Data Storage  
**Functional Requirements:** FR-232, FR-233, FR-234, FR-235, FR-236, FR-237, FR-238, FR-239, FR-240, FR-241, FR-242, FR-243

[[[PAGEBREAK]]]


---
[[[PAGEBREAK]]]


# Feature Group: Client UX/UI (FG-023)

## Business Objective

Deliver consistent, intuitive, and accessible user interfaces across web and mobile channels with multi-language support (Bengali/English), responsive design, and inclusive accessibility to maximize conversion and customer satisfaction.

**Business Value:**
- Reduce drop-off through intuitive, guided UX (target: <20% abandonment in purchase funnel)
- Multi-language support expands addressable market (Bengali-first for mass adoption)
- Responsive design works on low-end devices (cost barrier reduced)
- Accessibility compliance enables inclusive insurance

## Actors & Portals

| Actor | Portal(s) | Responsibilities |
|-------|-----------|------------------|
| Customer | Mobile App (iOS/Android), Web PWA | Primary user journeys (register, purchase, claims, support) |
| Agent | Agent Mobile App | Assisted sales, customer onboarding |
| Partner Admin | Partner Portal (Web) | Dashboard, reports, configuration |
| Business Admin | Admin Portal (Web) | Configuration, approvals, reports |

## User Stories

### US-FG-023-01: Mobile-First Responsive Design

**As a** customer on a low-end smartphone  
**I want** the app to load quickly and work smoothly  
**So that** I can complete transactions without frustration

**Acceptance Criteria:**
- App works on devices with 2GB RAM, Android 8+
- Page load time: <3s on 3G connection
- Offline capability: view cached policies, claim status
- Minimal data usage (images optimized, lazy loading)

![Mobile Responsive Design](images/ui_mobile_responsive.png)

**Related FRs:** FR-245, FR-246

### US-FG-023-02: Multi-Language Support (Bengali/English)

**As a** Bengali-speaking customer  
**I want** the entire app in Bengali  
**So that** I understand all terms and instructions

**Acceptance Criteria:**
- Language toggle in app settings (Bengali ↔ English)
- All UI text, buttons, labels translated
- Product descriptions, T&C, FAQs available in both languages
- Language preference saved per user
- Default: Bengali (auto-detect based on device locale)

![Language Toggle](images/ui_language_toggle.png)

**Related FRs:** FR-244, FR-245

### US-FG-023-03: Guided Purchase Flow (Step-by-Step)

**As a** customer new to insurance  
**I want** a guided purchase flow with clear steps  
**So that** I don't get lost or make mistakes

**Acceptance Criteria:**
- Purchase flow: Product Selection → Details Entry → Nominee Setup → Review → Payment → Confirmation
- Progress indicator shows current step (e.g., "Step 2 of 6")
- Each step has clear instructions + inline help
- Form validation with friendly error messages
- "Back" button allows correction without losing data

![Guided Purchase Flow](images/flow_guided_purchase.png)

**Related FRs:** FR-246

### US-FG-023-04: Accessibility (Screen Reader Support)

**As a** visually-impaired customer  
**I want** the app to work with screen readers  
**So that** I can use insurance services independently

**Acceptance Criteria:**
- All buttons, form fields have descriptive labels (ARIA)
- Screen reader announces page changes and form errors
- Color contrast meets WCAG AA standards (4.5:1 for text)
- Keyboard navigation works for all interactive elements
- Voice-over tested on iOS, TalkBack on Android

**Related FRs:** FR-247

### US-FG-023-05: Consistent Design System

**As a** developer  
**I want** a shared design system across all portals  
**So that** UX is consistent and development is faster

**Acceptance Criteria:**
- Design system includes: color palette, typography, button styles, form components, icons
- Documented in style guide (Storybook or similar)
- Reusable components for Customer App, Partner Portal, Admin Portal
- Dark mode support (future)

![Design System](images/design_system_components.png)

**Related FRs:** FR-248

## Business Rules

| Rule ID | Description |
|---------|-------------|
| BR-UX-01 | Bengali is default language (can be overridden by user) |
| BR-UX-02 | Mobile app works on Android 8+ and iOS 13+ |
| BR-UX-03 | Accessibility: WCAG 2.1 Level AA compliance |
| BR-UX-04 | Page load time: <3s on 3G (target) |
| BR-UX-05 | Critical user flows tested on low-end devices (2GB RAM) |

## Key Workflows

### Language Switching
1. User opens app → default language Bengali (or device locale)
2. User taps language toggle (flag icon or "Language" in settings)
3. User selects English
4. App reloads with English text
5. Preference saved → persists across sessions

### Guided Purchase (Mobile)
1. User taps "Buy Insurance"
2. Step 1: Product selection (cards with clear descriptions)
3. Step 2: Customer details (auto-filled if logged in)
4. Step 3: Nominee setup (inline validation)
5. Step 4: Review summary (edit option for each section)
6. Step 5: Payment (multi-channel options)
7. Step 6: Confirmation (download policy, share option)

### Accessibility Testing
1. QA team tests with screen readers (VoiceOver, TalkBack)
2. Tests keyboard navigation (tab order, focus indicators)
3. Tests color contrast (automated tool + manual review)
4. User testing with visually-impaired users (beta program)

## Data Model Notes

**User Preferences**
- user_id
- language (bn, en)
- theme (light, dark - future)
- accessibility_mode (boolean)

## Integration Touchpoints

| External System | Purpose | Failure Handling |
|----------------|---------|------------------|
| CDN | Serve static assets (images, fonts) | Fallback to local cache |
| Translation Service (future) | Dynamic translations | Pre-built translation files as fallback |

## Security & Privacy

- No sensitive data cached locally (policies encrypted if cached offline)
- User preferences stored securely (encrypted)
- Third-party UI libraries vetted for security vulnerabilities

## NFR Constraints

| NFR | Target |
|-----|--------|
| Page Load Time (3G) | <3s |
| App Size (APK/IPA) | <50MB (initial download) |
| Frame Rate | 60 FPS for animations |
| Accessibility Compliance | WCAG 2.1 Level AA |
| Multi-Language Coverage | 100% of user-facing text |

## Acceptance Criteria

- [ ] Mobile app works smoothly on low-end devices (2GB RAM)
- [ ] All text available in Bengali and English
- [ ] Purchase flow guided with clear step indicators
- [ ] Screen readers work for all critical flows
- [ ] Design system implemented and documented

## Traceability

**SRS Reference:** FG-023 — 4.23 User Interface Requirements  
**Functional Requirements:** FR-244, FR-245, FR-246, FR-247, FR-248

[[[PAGEBREAK]]]


---
[[[PAGEBREAK]]]


---

## Document Generation Information

**Generated:** 2025-12-19 10:40:12  
**Generator:** BRD V3.7 Merge Script  
**Source:** Modular sections in BRD/sections/  

---

**End of Document**
