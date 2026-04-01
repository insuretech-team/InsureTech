# Insurance Engine - Synchronization & Architecture Guide

This guide ensures that the Insurance Engine microservice remains production-ready, synchronized with gRPC contracts, and adheres to the Vertical Slice Architecture (VSA).

## 🏛️ Architectural Principles

1.  **Proto-First Development**:
    *   The `.proto` files in the `proto/` root are the **absolute source of truth**.
    *   Never modify C# DTOs or generated types manually. Always update the `.proto` and let the compiler regenerate the types.
2.  **Vertical Slice Architecture (VSA)**:
    *   Features are organized by domain slices (e.g., `InsuranceEngine.Policy`, `InsuranceEngine.Commission`).
    *   Each slice contains its own Handlers, Commands, and Queries (VSA).
    *   Dependencies between slices should be minimized and handled via the `SharedKernel`.

## 🔄 Synchronization Rules

### 1. MediatR Handlers vs. Proto Contracts
When a gRPC contract changes, update the corresponding `Command` or `Query` in the domain slice to match the **Request** and **Response** types exactly.

*Example:*
```csharp
// Handler response must match Proto-generated type
public class CreatePolicyHandler : IRequestHandler<CreatePolicyCommand, CreatePolicyResponse> { ... }
```

### 2. ApiHost Controllers
API Controllers (`InsuranceEngine.ApiHost`) should be lightweight wrappers that map REST requests to MediatR commands.

*   **Mapping Rules**:
    *   REST requests should map directly to Proto types where possible.
    *   Controller actions should return Proto-generated response types directly to ensure type safety.

### 3. Build Enforcement
Always run `dotnet build` after any change to ensure zero errors.
*   **Zero-Error Policy**: The `InsuranceEngine` solution MUST build with 0 errors to be considered production-ready.

## 🛠️ Maintenance Workflow

1.  **Modify Proto**: Update `.proto` files.
2.  **Regenerate Classes**: The build process handles this automatically.
3.  **Update Handlers**: Align MediatR handlers with new Proto signatures.
4.  **Update Controllers**: Align REST endpoints with MediatR changes.
5.  **Verify**: Run `dotnet build`.

---
*Created during the Insurance Engine Backend Stabilization Effort (March 2026).*
