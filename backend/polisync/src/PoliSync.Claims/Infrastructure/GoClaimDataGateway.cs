using Grpc.Core;
using Insuretech.Insurance.Services.V1;
using PoliSync.Infrastructure.Clients;
using ClaimApprovalEntity = Insuretech.Claims.Entity.V1.ClaimApproval;
using ClaimEntity = Insuretech.Claims.Entity.V1.Claim;
using ClaimDocumentEntity = Insuretech.Claims.Entity.V1.ClaimDocument;

namespace PoliSync.Claims.Infrastructure;

public sealed class GoClaimDataGateway : IClaimDataGateway
{
    private readonly InsuranceServiceClient _insuranceClient;

    public GoClaimDataGateway(InsuranceServiceClient insuranceClient)
    {
        _insuranceClient = insuranceClient;
    }

    public async Task<ClaimEntity> CreateClaimAsync(ClaimEntity claim, CancellationToken ct = default)
    {
        var r = await _insuranceClient.Client.CreateClaimAsync(new CreateClaimRequest { Claim = claim }, _insuranceClient.BuildCallOptions(ct));
        return r.Claim;
    }

    public async Task<ClaimDocumentEntity> CreateClaimDocumentAsync(ClaimDocumentEntity document, CancellationToken ct = default)
    {
        var r = await _insuranceClient.Client.CreateClaimDocumentAsync(
            new CreateClaimDocumentRequest { Document = document },
            _insuranceClient.BuildCallOptions(ct));
        return r.Document;
    }

    public async Task<ClaimApprovalEntity> CreateClaimApprovalAsync(ClaimApprovalEntity approval, CancellationToken ct = default)
    {
        var r = await _insuranceClient.Client.CreateClaimApprovalAsync(
            new CreateClaimApprovalRequest { Approval = approval },
            _insuranceClient.BuildCallOptions(ct));
        return r.Approval;
    }

    public async Task<ClaimEntity?> GetClaimAsync(string claimId, CancellationToken ct = default)
    {
        try
        {
            var r = await _insuranceClient.Client.GetClaimAsync(new GetClaimRequest { ClaimId = claimId }, _insuranceClient.BuildCallOptions(ct));
            return r.Claim;
        }
        catch (RpcException ex) when (ex.StatusCode == StatusCode.NotFound) { return null; }
    }

    public async Task<ClaimEntity> UpdateClaimAsync(ClaimEntity claim, CancellationToken ct = default)
    {
        var r = await _insuranceClient.Client.UpdateClaimAsync(new UpdateClaimRequest { Claim = claim }, _insuranceClient.BuildCallOptions(ct));
        return r.Claim;
    }

    public async Task<IReadOnlyList<ClaimEntity>> ListClaimsAsync(string customerId, string policyId, int page, int pageSize, CancellationToken ct = default)
    {
        var r = await _insuranceClient.Client.ListClaimsAsync(
            new ListClaimsRequest { CustomerId = customerId, PolicyId = policyId, Page = page, PageSize = pageSize },
            _insuranceClient.BuildCallOptions(ct));
        return r.Claims;
    }
}
