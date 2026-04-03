using Grpc.Core;
using InsuranceEngine.Grpc.Clients;
using Insuretech.Quoting.Entity.V1;
using Insuretech.Quoting.Services.V1;
using InsuranceEngine.Quoting;

namespace InsuranceEngine.Quoting.Infrastructure;

/// <summary>
/// Implementation of IQuotingDataGateway using gRPC calls to the Go backend's QuotingService.
/// </summary>
public sealed class GoQuotingDataGateway : IQuotingDataGateway
{
    private readonly InsuranceServiceClient _client;

    public GoQuotingDataGateway(InsuranceServiceClient client)
    {
        _client = client;
    }

    public async Task<Quote> GenerateQuoteAsync(GenerateQuoteRequest request, CancellationToken cancellationToken = default)
    {
        var response = await _client.Quoting.GenerateQuoteAsync(request, _client.BuildCallOptions(cancellationToken));
        if (response.Error != null)
            throw new Exception($"[QuotingService] {response.Error.Code}: {response.Error.Message}");
        
        return response.Quote;
    }

    public async Task<Quote> GetQuoteAsync(string quoteId, CancellationToken cancellationToken = default)
    {
        var response = await _client.Quoting.GetQuoteAsync(new GetQuoteRequest { QuoteId = quoteId }, _client.BuildCallOptions(cancellationToken));
        if (response.Error != null)
            throw new Exception($"[QuotingService] {response.Error.Code}: {response.Error.Message}");

        return response.Quote;
    }

    public async Task<Quote> GetQuoteByNumberAsync(string quoteNumber, CancellationToken cancellationToken = default)
    {
        var response = await _client.Quoting.GetQuoteByNumberAsync(new GetQuoteByNumberRequest { QuoteNumber = quoteNumber }, _client.BuildCallOptions(cancellationToken));
        if (response.Error != null)
            throw new Exception($"[QuotingService] {response.Error.Code}: {response.Error.Message}");

        return response.Quote;
    }

    public async Task<ListQuotesResponse> ListQuotesAsync(ListQuotesRequest request, CancellationToken cancellationToken = default)
    {
        return await _client.Quoting.ListQuotesAsync(request, _client.BuildCallOptions(cancellationToken));
    }

    public async Task<Quote> ReviseQuoteAsync(ReviseQuoteRequest request, CancellationToken cancellationToken = default)
    {
        var response = await _client.Quoting.ReviseQuoteAsync(request, _client.BuildCallOptions(cancellationToken));
        if (response.Error != null)
            throw new Exception($"[QuotingService] {response.Error.Code}: {response.Error.Message}");

        return response.Quote;
    }

    public async Task<ConvertQuoteToPolicyResponse> ConvertQuoteToPolicyAsync(ConvertQuoteToPolicyRequest request, CancellationToken cancellationToken = default)
    {
        return await _client.Quoting.ConvertQuoteToPolicyAsync(request, _client.BuildCallOptions(cancellationToken));
    }
}
