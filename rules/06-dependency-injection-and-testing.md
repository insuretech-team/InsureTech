# Rule 06: Dependency Injection & Frontend Testing Standards

**Scope:** iOS, Android, Frontend Web — API consumer side  
**Priority:** 🔴 CRITICAL  
**Goal:** The API contract must support DI and mock testing cleanly by guaranteeing:
1. a standard response envelope
2. success schemas free of embedded error objects
3. explicit per-endpoint security declarations
4. one canonical pagination contract
5. correct status codes for generic client behavior

---

## What "DI-Ready API" Means

A DI-ready REST API allows frontend/mobile clients to:
1. **Inject a single HTTP client** with auth, retry, logging as cross-cutting concerns
2. **Swap real API with mocks** in tests using the same interface
3. **Use one generic response decoder** for ALL endpoints
4. **Write one error handler** that works everywhere
5. **Generate type-safe SDKs** with zero manual boilerplate

---

## The DI Contract Our API MUST Fulfil

### 1. Consistent Envelope (enables single decoder)
```
ALL responses → ApiResponse<T> { success, data, error, meta }
```
If even ONE endpoint breaks this, the single generic decoder fails.

### 2. Correct Status Codes (enables generic interceptors)
```
2xx → unwrap data
4xx/5xx → unwrap error
201 + Location → cache invalidation / redirect
204 → no body expected
```

### 3. Per-endpoint Security Declaration (enables auth injection)
```
security: []            → no Authorization header
security: [BearerAuth]  → inject Authorization: Bearer <token>
security: [ApiKeyAuth]  → inject X-API-Key: <key>
```

### 4. Standard Error Shape (enables single error handler)
```
error.code          → programmatic switch/case
error.field_violations → form validation display
error.retryable     → retry logic
error.retry_after_seconds → backoff timer
```

### 5. Standard Pagination in meta (enables generic list component)
```
meta.pagination.page, total_pages, has_next → universal list UI
```

---

## TypeScript DI Setup (Frontend)

### Step 1: Define the Contract (generated from OpenAPI)

```typescript
// types/api.ts — auto-generated from openapi.yaml
export interface ApiResponse<T = unknown> {
  success: boolean;
  data: T | null;
  error: ApiError | null;
  meta: ResponseMeta | null;
}

export interface ApiError {
  code: string;
  message: string;
  field_violations?: FieldViolation[];
  error_id?: string;
  retryable?: boolean;
  retry_after_seconds?: number;
  http_status_code?: number;
}

export interface FieldViolation {
  field: string;
  message: string;
  code?: string;
  rejected_value?: string;
}

export interface ResponseMeta {
  request_id?: string;
  pagination?: PaginationMeta;
  timestamp?: string;
}

export interface PaginationMeta {
  page: number;
  page_size: number;
  total_pages: number;
  total_items: number;
  has_next: boolean;
  has_previous: boolean;
}
```

### Step 2: Define the Service Interface (injectable)

```typescript
// services/interfaces/IPolicyService.ts
export interface IPolicyService {
  createPolicy(req: PolicyCreationRequest): Promise<PolicyData>;
  getPolicy(policyId: string): Promise<PolicyData>;
  listPolicies(params: ListPoliciesParams): Promise<ListResult<PolicyData>>;
  cancelPolicy(policyId: string, reason: string): Promise<void>;
  renewPolicy(policyId: string): Promise<PolicyData>;
}

// services/interfaces/IAuthService.ts
export interface IAuthService {
  register(req: RegistrationRequest): Promise<RegistrationData>;
  login(req: LoginRequest): Promise<AuthTokenData>;
  refreshToken(refreshToken: string): Promise<AuthTokenData>;
  logout(): Promise<void>;
  sendOTP(phone: string): Promise<OTPData>;
  verifyOTP(otpId: string, code: string): Promise<void>;
}
```

### Step 3: Real Implementation (calls actual API)

```typescript
// services/impl/PolicyService.ts
import { apiClient } from '../http/apiClient';
import { IPolicyService } from '../interfaces/IPolicyService';

export class PolicyService implements IPolicyService {
  async createPolicy(req: PolicyCreationRequest): Promise<PolicyData> {
    // apiClient interceptor handles envelope unwrapping
    return apiClient.post<PolicyData>('/v1/policies', req);
  }

  async getPolicy(policyId: string): Promise<PolicyData> {
    return apiClient.get<PolicyData>(`/v1/policies/${policyId}`);
  }

  async listPolicies(params: ListPoliciesParams): Promise<ListResult<PolicyData>> {
    return apiClient.getList<PolicyData>('/v1/policies', params);
  }

  async cancelPolicy(policyId: string, reason: string): Promise<void> {
    return apiClient.post(`/v1/policies/${policyId}:cancel`, { reason });
  }
}
```

