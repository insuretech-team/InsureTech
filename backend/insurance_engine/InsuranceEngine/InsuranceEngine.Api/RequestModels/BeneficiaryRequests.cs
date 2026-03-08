using System.ComponentModel.DataAnnotations;
using System.Text.Json.Serialization;

namespace InsuranceEngine.Api.RequestModels;

public class BeneficiariesListingRequest
{
    [Required]
    [JsonPropertyName("type")]
    public string Type { get; set; } = string.Empty;

    [JsonPropertyName("status")]
    public string? Status { get; set; }

    [Required]
    [Range(1, int.MaxValue)]
    [JsonPropertyName("page")]
    public int Page { get; set; }

    [Required]
    [Range(1, 100)]
    [JsonPropertyName("page_size")]
    public int PageSize { get; set; }
}

public class BeneficiaryRetrievalRequest
{
    [Required]
    [JsonPropertyName("beneficiary_id")]
    public string BeneficiaryId { get; set; } = string.Empty;
}

public class BeneficiaryUpdateRequest
{
    [Required]
    [JsonPropertyName("beneficiary_id")]
    public string BeneficiaryId { get; set; } = string.Empty;

    [JsonPropertyName("mobile_number")]
    public string? MobileNumber { get; set; }

    [Required]
    [EmailAddress]
    [JsonPropertyName("email")]
    public string Email { get; set; } = string.Empty;

    [JsonPropertyName("address")]
    public string? Address { get; set; }
}

public class BusinessBeneficiaryCreationRequest
{
    [Required]
    [JsonPropertyName("user_id")]
    public string UserId { get; set; } = string.Empty;

    [JsonPropertyName("business_name")]
    public string? BusinessName { get; set; }

    [JsonPropertyName("trade_license_number")]
    public string? TradeLicenseNumber { get; set; }

    [JsonPropertyName("tin_number")]
    public string? TinNumber { get; set; }

    [JsonPropertyName("focal_person_name")]
    public string? FocalPersonName { get; set; }

    [JsonPropertyName("focal_person_mobile")]
    public string? FocalPersonMobile { get; set; }

    [Required]
    [JsonPropertyName("partner_id")]
    public string PartnerId { get; set; } = string.Empty;
}

public class IndividualBeneficiaryCreationRequest
{
    [Required]
    [JsonPropertyName("user_id")]
    public string UserId { get; set; } = string.Empty;

    [JsonPropertyName("full_name")]
    public string? FullName { get; set; }

    [JsonPropertyName("date_of_birth")]
    public string? DateOfBirth { get; set; }

    [JsonPropertyName("gender")]
    public string? Gender { get; set; }

    [JsonPropertyName("nid_number")]
    public string? NidNumber { get; set; }

    [JsonPropertyName("mobile_number")]
    public string? MobileNumber { get; set; }

    [Required]
    [EmailAddress]
    [JsonPropertyName("email")]
    public string Email { get; set; } = string.Empty;

    [Required]
    [JsonPropertyName("partner_id")]
    public string PartnerId { get; set; } = string.Empty;
}
