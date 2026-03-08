using System.ComponentModel.DataAnnotations;
using InsuranceEngine.Domain.Enums;

namespace InsuranceEngine.Domain.Entities;

public class BeneficiaryIndividual
{
    [Key]
    public Guid Id { get; set; }

    public Guid BeneficiaryId { get; set; }

    [MaxLength(200)]
    [Required]
    public string FullName { get; set; } = string.Empty;

    [MaxLength(200)]
    public string? FullNameBn { get; set; }

    [Required]
    public DateTime DateOfBirth { get; set; }

    [Required]
    public BeneficiaryGender Gender { get; set; }

    [MaxLength(100)]
    public string? NidNumber { get; set; }

    [MaxLength(100)]
    public string? PassportNumber { get; set; }

    [MaxLength(100)]
    public string? BirthCertificateNumber { get; set; }

    [MaxLength(100)]
    public string? TinNumber { get; set; }

    public MaritalStatus? MaritalStatus { get; set; }

    [MaxLength(100)]
    public string? Occupation { get; set; }

    public ContactInfo? ContactInfo { get; set; }

    public Address? PermanentAddress { get; set; }

    public Address? PresentAddress { get; set; }

    [MaxLength(200)]
    public string? NomineeName { get; set; }

    [MaxLength(200)]
    public string? NomineeRelationship { get; set; }

    public AuditInfo AuditInfo { get; set; } = new();

    public Beneficiary Beneficiary { get; set; } = null!;
}
