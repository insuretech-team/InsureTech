using Insuretech.Underwriting.Services.V1;

namespace InsuranceEngine.Underwriting;

public interface IUnderwritingDataGateway
{
    Task<RequestQuoteResponse> RequestQuoteAsync(RequestQuoteRequest request, CancellationToken ct = default);
    Task<GetQuoteResponse> GetQuoteAsync(string quoteId, CancellationToken ct = default);
    Task<SubmitHealthDeclarationResponse> SubmitHealthDeclarationAsync(SubmitHealthDeclarationRequest request, CancellationToken ct = default);
    Task<ApproveUnderwritingResponse> ApproveUnderwritingAsync(ApproveUnderwritingRequest request, CancellationToken ct = default);
    Task<RejectUnderwritingResponse> RejectUnderwritingAsync(RejectUnderwritingRequest request, CancellationToken ct = default);
    Task<ConvertQuoteToPolicyResponse> ConvertQuoteToPolicyAsync(ConvertQuoteToPolicyRequest request, CancellationToken ct = default);
}
