using MediatR;

namespace PoliSync.SharedKernel.CQRS;

/// <summary>Marker for commands that return a Result.</summary>
public interface ICommand<TResult> : IRequest<Result<TResult>> { }

/// <summary>Marker for commands with no return value.</summary>
public interface ICommand : IRequest<Result> { }

/// <summary>Handler for commands returning Result of T.</summary>
public interface ICommandHandler<TCommand, TResult> : IRequestHandler<TCommand, Result<TResult>>
    where TCommand : ICommand<TResult> { }

/// <summary>Handler for commands with no return value.</summary>
public interface ICommandHandler<TCommand> : IRequestHandler<TCommand, Result>
    where TCommand : ICommand { }
