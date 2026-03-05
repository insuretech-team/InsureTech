# Phase 6: Kafka Publisher & Async Worker Integration - COMPLETED

## Overview
Successfully integrated the Kafka publisher and async generation worker into the DocGen microservice server.

## Files Modified

### 1. `server.go` - DocGenServer Structure & Initialization
**Changes:**
- Created new `DocGenServer` struct with three fields:
  - `handler documentservicev1.DocumentServiceServer` - gRPC handler
  - `kafkaPublisher *kafka.Publisher` - Kafka event publisher
  - `genWorker *worker.GenerationWorker` - Async generation worker pool

- Updated `NewDocumentServer()` function:
  - Now loads config via `config.Load()`
  - Initializes `kafka.Publisher` if `KafkaBrokers` are configured
  - Logs warning if Kafka is not configured, continues without it
  - Injects Kafka publisher into service via `docService.SetKafkaPublisher()`
  - Creates `GenerationWorker` if `AsyncGeneration` config is true
  - Starts worker goroutines via `genWorker.Start(ctx)`
  - Returns `*DocGenServer` instead of just the handler

- Added `Handler()` method:
  - Returns the underlying gRPC `DocumentServiceServer`

- Added `Close()` method:
  - Gracefully stops the `GenerationWorker` if initialized
  - Closes the `KafkaPublisher` if initialized
  - Handles cleanup errors properly

**Imports Added:**
```go
"github.com/newage-saint/insuretech/backend/inscore/microservices/docgen/internal/config"
"github.com/newage-saint/insuretech/backend/inscore/microservices/docgen/internal/kafka"
"github.com/newage-saint/insuretech/backend/inscore/microservices/docgen/internal/worker"
"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
"context"
```

---

### 2. `internal/service/document_service.go` - Kafka Event Publishing
**Changes:**
- Added `kafkaPublisher *kafka.Publisher` field to `DocumentService` struct

- Added `SetKafkaPublisher()` method:
  - Allows dependency injection of Kafka publisher after service creation
  - Nil-safe (publisher can be nil)

- **Template Creation Event** (`CreateTemplate`):
  - After successful template creation, publishes `document.template.created` event
  - Non-blocking goroutine with nil guard
  - Logs warnings on publish failures

- **Template Update Event** (`UpdateTemplate`):
  - After successful template update, publishes `document.template.updated` event
  - Non-blocking goroutine with nil guard
  - Logs warnings on publish failures

- **Document Generation Success Event** (`GenerateDocument`):
  - After successful document generation, publishes `document.generated` event
  - Non-blocking goroutine with nil guard
  - Includes generation ID, template ID, tenant ID, entity info, and format

- **Document Generation Failure Events** (`GenerateDocument`):
  - When template fetch fails: publishes `document.generation.failed` with error details
  - When document save fails: publishes `document.generation.failed` with error details
  - Non-blocking goroutines with nil guards
  - Includes failure reason in event payload

**Imports Added:**
```go
"github.com/newage-saint/insuretech/backend/inscore/microservices/docgen/internal/kafka"
"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
```

---

### 3. `cmd/server/main.go` - Entry Point with Signal Handling (NEW FILE)
**Purpose:** Server bootstrap with graceful shutdown

**Features:**
- Database initialization
- DocGen server creation with all dependencies
- gRPC server setup and listening
- Signal handling for `SIGTERM`, `SIGINT`, and `os.Interrupt`
- Graceful shutdown sequence:
  1. Closes DocGen server resources (calls `docgenServer.Close()`)
  2. Stops gRPC server gracefully
  3. Logs shutdown completion

**Cleanup on Shutdown:**
- When receiving shutdown signal:
  - Calls `docgenServer.Close()` to clean up Kafka publisher and async worker
  - Calls `grpcServer.GracefulStop()` to wait for in-flight requests

---

## Configuration Dependencies

The wiring relies on the following config fields from `internal/config/config.go`:
- `KafkaBrokers []string` - Kafka broker addresses
- `KafkaDocgenTopic string` - Topic for publishing events
- `AsyncGeneration bool` - Enable/disable async generation worker
- `AsyncWorkerCount int` - Number of worker goroutines
- `MaxGenerationTimeout time.Duration` - Generation timeout

---

## Kafka Events Published

### 1. `document.template.created`
- Triggered: After successful template creation
- Payload: `TemplateID`, `TenantID`, `TemplateName`, `EventType`, `Timestamp`

### 2. `document.template.updated`
- Triggered: After successful template update
- Payload: `TemplateID`, `TenantID`, `EventType`, `Timestamp`

### 3. `document.generated`
- Triggered: After successful document generation
- Payload: `GenerationID`, `TemplateID`, `TenantID`, `EntityID`, `EntityType`, `Format`, `EventType`, `Timestamp`

### 4. `document.generation.failed`
- Triggered: When template fetch fails OR document save fails
- Payload: `GenerationID`, `TenantID`, `Reason`, `EventType`, `Timestamp`

---

## Non-Blocking Publishing Pattern

All Kafka event publishing is implemented as **non-blocking fire-and-forget**:

```go
go func() {
    if s.kafkaPublisher != nil {
        if err := s.kafkaPublisher.PublishEventName(...); err != nil {
            logger.Warnf("failed to publish event: %v", err)
        }
    }
}()
```

Benefits:
- Document generation doesn't block waiting for Kafka acknowledgment
- Service remains responsive even if Kafka is slow or unavailable
- Failures are logged but don't propagate to caller
- Gracefully handles nil publisher (Kafka optional)

---

## Graceful Shutdown Flow

```
SIGTERM/SIGINT
     ↓
docgenServer.Close()
     ├── genWorker.Stop() → Waits for in-flight async jobs
     └── kafkaPublisher.Close() → Flushes pending messages
     ↓
grpcServer.GracefulStop()
     └── Waits for in-flight gRPC requests
     ↓
Application exits
```

---

## Testing Checklist

- [ ] Verify Kafka publisher connects to brokers
- [ ] Verify events are published after template creation
- [ ] Verify events are published after template update
- [ ] Verify events are published after document generation
- [ ] Verify failure events are published on generation errors
- [ ] Verify service continues without Kafka if brokers not configured
- [ ] Verify async worker processes generation requests
- [ ] Verify graceful shutdown closes resources properly
- [ ] Verify signal handling (SIGTERM, SIGINT)

---

## Notes

1. **TenantID in Template Operations**: Currently empty string in `CreateTemplate` and `UpdateTemplate` as tenant context is not available in current method signatures. Can be enhanced by passing tenant context through the API.

2. **Async Worker**: The generator function in the worker is currently a stub logging placeholder. Actual implementation would invoke `DocumentService.GenerateDocument()` asynchronously.

3. **Backwards Compatibility**: The return type of `NewDocumentServer()` changed from `(documentservicev1.DocumentServiceServer, error)` to `(*DocGenServer, error)`. Callers must use `docgenServer.Handler()` to get the gRPC handler.

4. **Error Handling**: All Kafka publish operations are non-blocking and failures are logged as warnings, not returned to caller. This prevents cascading failures.
