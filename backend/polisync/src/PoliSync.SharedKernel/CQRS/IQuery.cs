using MediatR;

namespace PoliSync.SharedKernel.CQRS;

/// <summary>
/// Marker interface for queries (read operations)
/// </summary>
public interface IQuery<TResponse> : IRequest<Result<TResponse>>
{
}
