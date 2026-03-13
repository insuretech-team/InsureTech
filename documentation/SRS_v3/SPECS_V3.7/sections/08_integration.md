# 8. Integration Requirements

Requirements

### 8.1 External System Integrations

### 8.1.1 Payment Gateway Integrations

| Integration | Type | Protocol | Purpose | Priority |
|-------------|------|----------|---------|----------|
| bKash API | Payment Gateway | REST/Webhook | Premium payments, claim settlements | M1 |
| Nagad API | Payment Gateway | REST/Webhook | Premium payments, claim settlements | M2 |
| Bangladesh Bank | Payment Gateway | ISO 8583 | Bank transfers, regulatory reporting | M2 |
| NID Verification API | Government Service | REST/SOAP | Identity verification | M1 |
| Mobile Number Verification | Government Service | REST | Phone number validation | M1 |
| SMS Gateway | Communication | REST | Notifications and OTP delivery | M2 |
| Email Service | Communication | SMTP/API | Email notifications and documents | M1 |
| Hospital EHR Systems | Healthcare | HL7 FHIR | Medical record integration | D |
| Weather API | Risk Assessment | REST | Environmental risk monitoring | M3 |
| WhatsApp Business | Communication | API | Customer service and notifications | M3 |

---

### 8.1.2 Integration Details **[EXPANDED IN V3.7]**

#### **bKash Payment Integration**
- **API Endpoints:** Payment initiation, verification, refund, account inquiry
- **Authentication:** OAuth 2.0 with client credentials flow
- **Rate Limits:** 100 requests/minute per API key
- **Transaction Timeout:** 5 seconds connection, 15 seconds read
- **Webhook Support:** Payment confirmation callback to `/webhooks/bkash`
- **Error Handling:** Retry logic with exponential backoff (3 attempts: 1s, 3s, 9s)
- **Cost:** Transaction fee 1.5% (to be negotiated with merchant agreement)
- **Testing:** Sandbox environment available at `https://tokenized.sandbox.bka.sh`
- **Mock Service:** JSON-based mock required for local development
- **Fallback:** Queue payment for manual processing if API down >5 minutes

#### **Nagad & Rocket Integration**
- Similar specifications as bKash with provider-specific endpoints
- **Fallback Mechanism:** If primary MFS fails, auto-retry with alternate MFS within 30 seconds
- **Health Check:** Monitor MFS availability every 60 seconds
- **Load Balancing:** Distribute payment load across available providers

#### **NID Verification API**
- **Provider:** Bangladesh Election Commission API / Third-party aggregator (PORICHOY - TBD)
- **API Endpoints:** `/verify-nid`, `/get-nid-details`, `/face-match`
- **Authentication:** API key + IP whitelist + TLS 1.3
- **Rate Limits:** 1000 verifications/day (Basic tier), 10,000/day (Premium tier)
- **Response Time:** <3 seconds average, 10 seconds timeout
- **Data Returned:** Name (English/Bengali), Father/Mother name, DOB, Address, Photo (Base64)
- **Cost:** 5-10 BDT per verification (volume-based pricing)
- **Fallback:** Manual verification queue when API unavailable (notify admin within 2 minutes)
- **Mock Service:** JSON-based mock with 100 sample NID records for testing
- **Compliance:** Store verification logs for 20 years as per IDRA requirement
- **Data Privacy:** NID data encrypted at rest, access logged

#### **Hospital EHR (HL7 FHIR) Integration**
- **Standard:** HL7 FHIR R4 (Fast Healthcare Interoperability Resources)
- **FHIR Resources Used:** 
  - Patient (demographics)
  - Encounter (admission/discharge)
  - Condition (diagnoses)
  - Procedure (treatments)
  - MedicationRequest (prescriptions)
  - DiagnosticReport (lab results, imaging)
- **Authentication:** OAuth 2.0 + JWT tokens, refresh every 30 minutes
- **Endpoints:** 
  - `GET /Patient/{id}` - Patient lookup
  - `POST /Claim` - Submit cashless claim
  - `GET /Claim/{id}` - Check claim status
