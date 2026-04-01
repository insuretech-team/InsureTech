using System.Collections.Concurrent;
using Google.Protobuf;
using Google.Protobuf.WellKnownTypes;
using Insuretech.Quoting.Entity.V1;

namespace PoliSync.Quoting.Services;

public class InMemoryQuoteRepository : IQuoteRepository
{
    private readonly ConcurrentDictionary<string, Quote> _quotes = new();
    private readonly ILogger<InMemoryQuoteRepository> _logger;

    public InMemoryQuoteRepository(ILogger<InMemoryQuoteRepository> logger)
    {
        _logger = logger;
    }

    public Task<Quote?> GetByIdAsync(string id, CancellationToken cancellationToken = default)
    {
        _quotes.TryGetValue(id, out var quote);
        return Task.FromResult(quote);
    }

    public Task<Quote?> GetByNumberAsync(string number, CancellationToken cancellationToken = default)
    {
        var quote = _quotes.Values.FirstOrDefault(q => 
            q.QuoteNumber.Equals(number, StringComparison.OrdinalIgnoreCase));
        return Task.FromResult(quote);
    }

    public Task<IEnumerable<Quote>> GetByCustomerIdAsync(string customerId, CancellationToken cancellationToken = default)
    {
        var quotes = _quotes.Values
            .Where(q => q.CustomerId == customerId)
            .OrderByDescending(q => q.CreatedAt)
            .AsEnumerable();
        return Task.FromResult(quotes);
    }

    public Task<IEnumerable<Quote>> GetAllAsync(CancellationToken cancellationToken = default)
    {
        return Task.FromResult(_quotes.Values.AsEnumerable());
    }

    public Task<IEnumerable<Quote>> GetByFilterAsync(
        string? customerId, 
        string? productId, 
        QuoteStatus? status, 
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
        
        return Task.FromResult<IEnumerable<Quote>>(query.OrderByDescending(q => q.CreatedAt));
    }

    public Task<Quote> CreateAsync(Quote quote, CancellationToken cancellationToken = default)
    {
        if (string.IsNullOrEmpty(quote.QuoteId))
        {
            quote.QuoteId = Guid.NewGuid().ToString();
        }
        
        quote.CreatedAt = Timestamp.FromDateTime(DateTime.UtcNow);
        quote.UpdatedAt = Timestamp.FromDateTime(DateTime.UtcNow);
        
        _quotes[quote.QuoteId] = quote;
        _logger.LogInformation("Created quote: {QuoteId} - {QuoteNumber}", 
            quote.QuoteId, quote.QuoteNumber);
        
        return Task.FromResult(quote);
    }

    public Task<Quote?> UpdateAsync(Quote quote, CancellationToken cancellationToken = default)
    {
        if (!_quotes.ContainsKey(quote.QuoteId))
        {
            return Task.FromResult<Quote?>(null);
        }

        quote.UpdatedAt = Timestamp.FromDateTime(DateTime.UtcNow);
        _quotes[quote.QuoteId] = quote;
        
        _logger.LogInformation("Updated quote: {QuoteId}", quote.QuoteId);
        
        return Task.FromResult<Quote?>(quote);
    }

    public Task<bool> DeleteAsync(string id, bool permanent = false, CancellationToken cancellationToken = default)
    {
        if (permanent)
        {
            var result = _quotes.TryRemove(id, out _);
            if (result)
            {
                _logger.LogInformation("Permanently deleted quote: {QuoteId}", id);
            }
            return Task.FromResult(result);
        }
        else
        {
            if (_quotes.TryGetValue(id, out var quote))
            {
                quote.DeletedAt = Timestamp.FromDateTime(DateTime.UtcNow);
                quote.Status = QuoteStatus.Expired;
                _logger.LogInformation("Soft deleted quote: {QuoteId}", id);
                return Task.FromResult(true);
            }
            return Task.FromResult(false);
        }
    }
}
