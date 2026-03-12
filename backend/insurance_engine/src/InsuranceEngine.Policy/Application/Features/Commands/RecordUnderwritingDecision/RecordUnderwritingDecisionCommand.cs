using System;
using System.Collections.Generic;
using InsuranceEngine.Policy.Application.DTOs;
using InsuranceEngine.Policy.Domain.Enums;
using InsuranceEngine.SharedKernel.CQRS;
using MediatR;

namespace InsuranceEngine.Policy.Application.Features.Commands.RecordUnderwritingDecision;

public record RecordUnderwritingDecisionCommand(
    Guid QuoteId,
    DecisionType Decision,
    DecisionMethod Method,
    decimal RiskScore,
    RiskLevel RiskLevel,
    List<string>? RiskFactors,
    string? Reason,
    List<string>? Conditions,
    bool IsPremiumAdjusted,
    long? AdjustedPremiumAmount,
    Guid? UnderwriterId,
    string? UnderwriterComments
) : IRequest<Result<UnderwritingDecisionDto>>;
