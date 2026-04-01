using System.ComponentModel.DataAnnotations;
using System.ComponentModel.DataAnnotations.Schema;

namespace InsuranceEngine.SharedKernel.Persistence.Entities;

/// <summary>
/// EF Core entity for 'quotes' table in insurance_schema.
/// Aligned with the database schema managed by Go migrations.
/// </summary>
[Table("quotes", Schema = "insurance_schema")]
public class QuoteEntity
{
    [Key]
    [Column("quote_id")]
    public Guid QuoteId { get; set; }

    [Column("quote_number")]
    public string QuoteNumber { get; set; } = string.Empty;

    [Column("beneficiary_id")]
    public Guid BeneficiaryId { get; set; }

    [Column("insurer_product_id")]
    public Guid InsurerProductId { get; set; }

    [Column("status")]
    public string Status { get; set; } = "DRAFT";

    [Column("sum_assured")]
    public long SumAssured { get; set; }

    [Column("sum_assured_currency")]
    public string SumAssuredCurrency { get; set; } = "BDT";

    [Column("term_years")]
    public int TermYears { get; set; }

    [Column("premium_payment_mode")]
    public string PremiumPaymentMode { get; set; } = "YEARLY";

    [Column("base_premium")]
    public long BasePremium { get; set; }

    [Column("base_premium_currency")]
    public string BasePremiumCurrency { get; set; } = "BDT";

    [Column("rider_premium")]
    public long? RiderPremium { get; set; }

    [Column("tax_amount")]
    public long? TaxAmount { get; set; }

    [Column("total_premium")]
    public long TotalPremium { get; set; }

    [Column("total_premium_currency")]
    public string TotalPremiumCurrency { get; set; } = "BDT";

    [Column("premium_calculation")]
    public string? PremiumCalculation { get; set; } // JSONB

    [Column("selected_riders")]
    public string? SelectedRiders { get; set; } // JSONB

    [Column("applicant_age")]
    public int ApplicantAge { get; set; }

    [Column("applicant_occupation")]
    public string? ApplicantOccupation { get; set; }

    [Column("smoker")]
    public bool Smoker { get; set; }

    [Column("valid_until")]
    public DateTime ValidUntil { get; set; }

    [Column("converted_policy_id")]
    public Guid? ConvertedPolicyId { get; set; }

    [Column("converted_at")]
    public DateTime? ConvertedAt { get; set; }

    [Column("created_at")]
    public DateTime CreatedAt { get; set; }

    [Column("updated_at")]
    public DateTime UpdatedAt { get; set; }

    [Column("deleted_at")]
    public DateTime? DeletedAt { get; set; }

    // Navigation
    public ICollection<HealthDeclarationEntity> HealthDeclarations { get; set; } = new List<HealthDeclarationEntity>();
    public ICollection<UnderwritingDecisionEntity> Decisions { get; set; } = new List<UnderwritingDecisionEntity>();
}
