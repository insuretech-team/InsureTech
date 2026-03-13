# InsureTech TypeScript SDK - Comprehensive Analysis Report

## 1. Directory Structure

```
insuretech-typescript-sdk/
├── src/
│   ├── index.ts                    # Main entry point
│   ├── client.gen.ts               # Generated client factory
│   ├── client-wrapper.ts           # Client wrapper implementation
│   ├── sdk.gen.ts                  # Generated SDK service methods
│   ├── types.gen.ts                # Generated type definitions (39,777 lines)
│   ├── client/
│   │   ├── client.gen.ts           # Core client implementation
│   │   ├── index.ts                # Client exports
│   │   ├── types.gen.ts            # Client types
│   │   └── utils.gen.ts            # Client utilities
│   └── core/
│       ├── auth.gen.ts             # Authentication/authorization core
│       ├── bodySerializer.gen.ts   # Request body serialization
│       ├── params.gen.ts           # URL parameters serialization
│       ├── pathSerializer.gen.ts   # Path parameters serialization
│       ├── queryKeySerializer.gen.ts # Query key serialization
│       ├── serverSentEvents.gen.ts # Server-sent events support
│       ├── types.gen.ts            # Core type definitions
│       └── utils.gen.ts            # Core utilities
├── tests/
│   ├── unit/                       # Unit tests
│   ├── integration/                # Integration tests (organized by service)
│   │   ├── auth/                   # Authentication tests
│   │   ├── claim/                  # Claims service tests
│   │   ├── policy/                 # Policy service tests
│   │   └── product/                # Product service tests
│   ├── e2e/                        # End-to-end tests
│   └── helpers/                    # Test utilities
├── package.json                    # Dependencies and scripts
├── tsconfig.json                   # TypeScript configuration
├── vitest.config.ts                # Vitest configuration
├── .eslintrc.json                  # ESLint configuration
├── .prettierrc                      # Prettier configuration
└── README.md                        # Documentation
```

## 2. Authentication & Authorization - Key Generated Files

### 2.1 Core Auth (src/core/auth.gen.ts)

```typescript
export type AuthToken = string | undefined;

export interface Auth {
  /**
   * Which part of the request do we use to send the auth?
   * @default 'header'
   */
  in?: 'header' | 'query' | 'cookie';
  /**
   * Header or query parameter name.
   * @default 'Authorization'
   */
  name?: string;
  scheme?: 'basic' | 'bearer';
  type: 'apiKey' | 'http';
}

export const getAuthToken = async (
  auth: Auth,
  callback: ((auth: Auth) => Promise<AuthToken> | AuthToken) | AuthToken
): Promise<string | undefined> => {
  const token = typeof callback === 'function' ? await callback(auth) : callback;

  if (!token) {
    return;
  }

  if (auth.scheme === 'bearer') {
    return `Bearer ${token}`;
  }

  if (auth.scheme === 'basic') {
    return `Basic ${btoa(token)}`;
  }

  return token;
};
```

**Features:**
- Supports multiple authentication schemes: `apiKey`, `http` (basic/bearer)
- Flexible token placement: header, query, or cookie
- Async token resolution with callback support
- Automatic Bearer/Basic prefix handling

### 2.2 Authentication Message Types (Session & User)

#### Session Type
```typescript
/**
 * Session represents user authentication session (HYBRID: Server-side + JWT)
 * Maps to 'sessions' table in authn_schema
 */
export type Session = {
  session_id?: string;
  user_id?: string;
  session_type?: SessionType;  // SERVER_SIDE (web) or JWT (mobile)
  access_token_jti?: string;   // JWT ID for access token
  refresh_token_jti?: string;  // JWT ID for refresh token
  access_token_expires_at?: string;
  refresh_token_expires_at?: string;
  expires_at?: string;         // 12 hours (SERVER_SIDE) or 7 days (JWT)
  ip_address?: string;
  user_agent?: string;
  device_id?: string;
  device_name?: string;
  device_type?: AuthnDeviceType;  // WEB, MOBILE_ANDROID, MOBILE_IOS, API
  created_at?: string;
  last_activity_at?: string;
  is_active?: boolean;
};
```

