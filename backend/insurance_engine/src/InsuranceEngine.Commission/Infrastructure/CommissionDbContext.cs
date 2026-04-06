using Microsoft.EntityFrameworkCore;
using InsuranceEngine.SharedKernel.Persistence.Entities;

namespace InsuranceEngine.Commission.Infrastructure;

public class CommissionDbContext : DbContext
{
    public CommissionDbContext(DbContextOptions<CommissionDbContext> options) : base(options) { }

    public DbSet<CommissionEntity> Commissions => Set<CommissionEntity>();
    public DbSet<CommissionPayoutEntity> CommissionPayouts => Set<CommissionPayoutEntity>();

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        modelBuilder.HasDefaultSchema("insurance_schema");
        modelBuilder.Entity<CommissionEntity>().ToTable("commissions");
        modelBuilder.Entity<CommissionPayoutEntity>().ToTable("commission_payouts");
    }
}
