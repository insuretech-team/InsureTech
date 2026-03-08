using Grpc.Core;
using Insuretech.Insurance.Services.V1;
using PoliSync.Infrastructure.Clients;
using GracePeriodEntity = Insuretech.Renewal.Entity.V1.GracePeriod;
using RenewalReminderEntity = Insuretech.Renewal.Entity.V1.RenewalReminder;
using RenewalScheduleEntity = Insuretech.Renewal.Entity.V1.RenewalSchedule;

namespace PoliSync.Renewal.Infrastructure;

public sealed class GoRenewalDataGateway : IRenewalDataGateway
{
    private readonly InsuranceServiceClient _insuranceClient;

    public GoRenewalDataGateway(InsuranceServiceClient insuranceClient)
    {
        _insuranceClient = insuranceClient;
    }

    public async Task<RenewalScheduleEntity> CreateRenewalScheduleAsync(RenewalScheduleEntity schedule, CancellationToken cancellationToken = default)
    {
        var response = await _insuranceClient.Client.CreateRenewalScheduleAsync(
            new CreateRenewalScheduleRequest { Schedule = schedule },
            cancellationToken: cancellationToken);

        return response.Schedule;
    }

    public async Task<RenewalScheduleEntity?> GetRenewalScheduleAsync(string scheduleId, CancellationToken cancellationToken = default)
    {
        try
        {
            var response = await _insuranceClient.Client.GetRenewalScheduleAsync(
                new GetRenewalScheduleRequest { ScheduleId = scheduleId },
                cancellationToken: cancellationToken);

            return response.Schedule;
        }
        catch (RpcException ex) when (ex.StatusCode == StatusCode.NotFound)
        {
            return null;
        }
    }

    public async Task<RenewalScheduleEntity> UpdateRenewalScheduleAsync(RenewalScheduleEntity schedule, CancellationToken cancellationToken = default)
    {
        var response = await _insuranceClient.Client.UpdateRenewalScheduleAsync(
            new UpdateRenewalScheduleRequest { Schedule = schedule },
            cancellationToken: cancellationToken);

        return response.Schedule;
    }

    public async Task<IReadOnlyList<RenewalScheduleEntity>> ListRenewalSchedulesAsync(string policyId, CancellationToken cancellationToken = default)
    {
        var response = await _insuranceClient.Client.ListRenewalSchedulesAsync(
            new ListRenewalSchedulesRequest { PolicyId = policyId ?? string.Empty },
            cancellationToken: cancellationToken);

        return response.Schedules;
    }

    public async Task<RenewalReminderEntity> CreateRenewalReminderAsync(RenewalReminderEntity reminder, CancellationToken cancellationToken = default)
    {
        var response = await _insuranceClient.Client.CreateRenewalReminderAsync(
            new CreateRenewalReminderRequest { Reminder = reminder },
            cancellationToken: cancellationToken);

        return response.Reminder;
    }

    public async Task<IReadOnlyList<RenewalReminderEntity>> ListRenewalRemindersAsync(string scheduleId, CancellationToken cancellationToken = default)
    {
        var response = await _insuranceClient.Client.ListRenewalRemindersAsync(
            new ListRenewalRemindersRequest { ScheduleId = scheduleId },
            cancellationToken: cancellationToken);

        return response.Reminders;
    }

    public async Task<GracePeriodEntity> CreateGracePeriodAsync(GracePeriodEntity gracePeriod, CancellationToken cancellationToken = default)
    {
        var response = await _insuranceClient.Client.CreateGracePeriodAsync(
            new CreateGracePeriodRequest { GracePeriod = gracePeriod },
            cancellationToken: cancellationToken);

        return response.GracePeriod;
    }

    public async Task<GracePeriodEntity?> GetGracePeriodByPolicyAsync(string policyId, CancellationToken cancellationToken = default)
    {
        try
        {
            var response = await _insuranceClient.Client.GetGracePeriodByPolicyAsync(
                new GetGracePeriodByPolicyRequest { PolicyId = policyId },
                cancellationToken: cancellationToken);

            if (response.GracePeriod is null || string.IsNullOrWhiteSpace(response.GracePeriod.Id))
            {
                return null;
            }

            return response.GracePeriod;
        }
        catch (RpcException ex) when (ex.StatusCode == StatusCode.NotFound)
        {
            return null;
        }
    }

    public async Task<GracePeriodEntity> UpdateGracePeriodAsync(GracePeriodEntity gracePeriod, CancellationToken cancellationToken = default)
    {
        var response = await _insuranceClient.Client.UpdateGracePeriodAsync(
            new UpdateGracePeriodRequest { GracePeriod = gracePeriod },
            cancellationToken: cancellationToken);

        return response.GracePeriod;
    }
}
