# LabAid InsureTech - Insurance Engine সম্পূর্ণ ডকুমেন্টেশন

## 📋 সূচিপত্র

1. [সিস্টেম ওভারভিউ](#সিস্টেম-ওভারভিউ)
2. [আর্কিটেকচার](#আর্কিটেকচার)
3. [মডিউল সমূহ](#মডিউল-সমূহ)
4. [ডাটাবেস স্কিমা](#ডাটাবেস-স্কিমা)
5. [সার্ভিস ইন্টিগ্রেশন](#সার্ভিস-ইন্টিগ্রেশন)
6. [ওয়ার্কফ্লো এবং স্টেট মেশিন](#ওয়ার্কফ্লো-এবং-স্টেট-মেশিন)
7. [অনুমোদন ম্যাট্রিক্স](#অনুমোদন-ম্যাট্রিক্স)

---

## 🎯 সিস্টেম ওভারভিউ

### প্রজেক্ট তথ্য

- **নাম:** LabAid InsureTech Platform - Insurance Engine Microservice
- **টেকনোলজি:** .NET 10
- **পোর্ট:** 5001
- **আর্কিটেকচার:** Vertical Slice Architecture (VSA) with CQRS
- **ডাটাবেস:** PostgreSQL
- **SRS সংস্করণ:** v3.11

### প্রজেক্ট স্ট্রাকচার

```
D:\InsureTech\
├── proto\                                    (Shared proto repository)
├── gen\                                      (Generated code from proto)
└── backend\
    └── insurance_engine\
        └── src\
            ├── InsuranceEngine.SharedKernel\
            │   └── Persistence\
            │       ├── Entities\             (Core shared entities)
            │       └── Migrations\           (EF Core migrations)
            └── InsuranceEngine.[Module]\     (Per-module projects)
                ├── GrpcServices\             (gRPC service implementations)
                ├── Application\
                │   ├── Commands\             (MediatR command handlers)
                │   ├── Queries\              (MediatR query handlers)
                │   └── Validators\           (FluentValidation, may be in Commands)
                └── Domain\                   (Module-specific entities if any)
```

### মূল উদ্দেশ্য

Insurance Engine হল LabAid InsureTech প্ল্যাটফর্মের কেন্দ্রীয় মাইক্রোসার্ভিস যা:

- বিভিন্ন ধরনের insurance product পরিচালনা করে
- Policy lifecycle সম্পূর্ণভাবে manage করে
- Claims processing এবং fraud detection করে
- Underwriting এবং risk assessment পরিচালনা করে

---

## 🏗️ আর্কিটেকচার

### টেকনোলজি স্ট্যাক

```
┌─────────────────────────────────────┐
│     .NET 8 Web API                  │
├─────────────────────────────────────┤
│  Vertical Slice Architecture (VSA)  │
│  • CQRS Pattern (MediatR)           │
│  • Feature-based organization       │
├─────────────────────────────────────┤
│  ডাটা অ্যাক্সেস                      │
│  • Entity Framework Core            │
│  • Repository Pattern               │
├─────────────────────────────────────┤
│  ভ্যালিডেশন ও রেজিলিয়েন্স           │
│  • FluentValidation                 │
│  • Polly (Retry, Circuit Breaker)   │
├─────────────────────────────────────┤
│  ক্যাশিং                             │
│  • StackExchange.Redis              │
├─────────────────────────────────────┤
│  কমিউনিকেশন                          │
│  • gRPC (Proto-first)               │
│  • Kafka (Event-driven)             │
└─────────────────────────────────────┘
```

### প্রধান ডিজাইন প্যাটার্ন

1. **Vertical Slice Architecture:** প্রতিটি feature নিজস্ব folder এ isolated
2. **CQRS:** Command এবং Query আলাদা করা
3. **Domain Events:** Kafka এর মাধ্যমে event publishing
4. **Proto-first:** gRPC service definition প্রথমে

---

## 📦 মডিউল সমূহ

### 1️⃣ Product Catalog Module

**উদ্দেশ্য:** বিভিন্ন insurance product define এবং manage করা

#### দায়িত্ব:

- Product তৈরি, update, deactivate করা
- Product configuration পরিচালনা (coverage, limits, exclusions)
- Pricing rules এবং base premium define করা
- Product eligibility criteria সেট করা

#### Product Types:

- 🚗 **Motor Insurance** (Vehicle)
- 🏥 **Health Insurance** (Medical)
- ✈️ **Travel Insurance**
- 🏠 **Fire Insurance** (Property)
- 💼 **Life Insurance**

#### Key Entities:

```
products
├── product_id (PK)
├── product_type (enum: Motor, Health, Travel, Fire, Life)
├── name
├── description
├── base_premium
├── coverage_details (JSONB)
├── eligibility_rules (JSONB)
├── is_active
└── created_at
```

#### API Operations:

- `CreateProduct` - নতুন product তৈরি
- `UpdateProduct` - product update
- `GetProductById` - product details
- `ListProducts` - সব active products
- `DeactivateProduct` - product বন্ধ করা

---

### 2️⃣ Underwriting Module

**উদ্দেশ্য:** Risk assessment এবং policy approval/rejection decision

#### দায়িত্ব:

- Application review করা
- Risk scoring করা
- Premium calculation করা (base + risk factors)
- Medical/inspection requirements identify করা
- Approval/Rejection/Referral decision নেওয়া

#### Underwriting Process:

```
Application Received
    ↓
Risk Factors Analysis
├── Age, Gender
├── Health History (Health)
├── Vehicle Details (Motor)
├── Travel Destination (Travel)
└── Property Value (Fire)
    ↓
Risk Score Calculation
    ↓
Decision Logic
├── Auto-Approve (Low Risk)
├── Refer to Senior (Medium Risk)
└── Auto-Reject (High Risk)
    ↓
Premium Adjustment
    ↓
Final Decision
```

#### Key Entities:

```
underwriting_applications
├── application_id (PK)
├── product_id (FK)
├── applicant_info (JSONB)
├── risk_score (0-100)
├── premium_calculated
├── status (Pending, Approved, Rejected, Referred)
├── decision_reason
└── decided_at

risk_factors
├── factor_id (PK)
├── product_type
├── factor_name
├── weight
└── scoring_rules (JSONB)
```

#### API Operations:

- `SubmitApplication` - নতুন application জমা
- `AssessRisk` - risk scoring করা
- `CalculatePremium` - premium হিসাব
- `ApproveApplication` - application approve
- `RejectApplication` - application reject
- `GetApplicationStatus` - application status check

---

### 3️⃣ Policy Lifecycle Module

**উদ্দেশ্য:** Active policy পরিচালনা - issue থেকে termination পর্যন্ত

#### দায়িত্ব:

- Policy issuance (underwriting approval পরে)
- Policy information management
- Policy status tracking
- Coverage period management
- Premium collection tracking

#### Policy States:

```
Draft → Active → Renewed
              ↓
         Suspended → Cancelled
              ↓
           Lapsed
```

#### Key Entities:

```
policies
├── policy_id (PK)
├── policy_number (unique)
├── product_id (FK)
├── policyholder_id (FK)
├── premium_amount
├── coverage_amount
├── start_date
├── end_date
├── status (Active, Suspended, Cancelled, Lapsed)
├── payment_frequency (Monthly, Quarterly, Yearly)
└── next_premium_due

insured_assets (Polymorphic Bridge Table)
├── asset_id (PK)
├── policy_id (FK)
├── asset_type (Vehicle, Health, Travel, Property)
└── references specific tables below

vehicle_details (Motor Insurance)
├── vehicle_id (PK)
├── asset_id (FK)
├── make, model, year
├── registration_number
└── chassis_number

health_declarations (Health Insurance)
├── declaration_id (PK)
├── asset_id (FK)
├── pre_existing_conditions (JSONB)
└── medical_history (JSONB)

travel_details (Travel Insurance)
├── travel_id (PK)
├── asset_id (FK)
├── destination_countries
├── trip_start_date
└── trip_end_date

property_details (Fire Insurance)
├── property_id (PK)
├── asset_id (FK)
├── property_type
├── address
└── construction_type
```

#### API Operations:

- `IssuePolicy` - নতুন policy issue
- `GetPolicyDetails` - policy information
- `UpdatePolicyholderInfo` - policyholder তথ্য update
- `SuspendPolicy` - temporary suspend
- `ActivatePolicy` - suspended থেকে activate
- `CancelPolicy` - policy cancel

---

### 4️⃣ Renewal Module

**উদ্দেশ্য:** Policy renewal পরিচালনা

#### দায়িত্ব:

- Renewal notice generation
- Renewal premium calculation (claim history বিবেচনা করে)
- Renewal offer তৈরি
- Renewal processing
- Lapsed policy re-activation

#### Renewal Process:

```
60 Days Before Expiry
    ↓
Renewal Notice Sent
    ↓
Calculate New Premium
├── Base Premium
├── Claim History Impact
├── Age/Risk Factor Changes
└── Market Adjustments
    ↓
Renewal Offer
    ↓
Customer Decision
├── Accept → Process Renewal
├── Modify → Re-quote
└── Decline → Mark as Lapsed
```

#### Key Entities:

```
renewals
├── renewal_id (PK)
├── policy_id (FK)
├── old_premium
├── new_premium
├── renewal_date
├── status (Pending, Accepted, Declined, Processed)
└── offer_valid_until
```

#### API Operations:

- `GenerateRenewalOffer` - renewal offer তৈরি
- `ProcessRenewal` - renewal করা
- `DeclineRenewal` - renewal না করা
- `GetRenewalStatus` - renewal status

---

### 5️⃣ Endorsement Module

**উদ্দেশ্য:** Active policy তে changes করা

#### দায়িত্ব:

- Coverage increase/decrease
- Policyholder information update
- Insured asset changes
- Premium adjustment
- Policy term modification

#### Endorsement Types:

- **Coverage Change:** সীমা বৃদ্ধি/হ্রাস
- **Insured Asset Change:** গাড়ি পরিবর্তন, ঠিকানা পরিবর্তন
- **Beneficiary Update:** nominee পরিবর্তন
- **Payment Frequency Change:** monthly → yearly

#### Key Entities:

```
endorsements
├── endorsement_id (PK)
├── policy_id (FK)
├── endorsement_type
├── change_details (JSONB)
├── premium_adjustment
├── effective_date
├── status (Pending, Approved, Rejected)
└── approved_by
```

#### API Operations:

- `RequestEndorsement` - endorsement request
- `ApproveEndorsement` - approve করা
- `RejectEndorsement` - reject করা
- `GetEndorsementHistory` - history দেখা

---

### 6️⃣ Cancellation Module

**উদ্দেশ্য:** Policy বাতিল করা এবং refund পরিচালনা

#### দায়িত্ব:

- Cancellation request processing
- Refund calculation (pro-rata basis)
- Cancellation reason tracking
- Regulatory compliance

#### Cancellation Types:

- **Customer Request:** গ্রাহক নিজে cancel করতে চায়
- **Non-payment:** premium পরিশোধ না করা
- **Fraud Detection:** জালিয়াতি পাওয়া গেলে
- **Underwriting Decision:** post-issuance review পরে

#### Refund Calculation:

```
Refund = (Unused Premium) - (Cancellation Fee) - (Claims Paid)

Where:
Unused Premium = (Remaining Days / Total Days) × Annual Premium
```

#### Key Entities:

```
cancellations
├── cancellation_id (PK)
├── policy_id (FK)
├── cancellation_reason
├── requested_by
├── cancellation_date
├── refund_amount
├── refund_status (Pending, Processed, Rejected)
└── processed_at
```

#### API Operations:

- `RequestCancellation` - cancellation request
- `CalculateRefund` - refund amount হিসাব
- `ProcessCancellation` - cancel করা
- `GetCancellationDetails` - details দেখা

---

### 7️⃣ Claims Management Module

**উদ্দেশ্য:** Insurance claim পরিচালনা - filing থেকে settlement পর্যন্ত

#### দায়িত্ব:

- Claim registration
- Document collection
- Claim investigation
- Loss assessment
- Approval workflow (L1 → Board)
- Settlement processing

#### Claims Process:

```
Claim Filed
    ↓
Document Verification
    ↓
Initial Assessment (L1)
├── Minor Claims → Approve
├── Medium Claims → L2 Review
└── Major Claims → Senior/Board
    ↓
Investigation (if needed)
    ↓
Loss Assessment
    ↓
Approval Decision
    ↓
Settlement
```

#### Claim States:

```
Submitted → UnderReview → Investigating
                              ↓
                      Approved/Rejected
                              ↓
                         Settled/Closed
```

#### Key Entities:

```
claims
├── claim_id (PK)
├── policy_id (FK)
├── claim_number (unique)
├── incident_date
├── incident_description
├── claim_amount_requested
├── claim_amount_approved
├── status (Submitted, UnderReview, Approved, Rejected, Settled)
├── assigned_to (adjuster)
└── settlement_date

claim_documents
├── document_id (PK)
├── claim_id (FK)
├── document_type (Police Report, Medical Bill, Photos, etc.)
├── file_path
└── uploaded_at

claim_assessments
├── assessment_id (PK)
├── claim_id (FK)
├── assessor_id
├── assessment_notes
├── recommended_amount
└── assessed_at
```

#### API Operations:

- `FileClaim` - নতুন claim দাখিল
- `UploadDocument` - documents submit
- `AssignClaim` - adjuster assign করা
- `AssessClaim` - loss assessment
- `ApproveClaim` - claim approve
- `RejectClaim` - claim reject
- `SettleClaim` - payment process
- `GetClaimStatus` - status tracking

---

### 8️⃣ Fraud Detection Module

**উদ্দেশ্য:** জালিয়াতি সনাক্ত করা এবং প্রতিরোধ করা

#### দায়িত্ব:

- Suspicious pattern detection
- Risk scoring for claims
- Duplicate claim detection
- Historical fraud pattern matching
- Alert generation

#### Fraud Indicators:

- 🚨 **Multiple claims** একই সময়ে একই type এর
- 🚨 **Recent policy issuance** (within 30 days) এবং major claim
- 🚨 **Inconsistent information** application vs claim এ
- 🚨 **Previous fraud history** policyholder এর
- 🚨 **Claim amount** অস্বাভাবিক রকম বেশি

#### Fraud Risk Scoring:

```
Risk Score = Σ (Factor Weight × Factor Value)

Factors:
- Recent Policy: 25 points
- Multiple Claims: 30 points
- High Amount: 20 points
- Missing Documents: 15 points
- Inconsistency: 25 points

Risk Levels:
0-30: Low Risk (Auto-process)
31-60: Medium Risk (Manual Review)
61-100: High Risk (Investigation Required)
```

#### Key Entities:

```
fraud_alerts
├── alert_id (PK)
├── claim_id (FK)
├── policy_id (FK)
├── risk_score (0-100)
├── fraud_indicators (JSONB)
├── alert_status (New, UnderReview, Confirmed, FalsePositive)
├── investigated_by
└── resolution_notes

fraud_rules
├── rule_id (PK)
├── rule_name
├── rule_description
├── detection_logic (JSONB)
├── risk_weight
└── is_active
```

#### API Operations:

- `AnalyzeClaim` - claim fraud check
- `GetFraudAlerts` - pending alerts
- `InvestigateFraud` - investigation করা
- `ConfirmFraud` - fraud confirm
- `MarkFalsePositive` - false alarm mark করা

---

## 🗄️ ডাটাবেস স্কিমা

### Core Tables Overview

```sql
-- Products
products (product_id, product_type, name, base_premium, is_active)

-- Policyholders
policyholders (policyholder_id, name, email, phone, address)

-- Policies
policies (policy_id, policy_number, product_id, policyholder_id,
         premium_amount, coverage_amount, start_date, end_date, status)

-- Polymorphic Asset Pattern
insured_assets (asset_id, policy_id, asset_type)
  ├── vehicle_details (vehicle_id, asset_id, make, model, year)
  ├── health_declarations (declaration_id, asset_id, conditions)
  ├── travel_details (travel_id, asset_id, destinations, dates)
  └── property_details (property_id, asset_id, address, type)

-- Underwriting
underwriting_applications (application_id, product_id, risk_score,
                          status, premium_calculated)

-- Claims
claims (claim_id, policy_id, claim_number, incident_date,
       claim_amount_requested, status)

-- Renewals
renewals (renewal_id, policy_id, old_premium, new_premium, status)

-- Endorsements
endorsements (endorsement_id, policy_id, endorsement_type,
             change_details, status)

-- Cancellations
cancellations (cancellation_id, policy_id, cancellation_reason,
              refund_amount, refund_status)

-- Fraud Detection
fraud_alerts (alert_id, claim_id, risk_score, fraud_indicators, status)
```

### Polymorphic Asset Pattern সুবিধা:

✅ একাধিক product type একই structure এ
✅ Type-specific details আলাদা table এ
✅ Extensible - নতুন product type সহজে যোগ করা যায়
✅ Query optimization - শুধু প্রয়োজনীয় table join

---

## 🔗 সার্ভিস ইন্টিগ্রেশন

### gRPC Services (Proto-first)

```protobuf
// প্রতিটি module এর জন্য আলাদা service
service ProductService { }
service UnderwritingService { }
service PolicyService { }
service ClaimsService { }
service FraudDetectionService { }

// Proto থেকে C# classes generate হয়
// Location: root/gen/csharp/
```

### Kafka Event Topics

```
insurance.policy.issued          - নতুন policy issue হলে
insurance.policy.renewed         - renewal হলে
insurance.policy.cancelled       - cancellation হলে
insurance.policy.endorsed        - endorsement হলে
insurance.claim.filed            - নতুন claim submit হলে
insurance.claim.approved         - claim approve হলে
insurance.claim.settled          - settlement complete হলে
insurance.fraud.detected         - fraud alert হলে
insurance.premium.due            - premium due reminder
```

### External Service Dependencies

- **Payment Service:** Premium collection, refund processing
- **Document Service:** Policy documents, claim documents storage
- **Notification Service:** Email, SMS alerts
- **Customer Service:** Customer information management
- **Audit Service:** Audit trail logging

---

## ⚙️ ওয়ার্কফ্লো এবং স্টেট মেশিন

### Policy State Machine

```
┌─────────┐
│  Draft  │ (Application under review)
└────┬────┘
     ↓ (Underwriting approved)
┌─────────┐
│ Active  │ (Policy in force)
└────┬────┘
     ├─→ Renewed (Renewal processed)
     ├─→ Suspended (Non-payment, temporary)
     ├─→ Cancelled (Permanent termination)
     └─→ Lapsed (Expired, not renewed)
```

### Claim State Machine

```
┌───────────┐
│ Submitted │ (Initial filing)
└─────┬─────┘
      ↓
┌─────────────┐
│ UnderReview │ (Document verification)
└─────┬───────┘
      ↓
┌──────────────┐
│Investigating │ (Detailed assessment)
└──────┬───────┘
       ├─→ Approved → Settled
       └─→ Rejected → Closed
```

---

## 👥 অনুমোদন ম্যাট্রিক্স (Claims Approval Matrix)

### Level-based Approval

```
┌─────────┬──────────────────┬─────────────────────┐
│ Level   │ Claim Amount     │ Approver            │
├─────────┼──────────────────┼─────────────────────┤
│ L1      │ ≤ 50,000 BDT     │ Claims Officer      │
│ L2      │ 50,001-200,000   │ Senior Claims Mgr   │
│ L3      │ 200,001-500,000  │ Head of Claims      │
│ Board   │ > 500,000 BDT    │ Board of Directors  │
└─────────┴──────────────────┴─────────────────────┘
```

### Approval Workflow

```
Claim Amount Check
    ↓
Route to Appropriate Level
    ↓
Level Review & Decision
    ↓
(If rejected) → Closed
(If approved & amount > level limit) → Escalate
(If approved & within limit) → Settle
```

---

## 🔧 প্রযুক্তিগত বৈশিষ্ট্য

### Resilience Patterns (Polly)

- **Retry Policy:** Failed API calls পুনরায় চেষ্টা
- **Circuit Breaker:** Failing service isolate করা
- **Timeout Policy:** Long-running operations control

### Caching Strategy (Redis)

- Product catalog caching
- Premium calculation rules caching
- Fraud detection rules caching
- Session data storage

### Validation (FluentValidation)

- Request validation middleware
- Business rule validation
- Cross-field validation

### API Documentation

- Generated from Proto definitions
- HTML docs location: `doc/`
- Swagger/OpenAPI support

---

## 📊 Module Dependency Graph

```
┌─────────────────┐
│Product Catalog  │ (Foundation layer)
└────────┬────────┘
         ↓
┌────────────────┐
│ Underwriting   │ (Uses product rules)
└────────┬───────┘
         ↓
┌────────────────┐
│Policy Lifecycle│ (Uses underwriting decisions)
└────┬───────────┘
     ├─→ Renewal Module
     ├─→ Endorsement Module
     └─→ Cancellation Module

┌────────────────┐
│Claims Mgmt     │ (Linked to active policies)
└────────┬───────┘
         ↓
┌────────────────┐
│Fraud Detection │ (Monitors claims)
└────────────────┘
```

---

## 🎯 মূল বিজনেস রুলস

### Policy Issuance Rules

- Underwriting approval থাকতে হবে
- Premium payment confirmation লাগবে
- সব mandatory documents থাকতে হবে
- Policyholder verification complete হতে হবে

### Claims Rules

- Policy active অবস্থায় থাকতে হবে
- Incident date policy period এর মধ্যে হতে হবে
- Waiting period (যদি থাকে) শেষ হতে হবে
- Claim amount coverage limit এর মধ্যে হতে হবে

### Renewal Rules

- Policy expiry এর 60 দিন আগে থেকে renewal শুরু
- Claim history premium এ প্রভাব ফেলবে
- Lapsed policy 90 দিন পর্যন্ত renew করা যাবে

---

## 📁 সোর্স কোড স্ট্রাকচার

```
D:\InsureTech\
├── proto\                                          (Root-level shared proto repository)
│   └── insurance_engine\
│       ├── products.proto
│       ├── policies.proto
│       ├── claims.proto
│       ├── underwriting.proto
│       └── [other modules].proto
│
├── gen\                                            (Root-level generated code)
│   └── csharp\
│       └── InsuranceEngine\
│           ├── Products\
│           ├── Policies\
│           ├── Claims\
│           └── [other modules]\
│
└── backend\
    └── insurance_engine\
        └── src\
            ├── InsuranceEngine.SharedKernel\       (Shared infrastructure)
            │   └── Persistence\
            │       ├── Entities\                   (Core shared entities)
            │       │   ├── Product.cs
            │       │   ├── Policy.cs
            │       │   ├── Claim.cs
            │       │   └── [other entities]
            │       └── Migrations\                 (EF Core migrations)
            │           └── *.cs
            │
            └── InsuranceEngine.[Module]\           (Module-specific projects)
                ├── InsuranceEngine.Products\
                ├── InsuranceEngine.Policies\
                ├── InsuranceEngine.Claims\
                ├── InsuranceEngine.Underwriting\
                ├── InsuranceEngine.Renewals\
                ├── InsuranceEngine.Endorsements\
                ├── InsuranceEngine.Cancellations\
                └── InsuranceEngine.FraudDetection\

প্রতিটি Module Project এর স্ট্রাকচার:
InsuranceEngine.[Module]\
├── GrpcServices\                           (gRPC service implementations)
│   └── [Module]ServiceImpl.cs             (Proto-generated base class inherit করে)
├── Application\
│   ├── Commands\                          (MediatR command handlers)
│   │   ├── Create[Entity]\
│   │   │   ├── Create[Entity]Command.cs
│   │   │   ├── Create[Entity]Handler.cs
│   │   │   └── Create[Entity]Validator.cs (Or in Validators folder)
│   │   └── [other commands]
│   ├── Queries\                           (MediatR query handlers)
│   │   ├── Get[Entity]ById\
│   │   └── List[Entities]\
│   └── Validators\                        (If separate from Commands)
└── Domain\                                (Module-specific entities if needed)
```

---

## 📚 ডকুমেন্টেশন রেফারেন্স

- **SRS:** `doccumentation/SRS_V3/LabAid_InsureTech_SRS_v3.11.md`
- **API Docs:** `doc/` (HTML generated from Proto)
- **Spec File:** `insurance-engine-spec.md` (DDL, Proto, Kafka, State Machines)
- **Proto Generated Classes:** `root/gen/csharp/`

---

## 🔍 Audit Categories

1. Proto Compliance
2. API Contract Adherence
3. SRS Business Rule Compliance
4. CRUD Operation Completeness
5. Architecture Layering (VSA/CQRS)
6. Validation Rules
7. Error Handling
8. Event Publishing
9. Database Schema Consistency

---

## 🚀 ডেভেলপমেন্ট স্ট্যাটাস

**Current Phase:** Implementation Complete
**Next Steps:** Audit and verification using Antigravity agent

**Audit Workflow:**

1. Proto-generated C# classes verification
2. HTML API docs compliance check
3. SRS v3.11 business rules validation
4. Cross-module integration testing

---

_এই ডকুমেন্টেশন LabAid InsureTech Insurance Engine এর সম্পূর্ণ technical এবং functional overview প্রদান করে।_
