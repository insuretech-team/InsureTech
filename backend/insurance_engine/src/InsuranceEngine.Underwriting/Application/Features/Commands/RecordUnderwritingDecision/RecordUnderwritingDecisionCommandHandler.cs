using System;
using System.Collections.Generic;
using System.Linq;
using System.Text.Json;
using System.Threading;
using System.Threading.Tasks;
using InsuranceEngine.Underwriting.Application.DTOs;
using InsuranceEngine.Underwriting.Application.Interfaces;
using InsuranceEngine.Underwriting.Domain.Entities;
using InsuranceEngine.Underwriting.Domain.Enums;
using InsuranceEngine.SharedKernel.CQRS;
using InsuranceEngine.SharedKernel.DTOs;
using MediatR;

namespace InsuranceEngine.Underwriting.Application.Features.Commands.RecordUnderwritingDecision;

public class RecordUnderwritingDecisionCommandHandler : IRequestHandler<RecordUnderwritingDecisionCommand, Result<UnderwritingDecisionDto>>
{
    private readonly IUnderwritingRepository _repository;

    public RecordUnderwritingDecisionCommandHandler(IUnderwritingRepository repository)
    {
        _repository = repository;
    }

    public async Task<Result<UnderwritingDecisionDto>> Handle(RecordUnderwritingDecisionCommand request, CancellationToken cancellationToken)
    {
        var quote = await _repository.GetQuoteByIdAsync(request.QuoteId);
        if (quote == null)
            return Result.Fail<UnderwritingDecisionDto>(Error.NotFound("Quote", request.QuoteId.ToString()));

        var decision = new UnderwritingDecision
        {
            Id = Guid.NewGuid(),
            QuoteId = request.QuoteId,
            Decision = request.Decision,
            Method = request.Method,
            RiskScore = request.RiskScore,
            RiskLevel = request.RiskLevel,
            RiskFactorsJson = request.RiskFactors != null ? JsonSerializer.Serialize(request.RiskFactors) : null,
            Reason = request.Reason,
            ConditionsJson = request.Conditions != null ? JsonSerializer.Serialize(request.Conditions) : null,
            IsPremiumAdjusted = request.IsPremiumAdjusted,
            AdjustedPremiumAmount = request.AdjustedPremiumAmount ?? 0,
            AdjustedPremiumCurrency = quote.Currency,
            AdjustmentReason = request.IsPremiumAdjusted ? "Underwriting adjustment" : null,
            UnderwriterId = request.UnderwriterId,
            UnderwriterComments = request.UnderwriterComments,
            DecidedAt = DateTime.UtcNow,
            ValidUntil = DateTime.UtcNow.AddDays(30),
            CreatedAt = DateTime.UtcNow,
            UpdatedAt = DateTime.UtcNow
        };

        if (request.Decision == DecisionType.Approved || request.Decision == DecisionType.Conditional)
        {
            quote.Status = QuoteStatus.Approved;
        }
        else if (request.Decision == DecisionType.Rejected)
        {
            quote.Status = QuoteStatus.Rejected;
        }
        else if (request.Decision == DecisionType.Referred)
        {
            quote.Status = QuoteStatus.PendingUnderwriting;
        }

        quote.UpdatedAt = DateTime.UtcNow;

        await _repository.AddDecisionAsync(decision);
        await _repository.UpdateQuoteAsync(quote);

        return Result.Ok(new UnderwritingDecisionDto(
            decision.Id,
            decision.QuoteId,
            decision.Decision,
            decision.Method,
            decision.RiskScore,
            decision.RiskLevel,
            request.RiskFactors,
            decision.Reason,
            request.Conditions,
            decision.IsPremiumAdjusted,
            new MoneyDto(decision.AdjustedPremiumAmount, decision.AdjustedPremiumCurrency),
            decision.AdjustmentReason,
            decision.UnderwriterId,
            decision.UnderwriterComments,
            decision.DecidedAt,
            decision.ValidUntil
        ));
    }
}
