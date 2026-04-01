# LabaidInsuretech Insurer Dashboard Full Specification

**Document Version:** 2.0  
**Date:** 2026-03-18  
**Audience:** Product, Engineering, QA, Operations, Compliance, Insurer Stakeholders  
**Status:** Draft for alignment and implementation  
**Changelog:** v2.0 — incorporated field-level digitization reference from parsed insurer documents (PPTX, DOCX, XLSX, PDF) in `documentation/docs_forms/`

---

## 1. Purpose

This document defines the full specification for the **LabaidInsuretech Insurer Dashboard**, a dedicated web portal for insurer-side operations within the broader LabaidInsuretech platform.

The dashboard is intended to give insurers a secure operational workspace to:

- review and decide on incoming proposals
- manage insurer-owned document templates and product forms
- process and monitor claims
- route survey-required claims to surveyors
- manage TPA-assisted health claims
- monitor service levels, fraud indicators, and operational KPIs
- collaborate with Labaid Insuretech teams through structured digital workflows

This specification is aligned to the current platform direction documented across `documentation/About`, `documentation/BRD`, `documentation/core_plans`, and `documentation/SRS_v3`.

---

## 2. Strategic Positioning

The insurer dashboard is not a standalone product. It is an insurer-facing portal inside the **same LabaidInsuretech platform architecture** used by the B2B portal and other operational surfaces.

It must support the MD-aligned segment model:

- **Health insurance:** TPA-oriented operating model with provider-network and claims orchestration
- **Auto insurance:** hybrid model requiring claim assessment, survey, pricing discipline, and fraud review
- **Life insurance:** traditional insurer controls with digital proposal intake and servicing
- **Property & Casualty:** operationally efficient digital processing with flat-fee and automation support

For Pragati and similar non-life carriers, the insurer dashboard is the operational control plane for:

- insurer review of proposals submitted by Labaid Insuretech
- insurer-controlled template/document publishing
- non-life claim triage and settlement operations
- surveyor coordination for motor, fire/property, pet, and other field-assessment claims

---

## 3. Objectives

### 3.1 Business Objectives

- reduce insurer turnaround time for proposal review and claim decisions
- replace email/WhatsApp/manual spreadsheet coordination with auditable workflows
- standardize document collection and insurer-specific form completion
- improve insurer visibility into proposal pipeline, claim backlog, and settlement performance
- provide product-level operational intelligence across health, motor, travel, fire, pet, and commercial lines

### 3.2 Operational Objectives

- establish a single insurer workspace across proposals, documents, claims, survey, and reporting
- support partial live integration while allowing controlled fallback/mock operation during rollout
- create structured communication interfaces for chat, calls, notes, and document requests
- maintain full auditability for insurer actions

### 3.3 Technical Objectives

- reuse the platform’s existing auth/session/gateway/BFF patterns
- align with proto-first domain models and planned service/event orchestration
- remain deployable within the current Docker + nginx production strategy
- remain responsive and desktop-first for insurer operations teams

---

## 4. Scope

### 4.1 In Scope

- insurer authentication and role-based portal access
- insurer dashboard home and KPI overview
- proposal intake, review, approval, rejection, and feedback loop
- insurer document template library
- digital forms derived from insurer Excel/PDF source sheets
- claims queue, review, settlement, and exception handling
- dedicated Surveyor Desk
- TPA and claim matrix workspace
- insurer product/category/playbook visibility
- collaboration tools for operational communication
- audit trail, compliance visibility, and reporting exports

### 4.2 Out of Scope for This Spec

- core actuarial engine logic
- customer portal experience
- partner portal sales workflows except where they surface insurer-side review needs
- insurer accounting/GL replacement
- full regulatory filing engine implementation

---

## 5. Primary Users and Roles

### 5.1 Insurer Roles

| Role | Purpose | Typical Permissions |
|------|---------|---------------------|
| Insurer Admin | Owns insurer configuration and operational oversight | full insurer workspace access, template publish, role assignment, reporting |
| Underwriting Officer | Reviews and decides proposals | view proposal queue, request clarification, approve/reject, add terms |
| Claims Officer | Handles claim review and settlement | claim review, document verification, reserve/settlement actions |
| Survey Coordinator | Assigns surveyor and monitors field review | assign surveyor, track SLA, escalate delays |
| Surveyor | Performs field/asset assessment | view assigned claims, upload findings, add notes, conduct call/chat |
| TPA/Medical Operations User | Manages health-network claim collaboration | provider verification, cashless/reimbursement review, pre-auth handling |
| Compliance Reviewer | Reviews audit/fraud/compliance flags | masked data access, audit review, regulatory reporting |
| Read-only Executive | Monitors insurer operational KPI | dashboards, reports, no transactional changes |

### 5.2 Internal Counterpart Roles

- Labaid Insuretech operations team
- focal person
- business admin
- support team
- document operations
- finance/settlement counterpart

The insurer dashboard must support these cross-organization workflows without breaking tenant isolation.

---

## 6. Product Principles

The dashboard must follow these principles:

1. **Single workspace:** one insurer-facing portal for proposals, documents, and claims.
2. **Insurer-first operations:** the portal should reflect real insurer work queues, not generic admin tables.
3. **Document-heavy by design:** insurer workflows depend on templates, declarations, forms, evidence, and review notes.
4. **State-driven workflow:** every proposal and claim action must map to explicit statuses and transitions.
5. **Collaborative but auditable:** chat, calls, comments, and requests must be stored with timestamps and actor identity.
6. **Fallback-ready:** the UI should stay operational while live endpoints are phased in.
7. **Brand and insurer specificity:** LabaidInsuretech platform identity with insurer-specific content and templates.

---

## 7. Architectural Fit

### 7.1 Platform Alignment

The insurer dashboard shall be implemented as a **React/Next.js App Router portal** using the same portal/BFF pattern already documented for the platform:

- browser client
- Next.js portal frontend
- Next.js server routes/BFF layer
- gateway/API layer
- downstream microservices over REST/gRPC

### 7.2 Authentication and Session Model

The dashboard shall align with the existing portal session conventions documented in the architecture and API summaries:

- `session_token`
- `csrf_token`
- lightweight portal context cookies such as role, user id, business id, email, mobile

The BFF layer shall resolve portal headers for downstream service calls, including:

- `x-portal`
- `x-user-id`
- `x-business-id`
- `x-tenant-id`

### 7.3 Service Dependencies

The insurer dashboard is expected to integrate with or depend on:

- auth service
- authorization service
- insurance engine
- claims domain services
- document/media/storage services
- notification services
- analytics/reporting services
- fraud service
- workflow/tasking services
- partner/provider data where applicable

### 7.4 Event and Saga Alignment

The dashboard should be compatible with the platform’s event-driven direction, especially:

- insurer-facing proposal submission
- insurer approval/rejection events
- claim filed, approved, and settled events
- notification and audit side effects

The current architecture notes indicate the proposal schema/foundation exists, while end-to-end insurer decision saga wiring is still maturing. The dashboard must therefore support:

- live API mode where available
- hybrid mode where some actions are persisted locally or through interim BFF endpoints
- fallback/mock mode for incomplete downstream capabilities

---

## 8. Information Architecture

The insurer dashboard should expose the following primary navigation:

1. Dashboard
2. Proposals
3. Documents
4. Claim Settlement
5. Surveyor Desk
6. TPA & Claim Matrix
7. Product Categories
8. Plan Templates
9. Reports
10. Settings

### 8.1 Dashboard Landing Expectations

The dashboard home shall answer these questions immediately:

- how many proposals need decision today
- how many claims are pending review, survey, document completion, and settlement
- what is current claim TAT and settlement ratio
- which claims are overdue or flagged
- which document templates are missing or outdated
- what insurer/product categories are active
- what surveyor and TPA workloads require attention

---

## 9. Detailed Module Specification

### 9.1 Dashboard Home

#### Purpose

Provide an executive and operational summary for the insurer workspace.

#### Core Widgets

