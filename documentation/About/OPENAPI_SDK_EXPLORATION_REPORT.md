# InsureTech OpenAPI Specs & SDK Generation - Exploration Report

## Executive Summary

The InsureTech project has a comprehensive OpenAPI 3.1.0 specification with modular YAML-based path and schema definitions organized by domain. The TypeScript SDK is generated using `@hey-api/openapi-ts` with a centralized configuration. Billing, invoicing, payment, and B2B functionalities are fully documented and exposed via dedicated services.

---

## 1. OpenAPI Configuration File

### Location
`E:\Projects\InsureTech\sdks\sdk-generator\typescript\openapi-ts.config.ts`

### Configuration Details
```typescript
import { defineConfig } from '@hey-api/openapi-ts';

export default defineConfig({
  client: '@hey-api/client-fetch',
  input: '../../../api/openapi.yaml',                    // Points to main spec
  output: {
    path: '../../insuretech-typescript-sdk/src',        // Generated code output
    format: false,                                       // Formatting disabled in pipeline
    lint: false,                                         // Linting disabled in pipeline
  },
  types: {
    enums: 'javascript',
  },
  services: {
    asClass: true,
    name: '{{name}}Service',                            // Groups methods by OpenAPI tags
  },
  schemas: false,
});
```

**Key Points:**
- Uses `@hey-api/client-fetch` as the HTTP client
- Input: `../../../api/openapi.yaml` (main aggregated specification)
- Output path: `../../insuretech-typescript-sdk/src/`
- Services are generated as classes (e.g., `AuthService`, `PolicyService`, `BillingService`)
- Formatting/linting post-processed by `generator.go` for performance
- Manual formatting can be run with: `npx prettier --write src/` and `npx eslint src/ --fix`

---

## 2. API Directory Structure - All Files

### Root Files
```
E:\Projects\InsureTech\api\
├── openapi.yaml                          # Main OpenAPI 3.1 spec (96,636 lines)
├── README.md                             # Documentation & generation instructions
├── ENDPOINT_MAP.md                       # API endpoint mapping reference
├── analyze_report.py                     # Analysis script
├── check_refs.py                         # Reference validation script
├── validate.py                           # OpenAPI validation script
├── db_proto_comparison.json              # Proto-DB mapping data
├── proto_schema_summary.json             # Proto schema reference
├── schema_api_mapping.json                # Schema-to-API mapping
├── validation_report.html                # Validation report (HTML)
└── validation_report.json                # Validation report (JSON)
```

### Subdirectories

