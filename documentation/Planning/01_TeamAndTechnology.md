## Team Members, Technology Stack & Existing Assets

### Team Member Details

#### Current Team (Available December 2025)

| # | Name | Role | Primary Tech Stack | Availability | Hours/Week | Responsibilities |
|---|------|------|-------------------|--------------|------------|------------------|
| 1 | **CTO** | Technical Lead & Gateway Dev | Go, Architecture | 50% coding, 50% management | 24 hrs | API Gateway, Kafka orchestration, IoT Broker, Storage Manager, Technical oversight |
| 2 | **Mamoon** | Senior Full Stack Developer | Node.js, React, MongoDB | 40% (other commitments) | 19 hrs | Payment Service completion, Ticketing Service, Full-stack features |
| 3 | **Sujon Ahmed** | Mid-level Full Stack Developer | Node.js, React, PostgreSQL | 100% (exhaust) | 48 hrs | Backend services, API development, Frontend support |
| 4 | **Rumon** | UI/UX Designer | Figma, Adobe XD, Design Systems | 100% (exhaust) | 48 hrs | Design system, Wireframes, Mockups, User experience |
| 5 | **Nur Hossain** | Android Developer | Kotlin, Java, Android SDK | 100% (exhaust) | 48 hrs | Android app development, Mobile API integration |
| 6 | **Sojol Ahmed** | iOS Developer | Swift, iOS SDK, Xcode | 100% (exhaust) | 48 hrs | iOS app development, Mobile API integration |
| 7 | **QA** | QA/Test Engineer | Selenium, Postman, Jest | 100% (exhaust) | 48 hrs | Manual & automated testing, Quality assurance |
| 8 | **Sagor** | DevOps Engineer | Docker, K8s, Go, Terraform | 50% (other projects) | 24 hrs | Infrastructure, CI/CD, Monitoring, Deployment |
| 9 | **React Dev** | Frontend Developer | React, Next.js, TypeScript | 100% (exhaust, joins Dec 18) | 48 hrs | Web portals, Admin panels, Frontend architecture |

**December Team Capacity:** ~7.4 FTE (355 hrs/week - CTO 50%, Mamoon 40%, Sagor 50%, others 100%)

---

#### New Members Joining January 2026

| # | Name | Role | Primary Tech Stack | Join Date | Hours/Week | Responsibilities |
|---|------|------|-------------------|-----------|------------|------------------|
| 10 | **Project Manager** | Project Manager | Agile, Scrum, JIRA | Jan 1, 2026 | 48 hrs (exhaust) | Sprint planning, Team coordination, Stakeholder management |
| 11 | **Mr. Delowar** | **Senior C# Developer (Lead)** | C# .NET, gRPC, Microservices | Jan 15, 2026 | 48 hrs (exhaust) | Insurance Engine lead, Policy Service, Risk Management, Team mentoring |
| 12 | **C# Developer** | Mid-level C# Developer | C# .NET, gRPC, SQL Server | Jan 1, 2026 | 48 hrs (exhaust) | Partner/Agent Management, Analytics & Reporting services |
| 13 | **Python Dev 1** | Senior Python Developer | Python, FastAPI, gRPC | Jan 1, 2026 | 24 hrs (50% - M2/M3 focus) | AI Engine, LLM multi-agent network, AI assistant service |
| 14 | **Python Dev 2** | Python Developer | Python, FastAPI, gRPC | Jan 1, 2026 | 24 hrs (50% - M2/M3 focus) | MCP servers, AI integrations, Data processing |

**Full Team Capacity (from Jan 15):** ~11.4 FTE (547 hrs/week - Python devs 50% each, CTO 50%, Mamoon 40%, Sagor 50%, others 100%)

---

### Technology Stack by Service

#### Microservices Architecture

| Service | Language/Framework | Status | Developer(s) | Notes |
|---------|-------------------|--------|-------------|-------|
| **API Gateway** | Go | 50% Ready (Existing) | CTO | Reuse + 50% new work |
| **Authentication** | Go | ✅ 100% Ready | CTO | Proven, tested code |
| **Authorization** | Go | ✅ 100% Ready | CTO | Proven, tested code |
| **DBManager** | Go | ✅ 100% Ready | CTO | Proven, tested code |
| **Storage Manager** | Go | ✅ 100% Ready | CTO | Proven, tested code |
| **IoT Broker** | Go | 80% Ready (Existing) | CTO | +20% for new requirements |
| **Payment Service** | Node.js | 70% Ready (Existing) | Mamoon | Bkash integration done |
| **Insurance Engine** | C# .NET + gRPC | 🆕 New | Mr. Delowar + C# Dev | Policy, Contract, Risk, Fraud |
| **Partner/Agent Mgmt** | C# .NET + gRPC | 🆕 New | C# Developer | Partner onboarding, verification |
| **AI Engine** | Python + gRPC | 🆕 New | Python Dev 1 & 2 | LLM, AI assistant, MCP servers |
| **Kafka Orchestration** | Go + Kafka | 🆕 New | CTO | Event streaming, Notification |
| **Ticketing Service** | Node.js | 🆕 New | Mamoon + Sujon | Customer support system |
| **Analytics & Reporting** | C# .NET + gRPC | 🆕 New | C# Developer | Business intelligence, Reports |

