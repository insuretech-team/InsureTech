using Grpc.Core;
using InsuranceEngine.Grpc.Clients;
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

    public async Task<GetPolicyResponse> GetPolicyAsync(string policyId, CancellationToken ct = default)
    {
        try
        {
            var response = await _client.Policies.GetPolicyAsync(
                new GetPolicyRequest { PolicyId = policyId }, 
                _client.BuildCallOptions(ct));
            
            return response;
        }
        catch (RpcException ex) when (ex.StatusCode == StatusCode.NotFound)
        {
            return new GetPolicyResponse();
        }
    }

    public async Task<UpdatePolicyResponse> UpdatePolicyAsync(string policyId, List<Insuretech.Policy.Entity.V1.Nominee>? nominees, CancellationToken ct = default)
    {
        var request = new UpdatePolicyRequest
        {
            PolicyId = policyId
        };
        
        if (nominees != null)
        {
            request.Nominees.AddRange(nominees);
        }
        
        var response = await _client.Policies.UpdatePolicyAsync(
            request, 
            _client.BuildCallOptions(ct));
            
        return response;
    }
}