#### User Type
```typescript
/**
 * User represents a registered user in the system
 * Maps to 'users' table in public schema
 */
export type User = {
  user_id?: string;
  mobile_number?: string;      // Required, encrypted
  email?: string;               // Optional, encrypted
  status?: UserStatus;          // ACTIVE, SUSPENDED, LOCKED, DELETED, PENDING_VERIFICATION
  created_at?: string;
  updated_at?: string;
  last_login_at?: string;
  last_login_session_type?: string;
  preferred_auth_method?: string;
  created_by?: string;
  updated_by?: string;
  login_attempts?: number;
  locked_until?: string;        // Account lockout timestamp
  deleted_at?: string;
  username?: string;
  preferred_language?: string;
  notification_preference?: string;
  wallet_payment_method?: string;
  user_type?: UserType;         // Determines auth method (B2C, AGENT, BUSINESS, SYSTEM, etc.)
  email_verified?: boolean;
  email_verified_at?: string;
  email_login_attempts?: number;
  email_locked_until?: string;
  biometric_token_idx?: string; // HMAC-SHA256 blind index
  mobile_number_idx?: string;   // Encrypted index for mobile
  email_idx?: string;           // Encrypted index for email
  totp_enabled?: boolean;       // Two-factor authentication
  totp_secret_enc?: string;     // Encrypted TOTP secret
  active_policies_count?: number;
  pending_claims_count?: number;
  wallet_balance?: Money;
};
```

#### User Profile Type
```typescript
/**
 * UserProfile stores additional user information
 * Maps to 'user_profiles' table in public schema
 */
export type UserProfile = {
  user_id?: string;
  full_name?: string;
  date_of_birth?: string;
  gender?: AuthnGender;         // MALE, FEMALE, OTHER, UNSPECIFIED
  occupation?: string;
  address_line1?: string;
  address_line2?: string;
  city?: string;
  district?: string;
  division?: string;
  postal_code?: string;
  country?: string;
  profile_photo_url?: string;
  kyc_verified?: boolean;
  kyc_verified_at?: string;
  created_at?: string;
  updated_at?: string;
  marital_status?: string;
  employer?: string;
  permanent_address?: string;
  emergency_contact_name?: string;
  emergency_contact_number?: string;
  id_type?: string;
  id_upload_front_url?: string;
  id_upload_back_url?: string;
  photograph_selfie_url?: string;
  proof_of_address_url?: string;
  consent_privacy_acceptance?: boolean;
};
```

### 2.3 API Key & Access Control

```typescript
/**
 * API Key entity for insurer/partner authentication
 */
export type ApiKey = {
  id: string;
  name: string;
  owner_type: ApiKeyOwnerType;  // INSURER, PARTNER, etc.
  owner_id: string;
  scopes?: Array<string>;       // e.g., ["auth:read", "session:write"]
  status: ApiKeyStatus;         // ACTIVE, EXPIRED, REVOKED, SUSPENDED, ROTATING
  rate_limit_per_minute: number;
  expires_at?: string;
  last_used_at?: string;
  ip_whitelist?: Array<string>;
  audit_info: AuditInfo;
};

/**
 * API usage tracking (FR-207)
 */
export type ApiKeyUsage = {
  id: string;
  api_key_id: string;
  endpoint: string;
  http_method: string;
  status_code: number;
  response_time_ms?: number;
  request_ip?: string;
  user_agent?: string;
  request_payload?: string;
  response_payload?: string;
  trace_id?: string;
  timestamp: string;
};
```

## 3. B2B Services - Key Message Types

### 3.1 OTP (One-Time Password) for Multi-Channel Verification

