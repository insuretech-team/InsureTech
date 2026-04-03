using Insuretech.Policy.Services.V1;

namespace InsuranceEngine.Endorsements;

public interface IEndorsementDataGateway
{
    Task<GetPolicyResponse> GetPolicyAsync(string policyId, CancellationToken ct = default);
    Task<UpdatePolicyResponse> UpdatePolicyAsync(string policyId, List<Insuretech.Policy.Entity.V1.Nominee>? nominees, CancellationToken ct = default);
}
