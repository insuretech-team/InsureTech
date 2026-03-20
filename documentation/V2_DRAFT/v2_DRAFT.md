





 



	

Version: 2.0
Status: Draft
Control Level: A
________________________________________

  Software Requirements Specification (SRS) 
  For

 Labaid Insure Tech Platform 








Revision History

SL	Date	Revised By
01		
02	12_06_02025	Faruk Hannan





Table of Contents
1.	Introduction
2.	MarketContext
3.	Overall Description
4.	System Features & Functional Requirements
5.	External Interface Requirements
6.	Non-Functional Requirements
7.	Data Model & Storage Requirements
8.	Security & Compliance Requirements
9.	Performance & Scalability Requirements
10.	Operational Requirements & Support
11.	Acceptance Criteria & Test Summary
12.	Traceability Matrix & Change Control
13.	Appendices (Use Cases, Screens, File References)









1. Introduction
1.1 Purpose
This SRS documents detailed system-level requirements for the Labaid  Insure Tech platform (mobile apps, partner portals, admin portal, and backend services). Its purpose is to provide an unambiguous specification for design, development, integration, testing, and deployment and operations teams.
1.2 Scope
      The system will enable digital onboarding, product discovery, policy purchase, digital KYC, payment processing, claims submission and tracking, partner integrations, and admin/insurer workflows. This SRS covers Phase 1 (core digital capabilities) while flagging Phase 2/3 enhancements (AI underwriting, Voice aided guidance, IOT based connectivity) for future releases. See Business Plan for strategic context.

1.3 Market Context
       With a low penetration rate of insurance compared to other countries in the region and considering mass literacy level and lack of awareness about insurance, a very simplified onboarding flow with step-by-step explanations and visual queue are required for users UI/UX space. Also considering rural internet network constraints, overall speed (Minimum 3G), a minimum spec of edge devices and also lack of availability of local cloud storage support within the country for reducing latency, specific sets of requirements are set for overall app/technology performance and reactivity.
1.4 Definitions, Acronyms, and Abbreviations
●	IDRA — Insurance Development & Regulatory Authority (Bangladesh)
●	KYC- Know your customer.
●	KYB-Know your Business
●	FCR - Financial Condition Report (IDRA requirement)
●	CARAMELS - Capital Adequacy, Reinsurance Arrangements, Management, Earnings, Liquidity & Asset quality, Sensitivity 
●	SAR - Suspicious Activity Report (AML requirement)
●	BFIU - Bangladesh Financial Intelligence Unit
●	MFS — Mobile Financial Service (bKash, Nagad)
●	API — Application Programming Interface
●	EHR — Electronic Health Record
●	ZHCT — Zero Human Touch Claims
●	UBI — Usage Based Insurance
●	IAM-Identity and Access Management
●	ACL: Access Control List
●	RBAC-Role based Access Control
●	ABAC- Attribute based Access control
●	UAT: User Acceptance Testing
●	MOU- 











2. Overall Description
2.1 Product Perspective
The platform is cloud-native, secure-first, micro services-based, event-driven and mobile-friendly. 
Components:
•	Mobile apps (iOS/Android) — Customer experience (screens & flows per uploaded designs).
•	Partner/Agent Portal — onboarding & embedded flows.
•	Admin Portal — product, pricing, user and claims management.
•	Back-end Services — authentication, authorization, partner/agent management, policy engine, profiling service,  contract management, payment gateway adapters, risk management, fraud detection, renewal service, transactional debit-credit service, storage service, LLM multi agent network , ai assistant service, MCP servers, iot-broker service, orchestration service, notification service, ticketing service(customer support) , analytics, reporting .
•	Integrations — Insurer APIs, LabAid systems, MFS/telco APIs.
•	Flow reference: userflow diagram.















 2.2 User Classes & Characteristics
        Broader classification: Customers and Stakeholders
•	Customers : Mobile-first, diverse digital literacy 
•	Type1 : Urban and Digitally Literate  -> Self-service mobile app users
•	Type2 : Semi-Urban/Agent-assisted->  payment with agent
•	Type3: Rural/Agricultural -> voice assisted workflow             
•	Stakeholders:
•	Partners:  System on boarded, collaborative  organization (MFS, hospitals, e-commerce)  
•	Insurer Agents:   Staff at partner organizations  under ACL
•	Insurer Underwriters: Internal to partner’s organization insurers (API consumers).
•	Independent Agent: on boarded, certified payment agents in remote locations(dashboard)
•	Vendor : independent service provider to insure Tech as a client (dashboard)
•	System Admin :   Root level to cloud provider gateway ( IAM , ACL ,Policy)
•	Repository Admin: Root  Level to repo and droplet ( Pull, Merge ,Deploy access )
•	Database Admin: Root Level to Databases /Storage (  Full management access)
•	Business Admin: Root level to portal (policy management, product updates).
•	Focal Person: Root level  to Partner management ( partner verification, onboarding ,dispute  )
•	Partner Admin : Partner Side Root Tenant ( IAM, ACL control)
•	Dev : user level ( software development)          
•	Support/Call Centre: Assisted onboarding & claims help(dashboard)  














