using Insuretech.Actuarial.Entity.V1;
using System.Data;

namespace PoliSync.Actuarial.Services;

public class PremiumCalculator : IPremiumCalculator
{
    public PremiumCalculationResult CalculatePremium(PremiumCalculationInput input)
    {
        // Start with base rate calculation
        double baseRate = 0.001; // Default base rate (0.1% of sum insured)
        double basePremium = input.SumInsured * baseRate * input.CoveragePeriodMonths / 12.0;
        
        // Apply rating factors
        double totalFactor = 1.0;
        var factorBreakdown = new List<FactorBreakdown>();
        
        foreach (var factor in input.RatingFactors)
        {
            totalFactor *= factor.Value;
            factorBreakdown.Add(new FactorBreakdown
            {
                FactorName = factor.Key,
                FactorType = "RATING_FACTOR",
                FactorValue = factor.Value,
                Amount = basePremium * (factor.Value - 1.0),
                Description = $"Rating factor: {factor.Key}"
            });
        }
        
        double netPremium = basePremium * totalFactor;
        
        // Apply loadings
        double totalLoadings = 0;
        foreach (var loading in input.Loadings)
        {
            double loadingAmount = netPremium * 0.05; // Default 5% loading
            totalLoadings += loadingAmount;
            factorBreakdown.Add(new FactorBreakdown
            {
                FactorName = loading,
                FactorType = "LOADING",
                FactorValue = 1.05,
                Amount = loadingAmount,
                Description = $"Loading: {loading}"
            });
        }
        
        // Apply discounts
        double totalDiscounts = 0;
        foreach (var discount in input.Discounts)
        {
            double discountAmount = netPremium * 0.10; // Default 10% discount
            totalDiscounts += discountAmount;
            factorBreakdown.Add(new FactorBreakdown
            {
                FactorName = discount,
                FactorType = "DISCOUNT",
                FactorValue = 0.90,
                Amount = -discountAmount,
                Description = $"Discount: {discount}"
            });
        }
        
        double grossPremium = netPremium + totalLoadings - totalDiscounts;
        
        // Ensure minimum premium
        if (grossPremium < 100)
        {
            grossPremium = 100;
        }
        
        return new PremiumCalculationResult
        {
            BasePremium = basePremium,
            NetPremium = netPremium,
            GrossPremium = grossPremium,
            TotalLoadings = totalLoadings,
            TotalDiscounts = totalDiscounts,
            FactorBreakdown = { factorBreakdown },
            Currency = "USD"
        };
    }

    public double CalculatePurePremium(double expectedClaims, double claimSeverity, double exposureUnits, double riskAdjustmentFactor)
    {
        // Pure Premium = Expected Claim Frequency × Average Claim Severity × Risk Adjustment
        double purePremium = expectedClaims * claimSeverity * riskAdjustmentFactor;
        
        // Per exposure unit
        return purePremium / exposureUnits;
    }
}

public class FormulaEvaluator : IFormulaEvaluator
{
    public double Evaluate(string formulaExpression, Dictionary<string, double> variables)
    {
        // Replace variables with values
        string expression = formulaExpression;
        foreach (var variable in variables)
        {
            expression = expression.Replace($"{{{variable.Key}}}", variable.Value.ToString());
            expression = expression.Replace(variable.Key, variable.Value.ToString());
        }
        
        // Use DataTable.Compute for simple mathematical expressions
        // This is a simplified implementation - in production, use a proper expression evaluator
        try
        {
            using var dt = new DataTable();
            var result = dt.Compute(expression, "");
            return Convert.ToDouble(result);
        }
        catch
        {
            return 0;
        }
    }

    public bool ValidateExpression(string formulaExpression, out List<string> errors)
    {
        errors = new List<string>();
        
        try
        {
            using var dt = new DataTable();
            dt.Compute(formulaExpression.Replace("{", "").Replace("}", ""), "");
            return true;
        }
        catch (Exception ex)
        {
            errors.Add($"Invalid expression: {ex.Message}");
            return false;
        }
    }

    public List<string> ExtractVariables(string formulaExpression)
    {
        var variables = new List<string>();
        var parts = formulaExpression.Split(new[] { ' ', '+', '-', '*', '/', '(', ')', '{', '}' }, StringSplitOptions.RemoveEmptyEntries);
        
        foreach (var part in parts)
        {
            if (!double.TryParse(part, out _) && !IsOperator(part) && !IsFunction(part))
            {
                variables.Add(part.Trim());
            }
        }
        
        return variables.Distinct().ToList();
    }

    private static bool IsOperator(string value)
    {
        return value is "+" or "-" or "*" or "/" or "^" or "%";
    }

    private static bool IsFunction(string value)
    {
        var functions = new[] { "sin", "cos", "tan", "log", "ln", "exp", "sqrt", "abs", "min", "max", "pow" };
        return functions.Contains(value.ToLower());
    }
}
