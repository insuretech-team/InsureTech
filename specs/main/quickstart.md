# Quickstart: AuthN / AuthZ / B2B Completion

**Feature**: AuthN/AuthZ/B2B Completion  
**Branch**: main  
**Last Updated**: 2026-03-12

---

## Prerequisites

| Tool | Required Version | Check |
|------|-----------------|-------|
| Go | 1.22+ | `go version` |
| Buf CLI | 2.x | `buf --version` |
| Docker Desktop | 4.x+ | `docker version` |
| PowerShell | 7+ | `$PSVersionTable` |
| PostgreSQL client | 17 (psql) | optional, for inspection |

---

## Local Setup

```powershell
# Clone / pull latest
cd e:\Projects\InsureTech

# Install Go dependencies
go work sync

# Verify proto generation pipeline works
.\scripts\generate.ps1
```

---

## Step-by-Step: Apply Changes

### 1. Commit proto changes first (Principle I gate)

After editing any `.proto` file:
```powershell
buf lint
buf breaking --against '.git#branch=main'
.\scripts\generate.ps1
```

The generation pipeline produces:
- `gen/go/` — Go protobuf + gRPC stubs
- `gen/ts/` — TypeScript stubs
- `gen/csharp/` — C# stubs
- `api/openapi.yaml` — updated OpenAPI spec

### 2. Apply DB migrations

```powershell
# Dry run first (validates no forbidden keywords)
.\run_migration.ps1 --Target=primary --dry-run --strict

# Apply
.\run_migration.ps1 --Target=primary
```

### 3. Apply seed data

```powershell
# Seeds only (runs INSERT … ON CONFLICT files in db/seeds/)
.\run_migration.ps1 --Target=primary --seeds-only
```

### 4. Start services (docker-compose)

```powershell
docker compose up -d postgres redis kafka
docker compose up -d authn authz b2b
```

Or individually:
```powershell
cd backend/inscore/microservices/authn && go run ./cmd/server
cd backend/inscore/microservices/authz && go run ./cmd/server  
cd backend/inscore/microservices/b2b  && go run ./cmd/server
```

---

## Running Tests

### Unit tests (no Docker required)

```powershell
# AuthZ
go test ./backend/inscore/microservices/authz/... -short -v

# AuthN
go test ./backend/inscore/microservices/authn/... -short -v

# B2B
go test ./backend/inscore/microservices/b2b/... -short -v
```

### Integration tests (requires Docker for testcontainers)

```powershell
# Each service spins up testcontainers automatically
go test ./backend/inscore/microservices/authz/... -run Integration -count=1 -timeout 120s
go test ./backend/inscore/microservices/authn/... -run Integration -count=1 -timeout 120s
go test ./backend/inscore/microservices/b2b/...   -run Integration -count=1 -timeout 120s
```

### Live DB tests (requires running local DB)

These are marked with `testing.Short()` skip in CI:
```powershell
go test ./backend/inscore/microservices/authz/internal/repository/... -v -count=1
```

---

## Smoke Tests (after services are running)

### AuthZ — verify JWKS is DB-backed

```powershell
# Should return KID from token_configs table
curl http://localhost:8080/.well-known/jwks.json
```

### AuthN — verify GetMe

```powershell
# Login first
$resp = Invoke-RestMethod -Uri http://localhost:8080/v1/auth/login `
  -Method POST -ContentType "application/json" `
  -Body '{"mobile":"+8801XXXXXXXXX","password":"Test@1234"}'

$token = $resp.access_token

# GetMe
Invoke-RestMethod -Uri http://localhost:8080/v1/auth/me `
  -Headers @{ Authorization = "Bearer $token" }
```

Expected response:
```json
{
  "user_id": "...",
  "user_type": "B2B_ORG_ADMIN",
  "portal": "b2b",
  "tenant_id": "...",
  "roles": ["partner_admin"],
  "permissions": [...]
}
```

### B2B — purchase order lifecycle

```powershell
# Create PO
$po = Invoke-RestMethod -Uri http://localhost:8080/v1/b2b/purchase-orders `
  -Method POST -ContentType "application/json" `
  -Headers @{ Authorization = "Bearer $token" } `
  -Body @{
    department_id  = "<dept_id>"
    plan_id        = "<plan_id>"
    employee_count = 10
    coverage_amount = @{ amount = 500000; currency = "BDT"; decimal_amount = 5000 }
  } | ConvertTo-Json

$poId = ($po | ConvertFrom-Json).purchase_order.purchase_order_id

# Approve PO (admin token required)
Invoke-RestMethod -Uri "http://localhost:8080/v1/b2b/purchase-orders/$poId`:approve" `
  -Method POST -ContentType "application/json" `
  -Headers @{ Authorization = "Bearer $adminToken" } `
  -Body '{"approved_by":"<admin_user_id>","approver_notes":"Looks good"}'
```

### B2B — check Prometheus metrics

```powershell
curl http://localhost:2112/metrics | Select-String "grpc_server_handled_total"
```

---

## Troubleshooting

| Symptom | Likely Cause | Fix |
|---------|-------------|-----|
| `buf lint` fails on new proto | Wrong field number or missing import | Check proto numbering; run `buf lint --path proto/insuretech/<domain>` |
| `generate.ps1` fails on inject_gorm_tags | New field doesn't have json tag | Add `json:"field_name,omitempty"` to proto field comment or regenerate |
| Migration dry-run fails with "forbidden keyword" | SQL file has `CREATE TABLE` or `INSERT` | Move to seeds/ or rewrite as ALTER |
| `GetMe` returns `NotFound` | user_profiles not seeded for test user | Run `go run ./internal/seeder/...` or use seeder user |
| `GetJWKS` returns empty keys | `token_configs` table empty and no env var | Run `portal_seeder.SeedTokenConfig()` manually in test setup |
| B2B metrics endpoint not available | Port 2112 not mapped in docker-compose | Add `2112:2112` to b2b service in docker-compose.yml |
