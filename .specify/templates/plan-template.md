# Implementation Plan: [FEATURE]

**Branch**: `[###-feature-name]` | **Date**: [DATE] | **Spec**: [link]
**Input**: Feature specification from `/specs/[###-feature-name]/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/plan-template.md` for the execution workflow.

## Summary

[Extract from feature spec: primary requirement + technical approach from research]

## Technical Context

<!--
  ACTION REQUIRED: Replace the content in this section with the technical details
  for the project. The structure here is presented in advisory capacity to guide
  the iteration process.
-->

**Language/Version**: [e.g., Python 3.11, Swift 5.9, Rust 1.75 or NEEDS CLARIFICATION]  
**Primary Dependencies**: [e.g., FastAPI, UIKit, LLVM or NEEDS CLARIFICATION]  
**Storage**: [if applicable, e.g., PostgreSQL, CoreData, files or N/A]  
**Testing**: [e.g., pytest, XCTest, cargo test or NEEDS CLARIFICATION]  
**Target Platform**: [e.g., Linux server, iOS 15+, WASM or NEEDS CLARIFICATION]
**Project Type**: [e.g., library/cli/web-service/mobile-app/compiler/desktop-app or NEEDS CLARIFICATION]  
**Performance Goals**: [domain-specific, e.g., 1000 req/s, 10k lines/sec, 60 fps or NEEDS CLARIFICATION]  
**Constraints**: [domain-specific, e.g., <200ms p95, <100MB memory, offline-capable or NEEDS CLARIFICATION]  
**Scale/Scope**: [domain-specific, e.g., 10k users, 1M LOC, 50 screens or NEEDS CLARIFICATION]

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

Verify each gate before proceeding. Violations must be documented in the Complexity
Tracking section below with explicit justification.

| # | Principle | Gate Question | Status |
|---|-----------|---------------|--------|
| I | Proto-First | Is a `.proto` file committed and buf-lint-clean before any service/portal code? | [ ] |
| II | Polyglot Ownership | Does the feature stay within its designated runtime (Go / C# / TS / Swift / Kotlin)? | [ ] |
| III | REST API Standard | Do all new HTTP paths use plural nouns, no verbs, and conform to the OpenAPI 3.0+ envelope? | [ ] |
| IV | Event-Driven | Are cross-domain state changes communicated via Kafka with canonical topic names and outbox pattern? | [ ] |
| V | Security & Compliance | Are PII fields AES-256 encrypted, auth/authz enforced, SAST passing, and any IDRA/AML rules accounted for? | [ ] |
| VI | VSA & Tests | Are failing tests written first, coverage target ≥ 80%, and Testcontainers used for infra? | [ ] |
| VII | Observability | Are OTel traces, Prometheus RED metrics, health probes, and structured logs included? | [ ] |
| VIII | Multi-Tenancy | Does every DB row carry `tenant_id` and does every query filter by `tenant_id`? | [ ] |
| IX | Versioning | Are breaking changes gated behind a new major version with deprecation headers? | [ ] |
| X | Simplicity | Is there speculative complexity or new abstractions that can be removed? | [ ] |
| XI | Hybrid SQL Migration | If schema changes: do SQL files contain only ALTER/index, no CREATE TABLE or DML? Is proto-gen freshness gate satisfied? | [ ] |
| XII | Platform Surface | Is the feature scoped to the correct portal/app (system / insurer / partners / business / regulatory / customer-app / agent-app)? | [ ] |
| XIII | Response Envelope | Does every response use the canonical `{success, data, error, meta}` envelope? | [ ] |
| XIV | HTTP Status Codes | Are correct status codes used (201 for creates, 204 for deletes, 400 vs 422 distinction)? | [ ] |
| XV | Error Handling | Do errors use the standard Error schema with domain-prefixed `UPPER_SNAKE_CASE` codes? | [ ] |
| XVI | API Security Declaration | Does every endpoint declare its security scheme explicitly (BearerAuth / ApiKeyAuth / `security: []`)? | [ ] |
| XVII | Pagination & List | Do list endpoints use `PaginationMeta` in `meta.pagination` with a named `items` key in `data`? | [ ] |
| XVIII | DI-Ready Contract | Do responses satisfy the five contract invariants (envelope, status codes, security, error shape, pagination)? | [ ] |
| XIX | URL & Naming | Are paths kebab-case, path params snake_case, query params snake_case, JSON fields snake_case? | [ ] |
| XX | Null & Optional | Are all schema fields classified as `required` or `nullable: true`? No silent field omission? | [ ] |

## Project Structure

### Documentation (this feature)

```text
specs/[###-feature]/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)
<!--
  ACTION REQUIRED: Replace the placeholder tree below with the concrete layout
  for this feature. Delete unused options and expand the chosen structure with
  real paths (e.g., apps/admin, packages/something). The delivered plan must
  not include Option labels.
-->

```text
# [REMOVE IF UNUSED] Option 1: Single project (DEFAULT)
src/
├── models/
├── services/
├── cli/
└── lib/

tests/
├── contract/
├── integration/
└── unit/

# [REMOVE IF UNUSED] Option 2: Web application (when "frontend" + "backend" detected)
backend/
├── src/
│   ├── models/
│   ├── services/
│   └── api/
└── tests/

frontend/
├── src/
│   ├── components/
│   ├── pages/
│   └── services/
└── tests/

# [REMOVE IF UNUSED] Option 3: Mobile + API (when "iOS/Android" detected)
api/
└── [same as backend above]

ios/ or android/
└── [platform-specific structure: feature modules, UI flows, platform tests]
```

**Structure Decision**: [Document the selected structure and reference the real
directories captured above]

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| [e.g., 4th project] | [current need] | [why 3 projects insufficient] |
| [e.g., Repository pattern] | [specific problem] | [why direct DB access insufficient] |
