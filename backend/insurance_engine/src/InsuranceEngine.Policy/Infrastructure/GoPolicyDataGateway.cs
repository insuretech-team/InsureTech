using Grpc.Core;
using InsuranceEngine.Grpc.Clients;
using PolicyEntity = Insuretech.Policy.Entity.V1.Policy;
using Insuretech.Policy.Services.V1;

namespace InsuranceEngine.Policy.Infrastructure;

/// <summary>
/// Implementation of IPolicyDataGateway using gRPC calls to the Go backend's PolicyService.
/// </summary>
public sealed class GoPolicyDataGateway : IPolicyDataGateway
{
    private readonly InsuranceServiceClient _client;

    public GoPolicyDataGateway(InsuranceServiceClient client)
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

    public async Task<ListUserPoliciesResponse> ListUserPoliciesAsync(ListUserPoliciesRequest request, CancellationToken ct = default)
    {
        return await _client.Policies.ListUserPoliciesAsync(request, _client.BuildCallOptions(ct));
    }

    public async Task<CreatePolicyResponse> CreatePolicyAsync(CreatePolicyRequest request, CancellationToken ct = default)
    {
        return await _client.Policies.CreatePolicyAsync(request, _client.BuildCallOptions(ct));
    }

    public async Task<UpdatePolicyResponse> UpdatePolicyAsync(string policyId, List<Insuretech.Policy.Entity.V1.Nominee>? nominees, string? address, CancellationToken ct = default)
    {
        var request = new UpdatePolicyRequest
        {
            PolicyId = policyId
        };
        
        if (nominees != null)
        {
            request.Nominees.AddRange(nominees);
        }
        
        if (address != null)
        {
            request.Address = address;
        }
        
        return await _client.Policies.UpdatePolicyAsync(request, _client.BuildCallOptions(ct));
    }

    public async Task<CancelPolicyResponse> CancelPolicyAsync(CancelPolicyRequest request, CancellationToken ct = default)
    {
        return await _client.Policies.CancelPolicyAsync(request, _client.BuildCallOptions(ct));
    }

    public async Task<RenewPolicyTenureResponse> RenewPolicyAsync(RenewPolicyTenureRequest request, CancellationToken ct = default)
    {
        return await _client.Policies.RenewPolicyAsync(request, _client.BuildCallOptions(ct));
    }

    public async Task<GeneratePolicyDocumentResponse> GeneratePolicyDocumentAsync(string policyId, CancellationToken ct = default)
    {
        return await _client.Policies.GeneratePolicyDocumentAsync(
            new GeneratePolicyDocumentRequest { PolicyId = policyId }, 
            _client.BuildCallOptions(ct));
    }

    public async Task<IssuePolicyResponse> IssuePolicyAsync(IssuePolicyRequest request, CancellationToken ct = default)
    {
        return await _client.Policies.IssuePolicyAsync(request, _client.BuildCallOptions(ct));
    }

    public async Task<ApproveCancellationResponse> ApproveCancellationAsync(ApproveCancellationRequest request, CancellationToken ct = default)
    {
        return await _client.Policies.ApproveCancellationAsync(request, _client.BuildCallOptions(ct));
    }
}
