using InsuranceEngine.Domain.Entities;
using InsuranceEngine.Domain.Enums;
using Microsoft.EntityFrameworkCore;

namespace InsuranceEngine.Infrastructure.Data;

public class InsuranceEngineDbContext(DbContextOptions<InsuranceEngineDbContext> options) : DbContext(options)
{
    public DbSet<Beneficiary> Beneficiaries => Set<Beneficiary>();
    public DbSet<BeneficiaryIndividual> BeneficiaryIndividuals => Set<BeneficiaryIndividual>();
    public DbSet<BeneficiaryBusiness> BeneficiaryBusinesses => Set<BeneficiaryBusiness>();
    public DbSet<Error> Errors => Set<Error>();
    public DbSet<FieldViolation> FieldViolations => Set<FieldViolation>();

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        base.OnModelCreating(modelBuilder);

        modelBuilder.Entity<Beneficiary>(entity =>
        {
            entity.ToTable("beneficiaries");
            entity.HasKey(b => b.BeneficiaryId);
            entity.Property(b => b.BeneficiaryId).HasColumnName("beneficiary_id");
            entity.Property(b => b.UserId).HasColumnName("user_id").HasMaxLength(100).IsRequired();
            entity.Property(b => b.PartnerId).HasColumnName("partner_id").HasMaxLength(100);
            entity.Property(b => b.BeneficiaryCode).HasColumnName("beneficiary_code").HasMaxLength(50).IsRequired();
            entity.Property(b => b.PolicyId).HasColumnName("policy_id");
            entity.Property(b => b.Type).HasColumnName("type").HasMaxLength(32)
                .IsRequired()
                .HasConversion(
                    value => EnumMemberValueMapper.GetEnumMemberValue(value),
                    value => EnumMemberValueMapper.ParseEnumMemberValue<BeneficiaryType>(value));
            entity.Property(b => b.Status).HasColumnName("status").HasMaxLength(32)
                .IsRequired()
                .HasConversion(
                    value => EnumMemberValueMapper.GetEnumMemberValue(value),
                    value => EnumMemberValueMapper.ParseEnumMemberValue<BeneficiaryStatus>(value));
            entity.Property(b => b.KycStatus).HasColumnName("kyc_status").HasMaxLength(32).IsRequired();
            entity.Property(b => b.KycCompletedAt).HasColumnName("kyc_completed_at");
            entity.Property(b => b.RiskScore).HasColumnName("risk_score").HasMaxLength(50);
            entity.Property(b => b.ReferralCode).HasColumnName("referral_code").HasMaxLength(100);
            entity.Property(b => b.ReferredBy).HasColumnName("referred_by").HasMaxLength(100);

            entity.OwnsOne(b => b.AuditInfo, audit =>
            {
                audit.Property(a => a.CreatedAt).HasColumnName("created_at").IsRequired();
                audit.Property(a => a.UpdatedAt).HasColumnName("updated_at").IsRequired();
                audit.Property(a => a.CreatedBy).HasColumnName("created_by").HasMaxLength(100).IsRequired();
                audit.Property(a => a.UpdatedBy).HasColumnName("updated_by").HasMaxLength(100).IsRequired();
                audit.Property(a => a.DeletedAt).HasColumnName("deleted_at");
                audit.Property(a => a.DeletedBy).HasColumnName("deleted_by").HasMaxLength(100);
            });
            entity.Navigation(b => b.AuditInfo).IsRequired();

            entity.HasOne(b => b.IndividualDetails)
                .WithOne(i => i.Beneficiary)
                .HasForeignKey<BeneficiaryIndividual>(i => i.BeneficiaryId)
                .OnDelete(DeleteBehavior.Cascade);

            entity.HasOne(b => b.BusinessDetails)
                .WithOne(i => i.Beneficiary)
                .HasForeignKey<BeneficiaryBusiness>(i => i.BeneficiaryId)
                .OnDelete(DeleteBehavior.Cascade);
        });

