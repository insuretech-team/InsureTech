using Insuretech.Quoting.Entity.V1;
using Insuretech.Common.V1;

namespace PoliSync.Quoting.Services;

public interface IQuoteService
{
    Task<Quote> GenerateQuoteAsync(
        string productId,
        string customerId,
        QuoteParameters parameters,
        string? agentId = null,
        int validityDays = 30,
        CancellationToken cancellationToken = default);

    Task<Quote?> GetQuoteAsync(string quoteId, CancellationToken cancellationToken = default);
    Task<Quote?> GetQuoteByNumberAsync(string quoteNumber, CancellationToken cancellationToken = default);
    Task<IEnumerable<Quote>> GetQuotesByCustomerAsync(string customerId, CancellationToken cancellationToken = default);
    Task<IEnumerable<Quote>> ListQuotesAsync(string? customerId, string? productId, QuoteStatus? status, CancellationToken cancellationToken = default);
    
    Task<Quote> ReviseQuoteAsync(
        string quoteId,
        QuoteParameters newParameters,
        string reason,
        int validityDays = 30,
        CancellationToken cancellationToken = default);
    
    Task<bool> ConvertQuoteToPolicyAsync(string quoteId, string policyId, string? convertedBy = null, CancellationToken cancellationToken = default);
    Task<bool> ExpireQuoteAsync(string quoteId, CancellationToken cancellationToken = default);
    Task<bool> DeleteQuoteAsync(string quoteId, bool permanent = false, CancellationToken cancellationToken = default);
    
    Task<IEnumerable<Quote>> CompareQuotesAsync(List<string> quoteIds, CancellationToken cancellationToken = default);
    
    Task<(PremiumCalculation Calculation, List<Coverage> Coverages, List<Discount> Discounts)> CalculatePremiumAsync(
        string productId,
        QuoteParameters parameters,
        CancellationToken cancellationToken = default);
}

public interface IQuoteRepository
{
    Task<Quote?> GetByIdAsync(string id, CancellationToken cancellationToken = default);
    Task<Quote?> GetByNumberAsync(string number, CancellationToken cancellationToken = default);
    Task<IEnumerable<Quote>> GetByCustomerIdAsync(string customerId, CancellationToken cancellationToken = default);
    Task<IEnumerable<Quote>> GetAllAsync(CancellationToken cancellationToken = default);
    Task<IEnumerable<Quote>> GetByFilterAsync(string? customerId, string? productId, QuoteStatus? status, CancellationToken cancellationToken = default);
    Task<Quote> CreateAsync(Quote quote, CancellationToken cancellationToken = default);
    Task<Quote?> UpdateAsync(Quote quote, CancellationToken cancellationToken = default);
    Task<bool> DeleteAsync(string id, bool permanent = false, CancellationToken cancellationToken = default);
}

public interface IPricingEngine
{
    Task<(PremiumCalculation Calculation, List<Coverage> Coverages, List<Discount> Discounts)> CalculatePremiumAsync(
        string productId,
        QuoteParameters parameters,
        CancellationToken cancellationToken = default);
}

public interface IQuoteNumberGenerator
{
    string GenerateQuoteNumber();
}
