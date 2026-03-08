using System.Text.Json;
using Insuretech.Underwriting.Entity.V1;

namespace PoliSync.Underwriting.Domain;

public interface IUnderwritingRiskScorer
{
    UnderwritingRiskAssessment Evaluate(UnderwritingRiskProfile profile);
}

public sealed class UnderwritingRiskScorer : IUnderwritingRiskScorer
{
    public UnderwritingRiskAssessment Evaluate(UnderwritingRiskProfile profile)
    {
        var score = 50;
        score += AgePoints(profile.ApplicantAge);
        score += BmiPoints(profile.WeightKg, profile.HeightCm);
        score += profile.Smoker ? 15 : 0;
        score += Math.Min(30, CountItems(profile.PreExistingConditions) * 10);
        score += Math.Min(16, CountItems(profile.FamilyHistory) * 8);

        score = Math.Clamp(score, 0, 100);

        var recommendation = score switch
        {
            <= 40 => UnderwritingRecommendation.Approved,
            <= 60 => UnderwritingRecommendation.ApprovedWithLoading,
            <= 75 => UnderwritingRecommendation.Referred,
            _ => UnderwritingRecommendation.Declined
        };

        var riskLevel = score switch
        {
            <= 60 => RiskLevel.Low,
            <= 75 => RiskLevel.Medium,
            <= 90 => RiskLevel.High,
            _ => RiskLevel.VeryHigh
        };

        var loadingPercentage = recommendation == UnderwritingRecommendation.ApprovedWithLoading
            ? score switch
            {
                <= 50 => 10m,
                <= 55 => 15m,
                _ => 25m
            }
            : 0m;

        return new UnderwritingRiskAssessment(score, riskLevel, recommendation, loadingPercentage);
    }

    private static int AgePoints(int age) => age switch
    {
        <= 35 => 0,
        <= 50 => 10,
        <= 65 => 20,
        _ => 30
    };

    private static int BmiPoints(string weightKg, int heightCm)
    {
        var bmi = CalculateBmi(weightKg, heightCm);
        return bmi switch
        {
            >= 30m => 15,
            >= 25m => 5,
            _ => 0
        };
    }

    private static decimal CalculateBmi(string weightKg, int heightCm)
    {
        if (!decimal.TryParse(weightKg, out var weight) || weight <= 0 || heightCm <= 0)
        {
            return 0m;
        }

        var heightM = heightCm / 100m;
        if (heightM <= 0)
        {
            return 0m;
        }

        return weight / (heightM * heightM);
    }

    private static int CountItems(string payload)
    {
        if (string.IsNullOrWhiteSpace(payload))
        {
            return 0;
        }

        try
        {
            using var doc = JsonDocument.Parse(payload);
            var root = doc.RootElement;
            return root.ValueKind switch
            {
                JsonValueKind.Array => root.GetArrayLength(),
                JsonValueKind.Object => root.EnumerateObject().Count(),
                JsonValueKind.String => SplitCount(root.GetString()),
                _ => 0
            };
        }
        catch
        {
            return SplitCount(payload);
        }
    }

    private static int SplitCount(string? value)
    {
        if (string.IsNullOrWhiteSpace(value))
        {
            return 0;
        }

        return value
            .Split(',', StringSplitOptions.RemoveEmptyEntries | StringSplitOptions.TrimEntries)
            .Length;
    }
}

public sealed record UnderwritingRiskProfile(
    int ApplicantAge,
    int HeightCm,
    string WeightKg,
    bool Smoker,
    string PreExistingConditions,
    string FamilyHistory);

public sealed record UnderwritingRiskAssessment(
    int Score,
    RiskLevel RiskLevel,
    UnderwritingRecommendation Recommendation,
    decimal LoadingPercentage);

public enum UnderwritingRecommendation
{
    Approved = 1,
    ApprovedWithLoading = 2,
    Referred = 3,
    Declined = 4
}