```typescript
/**
 * One-Time Password for verification
 */
export type Otp = {
  otp_id?: string;
  user_id?: string;
  otp_hash?: string;            // Hashed for security
  purpose?: string;             // LOGIN, REGISTRATION, PASSWORD_RESET
  device_type?: string;
  ip_address?: string;
  created_at?: string;
  expires_at?: string;
  attempts?: number;
  verified?: boolean;
  verified_at?: string;
  
  // BTRC Compliance Fields - SMS Delivery Tracking
  provider_message_id?: string; // SMS provider tracking
  dlr_status?: string;          // Delivery report status
  sender_id?: string;
  carrier?: string;
  channel?: string;             // SMS, EMAIL, VOICE
  recipient?: string;
  dlr_received_at?: string;
  dlr_error_code?: string;
  dlr_updated_at?: string;
};
```

### 3.2 Document Management

```typescript
/**
 * DocumentType represents a catalog of required/allowed document types
 */
export type AuthnDocumentType = {
  document_type_id?: string;
  code?: string;                // e.g., "nid", "passport", "license"
  name?: string;
  description?: string;
  is_active?: boolean;
  created_at?: string;
  updated_at?: string;
};

/**
 * UserDocument represents documents uploaded by a user for KYC or Policies
 */
export type UserDocument = {
  user_document_id?: string;
  user_id?: string;
  document_type_id?: string;
  policy_id?: string;           // Optional: for policy-specific docs
  file_url?: string;
  verification_status?: string; // PENDING, VERIFIED, REJECTED
  verified_by?: string;
  verified_at?: string;
  created_at?: string;
  updated_at?: string;
};
```

## 4. Service Client Definitions

The SDK is auto-generated from OpenAPI/Protobuf specs with the following service categories:

### 4.1 Authentication Service (AuthService)
- `EmailLoginData` / `EmailLoginResponses` - Email OTP login
- `BiometricAuthenticateData` / `BiometricAuthenticateResponses` - Biometric auth
- `ChangePasswordData` / `ChangePasswordResponses` - Password management
- `CreateUserProfileData` / `CreateUserProfileResponses` - Profile creation
- `CreateApiKeyData` / `CreateApiKeyResponses` - API key generation
- `DisableTotpData` / `DisableTotpResponses` - 2FA management
- `EnableTotpData` / `EnableTotpResponses` - 2FA setup
- `ApproveKycData` / `ApproveKycResponses` - KYC approval workflow
- `InitiateKycSessionData` / `InitiateKycSessionResponses` - KYC initiation
- `CompleteKycSessionData` / `CompleteKycSessionResponses` - KYC completion
- `GetCurrentSessionData` / `GetCurrentSessionResponses` - Session info
- `GetJwksData` / `GetJwksResponses` - JWT key set retrieval

### 4.2 API Key Service (ApiKeyService)
- `GenerateApiKeyData` / `GenerateApiKeyResponses` - Create new API key
- `GetApiKeyData` / `GetApiKeyResponses` - Retrieve key details
- `ListApiKeysData` / `ListApiKeysResponses` - List keys with pagination
- `RevokeApiKeyData` / `RevokeApiKeyResponses` - Revoke key
- `RotateApiKeyData` / `RotateApiKeyResponses` - Key rotation
- `GetUsageStatsData` / `GetUsageStatsResponses` - Usage analytics

### 4.3 Audit Service (AuditService)
- `CreateAuditLogData` / `CreateAuditLogResponses` - Log operations
- `GetAuditLogsData` / `GetAuditLogsResponses` - Retrieve audit logs
- `GetAuditTrailData` / `GetAuditTrailResponses` - Full audit trail
- `CreateComplianceLogData` / `CreateComplianceLogResponses` - Compliance tracking
- `GetComplianceLogsData` / `GetComplianceLogsResponses` - Compliance reports
- `GenerateComplianceReportData` / `GenerateComplianceReportResponses` - Regulatory reports

### 4.4 Analytics Service (AnalyticsService)
- `CreateDashboardData` / `CreateDashboardResponses` - Dashboard creation
- `GetDashboardData` / `GetDashboardResponses` - Retrieve dashboard
- `GenerateReportData` / `GenerateReportResponses` - Report generation
- `GetMetricsData` / `GetMetricsResponses` - Metrics retrieval
- `RunQueryData` / `RunQueryResponses` - Custom SQL queries
- `ScheduleReportData` / `ScheduleReportResponses` - Scheduled reports

