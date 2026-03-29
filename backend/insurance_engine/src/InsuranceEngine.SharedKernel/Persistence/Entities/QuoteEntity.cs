namespace InsuranceEngine.SharedKernel.Persistence.Entities;

/// <summary>
/// EF Core entity for 'quotes' table in insurance_schema.
/// Aligned with insuretech.underwriting.entity.v1.Quote proto definition.
/// </summary>
public class QuoteEntity
{
    public Guid QuoteId { get; set; }
    public string QuoteNumber { get; set; } = string.Empty;
    public Guid BeneficiaryId { get; set; }
    public Guid InsurerProductId { get; set; }
    public string Status { get; set; } = "DRAFT";
    public long SumAssured { get; set; }
    public string SumAssuredCurrency { get; set; } = "BDT";
    public int TermYears { get; set; }
    public string PremiumPaymentMode { get; set; } = "YEARLY";
    public long BasePremium { get; set; }
    public string BasePremiumCurrency { get; set; } = "BDT";
    public long? RiderPremium { get; set; }
    public long? TaxAmount { get; set; }
    public long TotalPremium { get; set; }
    public string TotalPremiumCurrency { get; set; } = "BDT";
    public string? PremiumCalculation { get; set; } // JSONB
    public string? SelectedRiders { get; set; } // JSONB
    public int ApplicantAge { get; set; }
    public string? ApplicantOccupation { get; set; }
    public bool Smoker { get; set; }
    public DateTime ValidUntil { get; set; }
    public Guid? ConvertedPolicyId { get; set; }
    public DateTime? ConvertedAt { get; set; }
    public DateTime CreatedAt { get; set; }
    public DateTime UpdatedAt { get; set; }
    public DateTime? DeletedAt { get; set; }

    // Navigation
    public ICollection<HealthDeclarationEntity> HealthDeclarations { get; set; } = new List<HealthDeclarationEntity>();
    public ICollection<UnderwritingDecisionEntity> Decisions { get; set; } = new List<UnderwritingDecisionEntity>();
}
