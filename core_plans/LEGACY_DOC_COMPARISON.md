# Legacy Document Comparison

This file explains how the old `core_plans` set maps into the current consolidated set.

## 1. Why consolidation was needed

The old top-level set mixed:

- full references
- point-in-time plans
- exploration notes
- generated-style summaries

Many of those were still useful, but several described older assumptions that no longer match the current codebase.

## 2. Main overlap groups

### B2B overlap group

Old docs:

- `B2B_PORTAL_REFERENCE.md`
- `B2B_SERVICE_REFERENCE.md`
- portal-status notes

Current consolidated replacement:

- `B2B_REFERENCE.md`
- `IMPLEMENTATION_BASELINE.md`

### Proto and SDK overlap group

Old docs:

- `PROTO_INDEX.md`
- `PROTO_COMMON_TYPES.md`
- `PROTO_CORE_MODULES.md`
- `PROTO_FILES_SUMMARY.md`
- `COMPLETE_SDK_CONTENT_SUMMARY.md`
- `SDK_DOCUMENTATION.md`
- `PROTO_API_AND_SDK_REFERENCE.md`

Current consolidated replacement:

- `PROTO_AND_SDK_REFERENCE.md`

### Platform-architecture overlap group

Old docs:

- `INSURANCE_SERVICE_ARCHITECTURE.md`
- `POLICY_SYNC_ARCHITECTURE.md`
- portions of `POLISYNC_REFERENCE.md`
- `dbmanager.md`

Current consolidated replacement:

- `IMPLEMENTATION_BASELINE.md`
- `POLISYNC_REFERENCE.md`
- `PAYMENT_ORDER_IMPLEMENTATION_PLAN.md`

## 3. Important corrections

Statements that were accurate in older docs but are now stale on this branch include:

- payment is missing
- fraud is effectively absent
- partner lacks real gateway integration
- bulk employee upload is still only planned
- purchase orders are mostly roadmap-only

The current repo state is more mature than those statements imply.

## 4. Use rule

Use the old docs for history and fine-grained context.

Use the consolidated docs for the current platform picture.
