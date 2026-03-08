using FluentAssertions;
using Insuretech.Underwriting.Entity.V1;
using PoliSync.Underwriting.Domain;
using Xunit;

namespace PoliSync.Underwriting.Tests;

public class UnderwritingRiskScorerTests
{
    private readonly UnderwritingRiskScorer _scorer = new();

    [Theory]
    [MemberData(nameof(MatrixCases))]
    public void Evaluate_ReturnsExpectedRiskAssessment(
        int age,
        int heightCm,
        string weightKg,
        bool smoker,
        string preExistingConditions,
        string familyHistory,
        int expectedScore,
        RiskLevel expectedRiskLevel,
        UnderwritingRecommendation expectedRecommendation,
        decimal expectedLoadingPercentage)
    {
        var result = _scorer.Evaluate(new UnderwritingRiskProfile(
            ApplicantAge: age,
            HeightCm: heightCm,
            WeightKg: weightKg,
            Smoker: smoker,
            PreExistingConditions: preExistingConditions,
            FamilyHistory: familyHistory));

        result.Score.Should().Be(expectedScore);
        result.RiskLevel.Should().Be(expectedRiskLevel);
        result.Recommendation.Should().Be(expectedRecommendation);
        result.LoadingPercentage.Should().Be(expectedLoadingPercentage);
    }

    public static IEnumerable<object[]> MatrixCases()
    {
        return
        [
            new object[] { 30, 170, "65", false, JsonArray(0, "c"), JsonArray(0, "f"), 50, RiskLevel.Low, UnderwritingRecommendation.ApprovedWithLoading, 10m },
            new object[] { 40, 170, "65", false, JsonArray(0, "c"), JsonArray(0, "f"), 60, RiskLevel.Low, UnderwritingRecommendation.ApprovedWithLoading, 25m },
            new object[] { 55, 170, "65", false, JsonArray(0, "c"), JsonArray(0, "f"), 70, RiskLevel.Medium, UnderwritingRecommendation.Referred, 0m },
            new object[] { 70, 170, "65", false, JsonArray(0, "c"), JsonArray(0, "f"), 80, RiskLevel.High, UnderwritingRecommendation.Declined, 0m },
            new object[] { 30, 170, "80", false, JsonArray(0, "c"), JsonArray(0, "f"), 55, RiskLevel.Low, UnderwritingRecommendation.ApprovedWithLoading, 15m },
            new object[] { 30, 170, "95", false, JsonArray(0, "c"), JsonArray(0, "f"), 65, RiskLevel.Medium, UnderwritingRecommendation.Referred, 0m },
            new object[] { 30, 170, "65", true, JsonArray(0, "c"), JsonArray(0, "f"), 65, RiskLevel.Medium, UnderwritingRecommendation.Referred, 0m },
            new object[] { 30, 170, "65", false, JsonArray(1, "c"), JsonArray(0, "f"), 60, RiskLevel.Low, UnderwritingRecommendation.ApprovedWithLoading, 25m },
            new object[] { 30, 170, "65", false, JsonArray(3, "c"), JsonArray(0, "f"), 80, RiskLevel.High, UnderwritingRecommendation.Declined, 0m },
            new object[] { 30, 170, "65", false, JsonArray(5, "c"), JsonArray(0, "f"), 80, RiskLevel.High, UnderwritingRecommendation.Declined, 0m },
            new object[] { 30, 170, "65", false, JsonArray(0, "c"), JsonArray(1, "f"), 58, RiskLevel.Low, UnderwritingRecommendation.ApprovedWithLoading, 25m },
            new object[] { 30, 170, "65", false, JsonArray(0, "c"), JsonArray(2, "f"), 66, RiskLevel.Medium, UnderwritingRecommendation.Referred, 0m },
            new object[] { 30, 170, "65", false, JsonArray(0, "c"), JsonArray(5, "f"), 66, RiskLevel.Medium, UnderwritingRecommendation.Referred, 0m },
            new object[] { 40, 170, "80", true, JsonArray(1, "c"), JsonArray(1, "f"), 98, RiskLevel.VeryHigh, UnderwritingRecommendation.Declined, 0m },
            new object[] { 55, 170, "95", true, JsonArray(3, "c"), JsonArray(2, "f"), 100, RiskLevel.VeryHigh, UnderwritingRecommendation.Declined, 0m },
            new object[] { 35, 170, "65", false, JsonArray(0, "c"), JsonArray(0, "f"), 50, RiskLevel.Low, UnderwritingRecommendation.ApprovedWithLoading, 10m },
            new object[] { 36, 170, "65", false, JsonArray(0, "c"), JsonArray(0, "f"), 60, RiskLevel.Low, UnderwritingRecommendation.ApprovedWithLoading, 25m },
            new object[] { 50, 170, "65", false, JsonArray(0, "c"), JsonArray(0, "f"), 60, RiskLevel.Low, UnderwritingRecommendation.ApprovedWithLoading, 25m },
            new object[] { 51, 170, "65", false, JsonArray(0, "c"), JsonArray(0, "f"), 70, RiskLevel.Medium, UnderwritingRecommendation.Referred, 0m },
            new object[] { 65, 170, "65", false, JsonArray(0, "c"), JsonArray(0, "f"), 70, RiskLevel.Medium, UnderwritingRecommendation.Referred, 0m },
            new object[] { 66, 170, "65", false, JsonArray(0, "c"), JsonArray(0, "f"), 80, RiskLevel.High, UnderwritingRecommendation.Declined, 0m },
            new object[] { 30, 0, "invalid", false, JsonArray(0, "c"), JsonArray(0, "f"), 50, RiskLevel.Low, UnderwritingRecommendation.ApprovedWithLoading, 10m },
            new object[] { 30, 170, "65", false, "diabetes,hypertension,asthma,cardiac", JsonArray(0, "f"), 80, RiskLevel.High, UnderwritingRecommendation.Declined, 0m },
            new object[] { 30, 170, "65", false, JsonArray(0, "c"), "stroke,cancer,diabetes", 66, RiskLevel.Medium, UnderwritingRecommendation.Referred, 0m }
        ];
    }

    private static string JsonArray(int count, string prefix)
    {
        if (count <= 0)
        {
            return "[]";
        }

        var items = Enumerable.Range(1, count).Select(i => $"\"{prefix}{i}\"");
        return $"[{string.Join(",", items)}]";
    }
}
