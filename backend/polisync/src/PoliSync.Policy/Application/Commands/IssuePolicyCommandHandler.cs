using MediatR;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Logging;
using PoliSync.SharedKernel.CQRS;
using PoliSync.SharedKernel.Messaging;
using PoliSync.Policy.Domain;
using PoliSync.Infrastructure.Persistence;

namespace PoliSync.Policy.Application.Commands;

public sealed class IssuePolicyCommandHandler : IRequestHandler<IssuePolicyCommand, Result<string>>
{
    private readonly PoliSyncDbContext _dbContext;
    private readonly IEventBus _eventBus;
    private readonly ILogger<IssuePolicyCommandHandler> _logger;
    
    public IssuePolicyCommandHandler(
        PoliSyncDbContext dbContext,
        IEventBus eventBus,
        ILogger<IssuePolicyCommandHandler> logger)
    {
        _dbContext = dbContext;
        _eventBus = eventBus;
        _logger = logger;
    }
    
    public async Task<Result<string>> Handle(IssuePolicyCommand request, CancellationToken cancellationToken)
    {
        try
        {
            // Create policy aggregate
            var policyAggregate = PolicyAggregate.Create(
                request.CustomerId,
                request.ProductId,
                request.QuoteId,
                request.PremiumAmountPaisa,
                request.SumInsuredPaisa,
                request.TenureMonths,
                request.StartDate,
                request.EndDate);
            
            // Set optional fields
            if (!string.IsNullOrEmpty(request.PartnerId))
                policyAggregate.Policy.PartnerId = request.PartnerId;
            
            if (!string.IsNullOrEmpty(request.AgentId))
                policyAggregate.Policy.AgentId = request.AgentId;
            
            // Issue the policy
            policyAggregate.IssuePolicy();
            
            // Save to database (using raw SQL for proto entity)
            await SavePolicyToDatabase(policyAggregate.Policy, cancellationToken);
            
            // Publish domain events to Kafka
            foreach (var domainEvent in policyAggregate.DomainEvents)
            {
                await _eventBus.PublishAsync(domainEvent, cancellationToken);
            }
            
            _logger.LogInformation(
                "Policy issued successfully: {PolicyId}, Number: {PolicyNumber}",
                policyAggregate.PolicyId,
                policyAggregate.PolicyNumber);
            
            return Result.Ok(policyAggregate.PolicyId);
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to issue policy for customer {CustomerId}", request.CustomerId);
            return Result.Fail<string>("POLICY_ISSUANCE_FAILED", ex.Message);
        }
    }
    
    private async Task SavePolicyToDatabase(Insuretech.Policy.Entity.V1.Policy policy, CancellationToken cancellationToken)
    {
        // Insert policy using raw SQL since we're using proto entities directly
        var sql = @"
            INSERT INTO insurance_schema.policies (
                policy_id, policy_number, product_id, customer_id, partner_id, agent_id,
                quote_id, status, premium_amount, premium_currency, sum_insured, sum_insured_currency,
                tenure_months, start_date, end_date, issued_at, created_at, updated_at
            ) VALUES (
                @p0, @p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9, @p10, @p11, @p12, @p13, @p14, @p15, @p16, @p17
            )";
        
        await _dbContext.Database.ExecuteSqlRawAsync(sql,
            new object[]
            {
                policy.PolicyId,
                policy.PolicyNumber,
                policy.ProductId,
                policy.CustomerId,
                policy.PartnerId ?? (object)DBNull.Value,
                policy.AgentId ?? (object)DBNull.Value,
                policy.QuoteId,
                policy.Status.ToString(),
                policy.PremiumAmount.Amount,
                policy.PremiumAmount.Currency,
                policy.SumInsured.Amount,
                policy.SumInsured.Currency,
                policy.TenureMonths,
                policy.StartDate.ToDateTime(),
                policy.EndDate.ToDateTime(),
                policy.IssuedAt?.ToDateTime() ?? (object)DBNull.Value,
                policy.CreatedAt.ToDateTime(),
                policy.UpdatedAt?.ToDateTime() ?? (object)DBNull.Value
            },
            cancellationToken);
    }
}