## 5. Message Types for Org/Employee/User Management

### 5.1 Organization Management (B2B)
**Note:** Full organization types not found in current types.gen.ts. Likely in separate b2b service definitions.

### 5.2 User & Role Management

#### User Status Enumeration
```typescript
export type UserStatus =
  | 'USER_STATUS_UNSPECIFIED'
  | 'USER_STATUS_PENDING_VERIFICATION'
  | 'USER_STATUS_ACTIVE'
  | 'USER_STATUS_SUSPENDED'
  | 'USER_STATUS_LOCKED'
  | 'USER_STATUS_DELETED';
```

#### User Type Enumeration (Portal & Auth Method Routing)
```typescript
/**
 * UserType controls authentication method routing and portal assignment.
 * Portal mapping:
 * - B2C_CUSTOMER:         portal: "b2c"        (mobile SMS OTP + JWT)
 * - AGENT:                portal: "agent"      (mobile SMS OTP + JWT)
 * - BUSINESS_BENEFICIARY: portal: "business"   (email OTP + web server-side session)
 * - SYSTEM_USER:          portal: "system"     (email OTP + web server-side session)
 * - PARTNER:              portal: "b2b"        (email OTP + web server-side session)
 * - REGULATOR:            portal: "regulator"  (email OTP + web server-side session, read-only)
 * - BUSINESS_ADMIN:       Business admin portal
 * - B2B_ORG_ADMIN:        B2B organization admin
 */
export type UserType =
  | 'USER_TYPE_UNSPECIFIED'
  | 'USER_TYPE_B2C_CUSTOMER'
  | 'USER_TYPE_AGENT'
  | 'USER_TYPE_BUSINESS_BENEFICIARY'
  | 'USER_TYPE_SYSTEM_USER'
  | 'USER_TYPE_PARTNER'
  | 'USER_TYPE_REGULATOR'
  | 'USER_TYPE_BUSINESS_ADMIN'
  | 'USER_TYPE_B2B_ORG_ADMIN';
```

#### Session Type Enumeration
```typescript
export type SessionType =
  | 'SESSION_TYPE_UNSPECIFIED'
  | 'SESSION_TYPE_SERVER_SIDE'   // 12-hour expiry (web portals)
  | 'SESSION_TYPE_JWT';          // 7-day expiry (mobile apps)
```

#### Device Type Enumeration
```typescript
/**
 * Device type enum - auto-maps to session type
 */
export type AuthnDeviceType =
  | 'DEVICE_TYPE_UNSPECIFIED'
  | 'DEVICE_TYPE_WEB'
  | 'DEVICE_TYPE_MOBILE_ANDROID'
  | 'DEVICE_TYPE_MOBILE_IOS'
  | 'DEVICE_TYPE_API';
```

#### Gender Enumeration
```typescript
export type AuthnGender =
  | 'GENDER_UNSPECIFIED'
  | 'GENDER_MALE'
  | 'GENDER_FEMALE'
  | 'GENDER_OTHER';
```

### 5.3 Audit & Access Control

#### Audit Action Enumeration
```typescript
/**
 * Available actions for audit. Actions tracked for compliance.
 */
export type AuditAction =
  | 'AUDIT_ACTION_UNSPECIFIED'
  | 'AUDIT_ACTION_CREATE'
  | 'AUDIT_ACTION_READ'
  | 'AUDIT_ACTION_UPDATE'
  | 'AUDIT_ACTION_DELETE'
  | 'AUDIT_ACTION_LOGIN'
  | 'AUDIT_ACTION_LOGOUT'
  | 'AUDIT_ACTION_APPROVE'
  | 'AUDIT_ACTION_REJECT'
  | 'AUDIT_ACTION_EXPORT';
```

