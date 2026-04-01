# KYC-FLVE Alignment Plan

> **Last updated:** 2025-07-25
> **Status:** ✅ IMPLEMENTED (Phases 1–7 complete, build verified)
> **Projects:** `InsureTech` (proto-first insurance platform) + `LabaidAi-Retina` (FLVE face verification engine)
> **Priority stack:** InsureTech proto contracts > persistence (users, user_profiles, auth, kyc_verifications) > ApiResponse\<T\> envelope > FLVE as provider

---

## Goal

Align InsureTech KYC with the FLVE Hugging Face Space deployment from LabaidAi-Retina so that:

1. **InsureTech proto contracts remain the canonical product API.** Clients (mobile, web portals, BFF) only see InsureTech-shaped messages.
2. **`users`, `user_profiles`, `auth`, and `kyc_verifications` stay the system of record.** FLVE owns no identity data.
3. **The gateway wraps every HTTP response in the universal `ApiResponse<T>` envelope.** The recent commit standardizing this pattern is preserved and extended to KYC endpoints.
4. **FLVE Space becomes an implementation provider behind the KYC flow**, not the public contract. Raw FLVE payloads are stored in the DB but never surfaced as the top-level API shape.

---

## Architecture Overview

```
┌──────────────────────────────────────────────────────────────────┐
│  Mobile / Web Portal / BFF                                       │
│  (only speaks InsureTech proto-shaped JSON inside ApiResponse<T>)│
└────────────────────────┬─────────────────────────────────────────┘
                         │ HTTPS
                         ▼
┌────────────────────────────────────────────────────┐
│  InsureTech Gateway (Go net/http)                  │
│  ─ POST /v1/auth/users/{user_id}/kyc              │
│  ─ POST /v1/auth/users/{user_id}/kyc:submit-frame │
│  ─ POST /v1/auth/users/{user_id}/kyc:complete     │
│  ─ GET  /v1/auth/users/{user_id}/kyc              │
│  ─ POST /v1/auth/kyc/{kyc_id}:approve             │
│  All responses wrapped in ApiResponse<T>           │
└────────────────────────┬───────────────────────────┘
                         │ gRPC / in-process
                         ▼
┌────────────────────────────────────────────────────┐
│  AuthN Service  (kyc_orchestrator_service.go)      │
│  ─ owns kyc_verifications, user_profiles           │
│  ─ FLVEAdapter (HTTP client)                       │
└────────────────────────┬───────────────────────────┘
                         │ HTTPS (bearer token)
                         ▼
┌────────────────────────────────────────────────────┐
│  FLVE HuggingFace Space  (FastAPI, T4 GPU)        │
│  ─ POST /ekyc/start                               │
│  ─ POST /ekyc/frame     (multipart file upload)   │
│  ─ POST /ekyc/complete                             │
│  ─ GET  /ekyc/status/{session_id}                 │
│  ─ CDN upload (DigitalOcean Spaces)               │
└────────────────────────────────────────────────────┘
```

---

## Current State

### InsureTech user-facing KYC orchestration API

| Endpoint | Source |
|---|---|
| `POST /v1/auth/users/{user_id}/kyc` | `proto/insuretech/authn/services/v1/auth_service.proto` |
| `GET  /v1/auth/users/{user_id}/kyc` | `proto/insuretech/authn/services/v1/auth_service.proto` |
| `POST /v1/auth/users/{user_id}/kyc:submit-frame` | `proto/insuretech/authn/services/v1/core.proto` |
| `POST /v1/auth/users/{user_id}/kyc:complete` | `proto/insuretech/authn/services/v1/core.proto` |
| `POST /v1/auth/kyc/{kyc_id}:approve` | `proto/insuretech/authn/services/v1/core.proto` |

### InsureTech persistence fields that can hold provider state

**`kyc_verifications`** (proto: `insuretech.kyc.entity.v1.KYCVerification`):

| Field | Purpose |
|---|---|
| `id` | UUID primary key — canonical InsureTech KYC ID |
| `provider` | e.g. `"FLVE_HF"` |
| `provider_reference` | FLVE `session_id` |
| `documents` | JSONB, encrypted |
| `verification_result` | JSONB, encrypted — stores raw FLVE payloads |
| `status` | `PENDING → IN_PROGRESS → VERIFIED / REJECTED / EXPIRED` |

**`user_profiles`** (proto: `insuretech.authn.entity.v1.UserProfile`):

| Field | Purpose |
|---|---|
| `profile_photo_url` | CDN URL from FLVE capture step |
| `kyc_verified` | boolean flag |
| `kyc_verified_at` | timestamp |
| `id_upload_front_url` | NID front — potential reference image for FLVE |
| `id_upload_back_url` | NID back |
| `photograph_selfie_url` | selfie — potential reference image for FLVE |

