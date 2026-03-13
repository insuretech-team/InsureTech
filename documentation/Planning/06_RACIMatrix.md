## 5. Responsibility Distribution - RACI Matrix

### 5.1 RACI Legend
- **R = Responsible:** Person who does the work
- **A = Accountable:** Person who is ultimately answerable (only one per task)
- **C = Consulted:** Person whose input is sought
- **I = Informed:** Person who is kept up-to-date

---

### 5.2 M1 - Core Services RACI Matrix

#### User Service

| Task/Activity | Backend Lead | Backend Dev 1 | Backend Dev 2 | Frontend Dev | Mobile Dev 1 | Mobile Dev 2 | DevOps | QA Lead | QA Tester | UI/UX | Product Owner |
|---------------|-------------|---------------|---------------|--------------|-------------|-------------|--------|---------|-----------|-------|---------------|
| Authentication Design | A/R | C | C | I | I | I | C | I | I | I | C |
| JWT Implementation | A | R | R | I | I | I | C | I | I | - | I |
| RBAC Implementation | A | R | R | I | I | I | I | I | I | - | I |
| User Registration | C | A/R | R | I | I | I | I | I | I | - | I |
| Profile Management | C | A/R | R | I | I | I | - | I | I | - | I |
| Password Reset | C | R | A/R | I | I | I | C | I | I | - | I |
| MFA Implementation | A | R | R | I | I | I | C | I | I | - | C |
| API Documentation | A | R | R | C | C | C | I | I | I | - | I |
| Unit Testing | C | A/R | R | - | - | - | - | C | C | - | I |
| Integration Testing | C | R | R | - | - | - | C | A/R | R | - | I |

---

#### Policy Service

| Task/Activity | Backend Lead | Backend Dev 1 | Backend Dev 2 | Backend Dev 3 | Frontend Dev | DevOps | QA Lead | QA Tester | UI/UX | Product Owner |
|---------------|-------------|---------------|---------------|---------------|--------------|--------|---------|-----------|-------|---------------|
| Policy Schema Design | A/R | R | R | R | I | I | I | I | I | C |
| Policy CRUD Operations | A | R | R | R | I | I | I | I | - | I |
| Policy Search | C | A/R | R | C | C | I | I | I | - | C |
| Policy Renewal Logic | A | R | R | R | I | C | I | I | - | C |
| Premium Calculation | A/R | R | R | R | I | I | I | I | - | C |
| Document Generation | C | A/R | R | R | I | C | I | I | - | I |
| API Documentation | A | R | R | R | C | I | I | I | - | I |
| Integration Testing | C | R | R | R | - | C | A/R | R | - | I |

---

#### Claim Service

| Task/Activity | Backend Lead | Backend Dev 1 | Backend Dev 2 | Backend Dev 3 | Frontend Dev | DevOps | QA Lead | QA Tester | Product Owner |
|---------------|-------------|---------------|---------------|---------------|--------------|--------|---------|-----------|---------------|
| Claim Models Design | A/R | R | R | R | I | I | I | I | C |
| Claim Filing | C | A/R | R | R | C | I | I | I | C |
| Approval Workflow | A | R | R | R | I | C | I | I | C |
| Claim Assessment | A/R | R | R | R | I | I | I | I | C |
| Settlement Logic | A | R | R | R | C | C | I | I | C |
| Fraud Detection | A/R | R | R | R | I | I | I | I | C |
| Document Integration | C | R | A/R | R | I | C | I | I | I |
| Integration Testing | C | R | R | R | - | C | A/R | R | I |

---

#### Payment Service

| Task/Activity | Backend Lead | Backend Dev 1 | Backend Dev 2 | DevOps | QA Lead | QA Tester | Product Owner |
|---------------|-------------|---------------|---------------|--------|---------|-----------|---------------|
| Payment Gateway Setup | A | R | R | R | I | I | C |
| Payment Processing | A/R | R | R | C | I | I | C |
| Invoice Generation | C | A/R | R | I | I | I | I |
| Refund Processing | A | R | R | C | I | I | C |
| Payment Reconciliation | A/R | R | R | I | I | I | C |
| Failed Payment Handling | C | A/R | R | C | I | I | I |
| Security Implementation | A | R | R | R | C | C | C |
| Integration Testing | C | R | R | C | A/R | R | I |

---

#### Document Service

| Task/Activity | Backend Lead | Backend Dev 1 | Backend Dev 2 | DevOps | QA Lead | QA Tester | Product Owner |
|---------------|-------------|---------------|---------------|--------|---------|-----------|---------------|
| Storage Integration (S3) | A | R | R | R | I | I | I |
| File Upload Logic | C | A/R | R | C | I | I | I |
| Document Metadata | C | R | A/R | I | I | I | I |
| Access Control | A | R | R | C | I | I | C |
| OCR Integration | A/R | R | R | I | I | I | C |
| Document Versioning | C | A/R | R | I | I | I | I |
| Integration Testing | C | R | R | C | A/R | R | I |