- proposal pipeline summary
- claims by status
- claims by product line
- surveyor workload summary
- TPA and provider claim summary
- overdue item tracker
- fraud/risk alert summary
- document template completeness summary
- recent activity feed
- announcements and insurer notices

#### Required Behaviors

- default insurer context should load as **PRAGATI INSURANCE** unless explicitly changed
- all widgets should be filterable by date range, product line, claim mode, and status
- KPI cards should support drill-down into the relevant work queue
- tiles must support empty states and fallback data mode

#### KPI Examples

- proposals pending SLA
- proposals approved today
- proposals rejected today
- active claims
- claims pending documents
- claims awaiting survey
- claims pending dual approval
- settled amount this month
- average claim TAT
- fraud-flagged claim count

---

### 9.2 Proposals Module

#### Purpose

Support insurer review of proposals sent by Labaid Insuretech before policy issuance.

#### Proposal Lifecycle

Recommended visible states:

- Draft at Labaid side
- Submitted to Insurer
- Under Review
- Clarification Requested
- Additional Documents Requested
- Approved
- Approved with Conditions
- Rejected
- Expired

#### Functional Requirements

- list and search proposals by customer, plan, product category, amount, and SLA
- open a structured proposal detail view
- display proposer, insured members, risk summary, premium summary, and attachments
- display linked digital forms and completed insurer templates
- allow insurer notes and decision rationale
- allow approve, reject, or request clarification
- allow conditions such as premium loading, exclusions, rider adjustment, or medical review
- expose document completeness score

#### Product-Specific Proposal Fields

Based on actual insurer proposal forms (see Appendix A), the proposal detail view must support product-specific field rendering:

- **Motor Insurance:** Vehicle details (CC, make, engine/chassis, seating), value segregation (glass/non-glass/electrical/accessories), underwriting Q&A (12 questions), coverage type selection (comprehensive/act-only), NCB entitlement
- **Fire Insurance:** Property location, construction details, building occupancy, sum insured per asset type (building/machinery/furniture/stock)
- **OMP (Travel Mediclaim):** Travel plan (Schengen/Non-Schengen, Plan A/B), medical history (6 conditions), declaration checkboxes, itinerary
- **Cattle/Livestock:** Farm location, animal roster (ear tag, species, weight, gender, value)
- **Group Enrollment:** Bulk member upload with dependent rows (spouse, children), NID/passport validation

#### Integration Notes

The module should align with the planned insurer proposal saga described in the PoliSync notes:

- payment-confirmed order can create insurer-facing proposal
- insurer decision should govern policy issuance or refund path
- future event compatibility should include proposal submitted, approved, and rejected topics

---

### 9.3 Documents Module

#### Purpose

Provide the insurer-owned repository for templates, form packs, and digital completion workflows.

#### Business Rule

For many products and claims, **Labaid Insuretech fills insurer-required forms before sending proposals or claim packs**. Pragati and similar insurers upload the source templates; the dashboard digitizes and manages them.

#### Document Types

- proposal forms
- health declarations
- member census sheets
- rate cards
- motor proposal sheets
- commercial vehicle forms
- fire/property forms
- livestock/pet forms
- travel mediclaim forms
- claim reimbursement packs
- discharge/settlement forms
- survey forms
- regulatory declarations

#### Functional Requirements

- template library grouped by insurer and product category
- upload support for insurer-owned source templates
- version tracking and effective dates
- digital rendering of Excel/PDF-derived field structures
- draft save and resume
- status tags such as active, draft, archived, superseded
- mapping between forms and proposal/claim workflows
- export to PDF/print/download package
- template preview before publish

#### Pragati Workbook Digitization

The portal should support digital forms derived from the Pragati workbook and similar future insurer workbooks:

- each sheet becomes a structured digital form
- rows/fields should map to typed inputs where possible
- insurer-required mandatory fields must be enforced
- repeated sections should support dynamic row addition where needed
- completed forms should be attachable to proposals and claims

**Confirmed Pragati workbook sheets for digitization** (14 sheets parsed from `pragati.xlsx`):

1. OMP Proposal Form — proposer details and plan selection
2. OMP Medical History — medical questions and declaration
3. OMP Declaration and Schengen reference — signed declarations
4. OMP Rate Card (Non-Schengen) — Plan A and Plan B by age/period
5. OMP Children exclusions and CFT/Employment rates (Non-Schengen)
6. OMP Rate Card (Schengen) — Plan A and Plan B by age/period
7. OMP Children exclusions and CFT/Employment rates (Schengen)
8. Private Vehicle Proposal Form — full vehicle and underwriting fields
9. Private Vehicle Proposal Form — continued (questions 7–14)
10. Fire Insurance Proposal Form — property, construction, occupation
11. Commercial Vehicle Proposal Form — vehicle, permits, capacity
12. Cattle/Livestock Proposal Form — animal details
13. Group Health Enrollment with Dependents
14. Health Insurance Claim Form — expense breakdown and document checklist

Field-level mapping for all 14 sheets is provided in Appendix A.

---

### 9.4 Claim Settlement Module

#### Purpose

Provide the claims review and settlement workspace for insurer claims teams.

#### Claim State Machine

The module shall align with the SRS state direction:

- Submitted
- Under Review
- Documents Requested
- Approved
- Rejected
- Payment Initiated
- Settled
- Closed

Additional operational statuses are allowed if they map cleanly to the core state machine, such as:

- Pending Survey
- Survey Report Received
- Pending Dual Approval
- Flagged for Investigation

#### Functional Requirements

- claim queue with filters by status, product line, amount band, claim mode, SLA, and flag
- responsive status tabs that always remain contained within their card or layout container
- claim detail page/panel with claimant summary, policy summary, incident details, amount, and documents
- document verification checklist
- reserve/settlement recommendation capture
- approval or rejection with reason
- escalation for high-value or suspicious claims
- fraud indicators and triggered rule display
- payment readiness checklist
- status timeline and audit history

#### Product-Specific Claims Document Checklists

Based on actual insurer document requirements (see Appendix B), the claim detail view must present a pre-populated document checklist by product line:

- **Motor claims:** 10 required items (+ 3 additional for theft), including 3 repair estimates from different workshops, MVI report, driver statement, and GD/FIR tracking
- **Fire claims:** 17 required items, including fire brigade report, 90-day stock verification, electrical compliance, and factory personnel statements
- **Health claims:** 8 categories of required documents, with itemized bill verification (accommodation, consultant, medicines, surgical, ancillary)
- **OMP claims:** TPA-routed through Crisis 24 or Van Ameyd, with international assistance workflow

#### Claim Mode Support

The dashboard must support clear operational handling for:

- cashless claims
- reimbursement claims
- pre-auth or provider-assisted claims where applicable

#### Approval Matrix Support

The claim workflow should support approval routing by amount and risk, including:

- auto-approve path for small low-risk claims where configured
- L1, L2, and L3 review routing
- dual approval path for claims requiring Business Admin and Focal Person approval
- investigation path when fraud score breaches threshold

---

### 9.5 Surveyor Desk

#### Purpose

Provide a dedicated workspace for survey-required claims instead of mixing survey tasks into the general settlement queue.

#### Business Rule

For Pragati non-life operations, categories such as **auto**, **fire/property**, and **pet/livestock** may require surveyor involvement before claim decision.

#### Surveyor Desk Features

- separate navigation tab and queue
- assigned survey claims list
- claim and asset details
- incident summary and location
- policy coverage summary
- field visit scheduling
- survey report submission
- image/video evidence upload
- damage assessment and estimate capture
- chat workspace
- web call/video verification workspace
- additional document request workflow
- recommendation outcome such as payable, partially payable, or not payable

#### Surveyor Workflow

1. Claim identified as survey-required.
2. Survey coordinator assigns surveyor.
3. Surveyor reviews case file and contacts claimant or counterpart.
4. Surveyor conducts visit or remote verification.
5. Surveyor uploads findings, media, notes, and recommendation.
6. Claims officer reviews survey result and continues settlement process.

