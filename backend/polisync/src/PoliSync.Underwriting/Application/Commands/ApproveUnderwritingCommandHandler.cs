using Insuretech.Underwriting.Entity.V1;
using MediatR;
using Microsoft.Extensions.Logging;
using PoliSync.SharedKernel.CQRS;
using PoliSync.SharedKernel.Messaging;
using PoliSync.Underwriting.Domain;
using PoliSync.Underwriting.Events;
using PoliSync.Underwriting.Infrastructure;
using ProtoMoney = Insuretech.Common.V1.Money;

namespace PoliSync.Underwriting.Application.Commands;

public sealed class ApproveUnderwritingCommandHandler
    : IRequestHandler<ApproveUnderwritingCommand, Result<string>>
{
    private readonly IUnderwritingDataGateway _dataGateway;
    private readonly IUnderwritingRiskScorer _riskScorer;
    private readonly IEventBus _eventBus;
    private readonly ILogger<ApproveUnderwritingCommandHandler> _logger;

    public ApproveUnderwritingCommandHandler(
        IUnderwritingDataGateway dataGateway,
        IUnderwritingRiskScorer riskScorer,
        IEventBus eventBus,
        ILogger<ApproveUnderwritingCommandHandler> logger)
    {
        _dataGateway = dataGateway;
        _riskScorer = riskScorer;
        _eventBus = eventBus;
        _logger = logger;
    }

    public async Task<Result<string>> Handle(ApproveUnderwritingCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var quote = await _dataGateway.GetQuoteAsync(request.QuoteId, cancellationToken);
            if (quote is null)
                return Result.Fail<string>("QUOTE_NOT_FOUND", $"Quote {request.QuoteId} not found");

            var declaration = await _dataGateway.GetHealthDeclarationByQuoteAsync(request.QuoteId, cancellationToken);
            var riskAssessment = _riskScorer.Evaluate(new UnderwritingRiskProfile(
                ApplicantAge: (int)quote.ApplicantAge,
                HeightCm: declaration?.HeightCm ?? 0,
                WeightKg: declaration?.WeightKg ?? string.Empty,
                Smoker: declaration?.Smoker ?? quote.Smoker,
                PreExistingConditions: declaration?.PreExistingConditions ?? string.Empty,
                FamilyHistory: declaration?.FamilyHistory ?? string.Empty));

            var riskLevel = ParseRiskLevel(request.RiskLevel) == RiskLevel.Unspecified
                ? riskAssessment.RiskLevel
                : ParseRiskLevel(request.RiskLevel);

            var adjustedPremium = request.PremiumAdjusted && request.AdjustedPremiumPaisa > 0
                ? new ProtoMoney { Amount = request.AdjustedPremiumPaisa, Currency = "BDT" }
                : quote.TotalPremium;

            var decisionResult = UnderwritingDecisionAggregate.CreateApproved(
                quoteId: request.QuoteId,
                underwriterId: request.UnderwriterId,
                comments: request.Comments,
                conditionsJson: request.ConditionsJson ?? string.Empty,
                premiumAdjusted: request.PremiumAdjusted,
                adjustedPremium: adjustedPremium,
                riskAssessment: riskAssessment,
                riskLevel: riskLevel);

            if (decisionResult.IsFailure)
                return Result.Fail<string>(decisionResult.Error!.Code, decisionResult.Error.Message);

            var persisted = await _dataGateway.UpsertUnderwritingDecisionAsync(
                decisionResult.Value!.Decision, cancellationToken);

            quote.Status = Insuretech.Underwriting.Entity.V1.QuoteStatus.Approved;
            if (request.PremiumAdjusted && adjustedPremium is not null)
                quote.TotalPremium = adjustedPremium;
            await _dataGateway.UpdateQuoteAsync(quote, cancellationToken);

            foreach (var evt in decisionResult.Value.DomainEvents)
                await _eventBus.PublishAsync(evt, cancellationToken);

            await _eventBus.PublishAsync(new UnderwritingDecisionMadeEvent
            {
                QuoteId = request.QuoteId,
                DecisionId = persisted.Id,
                Decision = "APPROVED",
                RiskLevel = riskLevel.ToString(),
                PremiumAdjusted = request.PremiumAdjusted,
                QuotedAmount = adjustedPremium?.Amount ?? 0,
                Currency = adjustedPremium?.Currency ?? "BDT",
                Reason = string.Empty
            }, "insuretech.underwriting.decision_made.v1", cancellationToken);

            _logger.LogInformation("Underwriting approved for quote {QuoteId}. Decision: {DecisionId}", request.QuoteId, persisted.Id);
            return Result.Ok(persisted.Id);
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to approve underwriting for quote {QuoteId}", request.QuoteId);
            return Result.Fail<string>("APPROVE_UNDERWRITING_FAILED", ex.Message);
        }
    }

    private static RiskLevel ParseRiskLevel(string value)
    {
        if (string.IsNullOrWhiteSpace(value)) return RiskLevel.Unspecified;
        return Enum.TryParse<RiskLevel>(value, true, out var parsed) ? parsed : RiskLevel.Unspecified;
    }
}
