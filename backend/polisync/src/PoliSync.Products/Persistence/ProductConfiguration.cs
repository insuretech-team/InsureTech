using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Metadata.Builders;

namespace PoliSync.Products.Persistence;

public sealed class ProductRecordConfiguration : IEntityTypeConfiguration<ProductRecord>
{
    public void Configure(EntityTypeBuilder<ProductRecord> b)
    {
        b.ToTable("products", "insurance_schema");
        b.HasKey(x => x.ProductId);
        b.Property(x => x.ProductId).HasColumnName("product_id").HasDefaultValueSql("gen_random_uuid()");
        b.Property(x => x.ProductCode).HasColumnName("product_code").HasMaxLength(50).IsRequired();
        b.Property(x => x.ProductName).HasColumnName("product_name").HasMaxLength(255).IsRequired();
        b.Property(x => x.Category).HasColumnName("category").HasMaxLength(50).IsRequired();
        b.Property(x => x.Description).HasColumnName("description");
        b.Property(x => x.BasePremium).HasColumnName("base_premium").IsRequired();
        b.Property(x => x.BasePremiumCurrency).HasColumnName("base_premium_currency").HasMaxLength(3).HasDefaultValue("BDT").IsRequired();
        b.Property(x => x.MinSumInsured).HasColumnName("min_sum_insured").IsRequired();
        b.Property(x => x.MinSumInsuredCurrency).HasColumnName("min_sum_insured_currency").HasMaxLength(3).HasDefaultValue("BDT").IsRequired();
        b.Property(x => x.MaxSumInsured).HasColumnName("max_sum_insured").IsRequired();
        b.Property(x => x.MaxSumInsuredCurrency).HasColumnName("max_sum_insured_currency").HasMaxLength(3).HasDefaultValue("BDT").IsRequired();
        b.Property(x => x.MinTenureMonths).HasColumnName("min_tenure_months").IsRequired();
        b.Property(x => x.MaxTenureMonths).HasColumnName("max_tenure_months").IsRequired();
        b.Property(x => x.Exclusions).HasColumnName("exclusions").HasColumnType("text[]");
        b.Property(x => x.Status).HasColumnName("status").HasMaxLength(50).IsRequired();
        b.Property(x => x.ProductAttributes).HasColumnName("product_attributes").HasColumnType("jsonb");
        b.Property(x => x.CreatedBy).HasColumnName("created_by").IsRequired();
        b.Property(x => x.CreatedAt).HasColumnName("created_at").IsRequired();
        b.Property(x => x.UpdatedAt).HasColumnName("updated_at").IsRequired();
        b.Property(x => x.DeletedAt).HasColumnName("deleted_at");

        b.HasQueryFilter(x => x.DeletedAt == null); // soft delete global filter

        b.HasMany(x => x.Plans).WithOne(x => x.Product).HasForeignKey(x => x.ProductId);
        b.HasMany(x => x.Riders).WithOne(x => x.Product).HasForeignKey(x => x.ProductId);
        b.HasOne(x => x.PricingConfig).WithOne(x => x.Product).HasForeignKey<PricingConfigRecord>(x => x.ProductId);
    }
}

public sealed class ProductPlanRecordConfiguration : IEntityTypeConfiguration<ProductPlanRecord>
{
    public void Configure(EntityTypeBuilder<ProductPlanRecord> b)
    {
        b.ToTable("product_plans", "insurance_schema");
        b.HasKey(x => x.PlanId);
        b.Property(x => x.PlanId).HasColumnName("plan_id").HasDefaultValueSql("gen_random_uuid()");
        b.Property(x => x.ProductId).HasColumnName("product_id").IsRequired();
        b.Property(x => x.PlanName).HasColumnName("plan_name").HasMaxLength(255).IsRequired();
        b.Property(x => x.PlanDescription).HasColumnName("plan_description");
        b.Property(x => x.PremiumAmount).HasColumnName("premium_amount").IsRequired();
        b.Property(x => x.PremiumCurrency).HasColumnName("premium_currency").HasMaxLength(3).HasDefaultValue("BDT");
        b.Property(x => x.MinSumInsured).HasColumnName("min_sum_insured").IsRequired();
        b.Property(x => x.MinSumInsuredCurrency).HasColumnName("min_sum_insured_currency").HasMaxLength(3).HasDefaultValue("BDT");
        b.Property(x => x.MaxSumInsured).HasColumnName("max_sum_insured").IsRequired();
        b.Property(x => x.MaxSumInsuredCurrency).HasColumnName("max_sum_insured_currency").HasMaxLength(3).HasDefaultValue("BDT");
        b.Property(x => x.Attributes).HasColumnName("attributes").HasColumnType("jsonb");
        b.Property(x => x.CreatedAt).HasColumnName("created_at").IsRequired();
        b.Property(x => x.UpdatedAt).HasColumnName("updated_at").IsRequired();
    }
}

public sealed class RiderRecordConfiguration : IEntityTypeConfiguration<RiderRecord>
{
    public void Configure(EntityTypeBuilder<RiderRecord> b)
    {
        b.ToTable("product_riders", "insurance_schema");
        b.HasKey(x => x.RiderId);
        b.Property(x => x.RiderId).HasColumnName("rider_id").HasDefaultValueSql("gen_random_uuid()");
        b.Property(x => x.ProductId).HasColumnName("product_id").IsRequired();
        b.Property(x => x.RiderName).HasColumnName("rider_name").HasMaxLength(255).IsRequired();
        b.Property(x => x.Description).HasColumnName("description");
        b.Property(x => x.PremiumAmount).HasColumnName("premium_amount").IsRequired();
        b.Property(x => x.PremiumCurrency).HasColumnName("premium_currency").HasMaxLength(3).HasDefaultValue("BDT");
        b.Property(x => x.CoverageAmount).HasColumnName("coverage_amount").IsRequired();
        b.Property(x => x.CoverageCurrency).HasColumnName("coverage_currency").HasMaxLength(3).HasDefaultValue("BDT");
        b.Property(x => x.IsMandatory).HasColumnName("is_mandatory").HasDefaultValue(false).IsRequired();
        b.Property(x => x.CreatedAt).HasColumnName("created_at").IsRequired();
        b.Property(x => x.UpdatedAt).HasColumnName("updated_at").IsRequired();
    }
}

public sealed class PricingConfigRecordConfiguration : IEntityTypeConfiguration<PricingConfigRecord>
{
    public void Configure(EntityTypeBuilder<PricingConfigRecord> b)
    {
        b.ToTable("pricing_configs", "insurance_schema");
        b.HasKey(x => x.PricingConfigId);
        b.Property(x => x.PricingConfigId).HasColumnName("pricing_config_id").HasDefaultValueSql("gen_random_uuid()");
        b.Property(x => x.ProductId).HasColumnName("product_id").IsRequired();
        b.Property(x => x.Rules).HasColumnName("rules").HasColumnType("jsonb").IsRequired();
        b.Property(x => x.EffectiveFrom).HasColumnName("effective_from").IsRequired();
        b.Property(x => x.EffectiveTo).HasColumnName("effective_to");
        b.Property(x => x.CreatedAt).HasColumnName("created_at").IsRequired();
        b.Property(x => x.UpdatedAt).HasColumnName("updated_at").IsRequired();
    }
}
