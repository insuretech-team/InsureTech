using Grpc.Core;
using Grpc.Core.Interceptors;
using Microsoft.Extensions.Logging;

namespace InsuranceEngine.SharedKernel.Infrastructure;

/// <summary>
/// Global gRPC server interceptor for structured error handling.
/// Catches unhandled exceptions and returns proper gRPC status codes.
/// </summary>
public class GlobalGrpcErrorInterceptor : Interceptor
{
    private readonly ILogger<GlobalGrpcErrorInterceptor> _logger;

    public GlobalGrpcErrorInterceptor(ILogger<GlobalGrpcErrorInterceptor> logger)
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
            return await continuation(request, context);
        }
        catch (RpcException)
        {
            // Already a gRPC exception, re-throw as-is
            throw;
        }
        catch (ArgumentException ex)
        {
            _logger.LogWarning(ex, "Validation error in {Method}", context.Method);
            throw new RpcException(new Status(StatusCode.InvalidArgument, ex.Message));
        }
        catch (KeyNotFoundException ex)
        {
            _logger.LogWarning(ex, "Resource not found in {Method}", context.Method);
            throw new RpcException(new Status(StatusCode.NotFound, ex.Message));
        }
        catch (InvalidOperationException ex)
        {
            _logger.LogWarning(ex, "Invalid operation in {Method}", context.Method);
            throw new RpcException(new Status(StatusCode.FailedPrecondition, ex.Message));
        }
        catch (UnauthorizedAccessException ex)
        {
            _logger.LogWarning(ex, "Unauthorized access in {Method}", context.Method);
            throw new RpcException(new Status(StatusCode.PermissionDenied, ex.Message));
        }
        catch (OperationCanceledException)
        {
            _logger.LogInformation("Request cancelled for {Method}", context.Method);
            throw new RpcException(new Status(StatusCode.Cancelled, "Request was cancelled"));
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Unhandled exception in {Method}", context.Method);
            throw new RpcException(new Status(StatusCode.Internal, "An internal error occurred. Please try again later."));
        }
    }
}
