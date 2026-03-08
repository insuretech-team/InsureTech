using Grpc.Core;
using Insuretech.Authz.Services.V1;

namespace PoliSync.Infrastructure.GrpcClients;

/// <summary>
/// Typed wrapper for the Go AuthZ service gRPC client.
/// Used for RBAC/ABAC permission checks before sensitive operations.
/// </summary>
public sealed class AuthzGrpcClient
{
    private readonly GrpcClientFactory _factory;

    public AuthzGrpcClient(GrpcClientFactory factory) => _factory = factory;

    private AuthZService.AuthZServiceClient Client =>
        _factory.GetClient("AuthzService", ch => new AuthZService.AuthZServiceClient(ch));

    /// <summary>Checks if a user has permission. Returns true if allowed.</summary>
    public async Task<bool> IsAllowedAsync(
        string userId, string tenantId, string resource, string action,
        CancellationToken ct = default)
    {
        try
        {
            var resp = await Client.CheckAccessAsync(new CheckAccessRequest
            {
                UserId   = userId,
                TenantId = tenantId,
                Resource = resource,
                Action   = action,
            }, cancellationToken: ct);
            return resp.Allowed;
        }
        catch (RpcException ex)
        {
            // fail-closed: deny on error
            throw new InvalidOperationException(
                $"AuthZ check failed for {userId} {action} {resource}: {ex.Status.Detail}", ex);
        }
    }
}
