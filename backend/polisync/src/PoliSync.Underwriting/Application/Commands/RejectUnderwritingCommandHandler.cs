using Insuretech.Underwriting.Entity.V1;
using MediatR;
using Microsoft.Extensions.Logging;
using PoliSync.SharedKernel.CQRS;
using PoliSync.SharedKernel.Messaging;
using PoliSync.Underwriting.Domain;
using PoliSync.Underwriting.Events;
using PoliSync.Underwriting.Infrastructure;

namespace PoliSync.Underwriting.Application.Commands;

public sealed class RejectUnderwritingCommandHandler
    : IRequestHandler<RejectUnderwritingCommand, Result<string>>
{
    private readonly IUnderwritingDataGateway _dataGateway;
    private readonly IUnderwritingRiskScorer _riskScorer;
    private readonly IEventBus _eventBus;
    private readonly ILogger<RejectUnderwritingCommandHandler> _logger;

    public RejectUnderwritingCommandHandler(
        IUnderwritingDataGateway dataGateway,
        IUnderwritingRiskScorer riskScorer,
        IEventBus eventBus,
        ILogger<RejectUnderwritingCommandHandler> logger)
    {
        _dataGateway = dataGateway;
        _riskScorer = riskScorer;
        _eventBus = eventBus;
        _logger = logger;
    }

    public async Task<Result<string>> Handle(RejectUnderwritingCommand request, CancellationToken cancellationToken)
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

            var decisionResult = UnderwritingDecisionAggregate.CreateRejected(
                quoteId: request.QuoteId,
                underwriterId: request.UnderwriterId,
                reason: request.Reason,
                comments: request.Comments,
                riskAssessment: riskAssessment,
                riskLevel: riskLevel);

            if (decisionResult.IsFailure)
                return Result.Fail<string>(decisionResult.Error!.Code, decisionResult.Error.Message);

            var persisted = await _dataGateway.UpsertUnderwritingDecisionAsync(
                decisionResult.Value!.Decision, cancellationToken);

            quote.Status = QuoteStatus.Rejected;
            await _dataGateway.UpdateQuoteAsync(quote, cancellationToken);

            foreach (var evt in decisionResult.Value.DomainEvents)
                await _eventBus.PublishAsync(evt, cancellationToken);

            await _eventBus.PublishAsync(new UnderwritingDecisionMadeEvent
            {
                QuoteId = request.QuoteId,
                DecisionId = persisted.Id,
                Decision = "REJECTED",
                RiskLevel = riskLevel.ToString(),
                PremiumAdjusted = false,
                QuotedAmount = 0,
                Currency = "BDT",
                Reason = request.Reason
            }, "insuretech.underwriting.decision_made.v1", cancellationToken);

            _logger.LogInformation("Underwriting rejected for quote {QuoteId}. Decision: {DecisionId}", request.QuoteId, persisted.Id);
            return Result.Ok(persisted.Id);
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to reject underwriting for quote {QuoteId}", request.QuoteId);
            return Result.Fail<string>("REJECT_UNDERWRITING_FAILED", ex.Message);
        }
    }

    private static RiskLevel ParseRiskLevel(string value)
    {
        if (string.IsNullOrWhiteSpace(value)) return RiskLevel.Unspecified;
        return Enum.TryParse<RiskLevel>(value, true, out var parsed) ? parsed : RiskLevel.Unspecified;
    }
}
