BUSINESS REQUIREMENTS DOCUMENT (BRD)

Labaid InsureTech Company Limited

> 🟦 **REGIONAL BENCHMARK INSIGHT:** Compared to regional players (PolicyBazaar, Qoala), this plan is overly aggressive on technical complexity. Successful MVPs in this region typically start with simpler architectures and manual back-offices.

Document Purpose

This BRD defines the business needs, functional specifications and technical requirements for developing an advanced InsureTech platform that aligns with:

Labaid Group’s health ecosystem

IDRA regulatory framework

Modern InsureTech/HealthTech best practices

Bangladesh’s digital financial landscape

The 2026–2028 business roadmap

This BRD shall guide technology vendors, internal development teams, UX teams and business stakeholders throughout platform design, development, testing and deployment.

Business Overview

Labaid InsureTech aims to modernize insurance distribution through:

Digital-first experiences

AI-driven underwriting & claims automation 
> ⚠️ **RISK:** AI-driven underwriting in Phase 1 is highly complex.
> 🟦 **BENCHMARK:** **PolicyBazaar** ran on manual rules for years. **Digit** used AI only after significant scale.
> 🟢 **SUGGESTION:** Start with rule-based underwriting first.

Mobile app–based onboarding

Micro-insurance and inclusive products

Integration with LabAid health services 
> 🔴 **ISSUE:** Deep integration with legacy hospital systems is a massive timeline risk. 
> 🟢 **SUGGESTION:** Phase 1 should use loose coupling (e.g., visual verification of digital card) rather than real-time IT integration.

Partner-driven embedded insurance
> 🟦 **BENCHMARK:** **Qoala (Indonesia)** scaled via this model (embedding in Traveloka/Tokopedia). This is a higher-value channel than a standalone B2C app.

The platform will support Life, Health, General, Agriculture and Device Insurance through a unified, seamless mobile experience.

Business Objectives

Short-term (2026)

Launch core InsureTech platform

Enable digital onboarding, policy purchase and claims

Achieve 40,000+ active policies

Integrate with insurers and payment gateways 
> ⚠️ **RISK:** Integrating with multiple "Insurers" assumes they all have ready APIs. This is rarely true in the local market.
> 🟢 **SUGGESTION:** Build a "Mock" adapter layer. If insurers lack APIs, implement a portal for them to manually approve/issue policies triggered by the app.

Mid-term (2027)

Introduce AI underwriting

Automate 80% of claims

Achieve break-even (Q3 2027)

Launch Super-App 2.0

Long-term (2028)

Become Top 3 InsureTech platform in Bangladesh

Expand regionally (Nepal, Bhutan, Maldives)

Implement predictive risk scoring, behavioral pricing, IoT integration

Project Scope

In-Scope

Customer mobile app (iOS/Android) 
> 🟢 **SUGGESTION:** Focus primarily on this for Phase 1.

Web admin portal

Agent/Partner portal 
> 🔴 **ISSUE:** Building an Agent Portal simultaneously with Customer App splits development focus.
> 🟦 **BENCHMARK:** **PolicyStreet (Malaysia)** focused on B2B/Gig workers first. Splitting focus dilutes quality.
> 🟢 **SUGGESTION:** Move Agent Portal to Phase 1.5 or 2.

Policy purchase, renewal, cancellation

Digital KYC & verification

Claims submission & automation

Payment integrations (bKash, Nagad, cards)

Notification system (SMS, push, email)

Integration with insurer APIs

Integration with LabAid hospital systems 
> 🔴 **ISSUE:** As noted above, this is a blocker.
> 🟢 **ALTERNATIVE:** Manual verification for Phase 1.

Out-of-Scope (Phase 1)

Full AI-driven underwriting (Phase 2: 2027)

IoT/Telematics integration (Phase 3: 2028)

Cross-border insurance sales (2028)

User Groups

| User Type | Description |
| --- | --- |
| Customer/Policy Buyer | App users purchasing policies |
| Partner/Agent | Distributors (hospitals, MFS, telcos, e-commerce) |
| Insurer Underwriter | Validates policy and approves underwriting |
| Claims Officer | Processes and approves claims |
| Super Admin | Full platform control |
| Support/Call Center | Customer assistance and issue resolution |

High-Level User Journey Flows

Registration & Login

Signup/Login

OTP Verification

Personal Info

Nominee Info

Review & Complete Registration

Policy Discovery & Selection

Home screen shows product categories

Policy list with premium & benefits

Compare policies

View details

Policy Purchase Workflow

As shown in UI screens (multi-step forms)

