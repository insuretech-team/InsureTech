using Google.Protobuf.WellKnownTypes;
using Grpc.Core;
using Insuretech.Common.V1;
using Insuretech.Renewal.Entity.V1;
using Insuretech.Renewal.Services.V1;
using MediatR;
using Microsoft.Extensions.Logging;
using PoliSync.Renewal.Infrastructure;

namespace PoliSync.Renewal.GrpcServices;

public sealed class RenewalGrpcService : RenewalService.RenewalServiceBase
{
    private readonly IMediator _mediator;
    private readonly ILogger<RenewalGrpcService> _logger;
    private readonly IRenewalDataGateway _dataGateway;

    public RenewalGrpcService(
        IMediator mediator,
        ILogger<RenewalGrpcService> logger,
        IRenewalDataGateway dataGateway)
    {
        _mediator = mediator;
        _logger = logger;
        _dataGateway = dataGateway;
    }

    public override async Task<GetRenewalScheduleResponse> GetRenewalSchedule(GetRenewalScheduleRequest request, ServerCallContext context)
    {
        if (string.IsNullOrWhiteSpace(request.PolicyId))
        {
            return new GetRenewalScheduleResponse
            {
                Error = BuildError("VALIDATION_ERROR", "PolicyId is required")
            };
        }

        try
        {
            var schedule = await EnsureScheduleAsync(request.PolicyId, GetCancellationToken(context));
            var reminders = await _dataGateway.ListRenewalRemindersAsync(schedule.Id, GetCancellationToken(context));

            var response = new GetRenewalScheduleResponse
            {
                RenewalSchedule = schedule
            };
            response.Reminders.AddRange(reminders.OrderByDescending(x => x.ScheduledAt?.Seconds ?? 0));
            return response;
        }
        catch (RpcException ex)
        {
            _logger.LogError(ex, "Failed to get renewal schedule for policy {PolicyId}", request.PolicyId);
            return new GetRenewalScheduleResponse
            {
                Error = BuildError("UPSTREAM_ERROR", ex.Status.Detail)
            };
        }
    }

    public override async Task<ListUpcomingRenewalsResponse> ListUpcomingRenewals(ListUpcomingRenewalsRequest request, ServerCallContext context)
    {
        var daysAhead = request.DaysAhead <= 0 ? 30 : request.DaysAhead;
        var page = request.Page <= 0 ? 1 : request.Page;
        var pageSize = request.PageSize <= 0 ? 20 : request.PageSize;

        var status = ParseRenewalStatus(request.Status);
        var now = DateTime.UtcNow;
        var threshold = now.AddDays(daysAhead);

        try
        {
            var schedules = await _dataGateway.ListRenewalSchedulesAsync(string.Empty, GetCancellationToken(context));
            var query = schedules.Where(x => x.RenewalDueDate.ToDateTime() <= threshold);

            if (status != RenewalStatus.Unspecified)
            {
                query = query.Where(x => x.Status == status);
            }

            var ordered = query.OrderBy(x => x.RenewalDueDate?.Seconds ?? 0).ToList();
            var pageItems = ordered.Skip((page - 1) * pageSize).Take(pageSize).ToList();

            var response = new ListUpcomingRenewalsResponse { TotalCount = ordered.Count };
            response.RenewalSchedules.AddRange(pageItems);
            return response;
        }
        catch (RpcException ex)
        {
            _logger.LogError(ex, "Failed to list upcoming renewals");
            return new ListUpcomingRenewalsResponse
            {
                Error = BuildError("UPSTREAM_ERROR", ex.Status.Detail)
            };
        }
    }

