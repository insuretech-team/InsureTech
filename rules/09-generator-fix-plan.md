# Rule 09: Generator Output Requirements

**Scope:** OpenAPI generation, fixup scripts, validation, and SDK-producing pipeline stages
**Priority:** Critical

---

## Purpose

The generation pipeline must emit OpenAPI and SDK inputs that comply with Rules 01 through 08.

This document defines the required output behavior of the generators. It is not a live implementation-status document.

---

## Generator obligations

### 9.1 Response envelope output

Generators must emit responses that use the standard envelope from Rule 01.

Required properties:

- `success`
- `data`
- `error`
- `meta`

Success schemas must not embed a separate `error` object.

### 9.2 Success code classification

Generators must classify success responses semantically:

- `201` for creation of a persistent resource
- `200` for reads and action-style POSTs
- `204` for delete or no-body actions

### 9.3 Required error responses

Generated operations must declare the appropriate error responses described in Rule 03, including:

- `400`
- `401` when authentication is required
- `403` when authorization applies
- `404` for resource lookups
- `409` for create-conflict scenarios
- `422` for write operations
- `429` where rate limiting applies
- `500`

### 9.4 Security declaration

Each generated operation must explicitly declare its security requirements per Rule 04.

- public endpoints must emit `security: []`
- protected endpoints must emit the applicable security scheme entries

### 9.5 Pagination contract

List responses must use the canonical pagination contract from Rule 05.

Requirements:

- use one canonical pagination schema
- expose pagination through `meta.pagination`
- avoid competing pagination object shapes

### 9.6 Success schema hygiene

Generated success schemas must:

- exclude embedded `error` fields
- avoid thin `message`-only payloads
- declare required fields where applicable
- mark nullable fields explicitly

### 9.7 Examples and headers

Generators should emit:

- representative response examples
- `Location` headers on `201` responses
- stable schema names and references for SDK generation

### 9.8 Validation gate

The pipeline must validate that generated output still satisfies these rules before publishing docs or SDKs.

---

## Canonical schemas to preserve

Generators and fixup scripts must preserve these canonical concepts:

- `ApiResponse`
- `ResponseMeta`
- `PaginationMeta`
- standard `Error`

No fixup script should duplicate or fork these canonical shapes.

---

## Output checklist

Before a generated spec is treated as publishable, it must satisfy all of the following:

- every operation has a valid success code
- every operation has appropriate error responses
- every operation has explicit security behavior
- list endpoints use canonical pagination
- success schemas do not embed error objects
- nullable and required fields are explicit
- generated artifacts remain usable by docs and SDK generators

