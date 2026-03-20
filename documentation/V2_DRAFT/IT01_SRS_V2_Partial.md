Labaid Insure Tech

![](blob:vscode-webview://1d34jupcfnhdq80egt1ktcbkhjnpbjp23f4tj9sqcvtasr5qtiej/c406b0e3-c12f-4ceb-94e1-f7bcc68af13e)

| Version:2.0Status:DraftControl Level:A |  |
| -------------------------------------- | - |

  Software Requirements Specification (SRS)

For

Labaid Insure Tech Platform

Revision History

| SL | Date        | Revised By   |
| -- | ----------- | ------------ |
| 01 |             |              |
| 02 | 12_06_02025 | Faruk Hannan |

Table of Contents

Introduction

MarketContext

Overall Description

System Features & Functional Requirements

External Interface Requirements

Non-Functional Requirements

Data Model & Storage Requirements

Security & Compliance Requirements

Performance & Scalability Requirements

Operational Requirements & Support

Acceptance Criteria & Test Summary

Traceability Matrix & Change Control

Appendices (Use Cases, Screens, File References)

1. Introduction

1.1 Purpose

This SRS documents detailed system-level requirements for the LabaidInsureTech platform (mobile apps, partner portals, admin portal, and backend services). Its purpose is to provide an unambiguous specification for design, development, integration, testing, and deployment and operations teams.

1.2 Scope

The system will enable digital onboarding, product discovery, policy purchase, digital KYC, payment processing, claims submission and tracking, partner integrations, and admin/insurer workflows. This SRS covers Phase 1 (core digital capabilities) while flagging Phase 2/3 enhancements (AI underwriting, Voiceaided guidance, IOT based connectivity) for future releases. See Business Plan for strategic context.

1.3 Market Context

 With a low penetration rate of insurance compared to other countries in the region and considering mass literacy level and lack of awareness about insurance, a very simplified onboarding flow with step-by-step explanations and visual queue are required for users UI/UX space. Also considering rural internet network constraints, overall speed (Minimum 3G), a minimum spec of edge devices and also lack of availability of local cloud storage support within the country for reducing latency,specific sets of requirements are set for overall app/technology performance and reactivity.

1.4**Definitions, Acronyms, and Abbreviations**

IDRA—InsuranceDevelopment&RegulatoryAuthority(Bangladesh)

KYC- Know your customer.

KYB-Know your Business

**FCR** - Financial Condition Report (IDRA requirement)

CARAMELS - Capital Adequacy, Reinsurance Arrangements, Management, Earnings, Liquidity & Asset quality, Sensitivity

**SAR** - Suspicious Activity Report (AML requirement)

**BFIU** - Bangladesh Financial Intelligence Unit

MFS — Mobile Financial Service (bKash, Nagad)

API — Application Programming Interface

EHR — Electronic Health Record

ZHCT — Zero Human Touch Claims

UBI — Usage Based Insurance

IAM-Identity and Access Management

ACL: Access Control List

RBAC-Role based Access Control

ABAC-Attribute based Access control

UAT: User Acceptance Testing

MOU-

2. Overall Description

2.1 Product Perspective

Theplatformiscloud-native, secure-first, micro services-based, event-drivenandmobile-friendly.

Components:

Mobileapps(iOS/Android)—Customerexperience(screens&flowsperuploadeddesigns).

Partner/AgentPortal—onboarding&embeddedflows.

AdminPortal—product,pricing,userandclaimsmanagement.

Back-endServices—authentication, authorization,partner/agent management, policyengine, profiling service, contract management, paymentgatewayadapters,risk management, fraud detection, renewal service, transactional debit-credit service, storage service, LLM multi agent network , ai assistant service, MCP servers, iot-broker service, orchestration service, notification service, ticketing service(customer support) ,analytics, reporting .

Integrations—InsurerAPIs,LabAidsystems,MFS/telcoAPIs.

Flowreference:userflow diagram.

2.2User Classes & Characteristics

    Broader classification: Customers and Stakeholders

Customers:Mobile-first, diverse digital literacy

Type1 :Urban and Digitally Literate  ->Self-service mobile app users

Type2 :Semi-Urban/Agent-assisted->  payment with agent

Type3: Rural/Agricultural -> voice assisted workflow

Stakeholders:

Partners:  System on boarded, collaborative  organization (MFS, hospitals, e-commerce)

Insurer Agents:Staff at partner organizations  underACL

Insurer Underwriters: Internal to partner’s organization insurers (API consumers).

Independent Agent: on boarded, certified payment agents in remote locations(dashboard)

Vendor :independent service provider to insure Tech as a client (dashboard)

System Admin :Root level to cloud provider gateway ( IAM,ACL,Policy)

Repository Admin: Root  Level to repo and droplet( Pull, Merge ,Deployaccess )

Database Admin: Root Level to Databases /Storage (  Full management access)

Business Admin:Root level to portal(policy management, product updates).

Focal Person: Root level  toPartner management ( partner verification,onboarding ,dispute )

Partner Admin :Partner Side Root Tenant( IAM, ACL control)

Dev :user level ( software development)

Support/Call Centre:Assisted onboarding & claims help(dashboard)

OperatingEnvironment

Cloud:AWS/Azurepreferred;region:Bangladesh(orcompliantregion like Singapore not more than 80ms added latency).

Mobile:iOS13+/Android 9+. Minimum support with 4GB RAM

Network :Minimum 3G (EDGE will not be supported) with frequent network switching

App Download: Must have lowest download size as low as 10 Megabyte , gradual download and caching of data run up to maximum 100 Megabyte

Offline Capabilities: Must have minimum  persistent view of  critical data

Low Bandwidth Mode: With minimum feature with device level cached static data .

Low Resource Mode: Must run with minimum memory and low power mode.

Browsersforportals:Chrome,Brave, Firefox,Edge(latesttwoversions).

USSD/SMS fallbacks: will have retry and acknowledgement support

3System Features & Functional Requirements

Each requirement below has a Function Group unique id (e.g., FG-1 )and Functional requirements  unique ID (e.g., FR-1).

Functionality Deploy Phases: Phases 1,Phase 1.5, Phase 2 ,Phase 3

Priority levels: M = Mandatory (Phase 1), D = Desirable (Phase 1.5), F = Future (Phase 2/3).

3.1 Functional Requirements

| ID     | Requirement Description                                                                                                                                                                               | Priority |
| ------ | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| FG-001 | 3.11 -IAM Group( Authentication)                                                                                                                                                                      | M        |
| FR-001 | The system shall allowCustomerstoregisterand loginviamobile number withOTPverificationwith minimumdata(name, date of birth, mobile no, email(optional)                                                | M        |
| FR-002 | The system shall auto update OTPfromSMSwithone timereadpermissionfromdevice                                                                                                                           | D        |
| FR-003 | The system shall have SMSOTP  failurefallbackwith persistentformstate,retryloopand resend button                                                                                                      | M        |
| FR-004 | The system shall have image capture stepfor user verification step1                                                                                                                                   | M        |
| FR-005 | The system shall perform loginstep  viae-KYC  asper Office of the Controller of Certifying Authorities Bangladesh                                                                                     | D        |
| FR-006 | The system shall automate with voice-assisted workflow on landing page for Type 3 customer to guide to login.                                                                                         | F        |
| FR-007 | The system shall have option for logged in customer to link with social accountvia  oAUth2 (for Type1 ,Type2 customer (Google,Facebook)  tofinalize profile                                           | D        |
| FR-008 | The system shall have mandatory NID upload for user verification step 2                                                                                                                               | M        |
| FR-009 | The system shalldigitizedandProcess NID data                                                                                                                                                          | D        |
| FR-010 | The System shall show green check mark for user verification step                                                                                                                                     | D        |
| FR-011 | The systemshall have option ofInsuretech employedFocalpersonto onboard a pre-approved,MOU-signed partner on behalf of thepartner, setup basic ACL,andprovide onePartner Adminwith temporary password. | M        |
| FR-012 | The system shallhavestakeholdersregistration via OPEN ID Identity provider                                                                                                                            | M        |
| FR-013 | The system shall havestakeholdersregistrationvia SAML Identity provider                                                                                                                               | D        |
| FR-014 | The system shall havestakeholdersregistrationvia SSO provider                                                                                                                                         | D        |
| FR-015 | The system shall havestakeholdersregistrationvia email with email verification                                                                                                                        | M        |
| FR-016 | The system shalllockverifiedstakeholdersKYBfrom update                                                                                                                                                | M        |
| FR-017 | The system will allow all logged in stake holder to changepassword                                                                                                                                    | M        |
| FR-018 | The system shall have option of account recovery question                                                                                                                                             | D        |
| FR-019 | The system shall have 2FA authentication via mobile number SMS/OTP                                                                                                                                    | M        |
| FR-020 | The system shall have account recovery option via email /mobileand 2^nd^step with security question                                                                                                   | M        |
| FR-021 | The system willallow verified stakeholdersassisted  byFocal Person formigration  toa new account incase verification data needs to be updated                                                         | D        |
| FR-022 | The system shall maintain a dashboardfor stakeholdersto view all ongoingverificationworkflows categorized byaccount type                                                                              | M        |
| FR-023 | The system shall allow users to updatenonKYC/KYBuserdatae.gemployee data                                                                                                                              | M        |
| FR-024 | The system shalldefine mandatoryKYCfieldsforCustomeridentifyverification                                                                                                                              | M        |
| FR-025 | The System shall define mandatory KYBfields  andKYCfields( ifAdmin )forstakeholdersverification                                                                                                       | M        |
| FR-026 | The system will block duplicate account forKYCconflict                                                                                                                                                | M        |
| FR-027 | The system shall have automatic merge try forKYCconflict                                                                                                                                              | D        |
| FR-028 | The system shall have manual merge step forKYCconflict.                                                                                                                                               | M        |
| FR-029 | The system will only allow manual merge step forKYBconflict                                                                                                                                           | M        |
| FR-030 | The system will block user for multiple wrongattemptofloginfor a short period                                                                                                                         | M        |
| FR-031 | The system will ban for identities from login for lifetime for malicious attempts                                                                                                                     | M        |
| FR-032 | The system will allow APIKeyswithlimited lifetimefor 3^rd^party system                                                                                                                                | D        |

| ID     | Requirement Description                                                                                                                  | Priority |
| ------ | ---------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| FG-002 | 3.12Verification and Authorization                                                                                                       | M        |
| FR-033 | The system shall verify Customer with minimum data-captured image, NID, mobilenumber                                                     | M        |
| FR-034 | The system shall ask for additionaldata(e.gmedicalhistory)asper  theinsurance product provider policyrequirementswithdigital copy upload | M        |
| FR-035 | The system shall allow medical dataverificationreport with authorized government certified medical officer??                             | D        |
| FR-036 | The system shall have imagerecognition system foruserverification                                                                        | D        |
| FR-037 | The system shallperform OCR from uploaded document to makeautomatedEHR/Medical historyfor Customer                                       | F        |
| FR-038 | The system will take Customer geo locationdata ,IPtracking,IOTdatato avoid ID theft                                                      | F        |
| FR-039 | The system shall haveminimum KYB for partner as perInsuretech internal business policy(Likewhether toadd banksolvency etc.)              | M        |
| FR-040 | The system willcategorizedstakeholders as perfinancial and regulatorycompliance                                                          | M        |
| FR-041 | The system shallprovide partners sub domain                                                                                              | F        |
| FR-042 | The system shall provide partner own portal with dashboard                                                                               | M        |
| FR-043 | The system shallprovidepartnerto make internal ACL with eitherRBACorABACwithresource ,policy ,permission set                             | M        |
| FR-044 | The system shall haveInsure tech internal ACL as percasbinRBACwith domains/tenantsRulefor .NET(https://casbin.org/docs/supported-models) | M        |
| FR-045 | The system will providePartner Adminof Full CRUD access of partner domain data (employee, policy, product, price etc.)                   | D        |
| FR-046 | The system will provide access Partner to dedicated data storage( asperInsureTechPolicy)                                                 | D        |
| FR-047 | The system will followserver sidesession management for stakeholders                                                                     | M        |
| FR-048 | Thesystem  willfollow oauth2 and token refresh for Customer                                                                              | M        |
| FR-049 | The systemwill  tryto make random token refreshand keyrotationtodisallow misuse                                                          | D        |
| FR-050 | The system will only expose secure proxy server IP (Cloud flare) to user                                                                 | D        |
| FR-051 | The system shall invalidate user session/token immediately after logout                                                                  | M        |
|        |                                                                                                                                          |          |
|        |                                                                                                                                          |          |
|        |                                                                                                                                          |          |

| ID     | Requirement Description                                                                                                                                                                          | Priority |
| ------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | -------- |
| FG-003 | 3.13-Product Catalog & Policy Discovery                                                                                                                                                          | M        |
| FR-055 | The System shall allow Partners to upload Product catalog, Product with dataMetadata ,coverage, premiums, insurer, terms and conditions                                                          | M        |
| FR-056 | The system shall allow partners to list products for approval ofinsureTech                                                                                                                       | M        |
| FR-057 | The system shall allow partners to check product approval status and if unapproved will be prompted with disqualification reason                                                                 | M        |
| FR-058 | The system shall allow partners to pick live date on site for any specific product                                                                                                               | M        |
| FR-059 | The system shall allow partners toedit a live product but will go through same re approval process                                                                                               | M        |
| FR-060 | The system will not allow partners todeleteaproduct afterit islive ,partnerswill take approval as perInsureTechmanagement for such change and will settle any claim during that product purchase | M        |
| FR-061 | The system will categorized products as per partners provided tag and meta data                                                                                                                  | M        |
| FR-062 | The system will allow users to compare same category product side by side with premium and coverage                                                                                              | M        |
| FR-063 | The system will have filtering policy to show particular products/a product as per category,premiums ,coverage,                                                                                  | M        |
| FR-064 | The system will auto pick any product based onInsureTechinternal rating as featured, recommended                                                                                                 | M        |
| FR-065 | Product Detail: Full policy wording, exclusions, and illustrative premium breakdown                                                                                                              | M        |

| ID     | Requirement Description                                                                                                                                                                                                                                                              | Priority |
| ------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | -------- |
| FG-004 | 3.14-Policy Purchase & Issuance                                                                                                                                                                                                                                                      | M        |
| FR-066 | The system will have a minimum number of fixedstepfor any product purchase asminimum-steps                                                                                                                                                                                           | M        |
| FR-067 | If addition flow is required more thanminimum-stepsit willbe added lateralways maintaining same serial                                                                                                                                                                               | M        |
| FR-068 | The system product purchaseminimum stepsareaddpersonal details(one time only),add product tobuy ,addnominee (will be prompted for previously setnominee) ,adddocuments, additional document in always same serial with documenttypes  .then purchase order will set for verification | M        |
| FR-069 | The system will try immediate verification of documents with visual progress status (e.ggreencircle)                                                                                                                                                                                 | D        |
| FR-070 | The system will prompt customer for product purchase order verification with modal notification /voice and afterthat  takeon next payment page                                                                                                                                       | M        |
| FR-071 | The system will allow customer to makepayment  asper listedbKash, Nagad, card (PCI-DSS compliant), bank transfer (apiversion details during detail engineering)                                                                                                                      | M        |
| FR-072 | The system will allow customer to allow payment via agent and add transaction id and upload payment details                                                                                                                                                                          | D        |
| FR-073 | The system shall verify/show premium calculation as per partners provided product data                                                                                                                                                                                               | M        |
| FR-074 | The system shall show premium calculations as per inbuilt pricing engine                                                                                                                                                                                                             | D        |
| FR-075 | The system shall show premium calculation of Partner API as per API data schema                                                                                                                                                                                                      | M        |
| FR-076 | The system shall support promo code for product discount                                                                                                                                                                                                                             | M        |
| FR-077 | The system will support partners to add /remove discounts                                                                                                                                                                                                                            | M        |
| FR-078 | The system will maintaintransparent ,indestructibletransactionrecord ,debit-credit with history per partner per product percustomer  inpurpose builtdatabase(Tiger beetleorequivalent)calledasInsuranceRecord DB(IRD)                                                                | D        |
| FR-079 | The system will generate a digital policycertificateas artifacts as pdfwith QR codeand store in user account.                                                                                                                                                                        | M        |
| FR-080 | The digital certificate will have customereKYCdigital signinjected                                                                                                                                                                                                                   | D        |
| FR-081 | The system will have verification support for any required authority to verify the generated certificate                                                                                                                                                                             | M        |
| FR-082 | The system will make anadditional ,irrevocabledigital tri party contract withinsureTech,partner,Customer  foreach record inIRDwill be called asInsuranceContract                                                                                                                     | D        |
| FR-083 | The system shall generateweb3.0andblockchainbasedInsuranceContract                                                                                                                                                                                                                   | F        |
| FR-084 | The systemshall  maintainlinking record with all older purchasedproducts  forany  newpurchase  ofcustomers                                                                                                                                                                           | M        |
| FR-085 | The system will allow customer to purchaseproductonbehalf  offamily member                                                                                                                                                                                                           | D        |
| FR-086 | The system will not ask for any data that customer already uploaded and verified in history                                                                                                                                                                                          | M        |
| FR-087 | The system will allow customer to progress medical history by updating with new record                                                                                                                                                                                               | M        |
| FR-088 | The system shall produce proper detail message as systemSMSfor customer if any document verification fails                                                                                                                                                                           | M        |
| FR-089 | The system shall havepurpose builtAIChatbot to assist customer during product search, selection, purchase,verification ,paymentstage                                                                                                                                                 | F        |

| ID     | Requirement Description                                                                                                                    | Priority |
| ------ | ------------------------------------------------------------------------------------------------------------------------------------------ | -------- |
| FG-005 | 3.15Claims Management                                                                                                                      | M        |
| FR-090 | The system willhave  aForm of fixed steps of Insurance Claims -Policy prefill, claim reason selection, document upload (images, bills      | M        |
| FR-091 | The system will show success screen after initial screening of documents after submission.                                                 | M        |
| FR-092 | The system will make a digital hash to identify each unique submission and keep as internal record                                         | M        |
| FR-093 | The system will automatically notify partner for new claims andStartcommon  statusdashboard with customer                                  | M        |
| FR-094 | The system will show status tracking of the claim in the dashboard to customer                                                             | M        |
| FR-095 | The system will have a chat interfacewith customer and partner agent in dashboard                                                          | D        |
| FR-096 | Thesystemwill have a web RTC based video call support with customer and partner agent in dashboard                                         | D        |
| FR-097 | The system will haveagents note support in dashboard                                                                                       | M        |
| FR-098 | The system will allow partner to add verification details forinsureTechrecord keepingand agree onApproval Matrixtimeline                   | M        |
| FR-099 | The system will allow partner to make initial approval or rejection request toinscureTech                                                  | M        |
| FR-100 | The system will allowBusiness AdminandFocal Personjoint approve or reject upon internal verification                                       | M        |
| FR-101 | The system will automate payment process as per customer selected payment channel                                                          | M        |
| FR-102 | The system shall have auto verification, claim approval andpaymentforsmall claims with pre agreement as per Partner, Insurer andInsureTech | D        |
| FR-103 | The system will have Auto fraud detection system todetect  frequentclaims and changed uploads by customer                                  | D        |
| FR-104 | The system will provide system warning to customer as perInsureTechpolicy                                                                  | M        |
| FR-105 | The system willautorevoke access to customer as perInsureTechpolicy for suspicious activity                                                | D        |
| FR-106 | The system will maintain proper balance sheet on Customer, Partner, Agent andInsureTechlevel for various selected period oftime .          | D        |

Refenrence :Table Approval matrix (Example)

| CalimedAmount | Approval Level  | Maximum TAT |
| ------------- | --------------- | ----------- |
| BDT 0-10K     | L1 Auto/Officer | 24 Hours    |
| BDT 10K-50K   | L2 Manager      | 3 days      |
| BDT 50K-2L    | L3 Head         | 7 days      |
| BDT 2L+       | Board + Insurer | 15days      |

| ID     | Requirement Description                                                                                                                                            | Priority |
| ------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------ | -------- |
| FG-006 | 3.16Policy Management & Renewals                                                                                                                                   | M        |
| FR-107 | The system will provide customer a policy dashboard after first product purchase                                                                                   | M        |
| FR-108 | The system will show active and past policies, renewal prompts, policy pdfdownload  incustomer dashboard                                                           | M        |
| FR-109 | The system will providecustomertwooption for un-payedpoliciesrenewal :autoandmanual .                                                                              | M        |
| FR-110 | The system will guide customer in clear steps for manual policy renewal                                                                                            | M        |
| FR-111 | The system will allow customer to partially adjust policy data1.Currentaddress 2. Nominee                                                                          | M        |
| FR-112 |                                                                                                                                                                    | D        |
| ID     | Requirement Description                                                                                                                                            | Priority |
| FG-007 | 3.17Notifications & Communication                                                                                                                                  | M        |
| FR-113 | The system will have KAFKA (or equivalent) eventdrivennotificationandorchestrationsystemtriggering following channelslike:InSystem PUSH(portal & mobile),SMS,Email | M        |
| FR-114 | The system will auto trigger notification for verification, purchase confirmation, claims updates, renewal reminders                                               | M        |
| FR-115 | The system will allow user to be notified for additional events as per providedlist(e;gnew product launch)                                                         | D        |
| FR-116 | The system will provide partners to create secondary marketing notificationforcustomersby filteringage,sex, groups                                                 | D        |
| FR-117 | The system will allow customerswith opt in opt out policy for additional notification                                                                              | M        |
| FR-118 | The system will support customermobilemute modewith minimum text notification (avoiding on system push for edge devices)                                           | M        |
| FR-119 | The system will allow customer support call within app                                                                                                             | D        |
| FR-120 | The system will provide common Q&A within app for customer knowledge                                                                                               | M        |
| FR-121 | The system shall provide AI chatbot for customer support and automatic ticketing for unresolved issue                                                              | F        |
| FR-122 | The system will Providecommon question and answer form for customer support callcenter(Vendor)                                                                     | D        |
| FR-123 | The system will provide ticket form fill up option with auto recorded customer support call                                                                        | D        |

| ID     | Requirement Description                                                                                                                                    | Priority |
| ------ | ---------------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| FG-008 | 3.18-Admin & Reporting                                                                                                                                     | M        |
| FR-124 | The system will maintain separate dashboard for different admin as per definedrole ,resources, policies byInsureTechmanagement                             | M        |
| FR-125 | The system will maintain strict 2FA for Any admin level access                                                                                             | M        |
| FR-126 | The system will make dynamic content view based oncontext ,workflow, and data accessavailable  fordifferent admin                                          | M        |
| FR-127 | The system willprovide  commondashboard fordomianlevel, team level                                                                                         | D        |
| FR-128 | The system will allow internalapproval ,task issuing to internal users                                                                                     | D        |
| FR-129 | The system will havesseparate  usermanagement tab                                                                                                          | M        |
| FR-130 | The system will have separate product management tab                                                                                                       | M        |
| FR-131 | The system will have separate claim management tab                                                                                                         | M        |
| FR-132 | The system will have assigned task tab per user                                                                                                            | M        |
| FR-133 | The system will have option of reporting as per access to users                                                                                            | M        |
| FR-134 | The system will have minimum reporting options asDaily sales, claims ratio, partner performance, policy counts and KPIs (aligned tobusiness plan targets). | M        |

| ID     | Requirement Description                                                                                                                                                                                                                                                                                                                                                                                                                                    | Priority |
| ------ | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| FG-009 | 3.19-PartnerPortal(Hospitals,Ecommerce)                                                                                                                                                                                                                                                                                                                                                                                                                    | M        |
| FR-135 | The system will provide special dashboard for hospital partner to be able to initiate purchase insurance for any customer on behalf of them                                                                                                                                                                                                                                                                                                                | M        |
| FR-136 | The system shall have option to transfer customer record with partner providedapiand purchase id and authentication token                                                                                                                                                                                                                                                                                                                                  | D        |
| FR-137 | The system shall provideAPI  forEcommerce partner to embedproduct in page and special dashboard to purchase product fromInsureTech                                                                                                                                                                                                                                                                                                                         | M        |
| FR-138 | The system shall provide sandbox for 3^rd^part developer for authentication and testing payment and purchase process                                                                                                                                                                                                                                                                                                                                       | F        |
| FR-139 | The system will provide any partneronboarding analytics,leads, Commission statements in their dedicated dashboard                                                                                                                                                                                                                                                                                                                                          | M        |
| FR139B | The systemwill have internal Adminspecial data product asBusiness IntelligenceTool: TBD (Options:Metabase, Tableau, Power BI, custom dashboards)Data Source: Read replica of production databaseUpdate Frequency: Near real-time (5-minute lag acceptable)Dashboards:• Executive: Daily sales, policy count, claims ratio, revenue• Operations: System health, support tickets, processing times• Compliance: AML flags, IDRA report status, audit logs | F        |
| FR-140 | The system will provide any partneronboarding analytics,leads, Commission statements via API                                                                                                                                                                                                                                                                                                                                                               | M        |

| ID     | Requirement Description                                                                                                              | Priority |
| ------ | ------------------------------------------------------------------------------------------------------------------------------------ | -------- |
| FG-010 | 3.20-Audit& Logging                                                                                                                  | M        |
| FR-141 | The systemshall maintain immutablelogs for critical actions suchas :policy issue, claimapproval ,claimrejection ,payment and dispute | M        |
| FR-142 | The system shallmaintain data retention policy up to minimum of 20yearstomaintain records for regulatory compliance                  | D        |
| FR-143 | The system will track each logged in user for auxiliary actions and will haveadditional  datalogs                                    | D        |
| FR-144 | The system will allow partner to maintain additional data logs as per customer andInsureTechMOU                                      | F        |
| FR-145 | The system will provide special portal to regulatory body to access requested data asper  policyand law of regulatory bodies         | M        |

| ID     | Requirement Description                                                                                | Priority |
| ------ | ------------------------------------------------------------------------------------------------------ | -------- |
| FG-011 | 3.21-UserInterface                                                                                     | M        |
| FR-146 | The system shall maintain similar User interfacewith different operating system                        | M        |
| FR-147 | The system shallprovide smartdata  widgetfor mobile user                                               | D        |
| FR-148 | The systemshall  providevoice assisted workflow for type 3 user                                        | F        |
| FR-149 | The system shall provide desktop first web UI for portal                                               | M        |
| FR-150 | The systemshall  takeminimum permission for all service from user device (e’gcamera,one timemsg read) | M        |

| FG-012 | 3.22-APIDesignand Data Flow                                                                                                                                                                                  | M |
| ------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | - |
| FR-151 | The system will maintain two level of API private and public                                                                                                                                                 | M |
| FR-152 | The system will have threelevel of privateAPIascategory1,category2andcategory 3with secure middle layer                                                                                                      | D |
| FR-153 | The system will use protocol buffer andGRPCbased communicationforcategory 1APIbetweengateway and microservicesinhttpwithsystem adminmiddle layer                                                             | F |
| FR-154 | The system will useGraphqlbasedAPIasCategory2  forgateway and customer device withjwttoken system and oauthv2based  securitymiddle layer                                                                     | M |
| FR-155 | The system will useRestFulAPIwithJSONwith compliance ofOPEN APIstandardasCategory 3 for 3^rd^party Integration with server side auth and will provide user publicdocumentation,wiki,sandboxand a mock server | D |
| FR-156 | The  systemwill provide publicRestfulAPIwithJSON with compliance ofOPEN APIstandardfor productsearch ,product list andview  forlive productsetc                                                              | M |
| FR-157 | The system will only expose proxy server (Cloud flare) and entry node IP (NGINX )to public space                                                                                                             | M |
| FR-158 | The system willblock all accesstomicroservices for all available API except category 1                                                                                                                       | M |
| FR-159 | The system willuseInsureTechInternal protocol for IOT based data extraction and data binding                                                                                                                 | F |
| FR-160 | The system willconsolidate ,annotate ,add context and process and savedata  withinregulatory limitation to train internal AI agents                                                                          | F |
| FR-161 | The system will generatestatistics,prediction based on big data and provide special projection data topartners  givenspecial agreement withInsuretech                                                        | F |
| FR-162 | The systemwillImplementation: WebSocket connection for real-time updatesStorage:                                                                                                                             | D |

| ID      | Requirement Description                                                                                                                                                                                                                                                                                                                                 | Priority |
| ------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| FG-013  | 3.22-DataStorage                                                                                                                                                                                                                                                                                                                                        | M        |
| FR-165  | The system will store all uploaded object withawss3 or equivalent data storagesystem.                                                                                                                                                                                                                                                                   | M        |
| FR-166  | The system will store relational Data onself hostedPostgress V17 database(ACID compliance, JSONB support, Bengali full-textsearch)withConnection pooling:PgBouncer, max 100 connections per service,Read replicas: For reporting queriesQuery optimization: <100ms for 95% queriesIndexing: Comprehensive indexes on foreign keys, status fields, dates | M        |
| FR-167  | The system will store finance transaction record, debit-credit record,insurance  recordand contract on Specialpurpousebuilt database (likeTigerbeetleor Equivalent                                                                                                                                                                                      | D        |
| FR-168  | Thesystem will use cache databaseforfaster datadelieryProduct catalog: 5-minute TTLUser sessions: Redis with 15-minute expiry                                                                                                                                                                                                                           | M        |
| FR-169  | The system will process tokenized data on vector database likePgvector( Postgress) orPineCone                                                                                                                                                                                                                                                           | D        |
| FR-0170 | The system will store app nativeencrypteddata in user device in SQLite                                                                                                                                                                                                                                                                                  | M        |
| FR0171  | The systemwllstore product catalog and metadata ,unstructureddata  ina NoSQL Database (AWS Dynamo DB /MongoDB)                                                                                                                                                                                                                                          | M        |
| FR0172  | Upload data policyClient-side compression: 5MB → 1-2MB (JPEG 80% quality, 1920×1080 max resolution)Chunked upload: 1MB chunks with resume capability (tus.io protocol)PresignedS3 URLs: Direct upload, 30-minute expiry                                                                                                                               | M        |
| FR0173  | Backup: Daily full, 6-hour incremental, continuous transaction logs                                                                                                                                                                                                                                                                                     | F        |
| FR0174  | Retention: 20 years for regulatory compliance, tiered storage (hot/warm/cold)                                                                                                                                                                                                                                                                           | F        |

4. Non-Functional Requirements

| ID       | Requirement Description                                                                                       | Priority |
| -------- | ------------------------------------------------------------------------------------------------------------- | -------- |
| NFR-001  | The system shall ensure 99.5% uptime to meet critical business needs.                                         | M        |
| NFR-002  | The system shall encrypt(TLS 1.3 and AES-256)allsensitive data in transit and at rest.                        | M        |
| NFR-003  | The system shall handle up to1,000 concurrent users without performance degradation.                          | M        |
| NFR-004  | The system shall handle up to 5,000 concurrent users without performance degradation.                         | D        |
| NFR-005  | The system shall handle up to 10,000 concurrent users without performance degradation.                        | F        |
| NFR-006  | The system shall havea medianresponse time of less than1secondsand not more than 2 sfor 95% of user actions   | M        |
| NFR-007  | Disaster Recovery: RTO < 15 minutes for critical services, RPO < 1 hour                                       | M        |
| NFR-008  | The system will have category 1 API response time < 100ms                                                     | D        |
| NFR-009  | The system will have category 3 API response time < 200ms                                                     | D        |
| NFR-010  | The system will have category 2 API response time < 2 sec                                                     | D        |
| NFR-011  | The system will have public API response time < 1sec                                                          | D        |
| NFR -012 | The system will ensure minimum app startup time<  5sec                                                        | D        |
| NFR-013  | The system will support 1 million active user and 200K active policies in 24 months withauto  verticalscaling | F        |
| NFR-014  | The system will have automatic switch option for user to different mobile device with data migration          | F        |
| NFR-015  | The system will have multi language support (English+ Bengali)                                                | D        |

5 .Security & Compliance Requirements

| 1  | The system will use separate secretvault .AWS KMS/Azure Key Vault/HashiCorp, 90-daykeyrotation                                                                                                                                                                                                                                                                                                                                                                    | M |
| -- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | - |
| 2  | The system will useData Masking: NID (last 3 digits), phone (mask middle), email (mask username)                                                                                                                                                                                                                                                                                                                                                                  | M |
| 3  | The system will followPCI-DSS compliance for card flowsApproach: Hosted payment page (redirect model) - DO NOT store card dataLevel: SAQ-A (simplest, for redirecting merchants)Requirements: Annual SAQ, quarterly ASV scans, TLS 1.3Tokenization: Store only gateway tokens for recurring payments                                                                                                                                                              | M |
| 4  | The system will have AML/CFT detectionhooks .Transaction Monitoring: 20+ automated rules for AML detection- Rapid purchases (>3 policies in 7 days)- High-value premiums (>BDT 5 lakh)- Frequent cancellations- Mismatched nominees- Geographic/payment anomalies                                                                                                                                                                                                 | D |
| 5  | The system will have IDRA reporting capabilities following IDRA dataformat .Monthly Reports: Premium Collection (Form IC-1), Claims Intimation (Form IC-2)Quarterly Reports: Claims Settlement (IC-3), Financial Performance (IC-4)Annual Reports: FCR (Financial Condition Report), CARAMELS Framework ReturnsEvent-Based: Significant incidents (48hrs), fraud cases (7 days)Platform: Report generator with IDRA Excel templates, audit trail, 20-year archive | D |
| 6  | The system will have regular penetrationtest .Penetration Testing: Pre-launch + annually (SISA InfoSec or international firm)                                                                                                                                                                                                                                                                                                                                     | D |
| 7  | The system will have regular security audits from various security auditors and regulatory bodies and maintain compliance                                                                                                                                                                                                                                                                                                                                         | D |
| 8  | DAST: OWASP ZAP/Burp Suite (weekly on staging)                                                                                                                                                                                                                                                                                                                                                                                                                    | D |
| 9  | SAST: SonarQube/Checkmarx(every commit, block critical vulnerabilities)                                                                                                                                                                                                                                                                                                                                                                                           | D |
| 10 | Virus scanning:ClamAVon uploaded files                                                                                                                                                                                                                                                                                                                                                                                                                            | M |

8. Performance & Scalability Requirements

• Peak TPS (Transactions per second) for purchase flow: baseline 50 TPS, scale to 500 TPS via

auto-scaling.

• Concurrent active sessions target: 100k with <5% degradation.

• File upload: support 10 MB images with background upload/resume.

9. Operational Requirements & Support

9.1 Monitoring

Application Metrics: Request rate, response time (median, p95, p99), error rate, active users

Infrastructure Metrics: CPU, memory, disk, network, load balancer health, database replication lag

Business Metrics: Policy issuance rate (>98%), payment success rate (>95%), OTP delivery (>98%), claim TAT

Tools: Datadog/Prometheus+Grafana/ELK Stack 10. Acceptance Criteria & Test Summary

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

✓Multi-step form saves progress (resume capability)

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

✓ All API calls use HTTPS (TLS 1.3)(except category1)

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

| Business Objective                | BRD Ref | FR/NFR            | Test Cases       | Status           |
| --------------------------------- | ------- | ----------------- | ---------------- | ---------------- |
| Digital onboarding 90% completion |         | BRD-1.2           | FR-001 to FR-010 | TC-001 to TC-020 |
| 40K policies Year 1               | BRD-2.1 | FR-055 to FR-089  | TC-050 to TC-080 | Not Started      |
| Claims TAT <7 days                | BRD-3.3 | FR-090 to FR-106, | TC-100 to TC-120 | Not started      |

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

IDRA Evolution: Regulations modernizing (CARAMELS framework 2025-2026, digital reporting push)

12.4 Sign-off

By signing below stakeholders accept this SRS as the authoritative technical specification for Phase 1 development of the LabaidInsureTech platform.

• Product Owner / LabAidInsureTech: ____________________ Date: ______

• CTO / Technical Lead: ____________________ Date: ______

• Compliance / Legal Officer: ____________________ Date: ______

• Insurer Partners (2-3 representatives): ____________________ Date: ______

• Payment Gateway Representatives (bKash, Nagad, SSLCommerz): ____________________ Date: ______

* 
* [•••]()
* 
* Go to[ ] Page
