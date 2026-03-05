using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Metadata.Builders;
using PoliSync.Underwriting.Domain;

namespace PoliSync.Underwriting.Persistence;

public class QuoteConfiguration : IEntityTypeConfiguration<Quote>
{
    public void Configure(EntityTypeBuilder<Quote> builder)
    {
        builder.ToTable("quotes", "insurance_schema");
        builder.HasKey(q => q.QuoteId);

        builder.Property(q => q.QuoteId).HasColumnName("quote_id");
        builder.Property(q => q.QuoteNumber).HasColumnName("quote_number").HasMaxLength(50).IsRequired();
        builder.Property(q => q.BeneficiaryId).HasColumnName("beneficiary_id").IsRequired();
        builder.Property(q => q.InsurerProductId).HasColumnName("insurer_product_id").IsRequired();
        builder.Property(q => q.Status).HasColumnName("status").HasMaxLength(50).IsRequired().HasConversion<string>();
        
        builder.Property(q => q.SumAssured).HasColumnName("sum_assured").IsRequired();
        builder.Property(q => q.Currency).HasColumnName("currency").HasMaxLength(3).IsRequired();
        builder.Property(q => q.TermYears).HasColumnName("term_years").IsRequired();
        builder.Property(q => q.PremiumPaymentMode).HasColumnName("premium_payment_mode").HasMaxLength(50).IsRequired();
        
        builder.Property(q => q.BasePremiumAmount).HasColumnName("base_premium_amount").IsRequired();
        builder.Property(q => q.RiderPremiumAmount).HasColumnName("rider_premium_amount").IsRequired();
        builder.Property(q => q.TaxAmount).HasColumnName("tax_amount").IsRequired();
        builder.Property(q => q.TotalPremiumAmount).HasColumnName("total_premium_amount").IsRequired();
        
        builder.Property(q => q.ApplicantAgeDays).HasColumnName("applicant_age_days").IsRequired();
        builder.Property(q => q.IsSmoker).HasColumnName("is_smoker").IsRequired();
        
        builder.Property(q => q.CreatedAt).HasColumnName("created_at").IsRequired();
        builder.Property(q => q.UpdatedAt).HasColumnName("updated_at").IsRequired();
        builder.Property(q => q.ValidUntil).HasColumnName("valid_until").IsRequired();

        builder.Ignore(q => q.DomainEvents);

        // Relationships
        builder.HasOne(q => q.HealthDeclaration)
            .WithOne()
            .HasForeignKey<HealthDeclaration>(hd => hd.QuoteId);

        builder.HasOne(q => q.Decision)
            .WithOne()
            .HasForeignKey<UnderwritingDecision>(d => d.QuoteId);

        builder.HasIndex(q => q.QuoteNumber).IsUnique();
        builder.HasIndex(q => q.BeneficiaryId);
    }
}

public class HealthDeclarationConfiguration : IEntityTypeConfiguration<HealthDeclaration>
{
    public void Configure(EntityTypeBuilder<HealthDeclaration> builder)
    {
        builder.ToTable("health_declarations", "insurance_schema");
        builder.HasKey(hd => hd.DeclarationId);

        builder.Property(hd => hd.DeclarationId).HasColumnName("declaration_id");
        builder.Property(hd => hd.QuoteId).HasColumnName("quote_id").IsRequired();
        
        builder.Property(hd => hd.HeightCm).HasColumnName("height_cm").IsRequired();
        builder.Property(hd => hd.WeightKg).HasColumnName("weight_kg").IsRequired();
        builder.Property(hd => hd.Bmi).HasColumnName("bmi").IsRequired();
        
        builder.Property(hd => hd.IsSmoker).HasColumnName("is_smoker").IsRequired();
        builder.Property(hd => hd.ConsumesAlcohol).HasColumnName("consumes_alcohol").IsRequired();
        
        builder.Property(hd => hd.HasPreExistingConditions).HasColumnName("has_pre_existing_conditions").IsRequired();
        builder.Property(hd => hd.ConditionDetails).HasColumnName("condition_details").HasColumnType("jsonb");
        
        builder.Property(hd => hd.HasFamilyHistoryOfCriticalIllness).HasColumnName("has_family_history_of_critical_illness").IsRequired();
        builder.Property(hd => hd.OccupationRiskLevel).HasColumnName("occupation_risk_level").HasMaxLength(50);
        
        builder.Property(hd => hd.SubmittedAt).HasColumnName("submitted_at").IsRequired();

        builder.Ignore(hd => hd.DomainEvents);
        builder.HasIndex(hd => hd.QuoteId).IsUnique();
    }
}

public class UnderwritingDecisionConfiguration : IEntityTypeConfiguration<UnderwritingDecision>
{
    public void Configure(EntityTypeBuilder<UnderwritingDecision> builder)
    {
        builder.ToTable("underwriting_decisions", "insurance_schema");
        builder.HasKey(d => d.DecisionId);

        builder.Property(d => d.DecisionId).HasColumnName("decision_id");
        builder.Property(d => d.QuoteId).HasColumnName("quote_id").IsRequired();
        
        builder.Property(d => d.Decision).HasColumnName("decision").HasMaxLength(50).IsRequired().HasConversion<string>();
        builder.Property(d => d.Method).HasColumnName("method").HasMaxLength(50).IsRequired().HasConversion<string>();
        builder.Property(d => d.RiskScore).HasColumnName("risk_score").IsRequired();
        builder.Property(d => d.RiskLevel).HasColumnName("risk_level").HasMaxLength(50).IsRequired().HasConversion<string>();
        
        builder.Property(d => d.Reason).HasColumnName("reason");
        builder.Property(d => d.Conditions).HasColumnName("conditions").HasColumnType("jsonb");
        builder.Property(d => d.RiskFactors).HasColumnName("risk_factors").HasColumnType("jsonb");
        
        builder.Property(d => d.AdjustedPremiumAmount).HasColumnName("adjusted_premium_amount");
        builder.Property(d => d.UnderwriterId).HasColumnName("underwriter_id");
        
        builder.Property(d => d.DecidedAt).HasColumnName("decided_at").IsRequired();

        builder.Ignore(d => d.DomainEvents);
        builder.HasIndex(d => d.QuoteId).IsUnique();
    }
}
