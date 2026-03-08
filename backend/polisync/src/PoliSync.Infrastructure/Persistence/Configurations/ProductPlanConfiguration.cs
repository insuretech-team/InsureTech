using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Metadata.Builders;
using PoliSync.Products.Domain;

namespace PoliSync.Infrastructure.Persistence.Configurations;

/// <summary>
/// EF Core configuration for ProductPlan entity.
/// </summary>
public class ProductPlanConfiguration : IEntityTypeConfiguration<ProductPlan>
{
    public void Configure(EntityTypeBuilder<ProductPlan> builder)
    {
        builder.HasKey(p => p.Id);

        builder.Property(p => p.Id)
            .HasColumnName("id")
            .HasDefaultValueSql("gen_random_uuid()");

        builder.Property(p => p.ProductId)
            .HasColumnName("product_id");

        builder.Property(p => p.PlanCode)
            .HasColumnName("plan_code")
            .HasMaxLength(50)
            .IsRequired();

        builder.Property(p => p.PlanName)
            .HasColumnName("plan_name")
            .HasMaxLength(200)
            .IsRequired();

        builder.Property(p => p.Description)
            .HasColumnName("description")
            .HasMaxLength(2000);

        builder.Property(p => p.BasePremiumPaisa)
            .HasColumnName("base_premium_paisa");

        builder.Property(p => p.SumInsuredPaisa)
            .HasColumnName("sum_insured_paisa");

        builder.Property(p => p.Currency)
            .HasColumnName("currency")
            .HasMaxLength(3)
            .HasDefaultValue("BDT");

        builder.Property(p => p.Features)
            .HasColumnName("features")
            .HasColumnType("jsonb");

        builder.Property(p => p.IsActive)
            .HasColumnName("is_active")
            .HasDefaultValue(true);

        builder.Property(p => p.CreatedAt)
            .HasColumnName("created_at");

        builder.Property(p => p.UpdatedAt)
            .HasColumnName("updated_at");

        // Indices
        builder.HasIndex(p => p.ProductId);

        builder.HasIndex(p => new { p.ProductId, p.PlanCode })
            .IsUnique();
    }
}
