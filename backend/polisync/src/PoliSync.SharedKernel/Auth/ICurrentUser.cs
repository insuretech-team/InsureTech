namespace PoliSync.SharedKernel.Auth;

/// <summary>
/// Current authenticated user context from JWT token
/// </summary>
public interface ICurrentUser
{
    Guid UserId { get; }
    Guid TenantId { get; }
    Guid? PartnerId { get; }
    string[] Roles { get; }
    string Email { get; }
    bool IsAuthenticated { get; }
    
    bool IsInRole(string role);
    bool HasAnyRole(params string[] roles);
}
