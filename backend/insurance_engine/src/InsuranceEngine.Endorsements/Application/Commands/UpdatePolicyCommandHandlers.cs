using MediatR;
using Microsoft.Extensions.Logging;
using Insuretech.Policy.Services.V1;
using Insuretech.Common.V1;
using InsuranceEngine.Endorsements.Infrastructure;

namespace InsuranceEngine.Endorsements.Application.Commands;

public sealed class UpdatePolicyCommandHandler : IRequestHandler<UpdatePolicyCommand, UpdatePolicyResponse>
{
    private readonly IEndorsementDataGateway _gateway;
    private readonly ILogger<UpdatePolicyCommandHandler> _logger;

    public UpdatePolicyCommandHandler(
        IEndorsementDataGateway gateway,
        ILogger<UpdatePolicyCommandHandler> logger)
    {
        _gateway = gateway;
        _logger = logger;
    }

    public async Task<UpdatePolicyResponse> Handle(UpdatePolicyCommand request, CancellationToken cancellationToken)
    {
        try
        {
            _logger.LogInformation("Updating policy (Endorsement): {PolicyId}", request.PolicyId);

            var currentPolicyResponse = await _gateway.GetPolicyAsync(request.PolicyId, cancellationToken);
            if (currentPolicyResponse.Policy == null)
            {
                return new UpdatePolicyResponse { Error = new Error { Code = "NOT_FOUND", Message = "Policy not found" } };
            }

            var response = await _gateway.UpdatePolicyAsync(request.PolicyId, request.Nominees, cancellationToken);
            
            _logger.LogInformation("Policy endorsement processed successfully for {PolicyId}", request.PolicyId);

            return response;
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to update policy via gateway: {PolicyId}", request.PolicyId);
            return new UpdatePolicyResponse { Error = new Error { Code = "GATEWAY_ERROR", Message = ex.Message } };
        }
    }
}
