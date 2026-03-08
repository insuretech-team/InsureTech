using Google.Protobuf.WellKnownTypes;
using Insuretech.Claims.Entity.V1;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Logging;
using PoliSync.Infrastructure.Persistence;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;

namespace PoliSync.Infrastructure.Repositories;

public interface IFraudAlertRepository
{
    Task<FraudAlert> CreateAsync(FraudAlert entity, CancellationToken cancellationToken = default);
    Task<> GetByIdAsync(string id, CancellationToken cancellationToken = default);
    Task<List<FraudAlert>> GetByStatusAsync(AlertStatus status, CancellationToken cancellationToken = default);    Task<List<FraudAlert>> GetAllAsync(int page = 1, int pageSize = 50, CancellationToken cancellationToken = default);
    Task<FraudAlert> UpdateAsync(FraudAlert entity, CancellationToken cancellationToken = default);
    Task DeleteAsync(string id, CancellationToken cancellationToken = default);
}

public class FraudAlertRepository : IFraudAlertRepository
{
    private readonly PoliSyncDbContext _context;
    private readonly ILogger<FraudAlertRepository> _logger;

    public FraudAlertRepository(PoliSyncDbContext context, ILogger<FraudAlertRepository> logger)
    {
        _context = context;
        _logger = logger;
    }

    public async Task<FraudAlert> CreateAsync(FraudAlert entity, CancellationToken cancellationToken = default)
    {
        var now = Timestamp.FromDateTime(DateTime.UtcNow);
        
        if (string.IsNullOrEmpty(entity.AlertId))
        {
            entity.AlertId = Guid.NewGuid().ToString();
        }

        _context.FraudAlerts.Add(entity);
        await _context.SaveChangesAsync(cancellationToken);
        
        _logger.LogInformation("Created FraudAlert {Id}", entity.AlertId);
        return entity;
    }

    public async Task<> GetByIdAsync(string id, CancellationToken cancellationToken = default)
    {
        return await _context.FraudAlerts
            .FirstOrDefaultAsync(e => e.AlertId == id, cancellationToken);
    }

    public async Task<List<FraudAlert>> GetByStatusAsync(AlertStatus status, CancellationToken cancellationToken = default)
    {
        return await _context.FraudAlerts
            .Where(e => e.Status == status)
            .ToListAsync(cancellationToken);
    }
    public async Task<List<FraudAlert>> GetAllAsync(int page = 1, int pageSize = 50, CancellationToken cancellationToken = default)
    {
        return await _context.FraudAlerts
            .OrderByDescending(e => e.AlertId)
            .Skip((page - 1) * pageSize)
            .Take(pageSize)
            .ToListAsync(cancellationToken);
    }

    public async Task<FraudAlert> UpdateAsync(FraudAlert entity, CancellationToken cancellationToken = default)
    {
        _context.FraudAlerts.Update(entity);
        await _context.SaveChangesAsync(cancellationToken);
        
        _logger.LogInformation("Updated FraudAlert {Id}", entity.AlertId);
        return entity;
    }

    public async Task DeleteAsync(string id, CancellationToken cancellationToken = default)
    {
        var entity = await GetByIdAsync(id, cancellationToken);
        if (entity != null)
        {
            _context.FraudAlerts.Remove(entity);
            await _context.SaveChangesAsync(cancellationToken);
            _logger.LogInformation("Deleted FraudAlert {Id}", id);
        }
    }
}
