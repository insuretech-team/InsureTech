# LabAid InsureTech Platform
## System Requirements Specification (SRS)

| Field | Details |
|-------|---------|
| **Version** | 3.11 |
| **Date** | Feb 2026 |
| **Status** | FINAL_DRAFT |
| **Company** | LabAid InsureTech |
| **Technology Partner** | LifePlus |

---

## Revision History

| Version | Date | Revised By | Description |
|---------|------|------------|-------------|
| 1.0 | Dec 2025 | Director InsureTech | Initial SRS with core business requirements |
| 2.0 | Dec 2025 | Faruk Hannan | Technical architecture and detailed requirements with MD Sirs Feedback |
| 2.2 | Dec 2025 | AI Engine | Enhanced Formatting, Grammar, Fact check |
| 3.0 | Dec 2025 | Faruk Hannan | Final SPEC Draft with proto models, VSA architecture, and additional requirements |
| 3.1 | Dec 2025 | Faruk Hannan | Final SPEC with MD Sirs Feedback M1.5 AI features and 10% Buffer -iteration1 |
| 3.2 | Dec 2025 | Faruk Hannan | Final SPEC with MD Sirs Feedback M1.5 AI features and 10% Buffer -iteration2 |
| 3.3 | Dec 2025 | Faruk Hannan | Final SPEC with MD Sirs Feedback M1.5 AI features and 10% Buffer -iteration3 |
| 3.4 | Dec 2025 | AI Engine | Formatting, Diagrams, Enhancements |
| 3.5 | Dec 2025 | JOY, NOOR | Feedback and plan to spread out milestones |
| 3.6 | Dec 2025 | SABBIR | Formatting |
| 3.7 | Dec 2025 | FARUK HANNAN | Reorganised priorities and added missing proto service definitions |
| 3.8 | Dec 2025 | CEO (LifePlus), Director InsureTech, FARUK HANNAN, Dev Team, Business Head | Corrected some discussed points |
| 3.9 | Jan 2026 | CEO (LifePlus), Director InsureTech, FARUK HANNAN, Dev Team, Business Head | Revised some requirements |
| 3.10 | Jan 2026 | Sabbir | FR ID correction |
| 3.11 | Feb 2026 | Imtiaz, Sabbir | Adjusted Requirements in FR-011, Adjusted FR-023-A and FR-023-B in section 4.3, Adjusted FR-076 in section 4.7, Adjusted FR-032-A in section 4.4, Adjusted FR-070 in section 4.7 |

---

## Executive Summary

This System Requirements Specification (SRS) defines the functional and non-functional requirements of the LabAid InsureTech Platform — a cloud-native, microservices-based system enabling end-to-end digital insurance for the Bangladesh market. It covers onboarding and KYC, product discovery and quotation, policy lifecycle management, payments and reconciliation, claims management, reporting, and regulatory compliance.

The SRS is plan-agnostic and team-agnostic. It specifies what the system must do and the quality attributes it must meet, independent of delivery timelines, resource assignments, or milestone planning (captured in separate BRD/Planning documents).

### Key Themes

- **Digital-first experience:** Mobile and web channels with Bangladesh-optimized UX and language support.
- **Compliance by design:** IDRA/BFIU-aligned data, auditability, and reporting.
- **Interoperability:** Clear interfaces for identity/KYC, payments, messaging, and health systems.
- **Security and privacy:** Zero-trust posture, encryption, least-privilege authorization, and governed data handling.
- **Observability and reliability:** Logging, metrics, tracing, and service health for regulated operations.

### Key Changes in v3.7

- Consolidated and de-duplicated functional requirements with continuous FR IDs.
- Expanded security and compliance requirements with cross-references.
- Proto-first interfaces organized by domain and included in appendices with examples.
- Clear separation of integration details under dedicated Integration section and references from FRs.

### Phase Definitions

| Phase | Date | Description |
|-------|------|-------------|
| M1 | March 1, 2025 | Soft Launch - National Insurance Day |
| M2 | April 14, 2025 | Grand Launch with critical features |
| M3 | August 1, 2025 | Upgrade Release features — LEMONADE V0.5 |
| D | October 1, 2025 | Enhance Tech Release features |
| S | November 1, 2025 | Scaling features |
| F | January 1, 2027 | Expansion features |

---

## Table of Contents

