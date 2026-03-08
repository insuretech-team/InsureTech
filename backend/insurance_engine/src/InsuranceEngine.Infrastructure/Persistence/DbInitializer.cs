using Microsoft.Extensions.DependencyInjection;
using Microsoft.EntityFrameworkCore;
using InsuranceEngine.Domain.Entities;
using InsuranceEngine.Domain.Enums;
using InsuranceEngine.Infrastructure.Persistence;
using System;
using System.Threading.Tasks;

namespace InsuranceEngine.Infrastructure.Persistence;

public static class DbInitializer
{
    public static async Task Initialize(IServiceProvider serviceProvider)
    {
        using var scope = serviceProvider.CreateScope();
        var context = scope.ServiceProvider.GetRequiredService<InsuranceDbContext>();

        // Ensure database is created or migrated
        // await context.Database.MigrateAsync();

        try
        {
            if (await context.Insurers.AnyAsync()) return;

            var insurer = new Insurer
            {
                Id = Guid.NewGuid(),
                Name = "LifePlus Insurance Ltd",
                Code = "LP-001"
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
                CreatedAt = DateTime.UtcNow,
                UpdatedAt = DateTime.UtcNow
            };

            await context.Insurers.AddAsync(insurer);
            await context.Products.AddAsync(product);
            await context.SaveChangesAsync();
        }
        catch (Exception ex)
        {
            // Log the exception, but don't crash the app
            Console.WriteLine($"An error occurred while initializing the database: {ex.Message}");
        }
    }
}