---

#### Notification Service

| Task/Activity | Backend Lead | Backend Dev 1 | DevOps | QA Lead | QA Tester | Product Owner |
|---------------|-------------|---------------|--------|---------|-----------|---------------|
| Email Integration | A | R | C | I | I | I |
| SMS Integration | A | R | C | I | I | I |
| Push Notification Setup | A | R | R | I | I | I |
| Template Management | C | A/R | I | I | I | C |
| Notification Queue | A/R | R | R | I | I | I |
| Multi-language Support | C | A/R | I | I | I | C |
| Integration Testing | C | R | C | A/R | R | I |

---

#### Customer Service

| Task/Activity | Backend Lead | Backend Dev 1 | Frontend Dev | QA Lead | QA Tester | Product Owner |
|---------------|-------------|---------------|--------------|---------|-----------|---------------|
| Ticket System Design | A/R | R | C | I | I | C |
| Ticket CRUD | C | A/R | C | I | I | I |
| Workflow Engine | A | R | I | I | I | C |
| Communication Module | C | A/R | C | I | I | I |
| FAQ Management | C | A/R | C | I | I | C |
| Knowledge Base | C | A/R | C | I | I | C |
| Integration Testing | C | R | - | A/R | R | I |

---

### 5.3 Frontend Development RACI Matrix

#### Web Admin Panel

| Task/Activity | Frontend Dev | Mobile Dev 1 | Mobile Dev 2 | Backend Lead | UI/UX | QA Lead | QA Tester | Product Owner |
|---------------|-------------|-------------|-------------|--------------|-------|---------|-----------|---------------|
| Project Setup | A/R | I | I | C | I | I | I | I |
| Design System Integration | A/R | C | C | I | R | I | I | I |
| Authentication UI | A/R | C | C | C | C | I | I | I |
| Dashboard Development | A/R | I | I | C | C | I | I | C |
| User Management UI | A/R | C | C | C | C | I | I | I |
| Policy Management UI | A/R | R | R | C | C | I | I | C |
| Claim Management UI | A/R | R | R | C | C | I | I | C |
| Payment Management UI | A/R | R | R | C | C | I | I | C |
| Document Management UI | A/R | C | C | C | C | I | I | I |
| Customer Support UI | A/R | R | C | C | C | I | I | C |
| State Management | A/R | C | C | C | - | I | I | I |
| API Integration | A/R | R | R | C | - | I | I | I |
| Responsive Design | A/R | C | C | I | R | I | I | I |
| E2E Testing | C | C | C | - | - | A/R | R | I |

**Note:** Mobile Developers start supporting frontend from Sprint 4, Week 2 (around Feb 10, 2026)

---

### 5.4 Mobile Development RACI Matrix

#### Mobile Apps (Customer & Agent)

| Task/Activity | Mobile Dev 1 | Mobile Dev 2 | Backend Lead | UI/UX | QA Lead | QA Tester | Product Owner |
|---------------|-------------|-------------|--------------|-------|---------|-----------|---------------|
| Project Setup | A/R | R | C | I | I | I | I |
| Architecture Design | A/R | R | C | I | I | I | I |
| Authentication Flow | A/R | R | C | C | I | I | I |
| Customer App - Dashboard | A/R | R | C | C | I | I | C |
| Customer App - Policy View | A/R | R | C | C | I | I | C |
| Customer App - Claim Filing | A/R | R | C | C | I | I | C |
| Customer App - Payments | R | A/R | C | C | I | I | C |
| Customer App - Documents | R | A/R | C | C | I | I | I |
| Customer App - Notifications | R | A/R | C | C | I | I | I |
| Agent App - Dashboard | A/R | R | C | C | I | I | C |
| Agent App - Policy Mgmt | A/R | R | C | C | I | I | C |
| Agent App - Customer Mgmt | R | A/R | C | C | I | I | C |
| Agent App - Lead Tracking | R | A/R | C | C | I | I | C |
| Push Notifications | A/R | R | C | - | I | I | I |
| Offline Mode | R | A/R | C | - | I | I | I |
| API Integration | A/R | R | C | - | I | I | I |
| App Store Submission | A/R | R | - | - | I | I | C |
| Testing | R | R | - | - | A/R | R | I |

**Reassignment Period:** Sprint 4 Week 2 onwards - Mobile devs support Frontend

---

### 5.5 Infrastructure & DevOps RACI Matrix

