namespace InsuranceEngine.SharedKernel.Persistence.Entities;

/// <summary>
/// EF Core entity for 'fraud_checks' table in insurance_schema.
/// Aligned with insuretech.claims.entity.v1.FraudCheckResult proto definition.
/// One-to-one with Claim.
/// </summary>
public class FraudCheckEntity
{
    public Guid FraudCheckId { get; set; }
    public Guid ClaimId { get; set; }
    public double FraudScore { get; set; } // 0-100
    public string[] RiskFactors { get; set; } = [];
    public bool Flagged { get; set; }
    public Guid? ReviewedBy { get; set; }
    public DateTime? ReviewedAt { get; set; }
    public DateTime CreatedAt { get; set; }

    // Navigation
    public ClaimEntity Claim { get; set; } = null!;
}
