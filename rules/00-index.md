# InsureTech REST API Rules and Standards

**Version:** 1.0
**Date:** 2026-03-13
**Purpose:** Normative API rules for all REST services and all generated API contracts.

---

## How to use this folder

- `rules/` is for rules only.
- Implementation status and repo-state notes live in `../API_PIPELINE_STATUS.md`.
- Pipeline flow, artifacts, and SDK/portal dependencies live in `../API_PIPELINE_REFERENCE.md`.

Read these rules as the target contract that the generators, documentation, SDKs, and backend responses must satisfy.

---

## Rule files

| File | Title | Focus |
|------|-------|-------|
| `01-response-envelope.md` | Standard Response Envelope | canonical success/error envelope |
| `02-http-status-codes.md` | HTTP Status Code Standards | success and error status semantics |
| `03-error-handling.md` | Error Handling Standards | canonical error object and response behavior |
| `04-security-authentication.md` | Security and Authentication | endpoint auth declaration rules |
| `05-pagination-and-lists.md` | Pagination and List Endpoints | canonical list and pagination shape |
| `06-dependency-injection-and-testing.md` | DI and Frontend Testing | client-consumable contract guarantees |
| `07-naming-and-url-design.md` | URL Design and Naming | resource and action path conventions |
| `08-null-empty-and-optional-data.md` | Null, Empty, and Optional Data | nullable vs required payload rules |
| `09-generator-fix-plan.md` | Generator Output Requirements | generator obligations for producing compliant OpenAPI |
| `dbrules.md` | Database Rules | database naming and schema conventions |
| `ground_truth.md` | Ground Truth | authoritative baseline notes |

---

## Contract summary

Every REST response must comply with these baseline rules:

1. use the standard envelope from Rule 01
2. use semantically correct HTTP status codes from Rule 02
3. return errors only through the error contract in Rule 03
4. declare endpoint security explicitly per Rule 04
5. use the canonical list and pagination shape from Rule 05
6. preserve client and SDK usability per Rule 06
7. follow stable URL and naming conventions from Rule 07
8. make required and nullable data explicit per Rule 08
9. ensure generators emit outputs that satisfy all of the above per Rule 09

---

## Quick reference

### Canonical response envelope

```json
{
  "success": true,
  "data": { },
  "error": null,
  "meta": {
    "request_id": "req_123",
    "pagination": null
  }
}
```

### Baseline status code rules

- create resource: `201`
- read or action success: `200`
- delete with no body: `204`
- malformed request: `400`
- unauthenticated: `401`
- unauthorized: `403`
- not found: `404`
- conflict: `409`
- validation failure: `422`
- rate limit: `429`
- unexpected server error: `500`

---

## Companion documents

These are not rules, but they are the main operational companions:

- `../API_PIPELINE_REFERENCE.md`
- `../API_PIPELINE_STATUS.md`

