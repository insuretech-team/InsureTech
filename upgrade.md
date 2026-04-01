# InsureTech Platform — Upgrade Roadmap

Proto-native upgrade paths derived from analyzing the 14 reference projects in `.resources/`, mapped to InsureTech's existing architecture.

---

## Feature Gap Analysis

| Reference Project | What It Offers | Existing InsureTech Proto Domain | Gap Level |
|---|---|---|---|
| `AI-Insurance-Policy-Engine` | RAG-based policy Q&A with vector search | `proto/insuretech/ai/` | 🟡 Partial — AI domain exists, needs RAG service |
| `SmartClaim-RAG-Engine` | Clause-backed claim decisions with FAISS | `proto/insuretech/claims/` | 🟡 Partial — Claims exist, needs AI layer |
| `PolicyPal-Insurance-Reasoning-Engine` | LLM reasoning over policy documents | `proto/insuretech/policy/` | 🔴 Missing — No reasoning engine |
| `Gemini-PolicyPal-Insurance-Reasoning-Engine` | Multi-model policy reasoning (Gemini) | `proto/insuretech/ai/` | 🔴 Missing — No multi-model support |
| `RulesEngine` | Configurable business rules (Microsoft) | `proto/insuretech/insurance/` | 🟡 Partial — Rules in `rules/` dir, not as a service |
| `InsuranceFraudDetection` | ML fraud detection pipeline | `proto/insuretech/fraud/` | 🟡 Partial — Fraud domain exists, needs ML models |
| `Insurance-Quoting-Tool` | Dynamic premium quoting | `proto/insuretech/products/` | 🟡 Partial — Products exist, needs quoting engine |
| `InsuranceBuddy-ChatBot` | NLU chatbot (DialogFlow) | `proto/insuretech/support/` | 🔴 Missing — No chatbot service |
| `InsuranceProCRM` | Lead/client management dashboard | `proto/insuretech/partner/` + `b2b/` | 🟡 Partial — Partner domain exists, lacks CRM features |
| `VehicleInsuranceAPI` | Vehicle-specific insurance logic | `proto/insuretech/products/` | 🟡 Partial — No vehicle-specific types |
| `CalculusPrime-PnC` | P&C actuarial calculations | `proto/insuretech/underwriting/` | 🟡 Partial — Underwriting exists, needs actuarial engine |
| `MaxLife_Insurance_project` | Life insurance product management | `proto/insuretech/products/` | 🟡 Partial — No life-specific product types |
| `insurence_query_engine` | Insurance query/search engine | `proto/insuretech/ai/` | 🔴 Missing — No semantic search service |
| `life-insurance-pricing-estimator` | Life insurance pricing models | `proto/insuretech/underwriting/` | 🔴 Missing — No pricing estimator |

---

## Prioritized Upgrade Roadmap

### Phase 1 — High Impact, Leverages Existing Domains

#### 1.1 Rules Engine Service
**Source:** `RulesEngine` (Microsoft)
**Target domain:** `proto/insuretech/workflow/` + new `proto/insuretech/rules/`
**Business value:** 🔴 Critical — externalizes business rules from code

```protobuf
// proto/insuretech/rules/v1/rules.proto
syntax = "proto3";
package insuretech.rules.v1;

service RulesService {
  rpc EvaluateRules(EvaluateRulesRequest) returns (EvaluateRulesResponse);
  rpc CreateWorkflow(CreateWorkflowRequest) returns (CreateWorkflowResponse);
  rpc ListWorkflows(ListWorkflowsRequest) returns (ListWorkflowsResponse);
  rpc GetWorkflow(GetWorkflowRequest) returns (GetWorkflowResponse);
  rpc UpdateWorkflow(UpdateWorkflowRequest) returns (UpdateWorkflowResponse);
  rpc DeleteWorkflow(DeleteWorkflowRequest) returns (DeleteWorkflowResponse);
}

message Rule {
  string id = 1;
  string rule_name = 2;
  string expression = 3; // Lambda expression
  RuleExpressionType expression_type = 4;
  string success_event = 5;
  string error_message = 6;
}

enum RuleExpressionType {
  RULE_EXPRESSION_TYPE_UNSPECIFIED = 0;
  RULE_EXPRESSION_TYPE_LAMBDA = 1;
  RULE_EXPRESSION_TYPE_CUSTOM = 2;
}

message Workflow {
  string id = 1;
  string workflow_name = 2;
  repeated Rule rules = 3;
  google.protobuf.Timestamp created_at = 4;
  google.protobuf.Timestamp updated_at = 5;
}
```

