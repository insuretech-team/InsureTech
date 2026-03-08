using Insuretech.Products.Entity.V1;
using Microsoft.EntityFrameworkCore;

namespace PoliSync.Infrastructure.Persistence;

/// <summary>
/// Minimal EF Core context for active proto-based implementation.
/// </summary>
public class PoliSyncDbContext : DbContext
{
    public PoliSyncDbContext(DbContextOptions<PoliSyncDbContext> options) : base(options)
    {
    }

    public DbSet<Product> Products => Set<Product>();

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        base.OnModelCreating(modelBuilder);

        modelBuilder.HasDefaultSchema("insurance_schema");

        modelBuilder.Entity<Product>(entity =>
        {
            entity.ToTable("products");
            entity.HasKey(e => e.ProductId);

            entity.Property(e => e.ProductId).HasColumnName("product_id");
            entity.Property(e => e.ProductCode).HasColumnName("product_code");
            entity.Property(e => e.ProductName).HasColumnName("product_name");
            entity.Property(e => e.Category).HasColumnName("category").HasConversion<string>();
            entity.Property(e => e.Description).HasColumnName("description");
            entity.Property(e => e.Status).HasColumnName("status").HasConversion<string>();
            entity.Property(e => e.MinTenureMonths).HasColumnName("min_tenure_months");
            entity.Property(e => e.MaxTenureMonths).HasColumnName("max_tenure_months");
            entity.Property(e => e.CreatedBy).HasColumnName("created_by");
            entity.Property(e => e.BasePremiumCurrency).HasColumnName("base_premium_currency");
            entity.Property(e => e.MinSumInsuredCurrency).HasColumnName("min_sum_insured_currency");
            entity.Property(e => e.MaxSumInsuredCurrency).HasColumnName("max_sum_insured_currency");

            entity.Ignore(e => e.BasePremium);
            entity.Ignore(e => e.MinSumInsured);
            entity.Ignore(e => e.MaxSumInsured);
            entity.Ignore(e => e.CreatedAt);
            entity.Ignore(e => e.UpdatedAt);
            entity.Ignore(e => e.DeletedAt);
            entity.Ignore(e => e.AvailableRiders);
            entity.Ignore(e => e.PricingConfig);
            entity.Ignore(e => e.Plans);
            entity.Ignore(e => e.Exclusions);
        });
    }
}