using Insuretech.Claims.Entity.V1;
using Insuretech.Claims.Services.V1;

namespace InsuranceEngine.Claims;

public interface IClaimDataGateway
{
    Task<GetClaimResponse> GetClaimAsync(string claimId, CancellationToken ct = default);
    Task<SubmitClaimResponse> SubmitClaimAsync(SubmitClaimRequest request, CancellationToken ct = default);
    Task<ApproveClaimResponse> ApproveClaimAsync(string claimId, string notes, CancellationToken ct = default);
    Task<RejectClaimResponse> RejectClaimAsync(string claimId, string reason, CancellationToken ct = default);
    Task<SettleClaimResponse> SettleClaimAsync(string claimId, string paymentMethod, CancellationToken ct = default);
    Task<UploadDocumentResponse> UploadDocumentAsync(string claimId, string fileName, string documentType, string documentUrl, CancellationToken ct = default);
    Task<RequestMoreDocumentsResponse> RequestMoreDocumentsAsync(string claimId, string message, List<string> requiredDocumentTypes, CancellationToken ct = default);
    Task<DisputeClaimResponse> DisputeClaimAsync(string claimId, string disputeReason, string customerId, CancellationToken ct = default);
}
