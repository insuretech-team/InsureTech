# InsureTech SDK & B2B Portal - Complete Code Structure Summary

## 📁 Directory Structure

### SDK Location: `E:\Projects\InsureTech\sdks\insuretech-typescript-sdk\`

```
insuretech-typescript-sdk/
├── src/
│   ├── index.ts                 # Main SDK entry point
│   ├── types.ts                 # Core API response types
│   ├── errors.ts                # Error classes and interfaces
│   ├── client-wrapper.ts        # Custom client factory
│   ├── sdk.gen.ts               # Auto-generated service functions (4534 lines)
│   ├── types.gen.ts             # Auto-generated TypeScript types
│   ├── client.gen.ts            # Auto-generated client implementation
│   ├── client/
│   │   ├── index.ts             # Client exports
│   │   ├── client.gen.ts        # Client implementation
│   │   ├── types.gen.ts         # Client types
│   │   └── utils.gen.ts         # Client utilities
│   └── core/
│       ├── auth.gen.ts          # Authentication configuration
│       ├── bodySerializer.gen.ts    # Request body serialization
│       ├── params.gen.ts        # Parameter building
│       ├── pathSerializer.gen.ts    # Path parameter serialization
│       ├── queryKeySerializer.gen.ts # Query key serialization
│       ├── serverSentEvents.gen.ts  # SSE support
│       ├── types.gen.ts         # Core types
│       └── utils.gen.ts         # Utility functions
├── tests/
│   ├── e2e/
│   │   ├── complete-flow.test.ts
│   │   └── error-handling.test.ts
│   ├── integration/
│   │   ├── auth/
│   │   ├── claim/
│   │   ├── policy/
│   │   └── product/
│   ├── unit/
│   ├── helpers/
│   └── setup.ts
├── package.json
├── tsconfig.json
├── vitest.config.ts
└── README.md
```

### B2B Portal Location: `E:\Projects\InsureTech\b2b_portal\`

```
b2b_portal/
├── src/
│   ├── lib/
│   │   ├── sdk/                 # SDK wrapper clients
│   │   │   ├── index.ts
│   │   │   ├── auth-client.ts
│   │   │   ├── b2b-sdk-client.ts
│   │   │   ├── employee-client.ts
│   │   │   ├── organisation-client.ts
│   │   │   ├── department-client.ts
│   │   │   ├── purchase-order-client.ts
│   │   │   ├── docgen-client.ts
│   │   │   ├── docgen-sdk-client.ts
│   │   │   ├── dashboard-config.ts
│   │   │   ├── session-headers.ts
│   │   │   ├── api-helpers.ts
│   │   │   └── shared.ts
│   │   ├── auth/
│   │   │   ├── backend-auth.ts
│   │   │   ├── resolve-user-id.ts
│   │   │   ├── session-store.ts
│   │   │   └── session.ts
│   │   ├── proto-generated/    # Protocol Buffer generated files
│   │   ├── types/
│   │   │   ├── auth.ts
│   │   │   ├── b2b.ts
│   │   │   ├── employee-form.ts
│   │   │   └── ui.ts
│   │   └── index.ts
│   └── hooks/
│       ├── useCrudList.ts
│       ├── useEmployeeForm.ts
│       ├── useOrganisationForm.ts
│       └── useToast.ts
├── app/
│   ├── api/
│   │   ├── auth/
│   │   │   ├── login/route.ts
│   │   │   ├── logout/route.ts
│   │   │   ├── session/route.ts
│   │   │   ├── otp:send/route.ts
│   │   │   ├── otp:verify/route.ts
│   │   │   └── ... (other auth routes)
│   │   ├── dashboard/
│   │   ├── departments/
│   │   ├── employees/
│   │   ├── organisations/
│   │   └── purchase-orders/
│   ├── components/
│   │   ├── auth/
│   │   ├── dashboard/
│   │   ├── modals/
│   │   └── ui/
│   └── layout.tsx
├── components/
│   ├── auth/
│   ├── dashboard/
│   ├── modals/
│   ├── ui/
│   └── ... (feature components)
├── public/
├── package.json
├── next.config.ts
├── tsconfig.json
└── middleware.ts
```

---

## 🔑 Key Files Content Summary

### SDK Core Files

#### 1. **index.ts** (Entry Point)
Exports:
- Core types: `ApiResponse`, `ResponseMeta`, `PaginationMeta`, `Money`, `Address`
- Error classes: `InsureTechApiError`, `ApiError`
- Type guards: `isApiSuccess`, `isApiError`, `unwrapData`
- Client factory: `createInsureTechClient`
- All generated services and types from `sdk.gen` and `types.gen`

