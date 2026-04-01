using Google.Protobuf.WellKnownTypes;
using Insuretech.Renewal.Entity.V1;
using MediatR;
using Microsoft.Extensions.Logging;
using PoliSync.Renewal.Infrastructure;
using PoliSync.SharedKernel.CQRS;
using PoliSync.SharedKernel.Messaging;

namespace PoliSync.Renewal.Application.Commands;

public sealed class RenewPolicyCommandHandler : IRequestHandler<RenewPolicyCommand, Result<string>>
{
    private readonly IRenewalDataGateway _dataGateway;
    private readonly IEventBus _eventBus;
    private readonly ILogger<RenewPolicyCommandHandler> _logger;

    public RenewPolicyCommandHandler(
        IRenewalDataGateway dataGateway,
        IEventBus eventBus,
        ILogger<RenewPolicyCommandHandler> logger)
    {
        _dataGateway = dataGateway;
        _eventBus = eventBus;
        _logger = logger;
    }

    public async Task<Result<string>> Handle(RenewPolicyCommand request, CancellationToken cancellationToken)
    {
        try
        {
            if (string.IsNullOrWhiteSpace(request.PolicyId))
                return Result.Fail<string>("VALIDATION_ERROR", "PolicyId is required");

            // Ensure a renewal schedule exists for the policy
            var schedules = await _dataGateway.ListRenewalSchedulesAsync(request.PolicyId, cancellationToken);
            var schedule = schedules.OrderByDescending(s => s.RenewalDueDate?.Seconds ?? 0).FirstOrDefault();

            if (schedule is null)
            {
                var dueDate = DateTime.UtcNow.AddDays(30);
                schedule = await _dataGateway.CreateRenewalScheduleAsync(new RenewalSchedule
                {
                    Id = Guid.NewGuid().ToString("N"),
                    PolicyId = request.PolicyId,
                    RenewalDueDate = Timestamp.FromDateTime(dueDate),
                    RenewalPremium = new Insuretech.Common.V1.Money { Amount = 250_000, Currency = "BDT" },
                    RenewalType = RenewalType.Manual,
                    Status = RenewalStatus.Pending,
                    GracePeriodDays = 30,
                    GracePeriodEnd = Timestamp.FromDateTime(dueDate.AddDays(30))
                }, cancellationToken);
            }

            var newPolicyId = string.IsNullOrWhiteSpace(schedule.RenewedPolicyId)
                ? $"POL-{Guid.NewGuid():N}"[..16]
                : schedule.RenewedPolicyId;

            schedule.Status = RenewalStatus.Renewed;
            schedule.RenewedPolicyId = newPolicyId;
            schedule.RenewedAt = Timestamp.FromDateTime(DateTime.UtcNow);
            await _dataGateway.UpdateRenewalScheduleAsync(schedule, cancellationToken);

            // Update grace period if it exists
            var gracePeriod = await _dataGateway.GetGracePeriodByPolicyAsync(request.PolicyId, cancellationToken);
            if (gracePeriod is not null)
            {
                gracePeriod.Status = GracePeriodStatus.Revived;
                gracePeriod.CoverageActive = true;
                gracePeriod.DaysRemaining = 0;
                gracePeriod.RevivedAt = Timestamp.FromDateTime(DateTime.UtcNow);
                await _dataGateway.UpdateGracePeriodAsync(gracePeriod, cancellationToken);
            }

            await _eventBus.PublishAsync(
                new PolicyRenewedViaRenewalEvent(request.PolicyId, newPolicyId, schedule.Id),
                cancellationToken);

            _logger.LogInformation("Policy {PolicyId} renewed as {NewPolicyId}", request.PolicyId, newPolicyId);
            return Result.Ok(newPolicyId);
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to renew policy {PolicyId}", request.PolicyId);
            return Result.Fail<string>("RENEW_POLICY_FAILED", ex.Message);
        }
    }
}

public sealed record PolicyRenewedViaRenewalEvent(string OldPolicyId, string NewPolicyId, string ScheduleId)
    : PoliSync.SharedKernel.Domain.DomainEvent;

