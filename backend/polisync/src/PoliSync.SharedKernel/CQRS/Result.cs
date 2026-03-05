namespace PoliSync.SharedKernel.CQRS;

/// <summary>
/// Discriminated union result type. No domain exceptions — errors flow as data.
/// Maps to Grpc.Core.StatusCode in the gRPC interceptor layer.
/// </summary>
public sealed class Result<T>
{
    public T? Value { get; private init; }
    public ResultError? Error { get; private init; }
    public bool IsSuccess => Error is null;
    public bool IsFailure => !IsSuccess;

    private Result() { }

    public static Result<T> Ok(T value) => new() { Value = value };

    public static Result<T> Fail(string code, string message, ResultErrorKind kind = ResultErrorKind.DomainError)
        => new() { Error = new ResultError(code, message, kind) };

    public static Result<T> NotFound(string message)
        => Fail("NOT_FOUND", message, ResultErrorKind.NotFound);

    public static Result<T> Unauthorized(string message)
        => Fail("UNAUTHORIZED", message, ResultErrorKind.Unauthorized);

    public static Result<T> Conflict(string message)
        => Fail("CONFLICT", message, ResultErrorKind.Conflict);

    public Result<TOut> Map<TOut>(Func<T, TOut> mapper)
        => IsSuccess ? Result<TOut>.Ok(mapper(Value!)) : Result<TOut>.Fail(Error!.Code, Error.Message, Error.Kind);

    public T GetValueOrThrow()
        => IsSuccess ? Value! : throw new InvalidOperationException($"[{Error!.Code}] {Error.Message}");
}

/// <summary>Unit result for commands with no return value.</summary>
public sealed class Result
{
    public ResultError? Error { get; private init; }
    public bool IsSuccess => Error is null;
    public bool IsFailure => !IsSuccess;

    private Result() { }

    public static Result Ok() => new();

    public static Result Fail(string code, string message, ResultErrorKind kind = ResultErrorKind.DomainError)
        => new() { Error = new ResultError(code, message, kind) };

    public static Result NotFound(string message)
        => Fail("NOT_FOUND", message, ResultErrorKind.NotFound);

    public static Result Unauthorized(string message)
        => Fail("UNAUTHORIZED", message, ResultErrorKind.Unauthorized);
}

public sealed record ResultError(string Code, string Message, ResultErrorKind Kind);

public enum ResultErrorKind
{
    DomainError,
    NotFound,
    Unauthorized,
    Conflict,
    Validation,
    Internal
}
