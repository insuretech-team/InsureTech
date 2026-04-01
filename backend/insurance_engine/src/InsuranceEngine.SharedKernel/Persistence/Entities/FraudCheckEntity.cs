using System.ComponentModel.DataAnnotations;
using System.ComponentModel.DataAnnotations.Schema;

namespace InsuranceEngine.SharedKernel.Persistence.Entities;

/// <summary>
/// EF Core entity for 'fraud_checks' table in insurance_schema.
/// Aligned with the database schema managed by Go migrations.
/// One-to-one with Claim.
/// </summary>
[Table("fraud_checks", Schema = "insurance_schema")]
public class FraudCheckEntity
{
    [Key]
    [Column("fraud_check_id")]
    public Guid FraudCheckId { get; set; }

    [Column("claim_id")]
    public Guid ClaimId { get; set; }

    [Column("fraud_score")]
    public double FraudScore { get; set; } // 0-100

    [Column("risk_factors")]
    public string[] RiskFactors { get; set; } = []; // JSONB or text[]

    [Column("flagged")]
    public bool Flagged { get; set; }

    [Column("reviewed_by")]
    public Guid? ReviewedBy { get; set; }

    [Column("reviewed_at")]
    public DateTime? ReviewedAt { get; set; }

    [Column("created_at")]
    public DateTime CreatedAt { get; set; }

    // Navigation
    public ClaimEntity Claim { get; set; } = null!;
}
