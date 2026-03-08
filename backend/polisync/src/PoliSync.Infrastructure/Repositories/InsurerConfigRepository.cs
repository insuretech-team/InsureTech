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

public interface IInsurerConfigRepository
{
    Task<InsurerConfig> CreateAsync(InsurerConfig entity, CancellationToken cancellationToken = default);
    Task<> GetByIdAsync(string id, CancellationToken cancellationToken = default);
    Task<List<InsurerConfig>> GetAllAsync(int page = 1, int pageSize = 50, CancellationToken cancellationToken = default);
    Task<InsurerConfig> UpdateAsync(InsurerConfig entity, CancellationToken cancellationToken = default);
    Task DeleteAsync(string id, CancellationToken cancellationToken = default);
}

public class InsurerConfigRepository : IInsurerConfigRepository
{
    private readonly PoliSyncDbContext _context;
    private readonly ILogger<InsurerConfigRepository> _logger;

    public InsurerConfigRepository(PoliSyncDbContext context, ILogger<InsurerConfigRepository> logger)
    {
        _context = context;
        _logger = logger;
    }

    public async Task<InsurerConfig> CreateAsync(InsurerConfig entity, CancellationToken cancellationToken = default)
    {
        var now = Timestamp.FromDateTime(DateTime.UtcNow);
        
        if (string.IsNullOrEmpty(entity.ConfigId))
        {
            entity.ConfigId = Guid.NewGuid().ToString();
        }

        _context.InsurerConfigs.Add(entity);
        await _context.SaveChangesAsync(cancellationToken);
        
        _logger.LogInformation("Created InsurerConfig {Id}", entity.ConfigId);
        return entity;
    }

    public async Task<> GetByIdAsync(string id, CancellationToken cancellationToken = default)
    {
        return await _context.InsurerConfigs
            .FirstOrDefaultAsync(e => e.ConfigId == id, cancellationToken);
    }

    public async Task<List<InsurerConfig>> GetAllAsync(int page = 1, int pageSize = 50, CancellationToken cancellationToken = default)
    {
        return await _context.InsurerConfigs
            .OrderByDescending(e => e.ConfigId)
            .Skip((page - 1) * pageSize)
            .Take(pageSize)
            .ToListAsync(cancellationToken);
    }

    public async Task<InsurerConfig> UpdateAsync(InsurerConfig entity, CancellationToken cancellationToken = default)
    {
        _context.InsurerConfigs.Update(entity);
        await _context.SaveChangesAsync(cancellationToken);
        
        _logger.LogInformation("Updated InsurerConfig {Id}", entity.ConfigId);
        return entity;
    }

    public async Task DeleteAsync(string id, CancellationToken cancellationToken = default)
    {
        var entity = await GetByIdAsync(id, cancellationToken);
        if (entity != null)
        {
            _context.InsurerConfigs.Remove(entity);
            await _context.SaveChangesAsync(cancellationToken);
            _logger.LogInformation("Deleted InsurerConfig {Id}", id);
        }
    }
}
