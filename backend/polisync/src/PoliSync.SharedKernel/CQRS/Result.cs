namespace PoliSync.SharedKernel.CQRS;

/// <summary>
/// Result pattern for operations without return value
/// </summary>
public sealed class Result
{
    public Error? Error { get; }
    public bool IsSuccess => Error is null;
    public bool IsFailure => !IsSuccess;

    private Result(Error? error)
    {
        Error = error;
    }

    public static Result Ok() => new(null);

    public static Result Success() => Ok();
    
    public static Result Fail(string code, string message) => 
        new(new Error(code, message));

    public static Result Failure(string message) =>
        Fail("OPERATION_FAILED", message);
    
    public static Result Fail(Error error) => new(error);

    public static Result<T> Ok<T>(T value) => Result<T>.Ok(value);
    
    public static Result<T> Fail<T>(string code, string message) => 
        Result<T>.Fail(code, message);
}

/// <summary>
/// Result pattern for operations with return value
/// </summary>
public sealed class Result<T>
{
    public T? Value { get; }
    public Error? Error { get; }
    public bool IsSuccess => Error is null;
    public bool IsFailure => !IsSuccess;

    private Result(T? value, Error? error)
    {
        Value = value;
        Error = error;
    }

    public static Result<T> Ok(T value) => new(value, null);

    public static Result<T> Success(T value) => Ok(value);
    
    public static Result<T> Fail(string code, string message) => 
        new(default, new Error(code, message));

    public static Result<T> Failure(string message) =>
        Fail("OPERATION_FAILED", message);
    
    public static Result<T> Fail(Error error) => new(default, error);

    public TResult Match<TResult>(
        Func<T, TResult> onSuccess,
        Func<Error, TResult> onFailure)
    {
        return IsSuccess ? onSuccess(Value!) : onFailure(Error!);
    }
}

/// <summary>
/// Error record for Result pattern
/// </summary>
public sealed record Error(string Code, string Message)
{
    public static readonly Error None = new(string.Empty, string.Empty);
    
    // Common error codes
    public static Error NotFound(string entity, string id) => 
        new("NOT_FOUND", $"{entity} with id {id} not found");
    
    public static Error Validation(string message) => 
        new("VALIDATION_ERROR", message);
    
    public static Error Unauthorized(string message = "Unauthorized") => 
        new("UNAUTHORIZED", message);
    
    public static Error Forbidden(string message = "Forbidden") => 
        new("FORBIDDEN", message);
    
    public static Error Conflict(string message) => 
        new("CONFLICT", message);
    
    public static Error Internal(string message = "Internal server error") => 
        new("INTERNAL_ERROR", message);
}
