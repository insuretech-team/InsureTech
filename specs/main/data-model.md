# Data Model: AuthN / AuthZ / B2B Completion

**Phase**: 1 — Design  
**Branch**: main  
**Generated**: 2026-03-12  
**Feeds**: plan.md tasks, proto changes

---

## Scope

This document captures:
1. **New fields/messages added to existing proto entities** (no entity creation without
   constitution Principle I clearance).
2. **New DB columns** expressed as `ALTER TABLE … ADD COLUMN` migration candidates.
3. **New Kafka event schemas**.
4. **State machines** for entities that gain lifecycle transitions.

---

## 1. AuthN — New Proto Additions

### 1.1 `GetMe` RPC (new)

**File**: `proto/insuretech/authn/services/v1/auth_service.proto`

```protobuf
// ── GetMe ──────────────────────────────────────────────────────────────────

message GetMeRequest {}

message GetMeResponse {
  string user_id            = 1 [(google.api.field_behavior) = OUTPUT_ONLY];
  string user_type          = 2 [(google.api.field_behavior) = OUTPUT_ONLY];
  string email              = 3 [(google.api.field_behavior) = OUTPUT_ONLY];
  string phone              = 4 [(google.api.field_behavior) = OUTPUT_ONLY];
  string portal             = 5 [(google.api.field_behavior) = OUTPUT_ONLY];
  string tenant_id          = 6 [(google.api.field_behavior) = OUTPUT_ONLY];
  insuretech.authn.entity.v1.UserProfile profile = 7
      [(google.api.field_behavior) = OUTPUT_ONLY];
  repeated string roles     = 8 [(google.api.field_behavior) = OUTPUT_ONLY];
  repeated insuretech.authz.entity.v1.Permission permissions = 9
      [(google.api.field_behavior) = OUTPUT_ONLY];
}

// Add to AuthService:
rpc GetMe(GetMeRequest) returns (GetMeResponse) {
  option (google.api.http) = { get: "/v1/auth/me" };
}
```

**No new DB columns** — reads from existing `authn_schema.users`,
`authn_schema.user_profiles` and AuthZ `GetUserPermissions`.

---

## 2. AuthZ — New Proto Additions

### 2.1 `RotateTokenConfig` RPC (new)

**File**: `proto/insuretech/authz/services/v1/authz_service.proto`

```protobuf
// ── Token Configuration Management ────────────────────────────────────────

message RotateTokenConfigRequest {
  string new_kid                 = 1 [(google.api.field_behavior) = REQUIRED];
  string new_public_key_pem      = 2 [(google.api.field_behavior) = REQUIRED];
  string new_private_key_ref     = 3 [(google.api.field_behavior) = REQUIRED];
  // algorithm defaults to RS256 if empty
  string algorithm               = 4 [(google.api.field_behavior) = OPTIONAL];
}

message RotateTokenConfigResponse {
  insuretech.authz.entity.v1.TokenConfig active_config = 1
      [(google.api.field_behavior) = OUTPUT_ONLY];
  insuretech.authz.entity.v1.TokenConfig retired_config = 2
      [(google.api.field_behavior) = OUTPUT_ONLY];
  string message = 3 [(google.api.field_behavior) = OUTPUT_ONLY];
}

// Add to AuthZService:
rpc RotateTokenConfig(RotateTokenConfigRequest) returns (RotateTokenConfigResponse) {
  option (google.api.http) = {
    post: "/v1/authz/token-configs:rotate"
    body: "*"
  };
}
```

**DB migration** (ALTER only — conforms to Principle XI):
```sql
-- No new columns needed; token_configs table already has:
-- kid, algorithm, public_key_pem, private_key_ref, is_active, created_at, rotated_at
```

---

## 3. B2B — New Proto Additions

### 3.1 Purchase Order lifecycle RPCs (new)

**File**: `proto/insuretech/b2b/services/v1/b2b_service.proto`

