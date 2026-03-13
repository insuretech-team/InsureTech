# 1. Introduction
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
- Payment processing (bkash,nagad,rocket ,manual Phase 1, all channel automated +manual Phase 2)
- Claims submission and approval workflows (with basic check with AI Engine Phase 1, all channel automated +manual Phase 2)
- Partner management and tenant isolation
- Notification system (SMS, Email, Push)
- Reporting and analytics
- IDRA and BFIU compliance features
- Voice-assisted workflow (Bengali speech recognition, voice-guided policy purchase, voice claims submission)
- AI-based Claim Management (fraud detection, automated assessment, risk scoring)
- IoT-based Usage-Based Insurance (UBI) for vehicles and health tracking
- IOT based Tracking system


**Out of Scope:**
- Full AI-driven underwriting (Phase 2.5/3)
- Universal IoT/Telematics integration (Phase 2.5/3)
- Cross-border insurance (Future consideration)
- Blockchain-based smart contracts (Future consideration)
[[[PAGEBREAK]]]

### 1.3 Definitions, Acronyms & Abbreviations

| Term | Definition |
|------|------------|
| **IDRA** | Insurance Development & Regulatory Authority of Bangladesh |
| **BFIU** | Bangladesh Financial Intelligence Unit |
| **MFS** | Mobile Financial Services (bKash, Nagad, Rocket) |
| **KYC** | Know Your Customer |
| **AML** | Anti-Money Laundering |
| **CFT** | Combating the Financing of Terrorism |
| **gRPC** | Google Remote Procedure Call |
| **CQRS** | Command Query Responsibility Segregation |
| **VSA** | Vertical Slice Architecture |
| **Proto** | Protocol Buffers |
| **IoT** | Internet of Things |
| **AI** | Artificial Intelligence |
| **ML** | Machine Learning |
| **API** | Application Programming Interface |
| **SLA** | Service Level Agreement |
| **TAT** | Turn Around Time |
| **EHR** | Electronic Health Records |
| **OCR** | Optical Character Recognition |
| **SMS** | Short Message Service |
| **OTP** | One-Time Password |
| **JWT** | JSON Web Token |
| **RBAC** | Role-Based Access Control |

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

[[[PAGEBREAK]]]
