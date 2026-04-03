using Grpc.Core;
using InsuranceEngine.Grpc.Clients;
using Insuretech.Policy.Services.V1;

namespace InsuranceEngine.Cancellations.Infrastructure;

/// <summary>
/// Implementation of ICancellationDataGateway using gRPC calls to the Go backend's PolicyService.
/// Decoupled from the main Policy module.
/// </summary>
public sealed class GoCancellationDataGateway : ICancellationDataGateway
{
    private readonly InsuranceServiceClient _client;

    public GoCancellationDataGateway(InsuranceServiceClient client)
    {
        _client = client;
    }

    public async Task<CancelPolicyResponse> CancelPolicyAsync(CancelPolicyRequest request, CancellationToken ct = default)
    {
        // Use the Policies from the gRPC shared client
        return await _client.Policies.CancelPolicyAsync(request, _client.BuildCallOptions(ct));
    }

    public async Task<ApproveCancellationResponse> ApproveCancellationAsync(ApproveCancellationRequest request, CancellationToken ct = default)
    {
        return await _client.Policies.ApproveCancellationAsync(request, _client.BuildCallOptions(ct));
    }
}
