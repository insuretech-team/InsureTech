using System.ComponentModel.DataAnnotations;
using System.ComponentModel.DataAnnotations.Schema;

namespace InsuranceEngine.SharedKernel.Persistence.Entities;

[Table("fraud_checks", Schema = "insurance_schema")]
public class FraudCheckEntity
{
    [Key]
    [Column("fraud_check_id")]
    public string FraudCheckId { get; set; } = string.Empty;

    [Column("entity_id")]
    public string EntityId { get; set; } = string.Empty;

    [Column("entity_type")]
    public string EntityType { get; set; } = string.Empty;

    [Column("check_type")]
    public string? CheckType { get; set; }

    [Column("fraud_score")]
    public int FraudScore { get; set; }

    [Column("risk_level")]
    public string RiskLevel { get; set; } = "LOW";

    [Column("flagged")]
    public bool Flagged { get; set; }

    [Column("claim_id")]
    public string? ClaimId { get; set; }

    [Column("customer_id")]
    public string? CustomerId { get; set; }

    [Column("claim_type")]
    public string? ClaimType { get; set; }

    [Column("claim_amount")]
    public decimal? ClaimAmount { get; set; }

    [Column("recommendation")]
    public string? Recommendation { get; set; }

    [Column("checked_at")]
    public DateTime CheckedAt { get; set; }

    [Column("created_at")]
    public DateTime CreatedAt { get; set; }

    [Column("risk_factors")]
    public string[]? RiskFactors { get; set; }

    public ClaimEntity? Claim { get; set; }
}
