using Microsoft.Extensions.Logging;
using Insuretech.Commission.Services.V1;
using Insuretech.Common.V1;

namespace InsuranceEngine.Commission.Infrastructure;

public class SqlCommissionDataGateway : ICommissionDataGateway
{
    private readonly ILogger<SqlCommissionDataGateway> _logger;

    public SqlCommissionDataGateway(ILogger<SqlCommissionDataGateway> logger)
    {
        _logger = logger;
    }

    public Task<CalculateCommissionResponse> CalculateCommissionAsync(CalculateCommissionRequest request, CancellationToken ct = default)
    {
        _logger.LogInformation("SQL: Calculating commission for policy {PolicyId}", request.PolicyId);
        return Task.FromResult(new CalculateCommissionResponse { CommissionId = Guid.NewGuid().ToString() });
    }

    public Task<CreatePayoutResponse> CreatePayoutAsync(CreatePayoutRequest request, CancellationToken ct = default)
    {
        _logger.LogInformation("SQL: Creating payout for recipient {RecipientId}", request.RecipientId);
        return Task.FromResult(new CreatePayoutResponse { PayoutId = Guid.NewGuid().ToString(), PayoutNumber = $"PAY-{Guid.NewGuid().ToString()[..8].ToUpper()}" });
    }

    public Task<ProcessPayoutResponse> ProcessPayoutAsync(ProcessPayoutRequest request, CancellationToken ct = default)
    {
        _logger.LogInformation("SQL: Processing payout {PayoutId}", request.PayoutId);
        return Task.FromResult(new ProcessPayoutResponse { Message = "Payout processing" });
    }

    public Task<ListCommissionsResponse> ListCommissionsAsync(ListCommissionsRequest request, CancellationToken ct = default)
    {
        _logger.LogInformation("SQL: Listing commissions for recipient {RecipientId}", request.RecipientId);
        return Task.FromResult(new ListCommissionsResponse());
    }
}
