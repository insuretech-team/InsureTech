using Google.Protobuf.WellKnownTypes;
using Insuretech.Actuarial.Entity.V1;
using Microsoft.Extensions.Logging;
using System.Text.Json;

namespace PoliSync.Actuarial.Services;

public class ActuarialCalculationService : IActuarialService
{
    private readonly IActuarialCalculationRepository _calculationRepository;
    private readonly IPremiumCalculator _premiumCalculator;
    private readonly ILogger<ActuarialCalculationService> _logger;

    public ActuarialCalculationService(
        IActuarialCalculationRepository calculationRepository,
        IPremiumCalculator premiumCalculator,
        ILogger<ActuarialCalculationService> logger)
    {
        _calculationRepository = calculationRepository;
        _premiumCalculator = premiumCalculator;
        _logger = logger;
    }

    public async Task<ActuarialCalculation?> GetCalculationAsync(string calculationId, CancellationToken cancellationToken = default)
    {
        return await _calculationRepository.GetByIdAsync(calculationId, cancellationToken);
    }

    public async Task<ActuarialCalculation?> GetCalculationByReferenceAsync(string reference, CancellationToken cancellationToken = default)
    {
        return await _calculationRepository.GetByReferenceAsync(reference, cancellationToken);
    }

    public async Task<IEnumerable<ActuarialCalculation>> GetCalculationsAsync(
        ActuarialCalculationType? type = null,
        string? entityType = null,
        string? entityId = null,
        DateTime? from = null,
        DateTime? to = null,
        CancellationToken cancellationToken = default)
    {
        return await _calculationRepository.GetByFiltersAsync(type, entityType, entityId, from, to, cancellationToken);
    }

    public async Task<PremiumCalculationResult> CalculatePremiumAsync(PremiumCalculationInput input, CancellationToken cancellationToken = default)
    {
        _logger.LogInformation("Calculating premium for product {ProductId}", input.ProductId);
        
        var result = _premiumCalculator.CalculatePremium(input);
        
        // Save calculation
        var calculation = new ActuarialCalculation
        {
            CalculationId = Guid.NewGuid().ToString(),
            CalculationReference = $"ACT-{DateTime.UtcNow:yyyy}-{Guid.NewGuid().ToString()[..6].ToUpper()}",
            CalculationType = ActuarialCalculationType.Premium,
            EntityType = "PRODUCT",
            EntityId = input.ProductId,
            ParametersJson = JsonSerializer.Serialize(input),
            ResultsJson = JsonSerializer.Serialize(result),
            Status = CalculationStatus.Completed,
            CalculatedPremium = result.GrossPremium,
            CalculatedAt = Timestamp.FromDateTime(DateTime.UtcNow),
            EffectiveDate = Timestamp.FromDateTime(DateTime.UtcNow)
        };
        
        await _calculationRepository.CreateAsync(calculation, cancellationToken);
        
        return result;
    }

    public async Task<double> CalculatePurePremiumAsync(double expectedClaims, double claimSeverity, double exposureUnits, double riskAdjustmentFactor)
    {
        _logger.LogInformation("Calculating pure premium: expectedClaims={ExpectedClaims}, severity={Severity}", 
            expectedClaims, claimSeverity);
        
        return _premiumCalculator.CalculatePurePremium(expectedClaims, claimSeverity, exposureUnits, riskAdjustmentFactor);
    }
}
