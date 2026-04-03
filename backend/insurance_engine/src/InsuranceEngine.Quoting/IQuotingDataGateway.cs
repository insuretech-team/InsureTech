using Insuretech.Quoting.Entity.V1;
using Insuretech.Quoting.Services.V1;

namespace InsuranceEngine.Quoting;

public interface IQuotingDataGateway
{
    Task<Quote> GenerateQuoteAsync(GenerateQuoteRequest request, CancellationToken cancellationToken = default);
    Task<Quote> GetQuoteAsync(string quoteId, CancellationToken cancellationToken = default);
    Task<Quote> GetQuoteByNumberAsync(string quoteNumber, CancellationToken cancellationToken = default);
    Task<ListQuotesResponse> ListQuotesAsync(ListQuotesRequest request, CancellationToken cancellationToken = default);
    Task<Quote> ReviseQuoteAsync(ReviseQuoteRequest request, CancellationToken cancellationToken = default);
    Task<ConvertQuoteToPolicyResponse> ConvertQuoteToPolicyAsync(ConvertQuoteToPolicyRequest request, CancellationToken cancellationToken = default);
}
