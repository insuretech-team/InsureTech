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

public interface IEndorsementRepository
{
    Task<Endorsement> CreateAsync(Endorsement entity, CancellationToken cancellationToken = default);
    Task<> GetByIdAsync(string id, CancellationToken cancellationToken = default);
    Task<> GetByNumberAsync(string endorsementnumber, CancellationToken cancellationToken = default);    Task<List<Endorsement>> GetByStatusAsync(EndorsementStatus status, CancellationToken cancellationToken = default);    Task<List<Endorsement>> GetAllAsync(int page = 1, int pageSize = 50, CancellationToken cancellationToken = default);
    Task<Endorsement> UpdateAsync(Endorsement entity, CancellationToken cancellationToken = default);
    Task DeleteAsync(string id, CancellationToken cancellationToken = default);
}

public class EndorsementRepository : IEndorsementRepository
{
    private readonly PoliSyncDbContext _context;
    private readonly ILogger<EndorsementRepository> _logger;

    public EndorsementRepository(PoliSyncDbContext context, ILogger<EndorsementRepository> logger)
    {
        _context = context;
        _logger = logger;
    }

    public async Task<Endorsement> CreateAsync(Endorsement entity, CancellationToken cancellationToken = default)
    {
        var now = Timestamp.FromDateTime(DateTime.UtcNow);
        
        if (string.IsNullOrEmpty(entity.EndorsementId))
        {
            entity.EndorsementId = Guid.NewGuid().ToString();
        }

        _context.Endorsements.Add(entity);
        await _context.SaveChangesAsync(cancellationToken);
        
        _logger.LogInformation("Created Endorsement {Id}", entity.EndorsementId);
        return entity;
    }

    public async Task<> GetByIdAsync(string id, CancellationToken cancellationToken = default)
    {
        return await _context.Endorsements
            .FirstOrDefaultAsync(e => e.EndorsementId == id, cancellationToken);
    }

    public async Task<> GetByNumberAsync(string endorsementnumber, CancellationToken cancellationToken = default)
    {
        return await _context.Endorsements
            .FirstOrDefaultAsync(e => e.EndorsementNumber == endorsementnumber, cancellationToken);
    }
    public async Task<List<Endorsement>> GetByStatusAsync(EndorsementStatus status, CancellationToken cancellationToken = default)
    {
        return await _context.Endorsements
            .Where(e => e.Status == status)
            .ToListAsync(cancellationToken);
    }
    public async Task<List<Endorsement>> GetAllAsync(int page = 1, int pageSize = 50, CancellationToken cancellationToken = default)
    {
        return await _context.Endorsements
            .OrderByDescending(e => e.EndorsementId)
            .Skip((page - 1) * pageSize)
            .Take(pageSize)
            .ToListAsync(cancellationToken);
    }

    public async Task<Endorsement> UpdateAsync(Endorsement entity, CancellationToken cancellationToken = default)
    {
        _context.Endorsements.Update(entity);
        await _context.SaveChangesAsync(cancellationToken);
        
        _logger.LogInformation("Updated Endorsement {Id}", entity.EndorsementId);
        return entity;
    }

    public async Task DeleteAsync(string id, CancellationToken cancellationToken = default)
    {
        var entity = await GetByIdAsync(id, cancellationToken);
        if (entity != null)
        {
            _context.Endorsements.Remove(entity);
            await _context.SaveChangesAsync(cancellationToken);
            _logger.LogInformation("Deleted Endorsement {Id}", id);
        }
    }
}
