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

public interface IUnderwritingDecisionRepository
{
    Task<UnderwritingDecision> CreateAsync(UnderwritingDecision entity, CancellationToken cancellationToken = default);
    Task<> GetByIdAsync(string id, CancellationToken cancellationToken = default);
    Task<List<UnderwritingDecision>> GetByStatusAsync(DecisionStatus status, CancellationToken cancellationToken = default);    Task<List<UnderwritingDecision>> GetAllAsync(int page = 1, int pageSize = 50, CancellationToken cancellationToken = default);
    Task<UnderwritingDecision> UpdateAsync(UnderwritingDecision entity, CancellationToken cancellationToken = default);
    Task DeleteAsync(string id, CancellationToken cancellationToken = default);
}

public class UnderwritingDecisionRepository : IUnderwritingDecisionRepository
{
    private readonly PoliSyncDbContext _context;
    private readonly ILogger<UnderwritingDecisionRepository> _logger;

    public UnderwritingDecisionRepository(PoliSyncDbContext context, ILogger<UnderwritingDecisionRepository> logger)
    {
        _context = context;
        _logger = logger;
    }

    public async Task<UnderwritingDecision> CreateAsync(UnderwritingDecision entity, CancellationToken cancellationToken = default)
    {
        var now = Timestamp.FromDateTime(DateTime.UtcNow);
        
        if (string.IsNullOrEmpty(entity.DecisionId))
        {
            entity.DecisionId = Guid.NewGuid().ToString();
        }

        _context.UnderwritingDecisions.Add(entity);
        await _context.SaveChangesAsync(cancellationToken);
        
        _logger.LogInformation("Created UnderwritingDecision {Id}", entity.DecisionId);
        return entity;
    }

    public async Task<> GetByIdAsync(string id, CancellationToken cancellationToken = default)
    {
        return await _context.UnderwritingDecisions
            .FirstOrDefaultAsync(e => e.DecisionId == id, cancellationToken);
    }

    public async Task<List<UnderwritingDecision>> GetByStatusAsync(DecisionStatus status, CancellationToken cancellationToken = default)
    {
        return await _context.UnderwritingDecisions
            .Where(e => e.Status == status)
            .ToListAsync(cancellationToken);
    }
    public async Task<List<UnderwritingDecision>> GetAllAsync(int page = 1, int pageSize = 50, CancellationToken cancellationToken = default)
    {
        return await _context.UnderwritingDecisions
            .OrderByDescending(e => e.DecisionId)
            .Skip((page - 1) * pageSize)
            .Take(pageSize)
            .ToListAsync(cancellationToken);
    }

    public async Task<UnderwritingDecision> UpdateAsync(UnderwritingDecision entity, CancellationToken cancellationToken = default)
    {
        _context.UnderwritingDecisions.Update(entity);
        await _context.SaveChangesAsync(cancellationToken);
        
        _logger.LogInformation("Updated UnderwritingDecision {Id}", entity.DecisionId);
        return entity;
    }

    public async Task DeleteAsync(string id, CancellationToken cancellationToken = default)
    {
        var entity = await GetByIdAsync(id, cancellationToken);
        if (entity != null)
        {
            _context.UnderwritingDecisions.Remove(entity);
            await _context.SaveChangesAsync(cancellationToken);
            _logger.LogInformation("Deleted UnderwritingDecision {Id}", id);
        }
    }
}
