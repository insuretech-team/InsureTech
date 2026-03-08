using Grpc.Core;
using Grpc.Core.Interceptors;
using System.Diagnostics;

namespace PoliSync.ApiHost.Interceptors;

/// <summary>
/// gRPC interceptor for request/response logging and timing
/// </summary>
public sealed class LoggingInterceptor : Interceptor
{
    private readonly ILogger<LoggingInterceptor> _logger;

    public LoggingInterceptor(ILogger<LoggingInterceptor> logger)
    {
        _logger = logger;
    }

    public override async Task<TResponse> UnaryServerHandler<TRequest, TResponse>(
        TRequest request,
        ServerCallContext context,
        UnaryServerMethod<TRequest, TResponse> continuation)
    {
        var stopwatch = Stopwatch.StartNew();
        var method = context.Method;

        _logger.LogInformation("gRPC Request: {Method}", method);

        try
        {
            var response = await continuation(request, context);
            
            stopwatch.Stop();
            
            _logger.LogInformation(
                "gRPC Response: {Method} completed in {ElapsedMs}ms",
                method, stopwatch.ElapsedMilliseconds);

            return response;
        }
        catch (RpcException ex)
        {
            stopwatch.Stop();
            
            _logger.LogWarning(
                "gRPC Error: {Method} failed with {StatusCode} in {ElapsedMs}ms: {Message}",
                method, ex.StatusCode, stopwatch.ElapsedMilliseconds, ex.Status.Detail);
            
            throw;
        }
        catch (Exception ex)
        {
            stopwatch.Stop();
            
            _logger.LogError(ex,
                "gRPC Exception: {Method} threw exception in {ElapsedMs}ms",
                method, stopwatch.ElapsedMilliseconds);
            
            throw new RpcException(new Status(StatusCode.Internal, "Internal server error"));
        }
    }
}
