#!/usr/bin/env pwsh
# Test error handling with invalid data

Write-Host "Testing Error Handling..." -ForegroundColor Cyan
Write-Host ""

# Load .env
$envPath = "..\..\..\..\.env"
if (Test-Path $envPath) {
    Get-Content $envPath | ForEach-Object {
        if ($_ -match '^\s*([^#][^=]+)=(.*)$') {
            $key = $matches[1].Trim()
            $value = $matches[2].Trim().Trim('"').Trim("'")
            [Environment]::SetEnvironmentVariable($key, $value, "Process")
        }
    }
}

# Create a simple C# test inline
$testCode = @'
using Grpc.Net.Client;
using Insuretech.Insurance.Services.V1;
using Insuretech.Products.Entity.V1;
using Insuretech.Common.V1;
using Grpc.Core;

var channel = GrpcChannel.ForAddress("http://localhost:50115");
var client = new InsuranceService.InsuranceServiceClient(channel);

Console.WriteLine("Test 1: Invalid product_code format (should fail with InvalidArgument)");
try {
    var req = new CreateProductRequest {
        Product = new Product {
            ProductId = Guid.NewGuid().ToString(),
            ProductCode = "INVALID-FORMAT-123",  // Wrong format
            ProductName = "Test",
            Category = ProductCategory.Health,
            BasePremium = new Money { Amount = 100000, Currency = "BDT" },
            MinSumInsured = new Money { Amount = 1000000, Currency = "BDT" },
            MaxSumInsured = new Money { Amount = 5000000, Currency = "BDT" },
            MinTenureMonths = 12,
            MaxTenureMonths = 60,
            Status = ProductStatus.Active,
            CreatedBy = "ccca65ad-ae2c-4d42-8ccc-2122db78d617"
        }
    };
    await client.CreateProductAsync(req);
    Console.WriteLine("  ✗ FAILED: Should have thrown error");
} catch (RpcException ex) {
    Console.WriteLine($"  ✓ Got expected error: {ex.StatusCode} - {ex.Status.Detail}");
}

Console.WriteLine("\nTest 2: Duplicate product_code (should fail with AlreadyExists)");
try {
    // First create
    var productId = Guid.NewGuid().ToString();
    var req = new CreateProductRequest {
        Product = new Product {
            ProductId = productId,
            ProductCode = "DUP-999",
            ProductName = "Test",
            Category = ProductCategory.Health,
            BasePremium = new Money { Amount = 100000, Currency = "BDT" },
            MinSumInsured = new Money { Amount = 1000000, Currency = "BDT" },
            MaxSumInsured = new Money { Amount = 5000000, Currency = "BDT" },
            MinTenureMonths = 12,
            MaxTenureMonths = 60,
            Status = ProductStatus.Active,
            CreatedBy = "ccca65ad-ae2c-4d42-8ccc-2122db78d617"
        }
    };
    await client.CreateProductAsync(req);
    
    // Try duplicate
    req.Product.ProductId = Guid.NewGuid().ToString();
    await client.CreateProductAsync(req);
    Console.WriteLine("  ✗ FAILED: Should have thrown error");
    
    // Cleanup
    await client.DeleteProductAsync(new DeleteProductRequest { ProductId = productId });
} catch (RpcException ex) {
    Console.WriteLine($"  ✓ Got expected error: {ex.StatusCode} - {ex.Status.Detail}");
}

Console.WriteLine("\nError handling tests complete!");
'@

# Save and run
$testCode | Out-File -FilePath "temp-error-test.cs" -Encoding UTF8
dotnet-script temp-error-test.cs
Remove-Item temp-error-test.cs
