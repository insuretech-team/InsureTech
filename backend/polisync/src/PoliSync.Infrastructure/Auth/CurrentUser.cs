using PoliSync.SharedKernel.Auth;

namespace PoliSync.Infrastructure.Auth;

/// <summary>
/// Scoped service populated by AuthInterceptor from gRPC metadata.
/// One instance per request.
/// </summary>
public class CurrentUser : ICurrentUser
{
    public string UserId { get; private set; } = string.Empty;
    public string Role { get; private set; } = string.Empty;
    public string? TenantId { get; private set; }
    public bool IsAuthenticated => !string.IsNullOrEmpty(UserId);

    /// <summary>
    /// Called by the auth interceptor to populate user context from gRPC headers.
    /// </summary>
    public void Populate(string userId, string role, string? tenantId = null)
    {
        UserId = userId;
        Role = role;
        TenantId = tenantId;
    }
}
