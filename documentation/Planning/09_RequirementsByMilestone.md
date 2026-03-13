# Requirements Distribution by Milestone

## Overview

**Total Requirements in SRS V3.7:** 345 functional requirements

| Priority | Count | Percentage | Target Milestone |
|----------|-------|------------|------------------|
| **M1** | 103 | 29.9% | Beta Launch - March 1, 2026 |
| **M2** | 92 | 26.7% | Grand Launch - April 14, 2026 |
| **M3** | 81 | 23.5% | Complete Platform - August 1, 2026 |
| **D (Desirable)** | 53 | 15.4% | Post M3 / Future phases |
| **S (Should Have)** | 1 | 0.3% | As needed |
| **F (Future)** | 15 | 4.3% | Future releases |

---

## M1 - Beta Launch (March 1, 2026) - 103 Requirements

### Strategy: **Barebone Functional Platform**
- **Goal:** Prove the concept with minimal viable features
- **Target Users:** Early adopters, limited partner base (2-3 partners)
- **Scope:** End-to-end flow for 1 insurance product type

### M1 Core Features (Implemented)

#### Authentication & Authorization (FG-001, FG-002) - 16 Requirements
| FR ID | Requirement | Implementation Notes |
|-------|-------------|---------------------|
| FR-001 to FR-009 | Core authentication (mobile, email, password, OTP, session) | ✅ Using CTO's existing Go services (100% ready) |
| FR-014 to FR-018 | RBAC, ABAC, ACL | ✅ Using CTO's existing Authorization service |

**Effort Saved:** ~500 hours (existing code reuse)

---

#### Product & Policy Management (FG-003, FG-004) - 20 Requirements
| FR ID | Requirement | Scope for M1 |
|-------|-------------|--------------|
| FR-021, FR-022 | Product catalog with search | Single product category (Health) |
| FR-026 | Product CRUD by Business Admin | Basic operations only |
| FR-030 to FR-034 | End-to-end policy purchase flow | Simplified flow, single nominee |
| FR-039, FR-040 | Policy status tracking and dashboard | Basic statuses only |
| FR-091, FR-092 | Policy document and lifecycle tracking | PDF generation, basic audit |

