namespace PoliSync.Products.Domain;

/// <summary>
/// Static pricing engine for evaluating pricing rules against contextual factors.
/// </summary>
public static class PricingEngine
{
    /// <summary>
    /// Evaluate pricing rules against a context dictionary to compute final premium.
    /// </summary>
    /// <param name="rules">List of pricing rules ordered by priority</param>
    /// <param name="basePremiumPaisa">Base premium in paisa</param>
    /// <param name="context">Dictionary of contextual factors (e.g., age, sum_insured, product_category, tenure_months)</param>
    /// <returns>PricingResult with breakdown and final premium</returns>
    public static PricingResult Evaluate(List<PricingRule> rules, long basePremiumPaisa, Dictionary<string, string> context)
    {
        var breakdown = new List<PremiumBreakdownItem>();
        var appliedRules = new List<string>();
        decimal currentPremium = basePremiumPaisa;

        // Add base premium to breakdown
        breakdown.Add(new PremiumBreakdownItem(
            Description: "Base Premium",
            DeltaPaisa: basePremiumPaisa));

        // Sort rules by priority
        var sortedRules = rules.OrderBy(r => r.Priority).ToList();
        var ruleMatched = false;

        foreach (var rule in sortedRules)
        {
            // Evaluate conditions (AND logic)
            if (!EvaluateConditions(rule.Conditions, context))
                continue;

            // Rule matched
            ruleMatched = true;
            appliedRules.Add(rule.RuleName);

            // Apply action
            var (newPremium, deltaPaisa) = ApplyAction(currentPremium, rule.Action, basePremiumPaisa);
            currentPremium = newPremium;

            breakdown.Add(new PremiumBreakdownItem(
                Description: rule.RuleName,
                DeltaPaisa: deltaPaisa));

            // If ApplyAll is false, stop after first match
            if (!rule.ApplyAll)
                break;
        }

        // Ensure premium is never negative
        if (currentPremium < 0)
            currentPremium = 0;

        return new PricingResult(
            BasePremiumPaisa: basePremiumPaisa,
            FinalPremiumPaisa: (long)currentPremium,
            AppliedRules: appliedRules,
            Breakdown: breakdown);
    }

    /// <summary>
    /// Evaluate all conditions in AND logic.
    /// </summary>
    private static bool EvaluateConditions(List<RuleCondition> conditions, Dictionary<string, string> context)
    {
        if (!conditions.Any())
            return true;

        foreach (var condition in conditions)
        {
            if (!context.TryGetValue(condition.Field, out var contextValue))
                return false;

            if (!EvaluateCondition(contextValue, condition.Operator, condition.Value))
                return false;
        }

        return true;
    }

    /// <summary>
    /// Evaluate a single condition.
    /// </summary>
    private static bool EvaluateCondition(string contextValue, string op, string ruleValue)
    {
        return op.ToUpperInvariant() switch
        {
            "EQ" => contextValue == ruleValue,
            "NEQ" => contextValue != ruleValue,
            "GT" => TryParseDecimal(contextValue, out var ctx) && TryParseDecimal(ruleValue, out var rule) && ctx > rule,
            "GTE" => TryParseDecimal(contextValue, out var ctx1) && TryParseDecimal(ruleValue, out var rule1) && ctx1 >= rule1,
            "LT" => TryParseDecimal(contextValue, out var ctx2) && TryParseDecimal(ruleValue, out var rule2) && ctx2 < rule2,
            "LTE" => TryParseDecimal(contextValue, out var ctx3) && TryParseDecimal(ruleValue, out var rule3) && ctx3 <= rule3,
            "IN" => ruleValue.Split(',').Contains(contextValue),
            "BETWEEN" => EvaluateBetween(contextValue, ruleValue),
            _ => false
        };
    }

    /// <summary>
    /// Evaluate BETWEEN operator. Format: "min,max"
    /// </summary>
    private static bool EvaluateBetween(string contextValue, string rangeValue)
    {
        var parts = rangeValue.Split(',');
        if (parts.Length != 2)
            return false;

        if (!TryParseDecimal(contextValue, out var context) ||
            !TryParseDecimal(parts[0], out var min) ||
            !TryParseDecimal(parts[1], out var max))
            return false;

        return context >= min && context <= max;
    }

    /// <summary>
    /// Apply a pricing action and return the new premium and delta.
    /// </summary>
    private static (decimal NewPremium, long DeltaPaisa) ApplyAction(decimal currentPremium, RuleAction action, long basePremiumPaisa)
    {
        return action.ActionType.ToUpperInvariant() switch
        {
            "MULTIPLY" => 
                (currentPremium * action.Value, (long)(currentPremium * (action.Value - 1))),
            
            "ADD" => 
                (currentPremium + action.Value, (long)action.Value),
            
            "SET" => 
                (action.Value, (long)(action.Value - currentPremium)),
            
            "DISCOUNT" => 
                // Discount is a percentage reduction
                (currentPremium * (1 - action.Value / 100), -(long)(currentPremium * action.Value / 100)),
            
            _ => (currentPremium, 0)
        };
    }

    /// <summary>
    /// Helper to safely parse decimal values.
    /// </summary>
    private static bool TryParseDecimal(string value, out decimal result)
    {
        return decimal.TryParse(value, out result);
    }
}

/// <summary>
/// Result of pricing evaluation.
/// </summary>
public record PricingResult(
    long BasePremiumPaisa,
    long FinalPremiumPaisa,
    List<string> AppliedRules,
    List<PremiumBreakdownItem> Breakdown);

/// <summary>
/// Individual line item in premium breakdown.
/// </summary>
public record PremiumBreakdownItem(
    string Description,
    long DeltaPaisa);
