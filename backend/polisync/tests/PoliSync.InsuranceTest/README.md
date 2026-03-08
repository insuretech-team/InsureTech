# Insurance Service CRUD Test

This test application validates the Insurance Service gRPC API by testing CRUD operations on the Products table.

## Overview

The Insurance Service is a Go microservice that provides CRUD operations for all 30 tables in the `insurance_schema`. This C# test application connects to the service via gRPC and validates that the basic operations work correctly.

## What's Tested

Currently, the test validates:
- **CREATE**: Creating a new product with all required fields
- **READ**: Retrieving a product by ID
- **UPDATE**: Modifying product fields
- **LIST**: Listing products with pagination
- **DELETE**: Soft-deleting a product

## Prerequisites

1. **Go Insurance Service** must be running on `localhost:50115`
2. **PostgreSQL database** must be accessible with the insurance_schema
3. **Proto files** must be generated (`buf generate` from project root)

## Running the Tests

### 1. Start the Insurance Service

```powershell
cd backend/inscore/microservices/insurance
go run ./cmd/server
```

The service should start on port 50115.

### 2. Run the Test

```powershell
cd backend/polisync/tests/PoliSync.InsuranceTest
dotnet run
```

## Expected Output

```
==============================================
Insurance Service CRUD Test
==============================================

Connecting to Insurance Service at: http://localhost:50115

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
SIMPLE CRUD TEST - Products Table
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Testing Product CRUD Operations...

  ✓ CREATE Product
  ✓ READ Product (GetByID)
  ✓ UPDATE Product
  ✓ LIST Products
    Found 1 products (Total: 1)
  ✓ DELETE Product

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
TEST SUMMARY
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Total Tests: 5
Passed: 5
Failed: 0
Success Rate: 100.0%

✅ All tests passed successfully!
```

## Configuration

The service URL can be configured in `appsettings.json`:

```json
{
  "InsuranceService": {
    "Url": "http://localhost:50115"
  }
}
```

## Troubleshooting

### Connection Refused
- Ensure the Insurance Service is running on port 50115
- Check firewall settings

### NotFound Errors
- Ensure the database has the insurance_schema
- Check that migrations have been run

### Invalid Argument Errors
- Verify that the `created_by` user UUID exists in `authn_schema.users`
- The test uses `00000000-0000-0000-0000-000000000001` as the system user

## Next Steps

To test additional tables:
1. Check the actual proto entity definitions in `proto/insuretech/*/entity/v1/`
2. Add test methods to `SimpleTestRunner.cs` following the Product pattern
3. Ensure all required fields match the proto definitions

## Architecture

```
┌─────────────────────┐
│  C# Test App        │
│  (This Project)     │
└──────────┬──────────┘
           │ gRPC
           ▼
┌─────────────────────┐
│  Go Insurance       │
│  Service            │
│  (Port 50115)       │
└──────────┬──────────┘
           │ SQL
           ▼
┌─────────────────────┐
│  PostgreSQL         │
│  insurance_schema   │
└─────────────────────┘
```

## Files

- `Program.cs` - Test entry point
- `SimpleTestRunner.cs` - Test implementation
- `appsettings.json` - Configuration
- `PoliSync.InsuranceTest.csproj` - Project file
