using Google.Protobuf.WellKnownTypes;
using Insuretech.Renewal.Entity.V1;
using MediatR;
using Microsoft.Extensions.Logging;
using PoliSync.Renewal.Infrastructure;
using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Renewal.Application.Queries;

public sealed class GetRenewalScheduleQueryHandler
    : IRequestHandler<GetRenewalScheduleQuery, Result<RenewalSchedule>>
{
    private readonly IRenewalDataGateway _dataGateway;
    private readonly ILogger<GetRenewalScheduleQueryHandler> _logger;

    public GetRenewalScheduleQueryHandler(IRenewalDataGateway dataGateway, ILogger<GetRenewalScheduleQueryHandler> logger)
    {
        _dataGateway = dataGateway;
        _logger = logger;
    }

    public async Task<Result<RenewalSchedule>> Handle(GetRenewalScheduleQuery request, CancellationToken cancellationToken)
    {
        try
        {
            var schedules = await _dataGateway.ListRenewalSchedulesAsync(request.PolicyId, cancellationToken);
            var schedule = schedules.OrderByDescending(s => s.RenewalDueDate?.Seconds ?? 0).FirstOrDefault();

            if (schedule is null)
            {
                // Auto-create a schedule if none exists
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

            return Result.Ok(schedule);
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to get renewal schedule for policy {PolicyId}", request.PolicyId);
            return Result.Fail<RenewalSchedule>("GET_SCHEDULE_FAILED", ex.Message);
        }
    }
}
