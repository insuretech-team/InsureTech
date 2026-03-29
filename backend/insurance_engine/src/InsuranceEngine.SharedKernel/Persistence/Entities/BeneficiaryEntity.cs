using System;
using System.Collections.Generic;

namespace InsuranceEngine.SharedKernel.Persistence.Entities;

/// <summary>
/// EF Core entity for 'beneficiaries' table in insurance_schema.
/// Aligned with insuretech.beneficiary.entity.v1.Beneficiary proto.
/// </summary>
public class BeneficiaryEntity
{
    public Guid BeneficiaryId { get; set; }
    public Guid UserId { get; set; }
    public string Type { get; set; } = string.Empty; // INDIVIDUAL, BUSINESS
    public string Code { get; set; } = string.Empty; // BEN-XXXXXX
    public string Status { get; set; } = "PENDING_KYC";
    public string KycStatus { get; set; } = "NOT_STARTED";
    public DateTime? KycCompletedAt { get; set; }
    public string? RiskScore { get; set; }
    public string? ReferralCode { get; set; }
    public Guid? ReferredBy { get; set; }
    public Guid? PartnerId { get; set; }
    public string? AuditInfo { get; set; } // JSONB
    public DateTime CreatedAt { get; set; }
    public DateTime UpdatedAt { get; set; }
    public DateTime? DeletedAt { get; set; }

    // Navigation
    public IndividualBeneficiaryEntity? IndividualDetails { get; set; }
    public BusinessBeneficiaryEntity? BusinessDetails { get; set; }
}
