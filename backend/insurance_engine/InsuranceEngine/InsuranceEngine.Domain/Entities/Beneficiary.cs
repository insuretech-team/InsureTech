using System.ComponentModel.DataAnnotations;

using InsuranceEngine.Domain.Enums;

namespace InsuranceEngine.Domain.Entities;

public class Beneficiary
{
    [Key]
    public Guid BeneficiaryId { get; set; }

    [Required]
    [MaxLength(100)]
    public string UserId { get; set; } = string.Empty;

    [MaxLength(100)]
    public string? PartnerId { get; set; }

    [Required]
    [MaxLength(50)]
    public string BeneficiaryCode { get; set; } = string.Empty;

    public Guid? PolicyId { get; set; }

    [Required]
    public BeneficiaryType Type { get; set; }

    [Required]
    public BeneficiaryStatus Status { get; set; }

    [MaxLength(32)]
    [Required]
    public string KycStatus { get; set; } = "KYC_STATUS_UNSPECIFIED";

    public DateTime? KycCompletedAt { get; set; }

    [MaxLength(50)]
    public string? RiskScore { get; set; }

    [MaxLength(100)]
    public string? ReferralCode { get; set; }

    [MaxLength(100)]
    public string? ReferredBy { get; set; }

    public AuditInfo AuditInfo { get; set; } = new();

    public BeneficiaryIndividual? IndividualDetails { get; set; }

    public BeneficiaryBusiness? BusinessDetails { get; set; }
}