**Implementation:** Go service wrapping a rules evaluation engine, configurable via the system portal. Could drive underwriting decisions, claims automation, and policy eligibility checks.

---

#### 1.2 Fraud Detection ML Service
**Source:** `InsuranceFraudDetection`
**Target domain:** `proto/insuretech/fraud/`
**Business value:** 🔴 Critical — reduces fraudulent claims payouts

```protobuf
// proto/insuretech/fraud/v1/detection.proto (extend existing)

service FraudDetectionService {
  rpc AnalyzeClaim(AnalyzeClaimRequest) returns (AnalyzeClaimResponse);
  rpc TrainModel(TrainModelRequest) returns (TrainModelResponse);
  rpc GetModelStatus(GetModelStatusRequest) returns (GetModelStatusResponse);
  rpc ListFraudIndicators(ListFraudIndicatorsRequest) returns (ListFraudIndicatorsResponse);
}

message AnalyzeClaimResponse {
  string claim_id = 1;
  double fraud_probability = 2;
  FraudRiskLevel risk_level = 3;
  repeated FraudIndicator indicators = 4;
  string recommendation = 5;
}

enum FraudRiskLevel {
  FRAUD_RISK_LEVEL_UNSPECIFIED = 0;
  FRAUD_RISK_LEVEL_LOW = 1;
  FRAUD_RISK_LEVEL_MEDIUM = 2;
  FRAUD_RISK_LEVEL_HIGH = 3;
  FRAUD_RISK_LEVEL_CRITICAL = 4;
}
```

**Implementation:** Python ML microservice with gRPC interface, called by `insurance_engine` during claims processing. Model trained on historical claims data.

---

#### 1.3 Dynamic Quoting Engine
**Source:** `Insurance-Quoting-Tool`, `life-insurance-pricing-estimator`
**Target domain:** `proto/insuretech/products/`
**Business value:** 🟠 High — enables real-time premium calculation

```protobuf
// proto/insuretech/products/v1/quoting.proto

service QuotingService {
  rpc GenerateQuote(GenerateQuoteRequest) returns (GenerateQuoteResponse);
  rpc CompareQuotes(CompareQuotesRequest) returns (CompareQuotesResponse);
  rpc GetQuote(GetQuoteRequest) returns (GetQuoteResponse);
  rpc ListQuotes(ListQuotesRequest) returns (ListQuotesResponse);
}

message GenerateQuoteRequest {
  string product_id = 1;
  QuoteParameters parameters = 2;
  string customer_id = 3;
}

message QuoteParameters {
  CoverageType coverage_type = 1;
  string coverage_plan = 2;
  double asset_value = 3;
  repeated string optional_coverages = 4;
  google.protobuf.Duration coverage_duration = 5;
}

message GenerateQuoteResponse {
  string quote_id = 1;
  double premium_amount = 2;
  repeated PremiumBreakdown breakdown = 3;
  google.protobuf.Timestamp valid_until = 4;
}
```

**Implementation:** Extend `insurance_engine` with a quoting module. Calculation logic driven by the Rules Engine (1.1). Exposed to `customer_portal` and `b2b_portal`.

---

### Phase 2 — AI-Powered Features

