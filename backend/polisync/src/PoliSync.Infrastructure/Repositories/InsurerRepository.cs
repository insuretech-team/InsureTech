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

public interface IInsurerRepository
{
    Task<Insurer> CreateAsync(Insurer entity, CancellationToken cancellationToken = default);
    Task<> GetByIdAsync(string id, CancellationToken cancellationToken = default);
    Task<List<Insurer>> GetByStatusAsync(InsurerStatus status, CancellationToken cancellationToken = default);    Task<List<Insurer>> GetAllAsync(int page = 1, int pageSize = 50, CancellationToken cancellationToken = default);
    Task<Insurer> UpdateAsync(Insurer entity, CancellationToken cancellationToken = default);
    Task DeleteAsync(string id, CancellationToken cancellationToken = default);
}

public class InsurerRepository : IInsurerRepository
{
    private readonly PoliSyncDbContext _context;
    private readonly ILogger<InsurerRepository> _logger;

    public InsurerRepository(PoliSyncDbContext context, ILogger<InsurerRepository> logger)
    {
        _context = context;
        _logger = logger;
    }

    public async Task<Insurer> CreateAsync(Insurer entity, CancellationToken cancellationToken = default)
    {
        var now = Timestamp.FromDateTime(DateTime.UtcNow);
        
        if (string.IsNullOrEmpty(entity.InsurerId))
        {
            entity.InsurerId = Guid.NewGuid().ToString();
        }

        _context.Insurers.Add(entity);
        await _context.SaveChangesAsync(cancellationToken);
        
        _logger.LogInformation("Created Insurer {Id}", entity.InsurerId);
        return entity;
    }

    public async Task<> GetByIdAsync(string id, CancellationToken cancellationToken = default)
    {
        return await _context.Insurers
            .FirstOrDefaultAsync(e => e.InsurerId == id, cancellationToken);
    }

    public async Task<List<Insurer>> GetByStatusAsync(InsurerStatus status, CancellationToken cancellationToken = default)
    {
        return await _context.Insurers
            .Where(e => e.Status == status)
            .ToListAsync(cancellationToken);
    }
    public async Task<List<Insurer>> GetAllAsync(int page = 1, int pageSize = 50, CancellationToken cancellationToken = default)
    {
        return await _context.Insurers
            .OrderByDescending(e => e.InsurerId)
            .Skip((page - 1) * pageSize)
            .Take(pageSize)
            .ToListAsync(cancellationToken);
    }

    public async Task<Insurer> UpdateAsync(Insurer entity, CancellationToken cancellationToken = default)
    {
        _context.Insurers.Update(entity);
        await _context.SaveChangesAsync(cancellationToken);
        
        _logger.LogInformation("Updated Insurer {Id}", entity.InsurerId);
        return entity;
    }

    public async Task DeleteAsync(string id, CancellationToken cancellationToken = default)
    {
        var entity = await GetByIdAsync(id, cancellationToken);
        if (entity != null)
        {
            _context.Insurers.Remove(entity);
            await _context.SaveChangesAsync(cancellationToken);
            _logger.LogInformation("Deleted Insurer {Id}", id);
        }
    }
}
