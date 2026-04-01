using System.Collections.Concurrent;
using Google.Protobuf;
using Google.Protobuf.WellKnownTypes;
using Insuretech.Life.Entity.V1;

namespace PoliSync.LifeInsurance.Services;

public class InMemoryLifeQuoteRepository : ILifeQuoteRepository
{
    private readonly ConcurrentDictionary<string, LifeQuote> _quotes = new();
    private readonly ILogger<InMemoryLifeQuoteRepository> _logger;

    public InMemoryLifeQuoteRepository(ILogger<InMemoryLifeQuoteRepository> logger)
    {
        _logger = logger;
    }

    public Task<LifeQuote?> GetByIdAsync(string id, CancellationToken cancellationToken = default)
    {
        _quotes.TryGetValue(id, out var quote);
        return Task.FromResult(quote);
    }

    public Task<LifeQuote?> GetByNumberAsync(string number, CancellationToken cancellationToken = default)
    {
        var quote = _quotes.Values.FirstOrDefault(q => 
            q.QuoteNumber.Equals(number, StringComparison.OrdinalIgnoreCase));
        return Task.FromResult(quote);
    }

    public Task<IEnumerable<LifeQuote>> GetByCustomerIdAsync(string customerId, CancellationToken cancellationToken = default)
    {
        var quotes = _quotes.Values
            .Where(q => q.CustomerId == customerId)
            .OrderByDescending(q => q.CreatedAt)
            .AsEnumerable();
        return Task.FromResult(quotes);
    }

    public Task<IEnumerable<LifeQuote>> GetByFilterAsync(
        string? customerId, 
        string? productId, 
        LifeQuoteStatus? status, 
        CancellationToken cancellationToken = default)
    {
        var query = _quotes.Values.AsEnumerable();
        
        if (!string.IsNullOrEmpty(customerId))
        {
            query = query.Where(q => q.CustomerId == customerId);
        }
        
        if (!string.IsNullOrEmpty(productId))
        {
            query = query.Where(q => q.ProductId == productId);
        }
        
        if (status.HasValue)
        {
            query = query.Where(q => q.Status == status.Value);
        }
        
        return Task.FromResult<IEnumerable<LifeQuote>>(query.OrderByDescending(q => q.CreatedAt));
    }

    public Task<LifeQuote> CreateAsync(LifeQuote quote, CancellationToken cancellationToken = default)
    {
        if (string.IsNullOrEmpty(quote.QuoteId))
        {
            quote.QuoteId = Guid.NewGuid().ToString();
        }
        
        quote.CreatedAt = Timestamp.FromDateTime(DateTime.UtcNow);
        quote.UpdatedAt = Timestamp.FromDateTime(DateTime.UtcNow);
        
        _quotes[quote.QuoteId] = quote;
        _logger.LogInformation("Created life quote: {QuoteId} - {QuoteNumber}", 
            quote.QuoteId, quote.QuoteNumber);
        
        return Task.FromResult(quote);
    }

    public Task<LifeQuote?> UpdateAsync(LifeQuote quote, CancellationToken cancellationToken = default)
    {
        if (!_quotes.ContainsKey(quote.QuoteId))
        {
            return Task.FromResult<LifeQuote?>(null);
        }

        quote.UpdatedAt = Timestamp.FromDateTime(DateTime.UtcNow);
        _quotes[quote.QuoteId] = quote;
        
        _logger.LogInformation("Updated life quote: {QuoteId}", quote.QuoteId);
        
        return Task.FromResult<LifeQuote?>(quote);
    }
}
