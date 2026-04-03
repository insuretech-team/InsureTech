using MediatR;
using Microsoft.Extensions.Logging;
using Insuretech.Policy.Services.V1;
using Insuretech.Common.V1;
using InsuranceEngine.Cancellations.Infrastructure;

namespace InsuranceEngine.Cancellations.Application.Commands;

public sealed class CancelPolicyCommandHandler : IRequestHandler<CancelPolicyCommand, CancelPolicyResponse>
{
    private readonly ICancellationDataGateway _gateway;
    private readonly ILogger<CancelPolicyCommandHandler> _logger;

    public CancelPolicyCommandHandler(
        ICancellationDataGateway gateway,
        ILogger<CancelPolicyCommandHandler> logger)
    {
        _gateway = gateway;
        _logger = logger;
    }

    public async Task<CancelPolicyResponse> Handle(CancelPolicyCommand request, CancellationToken cancellationToken)
    {
        try
        {
            _logger.LogInformation("Initiating policy cancellation: {PolicyId}", request.PolicyId);

            var grpcRequest = new CancelPolicyRequest
            {
                PolicyId = request.PolicyId,
                Reason = request.Reason
            };

            var response = await _gateway.CancelPolicyAsync(grpcRequest, cancellationToken);
            
            if (response.Error != null)
            {
                _logger.LogWarning("Policy cancellation failed: {ErrorCode} - {ErrorMessage}", response.Error.Code, response.Error.Message);
            }
            else
            {
                _logger.LogInformation("Policy cancellation processed: {PolicyId}. Status: {Message}", request.PolicyId, response.Message);
            }

            return response;
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to cancel policy via gateway: {PolicyId}", request.PolicyId);
            return new CancelPolicyResponse { Error = new Error { Code = "GATEWAY_ERROR", Message = ex.Message } };
        }
    }
}
