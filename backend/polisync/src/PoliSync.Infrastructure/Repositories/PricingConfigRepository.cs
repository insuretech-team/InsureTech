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

public interface IPricingConfigRepository
{
    Task<PricingConfig> CreateAsync(PricingConfig entity, CancellationToken cancellationToken = default);
    Task<> GetByIdAsync(string id, CancellationToken cancellationToken = default);
    Task<List<PricingConfig>> GetAllAsync(int page = 1, int pageSize = 50, CancellationToken cancellationToken = default);
    Task<PricingConfig> UpdateAsync(PricingConfig entity, CancellationToken cancellationToken = default);
    Task DeleteAsync(string id, CancellationToken cancellationToken = default);
}

public class PricingConfigRepository : IPricingConfigRepository
{
    private readonly PoliSyncDbContext _context;
    private readonly ILogger<PricingConfigRepository> _logger;

    public PricingConfigRepository(PoliSyncDbContext context, ILogger<PricingConfigRepository> logger)
    {
        _context = context;
        _logger = logger;
    }

    public async Task<PricingConfig> CreateAsync(PricingConfig entity, CancellationToken cancellationToken = default)
    {
        var now = Timestamp.FromDateTime(DateTime.UtcNow);
        
        if (string.IsNullOrEmpty(entity.ConfigId))
        {
            entity.ConfigId = Guid.NewGuid().ToString();
        }

        _context.PricingConfigs.Add(entity);
        await _context.SaveChangesAsync(cancellationToken);
        
        _logger.LogInformation("Created PricingConfig {Id}", entity.ConfigId);
        return entity;
    }

    public async Task<> GetByIdAsync(string id, CancellationToken cancellationToken = default)
    {
        return await _context.PricingConfigs
            .FirstOrDefaultAsync(e => e.ConfigId == id, cancellationToken);
    }

    public async Task<List<PricingConfig>> GetAllAsync(int page = 1, int pageSize = 50, CancellationToken cancellationToken = default)
    {
        return await _context.PricingConfigs
            .OrderByDescending(e => e.ConfigId)
            .Skip((page - 1) * pageSize)
            .Take(pageSize)
            .ToListAsync(cancellationToken);
    }

    public async Task<PricingConfig> UpdateAsync(PricingConfig entity, CancellationToken cancellationToken = default)
    {
        _context.PricingConfigs.Update(entity);
        await _context.SaveChangesAsync(cancellationToken);
        
        _logger.LogInformation("Updated PricingConfig {Id}", entity.ConfigId);
        return entity;
    }

    public async Task DeleteAsync(string id, CancellationToken cancellationToken = default)
    {
        var entity = await GetByIdAsync(id, cancellationToken);
        if (entity != null)
        {
            _context.PricingConfigs.Remove(entity);
            await _context.SaveChangesAsync(cancellationToken);
            _logger.LogInformation("Deleted PricingConfig {Id}", id);
        }
    }
}
