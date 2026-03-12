using System;
using System.Threading;
using System.Threading.Tasks;
using MediatR;
using InsuranceEngine.Policy.Application.Interfaces;
using InsuranceEngine.Policy.Domain.Events;
using InsuranceEngine.SharedKernel.CQRS;
using InsuranceEngine.SharedKernel.Interfaces;

namespace InsuranceEngine.Policy.Application.Features.Commands.CancelPolicy;

public class CancelPolicyCommandHandler : IRequestHandler<CancelPolicyCommand, Result>
{
    private readonly IPolicyRepository _repo;
    private readonly IEventBus _eventBus;

    public CancelPolicyCommandHandler(IPolicyRepository repo, IEventBus eventBus)
    {
        _repo = repo;
        _eventBus = eventBus;
    }

    public async Task<Result> Handle(CancelPolicyCommand request, CancellationToken cancellationToken)
    {
        var policy = await _repo.GetByIdAsync(request.PolicyId);
        if (policy == null) return Result.Fail(Error.NotFound("Policy", request.PolicyId.ToString()));

        var result = policy.Cancel(request.Reason);
        if (result.IsFailure) return result;

        await _repo.UpdateAsync(policy);

        await _eventBus.PublishAsync("policy.events", new PolicyCancelledEvent(
            PolicyId: policy.Id, PolicyNumber: policy.PolicyNumber,
            CancelledAt: DateTime.UtcNow, Reason: request.Reason));

        return Result.Ok();
    }
}