#### Audit Log Type
```typescript
/**
 * Audit log (FR-153 to FR-158)
 */
export type AuditLog = {
  audit_log_id?: string;
  entity_type?: string;        // Resource type (User, Policy, Claim, etc.)
  entity_id?: string;          // Resource ID
  action?: AuditAction;        // Operation performed
  user_id?: string;            // Actor
  user_email?: string;
  user_role?: string;          // Actor's role
  old_values?: string;         // JSON: before values
  new_values?: string;         // JSON: after values
  changes?: string;            // JSON: specific changes
  ip_address?: string;
  user_agent?: string;
  trace_id?: string;           // Distributed tracing ID
  timestamp?: string;
};
```

#### API Key Status Enumeration
```typescript
export type ApiKeyStatus =
  | 'API_KEY_STATUS_UNSPECIFIED'
  | 'API_KEY_STATUS_ACTIVE'
  | 'API_KEY_STATUS_EXPIRED'
  | 'API_KEY_STATUS_REVOKED'
  | 'API_KEY_STATUS_SUSPENDED'
  | 'API_KEY_STATUS_ROTATING';
```

#### Approval Decision & Status
```typescript
/**
 * Enumeration of approval decision values
 */
export type ApprovalDecision =
  | 'APPROVAL_DECISION_UNSPECIFIED'
  | 'APPROVAL_DECISION_PENDING'
  | 'APPROVAL_DECISION_APPROVED'
  | 'APPROVAL_DECISION_REJECTED'
  | 'APPROVAL_DECISION_NEEDS_MORE_INFO';

/**
 * Status values for approval
 */
export type ApprovalStatus =
  | 'APPROVAL_STATUS_UNSPECIFIED'
  | 'APPROVAL_STATUS_PENDING'
  | 'APPROVAL_STATUS_APPROVED'
  | 'APPROVAL_STATUS_REJECTED'
  | 'APPROVAL_STATUS_CANCELLED';
```

#### Verification Types
```typescript
/**
 * Available methods for verification
 */
export type VerificationMethod =
  | 'VERIFICATION_METHOD_UNSPECIFIED'
  | 'VERIFICATION_METHOD_PORICHOY'   // Bangladesh national ID verification
  | 'VERIFICATION_METHOD_NID'        // National ID
  | 'VERIFICATION_METHOD_PASSPORT'
  | 'VERIFICATION_METHOD_MANUAL'
  | 'VERIFICATION_METHOD_TRADE_LICENSE';

/**
 * Type categorization for verification
 */
export type VerificationType =
  | 'VERIFICATION_TYPE_UNSPECIFIED'
  | 'VERIFICATION_TYPE_KYC'          // Know Your Customer
  | 'VERIFICATION_TYPE_KYB';         // Know Your Business

/**
 * Status values for verification
 */
export type VerificationStatus =
  | 'VERIFICATION_STATUS_UNSPECIFIED'
  | 'VERIFICATION_STATUS_PENDING'
  | 'VERIFICATION_STATUS_IN_PROGRESS'
  | 'VERIFICATION_STATUS_VERIFIED'
  | 'VERIFICATION_STATUS_REJECTED'
  | 'VERIFICATION_STATUS_EXPIRED';
```

## 6. Role & Permission Enum Definitions

### 6.1 Authentication Types
```typescript
/**
 * Type categorization for authentication
 */
export type AuthenticationType =
  | 'AUTHENTICATION_TYPE_UNSPECIFIED'
  | 'AUTHENTICATION_TYPE_API_KEY'
  | 'AUTHENTICATION_TYPE_OAUTH2'
  | 'AUTHENTICATION_TYPE_BASIC_AUTH'
  | 'AUTHENTICATION_TYPE_MUTUAL_TLS';
```

### 6.2 Compliance & Regulatory

#### Compliance Type
```typescript
export type ComplianceType =
  | 'COMPLIANCE_TYPE_UNSPECIFIED'
  | 'COMPLIANCE_TYPE_AUDIT'
  | 'COMPLIANCE_TYPE_REGULATION'
  | 'COMPLIANCE_TYPE_POLICY';
```

