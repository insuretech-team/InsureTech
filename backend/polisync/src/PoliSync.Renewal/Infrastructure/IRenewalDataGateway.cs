using GracePeriodEntity = Insuretech.Renewal.Entity.V1.GracePeriod;
using RenewalReminderEntity = Insuretech.Renewal.Entity.V1.RenewalReminder;
using RenewalScheduleEntity = Insuretech.Renewal.Entity.V1.RenewalSchedule;

namespace PoliSync.Renewal.Infrastructure;

public interface IRenewalDataGateway
{
    Task<RenewalScheduleEntity> CreateRenewalScheduleAsync(RenewalScheduleEntity schedule, CancellationToken cancellationToken = default);
    Task<RenewalScheduleEntity?> GetRenewalScheduleAsync(string scheduleId, CancellationToken cancellationToken = default);
    Task<RenewalScheduleEntity> UpdateRenewalScheduleAsync(RenewalScheduleEntity schedule, CancellationToken cancellationToken = default);
    Task<IReadOnlyList<RenewalScheduleEntity>> ListRenewalSchedulesAsync(string policyId, CancellationToken cancellationToken = default);

    Task<RenewalReminderEntity> CreateRenewalReminderAsync(RenewalReminderEntity reminder, CancellationToken cancellationToken = default);
    Task<IReadOnlyList<RenewalReminderEntity>> ListRenewalRemindersAsync(string scheduleId, CancellationToken cancellationToken = default);

    Task<GracePeriodEntity> CreateGracePeriodAsync(GracePeriodEntity gracePeriod, CancellationToken cancellationToken = default);
    Task<GracePeriodEntity?> GetGracePeriodByPolicyAsync(string policyId, CancellationToken cancellationToken = default);
    Task<GracePeriodEntity> UpdateGracePeriodAsync(GracePeriodEntity gracePeriod, CancellationToken cancellationToken = default);
}