- **Connection Timeout:** 5 seconds
- **Read Timeout:** 15 seconds
- **Fallback Behavior:** 
  - If timeout → Queue for manual verification
  - Send SMS to hospital focal point: "System unavailable, manual claim processing required"
  - Notify InsureTech support team via Slack/email
- **Data Mapping:** Custom FHIR profile for Bangladesh healthcare context
- **Consent Management:** Patient consent workflow (digital signature) before EHR access
- **Participating Hospitals:** 
  - **Phase M1:** LabAid Hospitals (5 locations)
  - **Phase M2+:** Expand to 20+ partner hospitals
- **Testing:** FHIR test server with synthetic patient data

#### **SMS Gateway Integration**
- **Provider:** Twilio (international) / Bangladesh SMS Provider (local - TBD)
- **Use Cases:** OTP delivery, payment confirmation, claim status updates, policy renewal reminders
- **Delivery Rate Target:** >95%
- **Delivery Status Tracking:** Webhook for delivery confirmation
- **Cost Optimization:** Batching, template caching, priority queuing
- **Rate Limits:** 100 messages/second

#### **WhatsApp Business API**
- **Message Templates:** Pre-approved templates for policy confirmation, claim updates
- **Opt-in/Opt-out Workflow:** User consent required, unsubscribe link in messages
- **Rate Limits:** 1000 conversations/day (business tier)
- **Cost:** Per-conversation pricing (varies by country)
- **Use Cases:** Rich media policy documents, claim photo uploads, customer support chat

---


### 8.2 Internal Service Communications **[ENHANCED IN V3.7]**

All microservices communicate via **gRPC** using Protocol Buffers for type safety, performance, and cross-language compatibility:

**Benefits of gRPC:**
- ✅ Strong typing with Protocol Buffers
- ✅ HTTP/2 multiplexing (multiple requests over single connection)
- ✅ Bi-directional streaming support
- ✅ Built-in load balancing and health checking
- ✅ Language-agnostic (Go, C#, Python, Node.js)
- ✅ 7-10x faster than REST JSON for internal communication

**Service Discovery:** Consul for service registry and health checks
**Load Balancing:** Client-side load balancing with round-robin strategy
**Timeout Configuration:** 30 seconds for standard RPCs, 5 minutes for long-running operations

```protobuf
// Insurance Engine Service
service InsuranceEngineService {
  rpc IssuePolicy(IssuePolicyRequest) returns (IssuePolicyResponse);
  rpc CalculatePremium(CalculatePremiumRequest) returns (CalculatePremiumResponse);
  rpc ProcessRenewal(ProcessRenewalRequest) returns (ProcessRenewalResponse);
  rpc SubmitClaim(SubmitClaimRequest) returns (SubmitClaimResponse);
}

// Payment Service
service PaymentService {
  rpc ProcessPayment(ProcessPaymentRequest) returns (ProcessPaymentResponse);
  rpc RefundPayment(RefundPaymentRequest) returns (RefundPaymentResponse);
  rpc GetPaymentStatus(GetPaymentStatusRequest) returns (GetPaymentStatusResponse);
}

// Partner Management Service
service PartnerManagementService {
  rpc OnboardPartner(OnboardPartnerRequest) returns (OnboardPartnerResponse);
  rpc RegisterAgent(RegisterAgentRequest) returns (RegisterAgentResponse);
  rpc CalculateCommission(CalculateCommissionRequest) returns (CalculateCommissionResponse);
}
```

### 8.3 Event-Driven Architecture

**Kafka Event Streaming:**
- **Policy Events:** PolicyIssued, PolicyRenewed, PolicyCancelled, PremiumPaid
- **Claim Events:** ClaimSubmitted, ClaimApproved, ClaimSettled, ClaimRejected
- **Payment Events:** PaymentProcessed, PaymentFailed, RefundIssued
- **User Events:** UserRegistered, KYCCompleted, ProfileUpdated

**Event Processing Patterns:**
- **Event Sourcing:** Critical business events for audit and replay
- **CQRS:** Separate command and query models for complex aggregates
- **Saga Pattern:** Distributed transaction management across services
- **Event Streaming:** Real-time analytics and monitoring

[[[PAGEBREAK]]]
