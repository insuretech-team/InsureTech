using System;

namespace InsuranceEngine.SharedKernel.Persistence.Entities;

/// <summary>
/// EF Core entity for 'business_beneficiaries' table in insurance_schema.
/// Aligned with insuretech.beneficiary.entity.v1.BusinessBeneficiary proto.
/// </summary>
public class BusinessBeneficiaryEntity
{
    public Guid BeneficiaryId { get; set; } // Primary key
    public Guid ParentBeneficiaryId { get; set; } // Link to beneficiaries table
    public string BusinessName { get; set; } = string.Empty;
    public string? BusinessNameBn { get; set; }
    public string TradeLicenseNumber { get; set; } = string.Empty;
    public DateTime? TradeLicenseIssueDate { get; set; }
    public DateTime? TradeLicenseExpiryDate { get; set; }
    public string TinNumber { get; set; } = string.Empty;
    public string? BinNumber { get; set; }
    public string BusinessType { get; set; } = string.Empty;
    public string? IndustrySector { get; set; }
    public int? EmployeeCount { get; set; }
    public DateTime? IncorporationDate { get; set; }
    public string? ContactInfo { get; set; } // JSONB
    public string? RegisteredAddress { get; set; } // JSONB
    public string? BusinessAddress { get; set; } // JSONB
    public string FocalPersonName { get; set; } = string.Empty;
    public string? FocalPersonDesignation { get; set; }
    public string? FocalPersonNid { get; set; }
    public string? FocalPersonContact { get; set; } // JSONB
    public string? RegistrationNumber { get; set; }
    public string? TaxId { get; set; }
    public string? PrimaryContact { get; set; } // JSONB
    public int TotalEmployeesCovered { get; set; }
    public int ActivePoliciesCount { get; set; }
    public long TotalPremiumAmount { get; set; } // Paisa
    public int PendingActionsCount { get; set; }
    public string? AuditInfo { get; set; } // JSONB
    public DateTime CreatedAt { get; set; }
    public DateTime UpdatedAt { get; set; }

    // Navigation
    public BeneficiaryEntity Beneficiary { get; set; } = null!;
}