### Current ExternalKYCClient (narrow interface)

File: `backend/inscore/microservices/authn/internal/service/kyc_external_client.go`

```go
type ExternalKYCClient interface {
    StartKYCVerification(ctx, *StartKYCVerificationRequest, ...grpc.CallOption) (*StartKYCVerificationResponse, error)
    UploadDocument(ctx, *UploadDocumentRequest, ...grpc.CallOption) (*UploadDocumentResponse, error)
    VerifyKYC(ctx, *VerifyKYCRequest, ...grpc.CallOption) (*VerifyKYCResponse, error)
}
```

This interface maps InsureTech's generic KYC service, **not** FLVE's richer `/ekyc/*` flow. The existing `flveExternalKYCClient` (`flve_external_client.go`) shoehorns FLVE calls into this interface by:
- Sending `StartKYCVerification` → FLVE `/ekyc/start` (loses `UserContext` fields)
- Sending `UploadDocument` → FLVE `/ekyc/frame` (wraps image as data URL)
- Sending `VerifyKYC` → FLVE `/ekyc/complete` (loses rich completion data)

### FLVE HuggingFace Space eKYC endpoints

Deployment: `LabaidAi-Retina/deployments/huggingface/`

| Endpoint | Method | Input | Output |
|---|---|---|---|
| `/ekyc/start` | POST JSON | `StartEKYCRequest` | `StartEKYCResponse` |
| `/ekyc/frame` | POST multipart | `session_id` + file upload | `EKYCFrameResponse` |
| `/ekyc/complete` | POST query | `session_id` | `CompleteEKYCResponse` |
| `/ekyc/status/{session_id}` | GET | path param | `GetEKYCStatusResponse` |

### FLVE StartEKYCRequest (already InsureTech-aligned)

```python
class StartEKYCRequest(BaseModel):
    user_id: str              # InsureTech user UUID (validated)
    tenant_id: Optional[str]
    user_type: Optional[str]
    portal: Optional[str]
    kyc_verification_id: Optional[str]  # InsureTech KYC record ID
    reference_image_url: Optional[str]  # for identity matching
    metadata: Optional[Dict[str, Any]]
```

### FLVE eKYC challenge sequence (4 steps)

| Step | Challenge | Detection logic | Thresholds |
|---|---|---|---|
| 1 | `BLINK` | Eye closed → reopened (EAR-based) | closed < 0.3, reopened > 0.7 |
| 2 | `LOOK_LEFT` | Head yaw | yaw < -25° |
| 3 | `LOOK_RIGHT` | Head yaw | yaw > 25° |
| 4 | `CAPTURE` | Frontal pose + liveness | \|yaw\| < 15°, \|pitch\| < 15°, liveness ≥ 0.7, confidence ≥ 0.8 |

Per-step timeout: 10s. Session timeout: 60s.

### FLVE CompleteEKYCResponse (rich output)

```python
class CompleteEKYCResponse(BaseModel):
    session_id: str
    success: bool
    state: str                           # "EKYC_SESSION_COMPLETED"
    profile_image_url: Optional[str]     # CDN URL of best captured frame
    profile_image_id: Optional[str]
    captured_image_base64: Optional[str] # base64-encoded best frame
    embedding: List[float]               # 512-dim face embedding
    liveness_confidence: float           # avg confidence across completed steps
    identity_match: bool                 # face match against reference (not yet computed)
    match_score: float                   # cosine similarity (not yet computed)
    summary: Optional[EKYCSessionSummary]
    completed_at: Optional[str]
    error: Optional[str]
```

### FLVE EKYCFrameResponse (per-frame feedback)

```python
class EKYCFrameResponse(BaseModel):
    session_id: str
    session_state: str
    current_step: Optional[EKYCStep]
    next_step: Optional[EKYCStep]
    step_completed: bool
    step_progress: float               # 0.0–1.0
    overall_progress: float            # 0.0–1.0
    detection: Optional[Dict]          # {detected, box}
    head_pose: Optional[HeadPose]      # {yaw, pitch, roll}
    eye_state: Optional[EyeState]      # {left_openness, right_openness, is_blinking, ...}
    eye_contours: Optional[Dict]
    liveness_score: float
    guidance: List[str]                # human-readable instructions
    error: Optional[str]
```

### ApiResponse\<T\> envelope (recent commit)

