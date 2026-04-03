using Microsoft.Extensions.Logging;
using Insuretech.Commission.Services.V1;
using InsuranceEngine.Grpc.Clients;

namespace InsuranceEngine.Commission.Infrastructure;

/// <summary>
/// Implementation of ICommissionDataGateway using gRPC calls to the Go backend.
/// </summary>
public sealed class GoCommissionDataGateway : ICommissionDataGateway
{
    private readonly InsuranceServiceClient _client;
    private readonly ILogger<GoCommissionDataGateway> _logger;

    public GoCommissionDataGateway(
        InsuranceServiceClient client,
        ILogger<GoCommissionDataGateway> logger)
    {
        _client = client;
        _logger = logger;
    }

    public async Task<CalculateCommissionResponse> CalculateCommissionAsync(CalculateCommissionRequest request, CancellationToken ct = default)
    {
        try
        {
            return await _client.Commissions.CalculateCommissionAsync(request, _client.BuildCallOptions(ct));
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "gRPC Error: CalculateCommission failed for policy {PolicyId}", request.PolicyId);
            throw;
        }
    }

    public async Task<CreatePayoutResponse> CreatePayoutAsync(CreatePayoutRequest request, CancellationToken ct = default)
    {
        try
        {
            return await _client.Commissions.CreatePayoutAsync(request, _client.BuildCallOptions(ct));
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "gRPC Error: CreatePayout failed for recipient {RecipientId}", request.RecipientId);
            throw;
        }
    }

    public async Task<ProcessPayoutResponse> ProcessPayoutAsync(ProcessPayoutRequest request, CancellationToken ct = default)
    {
        try
        {
            return await _client.Commissions.ProcessPayoutAsync(request, _client.BuildCallOptions(ct));
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "gRPC Error: ProcessPayout failed for payout {PayoutId}", request.PayoutId);
            throw;
        }
    }

    public async Task<ListCommissionsResponse> ListCommissionsAsync(ListCommissionsRequest request, CancellationToken ct = default)
    {
        try
        {
            return await _client.Commissions.ListCommissionsAsync(request, _client.BuildCallOptions(ct));
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "gRPC Error: ListCommissions failed for recipient {RecipientId}", request.RecipientId);
            throw;
        }
    }
}
