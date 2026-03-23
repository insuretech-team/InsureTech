using System;
using InsuranceEngine.Beneficiary.Domain.Enums;
using InsuranceEngine.SharedKernel.Domain;

namespace InsuranceEngine.Beneficiary.Domain.Entities;

public class IndividualBeneficiary : Entity<Guid>
{
    public Guid BeneficiaryId { get; private set; }
    
    public string FullName { get; private set; } = string.Empty;
    public string? FullNameBn { get; set; }
    public DateTime DateOfBirth { get; private set; }
    public Gender Gender { get; private set; }
    
    public string? NidNumber { get; set; }
    public string? PassportNumber { get; set; }
    public string? BirthCertificateNumber { get; set; }
    public string? TinNumber { get; set; }
    
    public MaritalStatus MaritalStatus { get; set; }
    public string? Occupation { get; set; }
    
    public string? NomineeName { get; set; }
    public string? NomineeRelationship { get; set; }

    public string? ContactInfoJson { get; set; }
    public string? PermanentAddressJson { get; set; }
    public string? PresentAddressJson { get; set; }
    public string? AuditInfoJson { get; set; }

    public DateTime CreatedAt { get; private set; }
    public DateTime UpdatedAt { get; private set; }

    // EF Core constructor
    public IndividualBeneficiary() { }

    public IndividualBeneficiary(
        Guid id,
        Guid beneficiaryId,
        string fullName,
        DateTime dob,
        Gender gender) : base(id)
    {
        BeneficiaryId = beneficiaryId;
        FullName = fullName;
        DateOfBirth = dob;
        Gender = gender;
        CreatedAt = DateTime.UtcNow;
        UpdatedAt = DateTime.UtcNow;
    }
}
