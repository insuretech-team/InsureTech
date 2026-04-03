using MediatR;
using Microsoft.Extensions.Logging;
using Insuretech.Policy.Services.V1;
using Insuretech.Policy.Entity.V1;
using InsuranceEngine.SharedKernel.Domain.Events;
using InsuranceEngine.SharedKernel.Infrastructure;

namespace InsuranceEngine.Policy.Application.Commands;

// ===== IssuePolicy =====
public sealed class IssuePolicyCommandHandler : IRequestHandler<IssuePolicyCommand, IssuePolicyResponse>
{
    private readonly IPolicyDataGateway _gateway;
    private readonly IKafkaPublisher _kafkaPublisher;
    private readonly ILogger<IssuePolicyCommandHandler> _logger;

    public IssuePolicyCommandHandler(
        IPolicyDataGateway gateway,
        IKafkaPublisher kafkaPublisher,
        ILogger<IssuePolicyCommandHandler> logger)
    {
        _gateway = gateway;
        _kafkaPublisher = kafkaPublisher;
        _logger = logger;
    }

    public async Task<IssuePolicyResponse> Handle(IssuePolicyCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var policyResponse = await _gateway.GetPolicyAsync(request.PolicyId, cancellationToken);
            if (policyResponse.Policy == null)
            {
                return new IssuePolicyResponse
                {
                    Error = new Insuretech.Common.V1.Error { Code = "POLICY_NOT_FOUND", Message = "Policy not found" }
                };
            }

            var policy = policyResponse.Policy;

            if (policy.Status != PolicyStatus.PendingPayment)
            {
                return new IssuePolicyResponse
                {
                    Error = new Insuretech.Common.V1.Error { Code = "INVALID_STATUS", Message = $"Policy cannot be issued from status '{policy.Status}'" }
                };
            }

            var issuedPolicy = policy.Clone();
            issuedPolicy.Status = PolicyStatus.Active;
            issuedPolicy.IssuedAt = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(DateTime.UtcNow);

            var updateResponse = await _gateway.UpdatePolicyAsync(request.PolicyId, issuedPolicy.Nominees.ToList(), null, cancellationToken);
            
            if (updateResponse.Error != null)
            {
                return new IssuePolicyResponse { Error = updateResponse.Error };
            }

            var kafkaPolicyResponse = await _gateway.GetPolicyAsync(request.PolicyId, cancellationToken);
            var kafkaPolicy = kafkaPolicyResponse.Policy!;

            var evt = new PolicyIssuedEvent(
                Guid.Parse(kafkaPolicy.PolicyId), 
                kafkaPolicy.PolicyNumber, 
                Guid.Parse(kafkaPolicy.CustomerId), 
                kafkaPolicy.PremiumAmount.Amount);
            
            await _kafkaPublisher.PublishAsync("insurance.policy.issued", evt);

            _logger.LogInformation("Policy issued via Go SSOT: {PolicyNumber}", kafkaPolicy.PolicyNumber);

            return new IssuePolicyResponse { Policy = kafkaPolicy, Message = "Policy issued successfully" };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to issue policy {PolicyId}", request.PolicyId);
            return new IssuePolicyResponse
            {
                Error = new Insuretech.Common.V1.Error { Code = "ISSUE_FAILED", Message = ex.Message }
            };
        }
    }
}

// ===== GeneratePolicyDocument =====
public sealed class GeneratePolicyDocumentCommandHandler : IRequestHandler<GeneratePolicyDocumentCommand, GeneratePolicyDocumentResponse>
{
    private readonly IPolicyDataGateway _gateway;
    private readonly IPdfGenerator _pdfGenerator;
    private readonly ILogger<GeneratePolicyDocumentCommandHandler> _logger;

    public GeneratePolicyDocumentCommandHandler(
        IPolicyDataGateway gateway,
        IPdfGenerator pdfGenerator,
        ILogger<GeneratePolicyDocumentCommandHandler> logger)
    {
        _gateway = gateway;
        _pdfGenerator = pdfGenerator;
        _logger = logger;
    }

    public async Task<GeneratePolicyDocumentResponse> Handle(GeneratePolicyDocumentCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var policyResponse = await _gateway.GetPolicyAsync(request.PolicyId, cancellationToken);
            if (policyResponse.Policy == null)
            {
                return new GeneratePolicyDocumentResponse
                {
                    Error = new Insuretech.Common.V1.Error { Code = "POLICY_NOT_FOUND", Message = "Policy not found" }
                };
            }

            var policy = policyResponse.Policy;

            var pdfBytes = await _pdfGenerator.GeneratePolicyDocumentAsync(
                policy.PolicyNumber,
                policy.CustomerId,
                policy.ProductId,
                (decimal)policy.PremiumAmount.Amount / 100m
            );

            var documentUrl = $"https://storage.insuretech.labaid.com/policies/{policy.PolicyNumber}.pdf";

            policy.PolicyDocumentUrl = documentUrl;

            _logger.LogInformation("Policy document generated for Go SSOT: {PolicyNumber}", policy.PolicyNumber);

            return new GeneratePolicyDocumentResponse
            {
                DocumentUrl = documentUrl,
                QrCode = $"https://insuretech.labaid.com/verify/{policy.PolicyNumber}"
            };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to generate document for policy {PolicyId}", request.PolicyId);
            return new GeneratePolicyDocumentResponse
            {
                Error = new Insuretech.Common.V1.Error { Code = "DOCUMENT_GENERATION_FAILED", Message = ex.Message }
            };
        }
    }
}
