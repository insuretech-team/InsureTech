# InsureTech Project Directory Structure - Complete Overview

## 1. PROJECT LAYOUT

The InsureTech project is organized as a monorepo with the following top-level structure:

```
E:\Projects\InsureTech\
├── Android/                    # Mobile app (Android)
├── api/                        # API documentation and OpenAPI schemas
├── backend/                    # Backend services and infrastructure
├── b2b_portal/                 # B2B portal (Next.js/TypeScript)
├── backups/                    # Database backups
├── customer_portal/            # Customer-facing portal
├── docs/                       # Documentation
├── documentation/              # Additional documentation
├── gen/                        # Generated code
├── .github/                    # GitHub workflows and configs
├── iOS/                        # Mobile app (iOS)
├── ops/                        # Operations and deployment scripts
├── partner_portal/             # Partner portal application
├── proto/                      # Protocol Buffer definitions (gRPC)
├── rules/                      # Business rules
├── scripts/                    # Utility scripts
├── sdks/                       # SDK packages (Go, TypeScript)
├── specs/                      # Specification documents
├── system_portal/              # System administration portal
├── web_shared/                 # Shared web components/utilities
├── .devcontainer/              # Dev container configuration
├── .specify/                   # Specify configuration
├── buf.yaml, buf.lock          # Buf proto tool config
├── docker-compose.yml          # Local development environment
├── go.mod, go.sum              # Go workspace files
├── .env.example                # Environment variables template
├── InsureTech.code-workspace   # VS Code workspace config
└── README.md, START_HERE.md    # Project documentation
```

---

## 2. BACKEND SERVICES & MICROSERVICES

### Backend Structure: `E:\Projects\InsureTech\backend\inscore\`

The backend is a Go-based monolithic application with microservices architecture. Services are organized by domain:

**Microservices (located in `microservices/` directory):**

1. **AI Service** - Artificial intelligence capabilities
2. **Analytics Service** - Analytics and reporting
3. **Audit Service** - Audit logging and compliance
4. **Authentication (authn)** - User authentication and identity management
5. **Authorization (authz)** - Access control and permissions
6. **B2B Service** - Business-to-business operations
7. **Beneficiary Service** - Beneficiary management
8. **Conference/WebRTC Service** - Video conferencing capabilities
9. **Document Generation (docgen)** - Document generation service
10. **Fraud Detection** - Fraud detection and prevention
11. **Insurance Service** - Core insurance logic
12. **IoT Service** - Internet of Things integration
13. **KYC Service** - Know Your Customer verification
14. **Media Service** - Media management and processing
15. **Notification Service** - Multi-channel notifications
16. **OCR Service** - Optical Character Recognition
17. **Orders Service** - Order management
18. **Partner Service** - Partner management
19. **Payment Service** - Payment processing
20. **Storage Service** - File storage management
21. **Support Service** - Customer support/ticketing
22. **Tenant Service** - Multi-tenant management
23. **Workflow Service** - Workflow orchestration

**Additional Backend Components:**
- **Gateway** - API Gateway
- **Database Operations (dbops)** - Migration and schema management
- **Database Tools (dbx, dbsql, dbfix)** - Database utilities

**Database Migrations:** Located in `backend/inscore/db/migrations/`
- ai_schema
- analytics_schema
- authn_schema
- authz_schema
- b2b_schema
- compliance_schema
- insurance_schema
- iot_schema
- media_schema
- notification_schema
- partner_schema
- payment_schema
- storage_schema
- support_schema
- tenant_schema
- workflow_schema

**Infrastructure (infra/):**
- **Docker** - Dockerfiles for each service
- **Nginx** - Reverse proxy and web server configuration
- **Observability** - Prometheus, Grafana, and tracing setup

**Legacy .NET Services (insurance_engine/, polisync/):**
- InsuranceEngine - C# insurance service
- PoliSync - C# policy synchronization service

---

## 3. SERVICES & PACKAGES SUMMARY

