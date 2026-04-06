using Microsoft.EntityFrameworkCore;
using InsuranceEngine.SharedKernel.Persistence.Entities;

namespace InsuranceEngine.Policy.Infrastructure;

public class PolicyDbContext : DbContext
{
    public PolicyDbContext(DbContextOptions<PolicyDbContext> options) : base(options) { }

    public DbSet<PolicyEntity> Policies => Set<PolicyEntity>();
    public DbSet<PolicyNomineeEntity> PolicyNominees => Set<PolicyNomineeEntity>();
    public DbSet<PolicyRiderEntity> PolicyRiders => Set<PolicyRiderEntity>();

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        modelBuilder.HasDefaultSchema("insurance_schema");

        modelBuilder.Entity<PolicyEntity>(entity =>
        {
            entity.HasKey(e => e.PolicyId);
            entity.Property(e => e.PolicyNumber).IsRequired().HasMaxLength(50);
            entity.Property(e => e.Status).HasMaxLength(50).HasDefaultValue("PENDING_PAYMENT");
            entity.Property(e => e.PremiumCurrency).HasMaxLength(3).HasDefaultValue("BDT");
            entity.Property(e => e.SumInsuredCurrency).HasMaxLength(3).HasDefaultValue("BDT");
            entity.HasIndex(e => e.PolicyNumber).IsUnique();
            entity.HasIndex(e => e.CustomerId);
            entity.HasIndex(e => e.Status);
        });

        modelBuilder.Entity<PolicyNomineeEntity>(entity =>
        {
            entity.HasKey(e => e.NomineeId);
            entity.Property(e => e.FullName).IsRequired().HasMaxLength(200);
            entity.Property(e => e.Relationship).HasMaxLength(50);
            entity.Property(e => e.SharePercentage).HasDefaultValue(100);
            entity.HasIndex(e => e.PolicyId);
        });

        modelBuilder.Entity<PolicyRiderEntity>(entity =>
        {
            entity.HasKey(e => e.RiderId);
            entity.Property(e => e.RiderName).HasMaxLength(100);
            entity.HasIndex(e => e.PolicyId);
        });
    }
}
