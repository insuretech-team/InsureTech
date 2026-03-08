using Microsoft.EntityFrameworkCore;
using InsuranceEngine.Domain.Entities;

namespace InsuranceEngine.Infrastructure.Persistence;

public class InsuranceDbContext : DbContext
{
    public InsuranceDbContext(DbContextOptions<InsuranceDbContext> options) : base(options) { }

    public DbSet<Product> Products => Set<Product>();
    public DbSet<Insurer> Insurers => Set<Insurer>();
    public DbSet<ProductPlan> ProductPlans => Set<ProductPlan>();
    public DbSet<PricingRule> PricingRules => Set<PricingRule>();
    public DbSet<RiskAssessmentQuestion> RiskAssessmentQuestions => Set<RiskAssessmentQuestion>();

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        modelBuilder.HasDefaultSchema("insurance_engine");

        modelBuilder.Entity<Product>(entity =>
        {
            entity.HasKey(e => e.Id);
            entity.Property(e => e.ProductCode).IsRequired().HasMaxLength(50);
            entity.HasOne(e => e.Insurer).WithMany(i => i.Products).HasForeignKey(e => e.InsurerId);
        });

        modelBuilder.Entity<RiskAssessmentQuestion>(entity =>
        {
            entity.HasKey(e => e.Id);
        });
    }
}
