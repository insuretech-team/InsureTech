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
    public DbSet<FraudCheckResult> FraudCheckResults { get; set; } = null!;

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

            // Money fields — stored as BIGINT (paisa)
            entity.Property(e => e.ClaimedAmount).HasColumnName("claimed_amount").IsRequired();
            entity.Property(e => e.ClaimedCurrency).HasColumnName("claimed_currency").HasMaxLength(3).HasDefaultValue("BDT");
            entity.Property(e => e.ApprovedAmount).HasColumnName("approved_amount");
            entity.Property(e => e.ApprovedCurrency).HasColumnName("approved_currency").HasMaxLength(3).HasDefaultValue("BDT");
            entity.Property(e => e.SettledAmount).HasColumnName("settled_amount");
            entity.Property(e => e.SettledCurrency).HasColumnName("settled_currency").HasMaxLength(3).HasDefaultValue("BDT");
            entity.Property(e => e.DeductibleAmount).HasColumnName("deductible_amount");
            entity.Property(e => e.DeductibleCurrency).HasColumnName("deductible_currency").HasMaxLength(3).HasDefaultValue("BDT");
            entity.Property(e => e.CoPayAmount).HasColumnName("co_pay_amount");
            entity.Property(e => e.CoPayCurrency).HasColumnName("co_pay_currency").HasMaxLength(3).HasDefaultValue("BDT");

            // Proto-aligned new fields
            entity.Property(e => e.BankDetailsForPayout).HasColumnName("bank_details_for_payout");
            entity.Property(e => e.AppealOptionAvailable).HasColumnName("appeal_option_available").HasDefaultValue(false);
            entity.Property(e => e.InAppMessages).HasColumnType("jsonb").HasColumnName("in_app_messages");
            entity.Property(e => e.ProcessorNotes).HasColumnName("processor_notes");

            // Indexes
            entity.HasIndex(e => e.PolicyId);
            entity.HasIndex(e => e.CustomerId);
            entity.HasIndex(e => e.Status);
            entity.HasIndex(e => e.IncidentDate);

            // Ignore Money convenience accessors (not mapped to DB)
            entity.Ignore(e => e.ClaimedMoney);
            entity.Ignore(e => e.ApprovedMoney);
            entity.Ignore(e => e.SettledMoney);
            entity.Ignore(e => e.DeductibleMoney);
            entity.Ignore(e => e.CoPayMoney);

            // Relationships
            entity.HasMany(e => e.Approvals)
                  .WithOne()
                  .HasForeignKey(a => a.ClaimId)
                  .OnDelete(DeleteBehavior.Cascade);

            entity.HasMany(e => e.Documents)
                  .WithOne()
                  .HasForeignKey(d => d.ClaimId)
                  .OnDelete(DeleteBehavior.Cascade);

            entity.HasOne(e => e.FraudCheck)
                  .WithOne()
                  .HasForeignKey<FraudCheckResult>(f => f.ClaimId)
                  .OnDelete(DeleteBehavior.Cascade);
        });

        modelBuilder.Entity<ClaimApproval>(entity =>
        {
            entity.ToTable("claim_approvals");
            entity.HasKey(e => e.Id);
            entity.Property(e => e.Decision).HasConversion<string>().HasMaxLength(50).IsRequired();
            entity.Property(e => e.ApprovedAmount).HasColumnName("approved_amount");
            entity.Property(e => e.ApprovedCurrency).HasColumnName("approved_currency").HasMaxLength(3).HasDefaultValue("BDT");
            entity.Property(e => e.ApproverRole).HasMaxLength(100);

            entity.HasIndex(e => e.ClaimId);
            entity.HasIndex(e => e.ApproverId);
        });

        modelBuilder.Entity<ClaimDocument>(entity =>
        {
            entity.ToTable("claim_documents");
            entity.HasKey(e => e.Id);
            entity.Property(e => e.DocumentType).HasMaxLength(100).IsRequired();
            entity.Property(e => e.FileHash).HasMaxLength(64); // SHA-256 hex

            entity.HasIndex(e => e.ClaimId);
            entity.HasIndex(e => e.FileHash);
        });

        modelBuilder.Entity<FraudCheckResult>(entity =>
        {
            entity.ToTable("fraud_checks");
            entity.HasKey(e => e.Id);
            entity.Property(e => e.FraudScore).HasColumnType("decimal(5,2)");
            entity.Property(e => e.Flagged).HasDefaultValue(false);

            entity.HasIndex(e => e.ClaimId).IsUnique();
            entity.HasIndex(e => e.Flagged);
        });
    }
}
