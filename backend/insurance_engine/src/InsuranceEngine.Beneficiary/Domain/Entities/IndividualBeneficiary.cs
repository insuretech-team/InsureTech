using InsuranceEngine.Beneficiary.Domain.Enums;
using InsuranceEngine.SharedKernel.Domain;
using InsuranceEngine.SharedKernel.Domain.ValueObjects;

namespace InsuranceEngine.Beneficiary.Domain.Entities;

public class IndividualBeneficiary : Entity<Guid>
{
    public Guid BeneficiaryId { get; private set; }
    
    public string FullName { get; private set; } = string.Empty;
    public string? FullNameBn { get; set; }
    public DateTime DateOfBirth { get; private set; }
    public BeneficiaryGender Gender { get; set; }
    
    public string? NidNumber { get; set; }
    public string? PassportNumber { get; set; }
    public string? BirthCertificateNumber { get; set; }
    public string? TinNumber { get; set; }
    
    public MaritalStatus MaritalStatus { get; set; }
    public string? Occupation { get; set; }
    
    public string? NomineeName { get; set; }
    public string? NomineeRelationship { get; set; }
 
    public ContactInfo ContactInfo { get; set; } = new();
    public Address PermanentAddress { get; set; } = new();
    public Address PresentAddress { get; set; } = new();
    public AuditInfo AuditInfo { get; private set; } = new();

    // EF Core constructor
    public IndividualBeneficiary() { }

    public IndividualBeneficiary(
        Guid id,
        Guid beneficiaryId,
        string fullName,
        DateTime dob,
        BeneficiaryGender gender) : base(id)
    {
        BeneficiaryId = beneficiaryId;
        FullName = fullName;
        DateOfBirth = dob;
        Gender = gender;
        AuditInfo = new AuditInfo();
    }
}
