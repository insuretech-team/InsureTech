using System.Text.Json.Serialization;
using InsuranceEngine.Api.DTOs;
using InsuranceEngine.Domain.Enums;

namespace InsuranceEngine.Api.ResponseModels;

public class CreateBeneficiaryResponseV1
{
    [JsonPropertyName("beneficiary_id")]
    public string? BeneficiaryId { get; set; }

    [JsonPropertyName("beneficiary_code")]
    public string? BeneficiaryCode { get; set; }

    [JsonPropertyName("message")]
    public string? Message { get; set; }

    [JsonPropertyName("error")]
    public Error? Error { get; set; }
}

public class UpdateBeneficiaryResponseV1
{
    [JsonPropertyName("message")]
    public string? Message { get; set; }

    [JsonPropertyName("error")]
    public Error? Error { get; set; }
}

public class ListBeneficiariesResponseV1
{
    [JsonPropertyName("beneficiaries")]
    public IReadOnlyList<BeneficiarySummaryV1> Beneficiaries { get; set; } = Array.Empty<BeneficiarySummaryV1>();

    [JsonPropertyName("total_count")]
    public int TotalCount { get; set; }

    [JsonPropertyName("error")]
    public Error? Error { get; set; }
}

public class BeneficiarySummaryV1
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

public class GetBeneficiaryResponseV1
{
    [JsonPropertyName("beneficiary")]
    public BeneficiaryProfileV1? Beneficiary { get; set; }

    [JsonPropertyName("individual_details")]
    public IndividualDetailsV1? IndividualDetails { get; set; }

    [JsonPropertyName("business_details")]
    public BusinessDetailsV1? BusinessDetails { get; set; }

    [JsonPropertyName("error")]
    public Error? Error { get; set; }
}

public class AuditInfoV1
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

public class ContactInfoV1
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

public class AddressV1
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

public class BeneficiaryProfileV1
{
    [JsonPropertyName("beneficiary_id")]
    public string? BeneficiaryId { get; set; }

    [JsonPropertyName("user_id")]
    public string? UserId { get; set; }

    [JsonPropertyName("type")]
    public BeneficiaryType Type { get; set; }

    [JsonPropertyName("code")]
    public string? Code { get; set; }

    [JsonPropertyName("status")]
    public BeneficiaryStatus Status { get; set; }

    [JsonPropertyName("kyc_status")]
    public string? KycStatus { get; set; }

    [JsonPropertyName("kyc_completed_at")]
    public DateTime? KycCompletedAt { get; set; }

    [JsonPropertyName("risk_score")]
    public string? RiskScore { get; set; }

    [JsonPropertyName("referral_code")]
    public string? ReferralCode { get; set; }

    [JsonPropertyName("referred_by")]
    public string? ReferredBy { get; set; }

    [JsonPropertyName("partner_id")]
    public string? PartnerId { get; set; }

    [JsonPropertyName("audit_info")]
    public AuditInfoV1? AuditInfo { get; set; }
}

public class IndividualDetailsV1
{
    [JsonPropertyName("id")]
    public string? Id { get; set; }

    [JsonPropertyName("beneficiary_id")]
    public string? BeneficiaryId { get; set; }

    [JsonPropertyName("full_name")]
    public string? FullName { get; set; }

    [JsonPropertyName("full_name_bn")]
    public string? FullNameBn { get; set; }

    [JsonPropertyName("date_of_birth")]
    public DateTime DateOfBirth { get; set; }

    [JsonPropertyName("gender")]
    public BeneficiaryGender Gender { get; set; }

    [JsonPropertyName("nid_number")]
    public string? NidNumber { get; set; }

    [JsonPropertyName("passport_number")]
    public string? PassportNumber { get; set; }

    [JsonPropertyName("birth_certificate_number")]
    public string? BirthCertificateNumber { get; set; }

    [JsonPropertyName("tin_number")]
    public string? TinNumber { get; set; }

    [JsonPropertyName("marital_status")]
    public MaritalStatus? MaritalStatus { get; set; }

    [JsonPropertyName("occupation")]
    public string? Occupation { get; set; }

    [JsonPropertyName("contact_info")]
    public ContactInfoV1? ContactInfo { get; set; }

    [JsonPropertyName("permanent_address")]
    public AddressV1? PermanentAddress { get; set; }

    [JsonPropertyName("present_address")]
    public AddressV1? PresentAddress { get; set; }

    [JsonPropertyName("nominee_name")]
    public string? NomineeName { get; set; }

    [JsonPropertyName("nominee_relationship")]
    public string? NomineeRelationship { get; set; }

    [JsonPropertyName("audit_info")]
    public AuditInfoV1? AuditInfo { get; set; }
}

public class BusinessDetailsV1
{
    [JsonPropertyName("id")]
    public string? Id { get; set; }

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
    public BusinessType BusinessType { get; set; }

    [JsonPropertyName("industry_sector")]
    public string? IndustrySector { get; set; }

    [JsonPropertyName("employee_count")]
    public int? EmployeeCount { get; set; }

    [JsonPropertyName("incorporation_date")]
    public DateTime? IncorporationDate { get; set; }

    [JsonPropertyName("contact_info")]
    public ContactInfoV1? ContactInfo { get; set; }

    [JsonPropertyName("registered_address")]
    public AddressV1? RegisteredAddress { get; set; }

    [JsonPropertyName("business_address")]
    public AddressV1? BusinessAddress { get; set; }

    [JsonPropertyName("focal_person_name")]
    public string? FocalPersonName { get; set; }

    [JsonPropertyName("focal_person_designation")]
    public string? FocalPersonDesignation { get; set; }

    [JsonPropertyName("focal_person_nid")]
    public string? FocalPersonNid { get; set; }

    [JsonPropertyName("focal_person_contact")]
    public ContactInfoV1? FocalPersonContact { get; set; }

    [JsonPropertyName("audit_info")]
    public AuditInfoV1? AuditInfo { get; set; }
}
