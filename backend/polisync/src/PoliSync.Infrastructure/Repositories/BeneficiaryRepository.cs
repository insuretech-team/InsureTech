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

public interface IBeneficiaryRepository
{
    Task<Beneficiary> CreateAsync(Beneficiary entity, CancellationToken cancellationToken = default);
    Task<> GetByIdAsync(string id, CancellationToken cancellationToken = default);
    Task<List<Beneficiary>> GetByStatusAsync(BeneficiaryStatus status, CancellationToken cancellationToken = default);    Task<List<Beneficiary>> GetAllAsync(int page = 1, int pageSize = 50, CancellationToken cancellationToken = default);
    Task<Beneficiary> UpdateAsync(Beneficiary entity, CancellationToken cancellationToken = default);
    Task DeleteAsync(string id, CancellationToken cancellationToken = default);
}

public class BeneficiaryRepository : IBeneficiaryRepository
{
    private readonly PoliSyncDbContext _context;
    private readonly ILogger<BeneficiaryRepository> _logger;

    public BeneficiaryRepository(PoliSyncDbContext context, ILogger<BeneficiaryRepository> logger)
    {
        _context = context;
        _logger = logger;
    }

    public async Task<Beneficiary> CreateAsync(Beneficiary entity, CancellationToken cancellationToken = default)
    {
        var now = Timestamp.FromDateTime(DateTime.UtcNow);
        
        if (string.IsNullOrEmpty(entity.BeneficiaryId))
        {
            entity.BeneficiaryId = Guid.NewGuid().ToString();
        }

        _context.Beneficiarys.Add(entity);
        await _context.SaveChangesAsync(cancellationToken);
        
        _logger.LogInformation("Created Beneficiary {Id}", entity.BeneficiaryId);
        return entity;
    }

    public async Task<> GetByIdAsync(string id, CancellationToken cancellationToken = default)
    {
        return await _context.Beneficiarys
            .FirstOrDefaultAsync(e => e.BeneficiaryId == id, cancellationToken);
    }

    public async Task<List<Beneficiary>> GetByStatusAsync(BeneficiaryStatus status, CancellationToken cancellationToken = default)
    {
        return await _context.Beneficiarys
            .Where(e => e.Status == status)
            .ToListAsync(cancellationToken);
    }
    public async Task<List<Beneficiary>> GetAllAsync(int page = 1, int pageSize = 50, CancellationToken cancellationToken = default)
    {
        return await _context.Beneficiarys
            .OrderByDescending(e => e.BeneficiaryId)
            .Skip((page - 1) * pageSize)
            .Take(pageSize)
            .ToListAsync(cancellationToken);
    }

    public async Task<Beneficiary> UpdateAsync(Beneficiary entity, CancellationToken cancellationToken = default)
    {
        _context.Beneficiarys.Update(entity);
        await _context.SaveChangesAsync(cancellationToken);
        
        _logger.LogInformation("Updated Beneficiary {Id}", entity.BeneficiaryId);
        return entity;
    }

    public async Task DeleteAsync(string id, CancellationToken cancellationToken = default)
    {
        var entity = await GetByIdAsync(id, cancellationToken);
        if (entity != null)
        {
            _context.Beneficiarys.Remove(entity);
            await _context.SaveChangesAsync(cancellationToken);
            _logger.LogInformation("Deleted Beneficiary {Id}", id);
        }
    }
}
