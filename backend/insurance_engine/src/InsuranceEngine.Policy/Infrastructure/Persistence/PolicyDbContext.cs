using InsuranceEngine.Policy.Domain.Entities;
using InsuranceEngine.Policy.Domain.Enums;
using InsuranceEngine.SharedKernel.Interfaces;
using Microsoft.EntityFrameworkCore;
using System;

namespace InsuranceEngine.Policy.Infrastructure.Persistence;

public class PolicyDbContext : DbContext
{
    private readonly Guid _tenantId;

    public PolicyDbContext(DbContextOptions<PolicyDbContext> options, ITenantService tenantService) : base(options)
    {
        _tenantId = tenantService.GetTenantId();
    }

    public DbSet<PolicyEntity> Policies { get; set; } = null!;
    public DbSet<Nominee> Nominees { get; set; } = null!;
    public DbSet<PolicyRider> PolicyRiders { get; set; } = null!;
    public DbSet<Quote> Quotes { get; set; } = null!;
    public DbSet<Beneficiary> Beneficiaries { get; set; } = null!;
    public DbSet<IndividualBeneficiary> IndividualBeneficiaries { get; set; } = null!;
    public DbSet<BusinessBeneficiary> BusinessBeneficiaries { get; set; } = null!;
    public DbSet<UnderwritingHealthDeclaration> HealthDeclarations { get; set; } = null!;
    public DbSet<UnderwritingDecision> UnderwritingDecisions { get; set; } = null!;
    public DbSet<Claim> Claims { get; set; } = null!;
    public DbSet<ClaimApproval> ClaimApprovals { get; set; } = null!;
    public DbSet<ClaimDocument> ClaimDocuments { get; set; } = null!;

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        modelBuilder.HasDefaultSchema("insurance_schema");

        // --- Policy ---
        modelBuilder.Entity<PolicyEntity>(entity =>
        {
            entity.ToTable("policies");
            entity.HasKey(e => e.Id);
            entity.HasQueryFilter(e => !e.IsDeleted);

            entity.Property(e => e.PolicyNumber).HasMaxLength(50).IsRequired();
            entity.HasIndex(e => e.PolicyNumber).IsUnique();

            entity.Property(e => e.Status).HasConversion<string>().HasMaxLength(50).IsRequired();

            // Money columns as bigint
            entity.Property(e => e.PremiumAmount).HasColumnName("premium_amount").IsRequired();
            entity.Property(e => e.PremiumCurrency).HasColumnName("premium_currency").HasMaxLength(3).HasDefaultValue("BDT");
            entity.Property(e => e.SumInsuredAmount).HasColumnName("sum_insured_amount").IsRequired();
            entity.Property(e => e.SumInsuredCurrency).HasColumnName("sum_insured_currency").HasMaxLength(3).HasDefaultValue("BDT");
            entity.Property(e => e.VatTaxAmount).HasColumnName("vat_tax_amount");
            entity.Property(e => e.ServiceFeeAmount).HasColumnName("service_fee_amount");
            entity.Property(e => e.TotalPayableAmount).HasColumnName("total_payable_amount");

            entity.Property(e => e.ProposerDetailsJson).HasColumnType("jsonb").HasColumnName("proposer_details");
            entity.Property(e => e.UnderwritingData).HasColumnType("jsonb").HasColumnName("underwriting_data");

            entity.Ignore(e => e.PremiumMoney);
            entity.Ignore(e => e.SumInsuredMoney);

            entity.HasIndex(e => e.CustomerId);
            entity.HasIndex(e => e.ProductId);
            entity.HasIndex(e => e.Status);

            entity.HasMany(e => e.Nominees)
                  .WithOne()
                  .HasForeignKey(n => n.PolicyId)
                  .OnDelete(DeleteBehavior.Cascade);

            entity.HasMany(e => e.Riders)
                  .WithOne()
                  .HasForeignKey(r => r.PolicyId)
                  .OnDelete(DeleteBehavior.Cascade);
        });