    public override async Task<RenewPolicyResponse> RenewPolicy(RenewPolicyRequest request, ServerCallContext context)
    {
        if (string.IsNullOrWhiteSpace(request.PolicyId))
        {
            return new RenewPolicyResponse
            {
                Error = BuildError("VALIDATION_ERROR", "PolicyId is required")
            };
        }

        try
        {
            var schedule = await EnsureScheduleAsync(request.PolicyId, GetCancellationToken(context));
            var newPolicyId = string.IsNullOrWhiteSpace(schedule.RenewedPolicyId)
                ? $"POL-{Guid.NewGuid():N}"[..16]
                : schedule.RenewedPolicyId;

            schedule.Status = RenewalStatus.Renewed;
            schedule.RenewedPolicyId = newPolicyId;
            schedule.RenewedAt = Timestamp.FromDateTime(DateTime.UtcNow);
            await _dataGateway.UpdateRenewalScheduleAsync(schedule, GetCancellationToken(context));

            var gracePeriod = await _dataGateway.GetGracePeriodByPolicyAsync(request.PolicyId, GetCancellationToken(context));
            if (gracePeriod is not null)
            {
                gracePeriod.Status = GracePeriodStatus.Revived;
                gracePeriod.CoverageActive = true;
                gracePeriod.DaysRemaining = 0;
                gracePeriod.RevivedAt = Timestamp.FromDateTime(DateTime.UtcNow);
                await _dataGateway.UpdateGracePeriodAsync(gracePeriod, GetCancellationToken(context));
            }

            return new RenewPolicyResponse
            {
                NewPolicyId = newPolicyId,
                Message = "Policy renewed successfully"
            };
        }
        catch (RpcException ex)
        {
            _logger.LogError(ex, "Failed to renew policy {PolicyId}", request.PolicyId);
            return new RenewPolicyResponse
            {
                Error = BuildError("UPSTREAM_ERROR", ex.Status.Detail)
            };
        }
    }

    public override async Task<SendRenewalReminderResponse> SendRenewalReminder(SendRenewalReminderRequest request, ServerCallContext context)
    {
        if (string.IsNullOrWhiteSpace(request.RenewalScheduleId))
        {
            return new SendRenewalReminderResponse
            {
                Error = BuildError("VALIDATION_ERROR", "RenewalScheduleId is required")
            };
        }

        try
        {
            var schedule = await _dataGateway.GetRenewalScheduleAsync(request.RenewalScheduleId, GetCancellationToken(context));
            if (schedule is null)
            {
                return new SendRenewalReminderResponse
                {
                    Error = BuildError("NOT_FOUND", "Renewal schedule not found")
                };
            }

            var now = DateTime.UtcNow;
            var dueDate = schedule.RenewalDueDate?.ToDateTime() ?? now;
            var reminder = new RenewalReminder
            {
                Id = Guid.NewGuid().ToString("N"),
                RenewalScheduleId = schedule.Id,
                DaysBeforeRenewal = Math.Max(0, (int)(dueDate - now).TotalDays),
                Channel = ParseReminderChannel(request.Channel),
                Status = ReminderStatus.Sent,
                ScheduledAt = Timestamp.FromDateTime(now),
                SentAt = Timestamp.FromDateTime(now),
                NotificationId = $"NTF-{Guid.NewGuid():N}"[..14]
            };

            var createdReminder = await _dataGateway.CreateRenewalReminderAsync(reminder, GetCancellationToken(context));

            if (schedule.Status == RenewalStatus.Pending)
            {
                schedule.Status = RenewalStatus.Reminded;
                await _dataGateway.UpdateRenewalScheduleAsync(schedule, GetCancellationToken(context));
            }

            return new SendRenewalReminderResponse
            {
                ReminderId = createdReminder.Id,
                Message = "Renewal reminder sent"
            };
        }
        catch (RpcException ex)
        {
            _logger.LogError(ex, "Failed to send reminder for schedule {ScheduleId}", request.RenewalScheduleId);
            return new SendRenewalReminderResponse
            {
                Error = BuildError("UPSTREAM_ERROR", ex.Status.Detail)
            };
        }
    }

