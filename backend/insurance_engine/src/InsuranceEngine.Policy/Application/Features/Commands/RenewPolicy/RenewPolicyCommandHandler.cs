using System;
using System.Threading;
using System.Threading.Tasks;
using MediatR;
using InsuranceEngine.Policy.Application.Interfaces;
using InsuranceEngine.Policy.Domain.Entities;
using InsuranceEngine.Policy.Domain.Enums;
using InsuranceEngine.Policy.Domain.Events;
using InsuranceEngine.Policy.Domain.Services;
using InsuranceEngine.SharedKernel.CQRS;
using InsuranceEngine.SharedKernel.Interfaces;

namespace InsuranceEngine.Policy.Application.Features.Commands.RenewPolicy;

public class RenewPolicyCommandHandler : IRequestHandler<RenewPolicyCommand, Result<RenewPolicyResponse>>
{
    private readonly IPolicyRepository _repo;
    private readonly PolicyNumberGenerator _policyNumberGenerator;
    private readonly IEventBus _eventBus;

    public RenewPolicyCommandHandler(IPolicyRepository repo, PolicyNumberGenerator gen, IEventBus eventBus)
    {
        _repo = repo;
        _policyNumberGenerator = gen;
        _eventBus = eventBus;
    }

    public async Task<Result<RenewPolicyResponse>> Handle(RenewPolicyCommand request, CancellationToken cancellationToken)
    {
        var oldPolicy = await _repo.GetByIdWithNomineesAsync(request.PolicyId);
        if (oldPolicy == null)
            return Result<RenewPolicyResponse>.Fail(Error.NotFound("Policy", request.PolicyId.ToString()));

        if (oldPolicy.Status != PolicyStatus.Active && oldPolicy.Status != PolicyStatus.Expired
            && oldPolicy.Status != PolicyStatus.GracePeriod)
            return Result<RenewPolicyResponse>.Fail(Error.InvalidStateTransition(
                "Policy can only be renewed from ACTIVE, EXPIRED, or GRACE_PERIOD status."));

        var productCode = await _repo.GetProductCodeAsync(oldPolicy.ProductId);
        var seqNum = await _repo.GetNextSequenceNumberAsync();
        var newPolicyNumber = _policyNumberGenerator.Generate(productCode ?? "UNK-000", seqNum);

        var newStartDate = oldPolicy.EndDate;
        var newEndDate = newStartDate.AddMonths(request.TenureMonths);

        var newPolicy = new PolicyEntity
        {
            Id = Guid.NewGuid(),
            PolicyNumber = newPolicyNumber,
            ProductId = oldPolicy.ProductId,
            CustomerId = oldPolicy.CustomerId,
            PartnerId = oldPolicy.PartnerId,
            AgentId = oldPolicy.AgentId,
            Status = PolicyStatus.PendingPayment,
            PremiumAmount = oldPolicy.PremiumAmount,
            PremiumCurrency = oldPolicy.PremiumCurrency,
            SumInsuredAmount = oldPolicy.SumInsuredAmount,
            SumInsuredCurrency = oldPolicy.SumInsuredCurrency,
            TenureMonths = request.TenureMonths,
            StartDate = newStartDate,
            EndDate = newEndDate,
            ProposerDetailsJson = oldPolicy.ProposerDetailsJson,
            ProviderName = oldPolicy.ProviderName,
            CreatedAt = DateTime.UtcNow,
            UpdatedAt = DateTime.UtcNow
        };

        var newPolicyId = await _repo.AddAsync(newPolicy);

        await _eventBus.PublishAsync("insurance.policy.v1", new PolicyRenewedEvent(
            OldPolicyId: oldPolicy.Id, NewPolicyId: newPolicyId,
            NewPolicyNumber: newPolicyNumber, RenewalDate: DateTime.UtcNow));

        return Result<RenewPolicyResponse>.Ok(new RenewPolicyResponse(newPolicyId, newPolicyNumber));
    }
}
