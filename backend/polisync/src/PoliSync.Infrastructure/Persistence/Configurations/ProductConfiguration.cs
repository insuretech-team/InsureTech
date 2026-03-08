using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Metadata.Builders;
using PoliSync.Products.Domain;

namespace PoliSync.Infrastructure.Persistence.Configurations;

/// <summary>
/// EF Core configuration for Product aggregate root.
/// </summary>
public class ProductConfiguration : IEntityTypeConfiguration<Product>
{
    public void Configure(EntityTypeBuilder<Product> builder)
    {
        builder.HasKey(p => p.Id);

        builder.Property(p => p.Id)
            .HasColumnName("id")
            .HasDefaultValueSql("gen_random_uuid()");

        builder.Property(p => p.TenantId)
            .HasColumnName("tenant_id")
            .IsRequired();

        builder.Property(p => p.PartnerId)
            .HasColumnName("partner_id")
            .IsRequired();

        builder.Property(p => p.ProductCode)
            .HasColumnName("product_code")
            .HasMaxLength(50)
            .IsRequired();

        builder.Property(p => p.ProductName)
            .HasColumnName("product_name")
            .HasMaxLength(200)
            .IsRequired();

        builder.Property(p => p.Description)
            .HasColumnName("description")
            .HasMaxLength(2000);

        builder.Property(p => p.Category)
            .HasColumnName("category")
            .HasMaxLength(50)
            .IsRequired();

        builder.Property(p => p.Status)
            .HasColumnName("status")
            .HasMaxLength(50)
            .IsRequired();

        builder.Property(p => p.BasePremiumPaisa)
            .HasColumnName("base_premium_paisa");

        builder.Property(p => p.Currency)
            .HasColumnName("currency")
            .HasMaxLength(3)
            .HasDefaultValue("BDT");

        builder.Property(p => p.SumInsuredMinPaisa)
            .HasColumnName("sum_insured_min_paisa");

        builder.Property(p => p.SumInsuredMaxPaisa)
            .HasColumnName("sum_insured_max_paisa");

        builder.Property(p => p.MinTenureMonths)
            .HasColumnName("min_tenure_months");

        builder.Property(p => p.MaxTenureMonths)
            .HasColumnName("max_tenure_months");

        builder.Property(p => p.Exclusions)
            .HasColumnName("exclusions")
            .HasColumnType("jsonb");

        builder.Property(p => p.Version)
            .HasColumnName("version")
            .IsConcurrencyToken();

        builder.Property(p => p.CreatedAt)
            .HasColumnName("created_at");

        builder.Property(p => p.UpdatedAt)
            .HasColumnName("updated_at");

        builder.Property(p => p.CreatedBy)
            .HasColumnName("created_by")
            .HasMaxLength(200);

        // Indices
        builder.HasIndex(p => new { p.TenantId, p.ProductCode })
            .IsUnique();

        builder.HasIndex(p => new { p.TenantId, p.Category });

        // Navigation properties
        builder.HasMany<ProductPlan>()
            .WithOne()
            .HasForeignKey("product_id")
            .OnDelete(DeleteBehavior.Cascade);

        builder.HasMany<Rider>()
            .WithOne()
            .HasForeignKey("product_id")
            .OnDelete(DeleteBehavior.Cascade);

        builder.HasOne(p => p.PricingConfig)
            .WithOne()
            .HasForeignKey<PricingConfig>("product_id")
            .OnDelete(DeleteBehavior.Cascade);
    }
}