```yaml
ApiResponse<T>:
  type: object
  required: [success]
  properties:
    success:
      type: boolean
    data:
      description: Present on success. MAY be null for no-data ops.
      nullable: true
    error:
      $ref: '#/components/schemas/Error'
      nullable: true
    meta:
      $ref: '#/components/schemas/ResponseMeta'
      nullable: true

ResponseMeta:
  type: object
  properties:
    request_id: { type: string }
    pagination: { $ref: PaginationMeta, nullable: true }
    timestamp: { type: string, format: date-time }
    api_version: { type: string }
```

---

## Key Gaps

### 1. Public API response messages are too small for FLVE output

| InsureTech message | Missing FLVE data |
|---|---|
| `InitiateKYCResponse` (kyc_id, status, message) | provider, provider_reference, session_state, steps, timeout |
| `SubmitKYCFrameResponse` (accepted, guidance, current_step, completed_steps, total_steps, liveness_confidence) | session_state, step_progress, overall_progress, detection, head_pose, eye_state, step type/instruction |
| `CompleteKYCSessionResponse` (kyc_id, status, success, liveness_confidence, profile_image_url, message) | provider_reference, embedding, captured_image, identity_match, match_score, summary, completed_at |

### 2. Provider ID leakage

`flveExternalKYCClient.StartKYCVerification` returns `session_id` as `kyc_verification_id`, coupling InsureTech IDs to the provider. Correct model:

- `kyc_verifications.id` = InsureTech UUID (canonical)
- `kyc_verifications.provider_reference` = FLVE `session_id`

### 3. Gateway route binding incomplete

The authn KYC handlers do not hydrate `user_id` or `kyc_id` from path parameters. Route shape and request body expectations are misaligned.

### 4. Premature auto-verification

`CompleteKYCSession` currently marks `kyc_verified=true` immediately. For insurance/regulatory KYC, the default should be `PENDING_REVIEW` unless straight-through rules are explicitly approved.

### 5. FLVE TS SDK is stale

`sdks/flve-ts-sdk/src/types.ts` (generated 2026-01-19) is behind the proto/HF deployment:

| SDK type | Issue |
|---|---|
| `StartEKYCRequest` | Still `{userId, challenges}` — missing `tenant_id`, `user_type`, `portal`, `kyc_verification_id`, `reference_image_url` |
| `StartEKYCResponse` | Has `expiresAt: Date` — should be `total_timeout_seconds: int`, `state: string` |
| `SubmitEKYCFrameResponse` | Only `{currentStep, allSteps, completed}` — missing detection, head_pose, eye_state, guidance, progress |
| `CompleteEKYCResponse` | Only `{success, finalScore, capturedImageUrl, stepResults}` — missing embedding, identity_match, match_score, summary |

**Do not make InsureTech depend on this SDK until regenerated.** InsureTech should talk to FLVE via its own adapter using hand-typed Go structs mapped from the authoritative `ekyc_schemas_proto.py` models.

### 6. FLVE identity matching not yet implemented

The FLVE proto and Pydantic models support `identity_match` and `match_score` in `CompleteEKYCResponse`. The HuggingFace `complete_ekyc` route populates them as `false` / `0.0`. The ArcFace embedding comparison logic exists in the Go backend (`internal/server/handlers.go`) but is not wired into the Python eKYC flow.

---

## Canonical Design Decisions

### Public contract

InsureTech authn KYC endpoints are the **only** frontend/mobile/BFF contract. Raw FLVE responses are never exposed to clients.

### Provider contract

FLVE Space is consumed **only** behind the AuthN service through a dedicated Go HTTP adapter (`FLVEAdapter`).

### Persistence contract

InsureTech tables are the source of truth:

| Table | Role |
|---|---|
| `authn_schema.kyc_verifications` | KYC record lifecycle, provider reference, raw result |
| `authn_schema.user_profiles` | kyc_verified flag, profile photo URL, selfie URL |
| `authn_schema.document_verifications` | optionally link NID/passport docs |

### Response contract

Every HTTP response uses `ApiResponse<T>`:

| Layer | Content |
|---|---|
| Outer envelope | `ApiResponse<T>` with `success`, `data`, `error`, `meta` |
| `data` | InsureTech proto-shaped payload |
| Raw FLVE payload | Stored in `verification_result` JSONB column, not in the response |

---

## Detailed Contract Mapping

### 1. Initiate KYC

**InsureTech request:** `InitiateKYCRequest { user_id }`

**Adapter builds FLVE `/ekyc/start` request:**

| FLVE field | Source |
|---|---|
| `user_id` | `req.user_id` (path param) |
| `tenant_id` | Auth context / JWT claim |
| `user_type` | Auth context (B2C_CUSTOMER, AGENT, etc.) |
| `portal` | Auth context (customer_portal, b2b_portal, etc.) |
| `kyc_verification_id` | Newly created local `kyc_verifications.id` |
| `reference_image_url` | `user_profiles.id_upload_front_url` or `photograph_selfie_url` (trusted source only, never from client input) |