        // --- PolicyRider ---
        modelBuilder.Entity<PolicyRider>(entity =>
        {
            entity.ToTable("policy_riders");
            entity.HasKey(e => e.Id);

            entity.Property(e => e.RiderName).HasMaxLength(255).IsRequired();
            entity.Property(e => e.PremiumAmount).HasColumnName("premium_amount").IsRequired();
            entity.Property(e => e.PremiumCurrency).HasColumnName("premium_currency").HasMaxLength(3).HasDefaultValue("BDT");
            entity.Property(e => e.CoverageAmount).HasColumnName("coverage_amount").IsRequired();
            entity.Property(e => e.CoverageCurrency).HasColumnName("coverage_currency").HasMaxLength(3).HasDefaultValue("BDT");

            entity.Ignore(e => e.Premium);
            entity.Ignore(e => e.Coverage);

            entity.HasIndex(e => e.PolicyId);
        });

        // --- Quote ---
        modelBuilder.Entity<Quote>(entity =>
        {
            entity.ToTable("quotes");
            entity.HasKey(e => e.Id);
            entity.HasQueryFilter(e => !e.IsDeleted);
            entity.Property(e => e.QuoteNumber).HasMaxLength(50).IsRequired();
            entity.HasIndex(e => e.QuoteNumber).IsUnique();
            entity.Property(e => e.Status).HasConversion<string>().HasMaxLength(50).IsRequired();
            
            entity.Property(e => e.SumAssuredAmount).HasColumnName("sum_assured_amount").IsRequired();
            entity.Property(e => e.SumAssuredCurrency).HasColumnName("sum_assured_currency").HasMaxLength(3).HasDefaultValue("BDT");
            entity.Property(e => e.BasePremiumAmount).HasColumnName("base_premium_amount").IsRequired();
            entity.Property(e => e.RiderPremiumAmount).HasColumnName("rider_premium_amount");
            entity.Property(e => e.TaxAmount).HasColumnName("tax_amount");
            entity.Property(e => e.TotalPremiumAmount).HasColumnName("total_premium_amount").IsRequired();
            entity.Property(e => e.Currency).HasColumnName("currency").HasMaxLength(3).HasDefaultValue("BDT");

            entity.Property(e => e.PremiumCalculationJson).HasColumnType("jsonb").HasColumnName("premium_calculation");
            entity.Property(e => e.SelectedRidersJson).HasColumnType("jsonb").HasColumnName("selected_riders");
            entity.Property(e => e.AuditInfoJson).HasColumnType("jsonb").HasColumnName("audit_info");

            entity.Ignore(e => e.SumAssured);
            entity.Ignore(e => e.BasePremium);
            entity.Ignore(e => e.TotalPremium);

            entity.HasIndex(e => e.BeneficiaryId);
            entity.HasIndex(e => e.InsurerProductId);
            entity.HasIndex(e => e.Status);
        });

        // --- UnderwritingHealthDeclaration ---
        modelBuilder.Entity<Beneficiary>(entity =>
        {
            entity.ToTable("beneficiaries", "insurance_schema");
            entity.HasKey(e => e.Id);
            entity.Property(e => e.Code).IsRequired().HasMaxLength(20);
            entity.Property(e => e.Type).HasConversion<string>();
            entity.Property(e => e.Status).HasConversion<string>();
            entity.Property(e => e.KycStatus).HasConversion<string>();
            entity.Property(e => e.AuditInfo).HasColumnType("jsonb");
            entity.HasQueryFilter(e => !e.IsDeleted);

            entity.HasOne(e => e.IndividualDetails)
                .WithOne(e => e.Beneficiary)
                .HasForeignKey<IndividualBeneficiary>(e => e.BeneficiaryId);

            entity.HasOne(e => e.BusinessDetails)
                .WithOne(e => e.Beneficiary)
                .HasForeignKey<BusinessBeneficiary>(e => e.BeneficiaryId);
        });

        modelBuilder.Entity<IndividualBeneficiary>(entity =>
        {
            entity.ToTable("individual_beneficiaries", "insurance_schema");
            entity.HasKey(e => e.BeneficiaryId);
            entity.Property(e => e.Gender).HasConversion<string>();
            entity.Property(e => e.MaritalStatus).HasConversion<string>();
            entity.Property(e => e.ContactInfoJson).HasColumnType("jsonb");
            entity.Property(e => e.PermanentAddressJson).HasColumnType("jsonb");
            entity.Property(e => e.PresentAddressJson).HasColumnType("jsonb");
            entity.Property(e => e.AuditInfo).HasColumnType("jsonb");
        });

