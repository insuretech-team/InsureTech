using System.Collections.Generic;
using System.Linq;
using System.Text.Json;
using System.Threading;
using System.Threading.Tasks;
using InsuranceEngine.Policy.Application.DTOs;
using InsuranceEngine.Policy.Application.Interfaces;
using InsuranceEngine.SharedKernel.CQRS;
using MediatR;

namespace InsuranceEngine.Policy.Application.Features.Queries.GetUnderwritingHistory;

public class GetUnderwritingHistoryQueryHandler : IRequestHandler<GetUnderwritingHistoryQuery, Result<List<UnderwritingDecisionDto>>>
{
    private readonly IUnderwritingRepository _repository;

    public GetUnderwritingHistoryQueryHandler(IUnderwritingRepository repository)
    {
        _repository = repository;
    }

    public async Task<Result<List<UnderwritingDecisionDto>>> Handle(GetUnderwritingHistoryQuery request, CancellationToken cancellationToken)
    {
        var (items, _) = await _repository.GetDecisionHistoryAsync(request.QuoteId);

        var dtos = items.Select(d => new UnderwritingDecisionDto(
            d.Id,
            d.QuoteId,
            d.Decision,
            d.Method,
            d.RiskScore,
            d.RiskLevel,
            d.RiskFactorsJson != null ? JsonSerializer.Deserialize<List<string>>(d.RiskFactorsJson) : null,
            d.Reason,
            d.ConditionsJson != null ? JsonSerializer.Deserialize<List<string>>(d.ConditionsJson) : null,
            d.IsPremiumAdjusted,
            new MoneyDto(d.AdjustedPremiumAmount, d.AdjustedPremiumCurrency),
            d.AdjustmentReason,
            d.UnderwriterId,
            d.UnderwriterComments,
            d.DecidedAt,
            d.ValidUntil
        )).ToList();

        return Result.Ok(dtos);
    }
}