**Persistence on start:**

```
1. INSERT kyc_verifications (
     id = new UUID,
     type = KYC,
     entity_type = USER,
     entity_id = user_id,
     method = FLVE_EKYC,
     provider = "FLVE_HF",
     status = IN_PROGRESS
   )
2. Call FLVE /ekyc/start
3. UPDATE kyc_verifications SET
     provider_reference = flve.session_id,
     verification_result = json(flve_start_response)
```

**InsureTech response (extended):**

```protobuf
message InitiateKYCResponse {
  string kyc_id = 1;                       // InsureTech UUID
  string status = 2;                       // "IN_PROGRESS"
  string message = 3;
  string provider = 4;                     // "FLVE_HF"  [NEW]
  string provider_reference = 5;           // FLVE session_id  [NEW]
  string session_state = 6;                // "EKYC_SESSION_ACTIVE"  [NEW]
  repeated KYCStep steps = 7;              // challenge steps  [NEW]
  int32 total_timeout_seconds = 8;         // session timeout  [NEW]
  insuretech.common.v1.Error error = 100;
}
```

### 2. Submit KYC Frame

**InsureTech request:** `SubmitKYCFrameRequest { user_id, session_id, image_data, frame_sequence }`

**Adapter behavior:**

```
1. Resolve local kyc_id from session_id (or session_id IS the kyc_id)
2. ensureKYCSessionOwner(ctx, kyc_id, user_id) — Redis fast path then DB fallback
3. Look up provider_reference from kyc_verifications
4. POST multipart to FLVE /ekyc/frame?session_id={provider_reference}
5. Map FLVE EKYCFrameResponse into InsureTech SubmitKYCFrameResponse
6. Persist progress snapshot into verification_result JSONB
```

**InsureTech response (extended):**

```protobuf
message SubmitKYCFrameResponse {
  bool accepted = 1;
  string guidance = 2;                     // primary guidance message
  string current_step = 3;                 // challenge type name
  int32 completed_steps = 4;
  int32 total_steps = 5;
  double liveness_confidence = 6;
  string message = 7;
  string session_state = 8;               // "EKYC_SESSION_ACTIVE"  [NEW]
  double step_progress = 9;               // 0.0–1.0  [NEW]
  double overall_progress = 10;           // 0.0–1.0  [NEW]
  repeated string guidance_messages = 11; // all guidance strings  [NEW]
  KYCDetection detection = 12;            // face detection box  [NEW]
  KYCHeadPose head_pose = 13;             // yaw/pitch/roll  [NEW]
  KYCEyeState eye_state = 14;             // blink state  [NEW]
  KYCStep current_step_detail = 15;       // full step info  [NEW]
  insuretech.common.v1.Error error = 100;
}

// Supporting messages (InsureTech-native, not imported from FLVE proto)
message KYCStep {
  int32 step_number = 1;
  string challenge_type = 2;              // "BLINK", "LOOK_LEFT", "LOOK_RIGHT", "CAPTURE"
  string state = 3;                       // "PENDING", "IN_PROGRESS", "COMPLETED", "FAILED"
  string instruction = 4;                 // "Please blink your eyes"
  string instruction_key = 5;             // "ekyc.instruction.blink" (i18n)
  int32 timeout_seconds = 6;
  double confidence = 7;
}

message KYCDetection {
  bool detected = 1;
  int32 x = 2;
  int32 y = 3;
  int32 width = 4;
  int32 height = 5;
}

message KYCHeadPose {
  double yaw = 1;
  double pitch = 2;
  double roll = 3;
}

message KYCEyeState {
  double left_openness = 1;
  double right_openness = 2;
  bool is_blinking = 3;
}
```

### 3. Complete KYC Session

**Adapter behavior:**

```
1. Validate session ownership
2. POST to FLVE /ekyc/complete?session_id={provider_reference}
3. Persist into kyc_verifications:
   ─ verification_result = full FLVE CompleteEKYCResponse JSON
   ─ status = PENDING_REVIEW (default) or VERIFIED (if straight-through approved)
4. Update user_profiles:
   ─ profile_photo_url = flve.profile_image_url
   ─ photograph_selfie_url = flve.profile_image_url
5. Do NOT set kyc_verified=true here (wait for approve step)
```

**InsureTech response (extended):**

