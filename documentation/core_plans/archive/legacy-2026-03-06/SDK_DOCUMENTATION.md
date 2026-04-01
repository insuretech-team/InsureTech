# InsureTech TypeScript SDK Documentation

## SDK Location
`E:/Projects/InsureTech/sdks/insuretech-typescript-sdk/`

## Package Information

**Name:** `@lifeplus/insuretech-sdk`  
**Version:** 0.1.0  
**Description:** Official TypeScript/JavaScript SDK for InsureTech API Platform  
**License:** MIT  
**Node Version:** >=16.0.0  
**Repository:** https://github.com/lifeplus/InsureTech

---

## Key Files in SDK

### 1. Entry Point: src/index.ts

```typescript
// Main SDK Entry Point - Auto-generated with custom enhancements

// Export all generated services and types
export * from './sdk.gen';
export * from './types.gen';

// Export custom client helper
export { createInsureTechClient } from './client-wrapper';
export type { InsureTechClientConfig } from './client-wrapper';
```

**Exports:**
- All service functions from `sdk.gen.ts`
- All type definitions from `types.gen.ts`
- `createInsureTechClient()` factory function
- `InsureTechClientConfig` interface

---

### 2. Client Wrapper: src/client-wrapper.ts

```typescript
// Custom Client Wrapper for InsureTech SDK
// Provides a configured client instance for use with generated services

import { createClient, createConfig } from './client';

export interface InsureTechClientConfig {
  /** API key for authentication */
  apiKey: string;
  /** Base URL for the API (optional, defaults to production) */
  baseUrl?: string;
  /** Additional headers to include in all requests */
  headers?: Record<string, string>;
}

/**
 * Create a configured client for the InsureTech API
 * 
 * @example
 * ```typescript
 * import { createInsureTechClient, AiService } from '@lifeplus/insuretech-sdk';
 * 
 * const client = createInsureTechClient({
 *   apiKey: 'your-api-key',
 *   baseUrl: 'https://api.insuretech.com'
 * });
 * 
 * // Use with any service method
 * const response = await AiService.aiServiceChat({
 *   client,
 *   body: { message: 'Hello' }
 * });
 * ```
 */
export function createInsureTechClient(config: InsureTechClientConfig) {
  return createClient(createConfig({
    baseUrl: config.baseUrl || 'https://api.insuretech.com',
    headers: {
      'Authorization': `Bearer ${config.apiKey}`,
      ...config.headers,
    },
  }));
}

// Re-export for convenience
export { createClient, createConfig } from './client';
```

---

### 3. Auto-Generated Files

#### src/sdk.gen.ts (6,958 lines)
Contains all service functions auto-generated from OpenAPI specification:
- Service pattern: `export const [servicePrefix][operationName] = (options) => ...`
- All services follow consistent pattern with type-safe options
- Examples: `aiServiceChat`, `claimServiceCreateClaim`, `policyServiceListPolicies`

#### src/types.gen.ts (34,737 lines)
Contains all type definitions:
- ClientOptions interface with supported base URLs
- Domain models (AiAgent, Conversation, etc.)
- Service-specific data types for each operation
- Enums for service states and types

---

## Department-Related SDK Methods

The B2B Portal uses these department methods from the SDK:

```typescript
// From @lifeplus/insuretech-sdk

export const b2bServiceListDepartments: (options) => Promise<...>;
export const b2bServiceCreateDepartment: (options) => Promise<...>;
export const b2bServiceGetDepartment: (options) => Promise<...>;
export const b2bServiceUpdateDepartment: (options) => Promise<...>;
export const b2bServiceDeleteDepartment: (options) => Promise<...>;
```

These are wrapped in the portal's SDK client factory (`b2b-sdk-client.ts`):

```typescript
// Usage in Next.js API handlers:
const sdk = makeSdkClient(request);
await sdk.listDepartments({ query: { page_size: 50 } });
await sdk.createDepartment({ body: { name, business_id } });
await sdk.getDepartment({ path: { department_id } });
await sdk.updateDepartment({ path: { department_id }, body: { name } });
await sdk.deleteDepartment({ path: { department_id } });
```

---

## SDK Structure

```
insuretech-typescript-sdk/
├── src/
│   ├── index.ts                 # Main entry point
│   ├── client-wrapper.ts        # Custom client configuration
│   ├── client.gen.ts            # Auto-generated client initialization
│   ├── sdk.gen.ts               # Auto-generated service functions (6,958 lines)
│   ├── types.gen.ts             # Auto-generated type definitions (34,737 lines)
│   ├── client/
│   │   ├── index.ts
│   │   ├── client.gen.ts        # Core HTTP client implementation (288 lines)
│   │   ├── types.gen.ts         # Client type definitions
│   │   └── utils.gen.ts
│   └── core/
│       ├── auth.gen.ts          # Authentication token handling
│       ├── bodySerializer.gen.ts # Request body serialization
│       ├── params.gen.ts        # Parameter building utilities
│       ├── pathSerializer.gen.ts # Path parameter serialization
│       └── utils.gen.ts         # URL building utilities
├── tests/
│   ├── unit/                    # Unit tests
│   ├── integration/             # Integration tests
│   └── e2e/                     # End-to-end tests
├── dist/                        # Build output (CJS + ESM + types)
├── package.json
├── tsconfig.json
├── .eslintrc.json
├── .prettierrc
├── vitest.config.ts
└── README.md
```

