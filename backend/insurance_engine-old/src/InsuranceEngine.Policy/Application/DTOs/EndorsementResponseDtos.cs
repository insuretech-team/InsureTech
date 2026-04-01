using System;
using System.Collections.Generic;
using System.Text.Json.Serialization;
using InsuranceEngine.SharedKernel.DTOs;
using InsuranceEngine.Policy.Domain.Enums;

namespace InsuranceEngine.Policy.Application.DTOs;

public record EndorsementDto(
    [property: JsonPropertyName("id")] Guid Id,
    [property: JsonPropertyName("endorsement_number")] string EndorsementNumber,
    [property: JsonPropertyName("policy_id")] Guid PolicyId,
    [property: JsonPropertyName("type")] EndorsementType Type,
    [property: JsonPropertyName("reason")] string Reason,
    [property: JsonPropertyName("status")] EndorsementStatus Status,
    [property: JsonPropertyName("premium_adjustment_amount")] decimal PremiumAdjustmentAmount,
    [property: JsonPropertyName("premium_adjustment_currency")] string PremiumAdjustmentCurrency,
    [property: JsonPropertyName("effective_date")] DateTime EffectiveDate,
    [property: JsonPropertyName("created_at")] DateTime CreatedAt
);

public record EndorsementsListingResponse(
    [property: JsonPropertyName("endorsements")] List<EndorsementDto> Endorsements,
    [property: JsonPropertyName("total_count")] int TotalCount,
    [property: JsonPropertyName("error")] ErrorDto? Error = null
);

public record EndorsementRetrievalResponse(
    [property: JsonPropertyName("endorsement")] EndorsementDto Endorsement,
    [property: JsonPropertyName("error")] ErrorDto? Error = null
);

public record EndorsementSubmissionResponse(
    [property: JsonPropertyName("endorsement_id")] Guid EndorsementId,
    [property: JsonPropertyName("endorsement_number")] string EndorsementNumber,
    [property: JsonPropertyName("message")] string Message,
    [property: JsonPropertyName("error")] ErrorDto? Error = null
);

public record EndorsementApprovalResponse(
    [property: JsonPropertyName("message")] string Message,
    [property: JsonPropertyName("error")] ErrorDto? Error = null
);
