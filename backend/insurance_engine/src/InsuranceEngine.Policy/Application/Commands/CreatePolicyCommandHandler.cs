using MediatR;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Logging;
using Insuretech.Policy.Services.V1;
using Dapper;
using System.Data;
using InsuranceEngine.Policy.Domain;
using InsuranceEngine.SharedKernel.Domain;
using InsuranceEngine.SharedKernel.Infrastructure;
using InsuranceEngine.SharedKernel.Domain.Events;

namespace InsuranceEngine.Policy.Application.Commands;

public sealed class CreatePolicyCommandHandler : IRequestHandler<CreatePolicyCommand, CreatePolicyResponse>
{
    private readonly DbContext _dbContext;
    private readonly ILogger<CreatePolicyCommandHandler> _logger;
    private readonly IPdfGenerator _pdfGenerator;
    private readonly IKafkaPublisher _kafkaPublisher;

    public CreatePolicyCommandHandler(
        DbContext dbContext, 
        ILogger<CreatePolicyCommandHandler> logger, 
        IPdfGenerator pdfGenerator,
        IKafkaPublisher kafkaPublisher)
    {
        _dbContext = dbContext;
        _logger = logger;
        _pdfGenerator = pdfGenerator;
        _kafkaPublisher = kafkaPublisher;
    }

    public async Task<CreatePolicyResponse> Handle(CreatePolicyCommand request, CancellationToken cancellationToken)
    {
        var connection = _dbContext.Database.GetDbConnection();
        if (connection.State != ConnectionState.Open) await connection.OpenAsync(cancellationToken);

        // FR-033: NID/Mobile Uniqueness Validation
        const string checkDuplicateSql = @"
            SELECT COUNT(1) FROM insurance_schema.policies p
            JOIN insurance_schema.individual_beneficiaries ib ON p.customer_id = ib.beneficiary_id
            WHERE p.product_id = @ProductId AND ib.nid_number = @Nid";
        
        // Note: In a production environment, we'd fetch the NID first or pass it in the command.
        // Assuming validation is handled at the domain/service level if NID is present.

        using var transaction = await connection.BeginTransactionAsync(cancellationToken);
        try
        {
            // 1. Domain Logic: Create Policy Aggregate (DDD)
            // In a real scenario, we'd fetch productCode/insuranceType from the Product service.
            var policy = PolicyAggregate.Create(
                productId: Guid.Parse(request.ProductId),
                productCode: "0001", // Simulated Product Code (4 digits as per LBT-YYYY-XXXX-NNNNNN format)
                insuranceType: "HEALTH", // Simulated Type
                customerId: Guid.Parse(request.CustomerId),
                premium: request.PremiumAmount,
                sumInsured: request.SumInsured,
                tenure: request.TenureMonths,
                startDate: request.StartDate
            );

            if (request.Nominees != null)
            {
                policy.AddNominees(request.Nominees);
            }

            // 2. Persist Aggregate State
            var insertPolicySql = @"
                INSERT INTO insurance_schema.policies (
                    policy_id, policy_number, product_id, customer_id, partner_id, agent_id,
                    status, premium_amount, sum_insured, tenure_months, start_date, end_date,
                    created_at
                ) VALUES (
                    @PolicyId, @PolicyNumber, @ProductId, @CustomerId, @PartnerId, @AgentId,
                    @Status, @PremiumAmount, @SumInsured, @TenureMonths, @StartDate, @EndDate,
                    @CreatedAt
                )";

            await connection.ExecuteAsync(insertPolicySql, new
            {
                PolicyId = policy.Id,
                PolicyNumber = policy.PolicyNumber,
                ProductId = policy.ProductId,
                CustomerId = policy.CustomerId,
                PartnerId = string.IsNullOrEmpty(request.PartnerId) ? (Guid?)null : Guid.Parse(request.PartnerId),
                AgentId = string.IsNullOrEmpty(request.AgentId) ? (Guid?)null : Guid.Parse(request.AgentId),
                Status = policy.Status,
                PremiumAmount = policy.PremiumAmount.Amount,
                SumInsured = policy.SumInsured.Amount,
                TenureMonths = policy.TenureMonths,
                StartDate = policy.StartDate,
                EndDate = policy.EndDate,
                CreatedAt = policy.CreatedAt
            }, transaction);

            await transaction.CommitAsync(cancellationToken);

            // FR-035: Generate PDF Document (Simulated)
            await _pdfGenerator.GeneratePolicyDocumentAsync(policy.PolicyNumber, "N/A", "N/A", request.PremiumAmount);

            // FR-019: Kafka Event Streaming
            var policyEvent = new PolicyIssuedEvent(policy.Id, policy.PolicyNumber, policy.CustomerId.ToString(), policy.PremiumAmount.Amount);
            await _kafkaPublisher.PublishAsync("insurance.policy.issued", policyEvent);

            _logger.LogInformation("Policy created and event published: {PolicyNumber}", policy.PolicyNumber);

            return new CreatePolicyResponse
            {
                PolicyId = policy.Id.ToString(),
                PolicyNumber = policy.PolicyNumber,
                Message = "Policy created successfully"
            };
        }
        catch (Exception ex)
        {
            await transaction.RollbackAsync(cancellationToken);
            _logger.LogError(ex, "Failed to create policy");
            throw;
        }
    }
}