    public override async Task<GetGracePeriodResponse> GetGracePeriod(GetGracePeriodRequest request, ServerCallContext context)
    {
        if (string.IsNullOrWhiteSpace(request.PolicyId))
        {
            return new GetGracePeriodResponse
            {
                Error = BuildError("VALIDATION_ERROR", "PolicyId is required")
            };
        }

        try
        {
            await EnsureScheduleAsync(request.PolicyId, GetCancellationToken(context));

            var gracePeriod = await _dataGateway.GetGracePeriodByPolicyAsync(request.PolicyId, GetCancellationToken(context));
            if (gracePeriod is null)
            {
                return new GetGracePeriodResponse
                {
                    Error = BuildError("NOT_FOUND", "Grace period not found")
                };
            }

            var changed = false;
            var now = DateTime.UtcNow;
            var endDate = gracePeriod.EndDate?.ToDateTime() ?? now;

            if (gracePeriod.Status != GracePeriodStatus.Revived)
            {
                var recalculatedDays = Math.Max(0, (int)(endDate.Date - now.Date).TotalDays);
                if (gracePeriod.DaysRemaining != recalculatedDays)
                {
                    gracePeriod.DaysRemaining = recalculatedDays;
                    changed = true;
                }

                if (now > endDate &&
                    (gracePeriod.Status != GracePeriodStatus.Lapsed || gracePeriod.CoverageActive || gracePeriod.DaysRemaining != 0))
                {
                    gracePeriod.Status = GracePeriodStatus.Lapsed;
                    gracePeriod.CoverageActive = false;
                    gracePeriod.DaysRemaining = 0;
                    changed = true;
                }
            }

            if (changed)
            {
                gracePeriod = await _dataGateway.UpdateGracePeriodAsync(gracePeriod, GetCancellationToken(context));
            }

            return new GetGracePeriodResponse
            {
                GracePeriod = gracePeriod
            };
        }
        catch (RpcException ex)
        {
            _logger.LogError(ex, "Failed to get grace period for policy {PolicyId}", request.PolicyId);
            return new GetGracePeriodResponse
            {
                Error = BuildError("UPSTREAM_ERROR", ex.Status.Detail)
            };
        }
    }

    public override async Task<RevivePolicyResponse> RevivePolicy(RevivePolicyRequest request, ServerCallContext context)
    {
        if (string.IsNullOrWhiteSpace(request.PolicyId))
        {
            return new RevivePolicyResponse
            {
                Error = BuildError("VALIDATION_ERROR", "PolicyId is required")
            };
        }

        try
        {
            var schedule = await EnsureScheduleAsync(request.PolicyId, GetCancellationToken(context));
            var gracePeriod = await _dataGateway.GetGracePeriodByPolicyAsync(request.PolicyId, GetCancellationToken(context));
            if (gracePeriod is null)
            {
                return new RevivePolicyResponse
                {
                    Error = BuildError("NOT_FOUND", "Grace period not found")
                };
            }

            if (gracePeriod.Status == GracePeriodStatus.Lapsed)
            {
                return new RevivePolicyResponse
                {
                    Error = BuildError("INVALID_STATE", "Policy cannot be revived after grace period lapse")
                };
            }

            gracePeriod.Status = GracePeriodStatus.Revived;
            gracePeriod.CoverageActive = true;
            gracePeriod.DaysRemaining = 0;
            gracePeriod.RevivedAt = Timestamp.FromDateTime(DateTime.UtcNow);
            await _dataGateway.UpdateGracePeriodAsync(gracePeriod, GetCancellationToken(context));

            if (schedule.Status != RenewalStatus.Renewed)
            {
                schedule.Status = RenewalStatus.Renewed;
                schedule.RenewedPolicyId = string.IsNullOrWhiteSpace(schedule.RenewedPolicyId)
                    ? $"POL-{Guid.NewGuid():N}"[..16]
                    : schedule.RenewedPolicyId;
                schedule.RenewedAt = Timestamp.FromDateTime(DateTime.UtcNow);
                await _dataGateway.UpdateRenewalScheduleAsync(schedule, GetCancellationToken(context));
            }

            return new RevivePolicyResponse
            {
                Message = "Policy revived successfully"
            };
        }
        catch (RpcException ex)
        {
            _logger.LogError(ex, "Failed to revive policy {PolicyId}", request.PolicyId);
            return new RevivePolicyResponse
            {
                Error = BuildError("UPSTREAM_ERROR", ex.Status.Detail)
            };
        }
    }

