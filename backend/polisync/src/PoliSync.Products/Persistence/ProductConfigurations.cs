using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Metadata.Builders;
using PoliSync.Products.Domain;

namespace PoliSync.Products.Persistence;

/// <summary>
/// EF Core configuration for Product table.
/// </summary>
public class ProductConfiguration : IEntityTypeConfiguration<Product>
{
    public void Configure(EntityTypeBuilder<Product> builder)
    {
        builder.ToTable("products", "insurance_schema");
        builder.HasKey(p => p.ProductId);

        builder.Property(p => p.ProductId).HasColumnName("product_id");
        builder.Property(p => p.ProductCode).HasColumnName("product_code").HasMaxLength(50).IsRequired();
        builder.Property(p => p.ProductName).HasColumnName("product_name").HasMaxLength(255).IsRequired();
        builder.Property(p => p.Category).HasColumnName("category").HasMaxLength(50).IsRequired()
            .HasConversion<string>();
        builder.Property(p => p.Description).HasColumnName("description");
        builder.Property(p => p.BasePremium).HasColumnName("base_premium").IsRequired();
        builder.Property(p => p.BasePremiumCurrency).HasColumnName("base_premium_currency").HasMaxLength(3).IsRequired();
        builder.Property(p => p.MinSumInsured).HasColumnName("min_sum_insured").IsRequired();
        builder.Property(p => p.MinSumInsuredCurrency).HasColumnName("min_sum_insured_currency").HasMaxLength(3).IsRequired();
        builder.Property(p => p.MaxSumInsured).HasColumnName("max_sum_insured").IsRequired();
        builder.Property(p => p.MaxSumInsuredCurrency).HasColumnName("max_sum_insured_currency").HasMaxLength(3).IsRequired();
        builder.Property(p => p.MinTenureMonths).HasColumnName("min_tenure_months").IsRequired();
        builder.Property(p => p.MaxTenureMonths).HasColumnName("max_tenure_months").IsRequired();
        builder.Property(p => p.Exclusions).HasColumnName("exclusions")
            .HasColumnType("text[]");
        builder.Property(p => p.Status).HasColumnName("status").HasMaxLength(50).IsRequired()
            .HasConversion<string>();
        builder.Property(p => p.ProductAttributes).HasColumnName("product_attributes").HasColumnType("jsonb");
        builder.Property(p => p.CreatedBy).HasColumnName("created_by").IsRequired();
        builder.Property(p => p.CreatedAt).HasColumnName("created_at").IsRequired();
        builder.Property(p => p.UpdatedAt).HasColumnName("updated_at").IsRequired();
        builder.Property(p => p.DeletedAt).HasColumnName("deleted_at");

        // Relationships
        builder.HasMany(p => p.Plans).WithOne(pp => pp.Product).HasForeignKey(pp => pp.ProductId);
        builder.HasMany(p => p.AvailableRiders).WithOne(r => r.Product).HasForeignKey(r => r.ProductId);
        builder.HasOne(p => p.PricingConfig).WithOne(pc => pc.Product).HasForeignKey<PricingConfig>(pc => pc.ProductId);

        // Indexes
        builder.HasIndex(p => p.ProductCode).IsUnique();
        builder.HasIndex(p => p.Category);
        builder.HasIndex(p => p.Status);

        // Soft delete query filter
        builder.HasQueryFilter(p => p.DeletedAt == null);

        // Ignore domain events (not mapped to DB)
        builder.Ignore(p => p.DomainEvents);
    }
}

public class ProductPlanConfiguration : IEntityTypeConfiguration<ProductPlan>
{
    public void Configure(EntityTypeBuilder<ProductPlan> builder)
    {
        builder.ToTable("product_plans", "insurance_schema");
        builder.HasKey(p => p.PlanId);

        builder.Property(p => p.PlanId).HasColumnName("plan_id");
        builder.Property(p => p.ProductId).HasColumnName("product_id").IsRequired();
        builder.Property(p => p.PlanName).HasColumnName("plan_name").HasMaxLength(255).IsRequired();
        builder.Property(p => p.PlanDescription).HasColumnName("plan_description");
        builder.Property(p => p.PremiumAmount).HasColumnName("premium_amount").IsRequired();
        builder.Property(p => p.MinSumInsured).HasColumnName("min_sum_insured").IsRequired();
        builder.Property(p => p.MaxSumInsured).HasColumnName("max_sum_insured").IsRequired();
        builder.Property(p => p.Attributes).HasColumnName("attributes").HasColumnType("jsonb");
        builder.Property(p => p.CreatedAt).HasColumnName("created_at").IsRequired();
        builder.Property(p => p.UpdatedAt).HasColumnName("updated_at").IsRequired();
    }
}

public class RiderConfiguration : IEntityTypeConfiguration<Rider>
{
    public void Configure(EntityTypeBuilder<Rider> builder)
    {
        builder.ToTable("product_riders", "insurance_schema");
        builder.HasKey(r => r.RiderId);

        builder.Property(r => r.RiderId).HasColumnName("rider_id");
        builder.Property(r => r.ProductId).HasColumnName("product_id").IsRequired();
        builder.Property(r => r.RiderName).HasColumnName("rider_name").HasMaxLength(255).IsRequired();
        builder.Property(r => r.Description).HasColumnName("description");
        builder.Property(r => r.PremiumAmount).HasColumnName("premium_amount").IsRequired();
        builder.Property(r => r.PremiumCurrency).HasColumnName("premium_currency").HasMaxLength(3).IsRequired();
        builder.Property(r => r.CoverageAmount).HasColumnName("coverage_amount").IsRequired();
        builder.Property(r => r.CoverageCurrency).HasColumnName("coverage_currency").HasMaxLength(3).IsRequired();
        builder.Property(r => r.IsMandatory).HasColumnName("is_mandatory").IsRequired();
        builder.Property(r => r.CreatedAt).HasColumnName("created_at").IsRequired();
        builder.Property(r => r.UpdatedAt).HasColumnName("updated_at").IsRequired();

        builder.HasIndex(r => r.ProductId);
    }
}

public class PricingConfigConfiguration : IEntityTypeConfiguration<PricingConfig>
{
    public void Configure(EntityTypeBuilder<PricingConfig> builder)
    {
        builder.ToTable("pricing_configs", "insurance_schema");
        builder.HasKey(pc => pc.PricingConfigId);

        builder.Property(pc => pc.PricingConfigId).HasColumnName("pricing_config_id");
        builder.Property(pc => pc.ProductId).HasColumnName("product_id").IsRequired();
        builder.Property(pc => pc.Rules).HasColumnName("rules").HasColumnType("jsonb").IsRequired();
        builder.Property(pc => pc.EffectiveFrom).HasColumnName("effective_from").IsRequired();
        builder.Property(pc => pc.EffectiveTo).HasColumnName("effective_to");
        builder.Property(pc => pc.CreatedAt).HasColumnName("created_at").IsRequired();
        builder.Property(pc => pc.UpdatedAt).HasColumnName("updated_at").IsRequired();

        builder.HasIndex(pc => pc.ProductId);
    }
}
