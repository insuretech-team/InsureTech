using Grpc.Core;
using InsuranceEngine.Grpc.Clients;
using Insuretech.Policy.Services.V1;

namespace InsuranceEngine.Renewals.Infrastructure;

/// <summary>
/// Implementation of IRenewalDataGateway using gRPC calls to the Go backend's PolicyService.
/// Decoupled from the main Policy module.
/// </summary>
public sealed class GoRenewalDataGateway : IRenewalDataGateway
{
    private readonly InsuranceServiceClient _client;

    public GoRenewalDataGateway(InsuranceServiceClient client)
    {
        _client = client;
    }

    public async Task<RenewPolicyTenureResponse> RenewPolicyAsync(RenewPolicyTenureRequest request, CancellationToken ct = default)
    {
        return await _client.Policies.RenewPolicyAsync(request, _client.BuildCallOptions(ct));
    }
}
