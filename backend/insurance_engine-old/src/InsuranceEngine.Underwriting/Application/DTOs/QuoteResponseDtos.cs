using System;
using System.Collections.Generic;
using System.Text.Json.Serialization;
using InsuranceEngine.SharedKernel.DTOs;

namespace InsuranceEngine.Underwriting.Application.DTOs;

public record RequestQuoteResponse(
    [property: JsonPropertyName("quote_id")] Guid QuoteId,
    [property: JsonPropertyName("quote_number")] string QuoteNumber,
    [property: JsonPropertyName("base_premium")] MoneyDto BasePremium,
    [property: JsonPropertyName("total_premium")] MoneyDto TotalPremium,
    [property: JsonPropertyName("valid_until")] DateTime ValidUntil,
    [property: JsonPropertyName("message")] string Message,
    [property: JsonPropertyName("error")] ErrorDto? Error = null
);

public record QuoteRetrievalResponse(
    [property: JsonPropertyName("quote")] QuoteDto Quote,
    [property: JsonPropertyName("error")] ErrorDto? Error = null
);

public record QuotesListingResponse(
    [property: JsonPropertyName("quotes")] List<QuoteDto> Quotes,
    [property: JsonPropertyName("total_count")] int TotalCount,
    [property: JsonPropertyName("error")] ErrorDto? Error = null
);

public record UnderwritingApprovalResponse(
    [property: JsonPropertyName("decision_id")] Guid DecisionId,
    [property: JsonPropertyName("message")] string Message,
    [property: JsonPropertyName("error")] ErrorDto? Error = null
);

public record UnderwritingDecisionRetrievalResponse(
    [property: JsonPropertyName("decision")] UnderwritingDecisionDto Decision,
    [property: JsonPropertyName("error")] ErrorDto? Error = null
);

public record HealthDeclarationRetrievalResponse(
    [property: JsonPropertyName("health_declaration")] UnderwritingHealthDeclarationResponseDto HealthDeclaration,
    [property: JsonPropertyName("error")] ErrorDto? Error = null
);
