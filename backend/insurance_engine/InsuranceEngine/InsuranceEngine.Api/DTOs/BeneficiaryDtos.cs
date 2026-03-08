using System.Text.Json.Serialization;
using InsuranceEngine.Domain.Enums;

namespace InsuranceEngine.Api.DTOs;

public class Beneficiary
{
    [JsonPropertyName("beneficiary_id")]
    public string? BeneficiaryId { get; set; }

    [JsonPropertyName("beneficiary_code")]
    public string? BeneficiaryCode { get; set; }

    [JsonPropertyName("user_id")]
    public string? UserId { get; set; }

    [JsonPropertyName("partner_id")]
    public string? PartnerId { get; set; }

    [JsonPropertyName("type")]
    public BeneficiaryType Type { get; set; }

    [JsonPropertyName("status")]
    public BeneficiaryStatus Status { get; set; }
}

public class ContactInfo
{
    [JsonPropertyName("mobile_number")]
    public string? MobileNumber { get; set; }

    [JsonPropertyName("email")]
    public string? Email { get; set; }

    [JsonPropertyName("alternate_mobile")]
    public string? AlternateMobile { get; set; }

    [JsonPropertyName("landline")]
    public string? Landline { get; set; }
}

public class Address
{
    [JsonPropertyName("address_line1")]
    public string? AddressLine1 { get; set; }

    [JsonPropertyName("address_line2")]
    public string? AddressLine2 { get; set; }

    [JsonPropertyName("city")]
    public string? City { get; set; }

    [JsonPropertyName("district")]
    public string? District { get; set; }

    [JsonPropertyName("division")]
    public string? Division { get; set; }

    [JsonPropertyName("postal_code")]
    public string? PostalCode { get; set; }

    [JsonPropertyName("country")]
    public string? Country { get; set; }

    [JsonPropertyName("latitude")]
    public decimal? Latitude { get; set; }

    [JsonPropertyName("longitude")]
    public decimal? Longitude { get; set; }
}

public class AuditInfo
{
    [JsonPropertyName("created_at")]
    public DateTime CreatedAt { get; set; }

    [JsonPropertyName("updated_at")]
    public DateTime UpdatedAt { get; set; }

    [JsonPropertyName("created_by")]
    public string? CreatedBy { get; set; }

    [JsonPropertyName("updated_by")]
    public string? UpdatedBy { get; set; }

    [JsonPropertyName("deleted_at")]
    public DateTime? DeletedAt { get; set; }

    [JsonPropertyName("deleted_by")]
    public string? DeletedBy { get; set; }
}

public class IndividualBeneficiary
{
    [JsonPropertyName("beneficiary_id")]
    public string? BeneficiaryId { get; set; }

    [JsonPropertyName("full_name")]
    public string? FullName { get; set; }

    [JsonPropertyName("full_name_bn")]
    public string? FullNameBn { get; set; }

    [JsonPropertyName("date_of_birth")]
    public string? DateOfBirth { get; set; }

    [JsonPropertyName("gender")]
    public string? Gender { get; set; }

    [JsonPropertyName("nid_number")]
    public string? NidNumber { get; set; }

    [JsonPropertyName("passport_number")]
    public string? PassportNumber { get; set; }

    [JsonPropertyName("birth_certificate_number")]
    public string? BirthCertificateNumber { get; set; }

    [JsonPropertyName("tin_number")]
    public string? TinNumber { get; set; }

    [JsonPropertyName("marital_status")]
    public string? MaritalStatus { get; set; }

    [JsonPropertyName("occupation")]
    public string? Occupation { get; set; }

    [JsonPropertyName("contact_info")]
    public ContactInfo? ContactInfo { get; set; }

    [JsonPropertyName("permanent_address")]
    public Address? PermanentAddress { get; set; }

    [JsonPropertyName("present_address")]
    public Address? PresentAddress { get; set; }

    [JsonPropertyName("nominee_name")]
    public string? NomineeName { get; set; }

    [JsonPropertyName("nominee_relationship")]
    public string? NomineeRelationship { get; set; }

    [JsonPropertyName("audit_info")]
    public AuditInfo? AuditInfo { get; set; }
}

public class BusinessBeneficiary
{
    [JsonPropertyName("beneficiary_id")]
    public string? BeneficiaryId { get; set; }

    [JsonPropertyName("business_name")]
    public string? BusinessName { get; set; }

    [JsonPropertyName("business_name_bn")]
    public string? BusinessNameBn { get; set; }

    [JsonPropertyName("trade_license_number")]
    public string? TradeLicenseNumber { get; set; }

    [JsonPropertyName("trade_license_issue_date")]
    public DateTime? TradeLicenseIssueDate { get; set; }

    [JsonPropertyName("trade_license_expiry_date")]
    public DateTime? TradeLicenseExpiryDate { get; set; }

    [JsonPropertyName("tin_number")]
    public string? TinNumber { get; set; }

    [JsonPropertyName("bin_number")]
    public string? BinNumber { get; set; }

    [JsonPropertyName("business_type")]
    public string? BusinessType { get; set; }

    [JsonPropertyName("industry_sector")]
    public string? IndustrySector { get; set; }

    [JsonPropertyName("employee_count")]
    public int? EmployeeCount { get; set; }

    [JsonPropertyName("incorporation_date")]
    public DateTime? IncorporationDate { get; set; }

    [JsonPropertyName("contact_info")]
    public ContactInfo? ContactInfo { get; set; }

    [JsonPropertyName("registered_address")]
    public Address? RegisteredAddress { get; set; }

    [JsonPropertyName("business_address")]
    public Address? BusinessAddress { get; set; }

    [JsonPropertyName("focal_person_name")]
    public string? FocalPersonName { get; set; }

    [JsonPropertyName("focal_person_designation")]
    public string? FocalPersonDesignation { get; set; }

    [JsonPropertyName("focal_person_nid")]
    public string? FocalPersonNid { get; set; }

    [JsonPropertyName("focal_person_contact")]
    public ContactInfo? FocalPersonContact { get; set; }

    [JsonPropertyName("audit_info")]
    public AuditInfo? AuditInfo { get; set; }
}
