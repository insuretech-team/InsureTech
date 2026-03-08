using Grpc.Core;
using Insuretech.Fraud.Services.V1;
using Microsoft.Extensions.Logging;

namespace PoliSync.Infrastructure.GrpcClients;

/// <summary>
/// Typed wrapper for the Go Fraud detection service gRPC client.
/// Called before policy issuance and claim approval (SRS FR-058, FR-139).
/// </summary>
public sealed class FraudGrpcClient
{
    private readonly GrpcClientFactory _factory;
    private readonly ILogger<FraudGrpcClient> _logger;

    public FraudGrpcClient(GrpcClientFactory factory, ILogger<FraudGrpcClient> logger)
    { _factory = factory; _logger = logger; }

    private FraudService.FraudServiceClient Client =>
        _factory.GetClient("FraudService", ch => new FraudService.FraudServiceClient(ch));

    public async Task<FraudCheckResult> CheckAsync(
        string entityType, string entityId,
        Dictionary<string, string> metadata,
        CancellationToken ct = default)
    {
        try
        {
            var req = new CheckFraudRequest { EntityType = entityType, EntityId = entityId };
            req.Metadata.Add(metadata);
            var resp = await Client.CheckFraudAsync(req, cancellationToken: ct);
            return new FraudCheckResult(resp.RiskScore, resp.IsFlagged, [.. resp.Flags]);
        }
        catch (RpcException ex)
        {
            _logger.LogWarning(ex, "Fraud check failed for {EntityType}/{EntityId} — fail-open", entityType, entityId);
            return new FraudCheckResult(0.0, false, []);
        }
    }
}

public sealed record FraudCheckResult(double RiskScore, bool IsFlagged, string[] Flags);