#### Communication Requirements

Chat and web call must open proper portal interfaces, not placeholder actions.  
Document request must open a structured request form with:

- requested document list
- due date
- reason
- requesting actor
- notification outcome

---

### 9.6 TPA & Claim Matrix Module

#### Purpose

Provide health-claim operating intelligence and category-specific claims handling guidance.

#### Scope

- TPA operating model overview
- provider-network logic
- cashless vs reimbursement handling
- claim document matrices per plan category
- approval matrix visibility
- exclusion and escalation guidance

#### Functional Requirements

- plan-category matrix by product and claim type
- required documents by claim scenario
- TPA/insurer responsibility split
- provider validation status
- co-pay and network rules visibility
- reimbursement steps and settlement SLA
- exportable claim matrix

#### Example Category Dimensions

- group health
- individual health
- travel medical
- motor private car
- motor commercial vehicle
- fire/property
- pet/livestock
- SME health

---

### 9.7 Product Categories and Plan Templates

#### Purpose

Allow insurers to understand or configure product/category-level playbooks and reusable plan structures.

#### Functional Requirements

- category cards with insurer relevance and operational notes
- template creation from existing insurer form packs
- category-specific claim and proposal rules
- underwriting notes and exclusions
- renewal or endorsement considerations
- draft vs published template management

This module is especially useful for aligning insurer operations with the MD-aligned rollout phases across health, motor, P&C, agriculture, and life products.

---

### 9.8 Reports and Analytics

#### Purpose

Provide insurer-side operational and compliance reporting.

#### Required Reports

- proposal conversion and rejection analysis
- claim intake and settlement trend
- average TAT by product line
- claim ratio by plan category
- pending > 3 days, > 7 days, > 15 days, > 30 days
- surveyor SLA performance
- document deficiency patterns
- fraud and exception report
- insurer operational workload report
- template usage and outdated template report

#### Compliance-Oriented Reports

- monthly claims intimation summary
- quarterly claims settlement summary
- audit action logs
- overdue compliance tasks

---

## 10. Cross-Cutting Collaboration Features

The portal must support structured collaboration across proposals and claims.

### 10.1 Chat

- contextual conversation threads tied to entity id
- participant labels by organization and role
- file attachment support
- read timestamps
- searchable history

### 10.2 Web Call

- claim or proposal linked call room
- participant list
- call notes and outcome summary
- recording flag and audit reference if enabled by policy

### 10.3 Activity and Audit Timeline

- who did what
- when action occurred
- before/after status
- reason/comment
- related document or call reference

---

## 11. Data and Entity Model Expectations

The insurer dashboard should operate on or project data from these primary entities:

- insurer
- user
- session
- product
- product plan
- pricing configuration
- quotation
- order
- insurance proposal
- policy
- claim
- claim document
- claim approval
- settlement
- refund
- document template
- generated document
- notification
- survey assignment
- survey report
- provider/network record

### Additional Needed Operational Fields

Based on the platform notes, insurer-side claims will benefit from explicit support for:

- `claim_mode` such as cashless or reimbursement
- `network_provider_id`
- `co_pay_percentage`
- fraud score and triggered rule set
- survey requirement flag
- survey assignment status

---

## 12. Security, Compliance, and Audit

The insurer dashboard must follow the platform’s zero-trust and compliance requirements.

### Required Controls

- RBAC and ABAC aligned access
- tenant isolation between insurers
- masked PII where full visibility is not required
- JWT/session expiry and refresh discipline
- CSRF protection for portal actions
- file validation and malware scanning
- audit logging for all mutations
- approval traceability
- retention rules for policies, claims, and audit evidence

### Compliance Considerations

- IDRA reporting traceability
- BFIU/AML review hooks for suspicious transactions and claims
- no silent data deletion for records under regulatory hold
- strong document custody and versioning

---

## 13. Non-Functional Requirements for the Portal

The insurer dashboard should conform to platform NFRs while emphasizing desktop operational usability.

### Experience Targets

- page load target under 2 seconds for primary screens under normal load
- queue filtering/search under 200 ms for common searches
- responsive desktop-first layout with usable tablet fallback
- WCAG 2.1 AA target
- Bengali and English readiness where required

### Reliability Targets

- graceful degradation when downstream services are unavailable
- fallback data or retry guidance instead of blank failure states
- autosave for forms and long review workflows
- safe handling of partially integrated modules

### Maintainability Targets

- moduleized React components
- BFF route separation from UI components
- typed contracts for insurer dashboard data
- observability hooks for critical interactions

---

## 14. UI and UX Expectations

### Branding

- platform branding must use **LabaidInsuretech**
- login and shell experience should be insurer-oriented, not generic React marketing
- insurer context should be visible in the header

### Layout

- desktop-first operational dashboard
- clear card/grid system
- no overflowing filter chips or status tabs
- high-density tables with readable spacing
- fast access to primary work queues

### Interaction Standards

- tab labels must remain contained within their parent card or responsive grid
- no placeholder alerts for core actions
- modal, drawer, or dedicated panel patterns should be used for document requests, call setup, and structured chat
- empty, loading, and error states must be purposeful and informative

---

## 15. Proposed Screen Set

1. Login
2. Dashboard Home
3. Proposal List
4. Proposal Detail
5. Documents Library
6. Digital Form Workspace
7. Enrollment & Census Workspace
8. Pricing & Commercials Workspace
9. Claim Settlement Queue
10. Claim Detail
11. Claims Checklist & Required Documents Workspace
12. Surveyor Desk Queue
13. Surveyor Claim Workspace
14. Travel Assistance & TPA Contacts
15. TPA & Claim Matrix
16. Knowledge Center / Underwriting Playbooks
17. Reports
18. Settings

---

## 16. Critical Gaps Identified From `documentation/docs_forms`

Review of the actual insurer source files in `documentation/docs_forms` shows that the current portal specification still misses several critical tabs and workflow areas. These gaps are not optional polish items; they are directly implied by the operational documents already being exchanged between Labaid Insuretech and Pragati.

### 16.1 Enrollment & Census Workspace

The files `Enrollment Format Alpha Force.xlsx`, `Enrollment format Prime Shine.xlsx`, and `pragati.xlsx` member schedule sheets imply a dedicated enrollment workspace, not just a generic documents tab.

This workspace should provide:

- employee/member census upload and inline editing
- dependent and nominee capture
- validation for missing member attributes before insurer submission
- versioned enrollment snapshots by proposal or client group
- bulk add, bulk correction, and insurer-ready export

Without this tab, group business onboarding remains fragmented across document previews rather than becoming an operational workflow.

### 16.2 Pricing & Commercials Workspace

The files `Financial Proposal LifePlus_Shanta-2026 (Group Insurance).pdf` and `Insurance Coverage & Income.xlsx` imply a dedicated commercial workspace for insurer-facing pricing preparation and comparison.

This workspace should provide:

- quote builder for group life and medical proposals
- benefit slabs, sum assured, premium-per-thousand, and annual premium views
- side-by-side comparison of draft commercial scenarios
- insurer proposal cover-letter / quotation generation
- approval-ready pricing summary before dispatch

The current spec mentions proposals and documents, but not a structured commercial desk for pricing conversations and financial proposal preparation.

### 16.3 Claims Checklist & Required Documents Workspace

The document `Documents are normally required for Claims.docx` is effectively a claim-document operations playbook. It shows that claim handling requires category-specific evidence collection and validation, especially for fire and motor claims.

This workspace should provide:

- line-of-business claim checklists by claim type
- checklist completion status per claim
- missing-document flagging before survey or settlement review
- document responsibility assignment to Labaid, client, bank, surveyor, or claimant
- rules for mandatory evidence by claim category

Without this tab, claim operations are visible, but document completeness management remains manual.

### 16.4 Travel Assistance & TPA Contacts Tab

The file `OMP New  Claim Process.docx` clearly shows that travel claims involve external assistance providers and emergency contact instructions, including Crisis24 and Van Ameyde acting as TPA / assistance providers.

