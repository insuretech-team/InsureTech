using Insuretech.Products.Entity.V1;
using Microsoft.EntityFrameworkCore;
using Product = Insuretech.Products.Entity.V1.Product;

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

        // Do NOT use HasDefaultSchema("insurance_schema") — Npgsql translates this to
        // SET search_path which pgBouncer (transaction mode) rejects with:
        // "unsupported startup parameter: search_path"
        // Instead use schema-qualified table names directly.

        modelBuilder.Entity<Product>(entity =>
        {
            entity.ToTable("products", schema: "insurance_schema");
            entity.HasKey(e => e.ProductId);

            // ── Mapped columns — force text cast for all uuid columns (Npgsql cannot read uuid as string)
            // product_id, created_by are uuid in DB but string in proto → use HasColumnType("text") workaround
            entity.Property(e => e.ProductId).HasColumnName("product_id").HasColumnType("text")
                .HasConversion(v => v, v => v);
            entity.Property(e => e.ProductCode).HasColumnName("product_code");
            entity.Property(e => e.ProductName).HasColumnName("product_name");
            entity.Property(e => e.Category).HasColumnName("category").HasConversion<string>();
            entity.Property(e => e.Description).HasColumnName("description");
            entity.Property(e => e.Status).HasColumnName("status").HasConversion<string>();
            entity.Property(e => e.MinTenureMonths).HasColumnName("min_tenure_months");
            entity.Property(e => e.MaxTenureMonths).HasColumnName("max_tenure_months");
            entity.Property(e => e.BasePremiumCurrency).HasColumnName("base_premium_currency");
            entity.Property(e => e.MinSumInsuredCurrency).HasColumnName("min_sum_insured_currency");
            entity.Property(e => e.MaxSumInsuredCurrency).HasColumnName("max_sum_insured_currency");

            // ── Ignored: uuid fields that can't be cast + Protobuf complex types EF cannot map ─
            entity.Ignore(e => e.CreatedBy);       // uuid in DB, string in proto — not needed for B2C listing
            entity.Ignore(e => e.BasePremium);     // Money proto type
            entity.Ignore(e => e.MinSumInsured);   // Money proto type
            entity.Ignore(e => e.MaxSumInsured);   // Money proto type
            entity.Ignore(e => e.CreatedAt);       // Timestamp proto type
            entity.Ignore(e => e.UpdatedAt);       // Timestamp proto type
            entity.Ignore(e => e.DeletedAt);       // Timestamp proto type
            entity.Ignore(e => e.AvailableRiders); // Repeated field
            entity.Ignore(e => e.PricingConfig);   // Complex proto type
            entity.Ignore(e => e.Plans);           // Repeated field
            entity.Ignore(e => e.Exclusions);      // Repeated string
            entity.Ignore(e => e.ProductAttributes); // jsonb — not in proto field map
        });
    }
}