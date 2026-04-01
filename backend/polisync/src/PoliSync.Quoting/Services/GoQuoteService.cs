using Insuretech.Quoting.Entity.V1;
using Insuretech.Quoting.Services.V1;
using PoliSync.Infrastructure.Clients;

namespace PoliSync.Quoting.Services;

public sealed class GoQuoteService : IQuoteService
{
    private readonly InsuranceServiceClient _client;

    public GoQuoteService(InsuranceServiceClient client) => _client = client;

    public async Task<Quote> GenerateQuoteAsync(
        string productId,
        string customerId,
        QuoteParameters parameters,
        string? agentId = null,
        int validityDays = 30,
        CancellationToken cancellationToken = default)
    {
        var response = await _client.QuotingClient.GenerateQuoteAsync(
            new GenerateQuoteRequest
            {
                ProductId = productId,
                CustomerId = customerId,
                AgentId = agentId ?? string.Empty,
                Parameters = parameters,
                ValidityDays = validityDays
            },
            _client.BuildCallOptions(cancellationToken));

        return response.Quote;
    }

    public async Task<Quote?> GetQuoteAsync(string quoteId, CancellationToken cancellationToken = default)
    {
        var response = await _client.QuotingClient.GetQuoteAsync(
            new GetQuoteRequest { QuoteId = quoteId },
            _client.BuildCallOptions(cancellationToken));
        return response.Error is null ? response.Quote : null;
    }

    public async Task<Quote?> GetQuoteByNumberAsync(string quoteNumber, CancellationToken cancellationToken = default)
    {
        var response = await _client.QuotingClient.GetQuoteByNumberAsync(
            new GetQuoteByNumberRequest { QuoteNumber = quoteNumber },
            _client.BuildCallOptions(cancellationToken));
        return response.Error is null ? response.Quote : null;
    }

    public async Task<IEnumerable<Quote>> GetQuotesByCustomerAsync(string customerId, CancellationToken cancellationToken = default)
    {
        var response = await _client.QuotingClient.ListQuotesAsync(
            new ListQuotesRequest { CustomerId = customerId, PageSize = 200 },
            _client.BuildCallOptions(cancellationToken));
        return response.Quotes;
    }

    public async Task<IEnumerable<Quote>> ListQuotesAsync(string? customerId, string? productId, QuoteStatus? status, CancellationToken cancellationToken = default)
    {
        var response = await _client.QuotingClient.ListQuotesAsync(
            new ListQuotesRequest
            {
                CustomerId = customerId ?? string.Empty,
                ProductId = productId ?? string.Empty,
                Status = status ?? QuoteStatus.Unspecified,
                PageSize = 200
            },
            _client.BuildCallOptions(cancellationToken));
        return response.Quotes;
    }

    public async Task<Quote> ReviseQuoteAsync(
        string quoteId,
        QuoteParameters newParameters,
        string reason,
        int validityDays = 30,
        CancellationToken cancellationToken = default)
    {
        var response = await _client.QuotingClient.ReviseQuoteAsync(
            new ReviseQuoteRequest
            {
                QuoteId = quoteId,
                NewParameters = newParameters,
                RevisionReason = reason,
                ValidityDays = validityDays
            },
            _client.BuildCallOptions(cancellationToken));

        return response.Quote;
    }

    public async Task<bool> ConvertQuoteToPolicyAsync(string quoteId, string policyId, string? convertedBy = null, CancellationToken cancellationToken = default)
    {
        var response = await _client.QuotingClient.ConvertQuoteToPolicyAsync(
            new ConvertQuoteToPolicyRequest
            {
                QuoteId = quoteId,
                PolicyId = policyId,
                ConvertedBy = convertedBy ?? string.Empty
            },
            _client.BuildCallOptions(cancellationToken));
        return response.Error is null;
    }

    public async Task<bool> ExpireQuoteAsync(string quoteId, CancellationToken cancellationToken = default)
    {
        var response = await _client.QuotingClient.ExpireQuoteAsync(
            new ExpireQuoteRequest { QuoteId = quoteId },
            _client.BuildCallOptions(cancellationToken));
        return response.Success && response.Error is null;
    }

    public async Task<bool> DeleteQuoteAsync(string quoteId, bool permanent = false, CancellationToken cancellationToken = default)
    {
        var response = await _client.QuotingClient.DeleteQuoteAsync(
            new DeleteQuoteRequest { QuoteId = quoteId, Permanent = permanent },
            _client.BuildCallOptions(cancellationToken));
        return response.Success && response.Error is null;
    }

    public async Task<IEnumerable<Quote>> CompareQuotesAsync(List<string> quoteIds, CancellationToken cancellationToken = default)
    {
        var quotes = new List<Quote>(quoteIds.Count);
        foreach (var quoteId in quoteIds)
        {
            var quote = await GetQuoteAsync(quoteId, cancellationToken);
            if (quote is not null)
            {
                quotes.Add(quote);
            }
        }

        return quotes;
    }

    public async Task<(PremiumCalculation Calculation, List<Coverage> Coverages, List<Discount> Discounts)> CalculatePremiumAsync(
        string productId,
        QuoteParameters parameters,
        CancellationToken cancellationToken = default)
    {
        var response = await _client.QuotingClient.CalculatePremiumAsync(
            new CalculatePremiumRequest
            {
                ProductId = productId,
                Parameters = parameters
            },
            _client.BuildCallOptions(cancellationToken));

        return (response.Calculation, response.Coverages.ToList(), response.Discounts.ToList());
    }
}
