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
            var evt = new PolicyIssuedEvent(policy.PolicyId, policy.PolicyNumber, policy.CustomerId.ToString(), policy.PremiumAmount);
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

// ===== CancelPolicy =====
public sealed class CancelPolicyCommandHandler : IRequestHandler<CancelPolicyCommand, CancelPolicyResponse>
{
    private readonly IRepository<PolicyEntity> _repository;
    private readonly IKafkaPublisher _kafkaPublisher;
    private readonly ILogger<CancelPolicyCommandHandler> _logger;

    public CancelPolicyCommandHandler(
        IRepository<PolicyEntity> repository,
        IKafkaPublisher kafkaPublisher,
        ILogger<CancelPolicyCommandHandler> logger)
    {
        _repository = repository;
        _kafkaPublisher = kafkaPublisher;
        _logger = logger;
    }

    public async Task<CancelPolicyResponse> Handle(CancelPolicyCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var policy = await _repository.GetByIdAsync(Guid.Parse(request.PolicyId), cancellationToken);
            if (policy == null)
            {
                return new CancelPolicyResponse
                {
                    Error = new Error { Code = "POLICY_NOT_FOUND", Message = "Policy not found" }
                };
            }

            if (policy.Status == "CANCELLED" || policy.Status == "EXPIRED")
            {
                return new CancelPolicyResponse
                {
                    Error = new Error { Code = "INVALID_STATUS", Message = $"Policy cannot be cancelled from status '{policy.Status}'" }
                };
            }

            // FR-038: Cooling-off period — 5 days from issuance for full refund
            long refundAmount = 0;
            if (policy.IssuedAt.HasValue)
            {
                var daysSinceIssuance = (DateTime.UtcNow - policy.IssuedAt.Value).TotalDays;
                if (daysSinceIssuance <= 5)
                {
                    // Full refund within cooling-off period
                    refundAmount = policy.PremiumAmount;
                }
                else
                {
                    // Proportional refund based on remaining tenure
                    var totalDays = (policy.EndDate - policy.StartDate).TotalDays;
                    var usedDays = (DateTime.UtcNow - policy.StartDate).TotalDays;
                    var remainingFraction = Math.Max(0, (totalDays - usedDays) / totalDays);
                    refundAmount = (long)(policy.PremiumAmount * remainingFraction * 0.9); // 10% admin fee
                }
            }

            policy.Status = "CANCELLED";
            policy.UpdatedAt = DateTime.UtcNow;
            await _repository.UpdateAsync(policy, cancellationToken);

            await _kafkaPublisher.PublishAsync("insurance.policy.cancelled", new { PolicyId = policy.PolicyId, Reason = request.Reason, RefundAmount = refundAmount });

            _logger.LogInformation("Policy cancelled: {PolicyNumber}, Refund: {RefundAmount}", policy.PolicyNumber, refundAmount);

            return new CancelPolicyResponse
            {
                Message = "Policy cancelled successfully",
                RefundAmount = new Money { Amount = refundAmount, Currency = "BDT" }
            };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to cancel policy {PolicyId}", request.PolicyId);
            return new CancelPolicyResponse
            {
                Error = new Error { Code = "CANCEL_FAILED", Message = ex.Message }
            };
        }
    }
}

// ===== RenewPolicy =====
public sealed class RenewPolicyCommandHandler : IRequestHandler<RenewPolicyCommand, RenewPolicyResponse>
{
    private readonly IRepository<PolicyEntity> _repository;
    private readonly IRepository<PolicyNomineeEntity> _nomineeRepository;
    private readonly InsuranceDbContext _dbContext;
    private readonly ILogger<RenewPolicyCommandHandler> _logger;

    public RenewPolicyCommandHandler(
        IRepository<PolicyEntity> repository,
        IRepository<PolicyNomineeEntity> nomineeRepository,
        InsuranceDbContext dbContext,
        ILogger<RenewPolicyCommandHandler> logger)
    {
        _repository = repository;
        _nomineeRepository = nomineeRepository;
        _dbContext = dbContext;
        _logger = logger;
    }

