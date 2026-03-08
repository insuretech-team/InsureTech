using Grpc.Core;
using Insuretech.Kyc.Services.V1;

namespace PoliSync.Infrastructure.GrpcClients;

/// <summary>
/// Typed wrapper for the Go KYC service gRPC client.
/// Verifies customer KYC status before policy issuance (SRS SEC-012).
/// </summary>
public sealed class KycGrpcClient
{
    private readonly GrpcClientFactory _factory;

    public KycGrpcClient(GrpcClientFactory factory) => _factory = factory;

    private KYCService.KYCServiceClient Client =>
        _factory.GetClient("KycService", ch => new KYCService.KYCServiceClient(ch));

    public async Task<bool> IsVerifiedAsync(string userId, CancellationToken ct = default)
    {
        try
        {
            var resp = await Client.GetKYCVerificationAsync(
                new GetKYCVerificationRequest { UserId = userId }, cancellationToken: ct);
            return resp.Verification?.Status == "APPROVED";
        }
        catch (RpcException ex) when (ex.StatusCode == StatusCode.NotFound)
        {
            return false;
        }
    }
}
