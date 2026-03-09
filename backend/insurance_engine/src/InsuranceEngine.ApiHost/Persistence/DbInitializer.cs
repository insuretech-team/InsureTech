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
            if (await productsContext.Insurers.AnyAsync() || await productsContext.Products.AnyAsync()) return;

            var tenantId = Guid.Parse("11111111-1111-1111-1111-111111111111");

            var insurer = new Insurer
            {
                Id = Guid.NewGuid(),
                Name = "LifePlus Insurance Ltd",
                Code = "LP-001",
                TenantId = tenantId
            };

            var product = new Product
            {
                Id = Guid.NewGuid(),
                ProductCode = "HP-001",
                ProductName = "Health Guard Plus",
                ProductNameBn = "হেলথ গার্ড প্লাস",
                Category = ProductCategory.Health,
                Status = ProductStatus.Active,
                MinSumInsured = 100000,
                MaxSumInsured = 1000000,
                MinAge = 18,
                MaxAge = 65,
                InsurerId = insurer.Id,
                TenantId = tenantId,
                CreatedAt = DateTime.UtcNow,
                UpdatedAt = DateTime.UtcNow
            };

            await productsContext.Insurers.AddAsync(insurer);
            await productsContext.Products.AddAsync(product);
            
            await productsContext.SaveChangesAsync();
            await productsContext.SaveChangesAsync();
        }
        catch (Exception ex)
        {
            Console.WriteLine($"An error occurred while initializing the database: {ex.Message}");
        }
    }
}


