using System.ComponentModel.DataAnnotations;

using InsuranceEngine.Domain.Enums;

namespace InsuranceEngine.Application.DTOs;

public class CreateBeneficiaryRequest
{
    [Required]
    public BeneficiaryType Type { get; set; }

    [Required]
    public BeneficiaryStatus Status { get; set; }

    public BeneficiaryGender? Gender { get; set; }

    [EmailAddress]
    [MaxLength(255)]
    public string? Email { get; set; }

    [MaxLength(50)]
    public string? Phone { get; set; }

    [Required]
    public Guid PolicyId { get; set; }

    public string? FirstName { get; set; }
    public string? LastName { get; set; }
    public DateTime? DateOfBirth { get; set; }
    public string? NationalIdNumber { get; set; }

    public string? BusinessName { get; set; }
    public string? BusinessRegistrationNumber { get; set; }
    public string? ContactPersonName { get; set; }
    public string? TaxIdentificationNumber { get; set; }
}

public class CreateIndividualBeneficiaryRequest
{
    [Required]
    public BeneficiaryStatus Status { get; set; }

    public BeneficiaryGender? Gender { get; set; }

    [EmailAddress]
    [MaxLength(255)]
    public string? Email { get; set; }

    [MaxLength(50)]
    public string? Phone { get; set; }

    [Required]
    public Guid PolicyId { get; set; }

    [Required]
    [MaxLength(100)]
    public string? FirstName { get; set; }

    [Required]
    [MaxLength(100)]
    public string? LastName { get; set; }

    public DateTime? DateOfBirth { get; set; }

    [MaxLength(100)]
    public string? NationalIdNumber { get; set; }
}

public class CreateBusinessBeneficiaryRequest
{
    [Required]
    public BeneficiaryStatus Status { get; set; }

    [EmailAddress]
    [MaxLength(255)]
    public string? Email { get; set; }

    [MaxLength(50)]
    public string? Phone { get; set; }

    [Required]
    public Guid PolicyId { get; set; }

    [Required]
    [MaxLength(255)]
    public string? BusinessName { get; set; }

    [MaxLength(100)]
    public string? BusinessRegistrationNumber { get; set; }

    [MaxLength(255)]
    public string? ContactPersonName { get; set; }

    [MaxLength(100)]
    public string? TaxIdentificationNumber { get; set; }
}

public class UpdateBeneficiaryRequest
{
    [Required]
    public BeneficiaryStatus Status { get; set; }

    public BeneficiaryGender? Gender { get; set; }

    [EmailAddress]
    [MaxLength(255)]
    public string? Email { get; set; }

    [MaxLength(50)]
    public string? Phone { get; set; }

    public string? FirstName { get; set; }
    public string? LastName { get; set; }
    public DateTime? DateOfBirth { get; set; }
    public string? NationalIdNumber { get; set; }

    public string? BusinessName { get; set; }
    public string? BusinessRegistrationNumber { get; set; }
    public string? ContactPersonName { get; set; }
    public string? TaxIdentificationNumber { get; set; }
}