---

### Existing Assets & Reusable Code

#### CTO's Proven Go Services (Production-Ready)

| Service | Status | Lines of Code | Test Coverage | Description |
|---------|--------|---------------|---------------|-------------|
| **Authentication Service** | ✅ 100% Ready | ~2,500 LOC | 90%+ | JWT-based auth, OAuth2, Session management |
| **Authorization Service** | ✅ 100% Ready | ~2,000 LOC | 85%+ | RBAC, Permission management, Policy engine |
| **DBManager Service** | ✅ 100% Ready | ~3,000 LOC | 85%+ | Database abstraction, Connection pooling, Query optimization |
| **Storage Manager** | ✅ 100% Ready | ~1,800 LOC | 80%+ | S3/Azure Blob integration, File management |
| **IoT Broker** | 80% Ready | ~4,000 LOC | 75%+ | MQTT broker, Device management, Data ingestion |
| **API Gateway** | 50% Ready | ~3,500 LOC | 70%+ | Routing, Rate limiting, API versioning |

**Total Existing Code:** ~16,800 LOC with high test coverage
**Estimated Savings:** ~1,200 development hours (equivalent to 6-8 weeks of work)

#### Mamoon's Payment Integration (70% Ready)

| Component | Status | Description |
|-----------|--------|-------------|
| **Bkash Integration** | ✅ Complete | Merchant account + Sandbox tested |
| **Payment Processing** | 70% Ready | Payment flow, Invoice generation |
| **Refund Logic** | 50% Ready | Partial refund handling |
| **Payment Gateway** | 70% Ready | Gateway abstraction, Multiple providers |

**Estimated Savings:** ~200 development hours

---

### Development Tools & Infrastructure

#### Already Configured

| Category | Tool/Service | Status |
|----------|-------------|--------|
| **Domain** | trendyco.insurance | ✅ Active |
| **Email** | Google Workspace | ✅ Configured |
| **Server** | Cloud Infrastructure | ✅ Setup |
| **Payment** | Bkash Merchant + Sandbox | ✅ Active |
| **Legal** | Trade License, BIN, TIN | ✅ Done |
| **Marketing** | Primary Website | ✅ Live |

#### Development Stack

| Category | Tools |
|----------|-------|
| **Backend Languages** | Go, C# .NET, Node.js, Python |
| **Frontend** | React, Next.js, TypeScript |
| **Mobile** | Kotlin (Android), Swift (iOS) |
| **Communication** | gRPC, REST API, WebSockets |
| **Message Queue** | Kafka, RabbitMQ |
| **Databases** | PostgreSQL, MongoDB, Redis |
| **Cloud** | AWS / Azure |
| **Containers** | Docker, Kubernetes |
| **CI/CD** | GitHub Actions, Jenkins |
| **Monitoring** | Prometheus, Grafana, Jaeger |
| **Testing** | Jest, Pytest, xUnit, Postman |

---

### Web Portals to Develop (React/Next.js)

**ALIGNED M1 SCOPE:** Only 2 portals for M1 to match capacity

| # | Portal Name | Primary Users | Priority | Notes |
|---|-------------|---------------|----------|-------|
| 1 | **Business Admin Portal** | Business Managers | M1 | Product/policy management, core admin functions |
| 2 | **Partner Portal** | Insurance Partners | M1 | Partner dashboard, commission tracking |
| 3 | **System Admin Portal** | System Administrators | M2 | User management, system config |
| 4 | **Agent Portal** | Insurance Agents | M2 | Field agent operations |
| 5 | **Customer Support Portal** | Support Staff | M2 | Ticketing, customer queries |
| 6 | **Focal Person Portal** | Internal Coordinators | M3 | Partner management, approvals (per SRS M1, but deferred due to capacity) |
| 7 | **DevOps Portal** | DevOps Team (Prometheus, Grafana) | M3 | Infrastructure monitoring |
| 8 | **Database Manager Portal** | DBAs | M3 | Database administration |
| 9 | **Partner Admin Portal** | Partner Administrators | M3 | Advanced partner features |
| 10 | **General Staff Portal** | General Employees | M3 | Internal operations |
| 11 | **Vendor Portal** | Third-party Vendors | M3 | Vendor integrations |
| 12 | **Marketing Page Admin** | Marketing Team | M3 | Content management |