This workspace should provide:

- emergency contact directory with phone, email, and address
- travel-claim routing guidance
- claimant communication instructions
- TPA handoff tracking
- assistance-provider document dispatch log

This is more specific than the current generic TPA & claim matrix view and deserves a dedicated operational tab for travel and overseas mediclaim support.

### 16.5 Knowledge Center / Underwriting Playbooks

The presentation `Motor insurance policy-LifePlus Bangladesh(Final)(Underwrite & Claims).pptx` is a practical underwriting and claims training asset. It implies the portal needs a knowledge-center layer, not only transactional screens.

This workspace should provide:

- insurer playbooks by product line
- underwriting guidance and exclusion references
- product-specific claims procedure notes
- training decks, SOPs, and internal cheat sheets
- searchable operational references for teams

This is critical for consistency across proposals, claims, and surveyor review.

### 16.6 Source Template Archive and Source-vs-Digital Comparison

The presence of paired PDFs, DOCX files, PPTX references, and workbook-driven digital forms implies a need for source-template management beyond simple document listing.

This workspace should provide:

- insurer-uploaded source document archive
- source preview alongside digital rendered version
- template version history
- mapping status between source template and digital form
- insurer-specific ownership and publishing controls

This is especially important because Pragati uploads templates while Labaid Insuretech operationally completes and forwards them.

### 16.7 Multi-Page Document Composition Requirement

The source forms show that many insurer documents are structured as true paper forms with header, body, footer, signatures, declarations, and multi-page flow. Examples include the overseas mediclaim proposal, motor proposal form, and financial proposal layout.

The specification should explicitly require:

- page-based document rendering
- reusable insurer header and footer blocks
- multi-page preview and print flow
- PDF-ready output preserving page structure
- signature, declaration, and stamp-ready layout zones

This should be treated as a first-class UX and document-generation requirement, not a cosmetic enhancement.

---

## 17. Rollout Recommendation

### Phase A: Operational MVP

- login and insurer workspace context
- dashboard KPIs
- proposals list/detail/decision
- documents library
- enrollment and census workspace
- pricing and commercials workspace
- core claim queue
- audit timeline

### Phase B: Claims Maturity

- Surveyor Desk
- claims checklist workspace
- travel assistance and TPA contacts
- structured document request workflow
- embedded chat and web call interfaces
- approval matrix routing
- fraud indicator views

### Phase C: Intelligence and Compliance

- TPA & claim matrix
- knowledge center and underwriting playbooks
- source template archive and digital-template governance
- advanced reports
- insurer template publishing lifecycle
- regulatory reporting exports
- event-driven proposal/claim orchestration integration

---

## 18. Open Integration Dependencies

The following dependencies should be tracked during delivery:

- final live insurer proposal approval API or event callback path
- claim service support for full approval matrix and fraud score normalization
- document generation and template persistence APIs
- PDF export / print rendering pipeline for multi-page insurer forms
- enrollment import/export pipeline for census workbooks
- pricing and quotation persistence model for commercial proposals
- provider/network source of truth for health claim validation
- external assistance-provider / TPA contact registry and escalation flow
- communication services for persistent chat and WebRTC call capability
- notification service integration for request-document and decision events

---

## 19. Success Criteria

The insurer dashboard should be considered successful when:

- insurers can review and decide proposals inside the portal
- insurer forms are digitized, well-formatted, and attached to operational workflows
- group enrollment and census handling no longer depends on standalone spreadsheets
- commercial pricing and financial proposals are managed in-system
- claim processing no longer depends on ad hoc spreadsheets and messaging
- claim-document completeness is tracked with category-specific checklist logic
- survey-required claims are handled in a dedicated surveyor workflow
- travel assistance and TPA contact guidance is available in-system
- TPA and claim-matrix knowledge is available in-system
- insurer source templates can be governed, versioned, and rendered as print/PDF-ready digital forms
- SLA, fraud, and compliance visibility are available to insurer operations teams
- the portal remains aligned with the existing platform architecture and deployment model

---

## 20. Source Basis

This specification was derived from the current documentation set, especially:

- `documentation/About/ARCHITECTURE_OVERVIEW.md`
- `documentation/About/API_ROUTES_SUMMARY.md`
- `documentation/About/POLISYNC_REFERENCE.md`
- `documentation/About/START_HERE.md`
- `documentation/BRD/ALIGNMENT_SUMMARY.md`
- `documentation/core_plans/ACTIVE_WORKSTREAMS.md`
- `documentation/SRS_v3/SPECS_V3.7/sections/03_architecture.md`
- `documentation/SRS_v3/SPECS_V3.7/sections/04_functional_requirements.md`
- `documentation/SRS_v3/SPECS_V3.7/sections/05_non_functional_requirements.md`
- `documentation/SRS_v3/SPECS_V3.7/sections/06_data_model.md`
- `documentation/SRS_v3/SPECS_V3.7/sections/07_security_compliance.md`
- `documentation/SRS_v3/SPECS_V3.7/sections/08_integration.md`
- `documentation/docs_forms/Documents are normally required for Claims.docx`
- `documentation/docs_forms/Enrollment Format  Alpha Force.xlsx`
- `documentation/docs_forms/Enrollment format Prime Shine.xlsx`
- `documentation/docs_forms/Insurance Coverage & Income.xlsx`
- `documentation/docs_forms/Financial Proposal LifePlus_Shanta-2026 (Group Insurance).pdf`
- `documentation/docs_forms/OMP New  Claim Process.docx`
- `documentation/docs_forms/OMP Proposal Form (New).pdf`
- `documentation/docs_forms/Motor Insurance Proposal Form.pdf`
- `documentation/docs_forms/Fire Insurance Proposal Form_20230622_0001.pdf`
- `documentation/docs_forms/Motor insurance policy-LifePlus Bangladesh(Final)(Underwrite & Claims).pptx`

### Note on Source Path

The request referenced `documentation/docs`, but the repository currently uses a consolidated documentation tree under `documentation/` with the relevant material distributed across `About`, `BRD`, `core_plans`, and `SRS_v3`. This specification is based on that actual structure.

---

## APPENDIX A: Insurer Policy Forms — Field-Level Digitization Reference

This appendix is derived from actual insurer source documents in `documentation/docs_forms/`, including Pragati Insurance proposal forms, LifePlus Bangladesh motor insurance training material, claims document checklists, enrollment templates, premium rating structures, and financial proposals. It provides the field-level detail needed for digital form rendering, proposal intake, claims processing, and template management within the insurer dashboard.

---

### A.1 Motor Insurance — Private Vehicle Proposal Form (Pragati)

**Source:** `Motor Insurance Proposal Form.pdf`, `pragati.xlsx` (Sheet8–Sheet9)

#### A.1.1 Proposer Information Fields

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| Full Name of Proposer | text | yes | block letters |
| Address | text | yes | block letters |
| Mobile / Phone Number | text | yes | |
| Email Address | email | yes | |
| Business or Profession | text | yes | |

#### A.1.2 Vehicle Details Fields

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| Registration Mark & No. | text | yes | |
| Make of Vehicle | text | yes | manufacturer name |
| Engine No. | text | yes | |
| Chassis No. | text | yes | |
| Type of Body | text | yes | |
| Cubic Capacity (CC) | number | yes | used for premium calculation |
| Horse Power | number | no | |
| Year of Manufacture | number | yes | |
| Seating Capacity (incl. driver) | number | yes | |
| Carrying Capacity | text | no | goods vehicles |
| Vehicle Purchase Invoice / Sum Insured Value | currency (BDT) | yes | |
| Present Estimate Market Value | currency (BDT) | conditional | if invoice not available |

#### A.1.3 Value Segregation Fields

| Field | Type |
|-------|------|
| Vehicle excluding glass items | currency (BDT) |
| Glass items | currency (BDT) |
| Electrical Appliances (TV, Radio, AC) | currency (BDT) |
| Accessories | currency (BDT) |
| Full Insured Value | currency (BDT) — computed total |

#### A.1.4 Underwriting Questions

