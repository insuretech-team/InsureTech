using Insuretech.Actuarial.Entity.V1;
using System.Collections.Concurrent;

namespace PoliSync.Actuarial.Services;

public class InMemoryActuarialCalculationRepository : IActuarialCalculationRepository
{
    private readonly ConcurrentDictionary<string, ActuarialCalculation> _calculations = new();

    public Task<ActuarialCalculation?> GetByIdAsync(string calculationId, CancellationToken cancellationToken = default)
    {
        _calculations.TryGetValue(calculationId, out var calculation);
        return Task.FromResult(calculation);
    }

    public Task<ActuarialCalculation?> GetByReferenceAsync(string reference, CancellationToken cancellationToken = default)
    {
        var calculation = _calculations.Values.FirstOrDefault(c => c.CalculationReference == reference);
        return Task.FromResult(calculation);
    }

    public Task<IEnumerable<ActuarialCalculation>> GetAllAsync(CancellationToken cancellationToken = default)
    {
        return Task.FromResult(_calculations.Values.AsEnumerable());
    }

    public Task<IEnumerable<ActuarialCalculation>> GetByFiltersAsync(
        ActuarialCalculationType? type,
        string? entityType,
        string? entityId,
        DateTime? from,
        DateTime? to,
        CancellationToken cancellationToken = default)
    {
        var query = _calculations.Values.AsEnumerable();
        
        if (type.HasValue)
        {
            query = query.Where(c => c.CalculationType == type.Value);
        }
        
        if (!string.IsNullOrEmpty(entityType))
        {
            query = query.Where(c => c.EntityType == entityType);
        }
        
        if (!string.IsNullOrEmpty(entityId))
        {
            query = query.Where(c => c.EntityId == entityId);
        }
        
        if (from.HasValue)
        {
            query = query.Where(c => c.CalculatedAt.ToDateTime() >= from.Value);
        }
        
        if (to.HasValue)
        {
            query = query.Where(c => c.CalculatedAt.ToDateTime() <= to.Value);
        }
        
        return Task.FromResult(query);
    }

    public Task<ActuarialCalculation> CreateAsync(ActuarialCalculation calculation, CancellationToken cancellationToken = default)
    {
        _calculations.TryAdd(calculation.CalculationId, calculation);
        return Task.FromResult(calculation);
    }
}

public class InMemoryReserveCalculationRepository : IReserveCalculationRepository
{
    private readonly ConcurrentDictionary<string, ReserveCalculation> _reserves = new();

    public Task<ReserveCalculation?> GetByIdAsync(string reserveId, CancellationToken cancellationToken = default)
    {
        _reserves.TryGetValue(reserveId, out var reserve);
        return Task.FromResult(reserve);
    }

    public Task<ReserveCalculation?> GetByClaimAsync(string claimId, CancellationToken cancellationToken = default)
    {
        var reserve = _reserves.Values.FirstOrDefault(r => r.ClaimId == claimId);
        return Task.FromResult(reserve);
    }

    public Task<IEnumerable<ReserveCalculation>> GetAllAsync(CancellationToken cancellationToken = default)
    {
        return Task.FromResult(_reserves.Values.AsEnumerable());
    }

    public Task<ReserveCalculation> CreateAsync(ReserveCalculation reserve, CancellationToken cancellationToken = default)
    {
        _reserves.TryAdd(reserve.ReserveId, reserve);
        return Task.FromResult(reserve);
    }

    public Task<ReserveCalculation?> UpdateAsync(ReserveCalculation reserve, CancellationToken cancellationToken = default)
    {
        if (_reserves.ContainsKey(reserve.ReserveId))
        {
            _reserves[reserve.ReserveId] = reserve;
            return Task.FromResult<ReserveCalculation?>(reserve);
        }
        return Task.FromResult<ReserveCalculation?>(null);
    }
}

public class InMemoryLossRatioCalculationRepository : ILossRatioCalculationRepository
{
    private readonly ConcurrentDictionary<string, LossRatioCalculation> _lossRatios = new();

    public Task<LossRatioCalculation?> GetByIdAsync(string lossRatioId, CancellationToken cancellationToken = default)
    {
        _lossRatios.TryGetValue(lossRatioId, out var lossRatio);
        return Task.FromResult(lossRatio);
    }

    public Task<IEnumerable<LossRatioCalculation>> GetAllAsync(CancellationToken cancellationToken = default)
    {
        return Task.FromResult(_lossRatios.Values.AsEnumerable());
    }

    public Task<IEnumerable<LossRatioCalculation>> GetByFiltersAsync(
        string? productId,
        string? lineOfBusiness,
        DateTime? from,
        DateTime? to,
        CancellationToken cancellationToken = default)
    {
        var query = _lossRatios.Values.AsEnumerable();
        
        if (!string.IsNullOrEmpty(productId))
        {
            query = query.Where(l => l.ProductId == productId);
        }
        
        if (!string.IsNullOrEmpty(lineOfBusiness))
        {
            query = query.Where(l => l.LineOfBusiness == lineOfBusiness);
        }
        
        if (from.HasValue)
        {
            query = query.Where(l => l.PeriodStart.ToDateTime() >= from.Value);
        }
        
        if (to.HasValue)
        {
            query = query.Where(l => l.PeriodEnd.ToDateTime() <= to.Value);
        }
        
        return Task.FromResult(query);
    }

    public Task<LossRatioCalculation> CreateAsync(LossRatioCalculation lossRatio, CancellationToken cancellationToken = default)
    {
        _lossRatios.TryAdd(lossRatio.LossRatioId, lossRatio);
        return Task.FromResult(lossRatio);
    }
}
