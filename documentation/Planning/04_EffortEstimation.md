## 3. Effort Estimation by Service/Component

### 3.1 Estimation Methodology
- **Unit:** Story Points converted to hours (1 SP = 6-8 hours)
- **Complexity Factors:** Technical complexity, dependencies, team experience
- **Estimation Approach:** Three-point estimation (Optimistic, Most Likely, Pessimistic)
- **Buffer:** 10% added for unknowns and integration testing (reduced from 20% due to existing proven code)

### 3.1.1 Existing Code Reuse Impact

**Proven Components (No Development Needed):**
- ✅ Authentication Service (Go) - Saves ~250 hours
- ✅ Authorization Service (Go) - Saves ~220 hours
- ✅ DBManager Service (Go) - Saves ~280 hours
- ✅ Storage Manager Service (Go) - Saves ~180 hours
- ✅ IoT Broker (Go) - 80% ready, Saves ~200 hours
- ✅ API Gateway (Go) - 50% ready, Saves ~150 hours
- ✅ Payment Service (Node.js) - 70% ready, Saves ~140 hours

**Total Savings from Existing Code:** ~1,420 hours
**Reduced Buffer Justification:** Existing tested code = lower risk = 10% buffer instead of 20%

---

### 3.2 M1 Services - Core Platform (Must Have for Soft Launch - March 1st)

**M1 SCOPE:** Minimal viable features for beta testing and demo on National Insurance Day
- Focus: User registration, basic policy purchase flow, partner portal foundation
- **NOT INCLUDED IN M1:** Claims processing (M2), Mobile apps (M2), Advanced features

#### 3.2.1 User Service (Authentication & Authorization)
**Priority:** M1 | **Owner:** CTO (Go Services - Existing Code 100% Ready)

| Feature | Complexity | Estimated Hours | Status |
|---------|-----------|----------------|--------|
| Authentication (Login/Logout/JWT) | Medium | ✅ 0 hrs | **100% Ready** |
| Authorization (RBAC) | High | ✅ 0 hrs | **100% Ready** |
| User Registration (Phone OTP) | Medium | ✅ 0 hrs | **100% Ready** |
| Profile Management (CRUD) | Low | ✅ 0 hrs | **100% Ready** |
| Password Reset & Recovery | Medium | ✅ 0 hrs | **100% Ready** |
| Session Management | Medium | ✅ 0 hrs | **100% Ready** |
| Multi-factor Authentication | High | ✅ 0 hrs | **100% Ready** |
| Insurance-specific User Roles (8 roles) | Low | 16 hrs | Role config only |
| API Documentation Update | Low | 8 hrs | Update for insurance domain |
| Integration Tests (Insurance context) | - | 16 hrs | Context-specific testing |
| **SUBTOTAL** | - | **40 hrs** | - |
| **Buffer (10%)** | - | **4 hrs** | - |
| **TOTAL** | - | **44 hrs** | - |

**Team Assignment:** CTO (40% time - 19 hrs/week) - Can complete in 3 days
**Savings:** ~400 hours from proven authentication/authorization system

---

