# Insurance Service Architecture

## Overview

The Insurance Service is a Go microservice that provides CRUD operations for all insurance_schema tables. This eliminates duplicate work by centralizing database operations in Go (which already has GORM working), while C# services call it via gRPC.

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    C# PoliSync Services                      │
│  (Product, Policy, Claims, Quote, Underwriting, etc.)       │
│                                                              │
│  Business Logic + Domain Models + gRPC APIs                 │
└──────────────────┬──────────────────────────────────────────┘
                   │ gRPC Calls
                   ▼
┌─────────────────────────────────────────────────────────────┐
│              Go Insurance Service (Port 50115)               │
│                                                              │
│  Generic CRUD for insurance_schema tables:                  │
│  - Products, ProductPlans, Riders, PricingConfigs           │
│  - Policies, Endorsements, Beneficiaries                    │
│  - Claims, Quotations, UnderwritingDecisions                │
│  - Renewals, Refunds, Fraud Detection                       │
└──────────────────┬──────────────────────────────────────────┘
                   │ GORM
                   ▼
┌─────────────────────────────────────────────────────────────┐
│              PostgreSQL (insurance_schema)                   │
│                                                              │
│  Schema managed by Go migrations                            │
│  Proto entities with GORM tags                              │
└─────────────────────────────────────────────────────────────┘
```

## Benefits

1. **No Duplicate Work**: Database mapping is done once in Go with GORM
2. **Single Source of Truth**: Go manages schema migrations and database access
3. **Proto-First**: Both Go and C# use the same proto definitions
4. **Separation of Concerns**: 
   - Go handles data persistence
   - C# handles business logic and domain rules

## Service Details

### Go Insurance Service

**Location**: `backend/inscore/microservices/insurance/`

**Port**: 50115 (gRPC)

**Responsibilities**:
- CRUD operations for all insurance_schema tables
- Direct GORM database access
- Proto message serialization/deserialization
- Basic validation

**Proto Definition**: `api/proto/insuretech/insurance/services/v1/insurance_service.proto`

### C# PoliSync Services

**Responsibilities**:
- Business logic and domain rules
- Complex workflows and orchestration
- Event publishing
- External API integration
- Calling Insurance Service for data operations

## Example Usage

### C# Calling Insurance Service

```csharp
// Create gRPC client
using var channel = GrpcChannel.ForAddress("http://localhost:50115");
var client = new InsuranceService.InsuranceServiceClient(channel);

// Create a product
var request = new CreateProductRequest
{
    Product = new Product
    {
        ProductId = Guid.NewGuid().ToString(),
        TenantId = "tenant-123",
        ProductCode = "LIFE-001",
        ProductName = "Term Life Insurance",
        // ... other fields
    }
};

var response = await client.CreateProductAsync(request);
```

## Running the Services

### 1. Start Insurance Service (Go)

```bash
cd backend/inscore/microservices/insurance
go run main.go
```

Or use the service manager:
```bash
cd backend/inscore
go run cmd/service-manager/main.go start insurance
```

### 2. Test from C#

```bash
cd backend/polisync/tests/PoliSync.DbTest
dotnet run
```

## Environment Variables

Add to `.env`:
```bash
INSURANCE_HOST=0.0.0.0
INSURANCE_GRPC_PORT=50115
INSURANCE_HTTP_PORT=50116
```

## Proto Generation

After modifying `insurance_service.proto`:

```bash
# Generate Go code
buf generate

# Generate C# code (already included in PoliSync.Proto project)
```

## Database Schema

The insurance_schema is managed by Go migrations in:
`backend/inscore/db/migrations/insurance_schema/`

All proto entities have GORM tags for database mapping.

## Next Steps

1. Generate proto files: `buf generate`
2. Build Go service: `cd backend/inscore/microservices/insurance && go build`
3. Build C# test: `cd backend/polisync/tests/PoliSync.DbTest && dotnet build`
4. Start Insurance Service
5. Run C# test to verify CRUD operations

## Files Created

### Go Service
- `backend/inscore/microservices/insurance/main.go`
- `backend/inscore/microservices/insurance/service/insurance_service.go`
- `backend/inscore/microservices/insurance/internal/repository/*.go`

### Proto
- `api/proto/insuretech/insurance/services/v1/insurance_service.proto`

### C# Client
- `backend/polisync/src/PoliSync.Infrastructure/Clients/InsuranceServiceClient.cs`
- `backend/polisync/tests/PoliSync.DbTest/Program.cs` (updated)

### Configuration
- `backend/inscore/configs/services.yaml` (updated)
- `.env` (updated)
