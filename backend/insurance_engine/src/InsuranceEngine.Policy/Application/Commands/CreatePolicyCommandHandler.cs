using MediatR;
using Microsoft.Extensions.Logging;
using Insuretech.Policy.Services.V1;
using Insuretech.Policy.Entity.V1;
using InsuranceEngine.Grpc.Gateways;
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
            // Note: Validation (product status, sequences, PDF generation, Kafka events)
            // is now handled by the Go backend (SSOT).
            
            var policy = new Insuretech.Policy.Entity.V1.Policy
            {
                ProductId = request.ProductId,
                CustomerId = request.CustomerId,
                PartnerId = request.PartnerId ?? string.Empty,
                AgentId = request.AgentId ?? string.Empty,
                QuoteId = request.QuoteId ?? string.Empty,
                Status = PolicyStatus.PendingPayment,
                PremiumAmount = new Insuretech.Common.V1.Money { Amount = (long)(request.PremiumAmount * 100), Currency = "BDT" },
                SumInsured = new Insuretech.Common.V1.Money { Amount = (long)(request.SumInsured * 100), Currency = "BDT" },
                TenureMonths = request.TenureMonths,
                StartDate = Timestamp.FromDateTime(request.StartDate.ToUniversalTime()),
                EndDate = Timestamp.FromDateTime(request.StartDate.ToUniversalTime().AddMonths(request.TenureMonths))
            };

            if (request.Nominees != null)
            {
                foreach (var nomineeReq in request.Nominees)
                {
                    policy.Nominees.Add(new Nominee
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

            var createdPolicy = await _gateway.CreatePolicyAsync(policy, cancellationToken);

            _logger.LogInformation("Policy created via Go SSOT: {PolicyNumber} for Customer: {CustomerId}", 
                createdPolicy.PolicyNumber, request.CustomerId);

            return new CreatePolicyResponse
            {
                PolicyId = createdPolicy.PolicyId,
                PolicyNumber = createdPolicy.PolicyNumber,
                Message = "Policy created successfully"
            };
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
