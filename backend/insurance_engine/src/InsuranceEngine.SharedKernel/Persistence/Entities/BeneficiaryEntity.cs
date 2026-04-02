using System.ComponentModel.DataAnnotations;
using System.ComponentModel.DataAnnotations.Schema;

namespace InsuranceEngine.SharedKernel.Persistence.Entities;

/// <summary>
/// EF Core entity for 'beneficiaries' table in insurance_schema.
/// Aligned with the database schema managed by Go migrations.
/// </summary>
[Table("beneficiaries", Schema = "insurance_schema")]
public class BeneficiaryEntity
{
    [Key]
    [Column("beneficiary_id")]
    public Guid BeneficiaryId { get; set; }

    [Column("user_id")]
    public Guid UserId { get; set; }

    [Column("type")]
    public string Type { get; set; } = string.Empty; // INDIVIDUAL, BUSINESS

    [Column("code")]
    public string Code { get; set; } = $"BEN-{Guid.NewGuid().ToString()[..8].ToUpper()}";

    [Column("status")]
    public string Status { get; set; } = "PENDING_KYC";

    [Column("kyc_status")]
    public string KycStatus { get; set; } = "NOT_STARTED";

    [Column("kyc_completed_at")]
    public DateTime? KycCompletedAt { get; set; }

    [Column("risk_score")]
    public string? RiskScore { get; set; }

    [Column("referral_code")]
    public string? ReferralCode { get; set; }

    [Column("referred_by")]
    public Guid? ReferredBy { get; set; }

    [Column("partner_id")]
    public Guid? PartnerId { get; set; }

    [Column("audit_info")]
    public string? AuditInfo { get; set; } // JSONB

    [Column("created_at")]
    public DateTime CreatedAt { get; set; }

    [Column("updated_at")]
    public DateTime UpdatedAt { get; set; }

    [Column("deleted_at")]
    public DateTime? DeletedAt { get; set; }

    // Navigation
    public IndividualBeneficiaryEntity? IndividualDetails { get; set; }
    public BusinessBeneficiaryEntity? BusinessDetails { get; set; }
}
