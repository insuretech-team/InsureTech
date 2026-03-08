using System.Text.Json;
using System.Text.Json.Serialization;
using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Metadata.Builders;
using PoliSync.Products.Domain;

namespace PoliSync.Infrastructure.Persistence.Configurations;

/// <summary>
/// EF Core configuration for PricingConfig entity.
/// Handles JSON serialization of complex PricingRule objects.
/// </summary>
public class PricingConfigConfiguration : IEntityTypeConfiguration<PricingConfig>
{
    private static readonly JsonSerializerOptions JsonOptions = new()
    {
        PropertyNamingPolicy = JsonNamingPolicy.CamelCase,
        DefaultIgnoreCondition = JsonIgnoreCondition.WhenWritingNull,
        WriteIndented = false
    };

    public void Configure(EntityTypeBuilder<PricingConfig> builder)
    {
        builder.HasKey(p => p.Id);

        builder.Property(p => p.Id)
            .HasColumnName("id")
            .HasDefaultValueSql("gen_random_uuid()");

        builder.Property(p => p.ProductId)
            .HasColumnName("product_id");

        builder.Property(p => p.Rules)
            .HasColumnName("rules")
            .HasColumnType("jsonb")
            .HasConversion(
                v => JsonSerializer.Serialize(v, JsonOptions),
                v => JsonSerializer.Deserialize<List<PricingRule>>(v, JsonOptions) ?? []);

        builder.Property(p => p.EffectiveFrom)
            .HasColumnName("effective_from");

        builder.Property(p => p.EffectiveTo)
            .HasColumnName("effective_to");

        builder.Property(p => p.Version)
            .HasColumnName("version");

        builder.Property(p => p.CreatedAt)
            .HasColumnName("created_at");

        builder.Property(p => p.UpdatedAt)
            .HasColumnName("updated_at");

        // Indices
        builder.HasIndex(p => p.ProductId)
            .IsUnique();
    }
}
