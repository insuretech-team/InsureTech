using PoliSync.SharedKernel.Domain;

namespace PoliSync.Beneficiaries.Domain;

public class IndividualBeneficiary : Entity
{
    public Guid BeneficiaryId { get; private set; } // PK and FK to Beneficiary
    public string FullName { get; private set; } = string.Empty;
    public string? FullNameBn { get; private set; }
    public DateTime DateOfBirth { get; private set; }
    public Gender Gender { get; private set; }
    public string? NidNumber { get; private set; }
    public string? PassportNumber { get; private set; }
    public string? BirthCertificateNumber { get; private set; }
    public string? TinNumber { get; private set; }
    public MaritalStatus MaritalStatus { get; private set; }
    public string? Occupation { get; private set; }
    public string? ContactInfo { get; private set; } // JSONB
    public string? PermanentAddress { get; private set; } // JSONB
    public string? PresentAddress { get; private set; } // JSONB
    public string? NomineeName { get; private set; }
    public string? NomineeRelationship { get; private set; }
    public string? AuditInfo { get; private set; } // JSONB

    private IndividualBeneficiary() { }

    public static IndividualBeneficiary Create(
        Guid beneficiaryId,
        string fullName,
        DateTime dateOfBirth,
        Gender gender,
        string? nidNumber = null,
        string? contactInfo = null,
        string? permanentAddress = null)
    {
        return new IndividualBeneficiary
        {
            BeneficiaryId = beneficiaryId,
            FullName = fullName,
            DateOfBirth = dateOfBirth,
            Gender = gender,
            NidNumber = nidNumber,
            ContactInfo = contactInfo ?? "{}",
            PermanentAddress = permanentAddress ?? "{}",
            AuditInfo = "{}"
        };
    }

    public void UpdateDetails(
        string? fullName = null,
        string? contactInfo = null,
        string? presentAddress = null)
    {
        if (fullName is not null) FullName = fullName;
        if (contactInfo is not null) ContactInfo = contactInfo;
        if (presentAddress is not null) PresentAddress = presentAddress;
    }
}
