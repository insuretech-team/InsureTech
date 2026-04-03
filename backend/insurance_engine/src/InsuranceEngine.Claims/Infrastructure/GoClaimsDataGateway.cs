using Grpc.Core;
using Insuretech.Claims.Services.V1;
using InsuranceEngine.Grpc.Clients;

namespace InsuranceEngine.Claims.Infrastructure;

public sealed class GoClaimsDataGateway : IClaimDataGateway
{
    private readonly InsuranceServiceClient _client;

    public GoClaimsDataGateway(InsuranceServiceClient client)
    {
        _client = client;
    }

    public async Task<GetClaimResponse> GetClaimAsync(string claimId, CancellationToken ct = default)
    {
        try
        {
            var response = await _client.Claims.GetClaimAsync(
                new GetClaimRequest { ClaimId = claimId }, 
                _client.BuildCallOptions(ct));
            
            return response;
        }
        catch (RpcException ex) when (ex.StatusCode == StatusCode.NotFound)
        {
            return new GetClaimResponse();
        }
    }

    public async Task<SubmitClaimResponse> SubmitClaimAsync(SubmitClaimRequest request, CancellationToken ct = default)
    {
        return await _client.Claims.SubmitClaimAsync(request, _client.BuildCallOptions(ct));
    }

    public async Task<ApproveClaimResponse> ApproveClaimAsync(string claimId, string notes, CancellationToken ct = default)
    {
        return await _client.Claims.ApproveClaimAsync(
            new ApproveClaimRequest { ClaimId = claimId, Notes = notes },
            _client.BuildCallOptions(ct));
    }

    public async Task<RejectClaimResponse> RejectClaimAsync(string claimId, string reason, CancellationToken ct = default)
    {
        return await _client.Claims.RejectClaimAsync(
            new RejectClaimRequest { ClaimId = claimId, Reason = reason },
            _client.BuildCallOptions(ct));
    }

    public async Task<SettleClaimResponse> SettleClaimAsync(string claimId, string paymentMethod, CancellationToken ct = default)
    {
        return await _client.Claims.SettleClaimAsync(
            new SettleClaimRequest { ClaimId = claimId, PaymentMethod = paymentMethod },
            _client.BuildCallOptions(ct));
    }

    public async Task<UploadDocumentResponse> UploadDocumentAsync(string claimId, string fileName, string documentType, string documentUrl, CancellationToken ct = default)
    {
        return await _client.Claims.UploadDocumentAsync(
            new UploadDocumentRequest 
            { 
                ClaimId = claimId, 
                FileName = fileName, 
                DocumentType = documentType
            },
            _client.BuildCallOptions(ct));
    }

    public async Task<RequestMoreDocumentsResponse> RequestMoreDocumentsAsync(string claimId, string message, List<string> requiredDocumentTypes, CancellationToken ct = default)
    {
        var request = new RequestMoreDocumentsRequest
        {
            ClaimId = claimId,
            Message = message
        };
        request.RequiredDocumentTypes.AddRange(requiredDocumentTypes);
        
        return await _client.Claims.RequestMoreDocumentsAsync(request, _client.BuildCallOptions(ct));
    }

    public async Task<DisputeClaimResponse> DisputeClaimAsync(string claimId, string disputeReason, string customerId, CancellationToken ct = default)
    {
        return await _client.Claims.DisputeClaimAsync(
            new DisputeClaimRequest 
            { 
                ClaimId = claimId, 
                DisputeReason = disputeReason,
                CustomerId = customerId
            },
            _client.BuildCallOptions(ct));
    }
}
