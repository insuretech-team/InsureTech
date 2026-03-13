using InsuranceEngine.Underwriting.Domain.Entities;
using InsuranceEngine.Underwriting.Domain.Enums;
using Microsoft.EntityFrameworkCore;

namespace InsuranceEngine.Underwriting.Infrastructure.Persistence;

public class UnderwritingDbContext : DbContext
{
    public UnderwritingDbContext(DbContextOptions<UnderwritingDbContext> options) : base(options)
    {
    }

    public DbSet<Quote> Quotes { get; set; } = null!;
    public DbSet<Beneficiary> Beneficiaries { get; set; } = null!;
    public DbSet<IndividualBeneficiary> IndividualBeneficiaries { get; set; } = null!;
    public DbSet<BusinessBeneficiary> BusinessBeneficiaries { get; set; } = null!;
    public DbSet<UnderwritingHealthDeclaration> HealthDeclarations { get; set; } = null!;
    public DbSet<UnderwritingDecision> UnderwritingDecisions { get; set; } = null!;

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        modelBuilder.HasDefaultSchema("insurance_schema");

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

            entity.HasIndex(e => (object)e.Decision);
        });

        modelBuilder.HasSequence<long>("quote_number_seq", "insurance_schema")
            .StartsAt(1)
            .IncrementsBy(1);
    }
}