```protobuf
message CompleteKYCSessionResponse {
  string kyc_id = 1;                      // InsureTech UUID
  string status = 2;                      // "PENDING_REVIEW"
  bool success = 3;
  double liveness_confidence = 4;
  string profile_image_url = 5;           // CDN URL
  string message = 6;
  string provider_reference = 7;          // FLVE session_id  [NEW]
  string session_state = 8;               // "EKYC_SESSION_COMPLETED"  [NEW]
  bool identity_match = 9;                // face match result  [NEW]
  double match_score = 10;                // cosine similarity  [NEW]
  KYCSessionSummary summary = 11;         // step-by-step results  [NEW]
  string completed_at = 12;              // ISO timestamp  [NEW]
  insuretech.common.v1.Error error = 100;
}

message KYCSessionSummary {
  int32 total_steps = 1;
  int32 completed_steps = 2;
  int32 failed_steps = 3;
  int32 total_frames_processed = 4;
  int32 elapsed_ms = 5;
  repeated KYCStepResult step_results = 6;
}

message KYCStepResult {
  string challenge_type = 1;
  string state = 2;
  double confidence = 3;
  int32 frames_processed = 4;
  int32 elapsed_ms = 5;
}
```

### 4. Get KYC Status

**InsureTech response (extended):**

```protobuf
message GetKYCStatusResponse {
  string kyc_id = 1;
  string status = 2;
  string rejection_reason = 3;
  google.protobuf.Timestamp submitted_at = 4;
  google.protobuf.Timestamp reviewed_at = 5;
  string provider = 6;                    // "FLVE_HF"  [NEW]
  string provider_reference = 7;          // FLVE session_id  [NEW]
  string session_state = 8;               // current FLVE state  [NEW]
  KYCStep current_step = 9;               // current challenge  [NEW]
  int32 completed_steps = 10;             // steps done  [NEW]
  int32 total_steps = 11;                 // total challenges  [NEW]
  double overall_progress = 12;           // 0.0–1.0  [NEW]
  int32 remaining_seconds = 13;           // time left  [NEW]
  insuretech.common.v1.Error error = 100;
}
```

### 5. Approve KYC

Stays fully in InsureTech. No FLVE interaction.

```
1. UPDATE kyc_verifications SET
     status = VERIFIED,
     verified_by = reviewer_id,
     verified_at = now()
2. UPDATE user_profiles SET
     kyc_verified = true,
     kyc_verified_at = now()
```

---

## Implementation Phases

### Phase 1: Fix gateway route binding (prerequisite)

**Files:**
- `backend/inscore/cmd/gateway/internal/handlers/authn_handler.go`

**Changes:**
- Hydrate `user_id` from `r.PathValue("user_id")` for all KYC endpoints
- Hydrate `kyc_id` from `r.PathValue("kyc_id")` for approve
- Support GET-without-body on KYC status
- Return correct HTTP status codes (201 on initiate, 200 on others)

### Phase 2: Extend InsureTech authn KYC proto messages

**Files:**
- `proto/insuretech/authn/services/v1/core.proto` — add new fields to response messages, add supporting messages (`KYCStep`, `KYCDetection`, `KYCHeadPose`, `KYCEyeState`, `KYCSessionSummary`, `KYCStepResult`)

**Rules:**
- Keep all existing field numbers for backward compatibility
- Add new fields starting from next available number
- Use InsureTech-native message names (not FLVE imports) — the public API should never reference FLVE types
- Add `VERIFICATION_METHOD_FLVE_EKYC` to `VerificationMethod` enum in `proto/insuretech/kyc/entity/v1/kyc_verification.proto`

### Phase 3: Build a proper FLVEAdapter

**New file:** `backend/inscore/microservices/authn/internal/service/flve_adapter.go`

Replace the forced-fit `ExternalKYCClient` interface with a purpose-built adapter:

```go
type FLVEAdapter interface {
    StartEKYC(ctx context.Context, req *FLVEStartRequest) (*FLVEStartResponse, error)
    SubmitEKYCFrame(ctx context.Context, sessionID string, imageData []byte) (*FLVEFrameResponse, error)
    CompleteEKYC(ctx context.Context, sessionID string) (*FLVECompleteResponse, error)
    GetEKYCStatus(ctx context.Context, sessionID string) (*FLVEStatusResponse, error)
}
```

Go struct types should be hand-mapped from the authoritative FLVE Pydantic models in `deployments/huggingface/src/models/ekyc_schemas_proto.py`:

