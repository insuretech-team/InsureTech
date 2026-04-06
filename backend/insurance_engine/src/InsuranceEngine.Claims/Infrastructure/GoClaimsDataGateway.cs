using InsuranceEngine.Grpc.Clients;
using Insuretech.Claims.Entity.V1;
using Insuretech.Claims.Services.V1;
using Microsoft.Extensions.Logging;

namespace InsuranceEngine.Claims.Infrastructure;

/// <summary>
/// Implementation of IClaimDataGateway that proxies all calls to the Go backend via gRPC.
/// Used when the Insurance Engine is configured to use the Go persistence layer.
/// </summary>
public class GoClaimsDataGateway : IClaimDataGateway
{
    private readonly InsuranceServiceClient _grpcClient;
    private readonly ILogger<GoClaimsDataGateway> _logger;

    public GoClaimsDataGateway(InsuranceServiceClient grpcClient, ILogger<GoClaimsDataGateway> logger)
    {
        _grpcClient = grpcClient;
        _logger = logger;
    }

    public async Task<GetClaimResponse> GetClaimAsync(string claimId, CancellationToken ct = default)
    {
        _logger.LogInformation("GoProxy: Getting claim {ClaimId}", claimId);
        return await _grpcClient.Claims.GetClaimAsync(new GetClaimRequest { ClaimId = claimId }, _grpcClient.BuildCallOptions(ct));
    }

    public async Task<SubmitClaimResponse> SubmitClaimAsync(SubmitClaimRequest request, CancellationToken ct = default)
    {
        _logger.LogInformation("GoProxy: Submitting claim for Policy {PolicyId}", request.PolicyId);
        return await _grpcClient.Claims.SubmitClaimAsync(request, _grpcClient.BuildCallOptions(ct));
    }

    public async Task<ApproveClaimResponse> ApproveClaimAsync(string claimId, string notes, CancellationToken ct = default)
    {
        _logger.LogInformation("GoProxy: Approving claim {ClaimId}", claimId);
        var request = new ApproveClaimRequest { ClaimId = claimId, Notes = notes };
        return await _grpcClient.Claims.ApproveClaimAsync(request, _grpcClient.BuildCallOptions(ct));
    }

    public async Task<RejectClaimResponse> RejectClaimAsync(string claimId, string reason, CancellationToken ct = default)
    {
        _logger.LogInformation("GoProxy: Rejecting claim {ClaimId}", claimId);
        var request = new RejectClaimRequest { ClaimId = claimId, Reason = reason };
        return await _grpcClient.Claims.RejectClaimAsync(request, _grpcClient.BuildCallOptions(ct));
    }

    public async Task<SettleClaimResponse> SettleClaimAsync(string claimId, string paymentMethod, CancellationToken ct = default)
    {
        _logger.LogInformation("GoProxy: Settling claim {ClaimId}", claimId);
        var request = new SettleClaimRequest { ClaimId = claimId, PaymentMethod = paymentMethod };
        return await _grpcClient.Claims.SettleClaimAsync(request, _grpcClient.BuildCallOptions(ct));
    }

    public async Task<UploadDocumentResponse> UploadDocumentAsync(string claimId, string fileName, string documentType, string documentUrl, CancellationToken ct = default)
    {
        _logger.LogInformation("GoProxy: Uploading document for claim {ClaimId}", claimId);
        var request = new UploadDocumentRequest 
        { 
            ClaimId = claimId, 
            FileName = fileName, 
            DocumentType = documentType, 
            DocumentUrl = documentUrl 
        };
        return await _grpcClient.Claims.UploadDocumentAsync(request, _grpcClient.BuildCallOptions(ct));
    }

    public async Task<RequestMoreDocumentsResponse> RequestMoreDocumentsAsync(string claimId, string message, List<string> requiredDocumentTypes, CancellationToken ct = default)
    {
        _logger.LogInformation("GoProxy: Requesting more documents for claim {ClaimId}", claimId);
        var request = new RequestMoreDocumentsRequest { ClaimId = claimId, Message = message };
        request.RequiredDocumentTypes.AddRange(requiredDocumentTypes);
        return await _grpcClient.Claims.RequestMoreDocumentsAsync(request, _grpcClient.BuildCallOptions(ct));
    }

    public async Task<DisputeClaimResponse> DisputeClaimAsync(string claimId, string disputeReason, string customerId, CancellationToken ct = default)
    {
        _logger.LogInformation("GoProxy: Filing dispute for claim {ClaimId}", claimId);
        var request = new DisputeClaimRequest { ClaimId = claimId, DisputeReason = disputeReason, CustomerId = customerId };
        return await _grpcClient.Claims.DisputeClaimAsync(request, _grpcClient.BuildCallOptions(ct));
    }
}
