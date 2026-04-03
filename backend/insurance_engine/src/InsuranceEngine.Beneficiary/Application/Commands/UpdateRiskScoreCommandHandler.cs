using MediatR;
using Microsoft.Extensions.Logging;
using Insuretech.Beneficiary.Services.V1;
using Insuretech.Common.V1;
using InsuranceEngine.Grpc.Gateways;

namespace InsuranceEngine.Beneficiary.Application.Commands;

public sealed class UpdateRiskScoreCommandHandler : IRequestHandler<UpdateRiskScoreCommand, UpdateRiskScoreResponse>
{
    private readonly IBeneficiaryDataGateway _gateway;
    private readonly ILogger<UpdateRiskScoreCommandHandler> _logger;

    public UpdateRiskScoreCommandHandler(
        IBeneficiaryDataGateway gateway,
        ILogger<UpdateRiskScoreCommandHandler> logger)
    {
        _gateway = gateway;
        _logger = logger;
    }

    public async Task<UpdateRiskScoreResponse> Handle(UpdateRiskScoreCommand request, CancellationToken cancellationToken)
    {
        try
        {
            _logger.LogInformation("Updating risk score for beneficiary: {BeneficiaryId}", request.BeneficiaryId);

            var grpcRequest = new UpdateRiskScoreRequest
            {
                BeneficiaryId = request.BeneficiaryId,
                RiskScore = request.RiskScore,
                Reason = request.Reason ?? "Manual update"
            };

            var response = await _gateway.UpdateRiskScoreAsync(grpcRequest, cancellationToken);
            
            if (response.Error != null)
            {
                _logger.LogWarning("Risk score update failed: {ErrorCode} - {ErrorMessage}", response.Error.Code, response.Error.Message);
            }
            else
            {
                _logger.LogInformation("Risk score updated successfully for beneficiary: {BeneficiaryId}", request.BeneficiaryId);
            }

            return response;
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to update risk score via gateway");
            return new UpdateRiskScoreResponse { Error = new Error { Code = "GATEWAY_ERROR", Message = ex.Message } };
        }
    }
}