```go
type FLVEStartRequest struct {
    UserID              string            `json:"user_id"`
    TenantID            string            `json:"tenant_id,omitempty"`
    UserType            string            `json:"user_type,omitempty"`
    Portal              string            `json:"portal,omitempty"`
    KYCVerificationID   string            `json:"kyc_verification_id,omitempty"`
    ReferenceImageURL   string            `json:"reference_image_url,omitempty"`
    Metadata            map[string]string `json:"metadata,omitempty"`
}

type FLVEStartResponse struct {
    SessionID           string      `json:"session_id"`
    Steps               []FLVEStep  `json:"steps"`
    TotalTimeoutSeconds int         `json:"total_timeout_seconds"`
    State               string      `json:"state"`
    Error               string      `json:"error,omitempty"`
}

type FLVEFrameResponse struct {
    SessionID       string         `json:"session_id"`
    SessionState    string         `json:"session_state"`
    CurrentStep     *FLVEStep      `json:"current_step,omitempty"`
    StepCompleted   bool           `json:"step_completed"`
    StepProgress    float64        `json:"step_progress"`
    OverallProgress float64        `json:"overall_progress"`
    Detection       map[string]any `json:"detection,omitempty"`
    HeadPose        *FLVEHeadPose  `json:"head_pose,omitempty"`
    EyeState        *FLVEEyeState  `json:"eye_state,omitempty"`
    LivenessScore   float64        `json:"liveness_score"`
    Guidance        []string       `json:"guidance"`
    Error           string         `json:"error,omitempty"`
}

type FLVECompleteResponse struct {
    SessionID          string              `json:"session_id"`
    Success            bool                `json:"success"`
    State              string              `json:"state"`
    ProfileImageURL    string              `json:"profile_image_url,omitempty"`
    ProfileImageID     string              `json:"profile_image_id,omitempty"`
    CapturedImageB64   string              `json:"captured_image_base64,omitempty"`
    Embedding          []float64           `json:"embedding"`
    LivenessConfidence float64             `json:"liveness_confidence"`
    IdentityMatch      bool                `json:"identity_match"`
    MatchScore         float64             `json:"match_score"`
    Summary            *FLVESessionSummary `json:"summary,omitempty"`
    CompletedAt        string              `json:"completed_at,omitempty"`
    Error              string              `json:"error,omitempty"`
}
```

**Adapter requirements:**
- HTTP client with configurable timeout (default 30s to handle Space cold starts)
- Bearer token auth via `Authorization` header
- Multipart file upload for `/ekyc/frame`
- Retry with exponential backoff (3 attempts, 2s base) for 5xx / connection errors
- Normalize FLVE errors into `status.Error` that the gateway can wrap in `ApiResponse.error`

### Phase 4: Separate local and provider IDs

**Implementation rule:**

| ID | Scope |
|---|---|
| `kyc_verifications.id` | InsureTech UUID — canonical, exposed to clients |
| `kyc_verifications.provider_reference` | FLVE `session_id` — internal only |

**Redis session mapping:**

```
kyc:session:owner:{kyc_id}     → user_id     (TTL: 30 min)
kyc:session:provider:{kyc_id}  → flve_session_id  (TTL: 30 min)  [NEW]
kyc:session:frames:{kyc_id}    → frame_count  (TTL: 30 min)
```

**Migration:** Tests that assert `kyc_id == external session_id` must be updated.

### Phase 5: Refactor AuthN KYC orchestration

**File:** `backend/inscore/microservices/authn/internal/service/kyc_orchestrator_service.go`

**Target flow:**

```
InitiateKYC(user_id):
  1. Create kyc_verifications row (id = new UUID, provider = "FLVE_HF", status = IN_PROGRESS)
  2. Build FLVEStartRequest with user context from JWT + kyc_verification_id
  3. Call adapter.StartEKYC()
  4. Store provider_reference = flve.session_id
  5. Store raw response in verification_result
  6. Cache session mapping in Redis
  7. Return extended InitiateKYCResponse

SubmitKYCFrame(user_id, kyc_id, image_data):
  1. ensureKYCSessionOwner(kyc_id, user_id)
  2. Look up provider_reference from Redis (fast) or DB (fallback)
  3. Call adapter.SubmitEKYCFrame(provider_reference, image_data)
  4. Map FLVE response → InsureTech SubmitKYCFrameResponse
  5. Persist progress snapshot in verification_result
  6. Increment frame counter in Redis

CompleteKYCSession(user_id, kyc_id):
  1. ensureKYCSessionOwner(kyc_id, user_id)
  2. Look up provider_reference
  3. Call adapter.CompleteEKYC(provider_reference)
  4. Store full response in verification_result
  5. Update user_profiles.profile_photo_url
  6. Set kyc_verifications.status = PENDING_REVIEW (default)
  7. Return extended CompleteKYCSessionResponse

GetKYCStatus(user_id):
  1. Load kyc_verifications by entity_id = user_id
  2. If provider_reference exists and status = IN_PROGRESS:
     call adapter.GetEKYCStatus(provider_reference) for live data
  3. Merge persisted + live data into GetKYCStatusResponse

ApproveKYC(kyc_id, reviewer_id):
  1. No FLVE interaction
  2. Set kyc_verifications.status = VERIFIED, verified_by, verified_at
  3. Set user_profiles.kyc_verified = true, kyc_verified_at
```

