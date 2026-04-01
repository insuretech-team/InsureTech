using System.Text.Json;
using Google.Protobuf.WellKnownTypes;
using Insuretech.Actuarial.Entity.V1;
using Insuretech.Actuarial.Services.V1;
using PoliSync.Infrastructure.Clients;

namespace PoliSync.Actuarial.Services;

public sealed class GoActuarialService :
    IActuarialService,
    IRatingFormulaService,
    IReserveCalculationService,
    ILossRatioService
{
    private readonly InsuranceServiceClient _client;

    public GoActuarialService(InsuranceServiceClient client) => _client = client;

    public async Task<ActuarialCalculation?> GetCalculationAsync(string calculationId, CancellationToken cancellationToken = default)
    {
        var response = await _client.ActuarialClient.GetCalculationAsync(
            new GetCalculationRequest { CalculationId = calculationId },
            _client.BuildCallOptions(cancellationToken));
        return response.Found ? response.Calculation : null;
    }

    public async Task<ActuarialCalculation?> GetCalculationByReferenceAsync(string reference, CancellationToken cancellationToken = default)
    {
        var response = await _client.ActuarialClient.GetCalculationAsync(
            new GetCalculationRequest { CalculationReference = reference },
            _client.BuildCallOptions(cancellationToken));
        return response.Found ? response.Calculation : null;
    }

    public async Task<IEnumerable<ActuarialCalculation>> GetCalculationsAsync(
        ActuarialCalculationType? type = null,
        string? entityType = null,
        string? entityId = null,
        DateTime? from = null,
        DateTime? to = null,
        CancellationToken cancellationToken = default)
    {
        var response = await _client.ActuarialClient.ListCalculationsAsync(
            new ListCalculationsRequest
            {
                CalculationType = type ?? ActuarialCalculationType.Unspecified,
                EntityType = entityType ?? string.Empty,
                EntityId = entityId ?? string.Empty,
                DateFrom = from is null ? null : Timestamp.FromDateTime(DateTime.SpecifyKind(from.Value, DateTimeKind.Utc)),
                DateTo = to is null ? null : Timestamp.FromDateTime(DateTime.SpecifyKind(to.Value, DateTimeKind.Utc)),
                PageSize = 200
            },
            _client.BuildCallOptions(cancellationToken));
        return response.Calculations;
    }

    public async Task<PremiumCalculationResult> CalculatePremiumAsync(PremiumCalculationInput input, CancellationToken cancellationToken = default)
    {
        var response = await _client.ActuarialClient.CalculatePremiumAsync(
            new CalculatePremiumRequest { Input = input, SaveCalculation = true },
            _client.BuildCallOptions(cancellationToken));
        return response.Result;
    }

    public async Task<double> CalculatePurePremiumAsync(double expectedClaims, double claimSeverity, double exposureUnits, double riskAdjustmentFactor)
    {
        var response = await _client.ActuarialClient.CalculatePurePremiumAsync(
            new CalculatePurePremiumRequest
            {
                ExpectedClaims = expectedClaims,
                ClaimSeverity = claimSeverity,
                ExposureUnits = exposureUnits,
                RiskAdjustmentFactor = riskAdjustmentFactor
            },
            _client.BuildCallOptions());
        return response.Result.GrossPremium;
    }

    public async Task<RatingFormula?> GetFormulaAsync(string formulaId, CancellationToken cancellationToken = default)
    {
        var response = await _client.ActuarialClient.GetRatingFormulaAsync(
            new GetRatingFormulaRequest { FormulaId = formulaId },
            _client.BuildCallOptions(cancellationToken));
        return response.Found ? response.Formula : null;
    }

    public async Task<RatingFormula?> GetFormulaByCodeAsync(string formulaCode, CancellationToken cancellationToken = default)
    {
        var response = await _client.ActuarialClient.GetRatingFormulaAsync(
            new GetRatingFormulaRequest { FormulaCode = formulaCode },
            _client.BuildCallOptions(cancellationToken));
        return response.Found ? response.Formula : null;
    }

    public async Task<IEnumerable<RatingFormula>> GetFormulasAsync(string? insuranceType = null, FormulaCategory? category = null, FormulaStatus? status = null, CancellationToken cancellationToken = default)
    {
        var response = await _client.ActuarialClient.ListRatingFormulasAsync(
            new ListRatingFormulasRequest
            {
                InsuranceType = insuranceType ?? string.Empty,
                Category = category ?? FormulaCategory.Unspecified,
                Status = status ?? FormulaStatus.Unspecified,
                PageSize = 200
            },
            _client.BuildCallOptions(cancellationToken));
        return response.Formulas;
    }

    public async Task<RatingFormula> CreateFormulaAsync(RatingFormula formula, CancellationToken cancellationToken = default)
    {
        var response = await _client.ActuarialClient.CreateRatingFormulaAsync(
            new CreateRatingFormulaRequest
            {
                FormulaCode = formula.FormulaCode,
                FormulaName = formula.FormulaName,
                Description = formula.Description,
                Category = formula.Category,
                InsuranceType = formula.InsuranceType,
                FormulaExpression = formula.FormulaExpression,
                Variables = { Deserialize<List<ActuarialVariable>>(formula.VariablesJson) ?? [] },
                SortOrder = formula.SortOrder,
                ValidFrom = formula.ValidFrom,
                ValidUntil = formula.ValidUntil,
                Metadata = { formula.Metadata }
            },
            _client.BuildCallOptions(cancellationToken));
        return response.Formula;
    }

    public async Task<RatingFormula?> UpdateFormulaAsync(string formulaId, RatingFormula formula, CancellationToken cancellationToken = default)
    {
        var response = await _client.ActuarialClient.UpdateRatingFormulaAsync(
            new UpdateRatingFormulaRequest
            {
                FormulaId = formulaId,
                FormulaName = formula.FormulaName,
                Description = formula.Description,
                Category = formula.Category,
                FormulaExpression = formula.FormulaExpression,
                Variables = { Deserialize<List<ActuarialVariable>>(formula.VariablesJson) ?? [] },
                SortOrder = formula.SortOrder,
                ValidUntil = formula.ValidUntil,
                Metadata = { formula.Metadata }
            },
            _client.BuildCallOptions(cancellationToken));
        return response.Success ? response.Formula : null;
    }

    public async Task<bool> DeleteFormulaAsync(string formulaId, CancellationToken cancellationToken = default)
    {
        var response = await _client.ActuarialClient.DeleteRatingFormulaAsync(
            new DeleteRatingFormulaRequest { FormulaId = formulaId },
            _client.BuildCallOptions(cancellationToken));
        return response.Success;
    }

    public async Task<RatingFormula?> ActivateFormulaAsync(string formulaId, CancellationToken cancellationToken = default)
    {
        var response = await _client.ActuarialClient.ActivateRatingFormulaAsync(
            new ActivateRatingFormulaRequest { FormulaId = formulaId },
            _client.BuildCallOptions(cancellationToken));
        return response.Success ? response.Formula : null;
    }

    public async Task<ReserveCalculation?> GetReserveAsync(string reserveId, CancellationToken cancellationToken = default)
    {
        var response = await _client.ActuarialClient.GetReserveCalculationAsync(
            new GetReserveCalculationRequest { ReserveId = reserveId },
            _client.BuildCallOptions(cancellationToken));
        return response.Found ? response.Reserve : null;
    }

    public async Task<ReserveCalculation?> GetReserveByClaimAsync(string claimId, CancellationToken cancellationToken = default)
    {
        var response = await _client.ActuarialClient.GetReserveCalculationAsync(
            new GetReserveCalculationRequest { ClaimId = claimId },
            _client.BuildCallOptions(cancellationToken));
        return response.Found ? response.Reserve : null;
    }

    public async Task<ReserveResult> CalculateReservesAsync(ReserveInput input, CancellationToken cancellationToken = default)
    {
        var response = await _client.ActuarialClient.CalculateReservesAsync(
            new CalculateReservesRequest
            {
                ClaimId = input.ClaimId,
                Input = input
            },
            _client.BuildCallOptions(cancellationToken));
        return response.Result;
    }

    public async Task<ReserveCalculation> SaveReserveAsync(ReserveCalculation reserve, CancellationToken cancellationToken = default)
    {
        var response = await _client.ActuarialClient.SaveReserveCalculationAsync(
            new SaveReserveCalculationRequest { Reserve = reserve },
            _client.BuildCallOptions(cancellationToken));
        return response.Reserve;
    }

    public async Task<LossRatioCalculation?> GetLossRatioAsync(string lossRatioId, CancellationToken cancellationToken = default)
    {
        var response = await _client.ActuarialClient.GetLossRatioCalculationAsync(
            new GetLossRatioCalculationRequest { LossRatioId = lossRatioId },
            _client.BuildCallOptions(cancellationToken));
        return response.Found ? response.LossRatio : null;
    }

    public async Task<IEnumerable<LossRatioCalculation>> GetLossRatiosAsync(string? productId = null, string? lineOfBusiness = null, DateTime? from = null, DateTime? to = null, CancellationToken cancellationToken = default)
    {
        var response = await _client.ActuarialClient.ListLossRatioCalculationsAsync(
            new ListLossRatioCalculationsRequest
            {
                ProductId = productId ?? string.Empty,
                LineOfBusiness = lineOfBusiness ?? string.Empty,
                PeriodStart = from is null ? null : Timestamp.FromDateTime(DateTime.SpecifyKind(from.Value, DateTimeKind.Utc)),
                PeriodEnd = to is null ? null : Timestamp.FromDateTime(DateTime.SpecifyKind(to.Value, DateTimeKind.Utc)),
                PageSize = 200
            },
            _client.BuildCallOptions(cancellationToken));
        return response.LossRatios;
    }

    public async Task<LossRatioResult> CalculateLossRatioAsync(LossRatioInput input, CancellationToken cancellationToken = default)
    {
        var response = await _client.ActuarialClient.CalculateLossRatioAsync(
            new CalculateLossRatioRequest
            {
                ProductId = input.ProductId,
                LineOfBusiness = input.LineOfBusiness,
                PeriodStart = input.PeriodStart,
                PeriodEnd = input.PeriodEnd,
                Input = input,
                SaveCalculation = true
            },
            _client.BuildCallOptions(cancellationToken));
        return response.Result;
    }

    public async Task<LossRatioCalculation> SaveLossRatioAsync(LossRatioCalculation lossRatio, CancellationToken cancellationToken = default)
    {
        var response = await _client.ActuarialClient.SaveLossRatioCalculationAsync(
            new SaveLossRatioCalculationRequest { LossRatio = lossRatio },
            _client.BuildCallOptions(cancellationToken));
        return response.LossRatio;
    }

    private static T? Deserialize<T>(string json)
    {
        return string.IsNullOrWhiteSpace(json) ? default : JsonSerializer.Deserialize<T>(json);
    }
}