---

## Build Scripts

From `package.json`:

```bash
# Build CJS and ESM with type definitions
npm run build

# Build in watch mode
npm run build:watch

# Run tests
npm test

# Run tests in watch mode
npm test:watch

# Generate coverage report
npm run test:coverage

# Lint TypeScript files
npm run lint

# Lint and auto-fix
npm run lint:fix

# Format with Prettier
npm run format

# Check formatting
npm run format:check

# Type check without emitting
npm run typecheck

# Prepare for publishing
npm run prepublishOnly
```

---

## Build Output

The build process generates:

```
dist/
├── index.js           # CommonJS entry point
├── index.mjs          # ES Module entry point
├── index.d.ts         # TypeScript definitions
├── **/*.js            # CommonJS modules
├── **/*.mjs           # ES Module modules
├── **/*.d.ts          # Module type definitions
└── **/*.js.map        # Source maps
```

---

## Key Features

### ✅ Comprehensive API Coverage
- 14+ service modules covering all InsureTech operations
- 100+ exported service functions
- Full type safety with 34,000+ lines of TypeScript definitions

### ✅ Modern Architecture
- Built with @hey-api/openapi-ts generation tools
- Fetch API-based HTTP client
- Full TypeScript support (v5.3.3)
- ESM and CommonJS dual builds

### ✅ Developer Experience
- Zero-config client initialization with `createInsureTechClient()`
- Type-safe API calls with full IDE autocompletion
- Comprehensive error handling and response typing
- Request/response interceptor middleware support

### ✅ Quality Assurance
- Vitest-based test suite (unit, integration, E2E)
- MSW for server mocking
- ESLint with TypeScript enforcement
- Prettier code formatting
- Full test coverage reporting

### ✅ Production Ready
- MIT licensed open source
- Semantic versioning
- Published to npm as `@lifeplus/insuretech-sdk`
- Node.js 16+ support

---

## Supported Environments

```typescript
baseUrl: 
  | 'https://api.labaidinsuretech.com'           // Production
  | 'https://staging-api.labaidinsuretech.com'   // Staging
  | (string & {})                                 // Custom URL
```

Default: `https://api.labaidinsuretech.com`

---

## Dependencies

### Production
- `@hey-api/client-fetch` ^0.1.0 (HTTP client)

### Development
- TypeScript 5.3.3
- Vitest 1.1.0 (Testing framework)
- ESLint & TypeScript ESLint 6.15.0
- Prettier 3.1.1
- tsup 8.0.1 (Build tool)
- MSW 2.0.0 (Mock Service Worker)
- @vitest/coverage-v8 1.1.0

---

## Service Modules

The SDK includes these major service modules:

1. **AI Service** - Chat, document analysis, fraud detection, risk assessment
2. **Analytics Service** - Dashboards, reports, metrics, queries
3. **Auth Service** - Login, password reset, KYC, user profiles, documents
4. **API Key Service** - Generate, revoke, rotate, list, usage stats
5. **Audit Service** - Audit logs, compliance logs, audit trails
6. **Claim Service** - Create, get, update, list, submit documents
7. **Policy Service** - List, get, create, update, renew policies
8. **Product Service** - List, get, create products
9. **Quote Service** - Create, get, calculate premium
10. **Task Service** - Create, get, update, assign, list tasks
11. **Tenant Service** - List, create, get tenants
12. **Underwriting Service** - Quotes, health declarations, decisions
13. **Voice Service** - Voice sessions, commands, transcripts
14. **Workflow Service** - Definitions, instances, tasks, history

---

## Integration with B2B Portal

The B2B Portal uses the SDK in two ways:

### 1. Server-Side (API Routes)
```typescript
// In app/api/departments/route.ts
const sdk = makeSdkClient(request);
const result = await sdk.listDepartments({ query: { page_size: 50 } });
```

### 2. Browser-Side (Through API Routes)
```typescript
// In browser code
const result = await departmentClient.list();
// This sends: GET /api/departments
// Which internally calls the SDK on the server
```

This architecture provides:
- **Security**: API key and session handling on server only
- **Type Safety**: Full TypeScript support end-to-end
- **Simplicity**: Browser code uses simple REST API routes
- **Flexibility**: SDK can be used directly in server code when needed

---

## Proto-Generated Types

The SDK is generated from Protocol Buffer definitions. The B2B Portal includes proto-generated types at:

`E:/Projects/InsureTech/b2b_portal/src/lib/proto-generated/`

Structure:
```
proto-generated/
├── google/api/               # Google API annotations
├── insuretech/
│   ├── b2b/                 # B2B domain (departments, employees, etc.)
│   ├── auth/                # Authentication types
│   ├── claims/              # Claims types
│   ├── policy/              # Policy types
│   ├── ai/                  # AI service types
│   └── [other domains]/
```

These provide the canonical data structure definitions that the SDK is built from.

---

## Summary

The **InsureTech TypeScript SDK** provides a complete, production-ready API client for all InsureTech platform services. It offers:

1. **Full coverage** of all 14+ service domains
2. **Type safety** with 34,000+ lines of TypeScript definitions
3. **Modern tooling** with latest TypeScript and build tools
4. **Easy integration** through `createInsureTechClient()` factory
5. **Enterprise quality** with comprehensive tests and documentation

The B2B Portal uses it for authenticated server-side access to department and other B2B operations, providing a secure and type-safe foundation for the application.
