using System;
using InsuranceEngine.Policy.Domain.Enums;

namespace InsuranceEngine.Policy.Domain.Entities;

/// <summary>
/// Business (B2B) beneficiary details.
/// Maps to 'business_beneficiaries' table in insurance_schema.
/// </summary>
public class BusinessBeneficiary
{
    public Guid Id { get; set; } // PK
    public Guid BeneficiaryId { get; set; } // Parent Beneficiary ID
    public Beneficiary? Beneficiary { get; set; }

    public string BusinessName { get; set; } = string.Empty;
    public string? BusinessNameBn { get; set; }
    public string TradeLicenseNumber { get; set; } = string.Empty;
    public DateTime? TradeLicenseIssueDate { get; set; }
    public DateTime? TradeLicenseExpiryDate { get; set; }
    public string TinNumber { get; set; } = string.Empty;
    public string? BinNumber { get; set; }
    public BusinessType BusinessType { get; set; }
    public string? IndustrySector { get; set; }
    public int EmployeeCount { get; set; }
    public DateTime? IncorporationDate { get; set; }

    // JSONB in proto
    public string? ContactInfoJson { get; set; }
    public string? RegisteredAddressJson { get; set; }
    public string? BusinessAddressJson { get; set; }

    public string FocalPersonName { get; set; } = string.Empty;
    public string? FocalPersonDesignation { get; set; }
    public string? FocalPersonNid { get; set; } // Encrypted at rest
    public string? FocalPersonContactJson { get; set; }

    public string? AuditInfo { get; set; }
    public string? RegistrationNumber { get; set; }
    public string? TaxId { get; set; }
    public string? PrimaryContactJson { get; set; }

    // Cached metrics
    public int TotalEmployeesCovered { get; set; }
    public int ActivePoliciesCount { get; set; }
    public long TotalPremiumAmount { get; set; } // In paisa
    public int PendingActionsCount { get; set; }

    public DateTime CreatedAt { get; set; }
    public DateTime UpdatedAt { get; set; }
}
