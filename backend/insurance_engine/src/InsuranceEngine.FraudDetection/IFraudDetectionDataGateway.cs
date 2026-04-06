using InsuranceEngine.SharedKernel.Persistence.Entities;

namespace InsuranceEngine.FraudDetection;

public interface IFraudDetectionDataGateway
{
    Task<FraudCheckEntity> CreateFraudCheckAsync(FraudCheckEntity fraudCheck, CancellationToken ct = default);
    Task<List<FraudCheckEntity>> GetFraudChecksByEntityAsync(string entityId, CancellationToken ct = default);
    Task<FraudAlertEntity?> GetAlertByIdAsync(string alertId, CancellationToken ct = default);
    Task<FraudAlertEntity> CreateAlertAsync(FraudAlertEntity alert, CancellationToken ct = default);
    Task<FraudAlertEntity> UpdateAlertAsync(FraudAlertEntity alert, CancellationToken ct = default);
    Task<List<FraudAlertEntity>> GetPendingAlertsAsync(int limit = 100, CancellationToken ct = default);
    Task<List<FraudCheckEntity>> GetRecentClaimsAsync(string customerId, int months, CancellationToken ct = default);
    Task<FraudDashboardSummaryEntity> GetTodaySummaryAsync(CancellationToken ct = default);
}