### Phase 6: Update OpenAPI / API specs

**Files:**
- `api/paths/insuretech/authn/services/v1/AuthService.yaml` — update KYC endpoint response schemas
- `api/schemas/insuretech/authn/services/v1/` — add new response schemas with extended fields

All responses must continue using `ApiResponse<T>`:

```json
{
  "success": true,
  "data": {
    "kyc_id": "...",
    "status": "IN_PROGRESS",
    "provider": "FLVE_HF",
    "provider_reference": "...",
    "session_state": "EKYC_SESSION_ACTIVE",
    "steps": [...],
    "total_timeout_seconds": 60
  },
  "error": null,
  "meta": {
    "request_id": "...",
    "timestamp": "...",
    "api_version": "v1"
  }
}
```

### Phase 7: Regenerate SDKs

**After proto changes are settled:**

| SDK | Action |
|---|---|
| InsureTech Go (gen/go) | `buf generate` |
| InsureTech TS (gen/typescript) | `buf generate` |
| InsureTech Postman collection | Regenerate from OpenAPI |
| FLVE TS SDK (LabaidAi-Retina/sdks/flve-ts-sdk) | Regenerate from FLVE proto — but InsureTech does NOT depend on this |

---

## FLVE TS SDK → InsureTech TS SDK Mapping (post-regen)

For InsureTech web portals that consume the gateway, the mapping is:

| Portal calls InsureTech API | Gateway wraps in ApiResponse\<T\> | No direct FLVE SDK usage |
|---|---|---|
| `POST /v1/auth/users/{id}/kyc` | `{ success, data: InitiateKYCResponse }` | N/A |
| `POST .../kyc:submit-frame` | `{ success, data: SubmitKYCFrameResponse }` | N/A |
| `POST .../kyc:complete` | `{ success, data: CompleteKYCSessionResponse }` | N/A |

The FLVE TS SDK (`flve-ts-sdk`) is only relevant if a frontend needs to talk **directly** to the FLVE Space (e.g., for testing/demos). For production InsureTech, all FLVE communication goes through the Go backend.

---

## Rollout Order

| # | Task | Depends on | Status |
|---|---|---|---|
| 1 | Fix gateway route/body binding for authn KYC endpoints | — | ✅ Done |
| 2 | Add `VERIFICATION_METHOD_FLVE_EKYC` to KYC proto enum | — | ✅ Done |
| 3 | Extend InsureTech authn KYC proto messages (new fields + supporting messages) | #2 | ✅ Done |
| 4 | Build `FLVEAdapter` (replaces narrow `flveExternalKYCClient`) | #3 | ✅ Done |
| 5 | Add Redis key for `kyc:session:provider:{kyc_id}` mapping | #4 | ✅ Done |
| 6 | Refactor `InitiateKYC` to use adapter + separate IDs | #4, #5 | ✅ Done |
| 7 | Refactor `SubmitKYCFrame` to use adapter + rich response mapping | #6 | ✅ Done |
| 8 | Refactor `CompleteKYCSession` to use adapter + PENDING_REVIEW default | #6 | ✅ Done |
| 9 | Refactor `GetKYCStatus` to merge persisted + live FLVE data | #6 | ✅ Done |
| 10 | Update OpenAPI specs for extended responses | #3 | ✅ Done |
| 11 | Regenerate InsureTech SDKs (`buf generate`) | #3 | ✅ Done |
| 12 | Regenerate Postman collection | #10 | ⬜ Manual |
| 13 | Integration tests against FLVE Space | #6–#9 | ⬜ Manual |
| 14 | Feature flag rollout (`KYC_SERVICE_ENABLED=true`) | #13 | ⬜ Manual |

---

## Testing Checklist

### Unit tests

- [ ] FLVEAdapter request/response mapping (start, frame, complete, status)
- [ ] Provider/session ID separation (`kyc_id` ≠ `flve_session_id`)
- [ ] Path param hydration in gateway KYC handlers
- [ ] `ApiResponse<T>` success/error wrapping for all KYC responses
- [ ] FLVE error normalization (timeout, 5xx, invalid session)

### Integration tests