1. [Introduction](#1-introduction)
2. [System Overview](#2-system-overview)
3. [System Architecture](#3-system-architecture)
4. [System Features & Functional Requirements](#4-system-features--functional-requirements)
5. [Non-Functional Requirements & Technical Constraints](#5-non-functional-requirements--technical-constraints)
6. [Data Model & Persistence](#6-data-model--persistence)
7. [Security & Compliance Requirements](#7-security--compliance-requirements)
8. [Integration Requirements](#8-integration-requirements)
9. [Performance & Monitoring](#9-performance--monitoring)
10. [Support & Maintenance](#10-support--maintenance)
11. [Acceptance Criteria & Test Requirements](#11-acceptance-criteria--test-requirements)
12. [Traceability Matrix & Change Control](#12-traceability-matrix--change-control)
13. [Appendices](#appendices)

---

## 1. Introduction

### 1.1 Purpose

This System Requirements Specification (SRS) document serves as the definitive technical blueprint for the LabAid InsureTech Platform. It provides comprehensive requirements for developers, testers, business analysts, and stakeholders to ensure successful delivery of a world-class digital insurance platform tailored for the Bangladesh market.

### 1.2 Scope

**In Scope:**
- Customer mobile application (Android, iOS)
- Customer web portal (React PWA)
- Partner portal (hospitals, MFS, e-commerce)
- Admin portal (multi-role: System Admin, Business Admin, Focal Person, Support)
- Backend microservices (Insurance Engine, Partner Management, AI Engine, Gateway, Kafka Orchestration, Ticketing, Analytics)
- Digital onboarding with KYC/KYB
- Product catalog and policy purchase
- Payment processing (bKash, Nagad, Rocket, manual Phase 1; all channel automated + manual Phase 2)
- Claims submission and approval workflows (with basic AI Engine check Phase 1; all channel automated + manual Phase 2)
- Partner management and tenant isolation
- Notification system (SMS, Email, Push)
- Reporting and analytics
- IDRA and BFIU compliance features
- Voice-assisted workflow (Bengali speech recognition, voice-guided policy purchase, voice claims submission)
- AI-based Claim Management (fraud detection, automated assessment, risk scoring)
- IoT-based Usage-Based Insurance (UBI) for vehicles and health tracking
- IoT-based Tracking system

**Out of Scope:**
- Full AI-driven underwriting (Phase 2.5/3)
- Universal IoT/Telematics integration (Phase 2.5/3)
- Cross-border insurance (Future consideration)
- Blockchain-based smart contracts (Future consideration)

### 1.3 Definitions, Acronyms & Abbreviations

| Term | Definition |
|------|------------|
| IDRA | Insurance Development & Regulatory Authority of Bangladesh |
| BFIU | Bangladesh Financial Intelligence Unit |
| MFS | Mobile Financial Services (bKash, Nagad, Rocket) |
| KYC | Know Your Customer |
| AML | Anti-Money Laundering |
| CFT | Combating the Financing of Terrorism |
| gRPC | Google Remote Procedure Call |
| CQRS | Command Query Responsibility Segregation |
| VSA | Vertical Slice Architecture |
| Proto | Protocol Buffers |
| IoT | Internet of Things |
| AI | Artificial Intelligence |
| ML | Machine Learning |
| API | Application Programming Interface |
| SLA | Service Level Agreement |
| TAT | Turn Around Time |
| EHR | Electronic Health Records |
| OCR | Optical Character Recognition |
| SMS | Short Message Service |
| OTP | One-Time Password |
| JWT | JSON Web Token |
| RBAC | Role-Based Access Control |

**Business Terms:**
- **Policyholders:** Individual customers who purchase insurance policies
- **Partners:** System-onboarded, collaborative organizations (MFS, hospitals, e-commerce)
- **Agents:** Individual sales representatives working with partners
- **Focal Persons:** Partner organization representatives managing agent networks
- **Sum Assured:** Maximum coverage amount for a policy
- **Premium:** Insurance policy cost paid by policyholder
- **Claims:** Requests for insurance payouts due to covered incidents
- **Underwriting:** Risk assessment process for policy approval
- **Reinsurance:** Insurance purchased by insurance company to limit risk

---

## 2. System Overview

### 2.1 Business Context

The LabAid InsureTech Platform addresses the significant insurance gap in Bangladesh, where traditional insurance penetration remains below 1% of GDP. The platform focuses on micro-insurance products (200–2,000 BDT premiums) to make insurance accessible to the mass market through digital channels and strategic partnerships.

**Market Opportunity:**
- **Target Market:** 165+ million Bangladeshi citizens with mobile phones
- **Insurance Gap:** 99% of population lacks adequate insurance coverage
- **Digital Adoption:** 50%+ smartphone penetration with growing digital payment usage
- **Regulatory Support:** IDRA's digital insurance initiatives and sandbox programs

**Business Model:**
- **B2B Approach:** Partner with MFS providers, hospitals, e-commerce platforms
- **B2C Approach:** Making Microinsurance products affordable for average Bangladeshi families
- **Commission-Based Revenue:** 15–25% commission on premium collections
- **Value-Added Services:** IoT risk assessment, AI-powered customer support
- **Data Monetization:** Anonymized insights for partners (IDRA compliant)

### 2.2 System Objectives

**Primary Objectives:**
- **Digital Insurance Ecosystem:** End-to-end digital insurance journey from discovery to settlement
- **Mass Market Accessibility:** Micro-insurance products affordable for average Bangladeshi families
- **Partner Network Growth:** Scalable platform supporting 100+ partners and 10,000+ agents
- **Regulatory Excellence:** Full IDRA and BFIU compliance with automated reporting
- **Operational Efficiency:** 90%+ automated processes with minimal manual intervention

**Success Metrics:**
- Policy Issuance: 10,000 policies/month by Month 6
- Claims Settlement: 95% claims settled within 72 hours
- Customer Satisfaction: 4.5+ star rating on app stores
- Partner Growth: 50+ active partners by Year 1
- Revenue Target: 10 Crore BDT annual premium by Year 2

### 2.3 Key Stakeholders

**Internal Stakeholders:**
- Business Executives: Strategic oversight and P&L responsibility
- Product Management: Feature prioritization and roadmap planning
- Technology Leadership: Architecture and development oversight
- Compliance Team: Regulatory adherence and audit management
- Customer Success: Partner relationship and customer experience

**External Stakeholders:**
- IDRA: Regulatory approval and ongoing compliance monitoring
- Partners: MFS providers, hospitals, e-commerce platforms
- Customers: Individual policyholders and their families
- Technology Vendors: Cloud providers, payment gateways, third-party services

### 2.4 System Boundaries

**System Includes:**
- Insurance policy lifecycle management
- Customer onboarding and KYC verification
- Premium collection and payment processing
- Claims submission, assessment, and settlement
- Partner and agent management
- Regulatory reporting and compliance
- Customer support and AI assistance

**System Interfaces:**
- Mobile applications (Android/iOS)
- Web portals (customer, partner, admin)
- REST and gRPC APIs for integrations
- Payment gateway connections
- Government database integrations (NID, mobile verification)
- Third-party services (SMS, email, cloud storage)

**External Dependencies:**
- Bangladesh government services (NID verification, mobile number validation)
- Payment providers (bKash, Nagad, banks)
- Cloud infrastructure (AWS/Azure)
- Regulatory reporting systems
- Partner systems and databases

---

## 3. System Architecture

### 3.1 Architectural Overview

The LabAid InsureTech Platform is built on a cloud-native, microservices architecture with Domain-Driven Design (DDD) principles. The system leverages Vertical Slice Architecture (VSA) pattern across all services for maximum cohesion and maintainability.

**Core Architectural Principles:**
- **Microservices First:** Independent, deployable services with single responsibilities
- **Event-Driven Architecture:** Asynchronous communication through domain events
- **API-First Design:** All services expose well-defined APIs (REST/gRPC)
- **Cloud-Native:** Built for containerization and orchestration
- **Security by Design:** Zero-trust security model with end-to-end encryption

### 3.2 Technology Stack

**Languages & Frameworks:**
- **Go:** Gateway, Authentication, Authorization, DBManager, Storage, IoT Broker, Kafka Services
- **C# .NET 8:** Insurance Engine, Partner Management, Analytics & Reporting
- **Node.js:** Payment Service, Ticketing/Support Service
- **Python:** AI Engine, OCR/PDF Processing
- **React:** All web portals and admin interfaces
- **React Native:** Mobile applications (Android/iOS)

**Data & Communication:**
- **Protocol Buffers:** All data models and service contracts
- **PostgreSQL:** Primary database for transactional data
- **MongoDB:** NoSQL for product catalogs and unstructured data
- **Apache Kafka:** Event streaming and service orchestration
- **gRPC:** Inter-service communication
- **REST APIs:** Client-facing and partner integrations

**Infrastructure:**
- **Docker:** Containerization
- **Kubernetes:** Container orchestration
- **AWS/Azure:** Cloud platform
- **Prometheus/Grafana:** Monitoring and alerting
- **Jaeger:** Distributed tracing
- **Redis:** Caching and session management

### 3.3 System Architecture — VSA Pattern

The LabAid InsureTech Platform adopts Vertical Slice Architecture (VSA) across ALL microservices, regardless of programming language:

- **Go Services:** Gateway, Auth, DBManager, Storage, IoT Broker, Kafka Orchestration
- **C# .NET Services:** Insurance Engine, Partner Management, Analytics & Reporting
- **Node.js Services:** Payment Service, Ticketing Service
- **Python Services:** AI Engine, OCR/PDF Service

**Key VSA Principles:**
- **High Cohesion:** Each slice contains all layers needed for one feature
- **Low Coupling:** Slices are independent and don't share logic
- **Feature-Focused:** Organized by business capability, not technical layer
- **Testability:** Each slice can be tested in isolation

### 3.4 Microservices Architecture

| Service | Language | Port | Responsibility |
|---------|----------|------|----------------|
| Gateway | Go | 8080 | API Gateway, routing, rate limiting |
| Auth Service | Go | 8081 | Authentication, JWT management |
| Authorization | Go | 8082 | RBAC, permissions, access control |
| DBManager | Go | 8083 | Database operations, migrations |
| Storage Service | Go | 8084 | File storage, S3 operations |
| IoT Broker | Go | 8085 | IoT device communication, MQTT |
| Kafka Service | Go | 8086 | Event orchestration, messaging |
| Insurance Engine | C# .NET | 5001 | Policy lifecycle, underwriting |
| Partner Management | C# .NET | 5002 | Partner onboarding, agent mgmt |
| Analytics & Reporting | C# .NET | 5003 | BI, dashboards, compliance reports |
| Payment Service | Node.js | 3001 | Payment processing, settlements |
| Ticketing Service | Node.js | 3002 | Customer support, help desk |
| AI Engine | Python | 4001 | LLM, chatbot, fraud detection |
| OCR Service | Python | 4002 | Document processing, KYC |

### 3.5 System Context Diagram

```
+----------------------------+            +----------------------------------+
|       External Systems     |            |     User Interfaces (Clients)    |
| - MFS (bKash/Nagad/Rocket) |<---------->| - Web (React PWA)                |
| - Hospital EHR (FHIR/HL7)  |  HTTPS/    | - Mobile (React Native)          |
| - NID Verification API     |  REST/     | - Partner/Admin Portals (React)  |
| - IDRA / BFIU Portals      |  gRPC/WS   | - Ops/Observability (Graf/Prom)  |
| - SMS / Email Gateway      |            +----------------------------------+
+----------------------------+
                 |
                 v
+--------------------------------------------------------------+
|                  API & Integration Layer                     |
|   - API Gateway (Go): routing, rate-limit, authN             |
|   - Auth Service (Go): OAuth2/OIDC, JWT, RBAC (authZ)        |
|   - Kafka (Go): domain events & async orchestration          |
+--------------------------------------------------------------+
                 |
                 | gRPC (Protobuf)  <-->  Kafka Events
                 v
+--------------------------------------------------------------+
|           Backend Microservices (VSA + gRPC)                 |
|   - Insurance Engine (C# .NET, CQRS/MediatR)                 |
|   - Partner/Agent Mgmt (C# .NET)                             |
|   - Analytics & Reporting (C# .NET)                          |
|   - Payment Service (Node.js)                                |
|   - Ticketing/Support (Node.js)                              |
|   - AI Engine (Python, LLM/Fraud)                            |
|   - OCR/PDF Service (Python)                                 |
|   - DBManager (Go)        - Storage (Go, S3/Object)          |
|   - IoT Broker (Go, MQTT/gRPC)                               |
+--------------------------------------------------------------+
                 |
                 v
+--------------------------------------------------------------+
|                        Data Layer                            |
|   - PostgreSQL 17 (Transactional, ACID)                      |
|   - MongoDB (Unstructured)  - Redis (Cache/Sessions)         |
|   - S3/Object Storage (Documents, Images)                    |
+--------------------------------------------------------------+
```

---

## 4. System Features & Functional Requirements

### Phase Definitions

| Phase | Date | Description |
|-------|------|-------------|
| M1 | March 1, 2025 | Must have for M1 launch (Soft Launch) |
| M2 | April 14, 2025 | Must have for M2 launch (Grand Launch) |
| M3 | August 1, 2025 | Must have for M3 Enhancement |
| D | October 1, 2025 | Desirable features |
| S | November 1, 2025 | Scalability |
| F | January 1, 2027 | Future enhancements |

---

### 4.1 User Management & Authentication (FG-001)

| FR ID | Requirement Description | Priority | Acceptance Criteria |
|-------|------------------------|----------|---------------------|
| FR-001 | The system shall support phone-based registration (Bangladesh mobile format: +880 1XXX XXXXXX) with OTP validation. User will get the option to login as Corporate ID and Individual phone number. User can use the same phone number against the corporate profile and individual profile. | M1 | OTP sent within 60s, 6-digit code valid for 5 minutes |
| FR-002 | The system shall send OTP via SMS within 60 seconds with 6-digit code valid for 5 minutes. We will not ask user for OTP every time. | M1 | 95% delivery success rate, retry on failure |
| FR-003 | The system shall allow maximum 3 OTP resend attempts per 15-minute window | M1 | Rate limiting enforced, user notified on limit |
| FR-004 | The system shall enforce unique mobile number per account and detect duplicate registrations. Mobile number will be unique. Family dependent max 6 members can be added. | M1 | Error message on duplicate, database constraint enforced |
| FR-005 | The system shall support email-based registration with email verification link (24-hour validity) | M2 | Verification email sent within 2 minutes, link expires after 24hrs |
| FR-006 | The system shall implement secure password policy: minimum 8 characters, 1 uppercase, 1 number, 1 special character | M1 | Password strength indicator shown, validation enforced |
| FR-007 | The system shall provide biometric authentication (fingerprint/face ID) for mobile users | D | Device biometric API integration, fallback to password |
| FR-008 | The system shall support password reset via OTP to registered mobile number | M1 | Reset OTP sent within 60s, new password saved securely |
| FR-009 | The system shall implement session management with Secure Token Service (15-minute access, 7-day refresh) | M1 | Token rotation implemented, refresh token stored securely |
| FR-010 | The system shall enforce account lockout after 5 failed login attempts for 30 minutes | M2 | Lockout triggered automatically, user notified via SMS |
| FR-011 | The system shall maintain user profile with: full name, date of birth, gender, occupation, income range, address. User can switch profile from Corporate to Individual and vice versa. | M1 | All mandatory fields validated, profile completeness indicator |
| FR-012 | The system shall support profile photo upload with validation (max 5MB, JPEG/PNG, face detection) | M3 | Image compressed to <2MB, face detection validates single face |
| FR-013 | The system shall have stakeholders registration via SAML Identity provider | D | SAML 2.0 integration with Azure AD/Okta, SSO enabled |

---

### 4.2 Authorization & Access Control (FG-002)

| FR ID | Requirement Description | Priority | Acceptance Criteria |
|-------|------------------------|----------|---------------------|
| FR-014 | The system shall implement Role-Based Access Control (RBAC) with predefined roles: System Admin, Business Admin, Focal Person, Partner Admin, Agent, Customer | M1 | Roles enforced at API gateway level, permissions validated on each request |
| FR-015 | The system shall enforce Attribute-Based Access Control (ABAC) for fine-grained permissions based on user attributes, resource type, and context | M1 | Dynamic policy evaluation <50ms, audit logs for all authorization decisions |
| FR-016 | The system shall implement tenant isolation for partner organizations with data segregation | M2 | Multi-tenant database architecture, row-level security enforced |
| FR-017 | The system shall enforce 2FA (Two-Factor Authentication) for all admin-level access | M3 | TOTP-based 2FA with 30-second rotation, backup codes provided |
| FR-018 | The system shall maintain Access Control Lists (ACL) for resource-level permissions | M1 | ACL stored in database, cached in Redis for performance |
| FR-019 | The system shall implement hierarchical role inheritance (Partner Admin > Agent > Customer) | D | Child roles inherit parent permissions, override capability available |
| FR-020 | The system shall provide permission audit trail for all sensitive operations | M3 | Immutable audit logs, queryable by role/user/action/timestamp |

---

### 4.3 Product Management & Catalog (FG-003)

| FR ID | Requirement Description | Priority | Acceptance Criteria |
|-------|------------------------|----------|---------------------|
| FR-021 | The system shall provide product catalog with categorization: Health, Life, Motor, Travel, Micro-insurance. Preference order: Health → Auto → Travel → Life | M1 | Products displayed by category, search and filter enabled |
| FR-022 | The system shall support product search by name, category, coverage type, and premium range | M1 | Search results <500ms, fuzzy matching for Bengali text |
| FR-023 | The system shall display product details: coverage, premium, tenure, exclusions, terms & conditions | M2 | All product information visible before purchase, PDF download available |
| FR-023-A | User will be able to purchase unit-wise plan. User can increase / decrease coverage amount | M2 | Coverage amount can be increased or decreased |
| FR-023-B | User should get multiple questions as risk assessment for all plans | M2 | Every plan should contain assessment questions |
| FR-024 | The system shall provide premium calculator with dynamic inputs (age, sum assured, tenure, riders) | M3 | Real-time calculation <2s, breakdown of premium components shown |
| FR-025 | The system shall support product comparison (side-by-side up to 2 products) | M3 | Comparison table with key features, coverage, and pricing |
| FR-026 | The system shall enable Business Admin to create, update, and deactivate products | M1 | Product CRUD operations, version history maintained |
| FR-027 | The system shall support product variants with configurable riders and add-ons. Coverage can be increased for add-ons existing with B2B product. | D | Base product + optional riders, dynamic pricing recalculation |
| FR-028 | The system shall cache product catalog in Redis with 5-minute TTL for performance | M3 | Cache hit rate >80%, automatic invalidation on product updates |
| FR-029 | The system shall support multi-language product descriptions (Bengali and English) | M3 | Language toggle in UI, content stored in i18n format |

---

### 4.4 Policy Lifecycle Management (FG-004)

| FR ID | Requirement Description | Priority | Acceptance Criteria |
|-------|------------------------|----------|---------------------|
| FR-030 | The system shall support end-to-end policy purchase flow: product selection → applicant details → nominee details → payment → policy issuance | M1 | Complete flow in <10 minutes, progress saved at each step |
| FR-031 | The system shall collect applicant information: full name, DOB, NID (optional), address, occupation, income, health declaration | M1 | All mandatory fields validated, conditional fields based on product type |
| FR-032 | The system shall support Single nominee/beneficiary | M1 | Only 1 nominee required |
| FR-032-A | Beneficiary income range should be optional | M1 | Beneficiary data should submit without income range |
| FR-033 | The system shall validate NID/Mobile Number uniqueness across policies to prevent duplicate insurance | M1 | Database constraint enforced, user notified of existing policies |
| FR-034 | The system shall generate unique policy number with format: `LBT-YYYY-XXXX-NNNNNN`. Example: InsuranceType-Insurer-NumericID(Product) | M1 | Sequential numbering, year-based prefix, collision prevention |
| FR-035 | The system shall issue digital policy document (PDF) with QR code for verification | M2 | PDF generated within 30s of payment confirmation, QR code scannable |
| FR-036 | The system shall send policy document via SMS link and email attachment | M2 | Delivery within 5 minutes, retry mechanism on failure |
| FR-037 | The system shall activate policy immediately upon payment confirmation for instant coverage (Non Life) | M2 | Policy status updated in real-time, customer notified |
| FR-038 | The system shall support policy cooling-off period (5 days from issuance) for full refund | M3 | Cancellation request processed within 24hrs, refund initiated |
| FR-039 | The system shall maintain policy status: Pending Payment, Active, Suspended, Cancelled, Lapsed, Expired | M1 | Status transitions logged with timestamp, notifications triggered |
| FR-040 | The system shall provide customer policy dashboard showing all active and past policies, renewal prompts, and premium payment history | M1 | Dashboard loads <3s, real-time status updates |
| FR-041 | Users will be able to see order history with options to view: Coverage details, Refer, Active Plans, Claimed Plans, and Expired plans. Max Referral limit 1. User will get option for purchase and download the certificate. | D | — |

---

### 4.5 Policy Management & Renewals (FG-005)

| FR ID | Requirement Description | Priority | Acceptance Criteria |
|-------|------------------------|----------|---------------------|
| FR-042 | The system shall implement 'Family Insurance Wallet' allowing users to group and manage policies for multiple family members under one account | D | Unified dashboard, single-click bulk payment, relationship management |
| FR-043 | The system shall send renewal reminders: 30 days, 7 days, 1 day before expiry via SMS, email, push notification | M2 | Notifications sent on schedule, delivery confirmation tracked |
| FR-044 | The system shall support manual policy renewal with one-click process reusing existing policy data | M2 | Renewal completed in <3 minutes, updated policy document issued |
| FR-045 | The system shall support automatic policy repurchase with stored payment method (opt-in by customer) | M3 | Customer consent recorded, auto-charge 7 days before expiry |
| FR-046 | The system shall allow customer to update policy details during renewal: current address, nominee information | M3 | Limited fields editable, verification required for major changes |
| FR-047 | The system shall implement grace period (30 days) for premium payment post-expiry with continued coverage | M2 | Policy status "Grace Period", coverage continues, daily reminders |
| FR-048 | The system shall auto-lapse policy after grace period if payment not received, with reinstatement option | M2 | Policy status "Lapsed", reinstatement within 90 days with penalty |
| FR-049 | The system shall provide policy document download (PDF) with version history for all renewals | M1 | All versions accessible, clearly marked with issue date |
| FR-050 | The system shall track policy lifecycle events: issuance, renewal, lapse, reinstatement, cancellation with audit trail | M1 | Immutable event log, queryable by date range and policy number |

#### 4.5.1 Policy Cancellation & Refund

| FR ID | Requirement Description | Priority | Acceptance Criteria |
|-------|------------------------|----------|---------------------|
| FR-051 | The system shall support policy cancellation workflow with cancellation request submission by customer/agent/admin | M1 | Request form with reason dropdown, attachment support |
| FR-052 | The system shall implement approval workflow for policy cancellation: Business Admin + Focal Person approval required for policies >30 days old | M1 | Approval routing, 48hr SLA |
| FR-053 | The system shall calculate pro-rata refund: (Premium Paid - Days Covered - Admin Fee - Cancellation Charge) with transparent breakdown | M1 | Refund calculator, configurable fees |
| FR-054 | The system shall process refund within 7 working days via MFS or bank transfer | M1 | Payment gateway integration, notifications |
| FR-055 | The system shall update policy status to CANCELLED and notify all stakeholders | M1 | Multi-channel notification, IDRA reporting |

#### 4.5.2 Policy Endorsement & Amendment

| FR ID | Requirement Description | Priority | Acceptance Criteria |
|-------|------------------------|----------|---------------------|
| FR-056 | The system shall support policy endorsement for: Address, Sum insured, Nominee, Contact changes | M1 | Amendment forms, validation |
| FR-057 | The system shall calculate additional premium for mid-term sum insured increases | D | Premium calculator, payment integration |
| FR-058 | The system shall calculate pro-rata refund for sum insured decreases | M2 | Credit to premium account |
| FR-059 | The system shall generate endorsement document with suffix (PLN-001/END-01) | M1 | PDF generation, version tracking |
| FR-060 | The system shall require approval for sum insured changes >10% | M1 | Approval workflow, threshold config |

---

### 4.6 Business Rules & Workflows (FG-006)

| FR ID | Requirement Description | Priority | Acceptance Criteria |
|-------|------------------------|----------|---------------------|
| FR-061 | The system shall implement premium calculation fallbacks: If insurer API fails, use cached rates (max 24hrs old); if unavailable, queue quote and notify customer within 2 hours | M1 | Fallback logic tested, cache validation, queue notification works |
| FR-062 | The system shall handle premium calculation edge cases: age-based loading, occupation risk factors, pre-existing conditions with clear messaging | M2 | All edge cases covered, messaging user-friendly, actuarial validation |
| FR-063 | The system shall implement duplicate policy detection: Block duplicate policy purchase for same product + same insured person within 30 days; allow cross-product purchases | M1 | Detection accurate, cross-product allowed, clear error message |
| FR-064 | The system shall enable policy merge workflow: Focal Person can merge duplicate accounts after verifying NID, transfer policies, consolidate claims history | M3 | Merge workflow tested, data integrity maintained, audit logged |
| FR-065 | The system shall define claim status state machine: Submitted → Under Review → Documents Requested → Approved/Rejected → Payment Initiated → Settled/Closed | M1 | State machine implemented, invalid transitions blocked, status tracking accurate |
| FR-066 | The system shall enforce claim status transition rules: Auto-move to "Documents Requested" if incomplete; require Business Admin+Focal Person approval for >BDT 50K | M1 | Transition rules enforced, approval routing correct, notifications sent |
| FR-067 | The system shall implement gamified renewal rewards program offering discounts or gift vouchers for early renewals | D | Points calculation engine, partner voucher integration, leaderboard |
| FR-068 | The system shall implement grace period logic: 30-day grace period post-expiry with coverage continued; auto-lapse if unpaid after grace period | M3 | Grace period enforced, coverage continued, auto-lapse works, customer notified |
| FR-069 | The system shall enable lapsed policy reinstatement: Allow reinstatement within 90 days of lapse with medical underwriting; require Focal Person approval | D | Reinstatement workflow, medical underwriting integrated, approval required |

---

### 4.7 Payment Processing (FG-007)

| FR ID | Requirement Description | Priority | Acceptance Criteria |
|-------|------------------------|----------|---------------------|
| FR-070 | The system shall support multiple payment methods: bKash, Nagad, Rocket, Bank Transfer, Credit/Debit Card (MFS / Banking only) | M1 | All MFS integrated, card via hosted payment page, manual verification |
| FR-071 | The system shall integrate bKash payment gateway with production credentials and sandbox for testing | M1 | Transaction success rate >99%, fallback to manual on failure |
| FR-072 | The system shall integrate Nagad and Rocket MFS with tokenization for recurring payments | M3 | Secure token storage, PCI-DSS Level SAQ-A compliance |
| FR-073 | The system shall support manual payment with proof upload (bank receipt, bKash screenshot) for verification | M1 | Image upload <5MB, admin verification within 24hrs |
| FR-074 | The system shall implement payment verification workflow: pending → verified → policy activated OR rejected → refund | M2 | Admin approval for manual payments, automated for MFS |
| FR-075 | The system shall generate payment receipt with transaction ID, amount, date, policy number | M2 | PDF receipt sent via SMS/email within 5 minutes |
| FR-076 | The system shall support partial payment and installment plans for high-premium policies (Monthly or Yearly applicable for M2) | M3 | Auto-reminders before due date, grace period 15 days |
| FR-077 | The system shall implement payment retry mechanism with exponential backoff for failed transactions | M2 | Max 3 retries, customer notified on each attempt |
| FR-078 | The system shall support refund processing for policy cancellation with configurable refund rules | M2 | Refund initiated within 7 days, credited to original payment method |
| FR-079 | The system shall integrate TigerBeetle for financial transaction recording with double-entry bookkeeping | M2 | All transactions recorded, real-time balance reconciliation |
| FR-080 | The system shall maintain payment audit trail with immutable logs for regulatory compliance | M1 | PostgreSQL + S3 storage, 20-year retention |

---

### 4.8 Claims Management (FG-008)

| FR ID | Requirement Description | Priority | Acceptance Criteria |
|-------|------------------------|----------|---------------------|
| FR-081 | The system shall provide fixed-step claim submission form: policy selection, incident details, claim reason, document upload (images, bills, reports). Claim tracker must be shown. | M1 | Form completion <5 minutes, draft saving at each step |
| FR-082 | The system shall validate claim eligibility: policy active, within coverage period, claim type covered, no duplicate submission | M1 | Validation in <3s, clear error messages on rejection |
| FR-083 | The system shall generate unique claim number with format: `CLM-YYYY-XXXX-NNNNNN` and digital hash for submission integrity | M1 | Collision-free numbering, SHA-256 hash for document integrity |
| FR-084 | The system shall automatically notify partner/insurer upon claim submission with shared status dashboard | M2 | Notification within 60s, dashboard accessible to all stakeholders |
| FR-085 | The system shall provide real-time claim status tracking. Customer should be notified through push notification and SMS. | M3 | Status updates visible in <5s, push notifications on status change |
| FR-086 | The system shall implement tiered approval workflow based on claim amount as per Approval Matrix | M3 | Auto-routing to correct approver, escalation on timeout |
| FR-087 | The system shall support document verification with image quality check, OCR extraction, and fraud detection | M3 | Image validation <10s, OCR accuracy >85%, duplicate detection |
| FR-088 | The system shall provide chat interface between customer, partner agent, and focal person for claim discussion | M3 | Real-time messaging, file attachment support, message history |
| FR-089 | The system shall support WebRTC video call for claim verification and inspection | D | HD video quality, screen sharing, call recording for audit |
| FR-090 | The system shall allow partner to add verification notes and approve/reject with reason | M2 | Notes timestamped, approval requires mandatory reason field |
| FR-091 | The system shall enforce joint approval by Business Admin and Focal Person for claims BDT 50K–2L (Only insurance company will provide Approval) | M3 | Both approvals required, timeout escalation after 5 days |
| FR-092 | The system shall automate payment process upon claim approval as per customer's selected payment channel | M3 | Payment initiated within 24hrs, confirmation sent to customer |
| FR-093 | The system shall support Zero Human Touch Claims (ZHTC) for auto-verification and payment of small claims (<BDT 10K) with partner pre-agreement | D | 95% automation rate, ML-based fraud check, instant settlement |
| FR-094 | The system shall implement fraud detection: frequent claims (>3 in 6 months), duplicate documents, rapid policy-to-claim (<48hrs) | M3 | Auto-flagging with risk score, manual review queue, customer warning system |
| FR-095 | The system shall auto-revoke customer access for confirmed fraud as per InsureTech policy | M3 | Account suspension after approval, appeal process available |
| FR-096 | The system shall maintain balance sheet on Customer, Partner, Agent, and InsureTech level for selected time periods | M3 | Daily, monthly, quarterly reconciliation, export to Excel/PDF |
| FR-097 | The system shall track Turn Around Time (TAT) per approval level and alert on SLA breach | M3 | Real-time TAT monitoring, email alerts on approaching deadline |
| FR-098 | The system shall provide claim history and analytics for risk assessment and premium adjustment | M3 | Claim frequency report, average claim amount, settlement ratio |

#### Claims Approval Matrix (Insurance Company Only)

| Claimed Amount | Approval Level | Approver(s) | Maximum TAT |
|----------------|----------------|-------------|-------------|
| BDT 0–10K | L1 Auto/Officer | System Auto-Approval OR Claims Officer | 24 Hours |
| BDT 10K–50K | L2 Manager | Claims Manager | 3 days |
| BDT 50K–2L | L3 Head | Business Admin + Focal Person (Joint) | 7 days |
| BDT 2L+ | Board | Board + Insurer Approval | 15 days |

#### 4.8.1 Claims Document Requirements & Processing

| FR ID | Requirement Description | Priority | Acceptance Criteria |
|-------|------------------------|----------|---------------------|
| FR-099 | The system shall enforce claims document requirements: PDF/JPG/PNG, max 5 MB per file, 25MB total per claim, 300 DPI minimum | M1 | Client-side validation, OCR quality check |
| FR-100 | The system shall calculate co-payment and deductibles: (Claim Amount - Deductible) × Co-payment % with annual deductible tracking | M1 | Product-level config, breakdown display |
| FR-101 | The system shall support claims reimbursement workflow with document review and bank/MFS transfer within 7–15 working days (handled by insurance company) | M1 | Document verification, payment processing, status notifications |

---

### 4.9 Partner & Agent Management (FG-009)

| FR ID | Requirement Description | Priority | Acceptance Criteria |
|-------|------------------------|----------|---------------------|
| FR-102 | The system shall support partner onboarding workflow: application submission, KYB verification, MOU upload, approval by Focal Person | M2 | Complete onboarding in <7 days, status tracking at each step |
| FR-103 | The system shall collect partner information: organization name, type, trade license, TIN, bank account, contact details | M2 | All mandatory fields validated, document verification required |
| FR-104 | The system shall implement KYB (Know Your Business) verification with trade license validation and credit check | M2 | Automated validation where possible, manual review for exceptions |
| FR-105 | The system shall provide dedicated partner portal with dashboard showing: leads, conversions, commissions, analytics | M2 | Dashboard loads <3s, real-time data updates, export functionality |
| FR-106 | The system shall calculate and track partner commissions based on configurable rates (acquisition, renewal, claims assistance) | M2 | Commission calculated on policy activation, monthly payout reports |
| FR-107 | The system shall support partner API integration for embedded insurance (e-commerce checkout, hospital admission) | M3 | RESTful API with sandbox, developer documentation, webhook support |
| FR-108 | The system shall enable partner to initiate policy purchase on behalf of customer with consent and authentication | M2 | Customer OTP verification required, policy linked to customer account |
| FR-109 | The system shall provide Focal Person portal for partner management: verification, approval, dispute resolution, performance monitoring | M1 | Full CRUD operations on partners, approval workflow, audit trail |
| FR-110 | The system shall support multi-level agent hierarchy under partners (Partner Admin > Regional Manager > Agent) | M3 | Hierarchical commission split, territory management, performance tracking |
| FR-111 | The system shall track partner performance metrics: policies sold, claim settlement ratio, customer satisfaction, fraud incidents | M2 | Weekly/monthly reports, performance scoring, alerts on anomalies |
| FR-112 | The system shall support partner suspension/termination with graceful policy transfer mechanism | M2 | Existing policies remain active, new sales blocked, customer notification |

#### 4.9.1 Stakeholder Hierarchy & Focal Person Role

```
System Admin (Cloud/IAM Root)
    │
    ├── Repository Admin (Code/Deploy)
    │
    ├── Database Admin (Data Management)
    │
    ├── Business Admin (Business Operations)
    │   │
    │   └── Focal Person (Partner Bridge) ★ KEY ROLE
    │       │
    │       └── Partner Admin (Tenant Root)
    │           │
    │           ├── Insurer Agents (Partner Staff)
    │           └── Insurer Underwriters (API Users)
    │
    ├── Dev (Development)
    │
    └── Support/Call Centre (Customer Assistance)
```

**Focal Person Role Requirements:**

| FR ID | Requirement Description | Priority | Acceptance Criteria |
|-------|------------------------|----------|---------------------|
| FR-113 | Focal Person shall have authority to verify and approve/reject partner applications within 3 business days | M1 | Decision recorded with reason, partner notified automatically |
| FR-114 | Focal Person shall monitor partner compliance and flag suspicious activities for investigation | M2 | Real-time dashboard with alerts, escalation to Business Admin |
| FR-115 | Focal Person shall resolve partner-customer disputes with documented decision trail | M2 | Dispute resolution within 7 days, audit log maintained |

---

### 4.10 Partner Portal & Business Intelligence (FG-010)

| FR ID | Requirement Description | Priority | Acceptance Criteria |
|-------|------------------------|----------|---------------------|
| FR-116 | The system shall provide hospital partners special dashboard to initiate insurance purchase on behalf of customers | M2 | Patient data prefill from hospital system, consent capture |
| FR-117 | The system shall support API for transferring customer records with authentication token and purchase ID | D | RESTful API with OAuth2, data mapping documentation |
| FR-118 | The system shall provide e-commerce partners embedded widget for insurance product display at checkout | M2 | JavaScript SDK, responsive design, cart integration |
| FR-119 | The system shall provide sandbox environment for 3rd party developers with test credentials and mock data | D | Isolated test environment, sample code, API documentation |
| FR-120 | The system shall provide partner analytics: leads generated, conversion rate, commission earned, customer feedback | M2 | Dashboard with filters, trend charts, export to Excel/PDF |
| FR-121 | The system shall provide partner API for retrieving analytics and commission statements programmatically | D | RESTful API, pagination support, webhook for new data |
| FR-122 | The system shall implement Business Intelligence tool (Metabase/Tableau/Power BI) for advanced analytics | F | Read replica connection, pre-built dashboards, scheduled reports |
| FR-123 | The system shall provide executive dashboard: daily sales, policy count, claims ratio, revenue, system health | M2 | Real-time data, drill-down capability, mobile-responsive |
| FR-124 | The system shall provide partner-specific branding capability for white-label insurance offerings | F | Custom logo, colors, domain mapping, isolated tenant data |
| FR-125 | The system shall enable partners to configure commission structures and incentive programs | D | Tiered commission, bonus rules, performance-based adjustments |
| FR-126 | The system shall log all API requests with payload, headers, timestamps | M2 | Structured logging, rotation, searchable |
| FR-127 | The system shall implement distributed tracing across microservices | D | Jaeger integration, trace ID propagation |

---

### 4.11 Customer Support & Helpdesk (FG-011)

| FR ID | Requirement Description | Priority | Acceptance Criteria |
|-------|------------------------|----------|---------------------|
| FR-128 | The system shall provide in-app FAQ section with searchable knowledge base (hotline: 8801766662772) | M1 | Search results <1s, categorized by topic, Bengali and English |
| FR-129 | The system shall support customer support call initiation from mobile app with call recording | M3 | Click-to-call integration, call routing to available agent |
| FR-130 | The system shall implement ticketing system for customer issues with unique ticket ID and status tracking | M2 | Ticket creation <30s, status updates via notification |
| FR-131 | The system shall provide support agent portal with ticket queue, customer history, and resolution templates | M2 | Agent dashboard loads <3s, SLA countdown visible |
| FR-132 | The system shall auto-record customer support calls and create ticket with call summary | M3 | Speech-to-text transcription, auto-tag issue category |
| FR-133 | The system shall track support metrics: average response time, resolution time, customer satisfaction score | M2 | Real-time dashboard, weekly reports to management |
| FR-134 | The system shall support escalation workflow: Tier 1 (Support) → Tier 2 (Technical) → Tier 3 (Engineering) | M2 | Auto-escalation after 24hrs unresolved, notification sent |
| FR-135 | The system shall provide customer feedback form after ticket resolution with 5-star rating | M2 | Feedback collected, low ratings flagged for review |

---

### 4.12 Notifications & Communication (FG-012)

| FR ID | Requirement Description | Priority | Acceptance Criteria |
|-------|------------------------|----------|---------------------|
| FR-136 | The system shall implement Kafka event-driven notification system with multiple channels: in-app push, SMS, email | M1 | Event published within 100ms, delivery to all channels coordinated |
| FR-137 | The system shall send notifications for: OTP, verification, purchase confirmation, claims updates, renewal reminders, payment confirmations | M1 | Template-based messages, personalized with customer data |
| FR-138 | The system shall support notification preferences with opt-in/opt-out for marketing and promotional messages | M2 | User preferences stored, GDPR-compliant consent management |
| FR-139 | The system shall implement customer mute mode with minimum text notification (avoiding push for low-end devices) | M2 | Device capability detection, graceful degradation |
| FR-140 | The system shall allow partners to create secondary marketing notifications filtered by: age, gender, location, policy type | M3/D | — |
| FR-141 | The system shall track notification delivery status: queued, sent, delivered, failed, bounced with retry mechanism | M2 | Real-time status tracking, max 3 retries with exponential backoff |
| FR-142 | The system shall support message templates with dynamic placeholders for personalization | M2 | Template engine with Bengali/English support, variable substitution |
| FR-143 | The system shall implement rate limiting for notifications to prevent spam (max 5 per hour per user) | M3 | Redis-based rate limiting, exception for critical alerts |
| FR-144 | The system shall provide notification history in customer dashboard with read/unread status | M3 | Last 90 days visible, older notifications archived |
| FR-145 | The system shall support rich push notifications with images, action buttons, and deep links | D | Platform-specific implementation (iOS/Android), click tracking |

---

### 4.13 IoT Integration & Usage-Based Insurance (FG-013)

| FR ID | Requirement Description | Priority | Acceptance Criteria |
|-------|------------------------|----------|---------------------|
| FR-146 | The system shall support IoT device integration for Usage-Based Insurance (UBI) via proprietary protocol | F | MQTT/CoAP protocol support, device authentication, encrypted communication |
| FR-147 | The system shall collect and process IoT data: location, speed, temperature, health vitals based on insurance type | D | Real-time data ingestion, time-series database storage |
| FR-148 | The system shall implement risk scoring based on IoT data patterns for dynamic premium adjustment | F | ML-based risk model, monthly recalculation, customer notification |
| FR-149 | The system shall provide customer dashboard showing IoT insights and risk score with improvement tips | F | Visualization with charts, gamification elements, personalized recommendations |
| FR-150 | The system shall support telematics integration for motor insurance with driving behavior analysis | D | Acceleration, braking, speed monitoring, trip history, safety score |
| FR-151 | The system shall integrate with wearable devices for health insurance with fitness tracking | D | Steps, heart rate, sleep quality monitoring, wellness rewards program |
| FR-152 | The system shall implement data privacy controls allowing customers to pause/resume IoT data collection | F | One-click toggle, data deletion option, privacy dashboard |
| FR-153 | The system shall integrate with IoT devices: GPS trackers (vehicles), health wearables (fitness bands), smart home sensors | M3 | MQTT/CoAP protocol support, device SDK documentation, API endpoints |
| FR-154 | The system shall support IoT device registration, provisioning, and lifecycle management with certificate-based authentication | M3 | X.509 certificates, device onboarding workflow, status tracking |
| FR-155 | The system shall process and store IoT telemetry data using MQTT broker with TimescaleDB | M3/D | Handle 10,000 devices, 1 msg/min/device average, 90 days hot / 2 years warm retention |
| FR-156 | The system shall generate real-time alerts based on IoT data thresholds | M3/D | Rule engine for threshold monitoring, push notifications, SMS alerts, configurable rules |
| FR-157 | The system shall support Usage-Based Insurance (UBI) pricing calculation based on IoT data | M3/D | Dynamic premium adjustment algorithm, monthly recalculation, transparent scoring dashboard |
| FR-158 | The system shall provide IoT device management portal for partners | M3/D | Real-time device status, data visualization charts, anomaly detection, bulk operations |
| FR-159 | The system shall support batch and real-time IoT data processing with configurable collection frequencies | M3/D | Stream processing (Kafka Streams), batch jobs, data quality checks, deduplication |
| FR-160 | The system shall maintain IoT device inventory with status tracking and metadata | M3/D | Device registry, heartbeat monitoring (5min timeout), auto-offline detection |

---

### 4.14 AI & Automation Features (FG-014)

| FR ID | Requirement Description | Priority | Acceptance Criteria |
|-------|------------------------|----------|---------------------|
| FR-161 | The system shall implement AI chatbot for customer assistance during product search, selection, purchase, and claims | F | Bengali NLP support, 80% query resolution, human handoff capability |
| FR-162 | The system shall implement LLM multi-agent network for intelligent document processing and validation | F | OCR integration, field extraction accuracy >90%, fraud detection |
| FR-163 | The system shall implement AI-powered fraud detection using pattern recognition and anomaly detection | D | ML model with continuous learning, risk scoring, false positive <10% |
| FR-164 | The system shall support predictive analytics for risk assessment and premium optimization | F | Historical data analysis, model retraining, A/B testing capability |
| FR-165 | The system shall implement voice-assisted workflow for Type 3 users (rural/low digital literacy) | F | Bengali speech recognition, step-by-step guidance, voice commands |
| FR-166 | The system shall provide AI-based document verification with face matching and NID validation | M3 | Liveness detection, face match confidence >95%, automated approval flow |

**AI Multi-Agent Architecture:**

```
┌─────────────────────────────────────────────────────────┐
│                  AI ENGINE (Python + gRPC)              │
├─────────────────────────────────────────────────────────┤
│  Agent 1: Document Processing                           │
│  - OCR & Text Extraction                               │
│  - NID/Document Validation                             │
│  - Medical Report Analysis                             │
├─────────────────────────────────────────────────────────┤
│  Agent 2: Customer Service                             │
│  - Bengali Language Processing                         │
│  - FAQ & Query Resolution                              │
│  - Escalation Decision Making                          │
├─────────────────────────────────────────────────────────┤
│  Agent 3: Risk Assessment                              │
│  - Fraud Pattern Detection                             │
│  - Behavioral Analysis                                 │
│  - Claim Risk Scoring                                  │
├─────────────────────────────────────────────────────────┤
│  Agent 4: Business Intelligence                        │
│  - Predictive Analytics                                │
│  - Market Trend Analysis                               │
│  - Customer Segmentation                               │
└─────────────────────────────────────────────────────────┘
```

---

### 4.15 Voice-Assisted Features (FG-015)

| FR ID | Requirement Description | Priority | Acceptance Criteria |
|-------|------------------------|----------|---------------------|
| FR-167 | The system shall support Bengali speech-to-text (STT) with 90%+ accuracy for standard dialects (Dhaka, Chittagong, Sylhet) | M2 | ASR model integration, <2s latency, multi-dialect support |
| FR-168 | The system shall provide voice-guided policy purchase workflow with step-by-step audio instructions in Bengali | M2 | Complete policy purchase via voice, TTS integration, progress tracking |
| FR-169 | The system shall support voice-based claims submission with automated transcription and field validation | M3 | Voice recording up to 5min, transcription accuracy >85%, auto-populate claim form |
| FR-170 | The system shall provide text-to-speech (TTS) for Bengali language with natural-sounding voice | M2 | Natural prosody, <1s response time, caching for common phrases, offline fallback |
| FR-171 | The system shall support voice navigation throughout mobile app for accessibility (elderly/visually impaired users) | D | Voice commands for all major functions, screen reader compatibility |
| FR-172 | The system shall provide voice command taxonomy: "buy policy", "file claim", "check status", "pay premium", "call agent" | M2 | Intent recognition with 85%+ accuracy, contextual understanding, error handling |
| FR-173 | The system shall support seamless fallback to human agent when voice recognition confidence is below 80% | M3 | Confidence scoring, automatic handoff with context transfer, queue management |
| FR-174 | The system shall log and analyze voice interactions for continuous improvement with user consent | D | Voice data collection opt-in, anonymization, model retraining pipeline |

---

### 4.16 Fraud Detection & Risk Controls (FG-016)

| FR ID | Requirement Description | Priority | Acceptance Criteria |
|-------|------------------------|----------|---------------------|
| FR-175 | The system shall flag claims submitted within 48hrs of policy purchase for manual review | M2 | Auto-flagging with notification to Claims Officer, review queue |
| FR-176 | The system shall detect same claim type >2 times in 12 months and flag for pattern analysis | M2 | Historical claim analysis, risk scoring, enhanced verification |
| FR-177 | The system shall flag claims where amount exactly matches policy limit (100% of coverage) | M2 | Suspicious pattern detection, additional document requirements |
| FR-178 | The system shall validate medical provider against approved network list and flag non-network claims | M2 | Provider database, real-time validation, approval workflow |
| FR-179 | The system shall implement device fingerprinting to detect multiple accounts from same device (>3 accounts) | M3 | Browser/mobile device ID tracking, IP analysis, account linking |
| FR-180 | The system shall provide fraud detection dashboard for Business Admin and Focal Person with drill-down capability | M2 | Real-time alerts, risk score visualization, action buttons |
| FR-181 | The system shall implement RACI for monitoring and incident escalation per defined roles | M1 | Responsibility matrix enforced, escalation triggers |

**Fraud Detection Rules:**

| Rule ID | Rule Description | Threshold | Action |
|---------|-----------------|-----------|--------|
| FR-182 | Rapid Policy-Claim: Policy purchase to claim submission | < 48 hours | Auto-flag + manual review |
| FR-183 | Frequent Claims: Same claim type repetition | >2 times in 12 months | Flag + pattern analysis |
| FR-184 | Amount Matching: Claim amount exactly matches coverage | 100% of coverage | Flag + enhanced verification |
| FR-185 | Network Violation: Medical provider not in approved list | Non-network provider | Flag + provider verification |
| FD-186 | Geographic Anomaly: Claim location vs registered address | >100 km distance | Flag + location verification |
| FD-187 | Device Fingerprinting: Multiple accounts from same device | >3 accounts | Flag + identity verification |
| FD-188 | Behavioral Pattern: Unusual activity patterns | ML-based scoring | Risk scoring + monitoring |

---

### 4.17 Admin & Reporting (FG-017)

| FR ID | Requirement Description | Priority | Acceptance Criteria |
|-------|------------------------|----------|---------------------|
| FR-189 | The system shall provide role-based admin dashboards for all admin roles | M1 | Dynamic content based on role, real-time data updates |
| FR-190 | The system shall enforce strict 2FA for all admin-level access with TOTP authentication | M1 | Google Authenticator/Authy compatible, backup codes provided |
| FR-191 | The system shall provide user management module: create, update, suspend, delete users with audit trail | M2 | Full CRUD operations, role assignment, activity logs |
| FR-192 | The system shall provide product management module: create, update, activate/deactivate insurance products | M1 | Version control, effective date management, pricing configuration |
| FR-193 | The system shall provide claims management dashboard with filtering: status, amount range, date, partner | M2 | Advanced search, bulk actions, export functionality |
| FR-194 | The system shall provide task management system with assignment to internal users and deadline tracking | D | Task creation, assignment, status updates, notification on overdue |
| FR-195 | The system shall generate standard reports: daily sales, claims ratio, partner performance, policy counts, revenue | M2 | — |
| FR-196 | The system shall provide custom report builder with drag-drop interface for business users | D | Visual query builder, chart generation, saved report templates |
| FR-197 | The system shall track KPIs aligned to business plan: policy acquisition rate, claim settlement ratio, customer retention | M3 | — |
| FR-198 | The system shall provide system health monitoring dashboard: server status, API response times, error rates | M2 | Integration with Prometheus/Grafana, alert configuration |

---

### 4.18 Analytics & Reporting (FG-018)

| FR ID | Requirement Description | Priority | Acceptance Criteria |
|-------|------------------------|----------|---------------------|
| FR-199 | The system shall track user behavior analytics: page views, feature usage, drop-off points, conversion funnel | D | Integration with analytics platform (Google Analytics/Mixpanel) |
| FR-200 | The system shall provide predictive analytics for customer churn, claim likelihood, policy renewal probability | F | — |
| FR-201 | The system shall generate customer segmentation reports: demographics, policy type, risk profile, lifetime value | D | Automated segmentation, export for marketing campaigns |
| FR-202 | The system shall provide geographic analytics: policy distribution by district, claims heatmap, agent performance by region | D | Map visualization, district-level drill-down, comparative analysis |
| FR-203 | The system shall provide geospatial risk visualization overlaying claims data on regional maps | D | Mapbox/Google Maps integration, district-level aggregation, color-coded risk zones |
| FR-204 | The system shall provide pre-built dashboards: Executive, Operations, Compliance with drill-down | D | Interactive charts, export capability, scheduled email delivery |
| FR-205 | The system shall track compliance metrics: AML flags, IDRA report status, audit logs access | D | Real-time compliance dashboard, alerts on violations |

---

### 4.19 Audit & Logging (FG-019)

| FR ID | Requirement Description | Priority | Acceptance Criteria |
|-------|------------------------|----------|---------------------|
| FR-206 | The system shall maintain immutable audit logs for critical actions: policy issue, claim approval, payment, dispute resolution | M1 | PostgreSQL with append-only tables, tamper detection |
| FR-207 | The system shall implement data retention policy with 20-year minimum for regulatory compliance | M2 | Tiered storage (hot/warm/cold), automated archival, retrieval SLA |
| FR-208 | The system shall track all logged-in user actions with IP address, device info, timestamp, action type | M3 | Comprehensive logging, queryable audit trail, GDPR compliance |
| FR-209 | The system shall allow partners to maintain additional logs as per MOU agreement with InsureTech | F | Partner-specific log tables, data isolation, access controls |
| FR-210 | The system shall provide regulatory portal for IDRA/BFIU to access requested data as per law | M2 | Secure portal, report generation, audit trail of data access |
| FR-211 | The system shall implement log aggregation and analysis with alerting on suspicious patterns | M2 | ELK stack/CloudWatch integration, anomaly detection, real-time alerts |

---

### 4.20 System Interface Architecture (FG-020)

| FR ID | Requirement Description | Priority | Acceptance Criteria |
|-------|------------------------|----------|---------------------|
| FR-212 | The system shall implement High-Performance Internal API for gateway-microservices communication | M1 | <100ms response time, circuit breaker pattern, retry logic |
| FR-213 | The system shall implement Client-Optimized API for gateway-customer device communication | M1 | <2s response time, query optimization, field-level authorization |
| FR-214 | The system shall implement Standard Integration API for 3rd party partners with comprehensive documentation | D | <200ms response time, standardized docs, sandbox environment |
| FR-215 | The system shall provide Public Discovery API for product search and listing with rate limiting | M1 | <1s response time, request limiting, caching enabled |
| FR-216 | The system shall expose only Cloudflare proxy and NGINX entry node to public, blocking direct microservice access | M1 | Firewall rules configured, internal IPs hidden, DDoS protection |
| FR-217 | The system shall implement Real-Time Connection capability for instant updates | D | Persistent connection management, automatic reconnection, heartbeat |
| FR-218 | The system shall use Efficient Binary Protocol for IoT data extraction and data binding | F | Custom binary formatting, data compression, low latency |
| FR-219 | The system shall consolidate, annotate and process data for AI agent training within regulatory limits | F | Data anonymization, consent management, audit trail |
| FR-220 | The system shall generate statistics and predictions based on big data for partner insights | F | ML pipeline, data lake architecture, API for insights delivery |
| FR-221 | The system shall implement Blockchain-based shared ledger for automated reinsurance settlements | D | Immutable ledger, transparency audit trail |
| FR-222 | The system shall implement AI-driven dynamic premium discounting based on real-time risk assessment | D | Risk model integration, real-time calculation, customer notification |
| FR-223 | The system shall integrate with SMS Gateway for OTP and notifications | M1 | Delivery rate >95%, delivery status tracking, cost optimization |
| FR-224 | The system shall integrate with Email Service for transactional and marketing emails | M1 | Template management, bounce handling, unsubscribe management |
| FR-225 | The system shall provide Webhook System for real-time event notifications to external systems | M2 | Event filtering, retry mechanism, authentication, payload signing |

**API Category Structure:**

| API Category | Protocol | Use Case | Security Layer | Performance Target |
|-------------|---------|---------|----------------|-------------------|
| Category 1 | Protocol Buffer + gRPC | Gateway ↔ Microservices | System Admin Middle Layer | < 100ms |
| Category 2 | GraphQL + JWT | Gateway ↔ Customer Device | JWT + OAuth v2 | < 2 seconds |
| Category 3 | RESTful + JSON (OpenAPI) | 3rd Party Integration | Server-side Auth | < 200ms |
| Public API | RESTful + JSON (OpenAPI) | Product Search/List | Public Access | < 1 second |

---

### 4.21 Integration (FG-021)

| FR ID | Requirement Description | Priority | Acceptance Criteria |
|-------|------------------------|----------|---------------------|
| FR-226 | The system shall provide API contract specification: All Category 3 APIs must provide OpenAPI 3.0 spec | M3 | OpenAPI spec complete, error codes documented, examples provided |
| FR-227 | The system shall define insurer API payloads: Premium Calculation API, Policy Issuance API | M1 | Payload formats defined, validation rules clear, sample payloads provided |
| FR-228 | The system shall define payment gateway payloads with HMAC-SHA256 signature validation | M1 | Payment payloads defined, signature validation implemented, security tested |
| FR-229 | The system shall implement retry logic with exponential backoff: 1s, 2s, 4s, 8s, 16s (max 5 retries) + circuit breaker | M1 | Retry logic tested, exponential backoff works, circuit breaker functional |
| FR-230 | The system shall implement idempotency: All payment and policy issuance APIs must accept Idempotency-Key header (UUID); Store keys for 24 hours | M1 | Idempotency enforced, key storage works, duplicate handling correct |
| FR-231 | The system shall implement callback security: Payment gateway webhooks must include HMAC-SHA256 signature | M2 | Signature validation works, invalid callbacks rejected, logging comprehensive |
| FR-232 | EHR Integration Option A (Preferred): LabAid FHIR API with Patient resource matching by NID/phone | S | FHIR API integrated, patient matching accurate, pre-auth workflow functional |
| FR-233 | EHR Integration Option B (Fallback): LabAid custom REST API with mutual TLS + API key | D | Custom API integrated, mTLS configured, API key management |
| FR-234 | The system shall handle EHR integration timeout: 5s connection, 15s read; queue for manual verification on timeout | D | Timeout handling works, manual queue functional, notifications sent |

---

### 4.22 Data Storage (FG-022)

| FR ID | Requirement Description | Priority | Acceptance Criteria |
|-------|------------------------|----------|---------------------|
| FR-235 | The system shall use PostgreSQL V17 for structured data with JSON support and full-text search | M1 | Primary database setup, performance optimization, localization |
| FR-236 | The system shall implement read replicas for reporting and analytics workloads | M3 | Read scaling, data consistency, performance monitoring |
| FR-237 | The system shall implement Graph Database (Neo4j/Amazon Neptune) for fraud relationship visualization | D | Graph schema defined, node relationship mapping, query performance <1s |
| FR-238 | The system shall use Redis for session management and high-frequency real-time data | M3 | Performance optimization, session management, cache strategies |
| FR-239 | The system shall implement data partitioning for policies and claims tables by month | M3 | Scalability, query performance, maintenance efficiency |
| FR-240 | The system shall use S3-compatible Object Storage for document files with encryption at rest | M1 | Secure document storage, lifecycle management, CDN integration |
| FR-241 | The system shall store product catalog and metadata in Document-Oriented NoSQL Database | M3 | Flexible schema, high availability, global distribution |
| FR-242 | Upload data policy: Client-side compression 5MB → 1-2MB (JPEG 80%, 1920x1080 max), Chunked upload 1MB chunks (tus.io), Presigned S3 URLs 30-min expiry | M1 | Check upload >5MB fails, <5MB passes |
| FR-243 | Backup: Daily full, 6-hour incremental, continuous transaction logs | M1 | Check new backup after 6 hours |
| FR-244 | The system shall store app native encrypted data in user device in SQLite | M2 | Check SQLite files |
| FR-245 | The system shall process tokenized data on Vector Database for AI embeddings | D | Similarity search latency check |
| FR-246 | The system shall implement Columnar Database (ClickHouse/Druid) for high-performance real-time analytics | D | OLAP query performance <500ms, data compression, scalability |

---

### 4.23 User Interface Requirements (FG-023)

#### 4.23.1 Mobile Application (Android/iOS)

**Customer Mobile App Requirements:**
- Platform Support: Android 8.0+ (API 26), iOS 13.0+
- Language Support: Bengali (primary), English (secondary)
- Offline Capability: Policy documents, basic information viewable offline
- Accessibility: WCAG 2.1 AA compliance for visually impaired users
- Performance: App startup < 3 seconds, screen transitions < 1 second

| FR ID | Requirement Description | Priority | Acceptance Criteria |
|-------|------------------------|----------|---------------------|
| FR-247 | The system shall maintain consistent UI across Android and iOS using React Native | M1 | Shared codebase >90% |
| FR-248 | The system shall provide smart data widgets for mobile users | D | Customizable dashboard |
| FR-249 | The system shall support desktop-first responsive design for portals | M1 | 1024px minimum width |
| FR-250 | The system shall request minimum device permissions | M1 | Camera, SMS read only |
| FR-251 | The system shall support Bengali and English with toggle | M1 | i18n framework implemented |

**Key Features:**
- User registration and KYC verification with document upload
- Product browsing and comparison
- Policy purchase and premium payment
- Claims submission with photo/video upload
- Policy document management and sharing
- Push notifications and in-app messaging
- Voice-assisted navigation for elderly users

#### 4.23.2 Web Portals

**Customer Web Portal:**
- Responsive design (desktop, tablet, mobile)
- Single-page application (SPA) architecture using React
- Multi-language support with language switcher
- Dashboard with policy overview, premium due dates, claims status
- Document management with secure download links
- Payment history and receipt downloads

**Partner Admin Portal:**
- Agent management and performance monitoring
- Commission calculation and payment tracking
- Sales analytics and reporting dashboards
- Product configuration and pricing management
- Customer support tools and escalation workflows
- Bulk operations for agent onboarding

**System Admin Portal:**
- User and role management
- System configuration and feature toggles
- Monitoring dashboards and system health metrics
- Regulatory reporting and compliance tracking
- Audit log viewing and analysis
- Business intelligence and analytics tools

#### 4.23.3 UI/UX Guidelines

**Design Principles:**
- Bangladesh-centric design with cultural sensitivity
- Mobile-first responsive design approach
- Accessibility compliance (WCAG 2.1 AA)
- Progressive Web App (PWA) capabilities
- Consistent color scheme and branding
- Bengali typography and text rendering optimization

**Interaction Patterns:**
- Intuitive navigation with minimal cognitive load
- Voice input support for Bengali language
- Gesture-based navigation on mobile devices
- Contextual help and guided tutorials
- Error prevention and graceful error handling
- Confirmation dialogs for critical actions

---

## 5. Non-Functional Requirements & Technical Constraints

### 5.1 Technology Constraints

| NFR ID | Constraint Area | Requirement | Measurement | Priority |
|--------|----------------|-------------|-------------|----------|
| NFR-252 | Database Technology | PostgreSQL V17 with JSONB support | ACID compliance tests | M1 |
| NFR-253 | Caching & Session | Redis for distributed caching and session management | Cache hit ratio monitoring | M1 |
| NFR-254 | API Protocol | gRPC with Protocol Buffers (Category 1) | Inter-service latency metrics | M1 |
| NFR-255 | Client API | REST (OpenAPI 3.0) with JWT authentication | Schema validation, Token checks | M1 |
| NFR-256 | Public Integration | RESTful APIs with OpenAPI 3.0 specifications (Category 3) | Swagger validator pass | D |
| NFR-257 | Search Engine | PostgreSQL Full-Text Search or dedicated engine | Query performance <200ms | M1 |
| NFR-258 | Object Storage | S3-compatible storage (AWS/DigitalOcean) | Upload/Download latency | M1 |
| NFR-259 | Message Broker | Apache Kafka | Throughput monitoring | M1 |
| NFR-260 | Time-Series Data | TimescaleDB for IoT telemetry | Ingestion rate monitoring | M2 |
| NFR-261 | Vector Database | Pgvector or Pinecone for AI embeddings | Similarity search latency | D |
| NFR-262 | Graph Database | Neo4j or Amazon Neptune for fraud visualization | Graph traversal depth/speed | D |
| NFR-263 | Columnar Database | ClickHouse or Druid for analytics | Analytical query speed | D |
| NFR-264 | Financial Ledger | TigerBeetle for double-entry bookkeeping | Ledger reconciliation check | M3 |
| NFR-265 | Mobile Framework | React Native | Code reuse >80% | M1 |
| NFR-266 | CDN & Security | Cloudflare proxy | WAF block rate | M1 |

### 5.2 Performance Requirements

| NFR ID | Requirement | Target | Measurement | Priority |
|--------|-------------|--------|-------------|----------|
| NFR-267 | API response time for policy operations | < 500ms (95th percentile) | Application performance monitoring | M1 |
| NFR-268 | Database query response time | < 100ms (average) | Database monitoring tools | M1 |
| NFR-269 | Mobile app startup time | < 3 seconds | App performance analytics | M1 |
| NFR-270 | Web portal page load time | < 2 seconds | Browser performance tools | M1 |
| NFR-271 | Payment processing time | < 10 seconds end-to-end | Payment gateway analytics | M1 |
| NFR-272 | Claim processing automation | 80% straight-through processing | Business process monitoring | M2 |
| NFR-273 | Report generation time | < 30 seconds for standard reports | Reporting system metrics | M2 |
| NFR-274 | Search functionality response | < 200ms for basic searches | Search performance monitoring | M2 |

### 5.3 Scalability Requirements

| NFR ID | Requirement | Target | Measurement | Priority |
|--------|-------------|--------|-------------|----------|
| NFR-275 | Concurrent user support | 10,000 active users | Load testing and monitoring | M1 |
| NFR-276 | Transaction throughput | 1,000 TPS (policies + claims) | Performance testing | M2 |
| NFR-277 | Database scalability | 100 million policy records | Database performance testing | M2 |
| NFR-278 | Auto-scaling capability | Scale out/in based on load | Infrastructure monitoring | M2 |
| NFR-279 | Peak load handling | 5x normal load during campaigns | Stress testing | M3 |
| NFR-280 | Storage scalability | 10TB+ document storage | Cloud storage metrics | M3 |

### 5.4 Availability & Reliability

| NFR ID | Requirement | Target | Measurement | Priority |
|--------|-------------|--------|-------------|----------|
| NFR-281 | System availability | 99.5% uptime (M1), 99.9% (M2) | Infrastructure monitoring | M1 |
| NFR-282 | Recovery Time Objective (RTO) | 4 hours maximum | Disaster recovery testing | M1 |
| NFR-283 | Recovery Point Objective (RPO) | 1 hour maximum data loss | Backup and recovery testing | M1 |
| NFR-284 | Mean Time To Recovery (MTTR) | < 2 hours | Incident response metrics | M2 |
| NFR-285 | Service degradation handling | Graceful degradation during outages | Chaos engineering testing | M2 |
| NFR-286 | Data backup frequency | Real-time replication + daily backups | Backup monitoring | M1 |

### 5.6 Usability Requirements

| NFR ID | Requirement | Target | Measurement | Priority |
|--------|-------------|--------|-------------|----------|
| NFR-287 | User satisfaction score | 4.5+ stars on app stores | User feedback and ratings | M2 |
| NFR-288 | Task completion rate | 95% for critical user journeys | User experience analytics | M1 |
| NFR-289 | Learning curve | New users complete first task < 5 minutes | User onboarding metrics | M2 |
| NFR-290 | Error recovery | Clear error messages with action guidance | Error tracking and analysis | M1 |
| NFR-291 | Accessibility compliance | WCAG 2.1 AA compliance | Accessibility testing tools | M2 |
| NFR-292 | Multi-language support | Bengali and English localization | Localization testing | M1 |

### 5.7 Maintainability & Operability

| NFR ID | Requirement | Target | Measurement | Priority |
|--------|-------------|--------|-------------|----------|
| NFR-293 | Code coverage | 80% unit test coverage | Automated testing reports | M2 |
| NFR-294 | Deployment frequency | Daily deployments capability | CI/CD pipeline metrics | M2 |
| NFR-295 | Mean Time To Deploy | < 30 minutes for hotfixes | Deployment automation metrics | M2 |
| NFR-296 | Monitoring coverage | 100% critical path monitoring | Observability platform | M1 |
| NFR-297 | Log aggregation | Centralized logging for all services | Logging platform metrics | M1 |
| NFR-298 | Documentation currency | API documentation auto-generated | Documentation automation | M2 |

---

## 6. Data Model & Persistence

### 6.1 Proto Schema Organization

```
proto/insuretech/
├── authn/                          # Authentication Domain
│   ├── entity/v1/
│   │   ├── user.proto
│   │   └── session.proto
│   ├── events/v1/
│   │   └── auth_events.proto
│   └── services/v1/
│       └── auth_service.proto
├── authz/                          # Authorization Domain
├── policy/                         # Policy Management Domain
│   ├── entity/v1/
│   │   └── policy.proto
│   ├── events/v1/
│   └── services/v1/
├── claims/                         # Claims Processing Domain
│   ├── entity/v1/
│   │   └── claim.proto
│   ├── events/v1/
│   └── services/v1/
├── payment/                        # Payment Processing Domain
├── partner/                        # Partner Management Domain
├── products/                       # Product Catalog Domain
├── notification/                   # Notification Domain
├── ai/                             # AI Engine Domain
├── analytics/                      # Analytics & BI Domain
└── iot/                            # IoT Integration Domain
```

### 6.2 Data Architecture Overview

**Data Storage Strategy:**
- **PostgreSQL 17+:** Primary transactional data (policies, claims, users, KYC)
- **TimescaleDB:** Time-series data (IoT telemetry, audit logs, analytics)
- **TigerBeetle:** Financial transactions with double-entry bookkeeping
- **DynamoDB:** Product catalog, configuration, session data
- **Redis 7.0+:** Caching, session management, real-time data
- **AWS S3:** Document storage (policy certificates, claims documents, images)
- **Apache Kafka:** Event streaming and audit logs
- **Pgvector:** Vector embeddings for AI/ML operations

### 6.3 Domain Models

**Authentication Domain:** User, Session, OTP

**Policy Domain:** Policy, Applicant, Nominee, Rider

**Claims Domain:** Claim, ClaimDocument, ClaimApproval, FraudCheck

**Payment Domain:** Payment, Transaction, Refund

**Partner Domain:** Partner, Agent, Commission

### 6.4 Proto-First Data Model Strategy

The platform adopts a Proto-First approach where Protocol Buffer definitions serve as the canonical data model across all layers:

```
Proto Definitions (Source of Truth)
    ├── Code Generation → Go/C#/Python/Node.js structs
    ├── Database Schema → PostgreSQL/TimescaleDB tables
    ├── API Contracts → gRPC/REST endpoints
    ├── Event Schemas → Kafka message formats
    └── Documentation → Auto-generated API docs
```

**Key Benefits:**
- ✅ Type Safety: Compile-time validation across all services
- ✅ Consistency: Same data structure in app, database, and APIs
- ✅ Versioning: Built-in support for backward/forward compatibility
- ✅ Multi-Language: Generate code for Go, C#, Python, Node.js from single source
- ✅ Performance: Efficient binary serialization
- ✅ Documentation: Self-documenting with comments in proto files

### 6.7 Data Retention & Archival

| Data Type | Hot Storage | Warm Storage | Cold Storage | Retention |
|-----------|-------------|--------------|--------------|-----------|
| Active Policies | PostgreSQL | — | — | Policy lifetime |
| Expired Policies | PostgreSQL (1 year) | S3 (5 years) | Glacier (20 years) | 20 years |
| Claims Data | PostgreSQL | S3 after settlement | Glacier (20 years) | 20 years |
| Audit Logs | TimescaleDB (90 days) | S3 (1 year) | Glacier (7 years) | 7 years |
| IoT Telemetry | TimescaleDB (90 days) | S3 (1 year) | Deleted | 1 year |
| User Sessions | Redis (7 days) | — | — | 7 days |

---

## 7. Security & Compliance Requirements

The LabAid InsureTech Platform implements a Zero Trust Security Model with defense-in-depth strategies.

**Core Security Principles:**
- **Never Trust, Always Verify:** All users and devices authenticated and authorized
- **Least Privilege Access:** Minimum required permissions for each role
- **Assume Breach:** Monitor and respond as if compromise has occurred
- **Encrypt Everything:** Data protection at all layers and states
- **Continuous Monitoring:** Real-time threat detection and response

### 7.1 Security Infrastructure & Key Management

| ID | Requirement Description | Priority |
|----|------------------------|----------|
| SEC-001 | Separate secret vault — AWS KMS/Azure Key Vault/HashiCorp, 90-day key rotation | M1 |
| SEC-002 | Data Masking: NID (last 3 digits), phone (mask middle), email (mask username) | M2 |
| SEC-003 | PCI-DSS compliance for card flows — Hosted payment page (redirect model), SAQ-A level, TLS 1.3, tokenization for recurring payments | M2 |
| SEC-004 | AML/CFT detection hooks with 20+ automated transaction monitoring rules | D |
| SEC-005 | IDRA reporting capabilities: Monthly (Form IC-1, IC-2), Quarterly (IC-3, IC-4), Annual FCR, 20-year archive | D |
| SEC-006 | Regular penetration testing — Pre-launch + annually | D |
| SEC-007 | Regular security audits from various security auditors and regulatory bodies | D |
| SEC-008 | DAST: OWASP ZAP/Burp Suite (weekly on staging) | D |
| SEC-009 | SAST: SonarQube/Checkmarx (every commit, block critical vulnerabilities) | D |
| SEC-010 | Virus scanning: ClamAV on uploaded files | M |
| SEC-011 | API rate limiting per user/IP: 1000 requests/hour authenticated, 100 requests/hour anonymous | M2 |
| SEC-012 | Separate encryption keys for different data types with hierarchical key management | M2 |
| SEC-013 | Real-time security incident response with automated threat isolation | M2 |
| SEC-014 | Continuous vulnerability assessment with automated patching for critical vulnerabilities | D |
| SEC-015 | Zero-trust network architecture with microsegmentation | D |

### 7.2 Enhanced IDRA Compliance

| ID | Requirement Description | Priority |
|----|------------------------|----------|
| SEC-016 | IDRA Monthly Reports: Form IC-1 (Premium Collection) by 10th of each month | M2 |
| SEC-017 | IDRA Monthly Reports: Form IC-2 (Claims Intimation) by 10th of each month | M2 |
| SEC-018 | IDRA Quarterly Reports: Form IC-3 (Claims Settlement) within 15 days of quarter-end | M2 |
| SEC-019 | IDRA Quarterly Reports: Form IC-4 (Financial Performance) within 20 days of quarter-end | M3 |
| SEC-020 | IDRA Annual FCR: Financial Condition Report within 90 days of year-end | M3 |
| SEC-021 | IDRA Event-Based Reporting: Significant incidents within 48 hours via IDRA portal | M3 |

### 7.3 Enhanced AML/CFT Compliance

| ID | Requirement Description | Priority |
|----|------------------------|----------|
| SEC-022 | AML/CFT Triggers: >3 policies in 7 days; Premium >BDT 5L without income proof; Nominee mismatch; Third-party payment; >2 cancellations in 30 days; Geographic anomaly; >3 failed KYC; PEP match | M3 |
| SEC-023 | SAR Workflow: Auto-flag → 24hr review → Escalate → Prepare SAR → Submit to BFIU within 3 business days → Enhanced monitoring (no customer notification — tipping off prohibited) | M3 |
| SEC-024 | Data Deletion Exceptions: Active policyholders, ongoing claims, SAR investigation, regulatory hold | M3 |
| SEC-025 | Right to Erasure Workflow: 30-day processing, anonymize PII while retaining transaction records, generate deletion certificate | D |

### 7.4 Data Protection & Encryption Standards

| Data Classification | Encryption Standard | Key Management | Access Control |
|--------------------|--------------------|--------------|----|
| PII | AES-256 | AWS KMS with 90-day rotation | Role-based with audit logging |
| Financial Transaction Data | AES-256 + Additional Hashing | TigerBeetle built-in encryption | Restricted access with MFA |
| KYC Documents | AES-256 with client-side encryption | End-to-end encryption | Compliance officer access only |
| Medical Records | AES-256 with additional anonymization | Healthcare-specific key management | Medical staff + consent-based |
| Audit Logs | AES-256 with immutable storage | Centralized key management | Read-only access for auditors |

### 7.6 IDRA Compliance Requirements

| IDRA ID | Requirement Description | Frequency | Priority |
|---------|------------------------|-----------|----------|
| IDRA-001 | Digital insurance product approval and registration | One-time + updates | M3 |
| IDRA-002 | Customer data protection and privacy compliance | Quarterly review | M3 |
| IDRA-003 | Policy issuance and documentation standards | Real-time compliance | M3 |
| IDRA-004 | Claims processing and settlement reporting | Monthly | M3 |
| IDRA-005 | Financial solvency and capital adequacy reporting | Quarterly | M3 |
| IDRA-006 | Agent licensing and training compliance | Ongoing | M3 |
| IDRA-007 | Marketing and sales practice compliance | Quarterly | M3 |
| IDRA-008 | Actuarial and risk management reporting | Annual | D |
| IDRA-009 | Audit trail and record keeping requirements | Ongoing | M3 |
| IDRA-010 | Regulatory change management and updates | As required | M3 |

### 7.7 BFIU Anti-Money Laundering (AML) Compliance

| BFIU ID | Requirement Description | Threshold | Priority |
|---------|------------------------|-----------|----------|
| BFIU-001 | Customer due diligence (CDD) for all policyholders | All customers | M3 |
| BFIU-002 | Enhanced due diligence (EDD) for high-value policies | >50,000 BDT sum assured | M3 |
| BFIU-003 | Suspicious transaction monitoring and reporting | Real-time analysis | M3 |
| BFIU-004 | Cash transaction reporting | >10,000 BDT | M3 |
| BFIU-005 | Wire transfer monitoring | >100,000 BDT | M3 |
| BFIU-006 | Politically exposed person (PEP) screening | All customers | M3 |
| BFIU-007 | Sanctions list screening | All parties | M3 |
| BFIU-008 | Record retention for AML purposes | 5 years minimum | M3 |
| BFIU-009 | AML training for employees and agents | Annual certification | M2 |
| BFIU-010 | AML audit and compliance reporting | Quarterly | M3 |

#### 7.7.2 Automated AML Monitoring Rules

| Rule ID | Rule Description | Threshold | Action |
|---------|-----------------|-----------|--------|
| TM-001 | Structuring: Multiple transactions just below reporting threshold | 3+ transactions of 9K–10K BDT in 7 days | Flag for review |
| TM-002 | Rapid Movement: Quick policy purchase and claim | Claim within 7 days of purchase | Flag + manual review |
| TM-003 | Geographic Anomaly | >100 km distance | Flag + location verification |
| TM-004 | Frequency Anomaly: Frequent claims | >3 claims in 6 months | Flag + pattern analysis |
| TM-005 | Amount Anomaly: Claim amount near coverage limit | >90% of coverage | Flag + document verification |
| TM-006 | Device Anomaly: Multiple accounts from same device | >3 accounts | Flag + fraud investigation |
| TM-007 | Payment Method Switch | >2 changes in 30 days | Flag + verification |
| TM-008 | Rapid Purchases: Multiple policies in short timeframe | >3 policies in 7 days | Flag + EDD |
| TM-009 | High-Value Premiums | >BDT 5 lakh | Enhanced due diligence |
| TM-010 | Frequent Cancellations | >2 cancellations in 3 months | Flag + investigation |
| TM-011 | Mismatched Nominees | Non-relative nominee | Flag + verification |
| TM-012 | Payment Source Anomaly | >2 different payers | Flag + source verification |
| TM-013 | Geographic Risk: High-risk location | Blacklisted areas | Enhanced monitoring |
| TM-014 | Age Anomaly | Outside typical range | Flag + verification |
| TM-015 | Occupation Risk: High-risk categories | PEP, cash-intensive business | Enhanced due diligence |
| TM-016 | Document Inconsistency | OCR verification failure | Flag + manual review |
| TM-017 | Refund Requests: Frequent refunds | >2 refunds in 6 months | Flag + pattern analysis |
| TM-018 | Beneficiary Changes | >2 changes in 12 months | Flag + verification |
| TM-019 | Third-Party Payments | Different payer than insured | Flag + source verification |
| TM-020 | Dormant Activation | No activity >6 months then sudden purchase | Flag + identity verification |
| TM-021–TM-030 | Additional rules including Time Anomaly (11PM–6AM), Velocity Check (>10/day), Round Amount patterns | Various | Flag + review |

---

## 8. Integration Requirements

### 8.1 External System Integrations

| Integration | Type | Protocol | Purpose | Priority |
|-------------|------|----------|---------|----------|
| bKash API | Payment Gateway | REST/Webhook | Premium payments, claim settlements | M1 |
| Nagad API | Payment Gateway | REST/Webhook | Premium payments, claim settlements | M2 |
| Bangladesh Bank | Payment Gateway | ISO 8583 | Bank transfers, regulatory reporting | M2 |
| NID Verification API | Government Service | REST/SOAP | Identity verification | M1 |
| Mobile Number Verification | Government Service | REST | Phone number validation | M1 |
| SMS Gateway | Communication | REST | Notifications and OTP delivery | M2 |
| Email Service | Communication | SMTP/API | Email notifications and documents | M1 |
| Hospital EHR Systems | Healthcare | HL7 FHIR | Medical record integration | D |
| Weather API | Risk Assessment | REST | Environmental risk monitoring | M3 |
| WhatsApp Business | Communication | API | Customer service and notifications | M3 |

**bKash Integration Details:**
- Authentication: OAuth 2.0 with client credentials flow
- Rate Limits: 100 requests/minute per API key
- Transaction Timeout: 5s connection, 15s read
- Webhook: Payment confirmation callback to `/webhooks/bkash`
- Error Handling: Retry with exponential backoff (3 attempts: 1s, 3s, 9s)
- Fallback: Queue payment for manual processing if API down >5 minutes

**NID Verification API:**
- Provider: Bangladesh Election Commission API / PORICHOY (TBD)
- Rate Limits: 1,000 verifications/day (Basic), 10,000/day (Premium)
- Response Time: <3 seconds average, 10 seconds timeout
- Data Returned: Name (EN/BN), Father/Mother name, DOB, Address, Photo (Base64)
- Cost: 5–10 BDT per verification
- Compliance: Store verification logs for 20 years

**Hospital EHR (HL7 FHIR) Integration:**
- Standard: HL7 FHIR R4
- FHIR Resources: Patient, Encounter, Condition, Procedure, MedicationRequest, DiagnosticReport
- Authentication: OAuth 2.0 + JWT tokens, refresh every 30 minutes
- Phase M1: LabAid Hospitals (5 locations)
- Phase M2+: Expand to 20+ partner hospitals

### 8.2 Internal Service Communications

All microservices communicate via gRPC using Protocol Buffers:

```protobuf
// Insurance Engine Service
service InsuranceEngineService {
  rpc IssuePolicy(IssuePolicyRequest) returns (IssuePolicyResponse);
  rpc CalculatePremium(CalculatePremiumRequest) returns (CalculatePremiumResponse);
  rpc ProcessRenewal(ProcessRenewalRequest) returns (ProcessRenewalResponse);
  rpc SubmitClaim(SubmitClaimRequest) returns (SubmitClaimResponse);
}
```

### 8.3 Event-Driven Architecture

**Kafka Topics:**
- **Policy Events:** PolicyIssued, PolicyRenewed, PolicyCancelled, PremiumPaid
- **Claim Events:** ClaimSubmitted, ClaimApproved, ClaimSettled, ClaimRejected
- **Payment Events:** PaymentProcessed, PaymentFailed, RefundIssued
- **User Events:** UserRegistered, KYCCompleted, ProfileUpdated

---

## 9. Performance & Monitoring

### 9.1 Performance Benchmarks

| Metric | Baseline Target | Peak Load Target | Measurement Method |
|--------|----------------|-----------------|-------------------|
| Category 1 API (gRPC) | < 100ms | < 150ms | APM tools (New Relic/Datadog) |
| Category 2 API (GraphQL) | < 2 seconds | < 3 seconds | GraphQL monitoring |
| Category 3 API (REST) | < 200ms | < 300ms | API gateway monitoring |
| Public API | < 1 second | < 1.5 seconds | Public endpoint monitoring |
| Mobile App Startup | < 5 seconds | < 7 seconds | Device testing |
| PostgreSQL Query | < 100ms for 95% | < 150ms for 95% | Database monitoring |
| TigerBeetle Transaction | < 10ms | < 20ms | Financial system monitoring |

### 9.2 Capacity Planning

| Component | Current Capacity | 12-Month Target | 24-Month Target | Scaling Strategy |
|-----------|-----------------|-----------------|-----------------|-----------------|
| Concurrent Users | 1,000 | 5,000 | 10,000 | Auto-scaling with CloudWatch |
| API Requests/Second | 100 | 1,000 | 5,000 | gRPC microservices scaling |
| Database Connections | 100 | 500 | 2,000 | PgBouncer connection pooling |
| TigerBeetle TPS | 1,000 | 10,000 | 50,000 | TigerBeetle cluster scaling |
| Storage (TB) | 1 | 10 | 50 | Auto-scaling object storage |
| Policy Documents | 10,000 | 500,000 | 2,000,000 | Distributed storage with archival |

### 9.4 Alerting & Incident Response

| Metric | Threshold | Alert Level | Action |
|--------|-----------|-------------|--------|
| API Response Time | > 1 second (95th percentile) | Warning | Auto-scale services |
| Database Connections | > 80% pool utilization | Warning | Scale database |
| Memory Usage | > 85% per service | Critical | Restart service |
| Disk Space | > 90% utilization | Critical | Add storage capacity |
| Authentication Failures | > 100 failed attempts/minute | Security | Block suspicious IPs |
| Payment Failures | > 5% failure rate | Critical | Alert payment team |
| System Downtime | > 5 minutes | Critical | Activate incident response |

---

## 10. Support & Maintenance

### 10.1 Support Model

| Support Level | Scope | Availability | Response SLA |
|--------------|-------|-------------|-------------|
| L1 - Customer Support | Basic inquiries, password resets | 24/7 | < 5 minutes |
| L2 - Technical Support | Application issues, payment problems | Business hours | < 30 minutes |
| L3 - Engineering Support | System bugs, performance issues | Business hours | < 2 hours |
| L4 - Critical Issues | System outages, security incidents | 24/7 | < 15 minutes |

**Support Channels:**
- Mobile App: In-app chat and support tickets
- Web Portal: Self-service help center and live chat
- Phone: Dedicated support hotline (Bengali/English)
- WhatsApp: Business account for basic inquiries
- Email: Support email with ticket tracking

### 10.2 Maintenance Windows

- **Daily:** Database optimization and log rotation (2:00 AM – 3:00 AM BST)
- **Weekly:** Security updates and patches (Sunday 1:00 AM – 3:00 AM BST)
- **Monthly:** Major updates and feature releases (First Saturday 10:00 PM – 2:00 AM BST)
- **Quarterly:** Infrastructure upgrades and capacity planning

---

## 11. Acceptance Criteria & Test Requirements

### 11.2 Test Types & Responsibilities

| Test Type | Coverage Target | Responsibility |
|-----------|----------------|----------------|
| Unit Testing | 80% code coverage | Development teams |
| Integration Testing | All service interfaces | QA team |
| API Testing | 100% endpoint coverage | Automation team |
| UI Testing | Critical user paths | QA team |
| Performance Testing | Load and stress scenarios | DevOps team |
| Security Testing | OWASP compliance | Security team |
| Accessibility Testing | WCAG 2.1 AA compliance | UX team |
| Compliance Testing | IDRA/BFIU requirements | Compliance team |

### 11.3 Critical Business Workflow Validation

| Workflow | Acceptance Criteria | Success Metrics |
|----------|--------------------|----|
| User Registration | Phone-based registration with OTP validation | >95% completion rate |
| KYC Verification | Document upload and verification within 5 minutes | >90% automated approval |
| Policy Purchase | End-to-end from product selection to policy issuance | >99% transaction success |
| Payment Processing | Multiple payment methods with real-time confirmation | >99.5% payment success |
| Claim Submission | Claim initiation with document upload and status tracking | <3 minutes submission time |
| Policy Renewal | Automated and manual renewal workflows | >95% renewal completion |

### 11.5 FR → Test Case Mapping

| FR-ID | Test Case ID | Test Scenario | Expected Result | Test Type |
|-------|-------------|--------------|----------------|-----------|
| FR-001 | TC-001 | Valid Bangladesh phone registration | OTP sent within 60s | Integration |
| FR-004 | TC-002 | Duplicate NID registration attempt | Error message displayed | Functional |
| FR-033 | TC-003 | End-to-end purchase with bKash | Policy issued within 30s | E2E |
| FR-051 | TC-004 | Joint approval (BizAdmin+Focal) | Claim approved only after both | Workflow |
| FR-129 | TC-005 | Insurer API failure during quote | Cached rate used + customer notified | Resilience |

---

## 12. Traceability Matrix & Change Control

### 12.1 Requirements Traceability

| Business Objective | Related Functional Requirements | Success Metrics |
|--------------------|--------------------------------|-----------------|
| Digital Onboarding: 40,000 policies by 2026 | FR-001 to FR-016, FR-033 to FR-040 | Monthly policy acquisition rate |
| API Performance Optimization | FR-107 to FR-118, NFR-008 to NFR-011 | API response time metrics |
| Financial Transaction Integrity | FR-121 (TigerBeetle), SEC-003 (PCI-DSS) | Transaction accuracy and speed |
| Regulatory Compliance | SEC-011 to SEC-020 (IDRA/AML/CFT) | Audit compliance score |
| Partner Management Excellence | FR-011 (Focal Person), FR-086 to FR-092 | Number of active partners |
| Claims Efficiency | FR-041 to FR-058, FR-133 to FR-137 | Average claim TAT |

### 12.2 Change Control Process

**Approval Hierarchy:**
1. Dev submits change request
2. Repository Admin reviews code changes
3. Database Admin reviews data model impact
4. System Admin reviews infrastructure impact
5. Business Admin approves business impact
6. Focal Person approves partner-related changes

---

## Appendices

### Appendix A — API Architecture Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                    CLOUDFLARE PROXY                         │
└─────────────────────┬───────────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────────┐
│                    NGINX GATEWAY                            │
└─────────────────────┬───────────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────────┐
│  PUBLIC API (REST/JSON) │ CATEGORY 2 (GraphQL + JWT)        │
│  - Product Search       │ - Customer Device                  │
│  - Product List         │ - Mobile Apps                      │
└─────────────────────────┼───────────────────────────────────┘
                         │
              ┌─────────▼─────────┐
              │  API GATEWAY      │
              │  (OAuth2 + JWT)   │
              └─────────┬─────────┘
                       │
    ┌─────────────────┼─────────────────┐
    │                 │                 │
┌───▼────┐ ┌─────────▼──────────┐ ┌─────▼─────────┐
│CAT 3   │ │     CATEGORY 1     │ │   INTERNAL    │
│REST API│ │  (gRPC + ProtoBuf) │ │  MICROSERVICES│
│3rd     │ │  Microservices     │ │   (Kafka      │
│Party   │ │  Communication     │ │   Orchestrated)│
└────────┘ └────────────────────┘ └───────────────┘
```

### Appendix B — Stakeholder Hierarchy

```
System Admin (Cloud/IAM Root)
    │
    ├── Repository Admin (Code/Deploy)
    ├── Database Admin (Data Management)
    ├── Business Admin (Business Operations)
    │   └── Focal Person (Partner Bridge) ★ KEY ROLE
    │       └── Partner Admin (Tenant Root)
    │           ├── Insurer Agents (Partner Staff)
    │           └── Insurer Underwriters (API Users)
    ├── Dev (Development)
    └── Support/Call Centre (Customer Assistance)
```

### Appendix C — Claims Approval Workflow

```
┌──────────────────┐
│ Customer submits │
│ claim via app    │
└────────┬─────────┘
         ▼
┌──────────────────┐
│ System validates │
│ and screens      │
└────────┬─────────┘
         ▼
    ┌────────┐
    │Amount? │
    └───┬────┘
        │
   ┌────┴─────┬──────────┬──────────┐
   │          │          │          │
   ▼          ▼          ▼          ▼
0-10K    10K-50K    50K-2L      2L+
   │          │          │          │
   ▼          ▼          ▼          ▼
Auto/L1    L2 Mgr   L3 Joint    Board
Officer              BA+FP
   │          │          │          │
   └────┬─────┴────┬─────┴──────────┘
        │          │
        ▼          ▼
   Approved    Rejected
        │          │
        ▼          ▼
   Payment    Notification
   Process    to Customer
```

### Appendix D — Summary Statistics

**Requirements Count:**
- Total Functional Requirements: 250+ (FR-001 to FR-251)
- Total Non-Functional Requirements: 47 (NFR-252 to NFR-298)
- Total IDRA Compliance Requirements: 10 (IDRA-001 to IDRA-010)
- Total BFIU Compliance Requirements: 10 (BFIU-001 to BFIU-010)
- Total AML Monitoring Rules: 30 (TM-001 to TM-030)

**By Priority:**
- M1 (Must Have - Phase 1): ~58 requirements
- M2 (Must Have - Phase 2): ~45 requirements
- D (Desirable): ~32 requirements
- F (Future): ~15 requirements

### Appendix E — Technology Stack Summary

**Programming Languages:** Go, C# .NET 8, Node.js, Python, TypeScript/React, React Native

**Data & Persistence:** Protocol Buffers, PostgreSQL, MongoDB, Redis, Apache Kafka

**Infrastructure:** Docker, Kubernetes, AWS/Azure, Prometheus/Grafana, Jaeger

### Appendix F — Compliance Checklist

**IDRA Requirements:**
- [ ] Digital product approval documentation
- [ ] Customer data protection policies
- [ ] Policy issuance standards compliance
- [ ] Claims processing procedures
- [ ] Financial reporting capabilities
- [ ] Agent licensing and training systems
- [ ] Marketing compliance monitoring
- [ ] Risk management frameworks
- [ ] Audit trail systems
- [ ] Regulatory change management

**BFIU AML Checklist:**
- [ ] Customer due diligence (CDD) procedures
- [ ] Enhanced due diligence (EDD) for high-risk customers
- [ ] Suspicious transaction monitoring systems
- [ ] Cash transaction reporting mechanisms
- [ ] Wire transfer monitoring capabilities
- [ ] PEP screening integration
- [ ] Sanctions list screening
- [ ] AML record retention policies
- [ ] Employee training programs
- [ ] Compliance reporting dashboards

### Appendix G — Integration Endpoints

```
bKash Merchant API:
  Base URL: https://api.bkash.com/v1/
  Authentication: Bearer token
  Key Endpoints:
    - POST /payments/create
    - POST /payments/execute
    - GET /payments/{paymentID}

Nagad Merchant API:
  Base URL: https://api.nagad.com.bd/api/dfs/
  Authentication: API key + signature
  Key Endpoints:
    - POST /check-out/initialize
    - POST /check-out/complete
    - GET /verify/payment/{paymentRef}

Government NID API:
  Base URL: https://api.nidw.gov.bd/v1/
  Authentication: Government issued certificates
  Key Endpoints:
    - POST /verify/nid
    - GET /citizen/details/{nid}
```

### Appendix H — Proto Schema Definitions

#### User Entity

```protobuf
syntax = "proto3";
package insuretech.authn.entity.v1;

message User {
  string user_id = 1;
  string mobile_number = 2;       // +880 1XXX XXXXXX
  string email = 3;
  string password_hash = 4;
  UserStatus status = 5;
  google.protobuf.Timestamp created_at = 6;
  google.protobuf.Timestamp updated_at = 7;
  google.protobuf.Timestamp last_login_at = 8;
  string created_by = 9;
  int32 login_attempts = 10;
  google.protobuf.Timestamp locked_until = 11;
}

enum UserStatus {
  USER_STATUS_UNSPECIFIED = 0;
  USER_STATUS_PENDING_VERIFICATION = 1;
  USER_STATUS_ACTIVE = 2;
  USER_STATUS_SUSPENDED = 3;
  USER_STATUS_LOCKED = 4;
  USER_STATUS_DELETED = 5;
}
```

#### Policy Entity

```protobuf
syntax = "proto3";
package insuretech.policy.entity.v1;

message Policy {
  string policy_id = 1;
  string policy_number = 2;         // LBT-YYYY-XXXX-NNNNNN
  string product_id = 3;
  string customer_id = 4;
  string partner_id = 5;
  string agent_id = 6;
  PolicyStatus status = 7;
  double premium_amount = 8;
  double sum_insured = 9;
  int32 tenure_months = 10;
  google.protobuf.Timestamp start_date = 11;
  google.protobuf.Timestamp end_date = 12;
  google.protobuf.Timestamp issued_at = 13;
  repeated Nominee nominees = 17;
  repeated Rider riders = 18;
}

enum PolicyStatus {
  POLICY_STATUS_UNSPECIFIED = 0;
  POLICY_STATUS_PENDING_PAYMENT = 1;
  POLICY_STATUS_ACTIVE = 2;
  POLICY_STATUS_GRACE_PERIOD = 3;
  POLICY_STATUS_LAPSED = 4;
  POLICY_STATUS_SUSPENDED = 5;
  POLICY_STATUS_CANCELLED = 6;
  POLICY_STATUS_EXPIRED = 7;
}
```

#### Claim Entity

```protobuf
syntax = "proto3";
package insuretech.claims.entity.v1;

message Claim {
  string claim_id = 1;
  string claim_number = 2;          // CLM-YYYY-XXXX-NNNNNN
  string policy_id = 3;
  string customer_id = 4;
  ClaimStatus status = 5;
  ClaimType type = 6;
  double claimed_amount = 7;
  double approved_amount = 8;
  double settled_amount = 9;
  google.protobuf.Timestamp incident_date = 10;
  string incident_description = 11;
  repeated ClaimDocument documents = 12;
  repeated ClaimApproval approvals = 13;
  FraudCheckResult fraud_check = 20;
}

enum ClaimStatus {
  CLAIM_STATUS_UNSPECIFIED = 0;
  CLAIM_STATUS_SUBMITTED = 1;
  CLAIM_STATUS_UNDER_REVIEW = 2;
  CLAIM_STATUS_PENDING_DOCUMENTS = 3;
  CLAIM_STATUS_APPROVED = 4;
  CLAIM_STATUS_REJECTED = 5;
  CLAIM_STATUS_SETTLED = 6;
  CLAIM_STATUS_DISPUTED = 7;
}
```

#### Payment Entity

```protobuf
syntax = "proto3";
package insuretech.payment.entity.v1;

message Payment {
  string payment_id = 1;
  string transaction_id = 2;
  string policy_id = 3;
  string claim_id = 4;
  PaymentType type = 5;
  PaymentMethod method = 6;
  PaymentStatus status = 7;
  double amount = 8;
  string currency = 9;              // BDT
  string payer_id = 10;
  string payee_id = 11;
  string gateway = 15;              // bKash, Nagad, SSLCommerz, etc.
  int32 retry_count = 18;
}

enum PaymentMethod {
  PAYMENT_METHOD_UNSPECIFIED = 0;
  PAYMENT_METHOD_BKASH = 1;
  PAYMENT_METHOD_NAGAD = 2;
  PAYMENT_METHOD_ROCKET = 3;
  PAYMENT_METHOD_CARD = 4;
  PAYMENT_METHOD_BANK_TRANSFER = 5;
  PAYMENT_METHOD_CASH = 6;
  PAYMENT_METHOD_CHEQUE = 7;
}
```

---

## Sign-off & Approval

| Role | Name | Signature | Date |
|------|------|-----------|------|
| Insuretech Director | — | _________________ | ___/___/2025 |
| Chief Executive Officer | — | _________________ | ___/___/2025 |
| Chief Technology Officer | — | _________________ | ___/___/2025 |
| Chief Financial Officer | — | _________________ | ___/___/2025 |
| Business Head InsureTech | — | _________________ | ___/___/2025 |
| Project Manager InsureTech | — | _________________ | ___/___/2025 |
| Senior Dev LifePlus | — | _________________ | ___/___/2025 |

---

**Document Status**
- Version: 3.11
- Date: February 2026
- Status: FINAL_DRAFT

---

*End of Document*
