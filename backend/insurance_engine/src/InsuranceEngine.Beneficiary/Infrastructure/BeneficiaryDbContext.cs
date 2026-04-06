using Microsoft.EntityFrameworkCore;
using InsuranceEngine.SharedKernel.Persistence.Entities;

namespace InsuranceEngine.Beneficiary.Infrastructure;

public class BeneficiaryDbContext : DbContext
{
    public BeneficiaryDbContext(DbContextOptions<BeneficiaryDbContext> options) : base(options) { }

    public DbSet<BeneficiaryEntity> Beneficiaries => Set<BeneficiaryEntity>();
    public DbSet<IndividualBeneficiaryEntity> IndividualBeneficiaries => Set<IndividualBeneficiaryEntity>();
    public DbSet<BusinessBeneficiaryEntity> BusinessBeneficiaries => Set<BusinessBeneficiaryEntity>();

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        modelBuilder.HasDefaultSchema("insurance_schema");

        modelBuilder.Entity<BeneficiaryEntity>(entity =>
        {
            entity.ToTable("beneficiaries");
            entity.HasKey(e => e.BeneficiaryId);
            entity.Property(e => e.BeneficiaryId).HasColumnName("beneficiary_id").HasDefaultValueSql("gen_random_uuid()");
            entity.Property(e => e.UserId).HasColumnName("user_id").IsRequired();
            entity.Property(e => e.Type).HasColumnName("type").HasMaxLength(50).IsRequired();
            entity.Property(e => e.Code).HasColumnName("code").HasMaxLength(50).IsRequired();
            entity.Property(e => e.Status).HasColumnName("status").HasMaxLength(50).HasDefaultValue("PENDING_KYC");
            entity.Property(e => e.KycStatus).HasColumnName("kyc_status").HasMaxLength(50).HasDefaultValue("NOT_STARTED");
            entity.Property(e => e.KycCompletedAt).HasColumnName("kyc_completed_at");
            entity.Property(e => e.RiskScore).HasColumnName("risk_score").HasMaxLength(20);
            entity.Property(e => e.ReferralCode).HasColumnName("referral_code").HasMaxLength(50);
            entity.Property(e => e.ReferredBy).HasColumnName("referred_by");
            entity.Property(e => e.PartnerId).HasColumnName("partner_id");
            entity.Property(e => e.AuditInfo).HasColumnName("audit_info");
            entity.Property(e => e.CreatedAt).HasColumnName("created_at").HasDefaultValueSql("CURRENT_TIMESTAMP").IsRequired();
            entity.Property(e => e.UpdatedAt).HasColumnName("updated_at").HasDefaultValueSql("CURRENT_TIMESTAMP").IsRequired();
            entity.Property(e => e.DeletedAt).HasColumnName("deleted_at");

            entity.HasIndex(e => e.Code).IsUnique().HasDatabaseName("idx_beneficiaries_code");
            entity.HasIndex(e => e.UserId).HasDatabaseName("idx_beneficiaries_user_id");
            entity.HasIndex(e => e.Status).HasDatabaseName("idx_beneficiaries_status");
        });

        modelBuilder.Entity<IndividualBeneficiaryEntity>(entity =>
        {
            entity.ToTable("individual_beneficiaries");
            entity.HasKey(e => e.BeneficiaryId);
            entity.Property(e => e.BeneficiaryId).HasColumnName("beneficiary_id").IsRequired();
            entity.Property(e => e.FullName).HasColumnName("full_name").HasMaxLength(200).IsRequired();
            entity.Property(e => e.FullNameBn).HasColumnName("full_name_bn").HasMaxLength(200);
            entity.Property(e => e.DateOfBirth).HasColumnName("date_of_birth").IsRequired();
            entity.Property(e => e.Gender).HasColumnName("gender").HasMaxLength(20).IsRequired();
            entity.Property(e => e.NidNumber).HasColumnName("nid_number").HasMaxLength(50);
            entity.Property(e => e.PassportNumber).HasColumnName("passport_number").HasMaxLength(50);
            entity.Property(e => e.BirthCertificateNumber).HasColumnName("birth_certificate_number").HasMaxLength(50);
            entity.Property(e => e.TinNumber).HasColumnName("tin_number").HasMaxLength(50);
            entity.Property(e => e.MaritalStatus).HasColumnName("marital_status").HasMaxLength(20);
            entity.Property(e => e.Occupation).HasColumnName("occupation").HasMaxLength(100);
            entity.Property(e => e.ContactInfo).HasColumnName("contact_info");
            entity.Property(e => e.PermanentAddress).HasColumnName("permanent_address");
            entity.Property(e => e.PresentAddress).HasColumnName("present_address");
            entity.Property(e => e.NomineeName).HasColumnName("nominee_name").HasMaxLength(200);
            entity.Property(e => e.NomineeRelationship).HasColumnName("nominee_relationship").HasMaxLength(50);
            entity.Property(e => e.AuditInfo).HasColumnName("audit_info");
            entity.Property(e => e.CreatedAt).HasColumnName("created_at").HasDefaultValueSql("CURRENT_TIMESTAMP").IsRequired();
            entity.Property(e => e.UpdatedAt).HasColumnName("updated_at").HasDefaultValueSql("CURRENT_TIMESTAMP").IsRequired();

            entity.HasOne(e => e.Beneficiary).WithOne(b => b.IndividualDetails).HasForeignKey<IndividualBeneficiaryEntity>(e => e.BeneficiaryId).OnDelete(DeleteBehavior.Cascade);
        });

