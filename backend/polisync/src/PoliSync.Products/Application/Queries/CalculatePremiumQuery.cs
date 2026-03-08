using Insuretech.Common.V1;
using Insuretech.Products.Services.V1;
using PoliSync.Products.Persistence;
using PoliSync.SharedKernel.CQRS;
using System.Text.Json;

namespace PoliSync.Products.Application.Queries;

/// <summary>
/// Stateless premium calculation — no DB write. Applies pricing rules from PricingConfig.
/// Rules stored as JSONB; evaluated left-to-right, first match wins.
/// </summary>
public sealed record CalculatePremiumQuery(CalculatePremiumRequest Request) : IQuery<CalculatePremiumResponse>;

public sealed class CalculatePremiumHandler : IQueryHandler<CalculatePremiumQuery, CalculatePremiumResponse>
{
    private readonly ProductRepository _repo;

    public CalculatePremiumHandler(ProductRepository repo) => _repo = repo;

    public async Task<Result<CalculatePremiumResponse>> Handle(CalculatePremiumQuery query, CancellationToken ct)
    {
        var req = query.Request;
        if (!Guid.TryParse(req.ProductId, out var productId))
            return Result<CalculatePremiumResponse>.Fail("INVALID", "Invalid product_id.");

        var record = await _repo.GetByIdAsync(productId, ct);
        if (record is null)
            return Result<CalculatePremiumResponse>.NotFound($"Product '{productId}' not found.");

        if (record.Status != "PRODUCT_STATUS_ACTIVE")
            return Result<CalculatePremiumResponse>.Fail("INVALID", "Premium can only be calculated for active products.");

        var sumInsured    = req.SumInsured?.Amount ?? 0;
        var tenureMonths  = req.TenureMonths;

        // Validate sum insured range
        if (sumInsured < record.MinSumInsured || sumInsured > record.MaxSumInsured)
            return Result<CalculatePremiumResponse>.Fail("INVALID",
                $"Sum insured must be between {record.MinSumInsured} and {record.MaxSumInsured} paisa.");

        // Validate tenure range
        if (tenureMonths < record.MinTenureMonths || tenureMonths > record.MaxTenureMonths)
            return Result<CalculatePremiumResponse>.Fail("INVALID",
                $"Tenure must be between {record.MinTenureMonths} and {record.MaxTenureMonths} months.");

        // Base premium (annualised, then scaled by tenure)
        var basePremiumAnnual = record.BasePremium;
        var basePremium       = (long)Math.Round(basePremiumAnnual * tenureMonths / 12.0);

        var breakdown = new List<PremiumBreakdown>
        {
            new() { Item = "Base Premium", Amount = new Money { Amount = basePremium, Currency = "BDT" },
                    Description = $"{tenureMonths} months" }
        };

        // Apply pricing rules
        long adjustedPremium = basePremium;
        if (record.PricingConfig?.Rules is { } rulesJson && rulesJson != "[]")
        {
            var rules = JsonSerializer.Deserialize<List<PricingRuleJson>>(rulesJson) ?? [];
            foreach (var rule in rules)
            {
                if (!EvaluateConditions(rule.Conditions, req.ApplicantData))
                    continue;

                long adjustment = rule.Action.Type switch
                {
                    "ACTION_TYPE_INCREASE_PERCENTAGE" =>
                        (long)Math.Round(adjustedPremium * rule.Action.Value / 100.0),
                    "ACTION_TYPE_DECREASE_PERCENTAGE" =>
                        -(long)Math.Round(adjustedPremium * rule.Action.Value / 100.0),
                    "ACTION_TYPE_FIXED_AMOUNT" =>
                        (long)rule.Action.Value,
                    _ => 0
                };

                if (adjustment != 0)
                {
                    adjustedPremium += adjustment;
                    breakdown.Add(new PremiumBreakdown
                    {
                        Item        = rule.RuleName,
                        Amount      = new Money { Amount = Math.Abs(adjustment), Currency = "BDT" },
                        Description = adjustment > 0 ? "Loading" : "Discount",
                    });
                }
                break; // first matching rule wins
            }
        }

        // Rider premiums
        long riderTotal = 0;
        foreach (var riderId in req.RiderIds)
        {
            var rider = record.Riders.FirstOrDefault(r => r.RiderId.ToString() == riderId);
            if (rider is null) continue;
            var riderPremium = (long)Math.Round(rider.PremiumAmount * tenureMonths / 12.0);
            riderTotal += riderPremium;
            breakdown.Add(new PremiumBreakdown
            {
                Item        = rider.RiderName,
                Amount      = new Money { Amount = riderPremium, Currency = "BDT" },
                Description = "Rider premium",
            });
        }

        var totalPremium = adjustedPremium + riderTotal;

        var response = new CalculatePremiumResponse
        {
            BasePremium   = new Money { Amount = adjustedPremium, Currency = "BDT" },
            RiderPremium  = new Money { Amount = riderTotal,      Currency = "BDT" },
            TotalPremium  = new Money { Amount = totalPremium,    Currency = "BDT" },
        };
        response.Breakdown.AddRange(breakdown);
        return Result<CalculatePremiumResponse>.Ok(response);
    }

    private static bool EvaluateConditions(
        List<PricingConditionJson> conditions,
        IDictionary<string, string> applicantData)
    {
        foreach (var cond in conditions)
        {
            if (!applicantData.TryGetValue(cond.Field, out var actual)) return false;
            var match = cond.Operator switch
            {
                ">="  => double.TryParse(actual, out var a) && double.TryParse(cond.Value, out var b) && a >= b,
                "<="  => double.TryParse(actual, out var a) && double.TryParse(cond.Value, out var b) && a <= b,
                ">"   => double.TryParse(actual, out var a) && double.TryParse(cond.Value, out var b) && a > b,
                "<"   => double.TryParse(actual, out var a) && double.TryParse(cond.Value, out var b) && a < b,
                "eq"  => string.Equals(actual, cond.Value, StringComparison.OrdinalIgnoreCase),
                "in"  => cond.Value.Split(',').Contains(actual, StringComparer.OrdinalIgnoreCase),
                _     => false,
            };
            if (!match) return false;
        }
        return true;
    }

    // ── Private JSON DTOs for JSONB rule deserialization ──────────────────
    private sealed record PricingRuleJson(
        string RuleId, string RuleName, string Type,
        List<PricingConditionJson> Conditions, PricingActionJson Action);

    private sealed record PricingConditionJson(string Field, string Operator, string Value);
    private sealed record PricingActionJson(string Type, double Value);
}
