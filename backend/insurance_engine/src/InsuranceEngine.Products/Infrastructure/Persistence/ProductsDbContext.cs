using InsuranceEngine.SharedKernel.Interfaces;
using InsuranceEngine.Products.Domain;
using InsuranceEngine.Products.Domain.Enums;
using Microsoft.EntityFrameworkCore;
using System;
using Newtonsoft.Json;

namespace InsuranceEngine.Products.Infrastructure.Persistence;

public class ProductsDbContext : DbContext
{
    private readonly Guid _tenantId;

    public ProductsDbContext(DbContextOptions<ProductsDbContext> options, ITenantService tenantService) : base(options)
    {
        _tenantId = tenantService.GetTenantId();
    }

    public DbSet<Product> Products { get; set; } = null!;
    public DbSet<ProductPlan> ProductPlans { get; set; } = null!;
    public DbSet<Rider> Riders { get; set; } = null!;
    public DbSet<PricingConfig> PricingConfigs { get; set; } = null!;
    public DbSet<RiskAssessmentQuestion> RiskAssessmentQuestions { get; set; } = null!;

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        modelBuilder.HasDefaultSchema("insurance_schema");

        // --- Product ---
        modelBuilder.Entity<Product>(entity =>
        {
            entity.ToTable("products");
            entity.HasKey(e => e.Id);
            entity.HasQueryFilter(e => e.TenantId == _tenantId && !e.IsDeleted);

            entity.Property(e => e.ProductCode).HasMaxLength(50).IsRequired();
            entity.HasIndex(e => e.ProductCode).IsUnique();

            entity.Property(e => e.ProductName).HasMaxLength(255).IsRequired();
            entity.Property(e => e.Category).HasConversion<string>().HasMaxLength(50).IsRequired();
            entity.Property(e => e.Status).HasConversion<string>().HasMaxLength(50).IsRequired();

            // Money columns stored as bigint
            entity.Property(e => e.BasePremiumAmount).HasColumnName("base_premium").IsRequired();
            entity.Property(e => e.BasePremiumCurrency).HasColumnName("base_premium_currency").HasMaxLength(3).IsRequired().HasDefaultValue("BDT");
            entity.Property(e => e.MinSumInsuredAmount).HasColumnName("min_sum_insured").IsRequired();
            entity.Property(e => e.MinSumInsuredCurrency).HasColumnName("min_sum_insured_currency").HasMaxLength(3).IsRequired().HasDefaultValue("BDT");
            entity.Property(e => e.MaxSumInsuredAmount).HasColumnName("max_sum_insured").IsRequired();
            entity.Property(e => e.MaxSumInsuredCurrency).HasColumnName("max_sum_insured_currency").HasMaxLength(3).IsRequired().HasDefaultValue("BDT");

            entity.Property(e => e.Exclusions).HasColumnType("text[]");
            entity.Property(e => e.ProductAttributes).HasColumnType("jsonb");

            entity.HasIndex(e => e.Category);
            entity.HasIndex(e => e.Status);

            // Ignore computed Money properties
            entity.Ignore(e => e.BasePremium);
            entity.Ignore(e => e.MinSumInsured);
            entity.Ignore(e => e.MaxSumInsured);

            // Relationships
            entity.HasMany(e => e.Plans)
                  .WithOne(p => p.Product)
                  .HasForeignKey(p => p.ProductId)
                  .OnDelete(DeleteBehavior.Cascade);

            entity.HasMany(e => e.AvailableRiders)
                  .WithOne(r => r.Product)
                  .HasForeignKey(r => r.ProductId)
                  .OnDelete(DeleteBehavior.Cascade);

            entity.HasOne(e => e.PricingConfig)
                  .WithOne(pc => pc.Product)
                  .HasForeignKey<PricingConfig>(pc => pc.ProductId)
                  .OnDelete(DeleteBehavior.Cascade);

            entity.HasMany(e => e.Questions)
                  .WithOne()
                  .HasForeignKey("ProductId")
                  .OnDelete(DeleteBehavior.Cascade);
        });

        // --- Rider ---
        modelBuilder.Entity<Rider>(entity =>
        {
            entity.ToTable("product_riders");
            entity.HasKey(e => e.Id);

            entity.Property(e => e.RiderName).HasMaxLength(255).IsRequired();
            entity.Property(e => e.PremiumAmount).HasColumnName("premium_amount").IsRequired();
            entity.Property(e => e.PremiumCurrency).HasColumnName("premium_currency").HasMaxLength(3).IsRequired().HasDefaultValue("BDT");
            entity.Property(e => e.CoverageAmount).HasColumnName("coverage_amount").IsRequired();
            entity.Property(e => e.CoverageCurrency).HasColumnName("coverage_currency").HasMaxLength(3).IsRequired().HasDefaultValue("BDT");

            entity.Ignore(e => e.Premium);
            entity.Ignore(e => e.Coverage);

            entity.HasIndex(e => e.ProductId);
        });

        // --- PricingConfig ---
        modelBuilder.Entity<PricingConfig>(entity =>
        {
            entity.ToTable("pricing_configs");
            entity.HasKey(e => e.Id);

            entity.Property(e => e.Rules)
                  .HasColumnType("jsonb")
                  .HasConversion(
                      v => JsonConvert.SerializeObject(v),
                      v => JsonConvert.DeserializeObject<List<PricingRule>>(v) ?? new());

            entity.HasIndex(e => e.ProductId);
        });

        // --- ProductPlan ---
        modelBuilder.Entity<ProductPlan>(entity =>
        {
            entity.ToTable("product_plans");
            entity.HasKey(e => e.Id);

            entity.Property(e => e.PlanName).HasMaxLength(255).IsRequired();
            entity.Property(e => e.PremiumAmount).HasColumnName("premium_amount").IsRequired();
            entity.Property(e => e.PremiumCurrency).HasColumnName("premium_currency").HasMaxLength(3).HasDefaultValue("BDT");
            entity.Property(e => e.MinSumInsuredAmount).HasColumnName("min_sum_insured").IsRequired();
            entity.Property(e => e.MinSumInsuredCurrency).HasColumnName("min_sum_insured_currency").HasMaxLength(3).HasDefaultValue("BDT");
            entity.Property(e => e.MaxSumInsuredAmount).HasColumnName("max_sum_insured").IsRequired();
            entity.Property(e => e.MaxSumInsuredCurrency).HasColumnName("max_sum_insured_currency").HasMaxLength(3).HasDefaultValue("BDT");
            entity.Property(e => e.Attributes).HasColumnType("jsonb");

            entity.Ignore(e => e.Premium);
            entity.Ignore(e => e.MinSumInsured);
            entity.Ignore(e => e.MaxSumInsured);
        });

        // --- RiskAssessmentQuestion ---
        modelBuilder.Entity<RiskAssessmentQuestion>(entity =>
        {
            entity.ToTable("risk_assessment_questions");
            entity.HasKey(e => e.Id);
        });
    }
}
