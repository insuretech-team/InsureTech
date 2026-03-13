"""
Comprehensive update script for SRS_V3_FINAL_DRAFT.md
Fixes all issues identified by user
"""

import re

# Read the current file
with open("SRS_V3_FINAL_DRAFT.md", "r", encoding="utf-8") as f:
    content = f.read()

# 1. Fix Approval Signatures - Replace with user's exact format
old_approval = r"""## Approval Signatures

\| Role \| Name \| Signature \| Date \|
\|------|------|-----------|------|
\| \*\*Project Sponsor\*\* \| Director \| _________________ \| ___/___/2025 \|
\| \*\*Chief Executive Officer\*\* \| CEO \| _________________ \| ___/___/2025 \|
\| \*\*Chief Technology Officer\*\* \| CTO \| _________________ \| ___/___/2025 \|
\| \*\*Chief Financial Officer\*\* \| CFO \| _________________ \| ___/___/2025 \|
\| \*\*Project Manager\*\* \| InsureTech \| _________________ \| ___/___/2025 \|
\| \*\*Senior Dev C#\*\* \| LifePlus \| _________________ \| ___/___/2025 \|
\| \*\*Senior Dev TS\*\* \| LifePlus \| _________________ \| ___/___/2025 \|
\| \*\*AI Lead Python\*\* \| LifePlus \| _________________ \| ___/___/2025 \|"""

new_approval = """## Approval Signatures

> **Note:** Please review and sign to approve this System Requirements Specification V3.0 FINAL DRAFT.

| Role | Name | Signature | Date |
|------|------|-----------|------|
| **Project Sponsor** | Director | _________________ | ___/___/2025 |
| **Chief Executive Officer** | CEO | _________________ | ___/___/2025 |
| **Chief Technology Officer** | CTO | _________________ | ___/___/2025 |
| **Chief Financial Officer** | CFO | _________________ | ___/___/2025 |
| **Project Manager** | InsureTech | _________________ | ___/___/2025 |
| **Senior Dev C#** | LifePlus | _________________ | ___/___/2025 |
| **Senior Dev TS** | LifePlus | _________________ | ___/___/2025 |
| **AI Lead Python** | LifePlus | _________________ | ___/___/2025 |"""

content = re.sub(old_approval, new_approval, content)

# 2. Add VSA image reference
vsa_image_section = """
### 4.1.1 VSA Architecture Diagram

![VSA Architecture](VSA.png)

*Figure 1: Vertical Slice Architecture - Language Agnostic Pattern for all microservices (Go, C#, Node.js, Python)*

The VSA pattern shown above is applied consistently across ALL microservices regardless of programming language:

- **Go Services:** Gateway, Auth, DBManager, Storage, IoT Broker, Kafka Orchestration
- **C# .NET Services:** Insurance Engine, Partner Management, Analytics & Reporting
- **Node.js Services:** Payment Service, Ticketing Service
- **Python Services:** AI Engine, OCR/PDF Service

Each service implements vertical slices where a single feature owns its entire stack from API endpoint through business logic to data access.

---
"""

# Find where to insert (after ### 4.1 Architectural Principles)
content = content.replace(
    "### 4.1 Architectural Principles",
    "### 4.1 Architectural Principles" + "\n\n" + vsa_image_section
)

