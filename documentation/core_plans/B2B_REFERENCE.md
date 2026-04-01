# B2B Reference

This is the consolidated deep-dive for the B2B vertical. It replaces the split between:

- `B2B_PORTAL_REFERENCE.md`
- `B2B_SERVICE_REFERENCE.md`
- `B2B_AND_PORTALS_STATUS.md`

## 1. Scope

The B2B vertical is currently the strongest end-to-end implementation path in the repository.

It covers:

- portal authentication and session handling
- organization scoping
- organization and member administration
- departments
- employees
- employee bulk upload
- purchase orders
- basic dashboard stats and activity
- document/template integration hooks

## 2. Current Status Summary

### Strongly implemented

- auth/session BFF routes
- organization CRUD and approval
- organization member management
- department CRUD
- employee CRUD
- employee bulk upload and CSV template download
- purchase-order list/create/catalog
- dashboard stats and activity
- organization profile settings

### Partially implemented

- settings workflow tab
- settings notification tab
- document workflows beyond route/client presence

### Mostly presentation-only today

- claims page
- payments page
- policies page
- insurance plans page
- billing/invoices page

## 3. Frontend Architecture

## 3.1 Stack

- Next.js 16
- React 19
- App Router
- local BFF pattern through `app/api`
- local SDK wrappers under `src/lib/sdk`

## 3.2 Folder Shape

Primary areas:

- `b2b_portal/app/`
- `b2b_portal/app/api/`
- `b2b_portal/components/`
- `b2b_portal/src/lib/sdk/`
- `b2b_portal/src/lib/auth/`
- `b2b_portal/src/hooks/`

## 3.3 Portal auth model

The portal does not directly trust browser state. It resolves session and context through:

- auth BFF routes in `app/api/auth`
- server-side session helpers under `src/lib/auth`
- request header propagation through `src/lib/sdk/session-headers.ts`

Important behavior:

- role is read from the current session
- super admin versus B2B admin changes organization scope behavior
- B2B admin is usually pinned to `organisations/me`

## 3.4 Organization scoping pattern

The most important B2B frontend pattern is the `organisations/me` resolution flow.

Why it matters:

- super admins can select among organizations
- B2B admins should not manually pick foreign organizations
- employee, department, and settings views depend on this resolution

Key evidence:

- `b2b_portal/app/api/organisations/me/route.ts`
- `b2b_portal/components/dashboard/employees/employees-table.tsx`
- `b2b_portal/components/dashboard/departments/Departments.tsx`
- `b2b_portal/components/dashboard/settings/partials/organization-form.tsx`

## 4. Portal BFF Surface

The B2B portal currently contains `38` route handlers under `app/api`.

## 4.1 Auth routes

- login
- logout
- refresh
- send OTP / verify OTP
- send email OTP / verify email
- change password
- session lookup
- session list / revoke
- TOTP
- profile
- profile photo URL

## 4.2 Organization routes

- list organizations
- get organization
- resolve `me`
- create/update/delete
- approve organization
- list members
- add member
- remove member
- assign admin
- create admin

## 4.3 Department routes

- list
- get
- create
- update
- delete

## 4.4 Employee routes

- list
- get
- create
- update
- delete
- bulk upload
- template download

## 4.5 Purchase-order routes

- list
- get
- create
- catalog

## 4.6 Dashboard routes

- stats
- activity

## 4.7 Document routes

- documents
- document templates
- document download

## 5. Key Portal Components

## 5.1 Organizations

Key components:

- `components/dashboard/organisations/Organisations.tsx`
- `components/organisations/org-detail-panel.tsx`
- `components/organisations/org-member-panel.tsx`
- `components/modals/add-organisation-modal.tsx`

What is real:

- CRUD list refresh
- open detail panel
- approval
- member list/add/remove
- admin assignment

## 5.2 Departments

Key components:

- `components/dashboard/departments/Departments.tsx`
- `components/modals/add-department-modal.tsx`

What is real:

- super-admin versus B2B-admin scope handling
- create/update/delete through BFF routes