### Step 4: Mock Implementation (for tests and dev)

```typescript
// services/mock/MockPolicyService.ts
import { IPolicyService } from '../interfaces/IPolicyService';
import { mockPolicies } from '../fixtures/policies';

export class MockPolicyService implements IPolicyService {
  private policies = [...mockPolicies];

  async createPolicy(req: PolicyCreationRequest): Promise<PolicyData> {
    const policy: PolicyData = {
      policy_id: `pol_${Date.now()}`,
      policy_number: `INS-TEST-${this.policies.length + 1}`,
      status: 'ACTIVE',
      ...req
    };
    this.policies.push(policy);
    return policy;
  }

  async getPolicy(policyId: string): Promise<PolicyData> {
    const policy = this.policies.find(p => p.policy_id === policyId);
    if (!policy) throw new ApiError({ code: 'POLICY_NOT_FOUND', message: 'Not found' });
    return policy;
  }

  async listPolicies(params: ListPoliciesParams): Promise<ListResult<PolicyData>> {
    return {
      items: this.policies.slice(0, params.page_size ?? 20),
      pagination: {
        page: 1, page_size: 20, total_pages: 1,
        total_items: this.policies.length,
        has_next: false, has_previous: false
      }
    };
  }

  async cancelPolicy(policyId: string): Promise<void> {
    const idx = this.policies.findIndex(p => p.policy_id === policyId);
    if (idx === -1) throw new ApiError({ code: 'POLICY_NOT_FOUND', message: 'Not found' });
    this.policies[idx].status = 'CANCELLED';
  }
}
```

### Step 5: DI Container Registration

```typescript
// di/container.ts
import { Container } from 'inversify';  // or tsyringe, awilix, etc.

const container = new Container();

if (process.env.NODE_ENV === 'test' || process.env.USE_MOCKS === 'true') {
  container.bind<IPolicyService>('PolicyService').to(MockPolicyService);
  container.bind<IAuthService>('AuthService').to(MockAuthService);
  container.bind<IPaymentService>('PaymentService').to(MockPaymentService);
  // ... all 34 services
} else {
  container.bind<IPolicyService>('PolicyService').to(PolicyService);
  container.bind<IAuthService>('AuthService').to(AuthService);
  container.bind<IPaymentService>('PaymentService').to(PaymentService);
}

export { container };
```

### Step 6: The HTTP Client (one client handles all auth + unwrapping)

```typescript
// http/apiClient.ts
import axios, { AxiosInstance } from 'axios';
import { tokenStore } from '../store/tokenStore';

class ApiClient {
  private http: AxiosInstance;

  constructor(baseURL: string) {
    this.http = axios.create({ baseURL, timeout: 30000 });

    // REQUEST: auto-inject Bearer token for protected endpoints
    this.http.interceptors.request.use(config => {
      const token = tokenStore.getAccessToken();
      if (token && !config.headers['X-Public']) {
        config.headers.Authorization = `Bearer ${token}`;
      }
      config.headers['X-Request-ID'] = crypto.randomUUID();
      return config;
    });

    // RESPONSE: unwrap envelope — single place for ALL endpoints
    this.http.interceptors.response.use(
      response => {
        const body = response.data as ApiResponse<unknown>;
        if (!body.success) throw new ApiError(body.error!);
        return body.data;           // always returns T, not ApiResponse<T>
      },
      async error => {
        const status = error.response?.status;
        const apiError = error.response?.data?.error as ApiError;

        if (status === 401) {
          // Auto-refresh token once
          try {
            await this.refreshToken();
            return this.http.request(error.config);
          } catch {
            tokenStore.clear();
            window.location.href = '/login';
          }
        }

        if (status === 429 && apiError?.retryable) {
          const delay = (apiError.retry_after_seconds ?? 5) * 1000;
          await sleep(delay);
          return this.http.request(error.config);
        }

        throw apiError ?? error;
      }
    );
  }

  async get<T>(path: string, params?: object): Promise<T> {
    return this.http.get(path, { params }) as Promise<T>;
  }

  async post<T>(path: string, data?: object): Promise<T> {
    return this.http.post(path, data) as Promise<T>;
  }

  async getList<T>(path: string, params?: object): Promise<ListResult<T>> {
    const data = await this.http.get<ListData<T>>(path, { params });
    // meta.pagination is extracted by interceptor and attached separately
    return data as unknown as ListResult<T>;
  }

  private async refreshToken(): Promise<void> {
    const refresh = tokenStore.getRefreshToken();
    const response = await this.http.post<AuthTokenData>('/v1/auth/token:refresh', 
      { refresh_token: refresh }, 
      { headers: { 'X-Public': 'true' } }  // skip auth injection for this call
    );
    tokenStore.setTokens(response);
  }
}

export const apiClient = new ApiClient(process.env.API_BASE_URL!);
```

