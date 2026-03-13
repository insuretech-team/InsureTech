# 3. System Architecture
### 3.1 Architectural Overview

The LabAid InsureTech Platform is built on a **cloud-native, microservices architecture** with **Domain-Driven Design (DDD)** principles. The system leverages **Vertical Slice Architecture (VSA)** pattern across all services for maximum cohesion and maintainability.

**Core Architectural Principles:**
1. **Microservices First:** Independent, deployable services with single responsibilities
2. **Event-Driven Architecture:** Asynchronous communication through domain events
3. **API-First Design:** All services expose well-defined APIs (REST/gRPC)
4. **Cloud-Native:** Built for containerization and orchestration
5. **Security by Design:** Zero-trust security model with end-to-end encryption

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

### 3.3 System Architecture - VSA Pattern

![VSA Architecture](VSA.png)

*Figure 1: Vertical Slice Architecture - Language-Agnostic Pattern*

The LabAid InsureTech Platform adopts **Vertical Slice Architecture (VSA)** across ALL microservices, regardless of programming language:

- **Go Services:** Gateway, Auth, DBManager, Storage, IoT Broker, Kafka Orchestration
- **C# .NET Services:** Insurance Engine, Partner Management, Analytics & Reporting  
- **Node.js Services:** Payment Service, Ticketing Service
- **Python Services:** AI Engine, OCR/PDF Service

**Key VSA Principles:**
1. **High Cohesion:** Each slice contains all layers needed for one feature
2. **Low Coupling:** Slices are independent and don't share logic
3. **Feature-Focused:** Organized by business capability, not technical layer
4. **Testability:** Each slice can be tested in isolation

[[[PAGEBREAK]]]

### 3.4 Microservices Architecture

**Service Inventory:**

| Service | Language | Port | Responsibility | Team Owner |
|---------|----------|------|----------------|------------|
| **Gateway** | Go | 8080 | API Gateway, routing, rate limiting | CTO |
| **Auth Service** | Go | 8081 | Authentication, JWT management | CTO |
| **Authorization** | Go | 8082 | RBAC, permissions, access control | CTO |
| **DBManager** | Go | 8083 | Database operations, migrations | CTO |
| **Storage Service** | Go | 8084 | File storage, S3 operations | CTO |
| **IoT Broker** | Go | 8085 | IoT device communication, MQTT | CTO |
| **Kafka Service** | Go | 8086 | Event orchestration, messaging | CTO |
| **Insurance Engine** | C# .NET | 5001 | Policy lifecycle, underwriting | C# Senior |
| **Partner Management** | C# .NET | 5002 | Partner onboarding, agent mgmt | C# Mid |
| **Analytics & Reporting** | C# .NET | 5003 | BI, dashboards, compliance reports | C# Senior |
| **Payment Service** | Node.js | 3001 | Payment processing, settlements | Mamoon |
| **Ticketing Service** | Node.js | 3002 | Customer support, help desk | Node.js Dev |
| **AI Engine** | Python | 4001 | LLM, chatbot, fraud detection | Python Dev 1 |
| **OCR Service** | Python | 4002 | Document processing, KYC | Python Dev 2 |

[[[PAGEBREAK]]]

### 3.5 System Context Diagram

```text
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





**Key Architectural Highlights:**

1. **Multi-Language Microservices:** Go, C#, Node.js, Python - all communicate via gRPC
2. **Event-Driven:** Kafka for asynchronous communication and event sourcing
3. **VSA Pattern:** Each service implements vertical slices internally
4. **CQRS:** Insurance Engine uses Command Query Responsibility Segregation
5. **Reusable Services:** 755 hours of production-tested code (Gateway, Auth, DBManager, Storage, IoT)

[[[PAGEBREAK]]]