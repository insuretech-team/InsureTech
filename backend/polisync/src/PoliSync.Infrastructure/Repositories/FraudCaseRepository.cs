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

public interface IFraudCaseRepository
{
    Task<FraudCase> CreateAsync(FraudCase entity, CancellationToken cancellationToken = default);
    Task<> GetByIdAsync(string id, CancellationToken cancellationToken = default);
    Task<> GetByNumberAsync(string casenumber, CancellationToken cancellationToken = default);    Task<List<FraudCase>> GetByStatusAsync(CaseStatus status, CancellationToken cancellationToken = default);    Task<List<FraudCase>> GetAllAsync(int page = 1, int pageSize = 50, CancellationToken cancellationToken = default);
    Task<FraudCase> UpdateAsync(FraudCase entity, CancellationToken cancellationToken = default);
    Task DeleteAsync(string id, CancellationToken cancellationToken = default);
}

public class FraudCaseRepository : IFraudCaseRepository
{
    private readonly PoliSyncDbContext _context;
    private readonly ILogger<FraudCaseRepository> _logger;

    public FraudCaseRepository(PoliSyncDbContext context, ILogger<FraudCaseRepository> logger)
    {
        _context = context;
        _logger = logger;
    }

    public async Task<FraudCase> CreateAsync(FraudCase entity, CancellationToken cancellationToken = default)
    {
        var now = Timestamp.FromDateTime(DateTime.UtcNow);
        
        if (string.IsNullOrEmpty(entity.CaseId))
        {
            entity.CaseId = Guid.NewGuid().ToString();
        }

        _context.FraudCases.Add(entity);
        await _context.SaveChangesAsync(cancellationToken);
        
        _logger.LogInformation("Created FraudCase {Id}", entity.CaseId);
        return entity;
    }

    public async Task<> GetByIdAsync(string id, CancellationToken cancellationToken = default)
    {
        return await _context.FraudCases
            .FirstOrDefaultAsync(e => e.CaseId == id, cancellationToken);
    }

    public async Task<> GetByNumberAsync(string casenumber, CancellationToken cancellationToken = default)
    {
        return await _context.FraudCases
            .FirstOrDefaultAsync(e => e.CaseNumber == casenumber, cancellationToken);
    }
    public async Task<List<FraudCase>> GetByStatusAsync(CaseStatus status, CancellationToken cancellationToken = default)
    {
        return await _context.FraudCases
            .Where(e => e.Status == status)
            .ToListAsync(cancellationToken);
    }
    public async Task<List<FraudCase>> GetAllAsync(int page = 1, int pageSize = 50, CancellationToken cancellationToken = default)
    {
        return await _context.FraudCases
            .OrderByDescending(e => e.CaseId)
            .Skip((page - 1) * pageSize)
            .Take(pageSize)
            .ToListAsync(cancellationToken);
    }

    public async Task<FraudCase> UpdateAsync(FraudCase entity, CancellationToken cancellationToken = default)
    {
        _context.FraudCases.Update(entity);
        await _context.SaveChangesAsync(cancellationToken);
        
        _logger.LogInformation("Updated FraudCase {Id}", entity.CaseId);
        return entity;
    }

    public async Task DeleteAsync(string id, CancellationToken cancellationToken = default)
    {
        var entity = await GetByIdAsync(id, cancellationToken);
        if (entity != null)
        {
            _context.FraudCases.Remove(entity);
            await _context.SaveChangesAsync(cancellationToken);
            _logger.LogInformation("Deleted FraudCase {Id}", id);
        }
    }
}
