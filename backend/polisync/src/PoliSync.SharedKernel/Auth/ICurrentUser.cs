namespace PoliSync.SharedKernel.Auth;

/// <summary>
/// Represents the currently authenticated user from gRPC metadata.
/// </summary>
public interface ICurrentUser
{
    string UserId { get; }
    string Role { get; }
    string? TenantId { get; }
    bool IsAuthenticated { get; }
}