**M1 Portals (with mock APIs):** Business Admin + Partner ONLY (343 + 132 = 475 hrs)
**M2 Portals:** System Admin + Agent + Customer Support (3 portals)
**M3 Portals:** Remaining 7 portals

**Note on Focal Person Portal:** SRS marks FR-066 as M1 priority, but due to capacity constraints (M1 already at 69% utilization), deferred to M3. Manual focal person processes sufficient for M1 beta.

---

### Mobile Applications

| App | Platform | Developer | Priority | Status |
|-----|----------|-----------|----------|--------|
| **Customer App** | Android | Nur Hossain | M1 | With mock APIs, real integration M2 |
| **Customer App** | iOS | Sojol Ahmed | M1 | With mock APIs, real integration M2 |
| **Agent App** | Android | - | **NOT IN SCOPE** | Agents use web portal on tablets |
| **Agent App** | iOS | - | **NOT IN SCOPE** | Agents use web portal on tablets |

**M1 Strategy:** Customer App built in M1 with mock APIs (634 hrs). Real API integration in M2 (458 hrs).
**Agent App Decision:** NO agent mobile app per SRS requirements. Agents use Agent Portal (web) on tablets/desktops. Saves 580 hrs.
**Reasoning:** SRS does not specify agent mobile app as requirement. Web portal on tablet sufficient for field agents.

---

### Capacity Calculation with Reduced Availability

#### December 2025 (Pre-January Hires)

| Team Member | Availability | Hours/Week | Effective Hours |
|-------------|--------------|------------|-----------------|
| CTO | 50% | 24 hrs | 24 hrs |
| Mamoon | 40% | 19 hrs | 19 hrs |
| Sujon Ahmed | 100% (exhaust) | 48 hrs | 48 hrs |
| Rumon | 100% (exhaust) | 48 hrs | 48 hrs |
| Nur Hossain | 100% (exhaust) | 48 hrs | 48 hrs |
| Sojol Ahmed | 100% (exhaust) | 48 hrs | 48 hrs |
| QA | 100% (exhaust) | 48 hrs | 48 hrs |
| Sagor | 50% | 24 hrs | 24 hrs |
| React Dev (from Dec 18) | 100% (exhaust) | 48 hrs | 48 hrs |
| **Total December** | - | **355 hrs/week** | **355 hrs/week** |

**December Sprint (2 weeks):** 710 hours total capacity

---

#### January 2026+ (Full Team)

| Team Member | Availability | Hours/Week | Effective Hours |
|-------------|--------------|------------|-----------------|
| Existing Team (9) | Mixed | 355 hrs | 355 hrs |
| Project Manager | 100% (exhaust) | 48 hrs | 48 hrs |
| Mr. Delowar (from Jan 15) | 100% (exhaust) | 48 hrs | 48 hrs (prorated) |
| C# Developer | 100% (exhaust) | 48 hrs | 48 hrs |
| Python Dev 1 | 50% | 24 hrs | 24 hrs |
| Python Dev 2 | 50% | 24 hrs | 24 hrs |
| **Total January+** | - | **547 hrs/week** | **547 hrs/week** |

**January Sprint (2 weeks):** 1,094 hours total capacity
**With Mr. Delowar joining mid-month (Jan 15):** His hours prorated for first sprint

---

### Buffer Strategy: 10% (Revised from 20%)

**Rationale for 10% Buffer:**
- Existing code reuse reduces unknowns
- Proven, tested components (Auth, DBManager, etc.)
- Experienced team members (CTO, Mamoon)
- Clear architecture already established

**Buffer Allocation:**
- Bug fixes & rework: 4%
- Integration issues: 3%
- Meetings & reviews: 2%
- Unexpected blockers: 1%

**Total Buffer:** 10% of estimated effort

---

### Risk Mitigation Through Existing Assets

| Risk | Mitigation via Existing Assets |
|------|-------------------------------|
| Authentication delays | ✅ 100% ready Go service (CTO) |
| Authorization complexity | ✅ 100% ready Go service (CTO) |
| Database performance | ✅ Proven DBManager service |
| Storage issues | ✅ Storage Manager ready |
| Payment integration | ✅ Bkash integration 70% done (Mamoon) |
| IoT requirements | ✅ IoT Broker 80% ready |
| API Gateway setup | ✅ 50% ready, proven architecture |

**Estimated Risk Reduction:** 40-50% due to existing proven components

---
