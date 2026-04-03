using Grpc.Core;
using InsuranceEngine.Grpc.Clients;
using PolicyEntity = Insuretech.Policy.Entity.V1.Policy;
using Insuretech.Policy.Services.V1;

namespace InsuranceEngine.Endorsements.Infrastructure;

/// <summary>
/// Implementation of IEndorsementDataGateway using gRPC calls to the Go backend's PolicyService.
/// Decoupled from the main Policy module.
/// </summary>
public sealed class GoEndorsementDataGateway : IEndorsementDataGateway
{
    private readonly InsuranceServiceClient _client;

    public GoEndorsementDataGateway(InsuranceServiceClient client)
    {
        _client = client;
    }

    public async Task<PolicyEntity?> GetPolicyAsync(string policyId, CancellationToken ct = default)
    {
        try
        {
            var response = await _client.Policies.GetPolicyAsync(
                new GetPolicyRequest { PolicyId = policyId }, 
                _client.BuildCallOptions(ct));
            
            return response.Policy;
        }
        catch (RpcException ex) when (ex.StatusCode == StatusCode.NotFound)
        {
            return null;
        }
    }

    public async Task<PolicyEntity> UpdatePolicyAsync(PolicyEntity policy, CancellationToken ct = default)
    {
        var response = await _client.Policies.UpdatePolicyAsync(
            new UpdatePolicyRequest { Policy = policy }, 
            _client.BuildCallOptions(ct));
            
        return response.Policy;
    }
}
