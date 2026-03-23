using System;
using System.Threading;
using System.Threading.Tasks;
using InsuranceEngine.Commission.Domain.Entities;
using InsuranceEngine.Commission.Domain.Enums;
using InsuranceEngine.Commission.Domain.Interfaces;
using InsuranceEngine.Policy.Application.Interfaces;
using InsuranceEngine.Policy.Domain.Events;
using MediatR;
using Microsoft.Extensions.Logging;

namespace InsuranceEngine.Commission.Application.Features.Events;

public class PolicyRenewedEventHandler : INotificationHandler<PolicyRenewedEvent>
{
    private readonly ICommissionRepository _commissionRepository;
    private readonly IPolicyRepository _policyRepository;
    private readonly ILogger<PolicyRenewedEventHandler> _logger;

    public PolicyRenewedEventHandler(
        ICommissionRepository commissionRepository, 
        IPolicyRepository policyRepository,
        ILogger<PolicyRenewedEventHandler> logger)
    {
        _commissionRepository = commissionRepository;
        _policyRepository = policyRepository;
        _logger = logger;
    }

    public async Task Handle(PolicyRenewedEvent notification, CancellationToken cancellationToken)
    {
        _logger.LogInformation("Processing commission for Renewed Policy: {PolicyId}", notification.NewPolicyId);

        var policy = await _policyRepository.GetByIdAsync(notification.NewPolicyId, cancellationToken);
        if (policy == null)
        {
            _logger.LogError("New Policy {PolicyId} not found for commission calculation", notification.NewPolicyId);
            return;
        }

        if (!policy.PartnerId.HasValue && !policy.AgentId.HasValue)
        {
            _logger.LogInformation("No partner or agent associated with Renewed Policy {PolicyId}. Skipping commission.", notification.NewPolicyId);
            return;
        }

        // Renewal Rate: 5%
        const double rate = 0.05;
        var amount = (long)Math.Round(notification.PremiumAmount * rate, MidpointRounding.AwayFromZero);

        var commission = Domain.Entities.Commission.Create(
            notification.NewPolicyId,
            policy.PartnerId,
            policy.AgentId,
            CommissionType.Renewal,
            amount);

        await _commissionRepository.CreateAsync(commission, cancellationToken);
        _logger.LogInformation("Renewal Commission of {Amount} created for Policy {PolicyId}", amount, notification.NewPolicyId);
    }
}
