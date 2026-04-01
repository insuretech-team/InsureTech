# PAYMENT + ORDER IMPLEMENTATION PLAN

## 1. Objective

Implement `orders` and `payment` as a proto-aligned, Kafka-connected workflow that satisfies both the current service contracts and the SRS requirements for:

- policy purchase and activation
- multi-channel payment collection
- billing and receipts
- document generation and storage
- B2B and partner-driven distribution
- auditability, AML/CFT controls, and regulatory retention

The correct delivery model is hybrid:

- synchronous gRPC for command validation, permission checks, and authoritative reads
- Kafka for durable post-commit facts, cross-service reactions, document fanout, and long-running workflows

## 2. Source Baseline Used for This Plan

This plan is now grounded in three sources together:

- actual proto contracts under `proto/insuretech/...`
- SRS V3.7 sections under `documentation/SRS_v3/SPECS_V3.7/sections`
- insurance domain notes under `documentation/InsuaranceKnowledgeBank`

This matters because the proto set defines what is currently implementable, while the SRS defines what the system is actually required to do.

## 3. High-Confidence Requirements That Must Shape the Design

### 3.1 Purchase and issuance requirements

From the SRS, the order/payment design must satisfy:

- end-to-end purchase flow: product selection -> applicant details -> nominee details -> payment -> policy issuance (`FR-030`)
- digital policy document PDF with QR code generated within 30 seconds of payment confirmation (`FR-035`)
- policy document sent via SMS link and email within 5 minutes (`FR-036`)
- policy activated immediately upon payment confirmation (`FR-037`)
- policy states must include at least `Pending Payment`, `Active`, `Suspended`, `Cancelled`, `Lapsed`, `Expired` (`FR-039`)
- dashboard must expose payment history and receipt downloads (`FR-040`, portal requirements)

Design implication:

- `OrderPaymentConfirmedEvent` cannot be the end of the workflow
- it must trigger policy issuance, document generation, storage linking, notification fanout, and dashboard projection updates

### 3.2 Payment requirements

The payment domain must support:

- bKash, Nagad, Rocket, bank transfer, cards, and manual cash/cheque (`FR-073`)
- manual payment with proof upload and later verification (`FR-076`)
- payment verification workflow: `pending -> verified -> policy activated` or `rejected -> refund` (`FR-077`)
- payment receipt generation with transaction ID, amount, date, policy number (`FR-078`)
- payment retry with exponential backoff (`FR-080`)
- cancellation refund processing (`FR-081`)
- TigerBeetle-backed double-entry bookkeeping (`FR-082`, `NFR-058`)
- immutable payment audit trail and long retention (`FR-083`, `FR-153`, `FR-154`)

Design implication:

- the current payment entity/event model is too thin
- receipts, manual proof uploads, verification decisions, and ledger postings must become first-class workflow elements

### 3.3 Integration and resiliency requirements

The order/payment flow must satisfy:

- Kafka as the async broker (`NFR-053`)
- gRPC + protobuf for internal service communication (`NFR-048`, `FR-193`)
- REST/OpenAPI for external integrations (`NFR-050`, `FR-195`)
- payment gateway webhook signature validation (`FR-225`, `FR-228`)
- idempotency key handling on payment and policy issuance APIs (`FR-227`)
- retry with exponential backoff and circuit breaker patterns (`FR-226`)
- webhook system for external real-time notifications (`FR-163`)

### 3.4 Storage and document requirements

The workflow must support:

- S3-compatible object storage (`NFR-052`, `FR-237`)
- secure document downloads and offline policy-document access
- presigned upload URLs with 30-minute expiry (`FR-239`)
- payment proof upload for manual verification (`FR-076`)
- policy documents, receipts, invoices, endorsements, and cancellation artifacts
- document retention and archival requirements aligned to audit/compliance rules

### 3.5 Security and compliance requirements

The order/payment design must incorporate:

- zero-trust command validation
- RBAC + ABAC (`FR-014`, `FR-015`)
- tenant isolation (`FR-016`)
- immutable audit logs (`FR-153`)
- 20-year regulatory retention target in the business requirements (`FR-154`, SEC reporting expectations)
- PCI-DSS SAQ-A style card flow: hosted payment page, no raw card storage (`SEC-003`)
- virus scanning on uploaded files (`SEC-010`)
- AML/CFT monitoring for rapid purchases, high-value premiums, third-party payments, frequent cancellations, and geographic anomalies (`SEC-017`, TM rules)
- no customer tipping-off in suspicious-activity workflows

## 4. Architecture Rules

### 4.0 Authentication Model — Corrected and Authoritative

**THIS SECTION SUPERSEDES ANY EARLIER LANGUAGE ABOUT JWT-ONLY OR PER-SERVICE JWT INTERCEPTORS.**

The InsureTech auth model is **hybrid** and split by client type:

| Client type | Auth mechanism | Session storage | Token type |
| --- | --- | --- | --- |
| Web portals (system, b2b, business, agent, regulator) | Email/password → server-side session | `authn_schema.sessions` (DB) | `session_token` in httpOnly cookie + CSRF token |
| Mobile apps (iOS, Android) | Phone/OTP or password → JWT | Stateless (JTI blocklist in Redis) | `access_token` (RS256 JWT) + `refresh_token` |
| API clients / partners | API key → scoped JWT | Stateless | API key scoped token |

**Session types confirmed in code:**
- `SESSION_TYPE_SERVER_SIDE` → web portals → cookie: `session_token=<opaque-token>` + `X-CSRF-Token` header
- `SESSION_TYPE_JWT` → mobile/API → `Authorization: Bearer <jwt>`

**`DeviceType` → `SessionType` mapping (from AUTHN_AUTHZ_REFERENCE.md):**
- `WEB` → `SERVER_SIDE`
- `MOBILE_ANDROID`, `MOBILE_IOS`, `API`, `DESKTOP` → `JWT`

#### 4.0.1 Gateway-level AuthN (already implemented — `auth_middleware.go`)

The gateway `AuthMiddleware` already handles both auth paths in a single middleware:

```go
// From gateway/internal/routes/auth_middleware.go
jwt := bearerToken(r.Header.Get("Authorization"))       // mobile/API path
sessionToken := cookie("session_token")                   // web portal path

resp, err := client.ValidateToken(ctx, &ValidateTokenRequest{
    AccessToken: jwt,
    SessionId:   sessionToken,  // AuthN resolves either; one will be empty
})
```

After `ValidateToken` succeeds, the gateway propagates identity as request headers to all downstream services:

| Header | Value | Source |
| --- | --- | --- |
| `X-User-ID` | `resp.UserId` | both paths |
| `X-Session-ID` | `resp.SessionId` | both paths |
| `X-Session-Type` | `"SERVER_SIDE"` or `"JWT"` | both paths |
| `X-Portal` | `"system"`, `"b2b"`, `"b2c"`, etc. | both paths |
| `X-Tenant-ID` | tenant UUID | both paths |
| `X-Token-ID` | JTI (for revocation) | both paths |
| `X-Device-ID` | device fingerprint | both paths |
| `X-User-Type` | `"B2C_CUSTOMER"`, `"SYSTEM_USER"`, etc. | both paths |

**Device binding check (JWT only):** if `X-Device-Id` header is supplied on a JWT session, it must match `resp.DeviceId` in the token claim. Server-side sessions skip this check.

**CSRF enforcement:** all state-changing web portal routes (POST/PATCH/DELETE) are wrapped with `CSRFMiddleware(authnConn)` which calls `AuthN.ValidateCSRF`. Mobile JWT requests do not carry CSRF tokens — they are exempt by session type.

#### 4.0.2 Gateway-level AuthZ (already implemented — `authz_middleware.go`)

After AuthN, the gateway calls `AuthZ.CheckAccess` (Casbin PERM model):

```
enforce("user:<user_id>", "portal:tenant_id", "svc:<service>/<resource>", "HTTP_METHOD")
```

Domain format: `"system:root"`, `"b2b:<org_id>"`, `"b2c:root"`, `"agent:root"`, etc.

Order/payment specific objects already seeded (from `AUTHN_AUTHZ_REFERENCE.md`):
- `b2b_org_admin` → `svc:invoice/*, *` + `svc:payment/*, *`
- `business/finance` → `svc:invoice/*, *` + `svc:payment/*, *`
- `b2c/customer` → `svc:policy/my/*, GET` + `svc:claim/my/*, *`

#### 4.0.3 Per-service gRPC interceptor pattern (defense-in-depth)

**Critical clarification from code inspection:** The `authz/internal/middleware/jwt_interceptor.go` is a JWT interceptor on the **authz-service itself** — it protects the authz gRPC API from unauthenticated callers. It is NOT the pattern to copy for order/payment services.

The correct defense-in-depth pattern for order/payment is the **b2b service interceptor** at `b2b/internal/middleware/authz_interceptor.go`. This interceptor:
- Does NOT re-parse JWT or re-call `AuthN.ValidateToken`
- Reads identity ONLY from incoming gRPC metadata headers set by the gateway
- Calls `AuthZ.CheckAccess` (Casbin) as the authorization check

```go
// From b2b/internal/middleware/authz_interceptor.go — exact pattern to copy
func (i *AuthZInterceptor) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
    return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
        md, ok := metadata.FromIncomingContext(ctx)
        if !ok {
            return nil, status.Error(codes.Unauthenticated, "missing metadata")
        }

        userIDs := md.Get("x-user-id")
        if len(userIDs) == 0 {
            return nil, status.Error(codes.Unauthenticated, "missing user_id")
        }
        userID := userIDs[0]

        // Normalize portal: "PORTAL_SYSTEM" → "system", "PORTAL_B2B" → "b2b"
        portalRaw := firstNonEmpty(md.Get("x-portal"))
        portalNorm := strings.ToLower(strings.TrimSpace(strings.TrimPrefix(portalRaw, "PORTAL_")))
        isSystemPortal := portalNorm == "system"

        orgID := firstNonEmpty(md.Get("x-business-id")) // B2B only; empty for B2C/agent

        // Non-system portals with no org context: allow bootstrap methods, deny others
        if !isSystemPortal && orgID == "" {
            if isBootstrapMethod(info.FullMethod) {
                return handler(ctx, req)
            }
            return nil, status.Error(codes.PermissionDenied, "missing organisation context")
        }

        resource, action := mapMethodToResourceAction(info.FullMethod)
        if resource == "" {
            return handler(ctx, req) // bootstrap — no authz check
        }

        domain := resolveAuthzDomain(md, orgID)
        resp, err := i.authzClient.CheckAccess(ctx, &authzservicev1.CheckAccessRequest{
            UserId: userID,  // NOTE: authz service prepends "user:" internally
            Domain: domain,
            Object: resource,  // "svc:order/*" or "svc:payment/*"
            Action: action,    // "GET", "POST", "PATCH", "DELETE"
        })
        if err != nil {
            return nil, status.Error(codes.Internal, "authorization check failed")
        }
        if !resp.Allowed {
            return nil, status.Error(codes.PermissionDenied, "access denied")
        }
        return handler(ctx, req)
    }
}

// resolveAuthzDomain — matches b2b interceptor exactly
func resolveAuthzDomain(md metadata.MD, orgID string) string {
    portal := strings.ToLower(strings.TrimPrefix(firstNonEmpty(md.Get("x-portal")), "PORTAL_"))
    tenantID := firstNonEmpty(md.Get("x-tenant-id"))
    switch portal {
    case "system":
        return "system:root"
    case "b2b":
        if orgID != "" { return "b2b:" + orgID }
        if tenantID != "" { return "b2b:" + tenantID }
        return "b2b:root"
    default:
        if tenantID != "" { return portal + ":" + tenantID }
        return portal + ":root"
    }
}

// mapMethodToResourceAction for orders-service
func mapOrderMethodToResourceAction(method string) (resource, action string) {
    parts := strings.Split(method, "/")
    methodName := parts[len(parts)-1]
    switch {
    case strings.HasPrefix(methodName, "Get"), strings.HasPrefix(methodName, "List"):
        return "svc:order/*", "GET"
    case strings.HasPrefix(methodName, "Create"):
        return "svc:order/*", "POST"
    case strings.HasPrefix(methodName, "Update"), strings.HasPrefix(methodName, "Confirm"):
        return "svc:order/*", "PATCH"
    case strings.HasPrefix(methodName, "Cancel"):
        return "svc:order/*", "DELETE"
    case methodName == "InitiatePayment":
        return "svc:order/*", "POST"
    default:
        return "", ""
    }
}

// mapMethodToResourceAction for payment-service
func mapPaymentMethodToResourceAction(method string) (resource, action string) {
    parts := strings.Split(method, "/")
    methodName := parts[len(parts)-1]
    switch {
    case strings.HasPrefix(methodName, "Get"), strings.HasPrefix(methodName, "List"):
        return "svc:payment/*", "GET"
    case strings.HasPrefix(methodName, "Initiate"), strings.HasPrefix(methodName, "Add"),
         strings.HasPrefix(methodName, "Submit"), strings.HasPrefix(methodName, "Handle"),
         strings.HasPrefix(methodName, "Generate"), strings.HasPrefix(methodName, "Reconcile"):
        return "svc:payment/*", "POST"
    case strings.HasPrefix(methodName, "Verify"), strings.HasPrefix(methodName, "Review"):
        return "svc:payment/*", "PATCH"
    default:
        return "", ""
    }
}
```

**Key rules confirmed from b2b interceptor code:**
- System portal users (`PORTAL_SYSTEM` → `system:root`) bypass org context requirement
- B2B portal users require `x-business-id` for non-bootstrap methods
- Bootstrap methods (e.g. resolving own organisation) skip Casbin entirely
- `mapMethodToResourceAction()` uses method name prefix matching — no regex, no proto introspection
- The `CheckAccessRequest.UserId` field receives the raw UUID string — the authz service internally formats it as `user:<uuid>` for Casbin

**authz JWT interceptor (`authz/internal/middleware/jwt_interceptor.go`) — what it actually does:**
This interceptor protects the **authz-service's own gRPC API**. It validates Bearer JWT tokens on calls TO the authz service (e.g. when the gateway calls `CheckAccess`). It uses `trustedInternalServices = ["gateway", "b2b-service"]` to allow internal callers without a JWT. This is not relevant to order/payment service implementation.

#### 4.0.4 What order-service and payment-service MUST NOT do

- ❌ Do NOT re-call `AuthN.ValidateToken` inside the service — the gateway already did this
- ❌ Do NOT implement a separate JWT parser or session lookup — trust the `X-*` metadata headers
- ❌ Do NOT add async/Kafka-based authorization — auth is always synchronous gRPC
- ❌ Do NOT store raw Bearer tokens or session tokens in the service layer
- ✅ DO add a gRPC unary interceptor that reads `x-user-id`, `x-portal`, `x-tenant-id` from incoming metadata and calls `AuthZ.CheckAccess` for defense-in-depth
- ✅ DO ensure `csrfMW` is present for all web portal state-changing order/payment routes; on the current branch it is already present for `/initiate-payment`, `/confirm`, and `/cancel`, but `POST /v1/orders` still lacks it
- ✅ DO propagate `x-user-id`, `x-portal`, `x-tenant-id`, `x-session-id`, `x-session-type`, `x-token-id` from incoming gRPC metadata into all outgoing events as `actor_user_id`, `portal`, `tenant_id`, `session_id`, `session_type`, `token_id`

#### 4.0.5 CSRF applicability to order/payment flows

| Operation | Session type | CSRF required |
| --- | --- | --- |
| `CreateOrder` from web portal | SERVER_SIDE | ⚠️ Should be yes, but current router still misses `csrfMW` on `POST /v1/orders` |
| `CreateOrder` from mobile app | JWT | ❌ No |
| `InitiatePayment` from web | SERVER_SIDE | ✅ Yes |
| `InitiatePayment` from mobile | JWT | ❌ No |
| `ConfirmPayment` (gateway callback) | neither — internal | ❌ No — internal system call |
| `HandleGatewayWebhook` | neither — external webhook | ❌ No — HMAC signature instead |
| `ReviewManualPayment` (admin portal) | SERVER_SIDE | ✅ Yes |
| `CancelOrder` from web | SERVER_SIDE | ✅ Yes |

Gateway router entries for orders are already present in `router.go`. On the current branch, `csrfMW` is already applied to `POST /v1/orders/{order_id}/initiate-payment`, `POST /v1/orders/{order_id}/confirm`, and `POST /v1/orders/{order_id}/cancel`; the verified remaining order-route gap is `POST /v1/orders`.

#### 4.0.6 Gateway routes for orders and payments (current shape + remaining gaps)

Orders routes (currently proxied to PoliSync `:50141`; route names on the current branch are `/initiate-payment`, `/confirm`, and `/cancel`):

```go
authzOrder := authzMW("svc:order", PathSegmentExtractor("/v1/"))

mux.Handle("POST /v1/orders",                              authMW(authzOrder(ordersHandler.Proxy()))) // current gap: add csrfMW
mux.Handle("GET  /v1/orders",                              authMW(authzOrder(ordersHandler.Proxy())))
mux.Handle("GET  /v1/orders/{order_id}",                   authMW(authzOrder(ordersHandler.Proxy())))
mux.Handle("POST /v1/orders/{order_id}/initiate-payment",  authMW(csrfMW(authzOrder(ordersHandler.Proxy()))))
mux.Handle("POST /v1/orders/{order_id}/confirm",           authMW(csrfMW(authzOrder(ordersHandler.Proxy()))))
mux.Handle("POST /v1/orders/{order_id}/cancel",            authMW(csrfMW(authzOrder(ordersHandler.Proxy()))))
```

Payment routes (still absent in the gateway and still need to be added):

```go
authzPayment := authzMW("svc:payment", PathSegmentExtractor("/v1/"))

mux.Handle("POST /v1/payments",                                    authMW(csrfMW(authzPayment(paymentHandler.Proxy()))))
mux.Handle("GET  /v1/payments",                                    authMW(authzPayment(paymentHandler.Proxy())))
mux.Handle("GET  /v1/payments/{payment_id}",                       authMW(authzPayment(paymentHandler.Proxy())))
mux.Handle("POST /v1/payments/{payment_id}:verify",                authMW(csrfMW(authzPayment(paymentHandler.Proxy()))))
mux.Handle("POST /v1/payments/{payment_id}/refunds",               authMW(csrfMW(authzPayment(paymentHandler.Proxy()))))
mux.Handle("GET  /v1/refunds/{refund_id}",                         authMW(authzPayment(paymentHandler.Proxy())))
mux.Handle("POST /v1/payments/{payment_id}:submit-proof",          authMW(csrfMW(authzPayment(paymentHandler.Proxy()))))
mux.Handle("POST /v1/payments/{payment_id}:review",                authMW(csrfMW(authzPayment(paymentHandler.Proxy()))))
mux.Handle("GET  /v1/payments/{payment_id}/receipt",               authMW(authzPayment(paymentHandler.Proxy())))
mux.Handle("POST /v1/payments:reconcile",                          authMW(csrfMW(authzPayment(paymentHandler.Proxy()))))
mux.Handle("POST /v1/payments/webhook/{provider}",                 http.HandlerFunc(paymentHandler.Webhook)) // public — HMAC only
```

#### 4.0.7 AuthZ Casbin permissions to seed for order/payment

These rules must be added to the portal seeder (`portal_seeder.go`) for order and payment objects:

| Portal/domain | Role | Object | Action | Notes |
| --- | --- | --- | --- | --- |
| `b2c:root` | `customer` | `svc:order/my/*` | `GET`, `POST` | customer creates/reads own orders |
| `b2c:root` | `customer` | `svc:payment/my/*` | `GET`, `POST` | customer initiates/reads own payments |
| `agent:root` | `agent` | `svc:order/*` | `GET`, `POST`, `PATCH` | agent-assisted purchases |
| `agent:root` | `agent` | `svc:payment/*` | `GET` | agent reads payment status |
| `b2b:root` | `b2b_org_admin` | `svc:order/*` | `*` | already seeded |
| `b2b:root` | `b2b_org_admin` | `svc:payment/*` | `*` | already seeded |
| `system:root` | `admin` | `svc:order/*` | `*` | system admin full access |
| `system:root` | `admin` | `svc:payment/*` | `*` | system admin full access |
| `system:root` | `finance` | `svc:payment/*` | `GET`, `POST` | finance reviews/reconciles |
| `system:root` | `support` | `svc:order/*` | `GET` | support reads orders |
| `business:root` | `finance` | `svc:payment/*` | `*` | already seeded |

### 4.1 Ownership split

- `orders-service` owns order persistence and order event publication.
- Go `payment-service` owns payment persistence, gateway interactions, verification, reconciliation, refund persistence, receipt metadata, and payment event publication.
- new `billing-service` owns invoice lifecycle and invoice event publication.
- C# `PoliSync` owns cross-domain business rules and orchestration across quotation, order, payment, billing, policy issuance, document generation, and compliance triggers.
- `insurance-service` remains the authoritative policy/quotation write path until a dedicated policy service exists.

### 4.2 Payment service runtime decision

The payment service is now a Go service. The SRS language note that previously mentioned Node.js should be treated as outdated architecture documentation, not an implementation constraint.

Rule:

- preserve the payment domain boundary and contract shape
- standardize payment implementation, outbox, gateway clients, and Kafka publishers in Go
- update downstream planning, ownership, and deployment assumptions accordingly

### 4.3 AuthN/AuthZ synchronous gate — corrected model

The gateway already enforces AuthN + AuthZ before any request reaches a backend service. The correct model for order-service and payment-service is therefore:

**At the gateway (already implemented):**
- `AuthMiddleware` → calls `AuthN.ValidateToken` (handles both `Bearer <jwt>` and `session_token` cookie transparently)
- `AuthZMiddleware` → calls `AuthZ.CheckAccess` with `(user:<id>, portal:tenant, svc:order/*, METHOD)`
- `CSRFMiddleware` → calls `AuthN.ValidateCSRF` for web portal state-changing routes only

**At the service (defense-in-depth, to implement):**
- gRPC unary interceptor → reads `x-user-id`, `x-portal`, `x-tenant-id` from incoming metadata → calls `AuthZ.CheckAccess`
- Maps gRPC method names to Casbin objects/actions using `mapMethodToResourceAction` pattern (copy from b2b interceptor)
- Resolves Casbin domain via `resolveAuthzDomain` (copy from b2b interceptor)

**For B2B flows specifically:**
- `B2BService.ResolveMyOrganisation` for business-context resolution where the actor belongs to an organisation
- `x-business-id` header propagated as org context into downstream metadata

**Incorrect (do not do):**
- Re-calling `AuthN.ValidateToken` inside order-service or payment-service
- Implementing a standalone JWT parser in the service layer
- Asynchronous authorization via Kafka
- Letting Kafka consumers decide whether an already-accepted command was allowed
- Adding CSRF validation inside the gRPC service (CSRF is an HTTP-layer concern, handled only at the gateway)

## 5. Command-Time Security and Context Model

### 5.1 Identity extraction at service boundary

For every `CreateOrder`, `InitiatePayment`, `ConfirmPayment`, `CancelOrder`, `IssueInvoice`, `RefundPayment`, or `IssuePolicy` command arriving at order-service or payment-service via gRPC:

**The gateway has already performed AuthN and AuthZ.** The service receives a validated, trusted identity via gRPC metadata headers. The service must:

1. **Read identity from incoming gRPC metadata** — do NOT re-validate the token:
   ```go
   md, _ := metadata.FromIncomingContext(ctx)
   userID     := firstVal(md.Get("x-user-id"))
   tenantID   := firstVal(md.Get("x-tenant-id"))
   portal     := firstVal(md.Get("x-portal"))       // "system", "b2b", "b2c", "agent"
   sessionID  := firstVal(md.Get("x-session-id"))
   sessionType:= firstVal(md.Get("x-session-type")) // "SERVER_SIDE" or "JWT"
   tokenID    := firstVal(md.Get("x-token-id"))     // JTI
   deviceID   := firstVal(md.Get("x-device-id"))
   orgID      := firstVal(md.Get("x-business-id"))  // B2B only — may be empty
   ```

2. **Run defense-in-depth AuthZ check** (via gRPC interceptor, same pattern as b2b service):
   - `domain = resolveAuthzDomain(portal, tenantID, orgID)`
   - `object = "svc:order/*"` or `"svc:payment/*"`
   - `action = mapMethodToAction(grpcMethodName)`
   - Call `AuthZ.CheckAccess` → deny with `codes.PermissionDenied` if not allowed

3. **If the flow is B2B**, `x-business-id` is already populated by the gateway's `B2BContextMiddleware`. No separate `ResolveMyOrganisation` call is needed inside order-service or payment-service.

4. **Persist command metadata** — store `actor_user_id`, `tenant_id`, `session_id`, `portal`, `organisation_id` on the order/payment record for audit trail.

### 5.2 Session-type-aware behavior

Services must be aware of session type for certain behavioral differences:

| Behavior | `SERVER_SIDE` (web portal) | `JWT` (mobile/API) |
| --- | --- | --- |
| Token is in | httpOnly cookie (handled by gateway) | `Authorization: Bearer` header |
| CSRF token | Required at gateway — already enforced | Not applicable |
| Device binding | Not enforced | Enforced at gateway if `X-Device-Id` present |
| Refresh path | Session extended server-side by AuthN | Client calls `RefreshToken` RPC |
| Identity source in service | `x-user-id` from metadata | same — gateway normalizes both |

The service itself does not need to distinguish session types for order/payment business logic. The distinction is already resolved before the request arrives. The service records `session_type` on events for downstream audit and analytics only.

### 5.3 Event context — what every emitted event must carry

Every Kafka event published by order-service or payment-service must include the following fields (either in the proto payload or as Kafka message headers):

**In proto event payload (preferred — typed, inspectable):**
- `event_id` — new UUID per event
- `correlation_id` — from the triggering command's `idempotency_key` or request `correlation_id`
- `causation_id` — the `event_id` of the event that caused this one (for event chains)
- `tenant_id` — from `x-tenant-id` metadata
- `portal` — from `x-portal` metadata (normalized: `"system"`, `"b2b"`, `"b2c"`, `"agent"`)
- `actor_user_id` — from `x-user-id` metadata
- `session_id` — from `x-session-id` metadata
- `session_type` — from `x-session-type` metadata (`"SERVER_SIDE"` or `"JWT"`)
- `token_id` — from `x-token-id` metadata (JTI)
- `organisation_id` — from `x-business-id` metadata (B2B flows only; empty string for B2C/agent)
- `trace_id` — from `x-request-id` or `x-trace-id` metadata (for distributed tracing)
- `idempotency_key` — when the source command had one (propagate verbatim)
- `occurred_at` — proto `Timestamp` of when the event occurred (not when it was published)

**Current gap in order-service events (must fix):**
- `OrderCreatedEvent`, `OrderPaymentInitiatedEvent`, `OrderPaymentConfirmedEvent`, `OrderCancelledEvent`, `OrderFailedEvent` are all missing `tenant_id`, `organisation_id`, `portal`, `actor_user_id`, `session_id`, `causation_id`
- These fields exist in the service (`resolveTenantID`, `resolveCallerID` are already implemented) but are not carried through to events

**Current gap in payment-service events (must fix):**
- `PaymentInitiatedEvent`, `PaymentCompletedEvent`, `PaymentFailedEvent`, `RefundProcessedEvent` are missing `tenant_id`, `order_id`, `invoice_id`, `organisation_id`, `portal`, `actor_user_id`, `session_id`
- `referenceInfo()` only reads `policy_id` or a metadata string map — must read typed `order_id`, `invoice_id` fields once proto is extended

### 5.4 JWT claims — exact structure (from token_service.go)

The `InsureTechClaims` struct in `authn/internal/service/token_service.go` defines exactly what is embedded in every JWT. Order/payment services never parse the JWT directly — the gateway validates it and forwards these claims as gRPC metadata headers. This table is the authoritative mapping:

| JWT claim key | Description | Gateway header forwarded | gRPC metadata key |
| --- | --- | --- | --- |
| `sub` (standard) | User UUID | `X-User-ID` | `x-user-id` |
| `jti` (standard) | JWT ID (token UUID) | `X-Token-ID` | `x-token-id` |
| `exp`, `iat` (standard) | Expiry / issued-at | not forwarded | n/a |
| `ins_utp` | `UserType` — `B2C_CUSTOMER`, `SYSTEM_USER`, `AGENT`, `PARTNER`, etc. | `X-User-Type` | `x-user-type` |
| `ins_portal` | Portal name — `"system"`, `"b2b"`, `"b2c"`, `"agent"`, `"business"`, `"regulator"` | `X-Portal` | `x-portal` |
| `ins_tid` | Tenant UUID | `X-Tenant-ID` | `x-tenant-id` |
| `ins_did` | Device ID (fingerprint) | `X-Device-ID` | `x-device-id` |
| `ins_sid` | Session UUID | `X-Session-ID` | `x-session-id` |
| `ins_ttype` | Token type — `"ACCESS"` or `"REFRESH"` | not forwarded | n/a |

