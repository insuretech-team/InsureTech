using PolicyEntity = Insuretech.Policy.Entity.V1.Policy;

namespace PoliSync.Policy.Infrastructure;

public interface IPolicyDataGateway
{
    Task<PolicyEntity> CreatePolicyAsync(PolicyEntity policy, CancellationToken cancellationToken = default);
    Task<PolicyEntity?> GetPolicyAsync(string policyId, CancellationToken cancellationToken = default);
    Task<PolicyEntity> UpdatePolicyAsync(PolicyEntity policy, CancellationToken cancellationToken = default);
    Task<IReadOnlyList<PolicyEntity>> ListPoliciesAsync(string customerId, int page, int pageSize, CancellationToken cancellationToken = default);
    Task DeletePolicyAsync(string policyId, CancellationToken cancellationToken = default);
}
