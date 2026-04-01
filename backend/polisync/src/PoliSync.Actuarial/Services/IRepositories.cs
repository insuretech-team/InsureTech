using Insuretech.Actuarial.Entity.V1;

namespace PoliSync.Actuarial.Services;

public interface IRatingFormulaRepository
{
    Task<RatingFormula?> GetByIdAsync(string formulaId, CancellationToken cancellationToken = default);
    Task<RatingFormula?> GetByCodeAsync(string formulaCode, CancellationToken cancellationToken = default);
    Task<IEnumerable<RatingFormula>> GetAllAsync(CancellationToken cancellationToken = default);
    Task<IEnumerable<RatingFormula>> GetByFiltersAsync(string? insuranceType, FormulaCategory? category, FormulaStatus? status, CancellationToken cancellationToken = default);
    Task<RatingFormula> CreateAsync(RatingFormula formula, CancellationToken cancellationToken = default);
    Task<RatingFormula?> UpdateAsync(RatingFormula formula, CancellationToken cancellationToken = default);
    Task<bool> DeleteAsync(string formulaId, CancellationToken cancellationToken = default);
}

public interface IActuarialCalculationRepository
{
    Task<ActuarialCalculation?> GetByIdAsync(string calculationId, CancellationToken cancellationToken = default);
    Task<ActuarialCalculation?> GetByReferenceAsync(string reference, CancellationToken cancellationToken = default);
    Task<IEnumerable<ActuarialCalculation>> GetAllAsync(CancellationToken cancellationToken = default);
    Task<IEnumerable<ActuarialCalculation>> GetByFiltersAsync(
        ActuarialCalculationType? type,
        string? entityType,
        string? entityId,
        DateTime? from,
        DateTime? to,
        CancellationToken cancellationToken = default);
    Task<ActuarialCalculation> CreateAsync(ActuarialCalculation calculation, CancellationToken cancellationToken = default);
}

public interface IReserveCalculationRepository
{
    Task<ReserveCalculation?> GetByIdAsync(string reserveId, CancellationToken cancellationToken = default);
    Task<ReserveCalculation?> GetByClaimAsync(string claimId, CancellationToken cancellationToken = default);
    Task<IEnumerable<ReserveCalculation>> GetAllAsync(CancellationToken cancellationToken = default);
    Task<ReserveCalculation> CreateAsync(ReserveCalculation reserve, CancellationToken cancellationToken = default);
    Task<ReserveCalculation?> UpdateAsync(ReserveCalculation reserve, CancellationToken cancellationToken = default);
}

public interface ILossRatioCalculationRepository
{
    Task<LossRatioCalculation?> GetByIdAsync(string lossRatioId, CancellationToken cancellationToken = default);
    Task<IEnumerable<LossRatioCalculation>> GetAllAsync(CancellationToken cancellationToken = default);
    Task<IEnumerable<LossRatioCalculation>> GetByFiltersAsync(
        string? productId,
        string? lineOfBusiness,
        DateTime? from,
        DateTime? to,
        CancellationToken cancellationToken = default);
    Task<LossRatioCalculation> CreateAsync(LossRatioCalculation lossRatio, CancellationToken cancellationToken = default);
}
