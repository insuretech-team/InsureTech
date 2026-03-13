using Microsoft.EntityFrameworkCore;
using InsuranceEngine.Fraud.Domain.Entities;

namespace InsuranceEngine.Fraud.Infrastructure.Persistence;

public class FraudDbContext : DbContext
{
    public FraudDbContext(DbContextOptions<FraudDbContext> options) : base(options) { }

    public DbSet<FraudCheck> FraudChecks => Set<FraudCheck>();

    protected override void OnModelCreating(ModelCreatingBuilder modelBuilder)
    {
        modelBuilder.HasDefaultSchema("insurance_schema");

        modelBuilder.Entity<FraudCheck>(entity =>
        {
            entity.ToTable("fraud_checks");
            entity.HasKey(e => e.Id);
            entity.Property(e => e.FindingsJson).HasColumnType("jsonb");
            entity.Property(e => e.CheckedRulesJson).HasColumnType("jsonb");
        });
    }
}