| # | Question | Answer Type |
|---|----------|-------------|
| 1a | Will the car be used solely for social, domestic, and pleasure purposes? | yes/no |
| 1b | If not, state other uses | text |
| 2 | Are you the owner and is it registered in your name? If not, state name/address of owner | yes/no + text |
| 3a | Date of purchase | date |
| 3b | Whether new or secondhand | select |
| 3c | Price paid | currency |
| 3d | Present estimate market value | currency |
| 4 | Do you or any known driver suffer from defective vision/hearing or any physical infirmity? | yes/no |
| 5 | Have you or any known driver been convicted of any motor offence in the past 5 years? | yes/no |
| 6 | How long have you been driving continuously? | text |
| 7 | Are you now or have you been insured for any motor vehicle? State name of underwriter | yes/no + text |
| 8 | Entitled to No Claim Bonus from previous insurer? Attach renewal notice | yes/no |
| 9a | Has any underwriter declined your proposal or cancelled/refused to renew? | yes/no |
| 9b | Required you to bear first cost of accident/loss? | yes/no |
| 9c | Imposed special conditions or increased premium? | yes/no |
| 10 | Particulars of accidents/losses in past 3 years | text |
| 11 | Coverage type required: Comprehensive / Act Only / Motor Vehicles Act limited | select |
| 12 | First-part cost of each accident/loss you wish to bear (excess) | currency |
| 13 | Insure rugs, coats, luggage (limited Tk 1500 per occurrence)? | yes/no |
| 14 | Any other benefits to insure? | text |

#### A.1.5 Policy Commencement and Declarations

| Field | Type |
|-------|------|
| Policy Commence Date | date |
| Proposer Signature | signature |
| Declaration warranty acceptance | checkbox |
| Bank Use: Branch Name | text |
| Bank Use: Account Number | text |
| Bank Use: Account Name | text |
| Bank Use: Officer Code | text |

---

### A.2 Motor Insurance — Commercial Vehicle Proposal Form (Pragati)

**Source:** `pragati.xlsx` (Sheet11)

#### A.2.1 Additional Fields Beyond Private Vehicle

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| Licensed Carrying Capacity — Goods (tons) | number | conditional | |
| Licensed Carrying Capacity — Passengers (excl. driver & cleaner) | number | conditional | |
| Trailer value | currency | no | |
| Other attachment value | currency | no | |
| Non-electrical accessories value | currency | no | |
| Fitted with dual rear wheels and double springs? | yes/no | yes | |
| Permit type: Private/Public/Stage/Contract Carrier | select | yes | |
| Operating area: Dhaka only or state locations | text | yes | |
| Usual garage location | text | yes | |
| Vehicle in perfect condition? | yes/no | yes | |
| Past accident history | text | yes | |
| Past third-party claims | text | yes | |
| Previous motor insurance? Name of corporation | text | yes | |

---

### A.3 Motor Insurance — Premium Calculation Model

**Source:** `Motor insurance policy-LifePlus Bangladesh(Final)(Underwrite & Claims).pptx` (Slides 12–19)

#### A.3.1 Comprehensive — Four Wheelers (Private) Tariff

| CC Range | Basic Rate (BDT) | Act Liability Rate (BDT) |
|----------|-------------------|--------------------------|
| Up to 1300 / Up to ½ Ton | 2,795 | 150 |
| 1301–1800 / ½ to 1 Ton | 2,873 | 250 |
| 1801–3000 / Up to 3 Ton | 2,925 | 350 |
| 3001+ / Above 3 Ton | 2,990 | 450 |

**Formula:** `Premium = Act Liability (by CC) + Basic Own Damage (by CC) + (Sum Insured × 2.65%) + 15% VAT`

Additional: Driver rate BDT 30, Per passenger BDT 45.

#### A.3.2 Perils Breakdown — Four Wheelers (2.65% total)

| Peril | Rate |
|-------|------|
| Fire | 0.50% |
| Theft | 0.50% |
| Earthquake | 0.25% |
| Flood & Cyclone | 0.25% |
| Riot & Strike | 0.50% |
| Others | 0.65% |
| **Total** | **2.65%** |

#### A.3.3 Comprehensive — Two Wheelers Tariff

| CC Range | Basic Rate (BDT) | Act Liability Rate (BDT) |
|----------|-------------------|--------------------------|
| Up to 150 | 200 | 100 |
| 151–250 | 275 | 130 |
| 251–350+ | 350 | 160 |

**Formula:** `Premium = Act Liability (by CC) + Basic Own Damage (by CC) + (Vehicle Value × 2.15%) + 15% VAT`

Additional: Rider rate BDT 50, Passenger BDT 45.

#### A.3.4 Perils Breakdown — Two Wheelers (2.15% total)

| Peril | Rate |
|-------|------|
| Fire and Theft | 1.00% |
| Riot & Strike | 0.50% |
| Flood & Cyclone | 0.25% |
| Earthquake | 0.25% |
| Others | 0.15% |
| **Total** | **2.15%** |

#### A.3.5 No Claim Bonus (NCB)

NCB applies to own damage premium only and is a renewal rebate.

| Year | Private Vehicle | Commercial Vehicle | Motorcycle |
|------|----------------|-------------------|------------|
| 1st year (no claim in preceding year) | 30% | 30% | 15% |
| 2nd year (no claim in preceding 2 years) | 40% | 40% | 20% |
| 3rd year (no claim in preceding 3 years) | 50% | 50% | 25% |

**Rule:** If insured makes a claim, one claim means two steps backward on the NCB schedule.

#### A.3.6 Claim Loading (Own Damage Section Only)

| Claims in Preceding Period | Private Vehicle | Commercial Vehicle |
|---------------------------|----------------|-------------------|
| One claim | 30% | 30% |
| Two claims | 40% | 40% |
| Three claims | 50% | 50% |

#### A.3.7 AVLS (Automatic Vehicle Location System) Discount

- 10% discount applicable on Own Damage Premium when vehicle is fitted with AVLS.
- Flood, Cyclone & Earthquake exclusion results in reduced perils total of 2.15%.

#### A.3.8 Act Liability Coverage Limits

| Liability | Amount (BDT) |
|-----------|-------------|
| Death | 20,000 |
| Severe Hurt | 10,000 |
| Any Other Hurt | 5,000 |
| Property Damage | 50,000 |

---

### A.4 Motor Insurance — Comprehensive Coverage Scope

**Source:** PPTX Slides 10–11

#### Covered Perils (Own Damage)

- Fire, explosion, self-ignition, or lightning
- Burglary, housebreaking, or theft
- Riot and strike including malicious and terrorism activities
- Earthquake (fire and shock damage)
- Flood, typhoon, hurricane, storm, tempest, inundation, cyclone, hailstorm, frost
- Accidental external means
- While in transit by road, rail, inland waterway, lift, elevator, or air

#### General Exclusions

- Accidents, loss, or liability outside the geographical area
- Claims arising from contractual liability
- Loss from violating the vehicle's stated use
- Being driven by uninsured/unauthorized person
- Nuclear weapons, war, invasion, foreign enemy, hostilities, civil war, mutiny, rebellion
- Ionizing radiation or contamination
- Terrorist activity (as per circular)

---

### A.5 Fire Insurance Proposal Form (Pragati)

**Source:** `pragati.xlsx` (Sheet10), `Fire Insurance Proposal Form_20230622_0001.pdf`

#### A.5.1 Proposer and Property Fields

| Field | Type | Required |
|-------|------|----------|
| Full name, address, trade/profession | text | yes |
| Term of insurance: from–to | date range | yes |
| **Sum Insured Breakdown** | | |
| — On Building | currency (BDT) | conditional |
| — On Machinery | currency (BDT) | conditional |
| — On Furniture and Effects | currency (BDT) | conditional |
| — On Merchandise / Stock-in-Trade | currency (BDT) | conditional |

Multiple buildings/locations supported (numbered items).

#### A.5.2 Location Details

