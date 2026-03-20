# System Requirements Specification (SRS)

**Project:** LabAid InsureTech Platform  
**Version:** 3.0 - Phased Delivery Edition  
**Date:** December 2025  
**Status:** Active - Baseline for Development  
**Control Level:** A  
**Classification:** Internal - Confidential

---

## Revision History

| Version | Date | Revised By | Description |
|---------|------|------------|-------------|
| 1.0 | Nov 2025 | Faruk Hannan | Initial draft with core requirements |
| 2.0 | Dec 06, 2025 | Faruk Hannan | Enhanced with novel features (Focal Person, Joint approvals, Voice assistance, IoT) |
| 2.1 | Dec 15, 2025 | Development Team | Professional formatting, comprehensive tables |
| 2.2 | Dec 16, 2025 | Development Team | MD feedback integrated, 42 new requirements, compliance enhanced |
| **3.0** | **Dec 18, 2025** | **Development Team** | **Phased delivery strategy with realistic priorities: M1 (Phase 1 - Mar 1), M2 (Phase 1.5 - May 1), S (Phase 2 - Nov 1), D/C/F (Phase 3 - Nov 2027)** |

---

## Approval Signatures

| Role | Name | Signature | Date |
|------|------|-----------|------|
| **Project Sponsor** | [Director Name] | _________________ | ___/___/2025 |
| **Project Manager** | Abid | _________________ | ___/___/2025 |
| **Chief Technology Officer** | [CTO Name] | _________________ | ___/___/2025 |
| **Chief Financial Officer** | [CFO Name] | _________________ | ___/___/2025 |
| **IDRA Compliance Officer** | [Compliance Name] | _________________ | ___/___/2025 |

---

## Document Change Control

Changes to this SRS require formal approval through RFC (Request For Change) process.

**Change Authority:**
- **Minor Changes** (clarifications, typos): PM approval
- **Major Changes** (new requirements, scope): Steering Committee approval
- **Phase Changes** (priority reassignment): CTO + CFO + Director approval

