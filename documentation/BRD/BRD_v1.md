BUSINESS REQUIREMENTS DOCUMENT (BRD)
Labaid InsureTech Company Limited


1.	Document Purpose
This BRD defines the business needs, functional specifications and technical requirements for developing an advanced InsureTech platform that aligns with:
•	Labaid Group’s health ecosystem
•	IDRA regulatory framework
•	Modern InsureTech/HealthTech best practices
•	Bangladesh’s digital financial landscape
•	The 2026–2028 business roadmap

This BRD shall guide technology vendors, internal development teams, UX teams and business stakeholders throughout platform design, development, testing and deployment.

2.	Business Overview
Labaid InsureTech aims to modernize insurance distribution through:
•	Digital-first experiences
•	AI-driven underwriting & claims automation
•	Mobile app–based onboarding
•	Micro-insurance and inclusive products
•	Integration with LabAid health services
•	Partner-driven embedded insurance
The platform will support Life, Health, General, Agriculture and Device Insurance through a unified, seamless mobile experience.

3.	Business Objectives
Short-term (2026)
•	Launch core InsureTech platform
•	Enable digital onboarding, policy purchase and claims
•	Achieve 40,000+ active policies
•	Integrate with insurers and payment gateways
Mid-term (2027)
•	Introduce AI underwriting
•	Automate 80% of claims
•	Achieve break-even (Q3 2027)
•	Launch Super-App 2.0
Long-term (2028)
•	Become Top 3 InsureTech platform in Bangladesh
•	Expand regionally (Nepal, Bhutan, Maldives)
•	Implement predictive risk scoring, behavioral pricing, IoT integration
 
4.	Project Scope
4.1	In-Scope
•	Customer mobile app (iOS/Android)
•	Web admin portal
•	Agent/Partner portal
•	Policy purchase, renewal, cancellation
•	Digital KYC & verification
•	Claims submission & automation
•	Payment integrations (bKash, Nagad, cards)
•	Notification system (SMS, push, email)
•	Integration with insurer APIs
•	Integration with LabAid hospital systems
4.2	Out-of-Scope (Phase 1)
•	Full AI-driven underwriting (Phase 2: 2027)
•	IoT/Telematics integration (Phase 3: 2028)
•	Cross-border insurance sales (2028)

5.	User Groups
User Type	Description
Customer / Policy Buyer	App users purchasing policies
Partner/Agent	Distributors (hospitals, MFS, telcos, e-commerce)
Insurer Underwriter	Validates policy and approves underwriting
Claims Officer	Processes and approves claims
Super Admin	Full platform control
Support/Call Center	Customer assistance and issue resolution


6.	High-Level User Journey Flows
6.1	Registration & Login
•	Signup/Login
•	OTP Verification
•	Personal Info
•	Nominee Info
•	Review & Complete Registration
6.2	Policy Discovery & Selection
•	Home screen shows product categories
•	Policy list with premium & benefits
•	Compare policies
•	View details
6.3	Policy Purchase Workflow
As shown in UI screens (multi-step forms)
1.	Select policy
2.	Enter personal details
 
3.	Upload documents
4.	Review summary
5.	Payment
6.	Confirmation
6.4	Claims Journey
•	Enter policy number
•	Upload claim documents
•	Receive updates
6.5	Policy History & Tracking
•	Past policies
•	Current active policies
•	Claim history
•	Renewal alerts

7.	Functional Requirements
7.1	Registration & Authentication Business Requirements
•	Users must verify identity via mobile OTP
•	KYC should follow IDRA guidelines
•	Nominee capture is mandatory for Life/Health insurance
Functional Requirements
•	OTP generation & validation
•	Duplicate number detection
•	Mandatory fields: name, DOB, NID/passport, nominee data
•	Optional: health declarations, lifestyle questions
7.2	Policy Marketplace Business Requirements
•	Users must browse and compare policy options with transparent benefits
•	Pricing must match insurer-approved rates
Functional Requirements
•	Product list API to fetch available plans
•	Sorting and filtering (premium, coverage type)
•	Compare up to 3 policies
•	Detailed policy information page

7.3	Policy Purchase Module Business Requirements
•	Paperless onboarding
•	Real-time underwriting for low-risk policies
•	Automated document verification
Functional Requirements
•	Multi-step purchase journey (aligned with screens)
•	Document upload: NID, photo, medical file (if needed)
 
•	Payment gateway integration
•	Auto-generation of digital policy document
7.4	Claims Module Business Requirements
•	Simple and fast claim initiation
•	Digital evidence submission
•	1–3 working days settlement (where possible)
Functional Requirements
•	Claim request form (policy number prefill)
•	Document upload
•	Claims history
•	Status notifications
•	Admin approval workflow
7.5	Policy Management Functional Requirements
•	View active policies
•	Renewals
•	Cancellation request
•	Download policy documents
•	Policy modification requests (address, nominee)

7.6	Notification System Requirements
•	Push, SMS, email
•	Renewal reminders
•	Claim updates
•	System alerts
•	Marketing campaigns
7.7	Admin Portal Key Capabilities
•	User management
•	Product & pricing management
•	Claims management
•	Dashboard & reporting
•	API integration controls
8.	Non-Functional Requirements (NFR)
8.1	Performance
•	App load time < 3 seconds
•	Payment processing < 10 seconds
•	Claim submission processing < 5 seconds
8.2	Security (IDRA-aligned)
•	AES-256 encryption for data at rest
 
•	TLS 1.3 for data in transit
•	MFA for admin users
•	Secure OCR storage for sensitive documents
8.3	Regulatory Compliance
•	IDRA reporting frameworks
•	Automated KYC validation
•	AML/CFT compliance
8.4	Scalability
•	Cloud-native (AWS/Azure)
•	Auto-scaling microservices
8.5	Availability
•	99.5% uptime
•	Disaster recovery with < 15 min RTO
9.	System Integrations Core Integrations
System	Purpose
Insurer API	Premium computation, underwriting, policy issuance
Payment Gateways	bKash, Nagad, Visa/Mastercard
LabAid Hospital Systems	Cashless IPD/OPD
CRM	Support tickets, communication
Notification Engine	SMS/Email/Push
Analytics Engine	KPI tracking, AI models

10.	Data Requirements Customer Data
•	Personal info, nominee info, contact info
•	KYC documents
•	Health declarations
Policy Data
•	Product details
•	Premium tables
•	Policy status & history
Claim Data
•	Claim type, documents, status trail
System Data
•	Audit logs
•	API logs
•	Notification logs
11.	Acceptance Criteria System-wide
•	All journeys must be fully digital (registration → purchase → claim)
•	0% paper dependency for standard products
 
•	98%+ successful OTP delivery rate
•	Claims module must allow successful upload of documents
Business
•	All policies must be issued with valid insurer authorization
•	Payment confirmation must auto-issue policy certificate
12.	Future Enhancements (2027–2028) AI & Predictive Systems
•	AI underwriting engine
•	Driving score & motor UBI models
•	Health scoring integration
IoT Integration
•	Device diagnostics
•	Vehicle telematics
•	Cattle & crop sensors
Super-App Features
•	Telemedicine
•	Health passport
•	Pharmacy ordering
•	Wellness rewards
Partner API Marketplace
•	E-commerce
•	Telcos
•	Ride-sharing
•	MFIs & rural networks

13.	BRD Sign-off
This Hybrid BRD reflects the complete functional, business and technical requirements for the modern InsureTech platform designed for Bangladesh’s evolving digital ecosystem.
Approval of this document authorizes the development team to proceed with:
•	UI/UX finalization
•	Technical architecture design
•	System development
•	Integration workstreams
