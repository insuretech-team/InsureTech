using Microsoft.EntityFrameworkCore;
using InsuranceEngine.Products.Domain.Entities;

namespace InsuranceEngine.Products.Infrastructure;

public class ProductsDbContext : DbContext
{
    public ProductsDbContext(DbContextOptions<ProductsDbContext> options) : base(options) { }

    public DbSet<ProductEntity> Products => Set<ProductEntity>();
    public DbSet<ProductRiderEntity> ProductRiders => Set<ProductRiderEntity>();

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        modelBuilder.HasDefaultSchema("insurance_schema");

        modelBuilder.Entity<ProductEntity>(entity =>
        {
            entity.ToTable("products");
            entity.HasKey(e => e.ProductId);
            entity.Property(e => e.ProductId).HasColumnName("product_id").HasDefaultValueSql("gen_random_uuid()");
            entity.Property(e => e.TenantId).HasColumnName("tenant_id").HasMaxLength(100).HasDefaultValue("default");
            entity.Property(e => e.ProductCode).HasColumnName("product_code").HasMaxLength(50).IsRequired();
            entity.Property(e => e.ProductName).HasColumnName("product_name").HasMaxLength(200).IsRequired();
            entity.Property(e => e.NameBn).HasColumnName("name_bn").HasMaxLength(200);
            entity.Property(e => e.ProductType).HasColumnName("product_type").HasMaxLength(50).HasDefaultValue("GENERAL");
            entity.Property(e => e.Category).HasColumnName("category").HasMaxLength(100).IsRequired();
            entity.Property(e => e.Description).HasColumnName("description");
            entity.Property(e => e.DescriptionBn).HasColumnName("description_bn");
            entity.Property(e => e.Status).HasColumnName("status").HasMaxLength(50).HasDefaultValue("DRAFT");
            entity.Property(e => e.IsActive).HasColumnName("is_active").HasDefaultValue(false);
            entity.Property(e => e.BasePremium).HasColumnName("base_premium").IsRequired();
            entity.Property(e => e.BasePremiumCurrency).HasColumnName("base_premium_currency").HasMaxLength(3).HasDefaultValue("BDT");
            entity.Property(e => e.MinSumInsured).HasColumnName("min_sum_insured").IsRequired();
            entity.Property(e => e.MinSumInsuredCurrency).HasColumnName("min_sum_insured_currency").HasMaxLength(3).HasDefaultValue("BDT");
            entity.Property(e => e.MaxSumInsured).HasColumnName("max_sum_insured").IsRequired();
            entity.Property(e => e.MaxSumInsuredCurrency).HasColumnName("max_sum_insured_currency").HasMaxLength(3).HasDefaultValue("BDT");
            entity.Property(e => e.UnitAmount).HasColumnName("unit_amount").HasDefaultValue(100000);
            entity.Property(e => e.MinAge).HasColumnName("min_age").IsRequired();
            entity.Property(e => e.MaxAge).HasColumnName("max_age").IsRequired();
            entity.Property(e => e.MinTenureMonths).HasColumnName("min_term_months").IsRequired();
            entity.Property(e => e.MaxTenureMonths).HasColumnName("max_term_months").IsRequired();
            entity.Property(e => e.TermsUrl).HasColumnName("terms_url").HasMaxLength(500);
            entity.Property(e => e.Questions).HasColumnName("questions");
            entity.Property(e => e.Exclusions).HasColumnName("exclusions");
            entity.Property(e => e.ProductAttributes).HasColumnName("product_attributes");
            entity.Property(e => e.CreatedBy).HasColumnName("created_by").HasMaxLength(100).IsRequired();
            entity.Property(e => e.UpdatedBy).HasColumnName("updated_by").HasMaxLength(100);
            entity.Property(e => e.CreatedAt).HasColumnName("created_at").HasDefaultValueSql("CURRENT_TIMESTAMP").IsRequired();
            entity.Property(e => e.UpdatedAt).HasColumnName("updated_at").HasDefaultValueSql("CURRENT_TIMESTAMP").IsRequired();
            entity.Property(e => e.DeletedAt).HasColumnName("deleted_at");
            entity.Property(e => e.Version).HasColumnName("version").HasDefaultValue(1);
            entity.Property(e => e.IsMandatory).HasColumnName("is_mandatory").HasDefaultValue(false);

            entity.HasIndex(e => e.ProductCode).IsUnique().HasDatabaseName("idx_products_code");
            entity.HasIndex(e => e.Category).HasDatabaseName("idx_products_category");
            entity.HasIndex(e => e.Status).HasDatabaseName("idx_products_status");
            entity.HasIndex(e => e.IsActive).HasDatabaseName("idx_products_active");

            entity.HasMany(e => e.Riders).WithOne(r => r.Product).HasForeignKey(r => r.ProductId).OnDelete(DeleteBehavior.Cascade);
        });

        modelBuilder.Entity<ProductRiderEntity>(entity =>
        {
            entity.ToTable("product_riders");
            entity.HasKey(e => e.RiderId);
            entity.Property(e => e.RiderId).HasColumnName("rider_id").HasDefaultValueSql("gen_random_uuid()");
            entity.Property(e => e.ProductId).HasColumnName("product_id").IsRequired();
            entity.Property(e => e.RiderName).HasColumnName("rider_name").HasMaxLength(100).IsRequired();
            entity.Property(e => e.NameEn).HasColumnName("name_en").HasMaxLength(100);
            entity.Property(e => e.NameBn).HasColumnName("name_bn").HasMaxLength(100);
            entity.Property(e => e.Description).HasColumnName("description");
            entity.Property(e => e.PremiumAmount).HasColumnName("premium_amount").IsRequired();
            entity.Property(e => e.PremiumCurrency).HasColumnName("premium_currency").HasMaxLength(3).HasDefaultValue("BDT");
            entity.Property(e => e.AdditionalPremium).HasColumnName("additional_premium").IsRequired();
            entity.Property(e => e.AdditionalPremiumCurrency).HasColumnName("additional_premium_currency").HasMaxLength(3).HasDefaultValue("BDT");
            entity.Property(e => e.CoverageAmount).HasColumnName("coverage_amount").IsRequired();
            entity.Property(e => e.CoverageCurrency).HasColumnName("coverage_currency").HasMaxLength(3).HasDefaultValue("BDT");
            entity.Property(e => e.AdditionalCoverage).HasColumnName("additional_coverage").IsRequired();
            entity.Property(e => e.AdditionalCoverageCurrency).HasColumnName("additional_coverage_currency").HasMaxLength(3).HasDefaultValue("BDT");
            entity.Property(e => e.IsMandatory).HasColumnName("is_mandatory").HasDefaultValue(false);
            entity.Property(e => e.CreatedAt).HasColumnName("created_at").HasDefaultValueSql("CURRENT_TIMESTAMP").IsRequired();
            entity.Property(e => e.UpdatedAt).HasColumnName("updated_at").HasDefaultValueSql("CURRENT_TIMESTAMP").IsRequired();

            entity.HasIndex(e => e.ProductId).HasDatabaseName("idx_product_riders_product_id");
        });
    }
}
