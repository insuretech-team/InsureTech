# Proto, API, And SDK Reference

This document replaces the previous cluster of proto summaries and SDK notes.

## 1. Source Of Truth

The project is proto-first.

The correct order of truth is:

1. `proto/`
2. gateway exposure and generated OpenAPI under `api/`
3. generated SDKs under `sdks/`
4. service and portal consumption code

Everything else is explanatory documentation.

## 2. Current Contract Surface

The repo currently contains about `180` `.proto` files.

The contract set covers at least these functional groups:

- identity and authorization
- B2B
- products and pricing
- quotations
- orders
- payments
- policies
- claims
- underwriting
- renewal
- refund
- partner
- fraud
- beneficiary
- document and storage
- workflow and task
- report and analytics
- tenant
- support
- voice and webRTC
- insurer and insurance CRUD

## 3. API Generation Reality

`api/` is not hand-authored primary design documentation. It is generated and derived from proto contracts plus gateway mapping.

Important files:

- `api/openapi.yaml`
- `api/ENDPOINT_MAP.md`
- `api/docs/`

Current state:

- endpoint coverage is broad and current
- the generated docs already include many newer domains and schemas
- the API layer should now be treated as a projection of the contract set, not a separate design stream

## 4. SDK Reality

Current generated SDKs include:

- `sdks/insuretech-typescript-sdk/`
- `sdks/insuretech-go-sdk/`

Current usage patterns:

- B2B portal consumes the packaged TypeScript SDK via local tarball dependency
- system portal also consumes the TypeScript SDK
- generated proto/type output is checked into portals to simplify frontend consumption

## 5. Portal Consumption Patterns

### B2B portal

The B2B portal has the most mature client-consumption layer:

- `b2b_portal/src/lib/sdk/auth-client.ts`
- `b2b_portal/src/lib/sdk/organisation-client.ts`
- `b2b_portal/src/lib/sdk/department-client.ts`
- `b2b_portal/src/lib/sdk/employee-client.ts`
- `b2b_portal/src/lib/sdk/purchase-order-client.ts`
- `b2b_portal/src/lib/sdk/docgen-client.ts`

It also has API-route mediation in `b2b_portal/app/api`, which means browser code does not directly speak to the gateway in most flows.

### System portal

The system portal has a lighter server-side access pattern:

- `system_portal/src/lib/server/api.ts`

This shows intent toward live integration, but the UI layer still often consumes local demo data instead of the API client.

## 6. Developer Rules

### If you need to change a field or RPC

- edit `proto/` first
- regenerate
- then update consuming code

### If you need to understand the HTTP API

- read `proto/`
- then confirm generated shape in `api/ENDPOINT_MAP.md`

### If you need to understand the frontend client surface

- read the relevant SDK client wrapper or generated SDK module

## 7. What Was Wrong With The Old Proto Docs

The previous structure had many separate files:

- index
- quick reference
- common types
- core modules
- support modules
- file summaries
- SDK summaries

That created duplication and made it too easy for explanations to diverge from the live proto tree.

The new rule is simpler:

- use one reference document
- point back to `proto/`
- treat generated material as downstream evidence

## 8. Canonical Statement

If a contract question matters, the answer must be verified in `proto/`, not inferred from a narrative plan document.
