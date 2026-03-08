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

public interface IRenewalScheduleRepository
{
    Task<RenewalSchedule> CreateAsync(RenewalSchedule entity, CancellationToken cancellationToken = default);
    Task<> GetByIdAsync(string id, CancellationToken cancellationToken = default);
    Task<List<RenewalSchedule>> GetByStatusAsync(RenewalStatus status, CancellationToken cancellationToken = default);    Task<List<RenewalSchedule>> GetAllAsync(int page = 1, int pageSize = 50, CancellationToken cancellationToken = default);
    Task<RenewalSchedule> UpdateAsync(RenewalSchedule entity, CancellationToken cancellationToken = default);
    Task DeleteAsync(string id, CancellationToken cancellationToken = default);
}

public class RenewalScheduleRepository : IRenewalScheduleRepository
{
    private readonly PoliSyncDbContext _context;
    private readonly ILogger<RenewalScheduleRepository> _logger;

    public RenewalScheduleRepository(PoliSyncDbContext context, ILogger<RenewalScheduleRepository> logger)
    {
        _context = context;
        _logger = logger;
    }

    public async Task<RenewalSchedule> CreateAsync(RenewalSchedule entity, CancellationToken cancellationToken = default)
    {
        var now = Timestamp.FromDateTime(DateTime.UtcNow);
        
        if (string.IsNullOrEmpty(entity.ScheduleId))
        {
            entity.ScheduleId = Guid.NewGuid().ToString();
        }

        _context.RenewalSchedules.Add(entity);
        await _context.SaveChangesAsync(cancellationToken);
        
        _logger.LogInformation("Created RenewalSchedule {Id}", entity.ScheduleId);
        return entity;
    }

    public async Task<> GetByIdAsync(string id, CancellationToken cancellationToken = default)
    {
        return await _context.RenewalSchedules
            .FirstOrDefaultAsync(e => e.ScheduleId == id, cancellationToken);
    }

    public async Task<List<RenewalSchedule>> GetByStatusAsync(RenewalStatus status, CancellationToken cancellationToken = default)
    {
        return await _context.RenewalSchedules
            .Where(e => e.Status == status)
            .ToListAsync(cancellationToken);
    }
    public async Task<List<RenewalSchedule>> GetAllAsync(int page = 1, int pageSize = 50, CancellationToken cancellationToken = default)
    {
        return await _context.RenewalSchedules
            .OrderByDescending(e => e.ScheduleId)
            .Skip((page - 1) * pageSize)
            .Take(pageSize)
            .ToListAsync(cancellationToken);
    }

    public async Task<RenewalSchedule> UpdateAsync(RenewalSchedule entity, CancellationToken cancellationToken = default)
    {
        _context.RenewalSchedules.Update(entity);
        await _context.SaveChangesAsync(cancellationToken);
        
        _logger.LogInformation("Updated RenewalSchedule {Id}", entity.ScheduleId);
        return entity;
    }

    public async Task DeleteAsync(string id, CancellationToken cancellationToken = default)
    {
        var entity = await GetByIdAsync(id, cancellationToken);
        if (entity != null)
        {
            _context.RenewalSchedules.Remove(entity);
            await _context.SaveChangesAsync(cancellationToken);
            _logger.LogInformation("Deleted RenewalSchedule {Id}", id);
        }
    }
}
