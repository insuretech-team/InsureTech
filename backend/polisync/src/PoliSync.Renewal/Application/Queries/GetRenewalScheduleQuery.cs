using PoliSync.SharedKernel.CQRS;
using RenewalScheduleEntity = Insuretech.Renewal.Entity.V1.RenewalSchedule;

namespace PoliSync.Renewal.Application.Queries;

public sealed record GetRenewalScheduleQuery(string PolicyId) : IQuery<RenewalScheduleEntity>;