- [ ] Full flow: initiate → frame × N → complete → approve
- [ ] Invalid owner cannot submit frame (ownership check)
- [ ] Expired FLVE session returns correct error
- [ ] HuggingFace Space cold start handling (retry/timeout)
- [ ] Manual review path (complete → PENDING_REVIEW → approve)
- [ ] Straight-through path (if policy allows)

### Persistence tests

- [ ] `provider_reference` stored correctly on initiate
- [ ] `verification_result` stores raw FLVE payload (start, frame snapshots, complete)
- [ ] `user_profiles.profile_photo_url` updated on complete
- [ ] `user_profiles.kyc_verified` NOT set on complete, only on approve
- [ ] Redis session mapping created and expired correctly

### Security tests

- [ ] FLVE bearer token not leaked in client responses or logs
- [ ] `reference_image_url` only sourced from trusted InsureTech storage, never from client input
- [ ] Captured image base64 and embedding stored encrypted in `verification_result`
- [ ] Cross-user session access blocked

---

## Risk Register

| Risk | Mitigation |
|---|---|
| FLVE Space cold start (30-60s on first request) | Adapter retry with 30s timeout; consider keeping Space warm via health check cron |
| FLVE `identity_match`/`match_score` not computed | Phase-gate: initially ignore these fields; wire ArcFace comparison into HF complete route later |
| Proto field additions break existing clients | All new fields are additive with higher field numbers; protobuf wire format is backward compatible |
| FLVE session timeout (60s) vs InsureTech Redis TTL (30min) | InsureTech TTL covers the orchestration lifetime; FLVE timeout covers the live challenge window. Both are needed. |
| CDN image ownership | FLVE uploads to `labaidai/profile/{user_id}/` — ensure InsureTech trusts this path prefix |

---

## File Reference

### InsureTech (canonical)

| File | Purpose |
|---|---|
| `proto/insuretech/authn/services/v1/auth_service.proto` | KYC RPC definitions |
| `proto/insuretech/authn/services/v1/core.proto` | KYC request/response messages |
| `proto/insuretech/authn/entity/v1/user.proto` | User entity |
| `proto/insuretech/authn/entity/v1/user_profile.proto` | User profile (KYC flags) |
| `proto/insuretech/kyc/entity/v1/kyc_verification.proto` | KYC verification entity |
| `proto/insuretech/kyc/entity/v1/document_verification.proto` | Document verification entity |
| `proto/insuretech/kyc/services/v1/kyc_service.proto` | KYC service messages |
| `backend/inscore/microservices/authn/internal/service/kyc_orchestrator_service.go` | KYC orchestation logic |
| `backend/inscore/microservices/authn/internal/service/kyc_external_client.go` | ExternalKYCClient interface |
| `backend/inscore/microservices/authn/internal/service/flve_external_client.go` | Current FLVE adapter (to be replaced) |
| `backend/inscore/cmd/gateway/internal/handlers/kyc_handler.go` | Gateway KYC handler |
| `rules/01-response-envelope.md` | ApiResponse\<T\> spec |
| `api/paths/insuretech/authn/services/v1/AuthService.yaml` | OpenAPI KYC paths |

### LabaidAi-Retina (provider)

| File | Purpose |
|---|---|
| `deployments/huggingface/src/routes/ekyc.py` | HF Space eKYC endpoints |
| `deployments/huggingface/src/models/ekyc_schemas_proto.py` | Authoritative Pydantic models (source of truth for Go adapter structs) |
| `deployments/huggingface/src/ekyc/state_machine.py` | Challenge detection logic |
| `deployments/huggingface/src/ekyc/session.py` | Session manager |
| `deployments/huggingface/src/storage.py` | CDN upload |
| `deployments/huggingface/src/core/engine.py` | YOLO + ArcFace + liveness engine |
| `proto/flve/v1/service.proto` | FLVE gRPC service definition |
| `proto/flve/identity/entity/v1/user.proto` | UserContext (InsureTech-aligned) |
| `sdks/flve-ts-sdk/src/types.ts` | TS SDK types (STALE — do not use until regenerated) |
- `user_profiles.kyc_verified` updated only on final approval if manual review is required

## Immediate Implementation Recommendation

If the goal is the fastest safe path, do this first:

1. Keep the public KYC API in `AuthService`.
2. Create a new FLVE-specific adapter instead of reusing the generic external KYC interface.
3. Keep `kyc_verifications.id` as the canonical local ID.
4. Store FLVE `session_id` in `provider_reference`.
5. Persist the raw FLVE complete payload in `verification_result`.
6. Do not auto-set `kyc_verified=true` on FLVE completion until the business confirms straight-through approval rules.

That gives you FLVE-backed KYC without letting the provider dictate InsureTech's public API or persistence model.
