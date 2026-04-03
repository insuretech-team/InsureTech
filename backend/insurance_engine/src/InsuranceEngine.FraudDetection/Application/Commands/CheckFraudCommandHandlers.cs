using MediatR;
using Microsoft.Extensions.Logging;
using Insuretech.Fraud.Services.V1;
using Insuretech.Common.V1;
namespace InsuranceEngine.FraudDetection.Application.Commands;

public sealed class CheckFraudCommandHandler : IRequestHandler<CheckFraudCommand, CheckFraudResponse>
{
    private readonly IFraudDetectionDataGateway _gateway;
    private readonly ILogger<CheckFraudCommandHandler> _logger;

    public CheckFraudCommandHandler(
        IFraudDetectionDataGateway gateway,
        ILogger<CheckFraudCommandHandler> logger)
    {
        _gateway = gateway;
        _logger = logger;
    }

    public async Task<CheckFraudResponse> Handle(CheckFraudCommand request, CancellationToken cancellationToken)
    {
        try
        {
            _logger.LogInformation("Performing fraud check for entity: {EntityId} ({EntityType})", request.EntityId, request.EntityType);

            var grpcRequest = new CheckFraudRequest
            {
                EntityId = request.EntityId,
                EntityType = request.EntityType,
                Data = request.Data
            };

            var response = await _gateway.CheckFraudAsync(grpcRequest, cancellationToken);
            
            if (response.Error != null)
            {
                _logger.LogWarning("Fraud check failed: {ErrorCode} - {ErrorMessage}", response.Error.Code, response.Error.Message);
            }
            else
            {
                _logger.LogInformation("Fraud check completed for {EntityId}: Detected: {Detected}, Score: {Score}, Risk: {Risk}", 
                    request.EntityId, response.IsFraudDetected, response.FraudScore, response.RiskLevel);
            }

            return response;
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to perform fraud check via gateway");
            return new CheckFraudResponse { Error = new Error { Code = "GATEWAY_ERROR", Message = ex.Message } };
        }
    }
}
