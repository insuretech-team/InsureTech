using Insuretech.Policy.Services.V1;

namespace InsuranceEngine.Cancellations;

public interface ICancellationDataGateway
{
    Task<ApproveCancellationResponse> ApproveCancellationAsync(ApproveCancellationRequest request, CancellationToken ct = default);
}
