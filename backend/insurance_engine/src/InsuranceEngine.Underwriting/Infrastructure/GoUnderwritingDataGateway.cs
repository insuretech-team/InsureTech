using Grpc.Core;
using InsuranceEngine.Grpc.Clients;
using Insuretech.Underwriting.Services.V1;
using InsuranceEngine.Underwriting;

namespace InsuranceEngine.Underwriting.Infrastructure;

/// <summary>
/// Implementation of IUnderwritingDataGateway using gRPC calls to the Go backend.
/// </summary>
public sealed class GoUnderwritingDataGateway : IUnderwritingDataGateway
{
    private readonly InsuranceServiceClient _client;

    public GoUnderwritingDataGateway(InsuranceServiceClient client)
    {
        _client = client;
    }

    public async Task<RequestQuoteResponse> RequestQuoteAsync(RequestQuoteRequest request, CancellationToken ct = default)
    {
        return await _client.Underwriting.RequestQuoteAsync(request, _client.BuildCallOptions(ct));
    }

    public async Task<GetQuoteResponse> GetQuoteAsync(string quoteId, CancellationToken ct = default)
    {
        var request = new GetQuoteRequest { QuoteId = quoteId };
        return await _client.Underwriting.GetQuoteAsync(request, _client.BuildCallOptions(ct));
    }

    public async Task<SubmitHealthDeclarationResponse> SubmitHealthDeclarationAsync(SubmitHealthDeclarationRequest request, CancellationToken ct = default)
    {
        return await _client.Underwriting.SubmitHealthDeclarationAsync(request, _client.BuildCallOptions(ct));
    }

    public async Task<ApproveUnderwritingResponse> ApproveUnderwritingAsync(ApproveUnderwritingRequest request, CancellationToken ct = default)
    {
        return await _client.Underwriting.ApproveUnderwritingAsync(request, _client.BuildCallOptions(ct));
    }

    public async Task<RejectUnderwritingResponse> RejectUnderwritingAsync(RejectUnderwritingRequest request, CancellationToken ct = default)
    {
        return await _client.Underwriting.RejectUnderwritingAsync(request, _client.BuildCallOptions(ct));
    }

    public async Task<ConvertQuoteToPolicyResponse> ConvertQuoteToPolicyAsync(ConvertQuoteToPolicyRequest request, CancellationToken ct = default)
    {
        return await _client.Underwriting.ConvertQuoteToPolicyAsync(request, _client.BuildCallOptions(ct));
    }
}
