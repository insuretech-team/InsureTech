# InsureTech Core Plans

This folder is now the curated documentation entrypoint for the repository as it exists today.

The previous `core_plans` folder had many overlapping documents that mixed three different things:

- architecture intent
- exploratory notes
- implementation status at a point in time

That made it hard to answer basic questions such as:

- what is already implemented?
- what is partially implemented?
- what still runs on mock/demo data?
- which old plan assumptions are now outdated?

This consolidated set keeps the detail, but reorganizes it around the current codebase.

## Read Order

1. `IMPLEMENTATION_BASELINE.md`
2. `B2B_REFERENCE.md`
3. `AUTHN_AUTHZ_REFERENCE.md`
4. `POLISYNC_REFERENCE.md`
5. `PAYMENT_ORDER_IMPLEMENTATION_PLAN.md`
6. `PROTO_AND_SDK_REFERENCE.md`
7. `DOCGEN_MEDIA_ENHANCED_PLAN_V3.md`
8. `FRAUD_PARTNER_IMPLEMENTATION_PLAN.md`
9. `ACTIVE_WORKSTREAMS.md`
10. `LEGACY_DOC_COMPARISON.md`

## Full Picture In Brief

> Last updated: 2026-03-06

The repo is a proto-first insurance platform split across contracts, Go platform services, C# domain services, and three web portals.

### Contracts

- `proto/` is the contract source of truth.
- The repository currently contains approximately `185` proto files.
- Proto field reference for payment, billing, orders: see `PROTO_AND_SDK_REFERENCE.md § 3.5`

### Go platform layer

- `backend/inscore/` contains the API gateway and most platform/infrastructure services.
- There are currently `29` command entrypoints under `backend/inscore/cmd`.
- There are currently `24` service folders under `backend/inscore/microservices`.
- `backend/inscore/configs/services.yaml` is the active service and port registry — includes `payment` (50190), `orders` (50142), `billing` (50195).

### C# domain layer

- `backend/polisync/` is a real, active .NET runtime, not just a design draft.
- `backend/polisync/src` currently contains `14` projects.
- `PoliSync.ApiHost` already maps live gRPC services for products, quotes, orders, policy, claims, commission, underwriting, endorsement, renewal, and refund.

### Gateway layer

- The HTTP gateway in `backend/inscore/cmd/gateway` is already broad.
- The router currently contains approximately `220` `mux.Handle` / `mux.HandleFunc` registrations.
- It fronts auth, B2B, media, document, storage, partner, fraud, product, quote, order, payment, billing, policy, claim, underwriting, commission, and related service surfaces.
- SSLCommerz callback routes: `POST /v1/payments/webhook/sslcommerz`, `/v1/payments/sslcommerz/success`, `/fail`, `/cancel` (public, no auth).

### Portals

- `b2b_portal/` is the most integrated portal path.
  - `15` page routes
  - `38` API route handlers
- `customer_portal/` is mostly UI/demo data today.
  - `10` page routes
- `system_portal/` has partial integration scaffolding but many dashboard pages still use local demo datasets.
  - `11` Svelte page routes

## What This Consolidation Keeps

- implementation detail
- comparison against earlier plans
- portal-by-portal maturity analysis
- contract/runtime ownership notes
- active gaps and architectural risks

## What Was Moved

The older fragmented documents were archived into:

- `core_plans/archive/legacy-2026-03-06/`

They remain useful as historical material, but they are no longer the canonical project description.
