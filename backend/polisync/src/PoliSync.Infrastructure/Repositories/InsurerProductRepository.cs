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

public interface IInsurerProductRepository
{
    Task<InsurerProduct> CreateAsync(InsurerProduct entity, CancellationToken cancellationToken = default);
    Task<> GetByIdAsync(string id, CancellationToken cancellationToken = default);
    Task<List<InsurerProduct>> GetAllAsync(int page = 1, int pageSize = 50, CancellationToken cancellationToken = default);
    Task<InsurerProduct> UpdateAsync(InsurerProduct entity, CancellationToken cancellationToken = default);
    Task DeleteAsync(string id, CancellationToken cancellationToken = default);
}

public class InsurerProductRepository : IInsurerProductRepository
{
    private readonly PoliSyncDbContext _context;
    private readonly ILogger<InsurerProductRepository> _logger;

    public InsurerProductRepository(PoliSyncDbContext context, ILogger<InsurerProductRepository> logger)
    {
        _context = context;
        _logger = logger;
    }

    public async Task<InsurerProduct> CreateAsync(InsurerProduct entity, CancellationToken cancellationToken = default)
    {
        var now = Timestamp.FromDateTime(DateTime.UtcNow);
        
        if (string.IsNullOrEmpty(entity.InsurerProductId))
        {
            entity.InsurerProductId = Guid.NewGuid().ToString();
        }

        _context.InsurerProducts.Add(entity);
        await _context.SaveChangesAsync(cancellationToken);
        
        _logger.LogInformation("Created InsurerProduct {Id}", entity.InsurerProductId);
        return entity;
    }

    public async Task<> GetByIdAsync(string id, CancellationToken cancellationToken = default)
    {
        return await _context.InsurerProducts
            .FirstOrDefaultAsync(e => e.InsurerProductId == id, cancellationToken);
    }

    public async Task<List<InsurerProduct>> GetAllAsync(int page = 1, int pageSize = 50, CancellationToken cancellationToken = default)
    {
        return await _context.InsurerProducts
            .OrderByDescending(e => e.InsurerProductId)
            .Skip((page - 1) * pageSize)
            .Take(pageSize)
            .ToListAsync(cancellationToken);
    }

    public async Task<InsurerProduct> UpdateAsync(InsurerProduct entity, CancellationToken cancellationToken = default)
    {
        _context.InsurerProducts.Update(entity);
        await _context.SaveChangesAsync(cancellationToken);
        
        _logger.LogInformation("Updated InsurerProduct {Id}", entity.InsurerProductId);
        return entity;
    }

    public async Task DeleteAsync(string id, CancellationToken cancellationToken = default)
    {
        var entity = await GetByIdAsync(id, cancellationToken);
        if (entity != null)
        {
            _context.InsurerProducts.Remove(entity);
            await _context.SaveChangesAsync(cancellationToken);
            _logger.LogInformation("Deleted InsurerProduct {Id}", id);
        }
    }
}
