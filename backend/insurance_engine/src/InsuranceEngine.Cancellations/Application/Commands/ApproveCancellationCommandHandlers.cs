using MediatR;
using Microsoft.Extensions.Logging;
using Insuretech.Policy.Services.V1;
using Insuretech.Common.V1;
using InsuranceEngine.Cancellations.Infrastructure;

namespace InsuranceEngine.Cancellations.Application.Commands;

public sealed record ApproveCancellationCommand(
    string PolicyId, 
    string Role, 
    string ApproverId, 
    string? Notes) : IRequest<ApproveCancellationResponse>;

public sealed class ApproveCancellationCommandHandler : IRequestHandler<ApproveCancellationCommand, ApproveCancellationResponse>
{
    private readonly ICancellationDataGateway _gateway;
    private readonly ILogger<ApproveCancellationCommandHandler> _logger;

    public ApproveCancellationCommandHandler(
        ICancellationDataGateway gateway,
        ILogger<ApproveCancellationCommandHandler> logger)
    {
        _gateway = gateway;
        _logger = logger;
    }

    public async Task<ApproveCancellationResponse> Handle(ApproveCancellationCommand request, CancellationToken cancellationToken)
    {
        try
        {
            _logger.LogInformation("Approving cancellation for policy: {PolicyId} by {Role}", request.PolicyId, request.Role);

            var grpcRequest = new ApproveCancellationRequest
            {
                PolicyId = request.PolicyId,
                Role = request.Role,
                ApproverId = request.ApproverId,
                Notes = request.Notes ?? string.Empty
            };

            var response = await _gateway.ApproveCancellationAsync(grpcRequest, cancellationToken);
            
            if (response.Error != null)
            {
                _logger.LogWarning("Policy cancellation approval failed: {ErrorCode} - {ErrorMessage}", response.Error.Code, response.Error.Message);
            }
            else
            {
                _logger.LogInformation("Policy cancellation approval processed for {PolicyId}. New Status: {Status}", request.PolicyId, response.Status);
            }

            return response;
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to approve policy cancellation via gateway: {PolicyId}", request.PolicyId);
            return new ApproveCancellationResponse { Error = new Error { Code = "GATEWAY_ERROR", Message = ex.Message } };
        }
    }
}
