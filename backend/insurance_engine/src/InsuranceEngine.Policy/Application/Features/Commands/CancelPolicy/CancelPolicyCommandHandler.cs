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

        // FR-053: Pro-rata refund calculation
        long refundAmount = 0;
        var now = DateTime.UtcNow;
        if (policy.StartDate <= now && policy.EndDate >= now && policy.PremiumAmount > 0)
        {
            var totalDays = (policy.EndDate - policy.StartDate).TotalDays;
            var daysUsed = (now - policy.StartDate).TotalDays;
            if (totalDays > 0 && daysUsed < totalDays)
            {
                var remainingDays = totalDays - daysUsed;
                var refundFactor = remainingDays / totalDays;
                refundAmount = (long)Math.Round(policy.PremiumAmount * refundFactor, MidpointRounding.AwayFromZero);
            }
        }
        else if (policy.StartDate > now)
        {
            // Full refund if cancelled before start
            refundAmount = policy.PremiumAmount;
        }

        await _repo.UpdateAsync(policy);

        await _eventBus.PublishAsync("insurance.policy.v1", new PolicyCancelledEvent(
            PolicyId: policy.Id, PolicyNumber: policy.PolicyNumber,
            CancelledAt: now, Reason: request.Reason,
            RefundAmount: refundAmount, RefundCurrency: policy.PremiumCurrency));

        return Result.Ok();
    }
}
