using FluentValidation;
using Grpc.Core;
using Grpc.Core.Interceptors;

namespace PoliSync.ApiHost.Interceptors;

/// <summary>
/// gRPC interceptor for FluentValidation integration
/// </summary>
public sealed class ValidationInterceptor : Interceptor
{
    private readonly IServiceProvider _serviceProvider;
    private readonly ILogger<ValidationInterceptor> _logger;

    public ValidationInterceptor(
        IServiceProvider serviceProvider,
        ILogger<ValidationInterceptor> logger)
    {
        _serviceProvider = serviceProvider;
        _logger = logger;
    }

    public override async Task<TResponse> UnaryServerHandler<TRequest, TResponse>(
        TRequest request,
        ServerCallContext context,
        UnaryServerMethod<TRequest, TResponse> continuation)
    {
        // Try to get validator for request type
        var validatorType = typeof(IValidator<>).MakeGenericType(typeof(TRequest));
        var validator = _serviceProvider.GetService(validatorType) as IValidator;

        if (validator != null)
        {
            var validationContext = new ValidationContext<TRequest>(request);
            var validationResult = await validator.ValidateAsync(validationContext);

            if (!validationResult.IsValid)
            {
                var errors = string.Join("; ", validationResult.Errors.Select(e => e.ErrorMessage));
                
                _logger.LogWarning(
                    "Validation failed for {Method}: {Errors}",
                    context.Method, errors);

                throw new RpcException(new Status(
                    StatusCode.InvalidArgument,
                    $"Validation failed: {errors}"));
            }
        }

        return await continuation(request, context);
    }
}
