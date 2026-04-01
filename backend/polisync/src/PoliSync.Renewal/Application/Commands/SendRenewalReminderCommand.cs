using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Renewal.Application.Commands;

public sealed record SendRenewalReminderCommand(
    string RenewalScheduleId,
    string Channel // SMS, EMAIL, PUSH
) : ICommand<string>; // Returns reminder_id