### API Layer
- **Location:** `E:\Projects\InsureTech\api\`
- **Files:** 
  - `openapi.yaml` - OpenAPI/Swagger specification
  - `ENDPOINT_MAP.md` - API endpoint mapping
  - Schema, path, and event definitions
  - Validation reports and analysis tools

### Protocol Buffers (gRPC Definitions)
- **Location:** `E:\Projects\InsureTech\proto\insuretech\`
- **Services Defined:** 40+ proto files organized by domain
- **Structure:**
  - `/entity/v1/` - Entity message definitions
  - `/events/v1/` - Event definitions
  - `/services/v1/` - Service definitions

### SDK Packages
- **Location:** `E:\Projects\InsureTech\sdks\`
- **Packages:**
  1. `insuretech-typescript-sdk` - TypeScript/JavaScript SDK
  2. `insuretech-go-sdk` - Go SDK
  3. `sdk-generator` - Generator tools for SDK creation

---

## 4. B2B PORTAL STRUCTURE

### Overview
**Path:** `E:\Projects\InsureTech\b2b_portal\`
**Type:** Next.js 13+ application (TypeScript/React)

### Directory Structure

```
b2b_portal/
├── app/                          # Next.js app router
│   ├── api/                      # API routes
│   │   ├── auth/                 # Authentication endpoints
│   │   ├── dashboard/            # Dashboard APIs
│   │   ├── departments/          # Department management
│   │   ├── documents/            # Document management
│   │   ├── employees/            # Employee management
│   │   ├── organisations/        # Organization management
│   │   └── purchase-orders/      # Purchase order management
│   ├── [feature-pages]/          # Feature pages
│   │   ├── billing-invoices/
│   │   ├── claims/
│   │   ├── departments/
│   │   ├── employees/
│   │   ├── insurance-plans/
│   │   ├── login/
│   │   ├── organisations/
│   │   ├── payments/
│   │   ├── policies/
│   │   ├── profile/
│   │   ├── purchase-orders/
│   │   ├── quotations/
│   │   ├── settings/
│   │   └── team/
│   └── root files (layout.tsx, page.tsx, globals.css)
│
├── components/                   # React components
│   ├── auth/                     # Authentication components
│   │   └── login-form.tsx
│   ├── charts/                   # Chart components
│   │   ├── bar-chart-overview.tsx
│   │   ├── policy-overview-chart.tsx
│   │   └── revenue-overview-chart.tsx
│   ├── dashboard/                # Dashboard layout & components
│   │   ├── billing-invoices/
│   │   ├── departments/
│   │   ├── employees/
│   │   ├── insurance-plans/
│   │   ├── organisations/
│   │   ├── overview-activity/
│   │   ├── policy-overview/
│   │   ├── profile/
│   │   ├── purchase-orders/
│   │   ├── quick-access/
│   │   ├── settings/ (with partials)
│   │   ├── stats-cards/
│   │   ├── upcoming-payments/
│   │   └── dashboard-layout.tsx, sidebar, header
│   ├── modals/                   # Modal dialogs
│   ├── organisations/            # Organization-specific components
│   ├── payments/                 # Payment components
│   ├── policies/                 # Policy components
│   ├── team/                     # Team management components
│   └── ui/                       # Base UI components (design system)
│       ├── avatar, badge, button, card, checkbox
│       ├── dialog, dropdown-menu, field, input, label
│       ├── select, separator, sort-header, status-badge
│       ├── table, tabs, toast-banner
│
├── lib/                          # Utility libraries
│   ├── sdk/                      # **SDK CLIENT LAYER** (see below)
│   ├── auth/                     # Authentication utilities
│   ├── proto-generated/          # Generated protobuf files
│   ├── types/                    # TypeScript type definitions
│   ├── invoices.ts, navigation.ts, notifications.ts
│   ├── payments.ts, purchase-orders.ts, stats-cards.ts
│   └── utils.ts, workflows.ts
│
├── public/                       # Static assets
│   ├── icons/, logos/, navbar-icons/, quotations/
│   ├── stats-cards/, insurance-plans/
│
├── src/                          # Additional source files
│   ├── hooks/                    # Custom React hooks
│   │   ├── useCrudList.ts
│   │   ├── useEmployeeForm.ts
│   │   ├── useOrganisationForm.ts
│   │   └── useToast.ts
│   ├── lib/
│   │   ├── sdk/                  # **SDK CLIENT LAYER** (see section below)
│   │   ├── auth/                 # Session management
│   │   ├── proto-generated/      # Generated protobuf TypeScript files
│   │   └── types/
│
├── Config files
│   ├── tsconfig.json             # TypeScript configuration
│   ├── tailwind.config.ts        # Tailwind CSS configuration
│   ├── next.config.ts            # Next.js configuration
│   ├── eslint.config.mjs         # ESLint configuration
│   ├── postcss.config.mjs        # PostCSS configuration
│   └── package.json              # Dependencies
│
└── Documentation & Errors
    ├── README.md
    ├── features.md
    ├── build_errors.txt
    ├── typescript_errors.txt
    └── middleware.ts
