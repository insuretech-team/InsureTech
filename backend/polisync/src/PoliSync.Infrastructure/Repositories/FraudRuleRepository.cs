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

public interface IFraudRuleRepository
{
    Task<FraudRule> CreateAsync(FraudRule entity, CancellationToken cancellationToken = default);
    Task<> GetByIdAsync(string id, CancellationToken cancellationToken = default);
    Task<List<FraudRule>> GetByStatusAsync(RuleStatus status, CancellationToken cancellationToken = default);    Task<List<FraudRule>> GetAllAsync(int page = 1, int pageSize = 50, CancellationToken cancellationToken = default);
    Task<FraudRule> UpdateAsync(FraudRule entity, CancellationToken cancellationToken = default);
    Task DeleteAsync(string id, CancellationToken cancellationToken = default);
}

public class FraudRuleRepository : IFraudRuleRepository
{
    private readonly PoliSyncDbContext _context;
    private readonly ILogger<FraudRuleRepository> _logger;

    public FraudRuleRepository(PoliSyncDbContext context, ILogger<FraudRuleRepository> logger)
    {
        _context = context;
        _logger = logger;
    }

    public async Task<FraudRule> CreateAsync(FraudRule entity, CancellationToken cancellationToken = default)
    {
        var now = Timestamp.FromDateTime(DateTime.UtcNow);
        
        if (string.IsNullOrEmpty(entity.RuleId))
        {
            entity.RuleId = Guid.NewGuid().ToString();
        }

        _context.FraudRules.Add(entity);
        await _context.SaveChangesAsync(cancellationToken);
        
        _logger.LogInformation("Created FraudRule {Id}", entity.RuleId);
        return entity;
    }

    public async Task<> GetByIdAsync(string id, CancellationToken cancellationToken = default)
    {
        return await _context.FraudRules
            .FirstOrDefaultAsync(e => e.RuleId == id, cancellationToken);
    }

    public async Task<List<FraudRule>> GetByStatusAsync(RuleStatus status, CancellationToken cancellationToken = default)
    {
        return await _context.FraudRules
            .Where(e => e.Status == status)
            .ToListAsync(cancellationToken);
    }
    public async Task<List<FraudRule>> GetAllAsync(int page = 1, int pageSize = 50, CancellationToken cancellationToken = default)
    {
        return await _context.FraudRules
            .OrderByDescending(e => e.RuleId)
            .Skip((page - 1) * pageSize)
            .Take(pageSize)
            .ToListAsync(cancellationToken);
    }

    public async Task<FraudRule> UpdateAsync(FraudRule entity, CancellationToken cancellationToken = default)
    {
        _context.FraudRules.Update(entity);
        await _context.SaveChangesAsync(cancellationToken);
        
        _logger.LogInformation("Updated FraudRule {Id}", entity.RuleId);
        return entity;
    }

    public async Task DeleteAsync(string id, CancellationToken cancellationToken = default)
    {
        var entity = await GetByIdAsync(id, cancellationToken);
        if (entity != null)
        {
            _context.FraudRules.Remove(entity);
            await _context.SaveChangesAsync(cancellationToken);
            _logger.LogInformation("Deleted FraudRule {Id}", id);
        }
    }
}