**Services Required:**
- Insurance Engine (C# .NET) - Mr. Delowar's team - 748 hrs
- Basic product catalog (10 products max)

---

#### Claims Management (FG-008) - **MOVED TO M2**
| FR ID | Requirement | Scope for M2 |
|-------|-------------|--------------|
| FR-041 to FR-043 | Claim submission and validation | Full workflow with approval |
| FR-044 to FR-058 | Claim tracking, approval, fraud detection | Complete claims system |
| FR-218 to FR-219 | Claim status state machine | Full state machine |

**Services Required (M2):**
- Claim Service (C# .NET) - 634 hrs (full implementation)
- Tiered approval workflow with Business Admin + Focal Person
- Fraud detection and document verification

**M1 Decision:** Claims REMOVED from M1 scope due to capacity constraints. M1 focuses on policy purchase demo only. Claims processing begins M2 (April 14, 2026).

---

#### Payment Processing (FG-007) - 8 Requirements
| FR ID | Requirement | Scope for M1 |
|-------|-------------|--------------|
| FR-073, FR-074 | bKash payment gateway | ✅ Mamoon's existing code (70% ready) |
| FR-076, FR-077 | Manual payment with verification | Basic workflow |
| FR-083 | Payment audit trail | PostgreSQL logging |

**Effort Saved:** ~357 hours (Mamoon's existing payment service)

---

#### Document Management (FG-005) - 6 Requirements
| FR ID | Requirement | Scope for M1 |
|-------|-------------|--------------|
| FR-232, FR-237 | File upload and S3 storage | ✅ CTO's Storage Manager (100% ready) |
| FR-239 | Upload policies and compression | Client-side validation |

**Effort Saved:** ~333 hours (CTO's existing Storage Manager)

---

#### Notifications (FG-012) - 8 Requirements
| FR ID | Requirement | Scope for M1 |
|-------|-------------|--------------|
| FR-114, FR-115 | Kafka-based notifications (SMS, email, push) | Basic templates only |
| FR-161, FR-162 | SMS and email integration | Essential notifications only |

**Services Required:**
- Notification Service (Go + Kafka) - CTO - 431 hrs

---

#### Partner Management (FG-009) - 10 Requirements (REDUCED SCOPE)
| FR ID | Requirement | Scope for M1 |
|-------|-------------|--------------|
| FR-066 | Focal Person portal for partner management | ⚠️ **DEFERRED to M2** |
| FR-192 | RACI for monitoring and escalation | Manual process in M1 |

**M1 Decision:** Manual partner onboarding, no automated portal

---

#### Customer Support (FG-011) - 6 Requirements
| FR ID | Requirement | Scope for M1 |
|-------|-------------|--------------|
| FR-106 | FAQ section with searchable knowledge base | Static FAQ page |
| FR-108, FR-109 | Ticketing system | Basic support ticketing |

**Services Required:**
- Ticketing Service (Node.js) - Mamoon + Sujon - 352 hrs

---

#### Admin & Reporting (FG-017) - 8 Requirements
| FR ID | Requirement | Scope for M1 |
|-------|-------------|--------------|
| FR-131, FR-132 | Admin dashboards with 2FA | System Admin + Partner portal only |
| FR-134, FR-135 | Product and claims management | Basic CRUD operations |
| FR-140 | System health monitoring | Basic Prometheus/Grafana |

---

#### Audit & Logging (FG-019) - 5 Requirements
| FR ID | Requirement | Scope for M1 |
|-------|-------------|--------------|
| FR-153, FR-154 | Immutable audit logs with 20-year retention | PostgreSQL + S3 archival |
| FR-083 | Payment audit trail | Integrated with payment service |

---

#### Security & Compliance (SEC) - 4 Requirements
| FR ID | Requirement | Scope for M1 |
|-------|-------------|--------------|
| SEC-001, SEC-002 | Zero Trust architecture, encryption at rest/transit | TLS 1.3, AES-256 |
| SEC-003 | PCI-DSS SAQ-A compliance for payments | bKash hosted payment page |
| FR-192 | RACI for monitoring | Basic monitoring setup |

---

### M1 Web Portals - **2 PORTALS ONLY (ALIGNED WITH 04_EffortEstimation.md)**

**DECISION:** Business Admin + Partner portals for M1 (matches 04_EffortEstimation.md line 233-256)

#### Portal 1: Business Admin Portal (Priority 1)
**Purpose:** Business managers manage products, policies, and core operations

| Feature | Complexity | Hours |
|---------|-----------|-------|
| Design System Setup (Tailwind/Shadcn) | Medium | 40 hrs |
| Authentication UI (Shared) | Low | 24 hrs |
| Navigation & Layout System | Medium | 32 hrs |
| State Management (Redux/Zustand) | Medium | 24 hrs |
| API Integration Layer | Medium | 32 hrs |
| Business Admin Portal (Product/Policy CRUD) | High | 72 hrs |
| Notification Center (Basic) | Low | 16 hrs |
| Responsive Design | Medium | 24 hrs |
| Error Handling & Validation | Low | 16 hrs |
| Testing (Unit + E2E) | - | 32 hrs |
| **SUBTOTAL** | - | **312 hrs** |
| **Buffer (10%)** | - | **31 hrs** |
| **TOTAL** | - | **343 hrs** |

---

#### Portal 2: Partner Portal (Priority 2)
**Purpose:** Insurance partners track leads, commissions, and performance

| Feature | Complexity | Hours |
|---------|-----------|-------|
| Partner Portal Dashboard | Medium | 48 hrs |
| Partner Profile Management | Low | 24 hrs |
| Performance Analytics (Basic) | Medium | 32 hrs |
| Commission Tracking | Low | 16 hrs |
| Testing | - | 12 hrs |
| **SUBTOTAL** | - | **132 hrs** |
| **Buffer (10%)** | - | **0 hrs** (included) |
| **TOTAL** | - | **132 hrs** |

---

#### **DEFERRED TO M2:**
- ❌ System Admin Portal (moved to M2)
- ❌ Agent Portal (moved to M2)
- ❌ Customer Support Portal (moved to M2)

#### **DEFERRED TO M3:**
- ❌ Focal Person Portal (SRS says M1, but capacity insufficient - manual process for M1)
- ❌ DevOps Portal
- ❌ Database Manager Portal
- ❌ 5 other portals

**Total M1 Portal Hours:** 475 hrs (343 + 132)
**Aligned with 04_EffortEstimation.md:** ✅ YES

---

### M1 Mobile Apps - **BASIC FEATURES ONLY**

#### Customer App (Android + iOS)
| Feature | Android | iOS | Total |
|---------|---------|-----|-------|
| Authentication | 20 hrs | 20 hrs | 40 hrs |
| Dashboard (policy list) | 24 hrs | 24 hrs | 48 hrs |
| Policy Purchase Flow | 40 hrs | 40 hrs | 80 hrs |
| Policy Details View | 24 hrs | 24 hrs | 48 hrs |
| Claim Filing (basic) | 40 hrs | 40 hrs | 80 hrs |
| Payment (bKash) | 32 hrs | 32 hrs | 64 hrs |
| Document Upload | 24 hrs | 24 hrs | 48 hrs |
| Push Notifications | 16 hrs | 16 hrs | 32 hrs |
| **SUBTOTAL** | **220 hrs** | **220 hrs** | **440 hrs** |
| **Buffer (10%)** | **22 hrs** | **22 hrs** | **44 hrs** |
| **TOTAL** | **242 hrs** | **242 hrs** | **484 hrs** |

---

#### Agent App (Android + iOS) - **DEFERRED TO M2**
- Agent dashboard, customer management, lead tracking moved to M2
- Agents can use Partner Portal on tablets for M1

**Total M1 Mobile Hours:** 484 hrs (down from 1,162 hrs)
**Savings:** 678 hours moved to M2

---

### M1 Infrastructure & DevOps - **ESSENTIAL ONLY**

| Component | Hours | Notes |
|-----------|-------|-------|
| Cloud Infrastructure (AWS/Azure) | 32 hrs | Reuse CTO's patterns |
| Docker Containerization | 20 hrs | Basic setup |
| Kubernetes Setup | 48 hrs | Simple cluster |
| CI/CD Pipeline | 40 hrs | GitHub Actions |
| Database Setup (PostgreSQL) | 24 hrs | Single primary |
| Kafka Setup | 32 hrs | For notifications |
| API Gateway Deployment | 24 hrs | 50% existing code |
| Monitoring (Prometheus/Grafana) | 40 hrs | Basic dashboards |
| Security (SSL, WAF) | 32 hrs | Cloudflare |
| Backup & DR | 24 hrs | Daily backups |
| **SUBTOTAL** | **316 hrs** | |
| **Buffer (10%)** | **32 hrs** | |
| **TOTAL** | **348 hrs** | Down from 493 hrs |

**Savings:** 145 hours (reduced complexity)

---

### M1 Effort Summary - FINAL ALIGNED

| Component | Hours | Team | Notes |
|-----------|-------|------|-------|
| User Service (Auth/AuthZ) | 44 hrs | CTO (50%) | ✅ 100% ready |
| Policy Service (Core) | 317 hrs | Delowar + C# Dev (90%) | New C# service |
| Payment Service | 106 hrs | Mamoon (40%) | 70% ready (Bkash) |
| Document Service | 66 hrs | CTO (50%) + Sagor | ✅ 100% ready |
| Notification Service | 176 hrs | CTO (50%) | 80% ready (IoT Broker) |
| Mock Server | 40 hrs | Sujon Ahmed (60%) | For parallel development |
| Business Admin Portal | 343 hrs | React Dev (70%) | With mocks |
| Partner Portal | 132 hrs | React Dev (70%) | With mocks |
| Customer Mobile App (Android) | 317 hrs | Nur Hossain (90%) | With mocks |
| Customer Mobile App (iOS) | 317 hrs | Sojol Ahmed (90%) | With mocks |
| DevOps & Infrastructure | 229 hrs | Sagor (50%) | Basic setup |
| QA & Testing | 414 hrs | 1 QA (60%) | Critical paths with mocks |
| UI/UX Design | 239 hrs | Rumon (60%) | Business Admin + Partner + Mobile |
| **TOTAL M1** | **2,740 hrs** | - | **✅ Matches capacity** |

**M1 Available Capacity:** 4,040 hours (corrected: +48 hrs from Delowar fix)
**M1 Required Effort:** 2,740 hours
**Utilization:** 68% ✅ **COMFORTABLE**

**❌ REMOVED FROM M1:**
- Claims Service → M2 (634 hrs)
- Ticketing Service → M2 (352 hrs)
- 10 Other Portals → M3 (800 hrs)
- Agent Mobile Apps → NOT IN SCOPE (per SRS, agents use web portal)

---

## M2 - Grand Launch (April 14, 2026) - 92 Requirements

### Strategy: **Partner Onboarding & Business Operations**
- **Goal:** Scale to 10-15 partners, add business intelligence
- **Target:** Full partner ecosystem, agent network activated

### M2 Additional Features

#### Partner Portal Enhancements (M2)
- Partner onboarding workflow (FR-059 to FR-062)
- KYB verification integration
- Commission calculation engine
- Partner analytics dashboard
- API integration for embedded insurance (FR-064)

**Effort:** 320 hrs

---

#### New Portals (M2)

**Portal 3: Business Admin Portal** - 280 hrs
- Business operations dashboard
- Policy approvals and overrides
- Financial reconciliation
- KPI tracking
- Executive reports

**Portal 4: Agent Portal** - 240 hrs  
- Agent dashboard
- Customer management
- Lead tracking
- Commission tracking
- Mobile-optimized

**Total M2 Portal Hours:** 520 hrs

---

#### Agent Mobile Apps (M2)

**Agent App (Android + iOS)** - 420 hrs
- Agent dashboard
- Customer onboarding
- Policy management
- Lead tracking
- Commission view

---

#### Enhanced Features (M2)

**Analytics & Reporting Service (C#)** - 440 hrs
- Business intelligence dashboards
- Predictive analytics
- Customer segmentation
- Geographic analytics (FR-152, FR-202)

**Commission Service (C#)** - 324 hrs
- Multi-level commission calculation
- Agent hierarchy management
- Commission payout processing
- Performance tracking

**Integration Service (C#)** - 364 hrs
- Third-party API framework
- Insurer integrations
- EHR/Hospital integrations (FR-229, FR-230)
- Webhook management

---

### M2 Effort Summary - FINAL ALIGNED

| Component | Hours | Team | Notes |
|-----------|-------|------|-------|
| Claims Service (Full) | 634 hrs | Delowar + C# Dev (90%) | New backend service |
| Customer Mobile Real API Integration | 320 hrs | Nur + Sojol (90%) | Connect to real APIs |
| Frontend Real API Integration | 200 hrs | React Dev (70%) | Business Admin + Partner |
| Push Notifications | 80 hrs | CTO (50%) + Mobile devs | FCM production |
| Commission tracking | 120 hrs | Mamoon (40%) + Sujon (60%) | Partner feature |
| End-to-End Testing | 200 hrs | 1 QA (60%) | Full integration |
| Bug Fixes & Polish | 87 hrs | All teams | Final touches |
| **TOTAL M2** | **1,641 hrs** | - | **✅ Matches exactly** |

**M2 Available Capacity:** 2,256 hours (5 weeks × 14 team members)
**M2 Required Effort:** 1,641 hours
**Status:** ✅ **73% utilization - COMFORTABLE**

**Note:** M2 focuses on Real API integration for Customer Mobile App + Claims backend + Push Notifications.

---

## M3 - Complete Platform (August 1, 2026) - 81 Requirements

### Strategy: **IoT, AI, and Advanced Features**
- **Goal:** Complete platform with all advanced capabilities
- **Focus:** IoT devices, AI agents, remaining portals, advanced analytics

### M3 Major Features

#### IoT Integration & UBI (FG-013) - **8 M3 + 5 Future + 1 Desirable = 14 IoT Requirements**

**IoT Infrastructure (Go)** - CTO + Sagor
| Feature | Hours | Notes |
|---------|-------|-------|
| IoT Device Registration & Provisioning (FR-179) | 80 hrs | X.509 certificates, device lifecycle |
| MQTT Broker Setup (FR-180) | 120 hrs | Handle 10,000 devices, TimescaleDB |
| IoT Telemetry Processing (FR-184) | 96 hrs | Kafka Streams, batch/real-time |
| Device Management Portal (FR-183) | 80 hrs | Real-time monitoring, device health |
| Real-time Alerts (FR-181) | 64 hrs | Rule engine, threshold monitoring |
| Usage-Based Insurance Pricing (FR-182) | 120 hrs | Dynamic premium algorithm, scoring |
| Device Inventory (FR-185) | 48 hrs | Registry, heartbeat monitoring |
| IoT Data Dashboard (FR-127) | 80 hrs | Customer-facing insights |
| **SUBTOTAL** | **688 hrs** | |
| **Buffer (10%)** | **69 hrs** | |
| **TOTAL IoT** | **757 hrs** | |

---

#### AI Engine & Multi-Agent Network (FG-014) - **Python Devs (50% time each)**

**AI Services Architecture:**
- Agent 1: Document Processing (OCR, NID validation)
- Agent 2: Customer Service (Bengali NLP, chatbot)
- Agent 3: Risk Assessment (Fraud detection, claim scoring)
- Agent 4: Business Intelligence (Predictive analytics)

| Feature | Hours | Notes |
|---------|-------|-------|
| AI Infrastructure Setup | 80 hrs | Python FastAPI, gRPC |
| Document Processing Agent (FR-165, FR-169) | 160 hrs | OCR, face matching, NID validation |
| Customer Service Chatbot (FR-164) | 120 hrs | Bengali NLP, 80% resolution rate |
| Fraud Detection ML Model (FR-166) | 140 hrs | Pattern recognition, anomaly detection |
| Predictive Analytics (FR-167) | 100 hrs | Risk assessment, premium optimization |
| Voice Assistant (FR-168, FG-015) | 180 hrs | Bengali STT/TTS, voice workflow |
| AI Training & Model Deployment | 80 hrs | ML pipeline, continuous learning |
| **SUBTOTAL** | **860 hrs** | |
| **Buffer (10%)** | **86 hrs** | |
| **TOTAL AI** | **946 hrs** | |

**Assignment:** Python Dev 1 (50%) + Python Dev 2 (50%) = ~1 FTE for M3

---

#### Remaining Portals (M3)

**Portal 5: Customer Support Portal** - 200 hrs
- Support ticket queue
- Call center integration
- Customer history view
- Resolution templates

**Portal 6: DevOps Portal** - 180 hrs
- Grafana/Prometheus integration
- System health monitoring
- Alert management
- Log aggregation

**Portal 7-12: Additional Portals** - 400 hrs
- Database Manager Portal
- Focal Person Portal
- Partner Admin Portal
- General Staff Portal
- Vendor Portal
- Marketing Page Admin

**Total M3 Portals:** 780 hrs

---

#### Advanced Features (M3)

**Performance Optimizations** - 320 hrs
- Database query optimization
- Read replicas setup
- Caching strategies (Redis)
- CDN configuration

**Advanced Security** - 280 hrs
- 2FA enforcement (FR-017)
- Enhanced fraud detection rules
- Penetration testing
- Security audit

**Advanced Analytics** - 240 hrs
- Churn prediction (FR-150)
- Customer lifetime value
- Geographic risk heatmaps (FR-202)
- Pre-built dashboards (FR-203)

---

### M3 Effort Summary

| Component | Estimated Hours |
|-----------|----------------|
| IoT Integration & UBI | 757 hrs |
| AI Engine & Multi-Agent Network | 946 hrs |
| Remaining Portals (6) | 780 hrs |
| Performance Optimizations | 320 hrs |
| Advanced Security | 280 hrs |
| Advanced Analytics | 240 hrs |
| Voice Assistant Features | 180 hrs (in AI) |
| Additional Features | 400 hrs |
| QA & Testing | 640 hrs |
| Documentation & Training | 320 hrs |
| **TOTAL M3** | **4,863 hrs** |

**M3 Available Capacity:** 8,535 hours
**M3 Status:** ✅ Comfortable (3,672 hrs buffer for polish & optimization)

---

## Summary: Final Project Capacity vs Effort (ALIGNED)

| Milestone | Available Hours | Required Hours | Utilization | Status |
|-----------|-----------------|----------------|-------------|--------|
| **M1 (Beta - Mar 1)** | 4,040 hrs | 2,740 hrs | 68% | ✅ **COMFORTABLE** |
| **M2 (Launch - Apr 14)** | 2,256 hrs | 1,641 hrs | 73% | ✅ **COMFORTABLE** |
| **M3 (Complete - Aug 1)** | 6,768 hrs | 4,863 hrs | 72% | ✅ **COMFORTABLE** |
| **TOTAL** | **13,064 hrs** | **9,244 hrs** | **71%** | ✅ **ACHIEVABLE** |

**Key Changes:**
- ✅ M1: Customer Mobile App (Android + iOS) + 2 Portals with mocks (69% utilization)
- ✅ M2: Real API integration + Claims backend (73% utilization)
- ✅ M3: IoT + AI + 10 Portals + Advanced features (72% utilization)
- ✅ NO Agent Mobile App - Per SRS, agents use web portal on tablets

**Team Allocation:**
- CTO 50%, Mamoon 40%, Sagor 50%, Python 50%, Sujon 60%, Rumon 60%, React 70%, QA 60%, Mobile/Backend 90%

---

## Key Decisions & Trade-offs

### M1 Reductions (to meet capacity):
1. ✅ **Only 2 portals** (Partner + System Admin) instead of 6 - Saves 378 hrs
2. ✅ **No Agent App** in M1 - Saves 678 hrs  
3. ✅ **Simplified Claims** (manual approval only) - Saves 184 hrs
4. ✅ **Basic Mobile App** (customer only, core features) - Saves 678 hrs
5. ✅ **Reduced DevOps** (essential infra only) - Saves 145 hrs

### M2 Additions (Grand Launch):
1. ✨ Business Admin Portal + Agent Portal
2. ✨ Agent Mobile Apps (Android + iOS)
3. ✨ Partner onboarding automation
4. ✨ Commission engine
5. ✨ Analytics & BI dashboards

### M3 Focus (as per SRS shift):
1. 🚀 **All IoT features** (14 requirements) - 757 hrs
2. 🤖 **AI Engine & Multi-Agent Network** - 946 hrs
3. 🎯 **Remaining 6 portals** - 780 hrs
4. ⚡ **Performance & Security** - 600 hrs
5. 📊 **Advanced Analytics** - 240 hrs

---

## Risk Mitigation for M1

**106-hour deficit solutions:**
1. ✅ **CTO 50% time** = 240 hrs (sufficient to cover deficit)
2. ✅ **Mamoon 40% time** + **Sujon 100%** can handle payment + ticketing
3. ✅ **Mobile devs** can help with portal components (weeks 7-10)
4. ✅ **Simplified flows** reduce complexity
5. ✅ **Existing code** (1,420 hrs worth) reduces unknowns

**Confidence Level:** HIGH (with existing code reuse and team flexibility)

---
