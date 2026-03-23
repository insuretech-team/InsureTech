using System;
using InsuranceEngine.Beneficiary.Domain.Enums;
using InsuranceEngine.SharedKernel.Domain;

namespace InsuranceEngine.Beneficiary.Domain.Entities;

public class BusinessBeneficiary : Entity<Guid>
{
    public Guid BeneficiaryId { get; private set; }
    
    public string BusinessName { get; private set; } = string.Empty;
    public string? BusinessNameBn { get; set; }
    public string TradeLicenseNumber { get; private set; } = string.Empty;
    public string TinNumber { get; private set; } = string.Empty;
    public string? ContactInfoJson { get; set; }
    public string? RegisteredAddressJson { get; set; }
    public string? BusinessAddressJson { get; set; }
    public BusinessType BusinessType { get; set; }
    public string? IndustrySector { get; set; }
    public int? EmployeeCount { get; set; }
    public DateTime? IncorporationDate { get; set; }
    
    public string? TradeLicenseIssueDate { get; set; }
    public string? TradeLicenseExpiryDate { get; set; }
    public string? BinNumber { get; set; }
    public string? RegistrationNumber { get; set; }
    public string? TaxId { get; set; }

    public string FocalPersonName { get; private set; } = string.Empty;
    public string? FocalPersonMobile { get; private set; }
    public string? FocalPersonDesignation { get; set; }
    public string? FocalPersonNid { get; set; }
    
    public string? AuditInfoJson { get; set; }

    public DateTime CreatedAt { get; private set; }
    public DateTime UpdatedAt { get; private set; }

    public void UpdateFocalPerson(string name, string mobile)
    {
        FocalPersonName = name;
        FocalPersonMobile = mobile;
        UpdatedAt = DateTime.UtcNow;
    }

    // EF Core constructor
    public BusinessBeneficiary() { }

    public BusinessBeneficiary(
        Guid id,
        Guid beneficiaryId,
        string businessName,
        string tradeLicense,
        string tin) : base(id)
    {
        BeneficiaryId = beneficiaryId;
        BusinessName = businessName;
        TradeLicenseNumber = tradeLicense;
        TinNumber = tin;
        CreatedAt = DateTime.UtcNow;
        UpdatedAt = DateTime.UtcNow;
    }
}
