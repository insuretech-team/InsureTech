# Shared `pkg` Architecture

`backend/inscore/pkg` is for cross-service mechanics that sit on hot paths and would otherwise be duplicated across microservices.

Good fits for `pkg`:
- transport helpers such as gRPC metadata extraction and request-context normalization
- runtime address normalization for Kafka, Redis, and service discovery
- delivery/provider clients that are infrastructure adapters rather than business workflows
- codecs and DTOs shared by multiple services, such as notification preference serialization
- reusable interceptors, retry helpers, idempotency helpers, and event envelope utilities

Bad fits for `pkg`:
- service-owned repositories that read or write a domain table as the source of truth
- orchestration that encodes business policy for one microservice
- cross-service shortcuts that bypass ownership boundaries

The intended pattern is:
1. The owning service keeps write responsibility and business invariants.
2. Other services reuse common transport/runtime/codec logic from `pkg`.
3. Read-heavy fanout paths can use local projections or read adapters, but not shared ownership of mutable domain state.

Current shared modules:
- `pkg/grpcmeta`: common metadata extraction and caller-context helpers
- `pkg/notifyprefs`: shared notification preference codec used by AuthN-backed user rows and notification dispatch
- `pkg/runtimeaddr`: runtime-safe Kafka and Redis address normalization across host, Docker, and WSL
