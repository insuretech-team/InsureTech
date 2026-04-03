using Insuretech.Policy.Services.V1;

namespace InsuranceEngine.Endorsements;

public interface IEndorsementDataGateway
{
    Task<UpdatePolicyResponse> UpdatePolicyAsync(UpdatePolicyRequest request, CancellationToken ct = default);
}
