using InsuranceEngine.SharedKernel.Interfaces;
using InsuranceEngine.Products.Domain;
using InsuranceEngine.Products.Domain.Enums;
using InsuranceEngine.Products.Domain.Enums;
using Microsoft.EntityFrameworkCore;
using System;

namespace InsuranceEngine.Products.Infrastructure.Persistence;

public class ProductsDbContext : DbContext
{
    private readonly Guid _tenantId;

    public ProductsDbContext(DbContextOptions<ProductsDbContext> options, ITenantService tenantService) : base(options) 
    { 
        _tenantId = tenantService.GetTenantId();
    }

    public DbSet<Product> Products { get; set; } = null!;
    public DbSet<Insurer> Insurers { get; set; } = null!;
    public DbSet<ProductPlan> ProductPlans { get; set; } = null!;
    public DbSet<PricingRule> PricingRules { get; set; } = null!;
    public DbSet<RiskAssessmentQuestion> RiskAssessmentQuestions { get; set; } = null!;

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        modelBuilder.Entity<Insurer>(entity =>
        {
            entity.HasKey(e => e.Id);
            entity.HasQueryFilter(e => e.TenantId == _tenantId);
            entity.HasIndex(e => e.Code).IsUnique();
        });
        modelBuilder.Entity<Product>(entity =>
        {
            entity.HasKey(e => e.Id);
            entity.HasQueryFilter(e => e.TenantId == _tenantId);
            entity.Property(e => e.Category).HasConversion<string>();
            entity.Property(e => e.Status).HasConversion<string>();
            
            // Relational constraints are minimal across slices
            entity.HasMany(e => e.Plans)
                  .WithOne()
                  .HasForeignKey("ProductId")
                  .OnDelete(DeleteBehavior.Cascade);

            entity.HasMany(e => e.PricingRules)
                  .WithOne()
                  .HasForeignKey("ProductId")
                  .OnDelete(DeleteBehavior.Cascade);

            entity.HasMany(e => e.Questions)
                  .WithOne()
                  .HasForeignKey("ProductId")
                  .OnDelete(DeleteBehavior.Cascade);
        });
    }
}

