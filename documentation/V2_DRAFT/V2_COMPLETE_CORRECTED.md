# System Requirements Specification (SRS)
**Project:** Labaid InsureTech Platform  
**Version:** V2.0 Complete (Corrected)  
**Date:** December 2024  
**Document Classification:** Bangladesh Tender Specification  
**Status:** Draft  
**Control Level:** A

---

## Table of Contents

1. [Introduction](#1-introduction)
2. [Market Context](#2-market-context)
3. [Overall Description](#3-overall-description)
4. [System Features & Functional Requirements](#4-system-features--functional-requirements)
5. [External Interface Requirements](#5-external-interface-requirements)
6. [Non-Functional Requirements](#6-non-functional-requirements)
7. [Data Model & Storage Requirements](#7-data-model--storage-requirements)
8. [Security & Compliance Requirements](#8-security--compliance-requirements)
9. [Performance & Scalability Requirements](#9-performance--scalability-requirements)
10. [AML/CFT Compliance Requirements](#10-amlcft-compliance-requirements)
11. [Operational Requirements & Support](#11-operational-requirements--support)
12. [Acceptance Criteria & Test Summary](#12-acceptance-criteria--test-summary)
13. [Traceability Matrix & Change Control](#13-traceability-matrix--change-control)
14. [Appendices](#14-appendices)

---

## 1. Introduction

### 1.1 Purpose
This SRS documents detailed system-level requirements for the LabAid InsureTech platform (mobile apps, partner portals, admin portal, and backend services). Its purpose is to provide an unambiguous specification for design, development, integration, testing, deployment, and operations teams.

### 1.2 Scope
The system will enable digital onboarding, product discovery, policy purchase, digital KYC, payment processing, claims submission and tracking, partner integrations, and admin/insurer workflows. This SRS covers Phase 1 (core digital capabilities) while flagging Phase 2/3 enhancements (AI underwriting, Voice-aided guidance, IoT based connectivity) for future releases. See Business Plan for strategic context.

### 1.3 Market Context
With a low penetration rate of insurance compared to other countries in the region and considering mass literacy level and lack of awareness about insurance, a very simplified onboarding flow with step-by-step explanations and visual cues are required for users in the UI/UX space. Also considering rural internet network constraints, overall speed (minimum 3G), minimum specifications of edge devices, and lack of availability of local cloud storage support within the country for reducing latency, specific optimization strategies are essential.

### 1.4 Definitions, Acronyms & Abbreviations

| Term | Definition |
|------|------------|
| **IDRA** | Insurance Development & Regulatory Authority (Bangladesh) |
| **KYC** | Know Your Customer |
| **KYB** | Know Your Business |
| **FCR** | Financial Condition Report (IDRA requirement) |
| **CARAMELS** | Capital Adequacy, Reinsurance Arrangements, Management, Earnings, Liquidity & Asset quality, Sensitivity |
| **SAR** | Suspicious Activity Report (AML requirement) |
| **BFIU** | Bangladesh Financial Intelligence Unit |
| **MFS** | Mobile Financial Service (bKash, Nagad, Rocket) |
| **API** | Application Programming Interface |
| **EHR** | Electronic Health Record |
| **ZHCT** | Zero Human Touch Claims |
| **UBI** | Usage Based Insurance |
| **IAM** | Identity and Access Management |
| **ACL** | Access Control List |
| **RBAC** | Role Based Access Control |
| **ABAC** | Attribute Based Access Control |
| **UAT** | User Acceptance Testing |
| **MOU** | Memorandum of Understanding |

---

## 2. Market Context

Target user categories for Bangladesh insurance market:
- **Type 1 Users:** Urban professionals with high digital literacy
- **Type 2 Users:** Semi-urban population with moderate digital literacy  
- **Type 3 Users:** Rural population with low digital literacy requiring voice-assisted workflows

---

## 3. Overall Description

### 3.1 Product Perspective
The platform is cloud-native, microservices-based, and mobile-first. Components:
- **Mobile apps (iOS/Android)** — Customer experience with Bangladesh-specific UI/UX considerations
- **Partner/Agent Portal** — Web-based onboarding & embedded flows for MFS, hospitals, e-commerce
- **Admin Portal** — Product, pricing, user and claims management with regulatory reporting
- **Backend Services** — Microservices with specialized databases and multi-protocol APIs
- **Third-party Integrations** — Insurer APIs, LabAid systems, MFS/telco APIs

### 3.2 User Classes & Characteristics
- **Customers:** Mobile-first, diverse digital literacy (urban youth → rural farmers)
- **Agent/Partners:** Staff at partner organisations (MFS, hospitals, e-commerce)
- **Insurer Underwriters:** Internal to partner insurers (API consumers)
- **Admin:** Business admins (policy management, product updates)
- **Support/Call Centre:** Assisted onboarding & claims help
- **Compliance Officers:** Regulatory compliance and AML/CFT monitoring

### 3.3 Operating Environment
- **Cloud:** AWS/Azure preferred; region: Bangladesh (or compliant region)
- **Mobile:** iOS 13+/Android 9+
- **Browsers for portals:** Chrome, Firefox, Edge (latest two versions)

---

## 4. System Features & Functional Requirements

Each requirement below has a unique ID (FR-XXX). Priority levels: M = Mandatory (Phase 1), D = Desirable (Phase 1), F = Future (Phase 2/3).

### 4.1 Authentication & User Management

| ID | Requirement Description | Priority |
|----|-------------------------|----------|
| FR-001 | The system shall support phone-based registration with OTP validation (SMS) and capture minimal profile fields (name, DOB, phone, email optional) | M |
| FR-002 | The system shall provide login: OTP login and password-based login; session management and refresh tokens | M |
| FR-003 | The system shall enable profile management: Update personal info, nominee details, and document uploads | M |
| FR-004 | The system shall implement duplicate prevention: Block duplicate accounts by national ID/phone; provide merge or support flow | M |

### 4.2 Digital KYC & Document Verification

| ID | Requirement Description | Priority |
|----|-------------------------|----------|
| FR-005 | The system shall provide document upload: Capture images/PDFs for NID/passport, photos, medical docs. UI flows per screens | M |
| FR-006 | The system shall implement OCR & validation: Extract NID/passport fields via OCR and pre-validate (format + checksum) | M |
| FR-007 | The system shall integrate eKYC integration: Integrate with third-party eKYC services for automated verification (Phase 1 if available) | D |
| FR-016 | The system shall lock verified stakeholder KYB data from unauthorized updates | M |

### 4.3 Product Catalog & Policy Discovery

| ID | Requirement Description | Priority |
|----|-------------------------|----------|
| FR-008 | The system shall maintain product catalog API: Serve product list with metadata (name, coverage, premiums, insurer, T&Cs) | M |
| FR-009 | The system shall provide filter & compare: Support filtering (category, premium, coverage) and side-by-side comparison (up to 3). UI screens show these flows | M |
| FR-010 | The system shall display product detail: Full policy wording, exclusions, and illustrative premium breakdown | M |

### 4.4 Policy Purchase & Issuance

| ID | Requirement Description | Priority |
|----|-------------------------|----------|
| FR-011 | The system shall implement multi-step purchase flow: Personal details, nominee, document upload, review, payment, confirmation. Flow follows provided mockups | M |
| FR-012 | The system shall provide premium calculation: Request insurer API or local pricing engine and show breakdown | M |
| FR-013 | The system shall support payment integration: Support bKash, Nagad, Rocket, card (PCI-DSS compliant), bank transfer | M |
| FR-014 | The system shall implement policy issuance: On successful payment, generate a digital policy certificate (PDF) and store in user account | M |
| FR-015 | The system shall support promo/discounts: Support coupon codes and partner-discount flows | D |

### 4.5 Claims Management

| ID | Requirement Description | Priority |
|----|-------------------------|----------|
| FR-016 | The system shall enable claim initiation: Policy prefill, claim reason selection, document upload (images, bills), and submission. UI Claim screen referenced | M |
| FR-017 | The system shall provide status tracking: Provide claim status updates and admin notes to user | M |
| FR-018 | The system shall implement admin workflow: Claims dashboard for triage, verification, approval, rejection, and payment initiation | M |
| FR-019 | The system shall provide automated triage: OCR, image verification and rule-based auto-accept for small claims (Phase 2 target) | D |

### 4.6 Policy Management & Renewals

| ID | Requirement Description | Priority |
|----|-------------------------|----------|
| FR-020 | The system shall provide policy dashboard: Active & past policies, download documents, renewal prompts. (UI page: Policy History) | M |
| FR-021 | The system shall implement renewals: Auto-renew option and manual renew flows, with reminders | M |
| FR-022 | The system shall support partial adjustments: Allow address, nominee change requests | D |

### 4.7 Notifications & Communication

| ID | Requirement Description | Priority |
|----|-------------------------|----------|
| FR-023 | The system shall implement notification engine: Trigger SMS/Push/Email for OTP, purchase confirmation, claims updates, renewal reminders | M |
| FR-024 | The system shall provide marketing opt-in/opt-out: Manage user preferences | D |

### 4.8 Admin & Reporting

| ID | Requirement Description | Priority |
|----|-------------------------|----------|
| FR-025 | The system shall provide admin portal: User management, product management, claim dashboards, and role-based access control | M |
| FR-026 | The system shall generate reports: Daily sales, claims ratio, partner performance, policy counts and KPIs (aligned to business plan targets) | M |

### 4.9 Partner / Agent Portal

| ID | Requirement Description | Priority |
|----|-------------------------|----------|
| FR-027 | The system shall provide embedded flow: Partner can initiate policy purchase for end-customer via API or partner portal | M |
| FR-028 | The system shall implement partner dashboard: Commission statements, leads, onboarding analytics | M |
| FR-076 | The system shall support promo code for product discount | M |

### 4.10 AI & Automation Features

| ID | Requirement Description | Priority |
|----|-------------------------|----------|
| FR-089 | The system shall have purpose built AI Chatbot to assist customer during product search, selection, purchase, verification, payment stage | F |
| FR-123 | The system shall provide ticket form fill up option with auto recorded customer support call | D |

### 4.11 Business Intelligence

| ID | Requirement Description | Priority |
|----|-------------------------|----------|
| FR-139B | The system shall have internal Admin special data product as Business Intelligence Tool: TBD (Options: Metabase, Tableau, Power BI, custom dashboards). Data Source: Read replica of production database. Update Frequency: Near real-time (5-minute lag acceptable). Dashboards: Executive (Daily sales, policy count, claims ratio, revenue), Operations (System health, support tickets, processing times), Compliance (AML flags, IDRA report status, audit logs) | F |

### 4.12 Audit & Logging

| ID | Requirement Description | Priority |
|----|-------------------------|----------|
| FR-141 | The system shall maintain immutable logs for critical actions such as: policy issue, claim approval, claim rejection, payment and dispute | M |
| FR-142 | The system shall maintain data retention policy up to minimum of 20 years to maintain records for regulatory compliance | D |
| FR-143 | The system shall track each logged in user for auxiliary actions and will have additional data logs | D |
| FR-144 | The system shall allow partner to maintain additional data logs as per customer and InsureTech MOU | F |
| FR-145 | The system shall provide special portal to regulatory body to access requested data as per policy and law of regulatory bodies | M |

### 4.13 User Interface Requirements

| ID | Requirement Description | Priority |
|----|-------------------------|----------|
| FR-146 | The system shall maintain similar User interface with different operating systems | M |
| FR-147 | The system shall provide smart data widget for mobile users | D |
| FR-148 | The system shall provide voice assisted workflow for type 3 users | F |
| FR-149 | The system shall provide desktop first web UI for portals | M |
| FR-150 | The system shall take minimum permissions for all services from user device (e.g., camera, one time message read) | M |

### 4.14 API Design and Data Flow

| ID | Requirement Description | Priority |
|----|-------------------------|----------|
| FR-151 | The system shall maintain two levels of API: private and public | M |
| FR-152 | The system shall have three levels of private API as Category 1, Category 2, and Category 3 with secure middle layer | D |
| FR-153 | The system shall use Protocol Buffer and gRPC based communication for Category 1 API between gateway and microservices in HTTP with system admin middle layer | F |
| FR-154 | The system shall use GraphQL based API as Category 2 for gateway and customer device with JWT token system and OAuth v2 based security middle layer | M |
| FR-155 | The system shall use RESTful API with JSON with compliance of OPEN API standard as Category 3 for 3rd party Integration with server side auth and shall provide user public documentation, wiki, sandbox and a mock server | D |
| FR-156 | The system shall provide public RESTful API with JSON with compliance of OPEN API standard for product search, product list and view for live products etc | M |
| FR-157 | The system shall only expose proxy server (Cloudflare) and entry node IP (NGINX) to public space | M |
| FR-158 | The system shall block all access to microservices for all available API except Category 1 | M |
| FR-159 | The system shall use InsureTech Internal protocol for IoT based data extraction and data binding | F |
| FR-160 | The system shall consolidate, annotate, add context and process and save data within regulatory limitation to train internal AI agents | F |
| FR-161 | The system shall generate statistics, prediction based on big data and provide special projection data to partners given special agreement with InsureTech | F |
| FR-162 | The system shall implement WebSocket connection for real-time updates | D |

### 4.15 Data Storage

| ID | Requirement Description | Priority |
|----|-------------------------|----------|
| FR-165 | The system shall store all uploaded objects with AWS S3 or equivalent data storage system | M |
| FR-166 | The system shall store relational data on self-hosted PostgreSQL V17 database (ACID compliance, JSONB support, Bengali full-text search) with Connection pooling: PgBouncer, max 100 connections per service, Read replicas for reporting queries, Query optimization: <100ms for 95% queries, Indexing: Comprehensive indexes on foreign keys, status fields, dates | M |
| FR-167 | The system shall store finance transaction records, debit-credit records, insurance records and contracts on special purpose built database (like TigerBeetle or equivalent) | D |
| FR-168 | The system shall use cache database for faster data delivery - Product catalog: 5-minute TTL, User sessions: Redis with 15-minute expiry | M |
| FR-169 | The system shall process tokenized data on vector database like Pgvector (PostgreSQL) or PineCone | D |
| FR-170 | The system shall store app native encrypted data in user device in SQLite | M |
| FR-171 | The system shall store product catalog and metadata, unstructured data in a NoSQL Database (AWS DynamoDB/MongoDB) | M |
| FR-172 | Upload data policy - Client-side compression: 5MB → 1-2MB (JPEG 80% quality, 1920×1080 max resolution), Chunked upload: 1MB chunks with resume capability (tus.io protocol), Presigned S3 URLs: Direct upload, 30-minute expiry | M |
| FR-173 | Backup: Daily full, 6-hour incremental, continuous transaction logs | F |
| FR-174 | Retention: 20 years for regulatory compliance, tiered storage (hot/warm/cold) | F |

---

## 5. External Interface Requirements

### 5.1 User Interfaces
- **Mobile app screens** consistent with provided prototype (purchase, review, payment, claim, history)
- **Web admin/partner portals** modern UI (desktop-first)
- **Voice-assisted UI** for Type 3 users (low digital literacy)

### 5.2 API Architecture (Your Exact Specifications)

**API Category Structure:**

| API Category | Protocol | Use Case | Security Layer | Performance Target |
|--------------|----------|----------|----------------|-------------------|
| **Category 1** | Protocol Buffer + gRPC | Gateway ↔ Microservices | System Admin Middle Layer | < 100ms |
| **Category 2** | GraphQL + JWT | Gateway ↔ Customer Device | JWT + OAuth v2 | < 2 seconds |
| **Category 3** | RESTful + JSON (OpenAPI) | 3rd Party Integration | Server-side Auth | < 200ms |
| **Public API** | RESTful + JSON (OpenAPI) | Product Search/List | Public Access | < 1 second |

### 5.3 Third-Party Systems
- **Payment Gateways:** bKash, Nagad, Rocket, Card processor (PCI scope)
- **Insurer APIs:** Premium/underwriting/policy issuance
- **LabAid EHR System:** For cashless verification (IPD) and pre-auth (secured HL7/FHIR endpoint or APIs)
- **SMS gateway provider** for OTP & alerts

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

## 7. Data Model & Storage Requirements (Your Exact Architecture)

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

### 7.2 Core Data Entities

**User Management:**
- User (UserID, name, phone, email, NID, DOB, KYC status)
- Profile (UserID, nominee details, documents)

**Policy Management:**
- Policy (PolicyID, productID, insured period, premium, status, insurer reference)
- Product (productID, coverage, T&Cs, premium table)

**Financial Transactions (TigerBeetle):**
- Transaction (TransactionID, amount, gateway, status)
- Account (AccountID, balance, currency)

**Claims Management:**
- Claim (ClaimID, policyID, date, documents, status, admin notes)
- AuditLog (event, user, timestamp)

---

## 8. Security & Compliance Requirements

### 8.1 Security Infrastructure (Your Exact Requirements)

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

### 8.2 Authentication & Authorization Framework
- JWT or OAuth2 for sessions
- Role-based access control for portals
- **JWT tokens expire after 15 minutes**

### 8.3 Data Protection
- PII encrypted at rest
- Masking of NID in logs and limited access to raw documents

---

## 9. Performance & Scalability Requirements

### 9.1 Performance Benchmarks (Your Exact Targets)

| API Category | Response Time Target | Load Target | Scaling Strategy |
|--------------|---------------------|-------------|------------------|
| **Category 1 (gRPC)** | < 100ms | High throughput | Microservices auto-scaling |
| **Category 2 (GraphQL)** | < 2 seconds | Mobile optimization | Connection pooling |
| **Category 3 (REST)** | < 200ms | 3rd party integration | API gateway scaling |
| **Public API** | < 1 second | Public access | CDN + caching |

### 9.2 Scalability Architecture
- Peak TPS (Transactions per second) for purchase flow: baseline 50 TPS, scale to 500 TPS via auto-scaling
- Concurrent active sessions target: 100k with <5% degradation
- File upload: support 10 MB images with background upload/resume

### 9.3 Database Performance
- **PostgreSQL:** <100ms for 95% queries, connection pooling with PgBouncer
- **TigerBeetle:** <10ms transaction processing
- **Redis:** <1ms cache access time

---

## 10. AML/CFT Compliance Requirements

### 10.1 Customer Due Diligence (CDD) Framework

**Mandatory CDD Requirements for Bangladesh:**
- Identity Verification: NID/Passport verification via approved eKYC (BFIU Guidelines)
- Address Verification: Utility bill or bank statement (MLPA Requirements)
- Photo Identification: Selfie with liveness detection (Enhanced CDD)
- Source of Funds: Income declaration for high-value policies (Risk-based approach)
- PEP Screening: Automated screening against watchlists (FATF Recommendations)

### 10.2 Risk-Based Customer Categorization

| Risk Level | Criteria | CDD Requirements | Monitoring Frequency |
|------------|----------|------------------|---------------------|
| **Low Risk** | Standard customers, low premium policies | Standard CDD | Annual review |
| **Medium Risk** | Higher premiums, multiple policies | Enhanced documentation | Quarterly review |
| **High Risk** | PEPs, large premiums, suspicious patterns | Enhanced Due Diligence (EDD) | Monthly monitoring |
| **Prohibited** | Sanctioned individuals, blocked entities | Transaction rejection | Real-time blocking |

### 10.3 Automated AML Monitoring Rules (Your Specifications)

**20+ Automated Rules for AML Detection:**
- Rapid purchases (>3 policies in 7 days)
- High-value premiums (>BDT 5 lakh)
- Frequent cancellations
- Mismatched nominees
- Geographic/payment anomalies

### 10.4 Record Keeping & Audit Trail

| Document Type | Retention Period | Storage Requirements | Access Controls |
|---------------|------------------|---------------------|-----------------|
| CDD Documentation | 5+ years after relationship end | Encrypted PostgreSQL + S3 | Compliance team only |
| Transaction Records | 7+ years | TigerBeetle + Archive | Audit and compliance |
| STR Documentation | 10+ years | Secured offline storage | Senior management |
| Training Records | 5+ years | HR system integration | HR and compliance |
| System Audit Logs | 20+ years | Immutable PostgreSQL logging | System administrators |

---

## 11. Operational Requirements & Support

### 11.1 System Monitoring & Alerting
- 24x7 monitoring with alerts (PagerDuty)
- Centralized logging (ELK/Datadog) and metrics
- Runbooks for incident management
- Support SLA: 1 hour for P1 incidents, 4 hours for P2

### 11.2 Maintenance & Updates
- Regular Maintenance: Every Sunday 2:00 AM - 4:00 AM Bangladesh Time
- Security Updates: Emergency patches within 24 hours of release
- Feature Releases: Monthly deployment schedule with staged rollout
- Database Maintenance: Quarterly optimization during low-traffic periods

---

## 12. Acceptance Criteria & Test Summary

### 12.1 High-Level Acceptance Criteria
- FR-critical workflows (Registration → Purchase → Payment → Policy Issuance → Claim Submission) operate end-to-end in the staging environment without manual intervention for standard products
- Payment callback reliability > 99.9%
- OTP deliverability > 98% and verification success > 95%
- Policy issuance within 30 seconds for auto-issuance products

### 12.2 API Performance Testing

| API Category | Load Profile | Success Criteria | Expected Performance |
|---------------|--------------|------------------|---------------------|
| **Category 1 (gRPC)** | 1000 concurrent requests | <100ms response time | High-throughput internal communication |
| **Category 2 (GraphQL)** | 500 concurrent requests | <2s response time | Mobile-optimized data fetching |
| **Category 3 (REST)** | 100 concurrent requests | <200ms response time | Standard 3rd party integration |
| **Public API** | 50 concurrent requests | <1s response time | Public product search |

### 12.3 Database Performance Testing

| Storage System | Test Scenario | Success Criteria | Performance Target |
|----------------|---------------|------------------|-------------------|
| **PostgreSQL V17** | 95% queries under load | <100ms response time | ACID compliance maintained |
| **TigerBeetle** | Financial transaction processing | <10ms latency | High-accuracy financial ops |
| **Redis Cache** | Session and catalog access | <1ms response time | 95%+ cache hit rate |

---

## 13. Traceability Matrix & Change Control

### 13.1 Requirements Traceability

**Business Objective Mapping:**
- **Digital Onboarding Target: 40,000 policies by 2026** → FR-001 to FR-007, FR-011 to FR-015
- **API Performance Optimization** → FR-151 to FR-162, NFR-008 to NFR-011
- **Financial Transaction Integrity** → FR-167 (TigerBeetle), SEC-003 (PCI-DSS)
- **Regulatory Compliance** → SEC-004 (AML/CFT), SEC-005 (IDRA reporting)

### 13.2 Change Control Process
- Changes must be submitted via RFC with business justification
- Technical impact analysis for API categories and database systems
- Security assessment for PCI-DSS and AML/CFT compliance
- Steering committee approval for scope changes affecting timeline or budget

---

## 14. Appendices

### 14.1 API Architecture Reference

**Your Exact API Design:**

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
│3rd     │ │  Microservices     │ │               │
│Party   │ │  Communication     │ │               │
└────────┘ └────────────────────┘ └───────────────┘
```

### 14.2 Database Architecture Reference

**Your Exact Database Strategy:**

```
┌──────────────────────────────────────────────────────────────┐
│                    APPLICATION LAYER                         │
└─────────────────────┬────────────────────────────────────────┘
                      │
    ┌─────────────────┼─────────────────┐
    │                 │                 │
┌───▼────────┐ ┌─────▼──────┐ ┌─────────▼─────────┐
│PostgreSQL  │ │TigerBeetle │ │   Redis Cache     │
│V17         │ │Financial   │ │   Session Store   │
│- User Data │ │Transactions│ │   - 15min TTL     │
│- Policies  │ │- ACID      │ │   - Product Cache │
│- Claims    │ │- <10ms     │ │   - 5min TTL      │
│- Audit     │ │            │ │                   │
└────────────┘ └────────────┘ └───────────────────┘
                                         │
    ┌────────────────────────────────────┼────────────────┐
    │                                    │                │
┌───▼─────────┐ ┌─────────────────┐ ┌───▼──────────┐ ┌───────────┐
│DynamoDB/    │ │    AWS S3       │ │SQLite        │ │Pgvector/  │
│MongoDB      │ │   Documents     │ │Mobile Local  │ │PineCone   │
│- Product    │ │   - KYC Docs    │ │- Offline     │ │Vector DB  │
│  Catalog    │ │   - Claims      │ │  Capability  │ │- AI Data  │
│- Metadata   │ │   - Policies    │ │- Encrypted   │ │           │
└─────────────┘ └─────────────────┘ └──────────────┘ └───────────┘
```

---

## Document Approval & Sign-off

**Technical Architecture Confirmation:**
By signing below, stakeholders confirm their acceptance of this System Requirements Specification including the specific technical decisions:
- **API Architecture:** Category 1 (gRPC), Category 2 (GraphQL), Category 3 (REST), Public (REST)
- **Database Strategy:** PostgreSQL V17 + TigerBeetle + Redis + DynamoDB/MongoDB + S3 + SQLite + Vector DB
- **Performance Targets:** <100ms (gRPC), <2s (GraphQL), <200ms (REST), <1s (Public API)
- **Security Framework:** Multi-layer with PCI-DSS compliance and AML/CFT integration

| Role | Name | Signature | Date |
|------|------|-----------|------|
| **Product Owner / LabAid InsureTech** | __________________ | __________________ | ______ |
| **Chief Technology Officer** | __________________ | __________________ | ______ |
| **Database Architect** | __________________ | __________________ | ______ |
| **API Architect** | __________________ | __________________ | ______ |
| **Compliance Officer** | __________________ | __________________ | ______ |

---

*This document contains proprietary and confidential information. Distribution is restricted to authorized personnel only.*

**END OF DOCUMENT**