| Field | Type |
|-------|------|
| Name of Building | text |
| Owner of Building | text |
| Plot Number | text |
| Holding Number | text |
| Name of Street | text |
| Town | text |
| District | text |

#### A.5.3 Construction Details

| Field | Type |
|-------|------|
| Number of Storeys | number |
| Walls / Boundary material | text |
| Roof material | text |
| Floors in each storey | text |
| Adjoining building details | text |
| Buildings within 50 feet | text |

#### A.5.4 Occupation and Usage

| Field | Type |
|-------|------|
| Name of Banker | text |
| Ground floor occupation | text |
| First floor occupation | text |
| Upper floor occupation | text |
| Building usage purpose | text |

#### A.5.5 Lighting, Heating and Power

| Field | Type |
|-------|------|
| How is building lighted? | text |
| How is building heated? | text |
| Full particulars of power used | text |

---

### A.6 Overseas Mediclaim Policy (OMP) Proposal Form (Pragati)

**Source:** `OMP Proposal Form (New).pdf`, `pragati.xlsx` (Sheets 1–7)

#### A.6.1 Proposer Fields

| Field | Type | Required |
|-------|------|----------|
| Name and status (as in passport) | text | yes |
| Title (Mr./Mrs./Miss/Master) | select | yes |
| Residence address | text | yes |
| Residence telephone & mobile | text | yes |
| Actual occupation | text | yes |
| Office name and address | text | no |
| Office telephone | text | no |
| Age (completed years) | number | yes |
| Passport number (copy attached) | text | yes |

#### A.6.2 Travel Plan Fields

| Field | Type | Required |
|-------|------|----------|
| Plan Type: Schengen / Non-Schengen | select | yes |
| Plan Level: Plan A (excl USA/Canada) / Plan B (incl USA/Canada) | select | yes |
| Purpose of trip (official/holiday conducted/holiday individual) | select | yes |
| Date of departure | date | yes |
| Number of days stay abroad | number | yes |
| Itinerary (countries, places, days at each) | text | yes |

#### A.6.3 Medical History Section

| # | Question | Answer Type |
|---|----------|-------------|
| 1 | In good health, free from physical and mental disease/infirmity? | yes/no |
| 2a | Nervous, mental, psychiatric disease, spinal disorder, fainting, blackout, paralysis? | yes/no + details |
| 2b | High blood pressure, heart disease, ischemic heart disease, piles, varicose veins, circulatory disorders, rheumatic fever? | yes/no + details |
| 2c | Hernia, rheumatic/joint disease, urinary disease, diabetes? | yes/no + details |
| 2d | Respiratory/allergic disease, stomach/bowel/gallbladder disorder? | yes/no + details |
| 2e | Any complaint requiring specialist consultation, surgery, or hospital treatment? | yes/no + details |
| 2f | Any complaint or tendency that may require such treatment in future? | yes/no + details |
| 3 | Additional facts affecting proposed insurance to disclose? | yes/no + details |
| 4 | Intention of engaging in winter sports or injury-liable pastimes? | yes/no |
| 5 | Illness, disease, or accident in 12 months preceding insurance | table (nature, date, practitioner details) |

#### A.6.4 Pre-existing Conditions Declaration

Table for known ailments that may require medical attention abroad (up to 4 entries).

#### A.6.5 Declarations (checkboxes)

1. Not travelling against physician advice
2. Not on waiting list for medical treatment
3. Not travelling for the purpose of obtaining medical treatment
4. No terminal prognosis received

#### A.6.6 Product Benefits and Limitations

| # | Benefit | Limit |
|---|---------|-------|
| 01 | Medical Expenses & Hospitalization (excl USA/Canada) | USD 50,000, Excess USD 100 |
| 02 | Medical Expenses & Hospitalization (incl USA/Canada) | USD 100,000, Excess USD 100 |
| 03 | Medical Expenses for Schengen Countries | EUR 30,000, Nil deductible |
| 04 | Transport/Repatriation (illness/accident) | Actual Expenses |
| 05 | Emergency Dental Care | USD 500, Excess USD 50 |
| 06 | Repatriation of Family Travelling with Insured | Actual Expenses |
| 07 | Repatriation of Mortal Remains | Actual Expenses |
| 08 | Travel of One Immediate Family Member | USD 100/day, max USD 1,000 |
| 09 | Emergency Return (death of close family member) | Actual Expenses |

**Exclusion:** Pre-existing medical condition, suicide/attempted suicide, mental illness, pregnancy/childbirth.

---

### A.7 OMP Premium Rating Structure (Pragati)

**Source:** `pragati.xlsx` (Sheets 4–7)

#### A.7.1 Non-Schengen — Plan A (Worldwide excl USA/Canada) — Sample Rates (BDT)

| Period (Days) | 0.5–40 yrs | 41–50 yrs | 51–55 yrs | 56–59 yrs | 60–65 yrs | 66–70 yrs | 71–76 yrs | 76–79 yrs | 80–85 yrs |
|---------------|-----------|-----------|-----------|-----------|-----------|-----------|-----------|-----------|-----------|
| 1–14 | 1,239 | 1,860 | 2,499 | 2,499 | 2,499 | 8,131 | 14,230 | 28,460 | 40,657 |
| 15–21 | 1,291 | 1,983 | 2,664 | 2,664 | 2,664 | 8,668 | 15,169 | 30,338 | 43,339 |
| 22–28 | 1,640 | 2,210 | 2,976 | 2,970 | 2,976 | 9,678 | 16,937 | 33,873 | 48,390 |
| 29–35 | 1,783 | 2,673 | 3,592 | 3,592 | 3,592 | 11,692 | 20,462 | 40,924 | 56,455 |
| 36–47 | 2,042 | 3,063 | 4,119 | 4,119 | 4,119 | 13,402 | 23,454 | 46,908 | 67,011 |
| 48–60 | 2,402 | 3,629 | 4,879 | 4,879 | 4,879 | 15,873 | 27,779 | 55,558 | 79,368 |
| 61–75 | 2,970 | 4,480 | 6,022 | 6,022 | 6,022 | 19,598 | 34,297 | 68,592 | No Cover |
| 76–90 | 3,550 | 5,311 | 7,140 | 7,140 | 7,140 | 23,232 | 40,657 | 81,312 | No Cover |
| 91–120 | 6,003 | 9,063 | 12,184 | 12,184 | 12,184 | No Cover | No Cover | No Cover | No Cover |
| 121–147 | 7,230 | 10,869 | 14,613 | 14,613 | 14,613 | No Cover | No Cover | No Cover | No Cover |
| 148–180 | 10,046 | 14,992 | 20,326 | 20,326 | 20,326 | No Cover | No Cover | No Cover | No Cover |

#### A.7.2 Corporate Frequent Travel (Annual Cover, Business Only)

| Age Band | Non-Schengen Premium (BDT) | Schengen Premium (BDT) |
|----------|---------------------------|------------------------|
| 18–40 | 11,704 | 14,630 |
| 41–59 | 27,511 | 34,390 |

Maximum period any one trip: 30 days.

#### A.7.3 Employment & Studies — Monthly Rates

**Employment (collected in foreign currency):**

| Plan | Age 18–40 | Age 41–59 |
|------|-----------|-----------|
| Plan C (excl USA/Canada) — participant/spouse | USD 64.56 (Non-Schengen) / USD 80.69 (Schengen) | USD 103.13 / USD 128.92 |
| Plan D (incl USA/Canada) — participant/spouse | USD 111.88 / USD 139.86 | USD 214.20 / USD 267.76 |
| Plan C — child (under 18) | USD 22.63 | — |
| Plan D — child (under 18) | USD 91.67 | — |

**Studies (collected in BDT):**

| Plan | Age 18–40 | Age 41–59 |
|------|-----------|-----------|
| Plan C | BDT 1,861.11 | BDT 2,810.00 |
| Plan D | BDT 2,995.56 | BDT 5,666.67 |

#### A.7.4 Deductibles

- Non-Schengen: USD 100
- Emergency Dental Care: USD 50
- Schengen: No deductible

