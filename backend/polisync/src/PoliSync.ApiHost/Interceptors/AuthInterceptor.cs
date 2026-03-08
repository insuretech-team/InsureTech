using Grpc.Core;
using Grpc.Core.Interceptors;
using PoliSync.Infrastructure.Auth;

namespace PoliSync.ApiHost.Interceptors;

/// <summary>
/// Populates CurrentUser from gRPC metadata injected by the Go gateway.
/// 
/// IMPORTANT: CurrentUser is scoped (per-request). This interceptor resolves it
/// from the request's IServiceScope — the SAME instance that MediatR handlers
/// receive via constructor injection. Do NOT create a new scope here.
/// 
/// The Go gateway (auth_middleware.go) validates JWT and injects:
///   x-user-id, x-tenant-id, x-partner-id, x-token-id,
///   x-user-type, x-portal, x-roles, x-request-id, x-session-id
/// </summary>
public sealed class AuthInterceptor : Interceptor
{
    private readonly IServiceProvider _rootProvider;

    public AuthInterceptor(IServiceProvider rootProvider) => _rootProvider = rootProvider;

    public override async Task<TResponse> UnaryServerHandler<TRequest, TResponse>(
        TRequest request, ServerCallContext context,
        UnaryServerMethod<TRequest, TResponse> continuation)
    {
        PopulateCurrentUser(context);
        return await continuation(request, context);
    }

    public override async Task ServerStreamingServerHandler<TRequest, TResponse>(
        TRequest request, IServerStreamWriter<TResponse> responseStream,
        ServerCallContext context, ServerStreamingServerMethod<TRequest, TResponse> continuation)
    {
        PopulateCurrentUser(context);
        await continuation(request, responseStream, context);
    }

    public override async Task<TResponse> ClientStreamingServerHandler<TRequest, TResponse>(
        IAsyncStreamReader<TRequest> requestStream, ServerCallContext context,
        ClientStreamingServerMethod<TRequest, TResponse> continuation)
    {
        PopulateCurrentUser(context);
        return await continuation(requestStream, context);
    }

    public override async Task DuplexStreamingServerHandler<TRequest, TResponse>(
        IAsyncStreamReader<TRequest> requestStream, IServerStreamWriter<TResponse> responseStream,
        ServerCallContext context, DuplexStreamingServerMethod<TRequest, TResponse> continuation)
    {
        PopulateCurrentUser(context);
        await continuation(requestStream, responseStream, context);
    }

    /// <summary>
    /// Resolves the scoped CurrentUser from the gRPC request scope and populates it.
    /// ASP.NET Core gRPC creates a new DI scope per request — we resolve from it
    /// via the HttpContext's RequestServices.
    /// </summary>
    private static void PopulateCurrentUser(ServerCallContext context)
    {
        // For Grpc.AspNetCore, ServerCallContext is actually HttpContextServerCallContext
        // which exposes the HttpContext (and thus the request-scoped IServiceProvider).
        var httpContext = context.GetHttpContext();
        if (httpContext?.RequestServices.GetService(typeof(CurrentUser)) is CurrentUser currentUser)
        {
            currentUser.Populate(context);
        }
    }
}
