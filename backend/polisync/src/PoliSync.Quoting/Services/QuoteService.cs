using System.Collections.Concurrent;
using System.Text.Json;
using Google.Protobuf;
using Google.Protobuf.WellKnownTypes;
using Insuretech.Quoting.Entity.V1;
using Insuretech.Common.V1;

namespace PoliSync.Quoting.Services;

public class QuoteService : IQuoteService
{
    private readonly IQuoteRepository _repository;
    private readonly IPricingEngine _pricingEngine;
    private readonly IQuoteNumberGenerator _numberGenerator;
    private readonly ILogger<QuoteService> _logger;

    public QuoteService(
        IQuoteRepository repository,
        IPricingEngine pricingEngine,
        IQuoteNumberGenerator numberGenerator,
        ILogger<QuoteService> logger)
    {
        _repository = repository;
        _pricingEngine = pricingEngine;
        _numberGenerator = numberGenerator;
        _logger = logger;
    }

    public async Task<Quote> GenerateQuoteAsync(
        string productId,
        string customerId,
        QuoteParameters parameters,
        string? agentId = null,
        int validityDays = 30,
        CancellationToken cancellationToken = default)
    {
        _logger.LogInformation(
            "Generating quote for product {ProductId}, customer {CustomerId}",
            productId, customerId);

        // Calculate premium
        var (calculation, coverages, discounts) = await _pricingEngine.CalculatePremiumAsync(
            productId, parameters, cancellationToken);

        // Create quote
        var quote = new Quote
        {
            QuoteId = Guid.NewGuid().ToString(),
            QuoteNumber = _numberGenerator.GenerateQuoteNumber(),
            ProductId = productId,
            CustomerId = customerId,
            AgentId = agentId ?? string.Empty,
            Status = QuoteStatus.Generated,
            ParametersJson = JsonSerializer.Serialize(parameters),
            PremiumCalculationJson = JsonSerializer.Serialize(calculation),
            CoveragesJson = JsonSerializer.Serialize(coverages),
            DiscountsJson = JsonSerializer.Serialize(discounts),
            TotalPremium = calculation.TotalPremium,
            ValidFrom = Timestamp.FromDateTime(DateTime.UtcNow),
            ValidUntil = Timestamp.FromDateTime(DateTime.UtcNow.AddDays(validityDays)),
            RevisionNumber = 1,
            CreatedAt = Timestamp.FromDateTime(DateTime.UtcNow),
            UpdatedAt = Timestamp.FromDateTime(DateTime.UtcNow)
        };

        var created = await _repository.CreateAsync(quote, cancellationToken);
        
        _logger.LogInformation(
            "Quote generated: {QuoteNumber} with premium {Premium}",
            created.QuoteNumber, created.TotalPremium);

        return created;
    }

    public Task<Quote?> GetQuoteAsync(string quoteId, CancellationToken cancellationToken = default)
    {
        return _repository.GetByIdAsync(quoteId, cancellationToken);
    }

    public Task<Quote?> GetQuoteByNumberAsync(string quoteNumber, CancellationToken cancellationToken = default)
    {
        return _repository.GetByNumberAsync(quoteNumber, cancellationToken);
    }

    public Task<IEnumerable<Quote>> GetQuotesByCustomerAsync(string customerId, CancellationToken cancellationToken = default)
    {
        return _repository.GetByCustomerIdAsync(customerId, cancellationToken);
    }

    public Task<IEnumerable<Quote>> ListQuotesAsync(
        string? customerId, 
        string? productId, 
        QuoteStatus? status, 
        CancellationToken cancellationToken = default)
    {
        return _repository.GetByFilterAsync(customerId, productId, status, cancellationToken);
    }

