using System;
using InsuranceEngine.Policy.Domain.Enums;

namespace InsuranceEngine.Policy.Domain.Entities;

/// <summary>
/// Base Beneficiary entity (primary customer record).
/// Maps to 'beneficiaries' table in insurance_schema.
/// </summary>
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

    // Relationship to specialized details
    public IndividualBeneficiary? IndividualDetails { get; set; }
    public BusinessBeneficiary? BusinessDetails { get; set; }

    // Audit info as JSONB (in proto) - we'll use properties here for simple mapping or a JSON string
    public string? AuditInfo { get; set; }

    public DateTime CreatedAt { get; set; }
    public DateTime UpdatedAt { get; set; }
    public bool IsDeleted { get; set; }
}
