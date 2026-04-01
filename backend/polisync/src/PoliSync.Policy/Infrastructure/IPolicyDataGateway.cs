using ProposalEntity = Insuretech.Policy.Entity.V1.InsuranceProposal;
using ProposalStatus = Insuretech.Policy.Entity.V1.ProposalStatus;
using PolicyEntity = Insuretech.Policy.Entity.V1.Policy;

namespace PoliSync.Policy.Infrastructure;

public interface IPolicyDataGateway
{
    Task<PolicyEntity> CreatePolicyAsync(PolicyEntity policy, CancellationToken cancellationToken = default);
    Task<PolicyEntity?> GetPolicyAsync(string policyId, CancellationToken cancellationToken = default);
    Task<PolicyEntity> UpdatePolicyAsync(PolicyEntity policy, CancellationToken cancellationToken = default);
    Task<IReadOnlyList<PolicyEntity>> ListPoliciesAsync(string customerId, int page, int pageSize, CancellationToken cancellationToken = default);
    Task DeletePolicyAsync(string policyId, CancellationToken cancellationToken = default);
    Task<ProposalEntity> CreateInsuranceProposalAsync(ProposalEntity proposal, CancellationToken cancellationToken = default);
    Task<ProposalEntity?> GetInsuranceProposalAsync(string proposalId, CancellationToken cancellationToken = default);
    Task<ProposalEntity> UpdateInsuranceProposalAsync(ProposalEntity proposal, CancellationToken cancellationToken = default);
    Task<IReadOnlyList<ProposalEntity>> ListInsuranceProposalsAsync(
        string? orderId,
        string? insurerId,
        string? customerId,
        ProposalStatus? status,
        int page,
        int pageSize,
        CancellationToken cancellationToken = default);
    Task DeleteInsuranceProposalAsync(string proposalId, CancellationToken cancellationToken = default);
}
