using System.ComponentModel.DataAnnotations;
using InsuranceEngine.Domain.Enums;

namespace InsuranceEngine.Domain.Entities;

public class BeneficiaryBusiness
{
    [Key]
    public Guid Id { get; set; }

    public Guid BeneficiaryId { get; set; }

    [MaxLength(255)]
    [Required]
    public string BusinessName { get; set; } = string.Empty;

    [MaxLength(255)]
    public string? BusinessNameBn { get; set; }

    [MaxLength(100)]
    [Required]
    public string TradeLicenseNumber { get; set; } = string.Empty;

    public DateTime? TradeLicenseIssueDate { get; set; }

    public DateTime? TradeLicenseExpiryDate { get; set; }

    [MaxLength(100)]
    [Required]
    public string TinNumber { get; set; } = string.Empty;

    [MaxLength(100)]
    public string? BinNumber { get; set; }

    [Required]
    public BusinessType BusinessType { get; set; }

    [MaxLength(150)]
    public string? IndustrySector { get; set; }

    public int? EmployeeCount { get; set; }

    public DateTime? IncorporationDate { get; set; }

    public ContactInfo? ContactInfo { get; set; }

    public Address? RegisteredAddress { get; set; }

    public Address? BusinessAddress { get; set; }

    [MaxLength(255)]
    [Required]
    public string FocalPersonName { get; set; } = string.Empty;

    [MaxLength(150)]
    public string? FocalPersonDesignation { get; set; }

    [MaxLength(100)]
    public string? FocalPersonNid { get; set; }

    public ContactInfo? FocalPersonContact { get; set; }

    public AuditInfo AuditInfo { get; set; } = new();

    public Beneficiary Beneficiary { get; set; } = null!;
}