**Current Baseline:** V3.0 - Approved for Phase 1 Development (

## 2. Market Context & Strategic Rationale

### 2.1 Bangladesh Insurance Market Analysis

#### Current Market State

**Market Size & Penetration:**
- Insurance penetration: 0.55% of GDP (2024)
- Regional comparison: India 4.2%, Pakistan 1.0%, Sri Lanka 1.2%
- Total insured population: ~2 million (1.2% of 160M population)
- Annual growth rate: 8-10% (below economic growth)
- Market opportunity: 98.8% population uninsured

**Market Challenges:**
1. **Awareness Gap:** 73% of population unaware of insurance benefits (IDRA Survey 2023)
2. **Trust Deficit:** 45% distrust insurance companies due to claim denial stories
3. **Distribution Bottleneck:** Only 5,000 licensed agents for 160M population
4. **Process Friction:** Average policy purchase takes 3-7 days with paper forms
5. **Claims Delays:** Average TAT 21 days (international standard: 7 days)

**Technology Adoption:**
- Smartphone penetration: 62% (99M users)
- Mobile internet users: 112M (70% of population)
- MFS accounts: 180M registered, 60M active
- Digital payment transactions: 45% CAGR
- E-commerce growth: 55% CAGR

**Regulatory Environment:**
- IDRA actively promoting digital insurance (Circular 2023/05)
- Mandatory e-KYC integration announced for 2026
- Open API framework for insurtech partnerships
- Sandboxing regime for innovation
- Fast-track approvals for digital-first insurers

#### Market Opportunity

**Target Segments:**

1. **Urban Middle Class (Primary):**
   - Population: 25M households
   - Income: 50,000-150,000 BDT/month
   - Smartphone ownership: 85%
   - Digital payment adoption: 75%
   - Insurance awareness: Medium
   - **Target:** 5% penetration = 1.25M policies

2. **SME & Gig Economy (Secondary):**
   - Population: 8M workers
   - Income: Variable, 15,000-80,000 BDT/month
   - Delivery riders, drivers, freelancers
   - High smartphone usage (95%)
   - Need: Accident, health, equipment coverage
   - **Target:** 3% penetration = 240K policies

3. **Rural & Agricultural (Tertiary):**
   - Population: 90M (rural)
   - Farmers: 15M households
   - Livestock owners: 8M
   - Need: Crop, livestock, weather insurance
   - Digital literacy: Low to medium
   - **Target:** 1% penetration = 230K policies (Phase 2/3)

**Competitive Landscape:**

**Traditional Insurers:**
- 46 registered companies (32 life, 14 general)
- Paper-based processes
- Agent-dependent distribution
- No mobile apps (90% of companies)
- Average claim TAT: 21 days
- **Our Advantage:** Digital-first, instant issuance

**Emerging Insurtech:**
- 3 digital insurance startups (pre-revenue)
- Limited product range
- No regulatory licenses yet
- Funding rounds ongoing
- **Our Advantage:** LabAid brand trust, compliance-ready, hospital network

**MFS Embedded Insurance:**
- bKash exploring micro-insurance (pilot stage)
- Nagad partnering with traditional insurers
- Limited product types (accident only)
- **Our Advantage:** Full product range, claims management, partner ecosystem

### 2.2 Target User Segments & Micro-Insurance Value Proposition

**LabAid InsureTech Core Value Proposition: Affordable Micro-Insurance (200-2,000 BDT Premium Range)**

Unlike traditional insurance targeting high-income segments with premiums of 10,000+ BDT, LabAid focuses on **micro-insurance products with premiums ranging from 200 to 2,000 BDT**, making insurance accessible to mass-market Bangladesh. This pricing democratizes financial protection for previously unserved segments.

#### Micro-Insurance Product Categories

**1. Livestock Protection (Cattle, Poultry)**
- **Premium:** 500-2,000 BDT per animal per year
- **Coverage:** Death due to disease, accident, natural disaster
- **Target:** 15 million livestock-owning households
- **Pain Point:** Single cattle loss = 50,000-100,000 BDT impact on family
- **Claim Process:** Photo evidence + veterinary certificate → payout within 48 hours
- **IoT Enhancement (Phase 2):** GPS tracking for theft prevention

**2. Crop Insurance (Flood, Drought, Pest)**
- **Premium:** 300-1,500 BDT per season per acre
- **Coverage:** Crop failure due to flood, drought, pest infestation
- **Target:** Rice, jute, wheat farmers in flood-prone areas
- **Pain Point:** Entire season's income lost due to unpredictable weather
- **Claim Trigger:** Satellite data + ground verification → automatic payout
- **Government Partnership:** Potential subsidy collaboration

**3. Device Protection (Laptop, Smartphone)**
- **Premium:** 200-800 BDT per year
- **Coverage:** Accidental damage, theft, liquid damage
- **Target:** Students, gig workers, freelancers
- **Pain Point:** Laptop/phone = livelihood tool, replacement cost 30,000-80,000 BDT
- **Claim Process:** Police report (theft) or repair estimate (damage) → approval within 24 hours
- **Partnership:** E-commerce platforms (device purchase insurance at checkout)

**4. Personal Accident (Low Premium)**
- **Premium:** 300-1,200 BDT per year
- **Coverage:** 50,000-200,000 BDT for death/disability
- **Target:** Delivery riders, rickshaw drivers, construction workers
- **Pain Point:** No savings, family dependent on daily income
- **Claim Process:** Hospital certificate → payout within 48 hours

**5. Hospitalization (Basic Coverage)**
- **Premium:** 800-2,000 BDT per year
- **Coverage:** 25,000-100,000 BDT hospitalization expenses
- **Target:** Low-income families, domestic workers
- **Pain Point:** Medical debt spiral from single hospital visit
- **Claim Process:** Hospital bill + discharge summary → cashless or reimbursement

**6. Shop/Inventory Protection**
- **Premium:** 500-1,500 BDT per year
- **Coverage:** Fire, theft, natural disaster affecting small business inventory
- **Target:** 8 million small shop owners
- **Pain Point:** Shop inventory = life savings, no buffer for loss

#### Customer Segments (Revised for Micro-Insurance Focus)

**Type 1: Gig Economy Workers & Low-Income Urban (40% of Phase 1)**

**Demographics:**
- Age: 22-45 years
- Location: Urban and peri-urban areas
- Income: 15,000-40,000 BDT/month (daily wage earners)
- Occupation: Delivery riders, drivers, freelancers, domestic workers

**Digital Behavior:**
- Smartphone ownership: 85%
- Daily usage: 2-4 hours (mostly social media, entertainment)
- MFS: Primary financial tool (salary, bill payments)
- Price-sensitive, value-conscious

**Micro-Insurance Needs:**
- Personal accident (200-500 BDT) - #1 priority
- Device protection (laptop/phone for gig work) - 200-400 BDT
- Basic hospitalization - 800-1,200 BDT

**Affordability Sweet Spot:** 300-800 BDT total annual premium across 1-2 products

**Type 2: Small Farmers & Livestock Owners (35% of Phase 1)**

**Demographics:**
- Age: 30-60 years
- Location: Rural Bangladesh
- Income: 10,000-35,000 BDT/month (seasonal)
- Assets: 1-5 cattle/goats, 1-3 acres of land

**Digital Behavior:**
- Basic smartphone or feature phone
- Limited internet usage
- MFS adoption growing (bKash/Nagad via agents)
- Agent-assisted transactions common

**Micro-Insurance Needs:**
- Livestock protection (500-1,500 BDT per cattle) - Critical
- Crop insurance (300-800 BDT per acre) - Essential
- Basic personal accident - 300-500 BDT

**Affordability Sweet Spot:** 500-2,000 BDT per season (harvest time payment)

**Type 3: Small Business Owners & Shop Keepers (25% of Phase 1)**

**Demographics:**
- Age: 25-55 years
- Location: District towns, market areas
- Income: 20,000-60,000 BDT/month
- Business: Grocery shops, tea stalls, tailoring, hardware

**Digital Behavior:**
- Smartphone with business apps
- Digital payments adoption (QR code, bKash merchant)
- Aware of insurance value but deterred by high premiums

**Micro-Insurance Needs:**
- Shop/inventory insurance (500-1,500 BDT) - High interest
- Personal accident - 400-800 BDT
- Device protection (business laptop/tablet) - 300-600 BDT

**Affordability Sweet Spot:** 800-2,000 BDT annually (willing to pay for business protection)

#### Why Micro-Insurance Works for Bangladesh

**Economic Reality:**
- 70% of population earns <30,000 BDT/month
- Savings rate: <10% of income
- Traditional insurance premiums (5,000-20,000 BDT) = 2-6 months savings
- Micro-insurance (200-2,000 BDT) = 1-7 days income → **AFFORDABLE**

**Behavioral Insight:**
- People understand catastrophic risk (cattle death, crop loss, accident)
- Unwilling to pay 10,000 BDT "just in case"
- Willing to pay 500 BDT for specific, tangible protection

**Market Validation:**
- Grameen Phone's micro-insurance pilot: 500,000 subscribers in 18 months
- BRAC's livestock insurance: 80% renewal rate (affordable premiums key factor)
- India's Pradhan Mantri Fasal Bima Yojana: 57 million farmers enrolled (subsidized crop insurance)

**LabAid Advantage:**
- Hospital network enables low-cost health insurance (direct billing, no intermediary)
- Digital-first = 60% cost reduction vs. traditional insurance (no agents, no paperwork)
- MFS integration = instant premium collection + claim payout (no bank account needed)

#### Stakeholder Types

**Partners:**
- Hospitals (cashless settlement, EHR integration)
- MFS Providers (payment processing)
- E-commerce Platforms (embedded insurance)
- IoT Device Vendors (livestock trackers, health wearables)
- Corporate Clients (group insurance)

**Internal Staff:**
- System Administrators
- Business Administrators
- Claims Adjusters (L1, L2, L3)
- Customer Support
- Compliance Officers
- Focal Persons (partner management)

**Regulators:**
- IDRA (Insurance Development & Regulatory Authority)
- BFIU (Bangladesh Financial Intelligence Unit)
- ICT Division (Digital Security)

### 2.3 Competitive Positioning

**Differentiation Strategy:**

1. **Hospital Network:** Leverage LabAid hospital brand and network for cashless settlements
2. **Digital-First:** Mobile app with instant policy issuance (minutes vs. days)
3. **Transparent Claims:** Real-time tracking, automated approvals, 48-hour TAT
4. **Simplified Products:** Easy-to-understand coverage, no jargon
5. **Flexible Payments:** Cash→MFS→Bank→Card, installments (Phase 2)
6. **IoT Integration:** Usage-based insurance for livestock, vehicles (Phase 2/3)
7. **Vernacular Support:** Bengali-first design, voice assistance (Phase 3)

**Competitive Advantages:**

| Factor | Traditional Insurers | Emerging Insurtech | LabAid InsureTech |
|--------|---------------------|-------------------|------------------|
| Brand Trust | High (established) | Low (new) | **High (LabAid hospital)** |
| Digital Experience | Low (paper-based) | High (mobile) | **High (mobile-first)** |
| Product Range | Wide (100+ products) | Narrow (2-3 products) | **Medium (10-15 products)** |
| Distribution | Agents (5,000) | Digital only | **Hybrid (digital + agents)** |
| Claims TAT | 21 days | Unknown | **48 hours (target)** |
| Hospital Network | Partnerships (limited) | None | **LabAid network (owned)** |
| Regulatory Status | Licensed | Sandbox/Pending | **Licensed** |
| Technology | Legacy systems | Modern stack | **Modern stack** |

### 2.4 Technical Environment Constraints

#### Network Infrastructure

**Bangladesh Internet Landscape:**
- Mobile internet: 3G (70%), 4G (28%), 2G (2%)
- Fixed broadband: 12M connections (urban only)
- Average speed: 3G ~2 Mbps, 4G ~15 Mbps
- Latency: 80-150ms (domestic), 200-400ms (international)
- Frequent network switching (2G↔3G↔4G)
- Rural areas: 2G/3G only

**System Design Implications:**
- **Image optimization:** Compress to <100KB for profile photos
- **Lazy loading:** Load content incrementally
- **Offline-first:** PWA with service workers, local storage
- **Resilient APIs:** Retry logic, request queuing
- **Graceful degradation:** Low-bandwidth mode
- **CDN:** Local caching (minimal international requests)

#### Device Constraints

**Target Device Specifications:**
- **Minimum:** 4GB RAM, Android 9.0+, iOS 13.0+
- **Typical:** 6GB RAM, Android 11, iOS 14
- **Screen:** 5.5-6.5 inch displays
- **Storage:** 64GB typical (app should be <100MB installed)
- **Battery:** Optimize for 3000-4000mAh batteries

**App Design Implications:**
- Small app size (<10MB initial download)
- Incremental loading (up to 100MB cached data)
- Dark mode (battery saving)
- Low memory mode (disable animations)
- Storage cleanup (auto-delete old data)

#### Cloud & Hosting

**Bangladesh Data Center Landscape:**
- No major cloud providers locally (AWS, Azure, GCP)
- Nearest regions: Singapore, Mumbai
- Local hosting: Limited, lower tier (99% SLA)
- Hybrid approach needed

**Hosting Strategy:**
- **Primary:** AWS Singapore (low latency to BD)
- **CDN:** Cloudflare (local PoPs in Dhaka)
- **Backup:** Azure Mumbai (DR)
- **Static assets:** Local CDN
- **Future:** Local cloud when available

#### Payment Infrastructure

**MFS Ecosystem:**
- bKash: 60M users, 70% market share
- Nagad: 50M users, 25% market share
- Rocket: 10M users, 5% market share
- Upay, SureCash: <1M users each

**Payment Gateway Challenges:**
- **Approval time:** 28-35 days for API access
- **Sandbox limitations:** Limited test scenarios
- **Webhook reliability:** 95% delivery (need retry logic)
- **Reconciliation:** Manual daily reconciliation required
- **Limits:** Per transaction, daily, monthly limits

**Phase 1 Strategy:** Manual payment workflow to eliminate dependency on API approvals. Add live integrations in Phase 1.5 post-approval.

#### Regulatory Technology Requirements

**IDRA Digital Requirements:**
- E-submission portal for FCR, CARAMELS
- API for policy data sharing (future)
- Real-time transaction monitoring
- Digital signature for policy documents
- Audit trail with 7-year retention

**BFIU AML/CFT Requirements:**
- Transaction monitoring system
- Risk scoring engine
- STR/SAR filing system (online portal)
- Customer due diligence (CDD) records
- Sanctions screening

---

## 3. Overall Description

### 3.1 Product Perspective

The LabAid InsureTech Platform is a **greenfield digital insurance ecosystem** built from the ground up as a modern, cloud-native, microservices-based solution. It operates within the following system context:

**System Context Diagram:**

┌─────────────────────────────────────────────────────────────────────┐
│                        External Systems                              │
├─────────────────────────────────────────────────────────────────────┤
│  - bKash/Nagad/Rocket (MFS Payment)                                 │
│  - Hospital EHR Systems (Claims verification)                       │
│  - NID Verification API (KYC)                                        │
│  - IDRA Portal (Regulatory reporting)                                │
│  - BFIU Portal (AML/CFT reporting)                                   │
│  - SMS Gateway (Notifications)                                       │
│  - Email Service (Notifications)                                     │
│  - IoT Device APIs (Livestock trackers, health wearables)           │
└─────────────────────────────────────────────────────────────────────┘
                                   ▲
                                   │ Integration Layer
                                   │
┌─────────────────────────────────────────────────────────────────────┐
│               LabAid InsureTech Platform (Core System)               │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  ┌──────────────────┐  ┌──────────────────┐  ┌─────────────────┐   │
│  │  API Gateway     │  │  Message Bus     │  │  Auth Service   │   │
│  │  (Go)            │  │  (Kafka)         │  │  (Go - Existing)│   │
│  └──────────────────┘  └──────────────────┘  └─────────────────┘   │
│           │                     │                      │            │
│  ┌────────────────────────────────────────────────────────────┐    │
│  │                Backend Microservices                        │    │
│  │  - Policy Service (Go)                                      │    │
│  │  - Claims Service (Go)                                      │    │
│  │  - Payment Service (Go - Manual + MFS)                      │    │
│  │  - Insurance Engine (C# - Premium calculation)              │    │
│  │  - Notification Service (Python)                            │    │
│  │  - OCR/PDF Service (Python - Existing)                      │    │
│  │  - Partner Service (Go)                                     │    │
│  │  - Analytics Service (Python)                               │    │
│  └────────────────────────────────────────────────────────────┘    │
│           │                     │                      │            │
│  ┌────────────────────────────────────────────────────────────┐    │
│  │                 Data Layer                                  │    │
│  │  - PostgreSQL 17 (Transactional data - Existing)           │    │
│  │  - MongoDB (Unstructured data, logs)                        │    │
│  │  - S3/Blob Storage (Documents, images - Existing)           │    │
│  │  - Redis (Session cache, rate limiting)                     │    │
│  └────────────────────────────────────────────────────────────┘    │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
                                   │
                                   │ Client Layer
                                   ▼
┌─────────────────────────────────────────────────────────────────────┐
│                        User Interfaces                               │
├─────────────────────────────────────────────────────────────────────┤
│  - Customer Web Portal (React PWA)                                   │
│  - Customer Mobile App (React Native/Flutter - Android, iOS)        │
│  - Partner Portal (React Web)                                        │
│  - Admin Portal (React Web)                                          │
└─────────────────────────────────────────────────────────────────────┘
**Reusable Components (Existing Assets):**

The platform leverages proven, production-tested components from prior projects:

1. **Database Service (Go):** PostgreSQL 17 management with migration CLI (dbmanager), backup/restore, seeding
2. **S3 Storage Service (Go):** Document/image storage abstraction, access control
3. **Auth Service (Go):** OpenID/OAuth2, RBAC (Casbin), session management, 2FA, audit logging
4. **OCR/PDF Service (Python):** Document processing, NID digitization, policy certificate generation
5. **Mobile OTP Service (Node.js - Mamoon):** SMS OTP, one-tap login, JWT issuance
6. **bKash Integration Code (Node.js - Mamoon):** Payment gateway adapter (for Phase 1.5)

**Estimated Savings:** 755 development hours (~5 weeks)

**New Components to Build:**

1. **API Gateway (Go):** Service orchestration, rate limiting, authentication proxy
2. **Policy Service (Go):** Product catalog, policy CRUD, underwriting
3. **Claims Service (Go):** Claims submission, approval workflows, document management
4. **Manual Payment Service (Go):** Payment proof upload, admin verification queue, reconciliation
5. **Insurance Engine (C#):** Premium calculation, risk scoring, actuarial formulas
6. **Partner Service (Go):** Partner onboarding, KYB, tenant management
7. **Admin Portal (React):** Policy/claims management, payment verification, reporting
8. **Customer Web Portal (React):** Product discovery, purchase, claims, account management
9. **Android App (React Native/Flutter):** Native mobile experience (Phase 1.5)

### 3.2 Product Functions

**Primary Functions:**

1. **Customer Acquisition:**
   - Self-service registration (mobile OTP)
   - Digital KYC (photo + NID upload)
   - Social login integration (Google, Facebook)
   - Agent-assisted onboarding

2. **Product Discovery:**
   - Browse insurance products (category, price filters)
   - Compare products side-by-side
   - Product recommendations (Phase 2)
   - Educational content (insurance basics)

3. **Policy Purchase:**
   - Multi-step purchase flow
   - Premium calculation (instant quotes)
   - Digital underwriting (automated risk assessment)
   - **Manual payment** (Phase 1) → **Live MFS** (Phase 1.5)
   - Instant policy issuance (PDF certificate)

4. **Claims Management:**
   - Claims submission (form + documents)
   - Photo/video upload (incident evidence)
   - Real-time status tracking
   - Manual approval workflow (Phase 1) → **Semi-automated** (Phase 1.5) → **Fully automated** (Phase 2)
   - Claim settlement (MFS payout)

5. **Partner Integration:**
   - Hospital cashless settlement (EHR integration - Phase 2)
   - MFS payment processing
   - E-commerce embedded insurance (Phase 2)
   - IoT device integration (Phase 2/3)

6. **Administration:**
   - User management (RBAC)
   - Policy management (view, edit, cancel)
   - **Manual payment verification** (Phase 1 key function)
   - Claims approval (matrix-based routing)
   - Reporting & analytics
   - Compliance monitoring

7. **Notifications:**
   - SMS (Tier-1 priority - Phase 1)
   - Email (Phase 1.5)
   - Push notifications (Phase 1.5)
   - WhatsApp Business (Phase 2)

8. **Reporting & Compliance:**
   - IDRA CARAMELS reports (quarterly)
   - BFIU AML/CFT reports
   - Business intelligence dashboards (Phase 2)
   - Audit logs (7-year retention)

### 3.3 User Classes and Characteristics

**User Class Hierarchy:**

┌───────────────────────────────────────────────────────────────┐
│                      System Admin                              │
│  - Full system access                                          │
│  - User provisioning                                           │
│  - System configuration                                        │
│  - Audit log access                                            │
└───────────────────────────────────────────────────────────────┘
                            │
        ┌───────────────────┼───────────────────┐
        │                   │                   │
┌───────────────┐  ┌────────────────┐  ┌───────────────────┐
│ Business Admin│  │ Database Admin │  │ Repository Admin  │
│ - Approvals   │  │ - DB access    │  │ - Code deploy     │
│ - Policies    │  │ - Backups      │  │ - CI/CD           │
└───────────────┘  └────────────────┘  └───────────────────┘
        │
┌───────────────────────────────────────────────────────────────┐
│                      Focal Person                              │
│  - Partner onboarding (exclusive)                              │
│  - Partner verification                                        │
│  - Dispute resolution                                          │
│  - Partner ACL setup                                           │
└───────────────────────────────────────────────────────────────┘
        │
┌───────────────────────────────────────────────────────────────┐
│                      Partner Admin                             │
│  - Tenant IAM (within scope)                                   │
│  - Partner users management                                    │
│  - Integration configuration                                   │
└───────────────────────────────────────────────────────────────┘
        │
        ├── Hospital Staff (claims verification)
        ├── MFS Operations (payment reconciliation)
        └── E-commerce Admins (product embedding)
**Customer Classes:**

| User Class | Tech Literacy | Device | Network | Assistance Needed |
|------------|---------------|--------|---------|-------------------|
| **Type 1: Urban Literate** | High | Smartphone (6GB RAM) | 4G | None (self-service) |
| **Type 2: Semi-Urban** | Medium | Smartphone (4GB RAM) | 3G/4G | Occasional (family/agent) |
| **Type 3: Rural** | Low | Feature/Basic phone | 2G/3G | Frequent (agent + voice) |

**Internal Staff Classes:**

| Role | Count (Phase 1) | Primary Functions | Access Level |
|------|----------------|-------------------|--------------|
| **System Admin** | 1 | Infrastructure, security | Full |
| **Business Admin** | 2 | Policy approval, compliance | High |
| **Focal Person** | 3 | Partner management | Medium |
| **L1 Claims Admin** | 5 | Daily operations, small claims | Low |
| **L2 Claims Admin** | 2 | Medium claims, escalations | Medium |
| **Customer Support** | 10 | Helpdesk, issue resolution | Read + ticket creation |

### 3.4 Operating Environment

**Client-Side Environment:**

**Mobile Applications:**
- **Android:** Version 9.0 (API 28) minimum, 11+ recommended
- **iOS:** Version 13.0 minimum, 15+ recommended
- **RAM:** 4GB minimum, 6GB+ recommended
- **Storage:** 500MB free space for app + cache
- **Network:** 3G minimum (2 Mbps), 4G recommended

**Web Browsers:**
- Chrome 90+ (primary target)
- Firefox 88+
- Safari 14+ (iOS)
- Edge 90+
- **Not supported:** IE11, older browsers

**Server-Side Environment:**

**Cloud Infrastructure:**
- **Provider:** AWS (primary), Azure (backup/DR)
- **Region:** Singapore (primary), Mumbai (secondary)
- **Services:** EC2/ECS, RDS, S3, CloudFront, Lambda

**Backend Services:**
- **Go Services:** Go 1.21+, containerized (Docker)
- **C# Insurance Engine:** .NET 8, containerized
- **Python Services:** Python 3.11+, containerized
- **Node.js Services:** Node 20 LTS (Mamoon's services)

**Databases:**
- **PostgreSQL:** Version 17, managed (AWS RDS)
- **MongoDB:** Version 7.0, managed (Atlas)
- **Redis:** Version 7.2, managed (ElastiCache)

**Message Bus:**
- **Kafka:** Version 3.6, managed (MSK or self-hosted)

**Monitoring & Observability:**
- **Metrics:** Prometheus + Grafana
- **Logs:** ELK Stack or CloudWatch
- **APM:** OpenTelemetry + Jaeger
- **Uptime:** Pingdom or UptimeRobot

### 3.5 Design and Implementation Constraints

**Regulatory Constraints:**
1. **IDRA Compliance:** All policy data must be reportable in CARAMELS format
2. **BFIU AML/CFT:** Transaction monitoring mandatory for amounts >10,000 BDT
3. **Data Residency:** Customer data must be stored within legal jurisdiction (current: no strict local requirement, but planned)
4. **Audit Trail:** 7-year retention for all transactions and policy records
5. **Digital Signature:** Policy certificates must have valid digital signature

**Technical Constraints:**
1. **Network:** Must work on 3G minimum (2 Mbps)
2. **Offline Mode:** Critical functions must work offline (view policy, track claims)
3. **Device Support:** Must run on 4GB RAM Android devices
4. **App Size:** Initial download <10MB, total installed <100MB
5. **Legacy Integrations:** Hospital EHR systems may use SOAP/XML (not REST/JSON)

**Security Constraints:**
1. **Encryption:** AES-256 at rest, TLS 1.3 in transit
2. **Authentication:** 2FA mandatory for admin access
3. **Session:** Server-side sessions only (no client-side JWT storage for admins)
4. **PCI DSS:** If card payments added, must be PCI DSS Level 2 compliant
5. **Penetration Testing:** Annual mandatory before IDRA license renewal

**Business Constraints:**
1. **Budget:** Phase 1 additional cost capped at ,000 USD
2. **Timeline:** Phase 1 beta launch March 1, 2026 (non-negotiable)
3. **Team:** 12-person core team, phased availability (see TeamDetails_Availability.md)
4. **Partnership Dependencies:** bKash API approval 28-35 days (hence manual payment Phase 1)
5. **Product Approval:** IDRA product filing 30-45 days per product

**Architectural Constraints:**
1. **Microservices:** Loosely coupled, independently deployable services
2. **API-First:** All business logic accessible via REST/gRPC APIs
3. **Event-Driven:** Kafka for async communication, eventual consistency
4. **Idempotency:** All payment and policy operations must be idempotent
5. **Backward Compatibility:** APIs must maintain backward compatibility (versioning)

### 3.6 Assumptions and Dependencies

**Assumptions:**

1. **Network Availability:** Target users have intermittent 3G/4G internet access
2. **Device Availability:** Target users own smartphones capable of running modern mobile applications
3. **Digital Literacy:** Primary user segments (Type 1 & 2) can navigate mobile applications with minimal assistance
4. **Payment Infrastructure:** Bangladesh MFS ecosystem is operational and accessible to target users
5. **Regulatory Stability:** IDRA regulations remain stable during system development and deployment
6. **Hospital Cooperation:** Partner hospitals are willing to integrate systems for cashless claim settlement
7. **Identity Verification:** National ID verification services are available and reliable
8. **Third-party APIs:** Payment gateway, SMS gateway, and integration APIs are available and documented

**External Dependencies:**

| Dependency | Provider | Impact | Mitigation Strategy |
|------------|----------|--------|---------------------|
| **bKash Payment Gateway** | bKash Limited | High - Required for live payments | Phase 1 uses manual payment workflow to eliminate dependency |
| **Nagad Payment Gateway** | Bangladesh Post Office | Medium - Alternative payment method | Add in Phase 1.5 after bKash |
| **NID Verification API** | Election Commission Bangladesh | Medium - KYC compliance | Manual verification fallback process |
| **SMS Gateway Service** | Multiple providers available | Low - Critical notifications | Multiple vendor options, redundancy |
| **Hospital EHR Systems** | Hospital IT departments | Medium - Cashless settlement | Manual verification fallback |
| **IDRA Product Approvals** | Insurance Development & Regulatory Authority | High - Legal requirement to operate | File applications early, maintain regular communication |
| **Cloud Infrastructure** | AWS/Azure/GCP | Low - Hosting dependency | Multi-cloud strategy, local data center fallback |

**Internal Dependencies:**

| Dependency | Description | Impact | Notes |
|------------|-------------|--------|-------|
| **Reusable Components** | Existing production-tested services (Database, S3, Auth, OCR) | High - Accelerates development | Services must be production-ready with minimal adaptation |
| **Technical Expertise** | Domain knowledge in insurance, regulatory compliance | Medium - Learning curve | Training and documentation required |
| **Design System** | Consistent UI/UX across all platforms | Medium - User experience | LabAid brand guidelines as foundation |
| **Compliance Documentation** | IDRA and BFIU regulatory templates | High - Legal requirement | Legal and compliance team involvement |

---

---



### 4.5 Payment Processing (FG-005)

**Function Group Description:** Payment workflows including **manual payment verification (Phase 1 key feature)**, MFS integration, reconciliation, and refunds.

| FR ID | Requirement Description | Phase | Acceptance Criteria | Dependencies |
|-------|------------------------|-------|---------------------|--------------|
| **FR-044** | The system shall implement **manual payment workflow** for Phase 1: (1) Customer completes policy details, (2) System generates Payment Reference Number (PRN-YYYYMMDD-XXXXX), (3) Customer sees payment instructions (bank account/bKash/Nagad merchant number with PRN), (4) Customer makes payment offline, (5) Customer uploads payment proof (bank slip photo or MFS screenshot, max 5MB), (6) Admin receives notification in verification queue, (7) Admin verifies payment and approves/rejects within 1 hour SLA, (8) Policy activates on approval. | **M1** | • PRN generation unique<br>• Upload accepts JPEG/PNG/PDF<br>• Admin queue real-time<br>• Approval flow <1 hour avg<br>• SMS notification on status<br>• Retry allowed on rejection | Database service, S3 service, Notification service |
| **FR-045** | The system shall provide admin payment verification dashboard with: pending queue (sorted by submission time), payment proof viewer (zoom, download), customer details panel, transaction amount, payment method dropdown (Bank/bKash/Nagad/Cash), approve/reject buttons, rejection reason field (mandatory if reject). | **M1** | • Dashboard loads <2 seconds<br>• Image zoom functional<br>• Approve creates policy<br>• Reject sends SMS with reason<br>• Audit log complete | FR-044, Admin portal |tal |
| **FR-046** | The system shall enforce payment verification SLA: pending >2 hours = yellow alert to supervisor, pending >4 hours = red alert to Business Admin, auto-escalation email. Performance metric tracked. | **M1** | • Alerts trigger on time<br>• Email notifications sent<br>• Dashboard shows SLA status<br>• Monthly SLA report<br>• Gamification (optional) | FR-044, FR-045, Notification service |
| **FR-047** | The system shall integrate with bKash Payment Gateway API v1.2 (checkout, execute, query) for live payment processing. Supported: single payment, installment (if available), refund. Webhook for async status updates. | **M2** | • Sandbox testing complete<br>• Production approval obtained<br>• Checkout flow <10 seconds<br>• Webhook retry logic<br>• Daily reconciliation | bKash API access, Mamoon's existing code |
| **FR-048** | The system shall integrate with Nagad Payment API v2.0 for live payment processing. Same capabilities as bKash. | **M2** | • API integration complete<br>• Payment success rate >98%<br>• Error handling comprehensive<br>• Reconciliation automated<br>• Support contact established | Nagad API access |
| **FR-049** | The system shall integrate with Rocket Payment API (Dutch-Bangla Bank) for payment processing. | **S** | • API integration complete<br>• Similar to bKash/Nagad<br>• Reconciliation automated | Rocket API access |
| **FR-050** | The system shall support card payments (Visa, Mastercard) via payment gateway (SSLCommerz or similar) with PCI DSS compliance. | **S** | • Gateway integration complete<br>• 3D Secure mandatory<br>• PCI DSS audit passed<br>• Refund process tested | Payment gateway contract, PCI compliance |
| **FR-051** | The system shall support installment payment plans (EMI) for premiums >5,000 BDT. Options: 3, 6, 12 months. Interest rate configurable. Auto-debit on due date. | **S** | • Installment calculator accurate<br>• Due date reminders work<br>• Auto-debit success >95%<br>• Default handling process<br>• Policy suspension on missed payment | FR-047 or FR-048, Scheduler service |
| **FR-052** | The system shall implement auto-debit for policy renewals with customer consent. Consent collected during purchase, revocable anytime. Notification 7 days before debit. | **S** | • Consent flow clear<br>• Pre-notification sent<br>• Debit success rate >90%<br>• Failure handling graceful<br>• Opt-out easy | FR-047, FR-048, FR-051 |
| **FR-053** | The system shall perform daily payment reconciliation: match payment records with bank/MFS statements, flag discrepancies, generate reconciliation report, alert finance team on mismatches. | **M2** | • Reconciliation runs 2 AM daily<br>• Mismatch detection accurate<br>• Report detailed<br>• Alert email sent<br>• Manual resolution workflow | FR-047, FR-048, Scheduler |
| **FR-054** | The system shall process refunds for cancelled policies (within cooling period) or claim settlements. Refund via original payment method. TAT: 7 working days. | **M2** | • Refund initiation <24 hours<br>• Status tracking available<br>• SMS notification on completion<br>• Finance approval workflow<br>• Audit trail | FR-047, FR-048, Notification service |

**Manual Payment Workflow Details (Phase 1):**

**Payment Instructions Display Requirements:** System shall display payment reference number, amount, policy details, and instructions for each available payment method (bank transfer, bKash, Nagad) with clear emphasis on reference number usage.

**
**Admin Verification Interface Requirements:** System shall provide queue-based interface showing pending verifications with payment proof display, customer details, verification form, approve/reject actions, and SLA indicators.

**
**Function Group Summary:**
- **Total Requirements:** 11
- **Phase 1 (M1):** 3 requirements (**manual payment workflow** - critical for Phase 1)
- **Phase 1.5 (M2):** 4 requirements (bKash/Nagad live, reconciliation, refunds)
- **Phase 2 (S):** 4 requirements (Rocket, card, installments, auto-debit)
- **Phase 1 Focus:** Manual payment is THE KEY STRATEGY to eliminate bKash dependency

---

### 4.6 Claims Management (FG-006)

**Function Group Description:** Claims submission, document upload, approval workflows, fraud detection, settlement processing.

| FR ID | Requirement Description | Phase | Acceptance Criteria | Dependencies |
|-------|------------------------|-------|---------------------|--------------|
| **FR-055** | The system shall provide claims submission form with fields: policy number (auto-populated if logged in), incident date, incident type (dropdown), incident location, description (500 char max), claimed amount, supporting documents upload (multiple files, max 10MB total). | **M1** | • Form validation complete<br>• Document upload to S3<br>• Submission confirmation<br>• SMS notification sent<br>• Claim ID generated | Database service, S3 service, Notification service |
| **FR-056** | The system shall allow customers to track claim status in real-time: Submitted → Under Review → Approved/Rejected → Settlement Initiated → Paid. Status visible in app/web portal. | **M1** | • Status updates real-time<br>• Timeline visualization<br>• Estimated TAT shown<br>• SMS on status change<br>• Chat option (Phase 2) | FR-055, Notification service |
| **FR-057** | The system shall implement manual claims approval workflow for Phase 1: L1 Admin reviews claim, verifies documents, checks policy validity, routes to appropriate approver based on amount (see matrix), approver approves/rejects with notes, SMS notification sent to customer. | **M1** | • Workflow routing correct<br>• Document viewer functional<br>• Approval notes captured<br>• TAT tracking automated<br>• Audit log complete | FR-055, Admin portal, FR-058 |
| **FR-058** | The system shall enforce claims approval matrix based on claimed amount: <5,000 BDT = L1 Admin (TAT 24h), 5,000-20,000 BDT = L2 Admin (TAT 48h), 20,001-50,000 BDT = L3 Admin (TAT 72h), >50,000 BDT = Business Admin + Director (TAT 5 days). Auto-routing based on amount. | **M1** | • Matrix configurable in admin<br>• Routing automated<br>• Escalation on TAT breach<br>• Multiple approvers if needed<br>• Parallel approval support | FR-057, Database service |
| **FR-059** | The system shall implement automated approval for small claims (<5,000 BDT) meeting criteria: policy active, premium paid, within coverage, no fraud flags, customer KYC complete, first claim on policy. Auto-approve and initiate settlement. | **M2** | • Criteria validation automated<br>• Auto-approval rate >80%<br>• Manual review queue for edge cases<br>• Audit log detailed<br>• Override option for admins | FR-057, FR-058, FR-060 |
| **FR-060** | The system shall implement basic fraud detection with checks: duplicate claim (same incident), claim amount vs coverage limit, frequency check (>3 claims in 6 months = flag), geo-location mismatch (incident vs customer location), NID verification status. Flagged claims routed to fraud team. | **M2** | • Rules engine configurable<br>• Fraud detection accuracy >85%<br>• False positive rate <10%<br>• Fraud queue separate<br>• Investigation workflow | FR-055, FR-057, Database service |
| **FR-061** | The system shall implement advanced fraud detection using ML model trained on: claim patterns, device fingerprinting, social network analysis, image forensics (detect tampered documents), behavior anomaly detection. | **S** | • ML model deployed<br>• Prediction accuracy >90%<br>• Model retraining quarterly<br>• Explainable AI for decisions<br>• Human-in-the-loop review | FR-060, ML pipeline, Python analytics service |
| **FR-062** | The system shall integrate with hospital EHR (Electronic Health Records) for health insurance claims: fetch patient admission details, diagnosis, treatment, billing, discharge summary. Auto-populate claim form. | **S** | • EHR API integration complete<br>• Data mapping accurate<br>• Consent workflow clear<br>• HIPAA-equivalent privacy<br>• Fallback to manual | Hospital partnerships, SOAP/REST adapters |
| **FR-063** | The system shall support cashless claim settlement at partner hospitals: customer shows policy QR code at hospital, hospital submits claim directly via partner portal, LabAid settles with hospital (net 30 days), customer pays zero upfront (within coverage limit). | **S** | • QR code verification instant<br>• Hospital portal functional<br>• Real-time limit check<br>• Settlement workflow automated<br>• Reconciliation monthly | FR-062, Partner portal, Payment service |
| **FR-064** | The system shall settle approved claims within TAT: small claims (<5,000) = 24 hours, medium claims (5,000-20,000) = 48 hours, large claims (>20,000) = 5 working days. Payment via MFS to registered mobile number. | **M2** | • TAT compliance >95%<br>• Auto-settlement for small claims<br>• Payment success rate >98%<br>• Failed payment retry logic<br>• SMS confirmation | FR-057, FR-059, FR-047, FR-048 |
| **FR-065** | The system shall send SMS notifications on claim lifecycle events: submission acknowledgment, document request (if any), under review, approved (with amount and TAT), rejected (with reason and appeal process), settlement initiated, payment completed. | **M1** | • All events covered<br>• SMS delivery rate >98%<br>• Content clear and actionable<br>• Bengali + English templates<br>• Opt-out option | Notification service, FR-055 |

**Claims Approval Matrix (Configurable):**

| Claimed Amount (BDT) | Approval Level | Approver Role | Maximum TAT | Auto-Approval Eligible |
|---------------------|----------------|---------------|-------------|----------------------|
| < 5,000 | Level 1 | Any L1 Admin | 24 hours | Yes (Phase 1.5) |
| 5,000 - 20,000 | Level 2 | L2 Admin | 48 hours | No |
| 20,001 - 50,000 | Level 3 | L3 Admin | 72 hours | No |
| 50,001 - 100,000 | Executive | Business Admin | 5 days | No |
| > 100,000 | Board | Business Admin + Director | 7 days | No |

**Function Group Summary:**
- **Total Requirements:** 11
- **Phase 1 (M1):** 5 requirements (submission, tracking, manual approval, approval matrix, SMS notifications)
- **Phase 1.5 (M2):** 3 requirements (auto-approval small claims, basic fraud detection, settlement TAT)
- **Phase 2 (S):** 3 requirements (ML fraud detection, EHR integration, cashless settlement)
- **Phase 1 Focus:** Manual claims workflow with matrix-based routing, basic document upload

---

### 4.7 Partner Management (FG-007)

**Function Group Description:** Partner onboarding, KYB verification, tenant isolation, portal access, integration management.

| FR ID | Requirement Description | Phase | Acceptance Criteria | Dependencies |
|-------|------------------------|-------|---------------------|--------------|
| **FR-066** | The system shall provide Focal Person exclusive access to partner onboarding workflow: (1) Verify MOU signed, (2) Collect partner details (company name, trade license, TIN, bank details, contact person), (3) Perform KYB (Know Your Business) verification, (4) Create partner tenant, (5) Setup ACL and data isolation, (6) Assign Partner Admin role with temporary password, (7) Send welcome email with portal access instructions. | **M1** | • Focal Person dashboard functional<br>• MOU upload required<br>• KYB checklist complete<br>• Tenant created isolated<br>• Email notification sent<br>• Audit log detailed | Existing Auth service RBAC/ACL, Database service, Notification service |
| **FR-067** | The system shall enforce multi-tenant data isolation: partner users can only access their tenant's data (policies, claims, customers within their scope), cross-tenant queries return 403 Forbidden, database queries filtered by tenant_id, S3 paths include tenant prefix. | **M1** | • Tenant ID in all queries<br>• Unit tests for isolation<br>• Penetration test passed<br>• Performance acceptable<br>• Audit log per tenant | FR-066, Database service, S3 service |
| **FR-068** | The system shall provide Partner Admin portal with capabilities: manage partner users (create/edit/disable), view partner-specific policies, view partner-specific claims, configure partner integration settings (webhook URLs, API keys), view commission reports (Phase 1.5), view performance metrics. | **M2** | • Portal responsive<br>• User management functional<br>• Data filtered by tenant<br>• API key generation secure<br>• Commission calc accurate | FR-066, FR-067, Partner portal (new) |
| **FR-069** | The system shall support hospital partner integration: webhook for cashless claim requests, API to query policy validity, API to submit claims, settlement report download (monthly). | **S** | • Webhook reliable<br>• API documented (OpenAPI)<br>• Auth via API key<br>• Rate limiting enforced<br>• Support SLA defined | FR-063, API Gateway |
| **FR-070** | The system shall support MFS partner integration: webhook for payment status updates, API to query payment reconciliation, settlement report download. | **M2** | • Webhook retry logic<br>• Reconciliation API accurate<br>• Report format agreed<br>• Error handling comprehensive | FR-047, FR-048, FR-053 |
| **FR-071** | The system shall support e-commerce partner integration: embedded insurance widget (iframe/SDK), API to create policy, webhook for policy status, commission calculation. | **S** | • Widget responsive<br>• SDK documented<br>• API rate limits set<br>• Commission automated<br>• Co-branding support | FR-066, FR-068 |
| **FR-072** | The system shall calculate and track partner commissions: percentage-based (configurable per partner), generated on policy activation, accumulated monthly, payable via bank transfer, commission report downloadable by partner. | **M2** | • Commission calc accurate<br>• Monthly statement generated<br>• Payment workflow defined<br>• Tax deduction support<br>• Audit trail | FR-066, Database service |

**Partner Types (Phase Roadmap):**

1. **Hospitals (Phase 2):**
   - Purpose: Cashless claim settlement, EHR integration
   - Count: 10 LabAid hospitals + 50 partner hospitals (by end Phase 2)
   - Integration: API + Partner Portal

2. **MFS Providers (Phase 1.5):**
   - Purpose: Payment processing
   - Count: bKash, Nagad (Phase 1.5), Rocket (Phase 2)
   - Integration: Payment Gateway API

3. **E-commerce Platforms (Phase 2):**
   - Purpose: Embedded insurance at checkout
   - Count: 5 major platforms (Daraz, Evaly equivalents)
   - Integration: Widget/SDK

4. **Corporate Clients (Phase 2):**
   - Purpose: Group insurance for employees
   - Count: 20 companies (target)
   - Integration: Bulk upload API, HR system integration

5. **IoT Device Vendors (Phase 2/3):**
   - Purpose: Usage-based insurance data
   - Count: 2-3 vendors (livestock trackers, health wearables)
   - Integration: IoT API, MQTT/HTTP webhooks

**Function Group Summary:**
- **Total Requirements:** 7
- **Phase 1 (M1):** 2 requirements (Focal Person onboarding, multi-tenant isolation)
- **Phase 1.5 (M2):** 3 requirements (Partner Admin portal, MFS integration, commission tracking)
- **Phase 2 (S):** 2 requirements (hospital integration, e-commerce integration)
- **Phase 1 Focus:** Basic partner onboarding by Focal Person, strict tenant isolation

---




### 4.11 Database & Storage (FG-011)

**Function Group Description:** Data persistence, storage management, backup, recovery, and data lifecycle requirements.

| FR ID | Requirement Description | Phase | Acceptance Criteria | Dependencies |
|-------|------------------------|-------|---------------------|--------------|
| **FR-099** | The system shall use relational database (PostgreSQL 17 or equivalent) for transactional data storage including: users, policies, claims, payments, audit logs, system configuration. Database shall support ACID transactions. | **M1** | • ACID compliance verified<br>• Referential integrity enforced<br>• Transaction isolation levels configurable<br>• Query performance <100ms (95th percentile)<br>• Connection pooling implemented | None |
| **FR-100** | The system shall use object storage (S3-compatible) for unstructured data: user photos (KYC), NID documents, policy certificates (PDF), claim documents, payment proofs, system backups. Storage shall support encryption at rest. | **M1** | • Encryption at rest enabled (AES-256)<br>• Access control lists functional<br>• Versioning enabled<br>• Lifecycle policies configured<br>• CDN integration for public assets | None |
| **FR-101** | The system shall implement automated database migration system with version control, rollback capability, schema change tracking, and zero-downtime deployment support. | **M1** | • Migration scripts versioned in Git<br>• Rollback tested for all migrations<br>• Schema changes logged<br>• Zero-downtime deployment verified<br>• Migration status dashboard | FR-099 |
| **FR-102** | The system shall perform automated daily database backups with 30-day retention for operational data and 7-year retention for compliance data (policies, claims, financial records). Backups encrypted and stored in geographically separate location. | **M1** | • Backup completion within 4-hour window<br>• Backup verification automated<br>• Encryption verified<br>• Geo-redundancy confirmed<br>• Restore procedure tested monthly | FR-099, FR-100 |
| **FR-103** | The system shall use NoSQL database (MongoDB or DynamoDB) for: product catalog metadata, analytics events, session data, notification logs, IoT device data. | **S** | • Schema flexibility verified<br>• Query performance acceptable<br>• Replication configured<br>• Backup strategy defined<br>• Data sync with relational DB | FR-099 |
| **FR-104** | The system shall implement database read replicas for reporting and analytics queries to avoid impacting transactional workload. Replication lag target: <5 minutes. | **S** | • Read replica configured<br>• Replication lag monitored<br>• Failover tested<br>• Load balancing implemented<br>• Query routing automated | FR-099, FR-103 |
| **FR-105** | The system shall implement data archival policy: archive inactive policies (>2 years old) to cold storage, archive processed claims (>1 year), maintain audit logs for 7 years (regulatory requirement), purge anonymized test data after 90 days. | **M2** | • Archival jobs scheduled<br>• Retrieval process tested<br>• Compliance verified<br>• Performance impact minimal<br>• Cost optimization achieved | FR-099, FR-100, FR-102 |
| **FR-106** | The system shall maintain data dictionary documenting all database tables, columns, data types, constraints, relationships, and business rules. Documentation auto-generated from schema annotations. | **M1** | • Documentation comprehensive<br>• Auto-generation functional<br>• Version controlled<br>• Searchable format<br>• Developer accessible | FR-099 |

**Database Schema - Core Entities:**

- **users:** Customer and stakeholder identity and authentication data
- **policies:** Insurance policy master data and status
- **policy_transactions:** Premium payments and adjustments
- **claims:** Claim submissions and processing workflow
- **claim_documents:** Supporting documentation for claims
- **payments:** Payment records and reconciliation data
- **payment_proofs:** Manual payment verification artifacts (Phase 1)
- **products:** Insurance product catalog and configuration
- **partners:** Partner organizations and integration settings
- **audit_logs:** Immutable audit trail for all system actions
- **notifications:** Notification history and delivery status
- **sessions:** User session management
- **system_config:** Application configuration and feature flags

**Function Group Summary:**
- **Total Requirements:** 8
- **Phase 1 (M1):** 5 requirements (relational DB, object storage, migrations, backups, data dictionary)
- **Phase 1.5 (M2):** 1 requirement (data archival)
- **Phase 2 (S):** 2 requirements (NoSQL, read replicas)

---

### 4.12 Integration Services (FG-012)

**Function Group Description:** External API integrations, webhooks, message queue, and inter-service communication.

| FR ID | Requirement Description | Phase | Acceptance Criteria | Dependencies |
|-------|------------------------|-------|---------------------|--------------|
| **FR-107** | The system shall provide RESTful API Gateway with capabilities: authentication (JWT, API key), rate limiting (per client, per endpoint), request validation, response transformation, error handling, API versioning (v1, v2 etc), request/response logging. | **M1** | • OpenAPI 3.0 specification complete<br>• Authentication enforced<br>• Rate limits configurable<br>• Logging comprehensive<br>• Documentation published | None |
| **FR-108** | The system shall implement gRPC APIs for internal service-to-service communication with capabilities: protobuf schemas, bi-directional streaming, load balancing, circuit breaker pattern, timeout configuration, retry logic with exponential backoff. | **M2** | • Protobuf schemas versioned<br>• Streaming tested<br>• Circuit breaker functional<br>• Performance benchmarked<br>• Service mesh integration (optional) | FR-107 |
| **FR-109** | The system shall integrate with SMS gateway provider API supporting: single SMS send, bulk SMS send (up to 10,000 recipients), delivery status callback, message templates, Unicode/Bengali language support, cost tracking. | **M1** | • SMS delivery rate >98%<br>• Callback handling reliable<br>• Unicode support verified<br>• Cost per SMS tracked<br>• Vendor redundancy configured | None |
| **FR-110** | The system shall integrate with email service provider (AWS SES, SendGrid, or equivalent) supporting: transactional email, HTML templates, attachments (up to 5MB), bounce/complaint handling, DMARC/SPF/DKIM authentication, unsubscribe management. | **M2** | • Email delivery rate >95%<br>• Templates rendering correctly<br>• Attachments working<br>• Bounce handling automated<br>• Authentication configured | None |
| **FR-111** | The system shall integrate with National ID (NID) verification API for KYC validation: submit NID number + DOB, receive verification status, retrieve NID holder name and photo (if available), log all verification attempts for audit. | **M2** | • API integration complete<br>• Response time <30 seconds<br>• Success rate >95%<br>• Fallback to manual verification<br>• Audit trail maintained | External API contract |
| **FR-112** | The system shall integrate with bKash Payment Gateway API including: tokenization for stored credentials, checkout flow (web, app), payment status query, refund processing, webhook for async notifications, daily reconciliation file download. | **M2** | • Checkout success rate >98%<br>• Webhook reliability >99%<br>• Reconciliation automated<br>• Security audit passed<br>• Error handling comprehensive | bKash API approval |
| **FR-113** | The system shall integrate with Nagad Payment API with similar capabilities as bKash integration. | **M2** | • Same criteria as FR-112<br>• Vendor-specific features supported | Nagad API approval |
| **FR-114** | The system shall integrate with hospital EHR systems for health claim verification: query patient admission records, retrieve diagnosis and treatment details, fetch billing information, verify discharge summary. Integration methods: REST API (preferred), SOAP API, HL7 FHIR (if available). | **S** | • Data mapping accurate<br>• Patient consent verified<br>• Privacy compliant<br>• Fallback to manual verification<br>• Multiple EHR vendors supported | Hospital partnerships |
| **FR-115** | The system shall provide webhook infrastructure for partner integrations: register webhook URLs, configure authentication (HMAC signature, API key), define retry policy, maintain delivery logs, provide webhook testing sandbox. | **M2** | • Registration API functional<br>• Signature verification working<br>• Retry logic tested<br>• Delivery logs queryable<br>• Sandbox environment available | FR-107 |
| **FR-116** | The system shall implement message queue (Kafka or RabbitMQ) for asynchronous processing: event streaming, notification orchestration, audit log collection, payment reconciliation jobs, report generation queue. | **M1** | • Message delivery guaranteed<br>• Consumer lag monitored<br>• Dead letter queue configured<br>• Performance acceptable<br>• Topic organization logical | None |

**API Design Principles:**

1. **RESTful Conventions:** Resources identified by URIs, HTTP methods semantic (GET, POST, PUT, DELETE), stateless communication
2. **Versioning:** URL-based versioning (api.example.com/v1/policies)
3. **Error Handling:** Standard HTTP status codes, detailed error responses with error codes and messages
4. **Pagination:** Cursor-based pagination for large result sets (limit, offset, cursor)
5. **Filtering & Sorting:** Query parameters for filtering (status=active) and sorting (sort=created_at:desc)
6. **Security:** HTTPS mandatory, authentication on all endpoints, authorization checks per resource
7. **Documentation:** OpenAPI/Swagger specification, interactive API explorer, code samples

**Function Group Summary:**
- **Total Requirements:** 10
- **Phase 1 (M1):** 3 requirements (API Gateway, SMS, message queue)
- **Phase 1.5 (M2):** 5 requirements (gRPC, email, NID verification, bKash/Nagad, webhooks)
- **Phase 2 (S):** 1 requirement (EHR integration)

---

### 4.13 IoT & Device Integration (FG-013)

**Function Group Description:** Integration with Internet of Things devices for usage-based insurance and real-time monitoring.

| FR ID | Requirement Description | Phase | Acceptance Criteria | Dependencies |
|-------|------------------------|-------|---------------------|--------------|
| **FR-117** | The system shall integrate with livestock tracking IoT devices (GPS collars, RFID tags) for cattle insurance: receive location updates (hourly), detect movement patterns, alert on geo-fence breach, monitor health metrics (if available), correlate with claims data. | **S** | • Device registration workflow<br>• Data ingestion reliable<br>• Alerts timely<br>• Dashboard visualization<br>• Battery level monitoring | IoT vendor partnerships |
| **FR-118** | The system shall support health wearable device integration (fitness trackers, smartwatches) for health insurance: sync activity data (steps, heart rate, sleep), calculate wellness score, offer premium discounts based on activity, privacy-preserving data aggregation. | **F** | • OAuth device authorization<br>• Data sync automated<br>• Wellness algorithm verified<br>• Privacy compliant<br>• User consent clear | Health wearable APIs |
| **FR-119** | The system shall provide IoT device management portal: register devices, map devices to policies, view device status, configure alert thresholds, download device data (CSV export), deactivate lost devices. | **S** | • Portal functional<br>• Bulk operations supported<br>• Real-time status updates<br>• Export working<br>• Role-based access | FR-117, FR-118 |
| **FR-120** | The system shall implement IoT data ingestion pipeline: receive MQTT messages from devices, validate and transform data, store time-series data, trigger alerts based on rules, provide query API for analytics. | **S** | • MQTT broker configured<br>• Message validation working<br>• Time-series DB performance<br>• Alert rules flexible<br>• Query API responsive | FR-117, Message broker |
| **FR-121** | The system shall calculate usage-based insurance premiums using IoT data: analyze historical usage patterns, apply actuarial models, adjust premiums at renewal, provide transparency to customers on calculation methodology. | **F** | • Calculation accuracy verified<br>• Actuarial approval obtained<br>• Customer communication clear<br>• Regulatory compliance<br>• Dispute resolution process | FR-117, FR-118, FR-120 |

**IoT Use Cases (Phased Approach):**

**Phase 2 (S) - Livestock Insurance Pilot:**
- Device Type: GPS collars with accelerometer
- Target: 1,000 cattle in 100 farms (pilot)
- Data: Location every hour, movement patterns, geo-fence alerts
- Business Value: Theft prevention, mortality verification, grazing pattern analysis

**Phase 3 (F) - Health Insurance UBI:**
- Device Type: Fitness trackers, smartwatches (user-owned)
- Integration: Apple Health, Google Fit, Samsung Health
- Data: Daily steps, active minutes, heart rate trends (aggregated)
- Business Value: Wellness incentives, risk assessment, customer engagement

**Function Group Summary:**
- **Total Requirements:** 5
- **Phase 2 (S):** 3 requirements (livestock IoT, device portal, data pipeline)
- **Phase 3 (F):** 2 requirements (health wearables, UBI premium calculation)
- **Phase 1:** No IoT features (deferred to validate core platform first)

---

### 4.14 Voice-Assisted Workflows (FG-014)

**Function Group Description:** Voice recognition and text-to-speech for Type 3 rural customers with limited digital literacy.

| FR ID | Requirement Description | Phase | Acceptance Criteria | Dependencies |
|-------|------------------------|-------|---------------------|--------------|
| **FR-122** | The system shall provide voice-assisted product browsing: customer speaks query in Bengali, system transcribes and searches products, system reads out product names and basic details, customer can ask follow-up questions, final selection via touch or voice. | **F** | • Speech recognition accuracy >90%<br>• Bengali language support<br>• Natural language understanding<br>• Response time <3 seconds<br>• Noise cancellation effective | ML speech models |
| **FR-123** | The system shall provide voice-guided policy purchase flow: system prompts for required information via voice, customer responds verbally, system confirms inputs, critical fields (payment amount) require visual confirmation, final approval via touch. | **F** | • Guided workflow complete<br>• Error handling graceful<br>• Fallback to visual mode<br>• Security for sensitive data<br>• Session timeout appropriate | FR-122 |
| **FR-124** | The system shall provide voice-based claims submission: customer describes incident verbally, system captures voice recording and transcribes, system prompts for required documents, customer uploads photos via camera, system reads back claim summary for confirmation. | **F** | • Recording quality acceptable<br>• Transcription accuracy >85%<br>• Bengali language support<br>• Audio stored securely<br>• Text searchable | FR-122, FR-123 |
| **FR-125** | The system shall provide USSD fallback for feature phone users (no smartphone): dial shortcode to access menu, navigate via numeric keypad, query policy status, submit basic claim notification (detailed submission requires agent assistance). | **F** | • USSD menu functional<br>• Mobile operator integration<br>• Session management working<br>• Response within 5 seconds<br>• Coverage all major operators | Telecom operator partnerships |

**Voice Assistant Technology Stack (Phase 3):**
- **Speech-to-Text:** Google Cloud Speech-to-Text or Azure Speech Services (Bengali language pack)
- **Text-to-Speech:** Natural-sounding Bengali voice synthesis
- **NLU (Natural Language Understanding):** Intent recognition, entity extraction
- **Dialog Management:** Multi-turn conversation handling
- **Fallback Strategy:** Visual mode always available, agent handoff option

**Function Group Summary:**
- **Total Requirements:** 4
- **Phase 3 (F):** All 4 requirements (deferred until core platform proven and Type 3 user adoption validated)

---

### 4.15 WebRTC Communication (FG-015)

**Function Group Description:** Real-time video and audio communication for claims verification and customer support.

| FR ID | Requirement Description | Phase | Acceptance Criteria | Dependencies |
|-------|------------------------|-------|---------------------|--------------|
| **FR-126** | The system shall provide video claim verification: customer initiates video call with claims adjuster, adjuster views incident site in real-time, adjuster captures screenshots during call, call recorded for audit purposes, customer consent obtained before recording. | **F** | • Video quality acceptable on 3G<br>• Recording storage secure<br>• Consent workflow clear<br>• Call duration tracked<br>• Bandwidth adaptive | WebRTC infrastructure |
| **FR-127** | The system shall provide customer support video chat: customer requests video assistance from help center, support staff responds via web portal, screen sharing capability for guidance, chat transcript saved, feedback collected post-call. | **F** | • Queue management working<br>• Screen sharing functional<br>• Mobile + web support<br>• Call quality monitoring<br>• Support metrics tracked | FR-126 |
| **FR-128** | The system shall implement WebRTC infrastructure: TURN/STUN servers for NAT traversal, signaling server for call setup, media server for recording, bandwidth detection and adaptation, network quality indicators visible to users. | **F** | • Connection success rate >95%<br>• Latency <300ms<br>• NAT traversal working<br>• Recording quality acceptable<br>• Cost optimized | Cloud infrastructure |

**Function Group Summary:**
- **Total Requirements:** 3
- **Phase 3 (F):** All 3 requirements (deferred as advanced feature requiring infrastructure investment)

---

### 4.16 Mobile Applications (FG-016)

**Function Group Description:** Native mobile applications for Android and iOS platforms.

| FR ID | Requirement Description | Phase | Acceptance Criteria | Dependencies |
|-------|------------------------|-------|---------------------|--------------|
| **FR-129** | The system shall provide Android mobile application with features: user registration and login, product browsing, policy purchase, payment processing, claims submission, policy document viewer, push notifications, offline mode for viewing policy details. | **M2** | • App size <10MB initial download<br>• Minimum Android 9.0 (API 28)<br>• Published on Google Play Store<br>• 4GB RAM device support<br>• Offline functionality working | All backend services |
| **FR-130** | The system shall provide iOS mobile application with same features as Android application. | **S** | • App size <10MB<br>• Minimum iOS 13.0<br>• Published on Apple App Store<br>• iPhone 8 and above support<br>• Feature parity with Android | FR-129 |
| **FR-131** | The system shall implement mobile app security: biometric authentication (fingerprint, face), secure storage for sensitive data, certificate pinning for API calls, jailbreak/root detection, session timeout on background. | **M2** | • Biometric setup optional<br>• Secure storage encrypted<br>• Certificate pinning tested<br>• Detection alerts user<br>• Security audit passed | FR-129 |
| **FR-132** | The system shall implement mobile app offline capabilities: cache policy documents, cache recent transactions, queue actions when offline (sync when online), offline-friendly UI (no broken images). | **M2** | • Offline mode functional<br>• Sync conflict resolution<br>• Storage management smart<br>• User informed of offline status<br>• Data consistency maintained | FR-129 |
| **FR-133** | The system shall optimize mobile app performance for low-end devices: lazy loading, image compression, pagination, memory management, battery optimization, dark mode support. | **M2** | • App responsive on 4GB RAM device<br>• Battery drain acceptable<br>• ANR (Application Not Responding) rate <0.1%<br>• Crash rate <1%<br>• Dark mode complete | FR-129 |

**Mobile App Technology Considerations:**
- **Framework:** React Native or Flutter (cross-platform) vs Native (platform-specific)
- **State Management:** Redux, MobX, or Context API
- **Networking:** Axios or Fetch with retry and timeout handling
- **Storage:** SQLite or Realm for offline data
- **Analytics:** Firebase Analytics or equivalent
- **Crash Reporting:** Crashlytics or Sentry

**Function Group Summary:**
- **Total Requirements:** 5
- **Phase 1.5 (M2):** 4 requirements (Android app with security, offline, performance optimization)
- **Phase 2 (S):** 1 requirement (iOS app)

---

### 4.17 Web Applications (FG-017)

**Function Group Description:** Responsive web portals for customers, partners, and administrators.

| FR ID | Requirement Description | Phase | Acceptance Criteria | Dependencies |
|-------|------------------------|-------|---------------------|--------------|
| **FR-134** | The system shall provide customer web portal with features: responsive design (mobile, tablet, desktop), progressive web app (PWA) installable on home screen, offline mode for viewing policies, fast loading (<3 seconds on 3G), accessibility compliant (WCAG 2.1 AA). | **M1** | • Mobile-first design verified<br>• PWA installation working<br>• Lighthouse score >90<br>• Accessibility audit passed<br>• Cross-browser compatible | All backend services |
| **FR-135** | The system shall provide partner web portal with tenant-specific branding, user management, policy view (scoped to partner), claims view (scoped to partner), integration configuration, reporting and analytics. | **M2** | • Multi-tenant isolation verified<br>• Branding customization working<br>• Data filtering accurate<br>• Performance acceptable<br>• Role-based access enforced | FR-066, FR-067 |
| **FR-136** | The system shall provide admin web portal with all administrative functions: dashboard, user management, policy management, payment verification, claims approval, product management, reporting, system configuration. | **M1** | • All features functional<br>• Performance optimized<br>• Responsive design<br>• Role-based UI rendering<br>• Audit logging integrated | All backend services |
| **FR-137** | The system shall implement web portal security: HTTPS only, HSTS headers, CSP (Content Security Policy), CSRF protection, XSS prevention, secure session cookies (httpOnly, secure flags), automatic logout after 30 minutes inactivity. | **M1** | • Security headers verified<br>• CSRF tokens working<br>• XSS tests passed<br>• Session management secure<br>• Security scan passed | None |
| **FR-138** | The system shall optimize web performance: code splitting, lazy loading, image optimization (WebP, lazy load), CDN for static assets, browser caching, compression (gzip/brotli), critical CSS inline. | **M1** | • Page load <3s on 3G<br>• Time to Interactive <5s<br>• First Contentful Paint <2s<br>• Lighthouse performance >90<br>• Bundle size optimized | None |

**Web Technology Stack:**
- **Frontend Framework:** React, Vue, or Angular
- **Styling:** Tailwind CSS, Material-UI, or custom design system
- **State Management:** Context API, Redux, or Zustand
- **Build Tool:** Vite or Webpack
- **Testing:** Jest, React Testing Library, Cypress
- **Hosting:** CDN (CloudFront, Cloudflare) + origin server

**Function Group Summary:**
- **Total Requirements:** 5
- **Phase 1 (M1):** 3 requirements (customer portal, admin portal, security, performance)
- **Phase 1.5 (M2):** 1 requirement (partner portal)
- **Phase 1 Focus:** Customer PWA and admin portal as minimum viable web presence

---

## Functional Requirements Summary

**Total Functional Requirements:** 138 (across 17 function groups)

**Phase Distribution:**
- **Phase 1 (M1):** 60 requirements - Core MVP for beta launch
- **Phase 1.5 (M2):** 45 requirements - Complete baseline for public launch
- **Phase 2 (S):** 21 requirements - Scaling and automation
- **Phase 3 (D/C/F):** 12 requirements - Innovation and advanced features

**Critical Phase 1 Innovations:**
1. **Manual Payment Workflow** - Eliminates external dependency on bKash API approval
2. **Reusable Component Strategy** - Leverages existing production services for 755-hour savings
3. **SMS-First Notifications** - Focuses on most reliable channel for Phase 1
4. **Web PWA Before Native Apps** - Delivers mobile experience without app store delays
5. **Manual Claims Approval** - Human oversight before automation (risk management)

---


## 5. External Interface Requirements

### 5.1 User Interfaces

**General UI Requirements:**

| Requirement | Specification | Rationale |
|-------------|---------------|-----------|
| **Responsive Design** | Support 320px (mobile) to 2560px (desktop) viewports | Multi-device accessibility |
| **Language Support** | Bengali (primary), English (secondary) | Target market requirement |
| **Accessibility** | WCAG 2.1 Level AA compliance | Inclusive design for disabilities |
| **Loading States** | Skeleton screens, progress indicators for operations >2 seconds | User experience feedback |
| **Error Messages** | Clear, actionable error messages in user's language | Error recovery guidance |
| **Form Validation** | Real-time validation with inline error messages | Reduce submission errors |
| **Mobile Optimization** | Touch-friendly controls (min 44x44px), gesture support | Mobile usability |
| **Network Resilience** | Graceful degradation on slow networks (3G) | Bangladesh network reality |

**User Interface Components:**

1. **Customer Mobile/Web Interface:**
   - Home/Dashboard: Quick actions, policy status, recent transactions
   - Product Catalog: Grid/list view, filters, comparison
   - Purchase Flow: Multi-step wizard with progress indicator
   - Payment Interface: Instructions display, proof upload
   - Policy Viewer: PDF viewer, download button, share options
   - Claims Interface: Form, document upload, status tracking
   - Profile Management: KYC status, edit profile, preferences

2. **Admin Portal Interface:**
   - Dashboard: Key metrics, alerts, quick actions
   - Payment Verification Queue: Image viewer, approve/reject workflow
   - Claims Management: Queue management, document viewer, approval workflow
   - User Management: Search, view, edit, suspend capabilities
   - Reporting: Report parameters, preview, download

3. **Partner Portal Interface:**
   - Dashboard: Partner-specific metrics
   - User Management: Tenant-scoped user administration
   - Policy/Claims View: Filtered by partner association
   - Integration Settings: API key management, webhook configuration

### 5.2 Hardware Interfaces

**Minimum Device Specifications:**

| Platform | Minimum Requirements | Recommended |
|----------|---------------------|-------------|
| **Android Mobile** | Android 9.0, 4GB RAM, 64GB storage, 5.5" screen | Android 11+, 6GB RAM, 128GB storage |
| **iOS Mobile** | iOS 13.0, 64GB storage, iPhone 8 or newer | iOS 15+, iPhone 11 or newer |
| **Web Browser** | Chrome 90+, Firefox 88+, Safari 14+, 2GB RAM | Latest versions, 4GB+ RAM |
| **Admin Workstation** | Windows 10/macOS 10.15, 8GB RAM, 1920x1080 display | Windows 11/macOS 12+, 16GB RAM, 2K display |

**Peripheral Requirements:**
- **Camera:** Minimum 5MP for document capture, 8MP+ recommended
- **GPS:** Required for geo-tagging claims (livestock insurance)
- **Biometric Sensor:** Fingerprint or Face ID for mobile app security (optional)
- **Printer:** For partner/admin printing of documents (optional)

### 5.3 Software Interfaces

**External System Integrations:**

| System | Interface Type | Data Format | Authentication | Purpose |
|--------|---------------|-------------|----------------|---------|
| **bKash Payment Gateway** | REST API (HTTPS) | JSON | API Key + OAuth2 | Payment processing |
| **Nagad Payment API** | REST API (HTTPS) | JSON | API Key + HMAC | Payment processing |
| **NID Verification (OCCA)** | REST API (HTTPS) | JSON | API Key | KYC compliance |
| **SMS Gateway** | REST API (HTTPS) | JSON | API Key | Notifications |
| **Email Service (SES/SendGrid)** | REST API (HTTPS) | JSON | API Key | Notifications |
| **Hospital EHR Systems** | REST/SOAP (HTTPS) | JSON/XML | OAuth2/API Key | Claim verification |
| **IDRA Portal** | Web Form Submission | Excel/PDF | Username/Password | Regulatory reporting |
| **BFIU Portal** | Web Form Submission | Excel/PDF | Username/Password | AML/CFT reporting |
| **IoT Devices (Phase 2)** | MQTT/HTTP | JSON | Device Certificate | Livestock tracking |

**Internal Service Communication:**

| Service | Interface | Protocol | Data Format | Purpose |
|---------|-----------|----------|-------------|---------|
| **API Gateway ↔ Microservices** | REST/gRPC | HTTPS/HTTP2 | JSON/Protobuf | Service orchestration |
| **Services ↔ Database** | Database Driver | TCP | SQL | Data persistence |
| **Services ↔ Message Queue** | Kafka Client | TCP | Binary | Event streaming |
| **Services ↔ Object Storage** | S3 API | HTTPS | Binary | File storage |
| **Services ↔ Cache** | Redis Protocol | TCP | Binary | Session/cache |

### 5.4 Communication Interfaces

**Network Protocols:**
- **HTTPS (TLS 1.3):** All client-server communication
- **HTTP/2:** gRPC internal service communication
- **WebSocket (WSS):** Real-time notifications (Phase 2)
- **MQTT (TLS):** IoT device communication (Phase 2)

**Data Exchange Formats:**
- **JSON:** Primary API data format (human-readable, widely supported)
- **Protocol Buffers:** Internal gRPC communication (performance-optimized)
- **XML:** Legacy system integration (hospital EHR systems)
- **CSV/Excel:** Bulk data export, regulatory reporting

**API Standards:**
- **RESTful Design:** Resource-based URLs, HTTP methods semantic
- **OpenAPI 3.0:** API specification and documentation
- **OAuth 2.0:** Third-party authentication
- **JWT (RFC 7519):** Stateless authentication tokens
- **Webhook:** Event notification to partners (HTTP POST with HMAC signature)

---

## 6. Non-Functional Requirements

### 6.1 Performance Requirements

| Metric | Requirement | Measurement Method | Priority |
|--------|-------------|-------------------|----------|
| **API Response Time** | 95th percentile <500ms for read operations, <1s for write operations | Application Performance Monitoring (APM) | M1 |
| **Page Load Time** | <3 seconds on 3G network (2 Mbps) | Lighthouse, WebPageTest | M1 |
| **Database Query** | 95th percentile <100ms for transactional queries | Database query logs | M1 |
| **Concurrent Users** | Support 10,000 concurrent users (Phase 1), 100,000 (Phase 2) | Load testing (JMeter, Gatling) | M1/S |
| **Peak Transactions** | 100 transactions/second (Phase 1), 1,000 TPS (Phase 2) | Load testing | M1/S |
| **Mobile App Startup** | <2 seconds cold start time | Firebase Performance | M2 |
| **Notification Delivery** | SMS within 30 seconds, Email within 5 minutes | Notification service logs | M1 |
| **Payment Processing** | Complete transaction within 30 seconds (excluding external gateway time) | Payment service metrics | M1 |
| **Report Generation** | <30 seconds for standard reports, <5 minutes for complex analytics | Reporting service logs | M1 |
| **Search Functionality** | <2 seconds for user/policy/claim search with filters | Search index metrics | M1 |

### 6.2 Safety Requirements

| Requirement | Specification | Mitigation Strategy | Priority |
|-------------|---------------|---------------------|----------|
| **Data Loss Prevention** | Zero data loss for committed transactions | ACID database transactions, daily backups with point-in-time recovery | M1 |
| **Payment Idempotency** | Duplicate payment requests result in same outcome | Idempotency keys, transaction deduplication | M1 |
| **Graceful Degradation** | System remains functional with degraded services | Circuit breaker pattern, fallback mechanisms | M1 |
| **Audit Trail** | All critical operations logged immutably | Append-only audit log, tamper detection | M1 |
| **Rollback Capability** | Database migrations reversible without data loss | Migration version control, rollback testing | M1 |
| **Disaster Recovery** | RPO (Recovery Point Objective) <1 hour, RTO (Recovery Time Objective) <4 hours | Geo-redundant backups, documented recovery procedures | M2 |

### 6.3 Security Requirements

*(Detailed in Section 8 - Security & Compliance Requirements)*

### 6.4 Software Quality Attributes

**Reliability:**
- **System Uptime:** 99.5% (Phase 1), 99.9% (Phase 2)
- **Mean Time Between Failures (MTBF):** >720 hours (30 days)
- **Mean Time To Recovery (MTTR):** <1 hour for critical issues
- **Error Rate:** <0.1% of transactions result in errors

**Maintainability:**
- **Code Coverage:** >70% unit test coverage for critical business logic
- **Code Quality:** SonarQube quality gate pass (no critical/blocker issues)
- **Documentation:** All APIs documented with OpenAPI specification
- **Deployment Frequency:** Daily deployments to staging, weekly to production (Phase 1), continuous to production (Phase 2)

**Scalability:**
- **Horizontal Scaling:** All stateless services support horizontal scaling
- **Database Scaling:** Read replicas for reporting (Phase 2), sharding strategy defined (Phase 3)
- **Auto-scaling:** Automatic scaling based on CPU/memory metrics (Phase 2)
- **Load Balancing:** Load distribution across multiple service instances

**Usability:**
- **Learnability:** New users can complete policy purchase within 10 minutes (no assistance)
- **Task Completion Time:** Policy purchase average <5 minutes, claim submission <3 minutes
- **Error Recovery:** Users can recover from errors with clear guidance (no data loss)
- **User Satisfaction:** Net Promoter Score (NPS) target >50 (Phase 1), >70 (Phase 2)

**Portability:**
- **Cloud Agnostic:** Application can run on AWS, Azure, or GCP with minimal changes
- **Database Portability:** SQL standard compliance, avoid vendor-specific features
- **Containerization:** All services containerized with Docker
- **Infrastructure as Code:** Terraform or equivalent for infrastructure provisioning

**Compatibility:**
- **Backward Compatibility:** API versions maintained for minimum 6 months after deprecation
- **Browser Compatibility:** Support last 2 major versions of Chrome, Firefox, Safari, Edge
- **Mobile OS Compatibility:** Android 9+ (API 28+), iOS 13+
- **Integration Compatibility:** Webhook payload versioning, graceful handling of unknown fields

### 6.5 Business Rules

**Policy Issuance Rules:**
- Customer must complete KYC (photo + NID) before first policy purchase
- Customer age must be 18-65 years for health insurance, 18-70 for accident insurance
- Maximum sum insured per customer: 1,000,000 BDT (Phase 1), configurable (Phase 2)
- Policy effective date: Day after payment verification (Phase 1), immediately (Phase 2 with live payment)
- Cooling period: 15 days for policy cancellation with full refund (regulatory requirement)

**Claims Processing Rules:**
- Claim must be submitted within 30 days of incident (health insurance), 7 days (accident insurance)
- Policy must be active (premium paid) at time of incident
- Waiting period: 30 days for health insurance (pre-existing condition exclusion), 0 days for accident
- Maximum claim amount cannot exceed sum insured
- Claimed amount must be supported by bills/documents

**Payment Rules:**
- Payment reference number unique per transaction (PRN-YYYYMMDD-XXXXX)
- Payment verification SLA: <1 hour during business hours (9 AM - 5 PM), <4 hours extended
- Payment rejection allows retry with same or different payment method
- Payment refund processed within 7 working days (regulatory requirement)
- Transaction limit (AML): 10,000 BDT per transaction without enhanced verification (Phase 1)

**Commission Rules (Phase 1.5):**
- Partner commission: Configurable percentage (default 10%) of premium
- Commission accrued on policy activation, payable monthly
- Minimum payout threshold: 5,000 BDT (accumulated commission)
- Commission payment via bank transfer within 15 days of month-end

---

## 7. Data Model & Storage Requirements

### 7.1 Logical Data Model

**Core Entity Relationships:**

USER (customer, admin, partner)
  ├─ has many ──> POLICY
  ├─ has many ──> CLAIM
  ├─ has many ──> PAYMENT
  └─ has many ──> AUDIT_LOG

PRODUCT
  ├─ has many ──> POLICY
  └─ has many ──> PRICING_RULE

POLICY
  ├─ belongs to ──> USER (customer)
  ├─ belongs to ──> PRODUCT
  ├─ has one ──> PAYMENT
  ├─ has many ──> CLAIM
  └─ has one ──> CERTIFICATE (PDF)

CLAIM
  ├─ belongs to ──> POLICY
  ├─ has many ──> CLAIM_DOCUMENT
  ├─ has many ──> APPROVAL_STEP
  └─ has one ──> SETTLEMENT_PAYMENT

PAYMENT
  ├─ belongs to ──> POLICY
  ├─ has one ──> PAYMENT_PROOF (Phase 1 manual)
  └─ has many ──> PAYMENT_RECONCILIATION_ENTRY

PARTNER
  ├─ has many ──> PARTNER_USER
  ├─ has many ──> POLICY (sold via partner)
  └─ has many ──> COMMISSION_STATEMENT
### 7.2 Data Dictionary

**Sample Critical Entities:**

**users table:**
- user_id (UUID, PK): Unique user identifier
- user_type (ENUM): customer, admin, partner
- mobile_number (VARCHAR 15, UNIQUE): Bangladesh mobile format
- email (VARCHAR 255): Optional for customers, required for admins
- full_name (VARCHAR 255): Legal name
- date_of_birth (DATE): For age validation
- nid_number (VARCHAR 20): National ID, encrypted
- kyc_status (ENUM): pending, verified, rejected
- created_at (TIMESTAMP): Registration timestamp
- updated_at (TIMESTAMP): Last update timestamp

**policies table:**
- policy_id (UUID, PK): Unique policy identifier
- policy_number (VARCHAR 50, UNIQUE): Human-readable format POL-YYYY-NNNNNN
- customer_id (UUID, FK → users): Policy holder
- product_id (UUID, FK → products): Insurance product
- sum_insured (DECIMAL 12,2): Coverage amount in BDT
- premium_amount (DECIMAL 10,2): Premium in BDT
- start_date (DATE): Policy effective date
- end_date (DATE): Policy expiry date
- status (ENUM): pending_payment, active, expired, cancelled, lapsed
- payment_id (UUID, FK → payments): Associated payment
- created_at (TIMESTAMP)

**claims table:**
- claim_id (UUID, PK)
- claim_number (VARCHAR 50, UNIQUE): CLM-YYYY-NNNNNN
- policy_id (UUID, FK → policies)
- incident_date (DATE): When incident occurred
- incident_type (VARCHAR 100): Category of claim
- claimed_amount (DECIMAL 12,2): Amount requested
- approved_amount (DECIMAL 12,2): Amount approved (nullable)
- status (ENUM): submitted, under_review, approved, rejected, paid
- assigned_to (UUID, FK → users): Claims adjuster
- created_at (TIMESTAMP)
- updated_at (TIMESTAMP)

### 7.3 Database Requirements

**Data Consistency:**
- ACID transactions for all policy and payment operations
- Foreign key constraints enforced
- Referential integrity maintained
- Optimistic locking for concurrent updates (version field)

**Data Partitioning (Phase 2):**
- Partition policies table by year (YYYY) for performance
- Partition claims table by year (YYYY)
- Partition audit_logs by month (YYYY-MM)

**Indexing Strategy:**
- Primary key indexes on all tables
- Unique indexes on business keys (policy_number, claim_number, PRN)
- Composite indexes on frequently queried combinations (user_id + status, policy_id + claim_date)
- Full-text search index on product catalog (Phase 2)

### 7.4 Data Migration Requirements

**Initial Data Load (Phase 1):**
- Product catalog (3 products)
- Admin users (System Admin, Business Admin, L1/L2 Admins)
- System configuration (payment methods, notification templates)
- Lookup tables (provinces, districts, insurance categories)

**Ongoing Data Sync:**
- NID verification results cached for 30 days (refresh on KYC update)
- Payment reconciliation data imported daily from MFS providers
- Partner commission calculated and updated monthly

### 7.5 Data Archival & Retention

| Data Type | Retention Period | Archival Strategy | Access Method |
|-----------|------------------|-------------------|---------------|
| **Active Policies** | Duration + 2 years | Online database | Real-time |
| **Expired Policies** | 7 years (regulatory) | Cold storage (S3 Glacier) | Retrieval within 24 hours |
| **Processed Claims** | 7 years (regulatory) | Cold storage after 1 year | Retrieval within 24 hours |
| **Audit Logs** | 7 years (regulatory) | Hot storage for 1 year, then cold | Real-time (1 year), batch retrieval (older) |
| **Payment Records** | 7 years (regulatory) | Hot storage for 1 year, then cold | Real-time (1 year), batch retrieval (older) |
| **User Data** | Until account deletion | Online database | Real-time |
| **System Logs** | 90 days | Rotate and delete | Real-time |
| **Backup Files** | 30 days operational, 7 years compliance | Geo-redundant storage | Point-in-time restore |

---


## 8. Security & Compliance Requirements

### 8.1 Security Requirements

| SEC ID | Requirement Description | Implementation Guidance | Priority |
|--------|------------------------|------------------------|----------|
| **SEC-001** | The system shall encrypt all data at rest using AES-256 encryption: database (transparent data encryption), object storage (server-side encryption), backups (encrypted before storage). | Database TDE feature, S3 server-side encryption (SSE-KMS), encrypted backup volumes | M1 |
| **SEC-002** | The system shall encrypt all data in transit using TLS 1.3 (minimum TLS 1.2 for legacy compatibility): client-server communication, service-to-service communication, external API calls. | TLS certificates from trusted CA, HTTP Strict Transport Security (HSTS) headers, certificate pinning for mobile apps | M1 |
| **SEC-003** | The system shall implement secure session management: server-side session storage, httpOnly and secure cookie flags, session timeout after 30 minutes inactivity, concurrent session limit (5 per user), session invalidation on logout. | Redis for session storage, secure cookie configuration, session middleware with timeout tracking | M1 |
| **SEC-004** | The system shall implement JWT-based API authentication with: RS256 algorithm (asymmetric), access token expiry 15 minutes, refresh token expiry 7 days, token revocation support, refresh token rotation. | JWT library with key management, token blacklist for revocation, refresh endpoint with secure token exchange | M1 |
| **SEC-005** | The system shall maintain comprehensive audit logs for all sensitive operations: user authentication (login, logout, password change), policy operations (create, update, cancel), payment operations (verification, approval, rejection), claims operations (submission, approval, settlement), admin actions (user management, configuration changes), access control changes. Logs immutable and tamper-evident. | Structured logging to append-only storage, log integrity verification with checksums, centralized log aggregation | M1 |
| **SEC-006** | The system shall implement input validation and sanitization per OWASP guidelines: validate all user inputs against whitelist, sanitize outputs to prevent XSS, parameterized queries to prevent SQL injection, file upload validation (type, size, content), rate limiting on all endpoints. | Input validation library, prepared statements for database, file type verification, rate limiter middleware | M1 |
| **SEC-007** | The system shall conduct regular security audits: annual penetration testing by third-party, quarterly vulnerability scans, continuous dependency scanning, code security review (SAST), runtime security monitoring (RASP). | Contract with security firm, automated scanning tools (Snyk, OWASP Dependency-Check), code review in CI/CD pipeline | M2 |
| **SEC-008** | The system shall implement Web Application Firewall (WAF) with rules for: SQL injection detection, XSS attack prevention, DDoS mitigation (rate limiting, IP blocking), bot detection, geolocation filtering (Bangladesh-only for Phase 1). | Cloud WAF service (AWS WAF, Cloudflare WAF), custom rule configuration, monitoring and alerting | S |
| **SEC-009** | The system shall implement API security best practices: OAuth 2.0 for third-party access, API key rotation every 90 days, request signing for webhooks (HMAC-SHA256), scope-based authorization, API versioning for backward compatibility. | OAuth provider integration, API key management interface, webhook signature verification, scoped permissions | M2 |
| **SEC-010** | The system shall implement secrets management: no secrets in source code, environment-based configuration, encrypted secrets storage, access control on secrets, secrets rotation policy (90 days). | Secrets manager service (AWS Secrets Manager, HashiCorp Vault), environment variable injection, automated rotation | M1 |
| **SEC-011** | The system shall implement rate limiting: 100 requests/minute per user for general APIs, 3 OTP requests/hour per mobile number, 10 login attempts/hour per user, adjustable limits per API endpoint. | Rate limiter with sliding window algorithm, Redis for rate limit counters, response headers indicating limits | M1 |
| **SEC-012** | The system shall implement CORS (Cross-Origin Resource Sharing) policy: whitelist allowed origins, restrict methods (GET, POST, PUT, DELETE), validate Origin header, no wildcard in production. | CORS middleware configuration, environment-specific allowed origins | M1 |
| **SEC-013** | The system shall implement Content Security Policy (CSP) headers: restrict script sources to same-origin and trusted CDNs, disallow inline scripts, restrict frame ancestors, report violations. | CSP headers in web server configuration, CSP violation reporting endpoint | M1 |
| **SEC-014** | The system shall implement password security (for admin/partner users): minimum 12 characters, complexity requirements (uppercase, lowercase, digit, special), bcrypt hashing (cost factor 12), password history (prevent reuse of last 5), password expiry 90 days. | Password validation library, bcrypt for hashing, password history table | M1 |
| **SEC-015** | The system shall implement account lockout policy: lock account after 5 failed login attempts, lockout duration 30 minutes or until admin unlock, CAPTCHA after 3 failed attempts, notification to user on lockout. | Failed login counter, lockout flag in user table, CAPTCHA integration | M1 |
| **SEC-016** | The system shall implement data masking for sensitive information: mask NID (show last 4 digits), mask mobile (show last 4 digits), mask bank account (show last 4 digits), mask email (show first char + domain), full access only for authorized roles. | Data masking utility functions, role-based field access control | M1 |
| **SEC-017** | The system shall implement security headers: X-Content-Type-Options: nosniff, X-Frame-Options: DENY, X-XSS-Protection: 1; mode=block, Referrer-Policy: strict-origin-when-cross-origin, Permissions-Policy (restrict camera, microphone, geolocation). | Web server or reverse proxy configuration, security header middleware | M1 |
| **SEC-018** | The system shall implement intrusion detection: monitor for suspicious patterns (rapid API calls, unusual access patterns, privilege escalation attempts), automated alerting on threshold breach, incident response playbook. | Log analysis with anomaly detection, alerting system (PagerDuty, Slack), documented incident response procedures | M2 |
| **SEC-019** | The system shall implement data breach response plan: incident detection and containment procedures, notification process (IDRA, BFIU, affected customers), forensic investigation guidelines, communication templates. | Documented playbook, contact list, communication templates, regular drills | M2 |
| **SEC-020** | The system shall implement security awareness: regular security training for developers, secure coding guidelines, security checklist in code review, vulnerability disclosure policy (responsible disclosure). | Training program, security documentation, code review checklist, public security policy page | M2 |

### 8.2 IDRA Compliance Requirements

**Insurance Development & Regulatory Authority (Bangladesh) - Digital Insurance Guidelines:**

| Requirement Area | Specification | SRS Reference |
|------------------|---------------|---------------|
| **Licensing** | Digital insurance license obtained before public launch | Business prerequisite |
| **Product Approval** | All insurance products must receive IDRA approval before offering (30-45 days per product) | FR-031 (Product management) |
| **KYC/AML Compliance** | Customer identification per IDRA circular, NID verification mandatory | FR-001, FR-004, FR-008 |
| **Policy Issuance** | Policy certificate must include: policy number, customer details, coverage, premium, T&Cs, digital signature | FR-042 |
| **Claims Processing** | Transparent claims process, defined TAT, customer notification at each stage | FR-055 to FR-065 |
| **Data Protection** | Customer data protection, consent-based data sharing, secure storage | SEC-001, SEC-002, SEC-016 |
| **Financial Reporting** | Quarterly CARAMELS report, annual Financial Condition Report (FCR) | FR-093 |
| **Complaint Handling** | Customer grievance mechanism, escalation process, regulatory reporting of complaints | Customer support requirements |
| **Agent Licensing** | If using agents, must be licensed (Phase 1 - no agents, direct-to-customer) | N/A Phase 1 |
| **Solvency Margin** | Maintain minimum solvency ratio as per IDRA norms | Business financial management |
| **System Audit** | Annual IT system audit for digital insurers | SEC-007 |
| **Business Continuity** | Disaster recovery plan, data backup, system redundancy | FR-102, Section 6.2 |

**IDRA Reporting Requirements:**

1. **Quarterly CARAMELS Report:** Capital adequacy, Reinsurance, Assets, Management, Earnings, Liquidity, Sensitivity analysis
2. **Annual FCR (Financial Condition Report):** Comprehensive financial health report
3. **Policy Register:** All policies issued (monthly submission in Phase 1)
4. **Claims Register:** All claims with status (monthly submission)
5. **Complaint Register:** Customer complaints and resolution

### 8.3 AML/CFT Compliance Requirements

**Bangladesh Financial Intelligence Unit (BFIU) - Anti-Money Laundering & Countering Financing of Terrorism:**

| Requirement | Specification | Implementation | Priority |
|-------------|---------------|----------------|----------|
| **Customer Due Diligence (CDD)** | Verify customer identity (NID), collect source of funds for high-value transactions | KYC process (FR-001, FR-004, FR-008), transaction declaration form | M1 |
| **Transaction Monitoring** | Monitor transactions for suspicious patterns, threshold: 10,000 BDT (Phase 1), 50,000 BDT (Phase 2) | Automated flagging system, manual review workflow | M1 |
| **High-Value Transaction Reporting** | Report cash transactions >1,000,000 BDT to BFIU | CTR (Cash Transaction Report) generation and submission | M2 |
| **Suspicious Activity Reporting** | File STR (Suspicious Transaction Report) for flagged activities within 7 days | STR submission workflow, secure channel to BFIU portal | S |
| **Record Keeping** | Maintain transaction records for 7 years | Data retention policy (Section 7.5) | M1 |
| **Risk Assessment** | Customer risk profiling (low, medium, high), enhanced due diligence for high-risk | Risk scoring algorithm, enhanced verification workflow | M2 |
| **Politically Exposed Persons (PEP)** | Identify and flag PEP customers, enhanced monitoring | PEP database integration or manual declaration | S |
| **Sanctions Screening** | Screen customers against UN/national sanctions lists | Sanctions list integration, screening on registration and periodically | S |
| **Training** | AML/CFT training for staff handling transactions | Training program, certification tracking | M2 |

**Transaction Monitoring Rules:**

- Single transaction >10,000 BDT: Enhanced verification required (source of funds declaration)
- Multiple transactions totaling >50,000 BDT in 24 hours: Flagged for review
- Policy with sum insured >100,000 BDT: Enhanced CDD required
- Rapid succession of claims (>3 in 6 months): Fraud investigation triggered

### 8.4 Privacy & Data Protection

**Data Privacy Principles (Based on GDPR and Bangladesh ICT Act):**

| Principle | Requirement | Implementation |
|-----------|-------------|----------------|
| **Lawfulness** | Process personal data only with customer consent or legal basis | Consent checkbox during registration, privacy policy acceptance |
| **Purpose Limitation** | Collect data only for specified purposes, no repurposing without consent | Data collection form explanations, purpose declaration |
| **Data Minimization** | Collect only necessary data, avoid excessive information | Review all form fields for necessity |
| **Accuracy** | Maintain accurate data, provide correction mechanism | Profile edit functionality, data verification workflows |
| **Storage Limitation** | Retain data only as long as necessary (7 years for compliance data) | Data retention policy (Section 7.5), automated purging |
| **Integrity & Confidentiality** | Protect data with appropriate security measures | Encryption (SEC-001, SEC-002), access control (FG-002) |
| **Accountability** | Demonstrate compliance with data protection principles | Privacy policy, audit logs, compliance reports |

**Customer Rights:**

1. **Right to Access:** Customer can request copy of their personal data (self-service via portal)
2. **Right to Correction:** Customer can correct inaccurate data (profile edit)
3. **Right to Deletion:** Customer can request account deletion (with regulatory retention exception)
4. **Right to Portability:** Customer can download their data in machine-readable format (CSV/JSON export)
5. **Right to Withdraw Consent:** Customer can opt-out of marketing communications (preferences management)
6. **Right to Complain:** Customer can file complaint with company and escalate to regulator

**Privacy Policy Contents:**

- What data is collected and why
- How data is used and shared
- Data retention period
- Security measures
- Third-party data sharing (payment processors, partners)
- Customer rights and how to exercise them
- Contact information for privacy officer
- Updates to privacy policy (notification mechanism)

---

## 9. Operational Requirements

### 9.1 Support Model

**Customer Support (Phase 1):**
- **Channels:** Phone (16999), Email (support@labaidinsurance.com), In-app chat (Phase 2)
- **Operating Hours:** 9 AM - 5 PM (Sunday-Thursday), Emergency support for critical issues
- **Response SLA:** Phone answered within 3 minutes, Email response within 4 hours (business hours)
- **Staffing:** 10 support agents (Phase 1), expanding based on volume
- **Knowledge Base:** Self-service FAQs, video tutorials (Bengali + English)

**Technical Support (Internal):**
- **L1 Support:** Application issues, user account issues, first-line troubleshooting
- **L2 Support:** Database issues, integration failures, performance degradation
- **L3 Support:** Architecture issues, security incidents, core system bugs
- **Escalation Path:** L1 → L2 (within 2 hours if unresolved) → L3 (within 4 hours)
- **On-Call:** Rotating on-call schedule for after-hours critical issues

**Partner Support:**
- **Dedicated Contact:** Partner success manager for each major partner
- **Technical Integration Support:** API documentation, sandbox environment, integration testing
- **Response SLA:** Critical issues <2 hours, general inquiries <24 hours

### 9.2 Maintenance Requirements

**Scheduled Maintenance:**
- **Weekly Maintenance Window:** Sundays 2 AM - 6 AM (lowest traffic period)
- **Monthly Security Patching:** Database, OS, dependencies updated monthly
- **Quarterly Major Updates:** Feature releases, significant changes (with rollback plan)
- **Advance Notice:** 7 days notice for scheduled downtime, communicated via email, SMS, in-app banner

**Unscheduled Maintenance:**
- **Emergency Patches:** Security vulnerabilities patched within 24 hours of discovery
- **Hot Fixes:** Critical bugs fixed and deployed within 4 hours
- **Configuration Changes:** No-downtime configuration updates via feature flags

**Monitoring & Alerting:**
- **System Health Monitoring:** CPU, memory, disk, network metrics per service
- **Application Monitoring:** Error rates, response times, transaction volumes
- **Business Metrics Monitoring:** Registrations, policy sales, claims submitted (real-time dashboard)
- **Alert Channels:** PagerDuty for critical, Slack for warnings, Email for informational

### 9.3 Backup & Recovery

**Backup Strategy:**

| Data Type | Frequency | Retention | Storage Location | RTO Target | RPO Target |
|-----------|-----------|-----------|------------------|-----------|-----------|
| **Database (Full)** | Daily (2 AM) | 30 days | S3 (geo-redundant) | 4 hours | 24 hours |
| **Database (Incremental)** | Every 6 hours | 7 days | S3 (geo-redundant) | 1 hour | 6 hours |
| **Object Storage** | Continuous (versioning) | 30 days | S3 (versioning enabled) | 2 hours | Near-zero |
| **Application Code** | On every commit | Infinite | Git repository | 30 minutes | Zero |
| **Configuration** | On every change | Infinite | Git repository | 30 minutes | Zero |
| **Audit Logs** | Real-time replication | 7 years | S3 (append-only) | 4 hours | Near-zero |

**Recovery Procedures:**

1. **Database Recovery:** Documented point-in-time recovery procedure, tested monthly
2. **Application Recovery:** Automated rollback to previous version, Docker image tagged releases
3. **Data Center Failover:** Documented failover to secondary region (Phase 2)
4. **Recovery Testing:** Full disaster recovery drill quarterly, results documented

### 9.4 Disaster Recovery

**Disaster Scenarios & Recovery:**

| Scenario | Impact | Recovery Procedure | RTO | RPO |
|----------|--------|-------------------|-----|-----|
| **Database Failure** | Transaction processing halted | Restore from latest backup to standby instance, redirect traffic | 1 hour | 6 hours |
| **Application Server Failure** | Service degradation | Auto-scaling launches new instances, load balancer redirects | 5 minutes | Zero |
| **Object Storage Failure** | Document upload/download unavailable | Failover to redundant storage region | 30 minutes | Zero (versioned) |
| **Data Center Outage** | Complete service outage | Failover to secondary region (Phase 2), restore from backups | 4 hours | 24 hours |
| **Network Partition** | Service unreachable from specific regions | Re-route traffic via alternative ISP/CDN | 30 minutes | Zero |
| **Security Breach** | Data compromise | Isolate affected systems, restore from clean backup, incident response | 2 hours | Variable |
| **Database Corruption** | Data integrity compromised | Point-in-time recovery to last known good state | 2 hours | Variable |

**Business Continuity Measures:**

- **Communication Plan:** Stakeholder notification tree (customers, partners, regulators, team)
- **Alternate Work Arrangements:** Remote work capability for all roles
- **Backup Payment Processing:** Manual payment workflow always available as fallback
- **Manual Claims Processing:** Paper-based fallback process documented
- **Regulatory Communication:** Pre-drafted templates for IDRA/BFIU notification

---

## 10. Acceptance Criteria & Testing

### 10.1 Phase 1 Acceptance Criteria (March 1, 2026 - Beta Launch)

**Functional Acceptance:**
- [ ] 50 beta users can register via mobile OTP with >98% success rate
- [ ] 3 insurance products visible in catalog with correct details
- [ ] Users can complete full policy purchase flow with manual payment
- [ ] Admin can verify uploaded payment proofs and activate policies within 1-hour SLA
- [ ] Users can submit claims with document upload
- [ ] Admin can approve/reject claims via approval matrix workflow
- [ ] SMS notifications sent for all critical events with >98% delivery rate
- [ ] Policy certificates generated as PDF within 30 seconds
- [ ] System operates on 3G network with acceptable performance (<3s page load)

**Non-Functional Acceptance:**
- [ ] System uptime >95% over 7-day beta period
- [ ] API response time 95th percentile <500ms
- [ ] Zero critical security vulnerabilities (security scan passed)
- [ ] All data encrypted at rest and in transit (verification test passed)
- [ ] Audit logs captured for all admin actions
- [ ] Backup and recovery procedure tested successfully

**Business Acceptance:**
- [ ] 25 policies sold during beta period
- [ ] 5 claims submitted and processed
- [ ] Average payment verification TAT <1 hour
- [ ] Average claims processing TAT <48 hours
- [ ] User satisfaction survey >4.0/5.0 rating
- [ ] Zero customer data breaches or security incidents

### 10.2 Phase 1.5 Acceptance Criteria (May 1, 2026 - Public Launch)

**Additional Functional Acceptance:**
- [ ] Android app published on Google Play Store with 4.0+ rating
- [ ] bKash live payment integration with >98% success rate
- [ ] Nagad payment integration with >98% success rate
- [ ] Email notifications delivered with >95% success rate
- [ ] Partner portal functional for 2 pilot hospital partners
- [ ] Automated small claim approval (<5,000 BDT) working
- [ ] Basic fraud detection flagging suspicious claims

**Additional Non-Functional Acceptance:**
- [ ] System uptime >99% over 30-day period
- [ ] Support 500 concurrent users with acceptable performance
- [ ] Daily reconciliation with payment providers automated
- [ ] IDRA provisional approval obtained

**Additional Business Acceptance:**
- [ ] 100 policies sold in first month post-launch
- [ ] 500 registered users
- [ ] 10 claims processed with <48h average TAT
- [ ] 2 hospital partners onboarded
- [ ] Net Promoter Score (NPS) >50

### 10.3 Phase 2 Acceptance Criteria (November 1, 2026 - Scale)

**Key Acceptance Criteria:**
- [ ] iOS app published on Apple App Store
- [ ] System supports 5,000 concurrent users
- [ ] Automated claim processing for 80% of small claims
- [ ] EHR integration with 10 hospitals operational
- [ ] IoT livestock tracking pilot with 1,000 devices
- [ ] BI dashboards functional for business insights
- [ ] System uptime >99.9%
- [ ] 1,000 policies sold per month
- [ ] 50 partner hospitals onboarded

### 10.4 Phase 3 Acceptance Criteria (November 1, 2027 - Innovation)

**Key Acceptance Criteria:**
- [ ] Voice-assisted workflows functional in Bengali
- [ ] WebRTC video claims verification operational
- [ ] USSD fallback for feature phones working
- [ ] System supports 50,000 concurrent users
- [ ] Multi-language support (Bengali, English) complete
- [ ] Usage-based insurance pricing operational

### 10.5 Test Summary

**Test Types & Coverage:**

| Test Type | Scope | Target Coverage | Responsibility | Phase |
|-----------|-------|-----------------|----------------|-------|
| **Unit Testing** | Individual functions/methods | >70% code coverage | Developers | M1 |
| **Integration Testing** | Service-to-service interactions | All integration points | Developers | M1 |
| **API Testing** | REST/gRPC endpoints | All endpoints | QA Team | M1 |
| **UI Testing** | User interfaces | Critical user flows | QA Team | M1 |
| **End-to-End Testing** | Complete user journeys | 5 critical flows | QA Team | M1 |
| **Performance Testing** | Load, stress, endurance | 10,000 concurrent users target | QA Team | M2 |
| **Security Testing** | Vulnerability assessment | OWASP Top 10 | Security Team | M2 |
| **Penetration Testing** | Ethical hacking | Full system | Third-party | M2 |
| **Accessibility Testing** | WCAG 2.1 AA compliance | All user-facing screens | QA Team | M2 |
| **Compatibility Testing** | Browser/device compatibility | Target matrix | QA Team | M2 |
| **Regression Testing** | Existing functionality after changes | Automated suite | CI/CD Pipeline | Ongoing |
| **User Acceptance Testing** | Business requirements validation | All functional requirements | Business Users | Pre-launch |

---

## 11. Appendices

### Appendix A: Priority Classification Guide

**Phase Assignment Decision Tree:**

Is the requirement essential for beta launch (50 users)?
├─ YES → Is it technically feasible without external dependencies?
│         ├─ YES → M1 (Phase 1)
│         └─ NO → M2 (Phase 1.5) [e.g., bKash API after approval]
│
└─ NO → Is it required for public launch compliance?
          ├─ YES → M2 (Phase 1.5)
          └─ NO → Does it support growth/scale (5,000+ users)?
                    ├─ YES → S (Phase 2)
                    └─ NO → D/C/F (Phase 3)
**Priority Code Definitions:**

- **M1:** Core functionality without which beta launch impossible (blocking)
- **M2:** Required for public launch compliance and market readiness (critical for May 1)
- **S:** Supports scaling to 5,000+ users, automation, performance optimization
- **D:** Desirable features that improve user experience or operational efficiency
- **C:** Could-have features if time/budget permits
- **F:** Future innovation features requiring market validation first

### Appendix B: Glossary

*(Refer to Section 1.4 for complete glossary of 70+ terms)*

### Appendix C: Regulatory References

**Bangladesh Laws & Regulations:**
1. Insurance Act 2010 - Primary legislation governing insurance industry
2. IDRA Circular 2023/05 - Digital Insurance Guidelines
3. BFIU Circular No. 37 - Anti-Money Laundering Guidelines
4. ICT Act 2006 (Amended 2013) - Information and Communication Technology regulations
5. Digital Security Act 2018 - Cybersecurity and digital forensics
6. Payment and Settlement Systems Regulations 2014 - Financial transaction regulations

**International Standards:**
7. ISO/IEC 27001:2013 - Information Security Management Systems
8. PCI DSS 3.2.1 - Payment Card Industry Data Security Standard (if card payments added)
9. OWASP Top 10 2021 - Web Application Security Risks
10. WCAG 2.1 Level AA - Web Content Accessibility Guidelines

### Appendix D: Phase Delivery Schedule

**Phase 1 (M1) - March 1, 2026:**
- Duration: 10 weeks (Dec 15, 2025 - Mar 1, 2026)
- Features: 60 functional requirements
- Milestone: Beta launch with 50 users
- Success Metric: 25 policies sold, 5 claims processed

**Phase 1.5 (M2) - May 1, 2026:**
- Duration: 8 weeks (Mar 1 - May 1, 2026)
- Features: 45 functional requirements (cumulative 105)
- Milestone: Public launch
- Success Metric: 500 users, 100 policies sold

**Phase 2 (S) - November 1, 2026:**
- Duration: 24 weeks (May 1 - Nov 1, 2026)
- Features: 21 functional requirements (cumulative 126)
- Milestone: Scale and automation
- Success Metric: 5,000 users, 1,000 policies/month

**Phase 3 (D/C/F) - November 1, 2027:**
- Duration: 52 weeks (Nov 1, 2026 - Nov 1, 2027)
- Features: 12 functional requirements (cumulative 138)
- Milestone: Innovation and market leadership
- Success Metric: 50,000 users, enterprise scale

---

## Document Control & Version Management

**Change Control Process:**
1. Changes must be submitted via RFC (Request For Change) template
2. RFC reviewed in weekly change control board meeting
3. Impact assessment completed (technical, schedule, budget)
4. Approval required from PM + CTO for minor changes, Steering Committee for major changes
5. Approved changes incorporated in next SRS version with revision history updated

**Distribution List:**
- Project Manager: Master copy owner
- CTO: Technical review copy
- Development Team Leads: Working copies
- QA Lead: Testing reference copy
- Business Analyst: Requirements traceability copy
- Compliance Officer: Regulatory compliance copy

**Document Storage:**
- **Master Copy:** Project repository (Git)
- **Published Version:** Project portal (read-only PDF)
- **Working Drafts:** Collaborative editing platform

---

## Acknowledgments

This System Requirements Specification was prepared through collaborative effort of:
- **Business Stakeholders:** Requirements gathering and prioritization
- **Technical Team:** Feasibility analysis and architecture input
- **Compliance Team:** Regulatory requirements integration
- **QA Team:** Test scenarios and acceptance criteria definition
- **External Consultants:** Industry best practices and benchmarking

**Special Recognition:**
- **Phased Delivery Strategy:** Pragmatic approach balancing ambition with feasibility, enabling realistic March 1 beta launch

---

## SRS Conclusion

This System Requirements Specification V3.0 provides a comprehensive, implementation-agnostic definition of the LabAid InsureTech Platform across four phased releases. The phased approach acknowledges:

1. **Resource Realities:** Leveraging existing production services for 755-hour development acceleration
2. **External Dependencies:** Manual payment workflow eliminates blocking dependency on payment gateway approvals
3. **Market Validation:** MVP approach validates core value proposition before investing in advanced features
4. **Risk Management:** Incremental delivery reduces risk of catastrophic failure compared to big-bang launch
5. **Regulatory Compliance:** Maintains minimum viable compliance in Phase 1, achieving full compliance by Phase 1.5

**Total Scope:**
- **Functional Requirements:** 138 (60 M1, 45 M2, 21 S, 12 D/C/F)
- **Non-Functional Requirements:** 40+ across performance, security, scalability, usability
- **Security Requirements:** 20 detailed security controls
- **Compliance Requirements:** IDRA, BFIU, Privacy regulations

**Success Criteria:**
- **Phase 1:** Functional beta with 50 users validating core workflows
- **Phase 1.5:** Public launch with 500 users and full compliance
- **Phase 2:** Scaled platform serving 5,000 users with automation
- **Phase 3:** Market-leading platform with 50,000 users and innovation features

This SRS serves as the authoritative specification for all design, development, testing, and acceptance activities. Any team implementing this system—internal or external—has complete requirements without dependency on specific people, project plans, or organizational context.

---

**END OF SYSTEM REQUIREMENTS SPECIFICATION V3.0**

**Document Prepared:** December 2025  
**Next Review:** Post-Phase 1 Launch (March 2026)  
**Document Owner:** Project Manager  
**Status:** Approved for Development

---



