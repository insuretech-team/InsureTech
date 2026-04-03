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

    public async Task<IReadOnlyList<PolicyEntity>> ListPoliciesAsync(int page, int pageSize, string customerId = "", CancellationToken ct = default)
    {
        var response = await _client.Policies.ListUserPoliciesAsync(
            new ListUserPoliciesRequest { Page = page, PageSize = pageSize, CustomerId = customerId }, 
            _client.BuildCallOptions(ct));
            
        return response.Policies.ToList();
    }

    public async Task<PolicyEntity> CreatePolicyAsync(PolicyEntity policy, CancellationToken ct = default)
    {
        var response = await _client.Policies.CreatePolicyAsync(
            new CreatePolicyRequest { Policy = policy }, 
            _client.BuildCallOptions(ct));
            
        return response.Policy;
    }

    public async Task<PolicyEntity> UpdatePolicyAsync(PolicyEntity policy, CancellationToken ct = default)
    {
        var response = await _client.Policies.UpdatePolicyAsync(
            new UpdatePolicyRequest { Policy = policy }, 
            _client.BuildCallOptions(ct));
            
        return response.Policy;
    }

    public async Task<CancelPolicyResponse> CancelPolicyAsync(CancelPolicyRequest request, CancellationToken ct = default)
    {
        return await _client.Policies.CancelPolicyAsync(request, _client.BuildCallOptions(ct));
    }

    public async Task<ApproveCancellationResponse> ApproveCancellationAsync(ApproveCancellationRequest request, CancellationToken ct = default)
    {
        return await _client.Policies.ApproveCancellationAsync(request, _client.BuildCallOptions(ct));
    }

    public async Task<RenewPolicyTenureResponse> RenewPolicyAsync(RenewPolicyTenureRequest request, CancellationToken ct = default)
    {
        return await _client.Policies.RenewPolicyAsync(request, _client.BuildCallOptions(ct));
    }
}
