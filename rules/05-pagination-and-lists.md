# Rule 05: Pagination & List Endpoints

**Scope:** All GET endpoints returning collections  
**Priority:** 🟠 HIGH

---

## Standardization Requirement

Two competing pagination schemas exist:

| Schema | Fields | Problem |
|--------|--------|---------|
| `PageResponse` | `page`, `page_size`, `total_pages`, `total_items` (**string**), `next_page_token` | `total_items` typed as `string` — should be `integer` |
| `PaginationResponse` | `total_items` (**int32**), `total_pages`, `current_page`, `page_size`, `has_next`, `has_previous` | Different field names (`current_page` vs `page`) |

**These two schemas MUST be merged into ONE standard.**

---

## The Single Pagination Standard

### Query Parameters (Request)

```yaml
parameters:
  - name: page
    in: query
    required: false
    schema:
      type: integer
      minimum: 1
      default: 1
    description: Page number (1-based)
  - name: page_size
    in: query
    required: false
    schema:
      type: integer
      minimum: 1
      maximum: 100
      default: 20
    description: Number of items per page (max 100)
  - name: sort_by
    in: query
    required: false
    schema:
      type: string
    description: Field name to sort by (e.g. "created_at", "name")
  - name: sort_order
    in: query
    required: false
    schema:
      type: string
      enum: [asc, desc]
      default: desc
    description: Sort direction
  - name: search
    in: query
    required: false
    schema:
      type: string
    description: Full-text search query (service-specific fields)
```

### Response Shape (via ResponseMeta.pagination)

All list responses MUST use the envelope from Rule 01 with `meta.pagination`:

```json
{
  "success": true,
  "data": {
    "items": [
      { "policy_id": "pol_1", ... },
      { "policy_id": "pol_2", ... }
    ]
  },
  "error": null,
  "meta": {
    "request_id": "req_abc",
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total_pages": 5,
      "total_items": 98,
      "has_next": true,
      "has_previous": false
    }
  }
}
```

### Canonical PaginationMeta Schema (DELETE PageResponse and PaginationResponse)

```yaml
PaginationMeta:
  type: object
  required:
    - page
    - page_size
    - total_pages
    - total_items
    - has_next
    - has_previous
  properties:
    page:
      type: integer
      description: Current page number (1-based)
    page_size:
      type: integer
      description: Number of items in this response
    total_pages:
      type: integer
      description: Total number of pages available
    total_items:
      type: integer
      format: int64
      description: Total number of items across all pages
    has_next:
      type: boolean
      description: Whether a next page exists
    has_previous:
      type: boolean
      description: Whether a previous page exists
    next_page_token:
      type: string
      nullable: true
      description: Cursor token for cursor-based pagination (optional)
```

---

## List Response Data Shape

Every list endpoint returns `data.items` — an array. Never return a bare array.

```yaml
# CORRECT — data always has a named key
PoliciesListData:
  type: object
  properties:
    items:
      type: array
      items:
        $ref: '#/components/schemas/Policy'

# WRONG — bare array
# data: [ {...}, {...} ]
```

**Why:** A bare array makes it impossible to add new top-level fields (e.g., `summary`, `aggregates`) in the future without breaking clients.

---

## Empty List Response

When no results found, return **200 OK** (NOT 404) with empty items:

```json
{
  "success": true,
  "data": {
    "items": []
  },
  "error": null,
  "meta": {
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total_pages": 0,
      "total_items": 0,
      "has_next": false,
      "has_previous": false
    }
  }
}
```

**404 is for a specific resource not found, not for empty collections.**

---

## Filtering Standards

Use query parameters for filtering. Never use POST body for GET queries.

```
GET /v1/policies?status=ACTIVE&product_id=prod_123&page=1&page_size=20
GET /v1/claims?status=PENDING&from_date=2024-01-01&to_date=2024-03-31
GET /v1/payments?user_id=usr_123&method=BKASH&min_amount=1000
```

### Standard Filter Parameters (add where applicable)

```yaml
- name: status
  in: query
  schema:
    type: string
  description: Filter by status
- name: from_date
  in: query
  schema:
    type: string
    format: date
  description: Filter records from this date (inclusive), format YYYY-MM-DD
- name: to_date
  in: query
  schema:
    type: string
    format: date
  description: Filter records to this date (inclusive), format YYYY-MM-DD
- name: search
  in: query
  schema:
    type: string
  description: Full-text search
```

---

## Representative Collection Endpoints

These GET endpoints illustrate the kinds of collection responses that must use the canonical pagination shape:

| Endpoint | Example schema | Required shape |
|----------|---------------|-----|
| `GET /v1/policies` | `UserPoliciesListingResponse` | Use `items[]` + `PaginationMeta` in meta |
| `GET /v1/claims` | various | Use `items[]` + `PaginationMeta` in meta |
| `GET /v1/payments` | `PaymentsListingResponse` | Use `items[]` + `PaginationMeta` in meta |
| `GET /v1/products` | `ProductsProductsListingResponse` | Use `items[]` + `PaginationMeta` in meta |
| `GET /v1/audit-logs` | `AuditLogsRetrievalResponse` | Use `items[]` + `PaginationMeta` in meta |
| `GET /v1/tickets` | `TicketsListingResponse` | Use `items[]` + `PaginationMeta` in meta |
| All other list endpoints | Various | Standardize |

---

## Client Implementation

```typescript
// TypeScript — generic list fetcher works for ALL list endpoints
async function fetchList<T>(
  endpoint: string,
  params: { page?: number; page_size?: number; [key: string]: any }
): Promise<{ items: T[]; pagination: PaginationMeta }> {
  const response = await api.get<ApiResponse<{ items: T[] }>>(endpoint, { params });
  return {
    items: response.data.data?.items ?? [],
    pagination: response.data.meta?.pagination!
  };
}

// Usage — identical for every list endpoint
const { items: policies, pagination } = await fetchList<Policy>('/v1/policies', { page: 1, page_size: 20 });
const { items: claims, pagination } = await fetchList<Claim>('/v1/claims', { status: 'PENDING' });
```

```swift
// Swift — generic list response decoder
func fetchList<T: Decodable>(_ endpoint: String, params: [String: Any]) async throws -> ListResult<T> {
    let response: ApiResponse<ListData<T>> = try await client.get(endpoint, params: params)
    return ListResult(items: response.data?.items ?? [], pagination: response.meta?.pagination)
}
```
