using System.Text.Json;
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

        modelBuilder.Entity<IndividualBeneficiary>(entity =>
        {
            entity.ToTable("individual_beneficiaries");
            entity.HasKey(e => e.Id);
            entity.Property(e => e.Id).HasColumnName("id");
            
            entity.Property(e => e.BeneficiaryId).HasColumnName("beneficiary_id").IsRequired();

            entity.Property(e => e.FullName).HasColumnName("full_name").HasMaxLength(255).IsRequired();
            entity.Property(e => e.FullNameBn).HasColumnName("full_name_bn").HasMaxLength(255);
            entity.Property(e => e.DateOfBirth).HasColumnName("date_of_birth").HasColumnType("date");
            entity.Property(e => e.Gender).HasConversion<string>().HasMaxLength(20);
            entity.Property(e => e.NidNumber).HasColumnName("nid_number").HasMaxLength(50);
            entity.Property(e => e.PassportNumber).HasColumnName("passport_number").HasMaxLength(50);
            entity.Property(e => e.BirthCertificateNumber).HasColumnName("birth_certificate_number").HasMaxLength(50);
            entity.Property(e => e.TinNumber).HasColumnName("tin_number").HasMaxLength(50);
            entity.Property(e => e.MaritalStatus).HasConversion<string>().HasMaxLength(50);
            entity.Property(e => e.Occupation).HasColumnName("occupation").HasMaxLength(100);
            entity.Property(e => e.NomineeName).HasColumnName("nominee_name").HasMaxLength(255);
            entity.Property(e => e.NomineeRelationship).HasColumnName("nominee_relationship").HasMaxLength(100);
            
            entity.Property(e => e.ContactInfo)
                .HasColumnName("contact_info")
                .HasColumnType("jsonb")
                .IsRequired()
                .HasConversion(
                    v => JsonSerializer.Serialize(v, (JsonSerializerOptions)null!),
                    v => JsonSerializer.Deserialize<SharedKernel.Domain.ValueObjects.ContactInfo>(v, (JsonSerializerOptions)null!) ?? new());

            entity.Property(e => e.PermanentAddress)
                .HasColumnName("permanent_address")
                .HasColumnType("jsonb")
                .IsRequired()
                .HasConversion(
                    v => JsonSerializer.Serialize(v, (JsonSerializerOptions)null!),
                    v => JsonSerializer.Deserialize<SharedKernel.Domain.ValueObjects.Address>(v, (JsonSerializerOptions)null!) ?? new());

            entity.Property(e => e.PresentAddress)
                .HasColumnName("present_address")
                .HasColumnType("jsonb")
                .HasConversion(
                    v => JsonSerializer.Serialize(v, (JsonSerializerOptions)null!),
                    v => JsonSerializer.Deserialize<SharedKernel.Domain.ValueObjects.Address>(v, (JsonSerializerOptions)null!) ?? new());

            entity.Property(e => e.AuditInfo)
                .HasColumnName("audit_info")
                .HasColumnType("jsonb")
                .IsRequired()
                .HasConversion(
                    v => JsonSerializer.Serialize(v, (JsonSerializerOptions)null!),
                    v => JsonSerializer.Deserialize<SharedKernel.Domain.ValueObjects.AuditInfo>(v, (JsonSerializerOptions)null!) ?? new());
        });

        modelBuilder.Entity<BusinessBeneficiary>(entity =>
        {
            entity.ToTable("business_beneficiaries");
            entity.HasKey(e => e.Id);
            entity.Property(e => e.Id).HasColumnName("id");

            entity.Property(e => e.BeneficiaryId).HasColumnName("beneficiary_id").IsRequired();

            entity.Property(e => e.BusinessName).HasColumnName("business_name").HasMaxLength(255).IsRequired();
            entity.Property(e => e.BusinessNameBn).HasColumnName("business_name_bn").HasMaxLength(255);
            entity.Property(e => e.TradeLicenseNumber).HasColumnName("trade_license_number").HasMaxLength(100).IsRequired();
            entity.Property(e => e.TradeLicenseIssueDate).HasColumnName("trade_license_issue_date");
            entity.Property(e => e.TradeLicenseExpiryDate).HasColumnName("trade_license_expiry_date");
            entity.Property(e => e.TinNumber).HasColumnName("tin_number").HasMaxLength(50).IsRequired();
            entity.Property(e => e.BinNumber).HasColumnName("bin_number").HasMaxLength(50);
            entity.Property(e => e.RegistrationNumber).HasColumnName("registration_number").HasMaxLength(100);
            entity.Property(e => e.TaxId).HasColumnName("tax_id").HasMaxLength(100);
            entity.Property(e => e.BusinessType).HasConversion<string>().HasMaxLength(100).IsRequired();
            entity.Property(e => e.IndustrySector).HasColumnName("industry_sector").HasMaxLength(100);
            entity.Property(e => e.EmployeeCount).HasColumnName("employee_count");
            entity.Property(e => e.IncorporationDate).HasColumnName("incorporation_date");
            entity.Property(e => e.FocalPersonName).HasColumnName("focal_person_name").HasMaxLength(255).IsRequired();
            entity.Property(e => e.FocalPersonDesignation).HasColumnName("focal_person_designation").HasMaxLength(100);
            entity.Property(e => e.FocalPersonNid).HasColumnName("focal_person_nid").HasMaxLength(50);
            
            entity.Property(e => e.ActivePoliciesCount).HasColumnName("active_policies_count").HasDefaultValue(0);
            entity.Property(e => e.PendingActionsCount).HasColumnName("pending_actions_count").HasDefaultValue(0);
            entity.Property(e => e.TotalEmployeesCovered).HasColumnName("total_employees_covered").HasDefaultValue(0);

            entity.OwnsOne(e => e.TotalPremiumAmount, money =>
            {
                money.Property(m => m.Amount).HasColumnName("total_premium_amount").HasDefaultValue(0);
                money.Property(m => m.CurrencyCode).HasColumnName("total_premium_currency").HasMaxLength(3).HasDefaultValue("BDT");
            });

            entity.Property(e => e.FocalPersonContact)
                .HasColumnName("focal_person_contact")
                .HasColumnType("jsonb")
                .IsRequired()
                .HasConversion(
                    v => JsonSerializer.Serialize(v, (JsonSerializerOptions)null!),
                    v => JsonSerializer.Deserialize<SharedKernel.Domain.ValueObjects.ContactInfo>(v, (JsonSerializerOptions)null!) ?? new());

            entity.Property(e => e.PrimaryContact)
                .HasColumnName("primary_contact")
                .HasColumnType("jsonb")
                .HasConversion(
                    v => JsonSerializer.Serialize(v, (JsonSerializerOptions)null!),
                    v => JsonSerializer.Deserialize<SharedKernel.Domain.ValueObjects.PrimaryContact>(v, (JsonSerializerOptions)null!) ?? new());

            entity.Property(e => e.ContactInfo)
                .HasColumnName("contact_info")
                .HasColumnType("jsonb")
                .HasConversion(
                    v => JsonSerializer.Serialize(v, (JsonSerializerOptions)null!),
                    v => JsonSerializer.Deserialize<SharedKernel.Domain.ValueObjects.ContactInfo>(v, (JsonSerializerOptions)null!) ?? new());

            entity.Property(e => e.RegisteredAddress)
                .HasColumnName("registered_address")
                .HasColumnType("jsonb")
                .IsRequired()
                .HasConversion(
                    v => JsonSerializer.Serialize(v, (JsonSerializerOptions)null!),
                    v => JsonSerializer.Deserialize<SharedKernel.Domain.ValueObjects.Address>(v, (JsonSerializerOptions)null!) ?? new());

            entity.Property(e => e.BusinessAddress)
                .HasColumnName("business_address")
                .HasColumnType("jsonb")
                .IsRequired()
                .HasConversion(
                    v => JsonSerializer.Serialize(v, (JsonSerializerOptions)null!),
                    v => JsonSerializer.Deserialize<SharedKernel.Domain.ValueObjects.Address>(v, (JsonSerializerOptions)null!) ?? new());

            entity.Property(e => e.AuditInfo)
                .HasColumnName("audit_info")
                .HasColumnType("jsonb")
                .IsRequired()
                .HasConversion(
                    v => JsonSerializer.Serialize(v, (JsonSerializerOptions)null!),
                    v => JsonSerializer.Deserialize<SharedKernel.Domain.ValueObjects.AuditInfo>(v, (JsonSerializerOptions)null!) ?? new());
        });

        modelBuilder.Entity<Domain.Entities.Beneficiary>(entity =>
        {
            entity.ToTable("beneficiaries");
            entity.HasKey(e => e.Id);
            entity.Property(e => e.Id).HasColumnName("beneficiary_id");
            entity.Property(e => e.UserId).HasColumnName("user_id").IsRequired();
            entity.Property(e => e.Code).HasColumnName("code").HasMaxLength(20).IsRequired();
            entity.HasIndex(e => e.Code).IsUnique();
            
            entity.Property(e => e.Type).HasConversion<string>().HasMaxLength(20).IsRequired();
            entity.Property(e => e.KycCompletedAt).HasColumnName("kyc_completed_at");
            entity.Property(e => e.RiskScore).HasColumnName("risk_score").HasMaxLength(20);
            
            entity.Property(e => e.Status)
                .HasColumnType("jsonb")
                .HasConversion(
                    v => JsonSerializer.Serialize(v, (JsonSerializerOptions)null!),
                    v => JsonSerializer.Deserialize<BeneficiaryStatusInfo>(v, (JsonSerializerOptions)null!) ?? new());

            entity.Property(e => e.KycStatus)
                .HasColumnType("jsonb")
                .HasConversion(
                    v => JsonSerializer.Serialize(v, (JsonSerializerOptions)null!),
                    v => JsonSerializer.Deserialize<KYCStatusInfo>(v, (JsonSerializerOptions)null!) ?? new());

            entity.Property(e => e.AuditInfo)
                .HasColumnName("audit_info")
                .HasColumnType("jsonb")
                .HasConversion(
                    v => JsonSerializer.Serialize(v, (JsonSerializerOptions)null!),
                    v => JsonSerializer.Deserialize<InsuranceEngine.SharedKernel.Domain.ValueObjects.AuditInfo>(v, (JsonSerializerOptions)null!) ?? new());

            // Explicitly ignore domain events
            entity.Ignore(e => e.DomainEvents);

            entity.Property(e => e.ReferralCode).HasColumnName("referral_code").HasMaxLength(50);
            entity.Property(e => e.ReferredBy).HasColumnName("referred_by");
            entity.Property(e => e.PartnerId).HasColumnName("partner_id");
            
            entity.HasOne(e => e.Individual)
                .WithOne()
                .HasForeignKey<IndividualBeneficiary>(e => e.BeneficiaryId)
                .OnDelete(DeleteBehavior.Cascade);

            entity.HasOne(e => e.Business)
                .WithOne()
                .HasForeignKey<BusinessBeneficiary>(e => e.BeneficiaryId)
                .OnDelete(DeleteBehavior.Cascade);
        });
    }
}
