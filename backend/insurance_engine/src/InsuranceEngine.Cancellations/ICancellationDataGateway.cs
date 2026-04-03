using Insuretech.Policy.Services.V1;

namespace InsuranceEngine.Cancellations;

public interface ICancellationDataGateway
{
    Task<CancelPolicyResponse> CancelPolicyAsync(CancelPolicyRequest request, CancellationToken ct = default);
    Task<ApproveCancellationResponse> ApproveCancellationAsync(ApproveCancellationRequest request, CancellationToken ct = default);
}
