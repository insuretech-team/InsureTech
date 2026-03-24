using MediatR;

namespace InsuranceEngine.SharedKernel.CQRS;

public interface IQuery<TResponse> : IRequest<Result<TResponse>>
{
}
