using Insuretech.Commission.Services.V1;
using PoliSync.Infrastructure.Clients;
using PartnerCommission = Insuretech.Partner.Entity.V1.Commission;

namespace PoliSync.Commission.Infrastructure;

/// <summary>
/// Calls the Go commission gRPC service for all commission persistence and calculation.
/// </summary>
public sealed class GoCommissionDataGateway : ICommissionDataGateway
{
    private readonly CommissionServiceGrpcClient _client;

    public GoCommissionDataGateway(CommissionServiceGrpcClient client)
    {
        _client = client;
    }

    public async Task<CalculateCommissionResponse> CalculateCommissionAsync(
        string policyId, string commissionType, string recipientType, string recipientId,
        CancellationToken cancellationToken = default)
    {
        return await _client.Client.CalculateCommissionAsync(new CalculateCommissionRequest
        {
            PolicyId = policyId,
            CommissionType = commissionType,
            RecipientType = recipientType,
            RecipientId = recipientId
        }, cancellationToken: cancellationToken);
    }

    public async Task<PartnerCommission?> GetCommissionAsync(
        string commissionId, CancellationToken cancellationToken = default)
    {
        var response = await _client.Client.GetCommissionAsync(
            new GetCommissionRequest { CommissionId = commissionId },
            cancellationToken: cancellationToken);
        return response.Error is not null && !string.IsNullOrWhiteSpace(response.Error.Code)
            ? null
            : response.Commission;
    }

    public async Task<(IReadOnlyList<PartnerCommission> Items, int TotalCount, long TotalAmount)> ListCommissionsAsync(
        string recipientType, string recipientId, string status,
        string startDate, string endDate, int page, int pageSize,
        CancellationToken cancellationToken = default)
    {
        var response = await _client.Client.ListCommissionsAsync(new ListCommissionsRequest
        {
            RecipientType = recipientType,
            RecipientId = recipientId,
            Status = status,
            StartDate = startDate,
            EndDate = endDate,
            Page = page,
            PageSize = pageSize
        }, cancellationToken: cancellationToken);
        return (response.Commissions, response.TotalCount, response.TotalAmount?.Amount ?? 0);
    }

    public async Task<CreatePayoutResponse> CreatePayoutAsync(
        string recipientType, string recipientId,
        string periodStart, string periodEnd,
        IEnumerable<string> commissionIds,
        CancellationToken cancellationToken = default)
    {
        var req = new CreatePayoutRequest
        {
            RecipientType = recipientType,
            RecipientId = recipientId,
            PeriodStart = periodStart,
            PeriodEnd = periodEnd
        };
        req.CommissionIds.AddRange(commissionIds);
        return await _client.Client.CreatePayoutAsync(req, cancellationToken: cancellationToken);
    }

    public async Task<ProcessPayoutResponse> ProcessPayoutAsync(
        string payoutId, string paymentMethod, string paymentReference,
        CancellationToken cancellationToken = default)
    {
        return await _client.Client.ProcessPayoutAsync(new ProcessPayoutRequest
        {
            PayoutId = payoutId,
            PaymentMethod = paymentMethod,
            PaymentReference = paymentReference
        }, cancellationToken: cancellationToken);
    }

    public async Task<GetCommissionStatementResponse> GetCommissionStatementAsync(
        string recipientId, string periodStart, string periodEnd,
        CancellationToken cancellationToken = default)
    {
        return await _client.Client.GetCommissionStatementAsync(new GetCommissionStatementRequest
        {
            RecipientId = recipientId,
            PeriodStart = periodStart,
            PeriodEnd = periodEnd
        }, cancellationToken: cancellationToken);
    }

    public async Task<GetRevenueShareReportResponse> GetRevenueShareReportAsync(
        string insurerId, string startDate, string endDate,
        CancellationToken cancellationToken = default)
    {
        return await _client.Client.GetRevenueShareReportAsync(new GetRevenueShareReportRequest
        {
            InsurerId = insurerId,
            StartDate = startDate,
            EndDate = endDate
        }, cancellationToken: cancellationToken);
    }
}
