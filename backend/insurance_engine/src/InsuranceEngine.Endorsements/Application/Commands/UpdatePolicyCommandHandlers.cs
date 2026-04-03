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

            // Fetch current policy state to ensure we have the full aggregate
            var currentPolicy = await _gateway.GetPolicyAsync(request.PolicyId, cancellationToken);
            if (currentPolicy == null)
            {
                return new UpdatePolicyResponse { Error = new Error { Code = "NOT_FOUND", Message = "Policy not found" } };
            }

            // Update nominees if provided
            if (request.Nominees != null && request.Nominees.Count > 0)
            {
                currentPolicy.Nominees.Clear();
                currentPolicy.Nominees.AddRange(request.Nominees);
            }

            var grpcRequest = new UpdatePolicyRequest
            {
                Policy = currentPolicy
            };

            var response = await _gateway.UpdatePolicyAsync(grpcRequest.Policy, cancellationToken);
            
            _logger.LogInformation("Policy endorsement processed successfully for {PolicyId}", request.PolicyId);

            return new UpdatePolicyResponse 
            { 
                Message = "Policy updated successfully",
                Policy = response
            };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to update policy via gateway: {PolicyId}", request.PolicyId);
            return new UpdatePolicyResponse { Error = new Error { Code = "GATEWAY_ERROR", Message = ex.Message } };
        }
    }
}
