# 14. Appendices

**Note:** Proto schema definitions and workflow examples are in Appendix A and B (auto-generated from actual files).

## 14.1 Glossary of Terms

### Business Terms
| Term | Definition |
|------|------------|
| **Policyholder** | Individual or entity that purchases and owns an insurance policy |
| **Insured** | Person or entity covered by the insurance policy (may differ from policyholder) |
| **Nominee** | Designated beneficiary who receives claim proceeds |
| **Premium** | Payment made for insurance coverage |
| **Sum Insured** | Maximum amount payable under the policy |
| **Deductible** | Amount policyholder must pay before insurance coverage begins |
| **Copay** | Fixed amount paid by insured for specific covered services |
| **Underwriting** | Process of evaluating risk and determining premium |
| **Claim** | Request for payment under an insurance policy |
| **Settlement** | Payment of approved claim amount |
| **Lapse** | Policy termination due to non-payment of premium |
| **Grace Period** | Extended time after premium due date before policy lapses |
| **Renewal** | Extension of policy coverage for additional term |
| **Rider** | Additional coverage option added to base policy |
| **Exclusion** | Circumstances or conditions not covered by policy |

### Technical Terms
| Term | Definition |
|------|------------|
| **Proto3** | Protocol Buffers version 3 - Google's data serialization format |
| **gRPC** | High-performance RPC framework using HTTP/2 |
| **CQRS** | Command Query Responsibility Segregation pattern |
| **Event Sourcing** | Storing state changes as sequence of events |
| **VSA** | Vertical Slice Architecture - organizing code by feature |
| **Hypertable** | TimescaleDB table optimized for time-series data |
| **Materialized View** | Precomputed query results stored as table |
| **Circuit Breaker** | Fault tolerance pattern preventing cascading failures |
| **Idempotency** | Operation producing same result when called multiple times |
| **JWT** | JSON Web Token for stateless authentication |

## 14.2 Reference Documents

**External References:**
- IDRA Insurance Act 2010
- Bangladesh Insurance Rules 2011
- Digital Security Act 2018
- PCI-DSS 3.2.1 Standards
- ISO 27001:2013 Information Security

### Appendix A - API Architecture Diagram

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

### Appendix B - Stakeholder Hierarchy

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

### Appendix C - Claims Approval Workflow

**Workflow Diagram:**

