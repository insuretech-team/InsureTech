using System.Globalization;
using System.Text.Json;
using Google.Protobuf.WellKnownTypes;
using Insuretech.Underwriting.Entity.V1;
using PoliSync.SharedKernel.CQRS;
using PoliSync.SharedKernel.Domain;
using ProtoMoney = Insuretech.Common.V1.Money;

namespace PoliSync.Underwriting.Domain;

public sealed class UnderwritingDecisionAggregate
{
    private readonly UnderwritingDecision _decision;
    private readonly List<DomainEvent> _domainEvents = [];

    private UnderwritingDecisionAggregate(UnderwritingDecision decision)
    {
        _decision = decision;
    }

    public UnderwritingDecision Decision => _decision;
    public IReadOnlyCollection<DomainEvent> DomainEvents => _domainEvents.AsReadOnly();

    public void ClearDomainEvents() => _domainEvents.Clear();

    public static Result<UnderwritingDecisionAggregate> CreateApproved(
        string quoteId,
        string underwriterId,
        string comments,
        string conditionsJson,
        bool premiumAdjusted,
        ProtoMoney? adjustedPremium,
        UnderwritingRiskAssessment riskAssessment,
        RiskLevel riskLevel)
    {
        var validation = ValidateCommon(quoteId, underwriterId, riskAssessment);
        if (validation.IsFailure)
        {
            return Result<UnderwritingDecisionAggregate>.Fail(validation.Error!);
        }

        if (premiumAdjusted && (adjustedPremium is null || adjustedPremium.Amount <= 0))
        {
            return Result.Fail<UnderwritingDecisionAggregate>(
                "INVALID_ADJUSTED_PREMIUM",
                "Adjusted premium must be set when premiumAdjusted is true");
        }

        var hasConditions = !string.IsNullOrWhiteSpace(conditionsJson) && conditionsJson != "{}";
        var now = DateTime.UtcNow;
        var normalizedAdjustedPremium = premiumAdjusted
            ? NormalizeMoney(adjustedPremium)
            : new ProtoMoney { Amount = 0, Currency = "BDT" };

        var decision = new UnderwritingDecision
        {
            Id = Guid.NewGuid().ToString("N"),
            QuoteId = quoteId,
            Decision = hasConditions ? DecisionType.Conditional : DecisionType.Approved,
            Method = DecisionMethod.Manual,
            RiskScore = riskAssessment.Score.ToString(CultureInfo.InvariantCulture),
            RiskLevel = riskLevel,
            RiskFactors = BuildRiskFactorsJson(riskAssessment),
            Reason = "Approved",
            Conditions = hasConditions ? conditionsJson : string.Empty,
            PremiumAdjusted = premiumAdjusted,
            AdjustedPremium = normalizedAdjustedPremium,
            AdjustmentReason = premiumAdjusted
                ? $"Manual underwriter adjustment; loading={riskAssessment.LoadingPercentage.ToString("0.##", CultureInfo.InvariantCulture)}%"
                : string.Empty,
            UnderwriterId = underwriterId,
            UnderwriterComments = comments ?? string.Empty,
            DecidedAt = Timestamp.FromDateTime(now),
            ValidUntil = Timestamp.FromDateTime(now.AddDays(30))
        };

        var aggregate = new UnderwritingDecisionAggregate(decision);
        aggregate._domainEvents.Add(new UnderwritingDecisionCreatedDomainEvent(
            decision.Id,
            decision.QuoteId,
            decision.Decision.ToString()));
        return Result.Ok(aggregate);
    }

    public static Result<UnderwritingDecisionAggregate> CreateRejected(
        string quoteId,
        string underwriterId,
        string reason,
        string comments,
        UnderwritingRiskAssessment riskAssessment,
        RiskLevel riskLevel)
    {
        var validation = ValidateCommon(quoteId, underwriterId, riskAssessment);
        if (validation.IsFailure)
        {
            return Result<UnderwritingDecisionAggregate>.Fail(validation.Error!);
        }

        if (string.IsNullOrWhiteSpace(reason))
        {
            return Result.Fail<UnderwritingDecisionAggregate>("INVALID_REASON", "Rejection reason is required");
        }

        var now = DateTime.UtcNow;
        var decision = new UnderwritingDecision
        {
            Id = Guid.NewGuid().ToString("N"),
            QuoteId = quoteId,
            Decision = DecisionType.Rejected,
            Method = DecisionMethod.Manual,
            RiskScore = riskAssessment.Score.ToString(CultureInfo.InvariantCulture),
            RiskLevel = riskLevel,
            RiskFactors = BuildRiskFactorsJson(riskAssessment),
            Reason = reason,
            Conditions = string.Empty,
            PremiumAdjusted = false,
            AdjustedPremium = new ProtoMoney { Amount = 0, Currency = "BDT" },
            AdjustmentReason = string.Empty,
            UnderwriterId = underwriterId,
            UnderwriterComments = comments ?? string.Empty,
            DecidedAt = Timestamp.FromDateTime(now),
            ValidUntil = Timestamp.FromDateTime(now.AddDays(30))
        };

        var aggregate = new UnderwritingDecisionAggregate(decision);
        aggregate._domainEvents.Add(new UnderwritingDecisionCreatedDomainEvent(
            decision.Id,
            decision.QuoteId,
            decision.Decision.ToString()));
        return Result.Ok(aggregate);
    }

    private static Result ValidateCommon(
        string quoteId,
        string underwriterId,
        UnderwritingRiskAssessment riskAssessment)
    {
        if (string.IsNullOrWhiteSpace(quoteId))
        {
            return Result.Fail("INVALID_QUOTE_ID", "QuoteId is required");
        }

        if (string.IsNullOrWhiteSpace(underwriterId))
        {
            return Result.Fail("INVALID_UNDERWRITER_ID", "UnderwriterId is required");
        }

        if (riskAssessment.Score is < 0 or > 100)
        {
            return Result.Fail("INVALID_RISK_SCORE", "Risk score must be between 0 and 100");
        }

        return Result.Ok();
    }

    private static string BuildRiskFactorsJson(UnderwritingRiskAssessment riskAssessment)
    {
        var payload = new
        {
            recommendation = riskAssessment.Recommendation.ToString(),
            loading_percentage = riskAssessment.LoadingPercentage,
            score = riskAssessment.Score
        };

        return JsonSerializer.Serialize(payload);
    }

    private static ProtoMoney NormalizeMoney(ProtoMoney? money)
        => new()
        {
            Amount = money?.Amount ?? 0,
            Currency = string.IsNullOrWhiteSpace(money?.Currency) ? "BDT" : money.Currency
        };
}

public sealed record UnderwritingDecisionCreatedDomainEvent(
    string DecisionId,
    string QuoteId,
    string Decision) : DomainEvent;
