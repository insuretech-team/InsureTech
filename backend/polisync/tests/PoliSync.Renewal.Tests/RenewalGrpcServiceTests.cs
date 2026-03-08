using FluentAssertions;
using Insuretech.Renewal.Entity.V1;
using Insuretech.Renewal.Services.V1;
using MediatR;
using Microsoft.Extensions.Logging.Abstractions;
using Moq;
using PoliSync.Renewal.GrpcServices;
using PoliSync.Renewal.Infrastructure;
using GracePeriodEntity = Insuretech.Renewal.Entity.V1.GracePeriod;
using RenewalReminderEntity = Insuretech.Renewal.Entity.V1.RenewalReminder;
using RenewalScheduleEntity = Insuretech.Renewal.Entity.V1.RenewalSchedule;
using Xunit;

namespace PoliSync.Renewal.Tests;

public class RenewalGrpcServiceTests
{
    private sealed class FakeRenewalDataGateway : IRenewalDataGateway
    {
        private readonly Dictionary<string, RenewalScheduleEntity> _schedules = new();
        private readonly Dictionary<string, RenewalReminderEntity> _reminders = new();
        private readonly Dictionary<string, GracePeriodEntity> _gracePeriodsByPolicy = new();

        public Task<RenewalScheduleEntity> CreateRenewalScheduleAsync(RenewalScheduleEntity schedule, CancellationToken cancellationToken = default)
        {
            _schedules[schedule.Id] = schedule;
            return Task.FromResult(schedule);
        }

        public Task<RenewalScheduleEntity?> GetRenewalScheduleAsync(string scheduleId, CancellationToken cancellationToken = default)
            => Task.FromResult(_schedules.TryGetValue(scheduleId, out var schedule) ? schedule : null);

        public Task<RenewalScheduleEntity> UpdateRenewalScheduleAsync(RenewalScheduleEntity schedule, CancellationToken cancellationToken = default)
        {
            _schedules[schedule.Id] = schedule;
            return Task.FromResult(schedule);
        }

        public Task<IReadOnlyList<RenewalScheduleEntity>> ListRenewalSchedulesAsync(string policyId, CancellationToken cancellationToken = default)
        {
            var items = _schedules.Values
                .Where(x => string.IsNullOrWhiteSpace(policyId) || x.PolicyId == policyId)
                .ToList();
            return Task.FromResult<IReadOnlyList<RenewalScheduleEntity>>(items);
        }

        public Task<RenewalReminderEntity> CreateRenewalReminderAsync(RenewalReminderEntity reminder, CancellationToken cancellationToken = default)
        {
            _reminders[reminder.Id] = reminder;
            return Task.FromResult(reminder);
        }

        public Task<IReadOnlyList<RenewalReminderEntity>> ListRenewalRemindersAsync(string scheduleId, CancellationToken cancellationToken = default)
        {
            var items = _reminders.Values
                .Where(x => x.RenewalScheduleId == scheduleId)
                .ToList();
            return Task.FromResult<IReadOnlyList<RenewalReminderEntity>>(items);
        }

        public Task<GracePeriodEntity> CreateGracePeriodAsync(GracePeriodEntity gracePeriod, CancellationToken cancellationToken = default)
        {
            _gracePeriodsByPolicy[gracePeriod.PolicyId] = gracePeriod;
            return Task.FromResult(gracePeriod);
        }

        public Task<GracePeriodEntity?> GetGracePeriodByPolicyAsync(string policyId, CancellationToken cancellationToken = default)
            => Task.FromResult(_gracePeriodsByPolicy.TryGetValue(policyId, out var gracePeriod) ? gracePeriod : null);

        public Task<GracePeriodEntity> UpdateGracePeriodAsync(GracePeriodEntity gracePeriod, CancellationToken cancellationToken = default)
        {
            _gracePeriodsByPolicy[gracePeriod.PolicyId] = gracePeriod;
            return Task.FromResult(gracePeriod);
        }
    }

    private static RenewalGrpcService CreateService()
        => new(Mock.Of<IMediator>(), NullLogger<RenewalGrpcService>.Instance, new FakeRenewalDataGateway());

    [Fact]
    public async Task GetScheduleAndSendReminder_ReturnsReminderInSchedule()
    {
        var service = CreateService();
        var policyId = $"pol-{Guid.NewGuid():N}";

        var schedule = await service.GetRenewalSchedule(new GetRenewalScheduleRequest
        {
            PolicyId = policyId
        }, null!);

        var reminder = await service.SendRenewalReminder(new SendRenewalReminderRequest
        {
            RenewalScheduleId = schedule.RenewalSchedule.Id,
            Channel = "REMINDER_CHANNEL_SMS"
        }, null!);

        var scheduleAfterReminder = await service.GetRenewalSchedule(new GetRenewalScheduleRequest
        {
            PolicyId = policyId
        }, null!);

        reminder.ReminderId.Should().NotBeNullOrWhiteSpace();
        scheduleAfterReminder.Reminders.Should().Contain(x => x.Id == reminder.ReminderId);
    }

    [Fact]
    public async Task RenewPolicy_UpdatesGracePeriodToRevived()
    {
        var service = CreateService();
        var policyId = $"pol-{Guid.NewGuid():N}";

        var renewed = await service.RenewPolicy(new RenewPolicyRequest
        {
            PolicyId = policyId,
            PaymentMethod = "card",
            PaymentReference = "pay-100"
        }, null!);

        var grace = await service.GetGracePeriod(new GetGracePeriodRequest
        {
            PolicyId = policyId
        }, null!);

        renewed.NewPolicyId.Should().NotBeNullOrWhiteSpace();
        grace.GracePeriod.Status.Should().Be(GracePeriodStatus.Revived);
        grace.GracePeriod.CoverageActive.Should().BeTrue();
    }
}
