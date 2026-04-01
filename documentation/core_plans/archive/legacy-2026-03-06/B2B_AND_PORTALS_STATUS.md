# B2B And Portals Status

This document compares the three web portals by integration depth.

## 1. Executive Summary

The portals are not at the same maturity level.

### Current ranking

1. `b2b_portal` - real integration path
2. `system_portal` - partial integration shell with substantial demo data
3. `customer_portal` - mostly static/demo UI

Any documentation that describes all three as equally implemented is wrong.

## 2. B2B Portal

## 2.1 Implementation footprint

- `15` page routes
- `38` API route handlers
- local SDK client layer under `b2b_portal/src/lib/sdk`
- session and portal-header propagation
- role-aware organization scoping

## 2.2 What is clearly live-backed

### Auth and session

- login
- logout
- refresh
- current session
- session listing/revocation
- OTP and TOTP related flows
- profile and profile-photo helpers

Evidence:

- `b2b_portal/app/api/auth/*`
- `b2b_portal/src/lib/sdk/auth-client.ts`

### Organizations and membership

- list organizations
- get organization
- resolve current organization (`/organisations/me`)
- create/update/delete
- approve pending organization
- list/add/remove members
- assign admin
- create admin user

Evidence:

- `b2b_portal/app/api/organisations/*`
- `b2b_portal/src/lib/sdk/organisation-client.ts`
- `b2b_portal/components/dashboard/organisations/Organisations.tsx`
- `b2b_portal/components/organisations/org-detail-panel.tsx`
- `b2b_portal/components/organisations/org-member-panel.tsx`

### Departments

- scoped list by organization
- create/update/delete
- super-admin versus B2B-admin organization scope handling

Evidence:

- `b2b_portal/app/api/departments/*`
- `b2b_portal/components/dashboard/departments/Departments.tsx`

### Employees

- list
- create/update/delete
- organization-scoped loading
- bulk upload
- CSV template download

Evidence:

- `b2b_portal/app/api/employees/*`
- `b2b_portal/components/dashboard/employees/employees-table.tsx`
- `b2b_portal/components/modals/bulk-upload-employee-modal.tsx`

### Purchase orders

- list
- create
- get catalog
- dashboard-level summary cards computed from live results

Evidence:

- `b2b_portal/app/api/purchase-orders/*`
- `b2b_portal/components/dashboard/purchase-orders/purchase-orders.tsx`

### Dashboard

- live stats
- live activity

Evidence:

- `b2b_portal/app/api/dashboard/stats/route.ts`
- `b2b_portal/app/api/dashboard/activity/route.ts`
- `b2b_portal/components/dashboard/stats-cards/stats-cards.tsx`
- `b2b_portal/components/dashboard/overview-activity/overview-activity.tsx`

### Document and template management

- document routes exist
- template routes exist
- docgen client exists

Evidence:

- `b2b_portal/app/api/documents/*`
- `b2b_portal/app/api/document-templates/*`
- `b2b_portal/src/lib/sdk/docgen-client.ts`
- `b2b_portal/src/lib/sdk/docgen-sdk-client.ts`

## 2.3 What is partially wired

### Settings

There is mixed maturity across settings tabs.

- Organization profile: live-backed
- Workflow rules: local data only
- Notification preferences: local data only

Evidence:

- live:
  - `b2b_portal/components/dashboard/settings/partials/organization-form.tsx`
- mock/static:
  - `b2b_portal/components/dashboard/settings/partials/workflow-form.tsx`
  - `b2b_portal/components/dashboard/settings/partials/notification-form.tsx`
  - `b2b_portal/lib/workflows.ts`
  - `b2b_portal/lib/notifications.ts`

## 2.4 What is still mostly static/mock UI

The following B2B sections exist as pages, but the current implementation is mostly hard-coded or presentation-only compared with the core B2B flows:

- payments
- claims
- policies
- billing/invoices
- insurance plans

Evidence:

- `b2b_portal/components/payments/payments-page.tsx`
- `b2b_portal/components/claims/claims-page.tsx`
- `b2b_portal/components/policies/policies-page.tsx`
- `b2b_portal/components/dashboard/billing-invoices/billing-invoices.tsx`

## 2.5 Important corrections to old B2B docs

Older B2B docs are stale in several places:

- bulk upload is no longer missing
- purchase-order create/list/catalog are no longer just planned
- document routes and template routes already exist
- settings are not uniformly mock; the organization profile tab is wired
- quotations route is not missing; it currently redirects into purchase-order flow

## 3. Customer Portal

## 3.1 Reality

The customer portal is mostly a visual/dashboard prototype today.

### Signals that it is still mostly static

- employee data is hard-coded in component files
- claims are hard-coded in component files
- payments are hard-coded in component files
- policies are hard-coded in component files
- quotations, invoices, notifications, payments, and workflows are stored in local `lib/*.ts` arrays

Representative evidence:

- `customer_portal/components/dashboard/employees/employees-table.tsx`
- `customer_portal/components/claims/claims-page.tsx`
- `customer_portal/components/payments/payments-page.tsx`
- `customer_portal/components/policies/policies-page.tsx`
- `customer_portal/lib/quotations.ts`
- `customer_portal/lib/payments.ts`
- `customer_portal/lib/invoices.ts`
- `customer_portal/lib/notifications.ts`
- `customer_portal/lib/workflows.ts`

## 3.2 Interpretation

This portal should be described as:

- a UI baseline
- a design prototype
- a candidate for future API integration

It should not be described as a live application in the same sense as the B2B portal.

## 4. System Portal

## 4.1 Reality

The system portal has more backend-adjacent scaffolding than the customer portal, but many major screens still use local demo datasets.

### What exists

- SvelteKit application structure
- generated types
- server-side API client helper
- dashboard routing structure

### What still uses demo data

- products page
- product detail page
- policies page
- claims page
- analytics page
- partner dashboards
- top-level dashboard summaries

Representative evidence:

- `system_portal/src/routes/dashboard/products/+page.svelte`
- `system_portal/src/routes/dashboard/products/[id]/+page.svelte`
- `system_portal/src/routes/dashboard/policies/+page.svelte`
- `system_portal/src/routes/dashboard/claims/+page.svelte`
- `system_portal/src/routes/dashboard/analytics/+page.svelte`
- `system_portal/src/routes/dashboard/partners/life/+page.svelte`
- `system_portal/src/routes/dashboard/partners/non-life/+page.svelte`
- `system_portal/src/lib/data_detailed/products_demo.ts`
- `system_portal/src/lib/data_detailed/policies_demo.ts`
- `system_portal/src/lib/data_detailed/claims_demo.ts`
- `system_portal/src/lib/data_detailed/analyticsData.ts`

## 4.2 Interpretation

The system portal is best described as:

- an admin shell
- with real SDK/generated-type foundations
- but still substantial demo-data rendering in key screens

## 5. Practical Documentation Rule

When documenting portal status:

- say B2B is the only strongly integrated portal
- say customer is mostly static/demo
- say system is partially scaffolded with many demo datasets

Anything softer than that is likely to hide real implementation risk.
