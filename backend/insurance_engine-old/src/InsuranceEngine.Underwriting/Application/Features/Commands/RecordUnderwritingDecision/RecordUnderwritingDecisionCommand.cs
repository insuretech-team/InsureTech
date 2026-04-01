using System;
using System.Collections.Generic;
using InsuranceEngine.Underwriting.Application.DTOs;
using InsuranceEngine.Underwriting.Domain.Enums;
using InsuranceEngine.SharedKernel.CQRS;
using MediatR;

namespace InsuranceEngine.Underwriting.Application.Features.Commands.RecordUnderwritingDecision;

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