## 5.3 Employees

Key components:

- `components/dashboard/employees/employees-table.tsx`
- `components/modals/add-employee-modal.tsx`
- `components/modals/bulk-upload-employee-modal.tsx`

What is real:

- organization-scoped loading
- CRUD refresh pattern
- bulk upload modal
- template download
- delete action

## 5.4 Purchase orders

Key components:

- `components/dashboard/purchase-orders/purchase-orders.tsx`
- `components/modals/add-purchase-order-modal.tsx`

What is real:

- live list load
- catalog load
- department load
- create flow
- summary cards derived from returned data

## 5.5 Settings

Key components:

- `components/dashboard/settings/settings.tsx`
- `components/dashboard/settings/partials/organization-form.tsx`
- `components/dashboard/settings/partials/workflow-form.tsx`
- `components/dashboard/settings/partials/notification-form.tsx`

Current reality:

- organization form is live-backed through `organisationClient.getMe()` and update calls
- workflow and notification tabs are still based on local arrays

## 6. Backend Service Architecture

## 6.1 Service location

- `backend/inscore/microservices/b2b`
- entrypoint: `backend/inscore/microservices/b2b/cmd/server/main.go`

## 6.2 Runtime wiring

The B2B service currently wires:

- config loading
- `services.yaml` port resolution
- DB manager initialization
- Kafka producer
- event publisher
- repository
- service layer
- AuthZ client
- AuthZ gRPC interceptor
- consumer startup for relevant topics

This is a mature service bootstrap, not a thin placeholder.

## 6.3 gRPC surface

The B2B backend reference identifies `21` methods grouped across:

- organization CRUD
- organization members/admins
- department CRUD
- employee CRUD
- purchase orders and catalog
- bootstrap flow

Important current caveat:

- gateway `PATCH` and `DELETE` for purchase orders are explicit `501` paths because those RPCs are not yet defined in the proto.

## 6.4 Authorization model

B2B depends on AuthZ at two levels:

- gateway middleware
- service-level gRPC interceptor

The interceptor resolves:

- action from method
- domain from portal/tenant/org context
- object/service prefix

This is one reason the B2B path is more trustworthy than mock-only portal areas.

## 6.5 Eventing

The B2B service already publishes and consumes Kafka events.

Published/consumed concerns include:

- organization creation
- admin assignment
- organization approval
- org member changes
- authn/authz-related reactions

The current service also explicitly avoids crashing when Kafka producer init fails, which is an important production-oriented detail.

## 7. Gateway Integration

The gateway already exposes B2B routes for:

- organizations
- members
- admins
- departments
- employees
- bulk upload
- purchase orders
- purchase-order catalog

Important consequence:

The B2B portal is not calling private service code directly. It follows the full contract chain:

portal UI -> portal BFF -> gateway -> AuthN/AuthZ -> B2B service

## 8. Real Versus Mock In The B2B Portal

## 8.1 Real end-to-end paths

- organizations
- org members
- departments
- employees
- employee bulk upload
- purchase orders
- dashboard stats
- dashboard activity
- organization settings

## 8.2 Partial or future paths

- document flows beyond route presence
- workflow settings persistence
- notification settings persistence

## 8.3 Mostly mock UI pages

- claims
- payments
- policies
- billing/invoices
- insurance plans

These pages exist, but they should not be described the same way as the core B2B management flows.

## 9. Important Corrections Against Older B2B Docs

Statements that are now outdated:

- bulk upload is missing
- purchase orders are only planned
- document/template routes are not present
- quotations page is missing
- settings are fully mock

More accurate replacements:

- bulk upload exists in BFF, UI, gateway, and backend form
- purchase-order create/list/catalog are already implemented
- document/template route surfaces exist
- quotations currently redirects into purchase-order flow
- settings are mixed: organization tab is live, workflow/notification are still local

## 10. Canonical Use

Use this document when you need the detailed B2B picture.

Use `IMPLEMENTATION_BASELINE.md` when you need to place B2B in the full platform context.
