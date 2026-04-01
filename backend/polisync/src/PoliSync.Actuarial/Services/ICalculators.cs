using Insuretech.Actuarial.Entity.V1;

namespace PoliSync.Actuarial.Services;

public interface IPremiumCalculator
{
    PremiumCalculationResult CalculatePremium(PremiumCalculationInput input);
    double CalculatePurePremium(double expectedClaims, double claimSeverity, double exposureUnits, double riskAdjustmentFactor);
}

public interface IFormulaEvaluator
{
    double Evaluate(string formulaExpression, Dictionary<string, double> variables);
    bool ValidateExpression(string formulaExpression, out List<string> errors);
    List<string> ExtractVariables(string formulaExpression);
}