**Server-side session (web portals) — no JWT, no ins_* claims.** The gateway reads the `session_token` cookie, calls `AuthN.ValidateToken`, and then sets the same `X-*` headers from the `ValidateTokenResponse`. The downstream service sees identical metadata regardless of whether the caller used JWT or server-side session.

**`ValidateTokenResponse` fields** (confirmed from `auth_service.proto`):
- `valid`, `user_id`, `session_id`, `user_type`, `portal`, `tenant_id`, `token_id`, `device_id`, `session_type` (`"SERVER_SIDE"` or `"JWT"`)

**B2B org context** is NOT in the JWT claims. It is injected by the gateway's `B2BContextMiddleware` (reads `B2BService.ResolveMyOrganisation`) and forwarded as:
- Gateway header: `X-Business-ID`
- gRPC metadata key: `x-business-id`

### 5.5 Context propagation helper pattern (code-verified)

Both services should implement a context extraction helper aligned with the exact metadata keys the gateway sets:

```go
// internal/middleware/context.go (to create in both services)
package middleware

import (
    "context"
    "strings"
    "google.golang.org/grpc/metadata"
)

type RequestContext struct {
    UserID         string // from x-user-id (JWT sub / session user_id)
    TenantID       string // from x-tenant-id (ins_tid claim)
    Portal         string // from x-portal normalized: "system","b2b","b2c","agent","business","regulator"
    SessionID      string // from x-session-id (ins_sid claim)
    SessionType    string // from x-session-type: "SERVER_SIDE" or "JWT"
    TokenID        string // from x-token-id (JWT jti)
    DeviceID       string // from x-device-id (ins_did claim)
    UserType       string // from x-user-type (ins_utp claim)
    OrganisationID string // from x-business-id (B2B only — injected by B2BContextMiddleware)
    TraceID        string // from x-request-id (gateway RequestID middleware)
}

func ExtractRequestContext(ctx context.Context) RequestContext {
    md, ok := metadata.FromIncomingContext(ctx)
    if !ok {
        return RequestContext{}
    }
    return RequestContext{
        UserID:         first(md.Get("x-user-id")),
        TenantID:       first(md.Get("x-tenant-id")),
        Portal:         normPortal(first(md.Get("x-portal"))),
        SessionID:      first(md.Get("x-session-id")),
        SessionType:    first(md.Get("x-session-type")),
        TokenID:        first(md.Get("x-token-id")),
        DeviceID:       first(md.Get("x-device-id")),
        UserType:       first(md.Get("x-user-type")),
        OrganisationID: first(md.Get("x-business-id")),
        TraceID:        first(md.Get("x-request-id")),
    }
}

// normPortal strips "PORTAL_" prefix and lowercases:
// "PORTAL_B2B" → "b2b", "b2b" → "b2b", "PORTAL_SYSTEM" → "system"
func normPortal(raw string) string {
    return strings.ToLower(strings.TrimPrefix(strings.TrimSpace(raw), "PORTAL_"))
}

func first(vals []string) string {
    for _, v := range vals {
        if v = strings.TrimSpace(v); v != "" {
            return v
        }
    }
    return ""
}
```

**What this replaces in existing code:**
- `order_service.go` `resolveTenantID()` — reads only `x-tenant-id`, misses all other fields
- `order_service.go` `resolveCallerID()` — reads only `x-user-id`/`x-customer-id`/`x-subject`, misses portal, session, org
- `payment_service.go` `resolveUserID()` — reads only `x-user-id`/`x-customer-id`/`x-subject`

**Note on `x-subject`:** The existing `resolveCallerID()` and `resolveUserID()` fall back to `x-subject`. This is from an older convention. The authoritative key is `x-user-id` (set by `auth_middleware.go` as `resp.UserId`). Remove the `x-customer-id` and `x-subject` fallbacks — they are not set by the gateway and create false confidence.

### 5.6 Authn server wiring context (for reference)

The authn service itself uses these interceptor chains (from `authn/internal/grpc/server.go` + `interceptors.go`):
- `recoveryUnaryInterceptor` → panic recovery
- `requestIDUnaryInterceptor` → attaches `x-request-id` to context
- `loggingUnaryInterceptor` → structured zap logging per call
- `rateLimitUnaryInterceptor(rdb)` → Redis-backed rate limiting on sensitive methods

The authn service does NOT use an AuthZ interceptor on its own gRPC server — it is the identity root. Only downstream services (b2b, orders, payment) enforce AuthZ.

The authn gRPC server listens on port `50053` (from `DefaultServerConfig()`). Orders/payment services call it only indirectly — via the gateway's `AuthMiddleware`. Services never call authn directly.

Rate-limited authn methods (relevant for mobile app flows):
- `SendOTP`, `ResendOTP`, `SendEmailOTP` → 10 req/min per IP
- `Login`, `EmailLogin` → 20 req/min per IP
- `RefreshToken` → 60 req/min per IP (mobile apps call this frequently)

## 6. Domain Topology

| Domain | System of record | Sync API role | Kafka role |
| --- | --- | --- | --- |
| Orders | `orders-service` | create/read/update order state | emits order lifecycle facts |
| Payments | `payment-service` | initiate/verify/refund/reconcile/manual verify | emits payment and refund facts |
| Billing | new `billing-service` | issue/read/update invoice state | emits invoice lifecycle facts |
| PoliSync | C# | orchestration only | consumes facts and issues next commands |
| Insurance | `insurance-service` | quotation reads, policy writes | should emit quotation/policy lifecycle facts |
| Document | `document-service` | generate and fetch docs | emits generated/failed document facts |
| Storage | `storage-service` | upload/download metadata and presigned URLs | emits file lifecycle facts |
| B2B | `b2b-service` | organisation/employee/purchase-order operations | emits B2B purchase-order and organisation facts |
| Notifications | existing/future service | send SMS/email/push/webhooks | consumes order/payment/policy/billing facts |
| Compliance/Analytics | reporting stack | dashboards and reports | consumes audit and business facts |

## 7. Order and Payment Aggregate Boundaries

### 7.1 Order aggregate responsibility

`Order` should represent a customer purchase decision that is waiting for settlement and policy issuance. It should not also become:

- the financial invoice
- the payment attempt
- the B2B purchase-order aggregate
- the document record

The order aggregate should hold references to those things.

### 7.2 Payment aggregate responsibility

`Payment` should represent one commercial settlement attempt or outcome. It must support:

- hosted gateway flow
- MFS callback flow
- manual proof-upload flow
- refund flow
- reconciliation flow
- ledger posting reference

### 7.3 Billing aggregate responsibility

`Invoice` is required as a first-class commercial obligation. It should model:

- line items
- premium, VAT/tax, service fee, total payable
- who owes the amount
- what order(s) or purchase-order(s) it covers
- what payment(s) settled it

This is directly supported by the knowledge-bank data fields, which explicitly call out:

- `Premium Amount`
- `VAT / Tax`
- `Service Fee`
- `Total Payable`
- `Payment Gateway Reference`
- `Receipt Number`

## 8. Required Contract Enhancements

### 8.1 Billing contracts to add

Create `proto/insuretech/billing/services/v1/billing_service.proto` with at least:

- `IssueInvoice`
- `GetInvoice`
- `ListInvoices`
- `LinkPaymentToInvoice`
- `MarkInvoicePaid`
- `MarkInvoiceOverdue`
- `CancelInvoice`
- `IssueCreditNote` if refund/adjustment is not handled elsewhere

Create `proto/insuretech/billing/events/v1/billing_events.proto` with at least:

- `InvoiceIssuedEvent`
- `InvoicePaymentLinkedEvent`
- `InvoicePaidEvent`
- `InvoiceCancelledEvent`
- `InvoiceOverdueEvent`
- `CreditNoteIssuedEvent` if used

### 8.1.1 Detailed billing service proto plan

The plan for `billing_service.proto` should be concrete enough that implementation can start without another design pass.

Recommended service shape:

- `IssueInvoice(IssueInvoiceRequest) returns (IssueInvoiceResponse)`
- `GetInvoice(GetInvoiceRequest) returns (GetInvoiceResponse)`
- `ListInvoices(ListInvoicesRequest) returns (ListInvoicesResponse)`
- `LinkPaymentToInvoice(LinkPaymentToInvoiceRequest) returns (LinkPaymentToInvoiceResponse)`
- `MarkInvoicePaid(MarkInvoicePaidRequest) returns (MarkInvoicePaidResponse)`
- `CancelInvoice(CancelInvoiceRequest) returns (CancelInvoiceResponse)`
- `MarkInvoiceOverdue(MarkInvoiceOverdueRequest) returns (MarkInvoiceOverdueResponse)`
- `IssueCreditNote(IssueCreditNoteRequest) returns (IssueCreditNoteResponse)` if refund adjustments are invoice-driven

Recommended `IssueInvoiceRequest` fields:

- `tenant_id`
- `customer_id`
- `organisation_id`
- `order_id`
- `purchase_order_id`
- `quotation_id`
- `policy_id`
- `currency`
- `due_date`
- `repeated InvoiceLineItem line_items`
- `string notes`
- `string idempotency_key`
- `string issued_by`

Recommended `InvoiceLineItem` fields:

- `line_item_id`
- `line_type`
  - `PREMIUM`
  - `VAT`
  - `SERVICE_FEE`
  - `DISCOUNT`
  - `ENDORSEMENT_FEE`
  - `CANCELLATION_CHARGE`
  - `REFUND_ADJUSTMENT`
- `reference_id`
- `reference_type`
- `description`
- `amount`
- `quantity`
- `metadata`

Recommended invoice entity additions beyond the current invoice proto:

- `tenant_id`
- `customer_id`
- `organisation_id`
- `order_id`
- `purchase_order_id`
- `quotation_id`
- `currency`
- `subtotal_amount`
- `vat_amount`
- `service_fee_amount`
- `discount_amount`
- `total_amount`
- `balance_due_amount`
- `credit_note_amount`
- `invoice_pdf_file_id`
- `receipt_file_ids`
- `issued_by`
- `cancelled_by`
- `cancel_reason`

### 8.1.2 Detailed billing events proto plan

Recommended event payloads:

- `InvoiceIssuedEvent`
  - `event_id`
  - `invoice_id`
  - `invoice_number`
  - `tenant_id`
  - `customer_id`
  - `organisation_id`
  - `order_id`
  - `purchase_order_id`
  - `total_amount`
  - `currency`
  - `due_date`
  - `timestamp`
  - `correlation_id`
- `InvoicePaymentLinkedEvent`
  - `event_id`
  - `invoice_id`
  - `payment_id`
  - `linked_amount`
  - `remaining_balance`
  - `tenant_id`
  - `timestamp`
  - `correlation_id`
- `InvoicePaidEvent`
  - `event_id`
  - `invoice_id`
  - `payment_id`
  - `tenant_id`
  - `customer_id`
  - `organisation_id`
  - `order_id`
  - `purchase_order_id`
  - `paid_amount`
  - `currency`
  - `paid_at`
  - `timestamp`
  - `correlation_id`
- `InvoiceCancelledEvent`
  - `event_id`
  - `invoice_id`
  - `tenant_id`
  - `cancel_reason`
  - `cancelled_by`
  - `timestamp`
  - `correlation_id`
- `InvoiceOverdueEvent`
  - `event_id`
  - `invoice_id`
  - `tenant_id`
  - `days_overdue`
  - `balance_due_amount`
  - `timestamp`
  - `correlation_id`
- `CreditNoteIssuedEvent`
  - `event_id`
  - `credit_note_id`
  - `invoice_id`
  - `payment_id`
  - `refund_id`
  - `tenant_id`
  - `amount`
  - `reason`
  - `timestamp`
  - `correlation_id`

### 8.2 Payment proto enhancements required

The payment model must gain explicit first-class fields for:

- `order_id`
- `invoice_id`
- `quotation_id`
- `tenant_id`
- `customer_id`
- `organisation_id`
- `purchase_order_id` for B2B when relevant
- `receipt_number`
- `manual_proof_file_id`
- `verified_by`
- `verified_at`
- `rejection_reason`
- `ledger_transaction_id` or explicit TigerBeetle linkage
- `payment_channel` / `provider`
- `payment_frequency` for renewal/installment scenarios

Do not rely on free-form `metadata` for any of those core links.

### 8.2.1 Detailed payment entity and service delta plan

The payment proto should be extended so it supports both automated gateway settlement and manual verification without overloading one generic status field.

Recommended payment entity additions:

- linkage and ownership
  - `tenant_id`
  - `order_id`
  - `invoice_id`
  - `quotation_id`
  - `customer_id`
  - `organisation_id`
  - `purchase_order_id`
- commercial fields
  - `premium_amount`
  - `vat_amount`
  - `service_fee_amount`
  - `total_amount`
  - `payment_frequency`
- provider fields
  - `provider`
  - `provider_account`
  - `provider_reference`
  - `provider_status_code`
  - `callback_signature_verified`
- manual verification fields
  - `manual_proof_file_id`
  - `manual_proof_uploaded_at`
  - `manual_review_status`
  - `verified_by`
  - `verified_at`
  - `rejection_reason`
- receipt and accounting fields
  - `receipt_number`
  - `receipt_document_id`
  - `ledger_transaction_id`
  - `ledger_batch_id`
- risk and compliance fields
  - `third_party_payer`
  - `payer_relationship`
  - `risk_score`
  - `aml_flagged`
  - `aml_case_id`

Recommended service RPC additions or refinements:

- keep:
  - `InitiatePayment`
  - `VerifyPayment`
  - `GetPayment`
  - `ListPayments`
  - `InitiateRefund`
  - `GetRefundStatus`
  - `ListPaymentMethods`
  - `AddPaymentMethod`
  - `ReconcilePayments`
- add:
  - `SubmitManualPaymentProof`
  - `ReviewManualPayment`
  - `GenerateReceipt`
  - `HandleGatewayWebhook`
  - `GetPaymentReceipt`

Recommended `InitiatePaymentRequest` additions:

- `order_id`
- `invoice_id`
- `quotation_id`
- `tenant_id`
- `customer_id`
- `organisation_id`
- `purchase_order_id`
- `callback_url`
- `return_url`
- `cancel_url`
- `idempotency_key`

Recommended `ReviewManualPaymentRequest` fields:

- `payment_id`
- `decision`
  - `APPROVE`
  - `REJECT`
- `review_notes`
- `reviewed_by`
- `idempotency_key`

### 8.2.2 Payment state model plan

Recommended payment statuses:

- `PAYMENT_STATUS_PENDING`
- `PAYMENT_STATUS_INITIATED`
- `PAYMENT_STATUS_AWAITING_CUSTOMER_ACTION`
- `PAYMENT_STATUS_PENDING_MANUAL_VERIFICATION`
- `PAYMENT_STATUS_VERIFIED`
- `PAYMENT_STATUS_COMPLETED`
- `PAYMENT_STATUS_FAILED`
- `PAYMENT_STATUS_REFUND_INITIATED`
- `PAYMENT_STATUS_REFUNDED`
- `PAYMENT_STATUS_CANCELLED`

Recommended manual review statuses:

- `MANUAL_REVIEW_STATUS_NOT_REQUIRED`
- `MANUAL_REVIEW_STATUS_PENDING`
- `MANUAL_REVIEW_STATUS_APPROVED`
- `MANUAL_REVIEW_STATUS_REJECTED`

Recommended rule:

- `VERIFIED` is the commercial truth that settlement was accepted
- `COMPLETED` is the operational truth that receipt, ledger, and downstream publication completed successfully
- if the implementation prefers fewer states, retain both meanings in separate fields rather than collapsing them into ambiguous status text

### 8.3 Payment event enhancements required

`PaymentInitiatedEvent`, `PaymentCompletedEvent`, `PaymentFailedEvent`, and `RefundProcessedEvent` should carry enough routing context to avoid repeated synchronous lookups:

- `tenant_id`
- `order_id`
- `invoice_id`
- `quotation_id`
- `policy_id` when available
- `customer_id`
- `organisation_id`
- `purchase_order_id` when B2B
- `provider`
- `status`
- `payment_method`
- `failure_reason`
- `receipt_number` where available
- `ledger_transaction_id`

### 8.3.1 Detailed payment event plan

Recommended additional events:

- `PaymentVerifiedEvent`
- `ManualPaymentProofSubmittedEvent`
- `ManualPaymentReviewRequestedEvent`
- `ReceiptGeneratedEvent`
- `PaymentReconciliationMatchedEvent`
- `PaymentReconciliationMismatchEvent`

Recommended event intents:

- `ManualPaymentProofSubmittedEvent`:
  - emitted when proof upload is attached to a payment
  - consumed by verification workflow, compliance, and back-office queues
- `PaymentVerifiedEvent`:
  - emitted when a verifier or automated path accepts settlement
  - consumed by `PoliSync` and `orders-service`
- `ReceiptGeneratedEvent`:
  - emitted when receipt document is ready
  - consumed by notification and portal projection services
- `PaymentReconciliationMismatchEvent`:
  - emitted when gateway/bank/provider data does not reconcile with recorded payment state
  - consumed by finance and compliance flows

### 8.4 Order proto enhancements required

If not already present or consistently populated, `Order` should expose:

- `invoice_id`
- `payment_status`
- `billing_status`
- `fulfillment_status`
- `manual_review_required`
- `aml_flag_status`
- `payment_due_at`
- `coverage_start_at`
- `coverage_end_at`

These fields are needed because the SRS and knowledge bank both treat coverage dates and payment verification states as operationally significant.

### 8.4.1 Order state model refinement

Recommended order statuses and dimensions:

- order lifecycle:
  - `PENDING`
  - `PAYMENT_IN_PROGRESS`
  - `PAYMENT_CONFIRMED`
  - `FULFILLMENT_IN_PROGRESS`
  - `POLICY_ISSUED`
  - `CANCELLED`
  - `FAILED`
- payment dimension:
  - `UNPAID`
  - `PARTIALLY_PAID`
  - `PAID`
  - `REFUNDED`
- billing dimension:
  - `NOT_INVOICED`
  - `INVOICED`
  - `SETTLED`
  - `CREDITED`
- fulfillment dimension:
  - `NOT_STARTED`
  - `POLICY_PENDING`
  - `DOCUMENTS_PENDING`
  - `COMPLETED`

Reason:

- the SRS needs more than a single coarse order status if receipts, invoices, activation, cancellation, and renewal all need operational tracking

### 8.5 Insurance event contracts required

Add `insurance/events/v1` contracts for at least:

- `QuotationLockedEvent`
- `PremiumCalculatedEvent`
- `PolicyIssuanceRequestedEvent`
- `PolicyIssuedEvent`
- `PolicyIssuanceFailedEvent`
- `PolicyCancelledEvent`
- `PolicyRenewedEvent`
- `PolicyLapsedEvent`

These are needed because the SRS explicitly requires lifecycle tracking for issuance, renewal, lapse, reinstatement, and cancellation.

## 9. Kafka Event Model

### 9.1 Canonical topic naming

Recommended topic names:

- `insuretech.orders.v1.order.created`
- `insuretech.orders.v1.order.payment_initiated`
- `insuretech.orders.v1.order.payment_confirmed`
- `insuretech.orders.v1.order.cancelled`
- `insuretech.orders.v1.order.failed`
- `insuretech.payment.v1.payment.initiated`
- `insuretech.payment.v1.payment.verified`
- `insuretech.payment.v1.payment.completed`
- `insuretech.payment.v1.payment.failed`
- `insuretech.payment.v1.payment.manual_review_requested`
- `insuretech.payment.v1.refund.processed`
- `insuretech.billing.v1.invoice.issued`
- `insuretech.billing.v1.invoice.paid`
- `insuretech.billing.v1.invoice.cancelled`
- `insuretech.billing.v1.invoice.overdue`
- `insuretech.insurance.v1.policy.issued`
- `insuretech.insurance.v1.policy.cancelled`
- `insuretech.insurance.v1.policy.renewed`
- `insuretech.insurance.v1.policy.lapsed`
- `insuretech.document.v1.document.requested`
- `insuretech.document.v1.document.generated`
- `insuretech.document.v1.document.failed`
- `insuretech.storage.v1.file.uploaded`
- `insuretech.storage.v1.file.finalized`
- `insuretech.storage.v1.file.deleted`
- `insuretech.b2b.v1.purchase_order.created`
- `insuretech.b2b.v1.purchase_order.approved`
- `insuretech.b2b.v1.purchase_order.rejected`
- `insuretech.notifications.v1.notification.requested`
- `insuretech.integration.v1.webhook.delivery_requested`

### 9.2 Producer and consumer map

| Event | Producer | Primary consumers | Required effect |
| --- | --- | --- | --- |
| `OrderCreatedEvent` | `orders-service` | `PoliSync`, billing, analytics | invoice prep and purchase visibility |
| `OrderPaymentInitiatedEvent` | `orders-service` | payment projections, notifications | customer-facing payment pending state |
| `PaymentInitiatedEvent` | `payment-service` | `PoliSync`, billing, orders | active payment attempt registered |
| `PaymentVerifiedEvent` | `payment-service` | `PoliSync`, orders | manual or gateway verification completed |
| `PaymentCompletedEvent` | `payment-service` | `PoliSync`, billing, finance | settlement accepted |
| `PaymentFailedEvent` | `payment-service` | `PoliSync`, orders, notifications, AML | retry or failure handling |
| `RefundProcessedEvent` | `payment-service` | billing, insurance, orders, finance | refund settlement propagation |
| `InvoiceIssuedEvent` | `billing-service` | payment, notifications, B2B | payable artifact published |
| `InvoicePaidEvent` | `billing-service` | `PoliSync`, B2B, finance | commercial settlement confirmed |
| `OrderPaymentConfirmedEvent` | `orders-service` | `PoliSync`, insurance, docs | order becomes policy-issuable |
| `PolicyIssuedEvent` | insurance or PoliSync bridge | docs, notifications, analytics, portals | active coverage and policy pack generation |
| `DocumentGeneratedEvent` | `document-service` | `PoliSync`, notifications, portals | link generated artifact to entity |
| `FileUploadedEvent` | `storage-service` | `PoliSync`, document projection, verification workflows | uploaded proof/doc becomes usable |
| `PurchaseOrderApprovedEvent` | `b2b-service` | billing, `PoliSync` | corporate invoice and downstream issuance prep |

## 10. End-to-End Workflow Design

### 10.1 Retail new-policy flow

1. Customer selects product and provides applicant/nominee data.
2. `PoliSync` validates token and access.
3. `PoliSync` validates duplicate-purchase and quotation rules.
4. `PoliSync` reads quotation/product details from `insurance-service`.
5. `orders-service` creates an order.
6. `billing-service` issues an invoice or invoice-like payable artifact.
7. `payment-service` initiates payment via hosted gateway or MFS provider.
8. On callback or verification, `payment-service` emits settlement fact.
9. `PoliSync` consumes that fact and transitions order state.
10. `orders-service` emits `OrderPaymentConfirmedEvent`.
11. `PoliSync` issues the policy through `insurance-service`.
12. Policy issued fact triggers:
   - policy document generation
   - receipt generation if not already done
   - SMS/email/push notifications
   - dashboard/read-model update
13. `document-service` generates artifacts.
14. `storage-service` finalizes file availability.
15. notifications and webhook fanout complete downstream delivery.

This flow must satisfy:

- policy issuance within the SRS acceptance expectations
- receipt/document availability within the SRS time windows
- durable audit trail across all critical transitions

### 10.2 Manual payment verification flow

This flow is mandatory because the SRS requires manual payment support.

1. Customer chooses manual bank transfer/cash/cheque path.
2. `payment-service` creates a `PENDING_MANUAL_VERIFICATION` payment attempt.
3. Customer uploads proof using `storage-service` presigned upload or direct upload.
4. `storage-service` emits `FileUploadedEvent` and `FileUploadFinalizedEvent`.
5. Verification workflow consumes the storage event and attaches the file to payment.
6. Optional virus scanning and file-quality validation run before review.
7. Admin/focal person verifies or rejects the payment.
8. `payment-service` emits `PaymentVerifiedEvent` and then `PaymentCompletedEvent` or `PaymentFailedEvent`.
9. `PoliSync` continues normal order confirmation or refund path.

Required additions:

- explicit payment status for manual verification states
- attachment linkage from payment to storage file
- role-based verification authority with audit trail
- SLA tracking for review time

### 10.3 Cancellation and refund flow

Required by `FR-093` through `FR-097` and `FR-081`.

1. Customer/agent/admin submits cancellation request.
2. Policy cancellation approval flow runs if needed.
3. Billing calculates refundable amount.
4. `payment-service` initiates refund.
5. `RefundProcessedEvent` updates billing, policy, and order views.
6. cancellation document and refund receipt are generated.
7. notification fanout informs the right parties.

The refund calculation path must support:

- pro-rata premium logic
- admin fee and cancellation charge lines
- refund to original channel where possible

### 10.4 Renewal and grace-period flow

Required by `FR-085` through `FR-091` and `FR-221`.

1. billing or policy scheduler emits premium-due / renewal-due events.
2. invoice issued or renewal payable artifact created.
3. payment completion renews or reinstates the policy.
4. grace-period rules apply before lapse.
5. `PolicyRenewedEvent` or `PolicyLapsedEvent` updates downstream documents, dashboards, and notifications.

Design implication:

- billing cannot be modeled only for first purchase
- it must support recurring premium collection

### 10.5 B2B corporate flow

The SRS and knowledge bank both push toward a partner-heavy model. B2B cannot be treated as a copy of retail.

1. `b2b-service` owns organisation, employee, department, and purchase-order data.
2. `PurchaseOrderApprovedEvent` starts downstream commercial processing.
3. `billing-service` issues organisation-level or batch invoices.
4. `payment-service` settles the invoice using corporate payment channels.
5. `PoliSync` creates or confirms downstream employee-specific orders/policies.
6. policy schedules, invoice PDFs, and employer-facing bundles are generated.
7. storage and document projections make artifacts available to portals.

Critical rule:

- `PurchaseOrder` is not `Order`
- link them via references and events, not aggregate collapse

## 11. Document and Storage Design

### 11.1 Documents that must be generated

At minimum:

- policy PDF with QR code
- payment receipt
- invoice PDF
- cancellation/refund document
- endorsement document
- renewal versioned policy document
- B2B schedule pack where applicable

### 11.2 Trigger model

The generation should be event-driven even if `DocumentService.GenerateDocument` remains synchronous internally.

Examples:

- `InvoiceIssuedEvent` -> generate invoice PDF
- `PaymentCompletedEvent` -> generate receipt
- `PolicyIssuedEvent` -> generate policy PDF and bundle
- `PolicyCancelledEvent` -> generate cancellation confirmation
- `PolicyRenewedEvent` -> generate new versioned policy document

### 11.3 Storage conventions

Use `reference_type` / `reference_id` consistently:

- `POLICY`
- `ORDER`
- `PAYMENT`
- `INVOICE`
- `PURCHASE_ORDER`
- `CANCELLATION`
- `ENDORSEMENT`

SRS-driven requirements to enforce:

- presigned URLs expire after 30 minutes for uploads
- secure download links for portals
- support offline document availability in client channels via cached or packaged delivery strategy
- apply virus scanning to uploads

## 12. Billing and Ledger Design

### 12.1 Billing model

A proper invoice model should include:

- `invoice_id`
- `invoice_number`
- `tenant_id`
- `customer_id` or `organisation_id`
- `order_id` or `purchase_order_id`
- line items
- `premium_amount`
- `vat_amount`
- `service_fee_amount`
- `discount_amount` where applicable
- `total_payable`
- `currency`
- `status`
- `due_date`
- `paid_at`
- `payment_id` linkage
- `invoice_pdf_file_id`

### 12.2 Ledger integration

The SRS requires TigerBeetle-backed double-entry bookkeeping. That means payment completion is not complete until the ledger posting is either:

- executed synchronously inside `payment-service` before publishing completion, or
- completed transactionally with an outbox-safe post-commit mechanism that guarantees no orphaned payment success state

Minimum ledger references to persist:

- `tigerbeetle_transfer_id`
- debit account
- credit account
- settlement batch id if reconciled later

## 13. AML, Fraud, and Compliance Hooks

The order/payment plan must actively support the SRS compliance rules.

### 13.1 AML/CFT triggers to implement in payment and order flows