```protobuf
// ── Purchase Order Lifecycle ───────────────────────────────────────────────

message ApprovePurchaseOrderRequest {
  string purchase_order_id = 1 [(google.api.field_behavior) = REQUIRED];
  string approved_by       = 2 [(google.api.field_behavior) = REQUIRED];
  string approver_notes    = 3 [(google.api.field_behavior) = OPTIONAL];
}

message ApprovePurchaseOrderResponse {
  insuretech.b2b.entity.v1.PurchaseOrder purchase_order = 1
      [(google.api.field_behavior) = OUTPUT_ONLY];
  string message = 2 [(google.api.field_behavior) = OUTPUT_ONLY];
}

message RejectPurchaseOrderRequest {
  string purchase_order_id = 1 [(google.api.field_behavior) = REQUIRED];
  string rejected_by       = 2 [(google.api.field_behavior) = REQUIRED];
  string reason            = 3 [(google.api.field_behavior) = REQUIRED];
}

message RejectPurchaseOrderResponse {
  insuretech.b2b.entity.v1.PurchaseOrder purchase_order = 1
      [(google.api.field_behavior) = OUTPUT_ONLY];
  string message = 2 [(google.api.field_behavior) = OUTPUT_ONLY];
}

message FulfillPurchaseOrderRequest {
  string purchase_order_id  = 1 [(google.api.field_behavior) = REQUIRED];
  string fulfilled_by       = 2 [(google.api.field_behavior) = REQUIRED];
  string payment_reference  = 3 [(google.api.field_behavior) = OPTIONAL];
}

message FulfillPurchaseOrderResponse {
  insuretech.b2b.entity.v1.PurchaseOrder purchase_order = 1
      [(google.api.field_behavior) = OUTPUT_ONLY];
  string message = 2 [(google.api.field_behavior) = OUTPUT_ONLY];
}

// Add to B2BService:
rpc ApprovePurchaseOrder(ApprovePurchaseOrderRequest)
    returns (ApprovePurchaseOrderResponse) {
  option (google.api.http) = {
    post: "/v1/b2b/purchase-orders/{purchase_order_id}:approve"
    body: "*"
  };
}

rpc RejectPurchaseOrder(RejectPurchaseOrderRequest)
    returns (RejectPurchaseOrderResponse) {
  option (google.api.http) = {
    post: "/v1/b2b/purchase-orders/{purchase_order_id}:reject"
    body: "*"
  };
}

rpc FulfillPurchaseOrder(FulfillPurchaseOrderRequest)
    returns (FulfillPurchaseOrderResponse) {
  option (google.api.http) = {
    post: "/v1/b2b/purchase-orders/{purchase_order_id}:fulfill"
    body: "*"
  };
}
```

**DB migration** (ALTER only — Principle XI):
```sql
-- purchase_orders table already has status column.
-- Add approved_by, rejected_by, reason, payment_reference columns:
ALTER TABLE b2b_schema.purchase_orders ADD COLUMN IF NOT EXISTS approved_by UUID;
ALTER TABLE b2b_schema.purchase_orders ADD COLUMN IF NOT EXISTS approved_at TIMESTAMPTZ;
ALTER TABLE b2b_schema.purchase_orders ADD COLUMN IF NOT EXISTS approver_notes TEXT;
ALTER TABLE b2b_schema.purchase_orders ADD COLUMN IF NOT EXISTS rejected_by UUID;
ALTER TABLE b2b_schema.purchase_orders ADD COLUMN IF NOT EXISTS rejected_at TIMESTAMPTZ;
ALTER TABLE b2b_schema.purchase_orders ADD COLUMN IF NOT EXISTS rejection_reason TEXT;
ALTER TABLE b2b_schema.purchase_orders ADD COLUMN IF NOT EXISTS payment_reference TEXT;
ALTER TABLE b2b_schema.purchase_orders ADD COLUMN IF NOT EXISTS fulfilled_at TIMESTAMPTZ;
```

**Proto entity additions** (`proto/insuretech/b2b/entity/v1/purchase_order.proto`):
```protobuf
// Add to PurchaseOrder message:
string  approved_by       = 14 [(insuretech.common.v1.column).name = "approved_by"];
google.protobuf.Timestamp approved_at   = 15 [...]
string  approver_notes    = 16 [(insuretech.common.v1.column).name = "approver_notes"];
string  rejected_by       = 17 [...];
google.protobuf.Timestamp rejected_at   = 18 [...];
string  rejection_reason  = 19 [...];
string  payment_reference = 20 [...];
google.protobuf.Timestamp fulfilled_at  = 21 [...];
```

### 3.2 BulkImportEmployees RPC (new)

**File**: `proto/insuretech/b2b/services/v1/b2b_service.proto`

