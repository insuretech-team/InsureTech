using Grpc.Core;
using Insuretech.Insurance.Services.V1;
using PoliSync.Infrastructure.Clients;
using PolicyEntity = Insuretech.Policy.Entity.V1.Policy;

namespace PoliSync.Policy.Infrastructure;

public sealed class GoPolicyDataGateway : IPolicyDataGateway
{
    private readonly InsuranceServiceClient _insuranceClient;

    public GoPolicyDataGateway(InsuranceServiceClient insuranceClient)
    {
        _insuranceClient = insuranceClient;
    }

    public async Task<PolicyEntity> CreatePolicyAsync(PolicyEntity policy, CancellationToken cancellationToken = default)
    {
        var response = await _insuranceClient.Client.CreatePolicyAsync(new CreatePolicyRequest { Policy = policy }, cancellationToken: cancellationToken);
        return response.Policy;
    }

    public async Task<PolicyEntity?> GetPolicyAsync(string policyId, CancellationToken cancellationToken = default)
    {
        try
        {
            var response = await _insuranceClient.Client.GetPolicyAsync(new GetPolicyRequest { PolicyId = policyId }, cancellationToken: cancellationToken);
            return response.Policy;
        }
        catch (RpcException ex) when (ex.StatusCode == StatusCode.NotFound)
        {
            return null;
        }
    }

    public async Task<PolicyEntity> UpdatePolicyAsync(PolicyEntity policy, CancellationToken cancellationToken = default)
    {
        var response = await _insuranceClient.Client.UpdatePolicyAsync(new UpdatePolicyRequest { Policy = policy }, cancellationToken: cancellationToken);
        return response.Policy;
    }

    public async Task<IReadOnlyList<PolicyEntity>> ListPoliciesAsync(string customerId, int page, int pageSize, CancellationToken cancellationToken = default)
    {
        var response = await _insuranceClient.Client.ListPoliciesAsync(new ListPoliciesRequest
        {
            CustomerId = customerId,
            Page = page,
            PageSize = pageSize
        }, cancellationToken: cancellationToken);

        return response.Policies;
    }

    public Task DeletePolicyAsync(string policyId, CancellationToken cancellationToken = default)
    {
        return _insuranceClient.Client.DeletePolicyAsync(new DeletePolicyRequest { PolicyId = policyId }, cancellationToken: cancellationToken).ResponseAsync;
    }
}
