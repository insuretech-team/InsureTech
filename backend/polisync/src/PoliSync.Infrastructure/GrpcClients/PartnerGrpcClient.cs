using Insuretech.Partner.Services.V1;

namespace PoliSync.Infrastructure.GrpcClients;

/// <summary>
/// Typed wrapper for the Go Partner service gRPC client.
/// </summary>
public sealed class PartnerGrpcClient
{
    private readonly GrpcClientFactory _factory;

    public PartnerGrpcClient(GrpcClientFactory factory) => _factory = factory;

    private PartnerService.PartnerServiceClient Client =>
        _factory.GetClient("PartnerService", ch => new PartnerService.PartnerServiceClient(ch));

    public async Task<bool> ExistsAsync(string partnerId, CancellationToken ct = default)
    {
        var resp = await Client.GetPartnerAsync(
            new GetPartnerRequest { PartnerId = partnerId }, cancellationToken: ct);
        return resp.Partner is not null;
    }

    public async Task<GetPartnerResponse> GetAsync(string partnerId, CancellationToken ct = default)
        => await Client.GetPartnerAsync(
            new GetPartnerRequest { PartnerId = partnerId }, cancellationToken: ct);
}
