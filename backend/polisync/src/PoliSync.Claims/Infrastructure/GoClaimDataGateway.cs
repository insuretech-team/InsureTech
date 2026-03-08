using Grpc.Core;
using Insuretech.Insurance.Services.V1;
using PoliSync.Infrastructure.Clients;
using ClaimEntity = Insuretech.Claims.Entity.V1.Claim;

namespace PoliSync.Claims.Infrastructure;

public sealed class GoClaimDataGateway : IClaimDataGateway
{
    private readonly InsuranceServiceClient _insuranceClient;

    public GoClaimDataGateway(InsuranceServiceClient insuranceClient)
    {
        _insuranceClient = insuranceClient;
    }

    public async Task<ClaimEntity> CreateClaimAsync(ClaimEntity claim, CancellationToken cancellationToken = default)
    {
        var response = await _insuranceClient.Client.CreateClaimAsync(
            new CreateClaimRequest { Claim = claim },
            cancellationToken: cancellationToken);

        return response.Claim;
    }

    public async Task<ClaimEntity?> GetClaimAsync(string claimId, CancellationToken cancellationToken = default)
    {
        try
        {
            var response = await _insuranceClient.Client.GetClaimAsync(
                new GetClaimRequest { ClaimId = claimId },
                cancellationToken: cancellationToken);

            return response.Claim;
        }
        catch (RpcException ex) when (ex.StatusCode == StatusCode.NotFound)
        {
            return null;
        }
    }

    public async Task<ClaimEntity> UpdateClaimAsync(ClaimEntity claim, CancellationToken cancellationToken = default)
    {
        var response = await _insuranceClient.Client.UpdateClaimAsync(
            new UpdateClaimRequest { Claim = claim },
            cancellationToken: cancellationToken);

        return response.Claim;
    }

    public async Task<IReadOnlyList<ClaimEntity>> ListClaimsAsync(string customerId, string policyId, int page, int pageSize, CancellationToken cancellationToken = default)
    {
        var response = await _insuranceClient.Client.ListClaimsAsync(new ListClaimsRequest
        {
            CustomerId = customerId,
            PolicyId = policyId,
            Page = page,
            PageSize = pageSize
        }, cancellationToken: cancellationToken);

        return response.Claims;
    }
}
