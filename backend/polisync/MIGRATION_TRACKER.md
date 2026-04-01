# PoliSync Migration Tracker

## Target Rules

1. `proto/` is the source of truth.
2. `backend/polisync` uses proto-generated contracts and business orchestration.
3. Database CRUD belongs to `backend/inscore/microservices/insurance`.
4. Public HTTP belongs to `backend/inscore/cmd/gateway`.
5. No compile-exclusion tricks for dead code. If a file is obsolete, delete it. If it is live, make it compile.

## Progress On 2026-03-27

### Completed

- Fixed the stale integration test fake so `backend/polisync/PoliSync.sln` builds again.
- Cleaned `src/PoliSync.Infrastructure` so it no longer depends on `Compile Remove` to stay buildable.

### `PoliSync.Infrastructure` cleanup details

- Restored `DependencyInjection.cs` to normal compilation.
- Restored `GrpcClients/*` to normal compilation.
- Updated the typed gRPC wrappers to match the current proto contracts.
- Deleted dead excluded EF configuration files from the old product-domain path.
- Deleted dead excluded repository files that represented direct DB access paths we do not want to revive.

### Verification

- `dotnet build backend/polisync/src/PoliSync.Infrastructure/PoliSync.Infrastructure.csproj`
- `dotnet build backend/polisync/PoliSync.sln`

Both pass after the cleanup.

## Remaining High-Risk Areas

### Direct DB access still present in live code

- `src/PoliSync.ApiHost/Program.cs`
- `src/PoliSync.Infrastructure/Persistence/*`
- `src/PoliSync.Policy/Application/Commands/IssuePolicyCommandHandler.cs`
- `src/PoliSync.Policy/Application/Queries/GetPolicyQueryHandler.cs`
- `src/PoliSync.Claims/Application/Commands/FileClaimCommandHandler.cs`

### HTTP still owned by PoliSync

- `src/PoliSync.ApiHost/Controllers/ProductsController.cs`
- `src/PoliSync.ApiHost/Controllers/OrdersController.cs`
- `src/PoliSync.ApiHost/Controllers/PoliciesController.cs`
- `src/PoliSync.ApiHost/Controllers/ClaimsController.cs`
- `src/PoliSync.ApiHost/Controllers/InsuranceProposalsController.cs`
- `src/PoliSync.ApiHost/Controllers/QuotationsController.cs`
- `src/PoliSync.ApiHost/Controllers/WorkflowController.cs`
- `src/PoliSync.ApiHost/Controllers/WorkflowTemplateAdminController.cs`
- `src/PoliSync.ApiHost/Controllers/WebhooksController.cs`

### In-memory or scaffold-heavy domains

- `src/PoliSync.Actuarial`
- `src/PoliSync.CRM`
- `src/PoliSync.LifeInsurance`
- `src/PoliSync.Quoting`
- `src/PoliSync.RulesEngine`
- `src/PoliSync.VehicleInsurance`

## Best Next Folders

1. `src/PoliSync.ApiHost`
   Remove 501-style placeholder endpoints and reduce HTTP ownership.
2. `src/PoliSync.Products`
   Finish plan/rider/pricing/product endpoint behavior and remove current product stubs.
3. `src/PoliSync.Policy`
   Replace raw SQL and direct insurance-schema access with Go `insurance` gRPC.
