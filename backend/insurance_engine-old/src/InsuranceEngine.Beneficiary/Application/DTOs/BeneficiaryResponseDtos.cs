using System;
using System.Collections.Generic;
using System.Text.Json.Serialization;
using InsuranceEngine.SharedKernel.DTOs;

namespace InsuranceEngine.Beneficiary.Application.DTOs;

public record IndividualBeneficiaryCreationResponse(
    [property: JsonPropertyName("beneficiary_id")] Guid BeneficiaryId,
    [property: JsonPropertyName("beneficiary_code")] string BeneficiaryCode,
    [property: JsonPropertyName("message")] string Message,
    [property: JsonPropertyName("error")] ErrorDto? Error = null
);

public record BusinessBeneficiaryCreationResponse(
    [property: JsonPropertyName("beneficiary_id")] Guid BeneficiaryId,
    [property: JsonPropertyName("beneficiary_code")] string BeneficiaryCode,
    [property: JsonPropertyName("message")] string Message,
    [property: JsonPropertyName("error")] ErrorDto? Error = null
);

public record BeneficiariesListingResponse(
    [property: JsonPropertyName("beneficiaries")] List<BeneficiaryDto> Beneficiaries,
    [property: JsonPropertyName("total_count")] int TotalCount,
    [property: JsonPropertyName("error")] ErrorDto? Error = null
);

public record BeneficiaryRetrievalResponse(
    [property: JsonPropertyName("beneficiary")] BeneficiaryDto? Beneficiary = null,
    [property: JsonPropertyName("individual_details")] IndividualBeneficiaryDto? IndividualDetails = null,
    [property: JsonPropertyName("business_details")] BusinessBeneficiaryDto? BusinessDetails = null,
    [property: JsonPropertyName("error")] ErrorDto? Error = null
);

public record BeneficiaryUpdateResponse(
    [property: JsonPropertyName("message")] string Message,
    [property: JsonPropertyName("error")] ErrorDto? Error = null
);

public record KYCCompletionResponse(
    [property: JsonPropertyName("kyc_status")] string KycStatus,
    [property: JsonPropertyName("message")] string Message,
    [property: JsonPropertyName("error")] ErrorDto? Error = null
);

public record QuotesListingResponse(
    [property: JsonPropertyName("quotes")] List<BeneficiaryQuoteDto> Quotes,
    [property: JsonPropertyName("total_count")] int TotalCount,
    [property: JsonPropertyName("error")] ErrorDto? Error = null
);

public record BeneficiaryQuoteDto(
    [property: JsonPropertyName("quote_id")] Guid QuoteId,
    [property: JsonPropertyName("quote_number")] string QuoteNumber,
    [property: JsonPropertyName("status")] string Status,
    [property: JsonPropertyName("total_premium")] MoneyDto TotalPremium,
    [property: JsonPropertyName("valid_until")] DateTime ValidUntil
);

public record RiskScoreUpdateResponse(
    [property: JsonPropertyName("message")] string Message,
    [property: JsonPropertyName("error")] ErrorDto? Error = null
);
