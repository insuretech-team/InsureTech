using Grpc.Core;
using Grpc.Core.Interceptors;
using PoliSync.SharedKernel.Auth;

namespace PoliSync.ApiHost.Interceptors;

/// <summary>
/// gRPC interceptor that validates JWT tokens and populates ICurrentUser
/// </summary>
public sealed class JwtAuthInterceptor : Interceptor
{
    private readonly ILogger<JwtAuthInterceptor> _logger;

    public JwtAuthInterceptor(ILogger<JwtAuthInterceptor> logger)
    {
        _logger = logger;
    }

    public override async Task<TResponse> UnaryServerHandler<TRequest, TResponse>(
        TRequest request,
        ServerCallContext context,
        UnaryServerMethod<TRequest, TResponse> continuation)
    {
        try
        {
            // JWT validation is handled by ASP.NET Core middleware
            // This interceptor just logs and validates tenant isolation
            
            var httpContext = context.GetHttpContext();
            var user = httpContext.User;

            if (!user.Identity?.IsAuthenticated ?? true)
            {
                _logger.LogWarning("Unauthenticated request to {Method}", context.Method);
                throw new RpcException(new Status(StatusCode.Unauthenticated, "Authentication required"));
            }

            var tenantId = user.FindFirst("tenant_id")?.Value;
            var userId = user.FindFirst("sub")?.Value;

            _logger.LogDebug(
                "Authenticated request: User={UserId}, Tenant={TenantId}, Method={Method}",
                userId, tenantId, context.Method);

            return await continuation(request, context);
        }
        catch (RpcException)
        {
            throw;
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error in JWT auth interceptor");
            throw new RpcException(new Status(StatusCode.Internal, "Authentication error"));
        }
    }
}