        modelBuilder.Entity<BusinessBeneficiary>(entity =>
        {
            entity.ToTable("business_beneficiaries", "insurance_schema");
            entity.HasKey(e => e.Id);
            entity.Property(e => e.BusinessType).HasConversion<string>();
            entity.Property(e => e.ContactInfoJson).HasColumnType("jsonb");
            entity.Property(e => e.RegisteredAddressJson).HasColumnType("jsonb");
            entity.Property(e => e.BusinessAddressJson).HasColumnType("jsonb");
            entity.Property(e => e.FocalPersonContactJson).HasColumnType("jsonb");
            entity.Property(e => e.AuditInfo).HasColumnType("jsonb");
            entity.Property(e => e.PrimaryContactJson).HasColumnType("jsonb");
        });

        modelBuilder.Entity<Nominee>(entity =>
        {
            entity.ToTable("policy_nominees", "insurance_schema");
            entity.HasKey(e => e.Id);
            entity.Property(e => e.Relationship).IsRequired().HasMaxLength(50);
            entity.HasQueryFilter(e => !e.IsDeleted);

            entity.HasOne(e => e.Beneficiary)
                .WithMany()
                .HasForeignKey(e => e.BeneficiaryId);
        });

        modelBuilder.Entity<UnderwritingHealthDeclaration>(entity =>
        {
            entity.ToTable("health_declarations");
            entity.HasKey(e => e.Id);
            entity.Property(e => e.WeightKg).HasPrecision(5, 2);
            entity.Property(e => e.Bmi).HasPrecision(5, 2);
            entity.Property(e => e.PreExistingConditionsJson).HasColumnType("jsonb").HasColumnName("pre_existing_conditions");
            entity.Property(e => e.FamilyHistoryJson).HasColumnType("jsonb").HasColumnName("family_history");
            entity.Property(e => e.MedicalExamResultsJson).HasColumnType("jsonb").HasColumnName("medical_exam_results");
            entity.Property(e => e.MedicalDocumentsJson).HasColumnType("jsonb").HasColumnName("medical_documents");
            entity.Property(e => e.AuditInfoJson).HasColumnType("jsonb").HasColumnName("audit_info");

            entity.HasIndex(e => e.QuoteId).IsUnique();
        });

        // --- UnderwritingDecision ---
        modelBuilder.Entity<UnderwritingDecision>(entity =>
        {
            entity.ToTable("underwriting_decisions");
            entity.HasKey(e => e.Id);
            entity.Property(e => e.Decision).HasConversion<string>().HasMaxLength(50).IsRequired();
            entity.Property(e => e.Method).HasConversion<string>().HasMaxLength(50).IsRequired();
            entity.Property(e => e.RiskLevel).HasConversion<string>().HasMaxLength(50);
            entity.Property(e => e.RiskScore).HasPrecision(5, 2);
            entity.Property(e => e.AdjustedPremiumAmount).HasColumnName("adjusted_premium_amount");
            entity.Property(e => e.AdjustedPremiumCurrency).HasColumnName("adjusted_premium_currency").HasMaxLength(3).HasDefaultValue("BDT");
            
            entity.Property(e => e.RiskFactorsJson).HasColumnType("jsonb").HasColumnName("risk_factors");
            entity.Property(e => e.ConditionsJson).HasColumnType("jsonb").HasColumnName("conditions");
            entity.Property(e => e.AuditInfoJson).HasColumnType("jsonb").HasColumnName("audit_info");

            entity.Ignore(e => e.AdjustedPremium);

            entity.HasIndex(e => e.Decision);
        });

        // --- Claim ---
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

        // --- ClaimApproval ---
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

        // --- ClaimDocument ---
        modelBuilder.Entity<ClaimDocument>(entity =>
        {
            entity.ToTable("claim_documents");
            entity.HasKey(e => e.Id);
            entity.Property(e => e.DocumentType).HasMaxLength(100).IsRequired();
            
            entity.HasIndex(e => e.ClaimId);
        });

        // --- DB Sequences ---
        modelBuilder.HasSequence<long>("policy_number_seq", "insurance_schema")
            .StartsAt(1)
            .IncrementsBy(1);

        modelBuilder.HasSequence<long>("quote_number_seq", "insurance_schema")
            .StartsAt(1)
            .IncrementsBy(1);

        modelBuilder.HasSequence<long>("claim_number_seq", "insurance_schema")
            .StartsAt(1)
            .IncrementsBy(1);
    }
}
