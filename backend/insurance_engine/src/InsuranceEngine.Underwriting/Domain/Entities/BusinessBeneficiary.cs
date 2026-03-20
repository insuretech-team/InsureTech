using System;
using InsuranceEngine.Underwriting.Domain.Enums;

namespace InsuranceEngine.Underwriting.Domain.Entities;

public class BusinessBeneficiary
{
    public Guid Id { get; set; }
    public Guid BeneficiaryId { get; set; }
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

    public string? ContactInfoJson { get; set; }
    public string? RegisteredAddressJson { get; set; }
    public string? BusinessAddressJson { get; set; }

    public string FocalPersonName { get; set; } = string.Empty;
    public string? FocalPersonDesignation { get; set; }
    public string? FocalPersonNid { get; set; }
    public string? FocalPersonContactJson { get; set; }
    public string? Industry { get; set; }
    public string? FocalPersonContact { get; set; }

    public string? AuditInfo { get; set; }
    public string? RegistrationNumber { get; set; }
    public string? TaxId { get; set; }
    public string? PrimaryContactJson { get; set; }

    public int TotalEmployeesCovered { get; set; }
    public int ActivePoliciesCount { get; set; }
    public long TotalPremiumAmount { get; set; }
    public int PendingActionsCount { get; set; }

    public DateTime CreatedAt { get; set; }
    public DateTime UpdatedAt { get; set; }
}
