using System;
using System.Collections.Generic;
using System.Text.Json.Serialization;
using InsuranceEngine.SharedKernel.DTOs;

namespace InsuranceEngine.Commission.Application.DTOs;

public record CommissionDto(
    Guid Id,
    Guid PolicyId,
    Guid? PartnerId,
    Guid? AgentId,
    string Type,
    long Amount,
    string Currency,
    string Status,
    DateTime CreatedAt
);

public record CommissionsListingResponse(
    [property: JsonPropertyName("commissions")] List<CommissionDto> Commissions,
    [property: JsonPropertyName("total_count")] int TotalCount,
    [property: JsonPropertyName("total_amount")] MoneyDto TotalAmount,
    [property: JsonPropertyName("error")] ErrorDto? Error = null
);

public record CommissionRetrievalResponse(
    [property: JsonPropertyName("commission")] CommissionDto Commission,
    [property: JsonPropertyName("error")] ErrorDto? Error = null
);

public record CommissionCalculationResponse(
    [property: JsonPropertyName("commission_id")] Guid CommissionId,
    [property: JsonPropertyName("commission_number")] string CommissionNumber,
    [property: JsonPropertyName("amount")] MoneyDto Amount,
    [property: JsonPropertyName("calculation_breakdown")] string CalculationBreakdown,
    [property: JsonPropertyName("error")] ErrorDto? Error = null
);
