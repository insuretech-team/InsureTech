# Phase 6 Implementation Checklist

## ✅ MODIFICATION 1: server.go
- [x] Added `DocGenServer` struct with three fields:
  - [x] `handler documentservicev1.DocumentServiceServer`
  - [x] `kafkaPublisher *kafka.Publisher`
  - [x] `genWorker *worker.GenerationWorker`
- [x] Imported required packages:
  - [x] `internal/config`
  - [x] `internal/kafka`
  - [x] `internal/worker`
  - [x] `pkg/logger`
  - [x] `context`
- [x] Updated `NewDocumentServer()` to:
  - [x] Load config via `config.Load()`
  - [x] Initialize Kafka publisher if brokers configured
  - [x] Log warning if Kafka not configured (continue without it)
  - [x] Inject publisher via `docService.SetKafkaPublisher()`
  - [x] Initialize worker if `AsyncGeneration` enabled
  - [x] Start worker in goroutine via `worker.Start(ctx)`
  - [x] Return `*DocGenServer` instead of handler directly
- [x] Added `Handler()` method
- [x] Added `Close()` method that:
  - [x] Stops genWorker if not nil
  - [x] Closes kafkaPublisher if not nil
  - [x] Handles cleanup errors

## ✅ MODIFICATION 2: internal/service/document_service.go
- [x] Added imports:
  - [x] `internal/kafka`
  - [x] `pkg/logger`
- [x] Added `kafkaPublisher` field to `DocumentService` struct
- [x] Added `SetKafkaPublisher()` method for dependency injection
- [x] **Template Creation** (`CreateTemplate`):
  - [x] After successful creation, publish `document.template.created`
  - [x] Non-blocking (goroutine)
  - [x] Nil guard for kafkaPublisher
  - [x] Log warning on error
- [x] **Template Update** (`UpdateTemplate`):
  - [x] After successful update, publish `document.template.updated`
  - [x] Non-blocking (goroutine)
  - [x] Nil guard for kafkaPublisher
  - [x] Log warning on error
- [x] **Document Generation Success** (`GenerateDocument`):
  - [x] After successful generation, publish `document.generated`
  - [x] Non-blocking (goroutine)
  - [x] Nil guard for kafkaPublisher
  - [x] Log warning on error
  - [x] Include: generationID, templateID, tenantID, entityID, entityType, format
- [x] **Document Generation Failure** (`GenerateDocument`):
  - [x] On template fetch error: publish `document.generation.failed`
  - [x] On document save error: publish `document.generation.failed`
  - [x] Non-blocking (goroutine)
  - [x] Nil guard for kafkaPublisher
  - [x] Log warning on error
  - [x] Include: generationID, tenantID, reason

## ✅ MODIFICATION 3: cmd/server/main.go (NEW FILE)
- [x] File created at correct path
- [x] Implemented signal handling:
  - [x] `os.Interrupt` handling
  - [x] `syscall.SIGTERM` handling
  - [x] `syscall.SIGINT` handling
- [x] Graceful shutdown sequence:
  - [x] Call `docgenServer.Close()` (cleanup Kafka & worker)
  - [x] Call `grpcServer.GracefulStop()` (wait for in-flight requests)
  - [x] Log shutdown completion
- [x] Database setup (placeholder)
- [x] gRPC server setup

## Code Quality Checks
- [x] All Kafka publishes are non-blocking
- [x] All Kafka publishes have nil guards
- [x] All errors are logged (not silently ignored)
- [x] Context is properly managed (new background context for async publishes)
- [x] No existing functionality rewritten
- [x] Minimal, surgical changes applied
- [x] All imports are correct and used
- [x] Error handling is appropriate

## Integration Points
- [x] `NewDocumentServer()` signature changed - callers must use `.Handler()`
- [x] `DocumentService` now requires Kafka injection (optional via nil)
- [x] Config loading is integrated
- [x] Worker pool is created and started
- [x] Graceful shutdown is implemented

## Files Modified
1. ✅ `server.go` - Core wiring
2. ✅ `internal/service/document_service.go` - Event publishing
3. ✅ `cmd/server/main.go` - Entry point with signal handling (NEW)

## Documentation
- ✅ `PHASE_6_WIRING_SUMMARY.md` - Complete documentation
- ✅ `IMPLEMENTATION_CHECKLIST.md` - This checklist

## Status: ✅ COMPLETE

All modifications for Phase 6 have been successfully implemented and wired together.
The system now:
- Publishes Kafka events for template and document operations
- Supports async document generation via worker pool
- Includes graceful shutdown handling
- Remains resilient if Kafka is unavailable
