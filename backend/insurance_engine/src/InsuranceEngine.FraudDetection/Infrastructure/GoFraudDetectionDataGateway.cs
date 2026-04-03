using Microsoft.Extensions.Logging;
using Insuretech.Fraud.Services.V1;
using InsuranceEngine.Grpc.Clients;

namespace InsuranceEngine.FraudDetection.Infrastructure;

/// <summary>
/// Implementation of IFraudDetectionDataGateway using gRPC calls to the Go backend.
/// </summary>
public sealed class GoFraudDetectionDataGateway : IFraudDetectionDataGateway
{
    private readonly InsuranceServiceClient _client;
    private readonly ILogger<GoFraudDetectionDataGateway> _logger;

    public GoFraudDetectionDataGateway(
        InsuranceServiceClient client,
        ILogger<GoFraudDetectionDataGateway> logger)
    {
        _client = client;
        _logger = logger;
    }

    public async Task<CheckFraudResponse> CheckFraudAsync(CheckFraudRequest request, CancellationToken ct = default)
    {
        try
        {
            return await _client.Fraud.CheckFraudAsync(request, _client.BuildCallOptions(ct));
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "gRPC Error: CheckFraud failed for entity {EntityId}", request.EntityId);
            throw;
        }
    }
}
