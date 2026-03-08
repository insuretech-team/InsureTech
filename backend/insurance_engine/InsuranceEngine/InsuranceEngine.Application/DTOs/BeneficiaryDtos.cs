using InsuranceEngine.Domain.Entities;
using InsuranceEngine.Domain.Enums;

namespace InsuranceEngine.Application.DTOs;

public abstract class BeneficiaryDto
{
    public Guid BeneficiaryId { get; set; }
    public string? UserId { get; set; }
    public BeneficiaryType Type { get; set; }
    public BeneficiaryStatus Status { get; set; }
    public string? Code { get; set; }
    public string? KycStatus { get; set; }
    public DateTime? KycCompletedAt { get; set; }
    public string? RiskScore { get; set; }
    public string? ReferralCode { get; set; }
    public string? ReferredBy { get; set; }
    public string? PartnerId { get; set; }
    public Guid? PolicyId { get; set; }
    public AuditInfo? AuditInfo { get; set; }
}

public class IndividualBeneficiaryDto : BeneficiaryDto
{
    public Guid Id { get; set; }
    public string? FullName { get; set; }
    public string? FullNameBn { get; set; }
    public DateTime? DateOfBirth { get; set; }
    public BeneficiaryGender? Gender { get; set; }
    public string? NidNumber { get; set; }
    public string? PassportNumber { get; set; }
    public string? BirthCertificateNumber { get; set; }
    public string? TinNumber { get; set; }
    public MaritalStatus? MaritalStatus { get; set; }
    public string? Occupation { get; set; }
    public ContactInfo? ContactInfo { get; set; }
    public Address? PermanentAddress { get; set; }
    public Address? PresentAddress { get; set; }
    public string? NomineeName { get; set; }
    public string? NomineeRelationship { get; set; }
    public AuditInfo? IndividualAuditInfo { get; set; }
}

public class BusinessBeneficiaryDto : BeneficiaryDto
{
    public Guid Id { get; set; }
    public string? BusinessName { get; set; }
    public string? BusinessNameBn { get; set; }
    public string? TradeLicenseNumber { get; set; }
    public DateTime? TradeLicenseIssueDate { get; set; }
    public DateTime? TradeLicenseExpiryDate { get; set; }
    public string? TinNumber { get; set; }
    public string? BinNumber { get; set; }
    public BusinessType? BusinessType { get; set; }
    public string? IndustrySector { get; set; }
    public int? EmployeeCount { get; set; }
    public DateTime? IncorporationDate { get; set; }
    public ContactInfo? ContactInfo { get; set; }
    public Address? RegisteredAddress { get; set; }
    public Address? BusinessAddress { get; set; }
    public string? FocalPersonName { get; set; }
    public string? FocalPersonDesignation { get; set; }
    public string? FocalPersonNid { get; set; }
    public ContactInfo? FocalPersonContact { get; set; }
    public AuditInfo? BusinessAuditInfo { get; set; }
}