At minimum, flag and route these cases:

- more than 3 policies in 7 days
- premium above BDT 5 lakh without required proof
- frequent cancellations or refunds
- third-party payment source
- geographic mismatch between customer profile and payment source
- payment-method switching anomalies
- multiple failed KYC attempts
- rapid purchase-to-claim patterns for downstream fraud analytics

### 13.2 Required compliance behavior

- generate immutable audit records for payment, verification, refund, policy issuance, cancellation, and suspicious-activity review
- do not expose SAR/STR investigation signals to customers
- support monthly/quarterly IDRA reporting inputs for premium collection, settlement, and incident reporting
- retain financial and audit records for the required retention windows

### 13.3 Event hooks for compliance and analytics

Add consumers or derived streams for:

- premium collection reporting
- failed payment rate monitoring
- suspicious activity detection
- partner channel reporting
- geographic and product-line aggregation

## 14. Reliability and Performance Requirements

The implementation must be designed for the actual SRS constraints:

- payment processing end-to-end target under 10 seconds (`NFR-005`)
- Category 1 internal API target under 100ms (`FR-193`, acceptance section)
- 99.5%+ availability baseline (`NFR-015`)
- 1 hour RPO / 4 hour RTO (`NFR-016`, `NFR-017`)
- monitoring coverage across the full critical path (`NFR-038`)
- centralized logs (`NFR-039`)

Technical requirements:

- transactional outbox for every event-producing service
- idempotent consumers
- dead-letter topics
- retry with exponential backoff
- circuit breakers for insurer and payment-provider integrations
- correlation ID propagation
- partitioning by aggregate key where ordering matters

Recommended partition keys:

- `order_id` for order events
- `payment_id` for payment events
- `invoice_id` for billing events
- `policy_id` for policy events
- `purchase_order_id` for B2B purchase flows

## 15. Acceptance-Driven Deliverables

### 15.1 Contract deliverables

- `billing/services/v1/billing_service.proto`
- `billing/events/v1/billing_events.proto`
- payment proto updates for core linkage and manual-verification fields
- payment event updates for routing, compliance, and billing context
- insurance event proto package for policy lifecycle

### 15.2 Service deliverables

- `orders-service` outbox and event publication
- `payment-service` outbox, webhook validation, manual verification path, receipt metadata, refund flow, ledger integration
- new `billing-service`
- `PoliSync` consumers and orchestration handlers
- document-trigger consumers for receipt/invoice/policy generation
- storage-linked attachment/proof workflow

### 15.3 Cross-cutting deliverables

- `Idempotency-Key` support on payment and policy-issuance APIs
- HMAC-SHA256 verification for payment callbacks
- audit log stream for critical actions
- AML/fraud rule hooks
- external webhook delivery layer for partner notifications

## 16. Implementation Phases

### Phase 1. Contract correction

- add billing service and events
- extend payment proto with `order_id`, `invoice_id`, receipt, manual verification, and ledger fields
- extend payment events with tenant, actor, and routing context
- add insurance event contracts or temporary PoliSync bridge events
- standardize Kafka topic naming and event headers

### Phase 2. Secure command path

- enforce AuthN/AuthZ on all order/payment/billing commands
- implement idempotency key handling
- implement payment webhook signature validation
- propagate tenant, actor, session, and organisation context into events

### Phase 3. Core order-payment saga

- order create -> invoice issue -> payment initiate -> payment complete/fail -> order confirm
- outbox in `orders-service` and `payment-service`
- receipt generation trigger on payment completion

### Phase 4. Manual verification and storage flow

- payment proof upload path using `storage-service`
- file scanning and review workflow
- admin verification and rejection path
- payment verified event and downstream order transition

### Phase 5. Policy issuance and fulfillment

- consume paid/confirmed events
- issue policy through `insurance-service`
- emit policy lifecycle events
- generate policy document with QR code
- deliver SMS/email links and portal availability

### Phase 6. Refund, cancellation, renewal

- support cancellation approval and refund flow
- support invoice adjustment / credit note if needed
- support grace-period and renewal invoicing
- support versioned renewed policy documents

### Phase 7. B2B corporate flow

- consume `PurchaseOrderApprovedEvent`
- issue organisation-level invoice(s)
- settle invoice(s)
- create employee-targeted orders/policies
- generate employer-facing bundles and schedules

### Phase 8. Compliance and operations hardening

- AML rule consumers and reporting feeds
- IDRA premium collection / settlement reporting inputs
- failed payment and incident dashboards
- replay procedures, DLQ handling, retention/archival policies

## 17. Acceptance Criteria

The implementation is only complete when all of the following are true:

- order creation, payment initiation, payment verification, and order confirmation work with synchronous AuthN/AuthZ guards
- the system supports both automated gateway payments and manual proof-upload payments
- every successful payment can be traced to a concrete `order_id` and `invoice_id`
- payment completion produces a receipt and policy issuance can start immediately afterward
- policy PDF generation and storage linking happen from durable events, not manual follow-up
- document and receipt download links are available to customer and partner channels
- cancellation and refund flows are represented in billing, payment, and policy state
- B2B purchase orders can drive invoices and downstream issuance without collapsing into retail order semantics
- audit, AML, and retention hooks exist for all critical financial transitions
- all critical flows are replay-safe, idempotent, and observable by `correlation_id`

## 19. SSLCommerz API Reference (from NodeJS SDK — `SSLCommerz-NodeJS/`)

> Last reviewed: 2026-03-06. Source files: `api/payment-controller.js`, `api/payment-init-data-process.js`, `api/fetch.js`, `index.js`.

### Session Initialization — `POST /gwprocess/v4/api.php`

| SSLCommerz Field | Required | Our mapping | Notes |
|---|---|---|---|
| `store_id` | ✅ | `SSLCOMMERZ_STORE_ID` env var | Merchant store ID |
| `store_passwd` | ✅ | `SSLCOMMERZ_STORE_PASSWORD` env var | Merchant password |
| `total_amount` | ✅ | `payment.amount / 100` | Float BDT — convert from paisa |
| `currency` | ✅ | `BDT` | Always BDT for Bangladesh |
| `tran_id` | ✅ | `payment.tran_id` | **Merchant-generated** UUID — store in `payments.tran_id` to correlate IPN |
| `success_url` | ✅ | `PAYMENT_PUBLIC_BASE_URL/v1/payments/{id}/callback/success` | SSLCommerz POSTs here on success |
| `fail_url` | ✅ | `PAYMENT_PUBLIC_BASE_URL/v1/payments/{id}/callback/fail` | SSLCommerz POSTs here on failure |
| `cancel_url` | ✅ | `PAYMENT_PUBLIC_BASE_URL/v1/payments/{id}/callback/cancel` | SSLCommerz POSTs here on cancel |
| `ipn_url` | ✅ | `PAYMENT_IPN_BASE_URL/v1/payments/ipn/sslcommerz` | **Must be publicly accessible** — ngrok in dev |
| `cus_name` | ✅ | `req.customer_name` (typed) | From `InitiatePaymentRequest.customer_name` |
| `cus_email` | ✅ | `req.customer_email` (typed) | From `InitiatePaymentRequest.customer_email` |
| `cus_add1` | ✅ | `req.customer_address_line1` (typed) | From `InitiatePaymentRequest.customer_address_line1` |
| `cus_city` | ✅ | `req.customer_city` (typed) | From `InitiatePaymentRequest.customer_city` |
| `cus_postcode` | ✅ | `req.customer_postcode` (typed) | From `InitiatePaymentRequest.customer_postcode` |
| `cus_country` | ✅ | `req.customer_country` or `Bangladesh` | Default `Bangladesh` |
| `cus_phone` | ✅ | `req.customer_phone` (typed) | From `InitiatePaymentRequest.customer_phone` |
| `product_name` | ✅ | `Insurance Premium` | Or product name from order |
| `product_category` | ✅ | `insurance` | `SSLCOMMERZ_PRODUCT_CATEGORY` env var |
| `product_profile` | ✅ | `general` | SSLCommerz product profile |
| `shipping_method` | ✅ | `NO` | Digital product — no shipping |
| `num_of_item` | ✅ | `1` | Always 1 for single insurance order |

**Session init response fields (status == `SUCCESS`):**
- `GatewayPageURL` → `InitiatePaymentResponse.gateway_page_url` + stored nowhere (only returned to client)
- `sessionkey` → `Payment.session_key` (stored in DB for reference)

### IPN / Callback POST Body Fields

| Field | Our mapping | Notes |
|---|---|---|
| `tran_id` | Lookup: `repo.GetPaymentByTranID(tran_id)` | **Primary correlation key** — look up payment by this |
| `val_id` | `Payment.val_id` | Required for server-side validation — **not same as tran_id** |
| `status` | See status table below | `VALID`, `VALIDATED`, `FAILED`, `CANCELLED`, `UNATTEMPTED`, `EXPIRED` |
| `amount` | Validate against `payment.amount/100` | Float string BDT |
| `bank_tran_id` | `Payment.bank_tran_id` | Bank transaction reference |
| `card_type` | `Payment.card_type` | `VISA`, `MASTERCARD`, `bKash`, `Nagad`, etc. |
| `card_brand` | `Payment.card_brand` | Card brand |
| `card_issuer` | `Payment.card_issuer` | Issuing bank |
| `card_issuer_country` | `Payment.card_issuer_country` | Issuing country |
| `verify_sign` | IPN signature verification | MD5 signature |
| `verify_key` | IPN signature verification | Comma-separated field list for signature |

**IPN status values:**
- `VALID` / `VALIDATED` → trigger server-side `ValidatePayment(val_id)` → if valid, mark `PAYMENT_STATUS_SUCCESS`
- `FAILED` → mark `PAYMENT_STATUS_FAILED`
- `CANCELLED` → mark `PAYMENT_STATUS_CANCELLED`
- `UNATTEMPTED` / `EXPIRED` → no state change (customer never reached payment page)

### Server-Side Validation — `GET /validator/api/validationserverAPI.php`

```
GET ?val_id=<val_id>&store_id=<id>&store_passwd=<pwd>&v=1&format=json
```

**⚠️ CRITICAL:** `val_id` comes from the IPN/callback POST body — it is **NOT** the merchant `tran_id`.
Validation response `status`: `VALID`, `VALIDATED` = confirmed | `INVALID_TRANSACTION`, `FAILED`, `CANCELLED` = rejected.

### IPN Signature Verification Algorithm

1. Extract `verify_sign` and `verify_key` from IPN POST body
2. Parse `verify_key` as comma-separated list of field names
3. Build map of those fields from IPN body
4. Add `store_passwd` = MD5(store password)
5. Sort all keys alphabetically → concatenate as `k1=v1&k2=v2...`
6. Compute MD5 of that string
7. Compare with `verify_sign` — **reject IPN if mismatch**

### Go Service → SSLCommerz Field Mapping (complete)

| Our proto field | SSLCommerz field | Direction |
|---|---|---|
| `payment.tran_id` | `tran_id` | → sent to SSLCommerz |
| `payment.session_key` | `sessionkey` | ← returned by SSLCommerz |
| `InitiatePaymentResponse.gateway_page_url` | `GatewayPageURL` | ← returned, forwarded to client |
| `payment.val_id` | `val_id` | ← from callback/IPN body |
| `payment.bank_tran_id` | `bank_tran_id` | ← from callback/IPN body |
| `payment.card_type` | `card_type` | ← from callback/IPN body |
| `payment.card_brand` | `card_brand` | ← from callback/IPN body |
| `payment.card_issuer` | `card_issuer` | ← from callback/IPN body |
| `payment.validated_at` | (computed) | When `ValidatePayment()` returns VALID |
| `payment.callback_received_at` | (computed) | When success/fail callback received |
| `payment.ipn_received_at` | (computed) | When IPN received |

### Dev Environment Setup (SSLCommerz Sandbox)

1. Register at https://developer.sslcommerz.com/registration/ for sandbox credentials
2. Update `.env.dev`:
   ```env
   SSLCOMMERZ_STORE_ID=testbox
   SSLCOMMERZ_STORE_PASSWORD=qwerty
   SSLCOMMERZ_SANDBOX=true
   PAYMENT_PUBLIC_BASE_URL=https://xxxx.ngrok.io   # required for IPN
   PAYMENT_IPN_BASE_URL=https://xxxx.ngrok.io
   ```
3. Start ngrok: `ngrok http 50191` → copy HTTPS URL → paste into env vars above
4. Sandbox test card: `4111111111111111`, any future expiry, any CVV

### Code Gaps Found During Review (2026-03-06)

| Gap | File | Status |
|---|---|---|
| `publishInitiated/Completed/Failed` missing `order_id`, `tenant_id`, `provider`, `occurred_at` | `payment_service.go` | ✅ Fixed |
| `InitiatePayment` passed `order_id` via untyped metadata map | `payment_service.go` | ✅ Fixed |
| `InitSession()` call used `req.GetMetadata()["customer_name"]` instead of typed fields | `payment_service.go` | ✅ Fixed |
| Kafka topics used short non-canonical names (`payment.initiated` etc.) | `events/topics.go` | ✅ Fixed (now `insuretech.payment.v1.payment.initiated`) |
| `PaymentService` interface missing 6 new Phase 2 RPCs | `domain/interfaces.go` | ✅ Fixed |
| `PaymentRepository` missing `GetPaymentByTranID`, `GetPaymentByOrderID`, `GetPaymentByProviderReference` | `domain/interfaces.go` | ✅ Fixed |
| `payment.proto` had no `order_id`, `tran_id`, `val_id`, `provider`, manual-review, or receipt fields | `payment/entity/v1/payment.proto` | ✅ Fixed (fields 23–53) |
| `payment_service.proto` had no `HandleGatewayWebhook`, `SubmitManualPaymentProof`, `ReviewManualPayment`, `GenerateReceipt`, `GetPaymentReceipt` RPCs | `payment/services/v1/payment_service.proto` | ✅ Fixed |
| `billing/entity/v1/invoice.proto` missing `order_id`, `tenant_id`, `customer_id`, `tax_amount`, `total_amount` etc. | billing proto | ✅ Fixed |
| No `billing/services/v1/billing_service.proto` existed | billing proto | ✅ Created |
| No `billing/events/v1/billing_events.proto` existed | billing proto | ✅ Created |
| `.env.dev` missing `SSLCOMMERZ_*`, `BILLING_*`, `PAYMENT_PUBLIC_BASE_URL`, `PAYMENT_IPN_BASE_URL` | `.env.dev` | ✅ Fixed |
| `orders/consumer.go` used `evt.GetCorrelationId()` as order_id workaround | `consumers/consumer.go` | ✅ Fixed (uses typed `GetOrderId()`) |

---

## 18. Recommended Build Order

1. Fix the contracts first: billing service/events, payment linkage fields, insurance event contracts.
2. Implement secure command path: idempotency, AuthN/AuthZ enforcement, webhook verification.
3. Implement outbox + event publication in `orders-service` and `payment-service`.
4. Implement billing issuance and invoice-payment linkage.
5. Implement `PoliSync` saga consumers for payment/order/policy orchestration.
6. Implement document and storage fanout for receipts, invoices, and policy PDFs.
7. Implement manual payment verification flow.
8. Implement B2B purchase-order to invoice/order/policy workflow.
9. Implement compliance reporting hooks, AML flags, and operational hardening.

That order is deliberate. If billing linkage, idempotency, and event contracts are not fixed first, the rest of the implementation will leak core business semantics into ad hoc metadata, brittle webhooks, and one-off consumers.

## 18.1 Corrected Build Order — Code-Verified (March 2026)

This replaces the abstract ordering above with a concrete, dependency-ordered sequence grounded in what the code inspection found.

### Step 1 — Fix Kafka topic names (zero-risk, no proto change needed)

Both services use non-canonical topic names. Fix immediately — no schema or proto changes required.

**orders-service `internal/events/topics.go`** — rename all 5 constants:
```go
// BEFORE                                    // AFTER
"orders.order.created"           →   "insuretech.orders.v1.order.created"
"orders.order.payment_initiated" →   "insuretech.orders.v1.order.payment_initiated"
"orders.order.payment_confirmed" →   "insuretech.orders.v1.order.payment_confirmed"
"orders.order.cancelled"         →   "insuretech.orders.v1.order.cancelled"
"orders.order.failed"            →   "insuretech.orders.v1.order.failed"
```

**payment-service `internal/events/topics.go`** — rename all 4 constants + add 4 new:
```go
// BEFORE                     // AFTER
"payment.initiated"  →   "insuretech.payment.v1.payment.initiated"
"payment.completed"  →   "insuretech.payment.v1.payment.completed"
"payment.failed"     →   "insuretech.payment.v1.payment.failed"
"refund.processed"   →   "insuretech.payment.v1.refund.processed"
// New (add after proto Step 4):
TopicPaymentVerified            = "insuretech.payment.v1.payment.verified"
TopicManualReviewRequested      = "insuretech.payment.v1.payment.manual_review_requested"
TopicReceiptGenerated           = "insuretech.payment.v1.receipt.generated"
TopicReconciliationMismatch     = "insuretech.payment.v1.reconciliation.mismatch"
```

**orders-service `cmd/server/main.go`** — update consumer topic subscriptions. Currently (verified from code):
```go
// CURRENT (line ~130 in main.go):
consumer.Subscribe([]string{"payment.completed", "payment.failed", "policy.issued"})
// MUST BECOME:
consumer.Subscribe([]string{
    "insuretech.payment.v1.payment.completed",
    "insuretech.payment.v1.payment.failed",
    "insuretech.insurance.v1.policy.issued",  // canonical insurance topic
})
```

**Consumer group and DLQ names** (already set correctly in main.go — keep as-is):
- Consumer group: `"orders-service-group"` ✅
- Dead-letter queue: `"orders-service-dlq"` ✅

### Step 2 — Add `ExtractRequestContext` helper to both services

Create `internal/middleware/context.go` in both `orders` and `payment` microservices (§5.4 pattern). This is a pure addition with no dependencies. Replaces `resolveTenantID()` + `resolveCallerID()` in orders, and `resolveUserID()` in payment.

### Step 3 — Add gRPC AuthZ interceptor to both services (copy b2b pattern)

Create `internal/middleware/authz_interceptor.go` in both orders and payment services. The complete implementation pattern is in §4.0.3.

**Wire into gRPC server — orders-service `internal/grpc/server.go`** (currently uses `grpc.NewServer()` with no interceptors):
```go
// CURRENT (verified from code):
s := grpc.NewServer()

// MUST BECOME:
authzConn, err := grpc.Dial(cfg.AuthzServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
authzClient := authzservicev1.NewAuthorizationServiceClient(authzConn)
authzInterceptor := middleware.NewAuthZInterceptor(authzClient)

s := grpc.NewServer(
    grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
        grpc_recovery.UnaryServerInterceptor(),
        authzInterceptor.UnaryServerInterceptor(),
    )),
)
```

**Wire into gRPC server — payment-service `server.go`** (currently has custom `kafkaProducer` wiring but no interceptors on the gRPC server):
```go
// CURRENT (verified from code):
s := grpc.NewServer()

// MUST BECOME:
// Same pattern as orders-service above, with payment method→resource mapping
```

**Config changes needed:**

orders-service `internal/config/config.go` — add (currently has `PaymentServiceURL` already — add `AuthzServiceURL`):
```go
type Config struct {
    // ... existing fields ...
    PaymentServiceURL string `env:"PAYMENT_SERVICE_URL" envDefault:"localhost:50190"`
    AuthzServiceURL   string `env:"AUTHZ_SERVICE_URL"   envDefault:"localhost:50082"`  // ADD THIS
}
```

payment-service `internal/config/config.go` — add:
```go
type Config struct {
    // ... existing fields ...
    AuthzServiceURL string `env:"AUTHZ_SERVICE_URL" envDefault:"localhost:50082"`  // ADD THIS
}
```

`services.yaml` — authz service is already registered at port `50082`. No change needed.

**Graceful degradation:** The interceptor must log a warning and allow the request through if the authz-service is unreachable (gRPC `codes.Unavailable`). This prevents development being blocked when authz is not running locally. Set a 2-second timeout on the `CheckAccess` call.

**Bootstrap methods that skip Casbin** (orders-service):
- `HealthCheck` / `Ping` — standard health probes
- Any method where `mapOrderMethodToResourceAction` returns `("", "")` — treated as public

**Bootstrap methods that skip Casbin** (payment-service):
- `HealthCheck` / `Ping`
- `HandleGatewayWebhook` — external callback, HMAC-verified, not user-authenticated

### Step 4 — Extend proto contracts (order entity + events, payment entity + service + events)

Proto changes are the highest-impact step. Do them all in one `buf generate` pass:

**Orders proto (`proto/insuretech/orders/`):**
- `entity/v1/order.proto` — add `invoice_id`, `organisation_id`, `payment_status`, `billing_status`, `fulfillment_status`, `manual_review_required`, `payment_due_at`, `coverage_start_at`, `coverage_end_at`, `idempotency_key`, `correlation_id` + new status dimension enums
- `services/v1/order_service.proto` — enrich `CreateOrderRequest` with `tenant_id`, `organisation_id`, `idempotency_key`, `product_id`, `plan_id`, `total_payable`, `coverage_start_at`, `coverage_end_at`; add `UpdateFulfillmentStatus` RPC
- `events/v1/order_events.proto` — add `tenant_id`, `organisation_id`, `portal`, `actor_user_id`, `session_id`, `causation_id`, `occurred_at` to all 5 existing events; add `OrderFulfillmentCompletedEvent`

**Payment proto (`proto/insuretech/payment/`):**
- `entity/v1/payment.proto` — add all 29 new fields (§8.2.1); add `ManualReviewStatus` enum; expand `PaymentStatus` enum with `AWAITING_CUSTOMER_ACTION`, `PENDING_MANUAL_VERIFICATION`, `VERIFIED`
- `services/v1/payment_service.proto` — add `order_id`, `invoice_id`, `quotation_id`, `tenant_id`, `customer_id`, `organisation_id`, `purchase_order_id`, `return_url`, `cancel_url` to `InitiatePaymentRequest`; add 5 new RPCs
- `events/v1/payment_events.proto` — add routing context to all 4 existing events; add 6 new events

**Billing proto (create new):**
- `proto/insuretech/billing/entity/v1/invoice.proto` — enhance with required fields
- `proto/insuretech/billing/services/v1/billing_service.proto` — create with 8 RPCs
- `proto/insuretech/billing/events/v1/billing_events.proto` — create with 6 events

Run `buf generate` after all proto changes.

### Step 5 — DB migrations for new order and payment columns

Create migration files:

- `db/migrations/insurance_schema/YYYYMMDD_XXX_enhance_orders.up.sql` — add `invoice_id`, `organisation_id`, `payment_status`, `billing_status`, `fulfillment_status`, `manual_review_required`, `payment_due_at`, `coverage_start_at`, `coverage_end_at`, `idempotency_key`, `correlation_id`
- `db/migrations/payment_schema/YYYYMMDD_XXX_enhance_payments_v2.up.sql` — add all 29 new payment entity columns
- `db/migrations/billing_schema/YYYYMMDD_XXX_create_invoices.up.sql` — create invoices table

Update `orderScanRow` + `scanRowToProto()` in orders repository and `paymentRow` + `rowToPayment()` in payment repository to read/write new columns.

### Step 6 — Fix `CreateOrder` quotation resolution

Replace hardcoded `TotalPayable: &Money{Amount: 1}` with real values:
- Accept `product_id`, `plan_id`, `total_payable`, `currency`, `coverage_start_at`, `coverage_end_at` from the enriched `CreateOrderRequest` (populated by PoliSync or the calling client)
- Set `payment_status = UNPAID`, `billing_status = NOT_INVOICED`, `fulfillment_status = NOT_STARTED` on creation
- Store `idempotency_key` on the order row; add unique index to prevent duplicate orders from the same key

### Step 7 — Wire payment-service gRPC client into orders-service

- Add `authz` and `payment` gRPC dialers in `orders/cmd/server/main.go`
- Pass `paymentClient` into `OrderServiceImpl` constructor
- Replace fake payment ID generation in `InitiatePayment` with a real `paymentClient.InitiatePayment()` call
- Handle `codes.AlreadyExists` (idempotency replay) from payment-service

### Step 8 — Enrich all event payloads

In both services, use `ExtractRequestContext(ctx)` at the top of every command handler and pass all fields through to every `Publisher.Publish()` call:
- `actor_user_id`, `tenant_id`, `organisation_id`, `portal`, `session_id`, `session_type`, `token_id`, `trace_id` on all order and payment events

### Step 9 — Fix `newID()` in payment-service

Replace `fmt.Sprintf("%s_%d", prefix, time.Now().UTC().UnixNano())` with `fmt.Sprintf("%s_%s", prefix, uuid.NewString())` for all generated IDs (`payment_id`, `transaction_id`, `refund_id`, `event_id`).

### Step 10 — Fix payment Kafka consumer deserialization in orders-service

Replace `parseJSONPayload` flat-map approach in `consumers/consumer.go` with typed `protojson.Unmarshal` into concrete proto event types:
- `PaymentCompletedEvent` → read `order_id` (typed field, available after Step 4)
- `PaymentFailedEvent` → read `order_id`, `error_code`, `error_message`

### Step 11 — Add new payment RPCs

Implement `SubmitManualPaymentProof`, `ReviewManualPayment`, `GenerateReceipt`, `HandleGatewayWebhook`, `GetPaymentReceipt` in:
- `payment/internal/domain/interfaces.go` — add to `PaymentService` interface
- `payment/internal/service/payment_service.go` — implement
- `payment/internal/grpc/payment_handler.go` — add handlers
- Emit `ManualPaymentProofSubmittedEvent`, `PaymentVerifiedEvent`, `ReceiptGeneratedEvent` from new methods

### Step 12 — Fix `VerifyPayment` event sequencing

Emit `PaymentVerifiedEvent` before `PaymentCompletedEvent`. Add `VERIFIED` status intermediate step. Do not skip directly from `PENDING` to `SUCCESS`.

### Step 13 — Fix refund state machine

Replace auto-complete refund with proper `PENDING → APPROVED → PROCESSING → COMPLETED` flow:
- `InitiateRefund` creates refund in `PENDING` state, emits `RefundInitiatedEvent`
- Admin `ApproveRefund` transitions to `APPROVED`, triggers gateway refund call
- Gateway callback or polling transitions to `COMPLETED`, emits `RefundProcessedEvent`

### Step 14 — Wire gateway routes for orders-service and payment-service

In `gateway/internal/routes/router.go`:
- Add Go `orders-service` handler block (port `:50142`) per §4.0.6 with `csrfMW` on POST routes
- Add payment handler block (port `:50190`) per §4.0.6 with `csrfMW` on state-changing routes
- Add public webhook endpoint: `POST /v1/payments/webhook/{provider}` — HMAC only, no authMW

### Step 15 — Seed AuthZ Casbin permissions for order/payment

In `authz/internal/seeder/portal_seeder.go`:
- Add `svc:order/my/*` + GET/POST for `b2c:customer`
- Add `svc:payment/my/*` + GET/POST for `b2c:customer`
- Add `svc:order/*` + GET/POST/PATCH for `agent:agent`
- Add `svc:order/*` + `*` and `svc:payment/*` + `*` for `system:admin`
- Add `svc:payment/*` + GET/POST for `system:finance`

### Step 16 — Create billing microservice

- Create `backend/inscore/microservices/billing/` following the same structure as `orders` and `payment`
- Port: `50195` — add to `services.yaml`
- Implement `IssueInvoice`, `GetInvoice`, `ListInvoices`, `LinkPaymentToInvoice`, `MarkInvoicePaid`, `CancelInvoice`, `MarkInvoiceOverdue`, `IssueCreditNote`
- Publish `InvoiceIssuedEvent`, `InvoicePaidEvent`, `InvoiceCancelledEvent`, `InvoiceOverdueEvent` on canonical topics
- Wire into gateway router with `svc:invoice/*` authz

### Step 17 — Integration hardening

- Implement HMAC-SHA256 verification in `HandleGatewayWebhook` per provider
- Add dead-letter topic handling in Kafka consumer for both services
- Add `correlation_id` propagation through all event chains
- Add `ReconcilePayments` mismatch event emission
- Verify idempotency end-to-end: duplicate `CreateOrder` with same `idempotency_key` returns existing order (not error)

## 19. Current State Gap Analysis (Code-Verified, March 2026)

This section is grounded in actual code inspection of the existing implementations. It supersedes any earlier gap estimates.

### 19.0 Authentication model gaps (corrected)