```protobuf
message EmployeeImportRow {
  string full_name      = 1 [(google.api.field_behavior) = REQUIRED];
  string mobile         = 2 [(google.api.field_behavior) = REQUIRED];
  string department_id  = 3 [(google.api.field_behavior) = REQUIRED];
  string plan_id        = 4 [(google.api.field_behavior) = REQUIRED];
  string nid            = 5 [(google.api.field_behavior) = OPTIONAL];
  string email          = 6 [(google.api.field_behavior) = OPTIONAL];
  string date_of_birth  = 7 [(google.api.field_behavior) = OPTIONAL];
}

message ImportResult {
  int32  row_index    = 1;
  bool   success      = 2;
  string employee_id  = 3;
  string error_code   = 4;
  string error_detail = 5;
}

message BulkImportEmployeesRequest {
  string business_id           = 1 [(google.api.field_behavior) = REQUIRED];
  repeated EmployeeImportRow rows = 2 [(google.api.field_behavior) = REQUIRED];
}

message BulkImportEmployeesResponse {
  repeated ImportResult results     = 1 [(google.api.field_behavior) = OUTPUT_ONLY];
  int32 total_rows                  = 2 [(google.api.field_behavior) = OUTPUT_ONLY];
  int32 success_count               = 3 [(google.api.field_behavior) = OUTPUT_ONLY];
  int32 failure_count               = 4 [(google.api.field_behavior) = OUTPUT_ONLY];
}

// Add to B2BService:
rpc BulkImportEmployees(BulkImportEmployeesRequest)
    returns (BulkImportEmployeesResponse) {
  option (google.api.http) = {
    post: "/v1/b2b/employees:bulk-import"
    body: "*"
  };
}
```

**No new DB columns** — uses existing `employees` table fields.

### 3.3 Kafka events (new)

**File**: `proto/insuretech/b2b/events/v1/` (new event messages):

```protobuf
message PurchaseOrderApprovedEvent {
  string event_id          = 1;
  string purchase_order_id = 2;
  string business_id       = 3;
  string department_id     = 4;
  string plan_id           = 5;
  string approved_by       = 6;
  string approver_notes    = 7;
  google.protobuf.Timestamp timestamp = 8;
}

message PurchaseOrderRejectedEvent {
  string event_id          = 1;
  string purchase_order_id = 2;
  string business_id       = 3;
  string reason            = 4;
  string rejected_by       = 5;
  google.protobuf.Timestamp timestamp = 6;
}

message PurchaseOrderFulfilledEvent {
  string event_id          = 1;
  string purchase_order_id = 2;
  string business_id       = 3;
  string department_id     = 4;
  string plan_id           = 5;
  int32  employee_count    = 6;
  string payment_reference = 7;
  google.protobuf.Timestamp timestamp = 8;
}
```

**Kafka topics**:
```
insuretech.b2b.v1.PurchaseOrderApproved
insuretech.b2b.v1.PurchaseOrderRejected
insuretech.b2b.v1.PurchaseOrderFulfilled
```

---

## 4. State Machines

### 4.1 PurchaseOrder status state machine

```
               ┌──────────────────────────────────────────────┐
               │                                              │
  CreatePO     ▼     :approve         :fulfill                │
   ─────►  SUBMITTED ──────►  APPROVED ──────►  FULFILLED   │
                  │                                          │
                  │          :reject                         │
                  └──────────────────────►  REJECTED         │
                                                             │
  (DRAFT is reserved for future draft-mode POs)             │
  (All invalid transitions return FAILED_PRECONDITION)       │
               └──────────────────────────────────────────────┘
```

### 4.2 Token key state machine (AuthZ)

```
              RotateTokenConfig
  ─────►    ACTIVE    ──────────►   INACTIVE   (rotated_at set)
                                        │
                          (kept for in-flight JWT validation; 
                           deleted only after all tokens expire)
```

---

## 5. Validation Rules

| Entity | Field | Rule |
|--------|-------|------|
| `PurchaseOrder` | `employee_count` | > 0 |
| `PurchaseOrder` | `coverage_amount` | > 0, currency = "BDT" |
| `PurchaseOrder` | status → APPROVED | caller must have `b2b:approve_purchase_order` permission |
| `PurchaseOrder` | status → REJECTED | `reason` is required (non-empty) |
| `EmployeeImportRow` | `mobile` | Bangladesh phone: `^01[3-9]\d{8}$` |
| `EmployeeImportRow` | `nid` | 10 or 17 digits if provided |
| `RotateTokenConfigRequest` | `new_public_key_pem` | Must parse as RSA public key |
| `RotateTokenConfigRequest` | `new_kid` | Must be unique across token_configs |
| `GetMe` | authentication | Bearer JWT required; no anonymous call |

---

## 6. Index candidates (migration only)

```sql
-- B2B purchase orders: status-filtered queries
CREATE INDEX IF NOT EXISTS idx_purchase_orders_business_status
    ON b2b_schema.purchase_orders (business_id, status);

CREATE INDEX IF NOT EXISTS idx_purchase_orders_department_status
    ON b2b_schema.purchase_orders (department_id, status);

-- AuthZ token configs: active lookup (unique partial index)
CREATE UNIQUE INDEX IF NOT EXISTS idx_token_configs_active
    ON authz_schema.token_configs (is_active)
    WHERE is_active = true;
```