#### A.7.5 Children Under 5 Exclusions

Excludes cover for: Mumps, Chicken Pox, Measles, German Measles, Spina Bifida, Whooping Cough, Diphtheria, Poliomyelitis, Meningitis, Scarlet Fever — and consequences thereof.

---

### A.8 Cattle/Livestock Insurance Proposal Form (Pragati)

**Source:** `pragati.xlsx` (Sheet12)

| Field | Type | Required |
|-------|------|----------|
| Name of Insured | text | yes |
| Address | text | yes |
| Occupation | text | yes |
| Farm Location | text | yes |
| Period of Insurance | date range | yes |

**Animal Details (repeating rows):**

| Field | Type |
|-------|------|
| Sl# | number |
| Ear Tag / Identification | text |
| Date of Birth / Purchase | date |
| Gender | select (Male/Female) |
| Color | text |
| Weight (kg) | number |
| Value / Sum Insured (BDT) | currency |
| Species | text |

---

### A.9 Health Insurance Claim Form (Pragati)

**Source:** `pragati.xlsx` (Sheet14)

#### A.9.1 Claim Header Fields

| Field | Type | Required |
|-------|------|----------|
| Name of Organization | text | yes |
| Name of Employee | text | yes |
| Name of Patient | text | yes |
| Relationship with Employee | select (Husband/Wife/Son/Daughter/Self) | yes |
| Date of Prior Intimation | date | yes |
| Membership No. | text | yes |
| Hospital / Clinic Name and Address | text | yes |
| Date of Admission | date | yes |
| Date of Discharge | date | yes |

#### A.9.2 Hospitalization Expense Breakdown

| Expense Category | Type |
|------------------|------|
| Hospital Accommodation | currency (BDT) |
| Consultant's Fee | currency (BDT) |
| Routine Investigations | currency (BDT) |
| Medicines / Drugs | currency (BDT) |
| Surgical Charges | currency (BDT) |
| Ancillary Services | currency (BDT) |
| Others | currency (BDT) |
| **Total** | currency (BDT) — computed |

#### A.9.3 Required Claim Documents Checklist

1. Copy of prior claim intimation record (or telephonic intimation date)
2. Doctor's prescription(s) with duration of complaints, diagnosis, and hospitalization advice (original). For maternity: LMP, EDD, and Gravida required
3. Discharge certificate with history, diagnosis, treatment/operation note, admission/discharge dates
4. Certificate from employer/educational institution regarding absence during illness
5. Photocopy of patient's treatment records while confined
6. Hospital bill supported by original money receipt
7. All diagnostic reports with original receipts, supported by doctor's advice
8. Original bills specifying:
   - a) Accommodation charges (daily charge × days)
   - b) Consultant's fee (doctor's bill/receipts with date)
   - c) Medicines/drugs (name, quantity, price with prescription)
   - d) Surgical charges (breakdown: surgeon, OT, anesthetist, assistants)
   - e) Ancillary services (labor room, post-op care, oxygen, ICU, blood transfusion, equipment, dressing, non-routine tests, ambulance)

---

## APPENDIX B: Claims Document Requirements by Product Line

**Source:** `Documents are normally required for Claims.docx`, PPTX Slides 30–32

### B.1 Fire Insurance — Required Claim Documents

1. Claim form duly filled and signed by insured and concerned bank
2. Fire brigade report (original)
3. Daily stock report for 90 days up to date of incident (countersigned by bank)
4. Monthly stock statement (countersigned by bank)
5. Stock register photocopy for last 6 months
6. Tally book photocopy for 6 months of affected godown
7. Purchase documents/invoices for affected/lost/damaged stock
8. Challan and bill for local purchases
9. Fire license copy
10. Trade license copy
11. PDB/REB permission letter (or BERC generator certificate for self-generator)
12. Periodic electrical test report by competent authority
13. Generator logbook
14. Statements from: Factory In-charge, Electrical Engineer, Godown Keeper, Guard on duty (attested by authority)
15. GD entry copy
16. Photocopy of related insurance policies from other companies
17. Layout plan of entire mill and godown

### B.2 Motor Insurance — Required Claim Documents

1. Claim intimation letter
2. Claim form duly filled, signed, and sealed (by insured and bank)
3. GD entry copy (sealed and signed by concerned police station)
4. Three repair estimates from different reputed motor workshops/garages
5. Driver's written statement on cause of accident (countersigned by insured)
6. Registration certificate copy
7. Tax token, fitness certificate, route permit copies
8. Motor Vehicle Inspector (MVI) report
9. Driver's license copy (attested by BRTA)
10. Original challan of carrying goods at time of accident

**Additional for theft claims:**
11. FIR to police
12. Police final investigation report (certified by court)
13. Survey report (original)

### B.3 OMP (Overseas Mediclaim) — Claims Process

**TPA:** Van Ameyd (UK)
- Address: Office G 18, Bromley Old Town Hall, 30 Tweedy Road, Bromley, BR1 3FE
- Telephone: +44 208 315 0732

**Emergency Assistance:** Crisis 24
- Address: 2 London Bridge, London SE1 9RA, UK
- Telephone: +44 207 902 7131
- Email: opsassist@crisis24.com / corporateteam@crisis24.com

**Claim Flow:**
1. In the event of illness or accident abroad requiring hospital treatment or trip curtailment, contact Crisis 24 immediately
2. Request claim form from Crisis 24 or Van Ameyd
3. Complete claim form and submit to Crisis 24 / Van Ameyd with insurance certificate and relevant documentation
4. To avoid reverse-charge calls, claimant provides their phone number for callback

---

## APPENDIX C: Group Enrollment Data Formats

**Source:** `Enrollment Format Alpha Force.xlsx`, `Enrollment format Prime Shine.xlsx`

### C.1 Standard Enrollment Format (Alpha Force Model)

| Column | Type | Required | Notes |
|--------|------|----------|-------|
| SL No. | number | yes | sequential |
| Name | text | yes | full name |
| Nominee Name | text | no | |
| Enrollment ID | text | auto-generated | platform-assigned |
| Date of Birth | date | yes | |
| Insured NID / Birth Certificate / Passport Number | text | yes | |
| Age | number | computed | from DOB |
| Gender | select (Male/Female) | yes | |
| Mobile No. | text | yes | |

### C.2 Extended Enrollment Format (Prime Shine Model)

Adds:

| Column | Type | Notes |
|--------|------|-------|
| Coverage Start Date | date | policy effective |
| Coverage End Date | date | policy expiry |

### C.3 Group Health Enrollment with Dependents

**Source:** `pragati.xlsx` (Sheet13)

Structure supports repeating groups per employee:

| Row Type | Fields |
|----------|--------|
| Employee | Name, Ear Tag/ID, Designation/Relation, Gender, DOB, Age |
| Spouse | Name, Designation/Relation, Gender, DOB, Age |
| Child-1 | Name, Designation/Relation, Gender, DOB, Age |
| Child-2 | Name, Designation/Relation, Gender, DOB, Age |

**Digitization note:** The dashboard should support dynamic dependent rows per member, with add/remove capability.

### C.4 Enrollment Data Volumes (Reference)

- Alpha Force: 339 enrolled members
- Prime Shine: 386 enrolled members

---

## APPENDIX D: Financial Proposal and Revenue Model Reference

**Source:** `Financial Proposal LifePlus_Shanta-2026 (Group Insurance).pdf`, `Insurance Coverage & Income.xlsx`

### D.1 Group Life & Medical Insurance — Sample Proposal (LifePlus / Shanta Life)

| Plan | Type of Benefit | Members | Total Sum Assured | Premium Rate (per Tk 1000 SA) | Annual Premium (BDT) |
|------|----------------|---------|-------------------|-------------------------------|---------------------|
| Life | Group Life (GL) 100K flat | 3,000 | 30,00,00,000 | 2.44 | 7,32,000 |
| Life | ADB 200K | — | — | 1.10 | 3,30,000 |
| **Life Total** | | | | **3.54** | **10,62,000** |
| Health | IPC-Hospitalization 20K (coverage 2,50,000/member/year) | 3,000 | — | 410/member/year | 12,30,000 |
| Health | OPC General 2K (coverage 50,000/member/year) | 3,000 | — | 550/member/year | 11,00,000 |
| **Health Total** | | | | | **23,30,000** |
| **Grand Total** | | | | | **33,92,000** |