        modelBuilder.Entity<BeneficiaryIndividual>(entity =>
        {
            entity.ToTable("beneficiary_individuals");
            entity.HasKey(i => i.Id);
            entity.Property(i => i.Id).HasColumnName("id");
            entity.Property(i => i.BeneficiaryId).HasColumnName("beneficiary_id").IsRequired();
            entity.Property(i => i.FullName).HasColumnName("full_name").HasMaxLength(200).IsRequired();
            entity.Property(i => i.FullNameBn).HasColumnName("full_name_bn").HasMaxLength(200);
            entity.Property(i => i.DateOfBirth).HasColumnName("date_of_birth").IsRequired();
            entity.Property(i => i.Gender).HasColumnName("gender").HasMaxLength(32).IsRequired()
                .HasConversion(
                    value => EnumMemberValueMapper.GetEnumMemberValue(value),
                    value => EnumMemberValueMapper.ParseEnumMemberValue<BeneficiaryGender>(value));
            entity.Property(i => i.NidNumber).HasColumnName("nid_number").HasMaxLength(100);
            entity.Property(i => i.PassportNumber).HasColumnName("passport_number").HasMaxLength(100);
            entity.Property(i => i.BirthCertificateNumber).HasColumnName("birth_certificate_number").HasMaxLength(100);
            entity.Property(i => i.TinNumber).HasColumnName("tin_number").HasMaxLength(100);
            entity.Property(i => i.MaritalStatus).HasColumnName("marital_status").HasMaxLength(32)
                .HasConversion(
                    value => value.HasValue ? EnumMemberValueMapper.GetEnumMemberValue(value.Value) : null,
                    value => string.IsNullOrWhiteSpace(value)
                        ? null
                        : EnumMemberValueMapper.ParseEnumMemberValue<MaritalStatus>(value));
            entity.Property(i => i.Occupation).HasColumnName("occupation").HasMaxLength(100);
            entity.Property(i => i.NomineeName).HasColumnName("nominee_name").HasMaxLength(200);
            entity.Property(i => i.NomineeRelationship).HasColumnName("nominee_relationship").HasMaxLength(200);

            entity.OwnsOne(i => i.ContactInfo, contact =>
            {
                contact.Property(c => c.MobileNumber).HasColumnName("mobile_number").HasMaxLength(50);
                contact.Property(c => c.Email).HasColumnName("email").HasMaxLength(255);
                contact.Property(c => c.AlternateMobile).HasColumnName("alternate_mobile").HasMaxLength(50);
                contact.Property(c => c.Landline).HasColumnName("landline").HasMaxLength(50);
            });
            entity.Navigation(i => i.ContactInfo).IsRequired();

            entity.OwnsOne(i => i.PermanentAddress, address =>
            {
                address.Property(a => a.AddressLine1).HasColumnName("permanent_address_line1").HasMaxLength(255);
                address.Property(a => a.AddressLine2).HasColumnName("permanent_address_line2").HasMaxLength(255);
                address.Property(a => a.City).HasColumnName("permanent_city").HasMaxLength(100);
                address.Property(a => a.District).HasColumnName("permanent_district").HasMaxLength(100);
                address.Property(a => a.Division).HasColumnName("permanent_division").HasMaxLength(100);
                address.Property(a => a.PostalCode).HasColumnName("permanent_postal_code").HasMaxLength(20);
                address.Property(a => a.Country).HasColumnName("permanent_country").HasMaxLength(100);
                address.Property(a => a.Latitude).HasColumnName("permanent_latitude");
                address.Property(a => a.Longitude).HasColumnName("permanent_longitude");
            });
            entity.Navigation(i => i.PermanentAddress).IsRequired();

            entity.OwnsOne(i => i.PresentAddress, address =>
            {
                address.Property(a => a.AddressLine1).HasColumnName("present_address_line1").HasMaxLength(255);
                address.Property(a => a.AddressLine2).HasColumnName("present_address_line2").HasMaxLength(255);
                address.Property(a => a.City).HasColumnName("present_city").HasMaxLength(100);
                address.Property(a => a.District).HasColumnName("present_district").HasMaxLength(100);
                address.Property(a => a.Division).HasColumnName("present_division").HasMaxLength(100);
                address.Property(a => a.PostalCode).HasColumnName("present_postal_code").HasMaxLength(20);
                address.Property(a => a.Country).HasColumnName("present_country").HasMaxLength(100);
                address.Property(a => a.Latitude).HasColumnName("present_latitude");
                address.Property(a => a.Longitude).HasColumnName("present_longitude");
            });
            entity.Navigation(i => i.PresentAddress).IsRequired();

            entity.OwnsOne(i => i.AuditInfo, audit =>
            {
                audit.Property(a => a.CreatedAt).HasColumnName("created_at").IsRequired();
                audit.Property(a => a.UpdatedAt).HasColumnName("updated_at").IsRequired();
                audit.Property(a => a.CreatedBy).HasColumnName("created_by").HasMaxLength(100).IsRequired();
                audit.Property(a => a.UpdatedBy).HasColumnName("updated_by").HasMaxLength(100).IsRequired();
                audit.Property(a => a.DeletedAt).HasColumnName("deleted_at");
                audit.Property(a => a.DeletedBy).HasColumnName("deleted_by").HasMaxLength(100);
            });
            entity.Navigation(i => i.AuditInfo).IsRequired();
        });

