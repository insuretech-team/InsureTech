using InsuranceEngine.SharedKernel.Persistence.Entities;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Logging;

namespace InsuranceEngine.FraudDetection.Infrastructure;

public class SqlFraudDetectionDataGateway : IFraudDetectionDataGateway
{
    private readonly FraudDetectionDbContext _context;
    private readonly ILogger<SqlFraudDetectionDataGateway> _logger;

    public SqlFraudDetectionDataGateway(FraudDetectionDbContext context, ILogger<SqlFraudDetectionDataGateway> logger)
    {
        _context = context;
        _logger = logger;
    }

    public async Task<FraudCheckEntity> CreateFraudCheckAsync(FraudCheckEntity fraudCheck, CancellationToken ct = default)
    {
        _context.FraudChecks.Add(fraudCheck);
        await _context.SaveChangesAsync(ct);
        _logger.LogInformation("Created fraud check {FraudCheckId}", fraudCheck.FraudCheckId);
        return fraudCheck;
    }

    public async Task<List<FraudCheckEntity>> GetFraudChecksByEntityAsync(string entityId, CancellationToken ct = default)
    {
        return await _context.FraudChecks
            .Where(f => f.EntityId == entityId)
            .OrderByDescending(f => f.CreatedAt)
            .ToListAsync(ct);
    }

    public async Task<FraudAlertEntity?> GetAlertByIdAsync(string alertId, CancellationToken ct = default)
    {
        return await _context.FraudAlerts
            .FirstOrDefaultAsync(a => a.AlertId == alertId, ct);
    }

    public async Task<FraudAlertEntity> CreateAlertAsync(FraudAlertEntity alert, CancellationToken ct = default)
    {
        _context.FraudAlerts.Add(alert);
        await _context.SaveChangesAsync(ct);
        _logger.LogInformation("Created fraud alert {AlertId}", alert.AlertId);
        return alert;
    }

    public async Task<FraudAlertEntity> UpdateAlertAsync(FraudAlertEntity alert, CancellationToken ct = default)
    {
        alert.UpdatedAt = DateTime.UtcNow;
        _context.FraudAlerts.Update(alert);
        await _context.SaveChangesAsync(ct);
        _logger.LogInformation("Updated fraud alert {AlertId}", alert.AlertId);
        return alert;
    }

    public async Task<List<FraudAlertEntity>> GetPendingAlertsAsync(int limit = 100, CancellationToken ct = default)
    {
        return await _context.FraudAlerts
            .Where(a => a.Status == "OPEN" || a.Status == "UNDER_REVIEW")
            .OrderByDescending(a => a.CreatedAt)
            .Take(limit)
            .ToListAsync(ct);
    }

    public async Task<List<FraudCheckEntity>> GetRecentClaimsAsync(string customerId, int months, CancellationToken ct = default)
    {
        var startDate = DateTime.UtcNow.AddMonths(-months);
        return await _context.FraudChecks
            .Where(f => f.CustomerId == customerId && f.CheckedAt >= startDate)
            .OrderByDescending(f => f.CheckedAt)
            .ToListAsync(ct);
    }

    public async Task<FraudDashboardSummaryEntity> GetTodaySummaryAsync(CancellationToken ct = default)
    {
        var today = DateTime.UtcNow.Date;
        var summary = await _context.FraudDashboardSummaries
            .FirstOrDefaultAsync(s => s.SummaryDate.Date == today, ct);

        if (summary == null)
        {
            summary = new FraudDashboardSummaryEntity
            {
                SummaryId = Guid.NewGuid().ToString(),
                SummaryDate = today,
                TotalFlagsToday = 0,
                HighRiskFlags = 0,
                MediumRiskFlags = 0,
                LowRiskFlags = 0,
                PendingReviewCount = 0,
                ResolvedCount = 0,
                AverageFraudScore = 0,
                GeneratedAt = DateTime.UtcNow
            };
        }

        return summary;
    }
}
