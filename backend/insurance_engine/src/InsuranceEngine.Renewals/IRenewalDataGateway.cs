using Insuretech.Policy.Services.V1;

namespace InsuranceEngine.Renewals;

public interface IRenewalDataGateway
{
    Task<RenewPolicyTenureResponse> RenewPolicyAsync(RenewPolicyTenureRequest request, CancellationToken ct = default);
}
