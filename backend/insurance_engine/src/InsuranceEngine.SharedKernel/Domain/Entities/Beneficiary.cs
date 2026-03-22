using System;
using InsuranceEngine.SharedKernel.Domain.Enums;

namespace InsuranceEngine.SharedKernel.Domain.Entities;

public class Beneficiary
{
    public Guid Id { get; set; }
    public Guid UserId { get; set; }
    public BeneficiaryType Type { get; set; }
    public string Code { get; set; } = string.Empty; // BEN-XXXXXX
    public BeneficiaryStatus Status { get; set; }
    public KYCStatus KycStatus { get; set; }
    public DateTime? KycCompletedAt { get; set; }
    public string? RiskScore { get; set; } // LOW, MEDIUM, HIGH
    public string? ReferralCode { get; set; }
    public Guid? ReferredBy { get; set; }
    public Guid? PartnerId { get; set; }

    // Unified fields from Underwriting
    public string Name { get; set; } = string.Empty;
    public string ContactNumber { get; set; } = string.Empty;
    public string Email { get; set; } = string.Empty;
    public string Address { get; set; } = string.Empty;

    // Relationship to specialized details
    public IndividualBeneficiary? IndividualDetails { get; set; }
    public BusinessBeneficiary? BusinessDetails { get; set; }

    public string? AuditInfo { get; set; }

    public DateTime CreatedAt { get; set; }
    public DateTime UpdatedAt { get; set; }
    public bool IsDeleted { get; set; }
}