| Task/Activity | DevOps | Backend Lead | Frontend Dev | QA Lead | Product Owner |
|---------------|--------|-------------|-------------|---------|---------------|
| Cloud Infrastructure Setup | A/R | C | I | I | I |
| Docker Containerization | A/R | C | C | I | I |
| Kubernetes Configuration | A/R | C | I | I | I |
| CI/CD Pipeline | A/R | C | C | C | I |
| Database Setup | A/R | C | - | I | I |
| Message Broker Setup | A/R | C | - | I | I |
| API Gateway Configuration | A/R | C | C | I | I |
| Monitoring & Logging | A/R | C | C | C | I |
| Security Configuration | A/R | C | C | C | C |
| Backup & DR | A/R | C | - | I | C |
| Load Testing | R | C | C | A/R | I |
| Performance Optimization | A/R | R | R | C | I |
| Production Deployment | A/R | C | C | C | C |
| Incident Management | A/R | C | C | C | I |

---

### 5.6 Quality Assurance RACI Matrix

| Task/Activity | QA Lead | QA Tester | Backend Lead | Frontend Dev | Mobile Dev 1 | Mobile Dev 2 | DevOps | Product Owner |
|---------------|---------|-----------|--------------|-------------|-------------|-------------|--------|---------------|
| Test Plan Development | A/R | R | C | C | C | C | I | C |
| Test Case Creation | A/R | R | C | C | C | C | I | C |
| Manual Testing | A | R | I | I | I | I | I | I |
| API Testing | A/R | R | C | - | - | - | C | I |
| UI Testing | A | R | I | C | C | C | I | I |
| Mobile Testing | A | R | I | - | C | C | I | I |
| Integration Testing | A/R | R | C | C | C | C | C | I |
| Performance Testing | A/R | R | C | C | - | - | C | I |
| Security Testing | A/R | R | C | C | - | - | C | C |
| Regression Testing | A | R | I | I | I | I | I | I |
| Test Automation | A/R | R | C | C | C | C | C | I |
| Bug Reporting | A | R | I | I | I | I | I | I |
| UAT Coordination | A/R | R | C | C | C | C | C | C |
| Test Sign-off | A | R | I | I | I | I | I | C |

---

### 5.7 UI/UX Design RACI Matrix

| Task/Activity | UI/UX Designer | Product Owner | Frontend Dev | Mobile Dev 1 | Mobile Dev 2 | Backend Lead | QA Lead |
|---------------|---------------|---------------|-------------|-------------|-------------|--------------|---------|
| User Research | A/R | C | I | I | I | I | I |
| Personas Creation | A/R | C | I | I | I | I | I |
| Information Architecture | A/R | C | C | I | I | C | I |
| Wireframes | A/R | C | C | C | C | I | I |
| High-Fidelity Mockups | A/R | C | C | C | C | I | I |
| Prototyping | A/R | C | C | C | C | I | I |
| Design System | A/R | I | R | R | R | I | I |
| Icon & Asset Creation | A/R | I | C | C | C | I | I |
| Design Handoff | A/R | I | R | R | R | I | I |
| Design QA | A/R | C | C | C | C | I | C |
| Design Documentation | A/R | I | C | C | C | I | I |
| Design Iterations | A/R | C | C | C | C | I | C |

---

### 5.8 Cross-Functional Activities RACI Matrix

| Activity | Product Owner | Backend Lead | Frontend Dev | Mobile Dev 1 | Mobile Dev 2 | DevOps | QA Lead | UI/UX | Stakeholders |
|----------|---------------|-------------|-------------|-------------|-------------|--------|---------|-------|--------------|
| Sprint Planning | A/R | R | R | R | R | R | R | R | I |
| Daily Standups | C | R | R | R | R | R | R | R | - |
| Backlog Refinement | A/R | C | C | C | C | C | C | C | I |
| Sprint Review | A/R | R | R | R | R | R | R | R | C |
| Sprint Retrospective | A | R | R | R | R | R | R | R | - |
| Release Planning | A/R | C | C | C | C | C | C | I | C |
| Architecture Decisions | C | A/R | C | C | C | C | I | I | I |
| Technical Documentation | C | A/R | R | R | R | R | I | I | I |
| Code Reviews | - | A/R | R | R | R | I | I | - | - |
| Deployment | C | R | I | I | I | A/R | R | - | I |
| Post-Mortem Analysis | A | R | R | R | R | R | R | I | C |

---

### 5.9 Decision-Making Authority Matrix

