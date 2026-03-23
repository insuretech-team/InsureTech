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

public class PolicyIssuedEventHandler : INotificationHandler<PolicyIssuedEvent>
{
    private readonly ICommissionRepository _commissionRepository;
    private readonly IPolicyRepository _policyRepository;
    private readonly ILogger<PolicyIssuedEventHandler> _logger;

    public PolicyIssuedEventHandler(
        ICommissionRepository commissionRepository, 
        IPolicyRepository policyRepository,
        ILogger<PolicyIssuedEventHandler> logger)
    {
        _commissionRepository = commissionRepository;
        _policyRepository = policyRepository;
        _logger = logger;
    }

    public async Task Handle(PolicyIssuedEvent notification, CancellationToken cancellationToken)
    {
        _logger.LogInformation("Processing commission for Issued Policy: {PolicyId}", notification.PolicyId);

        var policy = await _policyRepository.GetByIdAsync(notification.PolicyId, cancellationToken);
        if (policy == null)
        {
            _logger.LogError("Policy {PolicyId} not found for commission calculation", notification.PolicyId);
            return;
        }

        if (!policy.PartnerId.HasValue && !policy.AgentId.HasValue)
        {
            _logger.LogInformation("No partner or agent associated with Policy {PolicyId}. Skipping commission.", notification.PolicyId);
            return;
        }

        // Acquisition Rate: 15%
        const double rate = 0.15;
        var amount = (long)Math.Round(notification.PremiumAmount * rate, MidpointRounding.AwayFromZero);

        var commission = Domain.Entities.Commission.Create(
            notification.PolicyId,
            policy.PartnerId,
            policy.AgentId,
            CommissionType.Acquisition,
            amount);

        await _commissionRepository.CreateAsync(commission, cancellationToken);
        _logger.LogInformation("Commission of {Amount} created for Policy {PolicyId}", amount, notification.PolicyId);
    }
}