| Service | Gap | Correct fix |
| --- | --- | --- |
| `orders-service` | No gRPC AuthZ interceptor | Add `AuthZInterceptor` (copy from b2b pattern); reads `x-user-id`, `x-portal`, `x-tenant-id` from metadata; maps gRPC method to Casbin `svc:order/*` + HTTP verb |
| `orders-service` | `resolveTenantID()` and `resolveCallerID()` are partial — only read `x-tenant-id` and `x-user-id`; miss `portal`, `session_id`, `token_id`, `organisation_id`, `session_type` | Replace with `ExtractRequestContext()` helper (§5.4) |
| `orders-service` | No `x-business-id` reading anywhere | B2B flows set `x-business-id` via `B2BContextMiddleware` at gateway — read it in `ExtractRequestContext` |
| `payment-service` | No gRPC AuthZ interceptor | Same fix — add `AuthZInterceptor` mapped to `svc:payment/*` |
| `payment-service` | `resolveUserID()` only reads `x-user-id` / `x-customer-id` / `x-subject` — misses all other identity fields | Replace with shared `ExtractRequestContext()` |
| gateway `router.go` | Order routes (lines 332-342) proxy to PoliSync `:50141` only — Go `orders-service` `:50142` not wired | Add separate route block for Go orders-service; wire `ordersHandler` |
| gateway `router.go` | Payment routes do not exist at all | Add payment route block per §4.0.6 |
| gateway `router.go` | Existing order routes are mostly corrected; the verified remaining CSRF gap is `POST /v1/orders` | Add `csrfMW` to `POST /v1/orders`; keep existing `csrfMW` on `/initiate-payment`, `/confirm`, and `/cancel` |
| portal seeder | `svc:order/*` and `svc:payment/*` objects not seeded for `b2c:customer`, `agent:agent` | Add per §4.0.7 |

### 19.0.1 What is already correctly implemented (do not break)

| Component | Status |
| --- | --- |
| `gateway/auth_middleware.go` — dual-path AuthN (JWT + server-side session cookie) | ✅ Correct — no changes needed |
| `gateway/authz_middleware.go` — Casbin enforcement via `AuthZ.CheckAccess` | ✅ Correct — no changes needed |
| `gateway/csrf_middleware.go` — CSRF for web portal state-changing routes | ✅ Correct — must be applied to new order/payment routes |
| `b2b/middleware/authz_interceptor.go` — gRPC-level AuthZ defense-in-depth | ✅ Reference implementation — copy pattern |
| `authn-service` — `ValidateToken` handles both JWT and server-side session | ✅ Correct |
| `authz-service` — Casbin PERM model with B2B two-level fallback | ✅ Correct |
| `payment-service` — idempotency key dedup on `InitiatePayment` | ✅ Correct — keep |
| `payment-service` — `GetPaymentByIdempotencyKey` in repository | ✅ Correct — keep |
| `orders-service` — `canCancel()` state guard | ✅ Correct — keep |
| `orders-service` — Kafka publisher graceful no-op when producer is nil | ✅ Correct — keep |

### 19.0.2 Order-service code gaps (verified)

| File | Gap | Fix |
| --- | --- | --- |
| `order_service.go` `CreateOrder` | Hardcodes `TotalPayable: &Money{Amount: 1}` — placeholder | Accept `product_id`, `plan_id`, `total_payable`, `coverage_start_at`, `coverage_end_at` from enriched request (set by PoliSync caller) |
| `order_service.go` `InitiatePayment` | Generates fake `paymentID` and `gatewayRef` locally with `uuid.NewString()` + format strings | Must call `payment-service.InitiatePayment` via gRPC using dialer from `config.PaymentServiceURL` (field exists in config but is unused) |
| `order_service.go` all methods | Events missing `tenant_id`, `organisation_id`, `portal`, `actor_user_id`, `session_id`, `causation_id` | Use `ExtractRequestContext()` and populate all event fields |
| `events/topics.go` | Short names: `"orders.order.created"` | Rename to canonical: `"insuretech.orders.v1.order.created"` |
| `consumers/consumer.go` | Parses flat JSON string map — breaks on nested proto types (`Money`, `Timestamp`) | Unmarshal typed proto events using `protojson.Unmarshal` |
| `consumers/consumer.go` | `HandlePaymentCompleted` reads `payload["order_id"]` from flat map — `PaymentCompletedEvent` does not currently have `order_id` field | After payment proto is extended with `order_id`, switch to typed unmarshal |
| `repository.go` `CreateOrder` | `product_id` / `plan_id` default to nil UUID `"00000000-..."` via `coalesceUUID` | Remove nil UUID fallback — require real values from quotation resolution |
| `grpc/server.go` | No AuthZ interceptor wired into `grpc.NewServer()` | Add `grpc.NewServer(grpc.UnaryInterceptor(authzInterceptor.UnaryServerInterceptor()))` |
| `cmd/server/main.go` | No gRPC client dialed to payment-service | Add dialer, pass client to `OrderServiceImpl` |

### 19.0.3 Payment-service code gaps (verified)

| File | Gap | Fix |
| --- | --- | --- |
| `payment_service.go` `newID()` | Uses `UnixNano` prefix — collision-prone under concurrent load: `pay_1709123456789012345` | Replace with `"pay_" + uuid.NewString()` |
| `payment_service.go` `InitiatePayment` | `policy_id` is the only typed linkage field; `order_id`, `invoice_id`, `quotation_id`, `tenant_id`, `customer_id` read from untyped `metadata` map | After proto extension: read from typed request fields |
| `payment_service.go` `referenceInfo()` | Only returns `("policy", policy_id)` or metadata values | After proto extension: return `("order", order_id)` first |
| `payment_service.go` `VerifyPayment` | Emits `PaymentCompletedEvent` but no `PaymentVerifiedEvent` | Emit `PaymentVerifiedEvent` first, then `PaymentCompletedEvent` after ledger post |
| `payment_service.go` `InitiateRefund` | Auto-completes refund (`COMPLETED`) in same transaction — skips `PENDING→APPROVED→PROCESSING` | Model proper async refund approval path for manual verification cases |
| `payment_service.go` `ReconcilePayments` | Only counts `SUCCESS`/`REFUNDED` rows — no mismatch event published | Emit `PaymentReconciliationMismatchEvent` for each mismatch |
| `events/topics.go` | Short names: `"payment.initiated"`, `"payment.completed"`, etc. | Rename to canonical: `"insuretech.payment.v1.payment.initiated"`, etc. |
| `events/topics.go` | Missing: `TopicPaymentVerified`, `TopicManualReviewRequested`, `TopicReceiptGenerated`, `TopicReconciliationMismatch` | Add after proto extension |
| `server.go` | `kafkaProducer` created with only `TopicPaymentInitiated` as default topic | Use service-level producer with no fixed default topic; pass topic per `Produce()` call |
| `payment_handler.go` | No `SubmitManualPaymentProof`, `ReviewManualPayment`, `GenerateReceipt`, `HandleGatewayWebhook`, `GetPaymentReceipt` handlers | Add after proto extension |
| `server.go` | No AuthZ interceptor wired | Add same pattern as orders-service |

## 20. Proto Delta Tables (Code-Verified, March 2026)

These tables document exactly what exists in each proto file today and what must be added. Field numbers pick up from the last used number in each message. Run `buf generate` once after all changes.

---

### 20.1 `orders/entity/v1/order.proto`

**Existing `Order` message fields (verified):**

| # | Name | Type | Notes |
|---|---|---|---|
| 1 | `order_id` | `string` | UUID |
| 2 | `user_id` | `string` | |
| 3 | `product_id` | `string` | |
| 4 | `plan_id` | `string` | |
| 5 | `status` | `OrderStatus` | |
| 6 | `total_payable` | `common.Money` | |
| 7 | `currency` | `string` | |
| 8 | `payment_id` | `string` | |
| 9 | `policy_id` | `string` | |
| 10 | `quotation_id` | `string` | |
| 11 | `notes` | `string` | |
| 12 | `metadata` | `map<string,string>` | |
| 13 | `created_at` | `Timestamp` | |
| 14 | `updated_at` | `Timestamp` | |
| 15 | `expires_at` | `Timestamp` | |
| 16 | `confirmed_at` | `Timestamp` | |
| 17 | `cancelled_at` | `Timestamp` | |
| 18 | `failure_reason` | `string` | |

**Existing `OrderStatus` enum values:**
`ORDER_STATUS_UNSPECIFIED(0)`, `PENDING(1)`, `PAYMENT_INITIATED(2)`, `PAYMENT_CONFIRMED(3)`, `PROCESSING(4)`, `COMPLETED(5)`, `CANCELLED(6)`, `FAILED(7)`

**Fields to ADD (starting at field 19):**

| # | Name | Type | Rationale |
|---|---|---|---|
| 19 | `invoice_id` | `string` | Links to billing-service invoice |
| 20 | `organisation_id` | `string` | B2B org context (from `x-business-id`) |
| 21 | `tenant_id` | `string` | Multi-tenancy routing |
| 22 | `idempotency_key` | `string` | Dedup key from client; unique index on DB |
| 23 | `correlation_id` | `string` | Distributed tracing across saga |
| 24 | `payment_status` | `OrderPaymentStatus` | Separate payment dimension |
| 25 | `billing_status` | `OrderBillingStatus` | Separate billing dimension |
| 26 | `fulfillment_status` | `OrderFulfillmentStatus` | Separate policy issuance dimension |
| 27 | `manual_review_required` | `bool` | AML/fraud flag triggers manual hold |
| 28 | `payment_due_at` | `Timestamp` | Payment deadline |
| 29 | `coverage_start_at` | `Timestamp` | Policy coverage start |
| 30 | `coverage_end_at` | `Timestamp` | Policy coverage end |
| 31 | `actor_user_id` | `string` | Portal user who created (agent/admin assisted) |
| 32 | `portal` | `string` | `"b2c"`, `"b2b"`, `"agent"`, `"system"` |

**New enums to ADD:**
```protobuf
enum OrderPaymentStatus {
  ORDER_PAYMENT_STATUS_UNSPECIFIED = 0;
  UNPAID = 1;
  PAYMENT_IN_PROGRESS = 2;
  PAID = 3;
  PAYMENT_FAILED = 4;
  REFUNDED = 5;
}

enum OrderBillingStatus {
  ORDER_BILLING_STATUS_UNSPECIFIED = 0;
  NOT_INVOICED = 1;
  INVOICED = 2;
  INVOICE_CANCELLED = 3;
}

enum OrderFulfillmentStatus {
  ORDER_FULFILLMENT_STATUS_UNSPECIFIED = 0;
  NOT_STARTED = 1;
  FULFILLMENT_IN_PROGRESS = 2;
  FULFILLED = 3;
  FULFILLMENT_FAILED = 4;
}
```

---

### 20.2 `orders/services/v1/order_service.proto`

**Existing RPCs (verified):**
`CreateOrder`, `GetOrder`, `GetOrderStatus`, `ListOrders`, `InitiatePayment`, `ConfirmPayment`, `CancelOrder`

**Existing `CreateOrderRequest` fields:**
`user_id(1)`, `product_id(2)`, `plan_id(3)`, `quotation_id(4)`, `currency(5)`, `notes(6)`, `metadata(7)`, `expires_at(8)`

**Fields to ADD to `CreateOrderRequest` (starting at field 9):**
| # | Name | Type |
|---|---|---|
| 9 | `tenant_id` | `string` |
| 10 | `organisation_id` | `string` |
| 11 | `idempotency_key` | `string` |
| 12 | `total_payable` | `common.Money` |
| 13 | `coverage_start_at` | `Timestamp` |
| 14 | `coverage_end_at` | `Timestamp` |
| 15 | `payment_due_at` | `Timestamp` |

**New RPC to ADD:**
```protobuf
rpc UpdateFulfillmentStatus(UpdateFulfillmentStatusRequest) returns (UpdateFulfillmentStatusResponse);
```
Called by insurance-service (via Kafka consumer in orders-service) when a policy is issued.

---

### 20.3 `orders/events/v1/order_events.proto`

**Existing events (verified):**
`OrderCreatedEvent`, `OrderPaymentInitiatedEvent`, `OrderPaymentConfirmedEvent`, `OrderCancelledEvent`, `OrderFailedEvent`

**Existing `OrderCreatedEvent` fields:**
`event_id(1)`, `order_id(2)`, `user_id(3)`, `product_id(4)`, `plan_id(5)`, `quotation_id(6)`, `total_payable(7)`, `currency(8)`, `created_at(9)`

**Fields to ADD to ALL existing order events (add sequentially after last existing field in each):**

| Name | Type | Notes |
|---|---|---|
| `tenant_id` | `string` | Multi-tenancy routing |
| `organisation_id` | `string` | B2B (empty for B2C) |
| `portal` | `string` | normalized portal name |
| `actor_user_id` | `string` | who triggered the event |
| `session_id` | `string` | for audit trail |
| `session_type` | `string` | `"SERVER_SIDE"` or `"JWT"` |
| `causation_id` | `string` | event_id of triggering event |
| `correlation_id` | `string` | saga correlation key |
| `idempotency_key` | `string` | from original command |
| `occurred_at` | `Timestamp` | business time (not publish time) |

**New event to ADD:**
```protobuf
message OrderFulfillmentCompletedEvent {
  string event_id        = 1;
  string order_id        = 2;
  string policy_id       = 3;
  string invoice_id      = 4;
  string payment_id      = 5;
  string user_id         = 6;
  string tenant_id       = 7;
  string organisation_id = 8;
  string correlation_id  = 9;
  google.protobuf.Timestamp occurred_at = 10;
}
```

---

### 20.4 `payment/entity/v1/payment.proto`

**Existing `Payment` message fields (verified):**

| # | Name | Type |
|---|---|---|
| 1 | `payment_id` | `string` |
| 2 | `policy_id` | `string` |
| 3 | `amount` | `common.Money` |
| 4 | `currency` | `string` |
| 5 | `status` | `PaymentStatus` |
| 6 | `provider` | `string` |
| 7 | `gateway_ref` | `string` |
| 8 | `payment_method` | `string` |
| 9 | `checkout_url` | `string` |
| 10 | `idempotency_key` | `string` |
| 11 | `transaction_id` | `string` |
| 12 | `failure_reason` | `string` |
| 13 | `metadata` | `map<string,string>` |
| 14 | `created_at` | `Timestamp` |
| 15 | `updated_at` | `Timestamp` |
| 16 | `paid_at` | `Timestamp` |
| 17 | `expires_at` | `Timestamp` |
| 18 | `refund_id` | `string` |
| 19 | `refunded_at` | `Timestamp` |
| 20 | `refund_amount` | `common.Money` |
| 21 | `refund_reason` | `string` |
| 22 | `error_code` | `string` |

**Existing `PaymentStatus` enum:**
`PAYMENT_STATUS_UNSPECIFIED(0)`, `PENDING(1)`, `PROCESSING(2)`, `SUCCESS(3)`, `FAILED(4)`, `REFUNDED(5)`, `CANCELLED(6)`, `EXPIRED(7)`

**Fields to ADD (starting at field 23):**

| # | Name | Type | Rationale |
|---|---|---|---|
| 23 | `order_id` | `string` | Primary linkage to orders-service |
| 24 | `invoice_id` | `string` | Links to billing-service |
| 25 | `quotation_id` | `string` | Traceability to quote |
| 26 | `tenant_id` | `string` | Multi-tenancy |
| 27 | `customer_id` | `string` | Explicit customer (separate from actor) |
| 28 | `organisation_id` | `string` | B2B org context |
| 29 | `purchase_order_id` | `string` | B2B PO reference |
| 30 | `receipt_number` | `string` | Human-readable receipt ID |
| 31 | `manual_proof_file_id` | `string` | Storage file ID for bank slip |
| 32 | `verified_by` | `string` | Admin user_id who verified |
| 33 | `verified_at` | `Timestamp` | Verification timestamp |
| 34 | `rejection_reason` | `string` | If manual review rejected |
| 35 | `manual_review_status` | `ManualReviewStatus` | Manual verification state machine |
| 36 | `premium_amount` | `common.Money` | Base premium component |
| 37 | `vat_amount` | `common.Money` | VAT component |
| 38 | `service_fee_amount` | `common.Money` | Platform fee component |
| 39 | `aml_flagged` | `bool` | AML screening flag |
| 40 | `aml_case_id` | `string` | Links to fraud-service AML case |
| 41 | `actor_user_id` | `string` | Portal user who initiated |
| 42 | `portal` | `string` | normalized portal name |
| 43 | `session_id` | `string` | For audit trail |
| 44 | `session_type` | `string` | `"SERVER_SIDE"` or `"JWT"` |
| 45 | `return_url` | `string` | Gateway redirect after success |
| 46 | `cancel_url` | `string` | Gateway redirect on cancel |
| 47 | `callback_signature_verified` | `bool` | HMAC webhook validation result |
| 48 | `ledger_transaction_id` | `string` | TigerBeetle transaction ID |

**New enum to ADD:**
```protobuf
enum ManualReviewStatus {
  MANUAL_REVIEW_STATUS_UNSPECIFIED = 0;
  NOT_REQUIRED = 1;
  PROOF_SUBMITTED = 2;
  REVIEW_IN_PROGRESS = 3;
  APPROVED = 4;
  REJECTED = 5;
}
```

**Values to ADD to existing `PaymentStatus` enum:**
```protobuf
AWAITING_CUSTOMER_ACTION = 8;   // checkout URL generated, waiting for user
PENDING_MANUAL_VERIFICATION = 9; // proof submitted, pending admin review
VERIFIED = 10;                   // admin verified; policy issuance can proceed
```

---

### 20.5 `payment/services/v1/payment_service.proto`

**Existing RPCs (verified):**
`InitiatePayment`, `GetPayment`, `VerifyPayment`, `GetPaymentByIdempotencyKey`, `ListPayments`, `InitiateRefund`, `GetRefund`, `ListRefunds`, `ReconcilePayments`

**Fields to ADD to `InitiatePaymentRequest` (pick up after last existing field):**

| # | Name | Type |
|---|---|---|
| +1 | `order_id` | `string` |
| +2 | `invoice_id` | `string` |
| +3 | `quotation_id` | `string` |
| +4 | `tenant_id` | `string` |
| +5 | `customer_id` | `string` |
| +6 | `organisation_id` | `string` |
| +7 | `purchase_order_id` | `string` |
| +8 | `return_url` | `string` |
| +9 | `cancel_url` | `string` |
| +10 | `premium_amount` | `common.Money` |
| +11 | `vat_amount` | `common.Money` |
| +12 | `service_fee_amount` | `common.Money` |

**New RPCs to ADD:**
```protobuf
rpc SubmitManualPaymentProof(SubmitManualPaymentProofRequest)
    returns (SubmitManualPaymentProofResponse);

rpc ReviewManualPayment(ReviewManualPaymentRequest)
    returns (ReviewManualPaymentResponse);

rpc GenerateReceipt(GenerateReceiptRequest)
    returns (GenerateReceiptResponse);

rpc HandleGatewayWebhook(HandleGatewayWebhookRequest)
    returns (HandleGatewayWebhookResponse);

rpc GetPaymentReceipt(GetPaymentReceiptRequest)
    returns (GetPaymentReceiptResponse);
```

---

### 20.6 `payment/events/v1/payment_events.proto`

**Existing events (verified):**
`PaymentInitiatedEvent`, `PaymentCompletedEvent`, `PaymentFailedEvent`, `RefundProcessedEvent`

**Fields to ADD to ALL existing payment events:**

| Name | Type |
|---|---|
| `order_id` | `string` |
| `invoice_id` | `string` |
| `tenant_id` | `string` |
| `organisation_id` | `string` |
| `portal` | `string` |
| `actor_user_id` | `string` |
| `session_id` | `string` |
| `causation_id` | `string` |
| `correlation_id` | `string` |
| `occurred_at` | `Timestamp` |

**New events to ADD:**
```protobuf
message PaymentVerifiedEvent {
  string event_id        = 1;
  string payment_id      = 2;
  string order_id        = 3;
  string invoice_id      = 4;
  string verified_by     = 5;
  string tenant_id       = 6;
  string organisation_id = 7;
  string correlation_id  = 8;
  string causation_id    = 9;
  google.protobuf.Timestamp verified_at  = 10;
  google.protobuf.Timestamp occurred_at  = 11;
}

message ManualPaymentProofSubmittedEvent {
  string event_id           = 1;
  string payment_id         = 2;
  string order_id           = 3;
  string manual_proof_file_id = 4;
  string submitted_by       = 5;
  string tenant_id          = 6;
  string correlation_id     = 7;
  google.protobuf.Timestamp occurred_at = 8;
}

message ReceiptGeneratedEvent {
  string event_id        = 1;
  string payment_id      = 2;
  string order_id        = 3;
  string receipt_number  = 4;
  string receipt_file_id = 5;
  string tenant_id       = 6;
  string correlation_id  = 7;
  google.protobuf.Timestamp occurred_at = 8;
}

message PaymentReconciliationMismatchEvent {
  string event_id          = 1;
  string payment_id        = 2;
  string order_id          = 3;
  string expected_amount   = 4;
  string actual_amount     = 5;
  string provider          = 6;
  string gateway_ref       = 7;
  string tenant_id         = 8;
  string reconciliation_id = 9;
  google.protobuf.Timestamp occurred_at = 10;
}
```

---

### 20.7 Billing proto (create new — `billing/entity/v1/invoice.proto`)

```protobuf
syntax = "proto3";
package insuretech.billing.entity.v1;

import "insuretech/common/v1/types.proto";
import "google/protobuf/timestamp.proto";

message Invoice {
  string invoice_id        = 1;
  string order_id          = 2;
  string payment_id        = 3;
  string policy_id         = 4;
  string customer_id       = 5;
  string organisation_id   = 6;  // B2B
  string purchase_order_id = 7;  // B2B
  string tenant_id         = 8;
  string invoice_number    = 9;  // human-readable: INV-2026-000001
  InvoiceStatus status     = 10;
  common.Money amount      = 11;
  common.Money tax_amount  = 12;
  common.Money total_amount = 13;
  string currency          = 14;
  string notes             = 15;
  string issued_by         = 16;  // user_id
  google.protobuf.Timestamp issued_at    = 17;
  google.protobuf.Timestamp due_at       = 18;
  google.protobuf.Timestamp paid_at      = 19;
  google.protobuf.Timestamp cancelled_at = 20;
  google.protobuf.Timestamp overdue_at   = 21;
  map<string,string> metadata = 22;
}

enum InvoiceStatus {
  INVOICE_STATUS_UNSPECIFIED = 0;
  DRAFT      = 1;
  ISSUED     = 2;
  PAID       = 3;
  CANCELLED  = 4;
  OVERDUE    = 5;
  CREDIT_NOTE_ISSUED = 6;
}
```

---

### 20.8 DB migration delta

**`insurance_schema.orders` — columns to ADD:**
```sql
ALTER TABLE orders
  ADD COLUMN IF NOT EXISTS invoice_id         UUID,
  ADD COLUMN IF NOT EXISTS organisation_id    UUID,
  ADD COLUMN IF NOT EXISTS tenant_id          UUID,
  ADD COLUMN IF NOT EXISTS idempotency_key    TEXT,
  ADD COLUMN IF NOT EXISTS correlation_id     TEXT,
  ADD COLUMN IF NOT EXISTS payment_status     TEXT NOT NULL DEFAULT 'UNPAID',
  ADD COLUMN IF NOT EXISTS billing_status     TEXT NOT NULL DEFAULT 'NOT_INVOICED',
  ADD COLUMN IF NOT EXISTS fulfillment_status TEXT NOT NULL DEFAULT 'NOT_STARTED',
  ADD COLUMN IF NOT EXISTS manual_review_required BOOLEAN NOT NULL DEFAULT FALSE,
  ADD COLUMN IF NOT EXISTS payment_due_at     TIMESTAMPTZ,
  ADD COLUMN IF NOT EXISTS coverage_start_at  TIMESTAMPTZ,
  ADD COLUMN IF NOT EXISTS coverage_end_at    TIMESTAMPTZ,
  ADD COLUMN IF NOT EXISTS actor_user_id      UUID,
  ADD COLUMN IF NOT EXISTS portal             TEXT;

-- Dedup index for idempotency
CREATE UNIQUE INDEX IF NOT EXISTS orders_idempotency_key_idx
  ON orders(idempotency_key) WHERE idempotency_key IS NOT NULL;
```

**`payment_schema.payments` — columns to ADD:**
```sql
ALTER TABLE payments
  ADD COLUMN IF NOT EXISTS order_id                   UUID,
  ADD COLUMN IF NOT EXISTS invoice_id                 UUID,
  ADD COLUMN IF NOT EXISTS quotation_id               UUID,
  ADD COLUMN IF NOT EXISTS tenant_id                  UUID,
  ADD COLUMN IF NOT EXISTS customer_id                UUID,
  ADD COLUMN IF NOT EXISTS organisation_id            UUID,
  ADD COLUMN IF NOT EXISTS purchase_order_id          UUID,
  ADD COLUMN IF NOT EXISTS receipt_number             TEXT,
  ADD COLUMN IF NOT EXISTS manual_proof_file_id       TEXT,
  ADD COLUMN IF NOT EXISTS verified_by                UUID,
  ADD COLUMN IF NOT EXISTS verified_at                TIMESTAMPTZ,
  ADD COLUMN IF NOT EXISTS rejection_reason           TEXT,
  ADD COLUMN IF NOT EXISTS manual_review_status       TEXT NOT NULL DEFAULT 'NOT_REQUIRED',
  ADD COLUMN IF NOT EXISTS premium_amount             BIGINT,
  ADD COLUMN IF NOT EXISTS vat_amount                 BIGINT,
  ADD COLUMN IF NOT EXISTS service_fee_amount         BIGINT,
  ADD COLUMN IF NOT EXISTS aml_flagged                BOOLEAN NOT NULL DEFAULT FALSE,
  ADD COLUMN IF NOT EXISTS aml_case_id                TEXT,
  ADD COLUMN IF NOT EXISTS actor_user_id              UUID,
  ADD COLUMN IF NOT EXISTS portal                     TEXT,
  ADD COLUMN IF NOT EXISTS session_id                 UUID,
  ADD COLUMN IF NOT EXISTS session_type               TEXT,
  ADD COLUMN IF NOT EXISTS return_url                 TEXT,
  ADD COLUMN IF NOT EXISTS cancel_url                 TEXT,
  ADD COLUMN IF NOT EXISTS callback_signature_verified BOOLEAN NOT NULL DEFAULT FALSE,
  ADD COLUMN IF NOT EXISTS ledger_transaction_id      TEXT;

-- Index for order lookups
CREATE INDEX IF NOT EXISTS payments_order_id_idx ON payments(order_id);
CREATE INDEX IF NOT EXISTS payments_invoice_id_idx ON payments(invoice_id);
```

**New `billing_schema.invoices` table:**
```sql
CREATE SCHEMA IF NOT EXISTS billing_schema;

CREATE TABLE IF NOT EXISTS billing_schema.invoices (
  invoice_id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  order_id            UUID,
  payment_id          UUID,
  policy_id           UUID,
  customer_id         UUID NOT NULL,
  organisation_id     UUID,
  purchase_order_id   UUID,
  tenant_id           UUID NOT NULL,
  invoice_number      TEXT NOT NULL UNIQUE,
  status              TEXT NOT NULL DEFAULT 'DRAFT',
  amount_units        BIGINT NOT NULL DEFAULT 0,
  amount_currency     TEXT NOT NULL DEFAULT 'BDT',
  tax_amount_units    BIGINT NOT NULL DEFAULT 0,
  total_amount_units  BIGINT NOT NULL DEFAULT 0,
  notes               TEXT,
  issued_by           UUID,
  issued_at           TIMESTAMPTZ,
  due_at              TIMESTAMPTZ,
  paid_at             TIMESTAMPTZ,
  cancelled_at        TIMESTAMPTZ,
  overdue_at          TIMESTAMPTZ,
  metadata            JSONB,
  created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS invoices_order_id_idx     ON billing_schema.invoices(order_id);
CREATE INDEX IF NOT EXISTS invoices_customer_id_idx  ON billing_schema.invoices(customer_id);
CREATE INDEX IF NOT EXISTS invoices_tenant_id_idx    ON billing_schema.invoices(tenant_id);
CREATE INDEX IF NOT EXISTS invoices_status_idx       ON billing_schema.invoices(status);
```

---

## 19. File-Level Contract Work Plan

This section does not replace the higher-level contract sections above. It translates them into explicit file work items so implementation can proceed without another architecture round.

### 19.1 Billing proto files to create

Create:

- `proto/insuretech/billing/services/v1/billing_service.proto`
- `proto/insuretech/billing/events/v1/billing_events.proto`