```
┌──────────────────┐
│ Customer submits │
│ claim via app    │
└────────┬─────────┘
         │
         ▼
┌──────────────────┐
│ System validates │
│ and screens      │
└────────┬─────────┘
         │
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

### Appendix D - Requirements Traceability Matrix

**Business Objective to Requirements Mapping:**

| Business Objective | Related Functional Requirements | Success Metrics |
|--------------------|--------------------------------|-----------------|
| **Digital Onboarding: 40,000 policies by 2026** | FR-001 to FR-016, FR-033 to FR-040 | Monthly policy acquisition rate |
| **API Performance Optimization** | FR-107 to FR-118, NFR-008 to NFR-011 | API response time metrics |
| **Financial Transaction Integrity** | FR-121 (TigerBeetle), SEC-003 (PCI-DSS) | Transaction accuracy and speed |
| **Regulatory Compliance** | SEC-011 to SEC-020 (IDRA/AML/CFT) | Audit compliance score |
| **Partner Management Excellence** | FR-011 (Focal Person), FR-086 to FR-092 | Number of active partners |
| **Claims Efficiency** | FR-041 to FR-058, FR-133 to FR-137 | Average claim TAT |

**Change Control Process:**

**Approval Hierarchy:**
1. Dev submits change request
2. Repository Admin reviews code changes
3. Database Admin reviews data model impact
4. System Admin reviews infrastructure impact
5. Business Admin approves business impact
6. Focal Person approves partner-related changes

**Summary Statistics:**
- **Total Functional Requirements:** 150 (FR-001 to FR-150)
- **Total Non-Functional Requirements:** 45 (NFR-001 to NFR-045)
- **Total IDRA Compliance Requirements:** 10 (IDRA-001 to IDRA-010)
- **Total BFIU Compliance Requirements:** 10 (BFIU-001 to BFIU-010)
- **Total AML Monitoring Rules:** 10 (TM-001 to TM-010)

**By Priority:**
- **M1 (Must Have - Phase 1):** 58 requirements
- **M2 (Must Have - Phase 2):** 45 requirements
- **D (Desirable):** 32 requirements
- **F (Future):** 15 requirements

**By Function Group:**
- **User Management (FG-001):** 8 requirements
- **Product Management (FG-002):** 6 requirements
- **Policy Lifecycle (FG-003):** 8 requirements
- **Claims Management (FG-004):** 8 requirements
- **Partner & Agent Management (FG-005):** 8 requirements
- **Payment Processing (FG-006):** 8 requirements
- **Customer Support (FG-007):** 7 requirements
- **Notifications (FG-008):** 6 requirements
- **IoT Integration (FG-009):** 4 requirements
- **Analytics & Reporting (FG-010):** 5 requirements
- **Audit & Logging (FG-011):** 4 requirements
- **Integration & APIs (FG-012):** 5 requirements

### Appendix E - Technology Stack Summary

**Programming Languages:**
- **Go:** Microservices backbone (Gateway, Auth, Storage, IoT)
- **C# .NET 8:** Business logic (Insurance Engine, Partner Management)
- **Node.js:** Payment processing and customer support
- **Python:** AI/ML and document processing
- **TypeScript/React:** Web interfaces and admin portals
- **React Native:** Mobile applications

**Data & Persistence:**
- **Protocol Buffers:** All service contracts and data models
- **PostgreSQL:** Transactional data and business logic
- **MongoDB:** Product catalog and document metadata
- **Redis:** Caching and session management
- **Apache Kafka:** Event streaming and audit trails

**Infrastructure:**
- **Docker & Kubernetes:** Containerization and orchestration
- **AWS/Azure:** Cloud platform and managed services
- **Prometheus/Grafana:** Monitoring and alerting
- **Jaeger:** Distributed tracing and observability

### Appendix F - Compliance Checklist

**IDRA Requirements Checklist:**
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

### Appendix G - Integration Endpoints

**External API Integrations:**
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

---

**Document Control:**
- **Version:** 3.1 Final Draft
- **Last Updated:** January 2025
- **Next Review:** March 2025
- **Approved By:** Director,CEO,CTO, Project Manager,Senior Devs,Teamleads
- **Distribution:** Development Team, Business Stakeholders, Compliance Team

---

---

## Document Approval & Sign-off

**Technical Architecture Confirmation:**

By signing below, stakeholders confirm their acceptance of this System Requirements Specification V3.0 including:

- **Function Groups (FG-001 to FG-017)** with incremental Functional Requirements (FR-001 to FR-150)
- **SEC numbering (SEC-001 to SEC-020)** for security requirements
- **API Architecture:** Category 1 (gRPC), Category 2 (GraphQL), Category 3 (REST), Public (REST)
- **Database Strategy:** PostgreSQL V17 + TigerBeetle + Redis + DynamoDB/MongoDB + S3 + SQLite + Pgvector
- **Novel Features:** Focal Person role, Joint Approval, WebRTC video calls, Kafka orchestration, Claims Approval Matrix
- **User Categorization:** Type 1 (Urban), Type 2 (Semi-Urban), Type 3 (Rural/Voice-assisted)
- **Mobile Constraints:** Offline mode, Low bandwidth mode, Gradual download (10MB → 100MB)
- **Compliance:** IDRA reports, AML/CFT monitoring, BFIU compliance
- **VSA Architecture:** Vertical Slice Architecture with CQRS and MediatR pattern
- **Proto-First Approach:** Protocol Buffers for all data models
