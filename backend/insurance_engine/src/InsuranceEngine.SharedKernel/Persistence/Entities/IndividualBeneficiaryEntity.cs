using System;

namespace InsuranceEngine.SharedKernel.Persistence.Entities;

/// <summary>
/// EF Core entity for 'individual_beneficiaries' table in insurance_schema.
/// Aligned with insuretech.beneficiary.entity.v1.IndividualBeneficiary proto.
/// </summary>
public class IndividualBeneficiaryEntity
{
    public Guid BeneficiaryId { get; set; }
    public string FullName { get; set; } = string.Empty;
    public string? FullNameBn { get; set; }
    public DateTime DateOfBirth { get; set; }
    public string Gender { get; set; } = string.Empty;
    public string? NidNumber { get; set; }
    public string? PassportNumber { get; set; }
    public string? BirthCertificateNumber { get; set; }
    public string? TinNumber { get; set; }
    public string? MaritalStatus { get; set; }
    public string? Occupation { get; set; }
    public string? ContactInfo { get; set; } // JSONB
    public string? PermanentAddress { get; set; } // JSONB
    public string? PresentAddress { get; set; } // JSONB
    public string? NomineeName { get; set; }
    public string? NomineeRelationship { get; set; }
    public string? AuditInfo { get; set; } // JSONB
    public DateTime CreatedAt { get; set; }
    public DateTime UpdatedAt { get; set; }

    // Navigation
    public BeneficiaryEntity Beneficiary { get; set; } = null!;
}
