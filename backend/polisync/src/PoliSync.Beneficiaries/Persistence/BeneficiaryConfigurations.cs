using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Metadata.Builders;
using PoliSync.Beneficiaries.Domain;

namespace PoliSync.Beneficiaries.Persistence;

public class BeneficiaryConfiguration : IEntityTypeConfiguration<Beneficiary>
{
    public void Configure(EntityTypeBuilder<Beneficiary> builder)
    {
        builder.ToTable("beneficiaries", "insurance_schema");
        builder.HasKey(b => b.BeneficiaryId);

        builder.Property(b => b.BeneficiaryId).HasColumnName("beneficiary_id");
        builder.Property(b => b.UserId).HasColumnName("user_id").IsRequired();
        builder.Property(b => b.Type).HasColumnName("type").HasMaxLength(20).IsRequired()
            .HasConversion<string>();
        builder.Property(b => b.Code).HasColumnName("code").HasMaxLength(20).IsRequired();
        builder.Property(b => b.Status).HasColumnName("status").HasMaxLength(20).IsRequired()
            .HasConversion<string>();
        builder.Property(b => b.KycStatus).HasColumnName("kyc_status").HasMaxLength(20).IsRequired()
            .HasConversion<string>();
        builder.Property(b => b.KycCompletedAt).HasColumnName("kyc_completed_at");
        builder.Property(b => b.RiskScore).HasColumnName("risk_score").HasMaxLength(10);
        builder.Property(b => b.ReferralCode).HasColumnName("referral_code").HasMaxLength(20);
        builder.Property(b => b.ReferredBy).HasColumnName("referred_by");
        builder.Property(b => b.PartnerId).HasColumnName("partner_id");
        builder.Property(b => b.AuditInfo).HasColumnName("audit_info").HasColumnType("jsonb");

        // Ignore domain events
        builder.Ignore(b => b.DomainEvents);

        // One-to-one with details
        builder.HasOne(b => b.IndividualDetails)
            .WithOne()
            .HasForeignKey<IndividualBeneficiary>(ib => ib.BeneficiaryId);

        builder.HasOne(b => b.BusinessDetails)
            .WithOne()
            .HasForeignKey<BusinessBeneficiary>(bb => bb.BeneficiaryId);
    }
}

public class IndividualBeneficiaryConfiguration : IEntityTypeConfiguration<IndividualBeneficiary>
{
    public void Configure(EntityTypeBuilder<IndividualBeneficiary> builder)
    {
        builder.ToTable("individual_beneficiaries", "insurance_schema");
        builder.HasKey(ib => ib.BeneficiaryId);

        builder.Property(ib => ib.BeneficiaryId).HasColumnName("beneficiary_id");
        builder.Property(ib => ib.FullName).HasColumnName("full_name").HasMaxLength(255).IsRequired();
        builder.Property(ib => ib.FullNameBn).HasColumnName("full_name_bn").HasMaxLength(255);
        builder.Property(ib => ib.DateOfBirth).HasColumnName("date_of_birth").IsRequired();
        builder.Property(ib => ib.Gender).HasColumnName("gender").HasMaxLength(10).IsRequired()
            .HasConversion<string>();
        builder.Property(ib => ib.NidNumber).HasColumnName("nid_number").HasMaxLength(17);
        builder.Property(ib => ib.PassportNumber).HasColumnName("passport_number").HasMaxLength(20);
        builder.Property(ib => ib.BirthCertificateNumber).HasColumnName("birth_certificate_number").HasMaxLength(20);
        builder.Property(ib => ib.TinNumber).HasColumnName("tin_number").HasMaxLength(12);
        builder.Property(ib => ib.MaritalStatus).HasColumnName("marital_status").HasMaxLength(20)
            .HasConversion<string>();
        builder.Property(ib => ib.Occupation).HasColumnName("occupation").HasMaxLength(100);
        builder.Property(ib => ib.ContactInfo).HasColumnName("contact_info").HasColumnType("jsonb");
        builder.Property(ib => ib.PermanentAddress).HasColumnName("permanent_address").HasColumnType("jsonb");
        builder.Property(ib => ib.PresentAddress).HasColumnName("present_address").HasColumnType("jsonb");
        builder.Property(ib => ib.NomineeName).HasColumnName("nominee_name").HasMaxLength(255);
        builder.Property(ib => ib.NomineeRelationship).HasColumnName("nominee_relationship").HasMaxLength(50);
        builder.Property(ib => ib.AuditInfo).HasColumnName("audit_info").HasColumnType("jsonb");

        builder.Ignore(ib => ib.DomainEvents);
    }
}

public class BusinessBeneficiaryConfiguration : IEntityTypeConfiguration<BusinessBeneficiary>
{
    public void Configure(EntityTypeBuilder<BusinessBeneficiary> builder)
    {
        builder.ToTable("business_beneficiaries", "insurance_schema");
        builder.HasKey(bb => bb.BeneficiaryId);

        builder.Property(bb => bb.BeneficiaryId).HasColumnName("beneficiary_id");
        builder.Property(bb => bb.BusinessName).HasColumnName("business_name").HasMaxLength(255).IsRequired();
        builder.Property(bb => bb.TradeLicenseNumber).HasColumnName("trade_license_number").HasMaxLength(50);
        builder.Property(bb => bb.TinNumber).HasColumnName("tin_number").HasMaxLength(20);
        builder.Property(bb => bb.BusinessType).HasColumnName("business_type").HasMaxLength(50);
        builder.Property(bb => bb.FocalPersonName).HasColumnName("focal_person_name").HasMaxLength(100);
        builder.Property(bb => bb.FocalPersonMobile).HasColumnName("focal_person_mobile").HasMaxLength(20);
        builder.Property(bb => bb.PartnerId).HasColumnName("partner_id");
        builder.Property(bb => bb.AuditInfo).HasColumnName("audit_info").HasColumnType("jsonb");

        builder.Ignore(bb => bb.DomainEvents);
    }
}