| Decision Type | Primary Authority | Must Consult | Must Inform |
|--------------|-------------------|--------------|-------------|
| **Technical Architecture** | Backend Lead | DevOps, Frontend Dev | Product Owner, Team |
| **UI/UX Design** | UI/UX Designer | Product Owner, Frontend Dev | Team |
| **Feature Priority** | Product Owner | Backend Lead, QA Lead | Team |
| **Sprint Scope** | Product Owner + Backend Lead | Team Leads | Team, Stakeholders |
| **Technology Stack** | Backend Lead + DevOps | Team Leads | Product Owner |
| **Security Decisions** | Backend Lead + DevOps | QA Lead | Product Owner, Stakeholders |
| **Performance Standards** | Backend Lead + DevOps | QA Lead | Product Owner |
| **Release Decisions** | Product Owner | Backend Lead, QA Lead, DevOps | Team, Stakeholders |
| **Resource Allocation** | Product Owner + Backend Lead | Team Leads | Team |
| **Budget Decisions** | Product Owner | Backend Lead, DevOps | Stakeholders |

---

### 5.10 Escalation Matrix

| Issue Level | Primary Contact | Escalation Path | Response Time |
|------------|-----------------|-----------------|---------------|
| **Level 1: Minor** | Individual Developer | Team Lead | 4 hours |
| **Level 2: Moderate** | Team Lead | Backend Lead / Product Owner | 2 hours |
| **Level 3: Major** | Backend Lead / Product Owner | Steering Committee | 1 hour |
| **Level 4: Critical** | Product Owner | Executive Sponsor | Immediate |

**Issue Examples:**
- **Level 1:** Bug fixes, minor clarifications, local environment issues
- **Level 2:** Integration issues, scope clarifications, resource conflicts
- **Level 3:** Major technical blockers, timeline risks, budget overruns
- **Level 4:** Security breaches, complete system failure, legal issues

---

### 5.11 Communication Responsibility Matrix

| Communication Type | Owner | Participants | Frequency |
|-------------------|-------|--------------|-----------|
| **Daily Standups** | Backend Lead | All Team | Daily (15 min) |
| **Sprint Planning** | Product Owner | All Team | Every 2 weeks (4 hrs) |
| **Sprint Review** | Product Owner | Team + Stakeholders | Every 2 weeks (2 hrs) |
| **Sprint Retrospective** | Backend Lead | Team only | Every 2 weeks (1 hr) |
| **Backlog Refinement** | Product Owner | Team Leads | Weekly (2 hrs) |
| **Technical Sync** | Backend Lead | Developers, DevOps | Twice weekly (1 hr) |
| **Stakeholder Updates** | Product Owner | Stakeholders | Weekly (30 min) |
| **Architecture Review** | Backend Lead | Tech Team | As needed |
| **Design Review** | UI/UX Designer | Frontend, Mobile, PO | As needed |
| **Release Planning** | Product Owner | All Team + Stakeholders | Monthly (2 hrs) |

---

### 5.12 Team Member Roles & Primary Responsibilities

#### Backend Lead
- **Primary Responsibilities:**
  - Overall backend architecture
  - Technical decisions
  - Code review oversight
  - Team coordination
  - Risk management
- **Accountable For:** Backend service quality, API design, integration success

#### Backend Developers (3)
- **Primary Responsibilities:**
  - Service implementation
  - API development
  - Database design
  - Unit testing
  - Code reviews
- **Accountable For:** Assigned services completion, code quality

#### Frontend Developer
- **Primary Responsibilities:**
  - Web admin panel development
  - UI component development
  - State management
  - Frontend integration
  - Responsive design
- **Accountable For:** Web application quality, user experience

#### Mobile Developers (2)
- **Primary Responsibilities:**
  - Mobile app development
  - Cross-platform implementation
  - Push notifications
  - Offline mode
  - **Secondary (from Sprint 4):** Frontend support
- **Accountable For:** Mobile app quality, app store deployment

#### DevOps Engineer
- **Primary Responsibilities:**
  - Infrastructure management
  - CI/CD pipeline
  - Monitoring & logging
  - Security configuration
  - Production support
- **Accountable For:** System availability, deployment success, infrastructure security

#### QA Lead
- **Primary Responsibilities:**
  - Test strategy
  - Quality assurance
  - Test team coordination
  - UAT management
  - Sign-off decisions
- **Accountable For:** Overall product quality, test coverage

#### QA Tester
- **Primary Responsibilities:**
  - Test execution
  - Bug reporting
  - Test automation
  - Regression testing
  - Documentation
- **Accountable For:** Test case execution, bug tracking

#### UI/UX Designer
- **Primary Responsibilities:**
  - Design system
  - User research
  - Wireframes & mockups
  - Prototyping
  - Design QA
- **Accountable For:** Design quality, user experience consistency

#### Product Owner
- **Primary Responsibilities:**
  - Product vision
  - Backlog management
  - Stakeholder communication
  - Feature prioritization
  - Acceptance decisions
- **Accountable For:** Product success, ROI, stakeholder satisfaction

---