    private async Task<RenewalSchedule> EnsureScheduleAsync(string policyId, CancellationToken cancellationToken)
    {
        var existingSchedules = await _dataGateway.ListRenewalSchedulesAsync(policyId, cancellationToken);
        var existingSchedule = existingSchedules
            .OrderByDescending(x => x.RenewalDueDate?.Seconds ?? 0)
            .FirstOrDefault();

        if (existingSchedule is not null)
        {
            return existingSchedule;
        }

        var dueDate = DateTime.UtcNow.AddDays(30);
        var schedule = new RenewalSchedule
        {
            Id = Guid.NewGuid().ToString("N"),
            PolicyId = policyId,
            RenewalDueDate = Timestamp.FromDateTime(dueDate),
            RenewalPremium = NewMoney(250_000),
            RenewalType = RenewalType.Manual,
            Status = RenewalStatus.Pending,
            GracePeriodDays = 30,
            GracePeriodEnd = Timestamp.FromDateTime(dueDate.AddDays(30))
        };

        var createdSchedule = await _dataGateway.CreateRenewalScheduleAsync(schedule, cancellationToken);

        var gracePeriod = await _dataGateway.GetGracePeriodByPolicyAsync(policyId, cancellationToken);
        if (gracePeriod is null)
        {
            await _dataGateway.CreateGracePeriodAsync(new GracePeriod
            {
                Id = Guid.NewGuid().ToString("N"),
                PolicyId = policyId,
                StartDate = Timestamp.FromDateTime(dueDate),
                EndDate = Timestamp.FromDateTime(dueDate.AddDays(30)),
                DaysRemaining = 30,
                Status = GracePeriodStatus.Active,
                CoverageActive = true
            }, cancellationToken);
        }

        return createdSchedule;
    }

    private static RenewalStatus ParseRenewalStatus(string value)
    {
        return ParseEnum(value, RenewalStatus.Unspecified);
    }

    private static ReminderChannel ParseReminderChannel(string value)
    {
        return ParseEnum(value, ReminderChannel.Sms);
    }

    private static TEnum ParseEnum<TEnum>(string value, TEnum fallback) where TEnum : struct, System.Enum
    {
        if (string.IsNullOrWhiteSpace(value))
        {
            return fallback;
        }

        var token = value.Trim();
        if (System.Enum.TryParse<TEnum>(token, true, out var direct))
        {
            return direct;
        }

        var parts = token.Split('_', StringSplitOptions.RemoveEmptyEntries);
        for (var i = 0; i < parts.Length; i++)
        {
            var candidate = string.Concat(parts.Skip(i).Select(ToPascalCase));
            if (System.Enum.TryParse<TEnum>(candidate, true, out var parsed))
            {
                return parsed;
            }
        }

        return fallback;
    }

    private static string ToPascalCase(string segment)
    {
        if (string.IsNullOrWhiteSpace(segment)) return string.Empty;
        var lower = segment.ToLowerInvariant();
        return char.ToUpperInvariant(lower[0]) + lower[1..];
    }

    private static CancellationToken GetCancellationToken(ServerCallContext? context)
        => context?.CancellationToken ?? CancellationToken.None;

    private static Money NewMoney(long amount, string currency = "BDT")
        => new() { Amount = amount, Currency = currency };

    private static Error BuildError(string code, string message)
        => new() { Code = code, Message = message };
}
