using InsuranceEngine.Partners.Domain.Entities;
using Microsoft.EntityFrameworkCore;

namespace InsuranceEngine.Partners.Infrastructure.Persistence;

public class PartnerDbContext : DbContext
{
    public PartnerDbContext(DbContextOptions<PartnerDbContext> options) : base(options)
    {
    }

    public DbSet<Partner> Partners { get; set; } = null!;
    public DbSet<Agent> Agents { get; set; } = null!;

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        modelBuilder.HasDefaultSchema("partners");

        modelBuilder.Entity<Partner>(entity =>
        {
            entity.ToTable("partners");
            entity.HasKey(e => e.Id);
            
            entity.Property(e => e.Name).HasMaxLength(255).IsRequired();
            entity.Property(e => e.Code).HasMaxLength(50).IsRequired();
            entity.HasIndex(e => e.Code).IsUnique();
            
            entity.Property(e => e.Email).HasMaxLength(255).IsRequired();
            entity.Property(e => e.Status).HasConversion<string>().HasMaxLength(50);
            
            entity.HasMany(e => e.Agents)
                  .WithOne()
                  .HasForeignKey(a => a.PartnerId)
                  .OnDelete(DeleteBehavior.Cascade);
        });

        modelBuilder.Entity<Agent>(entity =>
        {
            entity.ToTable("agents");
            entity.HasKey(e => e.Id);
            
            entity.Property(e => e.Name).HasMaxLength(255).IsRequired();
            entity.Property(e => e.Code).HasMaxLength(50).IsRequired();
            entity.HasIndex(e => e.Code).IsUnique();
            
            entity.Property(e => e.Email).HasMaxLength(255).IsRequired();
            entity.Property(e => e.Status).HasConversion<string>().HasMaxLength(50);
        });
    }
}
