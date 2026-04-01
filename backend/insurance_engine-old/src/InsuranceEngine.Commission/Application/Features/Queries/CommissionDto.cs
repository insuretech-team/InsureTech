using System;

namespace InsuranceEngine.Commission.Application.Features.Queries;

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
