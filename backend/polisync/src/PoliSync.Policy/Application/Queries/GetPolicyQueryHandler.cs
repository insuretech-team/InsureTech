using MediatR;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Logging;
using PoliSync.SharedKernel.CQRS;
using PoliSync.Infrastructure.Persistence;
using Insuretech.Policy.Entity.V1;

namespace PoliSync.Policy.Application.Queries;

public sealed class GetPolicyQueryHandler : IRequestHandler<GetPolicyQuery, Result<Insuretech.Policy.Entity.V1.Policy>>
{
    private readonly PoliSyncDbContext _dbContext;
    private readonly ILogger<GetPolicyQueryHandler> _logger;
    
    public GetPolicyQueryHandler(
        PoliSyncDbContext dbContext,
        ILogger<GetPolicyQueryHandler> logger)
    {
        _dbContext = dbContext;
        _logger = logger;
    }
    
    public async Task<Result<Insuretech.Policy.Entity.V1.Policy>> Handle(
        GetPolicyQuery request,
        CancellationToken cancellationToken)
    {
        try
        {
            // Query from database using raw SQL
            var sql = @"
                SELECT 
                    policy_id, policy_number, product_id, customer_id, partner_id, agent_id,
                    quote_id, status, premium_amount, premium_currency, sum_insured, sum_insured_currency,
                    tenure_months, start_date, end_date, issued_at, policy_document_url,
                    created_at, updated_at, deleted_at
                FROM insurance_schema.policies
                WHERE policy_id = {0} AND deleted_at IS NULL";
            
            var result = await _dbContext.Database
                .SqlQueryRaw<PolicyDto>(sql, request.PolicyId)
                .FirstOrDefaultAsync(cancellationToken);
            
            if (result == null)
            {
                return Result.Fail<Insuretech.Policy.Entity.V1.Policy>(
                    "POLICY_NOT_FOUND",
                    $"Policy with ID {request.PolicyId} not found");
            }
            
            // Map to proto entity
            var policy = new Insuretech.Policy.Entity.V1.Policy
            {
                PolicyId = result.PolicyId,
                PolicyNumber = result.PolicyNumber,
                ProductId = result.ProductId,
                CustomerId = result.CustomerId,
                PartnerId = result.PartnerId ?? string.Empty,
                AgentId = result.AgentId ?? string.Empty,
                QuoteId = result.QuoteId,
                Status = Enum.Parse<PolicyStatus>(result.Status),
                PremiumAmount = new Insuretech.Common.V1.Money
                {
                    Amount = result.PremiumAmount,
                    Currency = result.PremiumCurrency
                },
                SumInsured = new Insuretech.Common.V1.Money
                {
                    Amount = result.SumInsured,
                    Currency = result.SumInsuredCurrency
                },
                TenureMonths = result.TenureMonths,
                StartDate = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(result.StartDate.ToUniversalTime()),
                EndDate = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(result.EndDate.ToUniversalTime()),
                IssuedAt = result.IssuedAt.HasValue 
                    ? Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(result.IssuedAt.Value.ToUniversalTime())
                    : null,
                PolicyDocumentUrl = result.PolicyDocumentUrl ?? string.Empty,
                CreatedAt = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(result.CreatedAt.ToUniversalTime()),
                UpdatedAt = result.UpdatedAt.HasValue
                    ? Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(result.UpdatedAt.Value.ToUniversalTime())
                    : null
            };
            
            return Result.Ok(policy);
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to get policy {PolicyId}", request.PolicyId);
            return Result.Fail<Insuretech.Policy.Entity.V1.Policy>("QUERY_FAILED", ex.Message);
        }
    }
}

// DTO for database mapping
internal sealed class PolicyDto
{
    public string PolicyId { get; set; } = string.Empty;
    public string PolicyNumber { get; set; } = string.Empty;
    public string ProductId { get; set; } = string.Empty;
    public string CustomerId { get; set; } = string.Empty;
    public string? PartnerId { get; set; }
    public string? AgentId { get; set; }
    public string QuoteId { get; set; } = string.Empty;
    public string Status { get; set; } = string.Empty;
    public long PremiumAmount { get; set; }
    public string PremiumCurrency { get; set; } = "BDT";
    public long SumInsured { get; set; }
    public string SumInsuredCurrency { get; set; } = "BDT";
    public int TenureMonths { get; set; }
    public DateTime StartDate { get; set; }
    public DateTime EndDate { get; set; }
    public DateTime? IssuedAt { get; set; }
    public string? PolicyDocumentUrl { get; set; }
    public DateTime CreatedAt { get; set; }
    public DateTime? UpdatedAt { get; set; }
    public DateTime? DeletedAt { get; set; }
}
