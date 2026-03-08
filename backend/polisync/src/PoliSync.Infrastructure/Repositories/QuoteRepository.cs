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

public interface IQuoteRepository
{
    Task<Quote> CreateAsync(Quote entity, CancellationToken cancellationToken = default);
    Task<> GetByIdAsync(string id, CancellationToken cancellationToken = default);
    Task<> GetByNumberAsync(string quotenumber, CancellationToken cancellationToken = default);    Task<List<Quote>> GetByStatusAsync(QuoteStatus status, CancellationToken cancellationToken = default);    Task<List<Quote>> GetAllAsync(int page = 1, int pageSize = 50, CancellationToken cancellationToken = default);
    Task<Quote> UpdateAsync(Quote entity, CancellationToken cancellationToken = default);
    Task DeleteAsync(string id, CancellationToken cancellationToken = default);
}

public class QuoteRepository : IQuoteRepository
{
    private readonly PoliSyncDbContext _context;
    private readonly ILogger<QuoteRepository> _logger;

    public QuoteRepository(PoliSyncDbContext context, ILogger<QuoteRepository> logger)
    {
        _context = context;
        _logger = logger;
    }

    public async Task<Quote> CreateAsync(Quote entity, CancellationToken cancellationToken = default)
    {
        var now = Timestamp.FromDateTime(DateTime.UtcNow);
        
        if (string.IsNullOrEmpty(entity.QuoteId))
        {
            entity.QuoteId = Guid.NewGuid().ToString();
        }

        _context.Quotes.Add(entity);
        await _context.SaveChangesAsync(cancellationToken);
        
        _logger.LogInformation("Created Quote {Id}", entity.QuoteId);
        return entity;
    }

    public async Task<> GetByIdAsync(string id, CancellationToken cancellationToken = default)
    {
        return await _context.Quotes
            .FirstOrDefaultAsync(e => e.QuoteId == id, cancellationToken);
    }

    public async Task<> GetByNumberAsync(string quotenumber, CancellationToken cancellationToken = default)
    {
        return await _context.Quotes
            .FirstOrDefaultAsync(e => e.QuoteNumber == quotenumber, cancellationToken);
    }
    public async Task<List<Quote>> GetByStatusAsync(QuoteStatus status, CancellationToken cancellationToken = default)
    {
        return await _context.Quotes
            .Where(e => e.Status == status)
            .ToListAsync(cancellationToken);
    }
    public async Task<List<Quote>> GetAllAsync(int page = 1, int pageSize = 50, CancellationToken cancellationToken = default)
    {
        return await _context.Quotes
            .OrderByDescending(e => e.QuoteId)
            .Skip((page - 1) * pageSize)
            .Take(pageSize)
            .ToListAsync(cancellationToken);
    }

    public async Task<Quote> UpdateAsync(Quote entity, CancellationToken cancellationToken = default)
    {
        _context.Quotes.Update(entity);
        await _context.SaveChangesAsync(cancellationToken);
        
        _logger.LogInformation("Updated Quote {Id}", entity.QuoteId);
        return entity;
    }

    public async Task DeleteAsync(string id, CancellationToken cancellationToken = default)
    {
        var entity = await GetByIdAsync(id, cancellationToken);
        if (entity != null)
        {
            _context.Quotes.Remove(entity);
            await _context.SaveChangesAsync(cancellationToken);
            _logger.LogInformation("Deleted Quote {Id}", id);
        }
    }
}
