using InsuranceEngine.Beneficiary.Domain.Entities;
using InsuranceEngine.SharedKernel.Domain.Events;
using Microsoft.EntityFrameworkCore;

namespace InsuranceEngine.Beneficiary.Infrastructure.Persistence;

public class BeneficiaryDbContext : DbContext
{
    public BeneficiaryDbContext(DbContextOptions<BeneficiaryDbContext> options) : base(options)
    {
    }

    public DbSet<Domain.Entities.Beneficiary> Beneficiaries { get; set; } = null!;
    public DbSet<IndividualBeneficiary> IndividualBeneficiaries { get; set; } = null!;
    public DbSet<BusinessBeneficiary> BusinessBeneficiaries { get; set; } = null!;

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        modelBuilder.HasDefaultSchema("insurance_schema");

        modelBuilder.Entity<Domain.Entities.Beneficiary>(entity =>
        {
            entity.ToTable("beneficiaries");
            entity.HasKey(e => e.Id);
            entity.Property(e => e.Id).HasColumnName("beneficiary_id");
            entity.Property(e => e.Code).HasMaxLength(50).IsRequired();
            entity.HasIndex(e => e.Code).IsUnique();
            entity.Property(e => e.Status).HasConversion<string>().HasMaxLength(50).IsRequired();
            entity.Property(e => e.Type).HasConversion<string>().HasMaxLength(50).IsRequired();
            entity.Property(e => e.KycStatus).HasConversion<string>().HasMaxLength(50).IsRequired();
            
            // Explicitly ignore domain events to avoid mapping issues
            entity.Ignore(e => e.DomainEvents);

            entity.Property(e => e.ReferralCode).HasColumnName("referral_code").HasMaxLength(20);
            entity.Property(e => e.ReferredBy).HasColumnName("referred_by");
            entity.Property(e => e.AuditInfoJson).HasColumnName("audit_info").HasColumnType("jsonb");

            entity.HasOne(e => e.IndividualDetails)
                .WithOne()
                .HasForeignKey<IndividualBeneficiary>(e => e.Id)
                .OnDelete(DeleteBehavior.Cascade);

            entity.HasOne(e => e.BusinessDetails)
                .WithOne()
                .HasForeignKey<BusinessBeneficiary>(e => e.Id)
                .OnDelete(DeleteBehavior.Cascade);
        });

        modelBuilder.Entity<IndividualBeneficiary>(entity =>
        {
            entity.ToTable("individual_beneficiaries");
            entity.HasKey(e => e.Id);
            entity.Property(e => e.Id).HasColumnName("beneficiary_id");
            
            entity.Ignore(e => e.BeneficiaryId);

            entity.Property(e => e.Gender).HasConversion<string>().HasMaxLength(20);
            entity.Property(e => e.MaritalStatus).HasConversion<string>().HasMaxLength(20);
            
            entity.Property(e => e.ContactInfoJson).HasColumnName("contact_info").HasColumnType("jsonb");
            entity.Property(e => e.PermanentAddressJson).HasColumnName("permanent_address").HasColumnType("jsonb");
            entity.Property(e => e.PresentAddressJson).HasColumnName("present_address").HasColumnType("jsonb");
            entity.Property(e => e.AuditInfoJson).HasColumnName("audit_info").HasColumnType("jsonb");
        });

        modelBuilder.Entity<BusinessBeneficiary>(entity =>
        {
            entity.ToTable("business_beneficiaries");
            entity.HasKey(e => e.Id);
            entity.Property(e => e.Id).HasColumnName("beneficiary_id");

            entity.Ignore(e => e.BeneficiaryId);

            entity.Property(e => e.BusinessType).HasConversion<string>().HasMaxLength(50);
            
            entity.Property(e => e.ContactInfoJson).HasColumnName("contact_info").HasColumnType("jsonb");
            entity.Property(e => e.RegisteredAddressJson).HasColumnName("registered_address").HasColumnType("jsonb");
            entity.Property(e => e.BusinessAddressJson).HasColumnName("business_address").HasColumnType("jsonb");
            entity.Property(e => e.AuditInfoJson).HasColumnName("audit_info").HasColumnType("jsonb");
        });
    }
}
