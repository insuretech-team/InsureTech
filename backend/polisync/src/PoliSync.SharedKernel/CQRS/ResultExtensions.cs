using Grpc.Core;

namespace PoliSync.SharedKernel.CQRS;

/// <summary>
/// Maps Error to gRPC RpcException with the appropriate StatusCode.
/// </summary>
public static class ResultExtensions
{
    public static RpcException ToRpcException(this Error error)
    {
        var status = error.Code switch
        {
            "NOT_FOUND"         => new Status(StatusCode.NotFound,          error.Message),
            "UNAUTHORIZED"      => new Status(StatusCode.Unauthenticated,   error.Message),
            "FORBIDDEN"         => new Status(StatusCode.PermissionDenied,  error.Message),
            "CONFLICT"          => new Status(StatusCode.AlreadyExists,     error.Message),
            "VALIDATION_ERROR"  => new Status(StatusCode.InvalidArgument,   error.Message),
            "INTERNAL_ERROR"    => new Status(StatusCode.Internal,          error.Message),
            _                   => new Status(StatusCode.Unknown,           error.Message),
        };
        return new RpcException(status);
    }
}
