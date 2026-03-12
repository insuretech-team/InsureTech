using Microsoft.Extensions.DependencyInjection;
using InsuranceEngine.Products.Infrastructure.Persistence;
using Microsoft.EntityFrameworkCore;
using InsuranceEngine.Products.Domain;
using InsuranceEngine.Products.Domain.Enums;
using System;
using System.Threading.Tasks;

namespace InsuranceEngine.ApiHost.Persistence;

public static class DbInitializer
{
    public static async Task Initialize(IServiceProvider serviceProvider)
    {
        using var scope = serviceProvider.CreateScope();
        var productsContext = scope.ServiceProvider.GetRequiredService<ProductsDbContext>();

        try
        {
            if (await productsContext.Products.AnyAsync()) return;

            var tenantId = Guid.Parse("11111111-1111-1111-1111-111111111111");
            var createdBy = Guid.Parse("00000000-0000-0000-0000-000000000001");

            var products = new[]
            {
                new Product
                {
                    Id = Guid.NewGuid(),
                    ProductCode = "HLT-001",
                    ProductName = "Health Guard Plus",
                    ProductNameBn = "হেলথ গার্ড প্লাস",
                    Category = ProductCategory.Health,
                    Status = ProductStatus.Active,
                    BasePremiumAmount = 500000, // 5000 BDT in paisa
                    MinSumInsuredAmount = 10000000, // 100,000 BDT
                    MaxSumInsuredAmount = 100000000, // 1,000,000 BDT
                    MinAge = 18,
                    MaxAge = 65,
                    MinTenureMonths = 12,
                    MaxTenureMonths = 36,
                    CreatedBy = createdBy,
                    TenantId = tenantId,
                    CreatedAt = DateTime.UtcNow,
                    UpdatedAt = DateTime.UtcNow
                },
                new Product
                {
                    Id = Guid.NewGuid(),
                    ProductCode = "LIF-001",
                    ProductName = "LabAid Life Shield",
                    ProductNameBn = "ল্যাবএইড লাইফ শিল্ড",
                    Category = ProductCategory.Life,
                    Status = ProductStatus.Active,
                    BasePremiumAmount = 1000000, // 10,000 BDT
                    MinSumInsuredAmount = 50000000, // 500,000 BDT
                    MaxSumInsuredAmount = 500000000, // 5,000,000 BDT
                    MinAge = 18,
                    MaxAge = 60,
                    MinTenureMonths = 60,
                    MaxTenureMonths = 360,
                    CreatedBy = createdBy,
                    TenantId = tenantId,
                    CreatedAt = DateTime.UtcNow,
                    UpdatedAt = DateTime.UtcNow
                },
                new Product
                {
                    Id = Guid.NewGuid(),
                    ProductCode = "TRV-001",
                    ProductName = "Travel Secure",
                    ProductNameBn = "ট্রাভেল সিকিউর",
                    Category = ProductCategory.Travel,
                    Status = ProductStatus.Active,
                    BasePremiumAmount = 200000, // 2,000 BDT
                    MinSumInsuredAmount = 5000000, // 50,000 BDT
                    MaxSumInsuredAmount = 50000000, // 500,000 BDT
                    MinAge = 1,
                    MaxAge = 70,
                    MinTenureMonths = 1,
                    MaxTenureMonths = 12,
                    CreatedBy = createdBy,
                    TenantId = tenantId,
                    CreatedAt = DateTime.UtcNow,
                    UpdatedAt = DateTime.UtcNow
                },
                new Product
                {
                    Id = Guid.NewGuid(),
                    ProductCode = "MTR-001",
                    ProductName = "Motor Shield",
                    ProductNameBn = "মোটর শিল্ড",
                    Category = ProductCategory.Motor,
                    Status = ProductStatus.Active,
                    BasePremiumAmount = 300000, // 3,000 BDT
                    MinSumInsuredAmount = 20000000, // 200,000 BDT
                    MaxSumInsuredAmount = 200000000, // 2,000,000 BDT
                    MinAge = 18,
                    MaxAge = 70,
                    MinTenureMonths = 12,
                    MaxTenureMonths = 12,
                    CreatedBy = createdBy,
                    TenantId = tenantId,
                    CreatedAt = DateTime.UtcNow,
                    UpdatedAt = DateTime.UtcNow
                }
            };

            await productsContext.Products.AddRangeAsync(products);
            await productsContext.SaveChangesAsync();
        }
        catch (Exception ex)
        {
            Console.WriteLine($"An error occurred while initializing the database: {ex.Message}");
        }
    }
}