2.3	Operating Environment
•	Cloud: AWS/Azure preferred; region: Bangladesh (or compliant region like Singapore not more than 80ms added latency).
•	Mobile: iOS 13+/Android 9+. Minimum support with 4GB RAM 
•	Network :Minimum 3G (EDGE will not be supported) with frequent network switching
•	App Download: Must have lowest download size as low as 10 Megabyte , gradual download and caching of data run up to maximum 100 Megabyte
•	Offline Capabilities: Must have minimum  persistent view of  critical data 
•	Low Bandwidth Mode: With minimum feature with device level cached static data .
•	Low Resource Mode: Must run with minimum memory and low power mode.
•	Browsers for portals: Chrome, Brave, Firefox, Edge (latest two versions).
•	     USSD/SMS fallbacks: will have retry and acknowledgement        support












  
3  System Features & Functional Requirements
Each requirement below has a Function Group unique id (e.g., FG-1 ) and Functional requirements  unique ID (e.g., FR-1).
Functionality Deploy Phases :  Phases 1,Phase 1.5, Phase 2 ,Phase 3
Priority levels: M = Mandatory (Phase 1), D = Desirable (Phase 1.5), F = Future (Phase 2/3).
3.1 Functional Requirements
ID	Requirement Description	Priority
FG-001	3.11 -IAM Group ( Authentication )	M
FR-001	The system shall allow Customers  to  register and login via mobile number with OTP verification with minimum data(name, date of birth, mobile no, email(optional)	M
FR-002	The system shall auto update OTP from SMS  with one time read permission  from device 	D
FR-003	The system shall have SMS OTP  failure fallback with persistent form state  ,retry loop and resend button 	M
FR-004	The system shall have image capture step for user verification step1 	M
FR-005	The system shall perform login step  via e-KYC  as per Office of the Controller of Certifying Authorities Bangladesh	D
FR-006	The system shall automate with voice-assisted workflow on landing page for Type 3 customer to guide to login.	F
FR-007	The system shall have option for logged in customer to link with social account via  oAUth2 (for Type1 ,Type2 customer (Google, Facebook)  to finalize profile	D
FR-008	The system shall have mandatory NID upload for user verification step 2	M
FR-009	The system shall digitized  and Process NID data  	D
FR-010	The System shall show green check mark for user verification step	D
FR-011	The system shall have option of Insure tech employed Focal person to onboard a pre-approved, MOU-signed partner on behalf of the partner, setup basic ACL, and provide one Partner Admin with temporary password. 	M
FR-012	The system shall have stakeholders registration via OPEN ID Identity provider 	M
FR-013	The system shall have stakeholders  registration via SAML Identity provider 	D
FR-014	The system shall have stakeholders  registration via SSO provider 	D
FR-015	The system shall have stakeholders  registration via email with email verification	M
FR-016	The system shall lock verified stakeholders  KYB from update 	M
FR-017	The system will allow all logged in stake holder to change password	M
FR-018	The system shall have option of account recovery question	D
FR-019	The system shall have 2FA authentication via mobile number SMS/ OTP 	M
FR-020	The system shall have account recovery option via email /mobile and 2nd step with security question	M
FR-021	The system will allow verified stakeholders assisted  by Focal Person for migration  to a new account incase verification data needs to be updated  	D
FR-022	The system shall maintain a dashboard for stakeholders to view all ongoing verification workflows categorized by account type	M
FR-023	The system shall allow users to update non KYC/KYB  user data e.g employee data	M
FR-024	The system shall define mandatory KYC  fields for Customer  identify verification 	M
FR-025	The System shall define mandatory KYB fields  and KYC fields ( if Admin ) for stakeholders verification	M
FR-026	The system will block duplicate account for KYC conflict 	M
FR-027	The system shall have automatic merge try for KYC conflict 	D
FR-028	The system shall have manual merge step for KYC conflict.	M
FR-029	The system will only allow manual merge step for KYB conflict	M
FR-030	The system will block user for multiple wrong attempt of login for a short period	M
FR-031	The system will ban for identities from login for lifetime for malicious attempts	M
FR-032	The system will allow API Keys  with limited lifetime for 3rd party system	D











ID	Requirement Description	Priority
FG-002	3.12 Verification and Authorization 	M
FR-033	The system shall verify Customer with minimum data -captured image, NID, mobile number	M
FR-034	The system shall ask for additional data(e.g medical history)  as per  the insurance product provider policy requirements  with digital copy upload 	M
FR-035	The system shall allow medical data verification report with authorized government certified medical officer ??	D
FR-036	The system shall have image recognition system for user  verification 	D
FR-037	The system shall perform OCR from uploaded document to make automated EHR /Medical history for Customer	F
FR-038	The system will take Customer geo location data , IP tracking,  IOT data to avoid ID theft	F
FR-039	The system shall have minimum KYB for partner as per Insure tech internal business policy (Like whether to add bank solvency etc.)	M
FR-040	The system will categorized stakeholders as per financial and regulatory compliance 	M
FR-041	The system shall provide partners sub domain 	F
FR-042	The system shall provide partner own portal with dashboard	M
FR-043	The system shall provide partner to make internal ACL with either RBAC or ABAC with resource , policy , permission set 	M
FR-044	The system shall have Insure tech internal ACL as per casbin  RBAC with domains/tenants Rule for .NET (https://casbin.org/docs/supported-models)	M
FR-045	The system will provide Partner Admin of Full CRUD access of partner domain data (employee, policy, product, price etc.) 	D
FR-046	The system will provide access Partner to dedicated data storage ( as per InsureTech Policy)	D
FR-047	The system will follow server side session management for stakeholders	M
FR-048	The system  will follow oauth2 and token refresh for Customer	M
FR-049	The system will  try to make random token refresh and key rotation  to disallow misuse	D
FR-050	The system will only expose secure proxy server IP (Cloud flare) to user  	D
FR-051	The system shall invalidate user session/token immediately after logout	M
		
		
		








ID	Requirement Description	Priority
FG-003	3.13-Product Catalog & Policy Discovery	M
FR-055	The System shall allow Partners to upload Product catalog, Product with data Metadata , coverage, premiums, insurer, terms and conditions	M
FR-056	The system shall allow partners to list products for approval of insureTech	M
FR-057	The system shall allow partners to check product approval status and if unapproved will be prompted with disqualification reason	M
FR-058	The system shall allow partners to pick live date on site for any specific product	M
FR-059	The system shall allow partners to edit a live product but will go through same re approval process	M
FR-060	The system will not allow partners to delete  a product after it is live ,partners will take approval as per InsureTech management for such change and will settle any claim during that product purchase  	M
FR-061	The system will categorized products as per partners provided tag and meta data	M
FR-062	The system will allow users to compare same category product side by side with premium and coverage	M
FR-063	The system will have filtering policy to show particular products/a product as per category, premiums ,coverage, 	M
FR-064	The system will auto pick any product based on InsureTech internal rating as featured, recommended 	M
FR-065	Product Detail: Full policy wording, exclusions, and illustrative premium breakdown	M


ID	Requirement Description	Priority
FG-004	3.14-Policy Purchase & Issuance	M
FR-066	The system will have a minimum number of fixed step for any product purchase as minimum-steps	M
FR-067	If addition flow is required more than minimum-steps it will be added later always maintaining same serial 	M
FR-068	The system product purchase minimum steps are add personal details (one time only) , add product to buy ,add nominee (will be prompted for previously set nominee) ,add documents, additional document in always same serial with document types  .then purchase order will set for verification	M
FR-069	The system will try immediate verification of documents with visual progress status (e.g green circle)   	D
FR-070	The system will prompt customer for product purchase order verification with modal notification /voice and after that  take on next payment page	M
FR-071	The system will allow customer to make payment  as per listed bKash, Nagad, card (PCI-DSS compliant), bank transfer (api version details during detail engineering)	M
FR-072	The system will allow customer to allow payment via agent and add transaction id and upload payment details	D
FR-073	The system shall verify/show premium calculation as per partners provided product data 	M
FR-074	The system shall show premium calculations as per inbuilt pricing engine	D
FR-075	The system shall show premium calculation of Partner API as per API data schema	M
FR-076	The system shall support promo code for product discount	M
FR-077	The system will support partners to add /remove discounts 	M
FR-078	The system will maintain transparent ,indestructible transaction record ,debit-credit with history per partner per product per customer  in purpose built database(Tiger beetle or equivalent )called as  Insurance Record DB (IRD)	D
FR-079	The system will generate a digital policy certificate as artifacts as pdf with QR code and store in user account. 	M
FR-080	The digital certificate will have customer eKYC digital sign injected 	D
FR-081	The system will have verification support for any required authority to verify the generated certificate	M
FR-082	The system will make an additional ,irrevocable  digital tri party contract with insureTech , partner, Customer  for each record in IRD will be called as Insurance  Contract	D
FR-083	The system shall generate web3.0 and  block chain based Insurance Contract	F
FR-084	The system shall  maintain linking record with all older purchased products  for any  new purchase  of customers	M
FR-085	The system will allow customer to purchase product  on behalf  of family member 	D
FR-086	The system will not ask for any data that customer already uploaded and verified in history	M
FR-087	The system will allow customer to progress medical history by updating with new record 	M
FR-088	The system shall produce proper detail message as system SMS for customer if any document verification fails	M
FR-089	The system shall have purpose built AI Chatbot to assist customer during product search, selection, purchase, verification ,payment stage	F

ID	Requirement Description	Priority
FG-005	3.15 Claims Management	M
FR-090	The system will have  a Form of fixed steps of Insurance Claims - Policy prefill, claim reason selection, document upload (images, bills	M
FR-091	The system will show success screen after initial screening of documents after submission.	M
FR-092	The system will make a digital hash to identify each unique submission and keep as internal record 	M
FR-093	The system will automatically notify partner for new claims and Start common  status dashboard with customer	M
FR-094	The system will show status tracking of the claim in the dashboard to customer 	M
FR-095	The system will have a chat interface with customer and partner agent in dashboard	D
FR-096	The system will have a web RTC based video call support with customer and partner agent in dashboard	D
FR-097	The system will have agents note support in dashboard	M
FR-098	The system will allow partner to add verification details for insureTech record keeping	M
FR-099	The system will allow partner to make initial approval or rejection request to inscureTech 	M
FR-100 	The system will allow Business Admin and Focal Person joint approve or reject upon internal verification	M
FR-101	The system will automate payment process as per customer selected payment channel 	M
FR-102	The system shall have auto verification, claim approval and payment for small claims with pre agreement as per Partner, Insurer and InsureTech	D
FR-103	The system will have Auto fraud detection system to detect  frequent claims and changed uploads by customer	D
FR-104	The system will provide system warning to customer as per InsureTech policy	M
FR-105	The system will auto revoke access to customer as per InsureTech policy for suspicious activity	D
FR-106 	The system will maintain proper balance sheet on Customer, Partner, Agent and InsureTech level for various selected period of time . 	D


ID	Requirement Description	Priority
FG-006	3.16 Policy Management & Renewals	M
FR-107	The system will provide customer a policy dashboard after first product purchase	M
FR-108	The system will show active and past policies, renewal prompts, policy pdf download  in customer dashboard	M
FR-109	The system will provide customer  two option for un-payed policies renewal :auto and manual .	M
FR-110	The system will guide customer in clear steps for manual policy renewal 	M
FR-111	The system will allow customer to partially adjust policy data 1.Current address 2. Nominee	M
FR-112		D
ID	Requirement Description	Priority
FG-007	3.17 Notifications & Communication	M
FR-113	The system will have KAFKA (or equivalent) event driven  notification and orchestration  system triggering following channels like  : In System PUSH (portal & mobile) ,  SMS, Email	M
FR-114	 The system will auto trigger notification for verification, purchase confirmation, claims updates, renewal reminders	M
FR-115	The system will allow user to be notified for additional events as per provided list(e;g new product launch)	D
FR-116	The system will provide partners to create secondary marketing notification for customers by filtering age,sex, groups	D
FR-117	The system will allow customers with opt in opt out policy for additional notification 	M
FR-118	The system will support customer mobile mute mode with minimum text notification (avoiding on system push for edge devices) 	M
FR-119	The system will allow customer support call within app	D
FR-120	The system will provide common Q&A within app for customer knowledge 	M
FR-121	The system shall provide AI chatbot for customer support and automatic ticketing for unresolved issue	F
FR-122	The system will Provide common question and answer form for customer support call center(Vendor)	D
FR-123	The system will provide ticket form fill up option with auto recorded customer support call 	D



ID	Requirement Description	Priority
FG-008	3.18-Admin & Reporting	M
FR-124	The system will maintain separate dashboard for different admin as per defined role , resources, policies by InsureTech management	M
FR-125	The system will maintain strict 2FA for Any admin level access	M
FR-126	The system will make dynamic content view based on context , workflow, and data access available  for different admin 	M
FR-127	The system will provide  common dashboard for domian level, team level	D
FR-128	The system will allow internal approval , task issuing to internal users	D
FR-129	The system will haves separate  user management tab	M
FR-130	The system will have separate product management tab	M
FR-131	The system will have separate claim management tab	M
FR-132	The system will have assigned task tab per user	M
FR-133	The system will have option of reporting as per access to users	M
FR-134	The system will have minimum reporting options as Daily sales, claims ratio, partner performance, policy counts and KPIs (aligned to business plan targets).	M





ID	Requirement Description	Priority
FG-009	3.19-  Partner Portal (Hospitals, Ecommerce) 	M
FR-135	The system will provide special dashboard for hospital partner to be able to initiate purchase insurance for any customer on behalf of them  	M
FR-136	The system shall have option to transfer customer record with partner provided api and purchase id and authentication token 	D
FR-137	The system shall provide API  for Ecommerce partner to embed product in page and special dashboard to purchase product from InsureTech	M
FR-138	The system shall provide sandbox for 3rd part developer for authentication and testing payment and purchase process 	F
FR-139	The system will provide any partner onboarding analytics, leads, Commission statements in their dedicated dashboard 	M
FR-140	The system will provide any partner onboarding analytics, leads, Commission statements via API  	M


ID	Requirement Description	Priority
FG-010	3.20-  Audit & Logging	M
FR-141	The system shall maintain immutable   logs for critical actions such as : policy issue, claim approval ,claim rejection , payment and dispute 	M
FR-142	The system shall maintain data retention policy up to minimum of 20 years  to maintain records for regulatory compliance	D
FR-143	The system will track each logged in user for auxiliary actions and will have additional  data logs	D
FR-144	The system will allow partner to maintain additional data logs as per customer and InsureTech MOU	F
FR-145	The system will provide special portal to regulatory body to access requested data as per  policy and law of regulatory bodies 	M

ID	Requirement Description	Priority
FG-011	3.21-  User Interface	M
FR-146	The system shall maintain similar User interface with different operating system 	M
FR-147	The system shall provide smart data  widget for mobile user 	D
FR-148	The system shall  provide voice assisted workflow for type 3 user 	F
FR-149	The system shall provide desktop first web UI for portal 	M
FR-150	The system shall  take minimum permission for all service from user device (e’g camera, one time msg read)	M









FG-012	3.22-  API Design and Data Flow 	M
FR-151	The system will maintain two level of API private and public 	M
FR-152	The system will have three level of private API  as category 1 ,category 2 and category 3 with secure middle layer 	D
FR-153	The system will use protocol buffer and GRPC based communication for category 1 API  between  gateway and micro services  in http with system admin middle layer	F
FR-154	The system will use Graphql based API  as Category 2  for gateway and customer device with jwt token system and oauthv2 based  security middle layer 	M
FR-155	The system will use RestFul API  with JSON with compliance of OPEN API standard  as Category 3 for 3rd party Integration with server side auth and will provide user public documentation , wiki , sandbox and a mock server	D
FR-156	The  system will provide public Restful  API with   JSON with compliance of OPEN API standard   for product search , product list and view  for live products etc	M
FR-157	The system will only expose proxy server (Cloud flare) and entry node IP (NGINX ) to public space	M
FR-158	The system will block all access to  micro services for all available API except category 1 	M
FR-159	The system will use InsureTech Internal protocol for IOT based data extraction and data binding 	F
FR-160	The system will consolidate , annotate , add context and process and save data  within regulatory limitation to train internal AI agents	F
FR-161	The system will generate statistics, prediction based on big data and provide special projection data to partners  given special agreement with Insuretech	F

4. Non-Functional Requirements
ID	Requirement Description	Priority
NFR-001	The system shall ensure 99.5% uptime to meet critical business needs.	M
NFR-002	The system shall encrypt (TLS 1.3 and AES-256)  all sensitive data in transit and at rest.	M
NFR-003	The system shall handle up to 1,000 concurrent users without performance degradation.	M
NFR-004	The system shall handle up to 5,000 concurrent users without performance degradation.	D
NFR-005	The system shall handle up to 10,000 concurrent users without performance degradation.	F
NFR-006	The system shall have a median response time of less than 1 seconds and not more than 2 s for 95% of user actions	M
NFR-007	Disaster Recovery: RTO < 15 minutes for critical services, RPO < 1 hour	M
NFR-008	The system will have category 1 API response time < 100ms 	D
NFR-009	The system will have category 3 API response time < 200 ms 	D
NFR-010	The system will have category 2 API response time < 2 sec	D
NFR-011	The system will have public API response time < 1sec 	D
NFR -012	The system will ensure minimum app startup time <  5 sec	D
NFR-013	The system will support 1 million active user and 200K active policies in 24 months with auto  vertical scaling	F
NFR-014	The system will use separate secret vault (azure,hashi corp) for critical api keys	D
NFR-015	The system will use auto key rotation and special techniques  for token based auth	M
NFR-016	The system will follow PCI-DSS compliance for card flows 	M
NFR-017 	The system will have AML/CFT detection hooks and SAR reporting	D
NFR-018	The system will have IDRA reporting capabilities following IDRA data format 	M
NFR-019	The system will have regular penetration test with whitelisted international security experts	D
NFR-020	The system will have regular security audits from various security auditors and regulatory bodies and maintain compliance	D
NFR-021	The system will have automatic switch option for user to different mobile device with data migration	F
NFR-022	The system will have multi language support (English+ Bengali)	D
NFR-023	The system shall support Bangladesh mobile network conditions (3G minimum, frequent network switching)	M
NFR-024	The system shall compress images client-side from 5MB to 1-2MB before upload	M
NFR-025	The system shall support chunked file upload with resume capability for files >1MB	M
NFR-026	The system shall work on low-end Android devices (2-4GB RAM, Android 9+)	M
NFR-027	The system mobile app size shall be <25MB initial download, progressive loading up to 100MB	M
NFR-028	The system shall support offline form completion with background sync when network available	D
NFR-029	The system shall mask PII in all logs (show only last 3 digits of NID, mask phone numbers)	M
NFR-030	The system shall maintain immutable audit logs for policy issuance, claims, payments for 20 years	M
NFR-031	The system shall achieve 80% unit test code coverage minimum	M
NFR-032	The system shall pass OWASP Top 10 security checks (via SAST/DAST tools)	M
NFR-033	The system shall support JWT token expiry of 15 minutes with refresh token mechanism	M
NFR-034	The system shall implement rate limiting (100 requests/minute per user for public APIs)	M
NFR-035	The system shall use database connection pooling with max 100 connections per service	D
NFR-036	The system shall implement database read replicas for reporting queries	D
NFR-037	The system shall cache frequently accessed data (product catalog) with 5-minute TTL	D
NFR-038	The system shall support blue-green deployment with one-command rollback	M
NFR-039	The system shall implement circuit breakers for all external API calls (insurers, payment gateways)	M
NFR-040	The system shall retry failed API calls with exponential backoff (max 3 retries)	M
NFR-041	The system shall timeout external API calls after 30 seconds	M
NFR-042	The system shall support idempotent payment API calls to prevent duplicate charges	M
NFR-043	The system shall validate Bangladesh phone number format (+880-1XXX-XXXXXX)	M
NFR-044	The system shall validate NID checksum for 10-digit and 13-digit smart NIDs	M
NFR-045	The system shall support Bengali Unicode (UTF-8) in all text fields	M
NFR-046	The system shall render Bengali text correctly on all supported devices	M
NFR-047	The system shall implement MFA (OTP via SMS) for all admin logins	M
NFR-048	The system shall lock accounts after 5 failed login attempts for 30 minutes	M
NFR-049	The system shall auto-logout inactive sessions after 15 minutes	M
NFR-050	The system shall support CORS for partner portal integrations with configurable allowed origins	D
NFR-051	The system shall implement API versioning with backward compatibility for at least 6 months	D
NFR-052	The system shall provide API documentation via Swagger/OpenAPI 3.0 specification	M
NFR-053	The system shall support webhook delivery with retry mechanism (3 attempts over 24 hours)	M
NFR-054	The system shall verify webhook signatures using HMAC-SHA256	M


5. External Interface Requirements

5.1 Payment Gateway Integrations

5.1.1 bKash Integration
Priority: M
API Version: bKash Merchant API v1.2.0-beta or latest stable
Endpoints:
•	Create Payment: POST /checkout/payment/create
•	Execute Payment: POST /checkout/payment/execute  
•	Query Payment: GET /checkout/payment/query/{paymentID}
•	Refund: POST /checkout/payment/refund

Authentication: Bearer token (Grant Token from /checkout/token/grant)
Request Timeout: 90 seconds (bKash can be slow)
Retry Logic: 3 attempts with 5-second intervals for network errors only (not for business errors)
Idempotency: Use merchantInvoiceNumber as idempotency key
Webhook: POST callback to our endpoint with signature verification
Webhook Security: 
•	IP whitelist: bKash provides static IPs
•	Signature verification using shared secret
Error Handling:
•	2051: Insufficient balance → Show user-friendly message
•	2062: Transaction limit exceeded → Suggest alternative payment
•	Network timeout → Retry with exponential backoff
Reconciliation: Daily settlement file (CSV) via SFTP, match against our transaction records
Refund SLA: 7-10 business days per bKash policy
MDR (Merchant Discount Rate): 2.5% (to be negotiated)
Test Environment: bKash Sandbox with test wallet credentials

5.1.2 Nagad Integration  
Priority: M
API Version: Nagad Merchant API v2.0
Endpoints:
•	Initialize Payment: POST /api/dfs/check-out/initialize
•	Complete Payment: POST /api/dfs/check-out/complete
•	Verify Payment: GET /api/dfs/verify/payment/{paymentRefId}
•	Refund: POST /api/dfs/refund

Authentication: PGP signature + timestamp validation
Request Timeout: 60 seconds
Retry Logic: 2 attempts for network errors
Idempotency: Use orderId as unique transaction identifier
Webhook: Callback URL provided during initialization
Webhook Security: PGP signature verification
Error Codes:
•	E1003: Invalid merchant → Contact Nagad support
•	E2001: Duplicate transaction → Return existing transaction status
Settlement: T+1 daily settlement
MDR: 1.99% (lower than bKash but less reliable historically)
Test Environment: Nagad Sandbox

5.1.3 Card Payment Gateway (SSLCommerz / AamarPay)
Priority: M
Primary Gateway: SSLCommerz
API Version: SSLCommerz REST API v3
PCI-DSS Compliance Approach: Hosted payment page (redirect model) - DO NOT store card details
Endpoints:
•	Session Init: POST /gwprocess/v3/api.php
•	Validation: POST /validator/api/validationserverAPI.php  
•	Transaction Query: POST /validator/api/merchantTransIDvalidationAPI.php
•	Refund: POST /validator/api/merchantTransIDvalidationAPI.php

Authentication: Store ID + Store Password
Tokenization: Use SSLCommerz token vault for recurring payments
Webhook: IPN (Instant Payment Notification) to our callback URL
Webhook Security: Verify hash signature (MD5 of transaction details + store password)
3D Secure: Mandatory for Visa/Mastercard transactions in Bangladesh
Error Handling:
•	FAILED: Payment declined by bank → Show retry option
•	CANCELLED: User cancelled → Return to checkout
•	UNATTEMPTED: Session expired → Re-initiate
Reconciliation: Daily transaction report via merchant panel
Refund SLA: 5-7 business days
MDR: 2.5-3% for local cards, 4-5% for international cards
Test Environment: SSLCommerz Sandbox with test card numbers

5.1.4 Bank Transfer (BEFTN/RTGS)
Priority: D (Phase 1.5)
Manual Process: Customer uploads payment screenshot + transaction ID
Verification: Admin manually verifies against bank statement
Auto-reconciliation: Future enhancement via bank API integration (if available)

5.2 Insurer API Integrations

5.2.1 Integration Scenarios (Bangladesh Reality)
Scenario A - Full API Integration (10-15% of insurers):
•	Premium Quote API: POST /api/quote
•	Policy Issuance API: POST /api/policy/issue
•	Policy Query: GET /api/policy/{policyNumber}
•	Claim Submission: POST /api/claim/submit
•	Claim Status: GET /api/claim/{claimNumber}/status

Scenario B - Email-Based Integration (60-70% of insurers):
•	Platform generates PDF application
•	Email to insurer at designated email address
•	Insurer processes manually and responds via email
•	Platform admin updates system manually
•	Status tracking via email parsing or manual entry

Scenario C - Portal/Manual Entry (15-25% of insurers):
•	Platform admin logs into insurer portal
•	Manually enters application details
•	Downloads policy certificate from portal
•	Uploads to our platform

API Specifications (for Scenario A):
Authentication: API Key or OAuth 2.0 client credentials
Request Format: JSON (RESTful) or XML (for legacy systems)
Response Format: JSON
Timeout: 10 seconds for quotes, 30 seconds for policy issuance
Retry Logic: 2 retries for network errors, no retry for business validation errors
Error Codes (Standard):
•	INS-001: Invalid product code
•	INS-002: Sum insured out of range  
•	INS-003: Age not eligible
•	INS-004: Medical underwriting required (manual process)
•	INS-005: Duplicate policy detected
Idempotency: Use platform-generated quoteId/applicationId
Fallback: If API fails, queue for manual processing by operations team

5.2.2 LabAid Insurer Specific Integration
Priority: M (strategic partner)
Integration Type: Full API (Scenario A)
Endpoints:
•	Product Catalog: GET /api/products (cached for 1 hour)
•	Premium Calculator: POST /api/calculate-premium
•	Health Questionnaire: GET /api/products/{productId}/questionnaire
•	Policy Issuance: POST /api/policies
•	Claim Pre-authorization: POST /api/claims/preauth (for cashless)
•	Claim Reimbursement: POST /api/claims/reimbursement

Authentication: OAuth 2.0 with client credentials grant
Token Refresh: Automatic refresh 5 minutes before expiry
EHR Integration: Separate endpoint for cashless health claims (see 5.4)

5.3 SMS Gateway Integration

5.3.1 Primary SMS Provider  
Provider: TBD (Options: Grameenphone SMS API, Banglalink, Multi-aggregator like Muthofun, SSL Wireless)
Priority: M
Use Cases:
•	OTP delivery (most critical)
•	Payment confirmation
•	Policy issuance notification
•	Claim status updates
•	Renewal reminders

API Specifications:
Endpoint: POST /api/v1/send-sms
Authentication: API Key in header
Request Format:
{
  "to": "+8801XXXXXXXXX",
  "message": "Your OTP is 123456. Valid for 5 minutes.",
  "sender_id": "LABAIDINS",
  "type": "text" // or "unicode" for Bengali
}

Response:
{
  "message_id": "MSG-123456",
  "status": "queued",
  "to": "+8801XXXXXXXXX"
}

Delivery Report Webhook:
POST /webhook/sms-status
{
  "message_id": "MSG-123456",
  "status": "delivered", // or "failed", "pending"
  "delivered_at": "2025-12-06T10:30:00Z",
  "error_code": null
}

Requirements:
•	Sender ID Masking: Register "LABAIDINS" or "LabaidIns" (max 11 characters)
•	Unicode Support: Bengali text requires unicode=true parameter
•	Rate Limiting: Handle rate limits gracefully (typical: 100 SMS/second)
•	Delivery Rate SLA: >98% delivery rate, <60 seconds delivery time for OTP
•	Retry Logic: If SMS fails, retry once after 30 seconds, then try voice OTP
•	Operator Routing: Smart routing based on phone number prefix (GP, Robi, Banglalink, Teletalk)
•	Redundancy: Fallback to secondary SMS gateway if primary fails

Error Codes:
•	1001: Invalid phone number format
•	1002: Insufficient balance
•	1003: Sender ID not registered
•	1004: Blacklisted number (DND - Do Not Disturb)
•	2001: Network error → Retry

5.3.2 Voice OTP (Fallback)
Provider: TBD (if supported by SMS provider)
Use Case: When SMS delivery fails (network issues, phone off)
Implementation: Text-to-speech engine reads OTP in Bengali/English

5.4 eKYC Integration (Election Commission NIDW)

5.4.1 Direct Integration with Election Commission (Preferred for Phase 1)
Priority: M
Integration Type: Direct API or Third-party middleware
Prerequisites:
•	Formal MOU with Election Commission
•	IDRA approval letter for accessing NIDW database
•	Bangladesh Bank/BFIU approval (for financial services eKYC)

API Specifications (indicative - actual specs from EC):
Endpoint: POST /api/nid/verify
Authentication: Client certificate + API key
Request:
{
  "nid_number": "1234567890", // 10 or 13 digit
  "date_of_birth": "1990-01-15",
  "name": "মোহাম্মদ আহমেদ", // Bengali or English
  "photo": "base64_encoded_selfie"
}

Response:
{
  "status": "verified", // or "not_verified", "mismatch"
  "match_score": 0.95, // 0-1 for photo matching
  "details": {
    "nid": "1234567890",
    "name_bn": "মোহাম্মদ আহমেদ",
    "name_en": "Mohammad Ahmed",
    "dob": "1990-01-15",
    "address": "...",
    "father_name": "...",
    "mother_name": "..."
  }
}

Photo Matching: 
•	Face liveness check required (prevent photo-of-photo fraud)
•	Match threshold: >85% for auto-approval, <85% for manual review
Performance: <5 seconds response time (target)
Fallback: If API unavailable, queue for manual verification

5.4.2 Third-Party eKYC Provider (Alternative for Phase 1)
Providers: Pixl, Faceki, LEADS VerifID, Shufti Pro, AccuraScan
Priority: D (if EC direct integration delayed)
Benefits: Faster implementation, EC integration already built
Considerations:
•	Verify provider has valid IDRA/BFIU approval
•	Per-transaction cost (vs EC direct might be cheaper at scale)
•	Data privacy: Ensure data doesn't leave Bangladesh without consent
Selection Criteria:
•	EC NIDW integration confirmed
•	Liveness detection supported
•	OCR accuracy >95% for smart NID
•	Pricing <BDT 10 per verification
•	SLA: 99% uptime

5.5 LabAid EHR Integration (for Cashless Health Claims)

Priority: D (Phase 1.5)
Use Case: Patient admitted to LabAid hospital → Cashless claim processing
Integration Standards: HL7 FHIR R4 (preferred) or LabAid custom API

Endpoints:
•	Patient Lookup: GET /fhir/Patient?identifier={nid}
•	Pre-authorization Request: POST /fhir/Claim (with status=draft)
•	Coverage Eligibility: GET /fhir/CoverageEligibilityRequest
•	Claim Submission: POST /fhir/Claim (with status=active)
•	Admission/Discharge Notification: Webhook from EHR to our platform

Pre-authorization Workflow:
1.	Patient admitted to LabAid hospital
2.	Hospital staff checks insurance coverage via our portal/API
3.	Hospital submits estimated cost for pre-authorization
4.	Our platform checks policy coverage limits
5.	Auto-approve if within limits and policy active, else manual review
6.	Hospital receives approval with max coverage amount
7.	Patient discharged, hospital submits actual bills
8.	Platform processes claim (auto-pay if <BDT 10K and rules satisfied)

Authentication: OAuth 2.0 or HL7 SMART on FHIR
Data Security: TLS 1.3, PHI encrypted at rest
Compliance: Ensure LabAid EHR is compliant with Bangladesh health data regulations

5.6 Notification Service (Internal/External)

5.6.1 Email Service
Provider: AWS SES, SendGrid, or local provider (mailgun, etc.)
Priority: M
Use Cases: Policy certificate delivery, claim updates, monthly statements
Requirements:
•	DKIM/SPF configured for deliverability
•	Bounce/complaint handling
•	Unsubscribe management
•	Template engine for Bengali/English
•	Attachment support (policy PDFs <5MB)
Rate Limit: 50 emails/second (AWS SES limit)

5.6.2 Push Notifications (Mobile App)
Provider: Firebase Cloud Messaging (FCM) for Android, APNS for iOS
Priority: M
Use Cases: OTP auto-fill, payment status, in-app alerts
Implementation: Event-driven via Kafka → Notification service → FCM/APNS
Delivery Rate: Best effort (not guaranteed, SMS fallback for critical)

5.6.3 In-App Notifications
Priority: M
Implementation: WebSocket connection for real-time updates
Storage: Notification history in database (30-day retention)
Read/Unread tracking: User-specific flags

5.7 Analytics & Reporting Integration

5.7.1 Business Intelligence
Tool: TBD (Options: Metabase, Tableau, Power BI, custom dashboards)
Data Source: Read replica of production database
Update Frequency: Near real-time (5-minute lag acceptable)
Dashboards:
•	Executive: Daily sales, policy count, claims ratio, revenue
•	Operations: System health, support tickets, processing times  
•	Compliance: AML flags, IDRA report status, audit logs

5.7.2 IDRA Reporting
Priority: M  
Format: Excel templates (IDRA prescribed formats)
Frequency: Monthly, quarterly, annually depending on report type
Submission: Manual upload to IDRA portal (no direct API as of Dec 2025)
Reports Required:
•	Monthly Premium Collection Statement
•	Quarterly Claims Settlement Report
•	Annual Financial Condition Report (FCR)
•	CARAMELS Framework Returns (phasing in 2025-2026)
•	Event-based: Significant incidents, complaints, fraud cases

Automation: Platform generates reports in required format, compliance officer reviews and submits

5.7.3 BFIU/AML Reporting
Priority: M
SAR (Suspicious Activity Report) Filing:
•	Detection: Platform flags suspicious transactions based on rules
•	Review: Compliance officer reviews flagged transactions  
•	Submission: Manual filing to BFIU portal (goAML platform)
•	Timeline: Within 7 days of detection per BFIU guidelines
CTR (Cash Transaction Report): If applicable for large cash premiums
Frequency: As required per BFIU regulations



6. Data Model & Storage Requirements

(Note: Due to document length constraints, this section provides a summary. Full entity definitions with 20+ tables including User, Customer, KYCDocument, Nominee, Partner, Product, Policy, Premium, Payment, Transaction, Claim, ClaimDocument, Commission, Notification, Audit Log, PromoCode are documented in the comprehensive data model specification.)

6.1 Core Entities Summary

User Management: User, Customer, KYCDocument, Nominee
Product & Policy: Partner, Product, Policy, Premium  
Payments & Transactions: Payment, Transaction (IRD - Insurance Record DB)
Claims: Claim, ClaimDocument, ClaimStatusHistory
Supporting: Commission, Notification, AuditLog, PromoCode, AMLFlag, SARReport

6.2 Key Data Model Requirements

Database: PostgreSQL 14+ (ACID compliance, JSONB support, Bengali full-text search)
Optional: TigerBeetle for Transaction ledger (immutable financial records)
Document Storage: AWS S3 or MinIO
Encryption: AES-256 for PII fields (NID, documents, certificates)
Retention: 20 years for regulatory compliance, tiered storage (hot/warm/cold)
Backup: Daily full, 6-hour incremental, continuous transaction logs
DR: RTO <15 min, RPO <1 hour


7. Security & Compliance Requirements

7.1 IDRA Compliance

Monthly Reports: Premium Collection (Form IC-1), Claims Intimation (Form IC-2)
Quarterly Reports: Claims Settlement (IC-3), Financial Performance (IC-4)  
Annual Reports: FCR (Financial Condition Report), CARAMELS Framework Returns
Event-Based: Significant incidents (48 hrs), fraud cases (7 days)
Platform: Report generator with IDRA Excel templates, audit trail, 20-year archive

7.2 BFIU/AML Compliance

Tiered KYC:
- Tier 1 (Simplified): Sum assured <BDT 300K, NID + photo matching
- Tier 2 (Regular): Sum assured ≥BDT 300K, biometric + NIDW check
- EDD (Enhanced): PEPs, high net worth >BDT 50 lakh

Transaction Monitoring: 20+ automated rules for AML detection
- Rapid purchases (>3 policies in 7 days)
- High-value premiums (>BDT 5 lakh)
- Frequent cancellations
- Mismatched nominees
- Geographic/payment anomalies

SAR Filing: 7-day timeline to BFIU goAML portal, compliance officer workflow

7.3 PCI-DSS Compliance

Approach: Hosted payment page (redirect model) - DO NOT store card data
Level: SAQ-A (simplest, for redirecting merchants)
Requirements: Annual SAQ, quarterly ASV scans, TLS 1.3
Tokenization: Store only gateway tokens for recurring payments

7.4 Data Protection

Encryption at Rest: NID, documents, certificates (AES-256)
Encryption in Transit: TLS 1.3, certificate pinning for mobile
Key Management: AWS KMS/Azure Key Vault/HashiCorp, 90-day rotation
Data Masking: NID (last 3 digits), phone (mask middle), email (mask username)

7.5 Security Testing

SAST: SonarQube/Checkmarx (every commit, block critical vulnerabilities)
DAST: OWASP ZAP/Burp Suite (weekly on staging)
Penetration Testing: Pre-launch + annually (SISA InfoSec or international firm)
Vulnerability Scanning: Monthly (Nessus/Qualys/AWS Inspector)


8. Performance & Scalability Requirements

8.1 TPS (Transactions Per Second) Analysis

Baseline Load: 0.8 TPS (realistic for Phase 1 Bangladesh market)
Peak Load: 8 TPS (10× multiplier for flash sales/campaigns)
Target Capacity: 500 TPS (future-proofing, 50× current need)

8.2 Concurrent Sessions

Phase 1: 5K-10K concurrent users (realistic)
Phase 2: 100K concurrent sessions target
Memory: 10 KB per session, Redis caching

8.3 File Upload Optimization

Client-side compression: 5MB → 1-2MB (JPEG 80% quality, 1920×1080 max resolution)
Chunked upload: 1MB chunks with resume capability (tus.io protocol)
Presigned S3 URLs: Direct upload, 30-minute expiry
Virus scanning: ClamAV on uploaded files

8.4 Database Performance

Connection pooling: PgBouncer, max 100 connections per service
Read replicas: For reporting queries
Query optimization: <100ms for 95% queries
Indexing: Comprehensive indexes on foreign keys, status fields, dates

8.5 Caching Strategy

Product catalog: 5-minute TTL
User sessions: Redis with 15-minute expiry
API responses: Conditional based on data freshness requirements


9. Operational Requirements & Support

9.1 Monitoring

Application Metrics: Request rate, response time (median, p95, p99), error rate, active users
Infrastructure Metrics: CPU, memory, disk, network, load balancer health, database replication lag
Business Metrics: Policy issuance rate (>98%), payment success rate (>95%), OTP delivery (>98%), claim TAT
Tools: Datadog/Prometheus+Grafana/ELK Stack

9.2 Logging

Levels: ERROR, WARN, INFO, DEBUG (DEBUG only in staging)
Format: Structured JSON for parsing
Aggregation: ELK / Datadog Logs
Retention: 30 days hot, 90 days warm
PII Masking: Mandatory in all logs

9.3 Incident Runbooks

P1 (System Down): 15-min response, 1-hour resolution target, communicate every 30 min
P2 (Degraded Service): 1-hour response, 4-hour resolution target
P3 (Minor Issues): 4-hour response, 48-hour resolution target
Procedures: Detection, communication, investigation, fix, post-incident report

9.4 Support Channels

Customers: In-app chat (9 AM-9 PM), phone 09610-111-222, email (24-hr response), FAQ
Partners: Dedicated account manager, portal tickets, technical Slack
IDRA: Compliance officer direct contact, quarterly reports, on-demand portal

9.5 Deployment

Strategy: Blue-green deployment with one-command rollback
Windows: Tuesday/Wednesday 2-4 AM BST (avoid Thursday night, weekends, holidays)
Process: Code review (2 approvals) → Staging  → QA (30 min) → Production → Monitor (30 min)
Emergency Hotfixes: Anytime for P1, abbreviated testing

9.6 Change Management

RFC Template: Change description, business justification, affected systems, impact analysis (timeline/cost/risk), testing plan, rollback plan, approval workflow


10. Acceptance Criteria & Test Summary

10.1 Functional Acceptance Criteria (Sample - Full mapping in test plan)

Registration (FR-001 to FR-004):
✓ Valid Bangladesh phone (+8801XXXXXXXXX) accepted
✓ OTP delivered <60 sec, valid 5 min
✓ Resend OTP after 30-sec cooldown
✓ 3 failed OTP attempts → 30-min lock
✓ Duplicate NID/phone blocked with helpful message
✓ Profile validates: DOB not future, name max 200 chars, Bengali/English accepted

KYC (FR-005 to FR-010):
✓ Upload NID front/back, passport, photo (JPG/PNG/PDF, max 5MB)
✓ OCR extracts NID number (>90% accuracy), DOB, name (>85%)
✓ User can manually correct OCR errors
✓ KYC status updates: submitted → verified/rejected
✓ Notification on approval/rejection

Product Discovery (FR-055 to FR-065):
✓ Catalog loads <2 sec
✓ Filter by category, sort by premium
✓ Side-by-side comparison (up to 3 products)
✓ Product detail shows full policy wording (Bengali + English)

Purchase (FR-066 to FR-089):
✓ Multi-step form saves progress (resume capability)
✓ Premium calculation real-time, breakdown shows base + 15% VAT + total
✓ Payment methods: bKash, Nagad, Card, Bank
✓ Policy certificate generated <30 sec after payment success
✓ Certificate has QR code, emailed as PDF attachment
✓ Coupon code applies discount correctly

Claims (FR-090 to FR-106):
✓ Claim only on active policies (expired/lapsed blocked)
✓ Upload supporting documents (1-10 files)
✓ Claim number generated on submission
✓ Real-time status tracking, SMS on status change
✓ Admin can request additional documents
✓ Small claims (<BDT 10K) auto-approved if rules pass
✓ Claim approval triggers payment within 7 days

10.2 NFR Acceptance Criteria

Performance:
✓ Login API <200ms median
✓ Product catalog <300ms (cached)
✓ Policy issuance <30 sec
✓ File upload supports resume for >1MB files
✓ App starts <3 sec on mid-range Android

Security:
✓ All API calls use HTTPS (TLS 1.3)
✓ JWT tokens expire after 15 min
✓ NID masked in logs (last 3 digits only)
✓ Admin requires MFA to login
✓ Failed login locked after 5 tries

Scalability:
✓ System handles 50 TPS without degradation
✓ Database queries <100ms for 95%
✓ Auto-scaling triggers at 70% CPU

Compliance:
✓ All transactions logged to immutable audit trail
✓ IDRA reports generated on-demand
✓ AML flags reviewed within 1 hour of detection
✓ KYC documents encrypted at rest

10.3 Testing Types

Unit Testing: 80% code coverage, Jest/pytest/JUnit, mock external dependencies
Integration Testing: Service-to-service interactions (Auth→User, Policy→Payment, Claim→Notification)
E2E Testing: Critical journeys automated (Selenium/Cypress/Playwright)
Performance Testing: Normal/Peak/Stress/Spike scenarios with JMeter/K6
Security Testing: SAST (every commit), DAST (weekly), Pen-test (annual)
UAT: 20-30 LabAid employees, 2-week period, structured feedback, 80% task completion w/o help


11. Traceability Matrix & Change Control

11.1 Requirements Traceability (Sample Mapping)

| Business Objective | BRD Ref | FR/NFR | Test Cases | Status |
|--------------------|---------|---------|------------|--------|
| Digital onboarding 90% completion | BRD-1.2 | FR-001 to FR-010 | TC-001 to TC-020 | In Progress |
| 40K policies Year 1 | BRD-2.1 | FR-055 to FR-089 | TC-050 to TC-080 | Not Started |
| Claims TAT <7 days | BRD-3.3 | FR-090 to FR-106, NFR-003 | TC-100 to TC-120 | Not Started |
| 99.5% uptime | BRD-4.1 | NFR-001, NFR-007 | TC-200 to TC-210 | Planned |
| IDRA compliance | BRD-5.1 | NFR-018, NFR-030 | TC-300 to TC-320 | Not Started |
| 3+ partner integrations Year 1 | BRD-6.2 | Section 5 APIs | TC-400 to TC-420 | Not Started |

11.2 Change Control Process

Small Changes (<2 days): Product Owner approval, JIRA documentation
Medium Changes (2-10 days): Product Owner + CTO approval, weekly planning presentation
Large Changes (>10 days): Full RFC, Steering Committee approval, update project plan
Urgent Changes (P1 hotfixes): Verbal PO + CTO approval, retrospective RFC within 48 hrs


12. Appendices

12.1 Detailed Use Cases

UC-01: User Registration
Main Flow: Phone entry → OTP send → OTP validation → Profile creation → Success
Alternatives: Invalid phone, OTP not received, wrong OTP (3 attempts lock), duplicate phone
Exceptions: Network error (save state, retry), server error (log with request ID)
NFRs: OTP delivery <60 sec (95th percentile), >98% success rate, Bengali/English UI

UC-02: Policy Purchase  
[Detailed flow with 15+ steps, alternatives for payment failures, network issues, document upload problems]

UC-03: Claim Submission
[Detailed flow covering cashless vs reimbursement, document requirements, status tracking]

12.2 Bangladesh Market Context

Insurance Penetration: 0.46% of GDP (vs 4.2% in India) - Most users never bought insurance
User Literacy: ~76% literacy rate, many rural users have low financial literacy
Network Infrastructure: 4G in cities, 3G/2G rural, frequent network switching
Device Landscape: 85% Android, mostly 2-4GB RAM devices (Xiaomi, Samsung A-series, Realme)
MFS Dominance: bKash 39.9%, Nagad 18.1%, Rocket 11.7% - digital payment preference
IDRA Evolution: Regul ations modernizing (CARAMELS framework 2025-2026, digital reporting push)

12.3 Glossary (Enhanced)

BFIU: Bangladesh Financial Intelligence Unit (AML/CFT regulator)
CARAMELS: Capital Adequacy, Reinsurance, Management, Earnings, Liquidity, Sensitivity (IDRA framework)
CTR: Cash Transaction Report (for large cash transactions >BDT 5 lakh)
EDD: Enhanced Due Diligence (for PEPs, high-value policies)
FCR: Financial Condition Report (annual IDRA submission)
MFS: Mobile Financial Service (bKash, Nagad, Rocket)
NIDW: National ID Wing of Election Commission (eKYC database)
PCI-DSS: Payment Card Industry Data Security Standard
SAR: Suspicious Activity Report (filed to BFIU within 7 days)
ZHCT: Zero Human Touch Claims (automated claim processing)

12.4 Sign-off

By signing below stakeholders accept this SRS as the authoritative technical specification for Phase 1 development of the Labaid InsureTech platform.

• Product Owner / LabAid InsureTech: ____________________ Date: ______

• CTO / Technical Lead: ____________________ Date: ______

• Compliance / Legal Officer: ____________________ Date: ______

• Insurer Partners (2-3 representatives): ____________________ Date: ______

• Payment Gateway Representatives (bKash, Nagad, SSLCommerz): ____________________ Date: ______


Document Version: 2.0 (DRAFT)
Date: December 6, 2025
Status: READY FOR STAKEHOLDER REVIEW
Next Review: After stakeholder feedback incorporation
