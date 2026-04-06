using Microsoft.EntityFrameworkCore;
using InsuranceEngine.SharedKernel.Persistence.Entities;

namespace InsuranceEngine.Quoting.Infrastructure;

public class QuotingDbContext : DbContext
{
    public QuotingDbContext(DbContextOptions<QuotingDbContext> options) : base(options) { }

    public DbSet<QuoteEntity> Quotes => Set<QuoteEntity>();

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        modelBuilder.HasDefaultSchema("insurance_schema");

        modelBuilder.Entity<QuoteEntity>(entity =>
        {
            entity.ToTable("quotes");
            entity.HasKey(e => e.QuoteId);
            entity.Property(e => e.QuoteId).HasColumnName("quote_id").HasDefaultValueSql("gen_random_uuid()");
            entity.Property(e => e.QuoteNumber).HasColumnName("quote_number").HasMaxLength(50).IsRequired();
            entity.Property(e => e.BeneficiaryId).HasColumnName("beneficiary_id").IsRequired();
            entity.Property(e => e.InsurerProductId).HasColumnName("insurer_product_id").IsRequired();
            entity.Property(e => e.Status).HasColumnName("status").HasMaxLength(50).HasDefaultValue("DRAFT");
            entity.Property(e => e.SumAssured).HasColumnName("sum_assured").IsRequired();
            entity.Property(e => e.SumAssuredCurrency).HasColumnName("sum_assured_currency").HasMaxLength(3).HasDefaultValue("BDT");
            entity.Property(e => e.TermYears).HasColumnName("term_years").IsRequired();
            entity.Property(e => e.PremiumPaymentMode).HasColumnName("premium_payment_mode").HasMaxLength(50).HasDefaultValue("YEARLY");
            entity.Property(e => e.BasePremium).HasColumnName("base_premium").IsRequired();
            entity.Property(e => e.BasePremiumCurrency).HasColumnName("base_premium_currency").HasMaxLength(3).HasDefaultValue("BDT");
            entity.Property(e => e.RiderPremium).HasColumnName("rider_premium");
            entity.Property(e => e.TaxAmount).HasColumnName("tax_amount");
            entity.Property(e => e.TotalPremium).HasColumnName("total_premium").IsRequired();
            entity.Property(e => e.TotalPremiumCurrency).HasColumnName("total_premium_currency").HasMaxLength(3).HasDefaultValue("BDT");
            entity.Property(e => e.PremiumCalculation).HasColumnName("premium_calculation");
            entity.Property(e => e.SelectedRiders).HasColumnName("selected_riders");
            entity.Property(e => e.ApplicantAge).HasColumnName("applicant_age").IsRequired();
            entity.Property(e => e.ApplicantOccupation).HasColumnName("applicant_occupation").HasMaxLength(100);
            entity.Property(e => e.Smoker).HasColumnName("smoker").HasDefaultValue(false);
            entity.Property(e => e.ValidUntil).HasColumnName("valid_until").IsRequired();
            entity.Property(e => e.ConvertedPolicyId).HasColumnName("converted_policy_id");
            entity.Property(e => e.ConvertedAt).HasColumnName("converted_at");
            entity.Property(e => e.CreatedAt).HasColumnName("created_at").HasDefaultValueSql("CURRENT_TIMESTAMP").IsRequired();
            entity.Property(e => e.UpdatedAt).HasColumnName("updated_at").HasDefaultValueSql("CURRENT_TIMESTAMP").IsRequired();
            entity.Property(e => e.DeletedAt).HasColumnName("deleted_at");

            entity.HasIndex(e => e.QuoteNumber).IsUnique().HasDatabaseName("idx_quotes_number");
            entity.HasIndex(e => e.BeneficiaryId).HasDatabaseName("idx_quotes_beneficiary_id");
            entity.HasIndex(e => e.Status).HasDatabaseName("idx_quotes_status");
        });
    }
}
