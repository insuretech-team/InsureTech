using MediatR;
using Microsoft.Extensions.Logging;
using Insuretech.Policy.Services.V1;
using Insuretech.Common.V1;
using InsuranceEngine.Renewals.Infrastructure;

namespace InsuranceEngine.Renewals.Application.Commands;

public sealed class RenewPolicyCommandHandler : IRequestHandler<RenewPolicyCommand, RenewPolicyTenureResponse>
{
    private readonly IRenewalDataGateway _gateway;
    private readonly ILogger<RenewPolicyCommandHandler> _logger;

    public RenewPolicyCommandHandler(
        IRenewalDataGateway gateway,
        ILogger<RenewPolicyCommandHandler> logger)
    {
        _gateway = gateway;
        _logger = logger;
    }

    public async Task<RenewPolicyTenureResponse> Handle(RenewPolicyCommand request, CancellationToken cancellationToken)
    {
        try
        {
            _logger.LogInformation("Renewing policy: {PolicyId} for {TenureMonths} months", request.PolicyId, request.TenureMonths);

            var grpcRequest = new RenewPolicyTenureRequest
            {
                PolicyId = request.PolicyId,
                TenureMonths = request.TenureMonths
            };

            // Map nominees if provided
            if (request.UpdateNominees && request.Nominees != null)
            {
                grpcRequest.Nominees.AddRange(request.Nominees);
            }

            var response = await _gateway.RenewPolicyAsync(grpcRequest, cancellationToken);
            
            if (response.Error != null)
            {
                _logger.LogWarning("Policy renewal failed: {ErrorCode} - {ErrorMessage}", response.Error.Code, response.Error.Message);
            }
            else
            {
                _logger.LogInformation("Policy renewed successfully. New Policy ID: {NewPolicyId}, Number: {NewPolicyNumber}", 
                    response.NewPolicyId, response.NewPolicyNumber);
            }

            return response;
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to renew policy via gateway: {PolicyId}", request.PolicyId);
            return new RenewPolicyTenureResponse { Error = new Error { Code = "GATEWAY_ERROR", Message = ex.Message } };
        }
    }
}