        modelBuilder.Entity<BeneficiaryBusiness>(entity =>
        {
            entity.ToTable("beneficiary_businesses");
            entity.HasKey(i => i.Id);
            entity.Property(i => i.Id).HasColumnName("id");
            entity.Property(i => i.BeneficiaryId).HasColumnName("beneficiary_id").IsRequired();
            entity.Property(i => i.BusinessName).HasColumnName("business_name").HasMaxLength(255).IsRequired();
            entity.Property(i => i.BusinessNameBn).HasColumnName("business_name_bn").HasMaxLength(255);
            entity.Property(i => i.TradeLicenseNumber).HasColumnName("trade_license_number").HasMaxLength(100).IsRequired();
            entity.Property(i => i.TradeLicenseIssueDate).HasColumnName("trade_license_issue_date");
            entity.Property(i => i.TradeLicenseExpiryDate).HasColumnName("trade_license_expiry_date");
            entity.Property(i => i.TinNumber).HasColumnName("tin_number").HasMaxLength(100).IsRequired();
            entity.Property(i => i.BinNumber).HasColumnName("bin_number").HasMaxLength(100);
            entity.Property(i => i.BusinessType).HasColumnName("business_type").HasMaxLength(32).IsRequired()
                .HasConversion(
                    value => EnumMemberValueMapper.GetEnumMemberValue(value),
                    value => EnumMemberValueMapper.ParseEnumMemberValue<BusinessType>(value));
            entity.Property(i => i.IndustrySector).HasColumnName("industry_sector").HasMaxLength(150);
            entity.Property(i => i.EmployeeCount).HasColumnName("employee_count");
            entity.Property(i => i.IncorporationDate).HasColumnName("incorporation_date");
            entity.Property(i => i.FocalPersonName).HasColumnName("focal_person_name").HasMaxLength(255).IsRequired();
            entity.Property(i => i.FocalPersonDesignation).HasColumnName("focal_person_designation").HasMaxLength(150);
            entity.Property(i => i.FocalPersonNid).HasColumnName("focal_person_nid").HasMaxLength(100);

            entity.OwnsOne(i => i.ContactInfo, contact =>
            {
                contact.Property(c => c.MobileNumber).HasColumnName("contact_mobile_number").HasMaxLength(50);
                contact.Property(c => c.Email).HasColumnName("contact_email").HasMaxLength(255);
                contact.Property(c => c.AlternateMobile).HasColumnName("contact_alternate_mobile").HasMaxLength(50);
                contact.Property(c => c.Landline).HasColumnName("contact_landline").HasMaxLength(50);
            });
            entity.Navigation(i => i.ContactInfo).IsRequired();

            entity.OwnsOne(i => i.RegisteredAddress, address =>
            {
                address.Property(a => a.AddressLine1).HasColumnName("registered_address_line1").HasMaxLength(255);
                address.Property(a => a.AddressLine2).HasColumnName("registered_address_line2").HasMaxLength(255);
                address.Property(a => a.City).HasColumnName("registered_city").HasMaxLength(100);
                address.Property(a => a.District).HasColumnName("registered_district").HasMaxLength(100);
                address.Property(a => a.Division).HasColumnName("registered_division").HasMaxLength(100);
                address.Property(a => a.PostalCode).HasColumnName("registered_postal_code").HasMaxLength(20);
                address.Property(a => a.Country).HasColumnName("registered_country").HasMaxLength(100);
                address.Property(a => a.Latitude).HasColumnName("registered_latitude");
                address.Property(a => a.Longitude).HasColumnName("registered_longitude");
            });
            entity.Navigation(i => i.RegisteredAddress).IsRequired();

            entity.OwnsOne(i => i.BusinessAddress, address =>
            {
                address.Property(a => a.AddressLine1).HasColumnName("business_address_line1").HasMaxLength(255);
                address.Property(a => a.AddressLine2).HasColumnName("business_address_line2").HasMaxLength(255);
                address.Property(a => a.City).HasColumnName("business_city").HasMaxLength(100);
                address.Property(a => a.District).HasColumnName("business_district").HasMaxLength(100);
                address.Property(a => a.Division).HasColumnName("business_division").HasMaxLength(100);
                address.Property(a => a.PostalCode).HasColumnName("business_postal_code").HasMaxLength(20);
                address.Property(a => a.Country).HasColumnName("business_country").HasMaxLength(100);
                address.Property(a => a.Latitude).HasColumnName("business_latitude");
                address.Property(a => a.Longitude).HasColumnName("business_longitude");
            });
            entity.Navigation(i => i.BusinessAddress).IsRequired();

            entity.OwnsOne(i => i.FocalPersonContact, contact =>
            {
                contact.Property(c => c.MobileNumber).HasColumnName("focal_person_mobile_number").HasMaxLength(50);
                contact.Property(c => c.Email).HasColumnName("focal_person_email").HasMaxLength(255);
                contact.Property(c => c.AlternateMobile).HasColumnName("focal_person_alternate_mobile").HasMaxLength(50);
                contact.Property(c => c.Landline).HasColumnName("focal_person_landline").HasMaxLength(50);
            });
            entity.Navigation(i => i.FocalPersonContact).IsRequired();

            entity.OwnsOne(i => i.AuditInfo, audit =>
            {
                audit.Property(a => a.CreatedAt).HasColumnName("created_at").IsRequired();
                audit.Property(a => a.UpdatedAt).HasColumnName("updated_at").IsRequired();
                audit.Property(a => a.CreatedBy).HasColumnName("created_by").HasMaxLength(100).IsRequired();
                audit.Property(a => a.UpdatedBy).HasColumnName("updated_by").HasMaxLength(100).IsRequired();
                audit.Property(a => a.DeletedAt).HasColumnName("deleted_at");
                audit.Property(a => a.DeletedBy).HasColumnName("deleted_by").HasMaxLength(100);
            });
            entity.Navigation(i => i.AuditInfo).IsRequired();
        });