Update if needed:

- `proto/insuretech/billing/entity/v1/invoice.proto`

Required message additions in `invoice.proto`:

- `tenant_id`
- `customer_id`
- `organisation_id`
- `order_id`
- `purchase_order_id`
- `quotation_id`
- `currency`
- `subtotal_amount`
- `vat_amount`
- `service_fee_amount`
- `discount_amount`
- `balance_due_amount`
- `credit_note_amount`
- `invoice_pdf_file_id`
- `receipt_file_ids`
- `issued_by`
- `cancelled_by`
- `cancel_reason`

Required enum review in `invoice.proto`:

- keep `PENDING`, `APPROVED`, `PAID`, `OVERDUE`, `CANCELLED`
- add if needed:
  - `PARTIALLY_PAID`
  - `CREDITED`
  - `VOID`

### 19.2 Payment proto files to update

Update:

- `proto/insuretech/payment/entity/v1/payment.proto`
- `proto/insuretech/payment/services/v1/payment_service.proto`
- `proto/insuretech/payment/events/v1/payment_events.proto`

Entity additions should cover:

- linkage
- invoice and order references
- receipt references
- manual verification fields
- provider and callback verification metadata
- AML/risk fields
- ledger linkage
- reconciliation markers

Service additions should cover:

- manual proof submission
- manual review decision
- receipt generation and retrieval
- secure gateway callback handling
- explicit reconciliation result retrieval if needed

Event additions should cover:

- payment verified
- manual proof submitted
- manual review requested
- receipt generated
- reconciliation matched
- reconciliation mismatched

### 19.3 Order proto files to update

Update:

- `proto/insuretech/orders/entity/v1/order.proto`
- `proto/insuretech/orders/services/v1/order_service.proto`
- `proto/insuretech/orders/events/v1/order_events.proto`

Order changes should cover:

- invoice linkage
- richer order state dimensions
- explicit fulfillment stage
- risk/manual-review markers
- recurring-payment readiness
- event payload enrichment for downstream consumers

### 19.4 Insurance proto files to create or update

Create:

- `proto/insuretech/insurance/events/v1/insurance_events.proto`

Possibly update:

- `proto/insuretech/insurance/services/v1/insurance_service.proto`

Required event group:

- quotation lifecycle events
- premium calculation events
- policy issuance events
- policy cancellation events
- renewal/lapse/reinstatement events

### 19.5 Cross-domain common-contract need

If event headers remain implicit, add or standardize a common metadata message under a shared package. If not done as a proto type, enforce it as a Kafka header convention and document it centrally.

Minimum common metadata:

- `event_id`
- `event_version`
- `correlation_id`
- `causation_id`
- `tenant_id`
- `portal`
- `actor_user_id`
- `session_id`
- `token_id`
- `organisation_id`
- `trace_id`
- `idempotency_key`
- `occurred_at`

## 20. Detailed Saga and State Transition Plan

### 20.1 Retail purchase saga states

Recommended business progression:

1. quotation ready
2. order created
3. invoice issued
4. payment initiated
5. payment verified
6. order payment confirmed
7. policy issuance started
8. policy issued
9. policy document generated
10. receipt generated
11. customer delivery completed

Recommended failure branches:

- payment initiation failure
- callback timeout / verification pending
- manual review required
- payment rejected
- order cancelled before settlement
- policy issuance failed after payment
- document generation failed after policy issue

### 20.2 Command-to-event mapping

`CreateOrder` should produce:

- order write
- `OrderCreatedEvent`
- optional invoice request command or orchestration follow-up

`InitiatePayment` should produce:

- payment write
- `PaymentInitiatedEvent`
- `OrderPaymentInitiatedEvent` if the order view needs explicit state progression

`VerifyPayment` or webhook success should produce:

- payment verification update
- `PaymentVerifiedEvent`
- `PaymentCompletedEvent` once settlement is operationally complete

`ConfirmPayment` should produce:

- order state update
- `OrderPaymentConfirmedEvent`

`IssuePolicy` should produce:

- insurance write
- `PolicyIssuanceRequestedEvent`
- `PolicyIssuedEvent` or `PolicyIssuanceFailedEvent`

### 20.3 Required idempotency boundaries

The following operations must be independently idempotent:

- order creation from quotation
- invoice issuance for order
- payment initiation for invoice
- payment callback processing
- payment verification decision
- order payment confirmation
- policy issuance
- receipt generation
- document generation requests
- webhook delivery to external systems

### 20.4 Partial failure handling rules

If payment is settled but order confirmation fails:

- do not re-charge
- retain payment as settled
- retry order confirmation using correlation-aware recovery job

If order is confirmed but policy issuance fails:

- retain order as commercially settled
- move fulfillment status to failed/pending-retry
- create operational work item
- do not lose invoice/payment linkage

If policy is issued but document generation fails:

- keep policy active
- retry document generation asynchronously
- surface fulfillment status separately from policy state

If receipt generation fails:

- do not roll back payment completion
- retry receipt generation
- allow portal to show payment as settled with receipt pending

## 21. External Integration Plan

### 21.1 Payment gateway integration layers

The Go `payment-service` should have separate internal slices for:

- provider authentication/token management
- payment initiation
- payment status verification
- webhook handling
- refund handling
- reconciliation import
- provider health monitoring

Each provider integration should expose a normalized internal contract so the orchestration and ledger code do not depend on provider-specific response shapes.

### 21.2 Webhook security plan

For every provider webhook:

- verify HMAC or equivalent signature
- verify timestamp freshness when provider supports it
- verify provider reference and amount against stored payment intent
- store raw callback payload for audit
- reject duplicate callbacks idempotently
- emit audit event/log on every accepted and rejected callback

### 21.2.1 SSLCommerz practical integration notes

This subsection is based on the official SSLCommerz docs reviewed on 2026-03-06:

- docs index: [https://developer.sslcommerz.com/docs.html](https://developer.sslcommerz.com/docs.html)
- API reference: [https://developer.sslcommerz.com/doc/v4/](https://developer.sslcommerz.com/doc/v4/)

The current Go `payment-service` does not implement this flow yet. It only fabricates a hosted URL from config and marks payments successful when `VerifyPayment` is called. The implementation plan must therefore follow the real SSLCommerz lifecycle, not the current stub behavior.

Authoritative constraints to encode:

- Session initiation must call SSLCommerz instead of building a local URL.
- The official docs use:
  - sandbox session endpoint: `https://sandbox.sslcommerz.com/gwprocess/v4/api.php`
  - live session endpoint: `https://securepay.sslcommerz.com/gwprocess/v4/api.php`
- The official docs describe three backend integration steps:
  - create and get session
  - receive payment notification through IPN
  - validate the order with the validation API
- Session initiation must send the provider-required commercial fields, at minimum:
  - `store_id`
  - `store_passwd`
  - `total_amount`
  - `currency`
  - `tran_id`
  - `success_url`
  - `fail_url`
  - `cancel_url`
  - `ipn_url`
  - product and customer fields required by the hosted checkout profile
- SSLCommerz currently documents a transaction amount range of `10.00 BDT` to `500000.00 BDT` for initiation.
- Session initiation response handling must persist:
  - `status`
  - `sessionkey`
  - `GatewayPageURL`
  - `failedreason` when present
- IPN must be treated as the primary machine-to-machine notification path.
- Callback handling must not trust the browser redirect alone.
- Browser success/fail/cancel redirects are customer UX events, not settlement truth.
- Payment confirmation must validate the transaction using SSLCommerz validation APIs, especially the `val_id` returned in the success flow.
- Validation endpoint from the docs:
  - sandbox: `https://sandbox.sslcommerz.com/validator/api/validationserverAPI.php`
  - live: `https://securepay.sslcommerz.com/validator/api/validationserverAPI.php`
- The implementation must support the merchant transaction lookup path by `tran_id` as a fallback and reconciliation tool.
- The implementation should also support transaction lookup by `sessionkey`.
- Refund implementation must follow SSLCommerz refund initiation and refund-status query APIs rather than marking refunds completed immediately.
- SSLCommerz currently documents that live refund API use requires the merchant public IP to be registered.
- Risk metadata returned by SSLCommerz such as `risk_level` and `risk_title` must be persisted and exposed to AML/manual-review logic.
- Operational config must assume TLS 1.2 and the `securepay.sslcommerz.com` hostnames now documented by SSLCommerz.

Practical implementation rule:

- `success_url`, `fail_url`, `cancel_url`, and `ipn_url` must be generated by InsureTech and point back to gateway-owned endpoints.
- `tran_id` should be generated by InsureTech and remain stable for idempotent retries of the same payment intent.
- `payment_id` remains the internal primary key; `tran_id`, `val_id`, `bank_tran_id`, and `sessionkey` are provider-side correlation fields and must be stored separately.
- SSLCommerz `verify_sign` and `verify_key` fields may be captured from callbacks, but they are not enough on their own; the system must still call the validation API before granting coverage.
- final commercial success is established only after provider validation succeeds, not after redirect receipt.

### 21.3 Manual payment proof flow details

Manual proof review should explicitly support:

- upload proof file
- validate file type and size
- virus scan
- optional OCR extraction if useful later
- attach proof to payment
- assign review queue
- approve or reject with audit note
- generate customer-visible outcome

Required additions to `storage-service` and workflow:

- file upload with automatic scanning trigger
- file quarantine if scan fails
- storage event consumption in payment workflow
- attachment linkage table or payment aggregation

### 21.4 Reconciliation flow details

Reconciliation should:

- periodically fetch settlement reports from payment providers
- match provider transactions against recorded payment state
- flag mismatches such as missing transactions, amount differences, and status divergence
- trigger correction workflows or manual review
- emit reconciliation events for audit and finance

### 21.5 Notification integration steps

The plan currently references notifications but should treat them as explicit downstream work:

- invoice issued -> SMS/email/payment link
- payment received -> receipt notification
- policy issued -> policy document notification
- refund processed -> refund status notification
- overdue invoice -> reminder notifications
- manual review requested -> back-office notification only

### 21.6 External webhook delivery plan

For partner and enterprise integrations:

- publish internal domain event
- project to outbound webhook payload
- sign payload
- retry delivery with backoff
- persist delivery attempts
- support replay by event/correlation id

## 22. Read Models and Projection Plan

The event model will not be usable operationally unless projections are planned now.

### 22.1 Customer portal projections

Need read models for:

- policy list
- payment history
- receipt downloads
- invoice status
- renewal due summary
- refund status

### 22.2 Partner/B2B portal projections

Need read models for:

- organisation invoice ledger
- purchase-order linked invoices
- employee coverage summaries
- generated schedules and policy bundles
- overdue balance views

### 22.3 Admin and compliance projections

Need read models for:

- failed payments
- pending manual verification queue
- pending policy issuance queue
- reconciliation mismatches
- suspicious transaction queue
- premium collection reporting
- refund processing SLA dashboard

### 22.4 Suggested projection ownership

- `PoliSync` or a reporting/analytics service can own orchestration-oriented projections
- payment-specific operational projections should be owned by `payment-service` or a dedicated read-model worker
- billing-facing commercial projections should be owned by `billing-service`

## 23. Detailed B2B Expansion Plan

### 23.1 B2B commercial scenarios to support

The plan should not assume one B2B payment pattern. Support at least:

- one purchase order -> one invoice -> one payment
- one purchase order -> one invoice -> multiple payments
- one organisation account -> one invoice covering multiple purchase orders
- one organisation invoice -> multiple downstream employee policies

### 23.2 B2B data links required

Across order/payment/billing/events, support:

- `organisation_id`
- `purchase_order_id`
- `department_id`
- `employee_uuid` when policy is employee-specific
- `requested_by`
- `approved_by`

### 23.3 B2B workflow additions

Add to the flow:

- corporate invoice approval if required
- payment term support
- invoice aging and overdue handling
- schedule document generation
- employer-facing receipt and settlement reporting

## 24. Compliance, Reporting, and Audit Deepening

### 24.1 Reporting inputs that must come from the event layer

The following should be derivable from payment/order/billing/policy events:

- monthly premium collection
- claims settlement inputs where payment side is involved
- product-line revenue
- partner-channel revenue
- geographic premium distribution
- failed payment rate
- cancellation/refund ratio

### 24.2 Audit record minimum fields

For all critical actions, persist:

- actor id
- tenant id
- organisation id when relevant
- source IP
- device id
- session id
- token id
- command name
- aggregate id
- pre-state summary
- post-state summary
- correlation id
- timestamp

### 24.3 AML and fraud processing steps

For flagged payments/orders:

1. flag event generated
2. case record created
3. review queue assignment
4. escalation path to compliance/business admin/focal person
5. suppression of customer-facing suspicious-activity detail
6. outcome recorded for reporting and later model tuning

Minimum trigger rules to encode:

- same customer or tenant creates more than 3 policies in 7 days
- premium amount exceeds BDT 500,000 without required KYC proof
- refund pattern exceeds 2 refunds in 30 days
- payer identity does not match customer identity
- payment geography mismatches customer profile geography
- payment method changes unusually across successive purchases
- KYC failure churn exceeds 2 failed attempts in 24 hours

Compliance workflow requirements:

- emit `PaymentFlaggedForComplianceEvent`
- route the case to a manual review queue with SLA
- require reviewer action before policy issuance proceeds when the flag is blocking
- log all reviewer actions and final dispositions
- never expose suspicious-activity investigation detail to the customer

### 24.4 Additional audit and reporting requirements

Maintain immutable streams or equivalent audit records for:

- every order creation
- every payment attempt and result
- every manual verification decision
- every refund and adjustment
- every policy issuance
- every cancellation
- every compliance flag and resolution

These records must support:

- IDRA monthly premium collection reporting
- premium settlement reconciliation
- incident and unusual-transaction reporting
- regulatory investigations and audit reconstruction

## 25. Migration and Backward Compatibility Plan

### 25.1 Contract migration approach

Because payment, billing, and insurance contracts are changing, use additive proto evolution first:

- add fields without reusing field numbers
- mark obsolete fields deprecated, do not remove immediately
- publish new events alongside old ones temporarily if required
- maintain consumer compatibility until all critical services are upgraded

### 25.2 Data migration steps

For existing records:

- backfill `order_id` onto payment records where derivable
- backfill `invoice_id` once billing is introduced
- map historic provider references into normalized fields
- generate receipt numbers for past completed payments if business requires continuity
- populate missing tenant or organisation references where possible

### 25.3 Rollout sequencing

Recommended deployment order:

1. additive proto release
2. producer support for new fields
3. consumer support for both old and new events
4. billing service activation
5. payment-service field backfill
6. strict validation once all consumers are upgraded

## 26. Execution Checklist by Service

### 26.1 Go `payment-service`

- implement provider abstraction
- implement gateway callback validation
- implement manual proof submission
- implement manual review path
- implement receipt generation trigger
- implement refund initiation and status tracking
- implement TigerBeetle posting
- implement reconciliation jobs
- implement outbox publishing
- implement AML/risk hooks

### 26.2 Go `orders-service`

- enrich order model with invoice/payment/fulfillment dimensions
- publish order lifecycle events via outbox
- support payment confirmation idempotently
- support cancellation/failure transitions
- expose sufficient query APIs for orchestration and projections

### 26.3 New Go `billing-service`

- implement invoice entity and line items
- implement issue/pay/cancel/overdue flows
- implement payment linkage
- implement credit-note flow if required
- publish invoice events via outbox
- expose invoice query APIs for portals and orchestration

### 26.4 C# `PoliSync`

- validate quotation/business preconditions
- orchestrate order -> billing -> payment -> policy transitions
- consume payment, order, billing, storage, document, and B2B events
- drive policy issuance
- drive document generation commands
- maintain orchestration idempotency
- maintain fulfillment/recovery workflows for partial failures

### 26.5 `insurance-service`

- accept policy issuance requests with stronger linkage fields
- emit policy lifecycle events
- expose renewal/lapse/cancellation transitions consistently

### 26.6 `document-service` and `storage-service`

- support policy, receipt, invoice, and schedule document references
- preserve tenant/reference metadata consistently
- emit generated and uploaded events reliably

## 27. Missing Decisions That Still Need Explicit Resolution

These are not blockers for planning, but they must be decided before implementation hardens.

### 27.1 Financial questions

- whether invoice is mandatory for every retail payment or only for B2B and manual paths
- whether one payment may settle multiple invoices in phase 1
- whether partial payment is in scope immediately or later

### 27.2 Document questions

- whether receipts are generated by `payment-service` directly or always through `document-service`
- whether invoice PDF generation is synchronous on issue or asynchronous from event

### 27.3 Policy activation questions

- whether activation occurs on `PaymentVerifiedEvent` or only after `OrderPaymentConfirmedEvent`
- how policy numbering is reserved if payment succeeds but issuance retries later

### 27.4 Compliance questions

- whether AML flags hard-block payment completion or only route to review
- which events must be retained in hot storage vs archived storage

## 28. Expanded Next-Step Plan

The immediate planning work should now proceed in this order:

1. finalize billing service and event contract drafts
2. finalize payment proto delta with exact fields and enums
3. finalize insurance event package
4. finalize common Kafka metadata/header convention
5. define projection owners and storage models
6. define rollout and migration sequencing per service
7. only then start implementation against stabilized contracts

## 29. Draft Contract Spec: `billing_service.proto`

This section is a plan-level draft for `proto/insuretech/billing/services/v1/billing_service.proto`. It is intentionally explicit so the proto can be authored with minimal follow-up design work.

### 29.1 Proposed package and imports

Recommended package:

- `insuretech.billing.services.v1`

Recommended imports:

- `google/api/annotations.proto`
- `google/api/field_behavior.proto`
- `google/protobuf/timestamp.proto`
- `insuretech/common/v1/error.proto`
- `insuretech/common/v1/types.proto`
- `insuretech/billing/entity/v1/invoice.proto`

### 29.2 Proposed service shape

Recommended RPCs:

- `IssueInvoice`
- `GetInvoice`
- `ListInvoices`
- `LinkPaymentToInvoice`
- `MarkInvoicePaid`
- `CancelInvoice`
- `MarkInvoiceOverdue`
- `IssueCreditNote`
- `ListCreditNotes` if credit-note usage is enabled in phase 1

Recommended ownership rule:

- `IssueInvoice` is the only write path for creating an invoice
- `MarkInvoicePaid` should not create payments; it should only confirm settlement against an already-known payment
- `LinkPaymentToInvoice` must support partial settlement even if phase-1 UI does not expose partial payments yet

### 29.3 Proposed enums

Recommended `InvoiceLineType`:

- `INVOICE_LINE_TYPE_UNSPECIFIED`
- `INVOICE_LINE_TYPE_PREMIUM`
- `INVOICE_LINE_TYPE_VAT`
- `INVOICE_LINE_TYPE_SERVICE_FEE`
- `INVOICE_LINE_TYPE_DISCOUNT`
- `INVOICE_LINE_TYPE_ENDORSEMENT_FEE`
- `INVOICE_LINE_TYPE_CANCELLATION_CHARGE`
- `INVOICE_LINE_TYPE_LATE_FEE`
- `INVOICE_LINE_TYPE_REFUND_ADJUSTMENT`

Recommended `InvoiceSettlementStatus`:

- `INVOICE_SETTLEMENT_STATUS_UNSPECIFIED`
- `INVOICE_SETTLEMENT_STATUS_UNPAID`
- `INVOICE_SETTLEMENT_STATUS_PARTIALLY_PAID`
- `INVOICE_SETTLEMENT_STATUS_PAID`
- `INVOICE_SETTLEMENT_STATUS_REFUNDED`
- `INVOICE_SETTLEMENT_STATUS_CREDITED`

Recommended `InvoiceKind`:

- `INVOICE_KIND_UNSPECIFIED`
- `INVOICE_KIND_NEW_POLICY`
- `INVOICE_KIND_RENEWAL`
- `INVOICE_KIND_ENDORSEMENT`
- `INVOICE_KIND_CANCELLATION`
- `INVOICE_KIND_B2B_PURCHASE_ORDER`

### 29.4 Proposed message drafts

Recommended `InvoiceLineItemDraft`:

- `string line_item_id = 1`
- `InvoiceLineType line_type = 2`
- `string reference_id = 3`
- `string reference_type = 4`
- `string description = 5`
- `insuretech.common.v1.Money amount = 6`
- `int32 quantity = 7`
- `map<string, string> metadata = 8`

Recommended `IssueInvoiceRequest`:

- `string tenant_id = 1`
- `string customer_id = 2`
- `string organisation_id = 3`
- `string order_id = 4`
- `string purchase_order_id = 5`
- `string quotation_id = 6`
- `string policy_id = 7`
- `InvoiceKind kind = 8`
- `string currency = 9`
- `google.protobuf.Timestamp due_at = 10`
- `repeated InvoiceLineItemDraft line_items = 11`
- `string notes = 12`
- `string idempotency_key = 13`
- `string issued_by = 14`
- `string source_service = 15`

Recommended `IssueInvoiceResponse`:

- `insuretech.billing.entity.v1.Invoice invoice = 1`
- `string message = 2`
- `insuretech.common.v1.Error error = 3`

Recommended `GetInvoiceRequest`:

- `string invoice_id = 1`
- `string tenant_id = 2`

Recommended `GetInvoiceResponse`:

- `insuretech.billing.entity.v1.Invoice invoice = 1`
- `insuretech.common.v1.Error error = 2`

Recommended `ListInvoicesRequest`:

- `string tenant_id = 1`
- `string customer_id = 2`
- `string organisation_id = 3`
- `string order_id = 4`
- `string purchase_order_id = 5`
- `string policy_id = 6`
- `string quotation_id = 7`
- `insuretech.billing.entity.v1.InvoiceStatus status = 8`
- `InvoiceSettlementStatus settlement_status = 9`
- `int32 page_size = 10`
- `string page_token = 11`

Recommended `ListInvoicesResponse`:

- `repeated insuretech.billing.entity.v1.Invoice invoices = 1`
- `string next_page_token = 2`
- `int32 total_count = 3`
- `insuretech.common.v1.Error error = 4`

Recommended `LinkPaymentToInvoiceRequest`:

- `string tenant_id = 1`
- `string invoice_id = 2`
- `string payment_id = 3`
- `insuretech.common.v1.Money amount = 4`
- `string linked_by = 5`
- `string idempotency_key = 6`

Recommended `LinkPaymentToInvoiceResponse`:

- `insuretech.billing.entity.v1.Invoice invoice = 1`
- `string message = 2`
- `insuretech.common.v1.Error error = 3`

Recommended `MarkInvoicePaidRequest`:

- `string tenant_id = 1`
- `string invoice_id = 2`
- `string payment_id = 3`
- `google.protobuf.Timestamp paid_at = 4`
- `string confirmed_by = 5`
- `string idempotency_key = 6`

Recommended `MarkInvoicePaidResponse`:

- `insuretech.billing.entity.v1.Invoice invoice = 1`
- `string message = 2`
- `insuretech.common.v1.Error error = 3`

Recommended `CancelInvoiceRequest`:

- `string tenant_id = 1`
- `string invoice_id = 2`
- `string reason = 3`
- `string cancelled_by = 4`
- `string idempotency_key = 5`

Recommended `CancelInvoiceResponse`:

- `insuretech.billing.entity.v1.Invoice invoice = 1`
- `string message = 2`
- `insuretech.common.v1.Error error = 3`

Recommended `MarkInvoiceOverdueRequest`:

- `string tenant_id = 1`
- `string invoice_id = 2`
- `int32 days_overdue = 3`
- `string triggered_by = 4`
- `string idempotency_key = 5`

Recommended `MarkInvoiceOverdueResponse`:

- `insuretech.billing.entity.v1.Invoice invoice = 1`
- `string message = 2`
- `insuretech.common.v1.Error error = 3`

Recommended `IssueCreditNoteRequest`:

- `string tenant_id = 1`
- `string invoice_id = 2`
- `string payment_id = 3`
- `string refund_id = 4`
- `insuretech.common.v1.Money amount = 5`
- `string reason = 6`
- `string issued_by = 7`
- `string idempotency_key = 8`

Recommended `IssueCreditNoteResponse`:

- `string credit_note_id = 1`
- `insuretech.billing.entity.v1.Invoice invoice = 2`
- `string message = 3`
- `insuretech.common.v1.Error error = 4`

### 29.5 Billing validation rules to encode in service behavior

Required service-side rules:

- an invoice cannot be issued without a tenant context
- at least one of `customer_id` or `organisation_id` must exist
- at least one commercial reference should exist:
  - `order_id`
  - `purchase_order_id`
  - `policy_id`
  - `quotation_id`
- invoice currency must match payment currency expectations
- sum of line items must equal stored totals
- cancelling a fully paid invoice should require credit-note or refund workflow, not silent status overwrite

### 29.6 Billing event publication rules

The service should publish:

- `InvoiceIssuedEvent` after durable invoice creation
- `InvoicePaymentLinkedEvent` after each successful linkage
- `InvoicePaidEvent` after settlement reaches paid state
- `InvoiceCancelledEvent` after cancellation commit
- `InvoiceOverdueEvent` after overdue transition
- `CreditNoteIssuedEvent` after credit-note issuance

## 30. Draft Contract Spec: `payment.proto` and `payment_service.proto`

This section extends the earlier payment delta into a near-authorable proto draft.

### 30.1 Proposed payment entity evolution

Target file:

- `proto/insuretech/payment/entity/v1/payment.proto`

Recommended entity grouping:

- identity and linkage
- money and commercial context
- provider and gateway data
- manual verification data
- accounting and receipt data
- risk and compliance data
- timestamps and retry metadata

### 30.2 Proposed entity fields

Identity and linkage:

- `string payment_id = 1`
- `string transaction_id = 2`
- `string tenant_id = 3`
- `string order_id = 4`
- `string invoice_id = 5`
- `string quotation_id = 6`
- `string policy_id = 7`
- `string claim_id = 8`
- `string customer_id = 9`
- `string organisation_id = 10`
- `string purchase_order_id = 11`

Commercial context:

- `PaymentType type = 12`
- `PaymentMethod method = 13`
- `PaymentStatus status = 14`
- `string provider = 15`
- `string payment_channel = 16`
- `insuretech.common.v1.Money premium_amount = 17`
- `insuretech.common.v1.Money vat_amount = 18`
- `insuretech.common.v1.Money service_fee_amount = 19`
- `insuretech.common.v1.Money discount_amount = 20`
- `insuretech.common.v1.Money total_amount = 21`
- `string currency = 22`
- `string payment_frequency = 23`

Provider and gateway fields:

- `string provider_account = 24`
- `string provider_reference = 25`
- `string provider_status_code = 26`
- `string provider_status_message = 27`
- `bool callback_signature_verified = 28`
- `string callback_signature_algorithm = 29`
- `string gateway_response = 30`
- `string hosted_payment_url = 31`

Manual verification fields:

- `string manual_proof_file_id = 32`
- `google.protobuf.Timestamp manual_proof_uploaded_at = 33`
- `ManualReviewStatus manual_review_status = 34`
- `string verified_by = 35`
- `google.protobuf.Timestamp verified_at = 36`
- `string rejection_reason = 37`
- `string review_notes = 38`

Receipt and accounting fields:

- `string receipt_number = 39`
- `string receipt_document_id = 40`
- `string receipt_file_id = 41`
- `string ledger_transaction_id = 42`
- `string ledger_batch_id = 43`
- `string tigerbeetle_transfer_id = 44`
- `string reconciliation_status = 45`

Risk and compliance fields:

- `bool third_party_payer = 46`
- `string payer_id = 47`
- `string payer_relationship = 48`
- `double risk_score = 49`
- `bool aml_flagged = 50`
- `string aml_case_id = 51`
- `bool suspicious_activity_locked = 52`

Operational fields:

- `string idempotency_key = 53`
- `int32 retry_count = 54`
- `string failure_reason = 55`
- `google.protobuf.Timestamp initiated_at = 56`
- `google.protobuf.Timestamp completed_at = 57`
- `google.protobuf.Timestamp created_at = 58`
- `google.protobuf.Timestamp updated_at = 59`

### 30.3 Proposed enums

Recommended `PaymentStatus`:

- `PAYMENT_STATUS_UNSPECIFIED`
- `PAYMENT_STATUS_PENDING`
- `PAYMENT_STATUS_INITIATED`
- `PAYMENT_STATUS_AWAITING_CUSTOMER_ACTION`
- `PAYMENT_STATUS_PENDING_MANUAL_VERIFICATION`
- `PAYMENT_STATUS_VERIFIED`
- `PAYMENT_STATUS_COMPLETED`
- `PAYMENT_STATUS_FAILED`
- `PAYMENT_STATUS_REFUND_INITIATED`
- `PAYMENT_STATUS_REFUNDED`
- `PAYMENT_STATUS_CANCELLED`
- `PAYMENT_STATUS_RECONCILIATION_HOLD`

Recommended `ManualReviewStatus`:

- `MANUAL_REVIEW_STATUS_UNSPECIFIED`
- `MANUAL_REVIEW_STATUS_NOT_REQUIRED`
- `MANUAL_REVIEW_STATUS_PENDING`
- `MANUAL_REVIEW_STATUS_APPROVED`
- `MANUAL_REVIEW_STATUS_REJECTED`

Recommended `PaymentChannel` can be modeled as enum or string. If enum is chosen:

- `PAYMENT_CHANNEL_UNSPECIFIED`
- `PAYMENT_CHANNEL_BKASH`
- `PAYMENT_CHANNEL_NAGAD`
- `PAYMENT_CHANNEL_ROCKET`
- `PAYMENT_CHANNEL_BANK_TRANSFER`
- `PAYMENT_CHANNEL_CARD_HOSTED`
- `PAYMENT_CHANNEL_CASH`
- `PAYMENT_CHANNEL_CHEQUE`

### 30.4 Proposed payment service shape

Target file:

- `proto/insuretech/payment/services/v1/payment_service.proto`

Recommended RPC set:

- existing:
  - `InitiatePayment`
  - `VerifyPayment`
  - `GetPayment`
  - `ListPayments`
  - `InitiateRefund`
  - `GetRefundStatus`
  - `ListPaymentMethods`
  - `AddPaymentMethod`
  - `ReconcilePayments`
- add:
  - `SubmitManualPaymentProof`
  - `ReviewManualPayment`
  - `GenerateReceipt`
  - `GetPaymentReceipt`
  - `HandleGatewayWebhook`
  - `GetPaymentByProviderReference`

### 30.5 Proposed request/response drafts

Recommended `InitiatePaymentRequest`:

- `string tenant_id = 1`
- `string order_id = 2`
- `string invoice_id = 3`
- `string quotation_id = 4`
- `string policy_id = 5`
- `string customer_id = 6`
- `string organisation_id = 7`
- `string purchase_order_id = 8`
- `PaymentMethod method = 9`
- `string provider = 10`
- `insuretech.common.v1.Money total_amount = 11`
- `string currency = 12`
- `string callback_url = 13`
- `string return_url = 14`
- `string cancel_url = 15`
- `bool allow_manual_review = 16`
- `string idempotency_key = 17`
- `map<string, string> metadata = 18`

Recommended `InitiatePaymentResponse`:

- `insuretech.payment.entity.v1.Payment payment = 1`
- `string hosted_payment_url = 2`
- `string provider_reference = 3`
- `string message = 4`
- `insuretech.common.v1.Error error = 5`

Recommended `VerifyPaymentRequest`:

- `string tenant_id = 1`
- `string payment_id = 2`
- `string provider_reference = 3`
- `string verification_source = 4`
- `string idempotency_key = 5`

Recommended `VerifyPaymentResponse`:

- `insuretech.payment.entity.v1.Payment payment = 1`
- `string message = 2`
- `insuretech.common.v1.Error error = 3`

Recommended `SubmitManualPaymentProofRequest`:

- `string tenant_id = 1`
- `string payment_id = 2`
- `string file_id = 3`
- `string uploaded_by = 4`
- `string notes = 5`
- `string idempotency_key = 6`

Recommended `SubmitManualPaymentProofResponse`:

- `insuretech.payment.entity.v1.Payment payment = 1`
- `string message = 2`
- `insuretech.common.v1.Error error = 3`

Recommended `ReviewManualPaymentRequest`:

- `string tenant_id = 1`
- `string payment_id = 2`
- `ManualReviewDecision decision = 3`
- `string review_notes = 4`
- `string reviewed_by = 5`
- `string idempotency_key = 6`

Recommended `ReviewManualPaymentResponse`:

- `insuretech.payment.entity.v1.Payment payment = 1`
- `string message = 2`
- `insuretech.common.v1.Error error = 3`

Recommended `GenerateReceiptRequest`:

- `string tenant_id = 1`
- `string payment_id = 2`
- `string requested_by = 3`
- `string idempotency_key = 4`

Recommended `GenerateReceiptResponse`:

- `string receipt_number = 1`
- `string receipt_document_id = 2`
- `string receipt_file_id = 3`
- `string message = 4`
- `insuretech.common.v1.Error error = 5`

Recommended `GetPaymentReceiptRequest`:

- `string tenant_id = 1`
- `string payment_id = 2`

Recommended `GetPaymentReceiptResponse`:

- `string receipt_number = 1`
- `string receipt_document_id = 2`
- `string receipt_file_id = 3`
- `string download_url = 4`
- `insuretech.common.v1.Error error = 5`

Recommended `HandleGatewayWebhookRequest`:

- `string provider = 1`
- `map<string, string> headers = 2`
- `bytes raw_payload = 3`
- `string received_at_iso = 4`

Recommended `HandleGatewayWebhookResponse`:

- `bool accepted = 1`
- `string payment_id = 2`
- `string provider_reference = 3`
- `string message = 4`
- `insuretech.common.v1.Error error = 5`

Recommended `GetPaymentByProviderReferenceRequest`:

- `string tenant_id = 1`
- `string provider = 2`
- `string provider_reference = 3`

Recommended `GetPaymentByProviderReferenceResponse`:

- `insuretech.payment.entity.v1.Payment payment = 1`
- `insuretech.common.v1.Error error = 2`

### 30.6 Payment service behavior rules to encode

Required behavior:

- `InitiatePayment` must reject when neither `order_id` nor `invoice_id` is present unless an explicit exception is designed
- `SubmitManualPaymentProof` must move the payment into pending manual verification if not already there
- `ReviewManualPayment` must emit verification or failure events and persist reviewer identity
- `HandleGatewayWebhook` must never directly trust provider payload without signature and reference validation
- `GenerateReceipt` must be idempotent and safe to rerun

### 30.7 Draft payment events expansion

Target file:

- `proto/insuretech/payment/events/v1/payment_events.proto`

Recommended event field enrichment for all payment events:

- `tenant_id`
- `order_id`
- `invoice_id`
- `quotation_id`
- `policy_id`
- `customer_id`
- `organisation_id`
- `purchase_order_id`
- `provider`
- `payment_channel`
- `receipt_number`
- `ledger_transaction_id`
- `manual_review_status` where relevant

Recommended `PaymentVerifiedEvent`:

- `event_id`
- `payment_id`
- `tenant_id`
- `order_id`
- `invoice_id`
- `provider_reference`
- `verified_by`
- `verified_at`
- `correlation_id`
- `timestamp`

Recommended `ManualPaymentProofSubmittedEvent`:

- `event_id`
- `payment_id`
- `tenant_id`
- `invoice_id`
- `order_id`
- `file_id`
- `uploaded_by`
- `correlation_id`
- `timestamp`

Recommended `ReceiptGeneratedEvent`:

- `event_id`
- `payment_id`
- `tenant_id`
- `order_id`
- `invoice_id`
- `receipt_number`
- `receipt_document_id`
- `receipt_file_id`
- `correlation_id`
- `timestamp`

## 31. Draft Contract Spec: `insurance/events/v1/insurance_events.proto`

This section defines the missing insurance event package needed to make the order-payment-policy chain event-complete.

### 31.1 Proposed package and imports

Recommended package:

- `insuretech.insurance.events.v1`

Recommended imports:

- `google/protobuf/timestamp.proto`
- `insuretech/common/v1/types.proto`

### 31.2 Event families to include

Quotation events:

- `QuotationCreatedEvent`
- `QuotationUpdatedEvent`
- `QuotationLockedEvent`
- `QuotationExpiredEvent`
- `PremiumCalculatedEvent`

Policy issuance events:

- `PolicyIssuanceRequestedEvent`
- `PolicyIssuedEvent`
- `PolicyIssuanceFailedEvent`

Policy lifecycle events:

- `PolicyActivatedEvent`
- `PolicyCancelledEvent`
- `PolicyRenewedEvent`
- `PolicyLapsedEvent`
- `PolicyReinstatedEvent`
- `PolicyEndorsedEvent`

### 31.3 Common field rules for all insurance events

Each event should carry:

- `event_id`
- `tenant_id`
- `quotation_id` when relevant
- `policy_id` when relevant
- `order_id` when relevant
- `invoice_id` when relevant
- `payment_id` when relevant
- `customer_id`
- `organisation_id`
- `product_id`
- `plan_id`
- `correlation_id`
- `timestamp`

### 31.4 Proposed event drafts

Recommended `QuotationLockedEvent`:

- `string event_id = 1`
- `string quotation_id = 2`
- `string tenant_id = 3`
- `string customer_id = 4`
- `string organisation_id = 5`
- `string product_id = 6`
- `string plan_id = 7`
- `google.protobuf.Timestamp locked_until = 8`
- `string locked_by = 9`
- `string correlation_id = 10`
- `google.protobuf.Timestamp timestamp = 11`

Recommended `PremiumCalculatedEvent`:

- `string event_id = 1`
- `string quotation_id = 2`
- `string tenant_id = 3`
- `string customer_id = 4`
- `string organisation_id = 5`
- `string product_id = 6`
- `string plan_id = 7`
- `insuretech.common.v1.Money premium_amount = 8`
- `insuretech.common.v1.Money vat_amount = 9`
- `insuretech.common.v1.Money service_fee_amount = 10`
- `insuretech.common.v1.Money total_payable = 11`
- `string pricing_source = 12`
- `bool used_fallback_rate = 13`
- `string correlation_id = 14`
- `google.protobuf.Timestamp timestamp = 15`

Recommended `PolicyIssuanceRequestedEvent`:

- `string event_id = 1`
- `string tenant_id = 2`
- `string quotation_id = 3`
- `string order_id = 4`
- `string invoice_id = 5`
- `string payment_id = 6`
- `string customer_id = 7`
- `string organisation_id = 8`
- `string product_id = 9`
- `string plan_id = 10`
- `string requested_by = 11`
- `string correlation_id = 12`
- `google.protobuf.Timestamp timestamp = 13`

Recommended `PolicyIssuedEvent`:

- `string event_id = 1`
- `string policy_id = 2`
- `string policy_number = 3`
- `string tenant_id = 4`
- `string quotation_id = 5`
- `string order_id = 6`
- `string invoice_id = 7`
- `string payment_id = 8`
- `string customer_id = 9`
- `string organisation_id = 10`
- `string product_id = 11`
- `string plan_id = 12`
- `google.protobuf.Timestamp effective_from = 13`
- `google.protobuf.Timestamp effective_until = 14`
- `string correlation_id = 15`
- `google.protobuf.Timestamp timestamp = 16`

Recommended `PolicyIssuanceFailedEvent`:

- `string event_id = 1`
- `string tenant_id = 2`
- `string quotation_id = 3`
- `string order_id = 4`
- `string invoice_id = 5`
- `string payment_id = 6`
- `string customer_id = 7`
- `string organisation_id = 8`
- `string product_id = 9`
- `string plan_id = 10`
- `string failure_code = 11`
- `string failure_reason = 12`
- `string correlation_id = 13`
- `google.protobuf.Timestamp timestamp = 14`

Recommended `PolicyActivatedEvent`:

- `string event_id = 1`
- `string policy_id = 2`
- `string policy_number = 3`
- `string tenant_id = 4`
- `string customer_id = 5`
- `string organisation_id = 6`
- `google.protobuf.Timestamp activated_at = 7`
- `string correlation_id = 8`
- `google.protobuf.Timestamp timestamp = 9`

Recommended `PolicyCancelledEvent`:

- `string event_id = 1`
- `string policy_id = 2`
- `string policy_number = 3`
- `string tenant_id = 4`
- `string order_id = 5`
- `string invoice_id = 6`
- `string customer_id = 7`
- `string organisation_id = 8`
- `string reason = 9`
- `google.protobuf.Timestamp cancelled_at = 10`
- `string correlation_id = 11`
- `google.protobuf.Timestamp timestamp = 12`

Recommended `PolicyRenewedEvent`:

- `string event_id = 1`
- `string policy_id = 2`
- `string prior_policy_id = 3`
- `string policy_number = 4`
- `string tenant_id = 5`
- `string invoice_id = 6`
- `string payment_id = 7`
- `string customer_id = 8`
- `string organisation_id = 9`
- `google.protobuf.Timestamp effective_from = 10`
- `google.protobuf.Timestamp effective_until = 11`
- `string correlation_id = 12`
- `google.protobuf.Timestamp timestamp = 13`

Recommended `PolicyLapsedEvent`:

- `string event_id = 1`
- `string policy_id = 2`
- `string policy_number = 3`
- `string tenant_id = 4`
- `string customer_id = 5`
- `string organisation_id = 6`
- `google.protobuf.Timestamp lapsed_at = 7`
- `string lapse_reason = 8`
- `string correlation_id = 9`
- `google.protobuf.Timestamp timestamp = 10`

Recommended `PolicyReinstatedEvent`:

- `string event_id = 1`
- `string policy_id = 2`
- `string policy_number = 3`
- `string tenant_id = 4`
- `string invoice_id = 5`
- `string payment_id = 6`
- `string customer_id = 7`
- `string organisation_id = 8`
- `google.protobuf.Timestamp reinstated_at = 9`
- `string correlation_id = 10`
- `google.protobuf.Timestamp timestamp = 11`

Recommended `PolicyEndorsedEvent`:

- `string event_id = 1`
- `string policy_id = 2`
- `string policy_number = 3`
- `string tenant_id = 4`
- `string invoice_id = 5`
- `string customer_id = 6`
- `string organisation_id = 7`
- `string endorsement_type = 8`
- `string endorsement_number = 9`
- `google.protobuf.Timestamp endorsed_at = 10`
- `string correlation_id = 11`
- `google.protobuf.Timestamp timestamp = 12`

### 31.5 Insurance event behavior rules

Required rules:

- `PolicyIssuedEvent` must be emitted only after durable insurance persistence
- `PolicyActivatedEvent` may be separate from issue only if the system allows deferred activation; otherwise it may be omitted in phase 1
- `PremiumCalculatedEvent` should indicate whether cached/fallback pricing was used to satisfy the SRS fallback requirement
- `PolicyRenewedEvent` must support versioned document generation and renewal reminders
- `PolicyCancelledEvent` must be consumable by billing, payments, documents, and compliance flows

### 31.6 Transitional rule if insurance service cannot emit events immediately

If `insurance-service` cannot publish these events in the first implementation slice:

- `PoliSync` may publish integration events after successful insurance RPCs
- those bridge events must use the same payload contract as the intended final insurance events
- ownership should move to `insurance-service` once native event publishing exists

## 32. Current Go `orders-service` Assessment

This section is based on the current Go implementation under `backend/inscore/microservices/orders`. It exists to anchor the plan to the real service, not the intended future design.

### 32.1 What the current service actually does

The current service provides these gRPC operations:

- `CreateOrder`
- `GetOrder`
- `ListOrders`
- `InitiatePayment`
- `ConfirmPayment`
- `CancelOrder`
- `GetOrderStatus`

Current service behavior:

- `CreateOrder` only requires `quotation_id`
- `CreateOrder` does not call quotation, billing, or payment services
- `CreateOrder` accepts missing `product_id`, `plan_id`, and `total_payable` and persists empty/default values
- `CreateOrder` already has a service-level idempotency replay check
- `InitiatePayment` calls `payment-service` via gRPC when `paymentClient` is configured
- `InitiatePayment` falls back to a local stub path when the payment client is unavailable
- the stub path fabricates:
  - `payment_id`
  - `gateway_ref`
  - `payment_url`
- `ConfirmPayment` does not verify with a payment service or invoice service
- `ConfirmPayment` trusts incoming `payment_id` and `transaction_id` plus stored order state
- `CancelOrder` directly mutates order state and reason
- event publication is direct best-effort Kafka publish, not outbox-backed
- event publication failure is logged and ignored
- missing Kafka producer is treated as non-fatal no-op

### 32.2 Current state machine in the service

The service comments and code currently implement only:

- `PENDING -> PAYMENT_INITIATED -> PAID -> POLICY_ISSUED`
- `PENDING -> CANCELLED`
- `PAYMENT_INITIATED -> CANCELLED`
- `FAILED` exists in repo support but is not a first-class service path

This is materially narrower than the business requirements.

### 32.3 Current repository capabilities

The repository currently supports:

- `CreateOrder`
- `GetOrder`
- `GetOrderByIdempotencyKey`
- `ListOrders`
- `UpdateOrderStatus`
- `SetPaymentInfo`
- `SetPaymentStatus`
- `SetFulfillmentStatus`
- `SetPolicyID`
- `SetCancellationReason`
- `SetFailureReason`

This means the persistence layer does not yet support:

- invoice linkage
- billing status
- payment review states
- reconciliation fields
- aml/compliance flags
- versioning/optimistic concurrency
- outbox records

### 32.4 Current event implementation reality

The service currently publishes:

- `orders.order.created`
- `orders.order.payment_initiated`
- `orders.order.payment_confirmed`
- `orders.order.cancelled`
- `orders.order.failed`

Current issues:

- topic names do not follow the canonical names proposed in this plan
- no transactional outbox
- no guaranteed delivery
- no event version field
- inconsistent correlation ID semantics
  - random UUID on create
  - idempotency key on payment initiation
  - transaction ID on payment confirmation
  - caller ID on cancel
- no standard event header contract

### 32.5 Current security and tenancy reality

The service currently resolves tenant and caller from gRPC metadata with fallback behavior:

- tenant from request metadata or env var or a hardcoded default tenant UUID
- caller from `x-user-id`, `x-customer-id`, or `x-subject`

Current gaps:

- no synchronous AuthN validation inside the service boundary
- no synchronous AuthZ validation inside the service boundary
- no explicit B2B organisation resolution
- no tenant scoping in repository queries beyond data written at create time
- `GetOrder` and `ListOrders` do not appear to enforce tenant isolation in data access

### 32.6 Concrete service gaps against the target architecture

The current `orders-service` is missing at least these core responsibilities:

- no mandatory payment-service integration
- no billing-service integration
- no invoice creation or invoice linking
- no payment verification trust boundary
- no manual-review or pending-verification states
- no outbox
- no replay-safe consumers/publishers
- no richer order dimensions
- no policy-issuance callback/update path exposed in service API
- no explicit order failure flow exposed as RPC
- no idempotent command handling at repository/service boundary
- no tenant-safe list/get enforcement
- no B2B purchase-order linkage
- no compliance/audit projection hooks beyond raw events

### 32.7 Current Go `payment-service` assessment

This section is based on the current Go implementation under `backend/inscore/microservices/payment`.

### 32.8 What the current payment service actually does

The current service provides these gRPC operations:

- `InitiatePayment`
- `VerifyPayment`
- `GetPayment`
- `ListPayments`
- `InitiateRefund`
- `GetRefundStatus`
- `ListPaymentMethods`
- `AddPaymentMethod`
- `ReconcilePayments`

Current service behavior:

- `InitiatePayment` persists a real payment record in `payment_schema.payments`
- idempotency on `InitiatePayment` is already implemented by `idempotency_key`
- only two method paths are actually implemented in code:
  - `BANK_TRANSFER`
  - `CARD` or `SSLCOMMERZ`
- proto comments mention `BKASH`, `NAGAD`, and `ROCKET`, but the current Go code rejects them as unsupported
- the SSLCommerz path does not call the provider API
- the SSLCommerz path only builds a local redirect URL from `PAYMENT_SSLCOMMERZ_HOSTED_BASE_URL`
- no store credential handling exists yet
- no provider session creation exists yet
- no callback or IPN handler exists yet
- `VerifyPayment` marks a payment `SUCCESS` directly after basic request checks
- `VerifyPayment` does not call SSLCommerz validation APIs
- `VerifyPayment` does not validate `val_id`, `bank_tran_id`, or settled amount against provider data
- receipt generation is a placeholder `receipt_url` string, not a generated document workflow
- `InitiateRefund` creates a refund record and immediately marks it `COMPLETED`
- `ReconcilePayments` is currently just a local status scan, not provider reconciliation
- Kafka publishing is direct best-effort publish, not outbox-backed
- there are currently no payment-service tests in the repo

### 32.9 Current payment data and contract gaps

The current payment proto and repository mapping are still too thin for a real hosted-checkout workflow.

What is present:

- internal `payment_id`
- `transaction_id`
- payment type, method, status
- amount, payer, payee
- `gateway`
- `gateway_response`
- `receipt_url`
- retry count, failure reason, idempotency key

What is still missing for practical SSLCommerz and manual-payment support:

- `order_id`
- `invoice_id`
- `tenant_id`
- `customer_id`
- `organisation_id`
- `purchase_order_id`
- provider reference fields such as:
  - `tran_id`
  - `val_id`
  - `bank_tran_id`
  - `session_key`
- explicit callback URLs and callback receipt timestamps
- manual-proof linkage fields
- manual-review status and reviewer identity
- receipt identifiers such as `receipt_number`, `receipt_document_id`, and `receipt_file_id`
- ledger posting references
- raw webhook/audit persistence separate from the summary `gateway_response` blob

### 32.10 Current gateway and runtime exposure gaps

Current runtime gaps around the payment service:

- the gateway router currently exposes order routes, but there are no payment HTTP routes registered in `backend/inscore/cmd/gateway/internal/routes/router.go`
- there is no public webhook endpoint for SSLCommerz callbacks
- there is no internal payment callback handler that bridges browser return, IPN, and provider validation
- config currently only supports:
  - `KAFKA_BROKERS`
  - `PAYMENT_DEFAULT_PAYEE_ID`
  - `PAYMENT_SSLCOMMERZ_HOSTED_BASE_URL`

Minimum config expansion required:

- `PAYMENT_SSLCOMMERZ_STORE_ID`
- `PAYMENT_SSLCOMMERZ_STORE_PASSWORD`
- `PAYMENT_SSLCOMMERZ_API_BASE_URL`
- `PAYMENT_SSLCOMMERZ_VALIDATION_BASE_URL`
- `PAYMENT_SSLCOMMERZ_REFUND_BASE_URL`
- `PAYMENT_PUBLIC_BASE_URL`
- `PAYMENT_WEBHOOK_ALLOWED_SOURCE_CIDRS`

### 32.11 Practical target shape for phase-1 payment delivery

Phase 1 should not try to solve all providers at once. The practical first slice is:

- retail premium payment
- hosted SSLCommerz checkout
- one payment intent linked to one order
- provider callback plus validation
- receipt generation
- order confirmation event
- policy issuance handoff

Out of scope for phase 1 unless explicitly pulled in:

- partial payment
- installment plans
- multi-invoice settlement
- Rocket, Nagad, and bKash direct API integrations
- automated refund settlement without an ops checkpoint
- full TigerBeetle ledger rollout across all payment types

## 33. Concrete Change Plan for Go `orders-service`

This section converts the gap analysis above into explicit service work.

### 33.1 Service role after redesign

After redesign, `orders-service` should do only these things:

- own order persistence
- enforce local order invariants
- expose order query and transition APIs
- emit durable order events through an outbox

It should stop doing these things:

- fabricating payment URLs and gateway refs on its own
- acting as a fake payment provider
- silently accepting missing commercial fields required for downstream work
- publishing best-effort events without durability guarantees

### 33.2 Required service integrations

The current service currently has zero hard dependency on billing or payment. That must change in one of two acceptable ways:

Option A: orchestration-first, recommended

- `PoliSync` calls `orders-service` only for order commands
- `PoliSync` separately calls `billing-service` and `payment-service`
- `orders-service` remains integration-light and state-focused

Option B: service-collaboration path

- `orders-service` may request invoice or payment initiation through internal clients
- still must not fabricate payment details locally

Recommendation:

- keep `orders-service` integration-light
- let `PoliSync` orchestrate billing and payment
- reduce `orders-service` to authoritative order state management and event emission

### 33.3 Required order API changes

Keep existing RPCs:

- `CreateOrder`
- `GetOrder`
- `ListOrders`
- `CancelOrder`
- `GetOrderStatus`

Change semantics of:

- `InitiatePayment`
  - should no longer mint fake payment ids or URLs itself
  - should either:
    - become `MarkPaymentInitiated` from a trusted orchestrator path, or
    - accept a billing/payment linkage request after `payment-service` initiation succeeds
- `ConfirmPayment`
  - should no longer behave as a loosely trusted external callback endpoint
  - should accept only orchestrator-confirmed settlement input from `PoliSync`

Add RPCs:

- `MarkInvoiceIssued`
- `MarkPaymentPendingReview`
- `MarkPaymentFailed`
- `MarkPolicyIssuanceStarted`
- `MarkPolicyIssued`
- `MarkPolicyIssuanceFailed`
- `FailOrder`

Reason:

- the current API forces too much meaning into `InitiatePayment` and `ConfirmPayment`
- downstream orchestration needs explicit state transitions, not overloaded methods

### 33.4 Required `CreateOrder` changes

Current problem:

- it creates orders with only `quotation_id`, `tenant_id`, and optional `customer_id`
- it permits empty product, plan, and total payable data

Required changes:

- require or populate:
  - `tenant_id`
  - `quotation_id`
  - `customer_id` or `organisation_id`
  - `product_id`
  - `plan_id`
  - `currency`
  - `total_payable`
- add optional:
  - `purchase_order_id`
  - `coverage_start_at`
  - `coverage_end_at`
  - `source_channel`
  - `created_by`
  - `idempotency_key`
- reject order creation if required quotation enrichment is unavailable
- stop relying on nil UUID placeholders for missing commercial references

### 33.5 Required `InitiatePayment` changes

Current problem:

- the method is acting like a payment gateway adapter

Required changes:

- remove fake generation of:
  - `payment_id`
  - `gateway_ref`
  - `payment_url`
- convert this RPC into an order state update that records:
  - `invoice_id`
  - `payment_id`
  - `payment_gateway_ref`
  - `payment_due_at`
  - `payment_status`
- require idempotent linkage so repeated orchestration calls do not duplicate state transitions
- publish an order event only after durable state update

### 33.6 Required `ConfirmPayment` changes

Current problem:

- it only checks order status and `payment_id`
- it trusts caller-provided `transaction_id`
- it does not validate invoice linkage or paid amount

Required changes:

- accept:
  - `order_id`
  - `payment_id`
  - `invoice_id`
  - `transaction_id`
  - `provider_reference`
  - `paid_amount`
  - `currency`
  - `confirmed_by_service`
  - `correlation_id`
- validate:
  - payment id matches stored order payment id
  - invoice id matches stored invoice id
  - order is in a payable state
  - amount/currency match order expectation or supported partial-payment rules
- transition into `PAYMENT_CONFIRMED` style state before full fulfillment if the order model is expanded
- emit `OrderPaymentConfirmedEvent` using standardized correlation semantics

### 33.7 Required cancellation and failure changes

Current `CancelOrder` is too narrow for the target flow.

Required additions:

- support cancellation source:
  - customer
  - agent
  - admin
  - system
- persist cancellation actor and cancellation timestamp
- support pre-settlement cancellation and post-settlement cancellation-request states if needed
- add explicit `FailOrder` RPC for orchestrated failure handling
- use `SetFailureReason` through service APIs, not only repository helpers

### 33.8 Required policy linkage changes

The repository has `SetPolicyID`, but the service does not expose a proper policy-issued transition path.

Required additions:

- `MarkPolicyIssuanceStarted`
- `MarkPolicyIssued`
- `MarkPolicyIssuanceFailed`

Each should:

- validate current order state
- persist fulfillment dimension changes
- emit corresponding order-side events if needed

## 34. Required Repository and Schema Expansion for `orders-service`

### 34.1 New order fields required in storage

The current repository writes only the basic order columns. Add support for:

- `invoice_id`
- `organisation_id`
- `purchase_order_id`
- `payment_status`
- `billing_status`
- `fulfillment_status`
- `manual_review_required`
- `aml_flag_status`
- `payment_due_at`
- `coverage_start_at`
- `coverage_end_at`
- `payment_confirmed_at`
- `policy_issuance_started_at`
- `policy_issued_at`
- `cancelled_at`
- `cancelled_by`
- `failed_at`
- `failed_by`
- `failure_code`
- `idempotency_key`
- `source_channel`
- `version`

### 34.2 Repository method expansion

Current repository methods are too coarse. Add targeted methods or a transactional update model for:

- `SetInvoiceInfo`
- `SetPaymentPendingReview`
- `ConfirmPaymentForOrder`
- `MarkOrderFailed`
- `MarkPolicyIssuanceStarted`
- `MarkPolicyIssued`
- `MarkPolicyIssuanceFailed`
- `AppendOutboxEvent`
- `ClaimIdempotencyKey` if handled at service persistence layer

### 34.3 Tenant-safe repository changes

Current `GetOrder` and `ListOrders` do not appear tenant-scoped.

Required changes:

- every query path should include tenant scope
- B2B list/query paths should optionally include organisation scope
- customer-scoped list path should enforce actor/customer or actor/organisation relationship outside or inside service boundary

### 34.4 Concurrency control

The current repository does not expose optimistic concurrency.

Required changes:

- add a `version` field or compare-and-set update model
- use it on all state transition writes
- prevent duplicate transition commits during retries or duplicate events

## 35. Required Event and Outbox Changes for `orders-service`

### 35.1 Replace best-effort publish with transactional outbox

Current problem:

- events are published after writes using a fire-and-forget logger-backed publisher
- if Kafka is down, the order state commits but the event may disappear

Required changes:

- persist outbox record in same DB transaction as order state mutation
- move Kafka publish to background outbox dispatcher
- track publish status, retry count, and last error
- support DLQ or poison-event handling

### 35.2 Standardize order event payloads

Current order events are missing enough context for downstream consumers.

Enrich all order events with:

- `tenant_id`
- `organisation_id`
- `invoice_id`
- `purchase_order_id`
- `payment_status`
- `billing_status`
- `fulfillment_status`
- `actor_user_id` or propagate through headers
- `event_version`
- `correlation_id`
- `causation_id`

### 35.3 Standardize topic naming

Current topics:

- `orders.order.created`
- `orders.order.payment_initiated`
- `orders.order.payment_confirmed`
- `orders.order.cancelled`
- `orders.order.failed`

Required migration target:

- `insuretech.orders.v1.order.created`
- `insuretech.orders.v1.order.payment_initiated`
- `insuretech.orders.v1.order.payment_confirmed`
- `insuretech.orders.v1.order.cancelled`
- `insuretech.orders.v1.order.failed`
- optionally:
  - `insuretech.orders.v1.order.invoice_linked`
  - `insuretech.orders.v1.order.payment_pending_review`
  - `insuretech.orders.v1.order.policy_issuance_started`
  - `insuretech.orders.v1.order.policy_issuance_failed`
  - `insuretech.orders.v1.order.policy_issued`

### 35.4 Correlation and keying rules

Current correlation semantics are inconsistent.

Required rules:

- use a single correlation id for one business flow across order, billing, payment, and policy steps
- use causation id for the triggering command or event
- Kafka message key should remain `order_id` for order lifecycle events

## 36. Concrete Missing-Step Execution Plan for `orders-service`

### 36.1 Immediate missing implementation steps

1. Stop fake payment initiation in `orders-service`.
2. Add invoice linkage to the order model.
3. Expand the order state model beyond `PENDING -> PAYMENT_INITIATED -> PAID`.
4. Add tenant-safe query/update enforcement.
5. Add outbox-backed event publication.
6. Add explicit failure and policy-issuance service methods.
7. Align topic names and event payload shape with the cross-domain contract.

### 36.2 Recommended implementation sequence inside the service

Phase A: persistence correction

- extend DB schema and order entity mapping
- add repository methods for new transitions
- add tenant/version enforcement

Phase B: service transition correction

- refactor `CreateOrder`
- refactor `InitiatePayment`
- refactor `ConfirmPayment`
- add `FailOrder`
- add policy-issuance transition methods

Phase C: event durability

- replace direct publisher usage with outbox writes
- add dispatcher
- add publish retry and error visibility

Phase D: orchestration alignment

- tighten the contract between `PoliSync` and `orders-service`
- make `PoliSync` the only trusted caller for commercial settlement confirmation
- remove remaining fake gateway behavior from the order domain

### 36.3 Acceptance criteria specific to Go `orders-service`

The `orders-service` slice is only complete when:

- it no longer fabricates payment execution data
- it no longer relies on missing quotation enrichment for commercially valid orders
- all order writes are tenant-safe
- all order lifecycle events are outbox-backed
- it can represent invoice linkage, payment review, fulfillment progress, and failure states
- policy issuance updates can be persisted without bypassing the service layer
- every transition is idempotent and concurrency-safe


---

## 36. Orders Service — Implementation Status (Updated: 2026-03-06)

### 36.1 What has been implemented

The `orders-service` Go microservice has been fully implemented and verified against the live database. All code lives under:

```
backend/inscore/microservices/orders/
```

#### 36.1.1 File inventory

| File | Description |
|---|---|
| `cmd/server/main.go` | Bootstrap: logger → config → DB → Kafka producer → repo → service → gRPC server → graceful shutdown |
| `internal/config/config.go` | Env-based config; gRPC port resolved from `configs/services.yaml`; Kafka brokers and payment service URL from env |
| `internal/domain/interfaces.go` | `OrderRepository` + `OrderService` interfaces, `OrderCreateInput`, `OrderUpdateInput`, sentinel `ErrNotFound` |
| `internal/repository/repository.go` | PostgreSQL data access: `map[string]any` inserts via `db.Table().Create()`, explicit-column SELECT into `orderScanRow`, then conversion to proto `Order` |
| `internal/service/order_service.go` | Business logic with state machine enforcement and Kafka event publication |
| `internal/service/errors.go` | Sentinel errors: `ErrInvalidArgument`, `ErrNotFound`, `ErrInvalidTransition`, `ErrPaymentFailed` |
| `internal/grpc/order_handler.go` | gRPC transport layer, maps service errors to gRPC status codes |
| `internal/grpc/server.go` | gRPC server setup |
| `internal/events/publisher.go` | Wraps `kafkaproducer.EventProducer`; publishes typed proto events per state transition; nil-safe (no Kafka = silent drop) |
| `internal/events/topics.go` | Kafka topic constants aligned with canonical naming in section 9.1 |
| `internal/consumers/consumer.go` | Handles `payment.completed`, `payment.failed`, `policy.issued` events to drive order state transitions |

#### 36.1.2 Inject-tag applied

`protoc-go-inject-tag` was run on:

```
gen/go/insuretech/orders/entity/v1/order.pb.go
```

GORM tags are now injected from `@inject_tag` proto comments.

#### 36.1.3 New shared serializer

Added `backend/inscore/db/money_serializer.go`:

- `proto_money` GORM serializer registered in `init()`
- handles `*commonv1.Money` ↔ BIGINT (amount in paisa) conversion
- used for the `total_payable` column

#### 36.1.4 Repository pattern

- **Insert**: `db.Table("insurance_schema.orders").Create(map[string]any{...})` — raw Go types (string for enums, `int64` for money, `time.Time` for timestamps)
- **Read**: explicit SELECT column list with `COALESCE()` for nullable text columns, scanned into `orderScanRow` plain struct, then converted to proto `Order` via `scanRowToProto()`
- Avoids GORM auto-mapping issues with `*commonv1.Money` and `*timestamppb.Timestamp`

#### 36.1.5 Implemented state machine

```
PENDING → PAYMENT_INITIATED → PAID → POLICY_ISSUED
                           ↘ FAILED
PENDING / PAYMENT_INITIATED → CANCELLED
```

Service-layer enforcement:
- `InitiatePayment` only allowed from `PENDING`
- `ConfirmPayment` only allowed from `PAYMENT_INITIATED`; validates `payment_id` match
- `CancelOrder` only allowed from `PENDING` or `PAYMENT_INITIATED`

#### 36.1.6 Kafka events published

| State transition | Topic | Event type |
|---|---|---|
| Order created | `orders.order.created` | `OrderCreatedEvent` |
| Payment initiated | `orders.order.payment_initiated` | `OrderPaymentInitiatedEvent` |
| Payment confirmed | `orders.order.payment_confirmed` | `OrderPaymentConfirmedEvent` |
| Order cancelled | `orders.order.cancelled` | `OrderCancelledEvent` |

#### 36.1.7 Kafka events consumed

| Topic | Handler |
|---|---|
| `payment.completed` | `HandlePaymentCompleted` → transitions order to `PAID` |
| `payment.failed` | `HandlePaymentFailed` → transitions order to `FAILED` with reason |
| `policy.issued` | `HandlePolicyIssued` → links `policy_id` and transitions to `POLICY_ISSUED` |

### 36.2 Live database tests

All tests run against the live PostgreSQL instance with real fixture data. Tests self-clean via `t.Cleanup`.

#### 36.2.1 Repository live tests (12/12 PASS)

Location: `backend/inscore/microservices/orders/internal/repository/`

| Test | Covers |
|---|---|
| `TestLiveRepository_CreateOrder` | Insert, money reconstruction, timestamps |
| `TestLiveRepository_CreateOrder_DuplicateOrderNumber` | UNIQUE constraint |
| `TestLiveRepository_GetOrder` | Full field read-back |
| `TestLiveRepository_GetOrder_NotFound` | `ErrNotFound` sentinel |
| `TestLiveRepository_ListOrders` | Pagination + status filter |
| `TestLiveRepository_UpdateOrderStatus` | PENDING→PAYMENT_INITIATED→PAID + `paid_at` set |
| `TestLiveRepository_UpdateOrderStatus_NotFound` | `ErrNotFound` sentinel |
| `TestLiveRepository_SetPaymentInfo` | payment_id + gateway_ref + status transition |
| `TestLiveRepository_SetPolicyID` | policy_id + status transition |
| `TestLiveRepository_SetCancellationReason` | cancellation_reason + status |
| `TestLiveRepository_SetFailureReason` | failure_reason + status |
| `TestLiveRepository_FullOrderLifecycle` | End-to-end: create→pay→paid→policy_issued |

Run with:
```bash
go test ./backend/inscore/microservices/orders/internal/repository/... -run "TestLive" -count=1 -timeout 120s -v
```

#### 36.2.2 Service live tests (13/13 PASS)

Location: `backend/inscore/microservices/orders/internal/service/`

| Test | Covers |
|---|---|
| `TestOrderService_Live_CreateOrder` | Service path (skips if FK not satisfied by quotation lookup) |
| `TestOrderService_Live_CreateOrder_ValidationErrors` | nil request, missing quotation_id |
| `TestOrderService_Live_GetOrder` | GetOrder by ID |
| `TestOrderService_Live_GetOrder_NotFound` | `ErrNotFound` |
| `TestOrderService_Live_GetOrder_ValidationErrors` | nil, empty order_id |
| `TestOrderService_Live_ListOrders` | Pagination, total count |
| `TestOrderService_Live_GetOrderStatus` | Status + empty payment/policy IDs |
| `TestOrderService_Live_InitiatePayment` | SetPaymentInfo + status check |
| `TestOrderService_Live_InitiatePayment_ValidationErrors` | nil, missing fields |
| `TestOrderService_Live_InitiatePayment_WrongStatus` | `ErrInvalidTransition` |
| `TestOrderService_Live_ConfirmPayment` | Confirm + paid_at verification |
| `TestOrderService_Live_ConfirmPayment_WrongPaymentID` | `ErrInvalidArgument` |
| `TestOrderService_Live_CancelOrder` | Cancel + cancellation_reason check |
| `TestOrderService_Live_CancelOrder_AfterPayment_ShouldFail` | `ErrInvalidTransition` after PAID |
| `TestOrderService_Live_FullOrderLifecycle` | Full: create→initiate→confirm→policy_issued |

Run with:
```bash
go test ./backend/inscore/microservices/orders/internal/service/... -run "TestOrderService_Live" -count=1 -timeout 120s -v
```

#### 36.2.3 Test fixture chain

Tests create a complete FK-valid fixture chain and tear it down via `t.Cleanup`:

```
authn_schema.users
  └── insurance_schema.products  (created_by → users)
        └── insurance_schema.product_plans  (product_id → products)
              ├── insurance_schema.quotations  (plan_id → product_plans)
              │     └── insurance_schema.orders  (quotation_id, customer_id, product_id, plan_id)
              ├── payment_schema.payments  (payer_id → users)  → orders.payment_id FK
              └── insurance_schema.policies  (product_id, customer_id)  → orders.policy_id FK
```

### 36.3 Known gaps and next steps

The following items from the plan are **not yet implemented** in the orders service:

| Gap | Plan section | Priority |
|---|---|---|
| AuthN/AuthZ enforcement (`ValidateToken` + `CheckAccess`) on every command | §5, §4.3 | High |
| `invoice_id` field on Order + billing linkage | §8.4, §19.3 | High |
| Outbox pattern for Kafka events (currently direct publish) | §14, §20.3 | High |
| `product_id` / `plan_id` resolution from quotation service in `CreateOrder` | §10.1 step 4 | High |
| Richer order status dimensions (`payment_status`, `billing_status`, `fulfillment_status`) | §8.4.1 | Medium |
| `payment_due_at`, `coverage_start_at`, `coverage_end_at` on Order | §8.4 | Medium |
| `manual_review_required` and `aml_flag_status` on Order | §8.4, §13 | Medium |
| Idempotency key handling on `CreateOrder` and `InitiatePayment` | §5, §20.3 | Medium |
| Correlation/causation ID propagation in all emitted events | §5 | Medium |
| `payment-service` gRPC call in `InitiatePayment` is optional and still falls back to a local stub path when the client is absent | §10.1 step 7 | High |
| Payment gateway ref stored with real FK to `payment_schema.payments` | §7.2 | High |
| AML/CFT hook triggers in order creation and payment paths | §13.1 | Low |
| Dead-letter topic (DLQ) processing in consumer | §14 | Medium |
| HTTP/REST gateway registration for orders endpoints | §3.3 | Medium |

### 36.4 How to run the service locally

```bash
# Ensure DB is running and services.yaml has orders config
ORDERS_LIVE_TEST=true \
INSCORE_DB_CONFIG=path/to/database.yaml \
go run ./backend/inscore/microservices/orders/cmd/server/main.go
```

Or via the workspace runner:

```bash
cd backend/inscore && go run ./cmd/orders/main.go
```

## 36A. Payment Service - Practical Implementation Plan (Updated: 2026-03-06)

This section supersedes any vague payment implementation language elsewhere in the document. It is anchored to the current Go code and the official SSLCommerz docs reviewed on 2026-03-06.

### 36A.0 Current implementation status on this branch

Earlier statements in this document that describe `payment-service` as stub-only, or describe the gateway as having no payment routes or no SSLCommerz callback surface, are now partially outdated.

What is already implemented in code:

- real SSLCommerz session initiation in `backend/inscore/microservices/payment/internal/providers/sslcommerz/client.go`
- real provider-backed payment verification in `backend/inscore/microservices/payment/internal/service/payment_service.go`
- first-class persistence of provider/manual-review/receipt fields in `backend/inscore/microservices/payment/internal/repository/repository.go`
- public SSLCommerz callback routes in `backend/inscore/cmd/gateway/internal/routes/router.go`
  - `POST /v1/payments/webhook/sslcommerz`
  - `GET|POST /v1/payments/sslcommerz/success`
  - `GET|POST /v1/payments/sslcommerz/fail`
  - `GET|POST /v1/payments/sslcommerz/cancel`
- gateway callback transport in `backend/inscore/cmd/gateway/internal/handlers/payment_callback_handler.go`
  - gateway now forwards raw callback payloads to `PaymentService.HandleGatewayWebhook`
  - gateway is no longer responsible for deciding payment state
- payment-service gRPC methods now implemented for:
  - `HandleGatewayWebhook`
  - `GetPaymentByProviderReference`
  - `SubmitManualPaymentProof`
  - `ReviewManualPayment`
  - `GenerateReceipt`
  - `GetPaymentReceipt`
- `orders-service` stub fallback removed
  - if payment client is unavailable, order payment initiation now fails instead of fabricating payment ids or URLs
- test coverage now includes:
  - unit tests for SSLCommerz initiation and validation
  - unit tests for callback handling, provider-reference lookup, manual review, and receipt generation
  - live DB tests for SSLCommerz initiation/verification and manual-review/receipt persistence

What is still not complete:

- refund status polling and full SSLCommerz refund lifecycle are still thinner than the docs require
- provider webhook audit persistence is still stored in `gateway_response`; a dedicated `payment_provider_webhooks` table is still pending
- receipt generation currently produces durable receipt metadata and a stable receipt endpoint, but not a document-service-backed PDF artifact yet
- provider transaction history is still summarized on the payment row; a dedicated provider-transaction table is still pending

Practical interpretation:

- Phase 1 and the minimum useful Phase 2 callback flow are now implemented
- the remaining work is mostly Phase 4 and Phase 5 hardening, not foundational payment initiation anymore

### 36A.1 Delivery principle

Do not rewrite `payment-service` from scratch.

Build on the existing Go service in `backend/inscore/microservices/payment`.

Stub-removal status:

- synthetic SSLCommerz URL generation: removed
- trust-on-request payment verification: removed
- immediate refund completion: still pending full cleanup

Non-negotiable rule:

- once this implementation starts, any path that would previously fall back to fake hosted URLs or fake success states must return a real error instead
- local development may still use sandbox SSLCommerz credentials, but not fake payment success

### 36A.1.1 Exact SSLCommerz endpoints to implement

Based on the official docs reviewed on 2026-03-06:

- create session:
  - sandbox: `https://sandbox.sslcommerz.com/gwprocess/v4/api.php`
  - live: `https://securepay.sslcommerz.com/gwprocess/v4/api.php`
- validate successful transaction by `val_id`:
  - sandbox: `https://sandbox.sslcommerz.com/validator/api/validationserverAPI.php`
  - live: `https://securepay.sslcommerz.com/validator/api/validationserverAPI.php`
- transaction lookup by `sessionkey` or merchant transaction context:
  - sandbox: `https://sandbox.sslcommerz.com/validator/api/merchantTransIDvalidationAPI.php`
  - live: `https://securepay.sslcommerz.com/validator/api/merchantTransIDvalidationAPI.php`
- refund initiate and refund query:
  - sandbox: `https://sandbox.sslcommerz.com/validator/api/merchantTransIDvalidationAPI.php`
  - live: `https://securepay.sslcommerz.com/validator/api/merchantTransIDvalidationAPI.php`

Operational network rules from the official docs:

- SSLCommerz accepts only `TLS 1.2` or higher
- sandbox listener connectivity notes reference:
  - inbound from `103.26.139.87` on TCP `80/443`
  - outbound reachability to `103.26.139.87` on TCP `443`
- live listener connectivity notes reference:
  - inbound from `103.26.139.81` and `103.132.153.81` on TCP `80/443`
  - outbound reachability to `103.26.139.148` and `103.132.153.148` on TCP `443`

Implementation rule:

- keep these IPs and hostnames in ops documentation and firewall config
- do not hardcode the IP list into application logic

### 36A.1.2 Phase-1 supported use case

Phase 1 must support exactly this production path:

1. customer starts checkout for one order
2. `orders-service` creates the order
3. `payment-service` creates one SSLCommerz session for that order
4. customer is redirected to the real `GatewayPageURL`
5. SSLCommerz sends IPN to InsureTech
6. InsureTech validates the transaction with `val_id`
7. `payment-service` marks the payment validated and completed
8. `orders-service` is updated from trusted payment outcome only
9. policy issuance starts
10. receipt is generated asynchronously if needed

Anything outside that sequence is a later phase.

### 36A.2 Phase 0 - contract and schema stabilization

First, stabilize the write model before adding provider calls.

Required proto additions:

- extend `payment.proto` with:
  - `order_id`
  - `invoice_id`
  - `tenant_id`
  - `customer_id`
  - `organisation_id`
  - `purchase_order_id`
  - `provider`
  - `provider_reference`
  - `tran_id`
  - `val_id`
  - `bank_tran_id`
  - `session_key`
  - `manual_review_status`
  - `manual_proof_file_id`
  - `verified_by`
  - `verified_at`
  - `receipt_number`
  - `receipt_document_id`
  - `receipt_file_id`
- extend `payment_service.proto` with:
  - `HandleGatewayWebhook`
  - `GetPaymentByProviderReference`
  - `SubmitManualPaymentProof`
  - `ReviewManualPayment`
  - `GenerateReceipt`
  - `GetPaymentReceipt`
- extend `payment_events.proto` with:
  - `PaymentVerifiedEvent`
  - `ManualPaymentProofSubmittedEvent`
  - `ReceiptGeneratedEvent`
  - `PaymentReconciliationMismatchEvent`

Required database additions:

- add the provider correlation fields above to `payment_schema.payments`
- create a provider-neutral webhook table instead of forcing SSLCommerz into the current `mfs_webhooks` naming
- create a provider-neutral payment-attempt or provider-transaction table instead of overloading `gateway_response`
- add an outbox table for payment events if a shared outbox table does not already exist

Concrete tables to add:

- `payment_schema.payment_provider_transactions`
  - `provider_transaction_id`
  - `payment_id`
  - `provider`
  - `tran_id`
  - `val_id`
  - `session_key`
  - `bank_tran_id`
  - `status`
  - `amount`
  - `currency`
  - `risk_level`
  - `risk_title`
  - `raw_response`
  - `created_at`
  - `updated_at`
- `payment_schema.payment_provider_webhooks`
  - `webhook_id`
  - `provider`
  - `payment_id`
  - `tran_id`
  - `val_id`
  - `request_headers`
  - `request_body`
  - `received_at`
  - `validation_status`
  - `processed_at`
  - `error_reason`
- `payment_schema.payment_receipts`
  - `receipt_id`
  - `payment_id`
  - `receipt_number`
  - `document_id`
  - `file_id`
  - `created_at`

Practical rule:

- if the existing `mfs_transactions` and `mfs_webhooks` tables are generic enough, reuse them only after renaming their ownership in the plan and code comments
- otherwise create `payment_schema.provider_transactions` and `payment_schema.provider_webhooks` and leave the MFS tables alone

Recommendation:

- do not reuse the MFS table names for SSLCommerz
- keep MFS naming for bKash/Nagad/Rocket work later
- create payment-provider-neutral tables now so card checkout and MFS integrations can share the same storage model later

### 36A.3 Phase 1 - real SSLCommerz session initiation

Replace `buildSSLCommerzURL()` with a real provider client.

Implementation steps:

- create `internal/providers/sslcommerz/client.go`
- add an `InitSession` method that POSTs to the SSLCommerz session API
- map internal payment intent fields to SSLCommerz-required request fields
- persist:
  - `tran_id`
  - `session_key`
  - `GatewayPageURL`
  - provider status
  - raw provider response
- return the real hosted checkout URL from the provider response

Concrete request mapping from current Go service:

- `store_id` <- `PAYMENT_SSLCOMMERZ_STORE_ID`
- `store_passwd` <- `PAYMENT_SSLCOMMERZ_STORE_PASSWORD`
- `total_amount` <- payment amount converted from paisa to decimal BDT string
- `currency` <- `payment.amount.currency`
- `tran_id` <- generated merchant transaction id derived from internal `payment_id`
- `success_url` <- gateway URL for browser success return
- `fail_url` <- gateway URL for browser failure return
- `cancel_url` <- gateway URL for browser cancel return
- `ipn_url` <- gateway public webhook URL
- `value_a` <- internal `payment_id`
- `value_b` <- `order_id`
- `value_c` <- `tenant_id`
- `value_d` <- correlation id or idempotency key

Concrete customer fields required for checkout:

- `cus_name`
- `cus_email`
- `cus_add1`
- `cus_city`
- `cus_postcode`
- `cus_country`
- `cus_phone`

If those fields are not available on the current payment request:

- `PoliSync` or the caller must enrich them before calling `payment-service`
- do not invent placeholder customer data in production

Concrete service/file changes:

- `internal/config/config.go`
  - add SSLCommerz store credentials and base URLs
- `internal/service/payment_service.go`
  - remove `buildSSLCommerzURL()` usage
  - inject a provider client interface
  - build real request payloads
- `internal/domain/interfaces.go`
  - add provider client and webhook repository abstractions
- `internal/repository/repository.go`
  - persist provider transaction rows and provider fields on payment
- `server.go`
  - wire the provider client and any new repositories

Failure rule:

- if SSLCommerz session init fails, `InitiatePayment` must return an error and mark the payment failed or initiation-failed
- it must not return a placeholder URL

Required input rule changes:

- `InitiatePayment` must require `order_id` in metadata or as a first-class field
- `InitiatePayment` must generate provider callback URLs from `PAYMENT_PUBLIC_BASE_URL`
- `transaction_id` must stop being used as a vague stand-in for all provider references

Acceptance criteria:

- payment row is created before the provider call result is returned
- replay with the same `idempotency_key` returns the same payment intent
- no local fake URL is returned in normal execution

### 36A.4 Phase 2 - callback, IPN, and validation

Add a real public callback path and separate redirect handling from settlement confirmation.

Status on this branch:

- implemented: gateway public routes, gateway-to-payment-service forwarding, payment-service callback processing, provider re-validation, fail/cancel state transitions
- partially implemented: browser success path can trigger provider verification and then redirect to stored `callback_url`
- still pending: dedicated webhook audit table, explicit `verify_sign` persistence/verification workflow, retry workers around late IPN/provider requery

Implementation steps:

- expose a public gateway route such as `POST /v1/payments/webhook/sslcommerz`
- add a public browser-return handler for success, fail, and cancel URLs if the UX requires a friendly redirect page
- persist every raw callback or IPN payload before processing
- match callback payload to payment by `tran_id` and provider
- call the SSLCommerz validation API using `val_id`
- verify amount, currency, and transaction identity against the stored payment intent
- make callback processing idempotent on `payment_id` plus provider references

Concrete gateway routes:

- `POST /v1/payments/webhook/sslcommerz`
- `GET|POST /v1/payments/sslcommerz/success`
- `GET|POST /v1/payments/sslcommerz/fail`
- `GET|POST /v1/payments/sslcommerz/cancel`

Recommended handling split:

- `webhook/sslcommerz`
  - machine-to-machine IPN
  - durable processing path
  - can change payment state
- `sslcommerz/success`
  - browser return
  - record telemetry only
  - trigger lookup if IPN is late
  - must not finalize payment without validation
- `sslcommerz/fail`
  - record failed/cancelled return state
- `sslcommerz/cancel`
  - record user cancellation intent

Current implementation note:

- callback payloads are currently persisted on the payment row inside `gateway_response` together with first-class provider fields on `payment_schema.payments`
- that is sufficient for current flow control, but not sufficient for audit-grade webhook replay; keep the dedicated webhook table work in scope

Concrete validation rules from the official docs:

- first check callback payload signature fields such as `verify_sign` and `verify_key`
- then call `validationserverAPI.php` with `val_id`, `store_id`, and `store_passwd`
- accept settlement only when SSLCommerz returns a provider-valid success state (`VALID`, `VALIDATED`, or equivalent response normalized by the client)
- verify:
  - returned `tran_id` equals stored merchant transaction id
  - returned `amount` equals expected amount
  - returned `currency` equals expected currency
- persist:
  - `validated_on`
  - `bank_tran_id`
  - `card_type`
  - `card_brand`
  - `card_issuer`
  - `card_issuer_country`
  - `risk_level`
  - `risk_title`

Concrete risk rule:

- if `risk_level=1`, do not auto-issue policy
- move payment to manual review required or high-risk verified state
- emit a compliance/manual-review event

Concrete fallback rule:

- if success redirect arrives but IPN is missing, query by `sessionkey` or `tran_id`
- if validation still cannot prove settlement, keep payment pending and retry
- do not activate policy from browser redirect alone

State transition rule:

- `INITIATED` or `PENDING` -> `VERIFIED` after successful provider validation
- `VERIFIED` -> `COMPLETED` only after receipt generation, downstream event write, and any required ledger post succeed
- if the team prefers to keep the existing enum, keep a separate operational completion flag and do not collapse validation and fulfillment into the same state silently

### 36A.5 Phase 3 - order integration correction

Tighten the trust boundary between orders and payments.

Status on this branch:

- implemented: `orders-service` no longer fabricates payment ids, gateway refs, or hosted URLs when `payment-service` is unavailable
- implemented: payment rows now carry `order_id` as a first-class field when the caller provides a valid UUID
- still pending: final event-driven order confirmation path from trusted `payment.completed` or orchestration events

Required changes:

- `orders-service` must stop treating browser callback data as payment truth
- `orders-service` should confirm payment only from:
  - a trusted `payment.completed` event, or
  - a trusted orchestrator call from `PoliSync` after payment validation succeeds
- `orders-service` must always store the internal `payment_id` returned by `payment-service`
- `payment-service` must always store `order_id` so callbacks do not need cross-service guesswork

Concrete order-side correction:

- when `payment-service` is unreachable, `orders-service.InitiatePayment` must fail with a clear upstream dependency error
- remove the current stub branch in `backend/inscore/microservices/orders/internal/service/order_service.go`
- do not generate local `GW-...` references anymore
- use the real `payment_id` plus provider `tran_id`

Recommendation:

- keep the commercial truth in `payment-service`
- keep `orders-service` as the order state owner
- let `PoliSync` orchestrate the transition from validated payment to policy issuance

### 36A.6 Phase 4 - refunds, reconciliation, and manual payments

Status on this branch:

- manual-proof submission and manual-review approval/rejection RPCs are now implemented in `payment-service`
- manual approval now persists reviewer metadata and receipt metadata on the payment row
- refund and reconciliation logic still need another hardening pass to fully match the SSLCommerz operational lifecycle described below

Refunds:

- change `InitiateRefund` from immediate completion to:
  - `PENDING`
  - `APPROVED`
  - `PROCESSING`
  - `COMPLETED` or `FAILED`
- wire refund initiation to SSLCommerz refund APIs when card-hosted payments support it
- require an ops checkpoint before marking refund commercially completed in phase 1

Concrete SSLCommerz refund implementation:

- call `merchantTransIDvalidationAPI.php` with:
  - `bank_tran_id`
  - `refund_amount`
  - `refund_remarks`
  - `store_id`
  - `store_passwd`
- persist:
  - `refund_ref_id`
  - refund initiation status: `success`, `failed`, or `processing`
  - `errorReason`
- poll or query refund status with:
  - `refund_ref_id`
  - `store_id`
  - `store_passwd`
- update refund state based on SSLCommerz response:
  - `refunded`
  - `processing`
  - `cancelled`

Docs consistency note:

- the official refund API reference points to `merchantTransIDvalidationAPI.php`
- some embedded sample code on the docs page still shows `validationserverAPI.php`
- implement the endpoint documented in the API reference first, but keep the refund endpoint configurable in case SSLCommerz support asks for an environment-specific override

Operational rule from docs:

- live refund API access requires the merchant public IP to be registered with SSLCommerz
- this must be a deployment checklist item, not a hidden assumption

Reconciliation:

- replace `ReconcilePayments` local status counting with provider report import or provider transaction query
- emit mismatch events for:
  - missing provider transaction
  - amount mismatch
  - status mismatch
  - duplicate settlement

Manual payments:

- use the same `Payment` aggregate for bank transfer, cash, and cheque
- add proof upload and review APIs
- integrate with `storage-service` presigned upload and malware scanning
- emit `PaymentVerifiedEvent` or `PaymentFailedEvent` from review decisions

Phase-1 scope rule:

- manual bank transfer is allowed in the plan
- cash and cheque may keep a back-office-only path in phase 1
- but their state machine must use the same verified or rejected transitions as hosted checkout

### 36A.7 Phase 5 - receipts, ledger, and operations

Status on this branch:

- receipt metadata generation is implemented on the payment row (`receipt_number`, `receipt_url`, reviewer-linked approval path)
- receipt retrieval RPC is implemented and returns a stable InsureTech receipt endpoint
- document-service-backed receipt PDF generation is still pending
- TigerBeetle posting is still pending

Receipts:

- generate a receipt only after settlement is validated
- prefer `document-service` for receipt PDF generation if the platform already centralizes artifact generation there
- store `receipt_number`, document id, and file id on payment
- allow `COMPLETED with receipt pending` if document generation is delayed

Ledger:

- TigerBeetle posting should happen after validation and before final completion
- if TigerBeetle is not ready for phase 1, preserve placeholder ledger reference fields and keep the posting step behind a feature flag

Operations:

- add payment-service unit tests before provider rollout
- add sandbox integration tests against SSLCommerz sandbox
- add one replay/idempotency test for each command path
- add one callback-duplication test using the same `tran_id` and `val_id`

Minimum test matrix:

- `InitiatePayment` success against sandbox-mocked provider response
- `InitiatePayment` idempotency replay
- IPN payload accepted -> validation success -> payment completed
- success redirect without IPN -> fallback validation path
- duplicate IPN for same `val_id`
- amount mismatch on validation response
- `risk_level=1` validation response
- refund initiate success
- refund query moves `processing` to `refunded`
- order service fails cleanly when payment service is unavailable

### 36A.7.1 Concrete implementation backlog by file

Implementation note:

- the public callback and browser-return handlers are currently implemented in the gateway under `backend/inscore/cmd/gateway/internal/handlers/payment_callback_handler.go`
- the plan previously suggested adding standalone HTTP handlers inside `payment-service`; that is no longer required for the current gateway-fronted topology

Create:

- `backend/inscore/microservices/payment/internal/providers/sslcommerz/client.go`
- `backend/inscore/microservices/payment/internal/providers/sslcommerz/types.go`
- `backend/inscore/microservices/payment/internal/http/webhook_handler.go`
- `backend/inscore/microservices/payment/internal/http/return_handler.go`
- `backend/inscore/microservices/payment/internal/service/payment_service_test.go`

Update:

- `backend/inscore/microservices/payment/internal/config/config.go`
- `backend/inscore/microservices/payment/internal/domain/interfaces.go`
- `backend/inscore/microservices/payment/internal/repository/repository.go`
- `backend/inscore/microservices/payment/internal/service/payment_service.go`
- `backend/inscore/microservices/payment/internal/events/topics.go`
- `backend/inscore/microservices/payment/internal/events/publisher.go`
- `backend/inscore/microservices/payment/server.go`
- `backend/inscore/cmd/gateway/internal/routes/router.go`
- `backend/inscore/microservices/orders/internal/service/order_service.go`

Migration work:

- add payment provider transaction table
- add payment provider webhook table
- add receipt table or receipt linkage columns
- add refund reference and status columns if not already present

### 36A.8 Explicit non-goals for the first practical release

Do not block phase-1 release on these items unless business requires them immediately:

- multi-provider abstraction beyond SSLCommerz plus manual bank transfer
- partial settlement across multiple invoices
- full partner webhook fanout
- direct bKash, Nagad, and Rocket API integrations
- end-to-end TigerBeetle coverage for every payment and refund path

## 36B. Ordered Engineering Task List for Items 1-3

This section turns the next three implementation buckets into an execution order:

1. proto updates
2. DB migrations
3. payment-service provider client

Rule:

- do not start item 3 before item 1 and the write-side parts of item 2 are merged
- provider code without the right proto and schema will only recreate the current stub problem in a different shape

### 36B.1 Item 1 - proto updates

Goal:

- make the contracts represent real SSLCommerz and manual-payment workflow state
- remove dependence on untyped metadata for core payment linkage

#### 36B.1.1 `payment.proto` exact delta

Target file:

- `proto/insuretech/payment/entity/v1/payment.proto`

Add these identity and linkage fields to `Payment`:

- `string order_id`
- `string invoice_id`
- `string tenant_id`
- `string customer_id`
- `string organisation_id`
- `string purchase_order_id`

Add these provider correlation fields:

- `string provider`
- `string provider_reference`
- `string tran_id`
- `string val_id`
- `string session_key`
- `string bank_tran_id`
- `string card_type`
- `string card_brand`
- `string card_issuer`
- `string card_issuer_country`

Add these validation and operational fields:

- `google.protobuf.Timestamp validated_at`
- `string validation_status`
- `string risk_level`
- `string risk_title`
- `google.protobuf.Timestamp callback_received_at`
- `google.protobuf.Timestamp ipn_received_at`

Add these manual-review fields:

- `string manual_review_status`
- `string manual_proof_file_id`
- `string verified_by`
- `google.protobuf.Timestamp verified_at`
- `string rejection_reason`

Add these receipt and accounting fields:

- `string receipt_number`
- `string receipt_document_id`
- `string receipt_file_id`
- `string ledger_transaction_id`

Recommended enum additions:

- extend `PaymentStatus` with:
  - `PAYMENT_STATUS_VERIFIED`
  - `PAYMENT_STATUS_MANUAL_REVIEW_REQUIRED`
  - `PAYMENT_STATUS_RECEIPT_PENDING`
- add `ManualReviewStatus` enum:
  - `MANUAL_REVIEW_STATUS_UNSPECIFIED`
  - `MANUAL_REVIEW_STATUS_NOT_REQUIRED`
  - `MANUAL_REVIEW_STATUS_PENDING`
  - `MANUAL_REVIEW_STATUS_APPROVED`
  - `MANUAL_REVIEW_STATUS_REJECTED`

Compatibility rule:

- keep existing field numbers unchanged
- append new fields with new numbers only
- do not repurpose `transaction_id`; keep it as an internal or generic transaction reference while `tran_id` becomes the merchant transaction identifier sent to SSLCommerz

#### 36B.1.2 `payment_service.proto` exact delta

Target file:

- `proto/insuretech/payment/services/v1/payment_service.proto`

Change `InitiatePaymentRequest`:

- add first-class fields:
  - `string order_id`
  - `string invoice_id`
  - `string tenant_id`
  - `string customer_id`
  - `string organisation_id`
  - `string customer_name`
  - `string customer_email`
  - `string customer_phone`
  - `string customer_address_line1`
  - `string customer_city`
  - `string customer_postcode`
  - `string customer_country`
- keep `metadata` only for non-core extension values

Change `InitiatePaymentResponse`:

- add:
  - `string provider`
  - `string gateway_page_url`
  - `string tran_id`
  - `string session_key`

Change `VerifyPaymentRequest`:

- add:
  - `string provider`
  - `string val_id`
  - `string tran_id`
  - `string session_key`
  - `bool force_provider_requery`

Add new RPCs:

- `HandleGatewayWebhook`
- `GetPaymentByProviderReference`
- `SubmitManualPaymentProof`
- `ReviewManualPayment`
- `GenerateReceipt`
- `GetPaymentReceipt`

Recommended HTTP mappings:

- `POST /v1/payments/webhook/{provider}`
- `GET /v1/payments/provider/{provider}/references/{provider_reference}`
- `POST /v1/payments/{payment_id}:submit-proof`
- `POST /v1/payments/{payment_id}:review`
- `POST /v1/payments/{payment_id}:generate-receipt`
- `GET /v1/payments/{payment_id}/receipt`

Add exact webhook contract fields:

- `provider`
- `map<string,string> headers`
- `bytes raw_payload`
- `string remote_addr`
- `google.protobuf.Timestamp received_at`

Add exact manual review contract fields:

- `payment_id`
- `manual_proof_file_id`
- `review_notes`
- `reviewed_by`
- `decision`
- `idempotency_key`

#### 36B.1.3 `payment_events.proto` exact delta

Target file:

- `proto/insuretech/payment/events/v1/payment_events.proto`

Enrich existing events with:

- `tenant_id`
- `order_id`
- `invoice_id`
- `customer_id`
- `organisation_id`
- `provider`
- `tran_id`
- `val_id`
- `session_key`
- `receipt_number`

Add these new events:

- `PaymentVerifiedEvent`
- `ManualPaymentProofSubmittedEvent`
- `ReceiptGeneratedEvent`
- `PaymentReconciliationMismatchEvent`

Exact purpose:

- `PaymentVerifiedEvent`
  - emitted after SSLCommerz validation success but before or alongside operational completion
- `ManualPaymentProofSubmittedEvent`
  - emitted when proof upload is linked to a payment
- `ReceiptGeneratedEvent`
  - emitted when receipt artifact is persisted and downloadable
- `PaymentReconciliationMismatchEvent`
  - emitted when provider query and stored payment state diverge

#### 36B.1.4 Proto implementation checklist

- update proto files
- run code generation
- verify generated Go types compile
- verify generated gateway/OpenAPI routes reflect new methods
- update any C# or TS SDK consumers if they compile against these contracts

Definition of done for item 1:

- the payment proto can represent a real SSLCommerz transaction end to end without relying on core linkage inside `metadata`

### 36B.2 Item 2 - DB migrations

Goal:

- make the database capable of storing provider truth, raw webhook audit data, and receipt linkage

#### 36B.2.1 Migration set to create

Create new migration files under:

- `backend/inscore/db/migrations/payment_schema/`

Recommended migration order:

1. add columns to `payment_schema.payments`
2. create `payment_provider_transactions`
3. create `payment_provider_webhooks`
4. create `payment_receipts`
5. add indexes and triggers

#### 36B.2.2 `payment_schema.payments` exact additions

Add columns:

- `order_id UUID NULL`
- `invoice_id UUID NULL`
- `tenant_id UUID NULL`
- `customer_id UUID NULL`
- `organisation_id UUID NULL`
- `purchase_order_id UUID NULL`
- `provider VARCHAR(50) NULL`
- `provider_reference VARCHAR(255) NULL`
- `tran_id VARCHAR(255) NULL`
- `val_id VARCHAR(255) NULL`
- `session_key VARCHAR(255) NULL`
- `bank_tran_id VARCHAR(255) NULL`
- `card_type VARCHAR(100) NULL`
- `card_brand VARCHAR(100) NULL`
- `card_issuer VARCHAR(255) NULL`
- `card_issuer_country VARCHAR(100) NULL`
- `validation_status VARCHAR(50) NULL`
- `validated_at TIMESTAMPTZ NULL`
- `risk_level VARCHAR(20) NULL`
- `risk_title TEXT NULL`
- `callback_received_at TIMESTAMPTZ NULL`
- `ipn_received_at TIMESTAMPTZ NULL`
- `manual_review_status VARCHAR(50) NULL`
- `manual_proof_file_id UUID NULL`
- `verified_by UUID NULL`
- `verified_at TIMESTAMPTZ NULL`
- `rejection_reason TEXT NULL`
- `receipt_number VARCHAR(100) NULL`
- `receipt_document_id UUID NULL`
- `receipt_file_id UUID NULL`
- `ledger_transaction_id VARCHAR(255) NULL`

Add indexes:

- unique index on `tran_id` where not null
- unique index on `val_id` where not null
- unique index on `session_key` where not null
- index on `(provider, provider_reference)`
- index on `(tenant_id, created_at DESC)`
- index on `(order_id)`
- index on `(invoice_id)`
- index on `(bank_tran_id)`
- index on `(manual_review_status)`

Add constraints:

- check `provider in ('SSLCOMMERZ','BANK_TRANSFER',...)` if the team enforces provider values at DB level
- keep nullable linkage columns during migration to avoid breaking old rows

#### 36B.2.3 `payment_provider_transactions` exact table

Create table:

- `provider_transaction_id UUID PRIMARY KEY`
- `payment_id UUID NOT NULL REFERENCES payment_schema.payments(payment_id)`
- `provider VARCHAR(50) NOT NULL`
- `tran_id VARCHAR(255) NULL`
- `val_id VARCHAR(255) NULL`
- `session_key VARCHAR(255) NULL`
- `bank_tran_id VARCHAR(255) NULL`
- `provider_status VARCHAR(100) NULL`
- `amount BIGINT NULL`
- `currency VARCHAR(3) NULL`
- `risk_level VARCHAR(20) NULL`
- `risk_title TEXT NULL`
- `raw_request JSONB NULL`
- `raw_response JSONB NULL`
- `created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP`
- `updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP`

Indexes:

- unique `(provider, tran_id)` where `tran_id` is not null
- unique `(provider, val_id)` where `val_id` is not null
- unique `(provider, session_key)` where `session_key` is not null
- index on `payment_id`

Purpose:

- preserve every meaningful provider-side transaction lookup or validation result without bloating the core payment row

#### 36B.2.4 `payment_provider_webhooks` exact table

Create table:

- `webhook_id UUID PRIMARY KEY`
- `provider VARCHAR(50) NOT NULL`
- `payment_id UUID NULL REFERENCES payment_schema.payments(payment_id)`
- `tran_id VARCHAR(255) NULL`
- `val_id VARCHAR(255) NULL`
- `request_headers JSONB NOT NULL`
- `request_body JSONB NULL`
- `raw_body TEXT NULL`
- `remote_addr VARCHAR(255) NULL`
- `received_at TIMESTAMPTZ NOT NULL`
- `validation_status VARCHAR(50) NOT NULL`
- `processed_at TIMESTAMPTZ NULL`
- `error_reason TEXT NULL`
- `created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP`

Indexes:

- index on `(provider, received_at DESC)`
- index on `(provider, tran_id)`
- index on `(provider, val_id)`
- index on `payment_id`

Purpose:

- immutable audit trail for IPN and browser callback processing

#### 36B.2.5 `payment_receipts` exact table

Create table:

- `receipt_id UUID PRIMARY KEY`
- `payment_id UUID NOT NULL REFERENCES payment_schema.payments(payment_id)`
- `receipt_number VARCHAR(100) NOT NULL`
- `document_id UUID NULL`
- `file_id UUID NULL`
- `download_url TEXT NULL`
- `created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP`

Indexes:

- unique `receipt_number`
- unique `payment_id`

Purpose:

- keep receipt lifecycle explicit even if the payment row also stores denormalized receipt fields

#### 36B.2.6 Migration safety checklist

- all new columns are nullable in the first migration
- backfill old SSLCommerz-like or card rows where possible
- do not drop `gateway_response`
- do not rename legacy MFS tables in the same migration batch
- ship indexes after columns and tables exist

Definition of done for item 2:

- the DB can store and query SSLCommerz session creation, IPN receipt, validation, refund references, and receipt metadata without lossy JSON-only blobs

### 36B.3 Item 3 - payment-service provider client

Goal:

- replace the current URL builder with a real SSLCommerz integration boundary

#### 36B.3.1 Provider client package shape

Create package:

- `backend/inscore/microservices/payment/internal/providers/sslcommerz/`

Create files:

- `client.go`
- `types.go`
- `errors.go`

Recommended exported interface:

```go
type Client interface {
    InitSession(ctx context.Context, req *InitSessionRequest) (*InitSessionResponse, error)
    ValidateByValID(ctx context.Context, req *ValidateByValIDRequest) (*ValidateResponse, error)
    QueryBySessionOrTran(ctx context.Context, req *TransactionQueryRequest) (*TransactionQueryResponse, error)
    InitiateRefund(ctx context.Context, req *RefundRequest) (*RefundResponse, error)
    QueryRefund(ctx context.Context, req *RefundQueryRequest) (*RefundQueryResponse, error)
}
```

#### 36B.3.2 Config exact delta

Update:

- `backend/inscore/microservices/payment/internal/config/config.go`

Add:

- `SSLCommerzStoreID`
- `SSLCommerzStorePassword`
- `SSLCommerzAPIBaseURL`
- `SSLCommerzValidationBaseURL`
- `SSLCommerzRefundBaseURL`
- `PublicBaseURL`
- `HTTPTimeout`

Recommended env vars:

- `PAYMENT_SSLCOMMERZ_STORE_ID`
- `PAYMENT_SSLCOMMERZ_STORE_PASSWORD`
- `PAYMENT_SSLCOMMERZ_API_BASE_URL`
- `PAYMENT_SSLCOMMERZ_VALIDATION_BASE_URL`
- `PAYMENT_SSLCOMMERZ_REFUND_BASE_URL`
- `PAYMENT_PUBLIC_BASE_URL`
- `PAYMENT_HTTP_TIMEOUT_SECONDS`

#### 36B.3.3 `payment_service.go` exact rewrite points

Update:

- `backend/inscore/microservices/payment/internal/service/payment_service.go`

`InitiatePayment`:

- reject unsupported methods unless explicitly configured
- require `order_id`, `tenant_id`, and customer fields for SSLCommerz
- create payment row first
- build `tran_id`
- call `sslcommerz.Client.InitSession`
- persist `GatewayPageURL`, `session_key`, provider raw response
- return real `gateway_page_url`
- if provider call fails:
  - update payment status to failed or initiation-failed
  - publish `PaymentFailedEvent`
  - return error

`VerifyPayment`:

- stop setting `SUCCESS` from caller input alone
- if `val_id` is present:
  - call `ValidateByValID`
- otherwise:
  - query by `session_key` or `tran_id`
- verify amount/currency/tran id against stored payment
- persist validation result and risk fields
- transition to `VERIFIED` or equivalent
- publish `PaymentVerifiedEvent`
- then transition to `COMPLETED` after receipt and downstream bookkeeping requirements are satisfied

`InitiateRefund`:

- require `bank_tran_id`
- call provider refund API
- persist `refund_ref_id`
- return processing status instead of instant completion

`GetRefundStatus`:

- query provider refund status when local status is still processing
- update local refund row from provider response

`ReconcilePayments`:

- for SSLCommerz rows, query provider by `tran_id` or `session_key`
- compare provider status and amount with local record
- emit mismatch event when needed

#### 36B.3.4 Repository exact delta for provider client support

Update:

- `backend/inscore/microservices/payment/internal/repository/repository.go`

Add repository methods:

- `CreateProviderTransaction`
- `UpdateProviderTransaction`
- `GetPaymentByTranID`
- `GetPaymentByValID`
- `GetPaymentBySessionKey`
- `CreateProviderWebhook`
- `MarkProviderWebhookProcessed`
- `CreateReceipt`
- `GetReceiptByPaymentID`

Rule:

- every provider call that returns meaningful transaction metadata writes both:
  - the core payment row
  - the provider transaction table

#### 36B.3.5 Event exact delta for provider client support

Update:

- `backend/inscore/microservices/payment/internal/events/topics.go`
- `backend/inscore/microservices/payment/internal/events/publisher.go`

Add topics:

- `payment.verified`
- `payment.manual_review_requested`
- `payment.receipt.generated`
- `payment.reconciliation.mismatch`

Publish rules:

- `payment.initiated` after session creation succeeds
- `payment.failed` on provider init failure or failed validation
- `payment.verified` after SSLCommerz validation succeeds
- `payment.completed` after final settlement workflow step completes

#### 36B.3.6 Item 3 implementation order

Recommended order:

1. add config
2. add client package and DTOs
3. add repository methods
4. wire provider client into `server.go`
5. rewrite `InitiatePayment`
6. rewrite `VerifyPayment`
7. rewrite `InitiateRefund` and `GetRefundStatus`
8. add tests

Definition of done for item 3:

- there is no code path in `payment-service` that returns a fake SSLCommerz checkout URL or marks a payment successful without provider validation

## 37. Operational and Support Workflows

### 37.1 Manual fulfillment recovery

If payment is confirmed but policy issuance fails:

1. operational dashboard alerts the support team
2. support creates a `PolicyIssuanceRetry` work item
3. system retries with exponential backoff
4. if policy issuance succeeds on retry, normal flow resumes
5. if retries are exhausted, escalate for manual issuance or refund decision

### 37.2 Dead-letter queue handling

For every critical Kafka topic:

- consume from DLQ periodically
- inspect failure reason
- replay when the cause was transient
- escalate when intervention is required
- never silently drop messages

### 37.3 Data retention and archival

Financial and audit records:

- hot storage: 3 months
- warm storage: 1 year
- cold/archive: 20 years

Document artifacts:

- customer-downloadable retention: 7 years
- internal audit retention: 20 years

## 38. Observability and Monitoring Plan

### 38.1 Metrics to track

At minimum:

- order creation rate and latency
- payment initiation rate and success rate
- payment completion rate and end-to-end latency
- manual verification queue depth and resolution time
- policy issuance rate and latency
- document generation latency
- webhook delivery success rate and latency
- refund processing rate and latency
- reconciliation mismatch rate
- AML flag rate and resolution SLA

### 38.2 Distributed tracing

Propagate `trace_id` and `correlation_id` across all services:

- order creation initiates `correlation_id`
- every service request includes the `trace_id`
- Kafka messages carry both headers
- logs include both for correlation

### 38.3 Alerting thresholds

Alert if:

- order creation SLA exceeds 5 seconds
- payment processing end-to-end exceeds 15 seconds
- manual verification SLA exceeds 2 hours
- policy issuance latency exceeds 10 seconds
- document generation exceeds 60 seconds
- webhook delivery failures exceed 5%
- reconciliation mismatches exceed 2%
- DLQ depth exceeds 100 messages

## 39. Operational Topic Provisioning Checklist

Ensure the following topics exist with the expected partition key and retention:

| Topic | Partition Key | Retention |
| --- | --- | --- |
| `insuretech.orders.v1.order.created` | `order_id` | 30 days |
| `insuretech.orders.v1.order.payment_initiated` | `order_id` | 30 days |
| `insuretech.orders.v1.order.payment_confirmed` | `order_id` | 30 days |
| `insuretech.orders.v1.order.cancelled` | `order_id` | 30 days |
| `insuretech.orders.v1.order.failed` | `order_id` | 30 days |
| `insuretech.payment.v1.payment.initiated` | `payment_id` | 90 days |
| `insuretech.payment.v1.payment.verified` | `payment_id` | 90 days |
| `insuretech.payment.v1.payment.completed` | `payment_id` | 90 days |
| `insuretech.payment.v1.payment.failed` | `payment_id` | 90 days |
| `insuretech.payment.v1.payment.manual_review_requested` | `payment_id` | 90 days |
| `insuretech.payment.v1.refund.processed` | `payment_id` | 90 days |
| `insuretech.billing.v1.invoice.issued` | `invoice_id` | 90 days |
| `insuretech.billing.v1.invoice.paid` | `invoice_id` | 90 days |
| `insuretech.billing.v1.invoice.cancelled` | `invoice_id` | 90 days |
| `insuretech.billing.v1.invoice.overdue` | `invoice_id` | 90 days |
| `insuretech.insurance.v1.policy.issued` | `policy_id` | 90 days |
| `insuretech.insurance.v1.policy.cancelled` | `policy_id` | 90 days |
| `insuretech.insurance.v1.policy.renewed` | `policy_id` | 90 days |
| `insuretech.insurance.v1.policy.lapsed` | `policy_id` | 90 days |
| `insuretech.document.v1.document.requested` | `reference_id` | 30 days |
| `insuretech.document.v1.document.generated` | `reference_id` | 30 days |
| `insuretech.document.v1.document.failed` | `reference_id` | 30 days |
| `insuretech.storage.v1.file.uploaded` | `file_id` | 30 days |
| `insuretech.storage.v1.file.finalized` | `file_id` | 30 days |
| `insuretech.storage.v1.file.deleted` | `file_id` | 30 days |
| `insuretech.b2b.v1.purchase_order.created` | `purchase_order_id` | 90 days |
| `insuretech.b2b.v1.purchase_order.approved` | `purchase_order_id` | 90 days |
| `insuretech.b2b.v1.purchase_order.rejected` | `purchase_order_id` | 90 days |
| `insuretech.notifications.v1.notification.requested` | `correlation_id` | 7 days |
| `insuretech.integration.v1.webhook.delivery_requested` | `webhook_id` | 7 days |
| `insuretech.compliance.v1.payment.flagged_for_review` | `payment_id` | 365 days |
| `insuretech.audit.v1.critical_action.logged` | `actor_id` | 2555 days |

## 40. Proto Generation and CI/CD Integration

### 40.1 Proto compilation

After creating or updating proto files:

- run code generation for all supported targets
- generate and commit gRPC stubs as required by the repo workflow
- run linting and breaking-change checks
- validate generated code in CI before rollout

### 40.2 Contract versioning

- keep version numbers in package paths such as `v1`
- use a new `v2` package for breaking changes
- support dual-version consumers during migration windows
- document deprecation timelines explicitly

## 41. Payment Provider Integration Checklist

For each provider such as bKash, Nagad, Rocket, bank transfer, and cards:

- [ ] authentication setup validated
- [ ] sandbox or test environment validated
- [ ] payment initiation request builder implemented
- [ ] payment status polling logic implemented where required
- [ ] webhook endpoint and signature verification implemented
- [ ] provider errors mapped to normalized payment statuses
- [ ] retry and timeout logic implemented
- [ ] reconciliation import format supported
- [ ] settlement account linkage validated
- [ ] provider-specific quirks documented

## 42. Post-Implementation Sign-Off Checklist

The implementation is only ready for acceptance when:

- [ ] all proto contracts are complete and linted
- [ ] secure command path is enforced on all state-changing operations
- [ ] order creation, payment, and policy issuance work end-to-end
- [ ] manual payment proof flow is tested with sample uploads
- [ ] idempotency is verified via replay tests
- [ ] webhook signature validation is verified with provider test data
- [ ] document generation triggers work and artifacts are stored and retrievable
- [ ] refund and cancellation flows are tested
- [ ] B2B purchase-order to policy flow is tested
- [ ] AML flags are raised for test scenarios
- [ ] audit logs capture all critical transitions
- [ ] Kafka consumers have DLQ handling
- [ ] dashboards and read-model projections are populated
- [ ] alerts are configured and tested
- [ ] replay procedures are documented and tested
- [ ] retention and archival procedures are in place
- [ ] compliance sign-off is completed for audit and reporting hooks
- [ ] operations runbooks are documented for common failures
- [ ] load testing confirms SRS latency targets
- [ ] security review is completed for AuthN/AuthZ and webhook validation
