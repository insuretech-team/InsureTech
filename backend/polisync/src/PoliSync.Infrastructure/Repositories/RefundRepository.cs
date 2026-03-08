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

public interface IRefundRepository
{
    Task<Refund> CreateAsync(Refund entity, CancellationToken cancellationToken = default);
    Task<> GetByIdAsync(string id, CancellationToken cancellationToken = default);
    Task<> GetByNumberAsync(string refundnumber, CancellationToken cancellationToken = default);    Task<List<Refund>> GetByStatusAsync(RefundStatus status, CancellationToken cancellationToken = default);    Task<List<Refund>> GetAllAsync(int page = 1, int pageSize = 50, CancellationToken cancellationToken = default);
    Task<Refund> UpdateAsync(Refund entity, CancellationToken cancellationToken = default);
    Task DeleteAsync(string id, CancellationToken cancellationToken = default);
}

public class RefundRepository : IRefundRepository
{
    private readonly PoliSyncDbContext _context;
    private readonly ILogger<RefundRepository> _logger;

    public RefundRepository(PoliSyncDbContext context, ILogger<RefundRepository> logger)
    {
        _context = context;
        _logger = logger;
    }

    public async Task<Refund> CreateAsync(Refund entity, CancellationToken cancellationToken = default)
    {
        var now = Timestamp.FromDateTime(DateTime.UtcNow);
        
        if (string.IsNullOrEmpty(entity.RefundId))
        {
            entity.RefundId = Guid.NewGuid().ToString();
        }

        _context.Refunds.Add(entity);
        await _context.SaveChangesAsync(cancellationToken);
        
        _logger.LogInformation("Created Refund {Id}", entity.RefundId);
        return entity;
    }

    public async Task<> GetByIdAsync(string id, CancellationToken cancellationToken = default)
    {
        return await _context.Refunds
            .FirstOrDefaultAsync(e => e.RefundId == id, cancellationToken);
    }

    public async Task<> GetByNumberAsync(string refundnumber, CancellationToken cancellationToken = default)
    {
        return await _context.Refunds
            .FirstOrDefaultAsync(e => e.RefundNumber == refundnumber, cancellationToken);
    }
    public async Task<List<Refund>> GetByStatusAsync(RefundStatus status, CancellationToken cancellationToken = default)
    {
        return await _context.Refunds
            .Where(e => e.Status == status)
            .ToListAsync(cancellationToken);
    }
    public async Task<List<Refund>> GetAllAsync(int page = 1, int pageSize = 50, CancellationToken cancellationToken = default)
    {
        return await _context.Refunds
            .OrderByDescending(e => e.RefundId)
            .Skip((page - 1) * pageSize)
            .Take(pageSize)
            .ToListAsync(cancellationToken);
    }

    public async Task<Refund> UpdateAsync(Refund entity, CancellationToken cancellationToken = default)
    {
        _context.Refunds.Update(entity);
        await _context.SaveChangesAsync(cancellationToken);
        
        _logger.LogInformation("Updated Refund {Id}", entity.RefundId);
        return entity;
    }

    public async Task DeleteAsync(string id, CancellationToken cancellationToken = default)
    {
        var entity = await GetByIdAsync(id, cancellationToken);
        if (entity != null)
        {
            _context.Refunds.Remove(entity);
            await _context.SaveChangesAsync(cancellationToken);
            _logger.LogInformation("Deleted Refund {Id}", id);
        }
    }
}