# 3. Fix System Context Diagram (4.2) - recreate properly
fixed_diagram = """### 4.2 System Context Diagram

```
┌─────────────────────────────────────────────────────────────────────┐
│                        External Systems                              │
├─────────────────────────────────────────────────────────────────────┤
│  • bKash/Nagad/Rocket (MFS Payment Gateways)                        │
│  • Hospital EHR Systems (FHIR/HL7)                                   │
│  • NID Verification API (Bangladesh)                                 │
│  • IDRA Portal (Regulatory Reporting)                                │
│  • BFIU Portal (AML/CFT Reporting)                                   │
│  • SMS Gateway & Email Service                                       │
│  • IoT Device APIs (Phase 2 - Livestock, Health Wearables)          │
└─────────────────────────────────────────────────────────────────────┘
                              ▲
                              │ REST/SOAP/WebSocket
                              │
┌─────────────────────────────────────────────────────────────────────┐
│              LabAid InsureTech Platform (Core)                       │
├─────────────────────────────────────────────────────────────────────┤
│                                                                       │
│  ┌───────────────────────────────────────────────────────────────┐  │
│  │              Integration Layer                                 │  │
│  ├───────────────────────────────────────────────────────────────┤  │
│  │  • API Gateway (Go) - Orchestration & Routing                 │  │
│  │  • Kafka Message Bus - Event-Driven Communication             │  │
│  │  • Auth Service (Go) - IAM/RBAC/OAuth2/JWT                    │  │
│  └───────────────────────────────────────────────────────────────┘  │
│                              │                                        │
│                              │ gRPC (Protocol Buffers)               │
│                              ▼                                        │
│  ┌───────────────────────────────────────────────────────────────┐  │
│  │          Backend Microservices (VSA + gRPC)                   │  │
│  ├───────────────────────────────────────────────────────────────┤  │
│  │                                                                │  │
│  │  ┌──────────────────┐  ┌──────────────────┐                  │  │
│  │  │ Insurance Engine │  │ Partner/Agent    │                  │  │
│  │  │ (C# .NET gRPC)   │  │ Management       │                  │  │
│  │  │ - CQRS/MediatR   │  │ (C# .NET gRPC)   │                  │  │
│  │  └──────────────────┘  └──────────────────┘                  │  │
│  │                                                                │  │
│  │  ┌──────────────────┐  ┌──────────────────┐                  │  │
│  │  │ AI Engine        │  │ Kafka            │                  │  │
│  │  │ (Python gRPC)    │  │ Orchestration    │                  │  │
│  │  │ - LLM Agents     │  │ (Go + gRPC)      │                  │  │
│  │  └──────────────────┘  └──────────────────┘                  │  │
│  │                                                                │  │
│  │  ┌──────────────────┐  ┌──────────────────┐                  │  │
│  │  │ Ticketing        │  │ Analytics &      │                  │  │
│  │  │ (Node.js gRPC)   │  │ Reporting        │                  │  │
│  │  │                  │  │ (C# .NET gRPC)   │                  │  │
│  │  └──────────────────┘  └──────────────────┘                  │  │
│  │                                                                │  │
│  │  ┌──────────────────┐  ┌──────────────────┐                  │  │
│  │  │ Payment Service  │  │ OCR/PDF Service  │                  │  │
│  │  │ (Node.js gRPC)   │  │ (Python gRPC)    │                  │  │
│  │  └──────────────────┘  └──────────────────┘                  │  │
│  │                                                                │  │
│  │  ┌──────────────────┐  ┌──────────────────┐                  │  │
│  │  │ DBManager        │  │ Storage Service  │                  │  │
│  │  │ (Go gRPC)        │  │ (Go gRPC)        │                  │  │
│  │  │ 100% ready       │  │ 100% ready       │                  │  │
│  │  └──────────────────┘  └──────────────────┘                  │  │
│  │                                                                │  │
│  │  ┌──────────────────┐                                         │  │
│  │  │ IoT Broker       │                                         │  │
│  │  │ (Go gRPC)        │                                         │  │
│  │  │ 80% ready        │                                         │  │
│  │  └──────────────────┘                                         │  │
│  │                                                                │  │
│  └───────────────────────────────────────────────────────────────┘  │
│                              │                                        │
│                              ▼                                        │
│  ┌───────────────────────────────────────────────────────────────┐  │
│  │                    Data Layer                                  │  │
│  ├───────────────────────────────────────────────────────────────┤  │
│  │  • PostgreSQL 17 (Transactional, ACID)                        │  │
│  │  • MongoDB (Unstructured Data, Application Logs)              │  │
│  │  • Redis Cluster (Cache, Sessions, Rate Limiting)             │  │
│  │  • AWS S3 (Object Storage - Documents, Images)                │  │
│  └───────────────────────────────────────────────────────────────┘  │
│                                                                       │
└─────────────────────────────────────────────────────────────────────┘
                              │
                              │ HTTPS/GraphQL/WebSocket
                              ▼
┌─────────────────────────────────────────────────────────────────────┐
│                    User Interfaces (Client Layer)                    │
├─────────────────────────────────────────────────────────────────────┤
│  • Customer Web Portal (React PWA)                                   │
│  • Customer Mobile App (React Native/Flutter)                        │
│  • Partner Portal (React)                                            │
│  • Admin Portal (React)                                              │
│  • DevOps Portal (Grafana, Prometheus, Jaeger)                      │
└─────────────────────────────────────────────────────────────────────┘
```

**Key Architectural Highlights:**

1. **Multi-Language Microservices:** Go, C#, Node.js, Python - all communicate via gRPC
2. **Event-Driven:** Kafka for asynchronous communication and event sourcing
3. **VSA Pattern:** Each service implements vertical slices internally
4. **CQRS:** Insurance Engine uses Command Query Responsibility Segregation
5. **Reusable Services:** 755 hours of production-tested code (Gateway, Auth, DBManager, Storage, IoT)

"""

# Replace the broken diagram
content = re.sub(
    r"### 4\.2 System Context Diagram.*?```.*?```",
    fixed_diagram.strip(),
    content,
    flags=re.DOTALL
)

print("Step 1: Fixed approval signatures")
print("Step 2: Added VSA.png image reference")  
print("Step 3: Fixed System Context Diagram")
print("\nWriting updated content...")

# Write back
with open("SRS_V3_FINAL_DRAFT.md", "w", encoding="utf-8") as f:
    f.write(content)

print("✅ Phase 1 updates complete!")
print("Next: Update proto structure and add CQRS/MediatR details")
