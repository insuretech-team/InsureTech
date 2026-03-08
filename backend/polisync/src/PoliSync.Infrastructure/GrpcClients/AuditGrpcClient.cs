using Insuretech.Audit.Services.V1;
using Microsoft.Extensions.Logging;

namespace PoliSync.Infrastructure.GrpcClients;

/// <summary>
/// Typed wrapper for the Go Audit service gRPC client.
/// Fire-and-forget — audit failures never block business operations.
/// </summary>
public sealed class AuditGrpcClient
{
    private readonly GrpcClientFactory _factory;
    private readonly ILogger<AuditGrpcClient> _logger;

    public AuditGrpcClient(GrpcClientFactory factory, ILogger<AuditGrpcClient> logger)
    { _factory = factory; _logger = logger; }

    private AuditService.AuditServiceClient Client =>
        _factory.GetClient("AuditService", ch => new AuditService.AuditServiceClient(ch));

    public async Task LogAsync(
        string userId, string tenantId, string action,
        string resource, string resourceId, string details,
        CancellationToken ct = default)
    {
        try
        {
            await Client.CreateAuditLogAsync(new CreateAuditLogRequest
            {
                UserId     = userId,
                TenantId   = tenantId,
                Action     = action,
                Resource   = resource,
                ResourceId = resourceId,
                Details    = details,
            }, cancellationToken: ct);
        }
        catch (Exception ex)
        {
            _logger.LogWarning(ex, "Audit log failed for {Action} {Resource}/{ResourceId}", action, resource, resourceId);
        }
    }
}