        modelBuilder.Entity<BusinessBeneficiaryEntity>(entity =>
        {
            entity.ToTable("business_beneficiaries");
            entity.HasKey(e => e.BeneficiaryId);
            entity.Property(e => e.BeneficiaryId).HasColumnName("beneficiary_id").IsRequired();
            entity.Property(e => e.ParentBeneficiaryId).HasColumnName("parent_beneficiary_id").IsRequired();
            entity.Property(e => e.BusinessName).HasColumnName("business_name").HasMaxLength(200).IsRequired();
            entity.Property(e => e.BusinessNameBn).HasColumnName("business_name_bn").HasMaxLength(200);
            entity.Property(e => e.TradeLicenseNumber).HasColumnName("trade_license_number").HasMaxLength(100).IsRequired();
            entity.Property(e => e.TradeLicenseIssueDate).HasColumnName("trade_license_issue_date");
            entity.Property(e => e.TradeLicenseExpiryDate).HasColumnName("trade_license_expiry_date");
            entity.Property(e => e.TinNumber).HasColumnName("tin_number").HasMaxLength(50).IsRequired();
            entity.Property(e => e.BinNumber).HasColumnName("bin_number").HasMaxLength(50);
            entity.Property(e => e.BusinessType).HasColumnName("business_type").HasMaxLength(50).HasDefaultValue("CORPORATE");
            entity.Property(e => e.IndustrySector).HasColumnName("industry_sector").HasMaxLength(100);
            entity.Property(e => e.EmployeeCount).HasColumnName("employee_count");
            entity.Property(e => e.IncorporationDate).HasColumnName("incorporation_date");
            entity.Property(e => e.ContactInfo).HasColumnName("contact_info");
            entity.Property(e => e.RegisteredAddress).HasColumnName("registered_address");
            entity.Property(e => e.BusinessAddress).HasColumnName("business_address");
            entity.Property(e => e.FocalPersonName).HasColumnName("focal_person_name").HasMaxLength(200).IsRequired();
            entity.Property(e => e.FocalPersonDesignation).HasColumnName("focal_person_designation").HasMaxLength(100);
            entity.Property(e => e.FocalPersonNid).HasColumnName("focal_person_nid").HasMaxLength(50);
            entity.Property(e => e.FocalPersonContact).HasColumnName("focal_person_contact");
            entity.Property(e => e.RegistrationNumber).HasColumnName("registration_number").HasMaxLength(100);
            entity.Property(e => e.TaxId).HasColumnName("tax_id").HasMaxLength(50);
            entity.Property(e => e.PrimaryContact).HasColumnName("primary_contact");
            entity.Property(e => e.TotalEmployeesCovered).HasColumnName("total_employees_covered").HasDefaultValue(0);
            entity.Property(e => e.ActivePoliciesCount).HasColumnName("active_policies_count").HasDefaultValue(0);
            entity.Property(e => e.TotalPremiumAmount).HasColumnName("total_premium_amount").HasDefaultValue(0);
            entity.Property(e => e.PendingActionsCount).HasColumnName("pending_actions_count").HasDefaultValue(0);
            entity.Property(e => e.AuditInfo).HasColumnName("audit_info");
            entity.Property(e => e.CreatedAt).HasColumnName("created_at").HasDefaultValueSql("CURRENT_TIMESTAMP").IsRequired();
            entity.Property(e => e.UpdatedAt).HasColumnName("updated_at").HasDefaultValueSql("CURRENT_TIMESTAMP").IsRequired();

            entity.HasOne(e => e.Beneficiary).WithOne(b => b.BusinessDetails).HasForeignKey<BusinessBeneficiaryEntity>(e => e.BeneficiaryId).OnDelete(DeleteBehavior.Cascade);
        });
    }
}
