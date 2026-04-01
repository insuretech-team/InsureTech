using MediatR;
using Microsoft.Extensions.Logging;
using Insuretech.Policy.Services.V1;
using Insuretech.Common.V1;
using InsuranceEngine.SharedKernel.Persistence;
using InsuranceEngine.SharedKernel.Persistence.Entities;
using InsuranceEngine.SharedKernel.Infrastructure;
using InsuranceEngine.SharedKernel.Domain.Events;
using InsuranceEngine.Policy.Domain;
using Microsoft.EntityFrameworkCore;
using Google.Protobuf.WellKnownTypes;

namespace InsuranceEngine.Policy.Application.Commands;

// ===== IssuePolicy =====
public sealed class IssuePolicyCommandHandler : IRequestHandler<IssuePolicyCommand, IssuePolicyResponse>
{
    private readonly IRepository<PolicyEntity> _repository;
    private readonly IKafkaPublisher _kafkaPublisher;
    private readonly ILogger<IssuePolicyCommandHandler> _logger;

    public IssuePolicyCommandHandler(
        IRepository<PolicyEntity> repository,
        IKafkaPublisher kafkaPublisher,
        ILogger<IssuePolicyCommandHandler> logger)
    {
        _repository = repository;
        _kafkaPublisher = kafkaPublisher;
        _logger = logger;
    }

    public async Task<IssuePolicyResponse> Handle(IssuePolicyCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var policy = await _repository.GetByIdAsync(Guid.Parse(request.PolicyId), cancellationToken);
            if (policy == null)
            {
                return new IssuePolicyResponse
                {
                    Error = new Error { Code = "POLICY_NOT_FOUND", Message = "Policy not found" }
                };
            }

            if (policy.Status != "PENDING_PAYMENT")
            {
                return new IssuePolicyResponse
                {
                    Error = new Error { Code = "INVALID_STATUS", Message = $"Policy cannot be issued from status '{policy.Status}'" }
                };
            }

            // FR-037: Non-Life policy activates immediately upon payment confirmation
            policy.Status = "ACTIVE";
            policy.IssuedAt = DateTime.UtcNow;
            policy.UpdatedAt = DateTime.UtcNow;

            await _repository.UpdateAsync(policy, cancellationToken);

            // Kafka event: PolicyIssued
            var evt = new PolicyIssuedEvent(policy.PolicyId, policy.PolicyNumber, policy.CustomerId, policy.PremiumAmount);
            await _kafkaPublisher.PublishAsync("insurance.policy.issued", evt);

            _logger.LogInformation("Policy issued: {PolicyNumber}", policy.PolicyNumber);

            var protoPolicy = MapToProto(policy);
            return new IssuePolicyResponse { Policy = protoPolicy, Message = "Policy issued successfully" };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to issue policy {PolicyId}", request.PolicyId);
            return new IssuePolicyResponse
            {
                Error = new Error { Code = "ISSUE_FAILED", Message = ex.Message }
            };
        }
    }

    private static Insuretech.Policy.Entity.V1.Policy MapToProto(PolicyEntity e)
    {
        var p = new Insuretech.Policy.Entity.V1.Policy
        {
            PolicyId = e.PolicyId.ToString(),
            PolicyNumber = e.PolicyNumber,
            ProductId = e.ProductId.ToString(),
            CustomerId = e.CustomerId.ToString(),
            PartnerId = e.PartnerId?.ToString() ?? "",
            AgentId = e.AgentId?.ToString() ?? "",
            TenureMonths = e.TenureMonths,
            PremiumAmount = new Money { Amount = e.PremiumAmount, Currency = e.PremiumCurrency },
            SumInsured = new Money { Amount = e.SumInsured, Currency = e.SumInsuredCurrency },
            PolicyDocumentUrl = e.PolicyDocumentUrl ?? ""
        };
        if (System.Enum.TryParse<Insuretech.Policy.Entity.V1.PolicyStatus>(e.Status, true, out var s)) p.Status = s;
        p.StartDate = Timestamp.FromDateTime(DateTime.SpecifyKind(e.StartDate, DateTimeKind.Utc));
        p.EndDate = Timestamp.FromDateTime(DateTime.SpecifyKind(e.EndDate, DateTimeKind.Utc));
        p.CreatedAt = Timestamp.FromDateTime(DateTime.SpecifyKind(e.CreatedAt, DateTimeKind.Utc));
        if (e.IssuedAt.HasValue) p.IssuedAt = Timestamp.FromDateTime(DateTime.SpecifyKind(e.IssuedAt.Value, DateTimeKind.Utc));
        return p;
    }
}





// ===== GeneratePolicyDocument =====
public sealed class GeneratePolicyDocumentCommandHandler : IRequestHandler<GeneratePolicyDocumentCommand, GeneratePolicyDocumentResponse>
{
    private readonly IRepository<PolicyEntity> _repository;
    private readonly IPdfGenerator _pdfGenerator;
    private readonly ILogger<GeneratePolicyDocumentCommandHandler> _logger;

    public GeneratePolicyDocumentCommandHandler(
        IRepository<PolicyEntity> repository,
        IPdfGenerator pdfGenerator,
        ILogger<GeneratePolicyDocumentCommandHandler> logger)
    {
        _repository = repository;
        _pdfGenerator = pdfGenerator;
        _logger = logger;
    }

    public async Task<GeneratePolicyDocumentResponse> Handle(GeneratePolicyDocumentCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var policy = await _repository.GetByIdAsync(Guid.Parse(request.PolicyId), cancellationToken);
            if (policy == null)
            {
                return new GeneratePolicyDocumentResponse
                {
                    Error = new Error { Code = "POLICY_NOT_FOUND", Message = "Policy not found" }
                };
            }

            // FR-035: Generate PDF with QR code (within 30s of payment confirmation)
            var pdfBytes = await _pdfGenerator.GeneratePolicyDocumentAsync(
                policy.PolicyNumber,
                policy.CustomerId.ToString(),
                policy.ProductId.ToString(),
                (decimal)policy.PremiumAmount / 100m
            );

            // In production, upload pdfBytes to S3 and get URL. Simulated here.
            var documentUrl = $"https://storage.insuretech.labaid.com/policies/{policy.PolicyNumber}.pdf";

            // Generate QR code (simulated — in production, use a QR library)
            var qrCode = $"https://insuretech.labaid.com/verify/{policy.PolicyNumber}";

            // Update policy with document URL
            policy.PolicyDocumentUrl = documentUrl;
            policy.UpdatedAt = DateTime.UtcNow;
            await _repository.UpdateAsync(policy, cancellationToken);

            _logger.LogInformation("Policy document generated ({Bytes} bytes): {PolicyNumber}", pdfBytes.Length, policy.PolicyNumber);

            return new GeneratePolicyDocumentResponse
            {
                DocumentUrl = documentUrl,
                QrCode = qrCode
            };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to generate document for policy {PolicyId}", request.PolicyId);
            return new GeneratePolicyDocumentResponse
            {
                Error = new Error { Code = "DOCUMENT_GENERATION_FAILED", Message = ex.Message }
            };
        }
    }
}