    public async Task<RenewPolicyResponse> Handle(RenewPolicyCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var existingPolicy = await _repository.GetByIdAsync(Guid.Parse(request.PolicyId), cancellationToken);
            if (existingPolicy == null)
            {
                return new RenewPolicyResponse
                {
                    Error = new Error { Code = "POLICY_NOT_FOUND", Message = "Policy not found" }
                };
            }

            if (existingPolicy.Status != "ACTIVE" && existingPolicy.Status != "EXPIRED")
            {
                return new RenewPolicyResponse
                {
                    Error = new Error { Code = "INVALID_STATUS", Message = $"Policy cannot be renewed from status '{existingPolicy.Status}'" }
                };
            }

            // Get new sequence number
            var connection = _dbContext.Database.GetDbConnection();
            if (connection.State != System.Data.ConnectionState.Open)
                await connection.OpenAsync(cancellationToken);

            using var cmd = connection.CreateCommand();
            cmd.CommandText = "SELECT nextval('insurance_schema.policy_number_seq')";
            var seqResult = await cmd.ExecuteScalarAsync(cancellationToken);
            var sequenceNumber = Convert.ToInt64(seqResult);

            // Derive product code for new policy number
            var productCode = existingPolicy.PolicyNumber.Split('-').Length >= 3 
                ? existingPolicy.PolicyNumber.Split('-')[2] 
                : "0001";

            var policyNumber = Domain.ValueObjects.PolicyNumber.Generate(productCode, sequenceNumber).Value;

            // Create renewed policy
            var newStartDate = existingPolicy.EndDate > DateTime.UtcNow ? existingPolicy.EndDate : DateTime.UtcNow;
            var renewedPolicy = new PolicyEntity
            {
                PolicyId = Guid.NewGuid(),
                PolicyNumber = policyNumber,
                ProductId = existingPolicy.ProductId,
                CustomerId = existingPolicy.CustomerId,
                PartnerId = existingPolicy.PartnerId,
                AgentId = existingPolicy.AgentId,
                Status = "PENDING_PAYMENT",
                PremiumAmount = existingPolicy.PremiumAmount,
                PremiumCurrency = existingPolicy.PremiumCurrency,
                SumInsured = existingPolicy.SumInsured,
                SumInsuredCurrency = existingPolicy.SumInsuredCurrency,
                TenureMonths = request.TenureMonths,
                StartDate = newStartDate,
                EndDate = newStartDate.AddMonths(request.TenureMonths),
                CreatedAt = DateTime.UtcNow,
                UpdatedAt = DateTime.UtcNow
            };

            await _repository.AddAsync(renewedPolicy, cancellationToken);

            // Copy nominees from existing policy or use new ones
            if (request.UpdateNominees && request.Nominees != null)
            {
                foreach (var nominee in request.Nominees)
                {
                    await _nomineeRepository.AddAsync(new PolicyNomineeEntity
                    {
                        NomineeId = Guid.NewGuid(),
                        PolicyId = renewedPolicy.PolicyId,
                        FullName = nominee.FullName,
                        Relationship = nominee.Relationship,
                        SharePercentage = nominee.SharePercentage,
                        DateOfBirth = nominee.DateOfBirth?.ToDateTime() ?? DateTime.UtcNow,
                        NidNumber = nominee.NidNumber,
                        PhoneNumber = nominee.PhoneNumber,
                        CreatedAt = DateTime.UtcNow,
                        UpdatedAt = DateTime.UtcNow
                    }, cancellationToken);
                }
            }
            else
            {
                // Copy existing nominees
                var existingNominees = await _nomineeRepository.FindAsync(n => n.PolicyId == existingPolicy.PolicyId, cancellationToken);
                foreach (var nominee in existingNominees)
                {
                    await _nomineeRepository.AddAsync(new PolicyNomineeEntity
                    {
                        NomineeId = Guid.NewGuid(),
                        PolicyId = renewedPolicy.PolicyId,
                        FullName = nominee.FullName,
                        Relationship = nominee.Relationship,
                        SharePercentage = nominee.SharePercentage,
                        DateOfBirth = nominee.DateOfBirth,
                        NidNumber = nominee.NidNumber,
                        PhoneNumber = nominee.PhoneNumber,
                        CreatedAt = DateTime.UtcNow,
                        UpdatedAt = DateTime.UtcNow
                    }, cancellationToken);
                }
            }

            _logger.LogInformation("Policy renewed: {OldPolicyNumber} → {NewPolicyNumber}", existingPolicy.PolicyNumber, renewedPolicy.PolicyNumber);

            return new RenewPolicyResponse
            {
                NewPolicyId = renewedPolicy.PolicyId.ToString(),
                NewPolicyNumber = renewedPolicy.PolicyNumber,
                PremiumAmount = new Money { Amount = renewedPolicy.PremiumAmount, Currency = "BDT" },
                Message = "Policy renewed successfully"
            };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to renew policy {PolicyId}", request.PolicyId);
            return new RenewPolicyResponse
            {
                Error = new Error { Code = "RENEW_FAILED", Message = ex.Message }
            };
        }
    }
}

