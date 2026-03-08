using Google.Protobuf.WellKnownTypes;
using Insuretech.Products.Entity.V1;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Logging;
using PoliSync.Infrastructure.Persistence;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;

namespace PoliSync.Infrastructure.Repositories;

public interface IRiderRepository
{
    Task<Rider> CreateAsync(Rider entity, CancellationToken cancellationToken = default);
    Task<> GetByIdAsync(string id, CancellationToken cancellationToken = default);
    Task<List<Rider>> GetByStatusAsync(RiderStatus status, CancellationToken cancellationToken = default);    Task<List<Rider>> GetAllAsync(int page = 1, int pageSize = 50, CancellationToken cancellationToken = default);
    Task<Rider> UpdateAsync(Rider entity, CancellationToken cancellationToken = default);
    Task DeleteAsync(string id, CancellationToken cancellationToken = default);
}

public class RiderRepository : IRiderRepository
{
    private readonly PoliSyncDbContext _context;
    private readonly ILogger<RiderRepository> _logger;

    public RiderRepository(PoliSyncDbContext context, ILogger<RiderRepository> logger)
    {
        _context = context;
        _logger = logger;
    }

    public async Task<Rider> CreateAsync(Rider entity, CancellationToken cancellationToken = default)
    {
        var now = Timestamp.FromDateTime(DateTime.UtcNow);
        
        if (string.IsNullOrEmpty(entity.RiderId))
        {
            entity.RiderId = Guid.NewGuid().ToString();
        }

        _context.Riders.Add(entity);
        await _context.SaveChangesAsync(cancellationToken);
        
        _logger.LogInformation("Created Rider {Id}", entity.RiderId);
        return entity;
    }

    public async Task<> GetByIdAsync(string id, CancellationToken cancellationToken = default)
    {
        return await _context.Riders
            .FirstOrDefaultAsync(e => e.RiderId == id, cancellationToken);
    }

    public async Task<List<Rider>> GetByStatusAsync(RiderStatus status, CancellationToken cancellationToken = default)
    {
        return await _context.Riders
            .Where(e => e.Status == status)
            .ToListAsync(cancellationToken);
    }
    public async Task<List<Rider>> GetAllAsync(int page = 1, int pageSize = 50, CancellationToken cancellationToken = default)
    {
        return await _context.Riders
            .OrderByDescending(e => e.RiderId)
            .Skip((page - 1) * pageSize)
            .Take(pageSize)
            .ToListAsync(cancellationToken);
    }

    public async Task<Rider> UpdateAsync(Rider entity, CancellationToken cancellationToken = default)
    {
        _context.Riders.Update(entity);
        await _context.SaveChangesAsync(cancellationToken);
        
        _logger.LogInformation("Updated Rider {Id}", entity.RiderId);
        return entity;
    }

    public async Task DeleteAsync(string id, CancellationToken cancellationToken = default)
    {
        var entity = await GetByIdAsync(id, cancellationToken);
        if (entity != null)
        {
            _context.Riders.Remove(entity);
            await _context.SaveChangesAsync(cancellationToken);
            _logger.LogInformation("Deleted Rider {Id}", id);
        }
    }
}