---

## Swift (iOS) DI Setup

```swift
// Protocol — the contract
protocol PolicyServiceProtocol {
    func createPolicy(_ request: PolicyCreationRequest) async throws -> PolicyData
    func getPolicy(id: String) async throws -> PolicyData
    func listPolicies(params: ListParams) async throws -> ListResult<PolicyData>
}

// Real implementation
class PolicyService: PolicyServiceProtocol {
    private let client: APIClient
    init(client: APIClient) { self.client = client }

    func createPolicy(_ request: PolicyCreationRequest) async throws -> PolicyData {
        try await client.post("/v1/policies", body: request)
    }
}

// Mock — for unit tests and SwiftUI previews
class MockPolicyService: PolicyServiceProtocol {
    func createPolicy(_ request: PolicyCreationRequest) async throws -> PolicyData {
        return PolicyData.mock()  // instant, no network
    }
}

// DI via environment — SwiftUI
struct ContentView: View {
    @EnvironmentObject var policyService: any PolicyServiceProtocol
}

// In tests
let view = ContentView().environmentObject(MockPolicyService())

// In production
let view = ContentView().environmentObject(PolicyService(client: APIClient.shared))
```

---

## Kotlin (Android) DI Setup with Hilt

```kotlin
// Interface
interface PolicyService {
    suspend fun createPolicy(request: PolicyCreationRequest): PolicyData
    suspend fun getPolicy(policyId: String): PolicyData
    suspend fun listPolicies(params: ListParams): ListResult<PolicyData>
}

// Real implementation
@Singleton
class PolicyServiceImpl @Inject constructor(
    private val api: InsureTechApi  // Retrofit interface
) : PolicyService {
    override suspend fun createPolicy(req: PolicyCreationRequest) =
        api.createPolicy(req).data!!
}

// Mock implementation
class MockPolicyService : PolicyService {
    override suspend fun createPolicy(req: PolicyCreationRequest) =
        PolicyData(policyId = "pol_test_${System.currentTimeMillis()}", ...)
}

// Hilt Module
@Module @InstallIn(SingletonComponent::class)
object ServiceModule {
    @Provides @Singleton
    fun providePolicyService(api: InsureTechApi): PolicyService =
        if (BuildConfig.USE_MOCKS) MockPolicyService()
        else PolicyServiceImpl(api)
}
```

---

## Mock Data Fixtures (shared across platforms)

The API team MUST provide official mock fixtures in the OpenAPI spec using `x-mock-response` or Apidog mock rules:

```yaml
# In each path definition — add example responses
/v1/policies:
  post:
    responses:
      '201':
        content:
          application/json:
            examples:
              success:
                summary: Policy created successfully
                value:
                  success: true
                  data:
                    policy_id: "pol_abc123"
                    policy_number: "INS-2024-001"
                    status: "ACTIVE"
                    premium_amount: { amount: "5000", currency: "BDT" }
                  error: null
                  meta:
                    request_id: "req_mock_001"
              error_422:
                summary: Validation failed
                value:
                  success: false
                  data: null
                  error:
                    code: "VALIDATION_FAILED"
                    message: "Invalid premium amount"
                    field_violations:
                      - field: "sum_insured"
                        message: "Must be between 100000 and 10000000"
                  meta:
                    request_id: "req_mock_002"
```

---

## Testing Checklist for Frontend Teams

Before any feature is considered complete:

- [ ] Unit test passes with `MockService` (no network)
- [ ] Integration test passes against staging API
- [ ] `ApiResponse<T>` generic decoder handles the endpoint
- [ ] Error scenarios tested: 401, 403, 422, 500
- [ ] Pagination works with generic list component
- [ ] Token refresh tested (expire token mid-session)
- [ ] Loading/empty/error states all handled
- [ ] Retry logic tested for retryable errors
