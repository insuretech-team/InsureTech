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
        try
        {
            // Migrations are now managed externally via Go/SQL scripts in backend/inscore
            // to align with the PoliSync architectural pattern.
            
            using var seedScope = serviceProvider.CreateScope();
            var seedContext = seedScope.ServiceProvider.GetRequiredService<ProductsDbContext>();
            if (await seedContext.Products.AnyAsync()) return;

            var tenantId = Guid.Parse("11111111-1111-1111-1111-111111111111");
            var createdBy = Guid.Parse("00000000-0000-0000-0000-000000000001");

            var products = new[]
            {
                new Product
                {
                    Id = Guid.NewGuid(),
                    ProductCode = "HLT-001",
                    ProductName = "Health Guard Plus",
                    Category = ProductCategory.Health,
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
                    Category = ProductCategory.Life,
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
                    Category = ProductCategory.Travel,
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
                    Category = ProductCategory.Motor,
                    MinTenureMonths = 12,
                    MaxTenureMonths = 12,
                    CreatedBy = createdBy,
                    TenantId = tenantId,
                    CreatedAt = DateTime.UtcNow,
                    UpdatedAt = DateTime.UtcNow
                }
            };

            await seedContext.Products.AddRangeAsync(products);
            await seedContext.SaveChangesAsync();
        }
        catch (Exception ex)
        {
            Console.WriteLine($"An error occurred while initializing the database: {ex.Message}");
        }
    }
}
