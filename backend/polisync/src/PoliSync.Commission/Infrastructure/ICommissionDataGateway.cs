using Insuretech.Commission.Services.V1;
using PartnerCommission = Insuretech.Partner.Entity.V1.Commission;
using CommissionPayout = Insuretech.Commission.Entity.V1.CommissionPayout;

namespace PoliSync.Commission.Infrastructure;

/// <summary>
/// Abstraction over the Go commission service gRPC calls.
/// </summary>
public interface ICommissionDataGateway
{
    Task<CalculateCommissionResponse> CalculateCommissionAsync(
        string policyId, string commissionType, string recipientType, string recipientId,
        CancellationToken cancellationToken = default);

    Task<PartnerCommission?> GetCommissionAsync(string commissionId, CancellationToken cancellationToken = default);

    Task<(IReadOnlyList<PartnerCommission> Items, int TotalCount, long TotalAmount)> ListCommissionsAsync(
        string recipientType, string recipientId, string status,
        string startDate, string endDate, int page, int pageSize,
        CancellationToken cancellationToken = default);

    Task<CreatePayoutResponse> CreatePayoutAsync(
        string recipientType, string recipientId,
        string periodStart, string periodEnd,
        IEnumerable<string> commissionIds,
        CancellationToken cancellationToken = default);

    Task<ProcessPayoutResponse> ProcessPayoutAsync(
        string payoutId, string paymentMethod, string paymentReference,
        CancellationToken cancellationToken = default);

    Task<GetCommissionStatementResponse> GetCommissionStatementAsync(
        string recipientId, string periodStart, string periodEnd,
        CancellationToken cancellationToken = default);

    Task<GetRevenueShareReportResponse> GetRevenueShareReportAsync(
        string insurerId, string startDate, string endDate,
        CancellationToken cancellationToken = default);
}
