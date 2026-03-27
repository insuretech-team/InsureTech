# Insurance Engine — Enterprise Backend Development Prompt

You are working on an enterprise backend project inside a multi-microservice ecosystem. This project follows a **proto-first architecture** to ensure compatibility across multiple technology teams and services.

---

## Architecture Rule — Proto First (Highest Priority)

Proto files are the **highest source of truth** because shared models must remain identical across services and cross-service compatibility must be preserved.

### Source of Truth Priority (Strict Order)

1. **Proto Files + Generated C# Classes** — Located in `root/gen/csharp`. These are the official source for entity structure, field names, field types, and model contracts. Never manually recreate a model if a generated C# class already exists. If documentation conflicts with proto, **proto wins**.

2. **Technical API Contract (HTML Docs)** — Files inside `doc/`. Use for controllers, endpoints, DTOs, request/response schemas, validation rules, CRUD requirements, and error/success responses. HTML defines the exact API contract only — it does not override the proto model.

3. **Business Definition (SRS)** — `doc/SRS_V3/LabAid_InsureTech_SRS_v3.11.md`. Use for business understanding, module purpose, workflow, module dependencies, feature scope, and version planning.

4. **Old Previous Working Project (`insurance_engine-old`)** — Use as a recovery source for reusable service, controller, and repository logic — but only when compatible with proto, documentation, and current architecture.

5. **Internal Architecture Reference** — `backend/policysync`. Use for folder structure, layer separation, service/repository patterns, validation style, migration style, and response middleware style.

**Conflict resolution:** Proto wins for model → HTML wins for API contract → SRS wins for business meaning.

---

## Mandatory Full Audit Before Coding

Before writing any code, perform a complete audit in this order:

1. Read proto-generated classes in `gen/csharp` and identify all available proto models
2. Map proto models to insurance-engine modules
3. Read the SRS
4. Read all HTML docs
5. Review old project code (`insurance_engine-old`)
6. Review current new insurance-engine
7. Review policysync architecture

Then compare: proto vs implementation, proto vs documentation, old vs new project, documentation vs current endpoint behavior.

**Produce a gap analysis covering:** missing modules, missing CRUD, missing endpoint details, response mismatches, missing old logic, and proto mismatches.

> ⚠️ No coding begins before the audit is complete.

---

## Required Audit Output

### Audit Summary
- Proto models found
- Useful `gen/csharp` classes
- HTML endpoint findings
- DTO findings
- SRS module findings
- Old project reusable logic findings
- Current mismatch findings

### Recovery Plan
- Which code is recoverable from `insurance_engine-old`
- Which code needs correction

### Development Priority
- Safest module implementation order
- Version mapping: M1 → M2 → M3

---

## Core Development Rules

### Model Rule
Use only generated proto C# classes. Do not create manual duplicate models, conflicting entity classes, or alternate schemas. If mapping is needed, use adapter/DTO mapping only.

### CRUD Rule
If HTML defines CRUD, implement the full Create, Read, Update, Delete flow using proto-compatible model flow.

### Response Contract Rule
All endpoints must follow the exact HTML response contract:
- **Success:** `message`, `status`, `data`, `meta`, `code`
- **Errors:** `400`, `401`, `403`, `404`, `409`, `500`
- No framework default responses.

### Migration Rule
Maintain migration safety at all times — no unstable migrations, no conflicting schemas, no proto-breaking changes.

### Rebuild Strategy
When logic is missing: (1) recover from `insurance_engine-old` → (2) adapt to proto model → (3) align endpoint to HTML doc → (4) align response to HTML doc → (5) preserve architecture.

---

## Regression Protection — Verified Working APIs Must Not Break

The following flows are currently verified working and are **protected behavior**:

- Individual Beneficiary Create
- Business Beneficiary Create
- Beneficiary Get
- Individual Beneficiary Get By ID

### Mandatory Safe Modification Strategy
Before modifying any beneficiary module or shared layer:
1. Identify all currently working APIs
2. Compare old working code vs current code
3. Keep working logic intact
4. Apply only the required correction
5. Mentally re-verify all existing flows before finalizing

### Internal Error Zero-Tolerance Rule
No existing working endpoint may become a `500 Internal Server Error`. If any change risks this, stop and fix the root cause before continuing.

Before finalizing any change, verify no regressions were introduced in:
- Null references, mapping breaks, proto conversion failures, DTO mismatches
- Missing service registrations, repository injection issues, enum conversion failures
- Serialization issues, response wrapper failures
- All POST, GET, and GET-by-ID endpoints

### Beneficiary Stability Checklist
Explicitly confirm all of the following remain stable after any change:

**Individual Beneficiary:** Create ✓ | Get ✓ | Get By ID ✓

**Business Beneficiary:** Create ✓ | Get ✓ | Get By ID ✓

---

## Final Principle

**Working functionality is production-value. Protect it first. Then improve architecture carefully.**

Architecture improvement must never compromise existing successful behavior. Regression introduced in the name of architectural purity is not acceptable.
