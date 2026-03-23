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
            entity.Property(e => e.Id).HasColumnName("partner_id");
            entity.HasQueryFilter(e => e.DeletedAt == null);

            entity.Property(e => e.OrganizationName).HasColumnName("organization_name").HasMaxLength(255).IsRequired();
            entity.Property(e => e.Type).HasColumnName("type").HasMaxLength(50);
            entity.Property(e => e.Code).HasColumnName("code").HasMaxLength(50).IsRequired();
            entity.HasIndex(e => e.Code).IsUnique();
            
            entity.Property(e => e.Email).HasColumnName("email").HasMaxLength(255).IsRequired();
            entity.Property(e => e.Phone).HasColumnName("phone").HasMaxLength(20);
            entity.Property(e => e.Status).HasConversion<string>().HasColumnName("status").HasMaxLength(50);
            
            entity.Property(e => e.TradeLicense).HasColumnName("trade_license").HasMaxLength(100);
            entity.Property(e => e.BankAccount).HasColumnName("bank_account").HasMaxLength(100);
            entity.Property(e => e.AcquisitionCommissionRate).HasColumnName("acquisition_commission_rate");
            entity.Property(e => e.RenewalCommissionRate).HasColumnName("renewal_commission_rate");
            entity.Property(e => e.OnboardedAt).HasColumnName("onboarded_at");
            entity.Property(e => e.FocalPersonId).HasColumnName("focal_person_id");
            
            entity.Property(e => e.CommissionJson).HasColumnType("jsonb").HasColumnName("commission");
            entity.Property(e => e.BenefitsJson).HasColumnType("jsonb").HasColumnName("benefits");
            entity.Property(e => e.DeletedAt).HasColumnName("deleted_at");

            entity.HasMany(e => e.Agents)
                  .WithOne()
                  .HasForeignKey(a => a.PartnerId)
                  .HasConstraintName("fk_agents_partners_partner_id")
                  .OnDelete(DeleteBehavior.Cascade);
        });

        modelBuilder.Entity<Agent>(entity =>
        {
            entity.ToTable("agents");
            entity.HasKey(e => e.Id);
            entity.Property(e => e.Id).HasColumnName("agent_id");
            entity.HasQueryFilter(e => e.DeletedAt == null);

            entity.Property(e => e.PartnerId).HasColumnName("partner_id").IsRequired();
            entity.Property(e => e.Name).HasColumnName("name").HasMaxLength(255).IsRequired();
            entity.Property(e => e.Code).HasColumnName("code").HasMaxLength(50).IsRequired();
            entity.HasIndex(e => e.Code).IsUnique();
            
            entity.Property(e => e.Email).HasColumnName("email").HasMaxLength(255).IsRequired();
            entity.Property(e => e.Phone).HasColumnName("phone").HasMaxLength(20);
            entity.Property(e => e.Status).HasConversion<string>().HasColumnName("status").HasMaxLength(50);
            entity.Property(e => e.DeletedAt).HasColumnName("deleted_at");
        });
    }
}
