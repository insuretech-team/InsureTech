using MediatR;
using System.Transactions;

namespace InsuranceEngine.SharedKernel.Behaviors;

public class TransactionBehavior<TRequest, TResponse> : IPipelineBehavior<TRequest, TResponse>
    where TRequest : IRequest<TResponse>
{
    public async Task<TResponse> Handle(TRequest request, RequestHandlerDelegate<TResponse> next, CancellationToken cancellationToken)
    {
        // For CQRS we might restrict this to Commands (implementing ICommand or by name convention)
        // But the document uses a generic TransactionBehavior<,>
        var requestName = typeof(TRequest).Name;
        
        if (requestName.EndsWith("Command"))
        {
            using var transaction = new TransactionScope(TransactionScopeAsyncFlowOption.Enabled);
            var response = await next();
            transaction.Complete();
            return response;
        }
        
        return await next();
    }
}
