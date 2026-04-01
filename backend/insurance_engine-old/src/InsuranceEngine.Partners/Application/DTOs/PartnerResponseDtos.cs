using System;
using System.Collections.Generic;
using System.Text.Json.Serialization;
using InsuranceEngine.SharedKernel.DTOs;

namespace InsuranceEngine.Partners.Application.DTOs;

public record PartnerDto(
    Guid Id,
    string OrganizationName,
    string Code,
    string Email,
    string? Phone,
    string Status
);

public record PartnerCreationResponse(
    [property: JsonPropertyName("partner_id")] Guid PartnerId,
    [property: JsonPropertyName("partner")] PartnerDto Partner,
    [property: JsonPropertyName("error")] ErrorDto? Error = null
);

public record PartnerRetrievalResponse(
    [property: JsonPropertyName("partner")] PartnerDto Partner,
    [property: JsonPropertyName("error")] ErrorDto? Error = null
);

public record PartnersListingResponse(
    [property: JsonPropertyName("partners")] List<PartnerDto> Partners,
    [property: JsonPropertyName("total_count")] int TotalCount,
    [property: JsonPropertyName("error")] ErrorDto? Error = null
);

public record PartnerUpdateResponse(
    [property: JsonPropertyName("message")] string Message,
    [property: JsonPropertyName("error")] ErrorDto? Error = null
);

public record PartnerStatusUpdateResponse(
    [property: JsonPropertyName("message")] string Message,
    [property: JsonPropertyName("error")] ErrorDto? Error = null
);

public record PartnerCredentialRetrievalResponse(
    [property: JsonPropertyName("credentials")] object Credentials,
    [property: JsonPropertyName("error")] ErrorDto? Error = null
);

public record PartnerCredentialRotationResponse(
    [property: JsonPropertyName("credentials")] object Credentials,
    [property: JsonPropertyName("message")] string Message,
    [property: JsonPropertyName("error")] ErrorDto? Error = null
);

public record PartnerVerificationResponse(
    [property: JsonPropertyName("verified")] bool Verified,
    [property: JsonPropertyName("message")] string Message,
    [property: JsonPropertyName("error")] ErrorDto? Error = null
);

public record PartnerCommissionRetrievalResponse(
    [property: JsonPropertyName("commission_config")] object CommissionConfig,
    [property: JsonPropertyName("error")] ErrorDto? Error = null
);

public record PartnerCommissionUpdateResponse(
    [property: JsonPropertyName("message")] string Message,
    [property: JsonPropertyName("error")] ErrorDto? Error = null
);
