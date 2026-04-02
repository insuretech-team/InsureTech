using System.ComponentModel.DataAnnotations;
using System.ComponentModel.DataAnnotations.Schema;

namespace InsuranceEngine.SharedKernel.Persistence.Entities;

/// <summary>
/// EF Core entity for 'underwriting_decisions' table in insurance_schema.
/// Aligned with the database schema managed by Go migrations.
/// </summary>
[Table("underwriting_decisions", Schema = "insurance_schema")]
public class UnderwritingDecisionEntity
{
    [Key]
    [Column("decision_id")]
    public Guid DecisionId { get; set; }

    [Column("quote_id")]
    public Guid QuoteId { get; set; }

    [Column("decision")]
    public string Decision { get; set; } = "PENDING"; // APPROVED, REJECTED, REFERRED, CONDITIONAL

    [Column("method")]
    public string Method { get; set; } = "MANUAL"; // AUTOMATIC, MANUAL, HYBRID

    [Column("risk_score")]
    public string? RiskScore { get; set; }

    [Column("risk_level")]
    public string? RiskLevel { get; set; }

    [Column("risk_factors")]
    public string? RiskFactors { get; set; } // JSONB

    [Column("reason")]
    public string? Reason { get; set; }

    [Column("conditions")]
    public string? Conditions { get; set; } // JSONB

    [Column("premium_adjusted")]
    public bool PremiumAdjusted { get; set; } = false;

    [Column("adjusted_premium")]
    public long? AdjustedPremium { get; set; }

    [Column("adjusted_premium_currency")]
    public string AdjustedPremiumCurrency { get; set; } = "BDT";

    [Column("adjustment_reason")]
    public string? AdjustmentReason { get; set; }

    [Column("underwriter_id")]
    public Guid UnderwriterId { get; set; }

    [Column("underwriter_comments")]
    public string? UnderwriterComments { get; set; }

    [Column("decided_at")]
    public DateTime DecidedAt { get; set; }

    [Column("comments")]
    public string? Comments { get; set; }

    [Column("rejection_reason")]
    public string? RejectionReason { get; set; }

    [Column("valid_until")]
    public DateTime? ValidUntil { get; set; }

    [Column("created_at")]
    public DateTime CreatedAt { get; set; }

    // Navigation
    public QuoteEntity Quote { get; set; } = null!;
}
