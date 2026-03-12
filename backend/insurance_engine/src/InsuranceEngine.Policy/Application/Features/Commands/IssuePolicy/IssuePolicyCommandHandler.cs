using System;
using System.Threading;
using System.Threading.Tasks;
using MediatR;
using InsuranceEngine.Policy.Application.Interfaces;
using InsuranceEngine.Policy.Domain.Events;
using InsuranceEngine.SharedKernel.CQRS;
using InsuranceEngine.SharedKernel.Interfaces;

namespace InsuranceEngine.Policy.Application.Features.Commands.IssuePolicy;

public class IssuePolicyCommandHandler : IRequestHandler<IssuePolicyCommand, Result>
{
    private readonly IPolicyRepository _repo;
    private readonly IEventBus _eventBus;

    public IssuePolicyCommandHandler(IPolicyRepository repo, IEventBus eventBus)
    {
        _repo = repo;
        _eventBus = eventBus;
    }

    public async Task<Result> Handle(IssuePolicyCommand request, CancellationToken cancellationToken)
    {
        var policy = await _repo.GetByIdAsync(request.PolicyId);
        if (policy == null) return Result.Fail(Error.NotFound("Policy", request.PolicyId.ToString()));

        var result = policy.Issue(DateTime.UtcNow);
        if (result.IsFailure) return result;

        await _repo.UpdateAsync(policy);

        await _eventBus.PublishAsync("policy.events", new PolicyIssuedEvent(
            PolicyId: policy.Id, PolicyNumber: policy.PolicyNumber, IssuedAt: policy.IssuedAt!.Value));

        return Result.Ok();
    }
}
