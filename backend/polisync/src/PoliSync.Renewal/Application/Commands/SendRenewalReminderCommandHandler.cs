using Google.Protobuf.WellKnownTypes;
using Insuretech.Renewal.Entity.V1;
using MediatR;
using Microsoft.Extensions.Logging;
using PoliSync.Renewal.Infrastructure;
using PoliSync.SharedKernel.CQRS;
using PoliSync.SharedKernel.Messaging;

namespace PoliSync.Renewal.Application.Commands;

public sealed class SendRenewalReminderCommandHandler : IRequestHandler<SendRenewalReminderCommand, Result<string>>
{
    private readonly IRenewalDataGateway _dataGateway;
    private readonly IEventBus _eventBus;
    private readonly ILogger<SendRenewalReminderCommandHandler> _logger;

    public SendRenewalReminderCommandHandler(
        IRenewalDataGateway dataGateway,
        IEventBus eventBus,
        ILogger<SendRenewalReminderCommandHandler> logger)
    {
        _dataGateway = dataGateway;
        _eventBus = eventBus;
        _logger = logger;
    }

    public async Task<Result<string>> Handle(SendRenewalReminderCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var schedule = await _dataGateway.GetRenewalScheduleAsync(request.RenewalScheduleId, cancellationToken);
            if (schedule is null)
                return Result.Fail<string>("SCHEDULE_NOT_FOUND", $"Renewal schedule {request.RenewalScheduleId} not found");

            var now = DateTime.UtcNow;
            var dueDate = schedule.RenewalDueDate?.ToDateTime() ?? now;
            var daysRemaining = Math.Max(0, (int)(dueDate - now).TotalDays);

            var channel = ParseChannel(request.Channel);
            var reminder = new RenewalReminder
            {
                Id = Guid.NewGuid().ToString("N"),
                RenewalScheduleId = schedule.Id,
                DaysBeforeRenewal = daysRemaining,
                Channel = channel,
                Status = ReminderStatus.Sent,
                ScheduledAt = Timestamp.FromDateTime(now),
                SentAt = Timestamp.FromDateTime(now),
                NotificationId = $"NTF-{Guid.NewGuid():N}"[..14]
            };

            var created = await _dataGateway.CreateRenewalReminderAsync(reminder, cancellationToken);

            if (schedule.Status == RenewalStatus.Pending)
            {
                schedule.Status = RenewalStatus.Reminded;
                await _dataGateway.UpdateRenewalScheduleAsync(schedule, cancellationToken);
            }

            await _eventBus.PublishAsync(
                new RenewalReminderSentEvent(created.Id, schedule.PolicyId, channel.ToString()),
                cancellationToken);

            _logger.LogInformation("Renewal reminder {ReminderId} sent for schedule {ScheduleId}",
                created.Id, request.RenewalScheduleId);
            return Result.Ok(created.Id);
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to send renewal reminder for schedule {ScheduleId}", request.RenewalScheduleId);
            return Result.Fail<string>("SEND_REMINDER_FAILED", ex.Message);
        }
    }

    private static ReminderChannel ParseChannel(string value)
    {
        if (string.IsNullOrWhiteSpace(value)) return ReminderChannel.Sms;
        return System.Enum.TryParse<ReminderChannel>(value, true, out var parsed) ? parsed : ReminderChannel.Sms;
    }
}

public sealed record RenewalReminderSentEvent(string ReminderId, string PolicyId, string Channel)
    : PoliSync.SharedKernel.Domain.DomainEvent;

