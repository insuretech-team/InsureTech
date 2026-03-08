using System.Text.Json.Serialization;
using InsuranceEngine.Api.DTOs;

namespace InsuranceEngine.Api.ResponseModels;

public class BeneficiariesListingResponse
{
    [JsonPropertyName("beneficiaries")]
    public IReadOnlyList<Beneficiary> Beneficiaries { get; set; } = Array.Empty<Beneficiary>();

    [JsonPropertyName("total_count")]
    public int TotalCount { get; set; }

    [JsonPropertyName("error")]
    public Error Error { get; set; } = Error.None();
}

public class BeneficiaryRetrievalResponse
{
    [JsonPropertyName("beneficiary")]
    public Beneficiary Beneficiary { get; set; } = new();

    [JsonPropertyName("individual_details")]
    public IndividualBeneficiary IndividualDetails { get; set; } = new();

    [JsonPropertyName("business_details")]
    public BusinessBeneficiary BusinessDetails { get; set; } = new();

    [JsonPropertyName("error")]
    public Error Error { get; set; } = Error.None();
}

public class BeneficiaryUpdateResponse
{
    [JsonPropertyName("message")]
    public string Message { get; set; } = string.Empty;

    [JsonPropertyName("error")]
    public Error Error { get; set; } = Error.None();
}

public class BusinessBeneficiaryCreationResponse
{
    [JsonPropertyName("beneficiary_id")]
    public string BeneficiaryId { get; set; } = string.Empty;

    [JsonPropertyName("beneficiary_code")]
    public string BeneficiaryCode { get; set; } = string.Empty;

    [JsonPropertyName("message")]
    public string Message { get; set; } = string.Empty;

    [JsonPropertyName("error")]
    public Error Error { get; set; } = Error.None();
}

public class IndividualBeneficiaryCreationResponse
{
    [JsonPropertyName("beneficiary_id")]
    public string BeneficiaryId { get; set; } = string.Empty;

    [JsonPropertyName("beneficiary_code")]
    public string BeneficiaryCode { get; set; } = string.Empty;

    [JsonPropertyName("message")]
    public string Message { get; set; } = string.Empty;

    [JsonPropertyName("error")]
    public Error Error { get; set; } = Error.None();
}