#### 2. **types.ts** (Core API Contract)
Defines:
- `ApiResponse<T>` - Unified envelope for all API responses
- `ResponseMeta` - Request metadata
- `PaginationMeta` - Pagination information
- `PaginationRequest` - Standard pagination parameters
- `Money` - Currency amount type
- `Address` - Address structure
- Helper functions for response unwrapping and type guards

#### 3. **errors.ts** (Error Handling)
Defines:
- `ApiErrorDetail` interface
- `FieldViolation` interface
- `InsureTechApiError` class with methods:
  - `isClientError()` - HTTP 4xx check
  - `isServerError()` - HTTP 5xx check
  - `isRetryable()` - Retry capability check
  - `hasFieldViolations()` - Field validation error check

#### 4. **client-wrapper.ts** (SDK Client Factory)
Provides:
- `InsureTechClientConfig` interface with:
  - `apiKey: string` (required)
  - `baseUrl?: string` (optional, defaults to https://api.insuretech.com)
  - `headers?: Record<string, string>` (optional)
- `createInsureTechClient()` function - Creates configured client instance

#### 5. **sdk.gen.ts** (Auto-Generated Services)
Contains ~150+ exported service functions covering:
- AI Services (Chat, Claims Evaluation, Document Analysis, Fraud Detection, Risk Assessment)
- Analytics (Dashboards, Metrics, Queries, Reports)
- API Key Management
- Audit Logging
- Authentication & Authorization
- B2B Operations (Orgs, Employees, Departments, Purchase Orders)
- Claims Management
- Document Generation & Templates
- Endorsements
- Fraud Detection & Case Management
- Insurance Products
- IoT Device Management
- KYC & Verification
- Media Management
- Mobile Financial Services (MFS)
- Notifications
- Orders & Payments
- Partner Management
- Policies & Quotations
- Refunds
- Renewals & Grace Periods
- Reports & Tasks
- Tenant Management
- Underwriting
- Workflows
- WebRTC Services
- Voice Biometrics

#### 6. **types.gen.ts** (Auto-Generated Types)
Exports:
- 1000+ TypeScript type definitions for all API request/response objects
- Enums for various domains (AgentType, AiAgentStatus, ClaimStatus, etc.)
- Entity types (User, Policy, Claim, Employee, Organisation, etc.)
- Event types for event-driven architecture
- Request/Response payload types for all endpoints

---

## 📦 Package Dependencies

### SDK Dependencies
- **@hey-api/client-fetch** ^0.1.0 - HTTP client with fetch

### SDK Dev Dependencies
- **@types/node** ^20.10.0
- **@typescript-eslint/** ^6.15.0
- **vitest** ^1.1.0 - Test runner
- **prettier** ^3.1.1 - Code formatter
- **tsup** ^8.0.1 - Build tool
- **msw** ^2.0.0 - Mock service worker

### B2B Portal Dependencies
- **next** ^16.1.6 - React framework
- **react** 19.2.3
- **@lifeplus/insuretech-sdk** (local tarball)
- **@tanstack/react-table** ^8.21.3
- **@radix-ui/** - UI components
- **recharts** ^3.7.0 - Charts
- **@bufbuild/protobuf** ^2.11.0
- **pg** ^8.19.0 - PostgreSQL client

---

## 🔌 SDK Installation Method

The B2B portal installs the SDK via local file path:
```json
"@lifeplus/insuretech-sdk": "file:../sdks/insuretech-typescript-sdk/lifeplus-insuretech-sdk-0.1.0.tgz"
```

The SDK is packaged as a `.tgz` tarball and referenced from the relative path.

---

## 🚀 Service Function Pattern

All generated service functions follow this pattern:

```typescript
export const serviceName = <ThrowOnError extends boolean = false>(
  options: Options<ServiceNameData, ThrowOnError>
) => (options.client ?? client).method<ServiceNameResponses, ServiceNameErrors, ThrowOnError>({
  security: [{ scheme: 'bearer', type: 'http' }],
  url: '/v1/endpoint-path/{param}',
  ...options,
  headers: {
    'Content-Type': 'application/json',
    ...options.headers
  }
});
```

**Key Features:**
- Type-safe request and response data
- Optional custom client instance
- Configurable error handling via `ThrowOnError` generic
- Automatic Bearer token authentication
- Support for path parameters, query parameters, and request bodies

---

## 🔐 Authentication Flow

1. **Client Creation:**
   ```typescript
   const client = createInsureTechClient({
     apiKey: 'your-api-key',
     baseUrl: 'https://api.insuretech.com'
   });
   ```

2. **Automatic Authorization:**
   - Client automatically adds `Authorization: Bearer {apiKey}` header
   - Applied to all authenticated endpoints marked with `security: [{ scheme: 'bearer', type: 'http' }]`

3. **Session Management:**
   - Sessions are managed server-side in B2B portal via `session-store.ts`
   - Cookies are used for session persistence
   - Middleware in `middleware.ts` validates sessions

---

## 📝 Response Handling

### Option 1: Type Guard Pattern
```typescript
const response = await aiServiceChat({ client, body: {} });

if (isApiSuccess(response)) {
  // response.data is guaranteed to be T
  const data: AiServiceChatResponses = response.data;
} else if (isApiError(response)) {
  // response.error is guaranteed to be ApiErrorDetail
  console.error(response.error.message);
}
```

### Option 2: Unwrap Pattern
```typescript
try {
  const data = unwrapData(response);
  // Use data directly
} catch (error) {
  if (error instanceof InsureTechApiError) {
    console.error('Code:', error.code);
    console.error('Status:', error.statusCode);
    console.error('Retryable:', error.retryable);
    console.error('Violations:', error.fieldViolations);
  }
}
```

---

## 🛠️ Build Process

### SDK Build
```bash
npm run build  # tsup src/index.ts --format cjs,esm --dts --clean
```

Outputs:
- `dist/index.js` (CommonJS)
- `dist/index.mjs` (ES Module)
- `dist/index.d.ts` (TypeScript Declarations)

### B2B Portal Build
```bash
npm run build  # next build
```

---

## 🧪 Testing

### SDK Tests
```bash
npm run test           # vitest run --passWithNoTests
npm run test:watch    # vitest
npm run test:coverage # vitest run --coverage
```

Test locations:
- `tests/unit/` - Unit tests
- `tests/integration/` - Integration tests
- `tests/e2e/` - End-to-end tests
- `tests/helpers/` - Test utilities and mock server

---

## 📊 Generated Code Volume

- **sdk.gen.ts**: 4534 lines of service function definitions
- **types.gen.ts**: 500+ lines sampled (complete file much larger)
- **Total Generated Code**: 5000+ lines

All generated files are marked with: `// This file is auto-generated by @hey-api/openapi-ts`

---

## 🔄 SDK Versioning

- **SDK Version**: 0.1.0
- **Package Name**: @lifeplus/insuretech-sdk
- **License**: MIT
- **Node Requirement**: >=16.0.0

---

## 📚 Key Type Definitions in types.gen.ts

```typescript
// Sample types from auto-generated file:
- AiAgent
- AiAgentStatus
- AiAgentCreatedEvent
- AccessGrantedEvent
- AccountLockedEvent
- ActionType
- AddInsurerProductRequest
- AddInsurerProductResponse
- AddOrgMemberRequest
- ... (1000+ more types)
```

---

## 🎯 B2B Portal SDK Integration Points

### 1. **API Route Handlers** (`app/api/`)
Each route uses SDK wrapper clients to call backend:
```typescript
// Example: app/api/auth/login/route.ts
import { authClient } from '@/src/lib/sdk';
const result = await authClient.login(credentials);
```

### 2. **React Hooks** (`src/hooks/`)
Custom hooks wrap SDK calls for component usage:
```typescript
// useCrudList.ts - Generic CRUD hook
// useEmployeeForm.ts - Employee-specific form logic
// useOrganisationForm.ts - Organisation form logic
```

### 3. **Components** (`components/`)
Components use hooks to interact with SDK:
```typescript
const { data, loading, error } = useCrudList('employees');
```

### 4. **Session Management** (`src/lib/auth/`)
Manages user session and authentication state across portal

---

## ✅ Summary

The InsureTech TypeScript SDK is a comprehensive, auto-generated API client that:
- **Type-Safe**: Full TypeScript support with 1000+ auto-generated types
- **Modular**: Organized into service categories
- **Flexible**: Supports custom clients, error handling modes, and headers
- **Production-Ready**: Includes error handling, validation, and pagination
- **Well-Integrated**: Seamlessly used throughout the B2B portal via wrapper clients
- **Maintainable**: Auto-generated from OpenAPI spec via @hey-api tooling

The B2B portal uses this SDK as its primary interface to the backend API, with additional wrapper layers for Next.js integration and React component usage.
