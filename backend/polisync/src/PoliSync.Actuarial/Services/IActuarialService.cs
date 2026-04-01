using Insuretech.Actuarial.Entity.V1;

namespace PoliSync.Actuarial.Services;

public interface IActuarialService
{
    Task<ActuarialCalculation?> GetCalculationAsync(string calculationId, CancellationToken cancellationToken = default);
    Task<ActuarialCalculation?> GetCalculationByReferenceAsync(string reference, CancellationToken cancellationToken = default);
    Task<IEnumerable<ActuarialCalculation>> GetCalculationsAsync(
        ActuarialCalculationType? type = null,
        string? entityType = null,
        string? entityId = null,
        DateTime? from = null,
        DateTime? to = null,
        CancellationToken cancellationToken = default);
    Task<PremiumCalculationResult> CalculatePremiumAsync(PremiumCalculationInput input, CancellationToken cancellationToken = default);
    Task<double> CalculatePurePremiumAsync(double expectedClaims, double claimSeverity, double exposureUnits, double riskAdjustmentFactor);
}

public interface IRatingFormulaService
{
    Task<RatingFormula?> GetFormulaAsync(string formulaId, CancellationToken cancellationToken = default);
    Task<RatingFormula?> GetFormulaByCodeAsync(string formulaCode, CancellationToken cancellationToken = default);
    Task<IEnumerable<RatingFormula>> GetFormulasAsync(
        string? insuranceType = null,
        FormulaCategory? category = null,
        FormulaStatus? status = null,
        CancellationToken cancellationToken = default);
    Task<RatingFormula> CreateFormulaAsync(RatingFormula formula, CancellationToken cancellationToken = default);
    Task<RatingFormula?> UpdateFormulaAsync(string formulaId, RatingFormula formula, CancellationToken cancellationToken = default);
    Task<bool> DeleteFormulaAsync(string formulaId, CancellationToken cancellationToken = default);
    Task<RatingFormula?> ActivateFormulaAsync(string formulaId, CancellationToken cancellationToken = default);
}

public interface IReserveCalculationService
{
    Task<ReserveCalculation?> GetReserveAsync(string reserveId, CancellationToken cancellationToken = default);
    Task<ReserveCalculation?> GetReserveByClaimAsync(string claimId, CancellationToken cancellationToken = default);
    Task<ReserveResult> CalculateReservesAsync(ReserveInput input, CancellationToken cancellationToken = default);
    Task<ReserveCalculation> SaveReserveAsync(ReserveCalculation reserve, CancellationToken cancellationToken = default);
}

public interface ILossRatioService
{
    Task<LossRatioCalculation?> GetLossRatioAsync(string lossRatioId, CancellationToken cancellationToken = default);
    Task<IEnumerable<LossRatioCalculation>> GetLossRatiosAsync(
        string? productId = null,
        string? lineOfBusiness = null,
        DateTime? from = null,
        DateTime? to = null,
        CancellationToken cancellationToken = default);
    Task<LossRatioResult> CalculateLossRatioAsync(LossRatioInput input, CancellationToken cancellationToken = default);
    Task<LossRatioCalculation> SaveLossRatioAsync(LossRatioCalculation lossRatio, CancellationToken cancellationToken = default);
}
