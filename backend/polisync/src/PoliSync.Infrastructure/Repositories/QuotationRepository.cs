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

public interface IQuotationRepository
{
    Task<Quotation> CreateAsync(Quotation entity, CancellationToken cancellationToken = default);
    Task<> GetByIdAsync(string id, CancellationToken cancellationToken = default);
    Task<> GetByNumberAsync(string quotationnumber, CancellationToken cancellationToken = default);    Task<List<Quotation>> GetByStatusAsync(QuotationStatus status, CancellationToken cancellationToken = default);    Task<List<Quotation>> GetAllAsync(int page = 1, int pageSize = 50, CancellationToken cancellationToken = default);
    Task<Quotation> UpdateAsync(Quotation entity, CancellationToken cancellationToken = default);
    Task DeleteAsync(string id, CancellationToken cancellationToken = default);
}

public class QuotationRepository : IQuotationRepository
{
    private readonly PoliSyncDbContext _context;
    private readonly ILogger<QuotationRepository> _logger;

    public QuotationRepository(PoliSyncDbContext context, ILogger<QuotationRepository> logger)
    {
        _context = context;
        _logger = logger;
    }

    public async Task<Quotation> CreateAsync(Quotation entity, CancellationToken cancellationToken = default)
    {
        var now = Timestamp.FromDateTime(DateTime.UtcNow);
        
        if (string.IsNullOrEmpty(entity.QuotationId))
        {
            entity.QuotationId = Guid.NewGuid().ToString();
        }

        _context.Quotations.Add(entity);
        await _context.SaveChangesAsync(cancellationToken);
        
        _logger.LogInformation("Created Quotation {Id}", entity.QuotationId);
        return entity;
    }

    public async Task<> GetByIdAsync(string id, CancellationToken cancellationToken = default)
    {
        return await _context.Quotations
            .FirstOrDefaultAsync(e => e.QuotationId == id, cancellationToken);
    }

    public async Task<> GetByNumberAsync(string quotationnumber, CancellationToken cancellationToken = default)
    {
        return await _context.Quotations
            .FirstOrDefaultAsync(e => e.QuotationNumber == quotationnumber, cancellationToken);
    }
    public async Task<List<Quotation>> GetByStatusAsync(QuotationStatus status, CancellationToken cancellationToken = default)
    {
        return await _context.Quotations
            .Where(e => e.Status == status)
            .ToListAsync(cancellationToken);
    }
    public async Task<List<Quotation>> GetAllAsync(int page = 1, int pageSize = 50, CancellationToken cancellationToken = default)
    {
        return await _context.Quotations
            .OrderByDescending(e => e.QuotationId)
            .Skip((page - 1) * pageSize)
            .Take(pageSize)
            .ToListAsync(cancellationToken);
    }

    public async Task<Quotation> UpdateAsync(Quotation entity, CancellationToken cancellationToken = default)
    {
        _context.Quotations.Update(entity);
        await _context.SaveChangesAsync(cancellationToken);
        
        _logger.LogInformation("Updated Quotation {Id}", entity.QuotationId);
        return entity;
    }

    public async Task DeleteAsync(string id, CancellationToken cancellationToken = default)
    {
        var entity = await GetByIdAsync(id, cancellationToken);
        if (entity != null)
        {
            _context.Quotations.Remove(entity);
            await _context.SaveChangesAsync(cancellationToken);
            _logger.LogInformation("Deleted Quotation {Id}", id);
        }
    }
}
