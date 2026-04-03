using MediatR;

namespace InsuranceEngine.SharedKernel.CQRS;

/// <summary>
/// Marker interface for commands (write operations)
/// </summary>
public interface ICommand : IRequest<Result>
{
}

/// <summary>
/// Marker interface for commands that return a value
/// </summary>
public interface ICommand<TResponse> : IRequest<Result<TResponse>>
{
}
