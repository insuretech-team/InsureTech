using System;
using InsuranceEngine.Policy.Domain.Enums;

namespace InsuranceEngine.Policy.Domain.Entities;

/// <summary>
/// Individual (B2C) beneficiary details.
/// Maps to 'individual_beneficiaries' table in insurance_schema.
/// </summary>
public class IndividualBeneficiary
{
    public Guid BeneficiaryId { get; set; } // PK and FK to Beneficiary.Id
    public Beneficiary? Beneficiary { get; set; }

    public string FullName { get; set; } = string.Empty;
    public string? FullNameBn { get; set; }
    public DateTime DateOfBirth { get; set; }
    public Gender Gender { get; set; }
    
    public string? NidNumber { get; set; } // Encrypted at rest
    public string? PassportNumber { get; set; } // Encrypted at rest
    public string? BirthCertificateNumber { get; set; } // Encrypted at rest
    public string? TinNumber { get; set; }
    
    public MaritalStatus MaritalStatus { get; set; }
    public string? Occupation { get; set; }

    // JSONB in proto
    public string? ContactInfoJson { get; set; }
    public string? PermanentAddressJson { get; set; }
    public string? PresentAddressJson { get; set; }

    public string? NomineeName { get; set; }
    public string? NomineeRelationship { get; set; }

    public string? AuditInfo { get; set; }
    
    public DateTime CreatedAt { get; set; }
    public DateTime UpdatedAt { get; set; }
}
