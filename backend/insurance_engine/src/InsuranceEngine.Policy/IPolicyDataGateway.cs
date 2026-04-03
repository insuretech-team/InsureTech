using Insuretech.Policy.Services.V1;

namespace InsuranceEngine.Policy;

public interface IPolicyDataGateway
{
    Task<CreatePolicyResponse> CreatePolicyAsync(CreatePolicyRequest request, CancellationToken ct = default);
    Task<GetPolicyResponse> GetPolicyAsync(string policyId, CancellationToken ct = default);
    Task<ListUserPoliciesResponse> ListUserPoliciesAsync(ListUserPoliciesRequest request, CancellationToken ct = default);
    Task<UpdatePolicyResponse> UpdatePolicyAsync(UpdatePolicyRequest request, CancellationToken ct = default);
    Task<CancelPolicyResponse> CancelPolicyAsync(CancelPolicyRequest request, CancellationToken ct = default);
    Task<RenewPolicyTenureResponse> RenewPolicyAsync(RenewPolicyTenureRequest request, CancellationToken ct = default);
    Task<GeneratePolicyDocumentResponse> GeneratePolicyDocumentAsync(string policyId, CancellationToken ct = default);
    Task<IssuePolicyResponse> IssuePolicyAsync(IssuePolicyRequest request, CancellationToken ct = default);
    Task<ApproveCancellationResponse> ApproveCancellationAsync(ApproveCancellationRequest request, CancellationToken ct = default);
}