#### Compliance Status
```typescript
export type ComplianceStatus =
  | 'COMPLIANCE_STATUS_UNSPECIFIED'
  | 'COMPLIANCE_STATUS_PENDING'
  | 'COMPLIANCE_STATUS_COMPLIANT'
  | 'COMPLIANCE_STATUS_NON_COMPLIANT'
  | 'COMPLIANCE_STATUS_NEEDS_REVIEW';
```

### 6.3 Related Request/Response Types

#### API Key Generation Request
```typescript
/**
 * Request payload for api key generation operation
 */
export type ApiKeyGenerationRequest = {
  name: string;
  owner_type?: string;           // INSURER, PARTNER
  owner_id: string;
  scopes?: Array<string>;        // e.g., ["auth:read", "session:write"]
  rate_limit_per_minute?: number;
  expires_in_days?: string;
  ip_whitelist?: Array<string>;
};

/**
 * Response payload for api key generation operation
 */
export type ApiKeyGenerationResponse = {
  api_key_id?: string;
  api_key?: string;              // Plain text - shown only once!
  message?: string;
  error?: Error;
};
```

#### User Creation Request
```typescript
/**
 * Request payload for user profile creation
 */
export type ApiKeyCreationRequest = {
  name: string;                  // Human-readable name
  owner_id: string;              // User or service ID
  owner_type?: string;           // USER, SERVICE
  scopes?: Array<string>;        // e.g., ["auth:read", "session:write"]
  rate_limit_per_minute?: number; // 0 = unlimited
  expires_at?: string;           // null = never expires
};
```

## 7. SDK Initialization & Configuration

### 7.1 Client Options
```typescript
export type ClientOptions = {
  baseUrl:
    | 'https://api.labaidinsuretech.com'
    | 'https://staging-api.labaidinsuretech.com'
    | (string & {});
};
```

### 7.2 Client Usage Pattern
```typescript
// Generated client.gen.ts creates a pre-configured client:
export const client = createClient(
  createConfig<ClientOptions2>({ 
    baseUrl: 'https://api.labaidinsuretech.com' 
  })
);
```

## 8. Key Features Summary

### Authentication & Authorization
- **Hybrid Session Model**: Server-side (web) + JWT (mobile)
- **Multi-method Auth**: Email OTP (business users), SMS OTP (customers), Biometric
- **Portal Routing**: 8 user types mapped to different portals
- **API Key Management**: Scoped access, rate limiting, IP whitelisting
- **2FA Support**: TOTP-based two-factor authentication

### Compliance & Audit
- **Comprehensive Audit Logging**: All CRUD operations + auth events tracked
- **Compliance Reporting**: Regulatory compliance log types
- **Verification Methods**: KYC/KYB with multiple ID types (NID, Passport, etc.)
- **Access Control**: Fine-grained action-based auditing

### Data Security
- **Encryption**: Mobile number, email, biometric tokens (encrypted indices)
- **Blind Indexing**: HMAC-SHA256 indices for encrypted field lookups
- **SMS Compliance**: BTRC delivery tracking for SMS OTP
- **Distributed Tracing**: Trace IDs for request correlation

### B2B Features
- **Multi-tenant Support**: Organization/Partner-scoped API keys
- **Role-based Access**: User roles determine API scopes
- **Portal Assignment**: Dynamic portal routing based on user type
- **Usage Analytics**: Track API usage per key, endpoint, status code

## 9. File Size & Generation Tool

- **types.gen.ts**: 39,777 lines (comprehensive type definitions)
- **sdk.gen.ts**: 7,316 lines (service method definitions)
- **Generator**: @hey-api/openapi-ts (auto-generates from OpenAPI/Protobuf specs)
- **Language**: TypeScript with strict typing
- **Package**: @lifeplus/insuretech-sdk v0.1.0

## 10. Build & Distribution

```json
{
  "main": "./dist/index.js",
  "module": "./dist/index.mjs",
  "types": "./dist/index.d.ts",
  "exports": {
    ".": {
      "types": "./dist/index.d.ts",
      "import": "./dist/index.mjs",
      "require": "./dist/index.js"
    }
  }
}
```

**Build Tools**: tsup, vitest, eslint, prettier, typescript v5.3.3