#### 2.1 AI Policy Analysis (RAG)
**Source:** `AI-Insurance-Policy-Engine`, `SmartClaim-RAG-Engine`, `PolicyPal-Insurance-Reasoning-Engine`
**Target domain:** `proto/insuretech/ai/`
**Business value:** 🟠 High — enables natural language policy understanding

```protobuf
// proto/insuretech/ai/v1/policy_analysis.proto

service PolicyAnalysisService {
  rpc AnalyzePolicy(AnalyzePolicyRequest) returns (AnalyzePolicyResponse);
  rpc AskQuestion(AskQuestionRequest) returns (AskQuestionResponse);
  rpc IndexDocument(IndexDocumentRequest) returns (IndexDocumentResponse);
  rpc ListIndexedDocuments(ListIndexedDocumentsRequest) returns (ListIndexedDocumentsResponse);
}

message AskQuestionRequest {
  string policy_id = 1;
  string question = 2;
  int32 max_clauses = 3; // number of relevant clauses to retrieve
}

message AskQuestionResponse {
  string answer = 1;
  AnswerType answer_type = 2; // COVERAGE, LIMITS, EXCLUSION, etc.
  repeated ClauseReference source_clauses = 3;
  double confidence = 4;
}

message ClauseReference {
  string clause_id = 1;
  string clause_text = 2;
  string section = 3;
  int32 page_number = 4;
  double relevance_score = 5;
}

enum AnswerType {
  ANSWER_TYPE_UNSPECIFIED = 0;
  ANSWER_TYPE_COVERAGE = 1;
  ANSWER_TYPE_LIMITS = 2;
  ANSWER_TYPE_EXCLUSION = 3;
  ANSWER_TYPE_CONDITIONS = 4;
  ANSWER_TYPE_FINANCIAL = 5;
  ANSWER_TYPE_REQUIREMENTS = 6;
}
```

**Implementation:** Separate AI microservice using vector search (FAISS/ChromaDB) + LLM generation. Ingests policy PDFs from the `document` domain, stores embeddings in a vector store, and serves answers via gRPC to all portals.

---

#### 2.2 Insurance Chatbot
**Source:** `InsuranceBuddy-ChatBot`
**Target domain:** `proto/insuretech/support/`
**Business value:** 🟡 Medium — automates customer support

```protobuf
// proto/insuretech/support/v1/chatbot.proto

service ChatbotService {
  rpc SendMessage(SendMessageRequest) returns (SendMessageResponse);
  rpc GetConversation(GetConversationRequest) returns (GetConversationResponse);
  rpc ListConversations(ListConversationsRequest) returns (ListConversationsResponse);
}

message SendMessageRequest {
  string conversation_id = 1;
  string message = 2;
  string customer_id = 3;
}

message SendMessageResponse {
  string response = 1;
  MessageIntent intent = 2;
  repeated SuggestedAction actions = 3;
}

enum MessageIntent {
  MESSAGE_INTENT_UNSPECIFIED = 0;
  MESSAGE_INTENT_COVERAGE_QUERY = 1;
  MESSAGE_INTENT_CLAIM_STATUS = 2;
  MESSAGE_INTENT_QUOTE_REQUEST = 3;
  MESSAGE_INTENT_SUPPORT_TICKET = 4;
  MESSAGE_INTENT_GENERAL = 5;
}
```

**Implementation:** Wraps the AI Policy Analysis service (2.1) with a conversational interface. Intent classification routes to the appropriate backend service (claims, quoting, support). Deployed on `customer_portal`.

---

### Phase 3 — Domain-Specific Enhancements

#### 3.1 Vehicle Insurance Module
**Source:** `VehicleInsuranceAPI`
**Target domain:** `proto/insuretech/products/`
**Business value:** 🟡 Medium — enables vehicle-specific product lines

New message types for vehicle-specific fields (VIN, make/model/year, coverage tiers, deductible schedules). Extend the `QuotingService` with vehicle-aware premium calculations.