#### 3.2.2 Insurance Engine - Policy Service (M1 - Core Only)
**Priority:** M1 | **Owner:** Mr. Delowar (C# .NET Lead) + C# Mid-level Developer

**M1 Scope:** Basic policy creation and management only. Advanced features moved to M2.

| Feature | Complexity | Estimated Hours | Dependencies | Priority |
|---------|-----------|----------------|--------------|----------|
| Insurance Domain Models (C#) | High | 48 hrs | DBManager (existing) | M1 |
| Policy CRUD Operations (gRPC) | Medium | 40 hrs | User Service (existing) | M1 |
| Basic Policy Search & Filters | Low | 24 hrs | Policy CRUD | M1 |
| Premium Calculation Engine (Simple) | High | 56 hrs | Actuarial formulas | M1 |
| Beneficiary Management | Medium | 24 hrs | Policy CRUD | M1 |
| Policy Status Management | Low | 16 hrs | Policy CRUD | M1 |
| gRPC API Design & Implementation | Medium | 32 hrs | Core features | M1 |
| Unit & Integration Tests (.NET) | - | 48 hrs | Core features | M1 |
| **SUBTOTAL** | - | **288 hrs** | - | - |
| **Buffer (10%)** | - | **29 hrs** | - | - |
| **TOTAL M1** | - | **317 hrs** | - | - |

**Moved to M2 (460 hrs):** Policy Renewal Logic, Cancellation, History/Versioning, Contract Management, Coverage/Riders, Risk Assessment

**Team Assignment:** Mr. Delowar (Jan 15+) + C# Mid Dev (Jan 1+)
- 2 devs × 96 hrs/week = can complete in 3.5 weeks
**Note:** Heavy use of existing DBManager (saves 280 hrs) and Auth services (saves 250 hrs)

---

#### 3.2.3 Claim Service
**Priority:** M2 (NOT M1!) | **Owner:** Backend Team (Delowar + C# Dev)

**⚠️ CONFIRMED M2 PRIORITY** - Claims moved to M2 after Grand Launch
- SRS Analysis: FR-041 to FR-058 show mixed M1/M2/M3 priorities
- **Decision:** Move ALL claims to M2 for capacity management
- M1 focus: Policy purchase flow for beta demo only
- Claims processing starts after Grand Launch (April 14, 2026)
- Reasoning: M1 already at 69% capacity with core features

| Feature | Complexity | Estimated Hours | Dependencies | Priority |
|---------|-----------|----------------|--------------|----------|
| Claim Models & Schema | Medium | 24 hrs | Policy Service | M2 |
| Claim Filing & Submission | Medium | 48 hrs | Policy, Document | M2 |
| Claim Status Tracking | Low | 24 hrs | Claim Filing | M3 |
| Claim Approval Workflow | High | 72 hrs | User Service | M3 |
| Claim Assessment Logic | High | 64 hrs | Policy Service | M3 |
| Claim Settlement | High | 56 hrs | Payment Service | M3 |
| Fraud Detection (Basic) | High | 64 hrs | Claim Models | M3 |
| Claim Document Management | Medium | 40 hrs | Document Service | M2 |
| Claim History & Reporting | Medium | 32 hrs | Claim Models | M3 |
| Claim Notifications | Low | 24 hrs | Notification Service | M2 |
| API Documentation | Low | 16 hrs | All features | M2 |
| Unit & Integration Tests | - | 64 hrs | All features | M2 |
| **SUBTOTAL** | - | **528 hrs** | - | - |
| **Buffer (20%)** | - | **106 hrs** | - | - |
| **TOTAL M2** | - | **634 hrs** | - | - |

**Team Assignment:** 3 Backend Devs in M2 phase (after March 1st)

---

#### 3.2.4 Payment Service
**Priority:** M1 | **Owner:** Mamoon (Node.js - 70% Existing Code)

| Feature | Complexity | Estimated Hours | Status | Priority |
|---------|-----------|----------------|--------|----------|
| Payment Gateway Integration (Bkash) | High | ✅ 0 hrs | **100% Ready** | M1 |
| Payment Processing (Premium) | High | 16 hrs | Minor adaptation | M1 |
| Invoice Generation | Medium | 8 hrs | 70% ready, polish | M1 |
| Receipt Generation | Low | 12 hrs | Customization | M1 |
| Payment History & Tracking | Medium | 8 hrs | 70% ready, adapt | M1 |
| Payment Notifications | Low | 8 hrs | Integration | M1 |
| Failed Payment Handling | Medium | 12 hrs | Enhancement | M1 |
| API Documentation Update | Low | 8 hrs | Update docs | M1 |
| Unit & Integration Tests | - | 24 hrs | Additional tests | M1 |
| **SUBTOTAL M1** | - | **96 hrs** | - | - |
| **Buffer (10%)** | - | **10 hrs** | - | - |
| **TOTAL M1** | - | **106 hrs** | - | - |

**Moved to M2 (171 hrs):**
- Payment Processing (Claims) - 32 hrs (Claims are M2)
- Refund Processing - 24 hrs (Policy cancellation is M2)
- Payment Reconciliation - 40 hrs (Advanced feature)
- Multi-gateway Support (Nagad/Cards) - 48 hrs (Nice to have)
- Buffer - 27 hrs

**Team Assignment:** Mamoon (50% time = 24 hrs/week) - Can complete in ~4.5 weeks
**Savings:** ~450 hours from existing Bkash integration and payment infrastructure

---

#### 3.2.5 Document Service (Storage Manager)
**Priority:** M1 | **Owner:** CTO (Go - Existing Code 100% Ready) + Sagor

| Feature | Complexity | Estimated Hours | Status | Priority |
|---------|-----------|----------------|--------|----------|
| File Upload (Multi-format) | Medium | ✅ 0 hrs | **100% Ready** | M1 |
| Cloud Storage Integration (S3/Azure) | High | ✅ 0 hrs | **100% Ready** | M1 |
| Document Metadata Management | Medium | ✅ 0 hrs | **100% Ready** | M1 |
| Document Retrieval & Download | Low | ✅ 0 hrs | **100% Ready** | M1 |
| Document Versioning | Medium | ✅ 0 hrs | **100% Ready** | M1 |
| Document Security & Access Control | High | 8 hrs | Integration with Auth | M1 |
| Insurance Document Types | Low | 12 hrs | Policy, KYC docs config | M1 |
| Document Tagging & Categorization | Low | 16 hrs | Basic categories | M1 |
| API Documentation Update | Low | 8 hrs | Update for insurance | M1 |
| Unit & Integration Tests | - | 16 hrs | Additional tests | M1 |
| **SUBTOTAL M1** | - | **60 hrs** | - | - |
| **Buffer (10%)** | - | **6 hrs** | - | - |
| **TOTAL M1** | - | **66 hrs** | - | - |

**Moved to M2/M3 (119 hrs):**
- Document Preview Generation - 32 hrs (Nice to have)
- OCR Integration (Basic) - 48 hrs (Advanced feature for claims)
- Advanced tagging - 8 hrs
- Buffer - 11 hrs

**Team Assignment:** CTO (40% time = 19 hrs/week) + Sagor - Can complete in ~2 weeks
**Savings:** ~450 hours from existing Storage Manager service (S3, versioning, security all ready)

---

#### 3.2.6 Notification Service (Kafka Orchestration)
**Priority:** M1 (Basic) | **Owner:** CTO (Go + Kafka - 80% Ready from IoT Broker)

| Feature | Complexity | Estimated Hours | Dependencies | Priority |
|---------|-----------|----------------|--------------|----------|
| Kafka Setup & Configuration | Medium | 24 hrs | IoT Broker (80% ready) | M1 |
| Email Integration (SMTP/SendGrid) | Medium | 24 hrs | Kafka | M1 |
| SMS Integration (Local provider) | Medium | 24 hrs | Kafka | M1 |
| Basic Template System | Low | 24 hrs | Email/SMS | M1 |
| Notification Queue | Medium | 16 hrs | Kafka (reuse broker) | M1 |
| Policy Purchase Event Triggers | Low | 16 hrs | Policy service | M1 |
| API Documentation | Low | 8 hrs | Core features | M1 |
| Unit & Integration Tests | - | 24 hrs | Core features | M1 |
| **SUBTOTAL M1** | - | **160 hrs** | - | - |
| **Buffer (10%)** | - | **16 hrs** | - | - |
| **TOTAL M1** | - | **176 hrs** | - | - |

**Moved to M2 (255 hrs):**
- Push Notification (FCM) - 40 hrs (Mobile apps are M2)
- Advanced Template Management - 16 hrs
- Notification Preferences - 24 hrs
- Notification History & Tracking - 24 hrs
- Multi-language Support - 32 hrs (Bengali for M2)
- Advanced Event Triggers (Claims, Renewals) - 32 hrs
- Additional docs & tests - 64 hrs
- Buffer - 23 hrs

**Team Assignment:** CTO (40% time = 19 hrs/week) - Can complete in ~9 weeks alongside other tasks
**Savings:** ~200 hours from existing IoT Broker/Kafka infrastructure

---

#### 3.2.7 Ticketing/Customer Service
**Priority:** M2 (NOT M1!) | **Owner:** Mamoon + Sujon Ahmed (Node.js)

**⚠️ MOVED TO M2** - Customer support tickets not critical for beta demo
- M1 focus: Policy purchase flow demonstration
- Support can be handled manually via phone/email in beta
- Full ticketing system after Grand Launch (April 14)

| Feature | Complexity | Estimated Hours | Dependencies | Priority |
|---------|-----------|----------------|--------------|----------|
| Ticket Management System | Medium | 48 hrs | User Service (existing) | M2 |
| Ticket Status Workflow | Medium | 32 hrs | Ticket CRUD | M2 |
| Customer Communication | Medium | 32 hrs | Notification Service | M2 |
| FAQ Management (Bengali/English) | Medium | 32 hrs | CMS-style | M2 |
| Knowledge Base | Low | 24 hrs | Document linking | M2 |
| Support Agent Assignment | Medium | 32 hrs | User Service (existing) | M2 |
| Ticket Priority & Escalation | Medium | 32 hrs | Workflow | M3 |
| SLA Tracking | Medium | 24 hrs | Time-based rules | M3 |
| Insurance Query Templates | Low | 16 hrs | Common queries | M2 |
| API Documentation | Low | 16 hrs | All features | M2 |
| Unit & Integration Tests | - | 32 hrs | All features | M2 |
| **SUBTOTAL** | - | **320 hrs** | - | - |
| **Buffer (10%)** | - | **32 hrs** | - | - |
| **TOTAL M2** | - | **352 hrs** | - | - |

**Team Assignment:** Mamoon + Sujon Ahmed in M2 phase (March-April)
**Note:** Basic contact form sufficient for M1 beta

---

#### 3.2.8 Web Admin Portals (Multiple Portals - React/Next.js)
**Priority:** M1 (2 Portals), M2 (1 Portal) | **Owner:** React Dev (starts Dec 18)

**REVISED SCOPE:** Minimal portals per phase
- **M1 (March 1):** Business Admin + Shared Infrastructure only
- **M2 (April 14):** Partner Portal
- **M3 (June+):** All other 10 portals

| Portal/Feature | Complexity | Estimated Hours | Priority | Notes |
|----------------|-----------|----------------|----------|-------|
| **M1 - Shared Components & Infrastructure** | | | | |
| Design System Setup (Tailwind/Shadcn) | Medium | 40 hrs | M1 | Reusable foundation |
| Authentication UI (Shared) | Low | 24 hrs | M1 | One time setup |
| Navigation & Layout System | Medium | 32 hrs | M1 | Reusable template |
| State Management (Redux/Zustand) | Medium | 24 hrs | M1 | Simplified |
| API Integration Layer | Medium | 32 hrs | M1 | Core only |
| **Business Admin Portal** | High | 72 hrs | M1 | Product/Policy mgmt CRUD |
| Notification Center (Basic) | Low | 16 hrs | M1 | Simple alerts |
| Responsive Design | Medium | 24 hrs | M1 | Mobile friendly |
| Error Handling & Validation | Low | 16 hrs | M1 | Core features |
| Testing (Unit + E2E) | - | 32 hrs | M1 | Critical paths |
| **SUBTOTAL M1** | - | **312 hrs** | - | - |
| **Buffer (10%)** | - | **31 hrs** | - | - |
| **TOTAL M1** | - | **343 hrs** | - | - |
| | | | | |
| **M2 - Partner Portal** | | | | |
| Partner Portal Dashboard | Medium | 48 hrs | M2 | View policies/commissions |
| Partner Profile Management | Low | 24 hrs | M2 | KYB info |
| Performance Analytics (Basic) | Medium | 32 hrs | M2 | Charts/reports |
| Testing | - | 16 hrs | M2 | Integration tests |
| **SUBTOTAL M2** | - | **120 hrs** | - | - |
| **Buffer (10%)** | - | **12 hrs** | - | - |
| **TOTAL M2** | - | **132 hrs** | - | - |

**Moved to M3 (10 portals = ~800 hrs):**
- System Admin Portal - 48 hrs
- Agent Portal - 80 hrs
- Customer Support Portal - 72 hrs
- DevOps Portal - 64 hrs
- Database Manager Portal - 64 hrs
- Focal Person Portal - 80 hrs
- Partner Admin Portal - 72 hrs
- General Staff Portal - 64 hrs
- Vendor Portal - 64 hrs
- Marketing Admin Portal - 56 hrs
- Advanced features - 136 hrs

**Team Assignment:** 
- React Dev (Dec 18-Mar 1): M1 portals = 343 hrs ÷ 48 hrs/week = 7.2 weeks ✓ Achievable
- React Dev (Mar 1-Apr 14): M2 Partner Portal = 132 hrs easily fits in 6 weeks

---

#### 3.2.9 Mobile Apps (Android & iOS - Native)
**Priority:** M1 (WITH MOCK SERVERS) | **Owner:** Nur Hossain (Android) + Sojol Ahmed (iOS)

**✅ M1 STRATEGY - PARALLEL DEVELOPMENT:**
- Backend team provides mock server (Sujon Ahmed - 40 hrs)
- Mobile apps develop against mock APIs in M1 (Customer App ONLY)
- Real API integration happens in M2
- Both teams work in parallel - no blocking
- **NO Agent Mobile App** - Per SRS, agents use web portal on tablets (saves 580 hrs)

**Note:** Native development (Kotlin for Android, Swift for iOS)

| Feature | Complexity | Android Hours | iOS Hours | Total Hours | Phase |
|---------|-----------|---------------|-----------|-------------|-------|
| Project Setup & Architecture | Low | 16 hrs | 16 hrs | 32 hrs | M1 |
| Authentication Flow (Mock API) | Low | 24 hrs | 24 hrs | 48 hrs | M1 |
| **Customer App** | | | | | |
| Dashboard (Mock data) | Medium | 32 hrs | 32 hrs | 64 hrs | M1 |
| Policy View & Details (Mock) | Medium | 40 hrs | 40 hrs | 80 hrs | M1 |
| Claim Filing (UI only, Mock) | High | 56 hrs | 56 hrs | 112 hrs | M1 |
| Payment Integration (Mock) | High | 48 hrs | 48 hrs | 96 hrs | M1 |
| Document Upload & Camera | Medium | 32 hrs | 32 hrs | 64 hrs | M1 |
| Push Notifications (Setup) | Low | 24 hrs | 24 hrs | 48 hrs | M1 |
| Profile Management (Mock) | Low | 16 hrs | 16 hrs | 32 hrs | M1 |
| **SUBTOTAL M1** | - | **288 hrs** | **288 hrs** | **576 hrs** | - |
| **Buffer (10%)** | - | **29 hrs** | **29 hrs** | **58 hrs** | - |
| **TOTAL M1** | - | **317 hrs** | **317 hrs** | **634 hrs** | - |

**M2 Integration (Real APIs):**
| Activity | Android | iOS | Total | Phase |
|----------|---------|-----|-------|-------|
| Real API Integration | 80 hrs | 80 hrs | 160 hrs | M2 |
| Bkash SDK Real Integration | 40 hrs | 40 hrs | 80 hrs | M2 |
| End-to-end Testing | 40 hrs | 40 hrs | 80 hrs | M2 |
| Bug Fixes | 32 hrs | 32 hrs | 64 hrs | M2 |
| App Store Submission | 16 hrs | 16 hrs | 32 hrs | M2 |
| **M2 SUBTOTAL** | **208 hrs** | **208 hrs** | **416 hrs** | - |
| **Buffer (10%)** | **21 hrs** | **21 hrs** | **42 hrs** | - |
| **TOTAL M2** | **229 hrs** | **229 hrs** | **458 hrs** | - |

**Team Assignment:** 
- **M1 (Dec 20 - Mar 1):** Nur & Sojol build full customer app UI with mock APIs (634 hrs)
- **M2 (Mar 2 - Apr 14):** Real API integration + testing + app store (458 hrs)
- **Timeline:** M1 = 12 weeks @ 48 hrs/week each = 576 hrs available ✅
- **No blocking:** Backend and Mobile develop in parallel

---

#### 3.2.10 Infrastructure & DevOps
**Priority:** M1 (Essential only) | **Owner:** Sagor (DevOps) + CTO Support

| Feature | Complexity | Estimated Hours | Status | Priority |
|---------|-----------|----------------|--------|----------|
| Cloud Infrastructure Setup (AWS/Azure) | Medium | 32 hrs | CTO patterns exist | M1 |
| Containerization (Docker) | Low | 16 hrs | Experience exists | M1 |
| Database Setup (PostgreSQL) | Medium | 24 hrs | DBManager ready | M1 |
| API Gateway Deployment | Medium | 24 hrs | 50% ready code | M1 |
| Basic CI/CD Pipeline (GitHub Actions) | Medium | 32 hrs | Simple deploy | M1 |
| SSL/Security Configuration | High | 32 hrs | Critical | M1 |
| Basic Monitoring (Logs) | Low | 16 hrs | CloudWatch/Basic | M1 |
| Backup Strategy | Medium | 16 hrs | Database snapshots | M1 |
| Documentation | Low | 16 hrs | Essential setup | M1 |
| **SUBTOTAL M1** | - | **208 hrs** | - | - |
| **Buffer (10%)** | - | **21 hrs** | - | - |
| **TOTAL M1** | - | **229 hrs** | - | - |

**Moved to M2 (264 hrs):**
- Container Orchestration (K8s) - 56 hrs (Not needed for beta, Docker Compose sufficient)
- Advanced CI/CD - 16 hrs
- Kafka Setup - 40 hrs (Handled by CTO with notification service)
- Monitoring & Logging (Prometheus/Grafana) - 48 hrs (Full observability stack)
- Advanced Security (WAF) - 40 hrs
- Load Testing & Optimization - 32 hrs
- Advanced documentation - 8 hrs
- Buffer - 24 hrs

**Team Assignment:** Sagor (100% - 48 hrs/week) = 229 hrs ÷ 48 = ~4.8 weeks
**Note:** Focus on getting services deployed and running, advanced DevOps in M2

---

#### 3.2.11 QA & Testing
**Priority:** M1 (Essential testing only) | **Owner:** QA Team

| Activity | Estimated Hours | Dependencies | Priority |
|----------|----------------|--------------|----------|
| Test Plan Development | 24 hrs | Requirements | M1 |
| Test Case Creation (Critical paths) | 48 hrs | Core Features | M1 |
| Manual Testing (Functional - Core) | 80 hrs | Dev Complete | M1 |
| API Testing (Postman - Core APIs) | 48 hrs | Backend Services | M1 |
| UI Testing (Business Admin portal) | 40 hrs | Frontend | M1 |
| Integration Testing (Critical flows) | 64 hrs | Services integration | M1 |
| Basic Security Testing | 24 hrs | Auth/Payment | M1 |
| Bug Reporting & Tracking | 32 hrs | Testing | M1 |
| **SUBTOTAL M1** | **360 hrs** | - | - |
| **Buffer (15%)** | **54 hrs** | - | - |
| **TOTAL M1** | **414 hrs** | - | - |

**Moved to M2/M3 (613 hrs):**
- Comprehensive Test Case Creation - 32 hrs
- Extended Manual Testing - 80 hrs
- Advanced API Testing - 32 hrs
- Mobile UI/UX Testing - 40 hrs
- Comprehensive Integration Testing - 56 hrs
- Performance Testing - 64 hrs
- Advanced Security Testing - 24 hrs
- Regression Testing - 80 hrs
- Test Automation Setup - 64 hrs
- Additional testing activities - 24 hrs
- Buffer - 117 hrs

**Team Assignment:** 1 QA × 10 weeks (48 hrs/week) = 480 hrs (sufficient for M1 with buffer)
**Note:** Focus on critical user journeys: Registration → Policy Purchase → Payment

---

#### 3.2.12 UI/UX Design
**Priority:** M1 (Core screens only) | **Owner:** Rumon (UI/UX Designer)

| Activity | Estimated Hours | Dependencies | Priority |
|----------|----------------|--------------|----------|
| Design System Creation (Tailwind-based) | 32 hrs | Brand Guidelines | M1 |
| User Research & Personas (Basic) | 16 hrs | Requirements | M1 |
| Information Architecture (Core flows) | 16 hrs | Requirements | M1 |
| Wireframes (Business Admin portal) | 32 hrs | IA | M1 |
| High-Fidelity Mockups (Business Admin) | 48 hrs | Wireframes | M1 |
| Prototyping (Critical flows) | 24 hrs | Mockups | M1 |
| Design Handoff & Documentation | 16 hrs | Designs | M1 |
| Design Reviews & Iterations | 24 hrs | Feedback | M1 |
| **SUBTOTAL M1** | **208 hrs** | - | - |
| **Buffer (15%)** | **31 hrs** | - | - |
| **TOTAL M1** | **239 hrs** | - | - |

**Moved to M2/M3 (347 hrs):**
- Additional User Research - 16 hrs
- Extended IA (all portals) - 16 hrs
- Wireframes (Partner Portal, Mobile Apps) - 48 hrs
- High-Fidelity Mockups (Web - all portals) - 48 hrs
- High-Fidelity Mockups (Mobile apps) - 80 hrs
- Advanced Prototyping - 24 hrs
- Additional documentation - 16 hrs
- Design iterations - 16 hrs
- Buffer - 67 hrs
- Customer Portal designs - 16 hrs

**Team Assignment:** Rumon (100% - 48 hrs/week) = 239 hrs ÷ 48 = ~5 weeks
**Note:** Focus on Business Admin portal only for M1, Partner portal designs in M2

---

### 3.3 M1 Summary - Total Effort (REVISED WITH PARALLEL DEVELOPMENT)

**M1 STRATEGY:** All teams work in parallel using mock servers. Backend provides mock APIs, Frontend/Mobile develop against mocks. Integration in M2.

| Component | Estimated Hours | Team Assignment | Development Mode |
|-----------|----------------|----------------|------------------|
| User Service (Auth/AuthZ) | 44 hrs | CTO (50% time) | ✅ 100% ready - Deploy early |
| Policy Service (Core only) | 317 hrs | Delowar + C# Dev | Real development |
| Payment Service (Bkash) | 106 hrs | Mamoon (40%) | Real development |
| Document Service | 66 hrs | CTO (50%) + Sagor | ✅ 100% ready - Deploy early |
| Notification Service (Basic) | 176 hrs | CTO (50%) | Real development |
| **Backend Mock Server** | 40 hrs | Sujon Ahmed | Mock APIs for frontend/mobile |
| **Business Admin Portal** | 343 hrs | React Dev | Develop with mocks |
| **Partner Portal** | 132 hrs | React Dev | Develop with mocks |
| **Customer Mobile App (Android)** | 317 hrs | Nur Hossain (90%) | Develop with mocks |
| **Customer Mobile App (iOS)** | 317 hrs | Sojol Ahmed (90%) | Develop with mocks |
| DevOps & Infrastructure | 229 hrs | Sagor (50%) | Real infrastructure |
| QA & Testing (Mock-based) | 414 hrs | 1 QA (60%) | Test with mocks |
| UI/UX Design (M1 screens) | 239 hrs | Rumon (60%) | Business Admin + Partner + Mobile |
| **TOTAL M1 EFFORT** | **2,740 hrs** | - | - |

**Available Capacity M1:** 4,040 hrs (corrected from 3,992)
**Required Effort M1:** 2,740 hrs  
**Utilization:** 68% ✅ **COMFORTABLE**
**Buffer Remaining:** 1,300 hrs (32%)

**M1 Deliverables:**
- ✅ All backend services complete (real)
- ✅ Business Admin Portal (with mocks)
- ✅ Partner Portal (with mocks)
- ✅ Customer Mobile App - Android + iOS (with mocks)
- ✅ Infrastructure deployed

**M2 Focus (Critical):** 
- Real API Integration for Customer Mobile App
- Frontend Real API Integration (Business Admin + Partner)
- Claims Service backend
- Push Notifications (FCM)
- Commission tracking
- Testing + Bug fixes

---

### 3.4 M2 Services - Desirable Features

#### 3.4.1 Analytics Service
| Feature | Estimated Hours |
|---------|----------------|
| Data Collection & Storage | 48 hrs |
| Dashboard Builder | 64 hrs |
| Standard Reports | 80 hrs |
| Custom Report Builder | 72 hrs |
| Data Visualization | 56 hrs |
| Export Functionality | 32 hrs |
| Testing & Buffer | 88 hrs |
| **TOTAL** | **440 hrs** |

#### 3.4.2 Commission Service
| Feature | Estimated Hours |
|---------|----------------|
| Commission Structure Models | 40 hrs |
| Commission Calculation Engine | 72 hrs |
| Agent Performance Tracking | 48 hrs |
| Commission Reporting | 56 hrs |
| Payment Integration | 40 hrs |
| Testing & Buffer | 68 hrs |
| **TOTAL** | **324 hrs** |

#### 3.4.3 Integration Service
| Feature | Estimated Hours |
|---------|----------------|
| Third-party API Framework | 48 hrs |
| Reinsurance Integration | 64 hrs |
| Medical Provider Integration | 64 hrs |
| Government Portal Integration | 72 hrs |
| Webhook Management | 40 hrs |
| Testing & Buffer | 76 hrs |
| **TOTAL** | **364 hrs** |

### M2 Summary
| Component | Estimated Hours |
|-----------|----------------|
| Analytics Service | 440 hrs |
| Commission Service | 324 hrs |
| Integration Service | 364 hrs |
| Enhanced Mobile Features | 280 hrs |
| Advanced Customer Service | 240 hrs |
| QA & Testing | 420 hrs |
| **TOTAL M2 EFFORT** | **2,068 hrs** |

**Available Capacity for M2:** 2,622 hrs ✓ Sufficient (Tight)

---

### 3.5 M3 Services - Should Have & Future

#### M3 Summary (High-Level)
| Category | Estimated Hours |
|----------|----------------|
| Advanced Analytics & AI/ML | 800 hrs |
| Complete Integration Ecosystem | 640 hrs |
| Performance Optimizations | 480 hrs |
| Advanced Security Features | 560 hrs |
| Additional Features | 720 hrs |
| QA & Testing | 640 hrs |
| Documentation & Training | 320 hrs |
| **TOTAL M3 EFFORT** | **4,160 hrs** |

**Available Capacity for M3:** 6,128 hrs ✓ Sufficient with buffer

---

### 3.6 Overall Project Summary (FINAL - REALISTIC CAPACITY)

| Phase | Available Hours | Required Hours | Utilization | Status |
|-------|----------------|----------------|-------------|---------|
| **M1 (Mar 1 Beta)** | 4,040 hrs | 2,740 hrs | 68% | ✅ **COMFORTABLE** |
| **M2 (Apr 14 Launch)** | 2,256 hrs | 1,641 hrs | 73% | ✅ **COMFORTABLE** |
| **M3 (Aug 1 Complete)** | 6,768 hrs | 4,863 hrs | 72% | ✅ **COMFORTABLE** |
| **TOTAL PROJECT** | **13,064 hrs** | **9,244 hrs** | **71%** | ✅ **ACHIEVABLE** |

**Capacity Corrections Applied:**
- M1: 3,992 → 4,040 hrs (+48 hrs from Delowar correction)
- M2: 2,256 hrs (unchanged)
- M3: 6,768 hrs (unchanged)
- **Total: 13,016 → 13,064 hrs**

**Aligned Across Documents:**
- ✅ 03_TeamCapacity.md: 13,064 hrs available (corrected)
- ✅ 04_EffortEstimation.md: 9,244 hrs required  
- ✅ 05_SprintPlanning.md: Sprint hours reduced to match capacity
- ✅ 09_RequirementsByMilestone.md: Components aligned
- ✅ Team allocation: CTO 50%, Mamoon 40%, Sagor 50%, Python 50%, Sujon 60%, Rumon 60%, React 70%, QA 60%, Others 90%

**M1 Components (2,740 hrs) - PARALLEL DEVELOPMENT:**
- User Service: 44 hrs (100% ready)
- Policy Service: 317 hrs (new C#)
- Payment Service: 106 hrs (70% ready)
- Document Service: 66 hrs (100% ready)
- Notification: 176 hrs (80% ready)
- Mock Server: 40 hrs (enables parallel work)
- Business Admin Portal: 343 hrs (with mocks)
- Partner Portal: 132 hrs (with mocks)
- Customer Mobile Apps: 634 hrs (with mocks, Customer App only)
- DevOps: 229 hrs
- QA: 414 hrs (mock-based testing)
- UI/UX: 239 hrs (M1 screens)
- **Python Devs M1 (192 hrs each = 384 hrs):** Data pipeline setup, API testing automation, AI infrastructure prep (NOT counted in 2,740 - support role)

**M2 Components (1,641 hrs) - INTEGRATION PHASE:**
- Claims Service: 634 hrs (new backend)
- Customer Mobile Real API Integration: 320 hrs (connect to real APIs)
- Frontend Real API Integration: 200 hrs (connect Business Admin + Partner)
- Push Notifications: 80 hrs (FCM production)
- Commission tracking: 120 hrs (partner feature)
- End-to-End Testing: 200 hrs (full integration)
- Bug Fixes & Polish: 87 hrs
**TOTAL: 1,641 hrs** ✅

**M3 Components (4,863 hrs):**
- IoT Integration: 757 hrs
- AI Engine: 946 hrs
- 10 Remaining Portals: 780 hrs
- Analytics Service: 440 hrs
- Commission Service: 324 hrs
- Integration Service: 364 hrs
- Performance: 320 hrs
- Security: 280 hrs
- QA: 640 hrs
- Docs: 320 hrs

**M2 Now Comfortable (69% utilization):**
- Mock-based development in M1 solves the capacity crisis
- M2 focuses only on integration, testing, and claims backend
- No overtime needed ✅
- Launch on time for Pohela Boishakh (April 14, 2026) ✅

---
