using MediatR;
using Microsoft.Extensions.Logging;
using Insuretech.Policy.Services.V1;
using Insuretech.Policy.Entity.V1;
using Google.Protobuf.WellKnownTypes;

namespace InsuranceEngine.Policy.Application.Commands;

public sealed class CreatePolicyCommandHandler : IRequestHandler<CreatePolicyCommand, CreatePolicyResponse>
{
    private readonly IPolicyDataGateway _gateway;
    private readonly ILogger<CreatePolicyCommandHandler> _logger;

    public CreatePolicyCommandHandler(
        IPolicyDataGateway gateway,
        ILogger<CreatePolicyCommandHandler> logger)
    {
        _gateway = gateway;
        _logger = logger;
    }

    public async Task<CreatePolicyResponse> Handle(CreatePolicyCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var createRequest = new CreatePolicyRequest
            {
                ProductId = request.ProductId,
                CustomerId = request.CustomerId,
                PartnerId = request.PartnerId ?? string.Empty,
                AgentId = request.AgentId ?? string.Empty,
                PremiumAmount = new Insuretech.Common.V1.Money { Amount = (long)(request.PremiumAmount * 100), Currency = "BDT" },
                SumInsured = new Insuretech.Common.V1.Money { Amount = (long)(request.SumInsured * 100), Currency = "BDT" },
                TenureMonths = request.TenureMonths
            };

            if (request.Nominees != null)
            {
                foreach (var nomineeReq in request.Nominees)
                {
                    createRequest.Nominees.Add(new Nominee
                    {
                        FullName = nomineeReq.FullName,
                        Relationship = nomineeReq.Relationship,
                        SharePercentage = (double)nomineeReq.SharePercentage,
                        DateOfBirth = nomineeReq.DateOfBirth,
                        NidNumber = nomineeReq.NidNumber,
                        PhoneNumber = nomineeReq.PhoneNumber
                    });
                }
            }

            var response = await _gateway.CreatePolicyAsync(createRequest, cancellationToken);

            if (response.Error != null)
            {
                _logger.LogError("Policy creation failed: {Error}", response.Error.Message);
                return new CreatePolicyResponse
                {
                    Error = response.Error
                };
            }

            _logger.LogInformation("Policy created via Go SSOT: {PolicyNumber} for Customer: {CustomerId}", 
                response.PolicyNumber, request.CustomerId);

            return response;
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to create policy for customer {CustomerId}", request.CustomerId);
            return new CreatePolicyResponse
            {
                Error = new Insuretech.Common.V1.Error { Code = "POLICY_CREATION_FAILED", Message = ex.Message }
            };
        }
    }
}
