using Grpc.Core;
using Microsoft.AspNetCore.Http;
using PoliSync.SharedKernel.Auth;
using System.Security.Claims;

namespace PoliSync.Infrastructure.Auth;

/// <summary>
/// Current user implementation from JWT claims
/// </summary>
public sealed class CurrentUser : ICurrentUser
{
    private readonly IHttpContextAccessor _httpContextAccessor;
    private Guid? _metadataUserId;
    private Guid? _metadataTenantId;
    private Guid? _metadataPartnerId;
    private string[]? _metadataRoles;

    public CurrentUser(IHttpContextAccessor httpContextAccessor)
    {
        _httpContextAccessor = httpContextAccessor;
    }

    private ClaimsPrincipal? User => _httpContextAccessor.HttpContext?.User;

    public Guid UserId => _metadataUserId ?? GetGuidClaim("sub") ?? Guid.Empty;

    public Guid TenantId => _metadataTenantId ?? GetGuidClaim("tenant_id") ?? Guid.Empty;

    public Guid? PartnerId => _metadataPartnerId ?? GetGuidClaim("partner_id");

    public string[] Roles => _metadataRoles ?? User?.FindAll(ClaimTypes.Role)
        .Select(c => c.Value)
        .ToArray() ?? Array.Empty<string>();

    public string Email => User?.FindFirst(ClaimTypes.Email)?.Value ?? string.Empty;

    public bool IsAuthenticated => User?.Identity?.IsAuthenticated ?? false;

    public bool IsInRole(string role)
    {
        return User?.IsInRole(role) ?? false;
    }

    public bool HasAnyRole(params string[] roles)
    {
        return roles.Any(IsInRole);
    }

    public void Populate(ServerCallContext context)
    {
        var headers = context.RequestHeaders;
        _metadataUserId = TryParseGuid(headers.GetValue("x-user-id"));
        _metadataTenantId = TryParseGuid(headers.GetValue("x-tenant-id"));
        _metadataPartnerId = TryParseGuid(headers.GetValue("x-partner-id"));

        var rawRoles = headers.GetValue("x-roles");
        if (!string.IsNullOrWhiteSpace(rawRoles))
        {
            _metadataRoles = rawRoles
                .Split(',', StringSplitOptions.RemoveEmptyEntries | StringSplitOptions.TrimEntries);
        }
    }

    private Guid? GetGuidClaim(string claimType)
    {
        var claim = User?.FindFirst(claimType)?.Value;
        return Guid.TryParse(claim, out var guid) ? guid : null;
    }

    private static Guid? TryParseGuid(string? value)
    {
        return Guid.TryParse(value, out var guid) ? guid : null;
    }
}
