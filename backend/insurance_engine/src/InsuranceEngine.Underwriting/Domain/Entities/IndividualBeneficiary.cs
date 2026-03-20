using System;
using InsuranceEngine.Underwriting.Domain.Enums;

namespace InsuranceEngine.Underwriting.Domain.Entities;

public class IndividualBeneficiary
{
    public Guid Id { get; set; }
    public Guid BeneficiaryId { get; set; }
    public Beneficiary? Beneficiary { get; set; }

    public string FullName { get; set; } = string.Empty;
    public string? FullNameBn { get; set; }
    public DateTime DateOfBirth { get; set; }
    public Gender Gender { get; set; }
    
    public string? NidNumber { get; set; }
    public string? PassportNumber { get; set; }
    public string? BirthCertificateNumber { get; set; }
    public string? TinNumber { get; set; }
    
    public MaritalStatus MaritalStatus { get; set; }
    public string? Occupation { get; set; }
    public string? FatherName { get; set; }
    public string? MotherName { get; set; }
    public decimal MonthlyIncome { get; set; }

    public string? ContactInfoJson { get; set; }
    public string? PermanentAddressJson { get; set; }
    public string? PresentAddressJson { get; set; }

    public string? NomineeName { get; set; }
    public string? NomineeRelationship { get; set; }

    public string? AuditInfo { get; set; }
    
    public DateTime CreatedAt { get; set; }
    public DateTime UpdatedAt { get; set; }
}
