using System.Text.Json;
using FluentAssertions;
using Insuretech.Common.V1;
using Insuretech.Underwriting.Entity.V1;
using PoliSync.Underwriting.Domain;
using Xunit;

namespace PoliSync.Underwriting.Tests;

public class UnderwritingDecisionAggregateTests
{
    [Fact]
    public void CreateApproved_WithConditions_BuildsConditionalDecision()
    {
        var result = UnderwritingDecisionAggregate.CreateApproved(
            quoteId: Guid.NewGuid().ToString("N"),
            underwriterId: "uw-1",
            comments: "Proceed with exclusions",
            conditionsJson: "{\"exclusion\":\"X\"}",
            premiumAdjusted: true,
            adjustedPremium: new Money { Amount = 45_000, Currency = "BDT" },
            riskAssessment: new UnderwritingRiskAssessment(55, RiskLevel.Low, UnderwritingRecommendation.ApprovedWithLoading, 15m),
            riskLevel: RiskLevel.Medium);

        result.IsSuccess.Should().BeTrue();
        result.Value.Should().NotBeNull();

        var decision = result.Value!.Decision;
        decision.Decision.Should().Be(DecisionType.Conditional);
        decision.RiskScore.Should().Be("55");
        decision.PremiumAdjusted.Should().BeTrue();
        decision.AdjustedPremium.Amount.Should().Be(45_000);

        using var doc = JsonDocument.Parse(decision.RiskFactors);
        doc.RootElement.GetProperty("loading_percentage").GetDecimal().Should().Be(15m);
    }

    [Fact]
    public void CreateApproved_WithoutAdjustedPremiumWhenFlagTrue_FailsValidation()
    {
        var result = UnderwritingDecisionAggregate.CreateApproved(
            quoteId: Guid.NewGuid().ToString("N"),
            underwriterId: "uw-1",
            comments: string.Empty,
            conditionsJson: string.Empty,
            premiumAdjusted: true,
            adjustedPremium: null,
            riskAssessment: new UnderwritingRiskAssessment(50, RiskLevel.Low, UnderwritingRecommendation.ApprovedWithLoading, 10m),
            riskLevel: RiskLevel.Low);

        result.IsFailure.Should().BeTrue();
        result.Error!.Code.Should().Be("INVALID_ADJUSTED_PREMIUM");
    }

    [Fact]
    public void CreateRejected_WithoutReason_FailsValidation()
    {
        var result = UnderwritingDecisionAggregate.CreateRejected(
            quoteId: Guid.NewGuid().ToString("N"),
            underwriterId: "uw-2",
            reason: string.Empty,
            comments: "reject",
            riskAssessment: new UnderwritingRiskAssessment(88, RiskLevel.High, UnderwritingRecommendation.Declined, 0m),
            riskLevel: RiskLevel.High);

        result.IsFailure.Should().BeTrue();
        result.Error!.Code.Should().Be("INVALID_REASON");
    }

    [Fact]
    public void CreateRejected_WithValidInput_BuildsRejectedDecision()
    {
        var result = UnderwritingDecisionAggregate.CreateRejected(
            quoteId: Guid.NewGuid().ToString("N"),
            underwriterId: "uw-2",
            reason: "Severe risk profile",
            comments: "declined",
            riskAssessment: new UnderwritingRiskAssessment(90, RiskLevel.High, UnderwritingRecommendation.Declined, 0m),
            riskLevel: RiskLevel.VeryHigh);

        result.IsSuccess.Should().BeTrue();
        var decision = result.Value!.Decision;
        decision.Decision.Should().Be(DecisionType.Rejected);
        decision.Reason.Should().Be("Severe risk profile");
        decision.PremiumAdjusted.Should().BeFalse();
        decision.RiskLevel.Should().Be(RiskLevel.VeryHigh);
    }
}