#### `/descriptions` - Markdown Descriptions for DTOs & Entities
- **dto/insuretech/** - DTO descriptions for request/response objects
  - ai/services/v1/ - AI service DTOs
  - analytics/services/v1/ - Analytics DTOs
  - common/v1/ - Common DTOs (PaginationRequest.md)
- **entity/insuretech/** - Entity descriptions
  - ai/entity/v1/ - AI entities (AIAgent, AIAnalysis, Conversation, etc.)
  - common/v1/ - Common DB options (ColumnOptions, ForeignKey, IndexOptions, TableOptions)
- **event/insuretech/** - Event descriptions
  - ai/events/v1/ - AI-related events (12 event types including fraud detection, decision making, etc.)

#### `/docs` - Documentation (empty placeholder)

#### `/enums` - Enumeration Definitions (empty placeholder)

#### `/events` - Event Definitions (empty placeholder)

#### `/generator` - Generation Scripts (empty placeholder)

#### `/input` - Input Configuration (empty placeholder)

#### `/postman` - Postman Collection Files (empty placeholder)

#### `/paths` - OpenAPI Path Definitions (Service Endpoints)

**All 34 Service Domains with YAML path definitions:**

| Service | YAML File | Purpose |
|---------|-----------|---------|
| AI | `paths/insuretech/ai/services/v1/AIService.yaml` | AI-powered features (chat, claims evaluation, fraud detection, risk assessment) |
| Analytics | `paths/insuretech/analytics/services/v1/AnalyticsService.yaml` | Dashboards, metrics, queries, reports |
| API Key | `paths/insuretech/apikey/services/v1/ApiKeyService.yaml` | API key management & validation |
| Audit | `paths/insuretech/audit/services/v1/AuditService.yaml` | Audit logging & tracking |
| Authentication | `paths/insuretech/authn/services/v1/AuthService.yaml` | Login, register, OTP, session management |
| Authorization | `paths/insuretech/authz/services/v1/AuthZService.yaml` | Roles, permissions, policies, access control |
| **B2B** | `paths/insuretech/b2b/services/v1/B2BService.yaml` | Organizations, departments, employees, purchase orders |
| Beneficiary | `paths/insuretech/beneficiary/services/v1/BeneficiaryService.yaml` | Beneficiary management & KYC |
| **Billing** | `paths/insuretech/billing/services/v1/BillingService.yaml` | Invoice creation, management, PDF generation |
| Claims | `paths/insuretech/claims/services/v1/ClaimService.yaml` | Claim submission, approval, rejection, settlement |
| Commission | `paths/insuretech/commission/services/v1/CommissionService.yaml` | Commission calculation & payouts |
| Document | `paths/insuretech/document/services/v1/DocumentService.yaml` | Document generation & management |
| Endorsement | `paths/insuretech/endorsement/services/v1/EndorsementService.yaml` | Endorsement approval/rejection |
| Fraud | `paths/insuretech/fraud/services/v1/FraudService.yaml` | Fraud detection & risk management |
| Insurer | `paths/insuretech/insurer/services/v1/InsurerService.yaml` | Insurer management |
| IoT | `paths/insuretech/iot/services/v1/IoTService.yaml` | IoT device integration |
| KYC | `paths/insuretech/kyc/services/v1/KYCService.yaml` | KYC verification |
| Media | `paths/insuretech/media/services/v1/MediaService.yaml` | Media/file handling |
| MFS | `paths/insuretech/mfs/services/v1/MFSService.yaml` | Mobile Financial Services (bKash, Nagad, etc.) |
| Notification | `paths/insuretech/notification/services/v1/NotificationService.yaml` | Notification delivery |
| Orders | `paths/insuretech/orders/services/v1/OrderService.yaml` | Order management |
| Partner | `paths/insuretech/partner/services/v1/PartnerService.yaml` | Partner integration |
| **Payment** | `paths/insuretech/payment/services/v1/PaymentService.yaml` | Payment processing, refunds, verification |
| Policy | `paths/insuretech/policy/services/v1/PolicyService.yaml` | Policy management & lifecycle |
| Product | `paths/insuretech/products/services/v1/ProductService.yaml` | Insurance product catalog |
| Refund | `paths/insuretech/refund/services/v1/RefundService.yaml` | Refund processing |
| Renewal | `paths/insuretech/renewal/services/v1/RenewalService.yaml` | Policy renewal |
| Report | `paths/insuretech/report/services/v1/ReportService.yaml` | Report generation |
| Support | `paths/insuretech/support/services/v1/SupportService.yaml` | Customer support & FAQs |
| Task | `paths/insuretech/task/services/v1/TaskService.yaml` | Task management |
| Tenant | `paths/insuretech/tenant/services/v1/TenantService.yaml` | Multi-tenancy management |
| Underwriting | `paths/insuretech/underwriting/services/v1/UnderwritingService.yaml` | Quote generation & approval |
| Voice | `paths/insuretech/voice/services/v1/VoiceService.yaml` | Voice interaction |
| Workflow | `paths/insuretech/workflow/services/v1/WorkflowService.yaml` | Business process workflows |

#### `/schemas` - OpenAPI Schema Definitions

**Organization:**
```
schemas/
├── google/                              # Google API extensions
│   └── api/ - HttpRule, CustomHttpPattern, Http definitions
└── insuretech/
    ├── ai/entity/v1/                   # AI domain entities
    ├── analytics/entity/v1/            # Analytics entities
    ├── apikey/entity/v1/               # API Key entities
    ├── audit/entity/v1/                # Audit entities
    ├── authn/entity/v1/                # Auth entities
    ├── authz/entity/v1/                # AuthZ entities
    ├── b2b/entity/v1/                  # B2B entities
    ├── beneficiary/entity/v1/          # Beneficiary entities
    ├── billing/entity/v1/              # **BILLING SCHEMAS**
    ├── claims/entity/v1/               # Claims entities
    ├── commission/entity/v1/           # Commission entities
    ├── common/v1/                      # Shared types:
    │   ├── Address.yaml
    │   ├── Money.yaml                  # Payment amount type
    │   ├── Document.yaml
    │   ├── Error.yaml & ErrorInfo.yaml
    │   ├── PageRequest/Response.yaml   # Pagination
    │   ├── Timestamps.yaml
    │   └── [22 total shared schemas]
    ├── document/entity/v1/             # Document entities
    ├── endorsement/entity/v1/          # Endorsement entities
    ├── fraud/entity/v1/                # Fraud detection entities
    ├── insurance/services/v1/          # Insurance services
    ├── insurer/entity/v1/              # Insurer entities
    ├── iot/entity/v1/                  # IoT entities
    ├── kyc/entity/v1/                  # KYC entities
    ├── media/entity/v1/                # Media entities
    ├── mfs/entity/v1/                  # MFS entities
    ├── notification/entity/v1/         # Notification entities
    ├── orders/entity/v1/               # Order entities
    ├── partner/entity/v1/              # Partner entities
    ├── payment/entity/v1/              # **PAYMENT SCHEMAS**
    ├── policy/entity/v1/               # Policy entities
    ├── products/entity/v1/             # Product entities
    ├── refund/entity/v1/               # Refund entities
    ├── renewal/entity/v1/              # Renewal entities
    ├── report/entity/v1/               # Report entities
    ├── services/entity/v1/             # Service entities
    ├── storage/entity/v1/              # Storage entities
    ├── support/entity/v1/              # Support entities
    ├── task/entity/v1/                 # Task entities
    ├── tenant/entity/v1/               # Tenant entities
    ├── underwriting/entity/v1/         # Underwriting entities
    ├── voice/entity/v1/                # Voice entities
    ├── webrtc/v1/                      # WebRTC entities
    └── workflow/entity/v1/             # Workflow entities
```

#### `/templates` - Template Files (empty placeholder)

---

## 3. Key Service Details: Billing, Payment, B2B

### **BILLING SERVICE**
**File:** `E:\Projects\InsureTech\api\paths\insuretech\billing\services\v1\BillingService.yaml`

**Endpoints:**
| Endpoint | Method | Operation | Purpose |
|----------|--------|-----------|---------|
| `/v1/invoices` | POST | `BillingService_CreateInvoice` | Create invoice for B2C orders or B2B purchase orders |
| `/v1/invoices` | GET | `BillingService_ListInvoices` | List invoices with optional filters |
| `/v1/invoices/{invoice_id}` | GET | `BillingService_GetInvoice` | Get single invoice by ID |
| `/v1/invoices/{invoice_id}:mark-paid` | POST | `BillingService_MarkInvoicePaid` | Mark invoice as paid (called by payment-service) |
| `/v1/invoices/{invoice_id}:cancel` | POST | `BillingService_CancelInvoice` | Cancel invoice (before PAID status) |
| `/v1/invoices/{invoice_id}:issue` | POST | `BillingService_IssueInvoice` | Transition DRAFT → ISSUED and send to customer/org |
| `/v1/invoices/{invoice_id}/pdf` | GET | `BillingService_GetInvoicePDF` | Get invoice PDF (pre-signed URL or file ID) |
| `/v1/invoices/{invoice_id}:generate-pdf` | POST | `BillingService_GenerateInvoicePDF` | Trigger async PDF generation |
| `/v1/orders/{order_id}/invoice` | GET | `BillingService_GetInvoiceByOrderId` | Get invoice by order ID (orders-service integration) |

---

### **PAYMENT SERVICE**
**File:** `E:\Projects\InsureTech\api\paths\insuretech\payment\services\v1\PaymentService.yaml`

**Endpoints:**
| Endpoint | Method | Operation | Purpose |
|----------|--------|-----------|---------|
| `/v1/payments` | POST | `PaymentService_InitiatePayment` | Initiate payment processing |
| `/v1/payments` | GET | `PaymentService_ListPayments` | List payments |
| `/v1/payments/{payment_id}:verify` | POST | `PaymentService_VerifyPayment` | Verify payment status |
| `/v1/payments/{payment_id}` | GET | `PaymentService_GetPayment` | Get payment details |
| `/v1/payments/{payment_id}/refunds` | POST | `PaymentService_InitiateRefund` | Initiate refund |
| `/v1/refunds/{refund_id}/status` | GET | `PaymentService_GetRefundStatus` | Get refund status |
| `/v1/users/{user_id}/payment-methods` | GET | `PaymentService_ListPaymentMethods` | List payment methods |
| `/v1/users/{user_id}/payment-methods` | POST | `PaymentService_AddPaymentMethod` | Add new payment method |
| `/v1/payments:reconcile` | POST | `PaymentService_ReconcilePayments` | Payment reconciliation |
| `/v1/payments/webhook/{provider}` | POST | `PaymentService_HandleGatewayWebhook` | Gateway webhook handler (SSLCommerz, bKash, Nagad) |
| `/v1/payments/provider/{provider}/references/{provider_reference}` | GET | `PaymentService_GetPaymentByProviderReference` | Lookup by provider reference |
| `/v1/payments/{payment_id}:submit-proof` | POST | `PaymentService_SubmitManualPaymentProof` | Submit manual bank transfer proof |
| `/v1/payments/{payment_id}:review` | POST | `PaymentService_ReviewManualPayment` | Admin reviews manual payment |
| `/v1/payments/{payment_id}:generate-receipt` | POST | `PaymentService_GenerateReceipt` | Trigger async receipt PDF generation |
| `/v1/payments/{payment_id}/receipt` | GET | `PaymentService_GetPaymentReceipt` | Retrieve generated receipt |

**Payment Providers Supported:**
- SSLCommerz
- bKash (Mobile banking)
- Nagad (Mobile banking)
- Manual bank transfer

---

### **B2B SERVICE**
**File:** `E:\Projects\InsureTech\api\paths\insuretech\b2b\services\v1\B2BService.yaml`

**Endpoints (Organization Management):**
| Endpoint | Method | Operation | Purpose |
|----------|--------|-----------|---------|
| `/v1/b2b/organisations` | POST | `B2BService_CreateOrganisation` | Create organisation (SuperAdmin only) |
| `/v1/b2b/organisations` | GET | `B2BService_ListOrganisations` | List organisations (role-based visibility) |
| `/v1/b2b/organisations/{organisation_id}` | GET | `B2BService_GetOrganisation` | Get single organisation |
| `/v1/b2b/organisations/{organisation_id}` | PATCH | `B2BService_UpdateOrganisation` | Update organisation profile |
| `/v1/b2b/organisations/{organisation_id}` | DELETE | `B2BService_DeleteOrganisation` | Soft-delete organisation & revoke memberships |
| `/v1/b2b/organisations/{organisation_id}/members` | GET | `B2BService_ListOrgMembers` | List organisation members |
| `/v1/b2b/organisations/{organisation_id}/members` | POST | `B2BService_AddOrgMember` | Add platform user as OrgMember |
| `/v1/b2b/organisations/{organisation_id}/admins` | POST | `B2BService_AssignOrgAdmin` | Assign user as OrgAdmin |
| `/v1/b2b/organisations/{organisation_id}/members/{member_id}` | DELETE | `B2BService_RemoveOrgMember` | Remove OrgMember |

**Endpoints (Department Management):**
| Endpoint | Method | Operation | Purpose |
|----------|--------|-----------|---------|
| `/v1/b2b/departments` | GET | `B2BService_ListDepartments` | List departments for authenticated org |
| `/v1/b2b/departments` | POST | `B2BService_CreateDepartment` | Create new department |
| `/v1/b2b/departments/{department_id}` | GET | `B2BService_GetDepartment` | Get single department |
| `/v1/b2b/departments/{department_id}` | PATCH | `B2BService_UpdateDepartment` | Update department name |
| `/v1/b2b/departments/{department_id}` | DELETE | `B2BService_DeleteDepartment` | Soft-delete department (no active employees) |

**Endpoints (Employee Management):**
| Endpoint | Method | Operation | Purpose |
|----------|--------|-----------|---------|
| `/v1/b2b/employees` | GET | `B2BService_ListEmployees` | List employees for authenticated org |
| `/v1/b2b/employees` | POST | `B2BService_CreateEmployee` | Create new employee |
| `/v1/b2b/employees/{employee_uuid}` | GET | `B2BService_GetEmployee` | Get single employee |
| `/v1/b2b/employees/{employee_uuid}` | PATCH | `B2BService_UpdateEmployee` | Update employee details |
| `/v1/b2b/employees/{employee_uuid}` | DELETE | `B2BService_DeleteEmployee` | Soft-delete employee |

**Endpoints (Purchase Orders):**
| Endpoint | Method | Operation | Purpose |
|----------|--------|-----------|---------|
| `/v1/b2b/purchase-orders/catalog` | GET | `B2BService_ListPurchaseOrderCatalog` | List purchasable product plans |
| `/v1/b2b/purchase-orders` | GET | `B2BService_ListPurchaseOrders` | List purchase orders for authenticated org |
| `/v1/b2b/purchase-orders` | POST | `B2BService_CreatePurchaseOrder` | Create purchase order for product plan |
| `/v1/b2b/purchase-orders/{purchase_order_id}` | GET | `B2BService_GetPurchaseOrder` | Get single purchase order |

---

## 4. SDK Structure

### TypeScript SDK
**Location:** `E:\Projects\InsureTech\sdks\insuretech-typescript-sdk\`

**Generated Files:**
```
insuretech-typescript-sdk/
├── src/
│   ├── client-wrapper.ts          # Custom client wrapper
│   ├── client.gen.ts              # Generated client
│   ├── sdk.gen.ts                 # Generated SDK
│   ├── types.gen.ts               # Generated types
│   ├── types.ts                   # Custom types
│   ├── errors.ts                  # Error handling
│   ├── index.ts                   # Main exports
│   ├── client/                    # Client implementations
│   └── core/                      # Core utilities
├── tests/
│   ├── e2e/                       # End-to-end tests
│   ├── integration/               # Integration tests
│   ├── unit/                      # Unit tests
│   ├── helpers/                   # Test helpers
│   └── setup.ts
├── package.json
├── tsconfig.json
├── vitest.config.ts
├── .eslintrc.json
├── .prettierrc
├── README.md
├── SDK_ANALYSIS_REPORT.md
├── TEST_STATUS.md
└── TEST_SUITE_SUMMARY.md
```

### Go SDK
**Location:** `E:\Projects\InsureTech\sdks\insuretech-go-sdk\`

**Structure:**
```
insuretech-go-sdk/
├── pkg/
│   ├── client/                    # Generated client
│   ├── models/                    # Generated models
│   ├── services/                  # Generated services
│   └── errors/                    # Error types
├── docs/
├── examples/
├── go.mod
└── README.md
```

### SDK Generator
**Location:** `E:\Projects\InsureTech\sdks\sdk-generator\`

**Components:**
```
sdk-generator/
├── typescript/
│   ├── openapi-ts.config.ts       # TypeScript config (see section 1)
│   ├── generator.go               # Generator CLI
│   ├── package.json
│   ├── generate.sh
│   ├── generate.ps1               # PowerShell generation script
│   ├── templates/                 # Code templates
│   └── README.md
├── go/
│   ├── generator.go               # Go generator CLI
│   ├── generate.sh
│   ├── generate.ps1
│   ├── templates/                 # Code templates
│   ├── go.mod
│   ├── go.sum
│   └── README.md
```

---

## 5. Common Schemas Used Across Services

**Location:** `E:\Projects\InsureTech\api\schemas\insuretech\common\v1\`

| Schema | Purpose |
|--------|---------|
| `Address.yaml` | Physical address type |
| `Money.yaml` | Monetary amount with currency |
| `Document.yaml` | Document metadata |
| `Error.yaml` | Error response schema |
| `ErrorInfo.yaml` | Error details |
| `FieldViolation.yaml` | Field-level validation errors |
| `PageRequest.yaml` | Pagination request |
| `PageResponse.yaml` | Pagination response |
| `PaginationRequest.yaml` | Alternative pagination request |
| `PaginationResponse.yaml` | Alternative pagination response |
| `Timestamps.yaml` | Created/updated timestamps |
| `TenantContext.yaml` | Tenant identification |
| `RequestContext.yaml` | Request metadata |
| `UUID.yaml` | UUID type |
| `ApprovalInfo.yaml` | Approval workflow info |
| `AuditInfo.yaml` | Audit trail info |
| `ContactInfo.yaml` | Contact details |
| `NIDInfo.yaml` | National ID info |
| `TINInfo.yaml` | Tax ID info |
| `ColumnOptions.yaml` | Database column options |
| `IndexOptions.yaml` | Database index options |
| `TableOptions.yaml` | Database table options |
| `ForeignKey.yaml` | Database foreign key |
| `FieldSorting.yaml` | Sort field specification |

---

## 6. Generation Process

### From Proto → OpenAPI → SDK

```
1. Proto Service Definitions
   └→ *_service.proto files with google.api.http annotations
   
2. OpenAPI Generation (Python)
   └→ python api/generator.py
   └→ Reads all proto services
   └→ Extracts HTTP annotations
   └→ Converts to OpenAPI schemas
   └→ Outputs: api/openapi.yaml
   
3. TypeScript SDK Generation
   └→ Input: api/openapi.yaml
   └→ Config: sdks/sdk-generator/typescript/openapi-ts.config.ts
   └→ Tool: @hey-api/openapi-ts
   └→ Output: sdks/insuretech-typescript-sdk/src/
   └→ Post-process: generator.go (formatting, linting)
   
4. Go SDK Generation
   └→ Similar process using Go generator
```

### Manual Generation Commands

```bash
# TypeScript SDK
cd sdks/sdk-generator/typescript
./generate.sh                          # Unix/Linux
./generate.ps1                         # PowerShell
npx @hey-api/openapi-ts generate      # Direct

# Go SDK
cd sdks/sdk-generator/go
./generate.sh
./generate.ps1
```

---

## 7. OpenAPI Specification Details

**File:** `E:\Projects\InsureTech\api\openapi.yaml`

**Metadata:**
- OpenAPI Version: 3.1.0
- Title: InsureTech API
- Servers:
  - Production: `https://api.labaidinsuretech.com`
  - Staging: `https://staging-api.labaidinsuretech.com`
- File Size: ~96,636 lines

**Security Schemes:**
- BearerAuth (JWT tokens)
- API Key authentication
- OAuth2

**Content Types:**
- `application/json` (Primary)
- `application/grpc+json` (gRPC fallback)

**Features:**
- Standardized error responses
- Request/response examples
- Tag-based organization (by service/domain)
- Proto-generated descriptions
- Role-based operation access (SuperAdmin, BizAdmin, etc.)

---

## 8. Key Findings: Billing, Payment & B2B Integration

### **Billing & Invoice System**
- ✅ Full invoice lifecycle management (DRAFT → ISSUED → PAID/CANCELLED)
- ✅ Support for both B2C orders and B2B purchase orders
- ✅ Async PDF generation with pre-signed URLs
- ✅ Integration point with Payment Service (mark-paid callback)
- ✅ Integration point with Orders Service (get-invoice-by-order)

### **Payment Processing**
- ✅ Multiple payment gateways: SSLCommerz, bKash, Nagad
- ✅ Manual bank transfer support with proof submission & review
- ✅ Payment verification and reconciliation
- ✅ Refund management and status tracking
- ✅ Payment method management (add, list)
- ✅ Receipt PDF generation
- ✅ Provider webhook handlers for async callbacks

### **B2B Functionality**
- ✅ Multi-level organizational structure:
  - Organizations (SuperAdmin creation)
  - Departments (within organizations)
  - Employees (department members)
  - Members & Admins (role-based)
- ✅ Purchase order catalog and creation
- ✅ Soft-delete with cascading permissions
- ✅ Role-based visibility (SuperAdmin sees all, BizAdmin sees own org only)

---

## 9. Additional Notes

### Naming Convention
- Services use `{{name}}Service` pattern (e.g., BillingService, PaymentService)
- Operations follow `Service_Operation` pattern (e.g., BillingService_CreateInvoice)
- State transitions use colon syntax: `:operation` (e.g., `:mark-paid`, `:issue`, `:cancel`)

### Error Handling
- Standard error schema with field violations
- HTTP status codes: 200, 201, 204, 400, 401, 403, 409, 422, 429
- Standardized error metadata with request_id tracking

### Testing
- Test suites available for TypeScript SDK
- E2E, integration, and unit tests configured
- Vitest framework

### Documentation
- Generated from proto comments
- Markdown descriptions available for all DTOs and entities
- Postman collection available for manual testing

---

## 10. File Inventory - Complete List

### All .yaml files in /paths:
- 34 service YAML files (one per domain)
- Located at: `E:\Projects\InsureTech\api\paths\insuretech\{service}\services\v1\{Service}Service.yaml`

### All schema directories in /schemas:
- 34 domain-specific schema directories
- 1 common directory with 23 shared types
- Located at: `E:\Projects\InsureTech\api\schemas\insuretech\`

### Configuration files:
- Main spec: `openapi.yaml`
- SDK config: `openapi-ts.config.ts`
- Generation scripts: `generate.sh`, `generate.ps1` (in each SDK generator subdirectory)

### Documentation:
- API README: `E:\Projects\InsureTech\api\README.md`
- SDK READMEs: `E:\Projects\InsureTech\sdks\*/README.md`
- Analysis reports: `*.md` files in SDK directories
- Test documentation: `TEST_STATUS.md`, `TEST_SUITE_SUMMARY.md`

---

**Report Generated:** Exploration of InsureTech OpenAPI specs and SDK generation configuration  
**Services Documented:** 34 API services with focus on Billing, Payment, and B2B