```

---

## 5. B2B PORTAL SDK STRUCTURE - DETAILED

### Location: `E:\Projects\InsureTech\b2b_portal\src\lib\sdk\`

**ALL FILES IN SDK DIRECTORY:**

```
sdk/
├── api-helpers.ts                  # HTTP request helper utilities
├── auth-client.ts                  # Authentication client for login/logout/auth flows
├── b2b-sdk-client.ts              # Main B2B SDK client (orchestrator)
├── dashboard-config.ts             # Dashboard configuration and setup
├── department-client.ts            # Department CRUD operations
├── docgen-client.ts                # Document generation client
├── docgen-sdk-client.ts            # Document generation SDK wrapper
├── employee-client.ts              # Employee management client
├── index.ts                        # SDK exports/barrel file
├── organisation-client.ts          # Organization management client
├── purchase-order-client.ts        # Purchase order client
├── session-headers.ts              # Session and authentication headers
├── shared.ts                       # Shared utilities and constants
```

**SDK Clients Structure (12 TypeScript files):**

| Client | Purpose |
|--------|---------|
| `auth-client.ts` | Handles authentication workflows (login, logout, OTP, password changes) |
| `b2b-sdk-client.ts` | Main SDK facade that orchestrates all other clients |
| `department-client.ts` | Department CRUD and management operations |
| `employee-client.ts` | Employee management (CRUD, bulk uploads) |
| `organisation-client.ts` | Organization management and member operations |
| `purchase-order-client.ts` | Purchase order lifecycle management |
| `docgen-client.ts` | Document generation operations |
| `docgen-sdk-client.ts` | Document generation SDK wrapper |
| `dashboard-config.ts` | Dashboard configuration, stats, and activity data |
| `api-helpers.ts` | HTTP utilities, request/response handling |
| `session-headers.ts` | Session management and header construction |
| `shared.ts` | Shared utilities, error handling, constants |
| `index.ts` | Public exports of SDK |

**Proto-Generated Files:**

The SDK uses proto-generated TypeScript files located in:
`E:\Projects\InsureTech\b2b_portal\src\lib\proto-generated\`

These include generated message definitions for:
- AI, Analytics, API Keys, Audit, Authentication, Authorization
- B2B entities (departments, employees, organizations, purchase orders)
- Beneficiary, Billing, Claims, Commission, Documents
- Fraud detection, Insurance, IoT, KYC, Media
- Notifications, Orders, Payments, Policies, Products
- Refunds, Renewals, Reports, Storage, Support
- Tasks, Tenants, Underwriting, Voice, WebRTC, Workflows

---

## 6. KEY TECHNOLOGIES & ARCHITECTURE

### Frontend
- **Framework:** Next.js 13+ (App Router)
- **Language:** TypeScript
- **Styling:** Tailwind CSS
- **UI Components:** Custom design system
- **State Management:** Custom hooks + server components

### Backend
- **Language:** Go
- **gRPC:** Protocol Buffers for service communication
- **Databases:** Multiple schemas per domain
- **Message Queue:** Kafka
- **Search/Analytics:** Implied (analytics service)
- **Legacy Services:** C# (.NET) - InsuranceEngine, PoliSync

### Infrastructure
- **Containerization:** Docker
- **Reverse Proxy:** Nginx
- **Observability:** Prometheus, Grafana
- **Deployment:** Docker Compose

### Mobile
- **Android:** Native Android app
- **iOS:** Native iOS app

### SDKs
- **TypeScript SDK:** For client-side and Node.js integration
- **Go SDK:** For Go backend integration

---

## 7. DOMAIN SERVICES MAPPING

| Service | Responsibility | Proto Location |
|---------|-----------------|-----------------|
| authn | User authentication, OTP, sessions | `proto/insuretech/authn/` |
| authz | Authorization, roles, permissions | `proto/insuretech/authz/` |
| b2b | B2B operations (org, employees, departments) | `proto/insuretech/b2b/` |
| beneficiary | Beneficiary management | `proto/insuretech/beneficiary/` |
| billing | Invoicing and billing | `proto/insuretech/billing/` |
| claims | Claims processing | `proto/insuretech/claims/` |
| commission | Commission calculations | `proto/insuretech/commission/` |
| document | Document generation and templates | `proto/insuretech/document/` |
| fraud | Fraud detection and alerts | `proto/insuretech/fraud/` |
| insurance | Core insurance operations | `proto/insuretech/insurance/` |
| iot | IoT device management | `proto/insuretech/iot/` |
| kyc | KYC verification | `proto/insuretech/kyc/` |
| media | Media storage and processing | `proto/insuretech/media/` |
| mfs | Mobile financial services | `proto/insuretech/mfs/` |
| notification | Multi-channel notifications | `proto/insuretech/notification/` |
| orders | Order management | `proto/insuretech/orders/` |
| partner | Partner management | `proto/insuretech/partner/` |
| payment | Payment processing | `proto/insuretech/payment/` |
| policy | Policy management | `proto/insuretech/policy/` |
| products | Product and plan management | `proto/insuretech/products/` |
| refund | Refund processing | `proto/insuretech/refund/` |
| renewal | Policy renewals | `proto/insuretech/renewal/` |
| report | Reporting and analytics | `proto/insuretech/report/` |
| storage | File storage | `proto/insuretech/storage/` |
| support | Customer support/ticketing | `proto/insuretech/support/` |
| task | Task management | `proto/insuretech/task/` |
| tenant | Multi-tenancy management | `proto/insuretech/tenant/` |
| underwriting | Underwriting decisions | `proto/insuretech/underwriting/` |
| voice | Voice interactions | `proto/insuretech/voice/` |
| webrtc | Real-time communication | `proto/insuretech/webrtc/` |
| workflow | Workflow orchestration | `proto/insuretech/workflow/` |

---

## 8. API ORGANIZATION

**API Routes in B2B Portal:** `E:\Projects\InsureTech\b2b_portal\app\api\`

Key API endpoint groups:
- `/api/auth/*` - Authentication endpoints (login, logout, OTP, sessions, profile)
- `/api/dashboard/*` - Dashboard data endpoints
- `/api/departments/*` - Department CRUD
- `/api/documents/*` - Document management
- `/api/employees/*` - Employee management
- `/api/organisations/*` - Organization management
- `/api/purchase-orders/*` - Purchase order management

Each endpoint is a Next.js API route handler using TypeScript.

---

## 9. PROTO-GENERATED FILES IN B2B PORTAL

**Location:** `E:\Projects\InsureTech\b2b_portal\src\lib\proto-generated\`

Auto-generated TypeScript files from Protocol Buffers for:
- Google API annotations
- 40+ domain-specific services with entity, event, and service definitions
- Message types, enums, and service stubs

---

## SUMMARY

**InsureTech** is a comprehensive insurance technology platform with:

1. **Multi-service backend** (23+ Go microservices + legacy .NET services)
2. **Modern frontend** (Next.js B2B portal with React components)
3. **Multiple portals** (B2B, Customer, Partner, System, Admin)
4. **Mobile apps** (iOS and Android)
5. **SDK packages** (TypeScript and Go)
6. **Comprehensive gRPC-based architecture** using Protocol Buffers
7. **Multi-domain database schema** supporting specialized business logic
8. **Infrastructure as Code** (Docker, Nginx, observability stack)

The project is well-organized with clear separation of concerns, service boundaries, and client SDKs for integration.
