using InsuranceEngine.Beneficiary.Domain.Enums;
using InsuranceEngine.SharedKernel.Domain;
using InsuranceEngine.SharedKernel.Domain.ValueObjects;

namespace InsuranceEngine.Beneficiary.Domain.Entities;

public class BusinessBeneficiary : Entity<Guid>
{
    public Guid BeneficiaryId { get; private set; }
    
    public string BusinessName { get; private set; } = string.Empty;
    public string? BusinessNameBn { get; set; }
    public string TradeLicenseNumber { get; private set; } = string.Empty;
    public string TinNumber { get; private set; } = string.Empty;
    
    public ContactInfo ContactInfo { get; set; } = new();
    public Address RegisteredAddress { get; set; } = new();
    public Address BusinessAddress { get; set; } = new();
    
    public BusinessType BusinessType { get; set; }
    public string? IndustrySector { get; set; }
    public int? EmployeeCount { get; set; }
    public DateTime? IncorporationDate { get; set; }
    
    public string? TradeLicenseIssueDate { get; set; }
    public string? TradeLicenseExpiryDate { get; set; }
    public string? BinNumber { get; set; }
    public string? RegistrationNumber { get; set; }
    public string? TaxId { get; set; }

    public string FocalPersonName { get; set; } = string.Empty;
    public ContactInfo FocalPersonContact { get; set; } = new();
    public string? FocalPersonDesignation { get; set; }
    public string? FocalPersonNid { get; set; }
    
    public int ActivePoliciesCount { get; set; }
    public int PendingActionsCount { get; set; }
    public int TotalEmployeesCovered { get; set; }
    
    public Money TotalPremiumAmount { get; set; } = Money.Zero;
    
    public AuditInfo AuditInfo { get; private set; } = new();

    public void UpdateFocalPerson(string name, string mobile)
    {
        FocalPersonName = name;
        FocalPersonContact = new ContactInfo(mobile);
        AuditInfo = AuditInfo with { UpdatedAt = DateTime.UtcNow };
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
        AuditInfo = new AuditInfo();
    }
}