#### 3.2 Life Insurance Pricing
**Source:** `MaxLife_Insurance_project`, `life-insurance-pricing-estimator`
**Target domain:** `proto/insuretech/products/` + `underwriting/`
**Business value:** 🟡 Medium — adds actuarial pricing models

New proto messages for life insurance parameters (age, health class, mortality tables). Pricing models implemented as a separate calculation service called by the quoting engine.

#### 3.3 P&C Actuarial Engine
**Source:** `CalculusPrime-PnC`
**Target domain:** `proto/insuretech/underwriting/`
**Business value:** 🟡 Medium — enables property & casualty actuarial analysis

Integrate actuarial calculation models into the underwriting service for loss ratio analysis, reserve calculations, and risk scoring.

#### 3.4 CRM Enhancement
**Source:** `InsuranceProCRM`
**Target domain:** `proto/insuretech/partner/` + `b2b/`
**Business value:** 🟡 Medium — improves agent/partner experience

Add lead management, pipeline tracking, and client relationship features to the partner and B2B portal. Extend existing partner proto with CRM-specific messages.

---

## Integration Patterns

All upgrades follow the same proto-first pattern:

```
1. Define proto schema  →  proto/insuretech/<domain>/v1/<feature>.proto
2. Run buf generate     →  gen/go/, gen/ts/, gen/csharp/
3. Implement Go service →  backend/<service>/internal/handler/
4. Wire in gRPC server  →  backend/<service>/cmd/
5. Add REST mapping     →  api/openapi.yaml + api/paths/
6. Update Docker        →  docker-compose.yml
7. Frontend integration →  <portal>/src/lib/api/
```

### Cross-Service Communication

```
                    ┌─────────────────┐
                    │   API Gateway   │
                    └───────┬─────────┘
                            │
          ┌─────────┬───────┼───────┬──────────┐
          │         │       │       │          │
     ┌────▼──┐ ┌───▼──┐ ┌──▼───┐ ┌▼──────┐ ┌─▼────┐
     │Policy │ │Claims│ │Quote │ │Fraud  │ │AI    │
     │Sync   │ │      │ │Engine│ │Detect │ │RAG   │
     └───┬───┘ └──┬───┘ └──┬───┘ └──┬────┘ └──┬───┘
         │        │        │        │          │
         └────────┴────────┼────────┘          │
                           │                   │
                    ┌──────▼──────┐    ┌───────▼──────┐
                    │ Rules Engine│    │ Vector Store  │
                    └─────────────┘    └──────────────┘
```

### Deployment
Each new service gets its own Docker container in `docker-compose.yml`, its own health check, and its own environment variables following the `<SERVICE_NAME>_<SETTING>` convention.

---

## Estimated Effort

| Phase | Feature | Effort | Dependencies |
|-------|---------|--------|-------------|
| 1.1 | Rules Engine | 2-3 weeks | None |
| 1.2 | Fraud Detection | 3-4 weeks | Training data |
| 1.3 | Dynamic Quoting | 2-3 weeks | Rules Engine (1.1) |
| 2.1 | AI Policy Analysis | 4-6 weeks | Document ingestion pipeline |
| 2.2 | Insurance Chatbot | 2-3 weeks | AI Policy Analysis (2.1) |
| 3.1 | Vehicle Insurance | 1-2 weeks | Quoting Engine (1.3) |
| 3.2 | Life Insurance | 2-3 weeks | Quoting Engine (1.3) |
| 3.3 | P&C Actuarial | 3-4 weeks | Underwriting domain data |
| 3.4 | CRM Enhancement | 2-3 weeks | Partner/B2B portals |

**Total estimated effort: 21-31 weeks** (parallelizable across teams)

---

## Getting Started

Use OpenCode with the `@upgrade-analyst` subagent to dive deeper into any specific upgrade:

```
@upgrade-analyst Analyze the SmartClaim-RAG-Engine reference and design detailed proto schemas for policy analysis
```

Or use the custom command:
```
/analyze-upgrade
```
