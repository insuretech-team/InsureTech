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
    public DbSet<UnderwritingHealthDeclaration> HealthDeclarations { get; set; } = null!;
    public DbSet<UnderwritingDecision> UnderwritingDecisions { get; set; } = null!;

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        modelBuilder.HasDefaultSchema("insurance_schema");

        modelBuilder.Entity<Quote>(entity =>
        {
            entity.ToTable("quotes");
            entity.HasKey(e => e.Id);
            entity.Property(e => e.Id).HasColumnName("quote_id");
            entity.HasQueryFilter(e => e.DeletedAt == null);
            entity.Property(e => e.QuoteNumber).HasMaxLength(50).IsRequired();
            entity.HasIndex(e => e.QuoteNumber).IsUnique();
            entity.Property(e => e.Status).HasConversion<string>().HasMaxLength(50).IsRequired();
            
            entity.Property(e => e.SumAssuredAmount).HasColumnName("sum_assured_amount");
            entity.Property(e => e.SumAssuredCurrency).HasColumnName("sum_assured_currency").HasMaxLength(3).HasDefaultValue("BDT");
            
            entity.Property(e => e.TermYears).HasColumnName("term_years").IsRequired();
            entity.Property(e => e.PremiumPaymentMode).HasColumnName("premium_payment_mode").HasMaxLength(50).IsRequired();

            entity.Property(e => e.BasePremiumAmount).HasColumnName("base_premium_amount");
            entity.Property(e => e.BasePremiumCurrency).HasColumnName("base_premium_currency").HasMaxLength(3).HasDefaultValue("BDT");
            entity.Property(e => e.RiderPremiumAmount).HasColumnName("rider_premium_amount");
            entity.Property(e => e.RiderPremiumCurrency).HasColumnName("rider_premium_currency").HasMaxLength(3).HasDefaultValue("BDT");
            entity.Property(e => e.TaxAmount).HasColumnName("tax_amount");
            entity.Property(e => e.TaxCurrency).HasColumnName("tax_currency").HasMaxLength(3).HasDefaultValue("BDT");
            entity.Property(e => e.TotalPremiumAmount).HasColumnName("total_premium_amount");
            entity.Property(e => e.TotalPremiumCurrency).HasColumnName("total_premium_currency").HasMaxLength(3).HasDefaultValue("BDT");

            entity.Property(e => e.PremiumCalculationJson).HasColumnType("jsonb").HasColumnName("premium_calculation");
            entity.Property(e => e.SelectedRidersJson).HasColumnType("jsonb").HasColumnName("selected_riders");
            
            entity.Property(e => e.ApplicantAge).HasColumnName("applicant_age");
            entity.Property(e => e.ApplicantOccupation).HasColumnName("applicant_occupation").HasMaxLength(100);
            entity.Property(e => e.IsSmoker).HasColumnName("smoker");

            entity.Property(e => e.ValidUntil).HasColumnName("valid_until");
            entity.Property(e => e.ConvertedPolicyId).HasColumnName("converted_policy_id");
            entity.Property(e => e.ConvertedAt).HasColumnName("converted_at");

            entity.Property(e => e.AuditInfoJson).HasColumnType("jsonb").HasColumnName("audit_info");
            entity.Property(e => e.DeletedAt).HasColumnName("deleted_at");

            entity.Ignore(e => e.SumAssured);
            entity.Ignore(e => e.BasePremium);
            entity.Ignore(e => e.RiderPremium);
            entity.Ignore(e => e.Tax);
            entity.Ignore(e => e.TotalPremium);

            entity.HasIndex(e => e.BeneficiaryId);
            entity.HasIndex(e => e.InsurerProductId);
            entity.HasIndex(e => e.Status);
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
