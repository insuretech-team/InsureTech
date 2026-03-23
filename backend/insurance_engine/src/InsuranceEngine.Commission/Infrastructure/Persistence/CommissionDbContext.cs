using InsuranceEngine.Commission.Domain.Entities;
using Microsoft.EntityFrameworkCore;

namespace InsuranceEngine.Commission.Infrastructure.Persistence;

public class CommissionDbContext : DbContext
{
    public CommissionDbContext(DbContextOptions<CommissionDbContext> options) : base(options)
    {
    }

    public DbSet<Domain.Entities.Commission> Commissions { get; set; } = null!;
    public DbSet<Payout> Payouts { get; set; } = null!;

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        modelBuilder.HasDefaultSchema("commissions");

        modelBuilder.Entity<Domain.Entities.Commission>(entity =>
        {
            entity.ToTable("commissions");
            entity.HasKey(e => e.Id);
            
            entity.Property(e => e.Type).HasConversion<string>().HasMaxLength(50);
            entity.Property(e => e.Status).HasConversion<string>().HasMaxLength(50);
            entity.Property(e => e.Currency).HasMaxLength(3);
            
            entity.HasIndex(e => e.PolicyId);
            entity.HasIndex(e => e.PartnerId);
            entity.HasIndex(e => e.AgentId);
            entity.HasIndex(e => e.Status);
        });

        modelBuilder.Entity<Payout>(entity =>
        {
            entity.ToTable("payouts");
            entity.HasKey(e => e.Id);
            
            entity.Property(e => e.Status).HasConversion<string>().HasMaxLength(50);
            entity.Property(e => e.Currency).HasMaxLength(3);
            entity.Property(e => e.PaymentReference).HasMaxLength(255);
            
            entity.HasIndex(e => e.RecipientId);
            entity.HasIndex(e => e.Status);
        });
    }
}
