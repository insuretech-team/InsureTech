using Google.Protobuf.WellKnownTypes;
using Insuretech.Actuarial.Entity.V1;
using Microsoft.Extensions.Logging;

namespace PoliSync.Actuarial.Services;

public class ReserveCalculationService : IReserveCalculationService
{
    private readonly IReserveCalculationRepository _repository;
    private readonly ILogger<ReserveCalculationService> _logger;

    public ReserveCalculationService(
        IReserveCalculationRepository repository,
        ILogger<ReserveCalculationService> logger)
    {
        _repository = repository;
        _logger = logger;
    }

    public async Task<ReserveCalculation?> GetReserveAsync(string reserveId, CancellationToken cancellationToken = default)
    {
        return await _repository.GetByIdAsync(reserveId, cancellationToken);
    }

    public async Task<ReserveCalculation?> GetReserveByClaimAsync(string claimId, CancellationToken cancellationToken = default)
    {
        return await _repository.GetByClaimAsync(claimId, cancellationToken);
    }

    public Task<ReserveResult> CalculateReservesAsync(ReserveInput input, CancellationToken cancellationToken = default)
    {
        _logger.LogInformation("Calculating reserves for claim {ClaimId} using method {Method}", 
            input.ClaimId, input.CalculationMethod);
        
        // Simplified reserve calculation
        // In production, this would use chain ladder, BF method, etc.
        
        double caseReserve = input.CaseReserve;
        double reportedClaims = input.ReportedClaims;
        double paidClaims = input.PaidClaims;
        
        // IBNR calculation (simplified)
        // IBNR = Expected Ultimate Loss - Reported Claims
        double expectedUltimate = reportedClaims * 1.2; // 20% development factor
        double ibnrReserve = Math.Max(0, expectedUltimate - reportedClaims);
        
        // IBNER (Incurred But Not Enough Reported)
        double ibnerReserve = caseReserve * 0.1; // 10% of case reserve
        
        // Expense reserve (typically 10-15% of loss reserve)
        double expenseReserve = (caseReserve + ibnrReserve) * 0.12;
        
        double totalReserve = caseReserve + ibnrReserve + ibnerReserve + expenseReserve;
        
        // Calculate confidence interval (simplified)
        double stdDev = totalReserve * 0.15; // 15% standard deviation
        double zScore = 1.96; // 95% confidence
        double lowerBound = Math.Max(0, totalReserve - zScore * stdDev);
        double upperBound = totalReserve + zScore * stdDev;
        
        var result = new ReserveResult
        {
            CaseReserve = caseReserve,
            IbnrReserve = ibnrReserve,
            IbnerReserve = ibnerReserve,
            ExpenseReserve = expenseReserve,
            TotalReserve = totalReserve,
            LowerBound = lowerBound,
            UpperBound = upperBound,
            MethodUsed = input.CalculationMethod
        };
        
        return Task.FromResult(result);
    }

    public async Task<ReserveCalculation> SaveReserveAsync(ReserveCalculation reserve, CancellationToken cancellationToken = default)
    {
        if (string.IsNullOrEmpty(reserve.ReserveId))
        {
            reserve.ReserveId = Guid.NewGuid().ToString();
            reserve.CreatedAt = Timestamp.FromDateTime(DateTime.UtcNow);
        }
        
        reserve.UpdatedAt = Timestamp.FromDateTime(DateTime.UtcNow);
        
        return await _repository.CreateAsync(reserve, cancellationToken);
    }
}
