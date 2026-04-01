using Google.Protobuf.WellKnownTypes;
using Insuretech.Actuarial.Entity.V1;
using Microsoft.Extensions.Logging;

namespace PoliSync.Actuarial.Services;

public class LossRatioService : ILossRatioService
{
    private readonly ILossRatioCalculationRepository _repository;
    private readonly ILogger<LossRatioService> _logger;

    public LossRatioService(
        ILossRatioCalculationRepository repository,
        ILogger<LossRatioService> logger)
    {
        _repository = repository;
        _logger = logger;
    }

    public async Task<LossRatioCalculation?> GetLossRatioAsync(string lossRatioId, CancellationToken cancellationToken = default)
    {
        return await _repository.GetByIdAsync(lossRatioId, cancellationToken);
    }

    public async Task<IEnumerable<LossRatioCalculation>> GetLossRatiosAsync(
        string? productId = null,
        string? lineOfBusiness = null,
        DateTime? from = null,
        DateTime? to = null,
        CancellationToken cancellationToken = default)
    {
        return await _repository.GetByFiltersAsync(productId, lineOfBusiness, from, to, cancellationToken);
    }

    public Task<LossRatioResult> CalculateLossRatioAsync(LossRatioInput input, CancellationToken cancellationToken = default)
    {
        _logger.LogInformation("Calculating loss ratio for product {ProductId}", input.ProductId);
        
        double earnedPremium = input.EarnedPremium;
        double writtenPremium = input.WrittenPremium;
        double incurredLosses = input.IncurredLosses;
        double lossAdjustmentExpenses = input.LossAdjustmentExpenses;
        double operatingExpenses = input.OperatingExpenses;
        
        // Loss Ratio = (Incurred Losses + LAE) / Earned Premium
        double totalIncurred = incurredLosses + lossAdjustmentExpenses;
        double lossRatio = earnedPremium > 0 ? totalIncurred / earnedPremium : 0;
        
        // Expense Ratio = Operating Expenses / Written Premium
        double expenseRatio = writtenPremium > 0 ? operatingExpenses / writtenPremium : 0;
        
        // Combined Ratio = Loss Ratio + Expense Ratio
        double combinedRatio = lossRatio + expenseRatio;
        
        // Underwriting Profit Margin = 1 - Combined Ratio
        double underwritingProfitMargin = 1 - combinedRatio;
        
        // Interpretation
        string interpretation;
        if (combinedRatio < 0.95)
        {
            interpretation = "PROFITABLE";
        }
        else if (combinedRatio <= 1.05)
        {
            interpretation = "BREAK_EVEN";
        }
        else
        {
            interpretation = "LOSS_MAKING";
        }
        
        var result = new LossRatioResult
        {
            LossRatio = Math.Round(lossRatio, 4),
            ExpenseRatio = Math.Round(expenseRatio, 4),
            CombinedRatio = Math.Round(combinedRatio, 4),
            UnderwritingProfitMargin = Math.Round(underwritingProfitMargin, 4),
            Interpretation = interpretation
        };
        
        return Task.FromResult(result);
    }

    public async Task<LossRatioCalculation> SaveLossRatioAsync(LossRatioCalculation lossRatio, CancellationToken cancellationToken = default)
    {
        if (string.IsNullOrEmpty(lossRatio.LossRatioId))
        {
            lossRatio.LossRatioId = Guid.NewGuid().ToString();
            lossRatio.CreatedAt = Timestamp.FromDateTime(DateTime.UtcNow);
        }
        
        lossRatio.UpdatedAt = Timestamp.FromDateTime(DateTime.UtcNow);
        
        return await _repository.CreateAsync(lossRatio, cancellationToken);
    }
}
