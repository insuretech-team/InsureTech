using InsuranceEngine.Claims.Domain.Entities;
using InsuranceEngine.Claims.Domain.Enums;
using Microsoft.EntityFrameworkCore;

namespace InsuranceEngine.Claims.Infrastructure.Persistence;

public class ClaimsDbContext : DbContext
{
    public ClaimsDbContext(DbContextOptions<ClaimsDbContext> options) : base(options)
    {
    }

    public DbSet<Claim> Claims { get; set; } = null!;
    public DbSet<ClaimApproval> ClaimApprovals { get; set; } = null!;
    public DbSet<ClaimDocument> ClaimDocuments { get; set; } = null!;

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        modelBuilder.HasDefaultSchema("insurance_schema");

        modelBuilder.Entity<Claim>(entity =>
        {
            entity.ToTable("claims");
            entity.HasKey(e => e.Id);
            entity.HasQueryFilter(e => !e.IsDeleted);

            entity.Property(e => e.ClaimNumber).HasMaxLength(50).IsRequired();
            entity.HasIndex(e => e.ClaimNumber).IsUnique();
            entity.Property(e => e.Status).HasConversion<string>().HasMaxLength(50).IsRequired();
            entity.Property(e => e.Type).HasConversion<string>().HasMaxLength(50).IsRequired();
            entity.Property(e => e.ProcessingType).HasConversion<string>().HasMaxLength(50).IsRequired();

            entity.Property(e => e.ClaimedAmount).HasColumnName("claimed_amount").IsRequired();
            entity.Property(e => e.ClaimedCurrency).HasColumnName("claimed_currency").HasMaxLength(3).HasDefaultValue("BDT");
            entity.Property(e => e.ApprovedAmount).HasColumnName("approved_amount");
            entity.Property(e => e.ApprovedCurrency).HasColumnName("approved_currency").HasMaxLength(3).HasDefaultValue("BDT");
            entity.Property(e => e.SettledAmount).HasColumnName("settled_amount");
            entity.Property(e => e.SettledCurrency).HasColumnName("settled_currency").HasMaxLength(3).HasDefaultValue("BDT");

            entity.Property(e => e.FraudCheckData).HasColumnType("jsonb").HasColumnName("fraud_check_data");

            entity.HasIndex(e => e.PolicyId);
            entity.HasIndex(e => e.CustomerId);
            entity.HasIndex(e => e.Status);

            entity.HasMany(e => e.Approvals)
                  .WithOne()
                  .HasForeignKey(a => a.ClaimId)
                  .OnDelete(DeleteBehavior.Cascade);

            entity.HasMany(e => e.Documents)
                  .WithOne()
                  .HasForeignKey(d => d.ClaimId)
                  .OnDelete(DeleteBehavior.Cascade);
        });

        modelBuilder.Entity<ClaimApproval>(entity =>
        {
            entity.ToTable("claim_approvals");
            entity.HasKey(e => e.Id);
            entity.Property(e => e.Decision).HasConversion<string>().HasMaxLength(50).IsRequired();
            entity.Property(e => e.ApprovedAmount).HasColumnName("approved_amount");
            entity.Property(e => e.ApprovedCurrency).HasColumnName("approved_currency").HasMaxLength(3).HasDefaultValue("BDT");
            
            entity.HasIndex(e => e.ClaimId);
            entity.HasIndex(e => e.ApproverId);
        });

        modelBuilder.Entity<ClaimDocument>(entity =>
        {
            entity.ToTable("claim_documents");
            entity.HasKey(e => e.Id);
            entity.Property(e => e.DocumentType).HasMaxLength(100).IsRequired();
            
            entity.HasIndex(e => e.ClaimId);
        });
    }
}
