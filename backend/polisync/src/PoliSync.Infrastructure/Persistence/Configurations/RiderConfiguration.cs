using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Metadata.Builders;
using PoliSync.Products.Domain;

namespace PoliSync.Infrastructure.Persistence.Configurations;

/// <summary>
/// EF Core configuration for Rider entity.
/// </summary>
public class RiderConfiguration : IEntityTypeConfiguration<Rider>
{
    public void Configure(EntityTypeBuilder<Rider> builder)
    {
        builder.HasKey(r => r.Id);

        builder.Property(r => r.Id)
            .HasColumnName("id")
            .HasDefaultValueSql("gen_random_uuid()");

        builder.Property(r => r.ProductId)
            .HasColumnName("product_id");

        builder.Property(r => r.RiderCode)
            .HasColumnName("rider_code")
            .HasMaxLength(50)
            .IsRequired();

        builder.Property(r => r.RiderName)
            .HasColumnName("rider_name")
            .HasMaxLength(200)
            .IsRequired();

        builder.Property(r => r.Description)
            .HasColumnName("description")
            .HasMaxLength(2000);

        builder.Property(r => r.PremiumAmountPaisa)
            .HasColumnName("premium_amount_paisa");

        builder.Property(r => r.SumInsuredPaisa)
            .HasColumnName("sum_insured_paisa");

        builder.Property(r => r.Currency)
            .HasColumnName("currency")
            .HasMaxLength(3)
            .HasDefaultValue("BDT");

        builder.Property(r => r.Category)
            .HasColumnName("category")
            .HasMaxLength(100);

        builder.Property(r => r.IsMandatory)
            .HasColumnName("is_mandatory");

        builder.Property(r => r.IsActive)
            .HasColumnName("is_active");

        builder.Property(r => r.CreatedAt)
            .HasColumnName("created_at");

        builder.Property(r => r.UpdatedAt)
            .HasColumnName("updated_at");

        // Indices
        builder.HasIndex(r => r.ProductId);

        builder.HasIndex(r => new { r.ProductId, r.RiderCode })
            .IsUnique();
    }
}
