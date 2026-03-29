namespace InsuranceEngine.SharedKernel.Persistence.Entities;

/// <summary>
/// EF Core entity for 'underwriting_decisions' table in insurance_schema.
/// Aligned with insuretech.underwriting.entity.v1.UnderwritingDecision proto.
/// </summary>
public class UnderwritingDecisionEntity
{
    public Guid DecisionId { get; set; }
    public Guid QuoteId { get; set; }
    public Guid UnderwriterId { get; set; }
    public string Decision { get; set; } = "PENDING"; // APPROVED, REJECTED
    public string? RiskLevel { get; set; }
    public bool PremiumAdjusted { get; set; }
    public long? AdjustedPremium { get; set; }
    public string AdjustedPremiumCurrency { get; set; } = "BDT";
    public string? Conditions { get; set; } // JSONB
    public string? Comments { get; set; }
    public string? RejectionReason { get; set; }
    public DateTime CreatedAt { get; set; }

    // Navigation
    public QuoteEntity Quote { get; set; } = null!;
}
