using Grpc.Core;
using Insuretech.Claims.Services.V1;
using ClaimEntity = Insuretech.Claims.Entity.V1.Claim;
using InsuranceEngine.Grpc.Clients;
using InsuranceEngine.Claims;

namespace InsuranceEngine.Claims.Infrastructure;

/// <summary>
/// Implementation of IClaimDataGateway using gRPC calls to the Go backend.
/// </summary>
public sealed class GoClaimsDataGateway : IClaimDataGateway
{
    private readonly InsuranceServiceClient _client;

    public GoClaimsDataGateway(InsuranceServiceClient client)
    {
        _client = client;
    }

    public async Task<ClaimEntity?> GetClaimAsync(string claimId, CancellationToken ct = default)
    {
        try
        {
            var response = await _client.Claims.GetClaimAsync(
                new GetClaimRequest { ClaimId = claimId }, 
                _client.BuildCallOptions(ct));
            
            return response.Claim;
        }
        catch (RpcException ex) when (ex.StatusCode == StatusCode.NotFound)
        {
            return null;
        }
    }

    public async Task<IReadOnlyList<ClaimEntity>> ListClaimsAsync(int page, int pageSize, string policyId = "", CancellationToken ct = default)
    {
        var response = await _client.Claims.ListClaimsAsync(
            new ListClaimsRequest { Page = page, PageSize = pageSize, PolicyId = policyId }, 
            _client.BuildCallOptions(ct));
            
        return response.Claims.ToList();
    }

    public async Task<ClaimEntity> CreateClaimAsync(ClaimEntity claim, CancellationToken ct = default)
    {
        var response = await _client.Claims.CreateClaimAsync(
            new CreateClaimRequest { Claim = claim }, 
            _client.BuildCallOptions(ct));
            
        return response.Claim;
    }

    public async Task<ClaimEntity> UpdateClaimAsync(ClaimEntity claim, CancellationToken ct = default)
    {
        var response = await _client.Claims.UpdateClaimAsync(
            new UpdateClaimRequest { Claim = claim }, 
            _client.BuildCallOptions(ct));
            
        return response.Claim;
    }
}
