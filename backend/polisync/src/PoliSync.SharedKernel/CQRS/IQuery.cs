using MediatR;

namespace PoliSync.SharedKernel.CQRS;

/// <summary>Marker for queries that return a Result.</summary>
public interface IQuery<TResult> : IRequest<Result<TResult>> { }

/// <summary>Handler for queries.</summary>
public interface IQueryHandler<TQuery, TResult> : IRequestHandler<TQuery, Result<TResult>>
    where TQuery : IQuery<TResult> { }