Select policy

Enter personal details

Upload documents

Review summary

Payment

Confirmation

Claims Journey

Enter policy number

Upload claim documents

Receive updates

Policy History & Tracking

Past policies

Current active policies

Claim history

Renewal alerts

Functional Requirements

Registration & Authentication Business Requirements

Users must verify identity via mobile OTP

KYC should follow IDRA guidelines

Nominee capture is mandatory for Life/Health insurance

Functional Requirements

OTP generation & validation

Duplicate number detection

Mandatory fields: name, DOB, NID/passport, nominee data

Optional: health declarations, lifestyle questions

Policy Marketplace Business Requirements

Users must browse and compare policy options with transparent benefits

Pricing must match insurer-approved rates

Functional Requirements

Product list API to fetch available plans

Sorting and filtering (premium, coverage type)

Compare up to 3 policies

Detailed policy information page

Policy Purchase Module Business Requirements

Paperless onboarding

Real-time underwriting for low-risk policies

Automated document verification

Functional Requirements

Multi-step purchase journey (aligned with screens)

Document upload: NID, photo, medical file (if needed)

Payment gateway integration

Auto-generation of digital policy document

Claims Module Business Requirements

Simple and fast claim initiation

Digital evidence submission

1–3 working days settlement (where possible)

Functional Requirements

Claim request form (policy number prefill)

Document upload

Claims history

Status notifications

Admin approval workflow

Policy Management Functional Requirements

View active policies

Renewals

Cancellation request

Download policy documents

Policy modification requests (address, nominee)

Notification System Requirements

Push, SMS, email

Renewal reminders

Claim updates

System alerts

Marketing campaigns

Admin Portal Key Capabilities

User management

Product & pricing management

Claims management

Dashboard & reporting

API integration controls

Non-Functional Requirements (NFR)

Performance

App load time < 3 seconds

Payment processing < 10 seconds

Claim submission processing < 5 seconds

Security (IDRA-aligned)

AES-256 encryption for data at rest

TLS 1.3 for data in transit

MFA for admin users

Secure OCR storage for sensitive documents

Regulatory Compliance

IDRA reporting frameworks

Automated KYC validation

AML/CFT compliance

Scalability

Cloud-native (AWS/Azure)

Auto-scaling microservices 
> 🔴 **ISSUE:** "Auto-scaling microservices" is massive Over-Engineering for an MVP.
> 🟦 **BENCHMARK:** **PolicyBazaar** relied on a monolithic PHP/MySQL stack for its first ~8 years (handling millions of users) before moving to Microservices in 2017. 
> 🟢 **ALTERNATIVE:** Use a **Modular Monolith** architecture. It is easier to build, deploy, and debug.

Availability

99.5% uptime

Disaster recovery with < 15 min RTO

System Integrations Core Integrations

| System | Purpose |
| --- | --- |
| Insurer API | Premium computation, underwriting, policy issuance | 
> ⚠️ **RISK:** Dependency on external API maturity.

| Payment Gateways | bKash, Nagad, Visa/Mastercard |
| LabAid Hospital Systems | Cashless IPD/OPD | 
> ⚠️ **RISK:** Legacy integration.

| CRM | Support tickets, communication |
| Notification Engine | SMS/Email/Push |
| Analytics Engine | KPI tracking, AI models |

Data Requirements Customer Data

Personal info, nominee info, contact info

KYC documents

Health declarations

Policy Data

Product details

Premium tables

Policy status & history

Claim Data

Claim type, documents, status trail

System Data

Audit logs

API logs

Notification logs

Acceptance Criteria System-wide

All journeys must be fully digital (registration → purchase → claim)

0% paper dependency for standard products

98%+ successful OTP delivery rate

Claims module must allow successful upload of documents

Business

All policies must be issued with valid insurer authorization

Payment confirmation must auto-issue policy certificate

Future Enhancements (2027–2028) AI & Predictive Systems

AI underwriting engine

Driving score & motor UBI models

Health scoring integration

IoT Integration

Device diagnostics

Vehicle telematics

Cattle & crop sensors

Super-App Features

Telemedicine

Health passport

Pharmacy ordering

Wellness rewards

Partner API Marketplace

E-commerce

Telcos

Ride-sharing

MFIs & rural networks

BRD Sign-off

This Hybrid BRD reflects the complete functional, business and technical requirements for the modern InsureTech platform designed for Bangladesh’s evolving digital ecosystem.

Approval of this document authorizes the development team to proceed with:

UI/UX finalization

Technical architecture design

System development

Integration workstreams
