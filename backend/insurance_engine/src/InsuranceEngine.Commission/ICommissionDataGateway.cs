using Insuretech.Commission.Services.V1;

namespace InsuranceEngine.Commission;

public interface ICommissionDataGateway
{
    Task<CalculateCommissionResponse> CalculateCommissionAsync(CalculateCommissionRequest request, CancellationToken ct = default);
    Task<CreatePayoutResponse> CreatePayoutAsync(CreatePayoutRequest request, CancellationToken ct = default);
    Task<ProcessPayoutResponse> ProcessPayoutAsync(ProcessPayoutRequest request, CancellationToken ct = default);
    Task<ListCommissionsResponse> ListCommissionsAsync(ListCommissionsRequest request, CancellationToken ct = default);
}