        modelBuilder.Entity<Error>(entity =>
        {
            entity.ToTable("errors");
            entity.HasKey(e => e.ErrorId);
            entity.Property(e => e.ErrorId).HasColumnName("error_id");
            entity.Property(e => e.Code).HasColumnName("code").HasMaxLength(64)
                .HasConversion(
                    value => EnumMemberValueMapper.GetEnumMemberValue(value),
                    value => EnumMemberValueMapper.ParseEnumMemberValue<ErrorCode>(value));
            entity.Property(e => e.Message).HasColumnName("message").HasMaxLength(500);
            entity.Property(e => e.Details).HasColumnName("details").HasColumnType("jsonb");
            entity.Property(e => e.Retryable).HasColumnName("retryable");
            entity.Property(e => e.RetryAfterSeconds).HasColumnName("retry_after_seconds");
            entity.Property(e => e.HttpStatusCode).HasColumnName("http_status_code");
            entity.Property(e => e.DocumentationUrl).HasColumnName("documentation_url").HasMaxLength(2048);

            entity.HasMany(e => e.FieldViolations)
                .WithOne(v => v.Error)
                .HasForeignKey(v => v.ErrorId)
                .OnDelete(DeleteBehavior.Cascade);
        });

        modelBuilder.Entity<FieldViolation>(entity =>
        {
            entity.ToTable("field_violations");
            entity.HasKey(v => v.FieldViolationId);
            entity.Property(v => v.FieldViolationId).HasColumnName("field_violation_id");
            entity.Property(v => v.ErrorId).HasColumnName("error_id");
            entity.Property(v => v.Field).HasColumnName("field").HasMaxLength(200);
            entity.Property(v => v.Code).HasColumnName("code").HasMaxLength(100);
            entity.Property(v => v.Description).HasColumnName("description").HasMaxLength(1000);
            entity.Property(v => v.RejectedValue).HasColumnName("rejected_value").HasMaxLength(500);
        });
    }
}