// ===== UpdatePolicy =====
public sealed class UpdatePolicyCommandHandler : IRequestHandler<UpdatePolicyCommand, UpdatePolicyResponse>
{
    private readonly IRepository<PolicyEntity> _repository;
    private readonly IRepository<PolicyNomineeEntity> _nomineeRepository;
    private readonly ILogger<UpdatePolicyCommandHandler> _logger;

    public UpdatePolicyCommandHandler(
        IRepository<PolicyEntity> repository,
        IRepository<PolicyNomineeEntity> nomineeRepository,
        ILogger<UpdatePolicyCommandHandler> logger)
    {
        _repository = repository;
        _nomineeRepository = nomineeRepository;
        _logger = logger;
    }

    public async Task<UpdatePolicyResponse> Handle(UpdatePolicyCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var policy = await _repository.GetByIdAsync(Guid.Parse(request.PolicyId), cancellationToken);
            if (policy == null)
            {
                return new UpdatePolicyResponse
                {
                    Error = new Error { Code = "POLICY_NOT_FOUND", Message = "Policy not found" }
                };
            }

            // Update nominees if provided
            if (request.Nominees != null && request.Nominees.Count > 0)
            {
                // Remove existing nominees
                var existingNominees = await _nomineeRepository.FindAsync(n => n.PolicyId == policy.PolicyId, cancellationToken);
                foreach (var existing in existingNominees)
                {
                    await _nomineeRepository.DeleteAsync(existing, cancellationToken);
                }

                // Add new nominees
                foreach (var nominee in request.Nominees)
                {
                    await _nomineeRepository.AddAsync(new PolicyNomineeEntity
                    {
                        NomineeId = Guid.NewGuid(),
                        PolicyId = policy.PolicyId,
                        FullName = nominee.FullName,
                        Relationship = nominee.Relationship,
                        SharePercentage = nominee.SharePercentage,
                        DateOfBirth = nominee.DateOfBirth?.ToDateTime() ?? DateTime.UtcNow,
                        NidNumber = nominee.NidNumber,
                        PhoneNumber = nominee.PhoneNumber,
                        CreatedAt = DateTime.UtcNow,
                        UpdatedAt = DateTime.UtcNow
                    }, cancellationToken);
                }
            }

            policy.UpdatedAt = DateTime.UtcNow;
            await _repository.UpdateAsync(policy, cancellationToken);

            _logger.LogInformation("Policy updated: {PolicyId}", request.PolicyId);

            return new UpdatePolicyResponse { Message = "Policy updated successfully" };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to update policy {PolicyId}", request.PolicyId);
            return new UpdatePolicyResponse
            {
                Error = new Error { Code = "UPDATE_FAILED", Message = ex.Message }
            };
        }
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
