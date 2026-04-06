using MediatR;
using Microsoft.Extensions.Logging;
using Insuretech.Fraud.Services.V1;

namespace InsuranceEngine.FraudDetection.Application.Commands;

public sealed class CheckFraudCommand : IRequest<CheckFraudResponse>
{
    public string EntityType { get; set; } = string.Empty;
    public string EntityId { get; set; } = string.Empty;
    public string? Data { get; set; }

    public CheckFraudCommand(string entityType, string entityId, string? data = null)
    {
        EntityType = entityType;
        EntityId = entityId;
        Data = data;
    }
}

public class CheckFraudCommandHandler : IRequestHandler<CheckFraudCommand, CheckFraudResponse>
{
    private readonly IFraudDetectionService _fraudService;
    private readonly ILogger<CheckFraudCommandHandler> _logger;

    public CheckFraudCommandHandler(
        IFraudDetectionService fraudService,
        ILogger<CheckFraudCommandHandler> logger)
    {
        _fraudService = fraudService;
        _logger = logger;
    }

    public async Task<CheckFraudResponse> Handle(CheckFraudCommand request, CancellationToken cancellationToken)
    {
        try
        {
            _logger.LogInformation("Performing fraud check for entity: {EntityId} ({EntityType})", request.EntityId, request.EntityType);

            var fraudRequest = new FraudCheckRequest
            {
                EntityId = request.EntityId,
                EntityType = request.EntityType
            };

            var result = await _fraudService.CheckForFraudAsync(fraudRequest, cancellationToken);
            
            _logger.LogInformation("Fraud check completed for {EntityId}: Detected: {Detected}, Score: {Score}, Risk: {Risk}", 
                request.EntityId, result.IsFraudDetected, result.FraudScore, result.RiskLevel);

            return new CheckFraudResponse
            {
                IsFraudDetected = result.IsFraudDetected,
                FraudScore = result.FraudScore,
                RiskLevel = result.RiskLevel
            };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to perform fraud check");
            return new CheckFraudResponse
            {
                Error = new Insuretech.Common.V1.Error { Code = "INTERNAL_ERROR", Message = ex.Message }
            };
        }
    }
}
