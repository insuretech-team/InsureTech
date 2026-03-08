using Google.Protobuf.WellKnownTypes;
using Insuretech.Policy.Entity.V1;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Logging;
using PoliSync.Infrastructure.Persistence;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;

namespace PoliSync.Infrastructure.Repositories;

public interface IGracePeriodRepository
{
    Task<GracePeriod> CreateAsync(GracePeriod entity, CancellationToken cancellationToken = default);
    Task<> GetByIdAsync(string id, CancellationToken cancellationToken = default);
    Task<List<GracePeriod>> GetByStatusAsync(GracePeriodStatus status, CancellationToken cancellationToken = default);    Task<List<GracePeriod>> GetAllAsync(int page = 1, int pageSize = 50, CancellationToken cancellationToken = default);
    Task<GracePeriod> UpdateAsync(GracePeriod entity, CancellationToken cancellationToken = default);
    Task DeleteAsync(string id, CancellationToken cancellationToken = default);
}

public class GracePeriodRepository : IGracePeriodRepository
{
    private readonly PoliSyncDbContext _context;
    private readonly ILogger<GracePeriodRepository> _logger;

    public GracePeriodRepository(PoliSyncDbContext context, ILogger<GracePeriodRepository> logger)
    {
        _context = context;
        _logger = logger;
    }

    public async Task<GracePeriod> CreateAsync(GracePeriod entity, CancellationToken cancellationToken = default)
    {
        var now = Timestamp.FromDateTime(DateTime.UtcNow);
        
        if (string.IsNullOrEmpty(entity.GracePeriodId))
        {
            entity.GracePeriodId = Guid.NewGuid().ToString();
        }

        _context.GracePeriods.Add(entity);
        await _context.SaveChangesAsync(cancellationToken);
        
        _logger.LogInformation("Created GracePeriod {Id}", entity.GracePeriodId);
        return entity;
    }

    public async Task<> GetByIdAsync(string id, CancellationToken cancellationToken = default)
    {
        return await _context.GracePeriods
            .FirstOrDefaultAsync(e => e.GracePeriodId == id, cancellationToken);
    }

    public async Task<List<GracePeriod>> GetByStatusAsync(GracePeriodStatus status, CancellationToken cancellationToken = default)
    {
        return await _context.GracePeriods
            .Where(e => e.Status == status)
            .ToListAsync(cancellationToken);
    }
    public async Task<List<GracePeriod>> GetAllAsync(int page = 1, int pageSize = 50, CancellationToken cancellationToken = default)
    {
        return await _context.GracePeriods
            .OrderByDescending(e => e.GracePeriodId)
            .Skip((page - 1) * pageSize)
            .Take(pageSize)
            .ToListAsync(cancellationToken);
    }

    public async Task<GracePeriod> UpdateAsync(GracePeriod entity, CancellationToken cancellationToken = default)
    {
        _context.GracePeriods.Update(entity);
        await _context.SaveChangesAsync(cancellationToken);
        
        _logger.LogInformation("Updated GracePeriod {Id}", entity.GracePeriodId);
        return entity;
    }

    public async Task DeleteAsync(string id, CancellationToken cancellationToken = default)
    {
        var entity = await GetByIdAsync(id, cancellationToken);
        if (entity != null)
        {
            _context.GracePeriods.Remove(entity);
            await _context.SaveChangesAsync(cancellationToken);
            _logger.LogInformation("Deleted GracePeriod {Id}", id);
        }
    }
}