    public async Task<Quote> ReviseQuoteAsync(
        string quoteId,
        QuoteParameters newParameters,
        string reason,
        int validityDays = 30,
        CancellationToken cancellationToken = default)
    {
        _logger.LogInformation("Revising quote: {QuoteId}", quoteId);

        var parentQuote = await _repository.GetByIdAsync(quoteId, cancellationToken);
        if (parentQuote == null)
        {
            throw new InvalidOperationException($"Quote {quoteId} not found");
        }

        // Calculate new premium
        var (calculation, coverages, discounts) = await _pricingEngine.CalculatePremiumAsync(
            parentQuote.ProductId, newParameters, cancellationToken);

        // Create revised quote
        var revisedQuote = new Quote
        {
            QuoteId = Guid.NewGuid().ToString(),
            QuoteNumber = _numberGenerator.GenerateQuoteNumber(),
            ProductId = parentQuote.ProductId,
            CustomerId = parentQuote.CustomerId,
            AgentId = parentQuote.AgentId,
            Status = QuoteStatus.Generated,
            ParametersJson = JsonSerializer.Serialize(newParameters),
            PremiumCalculationJson = JsonSerializer.Serialize(calculation),
            CoveragesJson = JsonSerializer.Serialize(coverages),
            DiscountsJson = JsonSerializer.Serialize(discounts),
            TotalPremium = calculation.TotalPremium,
            ValidFrom = Timestamp.FromDateTime(DateTime.UtcNow),
            ValidUntil = Timestamp.FromDateTime(DateTime.UtcNow.AddDays(validityDays)),
            RevisionNumber = parentQuote.RevisionNumber + 1,
            ParentQuoteId = parentQuote.QuoteId,
            RevisionReason = reason,
            CreatedAt = Timestamp.FromDateTime(DateTime.UtcNow),
            UpdatedAt = Timestamp.FromDateTime(DateTime.UtcNow)
        };

        // Update parent quote status
        parentQuote.Status = QuoteStatus.Expired;
        await _repository.UpdateAsync(parentQuote, cancellationToken);

        var created = await _repository.CreateAsync(revisedQuote, cancellationToken);
        
        _logger.LogInformation(
            "Quote revised: {NewQuoteNumber} (parent: {ParentQuoteNumber})",
            created.QuoteNumber, parentQuote.QuoteNumber);

        return created;
    }

    public async Task<bool> ConvertQuoteToPolicyAsync(
        string quoteId, 
        string policyId, 
        string? convertedBy = null,
        CancellationToken cancellationToken = default)
    {
        _logger.LogInformation("Converting quote {QuoteId} to policy {PolicyId}", quoteId, policyId);

        var quote = await _repository.GetByIdAsync(quoteId, cancellationToken);
        if (quote == null)
        {
            return false;
        }

        if (quote.Status != QuoteStatus.Generated && quote.Status != QuoteStatus.Sent && quote.Status != QuoteStatus.Viewed)
        {
            _logger.LogWarning("Quote {QuoteId} cannot be converted. Status: {Status}", quoteId, quote.Status);
            return false;
        }

        quote.Status = QuoteStatus.Converted;
        quote.ConvertedPolicyId = policyId;
        quote.ConvertedAt = Timestamp.FromDateTime(DateTime.UtcNow);
        quote.UpdatedAt = Timestamp.FromDateTime(DateTime.UtcNow);

        await _repository.UpdateAsync(quote, cancellationToken);
        
        _logger.LogInformation("Quote {QuoteNumber} converted to policy {PolicyId}", 
            quote.QuoteNumber, policyId);

        return true;
    }

    public async Task<bool> ExpireQuoteAsync(string quoteId, CancellationToken cancellationToken = default)
    {
        _logger.LogInformation("Expiring quote: {QuoteId}", quoteId);

        var quote = await _repository.GetByIdAsync(quoteId, cancellationToken);
        if (quote == null)
        {
            return false;
        }

        quote.Status = QuoteStatus.Expired;
        quote.UpdatedAt = Timestamp.FromDateTime(DateTime.UtcNow);

        await _repository.UpdateAsync(quote, cancellationToken);
        return true;
    }

    public Task<bool> DeleteQuoteAsync(string quoteId, bool permanent = false, CancellationToken cancellationToken = default)
    {
        return _repository.DeleteAsync(quoteId, permanent, cancellationToken);
    }

    public async Task<IEnumerable<Quote>> CompareQuotesAsync(List<string> quoteIds, CancellationToken cancellationToken = default)
    {
        var quotes = new List<Quote>();
        foreach (var id in quoteIds)
        {
            var quote = await _repository.GetByIdAsync(id, cancellationToken);
            if (quote != null)
            {
                quotes.Add(quote);
            }
        }
        return quotes;
    }

    public Task<(PremiumCalculation Calculation, List<Coverage> Coverages, List<Discount> Discounts)> CalculatePremiumAsync(
        string productId,
        QuoteParameters parameters,
        CancellationToken cancellationToken = default)
    {
        return _pricingEngine.CalculatePremiumAsync(productId, parameters, cancellationToken);
    }
}
