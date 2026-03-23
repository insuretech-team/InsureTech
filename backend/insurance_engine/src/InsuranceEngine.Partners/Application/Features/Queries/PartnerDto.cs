using System;
using System.Collections.Generic;

namespace InsuranceEngine.Partners.Application.Features.Queries;

public record PartnerDto(
    Guid Id,
    string Name,
    string Code,
    string Email,
    string? Phone,
    string? Address,
    string Status,
    List<AgentDto> Agents
);

public record AgentDto(
    Guid Id,
    string Name,
    string Code,
    string Email,
    string? Phone,
    string Status
);
