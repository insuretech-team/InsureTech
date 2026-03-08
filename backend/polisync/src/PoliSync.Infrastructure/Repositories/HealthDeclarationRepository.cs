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

public interface IHealthDeclarationRepository
{
    Task<HealthDeclaration> CreateAsync(HealthDeclaration entity, CancellationToken cancellationToken = default);
    Task<> GetByIdAsync(string id, CancellationToken cancellationToken = default);
    Task<List<HealthDeclaration>> GetAllAsync(int page = 1, int pageSize = 50, CancellationToken cancellationToken = default);
    Task<HealthDeclaration> UpdateAsync(HealthDeclaration entity, CancellationToken cancellationToken = default);
    Task DeleteAsync(string id, CancellationToken cancellationToken = default);
}

public class HealthDeclarationRepository : IHealthDeclarationRepository
{
    private readonly PoliSyncDbContext _context;
    private readonly ILogger<HealthDeclarationRepository> _logger;

    public HealthDeclarationRepository(PoliSyncDbContext context, ILogger<HealthDeclarationRepository> logger)
    {
        _context = context;
        _logger = logger;
    }

    public async Task<HealthDeclaration> CreateAsync(HealthDeclaration entity, CancellationToken cancellationToken = default)
    {
        var now = Timestamp.FromDateTime(DateTime.UtcNow);
        
        if (string.IsNullOrEmpty(entity.DeclarationId))
        {
            entity.DeclarationId = Guid.NewGuid().ToString();
        }

        _context.HealthDeclarations.Add(entity);
        await _context.SaveChangesAsync(cancellationToken);
        
        _logger.LogInformation("Created HealthDeclaration {Id}", entity.DeclarationId);
        return entity;
    }

    public async Task<> GetByIdAsync(string id, CancellationToken cancellationToken = default)
    {
        return await _context.HealthDeclarations
            .FirstOrDefaultAsync(e => e.DeclarationId == id, cancellationToken);
    }

    public async Task<List<HealthDeclaration>> GetAllAsync(int page = 1, int pageSize = 50, CancellationToken cancellationToken = default)
    {
        return await _context.HealthDeclarations
            .OrderByDescending(e => e.DeclarationId)
            .Skip((page - 1) * pageSize)
            .Take(pageSize)
            .ToListAsync(cancellationToken);
    }

    public async Task<HealthDeclaration> UpdateAsync(HealthDeclaration entity, CancellationToken cancellationToken = default)
    {
        _context.HealthDeclarations.Update(entity);
        await _context.SaveChangesAsync(cancellationToken);
        
        _logger.LogInformation("Updated HealthDeclaration {Id}", entity.DeclarationId);
        return entity;
    }

    public async Task DeleteAsync(string id, CancellationToken cancellationToken = default)
    {
        var entity = await GetByIdAsync(id, cancellationToken);
        if (entity != null)
        {
            _context.HealthDeclarations.Remove(entity);
            await _context.SaveChangesAsync(cancellationToken);
            _logger.LogInformation("Deleted HealthDeclaration {Id}", id);
        }
    }
}
