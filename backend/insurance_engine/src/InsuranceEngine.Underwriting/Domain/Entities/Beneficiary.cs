using System;
using InsuranceEngine.Underwriting.Domain.Enums;

namespace InsuranceEngine.Underwriting.Domain.Entities;

public class Beneficiary
{
    public Guid Id { get; set; }
    public Guid UserId { get; set; }
    public BeneficiaryType Type { get; set; }
    public string Code { get; set; } = string.Empty;
    public BeneficiaryStatus Status { get; set; }
    public KYCStatus KycStatus { get; set; }
    public DateTime? KycCompletedAt { get; set; }
    public string? RiskScore { get; set; }
    public string? ReferralCode { get; set; }
    public Guid? ReferredBy { get; set; }
    public Guid? PartnerId { get; set; }

    public IndividualBeneficiary? IndividualDetails { get; set; }
    public BusinessBeneficiary? BusinessDetails { get; set; }

    public string? AuditInfo { get; set; }

    public DateTime CreatedAt { get; set; }
    public DateTime UpdatedAt { get; set; }
    public bool IsDeleted { get; set; }
}