**Key terms:**
- PPD & PTD schedule compliant with First Schedule of the Labour Law of Bangladesh
- All applicable VAT borne by insurer (Shanta Life)
- Premium payable yearly in advance
- Offer valid for 2 months

### D.2 Platform Revenue Structure

| Component | Rate/Amount | Notes |
|-----------|------------|-------|
| Member Premium | BDT 550 (low) / 750 (standard) | per member |
| Platform Charge | 20% of premium | |
| VAS (LifePlus Service) | 25% of premium | |
| Others / MC | BDT 52–112 | |
| VAT | 15% | |

**Income to LabaidInsuretech (LICL):**

| Revenue Stream | Rate |
|----------------|------|
| Commission on Premium | 15% |
| Platform Charge | 100% retained |
| VAS Income from LPB | 20% |

**Channel Split:** B2B 50%, B2D 30%, B2C 20%

---

## APPENDIX E: Dashboard Digitization and Workflow Implications

This section maps findings from the parsed insurer documents to specific dashboard module requirements and digital form rendering strategy.

### E.1 Form Digitization Strategy

Based on the parsed source documents, the insurer form library must support:

| Source Format | Form Category | Digital Rendering |
|---------------|---------------|-------------------|
| PDF/XLSX | Private Vehicle Proposal | structured form with ~25 fields + underwriting Q&A |
| PDF/XLSX | Commercial Vehicle Proposal | extended vehicle form with permit/capacity fields |
| PDF/XLSX | Fire Insurance Proposal | multi-section form (property, construction, occupation, power) |
| PDF/XLSX | OMP Proposal | multi-page form with medical history, travel plan, declarations |
| XLSX | Cattle/Livestock Proposal | header + repeating animal detail rows |
| XLSX | Health Claim Form | expense breakdown table + document checklist |
| XLSX | Group Enrollment | bulk data capture with dependent rows |
| PPTX (training) | Motor Premium Calculator | embedded calculator widget for underwriting officers |

### E.2 Proposal Module Enrichment

The proposal detail view (Section 9.2) should now render:

- **Motor proposals:** Vehicle details, CC-based premium calculation, NCB/loading history, underwriting Q&A answers, value segregation
- **Fire proposals:** Property location, construction details, occupancy, sum insured breakdown by asset type
- **OMP proposals:** Travel plan, medical history, plan type/level selection, benefits summary
- **Livestock proposals:** Farm details, animal roster with identification

### E.3 Claims Module Enrichment

The claims workspace (Section 9.4) should now support:

- **Product-specific document checklists** pre-populated per claim type (fire: 17 items, motor: 10–13 items, health: 8 categories)
- **Motor claims**: Automated survey requirement flag based on claim type, damage estimates comparison (3 workshops), MVI report upload, GD/FIR tracking
- **Fire claims**: 90-day stock verification flow, fire brigade report, electricity compliance documents
- **Health claims**: Hospitalization expense breakdown entry, itemized bill verification, maternity-specific fields (LMP/EDD/Gravida)
- **OMP claims**: TPA routing to Crisis 24 or Van Ameyd, international assistance workflow

### E.4 Premium Calculator Requirements

The insurer dashboard should provide a built-in premium calculator for motor insurance underwriting:

**Inputs:**
- Vehicle type (private 4-wheeler / 2-wheeler / commercial)
- CC range
- Sum insured value
- Seating capacity
- Selected perils (with option to exclude flood/cyclone/earthquake for AVLS discount)
- NCB status (years without claim)
- Claim loading status (claims in preceding period)
- AVLS fitted (yes/no)

**Outputs:**
- Act liability premium (by CC)
- Basic own damage premium (by CC)
- Peril-based premium (sum insured × rate)
- NCB discount (if applicable)
- Claim loading surcharge (if applicable)
- AVLS discount (if applicable)
- Subtotal
- 15% VAT
- Total premium

### E.5 Surveyor Desk Enrichment

The surveyor desk (Section 9.5) should require for motor claims:
- Survey report upload (original)
- 3 repair estimates comparison workspace
- Driver statement review
- MVI report review
- Damage assessment linked to vehicle value segregation (glass vs non-glass vs electrical vs accessories)
- Recommendation with breakdown against each coverage component

### E.6 Document Template Library — Minimum Set

Based on parsed documents, the template library (Section 9.3) must include at minimum:

| Template | Product Line | Source |
|----------|-------------|--------|
| Private Vehicle Proposal Form | Motor | Pragati PDF/XLSX |
| Commercial Vehicle Proposal Form | Motor | Pragati XLSX |
| Fire Insurance Proposal Form | Fire/Property | Pragati XLSX |
| OMP Proposal Form | Travel Mediclaim | Pragati PDF/XLSX |
| Cattle/Livestock Proposal Form | Agriculture/Livestock | Pragati XLSX |
| Health Insurance Claim Form | Health | Pragati XLSX |
| Group Enrollment Format (Standard) | Group Life/Health | Alpha Force XLSX |
| Group Enrollment Format (Extended) | Group Life/Health | Prime Shine XLSX |
| Group Enrollment with Dependents | Group Health | Pragati XLSX |
| Motor Claims Document Checklist | Motor | DOCX |
| Fire Claims Document Checklist | Fire/Property | DOCX |
| Motor Premium Tariff Card | Motor | PPTX training |
| OMP Premium Rate Card | Travel Mediclaim | Pragati XLSX |
| Financial Proposal Template | Group Life/Medical | Shanta Life PDF |

### E.7 Schengen Country Reference Data

The OMP module requires a maintained list of Schengen countries for plan type validation:

Austria, Belgium, Denmark, Finland, France, Germany, Iceland, Italy, Greece, Luxembourg, Netherlands, Norway, Portugal, Spain, Sweden, Estonia, Latvia, Lithuania, Poland, Czech Republic, Slovakia, Hungary, Slovenia, Malta, Cyprus, Switzerland, Liechtenstein.

---

## 20. Document Source Reference (Updated)

This specification now additionally incorporates data from the following source documents in `documentation/docs_forms/`:

| Document | Type | Content |
|----------|------|---------|
| Motor insurance policy-LifePlus Bangladesh(Final)(Underwrite & Claims).pptx | PPTX (33 slides) | Motor insurance types, coverage scope, premium tariffs (4-wheeler & 2-wheeler), NCB, claim loading, AVLS discount, exclusions, required documents, claims procedure |
| Documents are normally required for Claims.docx | DOCX | Fire and motor claims document checklists |
| OMP New Claim Process.docx | DOCX | OMP claims process via Van Ameyd (TPA) and Crisis 24 |
| Enrollment Format Alpha Force.xlsx | XLSX | Group enrollment data (339 members) |
| Enrollment format Prime Shine.xlsx | XLSX | Group enrollment data (386 members) with coverage dates |
| Insurance Coverage & Income.xlsx | XLSX | Premium structure, income breakdown, channel revenue split |
| pragati.xlsx | XLSX (14 sheets) | OMP proposal form, OMP premium rating (Schengen/Non-Schengen), private/commercial vehicle proposal forms, fire proposal form, cattle/livestock proposal, group health enrollment, health claim form |
| Financial Proposal LifePlus_Shanta-2026 (Group Insurance).pdf | PDF | Group life + medical insurance financial proposal |
| Fire Insurance Proposal Form_20230622_0001.pdf | PDF | Scanned fire proposal (low OCR quality) |
| Motor Insurance Proposal Form.pdf | PDF | Pragati private vehicle proposal form |
| OMP Proposal Form (New).pdf | PDF | Pragati overseas mediclaim proposal form |
